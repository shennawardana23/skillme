---
name: continuous-learning-v2
description: Use when designing a fine-grained, hook-driven learning system that observes every tool call (not just session end), stores atomic "instincts" with confidence scores, and evolves clusters of related instincts into skills, commands, or agents — with project-scoped vs. global separation so React conventions don't leak into a Python repo. Trigger phrases include "observe every tool call for patterns", "score confidence on a learned behavior", "keep this pattern scoped to this project", or "cluster these instincts into a skill". For the simpler session-end, whole-skill approach, see continuous-learning.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Continuous Learning v2 (instinct-based)

A finer-grained alternative to `continuous-learning`: instead of
evaluating a whole session once at the end and writing out full skill
files, this approach observes every tool call via hooks and accumulates
small, atomic, confidence-scored "instincts" that later get clustered
into skills, commands, or agents. See `continuous-learning` for when the
simpler whole-session approach is enough — don't reach for this added
complexity unless you actually need per-behavior confidence tracking or
project/global separation.

## Why hooks instead of a Stop hook

A `Stop` hook only fires once, on sessions that end normally, and only
gets one pass over the whole transcript. `PreToolUse`/`PostToolUse` hooks
fire on every single tool call, deterministically — so no candidate
pattern is missed because the session ended unusually, and observation
happens incrementally rather than in one large end-of-session pass.

## The instinct model

An instinct is one atomic, learned behavior — one trigger, one action —
not a full skill:

```yaml
---
id: prefer-explicit-error-returns
trigger: "when writing a function that can fail"
confidence: 0.7
domain: "code-style"
source: "session-observation"
scope: project
project_id: "a1b2c3d4e5f6"
---

# Prefer Explicit Error Returns

## Action
Return an explicit error value rather than panicking, for any function
whose failure is an expected, recoverable condition.

## Evidence
- User corrected a panic-based approach to an error-return approach twice
- Pattern held across three subsequent sessions without correction
```

Properties that make this different from a whole extracted skill:

- **Atomic** — one trigger, one action, easy to evaluate independently
- **Confidence-weighted** — a number, not a boolean, that can move up or
  down as evidence accumulates
- **Domain-tagged** — `code-style`, `testing`, `git`, `debugging`, etc.
- **Evidence-backed** — records what observations produced it, so a human
  reviewing it later can see why it exists
- **Scope-aware** — `project` by default, promotable to `global`

## Confidence scoring

| Score | Meaning | Behavior |
|---|---|---|
| 0.3 | Tentative | Suggested, not enforced |
| 0.5 | Moderate | Applied when relevant |
| 0.7 | Strong | Applied by default |
| 0.9 | Near-certain | Treated as a core behavior |

Confidence rises when the pattern is repeatedly observed, when the user
doesn't correct the suggested behavior, or when independent evidence
agrees with it. Confidence falls when the user explicitly corrects the
behavior, when contradicting evidence appears, or when the pattern goes
unobserved for a long stretch — an instinct that decays on disuse avoids
permanently enforcing a preference the project has since moved past.

## Project scoping

Detect the current project (a git remote URL, or a repo path as a
fallback when no remote exists) and store instincts under that project's
identity rather than in one global pool. This prevents a React-specific
convention learned in one repo from silently applying inside an unrelated
Python service, which is the failure mode a single global instinct pool
has no way to prevent.

| Pattern type | Scope | Example |
|---|---|---|
| Language/framework convention | project | "Use hooks, not classes, in this repo" |
| File-structure preference | project | "Tests live under `internal/.../_test.go`" |
| Security practice | global | "Always validate external input" |
| General engineering practice | global | "Write the failing test first" |
| Tool workflow preference | global | "Read before Edit" |

**Promotion, project to global:** when the same instinct shows up
independently in two or more projects with consistently high confidence,
promote it to global scope — that's the signal it's a general practice
rather than something specific to one codebase's conventions. Don't
promote on a single project's evidence alone; that's exactly the
cross-project contamination this scoping exists to prevent.

## Evolution: instincts to skills/commands/agents

Individually, instincts are too granular to read through one at a time.
Periodically cluster related instincts (same domain, related triggers)
into a higher-level artifact:

- A cluster of related `code-style` instincts becomes a skill.
- A cluster describing a recurring multi-step workflow becomes a command.
- A cluster describing a consistent specialized role becomes an agent
  definition.

This evolution step is where the atomic, hard-to-review instinct layer
turns into something a human can actually read, approve, and maintain —
treat instincts as the working memory and the evolved artifact as the
thing that actually ships.

See `references/instinct-file-structure.md` for a worked example of the
on-disk layout (project-scoped vs. global directories, observation logs,
evolved-artifact directories) if you're implementing this system rather
than just applying the model conceptually.

## Gotchas

- Confidence scoring only has integrity if corrections actually decay
  it — a system that only ever raises confidence, never lowers it on
  contradiction, converges to enforcing whatever was observed first
  regardless of whether it's still correct.
- Storing everything globally by default (no project scoping) is the
  single biggest source of cross-project contamination — a convention
  correct for one stack actively wrong for another.
- Auto-promoting an instinct to global scope from evidence in a single
  project defeats the purpose of scoping in the first place; require
  independent evidence from multiple projects.
- Raw observations (every tool call, every prompt) are a much larger and
  more sensitive artifact than the instincts derived from them — keep
  observation logs local and only export/share the derived instincts,
  never the raw log.

## Real-world grounding

The instinct-and-confidence-score design generalizes a familiar idea from
recommendation and reputation systems: a small number of unweighted
observations shouldn't have the same authority as a pattern reinforced
across many independent instances, and a signal that stops recurring
should lose influence over time rather than being locked in permanently
from its first observation. Applying that discipline to learned
developer-preference behavior is what separates this from just writing
down the first pattern the model noticed and treating it as gospel.

## Verification

- [ ] Observation happens via PreToolUse/PostToolUse hooks, not a single end-of-session pass
- [ ] Every instinct has a confidence score that can rise and fall with evidence
- [ ] Instincts are scoped to a project by default; promotion to global requires multi-project evidence
- [ ] Related instincts are periodically clustered into a reviewable skill/command/agent
- [ ] Only derived instincts are ever exported or shared — never raw observation logs
