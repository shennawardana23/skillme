## What it does

Runs a full-site quality review across the same four categories Google
Lighthouse audits — Performance, Accessibility, SEO, and Best Practices —
and returns one prioritized report instead of four separate ones. The
non-obvious part is the severity model underneath it: a finding's
category doesn't decide its urgency, its actual impact does. A missing
CSP header (Best Practices) outranks a missing meta description (SEO)
every time, so the skill sorts by Critical/High/Medium/Low across all
four categories together, not category-by-category.

## When to reach for it

Reach for this when someone wants a holistic "how healthy is this page"
answer — before a launch, during a periodic review, or when a stakeholder
asks "how are we doing on quality" without specifying which axis. If the
ask is already narrowed to one axis — "why is my LCP slow," "is this
accessible," "fix my meta tags," "is this a security risk" — go straight
to `core-web-vitals`/`performance`, `accessibility`, `seo`, or
`best-practices` instead; each of those goes considerably deeper on its
one area than this skill's summary-level pass does.

## Common questions

- **"Lighthouse's JSON report doesn't have the audit names this skill
  mentions anymore — is the skill out of date?"** No — Lighthouse v13
  (October 2025+) migrated the Performance category from per-opportunity
  audits to consolidated Performance Insight audits (e.g. several CLS-
  related checks merged into a single `cls-culprits-insight`). The
  underlying advice is unchanged; only the report's shape moved. Treat a
  newer Lighthouse JSON output as a superset of the older one, not a
  contradiction.
- **"Is a Critical finding always a security issue?"** No — Critical
  spans complete failures too (a totally broken page, not just a
  vulnerability). The four levels are about blast radius and urgency,
  not which of the four categories a finding came from.

## It's working if

- The final report ranks findings by actual severity across all four
  categories, not grouped and equally weighted by category
- Every finding names a specific file/line and a concrete fix, not a
  restated Lighthouse rule name
- The summary states a real recommended order to fix things in, not just
  a bucketed list

## Where it fits

The entry point for the web-quality family — it's the one to reach for
first when scope is unclear, and its own categories are exactly
`performance`/`core-web-vitals`, `accessibility`, `seo`, and
`best-practices`, each a deeper reference for its own slice once the
audit has told you which one actually matters here.
