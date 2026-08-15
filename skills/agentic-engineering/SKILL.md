---
name: agentic-engineering
description: Use when planning how to hand an implementation task to an AI coding agent — deciding how to decompose the work, which model tier to route it to, when to define eval/regression criteria before starting, and where to focus human review of the agent's output. Trigger phrases include "break this down for the agent", "which model should handle this", "define done for this task", or "review this AI-generated change".
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Agentic Engineering

The engineering practice for a single unit of work when an AI agent does
most of the typing and a human sets the goal, the constraints, and the
acceptance bar. This is a per-task discipline — how you shape and check
*one* piece of work — not the team-wide process changes (see
`ai-first-engineering` for that).

## Define done before you start

Write the completion criteria before the agent starts implementing, not
after reviewing what it produced. If you can't state a concrete pass/fail
condition, you don't yet know what "done" means well enough to hand the
task off — decompose or clarify further first.

- A capability check: what new behavior must exist and how would you
  demonstrate it (a test, a manual repro, a script).
- A regression check: what existing behavior must *not* change.
- An explicit non-goal, if the task is easy to over-scope (e.g. "do not
  touch the payment retry logic in this pass").

## Task decomposition: the 15-minute unit

Break work into units small enough that each one is independently
verifiable. A good unit:

- has a **single dominant risk** — one thing that's actually hard about
  it, not three unrelated hard things bundled together
- has a **clear done condition** you can check without re-reading the
  whole diff
- can be **verified on its own** — you shouldn't need three other
  in-flight units to know whether this one is correct

If you can't articulate the unit's dominant risk in one sentence, it's
still two units pretending to be one — split it.

## Eval-first loop

1. Define the capability eval (what should now work) and the regression
   eval (what must keep working) before implementation starts.
2. Run both against the current code to get a baseline — capture what
   already fails and why, so you're not surprised by pre-existing gaps.
3. Hand off the implementation.
4. Re-run both evals and diff against the baseline. A capability eval
   that now passes plus a regression eval with no new failures is the
   actual completion signal — not the agent's own claim that it's done.

## Model routing by task shape

Route by what kind of reasoning the task needs, not by "how important" it
feels:

| Tier | Route here | Signal |
|---|---|---|
| Small/fast | Classification, boilerplate transforms, narrow single-file edits | The change is mechanical once you know what to do |
| Mid-size | Implementation and refactors within an established pattern | The task requires judgment but the shape is known |
| Frontier | Architecture decisions, root-cause analysis, changes with multi-file invariants | Getting it wrong silently breaks something elsewhere |

Escalate tier only when the lower tier fails with a *reasoning* gap
(missed an invariant, misunderstood the goal) — not when it fails from a
missing fact you could have supplied instead (wrong file path, unstated
convention). Escalating for a missing-fact failure just burns a bigger
budget on the same mistake.

## Session strategy

- Keep working in the same session for closely-coupled units where later
  work depends on decisions made in earlier steps.
- Start a fresh session after a major phase transition (research done,
  now implementing; implementation done, now reviewing) — carrying over
  unrelated exploration context adds noise without adding value.
- Compact at a milestone boundary, not mid-debugging — compacting while
  actively chasing a failure discards the exact state (which hypotheses
  were tried, what the last error was) the next step needs.

## Reviewing AI-generated code

Spend review time where the agent is structurally more likely to be
wrong, and stop spending it where automation already covers you:

**Prioritize:**
- invariants and edge cases (does this hold for the empty case, the
  concurrent case, the partitioned-by-`hotel_id` case)
- error boundaries — what happens when a dependency fails
- security and auth assumptions — does this change who can do what
- hidden coupling — does this change something another part of the
  system silently depended on

**Don't spend review cycles on style disagreements that automated
lint/format already enforce** — reviewing for both correctness and taste
in the same pass dilutes attention on the part that actually matters.

## Cost discipline

Track per task: model tier, token estimate, retry count, wall-clock time,
and success/failure. This turns "the agent is slow/expensive" from a
feeling into a number you can act on — e.g. discovering that most cost
comes from retries on one specific task shape, which is a decomposition
problem, not a model-choice problem.

## Gotchas

- A capability eval that only checks the happy path will pass even when
  the agent silently broke an edge case nobody wrote a regression test
  for — write the regression eval from the existing test suite's gaps,
  not just from "did the new feature demo work."
- Escalating model tier after a failure without changing anything else
  about the task often reproduces the same mistake at higher cost — check
  whether the task was under-specified before assuming it was under-powered.
- A large diff with no intermediate commits is a decomposition failure
  that already happened — you can't tell which part introduced the risk
  after the fact as easily as you could have prevented it up front.

## Real-world grounding

The core idea — spend the review budget on invariants and boundaries
rather than style, because tooling already covers style — mirrors how
code review at scale has evolved: static analysis and formatters
(`gofmt`, ESLint, Prettier) took over the mechanical checks specifically
so human (and now agent) review time concentrates on correctness and
security, the things automation can't yet judge.

## Verification

- [ ] Capability and regression criteria were written before implementation started
- [ ] Each unit of work has one dominant risk and an independent done condition
- [ ] Model tier matches the task's reasoning demand, not its perceived importance
- [ ] Review focused on invariants/security/coupling, not style already covered by lint
- [ ] Cost (retries, tokens, time) was tracked per task, not just per session
