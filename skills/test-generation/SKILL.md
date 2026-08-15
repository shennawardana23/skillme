---
name: test-generation
description: Generates comprehensive unit, table-driven, and concurrency tests for a given function or module — happy paths, error paths, boundary values, and mocked dependencies — following Go and TypeScript idioms. Use when a function or file needs test coverage, a bug fix needs a regression test, or existing tests only cover the happy path.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Test Generation

For any given code, produce all of the following — not just the happy-path test that's easiest to write.

## What to generate

**1. Unit tests, per exported function/method:**
- Happy path — normal inputs produce expected outputs
- Error paths — every error condition is triggered and its content verified (not just `err != nil`)
- Nil/zero inputs — behavior with empty or zero values
- Boundary values — min, max, empty slice, single element, max slice

**2. Table-driven tests (Go)** — for any function with multiple input variants:

```go
func TestFoo(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {"happy path", validInput, expectedOutput, false},
        {"nil input", nil, zero, true},
        {"empty string", "", zero, true},
    }
    for _, tc := range tests {
        tc := tc // capture range variable (needed pre-Go 1.22; see Gotchas)
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            got, err := Foo(tc.input)
            if (err != nil) != tc.wantErr {
                t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
            }
            if !tc.wantErr && got != tc.want {
                t.Errorf("got %v, want %v", got, tc.want)
            }
        })
    }
}
```

**3. Concurrency tests (Go)** — for any function touching shared state, channels, or goroutines:
- Design the test to run under `go test -race`, not just to pass without it.
- Use `sync.WaitGroup` to synchronize multiple goroutines exercising the code concurrently.
- Assert the invariant holds after concurrent access, not just that it doesn't panic — a race can corrupt data silently without crashing.

**4. Mock patterns** — for external dependencies (database, HTTP, filesystem):
- Mock through an interface (hand-written or `golang/mock`/`uber-go/mock`), not by monkey-patching a concrete type.
- Cover both the success and error return from the mock.
- Verify the mock was called the expected number of times when call count is part of the contract (e.g., "retries exactly twice").

## Language-specific rules

**Go:** `package foo_test` for black-box testing, `package foo` when internals must be reached; `t.Parallel()` on independent subtests; `t.Helper()` in test helper functions; standard library `testing` only, unless `testify` (or similar) is already a dependency in `go.mod` — don't introduce a new test-only dependency; run with `-race -count=1 -coverprofile=cover.out`.

**TypeScript/JavaScript:** Jest or Vitest, whichever the project already uses; `describe`/`it` with `expect` assertions; mock with `jest.fn()`/`vi.fn()`; test async code with `async/await`, not a raw `.then()` chain that can swallow a rejected assertion.

## Output format

```go
// <filename>_test.go
package ...

import (...)

func TestFoo_HappyPath(t *testing.T) { ... }
func TestFoo_NilInput(t *testing.T) { ... }
func TestFoo_TableDriven(t *testing.T) { ... }
func BenchmarkFoo(b *testing.B) { ... }
```

Always include the full file with `package` and `import` declarations — a snippet the caller has to assemble themselves gets skipped. Add at least one benchmark for any performance-sensitive function (anything in a hot path, anything the code review flagged for allocation concerns).

## Gotchas

- The `tc := tc` range-variable capture is only necessary on Go before 1.22 — as of Go 1.22, `for _, tc := range tests` creates a new `tc` per iteration automatically, and adding the redundant capture line in newer code is harmless but a signal the generator didn't check the project's Go version. Check `go.mod`'s `go` directive before deciding whether to include it.
- A table-driven test with `t.Parallel()` inside the subtest but *not* also called on the parent test function can still serialize unexpectedly if the parent test does setup between subtests — `t.Parallel()` only takes effect once the parent test function itself returns to the test runner, so subtests queued after a blocking parent operation don't get the parallelism the code implies.
- Testing "the mock was called" without asserting on *what it was called with* verifies plumbing, not behavior — a retry-twice mock assertion that only checks `Times(2)` will pass even if both calls used the wrong argument.
- A concurrency test that passes without `-race` and is never run with it is not evidence of thread-safety — it's evidence the test didn't check. Every concurrency test generated here should be paired with an explicit note that it must be run with `-race`, not left as an assumed default.

## Real-world grounding

Go's table-driven test pattern is the idiomatic style used throughout the Go standard library itself (e.g. `strconv`, `net/http`, `encoding/json` test files all follow this exact `tests := []struct{...}` shape) — it's not a stylistic preference invented for this skill, it's the convention the language's own maintainers use specifically because it makes adding a new input case a one-line diff rather than a new function.

## Verification

- [ ] Every exported function has a happy-path test and at least one error-path test
- [ ] Error-path tests assert on error *content*, not just non-nil
- [ ] Boundary values (empty, nil, single-element, max) are covered where the type allows them
- [ ] Concurrency-touching tests are run with `-race`, not just written
- [ ] Mocks assert on call arguments, not only call count
- [ ] The full test file (package + imports) is included, not a bare function snippet
