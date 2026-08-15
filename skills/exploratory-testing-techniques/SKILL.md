---
name: exploratory-testing-techniques
description: This skill should be used when the user asks to "do some exploratory testing", "write a test charter", "explore this feature for bugs", "run a testing session", or needs an unscripted, human-judgment-driven approach to finding bugs in an ambiguous or newly-built feature that isn't well-specified enough for scripted automation yet. Use for structured, charter-driven manual investigation — not for writing automated Playwright/unit tests (see e2e-testing, test-driven-development).
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Exploratory Testing Techniques

Exploratory testing (ET), as defined by James Bach and Michael Bolton's
Rapid Software Testing methodology, is simultaneous learning, test design,
and test execution — the tester designs the next test based on what the
last one revealed, rather than following a pre-written script. It is a
disciplined skill, not "just clicking around": the discipline comes from
charters, time-boxing, and oracles, not from a script.

## Session-Based Test Management (SBTM)

Structure exploration so its value is traceable, per Jonathan Bach's
Session-Based Test Management:

1. **Write a charter** before the session starts:
   `Explore <area> with <resources/tools> to discover <information>.`
   Example: "Explore the discount-stacking checkout flow with expired and
   near-expiry promo codes to discover incorrect price calculations."
2. **Time-box the session** — 60-120 minutes is typical. Longer sessions
   lose focus; shorter ones don't give exploration room to follow leads.
3. **Explore and take notes continuously** — not just when you find a bug.
   Note what you tried, what you expected, what happened, and anything
   that felt off but wasn't clearly wrong yet.
4. **Log bugs immediately**, with repro steps, while the state is still in
   front of you — see `bug-triage-and-severity-classification` for what
   happens to the report next.
5. **Write a session report / debrief**: charter, session duration, % of
   time on-charter vs. opportunity testing (following an unplanned but
   promising lead), bugs found, issues/questions raised, areas still
   uncovered. A session without a report is testing debt — the value of
   ET is in what's traceable afterward, not just in the tester's memory.

## Oracles: how you judge "is this actually wrong"

Without a full spec, use Michael Bolton's HICCUPPS(F) heuristic to decide
whether an observed behavior is a bug:

- **H**istory — does it behave differently than it used to?
- **I**mage — does it embarrass the company/product's reputation?
- **C**omparable products — do similar products/features behave
  differently in a way that suggests this is wrong?
- **C**laims — does it contradict documentation, specs, or marketing
  claims?
- **U**ser expectations — would a reasonable user be surprised?
- **P**roduct — is it internally consistent with the rest of the product?
- **P**urpose — does it defeat the purpose the feature was built for?
- **S**tatutes — does it violate a law/regulation (e.g. accessibility,
  data privacy)?
- **F**amiliar problems — does it resemble a known bug pattern/class?

## Whittaker's Tours — concrete charter starting points

James Whittaker's "tours" (from *Exploratory Software Testing* and *How
Google Tests Software*) are named exploration angles to draw charters from
when you don't know where to start:

- **Guidebook Tour** — follow the documented/happy-path user journey
  exactly as a manual would describe it.
- **Money Tour** — exercise only the features that generate revenue or
  that the business cares most about (for a booking platform: the actual
  reservation and payment path).
- **Landmark Tour** — visit every major screen/feature once, breadth over
  depth, to map what exists.
- **FedEx Tour** — follow a piece of data through the entire system from
  entry to its final destination (e.g., trace a single reservation from
  booking-engine submission through to PMS record).
- **Back Alley Tour** — deliberately visit the least-used, least-polished
  corners of the product — old settings pages, rarely-touched admin
  screens.
- **Obsessive-Compulsive Tour** — repeat the same action many times, or in
  rapid succession, or with tiny variations, looking for state corruption.

Pick two or three tours relevant to the feature under test and turn each
into a charter rather than trying to run all of them.

## Gotchas

- **Undirected clicking is not exploratory testing.** ET without a charter
  and without notes is indistinguishable from random poking — it produces
  no traceable record of what was and wasn't covered, and can't be
  defended in a retro or audit.
- **A session report is the deliverable, not the bug count.** A
  zero-bug-found session with a clear report ("covered X, Y, Z, found
  nothing wrong") is valuable evidence of coverage; a session with bugs
  found but no report is not reproducible by anyone else.
- **Oracles prevent "looks fine to me" false negatives.** Without an
  explicit heuristic like HICCUPPS, testers tend to only flag things that
  crash or throw — HICCUPPS deliberately surfaces "technically works but
  violates user expectations" issues that scripted assertions miss.
- **Exploratory testing complements automation, it doesn't replace it.**
  Once ET reveals a concrete bug or a well-defined behavior, that behavior
  becomes a candidate for a permanent scripted regression test (see
  `test-driven-development` / `e2e-testing`) — ET's job is finding the
  unknown unknowns, not re-verifying the same known behavior every release.

## Real-world grounding

Session-Based Test Management was created by Jonathan Bach (published as
"Session-Based Test Management," 2000) as a way to make exploratory testing
manageable and reportable within James Bach and Michael Bolton's broader
Rapid Software Testing methodology. James Whittaker's tours heuristic
originates from his book *Exploratory Software Testing* (2009) and is
described in detail in *How Google Tests Software* (2012), documenting how
Google's own test teams used named tours to structure exploratory sessions
at scale.
