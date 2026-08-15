## What it does

Provides the default idiomatic patterns for Go backend code in this
catalog's domain — error handling, context propagation, concurrency,
package layout, and table-driven tests. The defining constraint: these are
defaults for *new* code, not a mandate to rewrite existing code that
already follows a different but internally consistent convention — match
what's already there first.

## When to reach for it

Reach for this skill whenever writing or reviewing Go service code and
either there's no existing local convention to match, or the review is
specifically checking whether Go idioms (error wrapping, context handling,
goroutine ownership) are followed correctly. It's a general-purpose Go
skill, not tied to a specific framework — `postgres-patterns` and
`postgres-hotel-partitioning` cover the database-query layer this skill's
examples touch but don't teach themselves.

## Common questions

- **"Should this function panic or return an error?"** Return an error for
  anything an external actor can cause (bad input, a failed network call, a
  missing file) — reserve `panic` for programmer errors that indicate a
  broken invariant (a nil pointer from a logic bug, a `MustCompile`-style
  helper whose precondition was violated by the calling code itself, never
  by external input).
- **"Do I need the `id := id` loop-variable copy inside `g.Go(func() error
  { ... })`?"** Only if the code targets Go before 1.22 — since 1.22, each
  loop iteration gets its own variable, so the copy is unnecessary (though
  harmless) on current Go versions.
- **"Is it fine to wrap the same error again at every layer it passes
  through?"** No — wrap once with enough context to be useful, and log once
  at the boundary that actually handles the failure (an HTTP handler, a CLI
  entrypoint, a queue consumer). Re-wrapping "failed: %w" at every
  intermediate layer produces a noisy chain, not more information.
- **"Why prefer `errgroup` over a hand-rolled `WaitGroup` and error
  channel?"** `errgroup.WithContext` gives you cancellation propagation (the
  first failing goroutine cancels the shared context, so siblings can stop
  early) for free — a hand-rolled version usually either skips that or
  reimplements it inconsistently across a codebase.

## It's working if

- Every function that fails on external input returns an error; panics are
  reserved for broken invariants only
- Error wraps use `%w` and add real context, not a repeated generic message
  at every layer
- Every launched goroutine has a clear owner for its completion or error
- Table-driven tests use `t.Run` per case with meaningful failure output

## Where it fits

Standalone, general-purpose Go skill — the base layer other Go-adjacent
skills in this catalog build on (`postgres-patterns`,
`postgres-hotel-partitioning`, `adk-go-agent-builder`). Cross-reference
`error-handling-philosophy` for the deeper "why" behind the error-vs-panic
distinction this skill states as a rule of thumb.
