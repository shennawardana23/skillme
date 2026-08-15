---
name: erd-expert
description: Designs entity-relationship schemas for the hotel domain (reservations, guests, rooms, rate plans, folios) on PostgreSQL, following Archipelago International's Sentec platform conventions — hotel_id LIST partitioning, naming, and required audit columns. Use when designing a new table or ERD for Sentec PMS, Sentec Booking Engine, or Sentec EMS, when reviewing a proposed schema against org conventions, or when modeling the reservation lifecycle or room-availability queries.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# ERD Expert — Hotel Domain Schema Design

Entity-relationship design for the hotel domain, targeting PostgreSQL and
this organization's Sentec platform (Sentec PMS, Sentec Booking Engine,
Sentec EMS). For the mechanics of keeping queries and migrations
partition-safe once a table exists, see `postgres-hotel-partitioning` —
this skill covers designing the ERD itself; that one covers querying and
migrating it correctly afterward. For index/data-type choices beyond what's
listed here, see `postgres-patterns`.

## Core hotel-domain entities

| Entity | Partition key | Key relationships |
|---|---|---|
| `hotel` | — | root anchor; every partitioned table's `hotel_id` FKs here conceptually (not enforceable as a literal FK across partitions in all cases — validate in application logic) |
| `room` | `hotel_id` | belongs to a hotel, has a `room_type` |
| `reservation` | `hotel_id` | references guest, hotel, room; carries a status state machine |
| `guest` | `hotel_id` | PII-sensitive; has a `loyalty_tier` |
| `rate_plan` | `hotel_id` | pricing rules, scoped per channel |
| `channel` | — | small reference table: `OTA \| DIRECT \| GDS \| WHOLESALER` |
| `folio` | `hotel_id` | billing record linked to a reservation |

### Reservation state machine

```
PENDING → CONFIRMED → CHECKED_IN → CHECKED_OUT
        ↘ CANCELLED (from PENDING or CONFIRMED)
        ↘ NO_SHOW   (from CONFIRMED)
```

Model this as a `CHECK` constraint or enum on `status`, plus an audit
table (`reservation_status_history`) if the transition history itself
must be queryable — don't try to derive "when did this become CONFIRMED"
from `updated_at` on the mutable row alone.

## Required columns and conventions

Every table in the hotel domain carries these columns, matching the
convention already established for partitioned tables:

```sql
id          bigint      GENERATED ALWAYS AS IDENTITY,
hotel_id    integer     NOT NULL,   -- partition key (LIST partitioning)
created_at  timestamptz NOT NULL DEFAULT now(),
updated_at  timestamptz NOT NULL DEFAULT now(),
PRIMARY KEY (hotel_id, id)
```

`hotel_id` is `integer`, not `uuid` — this matches the partitioning
convention already in use (`postgres-hotel-partitioning`), where
partitions are declared per discrete hotel ID (`FOR VALUES IN (101)`).
Keep this consistent across every new table in the domain; introducing a
`uuid` hotel identifier on a new table would break the ability to join
it against every existing `hotel_id`-partitioned table without a lookup.

**Naming**: tables `snake_case` singular (`reservation`, not
`reservations`); primary keys `{table}_id` when referenced as a foreign
key elsewhere (`reservation.guest_id` referencing `guest.id`); every
table gets `created_at`/`updated_at`.

**Data types** (see `postgres-patterns` for the full rationale): rates
and money as `numeric(12,2)`; flexible per-channel metadata as `jsonb`;
status fields as `text` with a `CHECK` constraint (portable, easy to
extend) or a native `enum` (stricter, but every new value requires an
`ALTER TYPE` migration — prefer `CHECK` unless the value set is genuinely
frozen); soft delete via nullable `deleted_at` (`NULL` = active row).

### Partitioning rule

```sql
CREATE TABLE reservation (
    id         bigint GENERATED ALWAYS AS IDENTITY,
    hotel_id   integer NOT NULL,
    guest_id   bigint NOT NULL,
    room_id    bigint NOT NULL,
    status     text NOT NULL CHECK (status IN ('PENDING','CONFIRMED','CHECKED_IN','CHECKED_OUT','CANCELLED','NO_SHOW')),
    check_in   date NOT NULL,
    check_out  date NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (hotel_id, id)
) PARTITION BY LIST (hotel_id);

CREATE TABLE reservation_101 PARTITION OF reservation FOR VALUES IN (101);
CREATE TABLE reservation_102 PARTITION OF reservation FOR VALUES IN (102);
```

`hotel_id` must appear in the primary key (and in every unique
constraint) on a list-partitioned table — PostgreSQL requires the
partition key be part of any unique index. Every query against this
table must filter `hotel_id` explicitly for the planner to prune
partitions — see `postgres-hotel-partitioning` for the full rule and the
join-on-both-sides requirement when joining two partitioned tables (e.g.
`reservation` to `room`, both scoped by the same `hotel_id`).

### Index strategy for hotel-domain tables

- B-tree on foreign keys (`guest_id`, `room_id`) — composed with
  `hotel_id` as the leading column since that's in every query.
- BRIN on append-only time columns at scale (`created_at`, `check_in`,
  `check_out`) — cheap to maintain, effective once a partition has
  millions of rows in roughly time order.
- Partial index for hot operational paths: `WHERE status IN
  ('CONFIRMED','CHECKED_IN')` for the "currently active reservations"
  query that runs on every front-desk lookup.

## Common query patterns

**Active reservations for a hotel** (partition-pruned):

```sql
SELECT * FROM reservation
WHERE hotel_id = $1
  AND status IN ('CONFIRMED', 'CHECKED_IN')
  AND check_out >= now();
```

**Room availability** (no double-booking overlap check):

```sql
SELECT r.id, r.room_type
FROM room r
WHERE r.hotel_id = $1
  AND r.status = 'AVAILABLE'
  AND NOT EXISTS (
    SELECT 1 FROM reservation res
    WHERE res.hotel_id = $1
      AND res.room_id = r.id
      AND res.status IN ('CONFIRMED', 'CHECKED_IN')
      AND res.check_in < $3 AND res.check_out > $2
  );
```

The overlap predicate (`check_in < requested_check_out AND check_out >
requested_check_in`) is the standard interval-overlap test — get the
strict/non-strict inequality direction right, or a reservation ending
exactly on the requested check-in date will incorrectly block (or fail to
block) the new booking.

## Sentec product boundaries

- **Sentec PMS**: core property management — `hotel`, `room`,
  `reservation`, `folio`.
- **Sentec Booking Engine**: direct online booking — `channel`,
  `rate_plan`, availability queries against the same `room`/`reservation`
  tables.
- **Sentec EMS**: employee/staff management — a **separate schema**
  (staff, shifts, assignments); do not add EMS entities to the PMS ERD or
  vice versa, and do not assume EMS tables share the `hotel_id`
  partitioning scheme without confirming it against that schema's own
  design.

## Gotchas

- Forgetting `hotel_id` in a new table's primary key breaks list
  partitioning outright at `CREATE TABLE` time (Postgres raises an
  error) — this is a compile-time-equivalent catch, but only if the table
  is declared `PARTITION BY LIST` in the first place; a table someone
  forgot to partition won't error, it'll just silently violate the
  org-wide convention.
- A `channel` or other small reference table should generally **not** be
  partitioned by `hotel_id` — partitioning a table with a few dozen rows
  adds planning overhead for no pruning benefit.
- Room availability overlap queries are correctness-critical (a bug here
  double-books a room); write a test with an exact boundary case
  (requested stay starts exactly on an existing reservation's checkout
  date) before trusting the query.
- `deleted_at` soft-delete columns need every `SELECT` and every unique
  index to account for them (`WHERE deleted_at IS NULL`, or a partial
  unique index) — a plain unique constraint on `email` will reject a new
  guest reusing an email that belongs to a soft-deleted guest record.

## Real-world grounding

PostgreSQL's own documentation states the partition key must be part of
every unique, primary key, and exclusion constraint on a partitioned
table — this is not a house style choice, it is an engine-enforced
requirement, which is why every table in this ERD carries `hotel_id` as
the leading column of its primary key rather than an independent
surrogate key.

## Verification

- [ ] Every hotel-domain table has `id`, `hotel_id`, `created_at`, `updated_at`, and `PRIMARY KEY (hotel_id, id)`
- [ ] `hotel_id` is `integer`, consistent with existing partitioned tables — not introduced as `uuid`
- [ ] Money/rate columns are `numeric(12,2)`; flexible metadata is `jsonb`
- [ ] New tables that should be partitioned declare `PARTITION BY LIST (hotel_id)` with a partition per hotel
- [ ] Every query against the new table filters `hotel_id` explicitly (see `postgres-hotel-partitioning`)
- [ ] Reservation state machine transitions are constrained (CHECK or enum), not left as a free-text field
- [ ] EMS entities are kept in their own schema, not merged into the PMS/Booking Engine ERD
