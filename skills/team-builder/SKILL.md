---
name: team-builder
description: Use when a task would benefit from several different specialist perspectives run in parallel — composing an ad-hoc team from a directory of agent persona files, dispatching them simultaneously, and synthesizing their outputs into one report. Trigger phrases include "team builder", "get me a security and performance review together", "compose a team of agents for this", or "what would a few different specialists say about this". Not for a single task routed to one skill — see orchestrator for that.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Team Builder

An interactive picker for composing and dispatching a small team of
agent personas against one task, in parallel, then synthesizing their
outputs into a single report. Distinct from `orchestrator`, which routes
one task to the single best-matching skill — team-builder is for when
you deliberately want *more than one* perspective on the same task at
once.

## When to use this over orchestrator

- The task genuinely benefits from multiple independent takes (e.g. "is
  this ready to launch" wants security, performance, and UX perspectives,
  not just one).
- You have a directory of agent persona files (markdown identity/rules/
  workflow prompts) to choose from, and want to pick a subset rather than
  running all of them every time.
- You want the outputs synthesized — agreements, tensions, and a combined
  recommendation — rather than one skill's single answer.

If the task clearly belongs to one domain, route it with `orchestrator`
instead; running a parallel team for a single-domain task just multiplies
cost without multiplying insight.

## Discovering available agents

Agent persona files are markdown documents containing an identity, rules,
workflow, and expected deliverables — not skills, and not the `Agent`
tool's built-in subagent types. Discover them from wherever your project
or global configuration keeps them (commonly a project-local directory
alongside a user-level one, with project-local definitions taking
precedence when names collide).

For each file found:

- Extract the agent's name from its first heading; if there's no
  heading, derive a name from the filename.
- Extract a one-line summary from the first paragraph after the heading.
- Group agents into domains — either from a subdirectory the file lives
  in, or (for a flat layout) from a shared filename prefix used by two or
  more files. A file with a unique prefix has no group and falls under a
  catch-all "General" bucket.

If no agent files are found, say so plainly and list where you looked —
don't fabricate a plausible-sounding list of agents that don't exist.

## Presenting the choice

Show discovered domains and let the user pick by number, by name, or by
"all agents in domain X":

```
Available agent domains:
1. Engineering (2) — Software Architect, Security Engineer
2. Frontend (1) — Vue Component Reviewer
3. Data (1) — Query Performance Analyst

Pick domains or name specific agents (e.g. "1,3" or "security + performance"):
```

Cap a team at a small number (five is a reasonable default) — beyond
that, synthesis quality drops and token cost rises faster than insight
does. If a selection exceeds the cap, list the candidates alphabetically
and ask the user to narrow down rather than silently truncating the list
yourself.

## Dispatching the team

1. Read each selected agent's full persona file.
2. Get (or confirm) the task description.
3. Spawn all selected agents **in parallel** — each one independent, with
   no expectation that they communicate with each other during the run.
   Each agent's prompt is its persona content plus the task description.
4. If one agent fails, times out, or returns nothing useful, note the
   failure inline and continue synthesizing from the agents that
   succeeded — one failure shouldn't block the whole report.

Parallel dispatch of independent personas is the right tool here; reserve
any mechanism for agents that need to *debate or respond to each other*
for a case that genuinely needs it — that's a different interaction
pattern from simply running several independent takes side by side.

## Synthesizing results

- Present each agent's findings grouped by agent.
- Add a synthesis section calling out:
  - where multiple agents agree (a converging signal worth trusting more)
  - where they conflict (a genuine tension the user needs to resolve, not
    something to paper over by picking a side silently)
  - a recommended next step given the combined picture
- If only one agent was actually selected, skip the synthesis section —
  there's nothing to synthesize from a single perspective.

## Gotchas

- Hardcoding a fixed list of "known agents" defeats the purpose — new
  persona files added to the directory should appear in the menu
  automatically without an update to this skill.
- A flat-filename layout where a domain name itself contains a hyphen
  (e.g. `product-management-*.md`) will misparse against a naive
  "split at the first hyphen" grouping rule — use a subdirectory layout
  for multi-word domain names instead of relying on prefix-splitting.
- Presenting a conflict between two agents' recommendations as if one is
  simply "more correct," without surfacing the tradeoff to the user, hides
  the actual decision the user needed to make.
- Running the full team on a task only one domain actually applies to
  wastes the cost of every agent whose findings turn out to be "not
  applicable here."

## Real-world grounding

Convening several domain specialists for one review rather than routing
to a single reviewer is the same instinct behind a cross-functional
design review or an architecture review board: a security engineer, a
performance engineer, and a UX designer looking at the same change
surface different risks precisely because none of them is trained to
notice what the others catch — the value comes from the disagreement as
much as the agreement.

## Verification

- [ ] Agents were discovered dynamically from the actual directory, not hardcoded
- [ ] The team size stayed within a sane cap; oversized selections were narrowed, not silently truncated
- [ ] All selected agents were dispatched in parallel with no artificial dependency between them
- [ ] The synthesis explicitly calls out agreement and conflict, not just a merged list of findings
- [ ] A failed agent was noted and excluded, not allowed to silently blank out the whole report
