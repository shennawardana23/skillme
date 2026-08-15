## What it does

Explains *why* Go handles errors as ordinary return values instead of
exceptions, and gives the underlying rule for two decisions that come up
constantly: return-an-error vs. panic, and when to wrap vs. compare errors
directly. The defining constraint, from Rob Pike's "Errors are values": an
error is data to inspect and pass along with ordinary code, not a
control-flow event — that framing is what the rest of the skill's rules
follow from, not a list of unrelated style preferences.

## When to reach for it

Reach for this skill when the question is "should this be an error or a
panic" or "why does Go do it this way" rather than "how do I write this
specific error-handling code" — `go-service-idioms` covers the concrete
patterns (wrapping syntax, sentinel errors, goroutine ownership); this
skill covers the reasoning underneath those patterns, useful when a review
comment needs to explain *why* a pattern matters, not just flag that it's
missing.

## Common questions

- **"Is `err.Error() == "not found"` an acceptable way to check what went
  wrong?"** No, and this isn't a style preference — it breaks the moment
  the message wording changes, including from an upstream dependency
  version bump with no other code change. Use `errors.Is`/`errors.As`
  against a defined sentinel or typed error instead.
- **"My goroutine panicked and took down the whole process even though I
  have `recover()` elsewhere in the codebase — why didn't it catch it?"**
  `recover()` only works within the same goroutine where the panic
  occurred. A goroutine launched with `go func() {...}()` that must not be
  allowed to crash the process needs its own `recover()`, not one inherited
  from the goroutine that launched it.
- **"I returned a nil `*MyError` as the function's `error` return type, and
  `err != nil` is somehow true — is that a bug in my code?"** It's a real,
  well-documented Go pitfall, not a bug in the language: a typed nil
  pointer stored in an `error` interface value produces a non-nil
  interface. Always return a literal `nil`, not a nil-valued concrete type.

## It's working if

- Errors from external actors (bad input, network, filesystem, DB) are
  returned, never panicked
- Panics are reserved for broken invariants, not expected failure paths
- Code that branches on error identity uses `errors.Is`/`errors.As`, never
  string comparison
- Goroutines that must not crash the process carry their own `recover()`

## Where it fits

A conceptual/policy skill underneath `go-service-idioms` — reach for that
skill for the concrete patterns, this one when the "why" itself needs
stating, e.g. in a review comment or a team discussion about whether a new
piece of code should return an error or panic.
