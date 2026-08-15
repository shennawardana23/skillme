---
name: market-research
description: Conduct business market research, competitive analysis, investor due diligence, and market-sizing that supports a business decision rather than producing research theater. Use when sizing a market (TAM/SAM/SOM), comparing competitors, vetting an investor or fund before outreach, or evaluating a vendor/technology adoption decision.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Market Research

Produce research that changes a decision, not research that just describes a topic. This is the business/strategic-decision counterpart to `research` (codebase-technical questions) and `deep-research` (open-ended web reports) — use this one specifically when the reader is a founder, investor, or operator deciding whether to enter a market, fund a deal, or adopt a vendor.

## When to Activate

- Sizing a market (TAM/SAM/SOM) before a build or funding decision
- Comparing competitors or adjacent products before positioning
- Vetting an investor or fund before outreach (check-size fit, thesis fit, red flags)
- Evaluating a vendor or technology adoption decision (lock-in, cost, risk)
- Pressure-testing a thesis before committing real money or headcount to it

## Research Standards

1. Every important claim needs a source. If you can't source it, label it an estimate.
2. Prefer recent data; explicitly flag anything more than 12-18 months old as potentially stale.
3. Include contrarian evidence and the downside case, not just the case that supports the thesis.
4. Translate findings into a decision, not just a summary — the output is a recommendation, not a report for its own sake.
5. Separate fact, inference, and recommendation visibly — a reader should be able to tell which is which without re-reading.

## Research Modes

### Investor / Fund Diligence
Collect: fund size, stage, and typical check size; relevant portfolio companies; public thesis and recent activity (recent deals, public statements); reasons the fund is or isn't a fit; any obvious red flags (fund nearing end of life, no recent activity in the category, portfolio conflicts).

### Competitive Analysis
Collect: product reality from actually using or demoing the product, not marketing copy; funding and investor history if public; traction signals if public (hiring rate, review volume, app rankings — never invent a number you can't source); distribution and pricing clues; strengths, weaknesses, and the specific gap your positioning would exploit.

### Market Sizing
Use top-down estimates from named, dated reports or public datasets, cross-checked with a bottom-up estimate built from realistic customer-acquisition assumptions (price × addressable buyers × realistic penetration rate). State every assumption behind each leap in logic explicitly — a sizing estimate with hidden assumptions is not verifiable and should not be trusted, including by the person who made it.

### Technology / Vendor Research
Collect: how it actually works (not the pitch); trade-offs and real adoption signals (who else uses it, for what scale); integration complexity; lock-in, security, compliance, and operational risk if this vendor disappears or changes pricing.

## Output Format

1. Executive summary (3-5 sentences, states the recommendation up front)
2. Key findings
3. Implications
4. Risks and caveats (including the strongest counter-argument to your own recommendation)
5. Recommendation
6. Sources

## Gotchas

- A market-sizing number with no bottom-up cross-check is usually wrong in the optimistic direction — top-down TAM figures from industry reports are frequently built by multiplying a broad category size by a small percentage, which feels rigorous but compounds whatever bias was in the original report.
- "No red flags found" on investor diligence is not the same as "diligence complete" — if you didn't check public deal activity in the last 12 months, say that explicitly rather than letting silence imply a clean bill of health.
- Traction metrics quoted in press releases or pitch decks are marketing claims, not verified facts, until you can trace them to a primary source (a filing, a named report, an on-the-record statement) — treat them as inference-labeled, not fact-labeled.
- A vendor's public case studies are selected for success; absence of a published failure case is not evidence failures don't happen.

## Real-world grounding

WeWork's 2019 IPO prospectus is a widely documented cautionary case for this exact skill: public market analysts and journalists picked apart the company's self-reported market-sizing and community-adjusted-EBITDA metrics within days of the S-1 filing, and the scrutiny that surfaced (governance red flags, unverified TAM claims, insider conflicts) directly led to the withdrawn IPO. It is the standard example cited in startup and VC diligence writing for why unsourced sizing numbers and unexamined red flags don't survive contact with a real decision.

## Verification

Before delivering:
- [ ] All numbers are sourced or explicitly labeled as estimates
- [ ] Data older than ~12-18 months is flagged as potentially stale
- [ ] The recommendation follows visibly from the evidence presented
- [ ] The strongest counter-argument or downside case is included, not omitted
- [ ] The output makes the reader's decision easier, not just more informed
