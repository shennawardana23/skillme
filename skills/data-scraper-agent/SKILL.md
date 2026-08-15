---
name: data-scraper-agent
description: Build an automated data-collection agent that scrapes a public source on a schedule, enriches results with a cheap/free LLM, stores them in a database, and improves scoring from user feedback over time — runnable entirely on free-tier infrastructure. Use when the user wants to monitor, collect, or track any public data automatically (job boards, prices, news, GitHub repos, sports scores, listings).
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Data Scraper Agent

Build a scheduled agent that collects public data, enriches it with an LLM, stores it, and learns from the user's accept/reject decisions — without requiring paid infrastructure.

## When to Activate

- User wants to scrape or monitor a public website or API on a recurring schedule
- User says "build a bot that checks...", "monitor X for me", "collect data from..."
- User wants to track jobs, prices, news, repos, sports scores, events, or listings
- User wants automated collection that improves over time based on their decisions

## The Three-Layer Architecture

```
COLLECT → ENRICH → STORE
scraper     LLM      database
on a      scores/  (Notion/Sheets/
schedule  classifies  Supabase/etc.)
```

Two principles make this survive free-tier limits:

1. **Batch every LLM call.** Never call the model once per item — batch 5+ items into a single call. Scraping 30 items with one call each burns a rate limit instantly; batching the same 30 items into 6 calls of 5 stays inside almost any free tier.
2. **Cascade through a model fallback chain.** On quota exhaustion (HTTP 429) or an unavailable model (404), fall through to the next cheaper/faster model in a pre-ordered list rather than failing the run. Order fastest/cheapest-and-most-available first, since free tiers usually grant the smallest model the highest rate limit.

## Workflow

### Step 1: Understand the goal

Ask (or infer from an unambiguous request):

1. **What to collect** — URL / API / RSS / public endpoint?
2. **What to extract** — which fields matter (title, price, URL, date, score)?
3. **Where to store** — Notion, Google Sheets, Supabase, or a local file?
4. **How to enrich** — should the LLM score, summarize, classify, or match each item against user context?
5. **How often** — hourly, daily, weekly?

Common shapes to recognize: job boards scored against a resume, product prices with drop alerts, GitHub repos summarized on new release, news classified by topic/sentiment, sports results tracked in a running table, event listings filtered by interest.

### Step 2: Lay out the project

```
my-agent/
├── config.yaml         # user-tunable: keywords, filters, storage provider, AI settings
├── profile/context.md   # user context the LLM uses (resume, interests, criteria)
├── scraper/
│   ├── main.py          # orchestrator: scrape → enrich → store
│   ├── filters.py        # rule-based pre-filter, runs before the LLM
│   └── sources/*.py      # one file per data source
├── ai/
│   ├── client.py         # LLM REST client with model fallback
│   ├── pipeline.py       # batches items into LLM calls
│   └── memory.py         # turns user feedback into a scoring bias prompt
├── storage/*_sync.py     # push + dedup logic for the chosen backend
├── data/feedback.json    # accumulated user decisions (auto-updated)
└── .github/workflows/scraper.yml   # cron schedule (or any scheduler)
```

Read `references/implementation.md` for the full working Python implementation of every file above — the scraper client, batching pipeline, feedback memory, storage sync, config template, and a GitHub Actions cron workflow — before wiring up a real source. This architecture is language-agnostic: the same three-layer split (collect/enrich/store) and the same batching + fallback-chain principles apply if the user's project standardizes on Go or TypeScript instead of Python.

### Step 3: Pick a scraping pattern per source

| Pattern | When to use |
|---|---|
| **REST API** | Source exposes JSON — easiest and most reliable, prefer this whenever available |
| **HTML scraping** | No API exists; parse the rendered HTML with a CSS-selector-based parser |
| **RSS feed** | Source publishes a feed — usually the most stable option for news/blogs |
| **Paginated API** | Loop pages until the source signals no more results, respecting rate limits between calls |
| **JS-rendered page** | Content only appears after client-side JavaScript runs — requires a headless browser, not a plain HTTP client |

Every source function should return a list of dicts/structs with a consistent minimal schema (name, URL, source, date found) so the enrichment and storage layers don't need source-specific branching.

### Step 4: Enrich in batches, store with dedup

Feed each batch to the LLM with: the batch of items, the user's context file, the user's stated priorities, and (once available) a feedback-derived preference prompt built from past positive/negative decisions. Before pushing to storage, deduplicate against existing records (URL is usually the natural dedup key) so repeated runs don't create duplicate rows.

### Step 5: Close the feedback loop

After the user marks items as accepted/rejected in the storage layer (e.g. a Notion "Status" property), a separate sync step should read those statuses back into `data/feedback.json`. The next run's enrichment prompt includes a summary of liked vs. rejected items, biasing future scoring toward what the user actually acts on rather than what they said they wanted at setup time.

## Gotchas

- Fetching a JavaScript-rendered page with a plain HTTP client returns the pre-render HTML shell — an empty or near-empty result that looks like a bug but is actually a wrong-tool choice; check whether the content appears in the raw HTML before assuming scraping is broken.
- A `maxOutputTokens` limit set too low silently truncates a batched JSON response mid-object, producing a parse error that looks like a prompt problem — raise the token budget before debugging the prompt when batch responses fail to parse.
- Scraping without a delay between requests can get the source's IP-banned within minutes on stricter sites, even for read-only, low-volume collection — a fixed delay between requests is cheap insurance.
- `robots.txt` and a site's terms of service govern what automated collection is permitted, independent of whether the technical scrape is possible — check both before automating collection from any given source, not just for personal/hobby projects.
- Parsing RSS/XML feeds with a language's default XML parser (e.g. Python's stdlib `xml.etree.ElementTree`) is vulnerable to XXE and billion-laughs attacks against untrusted feed content — use a hardened parser (e.g. `defusedxml`) instead, since the feed source is untrusted external input.

## Anti-Patterns to Avoid

| Anti-pattern | Problem | Fix |
|---|---|---|
| One LLM call per item | Hits rate limits almost immediately | Batch 5+ items per call |
| Hardcoded keywords/config in code | Not reusable across sources or users | Move all tunables to `config.yaml` |
| No delay between requests | Risk of IP ban | Add a fixed delay between requests |
| Secrets committed to source | Security exposure | `.env` + secret manager only, never in code |
| No deduplication before storage push | Duplicate rows accumulate every run | Always check the natural key (usually URL) before pushing |
| Ignoring `robots.txt` / ToS | Legal and ethical risk | Check both before automating any source |
| JS-rendered site fetched with a plain HTTP client | Empty/incomplete result | Use a headless browser, or find the underlying API the page itself calls |

## Real-world grounding

Automated scraping of a platform's data is not just a technical question — it has real legal history. *hiQ Labs v. LinkedIn* (2017-2022) was a widely reported US case over whether scraping publicly viewable LinkedIn profile data violated anti-hacking law; the litigation ran for years and produced conflicting rulings at different stages. Whatever the legal outcome in a specific jurisdiction, the case is a concrete illustration that reviewing a source's terms of service and `robots.txt` before automating collection is a real risk-management step, not a formality — treat it as part of Step 1, not an afterthought.

## Quality Checklist

Before marking the agent complete:

- [ ] `config.yaml` controls all user-facing settings — nothing hardcoded
- [ ] `profile/context.md` holds user-specific context for LLM matching
- [ ] Dedup by natural key (URL) runs before every storage push
- [ ] The LLM client has a multi-model fallback chain
- [ ] Batch size is 5+ items per LLM call, never 1
- [ ] Secrets live only in `.env` / the scheduler's secret store, never in committed code
- [ ] The scheduled workflow persists updated feedback history after each run
- [ ] `robots.txt` and terms of service for each source have been checked
