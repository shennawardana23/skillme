---
name: planning-and-task-breakdown
description: Breaks an existing spec or clear requirement into a dependency-ordered, sized task list with checkpoints and parallelization notes. Use when a task feels too large to start, when scope needs estimating, or when work can be split across parallel agents/sessions. For a single lightweight plan document for one request, use planning instead.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Planning and Task Breakdown

Decompose an already-scoped body of work into small, verifiable tasks with
explicit acceptance criteria, ordered by what depends on what. Good task
breakdown is the difference between an agent that completes work reliably
and one that produces a tangled mess mid-feature — every task should be
small enough to implement, test, and verify in one focused session.

Use `skills/planning/` instead for a single lightweight plan document for
one request. Use this skill once a spec or clear requirement already
exists and the job is turning it into an ordered list of implementable
tasks — this is the mechanism `skills/spec-driven-development/`'s Plan and
Tasks phases delegate to.

## When to use

A spec exists and needs to become implementable units; a task feels too
large or vague to start; work needs to be parallelized across multiple
agents or sessions; the implementation order isn't obvious.

**Skip it for:** a single-file change with obvious scope, or a spec that
already contains well-defined tasks.

## The process

### 1. Plan in read-only mode first

Read the spec and the relevant codebase sections, identify existing
patterns and conventions, map dependencies between components, note risks
and unknowns — without writing any code yet. The output of this phase is a
plan document, not implementation.

### 2. Map the dependency graph

```
Database schema
    │
    ├── API models/types
    │       ├── API endpoints ── Frontend API client ── UI components
    │       └── Validation logic
    └── Seed data / migrations
```

Implementation order follows this graph bottom-up — build foundations
before what depends on them, not in whatever order feels most interesting.

### 3. Slice vertically, not horizontally

**Horizontal (avoid):** "build all of the schema, then all of the API,
then all of the UI, then connect it" — nothing works end-to-end until the
very last task lands.

**Vertical (prefer):** "user can create an account" (schema + API + UI for
registration), "user can log in," "user can create a task" — each task
delivers one complete, testable path through the whole stack.

### 4. Write each task to this structure

```markdown
## Task [N]: [Short descriptive title]

**Description:** one paragraph explaining what this accomplishes.

**Acceptance criteria:**
- [ ] Specific, testable condition
- [ ] Specific, testable condition

**Verification:**
- [ ] Tests pass: `<exact command>`
- [ ] Build succeeds: `<exact command>`
- [ ] Manual check: <what to verify, if anything isn't automatable>

**Dependencies:** [task numbers, or "None"]

**Files likely touched:**
- `path/to/file`

**Estimated scope:** [size — see table below]
```

### 5. Order tasks and add checkpoints

Order so dependencies are satisfied, each task leaves the system working,
and high-risk tasks come early enough to fail fast. Insert an explicit
checkpoint every 2-3 tasks:

```markdown
## Checkpoint: After Tasks 1-3
- [ ] All tests pass
- [ ] Application builds without errors
- [ ] Core user flow works end-to-end
- [ ] Review with human before proceeding
```

## Task sizing

| Size | Files | Example |
|---|---|---|
| XS | 1 | Add a validation rule |
| S | 1-2 | Add a new API endpoint |
| M | 3-5 | User registration flow |
| L | 5-8 | Search with filtering and pagination |
| XL | 8+ | **Too large — break it down further** |

Break a task down further when: it would take more than one focused
session (roughly 2+ hours of agent work); acceptance criteria can't be
described in 3 or fewer bullets; it spans two or more independent
subsystems (auth and billing); or the task title contains "and" — a
near-certain sign it's two tasks wearing one name.

## Parallelization

- **Safe to parallelize:** independent feature slices, tests for
  already-implemented features, documentation.
- **Must be sequential:** database migrations, shared state changes,
  anything in a dependency chain.
- **Needs coordination first:** features sharing an API contract — define
  the contract as its own task before parallelizing the two sides against
  it.

## Output files

Save the implementation plan to `tasks/plan.md` and the task checklist to
`tasks/todo.md`, creating `tasks/` if it doesn't exist yet — a stable
convention that downstream build/execute steps can rely on without asking
where the plan lives.

## Gotchas

- A task whose title contains "and" (e.g. "add validation and refactor the
  handler") is reliably two tasks — splitting it isn't optional cleanup,
  it's what makes the size estimate and acceptance criteria meaningful.
- Ordering tasks by "what seems most important" instead of by the
  dependency graph produces a plan where task 2 silently assumes something
  task 5 hasn't built yet — the dependency graph, not perceived priority,
  determines the order.
- A checkpoint that only says "review with human" without also listing
  concrete pass/fail conditions (tests pass, build clean) gives the human
  nothing to check against — write the same acceptance-criteria discipline
  into checkpoints that tasks get.
- An XL task that "feels fine because I understand it" is still XL for
  execution purposes — understanding the whole doesn't make each verified
  increment any smaller.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "I'll figure it out as I go" | That's how a plan becomes rework — a few minutes of dependency-mapping upfront saves hours of backtracking. |
| "The tasks are obvious" | Write them down anyway — an explicit list surfaces hidden dependencies and forgotten edge cases that "obvious" mental models skip. |
| "I can hold it all in my head" | Context windows and sessions are finite; a written task list survives a session boundary or compaction, a mental model doesn't. |
| "This L task is fine, I don't need to split it" | An L task that fails partway through leaves no smaller verified checkpoint to roll back to — splitting it is what keeps failure cheap. |

## Real-world grounding

Agile's INVEST criteria for a well-formed user story (Bill Wake, 2003) —
Independent, Negotiable, Valuable, Estimable, Small, Testable — describes
the same target this skill's task template aims for: a unit of work small
and self-contained enough to size honestly and verify independently.
"Small" and "Estimable" are exactly the two properties an XL task fails,
which is why the sizing table above forces a further split rather than
letting an XL estimate stand.

## Verification

- [ ] Every task has explicit, testable acceptance criteria
- [ ] Every task has a verification step with an exact command, not just
      "test it"
- [ ] Task order follows the dependency graph, not intuition
- [ ] No task is sized XL — anything that large has been split
- [ ] Checkpoints exist every 2-3 tasks with concrete pass/fail conditions
- [ ] The plan and task list are saved as files, not left only in chat
