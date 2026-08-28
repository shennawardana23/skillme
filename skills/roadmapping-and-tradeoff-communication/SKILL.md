---
name: roadmapping-and-tradeoff-communication
description: Guides building Now/Next/Later, outcome-based roadmaps and communicating tradeoff decisions to stakeholders who push for a specific feature with a hard date. Use when the user asks to "build a roadmap," "create a Now/Next/Later roadmap," "turn this feature list into an outcome roadmap," "respond to a stakeholder demanding a date," "push back on a feature request," "explain why we can't commit to this date," or is reviewing a roadmap that is really just a dated feature list.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Roadmapping and Tradeoff Communication

Build roadmaps around outcomes and horizons, not features and dates.
Communicate tradeoffs by naming the mechanism that would have to change,
not by softening a refusal.

## Why date-committed feature roadmaps backfire

A roadmap that lists "Feature X — ships March 15" makes a promise the team
usually cannot keep, for reasons that have nothing to do with execution
quality: estimates made months out are wrong by construction (unknown
unknowns compound over time), priorities shift as the market or the data
changes, and the roadmap gets treated as a contract the moment a
stakeholder builds their own plan on top of that date. When the date
slips — and on any roadmap spanning more than a few weeks, some date will
slip — the failure mode isn't the delay itself, it's that stakeholders
correctly learn the roadmap's dates were never reliable and stop trusting
the *next* roadmap too. The trust damage compounds across cycles, not
just within one.

The deeper problem: a date-committed feature list states false certainty
about two things at once — that this specific feature is the right
solution (not yet validated), and that it will take exactly this long
(not yet known). Bundling an unvalidated solution with a precise date
manufactures confidence the underlying work doesn't support.

## Now / Next / Later

Popularized by Janna Bastow (co-founder of ProdPad) as a response to the
date-driven Gantt roadmap, Now/Next/Later organizes roadmap items into
three horizons instead of a calendar:

- **Now** — actively being worked on. Specific enough to describe concrete
  scope; confidence is high because it's in progress or fully scoped.
- **Next** — validated as a priority and coming after Now, but not yet
  scoped in detail. Sequencing is fairly firm; timing is not.
- **Later** — directionally important, on the radar, but not yet
  validated or prioritized against everything else that could land there.
  Genuinely likely to change.

The horizon itself communicates confidence implicitly, without ever
writing a date — an item in "Later" tells the stakeholder "this is a
direction we believe in, not a commitment" without saying that sentence
out loud every time someone asks. This is the format's real mechanism: it
replaces an explicit date (which reads as a promise) with a positional
signal (which reads as a confidence level) that degrades gracefully
instead of breaking trust when it moves.

Practical rules for using it honestly:
- Moving an item from Later to Next should require an actual event
  (validated demand, dependency cleared, capacity freed) — not just time
  passing. If Later items age into Next purely because a quarter ended,
  the columns are secretly dates with extra steps.
- Don't let "Now" become a dumping ground of everything in flight with no
  ordering — Now items should still be sequenced by priority within the
  column.
- If a stakeholder asks "so when is Next?", answer with the dependency or
  condition that moves it, not a hedge-date ("probably Q3-ish") — a
  hedge-date reintroduces the exact commitment the format exists to avoid.

## Outcome-based roadmap items vs. feature-list items

A feature-list roadmap item is a solution someone already picked:
"Add bulk CSV export." It commits the team to that solution before
checking whether it solves anyone's actual problem, and it gives the team
no room to discover a better solution once they dig in. This is the core
of Marty Cagan's critique of the "feature factory" — teams ship a stream
of committed features and call it progress, without any mechanism forcing
a check on whether those features moved a real business or user outcome.

An outcome-based item names the problem or metric to move — with a
baseline, if one exists — and leaves the solution open until the team is
actually working on it. When a stakeholder requests a specific feature,
capture the problem behind their request as the roadmap item and keep
their proposed feature as one candidate solution under it, not as the
item itself.

Worked feature-list-to-outcome examples and the full writing-rule
checklist: [references/outcome-vs-feature-examples.md](references/outcome-vs-feature-examples.md)

## Responding to a stakeholder demanding a specific feature by a hard date

Don't respond with a flat no or with vague reassurance ("we'll definitely
look into it"). Instead: name the problem behind the request and confirm
you understood it; show where that problem sits on the roadmap using
horizon language rather than a flat rejection; state the specific
tradeoff a commitment would require — what would actually get displaced;
hand the reprioritization decision to whoever has the authority to make
it, rather than absorbing it unilaterally or caving silently; and never
promise a date under pressure just to end the conversation — an
undefensible date offered in a hallway becomes the next broken promise.

Full five-step script with example phrasing: [references/stakeholder-tradeoff-conversation.md](references/stakeholder-tradeoff-conversation.md)

## Gotchas

- **A Now/Next/Later board with no visible movement is worse than a
  dated roadmap.** If items sit in Next for two quarters unexplained,
  stakeholders stop trusting horizons the same way they'd stop trusting
  slipped dates — the format only holds trust if items visibly move for
  stated reasons.
- **"Later" quietly becomes "never," and stakeholders notice the pattern**
  even unstated. If something has sat in Later for multiple planning
  cycles, say so explicitly and explain why rather than letting it
  silently age out — "we're deprioritizing this, here's why" preserves
  more trust than a Later item that's actually dead but not labeled that
  way.
- **Outcome framing can dodge accountability if there's no baseline or
  target.** "Improve onboarding" with no number and no deadline for
  revisiting it is an excuse not to commit to anything, not an
  outcome-based item. Every outcome item needs a way to eventually check
  whether it worked, even without a fixed date.
- **Reframing a feature request as "the problem behind it" can read as
  dismissive if done poorly** — skipping straight to "here's the real
  problem you actually have" without confirming your read first tells the
  stakeholder their request was overridden, not understood. Always
  confirm the reframed problem with them first.
- **Executives and sales teams often need a specific commitment for
  contractual or external reasons** (a customer contract, a board
  deadline) a horizon can't satisfy. Don't force Now/Next/Later onto a
  genuinely date-bound external commitment — say so explicitly, treat it
  as a separate committed-date item, and be honest it's an exception
  rather than quietly running two incompatible roadmap systems.

## Real-world grounding

Now/Next/Later was popularized by Janna Bastow, co-founder of the
roadmapping tool ProdPad, as an explicit alternative to Gantt-chart-style
roadmaps with fixed ship dates — she has written and spoken widely on how
date-based roadmaps create false certainty and recurring stakeholder
distrust when dates slip. The outcome-based critique of feature-list
roadmaps is closely associated with Marty Cagan (SVPG), whose "outcomes
over outputs" argument describes teams that ship a steady stream of
committed features — a "feature factory" — without ever validating that
those features moved any real outcome, and argues roadmaps should commit
to problems worth solving rather than to specific solutions and dates.

## Verification

- [ ] No roadmap item promises a specific ship date more than one horizon out
- [ ] Every roadmap item is phrased as a problem/outcome, not a UI feature
- [ ] Outcome items have a stated baseline metric, even if rough
- [ ] Items that moved horizons this cycle have a stated reason for the move
- [ ] A stakeholder's feature request was translated into an underlying
      problem and confirmed with them before being placed on the roadmap
- [ ] A tradeoff response named the specific thing that would be
      displaced, not just "we don't have capacity"
