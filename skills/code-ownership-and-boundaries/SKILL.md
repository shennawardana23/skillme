---
name: code-ownership-and-boundaries
description: Use when the user asks to "set up CODEOWNERS", "who should review this", "define module boundaries", "this service keeps needing changes across three teams", "reorganize our repo structure", or is designing team/service boundaries, module ownership, or review-routing rules. Guides applying CODEOWNERS-based review routing and Conway's Law (system architecture mirrors org communication structure) to keep ownership boundaries matched to team structure.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Code Ownership and Boundaries

Unclear ownership shows up as either "nobody reviews this, it merges
unchecked" or "five teams have to approve every change, so nothing ships."
Both are ownership-boundary problems, not process problems to fix with
more approval steps — fix the boundary, not the process on top of it.

## Conway's Law: architecture mirrors communication structure

> "Organizations which design systems ... are constrained to produce
> designs which are copies of the communication structures of these
> organizations." (Melvin Conway, 1968)

The practical implication: if a single logical feature routinely requires
coordinated changes across three teams' codebases, that's evidence the
*module* boundaries don't match the *team* boundaries — not a sign that
the teams need better coordination processes. The fix is usually to move
the boundary (merge the modules under one team, or split the team along
the same line the module already has), not to add more cross-team
meetings. This also runs in reverse (the "inverse Conway maneuver"):
deliberately restructuring teams to match a desired target architecture,
used when you're confident in the target design and want the org to
converge toward it.

## CODEOWNERS as enforced, not aspirational, ownership

A `CODEOWNERS` file (GitHub/GitLab convention: a file listing path
patterns mapped to required reviewers/teams) turns "the payments team
owns `/services/payments/`" from a wiki page nobody checks into a rule the
platform enforces at merge time — a PR touching that path cannot merge
without an approval from that team, regardless of who authored it.

Guidelines for a `CODEOWNERS` file that actually works:

- **Map ownership to the same boundary the org chart uses**, per Conway's
  Law — if `CODEOWNERS` says team A owns a directory but the people who
  actually understand that code report to team B, the file is documenting
  a fiction and reviewers will rubber-stamp without real review.
- **Every path should resolve to exactly one clear owner** for the common
  case; shared/overlapping ownership across many paths is a sign the
  module boundary itself needs redrawing, not a routing problem to solve
  with more reviewers per rule.
- **Keep rules narrow enough that the named owner can actually review
  meaningfully** — a single team listed as owner of the entire repository
  root is not real ownership, it's a formality that guarantees shallow
  review.
- **Update `CODEOWNERS` when teams reorganize**, in the same change that
  reflects the reorg elsewhere — a stale file routes reviews to a team
  that moved on, producing exactly the rubber-stamp problem above.

## Procedure: fixing a boundary that's causing friction

1. **Name the actual pain**: is it "nobody owns this, so it rots" (no
   `CODEOWNERS` entry, or an entry pointing at a team that's stopped
   engaging) or "everyone owns this, so it's slow" (a shared module that
   forces multi-team sign-off on every change)?
2. **For the first case**: assign a single team as owner even if the
   assignment is imperfect — an owner who reviews shallowly is still
   better than a masterless file nobody reviews at all, and gives you a
   concrete team to push the "should this be split" conversation to.
3. **For the second case**: identify the sub-boundaries inside the shared
   module that map to distinct concerns (e.g., a shared "booking" service
   that actually contains pricing logic one team owns and availability
   logic another team owns) — split the module along that line, don't
   just add a "requires 2 of 3 teams" review rule on top of the merged
   module.
4. **Check whether the org structure or the module structure is the one
   that should move.** If the module split is clearly correct and stable,
   consider whether team structure should follow it (inverse Conway) — if
   the org structure is the more stable, fixed constraint, design the
   module boundary to match it instead.

## Gotchas

- A `CODEOWNERS` entry naming an individual person, rather than a team, is
  fragile — it breaks silently the day that person changes teams or
  leaves, with no forcing function to notice or update it. Name a team,
  not a person, wherever the platform supports it.
- Wildcard/catch-all `CODEOWNERS` rules at the repo root that name a
  "platform" or "infra" team as owner of everything by default often mask
  directories that actually have no engaged owner — treat a catch-all rule
  as a todo list of paths that need a real, specific owner assigned, not
  as a solved problem.
- Conway's Law predicts friction from a *mismatch*, not from
  cross-team work itself — some features genuinely require multiple
  teams' code to change together (e.g., a breaking API change and its
  consumers); the signal worth acting on is *recurring, routine* friction
  on largely independent changes, not occasional coordinated releases.
- Splitting a module to match team boundaries has a real cost (more
  cross-service calls, more deployment coordination, more infra) — verify
  the recurring friction is actually costing more than the split would
  cost before recommending it, rather than treating "Conway's Law says
  split it" as a mechanical rule with no cost side.

## Real-world grounding

Conway's Law comes from Melvin Conway's 1968 paper "How Do Committees
Invent?" and is independently, repeatedly cited in modern software
engineering literature — most notably in Matthew Skelton and Manuel
Pais's *Team Topologies* (2019), which formalizes the "inverse Conway
maneuver" (deliberately shaping team structure to produce a desired system
architecture) as a named, practiced strategy at organizations restructuring
around microservices.

## Verification

- [ ] Every actively maintained path has a `CODEOWNERS` entry naming a
      team, not an individual
- [ ] The named owner matches who actually understands and works in that
      code, not a stale org chart
- [ ] Recurring multi-team sign-off requirements are treated as a boundary
      problem to fix, not a process to route around
- [ ] `CODEOWNERS` is updated in the same change as any team reorg it
      affects
