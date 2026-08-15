---
name: go-service-idioms
description: This skill should be used when the user asks to "write a Go function", "implement this in Go", "structure this Go service", "add error handling in Go", "handle context cancellation in Go", "write a Go worker pool", "add a table-driven test", "review this Go code for idioms", or otherwise writes or reviews Go backend/service code. Provides idiomatic Go patterns for error handling, context propagation, concurrency, package layout, and testing.
metadata:
  version: "0.1.0"
  category: "go"
---

# Go Service Idioms

Provide idiomatic, production-grade Go patterns for backend services: error
handling, context propagation, concurrency, package layout, and testing.
Apply these defaults whenever writing or reviewing Go code in this catalog's
domain, unless the surrounding project's existing conventions say otherwise —
match existing patterns first, fall back to this skill's defaults for new
code.

## Error handling

Wrap errors with `%w` and enough context to debug without a stack trace,
never discard an error silently, and never use `panic` for expected failure
paths (bad input, network failure, not-found). Reserve `panic` for programmer
errors (nil deref, invariant violation) that should crash fast in
development.

```go
cfg, err := loadConfig(path)
if err != nil {
    return nil, fmt.Errorf("load config from %s: %w", path, err)
}
```

Define sentinel errors (`var ErrNotFound = errors.New("not found")`) or
typed errors when callers need to branch on failure kind, and check them
with `errors.Is` / `errors.As` — never string-match an error message.

Do not wrap an error and then also log it at every layer it passes through;
wrap once, log once at the boundary that handles it (an HTTP handler, a CLI
entrypoint, a queue consumer).

## Context propagation

Thread `context.Context` as the first parameter of any function that does
I/O (network call, DB query, file read over a slow mount) so callers can
cancel or bound it with a deadline. Never store a `context.Context` in a
struct field — pass it explicitly through the call chain.

```go
func fetchReservation(ctx context.Context, db *sql.DB, id string) (*Reservation, error) {
    row := db.QueryRowContext(ctx, "SELECT ... WHERE id = $1", id)
    ...
}
```

Check `ctx.Err()` (or select on `ctx.Done()`) inside loops that could run
long, and return promptly with `ctx.Err()` wrapped, rather than polling with
`time.Sleep`.

## Concurrency

Prefer `golang.org/x/sync/errgroup` over hand-rolled `sync.WaitGroup` +
manual error channel when fanning out concurrent work that can fail:

```go
g, ctx := errgroup.WithContext(ctx)
for _, id := range ids {
    id := id // no longer required since Go 1.22 (each loop iteration gets
             // its own variable), but harmless; only add if targeting <1.22.
    g.Go(func() error {
        return processOne(ctx, id)
    })
}
if err := g.Wait(); err != nil {
    return fmt.Errorf("process ids: %w", err)
}
```

Every goroutine a function starts must have a clear owner responsible for
letting it finish or observing its error — do not fire-and-forget a
goroutine whose panic or leak has nowhere to surface. Bound worker pools with
a fixed number of workers reading from a channel rather than spawning one
goroutine per item when the item count is unbounded (user input, a queue).

## Package layout

Follow the standard layout: `cmd/<binary>/main.go` for entrypoints,
`internal/` for code not meant to be imported by other modules, and
top-level packages only for code that is genuinely a public, reusable API.
Name packages for what they provide, not what they contain (`package hotel`,
not `package models` or `package utils`). Avoid import cycles by keeping
dependencies flowing in one direction (e.g., `internal/api` depends on
`internal/store`, never the reverse).

Accept interfaces, return concrete structs: a constructor returns `*Store`,
but a function that only needs to read data takes a narrow interface it
defines itself, scoped to what it needs — do not force callers to depend on
a large interface for one method.

## Testing

Use table-driven tests with `t.Run` per case so failures report the specific
case name, and `t.Parallel()` inside each subtest when cases are independent:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        want     int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -1, -2},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

Prefer the standard library's `testing` package and `cmp.Diff`
(`github.com/google/go-cmp`) for deep comparisons over adding a heavier
assertion library, unless the project already depends on one.

## Additional resources

For a longer discussion of these patterns with more examples, consult
`references/idioms.md`.
