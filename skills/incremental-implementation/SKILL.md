---
name: incremental-implementation
description: Builds in thin vertical slices -- implement one piece, test it, verify it, commit, then expand -- instead of writing a whole feature in one pass. Use when implementing any change touching more than one file, when a task feels too large to start, or before writing more than roughly 100 lines without running anything.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Incremental Implementation

Build in thin vertical slices: implement one piece, test it, verify it,
commit, then expand. Each increment leaves the system in a working,
testable state. This is the execution discipline that turns "implement the
whole feature" from one large, unverifiable leap into a sequence of small,
individually-checked steps.

## When to use

Any multi-file change; building a feature from a task breakdown;
refactoring existing code; any time about to write more than ~100 lines
before running anything.

**Skip it for:** a single-file, single-function change where the scope is
already minimal.

## The increment cycle

```
Implement ──→ Test ──→ Verify ──→ Commit ──→ next slice
```

1. **Implement** the smallest complete piece of functionality.
2. **Test** — run the suite, or write a test if none covers this yet.
3. **Verify** the slice works (tests pass, build succeeds, or a manual
   check for something not yet automatable).
4. **Commit** with a descriptive message covering this slice only.
5. **Move to the next slice** — carry forward, don't restart from scratch.

## Slicing strategies

**Vertical (preferred)** — one complete path through the whole stack per
slice:

```
Slice 1: create a task (DB + API + minimal UI) → user can create via the UI
Slice 2: list tasks (query + API + UI)         → user can see their tasks
Slice 3: edit a task (update + API + UI)       → user can modify tasks
Slice 4: delete a task (+ confirmation)        → full CRUD complete
```

Each slice delivers working, end-to-end functionality — never "all of the
database layer, then all of the API, then all of the UI," which leaves
nothing demonstrably working until the very last slice lands.

**Contract-first** — when backend and frontend need to proceed in
parallel: define the API contract (types/interfaces/OpenAPI) first, then
implement backend against it and frontend against mock data matching it,
then integrate.

**Risk-first** — tackle the riskiest or most uncertain piece first (e.g.
prove a WebSocket connection works) so a fundamental blocker surfaces
before slices 2 and 3 are built on top of an unproven foundation.

## Rules

**Simplicity first.** Before writing code, ask what the simplest thing
that could work is; after writing it, ask whether the abstraction is
earning its complexity or whether three similar lines would be clearer
than a generic framework built for one caller. Implement the obviously
correct, naive version first — optimize only once correctness is proven by
a passing test.

**Scope discipline.** Touch only what the task requires. Don't "clean up"
adjacent code, refactor unrelated imports, or add features that "seem
useful" while in the area. If something outside scope is worth fixing,
note it rather than fixing it:

```
NOTICED BUT NOT TOUCHING: src/utils/format.ts has an unused import,
unrelated to this task. Want a separate task for it?
```

**One thing at a time.** A single increment changes one logical thing —
adding a component, refactoring an existing one, and updating build config
belong in three separate commits, not one.

**Keep it compilable.** After each increment the project builds and
existing tests pass — never leave the codebase broken between slices.

**Feature flags for incomplete work.** If a feature isn't ready for users
but needs to merge in pieces, gate it behind a flag rather than exposing a
half-built path:

```go
if featureFlags.TaskSharingEnabled {
    // new sharing behavior
}
```

**Rollback-friendly commits.** Additive changes (new files/functions) are
easy to revert cleanly; avoid deleting something and replacing it with its
successor in the same commit — split so a revert of one doesn't require
manually reintroducing the other.

## Directing an agent incrementally

Be explicit about what's in scope for *this* increment and what waits for
the next one:

```
Implement just the database schema change and the API endpoint for
Task 3. Don't touch the UI yet -- that's the next increment.
After implementing, run the test suite and the build to verify
nothing else broke.
```

## Gotchas

- Re-running a build or test command that just passed, with no
  intervening edit, adds no information — it's a stalling tell, not
  extra rigor. Re-run after the next change, not as reassurance.
- A "small" refactor bundled into a feature commit makes both harder to
  review and to bisect later — if a bug appears, the commit can't
  distinguish whether the refactor or the feature caused it.
- Deleting old code and adding its replacement in the same commit means a
  revert has to manually resurrect the deleted half; split it into a
  removal commit and an addition commit when the change is risky enough
  to want that safety net.
- A feature flag added after the incomplete feature is already merged and
  reachable is a flag added one commit too late — gate it from the first
  commit that introduces the incomplete path.

## Definition of done (the bar every increment clears)

Beyond the per-increment checklist below, a task isn't done until: the
full test suite passes (not just tests for this slice), the build is
clean, the feature works end-to-end as specified, and no uncommitted
changes remain. Per-increment verification is the local check; this is the
standing bar the whole task is measured against regardless of how many
increments it took.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "I'll test it all at the end" | Bugs compound — an error in slice 1 makes every slice built on top of it wrong too. Test each slice as it lands. |
| "It's faster to do it all at once" | Feels faster until something breaks and there's no way to tell which of 500 changed lines caused it. |
| "These changes are too small to commit separately" | Small commits are free; large commits hide bugs and make reverts painful. |
| "This refactor is small enough to include" | Mixed refactor-plus-feature commits are harder to review and to bisect. Separate them. |

## Real-world grounding

Git bisect's usefulness is a direct function of commit granularity: a
project with one commit per logical change lets `git bisect` narrow a
regression to the exact change in a handful of steps; a project with large
mixed commits turns the same search into "which of these twelve unrelated
changes broke it," often forcing a manual re-read of the whole diff
instead of a binary search. Incremental, single-purpose commits are what
make that tool (and code review) actually work as intended.

## Verification

- [ ] Each increment was individually tested and committed
- [ ] The full test suite passes, not just tests for the current slice
- [ ] The build is clean
- [ ] No increment mixed an unrelated refactor with the feature change
- [ ] No uncommitted changes remain once the task is complete
