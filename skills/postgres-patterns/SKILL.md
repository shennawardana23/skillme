---
name: postgres-patterns
description: Guides PostgreSQL query design, indexing strategy, data types, and row-level security patterns — the query/schema layer, as distinct from the migration-execution safety covered by database-migrations and the hotel_id partition-pruning rules covered by postgres-hotel-partitioning. Use when choosing an index type, designing a composite or partial index, picking a column data type, writing an UPSERT or a queue-processing query, or setting up row-level security.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# PostgreSQL Patterns

This skill is the query- and schema-design layer: how to index a table
and shape a query well. For *how to ship* a schema change safely, see
`database-migrations`. For any table partitioned by `hotel_id` (most core
tables in this organization), the partition-pruning predicate rules in
`postgres-hotel-partitioning` apply on top of everything here — every
example below assumes that predicate is already present.

## Index selection

| Query pattern | Index type | Example |
|---|---|---|
| `WHERE col = value` | B-tree (default) | `CREATE INDEX idx ON t (col)` |
| `WHERE col > value` / range | B-tree | `CREATE INDEX idx ON t (col)` |
| `WHERE a = x AND b > y` | Composite, equality columns first | `CREATE INDEX idx ON t (a, b)` |
| `WHERE jsonb_col @> '{}'` | GIN | `CREATE INDEX idx ON t USING gin (jsonb_col)` |
| Full-text `tsv @@ query` | GIN | `CREATE INDEX idx ON t USING gin (tsv)` |
| Append-only time-series ranges | BRIN | `CREATE INDEX idx ON t USING brin (created_at)` |

**Composite index column order matters**: put equality-tested columns
first, range-tested columns last, so the index can seek on the equality
prefix and then scan the range within it:

```sql
-- Serves: WHERE status = 'pending' AND created_at > '2025-01-01'
CREATE INDEX idx_jobs_status_created ON jobs (status, created_at);
```

**Covering index** avoids a heap lookup entirely when the query only
needs indexed + included columns:

```sql
CREATE INDEX idx_users_email ON users (email) INCLUDE (name, created_at);
-- SELECT email, name, created_at FROM users WHERE email = $1 -- index-only scan
```

**Partial index** shrinks the index to the subset that's actually queried
frequently:

```sql
CREATE INDEX idx_users_active ON users (email) WHERE deleted_at IS NULL;
```

An unindexed foreign key is a common, silent bottleneck (every delete on
the parent table triggers a sequential scan of the child to check for
references). Find them:

```sql
SELECT conrelid::regclass AS table_name, a.attname AS column_name
FROM pg_constraint c
JOIN pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
WHERE c.contype = 'f'
  AND NOT EXISTS (
    SELECT 1 FROM pg_index i
    WHERE i.indrelid = c.conrelid AND a.attnum = ANY(i.indkey)
  );
```

Always confirm an index is actually used with `EXPLAIN (ANALYZE)` before
and after adding it — a planner that has a better option (or a table too
small to bother) will ignore an unused index while it still costs every
write a maintenance overhead.

## Data types

| Use case | Prefer | Avoid | Why |
|---|---|---|---|
| Primary keys | `bigint GENERATED ALWAYS AS IDENTITY` | `int`, random `uuid` as PK | `int` overflows at ~2.1B rows; random UUIDs fragment the B-tree and defeat locality |
| Strings | `text` | `varchar(n)` | `varchar(n)` buys no performance in Postgres, only a constraint-check cost; use a `CHECK` constraint if you truly need a length cap |
| Timestamps | `timestamptz` | `timestamp` (no zone) | `timestamp` silently drops zone info; every ambiguity it causes is a production bug about which offset "midnight" meant |
| Money / rates | `numeric(12,2)` | `float`, `double precision` | float arithmetic on currency accumulates rounding error |
| Flags | `boolean` | `varchar('Y'/'N')`, `int` | boolean is self-documenting and enables partial indexes cleanly |
| Flexible/variable attributes | `jsonb` | `json` | `jsonb` is stored decomposed and binary-indexable (GIN); `json` is stored as text and re-parsed on every access |

## Query patterns

**UPSERT** — insert-or-update in one round trip, avoiding a
check-then-act race:

```sql
INSERT INTO settings (user_id, key, value)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, key) DO UPDATE SET value = EXCLUDED.value;
```

**Cursor (keyset) pagination** — `O(1)` regardless of page depth, unlike
`OFFSET` which must scan and discard every preceding row:

```sql
-- BAD: OFFSET cost grows linearly with page number
SELECT * FROM products ORDER BY id LIMIT 20 OFFSET 10000;

-- GOOD: keyset pagination, constant cost regardless of depth
SELECT * FROM products WHERE id > $last_seen_id ORDER BY id LIMIT 20;
```

**Queue processing** — `SKIP LOCKED` lets multiple workers pull from the
same job table concurrently without blocking on each other's row locks:

```sql
UPDATE jobs SET status = 'processing', locked_at = now()
WHERE id = (
    SELECT id FROM jobs
    WHERE status = 'pending'
    ORDER BY created_at
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;
```

Without `SKIP LOCKED`, a second worker's `SELECT ... FOR UPDATE` blocks
on the first worker's lock instead of moving on to the next unlocked row
— turning concurrent workers into a serial queue.

## Row-level security

RLS enforces a filter at the database layer that no application bug can
bypass — a defense-in-depth complement to (never a replacement for)
application-level authorization checks and `hotel_id` predicates:

```sql
ALTER TABLE reservations ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON reservations
    USING (hotel_id = current_setting('app.current_hotel_id')::int);
```

**Wrap the session-context call in a scalar subquery-safe form** —
`current_setting(...)` re-evaluates per row unless the planner can prove
it's stable; on Postgres 15+ wrapping it as `(SELECT current_setting(...))`
lets the planner cache it once per statement instead of once per row:

```sql
CREATE POLICY tenant_isolation ON reservations
    USING (hotel_id = (SELECT current_setting('app.current_hotel_id'))::int);
```

Set the session variable per connection/request, scoped with
`set_config(..., true)` (`true` = local to the current transaction) so
one connection-pooled request can never leak its tenant context into the
next request reusing that connection:

```sql
SELECT set_config('app.current_hotel_id', $1, true);
```

## Anti-pattern detection queries

```sql
-- Slow queries (requires pg_stat_statements extension)
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
WHERE mean_exec_time > 100
ORDER BY mean_exec_time DESC;

-- Table bloat / vacuum falling behind
SELECT relname, n_dead_tup, last_vacuum, last_autovacuum
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC;
```

## Gotchas

- `CREATE INDEX` (without `CONCURRENTLY`) on an existing table blocks
  writes for the build's duration — that's a migration-safety concern,
  covered in `database-migrations`, not a query-design one, but it bites
  people who copy an index definition from here straight into a migration
  against a live table.
- A composite index `(a, b)` serves `WHERE a = x` and `WHERE a = x AND b
  = y`, but not `WHERE b = y` alone — the leading column must be present
  in the predicate for the index to be used at all.
- `numeric` without a precision/scale (bare `numeric`) accepts unbounded
  digits and defeats the point of a fixed-point money type — always
  specify `numeric(precision, scale)`.
- `SKIP LOCKED` skips rows locked by *any* transaction, including ones
  that will eventually roll back — a worker can legitimately claim a job
  another worker's failed transaction "held," which is correct behavior,
  not a bug, but surprises people expecting strict FIFO ordering.
- RLS policies are bypassed entirely for the table owner and any role
  with `BYPASSRLS` — including most default admin/migration roles. Verify
  the actual application connects as a role without that attribute, or
  the policy is decorative.

## Real-world grounding

`SELECT ... FOR UPDATE SKIP LOCKED` was added in PostgreSQL 9.5 (2016)
specifically to make Postgres viable as a job queue backend without a
dedicated message broker — it's the same mechanism referenced in
`backend-patterns` for durable background job processing, and is widely
used in production queue implementations (e.g. Ruby's `good_job`,
Python's `procrastinate`) as an alternative to Redis-backed queues when a
team wants queue and application data in the same transactional store.

## Verification

- [ ] Every index choice matches the actual query predicate shape (equality vs range, composite column order)
- [ ] Foreign keys have a supporting index unless proven unnecessary
- [ ] Money uses `numeric(p,s)`; timestamps use `timestamptz`; flexible attributes use `jsonb`
- [ ] List endpoints paginate by cursor/keyset, not `OFFSET`, once the table is large
- [ ] Queue-style concurrent workers use `FOR UPDATE SKIP LOCKED`
- [ ] RLS policies (if used) wrap session-context calls for planner caching, and the app role does not have `BYPASSRLS`
- [ ] `EXPLAIN (ANALYZE)` confirms an added index is actually chosen by the planner
