---
name: golang-concurrency
description: This skill should be used when the user asks to "review this Go code for goroutine leaks", "should I use a channel or a mutex", "use sync.WaitGroup vs errgroup", "audit this Go code for race conditions", "write a Go pipeline with fan-out/fan-in", "add a bounded worker pool in Go", or is choosing between sync primitives, debugging a goroutine leak, or designing a concurrent pipeline. For the one-paragraph "use errgroup for fanning out work" default in ordinary service code, use go-service-idioms instead; this skill is for the deeper decision (which primitive, how the goroutine exits, leak/race auditing) once concurrency itself is the question.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Go Concurrency

Go's concurrency model is a liability until proven necessary: every
goroutine you spawn is a resource with an owner, an exit condition, and a
failure mode you must account for. Reach for this skill once concurrency
itself is the design question — which primitive to use, whether a
goroutine can leak, whether shared state is actually protected — not for
the routine "fan this out" case, which `go-service-idioms` already covers
with `errgroup`.

## The one question to ask before spawning a goroutine

**How does it exit?** Every goroutine needs one of: a `context.Context`
whose cancellation it observes in a `select`, a channel it reads until
close, or a `sync.WaitGroup`/`errgroup` the caller waits on. A goroutine
with none of these is a leak the moment its one caller stops waiting for
it — `runtime.NumGoroutine()` climbing over time in production is the
symptom, not the cause.

## Choosing a primitive

| Scenario | Use | Why |
|---|---|---|
| Passing data between goroutines, transferring ownership | channel | Communicates "you own this now," not just "here's a pointer" |
| Waiting for goroutines, no error to collect | `sync.WaitGroup`, or Go 1.25+ `wg.Go(f)` | `wg.Go` starts the goroutine and calls `Add`/`Done` for you — no more forgetting `Add` before `go` |
| Waiting + collecting the first error | `golang.org/x/sync/errgroup` | `g.Go(func() error {...}); g.Wait()` returns the first non-nil error |
| Waiting + canceling siblings on first error | `errgroup.WithContext(ctx)` | The returned `ctx` cancels for every `g.Go` call the moment one returns an error |
| Bounding concurrency without a hand-rolled semaphore | `errgroup.Group.SetLimit(n)` | Same `Go`/`Wait` API, caps in-flight goroutines |
| Protecting a shared struct's fields | `sync.Mutex` / `sync.RWMutex` | Simple critical sections; never hold one across a blocking I/O call |
| A single counter or flag | `sync/atomic` typed atomics (`atomic.Int64`, `atomic.Bool`, Go 1.19+) | Lock-free, cheaper than a mutex for one word |
| A map read far more than it's written | `sync.Map` | Optimized for read-heavy access; a plain `map` with concurrent read+write panics at runtime, it does not silently corrupt |
| Running init exactly once | `sync.Once` (or Go 1.21+ `sync.OnceFunc`/`OnceValue`) | Safe lazy init without a manual `initialized bool` + mutex |
| Deduplicating identical concurrent calls | `golang.org/x/sync/singleflight` | Collapses N simultaneous callers into one underlying call |

## Channels: rules that prevent real bugs

- **Only the sender closes a channel.** A receiver that closes it, or a
  second sender that closes an already-closed channel, panics.
- **Send copies or immutable values, not pointers to mutable data** — a
  pointer on a channel is shared memory with an extra step, which defeats
  the point of communicating by channel instead of by mutex.
- **Declare channel direction in signatures** (`chan<- T` for send-only,
  `<-chan T` for receive-only) so the compiler catches misuse at the call
  site instead of at runtime.
- **Always include `ctx.Done()` in a `select`** alongside channel
  operations in a loop; without it, a goroutine blocked on a channel that
  will never receive again outlives the caller that gave up on it.
- **Don't call `time.After` inside a hot loop** — each call allocates a
  new timer that isn't garbage-collected until it fires. Create one
  `time.NewTimer` outside the loop and call `Reset` inside it instead.

## Fan-out/fan-in and worker pools

Bound the number of workers to a fixed count reading from a shared input
channel rather than spawning one goroutine per item when the item count is
unbounded (user input, a queue, a directory walk):

```go
func process(ctx context.Context, items []Item, workers int) error {
    g, ctx := errgroup.WithContext(ctx)
    in := make(chan Item)

    g.Go(func() error {
        defer close(in)
        for _, it := range items {
            select {
            case in <- it:
            case <-ctx.Done():
                return ctx.Err()
            }
        }
        return nil
    })

    for i := 0; i < workers; i++ {
        g.Go(func() error {
            for it := range in {
                if err := handle(ctx, it); err != nil {
                    return err
                }
            }
            return nil
        })
    }
    return g.Wait()
}
```

This shape gives every stage a clear exit (`close(in)` from its single
sender, `ctx.Done()` in the feeder's select, `range` draining the channel
in each worker until close) — the three failure modes below are all
missing one of those.

## Gotchas

- `recover()` only catches a panic within the same goroutine — a
  goroutine started with `go func() {...}()` that panics and isn't
  recovered *inside that function* crashes the whole process, even if the
  rest of the program handles panics correctly elsewhere.
- `wg.Add` must happen before the `go` statement, not inside the spawned
  goroutine — `Wait` can return before a late `Add` ever registers,
  making the wait silently incomplete. `wg.Go(f)` (Go 1.25+) removes this
  class of bug entirely by doing `Add`/`go`/`Done` atomically.
- `sync.RWMutex.RLock` held by one goroutine while another calls `Lock`
  deadlocks if the first goroutine tries to upgrade its own `RLock` to
  `Lock` without releasing it first — never upgrade a read lock in place.
- Reusing a `sync.Pool` object without resetting it before `Put` leaks
  stale data into the next `Get` caller; always zero or reset before
  returning an object to the pool.
- `go test -race` only catches races that actually execute during that
  run — a race-free test run is not proof of a race-free program, only
  evidence for the code paths the test exercised. Run `-race` in CI on
  every PR, not just once.
- Track leaks in tests with `go.uber.org/goleak`; a leaked goroutine from
  one test can silently keep running during unrelated later tests in the
  same process and skew their timing or resource use.

## Real-world grounding

`sync.WaitGroup.Go` and the typed atomics under `sync/atomic` are
documented in the Go standard library (`go doc sync.WaitGroup`,
pkg.go.dev/sync); `golang.org/x/sync/errgroup` and `singleflight` are the
Go team's own supplementary concurrency packages, not a third-party
convention layered on top.

## Cross-references

- For the routine "fan this out with errgroup" default in ordinary
  service code, see `go-service-idioms`.
- For `context.Context` cancellation and timeout composition feeding the
  `select` statements above, see `golang-context`.
- For allocation-driven performance work once a concurrent path is
  correct, see `golang-patterns` and `golang-testing` (benchmarking).
