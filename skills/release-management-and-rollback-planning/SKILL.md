---
name: release-management-and-rollback-planning
description: This skill should be used when the user asks to "plan a canary release", "design a rollout strategy", "write a rollback runbook", "should this go behind a feature flag", "how do we release this safely", or is deciding how a change reaches production and how to back out of it if something goes wrong. Use for the rollout/rollback strategy itself — not for the CI/CD pipeline mechanics or the code change being released.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Release Management and Rollback Planning

The rollback plan is written before the release, not improvised during an
incident. A release with no tested rollback path isn't a release plan —
it's a bet that nothing will go wrong.

## Progressive delivery

Ship changes to an increasing slice of traffic rather than all at once:

1. **Canary**: route a small percentage of traffic (often 1-5%) to the new
   version; compare error rate, latency, and business metrics against the
   baseline before continuing.
2. **Staged rollout**: widen in increments (e.g. 5% → 25% → 50% → 100%, or
   region by region / hotel-brand by hotel-brand), with a defined dwell
   time and defined go/no-go metrics at each step.
3. **Full rollout**: only after each prior stage cleared its metrics
   thresholds, not just "elapsed time with no complaints."

Config and data changes need the same discipline as code changes — a
"just a config update" mindset is exactly what removes the staged rollout
step that would have caught a bad config before it reached 100% of
traffic (see Real-world grounding).

## Feature flags

Feature flags decouple *deploy* (code reaches production, dormant) from
*release* (the behavior is turned on for users), which is what makes
canarying and instant rollback possible without a redeploy. Martin
Fowler's feature toggle taxonomy usefully separates flags by purpose:

- **Release toggles** — hide incomplete work behind a flag so it can merge
  to main continuously; short-lived, removed once fully released.
- **Ops toggles** — a kill switch for degrading gracefully under load or
  disabling a risky feature instantly if it misbehaves; can be long-lived.
- **Experiment toggles** — A/B testing; lifespan matches the experiment.
- **Permission toggles** — gate a feature by plan/tenant/role; often
  long-lived by design.

Flag hygiene matters: a release toggle left in the code long after full
rollout is tech debt (dead code paths, combinatorial testing burden) —
schedule its removal at rollout time, not "eventually."

## Rollback runbook

Write this before the release ships, as part of the release plan, not
during the incident:

- **Trigger thresholds**: the specific metrics/alerts that mean "roll
  back now" (error rate above X%, latency above Y ms, a specific alert
  firing) — decided in advance so the on-call isn't debating thresholds
  mid-incident.
- **Rollback mechanism**: flag flip (fastest, if behind a flag),
  redeploy of the previous version, or database rollback if a migration
  was involved — and which one applies for this specific release.
- **Migration compatibility**: if a schema migration shipped alongside
  the code, confirm the previous code version still runs correctly
  against the *new* schema (expand/contract migration pattern) — a
  rollback that un-deploys code but leaves an incompatible schema behind
  just creates a second incident.
- **Test the rollback path** in staging the same way the forward path was
  tested — an untested rollback is not a rollback plan, it's a hope.

## Gotchas

- **A config/rule change deserves the same staged rollout as a code
  change.** Treating configuration as exempt from progressive delivery
  because "it's not really a deploy" is precisely the gap that turns a
  small mistake into a global outage (see Cloudflare, below).
- **Canary metrics need a real baseline comparison, not just "no errors
  yet."** A canary that runs for 10 minutes with low traffic can look
  clean by chance; define minimum sample size/duration before declaring it
  safe to widen.
- **Rollback isn't automatically safe if a migration shipped with the
  release.** Rolling back code while the schema has already moved forward
  can be worse than the original bug if the old code can't handle the new
  schema — plan migrations to be backward-compatible for at least one
  release cycle.
- **Feature flags accumulate as debt if nobody owns removing them.** Track
  flag removal as a real follow-up task tied to the rollout, not an
  unowned "cleanup someday."
- **"It's a hotfix, we don't have time for canary" is exactly when canary
  matters most** — a rushed, untested global change under pressure is the
  highest-risk category of release, not an exception to the process.

## Real-world grounding

Cloudflare's July 2, 2019 global outage (documented in Cloudflare's own
public incident writeup, "Details of the Cloudflare outage on July 2,
2019") was caused by a single bad regular expression in a WAF rule that
was deployed globally in one step rather than through a staged/canary
rollout, causing CPU exhaustion across their edge network and roughly 27
minutes of widespread outage affecting a large share of Cloudflare-served
internet traffic. It's a well-documented example of exactly the gap
above: a change treated as "just a rule/config update" bypassing the
staged-rollout discipline applied to code deploys. Martin Fowler's 2017
article "Feature Toggles (aka Feature Flags)" is the widely-cited source
for the toggle-type taxonomy used above.
