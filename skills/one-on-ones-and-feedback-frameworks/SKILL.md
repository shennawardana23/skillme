---
name: one-on-ones-and-feedback-frameworks
description: This skill should be used when the user asks to "give feedback using SBI", "structure a 1:1", "how do I tell my report about this issue", "write feedback for a performance review", or needs to deliver behavioral feedback (positive or corrective) or structure a recurring one-on-one meeting. Use for the structure and delivery of feedback and 1:1s — not for formal performance-review calibration processes.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# One-on-Ones and Feedback Frameworks

Feedback that names a trait ("you're not detail-oriented") invites
defensiveness because it's a character judgment with no clear next step.
Feedback anchored to a specific, observable moment gives the person
something they can actually recognize and change.

## The SBI model

Situation-Behavior-Impact, developed by the Center for Creative Leadership,
structures feedback in three concrete parts:

1. **Situation** — the specific time and place, not a generalization.
   "In yesterday's incident review" not "you always."
2. **Behavior** — what was observed, described factually, not
   interpreted. "You interrupted the on-call engineer twice while they
   were explaining the timeline" not "you were dismissive."
3. **Impact** — the effect it had, on you, the team, or the business.
   "It made it harder for the team to get the full timeline before we
   moved to root-causing" not "it was unprofessional."

Full example: "In yesterday's incident review (Situation), you interrupted
the on-call engineer twice while they were walking through the timeline
(Behavior), which meant we had to backtrack twice to get details we'd
already skipped past (Impact)."

SBI works the same way for positive feedback — naming the specific
situation and behavior makes praise land as something to repeat, not just
a vague compliment: "In yesterday's postmortem (Situation), you asked
'what would have caught this sooner' before anyone else did (Behavior),
which shifted the whole discussion from blame to prevention (Impact)."

## Structuring a 1:1

- **The report owns the agenda**, not the manager — a 1:1 that's the
  manager running through their own status-check list isn't a 1:1, it's a
  status meeting with one person.
  - **Cadence**: recurring, protected time (weekly or biweekly, 30
  minutes is typical) — cancelling repeatedly signals it's not actually a
  priority, which reports notice.
- **Mix content types**: blockers/support needed, growth and career
  conversation, and feedback in both directions — not 100% status
  updates, which async docs/standups already cover more efficiently.
- **Deliver feedback close to the event**, in the 1:1 or as close to it as
  reasonable — feedback on something from three weeks ago has lost the
  specific Situation detail that makes SBI work.

## Gotchas

- **SBI without a real Situation collapses back into character
  judgment.** "You're not a team player, impact: it hurts morale" isn't
  SBI — there's no specific, checkable Situation or Behavior underneath
  it, so it's a label wearing the SBI format without doing SBI's actual
  job.
- **A 1:1 that's entirely status updates is a wasted 1:1** — if every
  agenda item could have been an async message, the recurring meeting slot
  is being spent on the wrong content.
- **Positive feedback deserves the same specificity as corrective
  feedback** — vague praise ("great job this week") is as unhelpful as
  vague criticism, because it doesn't tell the person what to repeat.
- **Public praise and private correction use the same SBI structure but
  different venues** — delivering corrective SBI feedback in a group
  setting undermines the model's intent even if the words are right.
- **Don't skip the Impact.** Situation + Behavior without stating the
  impact leaves the person to guess why it mattered, which is often where
  the actual disagreement or misunderstanding lives.

## Real-world grounding

The Situation-Behavior-Impact model was developed by the Center for
Creative Leadership and is one of the most widely taught feedback
frameworks in management training globally, precisely because it gives
managers (including first-time managers with no formal training) a
repeatable structure that avoids the trait-labeling failure mode common in
unstructured feedback.
