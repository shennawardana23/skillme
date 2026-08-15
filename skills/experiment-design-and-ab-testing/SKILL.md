---
name: experiment-design-and-ab-testing
description: Guides the product-experimentation layer of A/B testing — writing a falsifiable hypothesis, sizing a test with minimum detectable effect (MDE), reading statistical significance without falling into the peeking problem, accounting for the novelty effect, and choosing guardrail metrics. Use when the user asks to "design an A/B test", "how big a sample do I need", "can we stop this test early, it hit p<0.05", "is this result significant", "what guardrail metrics should we track", or wants to review an experiment plan or an in-flight experiment's results before deciding to ship. For the technical execution of browser/UI tests, see the e2e-testing skill instead — this skill does not cover that layer.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Experiment Design and A/B Testing

This skill covers the product-decision layer of experimentation: what
question the experiment answers, how big it needs to be, how to read the
result honestly, and what else to watch besides the metric you're
optimizing. It does not cover how to technically wire up a browser test,
mock network calls, or automate a UI flow — for that, use
`skills/e2e-testing/` (Playwright page objects, CI config, flaky-test
diagnosis). The two are complementary: e2e-testing verifies the code works
before a rollout; this skill decides whether the rollout's outcome is real.

## 1. State a falsifiable hypothesis before running anything

Write the hypothesis in a form that could be proven wrong, before looking
at any data:

> "Changing X will change [primary metric] by at least [MDE] for
> [population], measured over [duration], because [causal mechanism]."

A hypothesis like "let's test a new checkout flow and see what happens" is
not falsifiable — there's no result that could disprove it, so any outcome
gets rationalized as a win after the fact. Committing to the metric, the
minimum effect size, and the duration *before* the test starts is what
prevents the analysis from being reverse-engineered to match whatever the
data happened to show.

## 2. Minimum detectable effect (MDE) drives sample size — and the relationship is steep, not linear

MDE is the smallest true effect you want the test to be able to detect
reliably. The qualitative shape of the tradeoff: required sample size grows
roughly with the *square* of how small an effect you want to detect —
halving the MDE you're targeting roughly quadruples the sample size
needed, not doubles it. This is why "let's just detect any improvement, no
matter how small" is not a free request — chasing a 0.5% lift on a metric
that needs a 3% lift's worth of sample to detect reliably means the test
will run far longer than the team expects, or will report "no significant
difference" on an effect that's real but too small for the sample to see.

Practical implications:
- Pick the MDE based on what effect size would actually be worth shipping
  for — not the smallest effect that's theoretically possible — because
  that's what keeps the required sample size achievable in a reasonable
  timeframe.
- Practitioners size tests with a calculator rather than deriving the
  formula by hand each time — Evan Miller's sample size calculator is the
  standard, widely used example of this kind of tool. Feed it baseline
  conversion rate, MDE, and desired power/significance, and it returns the
  sample size (or, for a fixed sample, the smallest MDE that's detectable).
- Low-traffic metrics (e.g. a rare conversion event) need either a much
  longer test window or a larger MDE target — there is no shortcut around
  the sample-size requirement by wanting the answer sooner.

## 3. Statistical significance — what the 95%/p<0.05 convention actually means, and its most common misuse

A p-value is the probability of seeing a result at least this extreme if
there were truly no effect. The 95% confidence / p<0.05 convention means:
if there's genuinely no difference, a test run this way will still show a
"significant" result by chance about 1 time in 20 — that's the false
positive rate you're accepting by using this threshold, not a certainty
that a p<0.05 result is real.

**The peeking problem (optional stopping)**: checking results daily and
stopping as soon as p crosses 0.05 inflates the true false-positive rate
far above the nominal 5% — sometimes to 20-30%+ depending on how often you
peek — because you're effectively running many sequential tests and taking
the first one that happens to cross the threshold by chance, rather than
running the one test you pre-committed to. This is a well-documented
statistical trap in A/B testing, not a hypothetical edge case: a test that
"just crossed p<0.05 after 2 days" when it was sized for two weeks of
traffic is exhibiting exactly this failure mode, and should not be
stopped on that basis.

Mitigations:
- Pre-register the sample size and/or duration before starting (from
  step 2), and hold to it — don't stop early just because significance
  was hit, and don't keep running past the planned sample "just to be
  sure" either, since that's peeking in the other direction.
- If a business genuinely needs to check results continuously and act on
  them early, that requires a sequential-testing method designed for
  repeated looks (e.g. always-valid inference / sequential testing
  approaches) — not repeated ordinary significance tests, which is the
  actual bug in "checking daily and stopping when it's significant."

## 4. The novelty effect

A change frequently performs well in the first days simply because it's
new and different — users notice it, engage out of curiosity, or actively
try the new option — not because it's actually better. As the novelty
wears off, the metric commonly regresses toward (or below) the baseline.
Reading a novelty-driven spike as the durable effect is a common way
experiments get misread as clear wins.

Mitigation: run the test long enough to span the novelty decay window —
in practice this usually means multiple full weekly cycles (not just a
few days), since user behavior has a weekly rhythm and the novelty effect
needs time to fade before the metric stabilizes. If the effect size is
clearly larger in week 1 than in week 3-4 of the same test, treat that as
a signal of novelty rather than an accelerating win, and weight the later,
more stable weeks more heavily in the shipping decision.

## 5. Guardrail metrics

The primary metric tells you whether the variant won on the thing you're
optimizing for; guardrail metrics tell you whether it won *by breaking
something else*. Define guardrails before the test starts, alongside the
primary metric and hypothesis — not after seeing a surprising result,
which turns guardrail selection into after-the-fact justification.

Canonical example: a checkout redesign that increases conversion rate
(primary metric) but silently increases page load time or step-abandonment
on a specific device class (guardrail). Shipping on the primary metric
alone would ship a regression that a different metric would have caught.

Practical guidance:
- Pick guardrails from a different layer than the primary metric —
  performance (load time, error rate), a different funnel stage's
  metric, or a segment-level breakdown — so a manipulation of the primary
  metric that comes at another layer's expense gets caught.
- A guardrail regression doesn't automatically block shipping — it forces
  an explicit tradeoff conversation instead of a default "primary metric
  won, ship it."

## Gotchas

- Stopping a test the moment it crosses p<0.05 is the single most common
  practitioner mistake in this space — see the peeking problem in
  section 3. If someone asks "can we stop early, it just hit
  significance," the answer is to check whether the pre-registered
  sample size/duration has actually been reached, not whether p<0.05 has
  been reached.
- A wildly better-than-expected result in the first few days is more often
  a novelty effect or an instrumentation bug than a genuinely large
  effect — treat an unusually large early lift as a reason to keep running
  and check tracking, not a reason to ship immediately.
- Segmenting results after the fact to find *some* subgroup where the
  variant "won" (when the overall result was flat or negative) is a form
  of p-hacking — with enough subgroups checked, one will show significance
  by chance. Pre-specify the segments you'll examine, or treat post-hoc
  segment wins as a hypothesis for a new, dedicated test rather than a
  finding to ship on.
- A statistically significant result is not automatically a practically
  significant one — a 0.1% lift can be statistically significant with
  enough traffic while being far too small to be worth the added
  maintenance cost of the variant. Compare the observed effect to the MDE
  that was worth shipping for (step 1), not just to the p-value.
- Running the same experiment simultaneously with another team's
  overlapping experiment on the same users, without a tracked
  interaction, can contaminate both results — check for concurrent
  experiments touching the same surface or population before trusting a
  clean result.

## Real-world grounding

The peeking/optional-stopping problem is one of the most widely documented
statistical pitfalls in industry A/B testing practice — repeatedly checking
p-values and stopping on the first significant read is well known to
inflate the false-positive rate far above the nominal 5%, which is why
major experimentation platforms (and the broader online-experimentation
literature) specifically warn against it and build sequential-testing
safeguards to allow legitimate early stopping. Evan Miller's publicly
available sample size calculator (evanmiller.org) is a standard,
widely-referenced example of the kind of tool practitioners actually reach
for to translate a baseline rate and MDE into a required sample size,
rather than deriving the underlying power-analysis formula by hand each
time. The novelty effect and the practice of tracking guardrail metrics
alongside a primary success metric are both standard, well-documented parts
of mature online-experimentation practice at companies that run
experimentation at scale, used specifically to avoid shipping changes that
look like wins only because they are new, or that win on one metric while
quietly regressing another.

See `references/sample-size-worked-example.md` for a worked walkthrough of
sizing a test from a baseline rate and MDE.
