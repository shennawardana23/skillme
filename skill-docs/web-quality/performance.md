## What it does

Covers the broader performance surface — resource budgets, the critical
rendering path, image/font/JS delivery, caching, and runtime efficiency —
that sits above the three named Core Web Vitals metrics. The distinction
that matters: this skill is strategy and budget ("how much JS should this
page ship, how should fonts load, what's the caching policy"), while
`core-web-vitals` is diagnosis-and-fix for one specific metric once it's
already identified as the problem.

## When to reach for it

Reach for this for a broad "why is my site slow" or "how do I set a
performance budget" ask, or anything about resource loading strategy
(preconnect, prerender, code splitting, caching headers) rather than one
named metric. For Go backend profiling and measurement-first triage
(pprof, allocation counts, N+1 queries), use `performance-optimization`
instead — that skill's frontend section is intentionally the short,
measurement-first version, and this skill is its deeper frontend
reference for the actual resource-loading and delivery patterns.

## Common questions

- **"Our origin needs 400ms+ to assemble the HTML response — is there
  anything to do about TTFB besides speeding up the backend itself?"**
  Yes: HTTP 103 Early Hints. The origin can respond with a `103 Early
  Hints` status carrying `Link: rel=preload` headers for the LCP
  image/critical CSS before the final `200 OK` is ready, so the browser
  starts fetching those resources during the wait instead of after.
  Chromium-based browsers act on it; others just ignore the 103 and fall
  through to the 200 — safe to enable unconditionally.
- **"Do Speculation Rules ever make things worse?"** Yes, in two ways:
  bandwidth/CPU cost (each prerender is close to a full page load, so
  scope `where` carefully and avoid `eagerness: "immediate"` outside
  small sites), and side effects firing early — analytics and ad code
  that runs on page load will fire when the prerender starts, not when
  the user actually navigates, unless explicitly gated on
  `document.prerendering`/the `prerenderingchange` event.
- **"We import all of lodash for one function — does that actually
  matter?"** Yes for bundle size specifically: `import _ from 'lodash'`
  defeats tree-shaking and ships the whole library for one call; import
  the specific function (`import debounce from 'lodash/debounce'`)
  instead.

## It's working if

- Total page weight and per-resource-type budgets (JS, CSS, images,
  fonts, third-party) are being tracked, not just eyeballed
- The LCP image and critical fonts are preloaded, not just referenced
- JavaScript is deferred/async/split by default, not loaded eagerly and
  synchronously
- Cache-Control headers distinguish immutable hashed assets from
  content that actually changes

## Where it fits

The strategy layer above `core-web-vitals`'s per-metric diagnosis, and
the deep frontend reference `performance-optimization` defers to for its
own frontend section. Feeds into `web-quality-audit`'s Performance
category when a broader four-category review is what's actually needed.
