---
name: engineering-values-and-tradeoff-principles
description: Use when the user asks to "should we prioritize shipping fast or reliability here", "define an error budget", "how much risk can we take on this release", "we keep debating speed vs quality", "write down our engineering principles", or is deciding how to balance reliability/velocity/quality tradeoffs at a policy level rather than for one specific decision. Guides turning a philosophical tradeoff into an operational, measurable policy, grounded in Google SRE's error budget concept.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Engineering Values and Tradeoff Principles

Teams that debate "should we prioritize speed or reliability" as an
abstract values question tend to relitigate it every release, because
nobody wrote down a rule that ends the debate the same way each time. The
practical fix demonstrated by Google's Site Reliability Engineering
practice: convert the philosophical tradeoff into a measurable, agreed
threshold ahead of time, so in the moment it's a lookup, not a debate.

## The error budget: reliability vs. velocity as one number

Google SRE's error budget makes this concrete: no service is ever
required to be 100% reliable (100% is the wrong target — it eliminates
any room to ship changes, since every change carries some risk of an
incident, and it costs disproportionately more to chase each additional
"nine" of reliability). Instead:

1. **Set an explicit reliability target below 100%** (e.g., 99.9%
   availability over a rolling window) based on what the product actually
   needs — a target chosen by asking what unreliability users would
   actually notice and be harmed by, not by defaulting to the highest
   number that sounds impressive.
2. **The gap between the target and 100% is the error budget** — the
   amount of unreliability "spent" on this service is allowed within the
   window before it's considered a problem.
3. **Spend the budget on deliberate risk**: shipping features, running
   experiments, doing risky migrations — anything that might cause an
   incident is implicitly "spending" budget.
4. **When the budget is exhausted, the pre-agreed policy activates
   automatically** — commonly a freeze on further risky releases until
   reliability recovers back within budget. This is what makes it
   operational rather than aspirational: the rule was agreed before the
   pressure of "we really need to ship this" existed, so it isn't
   renegotiated under that pressure.

The key move, generalizable beyond SRE: **pick the tradeoff's threshold
while calm, write it down, and let the written threshold — not a fresh
argument each time — decide in the moment the tradeoff actually arises.**

## Applying the same pattern to other tradeoffs

The error budget is one instance of a general technique. Apply the same
three steps to any recurring speed-vs-quality-style debate in your
organization:

1. **Name the two things in tension** explicitly (e.g., "ship velocity"
   vs. "production stability," "developer time saved" vs. "runtime
   resource cost," "review thoroughness" vs. "time to merge").
2. **Pick a measurable proxy and a threshold for each side**, agreed on
   without an active deadline pressuring the number (e.g., "no more than
   2 production incidents per month before we pause new feature work and
   focus on stability," "PRs over 400 lines require two reviewers, under
   get one").
3. **Write down what happens automatically when the threshold is
   crossed**, and who has authority to grant an explicit, logged exception
   — an exception process is fine, but it should require a specific
   named person to approve and should be visible after the fact, not a
   silent bypass.

## Gotchas

- An error budget (or any tradeoff policy) that's violated without
  consequence a few times stops functioning as a real policy — if the
  agreed freeze "doesn't really apply this time" repeatedly, the team has
  reverted to relitigating the tradeoff live, which is the exact failure
  mode the policy was meant to prevent.
- 100% reliability is not a more "correct" target than 99.9% — it's
  usually the wrong target, since it implies zero acceptable risk for any
  change ever, which is incompatible with shipping anything. State the
  target as a deliberate choice with a stated reason, not as an aspiration
  to maximize.
- A reliability target set without input from whoever actually experiences
  the downtime (support team, on-call engineers, or customers directly)
  tends to be miscalibrated — too strict for a target that costs too much
  to hit, or too loose for what users actually tolerate.
- Tradeoff policies decay if nobody owns re-examining them — a target set
  two years ago for a product's usage pattern at the time may no longer
  fit; schedule a periodic review of the threshold itself, not just
  enforcement against it.
- Applying an error-budget-style freeze mechanically without judgment
  (e.g., freezing even a critical security patch because the budget is
  exhausted) misses the point — the policy exists to force a real
  conversation about risk, and a written, visible exception process is
  part of the policy design, not a loophole undermining it.

## Real-world grounding

Google's *Site Reliability Engineering* book (Beyer, Jones, Petoff,
Murphy, O'Reilly, 2016), specifically its chapter on error budgets, is the
publicly documented, widely cited source for this pattern — describing how
Google turned the abstract "reliability vs. velocity" tension into an
operational metric with an automatic consequence (a release freeze) once
the metric crosses a threshold, rather than leaving each release decision
to a fresh negotiation between product and infrastructure teams.

## Verification

- [ ] The tradeoff was named explicitly with two sides, not left implicit
- [ ] A measurable proxy and threshold were agreed while calm, not under
      active release pressure
- [ ] What happens automatically when the threshold is crossed is written
      down, including who can grant a visible, logged exception
- [ ] The threshold itself has an owner and a periodic re-examination,
      not just permanent enforcement
- [ ] The target chosen is deliberately short of the theoretical maximum
      (not "100%" or its equivalent) with a stated reason why
