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
  directly, not a vanity count). Gives the team one number to rally
  around, plus a small set of **input metrics** as actionable levers.
- **AARRR / "Pirate Metrics"** (Dave McClure) — a five-stage funnel
  (Acquisition, Activation, Retention, Referral, Revenue) for diagnosing
  *where* in the user journey a product is leaking. Tells the team where
  to look, not what to rally around.

Used together: AARRR breaks the journey into diagnosable stages; NSM gives
the org one number everyone agrees matters. Input metrics are the
connective tissue — each usually maps to one AARRR stage, which is what
makes "the North Star moved, here's why" traceable instead of a guess.

## Procedure: choosing a North Star Metric

1. **Start from the core value exchange, not what's easy to measure.** Ask:
   "what does the user get, in one sentence, that they'd pay for or miss
   if it vanished?" Airbnb: "a place to stay booked through us" → *nights
   booked*. Spotify: "music that fills my time" → *time spent listening*.
   "Monthly active users" is usually too generic — it captures *presence*,
   not *value delivered*.
2. **Prefer value received over action taken.** "Searches performed" is
   an action; "nights booked" is value received. If a metric can rise
   while users are actually failing — e.g. searches rising because search
   is broken and users keep retrying — it's measuring effort, not value;
   pick something further downstream.
3. **Check it's a leading indicator of revenue, not revenue itself.**
   Revenue and headcount-style metrics are lagging and easy to game short
   term (discounting, one-time promotions). The NSM should predict
   revenue while staying closer to the user experience.
4. **Verify it's one number the whole company can track weekly**, not a
   composite index or a ten-chart dashboard. If engineering, sales, and
   support would each describe "success" using a different metric, keep
   narrowing.
5. **Pressure-test for gameability.** Ask "how would a team hit this
   number in a way that makes the product worse?" (e.g. "time spent
   listening" inflated by removing a skip button). Pair the NSM with a
   guardrail metric (e.g. skip rate, churn) if an obvious bad-faith path
   exists.
6. **Select 2-4 input metrics** that causally drive the NSM and that a
   team can actually act on this quarter — ideally the multiplicative or
   additive components of the NSM itself (e.g. "nights booked" = active
   listings × search-to-book conversion rate × average length of stay).
   Vague inputs ("brand awareness") aren't usable — no team owns a lever
   to move them.
7. **Assign ownership of each input metric to a specific team.** Unowned
   inputs turn the North Star into a metric everyone watches and no one
   is accountable for moving.

## Procedure: diagnosing a problem with AARRR

1. **Place the symptom on the funnel before proposing a fix.** Five
   stages, in order: **Acquisition** (users arrive), **Activation** (good
   first experience / "aha" moment), **Retention** (users return),
   **Referral** (users bring others), **Revenue** (users pay). Vague
   complaints ("growth is stalling," "engagement is down") usually
   describe one specific stage, and the fix differs by stage.
2. **"Signups up, usage flat" is Activation, not Acquisition** —
   Acquisition is clearly working; the leak is between signup and the
   user experiencing real value. Don't respond by spending more on
   acquisition — that widens the top of the funnel while the same
   fraction leaks through the middle.
3. **"Users try it once and don't come back" is Retention**, even if
   Activation looked fine — a good first experience with no reason to
   return is a distinct failure mode from a bad one. Instrument a
   specific "come back by day N" metric, not one "engagement" blob.
4. **"We get users but they never invite anyone" is a Referral gap**,
   often neglected because it has the least existing instrumentation —
   teams frequently have Acquisition and Revenue dashboards but no
   Referral metric, so the problem goes undiagnosed by default.
5. **"Usage is healthy but nobody converts to paid" is Revenue**, not
   Retention just because both are "downstream" — check whether the
   issue is pricing, packaging, or a missing upgrade prompt.
6. **Once the stage is identified, pick the input metric that lives in
   that stage** and confirm the NSM's breakdown actually covers it. If it
   doesn't, that's a sign the NSM or its inputs need revisiting, not that
   the leak doesn't matter.

## Worked example

Both procedures applied end-to-end to a recipe-planning app — NSM
selection, input metrics, and an AARRR diagnosis of "signups up, usage
flat" — is in
[references/worked-example.md](references/worked-example.md).

## Gotchas

- **Picking "Monthly Active Users" as the North Star** is one of the most
  common mistakes — MAU keeps rising from pure acquisition even as value
  delivered per user falls, masking exactly the problem a North Star is
  supposed to surface.
- **A North Star with no guardrail metric** invites a team to optimize it
  in ways that damage the product (e.g. inflating "time spent" with
  addictive-but-low-value engagement loops). Always pair the NSM with a
  metric that would catch this.
- **Treating AARRR stages as strictly sequential** oversimplifies real
  behavior — users can refer others before converting to paid, and
  retention and referral interact. Use the stages as a diagnostic lens
  for "where's the biggest leak," not a rigid pipeline.
- **No Referral metric exists on many teams by default**, which silently
  biases diagnosis toward Acquisition/Revenue — the only stages with
  dashboards. Absence of data at a stage is not evidence it's healthy.
- **Changing the North Star Metric frequently** defeats its purpose of
  giving the team a stable, shared target over quarters. Swapping it
  whenever a dashboard disappoints signals it wasn't chosen carefully.
- **Input metrics chosen because they're easy to move, not because they
  causally drive the NSM**, produce a team that hits its numbers while
  the North Star itself stays flat — check the causal link before
  adopting one.

## Real-world grounding

Both frameworks trace to named sources (Sean Ellis/Amplitude for NSM,
Dave McClure for AARRR) — attribution is in
[references/real-world-grounding.md](references/real-world-grounding.md).

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
