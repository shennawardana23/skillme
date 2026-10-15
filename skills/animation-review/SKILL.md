---
name: animation-review
description: This skill should be used when the user asks to "review this animation", "check this transition/motion code", "grade this animation diff", or wants a craft-bar pass over CSS/JS motion (transitions, keyframes, Framer Motion, GSAP, WAAPI). Reviews a diff or snippet against a fixed motion-craft bar and returns a findings table plus a block/approve verdict. Does not review non-motion code — see code-reviewer for that — and does not search for missing animation or plan a repo-wide motion audit — see animation-opportunity-scan and animation-improvement-planning for those.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "web-quality"
---

# Animation Review

A scanner, not a workflow: walk every rule below against the diff or
snippet under review, record every violation, then emit the verdict. This
is the motion-craft counterpart to `code-reviewer` — same fixed-category
scan shape, but its category is animation feel, which none of this
catalog's other review skills cover.

Default posture: flag first, approve second. A transition that runs
without erroring but feels sluggish, starts from nowhere, fires on every
keystroke, or drops frames is a regression, not a pass.

## The rules

Every animated property change in the diff is checked against all of
these. A violation is a finding.

1. **Justified.** Every animation answers "why does this move?" — feedback,
   spatial continuity, state change, or bridging a jarring cut. "Looks
   nice" on something seen often is a finding, not a pass.
2. **Frequency-scaled.** Keyboard shortcuts and anything triggered 100+
   times/day get no animation. Tens-per-day gets subtle/fast only.
   Occasional (modals, toasts) gets standard motion. Rare/first-run
   (onboarding, success states) is where a longer, bouncier moment is
   allowed.
3. **Correct easing.** Entrances and exits use `ease-out` or an authored
   cubic-bezier. `ease-in` on a UI element is a finding — it delays the
   part of the motion the user is watching. Default browser easings read
   as weak; expect a real curve.
4. **Duration budget.** UI animation stays under 300ms. A `Modal`/`Drawer`
   at 200–500ms is fine; anything slower on an interactive element needs
   a stated reason in the diff or a finding.
5. **Physical origin.** A popover, dropdown, or tooltip scales from its
   trigger (`transform-origin` set to the anchor), not from center or
   `scale(0)`. Start from `scale(0.92–0.97)` plus opacity. Centered modals
   are the one exception.
6. **Interruptible.** Anything the user can re-trigger quickly (toggles,
   toasts, drags) needs a transition or spring that retargets from its
   current value, not a `@keyframes` animation that restarts from zero.
7. **GPU properties only.** Animate `transform` and `opacity`. Animating
   `width`, `height`, `margin`, `padding`, `top`, or `left` — or a Framer
   Motion `x`/`y`/`scale` shorthand under load — is a performance finding.
8. **Reduced motion + hover gating.** `prefers-reduced-motion` must soften
   (not zero out — keep opacity/color, drop movement). Hover-triggered
   motion needs `@media (hover: hover) and (pointer: fine)` so touch
   devices don't inherit a hover-only effect.
9. **Asymmetric weight.** A deliberate action (hold-to-confirm, a
   destructive confirm) animates slower than the system's response to it.
   Equal in/out timing on a press-and-release is a finding.
10. **Fits the product.** Motion should match the surface's personality —
    bouncier is fine on a playful surface, wrong on a dense dashboard. If
    unsure whether a given motion belongs, the strongest fix is often to
    delete it rather than tune it.

## Escalation triggers — flag on sight

- `transition: all`
- `scale(0)` or a fade-only entrance with no initial transform
- `ease-in` on a UI interaction
- Motion on a keyboard shortcut or a 100+/day action
- A UI duration over 300ms with no comment explaining why
- `transform-origin: center` on a trigger-anchored element
- `@keyframes` on a toast, toggle, or anything re-triggerable
- Layout-property animation (`width`/`height`/`margin`/`padding`/`top`/`left`)
- A CSS custom property updated on a parent purely to drive a child's
  transform (forces a style recalc on every frame)
- No `prefers-reduced-motion` handling on a movement-heavy addition
- `:hover` motion with no pointer/hover media gate

## Fix order (cheapest leverage first)

When a finding needs a fix, prefer the earliest option that resolves it:

1. Delete the animation (high-frequency, no stated purpose, or
   keyboard-triggered).
2. Shrink it — shorter duration, smaller transform delta, fewer animated
   properties.
3. Fix the easing (`ease-in` → `ease-out` or an authored curve).
4. Fix the origin/starting state (`transform-origin`, replace `scale(0)`).
5. Make it interruptible (keyframes → transition/spring).
6. Move it to GPU properties.
7. Split symmetric timing into asymmetric.
8. Polish (stagger a group entrance, blur-mask a crossfade).
9. Add reduced-motion and hover gating.

## Output

Two sections, always in this order.

**Findings table** — one row per issue, `file:line` cited:

| Before | After | Why |
|---|---|---|
| `transition: all 300ms` | `transition: transform 200ms ease-out` | names the animated property; `all` catches unintended layout properties |

**Verdict** — group remaining commentary by tier (omit empty tiers), then
close with one of:

- **Block** — any feel-breaking finding: wrong easing on UI, motion on a
  high-frequency/keyboard action, `scale(0)`, or a layout-property
  animation with an easy GPU-only fix.
- **Approve** — no feel-breaking finding, nothing that should be deleted,
  durations/easing in range, reduced-motion respected.

State every finding against a concrete `file:line`. Don't approximate a
duration or curve — if the diff doesn't already define a token for it,
say what value it should use and why.

## Boundaries

- Review only. Don't rewrite the animation yourself unless asked —
  point out the fix, let the diff's author or `animation-improvement-planning`
  own the change.
- Non-motion findings (a11y roles, unrelated logic bugs) belong to
  `code-reviewer` or `accessibility` — don't fold them into this report.
