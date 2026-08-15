---
name: frontend-patterns
description: Framework-agnostic client-side performance architecture — route/component code splitting, list virtualization mechanics, hydration cost, and asset loading strategy. Use when a page feels slow to load or interact with, when deciding whether/how to split a bundle or lazy-load a component, when a long list janks while scrolling, or when triaging SSR hydration cost. Complements vue-nuxt-frontend-patterns (reactivity/SSR fetching) and frontend-ui-engineering (component/accessibility patterns) rather than repeating them.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Frontend Performance Architecture

Client-side performance problems come from a small set of root causes:
shipping code the user doesn't need yet, rendering DOM nodes the user can't
see, redoing work the server already did, or loading assets the wrong way.
This skill is deliberately framework-agnostic — the mechanics below apply
whether the framework is Vue, React, or Svelte. Reach for
`vue-nuxt-frontend-patterns` for Vue reactivity/composables/SSR-fetching
specifics, and `frontend-ui-engineering` for component architecture and
accessibility; this skill is scoped to load-time and runtime performance.

## When to use

A page is slow to first paint or first interaction; deciding whether/how to
split a bundle; a long list stutters while scrolling; triaging why an SSR
page feels slow despite fast server response; choosing an image-loading
strategy.

## Code splitting: ship what's needed, when it's needed

Every byte of JavaScript shipped up front has to be downloaded, parsed, and
executed before the page is interactive — the cost is not just network
transfer, it's main-thread time. Splitting at natural boundaries defers that
cost until it's actually needed:

- **Route-based splitting** — each route's code loads only when the user
  navigates to it. This is the highest-leverage split because most users
  only ever visit a handful of routes per session; Nuxt and most modern
  meta-frameworks do this automatically per page.
- **Component-based splitting** — a heavy, rarely-used component (a chart
  library, a rich text editor, a map) loads only when it's about to render,
  via dynamic `import()` behind a loading boundary. The cost/benefit only
  makes sense when the component is both heavy and conditionally shown —
  splitting a 2KB component adds request overhead for no real gain.
- **Vendor splitting** — third-party dependencies that change rarely go in
  their own chunk, separate from application code that changes on every
  deploy, so returning users can keep the vendor chunk cached across
  releases instead of re-downloading it every time app code changes.

The failure mode to watch for is over-splitting: dozens of tiny chunks each
carry their own HTTP request overhead (even over HTTP/2 multiplexing,
there's per-chunk parse/compile cost) — split at meaningful boundaries
(route, heavy-and-conditional component), not reflexively per file.

## List virtualization: render only what's visible

A list of 5,000 DOM nodes costs the browser layout and paint time
proportional to node count, regardless of how many are actually visible in
the viewport — scrolling performance degrades because every scroll event
can trigger layout recalculation across all 5,000 nodes, not just the ~20
currently on screen.

Virtualization renders only the visible slice (plus a small overscan
buffer) and repositions/recycles that fixed-size pool of DOM nodes as the
user scrolls, so DOM node count stays roughly constant regardless of list
length — a list of 50,000 items costs the same as a list of 50 as far as
the DOM is concerned. `@tanstack/vue-virtual` (or `@tanstack/react-virtual`
on the React side) is the common implementation; the mechanics are the same
regardless of library:

```
virtualizer tracks: total item count, estimated/measured item size, scroll offset
on each scroll: compute which item indices intersect the viewport (+ overscan)
render: only those items, absolutely positioned at their computed offset
```

Reach for virtualization once a list is large enough that node count itself
is the bottleneck — typically triggered by scroll jank complaints or DOM
node counts in the thousands in a profiler, not as a default for every list.

## Hydration cost (SSR/SSG apps)

Server-rendered HTML appears instantly, but a client-side framework still
has to "hydrate" it — attach event listeners and reconcile its virtual
representation against the existing DOM — before it's interactive. A page
that looks done can still be unresponsive to the first click or keystroke
until hydration finishes; this gap is exactly what the Interaction to Next
Paint (INP) Core Web Vital measures.

Levers that reduce hydration cost:

- **Reduce hydrated JavaScript on the critical path** — content that's
  purely static (a footer, marketing copy) doesn't need to hydrate as an
  interactive component at all; islands/partial-hydration approaches hydrate
  only the interactive fragments of a page.
- **Defer non-critical hydration** — hydrate below-the-fold or
  rarely-interacted-with sections lazily (on visibility, on idle, on
  interaction) rather than all at once on load.
- **Match server and client render output exactly** — a hydration mismatch
  (server renders one DOM shape, client expects another) forces the
  framework to discard and re-render instead of attaching to existing
  nodes, which is strictly more expensive than a clean hydration and often
  shows as a console warning worth fixing immediately, not ignoring.

## Asset loading strategy

- **Images**: serve modern formats (WebP/AVIF) with a fallback, set explicit
  `width`/`height` (or `aspect-ratio`) so the browser reserves space before
  the image loads — this is what prevents the layout-shift component of
  Cumulative Layout Shift (CLS). Use `loading="lazy"` for below-the-fold
  images and skip it for the largest above-the-fold image, since lazy-loading
  the actual LCP candidate delays it.
- **Fonts**: `font-display: swap` (or `optional`) avoids an invisible-text
  flash blocking render while a webfont loads; `<link rel="preload">` the
  critical font file if it's not discoverable early enough from the
  stylesheet alone.
- **Preload vs. prefetch**: `preload` tells the browser "fetch this now,
  it's needed for the current page" (a critical font, the LCP image);
  `prefetch` tells it "fetch this at low priority, it'll likely be needed
  soon" (the bundle for a route the user is likely to navigate to next).
  Using `preload` for things the current page doesn't need competes with
  actually-critical requests for bandwidth.

## Gotchas

- Route-based code splitting can still ship a large *shared* chunk if
  common dependencies aren't separated into their own vendor chunk — check
  the shared-chunk size, not just per-route chunk size.
- Virtualizing a list changes its DOM structure (items are absolutely
  positioned, not naturally flowing), which breaks `Ctrl+F`/browser
  find-in-page and any test that queries by DOM order rather than by
  content — this is a real trade-off to flag, not just an implementation
  detail.
- A hydration mismatch warning in the console is not cosmetic — it means
  the framework silently threw away server-rendered work and re-rendered
  from scratch, which is often slower than not using SSR at all for that
  fragment.
- `loading="lazy"` on the image that turns out to be the page's LCP
  candidate actively hurts the LCP metric it's often assumed to help —
  lazy-loading defers exactly the image you most need to load first.
- Over-splitting into many tiny chunks trades download-size savings for
  per-request overhead; profile before assuming more splitting is strictly
  better.

## Real-world grounding

Google's Core Web Vitals — LCP, CLS, and INP (which replaced FID in March
2024) — are the published, industry-standard metrics this whole skill
targets improving; they're also a direct Google Search ranking signal, which
is why "shipping less JS up front" and "reserving image space" are treated
as load-bearing performance work rather than micro-optimization. React's
public framing of virtualization (via `react-window`/`react-virtual`,
mirrored by `@tanstack/vue-virtual` on the Vue side) as "render only the
visible slice" is the standard, well-documented technique referenced above.

## Verification checklist

- [ ] Route-level code is split; shared vendor chunk size checked, not just
      per-route size
- [ ] Any list over ~1,000 rows uses virtualization, or has been profiled
      and shown not to need it
- [ ] LCP candidate image is not `loading="lazy"`, has explicit dimensions
- [ ] No hydration-mismatch warnings in the console on SSR/SSG pages
- [ ] `preload` reserved for genuinely critical current-page resources only
