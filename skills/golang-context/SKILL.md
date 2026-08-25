---
name: golang-context
description: This skill should be used when the user asks to "add a timeout to this Go call", "why isn't context cancellation propagating", "should I store context in a struct", "use context.WithValue", "compose context timeouts", "detach a context from its parent's cancellation", or is designing how context.Context deadlines, cancellation, and values flow through a Go call chain. go-service-idioms already covers "thread ctx as the first parameter" as a basic default; this skill is for the deeper composition, anti-pattern, and leak questions once more than one context is in play.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Go Context Propagation

`context.Context` carries three things across a call chain: a
cancellation signal, an optional deadline, and request-scoped values. Each
is a distinct mechanism with its own misuse pattern — this skill covers
composing them correctly once a function does more than "accept ctx, pass
it down," which `go-service-idioms` already handles.

## The rules that don't have exceptions

- **Never store a `context.Context` in a struct field.** Pass it as the
  first parameter, named `ctx`, to every function that needs it. A
  context stashed in a struct at construction time goes stale — it
  carries the deadline and cancellation state from *whenever the struct
  was built*, not from the call that's actually running.
  (go.dev/blog/context-and-structs)
- **Never pass `nil`.** If there's genuinely no context yet (test setup,
  a `main` that hasn't wired one up), pass `context.TODO()` — it's a
  greppable marker meaning "this needs a real context," not a silent
  landmine like `nil` would be the first time something calls
  `ctx.Done()` on it.
- **Use `WithValue` only for request-scoped data that crosses API
  boundaries** — a trace ID, a request-scoped auth principal — never for
  passing optional parameters a function could just take as an argument.
  A function whose real inputs are hidden in context values is harder to
  test and harder to read than one that takes them explicitly.

## Composing cancellation and deadlines

`WithCancel`, `WithDeadline`, and `WithTimeout` each return a derived
child context and a `CancelFunc`. Canceling a parent cancels every
context derived from it; canceling a child never affects its parent or
siblings. **Always call the returned `cancel` function** — typically via
`defer cancel()` right after creation — even if the context reaches its
deadline naturally; skipping it leaks the context's internal timer and
goroutine until the parent is canceled, which for a long-lived parent
(e.g. `context.Background()`) may be never. `go vet` flags a `CancelFunc`
that isn't called on some path — don't ignore that warning.

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
if err := callDownstream(ctx); err != nil {
    return fmt.Errorf("call downstream: %w", err)
}
```

A **derived timeout can only shrink, never extend** the deadline it
composes with — `context.WithTimeout(ctx, 10*time.Second)` on a `ctx`
that already has a 2-second deadline is still bound by the 2-second
deadline; the shorter of the two always wins as cancellation propagates
outward. Design timeout budgets top-down (an HTTP handler's overall
deadline, then smaller per-downstream-call timeouts nested under it), not
by tacking a longer timeout onto an already-short one and expecting it to
take effect.

## Recording *why* a context was canceled

`WithCancelCause`, `WithDeadlineCause`, and `WithTimeoutCause` (Go 1.20+)
return a `CancelCauseFunc` that records an error as the cancellation
reason. `context.Cause(ctx)` retrieves it — useful when several different
callers could cancel the same context and a bare `ctx.Err() ==
context.Canceled` doesn't say which one, or why:

```go
ctx, cancel := context.WithCancelCause(parent)
defer cancel(nil)
// elsewhere:
cancel(fmt.Errorf("user canceled upload"))
// later:
if err := context.Cause(ctx); err != nil { ... }
```

## Detaching from cancellation without losing values

`context.WithoutCancel(ctx)` (Go 1.21+) returns a context that never
cancels and has no deadline, but still carries `ctx`'s values — the
correct tool for "this cleanup/audit-log write must finish even though the
request that triggered it was just canceled," instead of the common
workaround of passing `context.Background()` and losing the request's
trace ID and other values along with its cancellation.

## Running code when a context is done

`context.AfterFunc(ctx, f)` (Go 1.21+) arranges for `f` to run in its own
goroutine when `ctx` is canceled or times out, and returns a `stop` func
that prevents `f` from running if called before cancellation — useful for
attaching cleanup to a context's lifetime without a manual `select` on
`ctx.Done()` in every place that needs it.

## Gotchas

- A `select` on a channel operation without also selecting on
  `ctx.Done()` blocks forever past the point the caller gave up — always
  include `case <-ctx.Done(): return ctx.Err()` in any `select` inside a
  loop that could otherwise run long.
- Checking `ctx.Err() != nil` only tells you cancellation *already*
  happened; it doesn't stop work already dispatched to a downstream call
  that isn't itself context-aware (e.g., a blocking C binding, or a
  library call with no context parameter) — that call keeps running to
  completion regardless of what the caller decided.
- `WithValue` keys must not be a built-in type like `string` — an
  unexported custom type (`type ctxKey int`) prevents accidental
  collisions between packages that both use `"userID"` as a plain string
  key.
- Testing code that depends on a deadline: don't `time.Sleep` past it and
  hope — use `context.WithDeadline` with a controllable clock, or assert
  on `errors.Is(err, context.DeadlineExceeded)` directly rather than
  timing the test itself.
- Forgetting `defer cancel()` after `WithTimeout`/`WithCancel` doesn't
  crash anything visibly — it's a slow leak that only shows up as rising
  goroutine/memory counts under load, which makes it easy to miss in a
  code review that isn't specifically looking for it.

## Real-world grounding

`context.WithoutCancel`, `AfterFunc`, and the `*Cause` variants are
documented directly in the standard library (`go doc context`,
pkg.go.dev/context); the struct-field rule is the Go team's own published
guidance at go.dev/blog/context-and-structs, not a house convention.

## Cross-references

- For the baseline "thread ctx as the first parameter, check ctx.Err() in
  loops" default, see `go-service-idioms`.
- For `select`ing on `ctx.Done()` alongside channel operations in
  concurrent code, see `golang-concurrency`.
- For context-aware database calls (`QueryContext`, `ExecContext`), see
  `golang-database`.
