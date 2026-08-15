---
name: debug
description: Runs a tight four-phase loop — reproduce, isolate, diagnose, fix — to root-cause a single bug with evidence, not guesses. Use when given a specific error message, stack trace, or unexpected behavior for one bug and need a fast, mechanical session (reproduce → isolate → diagnose → fix → prevention). For broader failure-class triage across tests/builds/incidents, or a stop-the-line policy, use debugging-and-error-recovery instead.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Debug

## Iron law

**Do not guess. Do not apply a fix until you understand the root cause.** Every fix must follow from evidence. "It might be X" is not sufficient — "evidence Y proves it is X" is.

## Phase 1: reproduce

Establish a minimal reproduction before doing anything else:
1. What exact input or action triggers the issue?
2. What is the expected outcome? What is the actual outcome (exact error message, full stack trace)?
3. Is it deterministic or intermittent?
4. When did it start? What changed?

If you cannot reproduce it, you cannot fix it with confidence.

## Phase 2: isolate

Narrow down the location:
1. **Binary search** — which layer does the failure originate in?
2. **Remove variables** — disable features, mock dependencies one at a time.
3. **Check inputs** — is the problem in the data, not the code?
4. **Check environment** — does it fail in prod but not dev, and why specifically?

For Go: check whether `go test -race` reveals it; add `log.Printf("[DEBUG] state=%v", state)` at layer boundaries; use `dlv` (Delve) for interactive stepping when a print-based binary search stalls.

## Phase 3: diagnose

Identify the root cause with evidence: trace the execution path from input to failure; check every assumption (types, nil-ness, ordering, timing); for panics, read the stack trace top-to-bottom — the root cause is usually the deepest frame, not the one at the top where the panic surfaced.

## Phase 4: fix

1. Fix the root cause, not a symptom.
2. Write a test that would have caught this bug — it should fail without the fix and pass with it.
3. Check whether the same bug pattern exists in similar code elsewhere in the codebase.
4. Update documentation if the behavior was wrong but previously undocumented as correct.

## Common Go patterns

| Symptom | Likely cause | Check |
|---|---|---|
| nil pointer dereference | Unguarded pointer use | Add a nil check before dereference |
| goroutine leak | Missing stop signal | Add a `context.Context` or `done` channel |
| data race | Concurrent map/slice mutation | Add a mutex, or use `sync.Map` |
| `context deadline exceeded` | Missing timeout propagation | Pass `ctx` into every blocking call on the path |
| unexpected log output | Wrong logger instance | Check which logger got wired at construction |
| test flakiness | Time-dependent assertion | Use an eventually/poll pattern instead of a fixed sleep |

## Output format

```
## Root Cause
<One sentence stating the exact cause>

## Evidence
1. <Observation that proves the cause>
2. ...

## Fix
<Code change with explanation>

## Prevention
<Test or guard that prevents recurrence>
```

## Gotchas

- A stack trace's *top* frame is where the failure surfaced, not necessarily where it originated — for a nil pointer dereference three calls deep, the panic frame tells you where the nil was used, not where it should have been set; walk down to the frame where the value was supposed to be populated.
- "It's intermittent" is often a state or ordering bug wearing a random-looking disguise — before concluding something is truly non-deterministic, check for shared package-level state, map iteration order (Go deliberately randomizes it), or test execution order dependence, which produce results that look random but are fully deterministic given the hidden variable.
- A fix that makes the reported symptom disappear without a new failing-then-passing test is unverified, even if it looks obviously correct — a regression test written *after* confirming the fix, not derived from staring at the diff, is the only way to know the reproduction case is actually gone.
- Adding a nil-check or a `recover()` around a panic without asking why the value was nil in the first place converts a loud, visible bug into a silent one — that is a symptom fix, not a root-cause fix, even though Phase 4 was nominally followed.

## Real-world grounding

The `context deadline exceeded` row above is one of the most common real-world Go production symptoms precisely because a missing `ctx` parameter on one function several calls deep compiles fine and passes unit tests that don't exercise real network latency — the bug only appears under load, which is why Phase 2's "check environment: prod vs dev" step exists as a named phase rather than an afterthought.

## Verification

- [ ] The bug reproduces reliably before any fix is attempted
- [ ] The root cause is stated as a single sentence backed by cited evidence, not a hypothesis
- [ ] A regression test exists that fails without the fix and passes with it
- [ ] Similar code elsewhere was checked for the same bug pattern
- [ ] The fix addresses the cause, not just the symptom that was reported
