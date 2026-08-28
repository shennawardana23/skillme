# RICE Worked Example and Application Steps

Load this when you're about to run an actual RICE scoring pass on a real
backlog — a full worked calculation plus the step-by-step process for
applying RICE consistently across rows.

## Worked example

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

## Applying RICE

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
