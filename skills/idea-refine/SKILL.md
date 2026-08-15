---
name: idea-refine
description: Refine a vague or raw idea into a sharp, actionable concept through structured divergent-then-convergent ideation — expand into variations, stress-test the strongest directions, then ship a one-page brief. Use when an idea is still fuzzy, when the user wants options before committing to one direction, or when they say "ideate on this" or "stress-test my plan".
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Idea Refine

Open an idea up before narrowing it down. This is the option-generation counterpart to `interview-me`: interview-me extracts what a user already wants through one-question-at-a-time dialogue; this skill is for when the idea itself is still unformed and needs variations explored before there's a "what" to interview about.

## How It Works

Three phases, run as a conversation, not a template — adapt based on what the user says at each step rather than running every question mechanically.

### Phase 1: Understand & Expand (Divergent)

1. **Restate the idea** as a crisp "How Might We" problem statement — this forces clarity on what's actually being solved before generating options.
2. **Ask 3-5 sharpening questions, no more**: who is this for specifically; what does success look like; what are the real constraints (time, tech, resources); what's been tried before; why now. Don't proceed to variations until you know who this is for and what success looks like.
3. **Generate 5-8 variations** using these lenses, picking whichever fit the idea rather than running all of them mechanically:
   - **Inversion**: what if we did the opposite?
   - **Constraint removal**: what if budget/time/tech weren't factors?
   - **Audience shift**: what if this were for a different user?
   - **Combination**: what if we merged this with an adjacent idea?
   - **Simplification**: what's the version that's 10x simpler?
   - **10x version**: what would this look like at massive scale?
   - **Expert lens**: what would a domain expert find obvious that an outsider wouldn't?

If working inside a codebase, scan for existing architecture, patterns, and prior art before proposing variations — ground the options in what actually exists rather than inventing in a vacuum.

### Phase 2: Evaluate & Converge

After the user reacts to Phase 1 (which ideas resonate, what they push back on):

1. **Cluster** the ideas that resonated into 2-3 distinct directions — each direction should feel meaningfully different, not a variation on a theme.
2. **Stress-test** each direction against: user value (painkiller or vitamin — who benefits and how much), feasibility (what's the hardest part), and differentiation (would someone actually switch from their current solution).
3. **Surface hidden assumptions** for each direction explicitly: what you're betting is true but haven't validated, what could kill this idea, and what you're choosing to ignore for now and why that's acceptable. This is where most ideation fails — don't skip it.

Be honest, not supportive. A weak idea should be named as weak, with kindness. Push back on unnecessary complexity and question real value rather than validating whatever was proposed first.

### Phase 3: Sharpen & Ship

Produce a concrete markdown one-pager:

```markdown
# [Idea Name]

## Problem Statement
[One-sentence "How Might We" framing]

## Recommended Direction
[The chosen direction and why — 2-3 paragraphs max]

## Key Assumptions to Validate
- [ ] [Assumption — how to test it]

## MVP Scope
[The minimum version that tests the core assumption. What's in, what's out.]

## Not Doing (and Why)
- [Thing] — [reason]

## Open Questions
- [Question that needs answering before building]
```

The "Not Doing" list is arguably the most valuable section — focus means saying no to good ideas, not just bad ones, and this makes the trade-off explicit rather than implicit.

Offer to save the result to `docs/ideas/[idea-name].md` (or a location of the user's choosing). Only save if the user confirms.

## Anti-Patterns

- Generating 20+ shallow variations instead of 5-8 well-considered ones
- Being a yes-machine — accepting a weak idea without naming the weakness
- Skipping "who is this for" before generating variations
- Producing the Phase 3 brief without surfacing assumptions in Phase 2
- Ignoring existing codebase constraints when ideating inside a project
- Jumping straight to the one-pager without running the divergent and convergent phases first

## Gotchas

- A user reacting positively to every variation you generate is a signal to push harder in Phase 2, not a signal that all directions are equally good — enthusiasm at the divergent stage is not validation.
- The "Not Doing" list fails silently if it's generic ("we're not doing everything") — each entry needs a specific thing and a specific reason, or it isn't doing its job of making a trade-off visible.
- Grounding variations in an existing codebase's architecture is a constraint AND a source of ideas — don't only use it to rule things out; look for what the existing patterns make newly cheap to build.

## Real-world grounding

The divergent-then-convergent structure in this skill is a direct application of design thinking as popularized by IDEO and Stanford's d.school — a well-documented methodology that explicitly separates "expand the option space" (divergent) from "narrow to a decision" (convergent) as distinct phases, precisely because doing both at once causes premature convergence on the first plausible idea. Amazon's internal "working backwards" practice (writing the press release and FAQ for a product before building it) is a well-known real-world analog to this skill's Phase 3 one-pager: both force a team to commit a fuzzy idea to a short, concrete written artifact before resources are spent building it.

## Verification

- [ ] A clear "How Might We" problem statement exists
- [ ] The target user and success criteria are defined before variations were generated
- [ ] Multiple directions were explored, not just the first idea that came up
- [ ] Hidden assumptions are explicitly listed with a way to validate each
- [ ] A "Not Doing" list makes trade-offs explicit, with reasons, not just a list
- [ ] The output is a concrete one-page artifact, not just conversation
- [ ] The user confirmed the final direction before any implementation work started
