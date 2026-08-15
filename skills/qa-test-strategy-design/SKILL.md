---
name: qa-test-strategy-design
description: This skill should be used when the user asks to "design a test strategy", "figure out our testing approach for this feature", "how much should we test this", "what's our test plan for this release", or needs to decide test-level allocation and where to focus coverage before writing any actual tests. Use for planning what to test and how much, at the feature or release level — not for writing or debugging individual tests (see e2e-testing for Playwright execution, test-driven-development for red-green-refactor at the unit level).
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "qa"
---

# QA Test Strategy Design

Test strategy answers "what should we test, how much, and at what level" —
before any test gets written. This is upstream of execution: once the
strategy says "this needs an E2E test," `skills/e2e-testing` covers writing
it; once it says "this needs unit coverage," `skills/test-driven-development`
covers the red-green-refactor loop. This skill does not duplicate either —
it decides which one applies where, and how much.

## The core tool: risk-based prioritization, not blanket coverage

"Test everything equally" is not a strategy, it's an excuse not to have
one. Risk-based testing (ISTQB's risk-based testing approach, and Cem
Kaner's work on it) scores each area by:

```
risk = likelihood of a defect × impact if it escapes to production
```

Procedure:

1. **Enumerate risk areas** for the feature/release: new code paths, high
   churn files (`git log --since=... --name-only | sort | uniq -c | sort
   -rn`), integration points with other services, areas with prior defect
   density, and anything touching money, auth, or data integrity.
2. **Score each area** likelihood (1-3) × impact (1-3) → a 1-9 risk tier.
   Impact for a hotel platform means: does this touch rate calculation,
   reservation state, payment capture, or PII? Those are impact-3 by
   default regardless of how simple the code looks.
3. **Map risk tier to depth and level**:
   - Tier 7-9 (high): unit + integration + at least one E2E happy-path +
     explicit negative/edge cases, reviewed by a second person.
   - Tier 4-6 (medium): unit + integration, E2E only if it's a new user
     journey rather than a variation of an existing one.
   - Tier 1-3 (low): unit tests only, or rely on existing regression
     coverage; do not add new E2E tests for low-risk cosmetic changes.
4. **Decide automate vs. explore vs. skip.** Repeatable, deterministic
   checks → automate. Areas needing human judgment about "does this feel
   right" (new UI flow, ambiguous spec) → route to
   `skills/exploratory-testing-techniques` instead of trying to script
   something that isn't well-specified yet. Truly low-risk, low-churn code
   → skip; a strategy that recommends testing everything has not actually
   prioritized.
5. **Write the strategy down** (even briefly, in the PR description or a
   short doc): scope, explicit out-of-scope, the risk matrix, which level
   gets used where, and exit criteria (what has to pass before release).

## Allocating test-pyramid effort

Default toward Mike Cohn's test pyramid shape: many fast unit tests, fewer
integration tests, fewest end-to-end tests — because E2E tests are slow,
flaky-prone, and expensive to maintain (see `skills/e2e-testing` for why).
A commonly cited allocation is roughly 70% unit / 20% integration / 10% E2E,
but treat that ratio as a starting heuristic, not a target to hit exactly.

- A service that is mostly business logic (rate calculation, availability
  rules) should skew even further toward unit tests.
- A service that is mostly integration glue (e.g. a booking-engine adapter
  calling three downstream PMS endpoints) legitimately has a fatter
  integration layer — forcing it toward the "ideal" ratio produces tests
  that mock away the exact risk you're trying to cover.

## Gotchas

- **100% code coverage is not a risk signal.** A line can be "covered" by
  a test that asserts nothing meaningful. Coverage percentage measures
  execution, not verification — don't let it substitute for the risk
  matrix.
- **The pyramid ratio is a guideline, not a law.** Services with heavy
  external-integration surface (common in hotel-industry PMS/booking
  integrations) legitimately need a thicker integration layer than the
  textbook ratio suggests.
- **A release with only low-risk changes needs almost no new tests.**
  Reflexively adding regression E2E tests for every PR inflates the suite's
  runtime and flakiness without reducing real risk — that budget is better
  spent on the tier 7-9 areas.
- **Partitioned tables need their own risk line item.** For this
  codebase's PostgreSQL tables partitioned by `hotel_id`, a query missing
  the partition key isn't a performance nit — it can silently return wrong
  or incomplete data across hotels; treat any new query path against a
  partitioned table as impact-3 regardless of how small the change looks.
- **Don't confuse "hard to test" with "low risk."** Legacy code with no
  seams is often exactly where the risk is highest and the existing
  coverage is thinnest — that combination should raise the tier, not lower
  it because "we can't test that easily."

## Real-world grounding

The test pyramid comes from Mike Cohn's *Succeeding with Agile* (2009),
formalized further in Google's public testing guidance distinguishing
"small/medium/large" tests by scope and speed. Risk-based test
prioritization is documented in the ISTQB Foundation syllabus and in Cem
Kaner's writing on risk-based test management — the core idea (likelihood
× impact drives depth of testing) predates and outlives any specific
tooling.
