---
name: code-simplification
description: Simplifies working code for clarity without changing behavior — reduces nesting, renames vague identifiers, removes dead code and unnecessary abstraction. Use when code works but is harder to read, maintain, or extend than it should be, or when a review flagged accumulated complexity. Not for adding features or fixing bugs — pair with test-generation or debug for those.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Code Simplification

Simplify code by reducing complexity while preserving exact behavior. The goal is not fewer lines — it's code a new team member would understand faster than the original. Every change must pass that test.

**When NOT to use:** the code is already clean; you don't yet understand what it does (comprehend first); it's performance-critical and the simpler version would be measurably slower; you're about to rewrite the module entirely (don't polish throwaway code).

## The five principles

1. **Preserve behavior exactly.** Same output for every input, same error behavior, same side effects and ordering, all existing tests pass unmodified. If you're not sure a change preserves behavior, don't make it.
2. **Follow project conventions.** Simplification means matching the codebase's existing style for imports, error handling, naming, and type annotation depth — not imposing outside preferences. Simplification that breaks project consistency is churn, not simplification.
3. **Prefer clarity over cleverness.** Explicit code beats compact code whenever the compact version needs a mental pause to parse — e.g. replace a dense nested ternary chain with a small function of `if` returns; replace a chained `reduce` doing inline mutation with a named `Map` built in a loop.
4. **Maintain balance — don't over-simplify.** Watch for: inlining away a helper that gave a concept a name; merging two simple functions into one complex one; removing an abstraction that exists for extensibility or testability, not complexity; optimizing for line count instead of comprehension.
5. **Scope to what changed.** Default to simplifying recently modified code. Don't drive-by refactor unrelated code unless asked — it creates diff noise and regression risk.

## Process

### Step 1 — understand before touching (Chesterton's Fence)

Before changing or removing anything, answer: what is this code's responsibility? What calls it, what does it call? What are its edge cases and error paths? Why might it have been written this way — performance, a platform constraint, history? Check `git blame` for the original context. If you can't answer these, read more before simplifying.

### Step 2 — identify simplification opportunities

**Structural:**

| Pattern | Signal | Fix |
|---|---|---|
| Deep nesting (3+ levels) | Hard to follow control flow | Guard clauses / extract to helper |
| Long functions (50+ lines) | Multiple responsibilities | Split into focused, named functions |
| Nested ternaries | Requires a mental stack to parse | if/else chain, switch, or lookup map |
| Boolean flag params (`doThing(true, false, true)`) | Call site is unreadable without checking the signature | Options object or separate named functions |
| Repeated conditionals | Same `if` check duplicated | Extract a well-named predicate function |

**Naming:**

| Pattern | Signal | Fix |
|---|---|---|
| Generic names (`data`, `temp`, `result`) | No content hint | Name the content: `userProfile`, `validationErrors` |
| Misleading names | `get` that also mutates | Rename to reflect actual behavior |
| "What" comments | `// increment counter` above `count++` | Delete — code is already clear |
| "Why" comments | `// Retry: API is flaky under load` | Keep — carries intent the code can't express |

**Redundancy:**

| Pattern | Signal | Fix |
|---|---|---|
| Duplicated logic (5+ lines repeated) | Same block in multiple places | Extract shared function |
| Dead code | Unreachable branch, unused var, commented-out block | Remove, after confirming it's truly dead |
| Unnecessary wrapper | Adds no value over the thing it wraps | Inline it |
| Over-engineered pattern | Factory-for-a-factory, strategy with one strategy | Replace with the direct approach |

### Step 3 — apply incrementally

One simplification at a time; run tests after each. If tests pass, keep going or commit; if they fail, revert and reconsider — don't push through a failure by adjusting the test. Submit refactoring separately from feature/bug-fix work — a change that does both is two changes.

**Rule of 500:** if a refactor would touch more than ~500 lines by hand, invest in automation (codemod, `sed`, AST transform) instead — manual edits at that scale are error-prone and exhausting to review.

### Step 4 — verify the result

Compare before/after: is it genuinely easier to understand? Any new pattern inconsistent with the codebase? Is the diff clean and reviewable? If the "simplified" version is harder to follow or review, revert — not every attempt succeeds.

## Language examples

```typescript
// Unnecessary async wrapper
// Before: async function getUser(id) { return await userService.findById(id); }
// After:  function getUser(id) { return userService.findById(id); }

// Manual array building -> filter
// Before: const activeUsers = []; for (const u of users) if (u.isActive) activeUsers.push(u);
// After:  const activeUsers = users.filter((u) => u.isActive);
```

```python
# Nested conditionals -> guard clauses
# Before: if data is not None: if data.is_valid(): ... else: raise ValueError(...) else: raise TypeError(...)
# After:
def process(data):
    if data is None:
        raise TypeError("Data is None")
    if not data.is_valid():
        raise ValueError("Invalid data")
    return do_work(data)
```

Prop drilling through intermediate components, or replacing it with context/composition, is a judgment call — flag it, don't auto-refactor it.

## Gotchas

- A simplification that requires editing a test to make it pass is not a simplification — it changed behavior. Treat any test edit during a "pure" simplification pass as a stop sign, not a nuisance to work around.
- Types are not documentation of intent. A well-named function communicates *why* in a way a type signature never does — "this abstraction is well-typed" is not evidence it's also simple or that it's safe to leave in place.
- Inlining a one-line helper can make the call site *longer* to read once its descriptive name is gone — count comprehension cost, not lines removed, before inlining something that exists purely to name a concept.
- "The original author must have had a reason" is sometimes true and sometimes just the residue of iteration under deadline pressure with no reason at all — `git blame` and the surrounding commit message settle this faster than guessing either way.

## Real-world grounding

This process mirrors the public [Claude Code Simplifier plugin](https://github.com/anthropics/claude-plugins-official) pattern — a dedicated simplification pass run separately from feature work, gated on behavior-preservation and test-passing at every step — generalized here to be usable by any coding agent rather than tied to one plugin implementation.

## Verification

- [ ] All existing tests pass without modification
- [ ] Build succeeds with no new warnings; linter/formatter passes
- [ ] Each simplification is a separate, reviewable, incremental change
- [ ] No unrelated changes mixed into the diff
- [ ] No error handling was removed or weakened
- [ ] No dead code left behind (unused imports, unreachable branches)
- [ ] Simplified code follows this project's conventions (checked against CLAUDE.md or equivalent)
