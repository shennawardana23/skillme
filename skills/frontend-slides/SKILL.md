---
name: frontend-slides
description: Create zero-dependency, animation-rich HTML presentations that run entirely in a browser, from a topic/draft or by converting a PowerPoint file. Use when the user wants a talk deck, pitch deck, workshop deck, or wants a .ppt/.pptx converted to a web presentation. Helps users without a clear aesthetic preference discover one through visual previews rather than an abstract questionnaire.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Frontend Slides

Self-contained HTML presentations: one file, inline CSS/JS, no build step, no
dependencies. Every slide must fit one viewport with zero internal
scrolling — that constraint drives nearly every other decision in this
skill.

**Before generating any slide HTML**, read
`references/style-presets.md` for the mandatory viewport-fitting CSS base,
the density limits table, the preset catalog, and the CSS gotchas — those
details are not repeated in full here.

## When to use

Building a talk deck, pitch deck, workshop deck, or internal presentation;
converting `.ppt`/`.pptx` into HTML; improving an existing HTML deck's
layout, motion, or typography; exploring presentation styles with a user who
doesn't know their preference yet.

## Non-negotiables

1. **Zero dependencies** — one self-contained HTML file, inline CSS/JS,
   unless the user explicitly wants a multi-file project.
2. **Viewport fit is mandatory** — every slide fits one viewport, no
   internal scrolling, ever.
3. **Show, don't tell** — generate visual previews instead of asking
   abstract style questions.
4. **Distinctive design** — no generic purple-gradient, Inter-on-white,
   template-looking deck.
5. **Production quality** — commented, accessible, responsive, performant.

## Workflow

### 1. Detect mode

New presentation (topic/notes/draft) · PPT/PPTX conversion · enhancement of
an existing HTML deck.

### 2. Discover content

Ask only what's needed: purpose (pitch/teaching/talk/internal update),
length (short 5-10, medium 10-20, long 20+), content state (finished copy,
rough notes, topic only). If the user has content, get it before styling.

### 3. Discover style through preview, not questionnaire

If the user already names a preset from `references/style-presets.md`, use
it directly — skip previews.

Otherwise: ask what feeling the deck should create (impressed, energized,
focused, inspired), then generate **3 single-slide preview files** in
`.ecc-design/slide-previews/`, each self-contained and under ~100 lines of
slide content, showing typography/color/motion clearly. Ask which to keep
or what to mix. Use the mood-to-preset mapping in the reference file.

### 4. Build the presentation

Output `presentation.html` or `[name].html`. Add an `assets/` folder only if
the deck has extracted or user-supplied images.

Required structure: semantic slide sections; the viewport-safe CSS base from
`references/style-presets.md`; CSS custom properties for theme values; a
presentation controller class handling keyboard, wheel, and touch
navigation; Intersection Observer for reveal animations; `prefers-reduced-motion`
support.

### 5. Enforce viewport fit (hard gate)

Every `.slide` uses `height: 100vh; height: 100dvh; overflow: hidden`. All
type and spacing scale with `clamp()`. When content doesn't fit, split into
more slides — never shrink text below readable size or allow an internal
scrollbar to "solve" overflow.

### 6. Validate

Check the deck at 1920×1080, 1280×720, 768×1024, 375×667, and 667×375. If
browser automation is available, verify no slide overflows and that
keyboard navigation works (see `browser-testing-with-devtools`).

### 7. Deliver

Delete temporary preview files unless the user wants them kept. Open the
deck with the platform-appropriate command (`open` on macOS, `xdg-open` on
Linux, `start ""` on Windows). Summarize file path, preset used, slide
count, and the easy theme-customization points (the CSS custom properties).

## PPT/PPTX conversion

Prefer `python3` with `python-pptx` to extract text, images, and speaker
notes — cross-platform, unlike macOS-only extraction tools. If
`python-pptx` is unavailable, ask whether to install it or fall back to a
manual/export-based workflow. Preserve slide order and notes, then run the
same style-discovery workflow as a new presentation.

## Implementation requirements

**HTML/CSS**: inline unless the user wants multi-file. Fonts from Google
Fonts or Fontshare. Prefer atmospheric backgrounds, strong type hierarchy,
abstract shapes/gradients/grids/noise over illustrations.

**JavaScript**: keyboard navigation, touch/swipe navigation, mouse-wheel
navigation, a progress indicator or slide index, reveal-on-enter animation
triggers.

**Accessibility**: semantic structure (`main`, `section`, `nav`), readable
contrast, full keyboard-only navigation, and `prefers-reduced-motion`
support — this is a real accessibility feature (WCAG 2.3.3, "Animation from
Interactions"), not an optional nicety.

## Anti-patterns

Generic startup gradients with no visual identity; system-font decks unless
intentionally editorial; long bullet walls; code blocks that need
scrolling; fixed-height content boxes that break on short screens; negated
CSS functions like `-clamp(...)` (see the reference file's CSS gotcha —
browsers silently ignore these).

## Gotchas

- `right: -clamp(28px, 3.5vw, 44px)` is invalid CSS that browsers silently
  ignore rather than error on — a slide can look subtly wrong with no
  console warning. Always write `calc(-1 * clamp(...))` instead.
- `100vh` doesn't account for mobile browser chrome (address bar) resizing
  the viewport — pair it with `100dvh` so short-viewport mobile browsers
  don't clip the bottom of a slide.
- A preview file under ~100 lines is a design-exploration tool, not the
  final deck — don't let a preview's simplified structure leak into the
  full build's required controller/observer/reduced-motion scaffolding.
- Fitting content by shrinking font size below a readable minimum is not
  "viewport fit" — it's an accessibility failure disguised as a fix. Split
  the slide instead.
- If the user already named a specific preset, generating three unrelated
  previews anyway wastes a round trip — check for a named preset before
  defaulting to the preview workflow.

## Real-world grounding

The `prefers-reduced-motion` media feature this skill requires on every
generated deck is a W3C-standardized user preference (part of the Media
Queries Level 5 spec, implemented across all major browsers) that
represents a real accessibility need — vestibular disorders where motion
effects can trigger nausea or dizziness — not a cosmetic toggle.

## Related skills

- `frontend-ui-engineering` for component-level accessibility and
  interaction patterns feeding into deck UI chrome.
- `browser-testing-with-devtools` for automated viewport/overflow
  verification during validation.

## Deliverable checklist

- [ ] Runs from a local HTML file in a browser, zero dependencies
- [ ] Every slide fits its viewport with no scrolling
- [ ] Style is distinctive and intentional, not a generic template
- [ ] Animation is meaningful, not noisy, and respects reduced motion
- [ ] Validated across the required viewport sizes
- [ ] File path, preset, slide count, and customization points explained
