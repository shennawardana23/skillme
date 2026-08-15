---
name: user-research-synthesis
description: Guides synthesizing raw qualitative user research into defensible themes using affinity mapping (the KJ method) and conducting Jobs-to-be-Done "Switch" interviews about past purchase decisions. Use when the user asks to "affinity map these notes", "synthesize interview findings", "cluster these observations into themes", "run a JTBD interview", "find out why customers switched to us", or reviews raw interview transcripts, support tickets, or survey verbatims for patterns.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# User Research Synthesis

Two failure modes dominate bad research synthesis: forcing raw data into
categories decided in advance, and asking interviewees what they *want*
instead of reconstructing what they actually *did*. This skill covers two
named, well-established techniques that avoid both — affinity mapping (the
KJ method) for synthesizing a batch of observations into themes, and
Jobs-to-be-Done "Switch" interviewing for understanding a specific past
decision.

## Affinity mapping (the KJ method)

Developed by Japanese anthropologist Jiro Kawakita, and adopted widely in
UX research as a bottom-up alternative to top-down categorization. The
entire point is that themes emerge from the data instead of being decided
before you've looked at it.

### Procedure

1. **Atomize first.** Before mapping, break every source (interview
   transcript, support ticket, survey verbatim) into individual, single-idea
   observations — one sticky note (physical or virtual) per idea. A
   transcript quote covering three separate points becomes three notes, not
   one. Mapping breaks down immediately if a note carries more than one
   idea, because it can't be sorted into a single cluster.
2. **Silent individual clustering first.** Each participant in the session
   moves notes into groups **without talking**. This is not a formality —
   discussing groupings while forming them causes the most vocal or senior
   person's mental model to become everyone's mental model before the data
   has been seen fresh. Silence is what makes the bottom-up property real.
3. **Merge and negotiate second.** Once individuals have independently
   clustered, compare groupings as a group, discuss disagreements, and
   merge into a shared set of clusters.
4. **Name clusters last, after they exist.** Write the theme label only
   once a cluster has taken shape from the notes inside it. Naming a
   cluster before sorting notes into it (or pre-writing category headers
   before the session starts) is functionally indistinguishable from
   deductive coding against a preset framework — it defeats the reason to
   use this method over a spreadsheet with columns.
5. **Look for outlier notes.** A note that resists joining any cluster is
   informative, not a defect — it often flags a segment or edge case the
   rest of the data doesn't represent. Don't force it into the
   nearest-fitting group just to close the session.
6. **Write the synthesis from the clusters, not from memory.** The output
   is a named theme plus the observations that support it, so a reader can
   trace every claim back to raw data — a synthesis with no traceable
   observations underneath it is an opinion, not a finding.

### When to use it

Any time you have a batch of qualitative data larger than can be held in
one person's head at once — interview note dumps, support ticket exports,
open-ended survey responses, usability-test observations. It does not
require a facilitator credential; it requires discipline about the
ordering (atomize → sort silently → merge → name), which is the part teams
skip under time pressure.

## Jobs-to-be-Done "Switch" interviews

Clayton Christensen's Jobs-to-be-Done theory holds that customers don't buy
products, they "hire" them to make progress on a specific job. Bob Moesta
and Chris Spiek operationalized this into a concrete interview technique —
commonly called the "Switch" interview — that reconstructs the timeline of
one specific, real purchase or switch decision, rather than asking people
what they think they'd want in the future.

### Why this differs from a typical feature-request interview

A standard interview asks "what features would you want?" or "what
frustrates you about your current tool?" — both invite hypothetical,
future-tense answers that people are bad at predicting and that tend to
just confirm whatever's already on the roadmap. A Switch interview instead
anchors on a real, already-completed event: the moment someone actually
paid for (or actually stopped using) something. Memory of a real decision
is more reliable than a prediction about a hypothetical one, and it can't
be shaped by what the interviewer hopes to hear because it already
happened before the interview.

### The timeline structure

Walk the interviewee backward and forward through four timeframes around
one specific switch:

1. **First thought.** When did the idea of a different solution first
   cross their mind? What event triggered it? (Often a "struggling moment"
   — something broke, changed, or became newly painful.)
2. **Passive looking.** Between first thought and active search, what
   passively kept the idea alive (a colleague's comment, an ad noticed but
   not acted on)?
3. **Active looking / deciding.** What changed to make them actively
   evaluate options? What alternatives did they compare, and what pushed
   them away from the old way and pulled them toward the new one?
4. **Consuming / onboarding.** What happened right after they bought or
   switched — including any anxiety or second-guessing before commitment?

At each stage dig for two forces pushing them forward (push of the current
situation's pain, pull of the new solution's promise) and two forces
holding them back (anxiety about the new choice, habit/attachment to the
old one). The switch happens only when push + pull outweigh anxiety +
habit — that tension is the actual finding, not a feature list.

### Procedure

1. Recruit people who made a **real, recent, specific** switch — not
   people describing a hypothetical future purchase, and not people who
   are still using the old solution and merely dislike it.
2. Open with the concrete event: "Tell me about the day you decided to
   [switch/buy/cancel] ___." Get a date or approximate timeframe to anchor
   the timeline.
3. Move through first thought → passive looking → active looking →
   consuming, asking "what happened next" and "what were you feeling at
   that point" rather than "what would you want."
4. Probe for the emotional and social context, not just the functional
   one — who else was involved in the decision, what would people around
   them think, what was at stake if they chose wrong.
5. Never ask "what features do you want us to build" during the
   interview — that question belongs to a different exercise
   (prioritization), and asking it here collapses the interview back into
   ordinary feature-request gathering.

## Gotchas

- **Leading questions bias toward the existing roadmap.** "Would a
  dashboard like X help you?" invites agreement regardless of the truth,
  because interviewees tend to be polite and pattern-match to what they
  think the interviewer wants. Ask about past behavior and past decisions
  instead of pitching a solution and gauging reaction.
- **Sorting into pre-existing categories defeats affinity mapping.**
  Bringing a persona framework, a set of feature categories, or last
  quarter's theme names into the room before the sort has started causes
  every note to be filed under an existing label instead of surfacing what
  the data actually says. If the cluster names were decided before the
  session, it wasn't affinity mapping — it was deductive coding wearing
  affinity mapping's clothes.
- **One vocal user is not a segment.** A single detailed, articulate
  complaint (frequently from the interviewee most similar to the person
  running the research) gets remembered and repeated far more than its
  actual prevalence in the data warrants. Cross-check any strong claim
  against the cluster size and the number of independent sources behind
  it before treating it as representative — a theme with one supporting
  note is an anecdote, a theme with notes from five different interviews
  is a finding.
- **JTBD interviews degrade into ordinary interviews under time
  pressure.** The moment the conversation shifts from "tell me what
  happened" to "what do you wish existed," the push/pull/anxiety/habit
  structure is lost and the interview produces the same generic
  feature-request data the technique was designed to avoid.
- **Recency and vividness bias which switch stories get told.** People
  reconstruct memory, and recent or emotionally intense events dominate
  recall even when they're not typical. Ask for the specific timeframe and
  cross-check timeline details (what else was happening then) rather than
  accepting a smoothed-over retelling at face value.

## Real-world grounding

Affinity mapping's KJ method comes from Jiro Kawakita's work organizing
field research data in the 1960s; it was adopted into UX practice as the
standard technique for turning a pile of sticky notes from usability
studies and interviews into named themes without pre-imposed categories,
and remains a staple exercise in UX research training and workshops today.

Jobs-to-be-Done originates with Clayton Christensen's work on disruptive
innovation (notably illustrated by the "milkshake" study, where people
"hired" milkshakes for a boring-commute job unrelated to typical
breakfast-food competition). Bob Moesta and Chris Spiek, working with
Christensen, turned the theory into the practical "Switch" interview
method — mapping the four timeframes (first thought, passive looking,
active looking, consuming) and the push/pull/anxiety/habit forces — which
is now widely taught and used as a concrete interviewing technique distinct
from the broader JTBD theory itself.

See `references/jtbd-question-bank.md` for a longer set of stage-by-stage
sample questions if the four-timeframe structure above needs more detail
during an actual interview.
