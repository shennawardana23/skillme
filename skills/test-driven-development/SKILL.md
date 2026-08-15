---
name: test-driven-development
description: Drives development with tests written before implementation code, in Go, TypeScript, PHP, or any language. Use when implementing new logic or behavior, fixing a reported bug (write the reproduction test first), modifying existing functionality, or when asked to "add a test," "TDD this," or "prove this works." Use even when the user doesn't say "test" explicitly but describes a bug report or a behavior change.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Test-Driven Development

Write a failing test before writing the code that makes it pass. For bug
fixes, reproduce the bug with a test before attempting a fix. A test that
passes on its first run proves nothing — it may be testing the wrong thing.

## The cycle

1. **RED** — write a test for behavior that doesn't exist yet. Run it. It
   must fail. If it passes immediately, the test is wrong, not the code.
2. **GREEN** — write the minimum code to pass. Don't add behavior the test
   doesn't require.
3. **REFACTOR** — clean up with tests green. Re-run tests after every step.

## The Prove-It Pattern (bug fixes)

Do not edit the implementation until a reproduction test exists and its
failure has been shown. This applies even to a one-line fix that looks
obviously correct — the reproduction test is what proves the bug was real
and what proves the fix addresses it, not the author's confidence in the
diagnosis. Sequence, every time:

1. Write a test that reproduces the reported bug.
2. Run it and show that it fails, with the failure output.
3. Only then change the implementation.
4. Re-run the same test and show it now passes.

```go
// Bug report: "CompleteTask doesn't set CompletedAt"
func TestCompleteTask_SetsCompletedAt(t *testing.T) {
	task := createTask(t, "test")
	completed, err := completeTask(task.ID)
	if err != nil {
		t.Fatalf("completeTask: %v", err)
	}
	if completed.CompletedAt.IsZero() {
		t.Fatal("CompletedAt was not set") // fails first — confirms the bug
	}
}
```

## Go specifics

- Use table-driven tests for multiple input variations — one `t.Run` subtest
  per case, named for the behavior, not "case 1".
- Call `t.Parallel()` inside the subtest closure, and be aware that the
  common table-driven footgun is capturing the loop variable by reference
  across a `range` in Go versions before 1.22; if this repo's `go.mod`
  targets Go < 1.22, shadow the loop variable (`tt := tt`) before calling
  `t.Parallel()`, or Go 1.22+'s per-iteration loop variable semantics make
  it unnecessary.
- Prefer the standard library's `testing` package and `testify/require` for
  fatal assertions; reserve `testify/assert` for checks that should continue
  after failure to surface multiple problems in one run.

## Gotchas

- A green test suite with zero new tests for a bug fix is a signal the fix
  wasn't actually verified — not a signal the code was already correct.
- Mocking the function under test (rather than its dependencies) makes a
  test pass unconditionally; it verifies the mock, not the code.
- `t.Parallel()` subtests that share package-level state (a global cache, an
  env var) produce flaky, order-dependent failures — isolate state per test.
- A migration, refactor, or "safe" rename is not exempt from needing a test:
  if it changes observable behavior and nothing covers that behavior, it is
  untested regardless of how mechanical the change looked.

## Real-world grounding

Knight Capital's 2012 trading system deployed new code to only 7 of 8
production servers; the 8th ran dead code with a repurposed flag, executing
unintended trades for 45 minutes and costing $440M — a change that reached
production with automated tests would have caught the dead-flag reuse before
deployment, not after. The lesson generalizes past that one incident: an
untested code path is a liability precisely because nobody knows it's
untested until it runs in production.

## Test pyramid

Most tests should be small (pure logic, no I/O, milliseconds). Fewer should
be medium (crosses a process boundary — a real Postgres in a container, an
HTTP handler). Fewest should be large (end-to-end, real browser or staging
environment) — reserve these for critical user flows.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "I'll add tests after it works" | Tests written after the fact test the implementation, not the intended behavior — they can't catch a wrong intent. |
| "Too simple to need a test" | Simple code accumulates edge cases. The test is the specification of what "correct" means. |
| "I tested it manually" | Manual verification doesn't persist; the next change can silently break it with no signal. |

See `references/patterns-and-antipatterns.md` for DAMP-vs-DRY test style,
mock-vs-fake-vs-real guidance, and a full anti-pattern table.
