## What it does

Guides writing product requirements as a PR-FAQ (Amazon's working-backwards
format) or a lighter one-page PRD when a full PR-FAQ would be theater. The
defining constraint: the exercise's value comes from the *difficulty* of
writing a truthful press release and a real customer quote before any
design or engineering starts — a PR-FAQ that was easy to write, or approved
without pushback on the first draft, usually means the exercise wasn't
actually done.

## When to reach for it

Reach for this skill when starting a new customer-facing initiative, a
strategic bet, or anything where multiple teams need to align on *why*
before arguing about *how*. Skip the full PR-FAQ (use the lighter PRD path
this skill also covers) for well-understood extensions of an existing
product where there's no real ambiguity about customer value — forcing a
press release onto a settings toggle wastes a day and teaches nothing.

## Prerequisites

None to run the exercise itself, but it works best with at least
preliminary evidence (support tickets, user research, usage data) to ground
the internal FAQ's "what is the customer problem" answer — an internal FAQ
built entirely on assertion rather than evidence undermines the mechanism.

## Common questions

- **"How do we know if we need a full PR-FAQ or just a lighter PRD?"** If
  you can already write an honest, specific customer quote without
  stretching the truth, the idea may not need the full exercise — a lighter
  PRD (problem, goal, non-goals, requirements, open questions) is enough
  when there's no real ambiguity about customer value, only about
  implementation.
- **"Our customer quote sounds generic — does that matter?"** Yes, it's the
  part people fake most often. A quote like "this is exactly what we
  needed" could describe any product and proves nothing. A real quote
  names a specific before/after and should be falsifiable — someone should
  be able to say "that's not actually true yet" if it's wrong.
- **"We wrote a great press release but skipped a detailed FAQ to save
  time — is that okay?"** No — the FAQ, especially "why now" and "what's
  explicitly out of scope," is where disagreements that would otherwise
  surface mid-build get forced out while they're still cheap to resolve.
  Skipping it defeats the exercise's actual purpose.
- **"Our PR-FAQ got approved on the first read with no pushback — good
  sign?"** Usually the opposite. The mechanism depends on critical,
  repeated rewriting; first-draft approval with no pushback typically means
  the reviewer treated it as a rubber stamp rather than reading it to find
  what's unconvincing.

## It's working if

- The headline names a specific customer benefit, not a capability
- The customer quote is specific and falsifiable, not generic praise
- Both external and internal FAQ sections exist, including an explicit
  out-of-scope list and unresolved risks
- The document has been through at least one critical rewrite before
  anyone commits resources to it

## Where it fits

A standalone, upstream planning skill — typically the first artifact in a
larger initiative, ahead of any engineering-facing spec or ticket
breakdown. Distinct from `defining-done-and-acceptance-criteria` (which
operates once work is already scoped) and `feature-prioritization-frameworks`
(which decides between already-defined candidates, not whether one
candidate deserves to exist at all).
