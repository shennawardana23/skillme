---
name: deep-research
description: Produce a thorough, multi-source, cited research report on an open-ended topic using web search and page-fetch tools, optionally fanning out across parallel sub-agents for broad subjects. Use when the user wants a deep dive, investigation, or "what's the current state of X" report on any topic, not a narrow technical question or a business decision memo.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Deep Research

Produce a thorough, cited report from multiple web sources on an open-ended topic. This is the broadest of the research-family skills — `research` answers one narrow, verifiable technical question against this codebase, and `market-research` produces a business decision memo; this skill produces a long-form, themed report on any topic, technical or not, sourced from the open web.

## When to Activate

- The user wants to research any topic in depth: "research X", "deep dive into Y", "investigate Z", "what's the current state of W"
- Competitive, technology, or trend analysis that isn't narrowly scoped to a business decision
- Any question requiring synthesis across many independent sources

## Tooling

Use whatever web search and page-fetch tools are available in the current environment (e.g. `WebSearch` and `WebFetch`). If MCP servers such as firecrawl or exa are configured and available, prefer them for broader crawl/search coverage — they are not required.

## Workflow

### Step 1: Understand the Goal

Ask 1-2 quick clarifying questions: "What's your goal — learning, making a decision, or writing something?" and "Any specific angle or depth you want?" If the user says "just research it," skip ahead with reasonable defaults rather than blocking on clarification.

### Step 2: Plan the Research

Break the topic into 3-5 research sub-questions before searching. Example — topic "impact of AI on healthcare": what are the main applications today; what clinical outcomes have been measured; what are the regulatory challenges; who are the leading companies; what's the market trajectory. Decomposing first prevents the search from wandering.

### Step 3: Execute Multi-Source Search

For each sub-question, search with 2-3 different keyword variations, mixing general and recency-focused queries. Aim for 15-30 unique sources total across all sub-questions. Prioritize academic, official, and reputable-outlet sources over blogs, and blogs over forums.

### Step 4: Deep-Read Key Sources

Fetch full content for the 3-5 most promising URLs rather than relying only on search snippets — snippets are frequently misleading about what a source actually argues.

### Step 5: Synthesize and Write the Report

```markdown
# [Topic]: Research Report
*Generated: [date] | Sources: [N] | Confidence: [High/Medium/Low]*

## Executive Summary
[3-5 sentence overview of key findings]

## 1. [First Major Theme]
- Key point ([Source Name](url))
- Supporting data ([Source Name](url))

## 2. [Second Major Theme]
...

## Key Takeaways
- [Actionable insight]

## Sources
1. [Title](url) — [one-line summary]

## Methodology
Searched [N] queries. Analyzed [M] sources. Sub-questions investigated: [list]
```

### Step 6: Deliver

Short topics: post the full report in chat. Long reports: post the executive summary and key takeaways, save the full report to a file if the user wants it kept.

## Parallel Research with Subagents

For broad topics, fan out sub-questions across parallel agents (e.g. one agent per 2 sub-questions), each returning sourced findings, then synthesize the results into one report in the main session. Don't fan out for narrow topics — the coordination overhead isn't worth it below roughly 3 sub-questions.

## Quality Rules

1. Every claim needs a source. No unsourced assertions.
2. Cross-reference. If only one source says something, flag it as unverified rather than presenting it as settled.
3. Recency matters — prefer sources from the last 12 months for anything fast-moving.
4. Acknowledge gaps explicitly if a sub-question couldn't be answered well.
5. No hallucination — if the answer isn't known, say "insufficient data found" rather than filling the gap plausibly.
6. Separate fact from inference — label estimates, projections, and opinions clearly as such.

## Gotchas

- A search engine's top results are frequently SEO-optimized restatements of the same original source — cross-referencing five pages that all cite the same one study is not five independent confirmations.
- Snippets shown in search results are sometimes taken out of context by the search tool itself; always read the surrounding paragraph in the fetched page before citing a claim.
- "Recent" articles that recycle an older, unlabeled statistic will make stale data look current — check the original source's publication date, not just the citing article's.

## Real-world grounding

Wikipedia's own "citation needed" convention, and its broader sourcing policy, exist precisely because uncited claims accumulate credibility they haven't earned simply by sitting in a plausible-looking document — the same failure mode this skill's "every claim needs a source" rule guards against. On the cross-referencing point, retracted or corrected news stories (a routine, publicly logged occurrence at most major outlets, tracked by services like retractiondatabase.org for the scientific literature) are the standard illustration of why a single source, however reputable, is not the same as independent confirmation.

## Verification

- [ ] The topic was broken into 3-5 sub-questions before searching
- [ ] 15-30 unique sources were consulted across the full report
- [ ] At least 3-5 key sources were read in full, not just as snippets
- [ ] Every material claim carries an inline citation
- [ ] Gaps and single-sourced claims are explicitly flagged, not smoothed over
