## What it does

Decides *what* to test and how much, before any test gets written — a
risk-based allocation of unit/integration/E2E effort across a feature or
release. The defining constraint: "test everything equally" isn't a
strategy, it's the absence of one — this skill exists to force a
likelihood-times-impact judgment instead of uniform coverage.

## When to reach for it

Reach for this skill at the planning stage of a feature or release, when
the question is "what needs coverage and at what depth" — not while
writing an individual test. Once the strategy says "this needs an E2E
test," hand off to `e2e-testing` for the Playwright mechanics; once it says
"this needs unit coverage," hand off to `test-driven-development` for the
red-green-refactor loop itself. This skill decides which of those applies
where; it doesn't duplicate either.

## Common questions

- **"We have 100% code coverage on this module — are we done?"** No.
  Coverage measures execution, not verification — a line can be "covered"
  by an assertion-free test. Coverage percentage is not a risk signal and
  shouldn't substitute for the risk matrix this skill builds.
- **"Our booking-integration service doesn't fit the textbook 70/20/10
  unit/integration/E2E ratio — is that a problem?"** No — that ratio is a
  starting heuristic, not a law. A service that's mostly integration glue
  (adapter code calling several downstream endpoints) legitimately needs a
  thicker integration layer; forcing it toward the textbook ratio produces
  tests that mock away the exact risk being tested for.
- **"This is legacy code with no test seams — does that mean it's low
  priority?"** The opposite, usually. "Hard to test" and "low risk" are
  different axes; code that's both hard to test and has thin existing
  coverage is often where real risk concentrates, and that combination
  should raise its tier, not lower it because testing it is inconvenient.
- **"A release only touches low-risk cosmetic changes — do we still add
  new E2E regression tests?"** Generally no. Reflexively adding E2E tests
  to every PR inflates suite runtime and flakiness without reducing real
  risk; that effort is better spent on the tier 7-9 areas the risk matrix
  actually flags.

## It's working if

- Every feature/release has an explicit (even if brief) written risk
  matrix before tests are written
- Test-level allocation (unit/integration/E2E) is justified by the risk
  tier and the service's actual shape, not copied from a template ratio
- Any query path against a `hotel_id`-partitioned table is treated as
  impact-3 regardless of how small the change looks
- Low-risk changes are consciously left with minimal or no new test
  investment, not padded "for safety"

## Where it fits

A planning-stage, upstream skill — chains into `e2e-testing` and
`test-driven-development` for execution, and into
`exploratory-testing-techniques` for the parts of a feature that need human
judgment rather than a scripted check. Not a substitute for either; it
decides which applies.
