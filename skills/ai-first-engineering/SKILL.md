---
name: ai-first-engineering
description: Use when designing team-level process, architecture standards, review policy, or hiring/evaluation criteria for an engineering org where AI agents generate a large share of the implementation. Trigger phrases include "how should our review process change for AI-generated code", "what architecture makes this codebase agent-friendly", "what should we test for in AI-first hiring", or "set testing standards for agent-written code". Not for single-task workflow — see agentic-engineering for that.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# AI-First Engineering

The organizational operating model for a team where AI agents write a
large share of the code: what changes about architecture, review, testing
policy, and hiring when the median contributor to a diff is an agent
directed by an engineer. This is a process/org-design skill — for how to
run *one task* through an agent, see `agentic-engineering`.

## What actually shifts

1. **Planning quality outweighs typing speed.** When implementation is
   cheap, the bottleneck moves to whether the plan was correct — a
   precisely wrong plan gets implemented precisely and wrongly, fast.
2. **Eval coverage outweighs anecdotal confidence.** "I tried it and it
   worked" doesn't scale when the person who tried it and the system that
   wrote it share the same blind spots (see `ai-regression-testing`).
3. **Review focus shifts from syntax to system behavior.** Syntax and
   style are cheap for a model to get right; behavioral correctness under
   edge cases and load is not — review effort should follow the risk, not
   old habits about what code review used to catch.

## Architecture that agents can work in safely

Prefer architectures with properties a model can reason about locally,
without holding the whole system in context:

- **Explicit boundaries** — a service or package's contract is written
  down, not inferred from reading five other files.
- **Stable, typed contracts** — function signatures and API shapes that
  don't vary based on hidden runtime state.
- **Deterministic tests** — a test that's flaky independent of the code
  under test teaches an agent (and a human) the wrong lesson about
  whether its change broke something.

Avoid implicit behavior spread across unwritten conventions ("we always
do X in this codebase, it's just not documented anywhere") — a human
engineer picks this up by osmosis over months; an agent starting a fresh
session has no such history, and will violate the convention with full
confidence.

## Code review policy for AI-first teams

Review for, in priority order:

1. Behavior regressions — did existing functionality change unintentionally
2. Security assumptions — auth, input trust boundaries, secrets handling
3. Data integrity — especially any change touching partitioned or shared
   tables, where a missing filter silently scans (or misses) the wrong
   partition
4. Failure handling — what happens when a dependency errors or times out
5. Rollout safety — can this be deployed and rolled back independently

Minimize time spent re-litigating style choices that automated
formatting/linting already enforce — that's a symptom of review habits
built for a world where humans typed every line and style was a decent
proxy for care. It isn't anymore.

## Testing standard for generated code

Raise the bar relative to what felt sufficient for hand-written code,
specifically because the author and the reviewer role are structurally
correlated when both are the same model:

- Required regression coverage for any domain the agent touched, not just
  the new behavior it added.
- Explicit edge-case assertions (empty input, boundary values, concurrent
  access) rather than only the happy path the feature request described.
- Integration checks at interface boundaries — the seams between modules
  are exactly where "each file works, together they don't" failures hide.

## Hiring and evaluation signals

Strong AI-first engineers are distinguished less by typing speed and more
by:

- decomposing ambiguous work into verifiable units cleanly (see
  `agentic-engineering`'s 15-minute-unit framing)
- writing measurable acceptance criteria before work starts, not
  discovering them during review
- producing high-signal prompts and evals — the prompt is now a work
  artifact worth reviewing, not a throwaway
- holding risk controls under delivery pressure — not skipping the eval
  step when the deadline is close, which is exactly when skipping it is
  most costly

## Gotchas

- Treating "the agent's code passed CI" as equivalent to "a human wrote
  and reviewed this" understates the risk — CI catches what it was built
  to catch, and an agent's blind spots are systematically different from
  a human's (see `ai-regression-testing` for the specific failure
  pattern).
- Porting an old review checklist unchanged just relocates human
  attention to the same cheap-to-verify issues (formatting, naming) that
  were already the easiest part to automate, while the actually-shifted
  risk (behavioral edge cases, security assumptions) goes unreviewed.
- A codebase full of implicit, undocumented conventions doesn't fail
  loudly when handed to an agent — it fails quietly, with confidently
  wrong code that matches the *written* contract but violates the
  unwritten one.

## Real-world grounding

The shift described here — automated tooling absorbing the mechanical
parts of code quality (formatting, lint) so human review time
concentrates on behavior and security — is the same trajectory static
analysis and CI took across the industry over the last two decades; teams
that never updated their review checklists to match kept reviewing for
the parts CI already covered, and under-reviewed everything CI couldn't
catch. AI-first teams face the same trap at a faster clock speed.

## Verification

- [ ] Architecture decisions favor explicit boundaries and typed contracts over implicit convention
- [ ] Review policy names its priority order (behavior/security/data/failure/rollout) explicitly
- [ ] Testing standard requires regression coverage for touched domains, not just new features
- [ ] Hiring/evaluation criteria reward decomposition and acceptance-criteria writing, not just speed
- [ ] No policy assumes "CI passed" is equivalent to independent human review
