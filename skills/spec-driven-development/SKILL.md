---
name: spec-driven-development
description: Writes a structured, human-reviewed specification before any code, deriving implementation from that written contract rather than from a vague request. Use when starting a new project or feature with no spec yet, when requirements are ambiguous or only exist as an idea, or before an architectural decision on a change that would take more than about 30 minutes to implement.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Spec-Driven Development

Write the spec before the code. The spec is the shared source of truth
between agent and human — what's being built, why, and how "done" will be
recognized — and every downstream implementation decision derives from it,
not from re-interpreting the original request each time. Code without a
spec is guessing dressed up as progress.

## When to use

Starting a new project or feature; requirements are ambiguous or
incomplete; the change touches multiple files or modules; about to make an
architectural decision; the task would take more than ~30 minutes to
implement.

**Skip it for:** single-line fixes, typo corrections, or a change where
requirements are unambiguous and fully self-contained.

## The gated workflow

Four phases. Don't advance until the current one is validated by a human.

```
SPECIFY ──→ PLAN ──→ TASKS ──→ IMPLEMENT
   │          │        │          │
   ▼          ▼        ▼          ▼
 human      human    human      human
 reviews    reviews  reviews    reviews
```

### 1. Specify

Ask clarifying questions until requirements are concrete, and surface
assumptions before writing anything else:

```
ASSUMPTIONS:
1. Web application, not native mobile
2. Session-based auth (not JWT)
3. PostgreSQL (matches the existing schema)
→ Correct these now, or I'll proceed with them.
```

Don't silently fill in ambiguous requirements — the spec's whole purpose
is surfacing misunderstanding *before* code exists, and a silent
assumption is the most dangerous kind of misunderstanding because nobody
gets a chance to catch it.

Cover six areas:

1. **Objective** — what's being built, for whom, what success looks like.
2. **Commands** — full executable commands, not just tool names
   (`npm run build`, not "build it").
3. **Project structure** — where source, tests, and docs live.
4. **Code style** — one real snippet in the target style beats three
   paragraphs describing it.
5. **Testing strategy** — framework, test locations, which levels cover
   which concerns.
6. **Boundaries** — three tiers: **Always** (run tests before commit,
   validate inputs), **Ask first** (schema changes, new dependencies, CI
   config), **Never** (commit secrets, edit vendor directories, delete
   failing tests without approval).

```markdown
# Spec: [Feature Name]

## Objective
## Tech Stack
## Commands
## Project Structure
## Code Style
## Testing Strategy
## Boundaries
- Always: ...
- Ask first: ...
- Never: ...

## Success Criteria
## Open Questions
```

Reframe vague requirements as testable success criteria before treating
them as settled:

```
REQUIREMENT: "Make the dashboard faster"
REFRAMED: LCP < 2.5s on 4G; initial data load < 500ms; CLS < 0.1
→ Are these the right targets?
```

### 2. Plan

From the validated spec, produce a technical plan: major components and
their dependencies, build order, risks and mitigations, what can be
parallelized vs. must be sequential, and verification checkpoints between
phases. See `skills/planning-and-task-breakdown/` for the dependency-graph
and vertical-slicing mechanics behind this step — this phase is the
"derive a plan from an already-written spec" application of that skill.
The plan should be reviewable on its own: a human should be able to read
it and say "yes, that's the right approach" without re-deriving it.

### 3. Tasks

Break the plan into discrete tasks, each completable in a single focused
session, each with explicit acceptance criteria and a verification step,
ordered by dependency rather than perceived importance. See
`skills/planning-and-task-breakdown/` for task sizing and dependency
ordering — this phase applies that skill's task template to the tasks
implied by the spec.

### 4. Implement

Execute tasks one at a time, applying
`skills/incremental-implementation/` for the vertical-slice execution
discipline and `skills/test-driven-development/` (or the
framework-specific `skills/laravel-tdd/` / `skills/tdd-workflow/`) for the
test-first loop within each task.

## Keeping the spec alive

The spec is a living document, not a one-time artifact: update it when a
decision changes (discover the data model needs to differ — update the
spec, then implement, not the reverse), update it when scope changes,
commit it alongside the code, and reference the relevant spec section from
each PR that implements part of it.

## Gotchas

- A spec that lists requirements but never states success criteria isn't
  a gate — "add search" without a testable definition of "found" lets
  every implementation claim to satisfy it.
- Silently resolving an ambiguity instead of listing it as an assumption
  defeats the entire purpose of the SPECIFY phase — the human never gets
  the chance to say "no, that's wrong" until the wrong thing is built.
- A spec written and then never updated as scope shifts becomes actively
  misleading — worse for a future reader than no spec, because it looks
  authoritative while describing a system that no longer exists.
- "This is simple, skip the spec" and "I'll write the spec after coding it
  to document what I built" are the same failure with different timing —
  both skip the part where requirements get pinned down *before* the
  implementation choices that follow from them are made.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "This is simple, I don't need a spec" | Simple tasks don't need *long* specs, but they still need testable acceptance criteria — a two-line spec is fine. |
| "I'll write the spec after it works" | That's documentation, not specification — its value is forcing clarity *before* code, not describing code after. |
| "The spec will slow us down" | A short spec prevents rework; the cost of writing it upfront is smaller than the cost of debugging a wrong assumption after the fact. |
| "Requirements will change anyway" | That's why the spec is a living document, updated when they do — an outdated spec is still more informative than none. |

## Real-world grounding

The RFC process used to define internet standards (IETF) runs on a
comparable discipline, summarized in its own culture as "rough consensus
and running code": a written specification is debated and reviewed before
implementations are built against it, precisely because multiple
independent implementers building from a shared, agreed contract produce
interoperable systems, while implementers each guessing from an informal
description do not.

## Verification

- [ ] The spec covers all six core areas (objective, commands, structure,
      style, testing, boundaries)
- [ ] Assumptions were surfaced explicitly, not silently resolved
- [ ] A human reviewed and approved the spec before planning began
- [ ] Success criteria are specific and testable, not vague adjectives
- [ ] Boundaries (Always / Ask first / Never) are defined
- [ ] The spec is saved in the repository, not only in conversation
