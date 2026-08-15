---
name: postgres-hotel-partitioning
description: This skill should be used when the user asks to "write a query for a hotel-partitioned table", "review this SQL for partition safety", "add a hotel_id filter", "check partition pruning", "write a migration for a partitioned table", or writes/reviews Go or SQL code touching any PostgreSQL table partitioned by hotel_id. Provides patterns for keeping queries, joins, and migrations partition-safe.
metadata:
  version: "0.1.0"
---

# PostgreSQL Hotel-Partitioned Tables

Many tables in this organization's PostgreSQL databases are partitioned by
`hotel_id` (declarative `PARTITION BY LIST (hotel_id)`). PostgreSQL's planner
only prunes partitions it can prove are irrelevant from the query's `WHERE`
clause at plan time — a query missing an explicit `hotel_id` predicate scans
every partition, not just the one relevant to the request. This is the
single most common correctness-and-performance bug against these tables:
apply the checks below to every query, join, and migration.

## The core rule

Every query against a partitioned-by-`hotel_id` table must include
`hotel_id = $N` (or `hotel_id = ANY($N)` for a known, bounded set) directly in
its `WHERE` clause — not buried inside a subquery the planner can't push
down, not applied only in application code after fetching rows.

```go
// BAD: no hotel_id predicate — scans every hotel's partition to find status matches.
const q = `SELECT * FROM reservations WHERE status = $1`

// GOOD: hotel_id is explicit, sits alongside the other predicate.
const q = `SELECT * FROM reservations WHERE hotel_id = $1 AND status = $2`
```

Always bind `hotel_id` as a parameter (`$1`, `:hotel_id`), never interpolate
it or any other value into the query string with `fmt.Sprintf` — that is a
SQL-injection vector independent of the partitioning concern.

## Joins across partitioned tables

A join between two tables partitioned by `hotel_id` must filter `hotel_id` on
**both** sides (or rely on the same bound parameter for both), otherwise the
planner cannot prune either side:

```sql
-- GOOD: hotel_id constrained on both tables in the join.
SELECT r.id, g.name
FROM reservations r
JOIN guests g ON g.id = r.guest_id AND g.hotel_id = r.hotel_id
WHERE r.hotel_id = $1 AND r.status = $2;
```

## Go query helpers

Prefer a thin wrapper that makes `hotelID` a required, typed parameter so it
cannot be forgotten at the call site — push the constraint into the
function signature rather than trusting every call site to remember it:

```go
// ByHotelAndStatus is the only way to query reservations by status —
// hotelID is a required parameter, not an optional filter callers can skip.
func (s *Store) ByHotelAndStatus(ctx context.Context, hotelID int64, status string) ([]Reservation, error) {
    const q = `SELECT id, hotel_id, status, checkin, checkout
               FROM reservations
               WHERE hotel_id = $1 AND status = $2`
    var out []Reservation
    if err := s.db.SelectContext(ctx, &out, q, hotelID, status); err != nil {
        return nil, fmt.Errorf("query reservations for hotel %d: %w", hotelID, err)
    }
    return out, nil
}
```

## Migrations

Declare the parent table with `PARTITION BY LIST (hotel_id)` and create one
partition per hotel (or a bucketed range if the hotel count is large enough
that per-hotel partitions become unwieldy — confirm the actual hotel count
before choosing bucketing over one-partition-per-hotel):

```sql
CREATE TABLE stays (
    id         bigint GENERATED ALWAYS AS IDENTITY,
    hotel_id   integer NOT NULL,
    guest_id   bigint NOT NULL,
    checkin    date NOT NULL,
    checkout   date NOT NULL,
    PRIMARY KEY (hotel_id, id)
) PARTITION BY LIST (hotel_id);

CREATE TABLE stays_101 PARTITION OF stays FOR VALUES IN (101);
CREATE TABLE stays_102 PARTITION OF stays FOR VALUES IN (102);
```

Note `hotel_id` must be part of every unique index and the primary key on a
list-partitioned table — PostgreSQL requires the partition key be included
in any unique constraint.

## Verifying partition pruning

Confirm pruning actually happens rather than assuming it from the query
shape — run `EXPLAIN (ANALYZE)` and check the plan lists only the expected
partition(s), not `Append` over every child:

```sql
EXPLAIN (ANALYZE) SELECT * FROM reservations WHERE hotel_id = 101 AND status = 'confirmed';
-- Expect: Index Scan on reservations_101 ...
-- Not:    Append -> Seq Scan on reservations_101 -> Seq Scan on reservations_102 -> ...
```

## Additional resources

For a walkthrough of the missing-filter bug pattern and how it shows up in
code review, consult `references/partitioning.md`.
