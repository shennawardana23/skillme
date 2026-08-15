---
name: autonomous-loops
description: Use when choosing or designing an architecture for running a coding agent repeatedly with little or no human intervention — a scripted sequence of non-interactive agent invocations, a parallel-generation loop, an iterative PR-and-merge loop, or a spec-driven multi-agent pipeline. Trigger phrases include "set up an autonomous loop", "run the agent in a pipeline", "which loop pattern should I use", "generate N variations in parallel", or "orchestrate multiple agents against an RFC". For operating and recovering an already-running loop, see continuous-agent-loop.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Autonomous Loops

A catalog of architectures for running a coding agent repeatedly with
little or no human intervention between iterations, from a simple
scripted sequence of non-interactive invocations up to a full RFC-driven
multi-agent DAG. Pick the simplest pattern that fits the task — most
autonomous work needs pattern 1 or 2, not the most sophisticated one on
this list.

## Pattern spectrum

| Pattern | Complexity | Best for |
|---|---|---|
| Sequential pipeline | Low | Scripted daily-dev steps with a known order |
| Persistent session loop | Low | Interactive, session-aware iteration |
| Parallel generation loop | Medium | Many independent variations of one spec |
| Iterative PR loop | Medium | Multi-day projects with CI as the gate |
| De-sloppify pass | Add-on | Cleanup after any implementation step |
| RFC-driven multi-agent DAG | High | Large features, many interdependent units, real merge risk |

## 1. Sequential pipeline

The simplest loop: a script that invokes the agent non-interactively,
one focused step at a time, each building on the filesystem state the
previous step left behind.

```bash
#!/bin/bash
set -e

# 1. Implement, TDD
claude -p "Read docs/auth-spec.md. Implement OAuth2 login in internal/auth/. Write tests first."

# 2. Cleanup pass (see De-sloppify below) — separate context, separate concern
claude -p "Review files changed by the last commit. Remove tests that verify language/framework behavior rather than business logic. Keep real logic tests. Run go test ./... after."

# 3. Verify
claude -p "Run go build, go vet, and go test ./... . Fix any failures. Do not add new features."

# 4. Commit
claude -p "Create a conventional commit for the staged changes."
```

Design principles:

- **Each step is isolated.** A fresh invocation per step means no context
  bleed from an earlier, now-irrelevant step.
- **Order matters and is explicit.** The script is the source of truth
  for sequencing, not the model's judgment about what to do next.
- **Avoid negative instructions in the implementer step** ("don't write
  pointless tests") — telling a model what *not* to do inside a step that
  also has to do a lot of *positive* work makes it hesitant across the
  board. Add a separate cleanup step instead (see below).
- **Exit codes propagate.** `set -e` stops the pipeline on the first
  failure rather than compounding it into the next step.

Model routing composes naturally here: route a research/analysis step to
a stronger model and a mechanical transform step to a faster one, based
on what the step actually requires — see `agentic-engineering` for how to
make that call per task.

## 2. Persistent session loop

For interactive, session-aware iteration instead of one-shot scripted
steps: keep a running conversation/session (with its own history file or
equivalent state) and drive it with successive prompts, rather than
starting fresh each time. Useful for exploration where you want
accumulated context; poor fit for CI/CD because the growing context
itself becomes a variable you don't control run to run.

| | Persistent session | Sequential pipeline |
|---|---|---|
| Interactive exploration | Good fit | Poor fit |
| Scripted automation | Poor fit | Good fit |
| Context accumulation | Grows every turn | Fresh each step |
| CI/CD integration | Hard to reproduce | Straightforward |

## 3. Parallel generation loop

A two-role pattern for generating many independent variations from one
specification: an **orchestrator** role reads the spec, scans existing
output to avoid duplication, and assigns each of N sub-agents a distinct
creative direction and a specific, non-conflicting output slot; each
**generator** sub-agent gets the full spec plus its assignment and
produces one variation, in parallel with the others.

The key design point: don't rely on sub-agents to self-differentiate.
The orchestrator assigns the differentiation (a direction, a slot) up
front — otherwise parallel agents converge on similar outputs because
each is independently reaching for the "obvious" interpretation of the
spec.

Batch size: run 1-5 agents simultaneously without much thought; batch
larger counts in waves of 3-5 so each wave's output can inform the
uniqueness constraints for the next wave, rather than launching dozens at
once with no feedback loop between waves.

## 4. Iterative PR loop

For multi-day projects where CI is the quality gate: each iteration
creates a branch, runs the agent against a prompt, commits, opens a PR,
waits for CI, and merges on green — repeating until a stop condition is
hit.

```
create branch -> agent implements -> (optional reviewer pass)
  -> commit -> push -> open PR -> wait for CI
  -> CI red? -> agent fixes, re-push, re-wait
  -> CI green -> merge -> back to main -> repeat
```

Two things make this pattern safe rather than a runaway loop:

- **A cross-iteration notes file.** Since each iteration is a fresh
  invocation, persist a short "what's done / what's next" note in the
  repo (e.g. a markdown file the agent reads at the start of each
  iteration and updates at the end) to bridge context across iterations
  that otherwise share nothing but the filesystem.
- **An explicit stop condition** — a maximum iteration count, a cost
  ceiling, a time box, or a completion signal the agent emits when it
  judges the work done. Require several consecutive completion signals
  before actually stopping, so one overconfident "I'm done" doesn't end
  the loop early.

The public `continuous-claude` project (Anand Chowdhary) is a worked
example of this pattern with CI-failure auto-recovery built in — treat it
as a reference implementation to read, not a dependency to assume is
already installed in your environment.

## 5. The de-sloppify pass

An add-on for any of the above: a dedicated cleanup step that runs after
an implementation step, in its own context, rather than constraining the
implementer with negative instructions.

Why not just tell the implementer "don't over-test"? Negative
instructions bleed into unrelated judgment calls — a model told not to
write excessive tests tends to also skip legitimate edge-case tests,
because it can no longer tell where the line is. Two focused passes (one
thorough implementer, one focused cleanup reviewer) outperform one
implementer trying to self-limit while also being thorough.

```bash
claude -p "Implement with full TDD. Be thorough with tests."
claude -p "Cleanup pass: remove tests that verify language/framework
behavior rather than business logic, redundant checks the type system
already guarantees, and dead code. Keep all business-logic tests. Run
the test suite after to confirm nothing broke."
```

## 6. RFC-driven multi-agent DAG

The most sophisticated pattern: decompose a spec/RFC into a dependency
DAG of work units, run each unit through a pipeline of separate
agent stages (research, plan, implement, test, review), and land units
through a merge queue that evicts and retries on conflict. See
`references/rfc-dag-orchestration.md` for the full walkthrough, including
worked complexity tiers and the merge-queue eviction protocol. Use this
pattern only when work units are genuinely interdependent and parallel —
for a single-file change or a quick iteration, any simpler pattern above
gets there faster with less coordination overhead.

## Choosing a pattern

```
Single focused change?              -> sequential pipeline
Multi-day project, spec exists?
  need parallel implementation?     -> RFC-driven DAG
  otherwise                         -> iterative PR loop
Need many variations of one thing?  -> parallel generation loop
None of the above?                  -> sequential pipeline + de-sloppify
```

These compose: sequential-pipeline-plus-de-sloppify is the most common
combination in practice; an iterative PR loop can add a de-sloppify
directive as its reviewer-pass prompt; any loop should sit behind the
verification discipline described in `ai-regression-testing` before a
change is considered done.

## Anti-patterns

- **No exit condition.** A loop with no max-iteration count, cost
  ceiling, time box, or completion signal will keep running past the
  point where it's doing anything useful — see Real-world grounding.
- **No context bridge between iterations.** Each invocation starts fresh;
  without a shared notes file or equivalent persisted state, later
  iterations repeat or contradict earlier decisions.
- **Retrying the identical failure.** A failed iteration should feed its
  error context into the next attempt, not just re-run the same prompt
  and hope.
- **Negative instructions instead of a cleanup pass.** See pattern 5.
- **One context window for author and reviewer.** The reviewer should
  never be the same context that wrote the code — see the author-bias
  discussion in `references/rfc-dag-orchestration.md`.
- **Ignoring file overlap between parallel agents.** Two agents editing
  the same file without a merge strategy (sequential landing, rebase, or
  conflict resolution) will silently clobber each other's work.

## Real-world grounding

Knight Capital's 2012 trading-system incident is a widely documented
case of automation running without an effective kill switch: a
deployment error left old code active alongside new code, and the
system executed automated trades for roughly 45 minutes before anyone
stopped it, producing large losses. The lesson generalizes directly to
autonomous coding loops: an automated loop with no explicit stop
condition, no cost ceiling, and no fast path to halt it is a liability
regardless of how good the individual steps are.

## Verification

- [ ] The loop has an explicit stop condition (count, cost, time, or signal)
- [ ] Context bridges across iterations via persisted state, not agent memory
- [ ] A failed iteration's error context feeds the next attempt
- [ ] Cleanup is a separate pass, not negative instructions bolted onto the implementer
- [ ] Parallel agents have a defined merge strategy for any file they might both touch
- [ ] The reviewer stage runs in a separate context from the implementer stage
