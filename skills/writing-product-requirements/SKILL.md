---
name: writing-product-requirements
description: Guides writing product requirements as a PR-FAQ (press release plus frequently-asked-questions) using Amazon's working-backwards methodology, or as a lighter PRD when a PR-FAQ is overkill. Use when the user asks to "write a PR-FAQ", "write product requirements", "draft a PRD", "work backwards from the customer", "write a press release for a product that doesn't exist yet", or asks for a critique/review of an existing PR-FAQ or PRD draft.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "product-management"
---

# Writing Product Requirements

Write requirements by working backwards from the customer's experience,
not forwards from a feature list. The PR-FAQ format forces this: you write
the press release and customer quote for a product that doesn't exist yet,
*before* any design or engineering work starts, and the difficulty of
writing a truthful, compelling press release is itself the test of whether
the idea deserves resources.

## Step 1: Decide PR-FAQ vs. lighter PRD

Use a full PR-FAQ when:
- The initiative is a new product, a new customer-facing surface, or a
  strategic bet that will consume more than a few weeks of a team's time.
- Multiple teams or leadership need to align on *why* before anyone
  argues about *how*.
- The idea is genuinely uncertain — if you can't yet write an honest
  customer quote, that uncertainty is exactly what the exercise surfaces.

Use a lighter one-page PRD (problem, goal, non-goals, requirements, open
questions) when:
- The work is a well-understood extension of an existing product
  (a settings toggle, an internal tool, a known bug-driven feature).
- There is no real ambiguity about customer value — only about
  implementation.
- Forcing a press release would be theater, not a forcing function. A
  PR-FAQ written for a task nobody doubts wastes a day and teaches nothing.

Doing this step wrong is a failure mode on its own: writing a PR-FAQ for
a trivial internal feature (nobody reads it, it becomes a rubber stamp) or
skipping it for a genuinely risky bet (the team argues about implementation
before agreeing on customer value) are both common misapplications.

## Step 2: Write the press release first

Write it as if the product shipped today and a journalist is covering it.
Standard structure, roughly one page:

1. **Headline** — one sentence, names the product and the customer benefit
   in language a customer would actually use. Not "New Platform Capability
   Launched" — that tells no one what changed for them.
2. **Sub-headline** — one sentence naming the target customer and the
   primary benefit.
3. **Summary paragraph** — what the product is, who it's for, the problem
   it solves. Datelined as if announced.
4. **Problem paragraph** — the customer's problem *today*, described from
   their point of view, not the company's. Should be specific enough that
   a reader recognizes the pain without being told it's a pain.
5. **Solution paragraph** — how the product solves it, in plain language,
   no internal jargon or codenames.
6. **A quote attributed to a company spokesperson** — states why the
   company built this, tying it to strategy.
7. **A quote attributed to a customer** — describes their experience using
   it, in the customer's own voice, specific enough to be falsifiable
   (what they could do before vs. after). This is the line most drafts
   get lazy on — see Gotchas.
8. **A call to action** — how a customer gets started today.

Constraint: no feature list disguised as prose, and no metric in the
release that engineering hasn't at least sanity-checked as plausible. The
press release is aspirational about the *experience*; it should not be
fictional about *feasibility*.

## Step 3: Write the FAQ

Split into two sections, both required:

**External FAQs** (a customer would ask these):
- Pricing, availability, platform/region support
- How it compares to existing alternatives (including the status quo of
  "doing nothing")
- Limitations — what it explicitly does *not* do yet

**Internal FAQs** (the team and leadership need these):
- What is the customer problem, backed by evidence (data, research,
  support tickets, direct quotes) — not assertion
- Why now — why hasn't this been built before, what changed
- Key metrics that would tell you the launch worked or failed
- Technical approach at a level a non-engineer can follow, including
  the riskiest technical assumption
- Cost / resourcing required
- Go-to-market plan
- Explicit list of what's *out of scope* for v1
- Risks and open questions the team has not resolved yet — an internal
  FAQ with no unresolved questions is a sign no one interrogated it

A PR-FAQ with only external FAQs (marketing-flavored) or only internal
FAQs (spec-flavored) has skipped half the exercise.

## Step 4: Treat it as a draft that gets rewritten, not a document to approve

The mechanism only works if the document is revised repeatedly before
anyone commits resources. Plan for multiple review passes where the
reviewer's job is to find what's unconvincing — a vague headline, a
customer quote that could describe any product, an FAQ that dodges the
hard question — and send it back for a rewrite. A PR-FAQ approved on the
first pass with no pushback usually means the reviewer treated it as a
rubber stamp rather than reading it critically.

## Gotchas

- **The customer quote is the part people fake.** A generic quote like
  "This is exactly what we needed!" attributed to "a customer" proves
  nothing — it could be pasted onto any product. A real quote names a
  specific before/after ("I used to spend 20 minutes reconciling reports
  by hand; now it's automatic") and should be falsifiable enough that
  someone could object "that's not actually true yet."
- **A vague headline hides a vague product.** If the headline could apply
  to three different products the team might build, the team hasn't
  actually agreed on what they're building — the vagueness in the doc
  is a symptom, not a wording problem to fix later.
- **Skipping the FAQ to save time defeats the purpose.** The press release
  is the easy, feel-good part; the FAQ — especially "why now" and "what's
  explicitly out of scope" — is where disagreements that would otherwise
  surface mid-build get forced out early, when they're cheap to resolve.
- **Writing the PR-FAQ after the roadmap is already committed** turns it
  into after-the-fact marketing copy instead of a decision tool — the
  entire value is writing it *before* resources are committed, so the
  document can still kill the idea.
- **One-and-done drafting.** If the first draft is treated as final, the
  team skipped the actual mechanism, which is disagreement-driven
  rewriting. Expect and budget for several rounds.
- **Confusing internal ambition with customer benefit.** A solution
  paragraph that describes a new internal architecture or a strategic
  goal ("consolidates our platform") instead of what changes for the
  customer has drifted from working-backwards into working-forwards.

## Real-world grounding

This is Amazon's "working backwards" process, publicly described by former
Amazon executives (including in Colin Bryar and Bill Carr's book
*Working Backwards*) as the mechanism Amazon uses before greenlighting new
products: instead of pitching with slides, a team writes a mock press
release and FAQ document, circulates it in a meeting where attendees read
it silently first, and the group interrogates it — often producing many
rewritten drafts over weeks — before anyone commits engineering time. The
underlying principle is "work backwards from the customer" rather than
forwards from existing technology or internal capability, and the PR-FAQ
is the concrete artifact that operationalizes that principle: if you
cannot write an honest, specific press release and a real customer quote,
that's evidence the idea isn't ready to build.

## Verification

- [ ] Headline names a specific customer benefit, not a capability
- [ ] Customer quote is specific and falsifiable, not generic praise
- [ ] Both external and internal FAQ sections are present
- [ ] Internal FAQ states why now, success metrics, and explicit out-of-scope items
- [ ] No feature list disguised as press-release prose
- [ ] Document has gone through at least one critical rewrite pass, not approved on first draft
