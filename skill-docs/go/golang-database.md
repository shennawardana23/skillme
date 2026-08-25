## What it does

Guides connection pooling, context propagation, and result scanning for
`database/sql` and `pgx` in Go. The defining constraint: `sql.DB` is
already a pool, not one connection, and its zero-value defaults
(unlimited open connections, unlimited idle connections, no connection
lifetime) are wrong for production the moment real load hits it.

## When to reach for it

Reach for this skill for the Go driver/pool layer specifically — sizing
`SetMaxOpenConns`, choosing `database/sql` vs. `pgxpool`, handling
`sql.ErrNoRows`, scanning rows into structs. `postgres-patterns` and
`mysql-patterns` decide what SQL to write against a schema;
`database-migrations` decides how a schema changes safely; this skill
sits between either of those and your Go code. For a table partitioned
by `hotel_id`, this skill's driver-level guidance doesn't substitute for
`postgres-hotel-partitioning`'s query-filtering rules — pooling and
partition-pruning are independent concerns that both apply.

## Common questions

- **"I called `sql.Open` and I'm not setting any pool limits — is that a
  problem?"** Yes. `sql.Open` doesn't even connect until first use, and
  the pool's defaults are unbounded open/idle connections with no
  lifetime — set `SetMaxOpenConns`, `SetMaxIdleConns`, and
  `SetConnMaxLifetime` explicitly, sized against the database's actual
  `max_connections`.
- **"A lookup-by-ID query returns `sql.ErrNoRows` — is that a database
  error?"** No — it's the expected outcome for "nothing matched," not a
  failure. Check it explicitly with `errors.Is(err, sql.ErrNoRows)` and
  translate it into the caller's own not-found semantics, rather than
  wrapping and propagating it alongside real database errors.
- **"Should I use `database/sql` or `pgxpool` for a Postgres-only Go
  service?"** `pgxpool` when you're Postgres-exclusive and want
  pgx-native features (typed arrays, `COPY`, batch operations, and its
  own connection-pool config); `database/sql` when the code needs to
  stay database-agnostic or already depends on a `database/sql`-based
  library.
- **"Why did my `Scan` call fail with a conversion error on a column I
  know exists?"** A `NULL` value scanned into a plain `string`/`int`
  destination fails — use `sql.NullString`/a pointer/`sql.Null[T]` (Go
  1.22+), or make the column `NOT NULL` at the schema layer if `NULL`
  was never actually meaningful there.

## It's working if

- `SetMaxOpenConns`/`SetMaxIdleConns`/`SetConnMaxLifetime` are set
  explicitly, sized against the database's real connection limit
- Every query/exec call uses its `*Context` variant
  (`QueryContext`/`ExecContext`), not the non-context one
- `sql.ErrNoRows` (or `pgx.ErrNoRows`) is checked and translated
  separately from other errors
- Every `hotel_id`-partitioned table query still carries an explicit
  `hotel_id` filter regardless of pool/driver configuration

## Where it fits

Sits between `postgres-patterns`/`mysql-patterns` (schema/query layer)
and Go application code. Pairs with `postgres-hotel-partitioning` for
partitioned tables, `golang-context` for the context propagation this
skill's `*Context` calls depend on, and `database-migrations` for schema
evolution.
