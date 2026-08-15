---
name: planning
description: Converts a single feature request, bug report, or vague task into one short, ordered plan document with a goal, acceptance criteria, sized steps, and risks. Use for a single task or feature that needs structure before work starts. For decomposing an already-specced body of work into a dependency-ordered task list, use planning-and-task-breakdown instead.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Planning

Turn one request into a short, executable plan — a single document, not a
multi-phase process. Use this when a request needs structure before work
starts but doesn't warrant the full gated spec workflow: a feature request
that's vague, a bug report that needs a clear fix boundary, or an
architectural choice that needs its trade-offs written down before code
gets touched.

Use `skills/spec-driven-development/` instead when the work needs a
human-reviewed gate at each phase (new project, ambiguous multi-week
scope). Use `skills/planning-and-task-breakdown/` instead when a spec
already exists and the job is decomposing it into an ordered, sized task
list — this skill produces one plan document for one request, not a task
breakdown for an existing spec.

## When to use

A feature request needs to be broken into implementation steps; a task is
vague and needs structure before work begins; an architectural decision
needs its options and trade-offs written down.

## Plan structure

Every plan contains five parts.

**1. Goal (one sentence).** What will be verifiably true when this is
done — not "improve X," but a statement you could check against the
finished result.

**2. Context.** What exists today that this changes or depends on, what
constraints apply (time, compatibility, existing conventions), who the
stakeholders are — two to four sentences, not a design document.

**3. Acceptance criteria.** A numbered list of observable outcomes. Each
item is either "Given / When / Then" or a binary, checkable statement —
not "works correctly."

**4. Steps (ordered).** Each step has a clear deliverable, can be verified
independently of the others, carries a size estimate (S = hours, M = about
a day, L = 2+ days), and names its blockers or dependencies on other
steps.

**5. Risks and open questions.** What could go wrong and how it would be
detected; what's unknown and needs investigation before it can be sized;
what decisions are being deliberately deferred.

```markdown
## Goal
<one sentence>

## Context
<2-4 sentences>

## Acceptance Criteria
1. ...
2. ...

## Steps
1. [S] Step description — Deliverable: X
2. [M] Step description — Depends on: step 1
3. [L] Step description — Blocker: need API key from ops

## Risks
- Risk: X → Detection: Y → Mitigation: Z

## Open Questions
- ...
```

## Vague vs. concrete

| Vague | Concrete |
|---|---|
| "Add authentication" | "Add JWT validation middleware that returns 401 on a missing/invalid token" |
| "Fix the bug" | "Fix nil pointer in `user.GetProfile()` when the user has no profile row" |
| "Improve performance" | "Reduce p99 latency on `/api/search` from 2s to <500ms by adding a cache" |

A vague goal produces vague acceptance criteria, which produces a plan
nobody can verify against — sharpen the goal first; the rest of the
document follows from it.

## Gotchas

- A step with no size estimate is a step nobody has actually thought
  through — if you can't say S/M/L, you don't yet understand the step
  well enough to sequence it correctly relative to the others.
- "Risks" that list only things you already know you'll handle isn't a
  risk section, it's a restatement of the steps — a real risk section
  names something that might invalidate a step or its estimate.
- An acceptance criterion phrased as "the feature works" is unfalsifiable
  — nobody can point to a state of the system and say whether that
  criterion passed or failed.
- Skipping the Context section on the theory that "everyone already knows
  this" is exactly how a plan becomes wrong the moment it's picked up by
  someone (or some future session) without that shared context.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "The steps are obvious, no need to write them down" | Obvious to the person who just thought about it — not to whoever executes it after a context reset. |
| "I'll size the steps once I start" | Sizing forces you to actually think through the step; skipping it just defers finding out a step is really an L in disguise. |
| "One sentence for the goal is too little" | If the goal needs a paragraph, it's actually several goals — split the plan or accept that the paragraph belongs in Context, not Goal. |

## Real-world grounding

The "Given/When/Then" acceptance-criteria format recommended here comes
from Behavior-Driven Development (Dan North, mid-2000s), adopted
specifically because it forces a criterion into a testable shape at the
moment it's written, rather than leaving verification to be figured out
later — the same reasoning applies whether or not the project uses a BDD
test framework to execute it literally.

## Verification

- [ ] The goal is one verifiable sentence, not a vague direction
- [ ] Every acceptance criterion is a checkable, binary statement
- [ ] Every step has a size estimate and states its dependencies/blockers
- [ ] At least one real risk is named with a detection method, not just a
      restated step
- [ ] Open questions are listed rather than silently resolved
