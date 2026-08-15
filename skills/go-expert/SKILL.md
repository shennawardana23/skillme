---
name: go-expert
description: Use when writing or reviewing Go code that starts goroutines, holds a sync.Mutex/RWMutex, or closes channels - to enforce goroutine lifecycle ownership, lock-hygiene, and channel-close-side rules, and to wire up the detection tooling (go test -race, go vet, go.uber.org/goleak) that catches violations tests alone won't. Use when the user asks "review this Go concurrency code," "why is this goroutine leaking," or "is this Go code race-safe."
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Go Concurrency Correctness

Race conditions and goroutine leaks are the Go bugs that pass every test run
and only surface under load, months later. `go-service-idioms` covers
*composing* concurrent work with `errgroup`; this skill covers the
*lower-level ownership rules* that make individual goroutines, mutexes, and
channels safe to combine that way in the first place. Apply these rules
whenever code starts a goroutine, takes a lock, or closes a channel — then
run the detection tooling in the last section to confirm you didn't miss
one, since these bugs are usually invisible without instrumentation.

## Goroutine lifecycle: every goroutine needs an owner and a stop signal

Never start a goroutine without a way to observe when it's done or tell it
to stop — a `context.Context`, a `done` channel, or a `sync.WaitGroup`
someone actually waits on.

```go
// Bad: fire-and-forget, no one can wait for it or cancel it
go worker(items)

// Good: owned, cancelable, and its completion is observable
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
done := make(chan struct{})
go func() {
    defer close(done)
    worker(ctx, items)
}()
```

Pair every `wg.Add(1)` with a `defer wg.Done()` **on the next line inside the
goroutine**, not somewhere later in the function body — a `return` or
`panic` between `Add` and a non-deferred `Done` hangs every caller of
`wg.Wait()` forever.

```go
wg.Add(1)
go func() {
    defer wg.Done() // immediately, before any other logic
    process(item)
}()
```

## Mutex hygiene: lock and defer-unlock are one atomic thought

`defer mu.Unlock()` goes on the line immediately after `mu.Lock()` — never
lock, do conditional work, and defer the unlock later based on a branch.
A goroutine that panics while holding a lock without a deferred unlock
poisons every future caller with a permanent deadlock, since Go mutexes
don't auto-release on panic.

```go
// Bad: unlock is reachable but not guaranteed on every path
mu.Lock()
if invalid(state) {
    mu.Unlock()
    return errors.New("invalid")
}
mutate(state)
mu.Unlock()

// Good: unlock guaranteed regardless of how the function returns
mu.Lock()
defer mu.Unlock()
if invalid(state) {
    return errors.New("invalid")
}
mutate(state)
```

Prefer `sync.RWMutex` over `sync.Mutex` for read-heavy shared state — many
concurrent readers can hold `RLock()` simultaneously, only a writer needs
exclusive `Lock()`. Never copy a struct containing a `sync.Mutex` by value
(passing it to a function, storing it in a slice that reallocates) — `go
vet`'s copylocks check catches this; see Detection below.

## Channel ownership: only the sender closes

Close a channel from the goroutine that sends on it, never from a receiver
— closing from the receiver side, or closing twice from anywhere, panics at
runtime (`close of closed channel` / `send on closed channel`).

```go
// Good: sender owns the channel's lifecycle
func produce(ctx context.Context, out chan<- int) {
    defer close(out) // only the producer closes
    for i := 0; i < 10; i++ {
        select {
        case out <- i:
        case <-ctx.Done():
            return
        }
    }
}
```

Use a buffered channel (`chan struct{}`) as a semaphore to bound
concurrency rather than spawning a goroutine pool with its own internal
queue — it's simpler and its capacity *is* the bound, visible at the
declaration site:

```go
sem := make(chan struct{}, maxConcurrent)
for _, item := range items {
    sem <- struct{}{} // acquire
    go func(item Item) {
        defer func() { <-sem }() // release
        process(item)
    }(item)
}
```

Never use `time.Sleep` to "wait long enough" for another goroutine to finish
or to synchronize state — it is not a synchronization primitive and is
inherently racy under load; use a channel receive, `sync.WaitGroup.Wait()`,
or `sync.Cond` instead.

## Detection: prove the rules held, don't just follow them by eye

These bugs are frequently invisible in code review and in a normal
`go test` run because they depend on scheduler timing. Always pair the
rules above with tooling that actually exercises the race:

```bash
go test -race ./...          # detects concurrent unsynchronized access
go vet ./...                 # catches sync.Mutex copies (copylocks), among others
```

For goroutine leaks specifically — a goroutine that outlives the test that
started it, holding a channel or context reference — add
[`go.uber.org/goleak`](https://pkg.go.dev/go.uber.org/goleak) to the test's
`TestMain`:

```go
func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

`-race` only detects a data race if the racing goroutines actually run
concurrently *during that specific test execution* — a race that only
manifests under production load or a specific interleaving can pass
`-race` cleanly in CI. Treat a clean `-race` run as evidence, not proof.

## Gotchas

- A `select` with a `default` case turns a blocking channel op into a
  non-blocking poll — if you meant to wait, dropping `default` back in is
  usually the fix for a goroutine that reports success without actually
  doing the work.
- `sync.WaitGroup.Add` called *inside* the goroutine it's tracking (instead
  of before `go func(){}()`) races with `Wait()` potentially returning
  before the `Add` even runs — call `Add` in the parent goroutine, before
  spawning.
- A `nil` channel blocks forever on both send and receive — this is
  sometimes used deliberately in a `select` to disable a case, but if it's
  accidental (a channel field never initialized), the symptom is a goroutine
  that just hangs with no error, easily mistaken for a deadlock elsewhere.

## Real-world grounding

Uber open-sourced `go.uber.org/goleak` specifically because goroutine leaks
pass ordinary test suites silently — a leaked goroutine holds no assertion
that fails, it just accumulates, and the industry-standard way to catch it
is to assert on goroutine count at test teardown rather than to try to spot
it by reading code.
