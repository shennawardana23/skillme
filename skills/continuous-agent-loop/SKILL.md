---
name: continuous-agent-loop
description: Use when a coding-agent loop is already running unattended and you need to define its quality gates, detect that it's stuck (churning without progress, retrying the same failure, cost drifting up), and recover it safely. Trigger phrases include "the loop has run 20 times with no progress", "the agent keeps failing the same way", "our autonomous run is costing more than expected", or "how do I add a quality gate to this loop". For choosing which loop architecture to build in the first place, see autonomous-loops.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Continuous Agent Loop

Operating discipline for a loop that is already running — as opposed to
`autonomous-loops`, which covers choosing the loop's architecture in the
first place. This skill is about the gates that decide whether an
iteration's output is good enough to keep, and about recognizing and
recovering from a loop that has quietly gone wrong.

## Quality gates per iteration

Every iteration of a continuous loop should pass through gates before its
output is accepted, not just "the agent said it's done":

1. **Build/compile gate** — the code compiles and the type checker (or
   equivalent) is clean.
2. **Test gate** — the existing test suite passes, plus any regression
   test written for a bug this iteration fixed (see
   `ai-regression-testing`).
3. **Acceptance-criteria gate** — the specific, concrete criteria defined
   before the iteration started are met, not a general "does this seem
   right" judgment made after the fact.

A loop with only an implicit gate ("the agent's summary looked
reasonable") has no gate at all — it will accept regressions the agent
didn't notice it introduced.

## Recognizing failure modes

A continuous loop can look "alive" — consuming compute, producing diffs,
committing — while making no real progress. Watch for:

- **Loop churn without measurable progress** — iteration count climbs but
  the acceptance-criteria gate never gets closer to passing; each
  iteration touches different code without converging.
- **Repeated retries with the same root cause** — the same test fails,
  the same build error recurs, across iterations that don't look
  identical on the surface but share a diagnosis. This means the loop is
  retrying blindly rather than incorporating what the last failure taught
  it.
- **Merge or landing stalls** — for loops that merge work (a PR loop or a
  multi-unit DAG), units repeatedly failing to land cleanly signals a
  conflict the loop isn't equipped to resolve on its own.
- **Cost drift from unbounded escalation** — cost per iteration creeping
  up, usually from repeatedly escalating to a stronger/more expensive
  model on the same failure instead of fixing the underlying
  under-specification (see the escalation guidance in
  `agentic-engineering`).

The common thread: all four are detectable mechanically (iteration count
vs. gate status, repeated error signatures, landing failure count, cost
per iteration) — you don't need to read every diff to notice the loop is
stuck, if you're tracking these signals as the loop runs.

## Recovery procedure

When a failure mode is detected, don't just let the loop keep running
hoping it self-corrects:

1. **Freeze the loop.** Stop new iterations from starting before doing
   anything else — a diagnosis performed while the loop keeps mutating
   state underneath it is unreliable.
2. **Audit what actually happened.** Read the iteration history: what
   changed, what the gates reported, what the repeated failure signature
   actually is. This is a diagnosis step, not a "just try again" step —
   see `superpowers:systematic-debugging` or an equivalent structured
   debugging approach if the root cause isn't obvious from the history
   alone.
3. **Reduce scope to the failing unit.** Don't resume the full loop;
   isolate the smallest unit of work that reproduces the stuck state and
   work that in isolation.
4. **Replay with explicit acceptance criteria.** Re-enter the loop (or
   just run one manual iteration) with the acceptance criteria stated
   explicitly and the diagnosed root cause included as context — a replay
   with no new information reproduces the same stuck state.

## Gotchas

- A loop that commits on every iteration regardless of gate status makes
  the history noisy and the eventual audit harder — gate the commit, not
  just the "did it try" step.
- "It's still running" is not evidence of progress — check the
  acceptance-criteria gate's actual status, not just whether the process
  is alive.
- Cost drift is easy to miss until a bill arrives, because each individual
  escalation looks locally reasonable ("this failed, let's try a stronger
  model"). Track cost per iteration as a first-class metric while the
  loop runs, not after.
- Resuming a frozen loop with the exact same prompt as before the freeze
  guarantees the same stuck state — the replay step must add the
  diagnosis, not just remove the freeze.

## Real-world grounding

The discipline of gating automated changes on measurable criteria rather
than "it looked fine" mirrors why CI systems adopted mandatory,
non-bypassable status checks instead of relying on the committer's own
judgment about whether their change was safe — a policy exists precisely
because self-assessment under time pressure is unreliable, whether the
committer is a person under deadline pressure or an agent iterating
toward a completion signal.

## Verification

- [ ] Every iteration is gated on build/test/acceptance-criteria, not on the agent's own summary
- [ ] Iteration count, repeated-error signatures, and cost-per-iteration are tracked as the loop runs
- [ ] A detected failure mode triggers a freeze before any further diagnosis or retry
- [ ] Recovery reduces scope to the failing unit rather than resuming the full loop blind
- [ ] The replayed attempt includes the diagnosed root cause, not just a repeat of the original prompt
