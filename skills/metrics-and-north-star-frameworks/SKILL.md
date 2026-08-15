---
name: metrics-and-north-star-frameworks
description: This skill should be used when the user asks to define a North Star Metric, choose input/driver metrics for a product, diagnose a funnel using AARRR or "Pirate Metrics" (Acquisition, Activation, Retention, Referral, Revenue), or asks things like "what metric should we rally the team around", "where in the funnel are we losing users", "signups are up but usage isn't, what does that mean", or "help me build a metrics tree for this product". Combines the North Star Metric framework with AARRR to connect one company-wide metric to funnel-stage diagnosis.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Metrics and North Star Frameworks

Two named frameworks solve two different problems, and using only one
leaves a gap:

- **North Star Metric (NSM)** — the single metric that best captures the
  core value a product delivers to customers *right now* (not revenue
  directly, and not a vanity count). It gives the whole team one number to
  rally around, with a small set of **input metrics** as the actionable
  levers that move it.
- **AARRR / "Pirate Metrics"** (Dave McClure) — a five-stage funnel
  (Acquisition, Activation, Retention, Referral, Revenue) for diagnosing
  *where* in the user journey a product is leaking. It doesn't tell the
  team what to rally around; it tells the team where to look when
  something's wrong.

Used together: AARRR breaks the user journey into stages you can instrument
and diagnose. NSM gives the team one number the whole org agrees matters.
Input metrics are the connective tissue — each input metric to the NSM
usually maps to one AARRR stage, which is what makes "the North Star moved,
here's why" a traceable statement instead of a guess.

## Procedure: choosing a North Star Metric

1. **Start from the core value exchange, not from what's easy to measure.**
   Ask: "what does the user get, in one sentence, that they'd pay for or
   miss if it vanished?" Airbnb's answer was "a place to stay booked
   through us" → *nights booked*. Spotify's was "music that fills my time"
   → *time spent listening*. A metric like "monthly active users" is
   usually too generic to be a real North Star because it doesn't capture
   *value delivered*, only *presence*.
2. **Prefer a metric that reflects value received, not just an action
   taken.** "Searches performed" is an action; "nights booked" is value
   received (a completed transaction implying the search worked). If the
   candidate metric can go up while users are actually failing (e.g.
   "searches performed" rising because search is broken and users keep
   retrying), it's measuring effort, not value — pick a metric further
   downstream.
3. **Check it's a leading indicator of revenue, not revenue itself.**
   Revenue and headcount-style metrics are lagging and easy to game short
   term (discounting, one-time promotions) without fixing the product.
   The NSM should predict revenue over time while staying closer to the
   user experience, so the team can act on it before revenue moves.
4. **Verify it's a single number the whole company can understand and
   track weekly**, not a composite index or dashboard of ten charts. If
   engineering, sales, and support would each describe "success" using a
   different metric, it isn't a North Star yet — keep narrowing.
5. **Pressure-test for gameability.** Ask "how would a team hit this
   number in a way that makes the product worse?" (e.g. "time spent
   listening" could be inflated by removing a skip button). If an obvious
   bad-faith path exists, pair the NSM with a guardrail metric (e.g.
   skip rate, churn) that the team also watches so the NSM can't be gamed
   in isolation.
6. **Select 2-4 input metrics** — the metrics that causally drive the
   NSM and that a team can actually act on this quarter. Good input
   metrics are the multiplicative or additive components of the NSM
   (e.g. for "nights booked": number of active listings × search-to-book
   conversion rate × average length of stay). Vague inputs ("brand
   awareness") aren't usable because no team owns a lever to move them
   directly.
7. **Assign ownership of each input metric to a specific team.** A North
   Star with no owned inputs turns into a metric everyone watches and no
   one is accountable for moving.

## Procedure: diagnosing a problem with AARRR

1. **Place the reported symptom on the funnel before proposing a fix.**
   The five stages, in order: **Acquisition** (users arrive), **Activation**
   (users have a good first experience / reach an "aha" moment),
   **Retention** (users come back), **Referral** (users bring others),
   **Revenue** (users pay). Most vague complaints ("growth is
   stalling," "engagement is down") actually describe one specific stage,
   and the fix is different for each.
2. **"Signups are up but usage isn't" is an Activation problem, not an
   Acquisition problem** — Acquisition (getting people to sign up) is
   clearly working; the leak is between signup and the user experiencing
   real value. Don't respond to this symptom by spending more on
   acquisition channels — it makes the funnel wider at the top while the
   same fraction leaks through the middle.
3. **"Users try it once and don't come back" is a Retention problem** even
   if Activation (their first session) looked fine — a good first
   experience with no reason to return is a distinct failure mode from a
   bad first experience. Instrument a specific "come back by day N" metric
   rather than treating "engagement" as one blob.
4. **"We get users but they never invite anyone" is a Referral gap**, and
   is often neglected because it's the stage with the least existing
   instrumentation — teams frequently have Acquisition and Revenue
   dashboards but no Referral metric at all, which means Referral problems
   go undiagnosed by default, not because they don't exist.
5. **"Usage is healthy but nobody converts to paid" is a Revenue-stage
   problem**, and should not be treated as a Retention problem just
   because both stages are "downstream" — check whether the issue is
   pricing, packaging, or a missing upgrade prompt, not engagement.
6. **Once the stage is identified, pick the input metric that lives in
   that stage** and confirm the NSM's input-metric breakdown actually has
   a metric there. If the NSM's inputs don't cover the stage where the
   real leak is, that's a sign the NSM or its inputs need to be revisited,
   not that the leak doesn't matter.

## Worked example

**Product**: a recipe-and-meal-planning app that generates a weekly grocery
list from saved recipes.

- **North Star Metric**: *meal plans completed per active user per week*
  (a "completed" plan = recipes selected + grocery list generated). This
  reflects the core value exchange (turning recipe browsing into an
  actual, actionable plan) rather than a proxy like "recipes viewed,"
  which could rise even if nobody ever finishes a plan.
- **Input metrics**:
  1. *% of new users who complete their first meal plan within 7 days*
     (Activation stage) — the biggest lever on whether a user ever
     experiences the core value at all.
  2. *Average number of recipes saved per user per week* (feeds Retention
     — a user with a growing recipe library has more reason to return).
  3. *% of completed plans that generate a grocery list* (a friction
     metric inside the core loop itself, closest to the NSM).
- **Mapping a specific problem**: leadership reports "signups grew 40%
  this quarter, but weekly active users barely moved." Using AARRR: signup
  growth confirms Acquisition is healthy. Flat WAU despite growing signups
  points at Activation — new users aren't reaching "first completed meal
  plan." The relevant input metric is #1 above; the fix is in onboarding
  (e.g. prompting a first plan during signup), not in acquisition spend or
  in Referral/Revenue features.

## Gotchas

- **Picking "Monthly Active Users" as the North Star** is one of the most
  common mistakes — MAU is a vanity metric that keeps rising from pure
  acquisition even while the actual value delivered per user is falling,
  so it can mask exactly the problem a North Star is supposed to surface.
- **A North Star with no guardrail metric** invites a team to optimize the
  number in a way that damages the product (e.g. maximizing "time spent"
  by adding addictive-but-low-value engagement loops). Always pair the
  NSM with at least one metric that would catch this.
- **Treating AARRR stages as strictly sequential for every user** oversimplifies
  real behavior — some users refer others before converting to paid, and
  retention and referral interact. Use the stages as a diagnostic lens for
  "where's the biggest leak," not as a rigid pipeline every single user
  must follow in order.
- **No metric exists for the Referral stage** on many teams by default,
  which silently biases diagnosis toward Acquisition/Revenue problems
  because those are the only stages with dashboards. Absence of data at a
  stage is not evidence the stage is healthy.
- **Changing the North Star Metric frequently** defeats its purpose — it
  exists to give the team a stable, shared target over quarters. Swapping
  it every time a dashboard looks disappointing is a sign the metric
  wasn't chosen carefully in the first place, not a reason to keep
  swapping.
- **Input metrics chosen because they're easy to move, not because they
  causally drive the NSM**, produce a team that hits its numbers while the
  North Star itself stays flat — always check the causal link (ideally
  via the NSM's own formula components) before adopting an input metric.

## Real-world grounding

The North Star Metric framework was popularized by growth expert Sean
Ellis and later formalized into Amplitude's widely-read "North Star
Playbook," which documents the now-standard examples of Airbnb's *nights
booked* and Spotify's *time spent listening* as metrics chosen because
they represent the core value exchange between product and customer,
supported by a small set of input metrics that teams can directly act on.
AARRR, commonly called "Pirate Metrics" because of its acronym, was
created by startup investor and advisor Dave McClure as a simple, ordered
funnel — Acquisition, Activation, Retention, Referral, Revenue — for
diagnosing where in the customer lifecycle a startup is losing users, and
remains a standard framework taught across startup accelerators and growth
teams.

## Verification

- [ ] The North Star Metric reflects value delivered to the user, not
      just an action taken or a vanity count
- [ ] The North Star Metric has a paired guardrail metric to prevent
      gaming
- [ ] 2-4 input metrics are identified, each causally linked to the North
      Star and owned by a specific team
- [ ] Every reported symptom is placed on a specific AARRR stage before a
      fix is proposed
- [ ] Each AARRR stage has at least one instrumented metric, including
      Referral, which is the stage most often left unmeasured
