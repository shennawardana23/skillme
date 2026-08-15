---
name: continuous-learning
description: Use when setting up end-of-session extraction of reusable patterns (error resolutions, user corrections, workarounds, project conventions) from a Claude Code session into a saved skill file, via a Stop hook. Trigger phrases include "save what we learned this session", "extract a skill from this conversation", "set up a Stop hook for pattern extraction", or "turn this fix into a reusable skill". For a finer-grained, hook-driven approach with confidence scoring, see continuous-learning-v2.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Continuous Learning (session-end extraction)

A lightweight approach to turning a Claude Code session into a reusable
skill: at the end of a session, evaluate the transcript for extractable
patterns and save the useful ones as new skill files. This is the coarse,
whole-skill version of the idea — see `continuous-learning-v2` for an
atomic, confidence-scored alternative and how the two differ.

## When to use this over v2

- You want something simple to reason about: one evaluation pass, one
  output artifact (a skill file), no confidence machinery.
- You're comfortable reviewing and approving what gets extracted rather
  than having it accumulate automatically.
- You don't need project-scoped vs. global separation of what's learned.

If you want atomic, per-behavior learning with confidence scores that
strengthen or decay over time, and separate handling for project-specific
vs. universal patterns, that's what `continuous-learning-v2` is for —
don't try to bolt confidence scoring onto this simpler approach; use v2
directly instead.

## How it works

This runs as a **Stop hook** — a hook that fires once, at the end of a
session, rather than on every tool call:

1. **Session-length check** — skip sessions too short to contain a real
   pattern (a reasonable default is 10+ substantive exchanges; a
   two-message session rarely has anything worth extracting).
2. **Pattern detection** — scan the transcript for the pattern types
   below.
3. **Extraction** — write useful patterns out as skill files in a
   dedicated location (e.g. `~/.claude/skills/learned/`), so they're
   available to future sessions the same way any other skill is.

### Why a Stop hook specifically

- **Lightweight** — it runs once, not on every message, so it adds no
  per-turn latency.
- **Complete context** — by the time it fires, it has the full session
  transcript to reason over, rather than a partial view.
- The tradeoff (see v2's rationale) is that a Stop hook only sees
  sessions that end normally and only gets one shot at the whole
  transcript — there's no per-tool-call observation to catch behaviors
  mid-session.

## Pattern types to detect

| Pattern | What it captures |
|---|---|
| Error resolution | How a specific error was actually fixed |
| User corrections | A place the user redirected an incorrect approach |
| Workarounds | A quirk of a framework/library and how it was worked around |
| Debugging techniques | An approach that successfully isolated a hard bug |
| Project-specific conventions | An unwritten convention discovered this session |

Configure what to ignore as deliberately as what to detect — extracting
"skills" for one-time typo fixes or external API flakiness just adds
noise that future sessions have to load and discard:

```json
{
  "min_session_length": 10,
  "extraction_threshold": "medium",
  "auto_approve": false,
  "patterns_to_detect": [
    "error_resolution", "user_corrections", "workarounds",
    "debugging_techniques", "project_specific"
  ],
  "ignore_patterns": [
    "simple_typos", "one_time_fixes", "external_api_issues"
  ]
}
```

`auto_approve: false` is a deliberate default — review what gets promoted
to a persistent skill before it starts shaping future sessions. An
over-eager pattern extracted from one unusual session can otherwise steer
every later session in the wrong direction.

## Hook setup

Wire the extraction as a `Stop` hook in your Claude Code settings:

```json
{
  "hooks": {
    "Stop": [{
      "matcher": "*",
      "hooks": [{
        "type": "command",
        "command": "path/to/your/evaluate-session-script"
      }]
    }]
  }
}
```

The script's job is exactly the three steps above: check session length,
detect patterns per your configured types, write out any that pass the
threshold as skill files.

## Gotchas

- A skill extracted from a single unusual session and auto-approved
  without review can quietly misdirect every future session that loads
  it — keep `auto_approve` off until you trust the extraction quality for
  your own sessions.
- Extracting a "workaround" for what was actually a one-off external API
  outage produces a stale skill that tells future sessions to work around
  a problem that no longer exists.
- A Stop hook only fires on sessions that end normally — a session that's
  killed or crashes produces no extraction, even if it contained a
  genuinely valuable pattern.

## Real-world grounding

The underlying practice — capturing a hard-won fix as a durable artifact
so the next person (or the next session) doesn't have to rediscover it —
is the same motivation behind postmortem "lessons learned" documents and
runbooks in incident response: the value isn't in having the incident,
it's in making sure the resolution outlives the person who found it.

## Verification

- [ ] The Stop hook is wired and fires at session end
- [ ] `min_session_length` filters out sessions too short to contain a real pattern
- [ ] `ignore_patterns` excludes one-off fixes and external-service flukes
- [ ] `auto_approve` stays off until extraction quality has been reviewed
- [ ] Extracted skills are periodically reviewed and stale ones removed
