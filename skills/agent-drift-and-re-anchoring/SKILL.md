---
name: agent-drift-and-re-anchoring
description: Use when a long-running agent session has drifted from its original goal, when reviewing why an agent kept working on something it wasn't asked to do, when designing checkpoint or re-anchoring behavior for an autonomous/long-horizon agent loop, or when asked about "agent focus", "agent losing the plot", or preventing an agent from going down a tangent.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "agentic-engineering"
---

# Agent Drift and Re-Anchoring

A long-running agent doesn't fail by crashing — it fails by quietly stopping
being the agent that started the task. This skill names the distinct ways
that happens ("drift" is the general term; there is no single fix) and the
specific, mechanism-matched mitigation for each, because a bigger context
window delays when drift becomes visible but does not prevent any of it.

## The six drift mechanisms

Treat these as distinct failure modes needing distinct fixes, not one
generic "agent got confused" bucket:

1. **Goal drift** — the agent starts pursuing a plausible-sounding side
   quest instead of the original ask (asked to fix a flaky test, ends up
   refactoring the whole test harness). *Mitigation*: periodic re-anchoring
   — explicitly restate the original objective and check the current
   action against it, especially before a costly or hard-to-reverse step.
2. **Context drift** — reasoning quality on unchanged information degrades
   as the context window fills; content near the middle/start becomes
   functionally less influential even though it's still technically
   present. *Mitigation*: active compression/pruning of stale detail, not
   just relying on a larger window to "fit everything."
3. **Stale context** — the agent reasons over an old snapshot of something
   that has since changed (a file edited by another process, a decision
   the user has since revised) without re-checking it's still current.
   *Mitigation*: externalize state that can change into a durable,
   re-readable store, and re-read it at decision points rather than
   trusting a cached mental model of it.
4. **Role drift** — a system prompt or initial framing loses influence
   over a long conversation as more recent content dominates.
   *Mitigation*: periodically re-inject the system prompt (or a condensed
   form of it) near the end of a long context, not just at the start.
5. **Tool-use drift** — the agent over-relies on one tool whose output is
   immediately visible and self-reinforcing (e.g., calling search
   repeatedly) while under-using a tool whose output is easy to lose track
   of once buried in history. *Mitigation*: explicit per-tool call budgets,
   and summarizing a tool's output before continuing rather than letting
   raw output pile up unprocessed.
6. **Plan decay** — the agent keeps executing steps of a plan that an
   earlier step has already made obsolete, because the plan was treated as
   a fixed script rather than a living hypothesis. *Mitigation*: treat the
   plan as mutable state re-evaluated at checkpoints, not as an execution
   log to be marched through regardless of what's been learned since.

A seventh, related failure — **hallucination cascades** — compounds any of
the above: an earlier guess gets cited later as an "established fact"
without a marker distinguishing verified-tool-output from inference.
*Mitigation*: tag context entries by provenance (system instruction /
verified tool output / model inference) so a later step can tell the
difference before building on one.

## Re-anchoring in practice

At natural checkpoints (before an irreversible action, after a significant
subtask completes, when context usage crosses a threshold):

1. Restate the original objective in one sentence, from the actual
   original request — not from a paraphrase built up over the session.
2. Compare the current action against that restatement. If it doesn't
   serve the original objective, either justify the deviation explicitly
   (a genuinely necessary detour) or stop and return to the objective.
3. Re-verify any state assumed-but-not-recently-checked (a file's current
   contents, a task's current status) rather than trusting an
   earlier-session mental snapshot of it.

## Gotchas

- **A bigger context window is not a fix for any of the six mechanisms** —
  it delays the point at which drift becomes visible, which can make the
  eventual failure larger and harder to trace back to its origin, not
  smaller.
- **Goal drift often looks like initiative, not failure**, in the moment —
  refactoring the test harness while fixing a flaky test can look like
  good engineering judgment; the tell is that it wasn't what was asked,
  and re-anchoring against the original request is what catches it, not
  a vague sense that "something feels off."
- **Tool-use drift is self-reinforcing**: a tool whose output is
  immediately visible in-context gets called again more easily than one
  whose output has to be recalled from earlier in a long session — this
  produces lopsided tool usage that has nothing to do with which tool is
  actually more useful for the task.
- **Treating the plan as an execution log rather than mutable state** is
  the single most common cause of plan decay — a step that made sense
  when the plan was written can be actively wrong by the time execution
  reaches it, if something learned in between hasn't updated the plan.
- **Provenance-tagging is cheap and rarely done** — distinguishing "the
  user said this," "a tool verified this," and "I inferred this earlier"
  costs almost nothing to track and is the single fix that prevents an
  early guess from being treated as settled fact several steps later.

## Real-world grounding

Chroma's research into long-context reasoning found measurable accuracy
degradation on unchanged tasks purely as a function of input length — a
documented effect sometimes summarized as "context rot," distinct from any
model capability limit, and the direct evidence behind why context drift
needs active management rather than "just use a bigger window." Anthropic's
own published guidance on agent design similarly treats persistent state
and periodic self-checks as first-class design concerns for any agent loop
expected to run for many steps, not an afterthought bolted onto a simple
prompt-response loop.

## Verification

- [ ] The agent's current action can be justified against the original
      request in one sentence, not just against the most recent subtask
- [ ] State that could have changed since it was last read is re-verified
      before being relied on, not trusted from an earlier snapshot
- [ ] A long session's plan has been revisited and updated at least once,
      not executed unchanged from its first draft
- [ ] Tool usage roughly matches actual task need, not just which tool's
      output happens to be easiest to see in context
- [ ] Claims built on an earlier inference are distinguishable from claims
      built on verified tool output
