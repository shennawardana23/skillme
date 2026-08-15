---
name: defining-done-and-acceptance-criteria
description: This skill should be used when the user asks to write acceptance criteria for a user story, define or review a team's Definition of Done, evaluate whether a story is "ready" or "done", or asks things like "write Given/When/Then criteria for this story", "is this Definition of Done good enough", "why do we keep disagreeing about whether this is done", or "review this ticket before we start it". Distinguishes the team-wide Definition of Done from per-story acceptance criteria and applies INVEST and Gherkin-style Given/When/Then to make both testable.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Defining Done and Acceptance Criteria

Two different tools solve two different problems. Confusing them is the
single most common cause of "done" meaning different things to different
people on the same team.

- **Definition of Done (DoD)** — one checklist, owned by the whole Scrum
  team, applied to *every* Product Backlog Item and Increment. It answers
  "is this shippable, full stop" (code reviewed, tests passing, deployed to
  staging, docs updated, no known regressions). It rarely mentions the
  feature itself.
- **Acceptance Criteria (AC)** — specific to *one* story, written during
  refinement, before or as work starts. It answers "does this particular
  story do the thing the user needs" (e.g. "search returns results in
  under 500ms for a 3-word query"). It says nothing about code review or
  deployment — that's the DoD's job.

A story can satisfy its acceptance criteria and still not be Done (it
wasn't code-reviewed). A story can pass every DoD checklist item and still
be wrong (it was reviewed, tested, and deployed — against the wrong
requirement). Both checks are required; neither substitutes for the other.

## Procedure: writing acceptance criteria for a story

1. **Restate the story as a user need**, not a task list. If the story
   reads like a technical to-do ("add Redis cache to search endpoint"),
   push back and ask what observable behavior changes for the user — AC
   are written against behavior, not implementation.
2. **Apply INVEST to the story first** (Bill Wake's heuristic for a good
   user story — AC quality is capped by story quality):
   - **Independent** — can be built and shipped without waiting on another
     unfinished story. If AC keep referencing "after story X is done,"
     the stories should probably be split differently or merged.
   - **Negotiable** — a placeholder for a conversation, not a frozen spec.
     If the AC already dictate the database schema or the exact button
     copy, the story has quietly become a technical spec, not a story.
   - **Valuable** — delivers something a user or the business actually
     cares about, not just a technical stepping stone. "Refactor the auth
     module" has no user-facing AC because it isn't a user story.
   - **Estimable** — the team can size it. If AC are so vague the team
     can't agree on small/medium/large, the story needs more detail or a
     spike first.
   - **Small** — fits in a sprint, ideally a few days. If AC run past 5-6
     Given/When/Then blocks, the story is probably two stories.
   - **Testable** — every AC has an observable pass/fail. If an AC can't be
     checked by a test or a manual QA step, rewrite it or cut it.
3. **Write each criterion as Given/When/Then** (Gherkin-style, from Dan
   North's Behavior-Driven Development): `Given` the starting state,
   `When` the user or system does something, `Then` the observable
   outcome. This format forces precondition, action, and expected result
   to be separated instead of blurred into one vague sentence.
   ```
   Given a logged-in user with 0 saved searches
   When they submit a search query of 3+ characters
   Then results appear within 500ms
   And a "no results" state is shown if nothing matches
   ```
4. **Write the negative and edge cases explicitly**, not just the happy
   path. "Search returns results" has no AC for empty query, query with
   only special characters, or the backend timing out — each of those is
   a separate Given/When/Then, and each is where bugs actually hide.
5. **Read every criterion back and ask "could two people disagree on
   pass/fail after reading this?"** If yes, it's not actually testable —
   replace subjective language ("fast", "intuitive", "handles errors
   gracefully") with a number, a specific error message, or a named state.
6. **Check the criteria don't smuggle in implementation.** "Then a Redis
   cache entry is created" is an implementation detail masquerading as an
   acceptance criterion — nobody outside engineering can verify it, and it
   locks in a technical approach the team may want to change later.
   Rewrite as the user-observable effect: "Then the second identical
   search returns in under 50ms."

## Procedure: writing or reviewing a Definition of Done

1. **Draft it as a team, not top-down.** Per the Scrum Guide, if the
   organization has no standing organizational Definition of Done, the
   Scrum Team creates one appropriate to the product, and the Developers
   are then required to conform to it — a DoD imposed unilaterally by one
   person or by management tends to get silently ignored under deadline
   pressure because nobody bought in.
2. **Make every line independently verifiable**, ideally by tooling, not
   memory: "linter passes," "test coverage does not decrease," "deployed
   to staging and smoke-tested," "changelog entry added." Avoid entries
   that require a person to *remember* to check something with no
   artifact proving it happened.
3. **Keep it applicable to every single item**, no exceptions carved out
   per-story. The moment "well, this one doesn't really need tests" is
   said out loud, the DoD has stopped being a Definition of Done and
   started being a suggestion.
4. **Expect it to grow as the team matures** — the Scrum Guide describes
   DoD evolving over time. A new team's DoD might be "code reviewed,
   builds, deployed to staging." A mature team's DoD might add
   "accessibility checked," "feature-flagged," "monitoring dashboard
   updated." Revisit the DoD at retrospectives, not just once at project
   start.
5. **Enforce it in the actual workflow**, not just in a wiki page. If
   "deployed to staging" is on the DoD but nothing blocks a story from
   moving to Done without it, the DoD is aspirational, not real. Tie it to
   the board's column-exit criteria or a PR merge check where possible.
6. **Distinguish DoD from a Definition of Ready (DoR).** DoR gates work
   *starting* (story is small enough, AC are written, dependencies known);
   DoD gates work *finishing*. Some teams conflate these and end up with a
   DoD that's actually a mix of both, which makes it unclear when in the
   workflow each item should be checked.

## Gotchas

- **"Done" silently means different things per person** even on teams
  that have a written DoD, because the DoD was written once and never
  re-read. A developer means "code merged"; QA means "tested"; the PM
  means "in production." Ask each role what "done" means for the same
  ticket — if the answers differ, the DoD isn't actually shared, it's
  just documented.
- **Acceptance criteria written after the code is built** almost always
  describe what the code does, not what the user needed — they become a
  test-after-the-fact rationalization rather than a spec. Insist AC exist
  *before* implementation starts, even if they change during development.
- **A DoD with "write tests" but no coverage or CI gate** gets skipped
  first whenever a deadline is tight, precisely because nothing enforces
  it — the DoD becomes theater. An unenforced checklist item is worse
  than no item, because it creates false confidence that the check
  happens.
- **AC that reference internal component or variable names** ("the
  `SearchCache` service returns...") instead of user-observable behavior
  break the moment the implementation is refactored, even though the
  actual user-facing behavior hasn't changed — the AC then has to be
  rewritten for reasons unrelated to the feature, which is a sign it was
  written at the wrong level of abstraction.
- **Treating INVEST as a checklist to satisfy after the fact** rather than
  a lens to write the story through in the first place produces stories
  that technically pass each letter but still don't hang together — e.g.
  "small" by word count but still coupled to three other unfinished
  stories, so not actually "independent."
- **One universal DoD applied to fundamentally different work types**
  (a UI tweak vs. a data-migration script) sometimes doesn't fit either
  well — a common fix is a shared baseline DoD plus explicit additional
  criteria per story type, not silently ignoring the parts that don't
  apply.

## Real-world grounding

The Definition of Done comes directly from the official Scrum Guide (Ken
Schwaber and Jeff Sutherland): it describes DoD as a formal description of
the state of the Increment when it meets the quality measures required for
the product. If no organization-wide standard exists, the Scrum Team
creates its own DoD and the Developers are required to conform to it; the
Guide explicitly notes DoD evolves as a Scrum Team matures — a newer
team's DoD is typically less demanding than an experienced team's. INVEST is Bill Wake's widely-adopted mnemonic
(Independent, Negotiable, Valuable, Estimable, Small, Testable) for
judging whether a user story is well-formed enough to work from.
Given/When/Then comes from Behavior-Driven Development, introduced by Dan
North as a way to write specifications in a structured, testable format
that both technical and non-technical stakeholders can read — it is the
format underlying tools like Cucumber and Gherkin syntax.

## Verification

- [ ] Every acceptance criterion is written as Given/When/Then (or an
      equivalent structured, testable statement)
- [ ] Acceptance criteria describe observable behavior, not implementation
- [ ] Negative and edge cases have their own criteria, not just the happy
      path
- [ ] The story itself passes INVEST before criteria are finalized
- [ ] The team's Definition of Done is written down, applies to every
      item with no exceptions, and each line is independently verifiable
- [ ] Every DoD item is actually enforced somewhere in the workflow (CI
      gate, board rule, merge check) — not just documented
