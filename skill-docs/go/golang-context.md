## What it does

Guides composing `context.Context` cancellation, deadlines, and values
correctly once more than one context is in play. The defining
constraint: a context stashed in a struct field goes stale relative to
whatever call is actually running — it must be passed explicitly as a
parameter to every function that needs it, never stored at construction
time.

## When to reach for it

`go-service-idioms` already states the baseline — thread `ctx` as the
first parameter, check `ctx.Err()` in loops. Reach for this skill once
the question moves past that baseline: composing nested timeouts,
recording *why* something was canceled, detaching a cleanup task from a
canceled parent's lifetime, or reviewing code that stores a context
somewhere it shouldn't.

## Common questions

- **"I set a 10-second timeout on a context that already had a 2-second
  deadline from its parent — why did it still cancel at 2 seconds?"** A
  derived timeout can only shrink the effective deadline, never extend
  it — the shorter of any two composed deadlines always wins. Design
  timeout budgets top-down from the outermost caller, not by tacking a
  longer timeout onto an already-short one.
- **"How do I know *why* a context was canceled when several different
  callers could have canceled it?"** `context.WithCancelCause` (and the
  `Deadline`/`Timeout` variants, Go 1.20+) let you attach an error as the
  cancellation reason, retrievable via `context.Cause(ctx)` — a bare
  `ctx.Err() == context.Canceled` can't distinguish who canceled it or
  why.
- **"A cleanup task needs to finish even though the request context that
  triggered it was just canceled, but it still needs the request's trace
  ID — what do I pass it?"** `context.WithoutCancel(ctx)` (Go 1.21+) — it
  keeps `ctx`'s values but is immune to its cancellation, unlike swapping
  to `context.Background()`, which loses the values too.
- **"Why does `go vet` warn about my `cancel` function?"** `go vet`
  specifically checks that a `CancelFunc` from `WithCancel`/`WithTimeout`/
  `WithDeadline` is called on every control-flow path — an unreferenced
  or conditionally-called one leaks the context's internal timer and
  goroutine until the parent context cancels, which for a long-lived
  parent may be never.

## It's working if

- No `context.Context` is ever stored as a struct field
- Every `WithCancel`/`WithTimeout`/`WithDeadline` has an unconditional
  `defer cancel()` right after creation
- `WithValue` carries only request-scoped data (trace IDs, auth
  principals) that crosses an API boundary, never an optional function
  parameter
- A `select` inside any loop that could run long includes `ctx.Done()`

## Where it fits

Deepens `go-service-idioms`'s baseline context-propagation rule. Feeds
`golang-concurrency`'s `select`-based goroutine shutdown patterns and
`golang-database`'s context-aware query calls (`QueryContext`,
`ExecContext`).
