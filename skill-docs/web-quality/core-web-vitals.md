## What it does

Goes deep on exactly three metrics — LCP, INP, and CLS — because these
three, and only these three, are the ones Google measures at the 75th
percentile of real visits and factors into page-experience ranking. The
non-obvious constraint: "Good" isn't an average, it's a percentile bar —
75% of a page's actual visits have to clear the threshold, so a page that
looks fine in a single lab test (Lighthouse, one run) can still fail in
the field if the slow tail of real users is large enough.

## When to reach for it

Reach for this when the ask is specifically about one of LCP/INP/CLS by
name, or a Search Console "Core Web Vitals" report is flagging a metric —
this skill has the actual diagnosis-and-fix detail (preload strategies,
scheduler-yielding patterns, layout-shift causes) that `performance`
covers only at the budget/strategy level. Reach for `performance` instead
when the concern is broader ("the whole site feels slow," bundle size,
caching strategy) rather than one specific metric.

## Common questions

- **"My LCP is fine on the page someone lands on, but users still say
  navigation feels slow — what's actually happening?"** For most sites,
  the LCP a user experiences most often is the *next* page they click
  into, not the one they landed on. The Speculation Rules API
  (`<script type="speculationrules">` with `eagerness: "moderate"`) lets
  the browser prerender likely-next pages on hover, collapsing that
  second LCP to near-zero — but analytics and other on-load side effects
  will fire when the prerender starts, not when the user actually
  navigates, unless they're gated on the `prerenderingchange` event.
- **"INP measurements look bad but no single click handler seems slow in
  testing — why?"** INP reports the *worst* interaction across the whole
  visit (98th percentile on high-traffic pages), not an average, so a
  handler that's fast 95% of the time but occasionally blocks on a slow
  third-party script or a large re-render is exactly what tanks the
  number even though most interactions felt fine.
- **"Why does my CLS show layout shifts I can't visually reproduce?"** A
  shift that happens during page load before the user has a chance to
  look — a font swap, a late-arriving ad — still counts even if nobody
  actually saw it happen; CLS is measured from frame data, not from what
  a human noticed.

## It's working if

- LCP's largest element is identifiable and loads without a client-side
  fetch delay
- No interaction handler blocks the main thread long enough to miss the
  next paint
- Every image, embed, and ad has reserved space before it loads
- All three metrics are being evaluated at the 75th percentile of real
  field data, not a single lab run

## Where it fits

The detailed reference underneath `web-quality-audit`'s Performance
category and `performance`'s Core Web Vitals section — reach for this
skill once you know which of the three metrics is the actual complaint.
