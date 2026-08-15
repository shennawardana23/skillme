---
name: error-handling-philosophy
description: Use when the user asks to "review error handling in this code", "should this panic or return an error", "design an error type", "how should this function report failure", or is writing Go code with error returns, panics, or wrapping. Guides applying Go's "errors are values" philosophy (Rob Pike) — explicit error returns over exceptions, wrapping with %w for context, and reserving panic for truly unrecoverable programmer errors.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "philosophy"
---

# Error Handling Philosophy

Go's error handling is a deliberate design choice, not an accident of a
missing exception mechanism. "Errors are values" (Rob Pike, 2015,
go.dev/blog/errors-are-values) means an error is a normal return value you
inspect, transform, and pass along using ordinary code — not a control-flow
event that unwinds the stack past code that didn't ask to handle it.

## The core distinction: error values vs. exceptions

Exception-based languages let a `throw` skip past every intermediate frame
until a `catch` matches — the calling code between throw and catch doesn't
need to mention the possibility of failure at all. This is convenient but
means you can't tell from a function signature whether it can fail, or
what could go wrong, without reading its full implementation (and its
callees', transitively).

Go's explicit `(result, error)` return makes failure a first-class,
visible part of every function signature. The tradeoff: more boilerplate
(`if err != nil { return err }` at every call site) in exchange for a
caller that can never accidentally ignore an error path — the compiler
won't stop you from ignoring it, but the pattern makes ignoring it a
visibly deliberate act (an unused `_` or a dropped return value) rather
than an invisible omission.

## When to return an error vs. when to panic

- **Return an error** for anything an *external* actor can cause: bad
  input, a network failure, a missing file, a database constraint
  violation, a downstream API returning 4xx/5xx. These are expected,
  recoverable conditions the caller should decide how to handle (retry,
  surface to a user, log and continue).
- **Panic** only for programmer errors that indicate the program's
  invariants are already broken and continuing would produce worse,
  silent corruption: an index genuinely out of bounds due to a logic bug,
  a nil pointer dereference from a broken invariant, a `must`-prefixed
  helper whose precondition was violated by calling code (e.g.,
  `regexp.MustCompile` on a pattern that's a compile-time constant, so a
  failure can only mean a bug in *this* code, never bad input).
- **Never panic on user input, network responses, or file contents** —
  those are exactly the "external actor" cases errors exist for. A JSON
  parse failure on a request body is an error, not a panic.
- **`recover()` is a last-resort safety net at a boundary** (e.g., an HTTP
  server's top-level handler wrapper, so one panicking request doesn't
  take down the whole process) — not a substitute for proper error
  returns in the code beneath that boundary.

## Wrapping and inspecting errors

- **Wrap with `%w`, not `%v`, when the caller might need to check the
  underlying error type**: `fmt.Errorf("loading hotel %d: %w", id, err)`.
  `%w` preserves the chain so `errors.Is`/`errors.As` can still find the
  original sentinel or type further up the call stack; `%v` flattens it to
  a string and destroys that ability permanently.
- **Add context on the way up, don't just repeat "error occurred."** Each
  wrap should say what this layer was doing when it failed — `"querying
  reservations for hotel %d: %w"` is useful in a log; `"failed: %w"`
  repeated at every layer is not.
- **Use sentinel errors (`var ErrNotFound = errors.New(...)`) or typed
  errors** when callers need to branch on *which* error occurred, checked
  via `errors.Is`/`errors.As` — never via string-matching `err.Error()`,
  which breaks the moment the message wording changes.
- **Don't wrap and re-wrap the same information at every single layer** —
  if five layers all just add "failed to X: %w" with no new information
  each time, the resulting chain is noise, not context.

## Gotchas

- Comparing errors with `err.Error() == "some string"` or
  `strings.Contains(err.Error(), "not found")` breaks silently the moment
  the error message text changes (including from an upstream dependency
  bump) — always prefer `errors.Is`/`errors.As` against a defined sentinel
  or type.
- Swallowing an error with `_ = someCall()` because "it probably won't
  fail" is a decision that should be visible in code review, not silent —
  if it's truly safe to ignore, a one-line comment saying why belongs next
  to the `_`.
- A goroutine that panics and isn't recovered inside that specific
  goroutine crashes the entire process, even if the rest of the program
  has proper error handling — `recover()` only works within the same
  goroutine where the panic occurred, so any goroutine launched with `go
  func() {...}()` needs its own recover if it must not be allowed to take
  the whole process down.
- Returning a typed nil inside an `error` interface value (e.g., returning
  a nil `*MyError` pointer as the `error` return type) produces a non-nil
  interface that fails `err != nil` checks in confusing ways — this is a
  well-known Go pitfall, not a hypothetical one; always return a literal
  `nil`, not a nil-valued concrete type, when there's no error.
- Retrying a wrapped error without checking whether the *underlying* cause
  is retryable (e.g., retrying a wrapped `context.Canceled` or a wrapped
  validation error) wastes cycles and can mask a bug — check the
  underlying error's nature via `errors.Is`, not just "an error happened,
  so retry."

## Real-world grounding

Rob Pike's 2015 Go blog post "Errors are values" (go.dev/blog/errors-are-values)
is the canonical, publicly documented articulation of this philosophy from
one of Go's original designers, explicitly framing Go's approach as a
deliberate alternative to exception-based error handling rather than a
limitation to work around.

## Verification

- [ ] Errors from external actors (input, network, filesystem, DB) are
      returned, never panicked
- [ ] Panics are reserved for broken invariants/programmer errors only
- [ ] Wrapping uses `%w`, and callers that need to branch use
      `errors.Is`/`errors.As`, never string comparison on `Error()`
- [ ] Every ignored error (`_ = call()`) has a reason visible in the code
- [ ] Goroutines that must not crash the process have their own `recover()`
