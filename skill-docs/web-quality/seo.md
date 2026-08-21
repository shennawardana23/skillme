## What it does

Covers technical SEO (crawlability, sitemaps, canonical URLs), on-page
SEO (titles, meta descriptions, heading hierarchy), and structured data
(JSON-LD for articles, products, FAQs, breadcrumbs). The framing that
matters: on-page and technical SEO together are a comparatively small
share of what actually moves search ranking (content quality/relevance
and backlinks dominate), so this skill is explicit about what it can and
can't influence rather than implying that fixing meta tags will move
rankings on its own.

## When to reach for it

Reach for this for anything about search visibility specifically —
crawlability, meta tags, structured data, sitemap issues, or "why isn't
this page showing up in search." Page-experience-driven ranking factors
(Core Web Vitals) live in `core-web-vitals`; this skill only notes that
they matter, it doesn't cover the metrics themselves.

## Common questions

- **"Should we block AI crawlers (ChatGPT search, Perplexity, Gemini
  Overviews) in robots.txt?"** Not wholesale — each has its own
  `robots.txt` user-agent (`OAI-SearchBot`, `PerplexityBot`,
  `GoogleOther`, `Google-Extended`, `ClaudeBot`, etc.), and a blanket
  `Disallow` removes the site from that specific bot's citations. Decide
  per-bot, and note that this whole area has no confirmed ranking
  signals as of 2026 — it's citation visibility, not search ranking.
- **"Is `llms.txt` worth adding?"** It's a proposed convention (a
  Markdown index of important pages at `/llms.txt`), adoption is around
  0.015% of sites, and no major AI vendor has confirmed reading it as of
  mid-2026. Treat it as a five-minute speculative add for a content
  site, not something to reorganize content around.
- **"Our sitemap has 60,000 URLs in one file — is that a problem?"**
  Yes — sitemaps are capped at 50,000 URLs or 50MB each; a site past that
  needs a sitemap index referencing multiple child sitemaps instead of
  one oversized file.

## It's working if

- Every page has a unique title (50-60 chars) and meta description
  (150-160 chars)
- There's exactly one `<h1>` per page and no skipped heading levels
- `robots.txt` and meta robots tags don't accidentally block important
  pages or their required resources
- Structured data validates against Google's Rich Results Test where
  used

## Where it fits

One of the four categories under `web-quality-audit`; for the
page-experience ranking factor it references but doesn't cover, see
`core-web-vitals`.
