## What it does

Guides choosing among Go's concurrency primitives (channels, mutexes,
atomics, `errgroup`) and auditing concurrent code for goroutine leaks and
races. The defining constraint: every goroutine needs a clear owner and
exit condition before it's spawned — a goroutine with no context,
channel, or `WaitGroup`/`errgroup` tracking it is a leak from the moment
its one caller stops waiting.

## When to reach for it

Reach for this skill once concurrency itself is the open question — which
primitive fits, whether a piece of code can leak a goroutine, whether
shared state is actually protected — not for the routine "fan this out"
case. `go-service-idioms` already covers that routine case with one
`errgroup` example; this skill is the deeper reference once more than
that one pattern is in play.

## Choosing among primitives

Most confusion here isn't "how do channels work" — it's picking the
right tool among several that all technically work: a `sync.Mutex` vs.
`sync/atomic` vs. `sync.Map` for shared state, or `sync.WaitGroup` vs.
`errgroup` for waiting on goroutines. The skill's primitive-selection
table exists because each of these has a narrow correct use case, and
picking the wrong one compiles and often runs fine under light load
before failing under contention or scale.

## Common questions

- **"Do I still need `id := id` inside a `for _, id := range ids { go
  func() {...}() }` loop?"** Only if targeting Go before 1.22 — since
  1.22, each loop iteration gets its own variable, so the classic
  loop-variable-capture bug no longer applies to code on current Go.
- **"Is `sync.WaitGroup.Go` the same as calling `Add(1)` then `go`?"**
  Yes, but it does both atomically (Go 1.25+) — it removes the specific
  bug where `Add` gets called inside the spawned goroutine instead of
  before it, letting `Wait()` return before every goroutine has
  registered.
- **"Why does a concurrent map read/write panic instead of just
  corrupting data?"** Go's runtime actively detects concurrent map
  access and crashes the program on purpose — it's a deliberate fail-fast
  design, not an accident, precisely because silent corruption would be
  far harder to debug than an immediate crash.
- **"My `go test -race` run passed — is the code race-free?"** Only for
  the code paths that test actually executed. The race detector proves
  the absence of a race in what ran, not in the program in general — it
  needs to run in CI on every PR, not once, to keep being useful.

## It's working if

- Every spawned goroutine has an identifiable exit path (context
  cancellation, channel close, or a `WaitGroup`/`errgroup` the caller
  waits on)
- Channels are only ever closed by their sender
- `go test -race` runs in CI, not just locally on demand
- A bounded worker count is used instead of one goroutine per item when
  the item count is unbounded (user input, a queue)

## Where it fits

Deepens `go-service-idioms`'s one-paragraph concurrency default.
Feeds `select`-on-`ctx.Done()` patterns from `golang-context`, and hands
off to `golang-testing` once a concurrent path needs benchmark-driven
performance work rather than correctness review.
