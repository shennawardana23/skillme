---
name: feature-prioritization-frameworks
description: Guides prioritizing a backlog of features or requests using RICE scoring (Reach x Impact x Confidence / Effort) as the primary quantitative method, with MoSCoW (Must/Should/Could/Won't) for stakeholder negotiation and the Kano model for classifying features by satisfaction-curve shape. Use when the user asks to "prioritize these features", "rank this backlog", "score with RICE", "do a MoSCoW", "apply Kano", or gives a list of feature requests and asks which to build first.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Feature Prioritization Frameworks

Prioritization frameworks answer different questions. RICE answers "which
of these comparable features gives the most impact per unit of effort."
MoSCoW answers "which of these committed-scope items can we cut under
time pressure." Kano answers "what shape of satisfaction does this feature
produce, and does building more of it even keep helping." Picking the
wrong one for the situation produces a confident-looking but wrong answer
— pick based on the question being asked, not habit.

## RICE scoring (primary, quantitative)

RICE = **(Reach × Impact × Confidence) / Effort**

Compute each factor before combining them — don't eyeball the final score.

- **Reach**: how many users/customers this affects in a fixed time period
  (e.g., "per quarter"). A count, not a percentage — use actual or
  estimated numbers (e.g., 400 users/month), so reach isn't silently
  double-weighted against impact.
- **Impact**: how much it moves the needle *per user reached*, scored on
  a discrete scale, not a continuum, because false precision here is the
  most common RICE mistake:
  - 3 = massive impact
  - 2 = high impact
  - 1 = medium impact
  - 0.5 = low impact
  - 0.25 = minimal impact
- **Confidence**: how sure you are about the Reach and Impact estimates,
  as a percentage, reflecting evidence quality:
  - 100% = backed by data (analytics, experiment results)
  - 80% = backed by partial data or strong qualitative signal
  - 50% = a guess with some reasoning behind it
  - Below 50% — the estimate is too weak to score; go get more evidence
    or explicitly flag the score as low-confidence in the output, don't
    silently treat it as equal to a data-backed guess.
- **Effort**: total person-time to ship, in a consistent unit (e.g.,
  "person-months"), including design/QA/rollout, not just the coding
  estimate — effort estimates that only count implementation time
  systematically overrate features with a hidden testing or migration
  cost.

### Worked example

Backlog of five features for a hotel booking product, one quarter reach:

| Feature | Reach (users/qtr) | Impact | Confidence | Effort (person-months) | RICE |
|---|---|---|---|---|---|
| A: One-click rebooking | 4,000 | 2 | 80% | 2 | (4000×2×0.8)/2 = **3,200** |
| B: Loyalty tier badges | 8,000 | 0.5 | 100% | 1 | (8000×0.5×1.0)/1 = **4,000** |
| C: AI itinerary chatbot | 1,000 | 3 | 50% | 6 | (1000×3×0.5)/6 = **250** |
| D: Guest review reminders | 6,000 | 1 | 80% | 0.5 | (6000×1×0.8)/0.5 = **9,600** |
| E: Multi-currency pricing | 2,500 | 2 | 50% | 3 | (2500×2×0.5)/3 = **833** |

Ranked by RICE: **D (9,600) > B (4,000) > A (3,200) > E (833) > C (250)**.

Notice C looks exciting narratively ("AI chatbot") but ranks last — high
effort and low confidence overwhelm a high impact score. This is RICE's
main value: it makes an intuitively-appealing but weakly-evidenced,
expensive bet lose to a boring, cheap, well-understood one on paper, and
forces the team to argue about the *inputs* (is confidence really only
50%? is effort really 6 months?) rather than about gut feel.

### Applying RICE

1. Score every feature on the same reach time window and the same effort
   unit — a common error is mixing "reach this month" with "reach this
   quarter" across rows, which silently distorts the ranking.
2. Have the same person or a small group score all rows in one sitting.
   Scoring different features on different days invites scope and
   optimism drift between them.
3. Use the ranked list as a starting point for discussion, not a final
   verdict — RICE doesn't know about dependencies (feature E might be a
   prerequisite for a future feature not yet in the backlog) or strategic
   commitments made to a specific customer.
4. Re-score when new evidence arrives (an experiment result, a sales
   commitment) rather than treating the first score as permanent.

## MoSCoW (lightweight, for stakeholder negotiation)

Sorts already-in-scope items into four buckets for a specific release or
deadline:

- **Must have** — the release is not viable without this; if it slips,
  the release date slips.
- **Should have** — important, painful to cut, but the release survives
  without it.
- **Could have** — desirable, cut first under time pressure, no real pain.
- **Won't have (this time)** — explicitly out of scope for *this* release,
  named so it stops coming up in every planning conversation, not
  forgotten forever.

Use MoSCoW instead of RICE when:
- The conversation is about a fixed deadline or fixed release, and the
  real question is "what do we cut," not "what has the best ROI" — MoSCoW
  has no numerator/denominator, so it can't rank within a tier, only sort
  into tiers.
- You need fast, low-friction stakeholder alignment in a room (a single
  meeting) rather than a defensible numeric artifact — MoSCoW is a
  negotiation tool, not an analytical one.
- Items are not really comparable on reach/impact (e.g., a legal
  compliance requirement vs. a UX polish item) — forcing both into RICE's
  numeric scale produces a meaningless comparison; MoSCoW lets "Must" absorb
  the compliance item without pretending to quantify it against the polish item.

The common failure: letting "Must have" become the default answer for
everything a loud stakeholder wants. Guard it by requiring a stated
consequence for each Must ("if this slips, we cannot launch because
___") — if no one can finish that sentence, it's not a Must.

## Kano model (classification, not scoring)

Kano classifies a feature by the *shape* of the relationship between how
much of it you build and how satisfied customers are — it does not
produce a rank-ordered list, so don't try to force a single Kano-derived
number next to a RICE score.

- **Basic (must-be)**: absence causes strong dissatisfaction, presence
  causes no delight, just the absence of complaint (hot water in a hotel
  room). Investment here has a ceiling — more doesn't make people happier
  past a point.
- **Performance**: satisfaction scales roughly linearly with how well it's
  done (Wi-Fi speed, checkout speed). More is proportionally better.
- **Delight (attractive/excitement)**: absence isn't noticed or missed,
  but presence creates disproportionate satisfaction (an unexpected
  upgrade). These often become tomorrow's Basic features once customers
  come to expect them.
- **Indifferent**: customers don't care either way; building it is
  usually wasted effort regardless of what RICE-style impact estimate was
  assigned to it.
- **Reverse**: some customers are actively less satisfied when the
  feature is present (e.g., a feature that adds a step for a segment that
  wants simplicity) — worth checking for before a broad rollout,
  especially on features aimed at power users vs. casual users.

Use Kano instead of RICE when the real question is "what kind of
investment is this" rather than "which of these comparable items wins" —
e.g., deciding whether a request is a Basic expectation you must meet
regardless of ROI math, or a Delight feature where diminishing returns
mean the fourth iteration isn't worth building even though the first one
tested well. Kano is typically informed by structured customer surveys
(the paired functional/dysfunctional question format from Kano's original
method), which is heavier to run than RICE or MoSCoW — reserve it for
features where the category (basic vs. delight) is genuinely unclear and
worth the research cost, not for routine backlog grooming.

## Choosing between the three

- Comparable features competing for the same engineering time, need a
  defensible numeric ranking → **RICE**.
- Fixed release, fixed deadline, need fast stakeholder agreement on what
  ships and what's cut → **MoSCoW**.
- Need to know whether a feature is table-stakes, a scaling lever, or a
  novelty that will fade → **Kano**.
- They're complementary, not exclusive: Kano can classify a feature as
  Basic before RICE ranks it against other Basics; MoSCoW can turn a
  RICE-ranked shortlist into a committed release scope.

## Gotchas

- RICE scores feel more objective than they are — every input (Impact,
  Confidence) is still a human judgment call. Two teams scoring the same
  backlog honestly can land materially different rankings; the value is
  in forcing the judgment calls to be explicit and comparable, not in
  producing an objectively "correct" number.
- Mixing reach time windows (monthly vs. quarterly) across rows in the
  same RICE table is the single most common silent error — it inflates
  or deflates scores in a way that isn't visible just by looking at the
  final numbers.
- Treating "Must have" in MoSCoW as "everything important" rather than
  "the release literally cannot ship without this" collapses the
  framework back into an undifferentiated priority list.
- Kano categories drift over time — a Delight feature (originally a
  differentiator) frequently decays into a Basic expectation as the market
  catches up (e.g., mobile check-in at hotels). A Kano classification
  done two years ago should not be assumed to still hold.
- Effort estimates in RICE that only cover the "happy path" build (no
  QA, rollout, support docs, or migration time) understate effort and
  systematically bias RICE toward features that look deceptively cheap.

## Real-world grounding

RICE was developed and publicly documented by Intercom's product team
(the framework write-up is commonly attributed to Intercom PM Sean
McBride, published on Intercom's blog) as a way to compare feature
requests on a shared quantitative basis rather than by whoever argued
loudest in the room. MoSCoW originated in the 1990s within the DSDM
(Dynamic Systems Development Method) agile framework as a way to
negotiate fixed-timebox scope with stakeholders. The Kano model comes from
Noriaki Kano's 1984 research on customer satisfaction, which used a
paired-question survey technique (asking how a customer would feel both
with and without a feature) to show that satisfaction is not always
linear in feature investment — some features have ceilings, some have
open-ended payoff, and some don't matter at all.

## Verification

- [ ] Every RICE row uses the same reach time window and same effort unit
- [ ] Impact and Confidence use the discrete scales, not arbitrary decimals
- [ ] Effort includes design/QA/rollout time, not just implementation
- [ ] MoSCoW "Must have" items each have a stated ship-blocking consequence
- [ ] Kano classification isn't assumed permanent for features exposed to a changing market
- [ ] The chosen framework matches the actual question being asked (rank vs. cut vs. classify)
