# Sample size worked example (qualitative, no formula derivation)

Read this when you need to walk a team through *why* a requested MDE
produces the sample size a calculator (e.g. Evan Miller's) returns, without
deriving the underlying statistics from scratch.

## Setup

Say a checkout page currently converts at a 10% baseline rate. The team
wants to detect an improvement. Two scenarios, same baseline:

- **Scenario A**: detect a 2 percentage-point absolute lift (10% → 12%,
  a 20% relative improvement). Plugging this into a standard two-sample
  proportion calculator at 95% significance / 80% power returns a sample
  size in the low thousands per variant.
- **Scenario B**: detect a 0.5 percentage-point absolute lift (10% →
  10.5%, a 5% relative improvement) — a smaller, more realistic effect
  for a mature, already-optimized checkout flow. The same calculator
  returns a sample size roughly 16x larger than Scenario A's, not 4x —
  because required sample size scales with roughly the inverse square of
  the effect size being targeted, so a 4x smaller MDE needs roughly a 16x
  larger sample.

## What this means in a planning conversation

- If a stakeholder says "let's detect any improvement, even a tiny one,"
  translate that into a concrete MDE and run it through a calculator
  before agreeing to a timeline — "any improvement" often implies
  Scenario-B-sized effects, which can turn a two-week test into a
  multi-month one at existing traffic levels.
- If traffic is fixed and the timeline is fixed, the honest response is to
  compute what MDE *is* detectable in that time/traffic budget, and check
  whether that MDE is still a meaningful, worth-shipping effect size — if
  it isn't, the test as scoped can't answer a useful question regardless
  of how long it runs.
- Low-baseline-rate metrics (e.g. a rare conversion event at 0.5% baseline)
  need dramatically more total traffic to hit the same absolute sample
  size requirement as a higher-baseline metric — this is a second,
  separate lever from MDE that's worth checking before committing to a
  test design.

## Practical workflow

1. Get the current baseline rate for the primary metric from existing
   data (not a guess).
2. Decide the smallest effect that would actually be worth shipping for —
   this is the MDE, and it's a product/business judgment call, not a
   statistical one.
3. Feed baseline rate, MDE, and desired significance/power into a sample
   size calculator (Evan Miller's calculator is the commonly used
   example) to get the required sample size per variant.
4. Divide by expected daily/weekly traffic per variant to get the
   required test duration.
5. Check the duration against the novelty-effect guidance (multiple full
   weekly cycles minimum) — if the sample-size-driven duration is shorter
   than that, extend to cover the novelty window; running the numerically
   "sufficient" sample in less time than that risks reading a novelty
   spike as the final result.
