---
name: strategic-compact
description: Use when deciding whether to manually compact conversation context now versus later — at a phase boundary (research done, plan written, milestone shipped), before switching to an unrelated task, or when responses feel slower/less coherent from context pressure. Trigger phrases include "should I compact now", "context is getting full", "switching to a different task", or "we just finished planning, what next". Covers what survives compaction and what doesn't, so the decision is informed rather than reflexive.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Strategic Compact

Manual `/compact` at a chosen point beats compaction triggered
automatically at an arbitrary token threshold, because automatic
compaction has no notion of where in the task you actually are — it can
fire mid-debugging and discard the exact state (which hypotheses were
already tried, what the last error looked like) the next step needs.
Compacting *deliberately*, at a logical boundary, keeps what matters and
discards what doesn't.

## The decision, by phase transition

| Transition | Compact? | Why |
|---|---|---|
| Research → planning | Yes | Research context is bulky; the plan is the distilled output worth keeping |
| Planning → implementation | Yes | Once the plan is written down (a todo list or a file), free the context it took to reach it |
| Implementation → testing | Maybe | Keep if tests reference details from the code just written; compact if switching focus entirely |
| Debugging → next feature | Yes | Debug traces and dead-end hypotheses pollute context for unrelated work |
| Mid-implementation | No | Losing variable names, file paths, and partial state mid-task is costly and hard to reconstruct |
| After a failed approach, before trying another | Yes | Clear the dead-end reasoning trail before starting the next attempt with a clean slate |

The pattern across all of these: compact when the *next* step doesn't
need the *reasoning process* that got you here, only its output. Don't
compact when the next step is a continuation of the same reasoning.

## What actually survives compaction

Compacting isn't all-or-nothing — some things persist regardless, which
should shape what you do *before* compacting:

| Persists | Lost |
|---|---|
| Written instructions/config files | Intermediate reasoning and analysis |
| A task list written to a file or tracked externally | Contents of files you read earlier but didn't need to re-read |
| Anything written to disk | Multi-step conversational context |
| Git history (commits, branches) | Tool-call history and counts |
| Files on disk generally | Nuanced preferences stated verbally but never written down |

The practical implication: **write down what you'd lose before you
compact**, not after. If a preference, a decision, or a plan only exists
in conversation, put it in a file or a tracked task list first — then
compacting loses nothing that mattered.

## Best practices

1. **Compact right after planning is finalized** — once the plan lives in
   a durable place (a todo list, a plan file), the context that produced
   it has done its job.
2. **Compact after debugging concludes** — clear the error-resolution
   trail before moving to unrelated work; carrying it forward adds noise
   without adding value to a task that doesn't touch the same code.
3. **Don't compact mid-implementation** — you'll lose exactly the state
   (what file, what variable, what was already tried) the next few steps
   need.
4. **Write before compacting**, not after — if it isn't in a file or a
   task list yet, it won't survive.
5. **Give `/compact` a summary directive** when you use it, so the fresh
   context starts pointed at the right next step rather than needing to
   rediscover it (e.g. "focus on implementing the auth middleware next").

## Reducing what needs compacting in the first place

Compaction is a mitigation; the cheaper fix is loading less in the first
place:

- **Load guidance on demand, not up front.** A skill or reference file
  should load when its trigger condition is actually met, not sit loaded
  from session start "just in case" it's needed later.
- **Watch what's actually consuming context**: always-loaded project
  instructions, currently-loaded skills, growing conversation history,
  and bulky tool results (large file reads, large search results) are the
  usual suspects, in roughly that order of how easy each is to control.
- **Check for duplicated instructions** — the same rule stated in more
  than one loaded file (a global instruction file and a project-level one
  both stating the same convention) costs context twice for the same
  information.

## Gotchas

- Auto-compaction firing mid-tool-call or mid-debugging is the exact
  scenario deliberate compaction avoids — if you're waiting for the
  automatic trigger, you've given up the choice this skill is about.
- A `/compact` with no summary directive leaves the fresh context to
  reconstruct "what's next" from whatever residual signal survives,
  which is strictly worse than just telling it directly.
- Assuming a verbally-stated preference "will just be remembered" across
  a compaction is the single most common way strategic compaction goes
  wrong — if it isn't written down, treat it as already lost.

## Real-world grounding

The same "checkpoint at a coherent boundary, not at an arbitrary
interval" logic underlies why version control commits are structured
around logical units of change rather than being triggered by, say, a
fixed number of keystrokes: a commit (or a compaction) is useful in
proportion to how well it captures a coherent, resumable state — and an
arbitrary trigger has no way to know whether it's landing on a coherent
boundary or in the middle of something.

## Verification

- [ ] The compaction decision was tied to a phase boundary, not a token count
- [ ] Anything that would be lost (plans, preferences, decisions) was written down first
- [ ] `/compact` was given a directive pointing at the next step
- [ ] Guidance that isn't currently relevant was left unloaded rather than pre-loaded
- [ ] No instruction is duplicated across more than one always-loaded file
