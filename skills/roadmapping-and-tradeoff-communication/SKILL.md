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
slip — the roadmap's actual failure mode isn't the delay itself, it's that
stakeholders correctly learn the roadmap's dates were never reliable, and
stop trusting the *next* roadmap too. The trust damage compounds across
cycles, not just within one.

The deeper problem is that a date-committed feature list states false
certainty about two different things at once: that this specific feature
is the right solution (not yet validated), and that it will take exactly
this long (not yet known, especially for anything not already in
detailed design). Bundling an unvalidated solution with a precise date
manufactures a confidence level the underlying work doesn't support.

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

The horizon itself communicates confidence — implicitly, without ever
writing a date. An item in "Later" tells the stakeholder "this is a
direction we believe in, not a commitment," without the awkwardness of
saying that sentence out loud every time someone asks. This is the
format's real mechanism: it replaces an explicit date (which reads as a
promise) with a positional signal (which reads as a confidence level) —
stakeholders learn to read horizon-as-confidence the same way they'd read
a date, but the signal degrades gracefully instead of breaking trust when
it moves.

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

An outcome-based item names the problem or metric to move, and leaves the
solution open until the team is actually working on it:

| Feature-list item | Outcome-based item |
|---|---|
| "Add bulk CSV export" | "Reduce time finance spends reconciling monthly bookings (currently ~3 hrs/month) — export is one candidate solution" |
| "Build a notifications center" | "Reduce missed-action rate on pending approvals" |
| "Redesign onboarding flow" | "Increase week-1 activation rate for new hotel-admin accounts" |

Writing rules for outcome items:
- Name the metric or problem, not the UI. If the item names a screen or a
  button, it's a feature item wearing an outcome label.
- Include the current baseline if you have it, even roughly — "reduce
  reconciliation time" without a baseline gives nobody a way to know if a
  shipped solution actually worked.
- It's fine, even expected, for the eventual shipped solution to differ
  from whatever solution prompted the roadmap conversation — that's the
  point of not pre-committing to the feature.
- Stakeholders who request a specific feature are usually pointing at a
  real problem through their preferred solution. Capture the problem
  behind their request as the roadmap item; keep their proposed feature as
  one candidate solution under it, not as the item itself.

## Responding to a stakeholder demanding a specific feature by a hard date

Don't respond with a flat no or with vague reassurance ("we'll definitely
look into it"). Use this sequence:

1. **Name the problem behind the request**, out loud, to confirm you
   understood it correctly: "It sounds like the goal is to stop
   double-booking suites during peak season — is that the core issue, or
   is there something more specific?" This does two things: it shows
   you're taking the concern seriously, and it opens the door to a
   different solution than the one they asked for.
2. **Show them where the underlying problem sits on the roadmap**, using
   the horizon language, not a flat rejection: "Reducing double-bookings
   is on our Next horizon — it's a validated priority, we just haven't
   scoped the specific approach yet." If the problem isn't on the roadmap
   at all, say that plainly rather than implying it secretly is.
3. **State the actual tradeoff being asked for, specifically** — name what
   would have to move: "Committing to ship this by the 15th means pulling
   the two engineers currently on [Now-horizon item], which would push
   that item out by roughly the same amount. Is that the tradeoff you want
   me to make?" This reframes "no" as a resourcing decision the stakeholder
   can see the cost of, rather than a personal refusal.
4. **Offer the decision, don't just absorb it.** If the stakeholder has
   authority to reprioritize, hand them the choice explicitly rather than
   deciding unilaterally and rather than caving silently: "I can make this
   the priority if you want — I want to be clear about what it displaces
   first." If they don't have that authority, say who does.
5. **Never promise a date to make the immediate conversation easier.** A
   date offered under pressure in a hallway conversation is exactly the
   kind of commitment that turns into a broken promise later — if you
   genuinely don't know, say "I don't have a defensible date for this yet"
   rather than producing one to end the conversation.

## Gotchas

- **A Now/Next/Later board with no visible movement is worse than a
  dated roadmap.** If items sit in Next for two quarters with no
  explanation, stakeholders stop trusting the horizons the same way they'd
  stop trusting slipped dates — the format only holds trust if items
  visibly move for stated reasons.
- **"Later" quietly becomes "never," and stakeholders notice the pattern**
  even if no one states it. If something has sat in Later for multiple
  planning cycles, say so explicitly and explain why, rather than letting
  it silently age out — an honest "we're deprioritizing this, here's why"
  preserves more trust than a Later item that's actually dead but not
  labeled that way.
- **Outcome framing can be used to dodge accountability if there's no
  baseline or target.** "Improve onboarding" with no number and no
  deadline for revisiting it is not an outcome-based roadmap item, it's an
  excuse not to commit to anything. Every outcome item needs a way to
  eventually check whether it worked, even if that check has no fixed date.
- **Translating a stakeholder's feature request into "the problem behind
  it" can come across as dismissive if done poorly** — skipping straight
  to "here's the real problem you actually have" without confirming your
  read first tells the stakeholder their request was overridden, not
  understood. Always confirm the reframed problem with them before using
  it to justify a different solution.
- **Executives and sales teams often need a specific commitment for
  contractual or external reasons** (a customer contract, a board
  deadline) that a horizon can't satisfy. Don't force Now/Next/Later onto
  a genuinely date-bound external commitment — say so explicitly, treat it
  as a separate committed-date item, and be honest that it's an exception
  rather than quietly running two incompatible roadmap systems.

## Real-world grounding

Now/Next/Later was popularized by Janna Bastow, co-founder of the
roadmapping tool ProdPad, as an explicit alternative to Gantt-chart-style
roadmaps with fixed ship dates — she has written and spoken widely on how
date-based roadmaps create false certainty and become a recurring source
of stakeholder distrust when dates inevitably slip. The outcome-based
critique of feature-list roadmaps is closely associated with Marty Cagan
(SVPG), whose "outcomes over outputs" argument describes teams that ship a
steady stream of committed features — a "feature factory" — without ever
validating that those features moved any real user or business outcome,
and argues that roadmaps should commit to problems worth solving rather
than to specific solutions and dates.

## Verification

- [ ] No roadmap item promises a specific ship date more than one horizon out
- [ ] Every roadmap item is phrased as a problem/outcome, not a UI feature
- [ ] Outcome items have a stated baseline metric, even if rough
- [ ] Items that moved horizons this cycle have a stated reason for the move
- [ ] A stakeholder's feature request was translated into an underlying
      problem and confirmed with them before being placed on the roadmap
- [ ] A tradeoff response named the specific thing that would be
      displaced, not just "we don't have capacity"
