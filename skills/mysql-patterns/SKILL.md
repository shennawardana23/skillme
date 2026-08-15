---
name: mysql-patterns
description: Guides MySQL/MariaDB query design, indexing strategy, data types, and concurrency patterns — the InnoDB-specific counterpart to postgres-patterns. Use when choosing a primary key strategy, picking a charset, writing an UPSERT, deciding on isolation level, or porting a query pattern from PostgreSQL to MySQL/MariaDB (or the reverse) and needing to know what actually differs.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "database"
---

# MySQL / MariaDB Patterns

This is the query- and schema-design layer for MySQL (InnoDB) and MariaDB,
written for a team whose default database is PostgreSQL (see
`postgres-patterns`, `postgres-hotel-partitioning`) but that also operates
MySQL/MariaDB-backed systems (commonly paired with CodeIgniter/Laravel
legacy stacks — see `php-codeigniter-patterns`, `laravel-patterns`). Every
section below calls out where MySQL's real, documented behavior differs from
Postgres, because that difference — not the syntax — is where migrated
assumptions actually break.

## Primary keys: the clustered-index difference that changes everything else

InnoDB stores a table's rows physically ordered by primary key — the PK
**is** the clustering key, and every secondary index stores the PK value as
its row pointer (not a physical row address, the way Postgres's heap-based
secondary indexes work). Two direct consequences:

- **A non-monotonic primary key (e.g. a random UUID) causes page splits and
  index fragmentation** on every insert, because new rows aren't appended at
  the end of physical storage — they're inserted wherever the PK value's
  ordering says they belong. This is a far more serious performance concern
  in InnoDB than in Postgres, where a random UUID PK costs only a
  (already-real) B-tree-locality penalty on that one index, not the
  clustering order of the whole table.
- **Every secondary index lookup requires a second lookup by PK value**
  (unless the query is covered), because a secondary index only stores the
  PK, not the full row — a wide/expensive PK type inflates every secondary
  index's storage cost, not just the primary key's own.

Prefer an auto-increment (or otherwise monotonically increasing) `BIGINT`
primary key for InnoDB tables. If a UUID is required for external
referencing, keep it as a separate unique-indexed column rather than the
clustering PK, or use a time-ordered UUID variant (UUIDv7-style) specifically
to preserve insertion order.

## Charset: `utf8mb4`, never `utf8`

MySQL's `utf8` charset is a **3-byte-max legacy encoding** that cannot store
most emoji or a meaningful share of CJK supplementary-plane characters — a
well-documented MySQL footgun, not a hypothetical one. `utf8mb4` is the real,
complete UTF-8 implementation (up to 4 bytes/codepoint). Every new table,
and every column added to an existing one:

```sql
CREATE TABLE guests (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

MySQL 8.0 made `utf8mb4` the server default (earlier versions defaulted to
`latin1`), but a database or table created under an older default, or
migrated from a MySQL 5.7 dump without an explicit charset conversion, can
still silently carry `utf8`/`latin1` — check `SHOW CREATE TABLE` on any
inherited schema rather than assuming a "MySQL 8" deployment means every
existing table is already `utf8mb4`.

## Query patterns

**UPSERT** — MySQL's insert-or-update, functionally parallel to Postgres's
`ON CONFLICT`:

```sql
INSERT INTO settings (user_id, `key`, value)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE value = VALUES(value);
```

Requires a unique index or primary key on the conflicting column(s)
(`user_id, key` here) for MySQL to detect the duplicate at all — without one,
this is just a plain `INSERT` that fails on any other constraint violation.

**Cursor (keyset) pagination** — same principle as Postgres, same reason:
`OFFSET` cost grows with page depth regardless of database engine.

```sql
-- BAD: cost grows linearly with page number, same problem as Postgres
SELECT * FROM products ORDER BY id LIMIT 20 OFFSET 10000;

-- GOOD: keyset pagination
SELECT * FROM products WHERE id > ? ORDER BY id LIMIT 20;
```

**Queue processing with `SKIP LOCKED`** — available, but **version-gated**
in a way that genuinely differs between the two engines this catalog cares
about:

```sql
SELECT id FROM jobs
WHERE status = 'pending'
ORDER BY created_at
LIMIT 1
FOR UPDATE SKIP LOCKED;
```

`SKIP LOCKED` (and `NOWAIT`) syntax was added to **MySQL in 8.0** — it does
not exist in MySQL 5.7. MariaDB's story is less clean: the syntax appears
earlier, but functioning skip-locked-row behavior wasn't reliably available
until **MariaDB 10.6**. Confirm the actual deployed engine and version
before relying on `SKIP LOCKED` — a query that parses successfully in an
older MariaDB does not guarantee it behaves as skip-locked; verify against
the specific version in the deployment target, not against "it's MySQL/
MariaDB so it should work like Postgres 9.5+."

## Isolation level: `REPEATABLE READ` is the real default, not `READ COMMITTED`

This is the single most consequential semantic difference from Postgres.
InnoDB's default transaction isolation level is `REPEATABLE READ`; Postgres's
default is `READ COMMITTED`. Under `REPEATABLE READ`, a transaction's
consistent (non-locking) reads are snapshotted at the *first* read in that
transaction and stay fixed for its duration — a second `SELECT` for the same
row inside the same still-open transaction won't see another transaction's
committed change in between, which is different from Postgres's
`READ COMMITTED` default, where each individual statement sees the latest
committed data.

```sql
SELECT @@transaction_isolation;  -- MySQL 8.0+
SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED;  -- if the app logic
                                                          -- was written assuming
                                                          -- Postgres semantics
```

A locking read (`SELECT ... FOR UPDATE`) always sees the latest committed
row regardless of isolation level — the snapshot behavior above applies to
plain, non-locking `SELECT`s. Code ported from a Postgres codebase that
assumes "every `SELECT` sees the latest commit" can silently behave
differently under MySQL's default isolation level in exactly the read-then-
read-again patterns that look harmless in review.

## Gotchas

- **A random (non-sequential) primary key is a much bigger problem in
  InnoDB than in Postgres** — it fragments the clustered index that *is*
  the table's physical storage, not just one secondary structure.
- **`utf8` is not UTF-8.** It is a 3-byte legacy encoding that silently
  truncates or rejects 4-byte characters (most emoji, some CJK). Always use
  `utf8mb4`.
- **`SKIP LOCKED` is unavailable on MySQL 5.7** and had an inconsistent
  rollout on MariaDB before 10.6 — check the actual deployed version, don't
  assume parity with Postgres's long-stable `SKIP LOCKED` support (since
  9.5, 2016).
- **`REPEATABLE READ`'s snapshot semantics can mask a lost-update bug**
  that Postgres's `READ COMMITTED` default would have surfaced differently
  — a non-locking read-modify-write pattern that "worked" under Postgres's
  default isolation may need an explicit `FOR UPDATE` or a
  `READ COMMITTED` session override under MySQL to behave the same way.
- **`ON DUPLICATE KEY UPDATE` silently does nothing meaningful without a
  real unique constraint on the conflicting columns** — it is not a
  general-purpose "upsert this row" statement independent of schema design,
  the way it might read at a glance.

## Real-world grounding

InnoDB's clustered-index-by-primary-key design is documented directly in
the MySQL 8.0 Reference Manual's InnoDB architecture chapter, and is the
standard, widely-cited reason UUID primary keys are treated as a much more
serious anti-pattern on MySQL/InnoDB than on Postgres — the same UUID-as-PK
choice that costs Postgres only B-tree locality on one index reorders the
physical storage of the entire InnoDB table. The `utf8`-vs-`utf8mb4`
distinction is one of MySQL's most repeatedly documented historical
footguns, significant enough that MySQL 8.0 changed the server default
charset specifically to close it going forward — but that default change
does not retroactively fix schemas created under earlier defaults.

## Verification

- [ ] Primary keys on InnoDB tables are monotonic (auto-increment or
      time-ordered), not random UUIDs used directly as the clustering key
- [ ] Every table and text/varchar column uses `utf8mb4`, confirmed via
      `SHOW CREATE TABLE` rather than assumed from server version
- [ ] `ON DUPLICATE KEY UPDATE` targets a column set with a real unique
      constraint
- [ ] Any `SKIP LOCKED` usage is confirmed against the actual deployed
      MySQL/MariaDB version, not assumed available
- [ ] Code that assumes "every SELECT sees the latest commit" is reviewed
      against InnoDB's `REPEATABLE READ` default, with an explicit
      `FOR UPDATE` or isolation-level override where that assumption
      actually matters
