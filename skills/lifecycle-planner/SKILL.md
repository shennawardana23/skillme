---
name: lifecycle-planner
description: Execute a full development lifecycle for a feature or bug fix in one pass — requirements, technical design, implementation, self-review, and a gated approval — rather than jumping straight to code. Use for feature requests, bug fixes, or technical changes that need multiple coordinated steps, especially in Go services.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Lifecycle Planner

Run a feature or fix through the full lifecycle in one execution, ending in a real quality gate rather than "looks done." This is the only skill in this cluster that ends in working code: `idea-refine` narrows down which idea to build, `interview-me` extracts what's actually wanted; this skill takes an already-scoped ask and executes it end to end.

## Phase 1: Requirements Analysis

Transform the request into structured acceptance criteria before touching code:

```
Feature: <name>
User story: As a <persona>, I want <capability> so that <benefit>.

Acceptance criteria:
1. Given <context>, when <action>, then <outcome>
2. ...

Edge cases:
- What happens when input is empty?
- What happens on network or dependency failure?
- What are the concurrency implications?

Out of scope:
- Explicit list of what NOT to build
```

## Phase 2: Technical Design

Before writing code, decide: the interface (public function signatures and types), the data flow (how data moves through the system), the error model (what errors can occur and how they're surfaced), the dependencies (what existing code this touches), and the test strategy (what will be tested and how).

## Phase 3: Implementation

This org defaults to Go. Follow:

- Clean layering: domain → service → repository → handler
- Dependency injection via constructor parameters, not globals
- `context.Context` as the first parameter on every blocking call
- Errors wrapped with `fmt.Errorf("operation: %w", err)` so callers can `errors.Is`/`errors.As`
- No global mutable state
- Every exported symbol documented with a doc comment

Implementation checklist:
- [ ] Interface defined before implementation
- [ ] Error handling on every path, not just the happy path
- [ ] Context cancellation respected in long-running or blocking calls
- [ ] No hardcoded configuration that should be injected
- [ ] Nil/zero values handled gracefully, not assumed away

## Phase 4: Self-Review Loop

After implementation, review before calling it done:

1. **Correctness**: does it satisfy every acceptance criterion from Phase 1?
2. **Security**: any new injection surface, auth bypass, or secret exposure?
3. **Performance**: any new O(n²) pattern or unnecessary allocation introduced?
4. **Tests**: is the new behavior tested, including error paths, not just the happy path?

Fix anything the self-review finds before declaring completion — a self-review that finds nothing to fix on the first pass is worth a second, more skeptical look.

## Phase 5: Output

```json
{
  "plan": "...",
  "implementation": "...",
  "quality_status": "approved | needs_revision: <issues>",
  "test_coverage": ["list of what is tested"]
}
```

## Quality Gate

The implementation is "approved" only when: all acceptance criteria are met, there are no CRITICAL or HIGH security issues, error paths are handled (not just the happy path), and there is at least one test per new function. A gate that can be talked around isn't a gate — if any of these fail, the status is `needs_revision`, not `approved with caveats`.

## Gotchas

- "Approved" is a binary status in Phase 5, not a spectrum — resist the pull to mark something approved with a caveat listed in the same breath; a caveat means `needs_revision`.
- Edge cases identified in Phase 1 (empty input, network failure, concurrency) are easy to list and easy to silently drop by Phase 3 — the self-review in Phase 4 should re-check the Phase 1 edge-case list explicitly, not just "does it compile and pass the happy-path test."
- Wrapping every error with `fmt.Errorf("operation: %w", err)` is only useful if callers actually use `errors.Is`/`errors.As` downstream — wrapping without ever unwrapping just adds string noise to logs.

## Real-world grounding

The acceptance-criteria-before-code and gated-approval structure here mirrors Scrum's well-documented "Definition of Ready" and "Definition of Done" practices: a Definition of Ready requires a story to have clear acceptance criteria before work starts (this skill's Phase 1), and a Definition of Done requires an explicit, checkable bar — tests passing, no known critical defects — before work is called complete (this skill's Phase 5 quality gate), rather than leaving "done" to individual judgment call by call.

## Verification

- [ ] Acceptance criteria and edge cases were written down before implementation started
- [ ] The technical design (interface, data flow, error model) was decided before coding
- [ ] Implementation follows the Go layering and error-wrapping conventions
- [ ] The self-review loop ran and any issues it found were fixed, not just noted
- [ ] The final quality_status is `approved` only if every gate condition is actually met
