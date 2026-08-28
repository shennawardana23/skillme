# Worked Example: NSM + AARRR Applied Together

Load this when you want to see the full procedure applied end-to-end to a
concrete product, not needed when you already know how to apply the
procedures to your own product.

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
