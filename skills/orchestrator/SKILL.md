---
name: orchestrator
description: Use when a request could be handled by more than one skill in your catalog, or the task type isn't obvious from the wording (e.g. it mixes a code review with a security question, or a feature request with test generation). Trigger phrases include "not sure which skill applies", "this touches security and code style", "route this to the right skill", or "what's the process for a request like this". Teaches how to build and use a signal-to-skill routing table rather than guessing.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Orchestrator

A routing discipline for picking which skill in a catalog actually
applies to an incoming request, and for handling requests that span more
than one. The core failure this prevents: improvising a generic answer
because the exact matching skill wasn't obviously the top hit, when a
specific skill with real domain knowledge exists and should be loaded and
followed instead.

## Build a routing table, don't route from memory

Maintain an explicit table of task signal to skill, and consult it before
answering — don't rely on recalling what skills exist from general
impression. A table survives catalog growth; memory doesn't.

| Signal in the request | Route to | Example phrasing |
|---|---|---|
| "is this correct", a diff, a PR to review | your code-review skill | "review this PR", "is this a good change" |
| "is this secure", auth/input-handling changes | your security-review skill | "check for vulnerabilities", "is this safe to expose" |
| A schema change, a new table, partitioned tables | your database-migration skill | "add a column", "migrate this table" |
| "implement", "build", "add feature" | your feature-implementation / TDD skill | "add support for X" |
| "write tests", "add coverage" | your test-driven-development or e2e-testing skill | "test this", "add E2E coverage" |
| An error message, a stack trace | your debugging skill | "why is this failing" |
| Go-specific idiom questions | your go-service-idioms skill | "is this idiomatic Go" |
| Framework-specific UI questions (Vue/Nuxt, Laravel) | your framework-pattern skill | "how should this component be structured" |
| A request spanning code generation *and* review | your agentic-engineering / ai-first-engineering skills | "decompose this for an agent", "what should the review focus on" |

Adapt the right-hand column to your catalog's actual skill names — the
point is the shape of the table (a signal maps to exactly one primary
skill), not these specific entries.

## Process

1. **Identify the task type** from the signals in the request, using the
   table above — not from a first impression of what the request "feels
   like."
2. **Load the matched skill in full** and follow its instructions as
   written. Don't improvise a plausible-sounding process instead of
   reading the actual skill — the whole point of having a skill catalog
   is that the specific skill encodes domain knowledge a generic response
   won't reproduce.
3. **For multi-step tasks, follow the skill's own phase structure** in
   the order it specifies, rather than reordering steps because a
   different order seems more natural.

## Handling requests that span multiple skills

A request touching more than one domain (a security review that also
raises API design questions, a feature request that should also get test
coverage) doesn't get answered by picking the closest single match and
ignoring the rest:

- **Start with the highest-risk skill first.** Security concerns get
  addressed before style concerns; a data-integrity question on a
  partitioned table gets addressed before a naming-convention question.
- **Load each relevant skill for its part of the answer**, rather than
  answering the whole thing from only the first-matched skill's
  perspective.
- **Ask exactly one clarifying question** if the request is genuinely
  ambiguous about scope — not zero (which risks answering the wrong
  question) and not several (which stalls a request that was actually
  answerable with reasonable defaults).

## Tone and process discipline

- Be direct: state what you're about to do before doing it, rather than
  narrating after the fact.
- Use numbered steps for anything with more than one stage.
- Flag assumptions explicitly when you make one instead of routing
  ("Assuming this is the production database because none was
  specified.") — a silent assumption is much harder to catch and correct
  than a stated one.

## Gotchas

- A request that matches two rows in the table with comparable strength
  is a sign the table needs a tiebreaker rule (usually: risk order), not
  a sign to pick whichever skill you thought of first.
- Loading a skill's title/description but not its full body and then
  answering from the description alone reproduces the exact "generic
  answer instead of domain knowledge" failure this pattern exists to
  prevent.
- A routing table that's never updated as the catalog grows silently
  degrades into "route everything to the skill I happen to remember" —
  treat the table itself as something to maintain, not a one-time setup.

## Real-world grounding

This is the same shape as the "front controller" pattern used by most web
frameworks: a single entry point inspects the incoming request and
dispatches it to the specific handler that owns that route, rather than
every handler trying to also parse and interpret requests meant for
other handlers. The value is the same in both cases — one place decides
"what kind of request is this," so the specialized handlers can each
assume they're only ever invoked for the case they're built for.

## Verification

- [ ] The request was matched against an explicit table, not routed from memory
- [ ] The matched skill was loaded and followed in full, not paraphrased from its description
- [ ] A multi-domain request had its highest-risk skill addressed first
- [ ] At most one clarifying question was asked, only if genuinely ambiguous
- [ ] Assumptions made during routing were stated explicitly, not left implicit
