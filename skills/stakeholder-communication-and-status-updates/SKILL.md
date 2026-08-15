---
name: stakeholder-communication-and-status-updates
description: Guides writing narrative-style stakeholder status updates and strategy memos instead of bullet-point decks, grounded in Amazon's 6-page narrative memo discipline. Use when the user asks to "write a status update," "draft a weekly update for stakeholders," "write a project update email," "prepare a 6-pager," "turn these bullets into a narrative," "write a stakeholder memo," or is reviewing a status update/deck that seems to be hiding risk or fuzzy thinking behind bullet points.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Stakeholder Communication and Status Updates

Write status updates and strategy memos as connected prose, not bullet
fragments. Full sentences force you to state the causal or logical
relationship between two facts; a bullet point lets you juxtapose them and
imply a connection without ever having to defend it.

## Why narrative beats bullets

A bullet list is a set of claims with no stated relationship between them.
Two adjacent bullets —

```
- Signups grew 12% this month
- We shipped the new onboarding flow
```

— invite the reader to assume the second caused the first. Maybe it did.
Maybe signups grew because of a marketing campaign and the onboarding flow
actually hurt conversion. The bullet format never forces the writer to make
the claim explicit, so it never gets checked. A sentence forces the choice:

> "Signups grew 12% this month, driven primarily by the referral campaign;
> the new onboarding flow shipped mid-month and it's too early to isolate
> its effect on conversion."

That sentence is falsifiable — a reader can push back on "driven primarily
by." The bullet pair was not falsifiable; it just sat there implying
causation. This is the core mechanism, not a stylistic preference: writing
prose exposes gaps in reasoning that the author would otherwise never
notice, because a bullet never demands the connective tissue that a
sentence grammatically requires.

The same mechanism catches fuzzy thinking at the whole-document level.
Writing "we should prioritize X" as a sentence with a "because" clause
forces you to actually have a reason. A slide titled "Prioritization" with
X, Y, Z as sub-bullets underneath implies X, Y, Z were compared and X won —
without ever showing the comparison.

## Structure of a status/strategy narrative

Use this order for anything longer than a quick weekly note (a strategy
memo, a quarterly review, a project deep-dive):

1. **Context** — what this document is about and why the reader is
   reading it now. One paragraph. Assume the reader has forgotten details
   from last time; don't assume they remember your last update.
2. **What happened** — the events and decisions since the last update,
   told as a sequence, not a status board. Include what you decided NOT to
   do and why, if it's relevant — omitted alternatives are exactly the kind
   of thing a bullet list quietly drops.
3. **Data** — the numbers that support or complicate the "what happened"
   narrative, with enough surrounding sentence to say what the number
   means and whether it's good, bad, or ambiguous. A number without an
   interpretation is not information; state whether it's above, below, or
   in line with expectation.
4. **What's next** — concrete next steps, owners, and rough timing,
   written as commitments, not aspirations ("we will" not "we hope to").
5. **Risks** — see below. This section exists even when there's nothing
   alarming to report; "no material risks this period" is itself a claim
   worth stating explicitly rather than omitting.

Order matters: risks go near the end structurally but must never be the
thing that gets cut when the document runs long or when the writer is
rushed. See the Gotchas section for how risk sections quietly disappear.

## Writing a narrative status update, not a 6-pager

A full 6-page narrative is for a decision meeting or a quarterly deep dive
— not for a routine weekly or biweekly stakeholder check-in. Scale the same
discipline down instead of switching to bullets:

- Keep all five sections above, but compress each to 1–3 sentences instead
  of a paragraph. A weekly update can be 10–15 sentences total and still be
  narrative.
- Do not compress by turning sentences back into bullets. Compress by
  cutting detail, not by cutting the sentence structure that forces
  reasoning to be explicit.
- Skip "Context" entirely only if the update goes to the same recurring
  audience and nothing about the project's purpose has changed — otherwise
  keep one sentence of it.
- If a section has nothing to report ("Risks: none new this period"), say
  so in one sentence rather than deleting the section — a missing section
  is indistinguishable from a forgotten one, to the reader.

## Surfacing risk early instead of only good news

Status updates drift toward good-news-only reporting because bad news is
uncomfortable to write and because good news is what got done, so it's
naturally what comes to mind first while writing. Counter this
structurally:

1. **Write risks before writing accomplishments.** Draft the risk section
   first, even if it appears later in the final document. If you write
   accomplishments first, the risk section becomes an afterthought
   squeezed in at the end, and afterthoughts get watered down.
2. **State risk as a specific, falsifiable claim with a timeframe**, not a
   mood. "There's some concern about the timeline" is not a risk
   statement; a reader can't act on it or check it later. "If the vendor
   API integration isn't confirmed working by the 14th, launch slips past
   the end of the month" is a risk statement — it names the trigger, the
   consequence, and the date.
3. **Separate "risk" from "issue already causing harm."** A risk is a
   thing that might happen; an issue is a thing that is happening. Burying
   an active issue inside a "risks" section labeled as hypothetical is a
   common way status updates soften bad news without technically lying.
4. **Report risk trend, not just risk existence.** "This risk is new this
   week" and "this risk has been open for three weeks and is not shrinking"
   are different signals; a recurring risk section that never changes its
   wording is often masking stalled progress on mitigation.

## Gotchas

- **The risk section is the first thing cut under time pressure**, and the
  cutting is invisible to the reader — they don't know a risk section
  existed and was removed, they just see a clean update. Treat "did I
  write the risk section, and does it name a specific trigger and date"
  as a non-negotiable check before sending, not an optional nicety.
- **A slide deck's bullet hierarchy substitutes for actual argument
  structure.** Indenting sub-bullets under a bold header implies "these
  support the header" without ever stating how. When converting a deck to
  a narrative, don't just prose-ify each bullet in place — ask what
  relationship was implied by the indentation and state it explicitly;
  often there wasn't one, which is itself the finding.
- **Silent reading only works if the memo is actually finishable in the
  time allotted.** Amazon's practice allocates real silent-reading time in
  the meeting (commonly 15–20 minutes for a 6-pager) before any
  discussion starts. A narrative memo assigned as "pre-read homework" that
  people are expected to have skimmed beforehand reintroduces the exact
  problem — a partially-read document — that silent in-room reading was
  designed to prevent.
- **"Narrative" does not mean "no data."** A memo that replaces every
  number with adjectives ("strong growth," "some churn") is prose-shaped
  bullet-point vagueness — it has sentences but still lacks the falsifiable
  claims that make prose useful. Every quantitative claim should have an
  actual number and a stated comparison point (vs. last period, vs. target).
- **Padding a narrative to hit a page target defeats the purpose.**
  Amazon's discipline works because writing tightly-argued prose is harder
  than bullets, not because 6 pages is a magic length — a rambling 6-pager
  with the same fuzzy thinking bullets would have hidden is worse than a
  tight 2-pager. Compress by cutting, never by padding to match a template.
- **Don't ban all visuals reflexively.** A chart showing an actual trend
  line is data, not a bullet — it's fine and often clearer than a sentence
  describing the same trend. What the discipline eliminates is prose being
  replaced by unsupported fragments, not all visual aids categorically.

## Real-world grounding

Amazon has, since roughly the mid-2000s, run S-team and other planning
meetings by having attendees silently read a narrative memo — famously
capped at six pages, with no slide deck — for the first part of the
meeting before any discussion begins. Jeff Bezos has described the
reasoning in shareholder letters and internal communications reported by
multiple former Amazon executives: a full-sentence, full-paragraph memo
forces the author to actually work out cause-and-effect and priority,
because prose has nowhere to hide a non-sequitur the way a bullet list
does. Silent reading in the room (rather than a pre-read) exists because
people skim pre-reads unevenly — some read closely, some skim the night
before, some don't open it at all — so the meeting starts with an uneven
base of understanding; reading together for a fixed block guarantees
everyone enters the discussion from the same material.

## Verification

- [ ] Every claim connecting two facts is stated as a sentence with the
      relationship explicit (causes, correlates with, contradicts) — not
      left as adjacent bullets
- [ ] Risk section was drafted before the accomplishments section
- [ ] Every risk names a specific trigger condition and a date, not just a mood
- [ ] Active issues are labeled as issues, not folded into "risks"
- [ ] Every number has a stated comparison point (vs. last period, vs. target)
- [ ] The update was compressed by cutting detail, not by reverting to bullets
