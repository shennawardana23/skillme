---
name: raw-sql-querybuilder-testing-patterns
description: Guides Go services that use database/sql directly with a lightweight in-house query builder (not an ORM like GORM or a codegen tool like sqlc), including how to make repositories unit-testable via mock-generated interfaces. Use when writing or reviewing a repository method built on database/sql, adding a new repository interface for mocking, or debugging a nil-scan/mock-drift test failure.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Raw database/sql + Query Builder + Mock-Generated Repositories

A common Go data-access shape that isn't an ORM (GORM), a codegen tool
(sqlc), or a fluent third-party builder (squirrel): a small, in-house query
builder wrapping `database/sql` directly, plus interfaces over `*sql.DB`/
`*sql.Tx` specifically so `mockgen`-generated mocks can stand in for the
database in unit tests. Distinct from `postgres-patterns` (which assumes
`pgx`/idiomatic Postgres tooling) and `database-migrations` (which assumes
a specific migration framework) — this skill is about the repository/
data-access layer itself when the stack is exactly this shape.

## Interfaces exist for testability, not abstraction for its own sake

```go
type Q interface {
    QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

//go:generate mockgen -destination=mocks/mock_db.go -package=mocks . Q
```

The point of wrapping `*sql.DB`/`*sql.Tx` behind a narrow interface like
this is exclusively to let a repository be unit-tested against a generated
mock instead of a real database — it is not a general abstraction layer to
extend for its own sake. Keep the interface as narrow as the repository
actually needs (the specific `*Context` methods it calls), not the full
`*sql.DB` surface.

## Translate driver errors at the repository boundary

```go
func (r *RoomRepository) FindByID(ctx context.Context, id int64) (*Room, error) {
    row := r.db.QueryRowContext(ctx, "SELECT ... FROM rooms WHERE id = ?", id)
    var room Room
    if err := row.Scan(&room.ID, &room.Name /* ... */); err != nil {
        return nil, translateDBError(ctx, err)
    }
    return &room, nil
}
```

Callers above the repository layer should see a consistent, translated
error type — not a raw driver error (`sql.ErrNoRows`, a MySQL-specific
error code, a Postgres-specific error code) leaking upward. Centralize the
translation once, at the boundary, rather than letting each repository
method interpret driver errors its own way.

## Nullable columns need nullable scan targets

A column that can be `NULL` (including a column added after existing rows
already exist, which start out `NULL` regardless of the new column's
eventual intended default) must be scanned into a nullable type
(`sql.NullString`, `sql.NullInt64`, or a custom `Scan`/`Value`-implementing
type), never a bare primitive:

```go
var displayName sql.NullString
row.Scan(&room.ID, &displayName)
if displayName.Valid {
    room.DisplayName = displayName.String
}
```

Scanning `NULL` into a non-nullable Go type doesn't just leave that one
field wrong — it fails the entire `Scan` call (and therefore the whole
row), which can look like a query bug rather than the actual cause (a
column that legitimately has `NULL` rows and a scan target that can't
represent that).

## Keeping the mock registry current

A manually-maintained list of "which repository interfaces get mocked"
(a Makefile variable, a `//go:generate` directive per interface) will
drift from the actual set of repository interfaces in the codebase unless
adding one is treated as a required step of adding the interface itself,
not an afterthought. A repository interface missing from that list means
tests either fail to compile against a stale mock or silently use an
outdated mock that no longer matches the real interface's method set.

## Gotchas

- **A DB-wrapping interface should stay as narrow as the repository
  actually calls** — resist widening it to the full `*sql.DB`/`*sql.Tx`
  surface just because it's convenient; a narrower interface is both
  easier to mock meaningfully and a clearer statement of what a
  repository actually depends on.
- **Untranslated driver errors leaking past the repository boundary**
  force every caller to know database-specific error shapes — centralize
  translation once at the repository boundary instead.
- **Scanning a nullable column into a non-nullable Go type fails the
  whole row's `Scan` call**, not just that field — this is a very common
  cause of a "why does this query suddenly fail" bug the moment a new
  nullable column is added to an existing table with existing rows.
- **A manually-maintained mock-generation list drifts from the actual
  repository interface set** unless updating it is a required, checked
  step of adding a new repository interface — a stale mock either fails
  to compile or silently tests against an outdated interface shape.
- **A hand-rolled query builder's chaining API is not interchangeable
  with a third-party one** (squirrel, sqlx, an ORM's query DSL) — code
  reviewed or generated with the wrong builder's idioms in mind won't
  compile against this shape, so confirm which query-building approach a
  given codebase actually uses before writing new repository code.

## Real-world grounding

Wrapping `database/sql` behind a narrow, mockable interface specifically
for unit-testing repositories (rather than integration-testing against a
real database for every test) is a widely-used pattern precisely because
`database/sql`'s own types (`*sql.DB`, `*sql.Rows`) aren't natively
mockable — `sql.Rows` in particular has no exported constructor, which is
why teams wrap the methods they call behind their own interface instead of
trying to mock the standard library's concrete types directly. The
NULL-scanning failure mode is a direct, well-documented consequence of Go's
static typing meeting SQL's nullable columns: `database/sql`'s `Scan`
method explicitly requires a `sql.Null*`-shaped destination (or a type
implementing `sql.Scanner`) for any column that can return `NULL`.

## Verification

- [ ] DB-wrapping interfaces expose only the methods a repository actually
      calls, not the full underlying type's surface
- [ ] Driver errors are translated to a consistent error type at the
      repository boundary, not left as raw driver errors for callers to
      interpret
- [ ] Every nullable column is scanned into a nullable type or a
      `Scanner`-implementing type, never a bare primitive
- [ ] Adding a new repository interface includes adding it to the
      mock-generation list in the same change, not as a follow-up
- [ ] New repository code matches this codebase's actual query-building
      approach, confirmed rather than assumed from a different Go data-
      access convention
