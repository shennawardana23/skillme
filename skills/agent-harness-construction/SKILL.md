---
name: agent-harness-construction
description: Use when designing or reviewing the tool/function interface an LLM agent calls through — naming and scoping tools, shaping observation payloads, writing error-recovery contracts, or deciding between ReAct-style and structured function-calling loops. Trigger phrases include "design a tool for the agent", "the agent keeps calling the wrong tool", "what should this tool return on error", "too many tools", or "the agent doesn't recover from failures".
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Agent Harness Construction

The "harness" is everything between the model and the real system: tool
schemas, the text the model sees back from a tool call, and the rules for
when to retry versus stop. Agent quality is bounded by this interface as
much as by the underlying model — a capable model given ambiguous tools or
opaque errors will still loop, guess, or give up.

## The four levers

1. **Action space** — which tools exist, how they're named, how narrow
   their inputs are.
2. **Observation shape** — what comes back from a tool call.
3. **Recovery contract** — what the model is told to do when a call fails.
4. **Context budget** — how much of the harness's own guidance and history
   the model has to hold in its context window at once.

Treat these as one design surface, not four independent decisions — a
vague action space usually shows up downstream as a recovery problem
("the model called the wrong tool, got a cryptic error, and retried the
same wrong call").

## Action space design

- **Stable, explicit names.** `create_reservation`, not `handle_request`.
  A name should tell the model when *not* to use it as much as when to.
- **Schema-first, narrow inputs.** Prefer several small typed parameters
  over one free-text field the model has to format correctly from memory.
  Every extra way to phrase the same input is a way for the model to get
  it wrong.
- **Deterministic output shape.** The same tool call against the same
  state should return the same shape every time — varying which fields
  are present based on a hidden internal state forces the model to guess.
- **Avoid catch-all tools** (`run_command`, `execute_sql`) unless there is
  genuinely no way to scope the operation — a catch-all tool pushes all
  validation and safety judgment onto the model at call time, which is
  exactly the point where mistakes are hardest to catch.

### Granularity

Match tool grain to risk and frequency, not to what's convenient to
implement:

| Grain | Use for | Example |
|---|---|---|
| Micro | High-risk, hard-to-undo operations | `deploy_to_production`, `run_migration`, `grant_permission` |
| Medium | The common read/edit/search loop | `read_file`, `search_code`, `apply_patch` |
| Macro | Only when round-trip latency dominates cost | `run_full_test_suite` |

Micro-granularity for risky operations isn't caution for its own sake — it
creates a place to put a confirmation step, a dry-run mode, or an
irreversibility warning that a macro tool would swallow silently.

## Observation design

Every tool response the model reads back should let it answer three
questions without another round trip: *did it work, what do I do next,
and where is the evidence.* A minimal consistent shape:

```json
{
  "status": "success | warning | error",
  "summary": "one line, human-readable result",
  "next_actions": ["what a caller could reasonably do next"],
  "artifacts": ["file paths, IDs, or URLs produced"]
}
```

Returning a bare string (`"ok"`, or a stack trace dumped as text) forces
the model to infer status from prose, which it does inconsistently across
calls — this is the single most common cause of an agent misreading a
failed call as successful and moving on.

## Error recovery contract

Every failure path a tool can hit needs three things, not just an error
message:

- **Root-cause hint** — what class of thing went wrong (`auth`,
  `not_found`, `conflict`, `rate_limited`), not just the raw exception
  text.
- **Safe retry instruction** — is retrying this call ever correct, and
  under what change (e.g. "retry after re-reading the file; the
  precondition changed").
- **Explicit stop condition** — what tells the model to stop retrying and
  surface the problem instead of looping (e.g. "if this fails twice with
  the same root cause, stop and report").

An error string alone answers none of these — the model is left choosing
between blind retry and giving up, and it doesn't have a principled way
to pick.

## Context budgeting

- Keep the system prompt / tool-list minimal and stable across calls —
  volatile system-level content forces re-reasoning about things that
  didn't change.
- Move long-tail guidance into an on-demand skill or reference file rather
  than inlining it; load it only when the trigger condition is real.
- Prefer a reference *to* a file over pasting the file's content into the
  prompt when the model can just read it when needed.
- Compact or summarize at phase boundaries (research done → plan written,
  plan done → implementation starting), not at an arbitrary token count —
  an arbitrary threshold can fire mid-tool-call and drop state the next
  step depends on.

## Architecture pattern choice

- **ReAct-style (reason, then act, then observe, repeat)** — best when the
  path to the goal is genuinely uncertain and each observation should
  change the plan (open-ended debugging, exploration).
- **Structured function-calling** — best for flows with a known shape and
  typed contract (form-filling, deterministic pipelines).
- **Hybrid (usually the right default)** — ReAct-style reasoning to decide
  *what* to do next, typed function calls to actually *do* it. This keeps
  the flexibility of open-ended reasoning without letting free-text
  drift into the parts of the system that need a strict contract.

## Benchmarking a harness

Track these per task, not just "did it eventually work":

- completion rate (finished without human intervention)
- retries per task (a proxy for how often the model got confused)
- pass@1 vs pass@3 (does it need multiple attempts to succeed once)
- cost per *successful* task, not per attempt

A harness change that raises pass@3 but not pass@1 usually means the
recovery contract improved, not the action space — useful for knowing
which lever to pull next.

## Gotchas

- Adding "helpful" extra fields to every observation (verbose logs,
  internal IDs nobody asked for) grows the context cost of every single
  tool call in the session — audit what a response actually needs to
  answer "did it work / what next", and cut the rest.
- Two tools with overlapping semantics (`update_user` and `patch_user`
  that do almost the same thing) don't average out — the model
  inconsistently picks one, and your error-handling code now has to
  support both call patterns forever.
- A tool that silently succeeds on a no-op (e.g. "delete" on an
  already-deleted resource returns `status: success`) hides real failures
  behind identical-looking successful calls elsewhere in the session.

## Real-world grounding

The ReAct paper (Yao et al., 2022, "ReAct: Synergizing Reasoning and
Acting in Language Models") showed that interleaving explicit reasoning
traces with tool calls and observations measurably reduced hallucinated
actions compared to acting without an intermediate reasoning step — the
core justification for treating the observation text an agent reads back
as a first-class design object, not an afterthought bolted onto whatever
the underlying function happened to return.

## Verification

- [ ] Every tool name states what it does and implies what it doesn't
- [ ] No two tools have overlapping semantics for the same operation
- [ ] Every response has a consistent status/summary/next_actions shape
- [ ] Every error path has a root-cause hint, a retry rule, and a stop condition
- [ ] High-risk operations are their own narrow (micro-grained) tools
- [ ] System-level guidance is stable; long-tail guidance loads on demand

See `references/observation-schemas.md` for worked examples of
observation payloads across a few common tool categories.
