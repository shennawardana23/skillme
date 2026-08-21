## What it does

Applies WCAG 2.2's four principles — Perceivable, Operable,
Understandable, Robust (POUR) — to make content usable by people relying
on assistive technology, keyboards, or other non-mouse input. The
constraint underneath most of the specific rules: prefer a native
interactive element (`<button>`, `<a href>`, form controls) over adding
ARIA roles and manual keyboard handling to a `<div>` — native elements
get focus, keyboard activation, and correct assistive-tech semantics for
free, and hand-rolled equivalents are easy to get subtly wrong (e.g.
double-firing a click on both a native button's built-in Enter/Space
handling and a manually added `keydown` listener).

## When to reach for it

Reach for this for anything about keyboard navigation, screen readers,
color contrast, WCAG conformance, or "is this accessible." AA is the
practical target (legal requirement in many jurisdictions); AAA is
enhanced, not the default bar to aim for.

## Common questions

- **"We removed the default focus outline for the design — is that
  okay?"** Not by itself. `*:focus { outline: none; }` with nothing
  replacing it removes the only visual signal a keyboard user has for
  where they are on the page. Use `:focus-visible` to restyle rather
  than remove the outline (it only applies for keyboard-driven focus,
  not mouse clicks, so it doesn't affect the visual design mouse users
  see).
- **"What's actually new in WCAG 2.2 versus 2.1?"** Several criteria this
  skill covers are new in 2.2 specifically: Focus Not Obscured (2.4.11 —
  a focused element can't be entirely hidden by a sticky header/footer),
  Target Size (2.5.8 — interactive targets need at least 24×24 CSS
  pixels at AA), Dragging Movements (2.5.7 — any drag interaction needs a
  single-pointer alternative), Consistent Help (3.2.6), Redundant Entry
  (3.3.7), and Accessible Authentication (3.3.8 — login flows can't rely
  solely on a cognitive test like remembering a password unless an
  alternative or autofill path exists).
- **"Does `aria-hidden="true"` on a decorative icon fully solve the
  accessible-name problem for an icon button?"** No — it only stops the
  icon itself from being announced. The button still needs its own
  accessible name via `aria-label` or visually-hidden text; hiding the
  icon and skipping the label leaves the button with no name at all.

## It's working if

- Every interactive control is reachable and operable by keyboard alone,
  with no traps
- Icon-only controls have an accessible name, not just a hidden decorative icon
- Normal text meets 4.5:1 contrast (3:1 for large text), and focus
  indicators meet 3:1 against their background
- Forms announce errors to assistive tech (`aria-live`/`role="alert"`)
  and associate them with the specific invalid field

## Where it fits

One of the four categories under `web-quality-audit`, and one of the
more heavily source-grounded skills in the web-quality family — its
`references/WCAG.md` and `references/A11Y-PATTERNS.md` carry the full
criteria list and code patterns (modal focus trap, skip link, ARIA tabs,
live regions) this page only summarizes.
