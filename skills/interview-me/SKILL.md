---
name: interview-me
description: Extract what a user actually wants, rather than what convention makes them ask for, via a one-question-at-a-time interview with a stated confidence level until reaching roughly 95% confidence in the underlying intent. Use when a request is underspecified (missing who/why/success/constraint), when the user says "interview me" or "are we sure?", or before any plan/spec/code exists for a non-trivial ask.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Interview Me

Close the gap between what people ask for and what they actually want, before it's expensive to close. Someone asks for "a dashboard" because that's the conventional answer, not because a dashboard solves their problem. This skill is the pre-decision counterpart to `idea-refine`: idea-refine generates variations once you know roughly what's wanted; this skill is for when you don't yet know what's wanted at all.

## When to Use

- The ask is missing at least one of: **who** the user is, **why** they want it, what **success** looks like, what the binding **constraint** is
- The request is conventional rather than specific ("build me X", "make it faster") and the convention can't be unpacked without guessing
- Two reasonable values are in tension (simplicity vs. flexibility, cost vs. speed) and the user hasn't said which one they're optimizing for
- The user explicitly invokes it: "interview me", "grill me", "before we start, are we sure?"

**Do not use** when the ask is unambiguous ("rename this variable", "fix this typo"), the user has explicitly asked for speed over verification, it's a pure information request ("how does X work?"), or confidence is already ≥95% — re-read the stop condition below before assuming it isn't.

## Loading Constraints

This skill needs a live, responsive user. Do not invoke it in non-interactive contexts (CI pipelines, scheduled runs, autonomous loops). If the ask is underspecified there, flag it as a blocker instead of guessing or interviewing an absent user.

## The Process

### Step 1: Hypothesize, with a confidence number

Before asking anything, write your current best read of what the user wants in one sentence, plus an honest confidence number (0-100%):

```
HYPOTHESIS: You want a way to answer "how are we doing?" in standup, and "dashboard" was the convention that came to mind.
CONFIDENCE: ~30% — missing: who it's for, what "metrics" means, what success looks like
```

Below ~70% confidence, append the reason on the same line — it tells the user exactly what the interview needs to surface.

### Step 2: Ask one question at a time, each with a guess attached

```
Q: <one focused question>
GUESS: <your hypothesis for the answer, with the reasoning behind it>
```

Wait for the reaction before asking the next question. One at a time, not a batch: a batch invites skim-reading, later questions often depend on earlier answers, and attaching a guess is faster for the user to react to than generating an answer from scratch — it also commits the agent to a hypothesis it can be visibly wrong about. The risk is a polite user agreeing just to be agreeable; mitigate by being visibly willing to be wrong, including guessing in directions expected to draw pushback.

### Step 3: Listen for "want vs. should want"

Watch for answers that pattern-match best-practice talk without specifics ("I want it to be scalable"), defer to convention ("the way most apps do it"), or use buzzwords as goals ("modern", "robust"). When you hear these, ask: *"If you didn't have to justify this to anyone, what would you actually want?"* — this single question often does more work than the previous five.

### Step 4: Restate intent in the user's own words

At high confidence, write back a tight restate the user can confirm or correct line by line:

```
- Outcome:      <one line>
- User:         <one line — who benefits>
- Why now:      <one line — what changed>
- Success:      <one line — how we know it worked>
- Constraint:   <one line — the binding limit>
- Out of scope: <one line — what we're explicitly not doing>

Yes / no / refine?
```

"Out of scope" is non-negotiable — half of misalignment is silent disagreement about what is *not* being built.

### Step 5: Confirm — explicit yes, not "whatever you think"

None of these count as confirmation: "whatever you think is best" (delegation, not confidence — re-ask with two concrete options); "sounds good" or "sure, let's go" (ambiguous or a polite exit — ask "anything you'd refine?"); silence followed by "okay let's start" (the user gave up on the interview, not converged — stop and ask what was missed). If corrected, fold the correction in and restate again. Loop until an explicit yes.

### The 95% Confidence Stop

Done when the answer to *"can I predict the user's reaction to the next three questions I would ask?"* is yes. If no, ask the next question. This has a floor: after several rounds without rising confidence, that's information about the ask, not a reason to keep grinding — say so and offer to step back.

## Output

The deliverable is the confirmed statement of intent from Step 4 plus the explicit yes from Step 5 — not a spec, plan, or task list; those are downstream and consume this intent. If it should persist across sessions, offer to save it to `docs/intent/[topic].md`, only after confirmation.

## Gotchas

- A confidence number with no reason attached below ~70% is not useful to the user — it doesn't tell them what to help close, so it functions as decoration rather than signal.
- Three or more rounds without confidence visibly rising means the wrong questions are being asked, not that more rounds are needed — reframe rather than continuing the same line.
- Saving the intent doc before the user confirms implies a yes that hasn't been given — never write the persisted artifact before Step 5 completes.

## Real-world grounding

The one-question-at-a-time, hypothesis-first structure here is the same idea behind Toyota's "five whys" technique from the Toyota Production System: rather than asking one broad question and accepting the first answer, each question builds on the previous answer and is asked separately, because a batch of questions invites surface-level answers the way a single root-cause question does. Basecamp/37signals has written publicly and repeatedly about pushing back on vague feature requests until the actual underlying need is named before any design work starts — the same "convention-signaling answer" problem this skill's Step 3 is built to catch.

## Verification

- [ ] An explicit hypothesis with a confidence number was stated in the first turn
- [ ] Every confidence number below ~70% carried a one-line reason
- [ ] Questions were asked one at a time, each with the agent's own guess attached
- [ ] At least one "what would you actually want if you didn't have to justify it?" probe ran when warranted
- [ ] A concrete restate (Outcome/User/Why now/Success/Constraint/Out of scope) was written back
- [ ] The user gave an explicit yes — not "whatever you think," not "sounds good," not silence
- [ ] At the stop point, the agent could predict reactions to its next three questions
