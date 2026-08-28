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

### Worked example and application steps

A full worked RICE calculation (a five-feature hotel-booking backlog,
computed scores, and why a narratively-exciting feature can still rank
last) plus the step-by-step process for scoring a real backlog
consistently are in
[references/rice-worked-example.md](references/rice-worked-example.md)
— load it when you're about to run an actual scoring pass, not for a
quick refresher on what the factors mean.

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

Five categories, briefly: **Basic** (absence causes dissatisfaction,
presence isn't noticed — has a satisfaction ceiling, e.g. hot water in a
hotel room), **Performance** (satisfaction scales roughly linearly with
execution quality, e.g. Wi-Fi speed), **Delight** (absence isn't missed
but presence disproportionately pleases — tends to decay into a Basic
expectation over time), **Indifferent** (customers don't care either
way; building it is usually wasted effort regardless of what impact
score it might have gotten in RICE), and **Reverse** (some customers are
actively less satisfied when the feature is present — worth checking
before a broad rollout). Full definitions, examples, and the survey
method used to classify a feature are in
[references/kano-categories-and-framework-origins.md](references/kano-categories-and-framework-origins.md).

Use Kano instead of RICE when the real question is "what kind of
investment is this" rather than "which of these comparable items wins" —
e.g., deciding whether a request is a Basic expectation you must meet
regardless of ROI math, or a Delight feature where diminishing returns
mean the fourth iteration isn't worth building even though the first one
tested well. Classifying a feature this way is heavier than scoring it
with RICE or sorting it with MoSCoW — reserve it for features where the
category is genuinely unclear and worth the research cost, not for
routine backlog grooming.

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

RICE comes from Intercom's product team, MoSCoW from 1990s DSDM agile
practice, and Kano from Noriaki Kano's 1984 customer-satisfaction
research. Full attribution and the reasoning behind each origin are in
[references/kano-categories-and-framework-origins.md](references/kano-categories-and-framework-origins.md).

## Verification

- [ ] Every RICE row uses the same reach time window and same effort unit
- [ ] Impact and Confidence use the discrete scales, not arbitrary decimals
- [ ] Effort includes design/QA/rollout time, not just implementation
- [ ] MoSCoW "Must have" items each have a stated ship-blocking consequence
- [ ] Kano classification isn't assumed permanent for features exposed to a changing market
- [ ] The chosen framework matches the actual question being asked (rank vs. cut vs. classify)
