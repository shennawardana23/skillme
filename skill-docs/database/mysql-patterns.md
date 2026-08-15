## What it does

Guides schema and query design for MySQL/MariaDB (InnoDB) — the same job
`postgres-patterns` does for PostgreSQL, but for an engine whose defaults
genuinely differ in ways that break assumptions carried over from Postgres.
The one constraint everything else here follows from: InnoDB's primary key
*is* the physical row order (a clustered index), not just one more index
among several, so a primary-key choice affects storage and every secondary
index, not just lookups on that one column.

## When to reach for it

Reach for this skill when a table, query, or migration targets MySQL or
MariaDB specifically — choosing a primary key strategy, picking a charset,
writing an upsert, or deciding whether a concurrency pattern that works on
Postgres actually works the same way here. If the target database is
PostgreSQL, use `postgres-patterns` instead; the two skills cover the same
kind of decision but the concrete engine behavior is not interchangeable,
and this skill exists specifically to flag the places it isn't.

## The porting trap

Most of the value here is not new SQL syntax — it's catching an assumption
that was true on Postgres and is silently false on MySQL: that a random
primary key only costs one index's locality (false — it costs the whole
table's physical layout), that `SKIP LOCKED` is available the moment you
need it (false below MySQL 8.0, and inconsistently true on older MariaDB),
that every plain `SELECT` sees the latest committed row (false under
InnoDB's default isolation level). None of these are exotic edge cases —
they're the first three things a team discovers the hard way when a service
that was designed against Postgres gets deployed against a MySQL-backed
legacy system instead.

## Common questions

- **"We used a UUID as the primary key to avoid guessable IDs — is that a
  problem on MySQL specifically?"** Yes, more so than on Postgres. Because
  InnoDB clusters the table by primary key, a non-sequential PK causes page
  splits across the table's actual physical storage on every insert, not
  just fragmentation of one secondary index. Keep the UUID as a separate
  unique-indexed column and use a monotonic key (auto-increment, or a
  time-ordered UUID variant) as the actual primary key.
- **"Our queue table uses `FOR UPDATE SKIP LOCKED` and it's not behaving
  like it does on our other MySQL box — why?"** Version. `SKIP LOCKED` was
  added in MySQL 8.0; it isn't available on 5.7. MariaDB's rollout was
  inconsistent — treat MariaDB versions before 10.6 as unreliable for this
  specific behavior, and confirm the deployed version rather than assuming
  parity with Postgres (which has had this since 9.5).
- **"A row I just updated in another connection doesn't show up in a
  `SELECT` inside an open transaction — is that a bug?"** No — it's InnoDB's
  default `REPEATABLE READ` isolation level, which snapshots a transaction's
  consistent reads at first read. Postgres defaults to `READ COMMITTED`,
  where each statement sees the latest commit; code migrated from a
  Postgres-assumption codebase can need an explicit `FOR UPDATE` or a
  session-level isolation override to behave the same way on MySQL.
- **"Do we need to worry about charset if we're already on MySQL 8?"** Only
  for schemas or dumps inherited from before the 8.0 default change (which
  moved the server default to `utf8mb4`). A table created under an older
  default, or restored from an old 5.7 dump without explicit conversion, can
  still be sitting on the legacy 3-byte `utf8` charset that silently can't
  store most emoji — check `SHOW CREATE TABLE`, don't assume from the
  server version alone.

## It's working if

- A new InnoDB table's primary key is monotonic, with any UUID need served
  by a separate unique-indexed column
- Every table and text column is confirmed `utf8mb4` via `SHOW CREATE TABLE`
- Any `SKIP LOCKED` usage has been checked against the actual deployed
  engine and version, not assumed
- Code that reads-then-writes inside a transaction has been reviewed
  against InnoDB's `REPEATABLE READ` default, not just tested happy-path

## Where it fits

Standalone reference skill, the direct engine-specific counterpart to
`postgres-patterns`. Commonly paired with `php-codeigniter-patterns` or
`laravel-patterns` when the application layer sitting on top of a MySQL
database is a CodeIgniter or Laravel app — this skill covers the database
side of that stack, those skills cover the framework side.
