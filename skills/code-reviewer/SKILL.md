---
name: code-reviewer
description: Runs a fixed-category severity scan over a diff, file, or snippet — correctness, security, error handling, concurrency, performance, readability, tests — and emits a CRITICAL/HIGH/MEDIUM/LOW/INFO-graded report. Use when asked to "scan this code", grade a change by severity, or review Go code for concurrency bugs (data races, goroutine leaks, deadlocks) alongside the usual review axes.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Code Reviewer

A mechanical, checklist-driven scan: walk every category below in order, record every match with a severity, then emit the report. This is a scanner, not a workflow — it does not cover change sizing, review etiquette, or process (see `code-review-and-quality` for that). It does not walk a security-priority checklist end to end (see `security-review` for that). Its distinguishing category is concurrency, which the other review skills in this catalog do not cover.

## Category checklist

Work through every category for every changed file. Don't stop at the first finding — a diff with one CRITICAL often has MEDIUM issues too.

### Correctness
- Logic bugs, off-by-one errors, wrong boundary conditions
- Missing nil/zero-value checks and unhandled edge cases
- Incorrect assumptions about a library's behavior — check the actual doc/signature, not intuition (e.g., does `strings.TrimSuffix` no-op silently if the suffix doesn't match? yes — a common mis-assumption)

### Security
- SQL/command/template injection via unsanitized input
- Hardcoded secrets, credentials, or API keys
- Insecure cryptographic choices (MD5/SHA1 for anything security-sensitive, ECB mode, `math/rand` for tokens)
- Missing authentication or per-resource authorization on new endpoints
- Path traversal, open redirects, SSRF-shaped calls (server fetches a user-influenced URL)

### Error handling
- Errors silently swallowed (`_ = err`) or logged but not propagated
- Missing cleanup on error paths — `defer f.Close()` placed *before* the error check, not after
- Panic recovery missing in goroutines (an unrecovered panic in a goroutine crashes the whole process, not just that goroutine)
- Errors wrapped with `%w` losing sentinel-comparability, or not wrapped at all, losing context

### Concurrency (Go)
- **Data races**: map or slice read/write from multiple goroutines without a mutex or `sync.Map`. Don't just read the code — run `go test -race` and check the diff's package is covered.
- **Goroutine leaks**: a goroutine started with no way for it to ever return — no `context.Context`, no `done` channel, no bound on the channel it's blocked sending to.
- **Deadlocks**: locks acquired in inconsistent order across two code paths; a channel send/receive with no counterpart and no `select` with a `ctx.Done()` case.
- `mu.Lock()` not immediately followed by `defer mu.Unlock()` on the next line — anything between them is a window where an early return skips the unlock.
- A `WaitGroup.Add` call issued inside the goroutine it's supposed to be counting, instead of before `go func(){...}()` — a benign-looking reorder that reintroduces the race `go vet`'s `loopclosure`/`-race` won't always catch on the first run.

### Performance
- Unnecessary allocations in a hot path (building a new slice/map per call where reuse is possible)
- O(n²) via a nested loop over the same slice where a map/set lookup gives O(n)
- Missing `context.Context` cancellation propagation into blocking calls
- Large structs passed by value in a hot path instead of by pointer

### Readability & maintainability
- Function length over ~60 lines with multiple responsibilities — should be decomposed
- Names that don't convey intent (`x`, `tmp`, `data`, `result` with no qualifying context)
- Missing doc comment on exported symbols
- Duplicated logic that should be extracted to a shared helper

### Test coverage
- Missing happy-path or error-path test for new behavior
- Missing table-driven tests where the function has multiple input variants
- A test that doesn't assert on error *content*, only `err != nil`
- No `-race` coverage for new concurrent code

## Severity levels

| Level | Meaning | Required before merge |
|---|---|---|
| CRITICAL | Data loss, security breach, crash, confirmed data race | Yes — block merge |
| HIGH | Incorrect behavior, likely race condition, goroutine leak | Yes — block merge |
| MEDIUM | Performance regression, missing test, error-handling gap | Strongly recommended |
| LOW | Style, naming, minor readability | Nice-to-have |
| INFO | Observation, question, non-blocking suggestion | No action required |

## Output format

```
## Summary
<2-3 sentence overall assessment>

## Critical Issues
<file:line, severity, description, concrete fix — every entry needs a fix>

## Improvements
<MEDIUM/HIGH items that should be fixed>

## Suggestions
<LOW/INFO items — optional>

## Quality Score: N/10
<brief justification tied to the categories above>
```

Be direct: name the category and severity for every finding, always propose a concrete fix (never just "this looks off"), and acknowledge at least one thing done well.

## Gotchas

- `go vet` and `staticcheck` do not detect most data races — only `go test -race` runs the actual race detector against real execution. A clean `go vet` pass is not evidence of race-freedom; always ask "was this run under `-race`?" before clearing a concurrency finding.
- `sync.WaitGroup.Add` must happen in the goroutine that starts the child goroutines, before the `go` statement — moving it inside the new goroutine (a change that "looks" equivalent) creates a race between `Add` and `Wait` that only manifests under scheduling pressure, so it usually passes CI and fails in production.
- A large struct passed by value isn't automatically a performance bug — for structs under ~3 machine words, the Go compiler often keeps them in registers and copying is cheaper than the pointer indirection. Flag it only when profiling data or a genuinely large struct (10+ fields, embedded slices/maps) supports it — don't cite this pattern from code reading alone.
- Buffered channels can mask a goroutine leak in testing (the leaked goroutine's send never blocks because the buffer absorbs it) while still leaking in production once the buffer fills — check the channel's actual capacity against the worst-case fan-in, not just "it's buffered so it's fine."

## Real-world grounding

Uber's public "go-torch"/production incident write-ups and the Go team's own `sync.WaitGroup` documentation both call out the `Add`-inside-goroutine ordering bug as the most common WaitGroup misuse in real codebases — it compiles, passes most test runs, and only shows up under load, which is exactly the profile of a concurrency defect a syntactic checklist alone will miss without also asking "was `-race` run."

## Verification

- [ ] Every CRITICAL/HIGH finding names a category and cites file:line
- [ ] Every finding has a concrete fix, not just a description of the problem
- [ ] Concurrency-touching diffs were checked against `-race` output, not just read
- [ ] At least one positive observation is included
- [ ] The quality score is justified by the findings above it, not asserted alone
