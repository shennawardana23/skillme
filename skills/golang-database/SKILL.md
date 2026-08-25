---
name: golang-database
description: This skill should be used when the user asks to "configure a Go database connection pool", "use database/sql in Go", "set up pgxpool", "why is my Go app running out of connections", "scan rows into a Go struct", "handle sql.ErrNoRows", or writes Go code that opens, pools, queries, or scans a PostgreSQL/MySQL connection via database/sql or pgx. postgres-patterns and mysql-patterns cover the SQL/schema layer (indexes, data types, query shape); database-migrations covers schema evolution; this skill covers the Go driver/pool mechanics between your code and either one. For hotel_id-partitioned tables, apply postgres-hotel-partitioning's filtering rules to every query this skill helps you write.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Go Database Driver Mechanics

`postgres-patterns`/`mysql-patterns` decide what SQL to write;
`database-migrations` decides how schema changes ship safely; this skill
covers the layer between them and your Go code — connection pooling,
context propagation into queries, and scanning results without silent
bugs.

## Connection pooling: `database/sql`

`sql.DB` is already a connection pool, not one connection — `sql.Open`
doesn't even connect until the first use. Configure it explicitly; the
zero-value defaults (unlimited open connections, unlimited idle
connections, no connection lifetime) are wrong for production:

```go
db.SetMaxOpenConns(25)          // cap total connections the pool can hold
db.SetMaxIdleConns(25)          // keep idle conns ready instead of reopening
db.SetConnMaxLifetime(5 * time.Minute)  // force rotation past a load balancer/DB restart
db.SetConnMaxIdleTime(1 * time.Minute)  // release idle conns the pool doesn't need
```

`SetMaxOpenConns` with no `SetMaxIdleConns` set can thrash — connections
open under load then get closed immediately once idle because the idle
limit defaults lower, then reopen on the next request. Set both, and size
them against the database's actual `max_connections`, not just your
app's guess at concurrency.

## Connection pooling: `pgxpool` (jackc/pgx/v5)

`pgxpool.New(ctx, connString)` or `pgxpool.NewWithConfig(ctx, cfg)` parse
their own equivalents of the settings above onto a `*pgxpool.Config`:
`MaxConns`, `MinConns`, `MaxConnLifetime`, `MaxConnIdleTime`,
`HealthCheckPeriod`. Prefer `pgxpool` over `database/sql` + a `pgx`
driver shim when you're on Postgres exclusively and want pgx-native
features (typed arrays, `COPY`, batch, native `context` cancellation
mid-query) — reach for `database/sql` when the code must stay
database-agnostic or already depends on a `database/sql`-based library.

## Always pass context through to the call that does I/O

Every query, exec, and prepare has a `*Context` variant —
`QueryContext`, `ExecContext`, `QueryRowContext`, `PrepareContext` (pgx:
`Query`, `Exec`, `QueryRow` already take `ctx` as their first argument).
Using the non-context variant (`db.Query` instead of
`db.QueryContext(ctx, ...)`) means the caller's cancellation or timeout
has no way to reach the database driver — the query runs to completion
regardless of what upstream gave up on.

```go
row := db.QueryRowContext(ctx, `SELECT name FROM guests WHERE id = $1 AND hotel_id = $2`, id, hotelID)
```

## Handling "no rows" without an error-shaped bug

`QueryRow(...).Scan(...)` returns `sql.ErrNoRows` when nothing matched —
this is an expected outcome for a lookup-by-ID, not a database failure.
Check it explicitly with `errors.Is(err, sql.ErrNoRows)` and translate it
to whatever "not found" means in the caller's domain; don't let it
propagate as a generic `%w`-wrapped 500 alongside real database errors.
`pgx` returns the equivalent `pgx.ErrNoRows`.

```go
var name string
err := db.QueryRowContext(ctx, q, id).Scan(&name)
switch {
case errors.Is(err, sql.ErrNoRows):
    return nil, ErrGuestNotFound
case err != nil:
    return nil, fmt.Errorf("query guest %s: %w", id, err)
}
```

## Scanning rows into structs

`database/sql` has no built-in struct scanning — `rows.Scan(&a, &b, &c)`
must list destinations in the exact column order of the `SELECT`. `pgx`
v5 adds `pgx.CollectRows(rows, pgx.RowToStructByName[T])`, which matches
columns to struct fields by name (via `db` struct tags) instead of
position — safer against a `SELECT *` column reorder, but still requires
every selected column to have a matching tagged field or it errors.

```go
type Guest struct {
    ID   string `db:"id"`
    Name string `db:"name"`
}
rows, err := pool.Query(ctx, `SELECT id, name FROM guests WHERE hotel_id = $1`, hotelID)
if err != nil { return nil, err }
guests, err := pgx.CollectRows(rows, pgx.RowToStructByName[Guest])
```

## Transactions: commit or rollback on every path

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback() // no-op if Commit already succeeded

if _, err := tx.ExecContext(ctx, insertSQL, args...); err != nil {
    return fmt.Errorf("insert: %w", err)
}
return tx.Commit()
```

`tx.Rollback()` after a successful `Commit()` is documented as a safe
no-op — the `defer tx.Rollback()` pattern covers every early-return path
without needing an explicit rollback call at each one.

## Gotchas

- `sql.DB.Close()` is almost never what you want mid-request — `sql.DB`
  represents the pool for the app's lifetime; closing it after one
  request tears down the pool for every other in-flight request too.
- A `NULL` column scanned into a plain `string`/`int` (not
  `sql.NullString`/a pointer/`sql.Null[T]`, Go 1.22+) fails the `Scan`
  call with a conversion error — either use the `sql.Null*` types, scan
  into a pointer, or make the column `NOT NULL` at the schema layer if
  `NULL` was never actually meaningful.
- `rows.Next()`/`rows.Scan()` loops must call `rows.Close()` (or exhaust
  `Next()` to `false`, which closes it automatically) and check
  `rows.Err()` after the loop — a `Scan` error inside the loop doesn't
  stop iteration by itself, and a loop that `break`s early on some other
  condition without closing leaks the underlying connection back to the
  pool later than necessary.
- `db.Prepare` outside a transaction on `database/sql` is a pooled,
  per-connection prepared statement, not a single global one — re-`Prepare`ing
  the same query string on every request defeats its purpose and adds
  overhead; prepare once at startup for hot queries instead, or rely on
  the driver's own statement cache (pgx caches automatically by default).
- Every query against a table partitioned by `hotel_id` still needs an
  explicit `hotel_id` filter regardless of how the pool or context is
  configured — the driver mechanics in this skill don't substitute for
  the partition-pruning rules in `postgres-hotel-partitioning`.

## Real-world grounding

`database/sql`'s pool tuning methods (`SetMaxOpenConns`,
`SetConnMaxLifetime`, etc.) and `sql.ErrNoRows` are documented directly in
the standard library (`go doc database/sql.DB`); `pgxpool.Config`'s
`MaxConns`/`MinConns`/`HealthCheckPeriod` fields and `pgx.CollectRows` are
documented on pkg.go.dev under `github.com/jackc/pgx/v5/pgxpool` and
`github.com/jackc/pgx/v5`.

## Cross-references

- For index/data-type/query-shape decisions once the connection works,
  see `postgres-patterns` or `mysql-patterns`.
- For `hotel_id`-partitioned tables, apply `postgres-hotel-partitioning`
  to every query this skill helps you write.
- For schema changes, see `database-migrations`.
- For propagating cancellation into `*Context` calls, see
  `golang-context`.
