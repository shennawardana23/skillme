---
name: database-migrations
description: Guides safe, reversible PostgreSQL schema and data migrations using golang-migrate, Prisma, Drizzle, or Django migrations. Use when creating or altering tables, adding or removing columns or indexes, planning a zero-downtime schema change, or reviewing a migration before it runs against production. For tables partitioned by hotel_id, use the postgres-hotel-partitioning skill alongside this one.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Database Migrations

Every schema change is a migration file, never a manual `psql` session
against production. Migrations are forward-only once deployed — a mistake
is fixed with a new migration, not by editing or deleting the one that ran.

## Core rules

1. Schema changes (DDL) and data backfills (DML) are **separate**
   migrations — mixing them makes a migration slow, lock-prone, and hard to
   roll back independently.
2. `NOT NULL` on a new column requires a default, or a nullable column
   followed by a backfill and a later constraint migration — adding
   `NOT NULL` with no default to an existing table locks it for a full
   rewrite.
3. Indexes on existing tables use `CREATE INDEX CONCURRENTLY` — a plain
   `CREATE INDEX` blocks writes for the duration of the build.
4. Test migrations against a production-sized copy, not a 100-row dev
   database — lock duration and rewrite cost scale with table size, and a
   migration that's instant at 100 rows can lock for minutes at 10M.

## Zero-downtime pattern (expand → migrate → contract)

```
EXPAND:   add new column/table (nullable or defaulted); app writes to BOTH old and new
MIGRATE:  backfill; app reads from NEW, still writes to BOTH; verify consistency
CONTRACT: app uses NEW only; drop the OLD column/table in a separate, later migration
```

Renaming a column in production always goes through this pattern — never a
direct `ALTER TABLE ... RENAME COLUMN` on a table any running application
code still references by the old name.

## Go (golang-migrate) — this org's default

```bash
migrate create -ext sql -dir migrations -seq add_reservation_notes
migrate -path migrations -database "$DATABASE_URL" up
migrate -path migrations -database "$DATABASE_URL" down 1
```

```sql
-- 000004_add_reservation_notes.up.sql
ALTER TABLE reservations ADD COLUMN notes TEXT;
CREATE INDEX CONCURRENTLY idx_reservations_notes ON reservations (notes) WHERE notes IS NOT NULL;

-- 000004_add_reservation_notes.down.sql
DROP INDEX IF EXISTS idx_reservations_notes;
ALTER TABLE reservations DROP COLUMN IF EXISTS notes;
```

## Gotchas

- `CREATE INDEX CONCURRENTLY` cannot run inside a transaction block — most
  migration tools wrap each migration in a transaction by default, so this
  needs explicit non-transactional handling (golang-migrate's `-x` flag
  behavior varies by driver; verify your tool's specific mechanism before
  relying on it).
- For any table partitioned by `hotel_id`, a migration that adds an index
  or constraint must apply it per-partition awareness — see
  `postgres-hotel-partitioning` for the specific patterns; a migration
  written as if the table were unpartitioned can silently miss partitions
  created after the migration ran.
- Large batch backfills should use `LIMIT` + `FOR UPDATE SKIP LOCKED` in a
  loop with periodic `COMMIT`, not a single unbounded `UPDATE` — the latter
  holds one long transaction and locks the table for its full duration.

## Real-world grounding

GitLab's 2017 production incident: an engineer, attempting to fix
replication lag, ran `rm -rf` against what they believed was the secondary
database directory but was actually production — and every one of the five
backup mechanisms in place (regular backups, disk snapshots, LVM snapshots,
S3 backups, and a separate replica) had been silently failing or
misconfigured for weeks, discovered only when the team tried to actually
restore from one. The lesson isn't "have backups" — it's that a backup
strategy that has never been tested with an actual restore is not a backup
strategy, it's an assumption.

## Verification

- [ ] Migration has both an up and a down (or is explicitly irreversible, documented as such)
- [ ] No `NOT NULL` added to an existing table without a default
- [ ] Indexes on existing tables use `CONCURRENTLY`
- [ ] Schema change and data backfill are separate migrations
- [ ] Tested against production-sized data, not a small dev sample
- [ ] hotel_id-partitioned tables checked against `postgres-hotel-partitioning`
