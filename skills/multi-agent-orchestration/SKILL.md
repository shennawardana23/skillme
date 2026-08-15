---
name: multi-agent-orchestration
description: Use when designing a system where one agent delegates to multiple subagents/workers, deciding whether a task should fan out to multiple agents at all, decomposing a task without creating a "telephone game" handoff chain, or reviewing why an orchestrator's subagents are duplicating work, leaving gaps, or declaring success too early. Distinct from the orchestrator skill (single-session skill routing) - this skill covers dispatching work to separate agent instances.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "agentic-engineering"
---

# Multi-Agent Orchestration

Fanning work out to multiple agents is a real capability with a real cost —
most orchestrator failures trace back to vague task specs and wrong
decomposition boundaries, not to the workers being individually incapable.
This skill covers when fan-out is worth its cost and how to decompose,
brief, and verify subagents so the pattern actually pays off.

## Decide whether to fan out at all

Multi-agent systems run materially more tokens than a single agent doing
the same work sequentially — treat this as a real cost to justify, not a
free capability upgrade. Reach for fan-out when:

- The task is genuinely parallelizable (independent research angles,
  independent files/modules) — not a pipeline where each stage needs the
  previous stage's full context to do its job.
- The task's scale would otherwise overflow a single agent's productive
  context, and the sub-pieces are independently verifiable.

Prefer a single agent when:

- Sub-tasks have heavy interdependency or need shared, evolving context —
  splitting these produces a handoff chain that degrades fidelity at every
  step, the same way a message degrades passed through several retellers.
- The extra cost doesn't clearly buy proportionally better output — most
  coding tasks have fewer genuinely independent parallelizable pieces than
  open-ended research does.

## Decompose by context boundary, not by pipeline stage

Splitting "one agent plans, one implements, one tests" as three sequential
handoffs is the default a task looks like it wants, and is usually the
wrong cut — each handoff loses fidelity, and the implementer often needs
context the planner had but didn't fully transfer. Split instead where a
subtask can be handed a self-contained brief and genuinely doesn't need the
rest of the system's evolving context to do its job well — group anything
that needs shared understanding into one agent rather than splitting it for
its own sake.

## Brief every subagent explicitly

A vague task description is the single most common cause of orchestrator
failure — not the subagents lacking skill. A brief that's missing an
explicit objective, output format, source/tool guidance, or scope boundary
produces either duplicated work (two subagents independently cover the same
angle) or gaps (an angle nobody covers), and looks like subagent
incompetence when it's actually a specification problem. Every subagent
brief needs: the specific objective (not the parent task's whole goal),
the expected output shape, which tools/sources to use, and an explicit
boundary against what a sibling subagent is covering.

## Scale effort deliberately

Match the number of subagents and their call budget to the task's actual
complexity, stated as an explicit rule the orchestrator follows rather than
decided ad hoc per task: a small number of calls for simple fact-finding, a
handful of subagents for a comparison-style task, many subagents with
explicitly divided responsibilities only for genuinely complex,
wide-scoped work. Without an explicit rule, orchestrators tend to both
over-spawn (many agents for a task one could handle) and under-scope (too
few for something genuinely wide).

## Verifying subagent output

A subagent asked to verify another's work will tend to run a small,
convenient check and declare success — the same shortcut a rushed human
reviewer takes. Counter this with an explicit, concrete success criterion
in the verification brief (not "check this is correct" but "run the full
suite and confirm N specific things"), and require negative-path checking
(does it correctly reject the bad case, not just accept the good one), not
just a single happy-path confirmation.

## Gotchas

- **Vague task descriptions cause duplicate or gapped work, not
  incompetent subagents** — before concluding subagents performed poorly,
  check whether the brief they were given actually specified a distinct,
  bounded objective.
- **Decomposing by pipeline stage (plan/implement/test as three separate
  agents) is the intuitive cut and usually the wrong one** — it creates a
  handoff chain where each transfer loses fidelity; decompose by what's
  genuinely independently understandable instead.
- **A verifier subagent will declare early victory by default** unless
  given an explicit, exhaustive check and required negative-path testing —
  a generic "verify this" brief gets a generic, shallow verification back.
- **Multi-agent systems cost several times the tokens of a single agent
  doing the same work** — this is a real, non-marginal cost; justify
  fan-out against it rather than treating parallelism as free.
- **Tasks with heavy inter-agent dependency or shared evolving context are
  an explicit anti-fit for fan-out**, not just a suboptimal choice — the
  coordination overhead this pattern trades for isolation is exactly what
  a heavily-interdependent task needs to avoid.
- **A single slow subagent can block an entire batch** if the orchestration
  is synchronous (no ability to steer or reassign an in-flight worker) —
  design around this constraint rather than assuming subagents can be
  redirected mid-flight the way a single agent's own reasoning can.

## Real-world grounding

Anthropic's own published account of building a multi-agent research
system documents exactly this failure pattern: a vague instruction like
"research the semiconductor shortage" caused two subagents to independently
duplicate the same angle while a third explored an unrelated period,
purely from under-specification — not from any subagent being individually
weak. The same writeup states multi-agent systems can run several times
the token cost of a single-agent approach and explicitly recommends against
reaching for the pattern on tasks needing heavy shared context or most
coding tasks, where the number of genuinely parallelizable subtasks is
usually smaller than in open-ended research.

## Verification

- [ ] Fan-out was chosen because the task is genuinely parallelizable and
      the cost is justified, not by default
- [ ] Decomposition boundaries follow independent context, not just
      pipeline stage
- [ ] Every subagent brief states its distinct objective, output shape,
      and boundary against sibling subagents
- [ ] Effort (subagent count, call budget) is scaled deliberately to task
      complexity via a stated rule
- [ ] Any verification subagent is given an explicit exhaustive check and
      a negative-path requirement, not a vague "check this"
