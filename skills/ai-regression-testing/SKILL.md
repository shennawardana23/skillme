---
name: ai-regression-testing
description: Use when an AI coding agent has fixed a bug or modified API/backend logic and you need to prevent the same class of bug from recurring — especially in codebases with two parallel code paths (sandbox/mock mode vs production, or feature-flagged variants). Trigger phrases include "the agent keeps missing the same bug", "write a regression test for this fix", "check sandbox and production return the same fields", or "set up a bug-check workflow". Covers Go with the standard testing package and table-driven tests.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# AI Regression Testing

Testing patterns for AI-assisted development, aimed at a specific failure
mode: when the same model writes a fix and then reviews that fix, it
carries the same blind spot into both steps. The review says "looks
correct" and the bug survives.

```
AI writes fix -> AI reviews fix -> AI says "looks correct" -> bug still there
```

This isn't a reason to distrust AI-written code generally — it's a reason
to make the check *mechanical* instead of another round of judgment from
the same source. A test either passes or it doesn't; it has no blind
spots to share with the author.

## The pattern to watch for: parallel-path drift

The single most common AI-introduced regression in codebases that have
two code paths returning "the same" data — a sandbox/mock mode next to a
real database path, or a feature-flagged variant next to the default — is
that a fix lands on one path and not the other. The model reasons about
the path it's looking at, patches it correctly, and doesn't notice the
sibling path exists unless something forces it to check.

```go
// BAD: sandbox path returns a different field set than production
func GetProfile(w http.ResponseWriter, r *http.Request) {
    if isSandboxMode(r) {
        writeJSON(w, sandboxProfile{ID: "u1", Email: "a@b.com", Name: "A"})
        return // forgot NotificationSettings here
    }
    p := loadProfileFromDB(r.Context())
    writeJSON(w, profile{ID: p.ID, Email: p.Email, Name: p.Name,
        NotificationSettings: p.NotificationSettings})
}
```

A human reviewing "did the fix work" by hitting the production endpoint
will see it work — the sandbox path's drift is invisible unless something
tests both paths against the same contract.

## Test the contract, not the implementation

Define the response contract once, then assert both paths satisfy it:

```go
package profile_test

import (
    "encoding/json"
    "net/http/httptest"
    "testing"
)

// requiredFields is the contract every code path must satisfy.
// notification_settings was added here after a bug shipped without it.
var requiredFields = []string{
    "id", "email", "name", "notification_settings",
}

func TestGetProfile_AllPaths_SatisfyContract(t *testing.T) {
    cases := []struct {
        name    string
        sandbox bool
    }{
        {"production path", false},
        {"sandbox path", true},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            rr := httptest.NewRecorder()
            req := newTestRequest(t, tc.sandbox)

            GetProfile(rr, req)

            var body map[string]any
            if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
                t.Fatalf("decode response: %v", err)
            }
            for _, field := range requiredFields {
                if _, ok := body[field]; !ok {
                    t.Errorf("%s: missing required field %q", tc.name, field)
                }
            }
        })
    }
}
```

Running both subtests against one shared `requiredFields` list is what
catches drift — a test that only exercises the production path (the
"obvious" one to test) will pass right through a sandbox-only regression.

## Regression tests as targeted memory, not coverage goals

Don't aim for blanket coverage. Write a regression test the moment a bug
is found, named after the bug it prevents, so the exact failure can never
silently reappear:

```go
// TestGetProfile_NotificationSettingsPresent guards against a regression
// where notification_settings was present in the production response but
// missing from the sandbox response after a schema change.
func TestGetProfile_NotificationSettingsPresent(t *testing.T) {
    // ... same table-driven shape as above, scoped to this one field
}
```

This strategy composes well with AI-assisted development specifically
because a model tends to make the same *category* of mistake across
similar changes — once one instance of "parallel-path drift on field X"
has a test, that category is closed for that field, and the next fix
attempt gets a mechanical check instead of another self-review pass.

## Making the check mandatory, not optional judgment

Wire tests and build/type checks as a **mandatory first step** before any
AI code review pass, so mechanically-detectable bugs never depend on the
review step noticing them:

```
Step 1 (mandatory, cannot be skipped by the reviewer's judgment):
    go test ./...
    go vet ./...
  -> failing test or vet error = highest-priority bug, stop here

Step 2 (only after step 1 passes):
    AI code review pass, focused on what tests can't check:
      - parallel-path field parity
      - error handling and rollback
      - race conditions in concurrent access
```

## Other recognizable AI regression shapes

- **Column/field selection omission** — a new field is added to a struct
  and the response payload, but not to the underlying `SELECT` (or ORM
  projection), so the field is silently always its zero value.
- **Error state leaves stale data visible** — a handler sets an error
  message on failure but doesn't clear data left over from before the
  failed call, so the UI/response shows old data next to an error.
- **Optimistic update without rollback** — client-side or cache state is
  updated before confirming the backend call succeeded, with no path to
  revert it if the call fails.

Each of these has the same test shape: assert the invariant that should
hold (all required fields present, stale state cleared on error,
state reverted on failure) rather than asserting today's exact output.

## Gotchas

- A test that only checks the path you just fixed proves the fix works
  *there* — it says nothing about the sibling path, which is exactly
  where this class of bug hides.
- Don't trust "the agent reviewed its own fix and confirmed it's correct"
  as a substitute for running the test suite — that's the same blind spot
  applied twice, not independent verification.
- Writing a test for every function regardless of history dilutes the
  signal; write tests where bugs were actually found, so test count grows
  with real risk instead of arbitrary coverage targets.

## Real-world grounding

The underlying issue — a single reasoning process being unreliable at
catching its own mistakes — is why software engineering has long
separated authorship from verification (independent code review, QA
separate from development). Self-review by an LLM reproduces the same
structural weakness a lone author reviewing their own diff has always
had: shared assumptions mean shared blind spots. Automated, deterministic
tests are the check that doesn't share the author's assumptions.

## Verification

- [ ] Every fixed bug gets a regression test named after the bug, before moving on
- [ ] Tests assert the contract (required fields, invariants) across all parallel code paths, not just the one just touched
- [ ] `go test` / `go vet` (or your stack's equivalent) run as a mandatory first step, before any AI review pass
- [ ] No completion claim is accepted on "the agent reviewed it and it looks correct" alone
