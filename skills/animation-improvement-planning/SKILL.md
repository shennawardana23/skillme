---
name: animation-improvement-planning
description: This skill should be used when the user asks to "improve the animations", "audit the motion across the app", "make this feel better" without a specific diff to review, or wants a prioritized roadmap of animation fixes rather than a single verdict. Surveys a codebase's motion, produces a severity-ranked findings table, then writes self-contained implementation plans an executor with zero context can follow. Read-only on source code — it never edits application code itself. Use animation-review instead for a single diff, and animation-opportunity-scan instead for finding missing motion with no existing code to fix.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "web-quality"
---

# Animation Improvement Planning

An advisor skill: spend judgment on understanding the codebase's motion
and deciding what's worth fixing, then hand off execution as a plan
precise enough for any executor — including a cheaper model or a
different agent — to run without asking a follow-up question.

Read-only. The only files this skill writes live under `plans/` (or
`animation-plans/` if `plans/` is already used for something else). If
asked to "just fix it," decline and point at running the resulting plan.

## Workflow

### Phase 1 — Recon

Map the motion surface before judging it: framework and motion library
(CSS transitions, Framer Motion/Motion, Vue `<Transition>`, GSAP, WAAPI),
where tokens live (`--ease-*`, `--duration-*`, a Tailwind config,
`@keyframes` blocks), the product's personality (playful vs. dense
dashboard — this changes what "too much" means), and a rough frequency
map (which animated surfaces are hit 100+ times/day vs. occasionally vs.
rarely). Plans must extend existing easing/duration tokens, never invent
parallel ones.

### Phase 2 — Audit

Check every animated surface against these eight categories: purpose &
frequency, easing & duration, physicality & origin, interruptibility,
performance (GPU-only properties), accessibility (reduced motion, hover
gating), cohesion with existing tokens, and missed opportunities. For a
large codebase, split the categories across parallel read-only
subagents — each one gets the recon facts, the category to check, and an
instruction to return `file:line` findings only, no fixes.

Depth scales with effort:

| Effort | Coverage | Findings |
|---|---|---|
| `quick` | High-traffic components only | ~5, high severity only |
| `standard` (default) | All interactive UI | full table |
| `deep` | Whole repo including marketing pages | full table + low-severity polish items |

### Phase 3 — Vet and prioritize

Re-read every cited line yourself before reporting it. Drop anything
that's deliberate (documented in a comment or design note), already
fixed, or duplicated across subagent reports. Present one table, ordered
by leverage (impact ÷ effort):

| # | Severity | Category | Location | Finding | Fix summary |
|---|---|---|---|---|---|

Severity: **HIGH** — feel-breaking (wrong easing on UI, motion on a
keyboard/high-frequency action, `scale(0)`, dropped frames). **MEDIUM** —
noticeably off (wrong transform-origin, non-interruptible dynamic UI,
missing reduced-motion). **LOW** — polish (stagger, blur-masked
crossfade, token consolidation).

List missed opportunities (motion that should exist but doesn't)
separately below the table — they're additive, not corrective.

Then stop and let the user pick which findings become plans. If running
non-interactively, default to the top 3–5 by leverage.

### Phase 4 — Write plans

One plan per selected finding, written to `plans/NNN-short-slug.md`
(monotonic numbering, respecting any plans already there). Each plan is
fully self-contained — the executor has no memory of this conversation:
inline the exact cubic-bezier, exact duration, exact file path and
current-code excerpt, never "use the curve discussed above." Include a
verification step describing how to feel-check the result (slow motion,
frame-by-frame, a real device for a gesture). Stamp the plan with the
current commit short hash. Update or create `plans/README.md` with
execution order, dependencies between plans, and a status column.

## Invocation shape

| Ask | Behavior |
|---|---|
| bare | Full workflow: recon → audit all categories → vet → confirm → plans |
| a category name (e.g. "performance", "accessibility") | Recon plus that category only |
| "plan: <description>" | Skip the audit, recon just enough to write one plan for the described fix |
| "reconcile" | Re-check `plans/` against current code — mark done plans done, refresh stale `file:line` references |

## Boundaries

- No installs, no builds with side effects, no commits, no formatters —
  read-only analysis and plan-writing only.
- "The motion here is already right" is a valid, complete audit result.
  Don't pad the table to look thorough.
- When feel can't be judged from code alone (a spring's bounce, a
  crossfade's timing), say so and put a feel-check step in the plan
  instead of guessing a value.
