---
name: engineering-onboarding-design
description: This skill should be used when the user asks to "design an onboarding plan", "write a 30/60/90 day plan for a new engineer", "set up onboarding for a new hire", "what should a new engineer do in their first week", or is structuring how a new engineer ramps up on a team/codebase. Use for the structure and cadence of onboarding a new engineer — not for the content of technical documentation itself.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Engineering Onboarding Design

Good onboarding has one job in the first week: prove, with a real (small,
safe) production change, that the new engineer's access, environment, and
mental model all actually work end to end. Everything else — reading docs,
shadowing, architecture overviews — supports that goal; none of it
substitutes for it.

## Structure by phase

**Day 1**: accounts, repo access, local environment running, a assigned
onboarding buddy (not their manager — someone who remembers what it's like
to not know anything yet). Confirm they can actually run the app/tests
locally before the day ends; environment setup that "should work" but
hasn't been verified is the most common silent blocker for week one.

**Week 1**: ship one small, real, reviewed change to production. Not a
toy repo, not a doc-only PR padded to look substantial — a genuinely small
production change (a copy fix, a small bug fix, a well-scoped test
addition) that exercises the actual review/CI/deploy pipeline. This
validates access, tooling, and process simultaneously, and gives the new
hire a concrete win instead of a week of passive reading.

**First 30 days**: own one task end-to-end — design, implement, test,
ship, and (if applicable) monitor after release — with support available
but not doing it for them. Pair this with directed codebase orientation:
have them trace one real user request through the system (a "FedEx tour,"
see `exploratory-testing-techniques`) rather than reading an architecture
diagram cold.

**First 60 days**: contribute to a project or feature independently,
including making a design decision of their own (see
`technical-decision-records` for how that decision should get written
down) and participating in on-call/incident response in a shadow or
paired capacity if the team carries on-call.

**First 90 days**: full productivity — owns an area or component,
participates in code review for others (not just receiving it), and gives
feedback on the onboarding process itself while it's still fresh (see
below).

## Checklist categories to cover across the ramp

- **Access & tooling**: repo, CI/CD, secrets/credentials process,
  deployment permissions, monitoring/alerting dashboards.
- **Codebase orientation**: how to run tests, where the domain logic
  lives, how data flows through the system, what's partitioned/sharded
  and why (if relevant to this codebase).
- **Domain knowledge**: the business the software serves — for a hotel
  management platform, that means understanding what a PMS actually does
  operationally, not just the API surface.
- **Culture & process**: how decisions get made, how code review norms
  work, how incidents get handled, where the org's existing ADRs/decision
  records live.
- **First contributions**: a queued, appropriately-sized first PR ready
  on day one or two — don't leave this to be improvised mid-week.

## Gotchas

- **Onboarding docs rot fast, and the fix is dogfooding, not more
  writing.** Have the new hire fix the onboarding doc as one of their
  first tasks — they're the only person for whom every gap and stale step
  is still visible. Waiting for a "someday" doc audit rarely happens.
- **Pure shadowing/reading with no real task by day 5 is a warning
  sign**, not a sign of a thorough plan — it usually means the team
  hasn't found (or hasn't bothered to carve out) a real starter task, and
  the new hire is quietly losing momentum.
- **Trivial busywork disguised as a "first task" defeats the point.** The
  goal of the week-one PR is proving the real pipeline works; a fake task
  routed around normal review/CI doesn't prove anything and often has to
  be redone properly later.
- **Onboarding buddy ≠ manager.** Using the manager as the day-to-day
  buddy conflates "who do I ask dumb questions" with "who evaluates my
  performance," which suppresses the questions that actually need asking.
- **Ask for onboarding feedback while it's fresh** (end of week 1, end of
  30 days) — after 90 days, the new hire has adapted to the gaps and can
  no longer see them clearly enough to report them usefully.

## Real-world grounding

GitLab publishes its entire onboarding process publicly in its team
handbook (about.gitlab.com/handbook), including the explicit practice of
having new hires update onboarding documentation as one of their first
contributions — a rare case of a company's real onboarding process being
fully open for inspection rather than described secondhand. The broader
30/60/90-day structure is widely used across the tech industry as a
standard onboarding cadence, independent of any single company's specific
program.
