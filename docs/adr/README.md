# Architecture decision records

Short records of the decisions that shape this repo, in the format from Michael Nygard's
"Documenting Architecture Decisions": context, decision, alternatives, consequences.

Rules:

- One decision per file, numbered `NNNN-short-title.md`, never renumbered.
- Status moves forward only: Proposed, then Accepted or Rejected, later Deprecated or
  "Superseded by ADR-NNNN". To change a decision, write a new ADR that supersedes the old one;
  do not edit an accepted record's Context or Decision.
- Write the ADR in the same PR as the decision.

ADRs 0001 to 0005 were recorded retroactively on 2026-10-15 from the code, commit history and docs.
Each names its evidence, and its date is when the decision first appears in history.

| ADR | Title | Status | Date |
| --- | --- | --- | --- |
| [0001](0001-agent-skills-format.md) | Use the open Agent Skills format as the only content format | Accepted | 2026-08-15 |
| [0002](0002-own-deterministic-eval-runner.md) | Build our own eval runner with deterministic grading, no LLM judge | Accepted | 2026-08-15 |
| [0003](0003-distribution-plugin-and-skills-cli.md) | Distribute as a Claude Code plugin and through the `npx skills` CLI | Accepted | 2026-08-15 |
| [0004](0004-isolate-each-eval-case.md) | Run every eval case in an isolated workspace with one skill | Accepted | 2026-08-15 |
| [0005](0005-host-slides-workspace-in-repo.md) | Host the open-slide deck workspace in this repo | Accepted | 2026-10-15 |

Template for the next one:

```markdown
# NNNN. Title stated as the decision

## Status
Proposed

## Date
YYYY-MM-DD

## Context
The forces that make a decision necessary, stated neutrally.

## Decision
We will ...

## Alternatives considered
- Option: pros, cons, why rejected.

## Consequences
What gets easier and what gets harder, including the downsides.
```
