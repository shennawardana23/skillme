---
name: project-guidelines-example
description: Shows what a project-specific guidelines document (architecture map, directory contract, response/error shape, exact commands, deployment steps, Always/Ask-first/Never boundaries) needs to contain to be useful to an agent working in that codebase, with one worked example. Use when asked to write a CLAUDE.md-style project guide, a "house rules" doc for a specific repo, or a project-onboarding reference -- not when authoring a new reusable skill for this catalog (see skill-catalog-authoring for that).
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Project Guidelines Document

A generic coding-standards skill tells an agent how to write good code in
general. A project guidelines document tells it the specific, otherwise
undiscoverable facts about *this* codebase — which database, which error
envelope, which exact command runs the tests — that no amount of general
knowledge would supply. This skill is about what such a document needs to
contain, not about how to package it as a reusable skill (for that, see
`skills/skill-catalog-authoring/`, which covers this repo's own SKILL.md
and eval conventions).

## When to use

Asked to write a project-onboarding doc, a "house rules" reference, or a
`CLAUDE.md`/`AGENTS.md`-style guide for a specific codebase the agent will
keep working in across sessions.

**Skip it for:** authoring a new general-purpose skill meant to be reused
across many projects — that's a different artifact with a different
audience (see `skills/skill-catalog-authoring/`).

## What it must contain

A guidelines document earns its keep only if it answers questions a
fresh agent session would otherwise have to guess at or rediscover by
reading code. Six things, in order of how often they're actually needed:

### 1. Architecture overview

The services/processes involved and how they talk to each other — not a
generic layered-architecture diagram, but the specific pieces of *this*
system:

```
┌─────────────────────────────┐
│           API                │  Go, net/http + chi router
│  Deployed: Cloud Run          │
└─────────────────────────────┘
              │
   ┌──────────┼──────────┐
   ▼          ▼          ▼
┌────────┐ ┌────────┐ ┌────────┐
│Postgres│ │ Redis  │ │  S3    │
│(tenant_id│ │ cache  │ │ assets │
│partition)│ │        │ │        │
└────────┘ └────────┘ └────────┘
```

### 2. Directory contract

Where things go, stated as a rule an agent can follow when creating a new
file, not just a snapshot of what already exists:

```
cmd/<service>/      → entry points, one per deployable binary
internal/domain/    → core types, no framework import
internal/service/   → business logic
internal/repository/→ persistence, behind an interface
internal/handler/   → HTTP boundary only — no business logic here
```

### 3. Response/error shape

The exact envelope every endpoint uses — this is pure convention, entirely
undiscoverable from general knowledge, and inconsistency here breaks every
client at once:

```go
type APIResponse[T any] struct {
    Success bool   `json:"success"`
    Data    T      `json:"data,omitempty"`
    Error   string `json:"error,omitempty"`
}
```

### 4. Exact commands

Full, copy-pasteable commands, not tool names:

```
Build: make build
Test:  make test
Lint:  make lint
Run:   make run
```

A command name without its flags ("run the tests") forces a guess; the
guess is sometimes wrong in a way that silently skips coverage or a
required env var.

### 5. Testing requirements

Framework, where tests live, and any hard threshold that's actually
enforced in CI — stated as *this project's* number, not a generic
recommendation:

```
Framework: go test + testify/require
Location:  _test.go next to the file under test
Threshold: 80% enforced by `make test-coverage` in CI (not a suggestion)
```

### 6. Boundaries: Always / Ask first / Never

The same three-tier boundary structure used in
`skills/spec-driven-development/`, but populated with *this project's*
specifics rather than generic examples:

```
Always:    run `make lint` and `make test` before committing
Ask first: schema changes, new external dependencies, CI config changes
Never:     commit .env files, bypass hotel_id partition filtering in
           queries, deploy directly to prod without the staging gate
```

## A worked example

For an illustrative Go + PostgreSQL service (not a specific named product
— the point is the shape, reusable for any real project):

```markdown
# Project Guidelines: Orders Service

## Architecture
Go 1.23, net/http + chi. PostgreSQL (partitioned by tenant_id), Redis
for caching, deployed to Cloud Run.

## Directory Contract
cmd/orders/       entry point
internal/domain/  Order, LineItem types — no sql/http imports
internal/service/ business logic, depends on repository interfaces
internal/repo/    Postgres queries — every query filters tenant_id explicitly

## Response Shape
{"success": bool, "data": T, "error": string} — every handler, no exceptions.

## Commands
Build: make build   Test: make test   Lint: make lint

## Testing
go test + testify. 80% coverage enforced in CI via make test-coverage.

## Boundaries
Always: run make lint && make test before commit.
Ask first: new migrations, new external dependencies.
Never: query a partitioned table without an explicit tenant_id filter.
```

## Gotchas

- A guidelines doc that restates general best practice ("write clean
  code," "handle errors properly") instead of this project's specific
  facts is dead weight — an agent already knows the generic advice; what
  it doesn't know is which envelope shape, which partition column, which
  exact make target.
- Listing the directory structure as a snapshot ("here's what exists
  today") without stating it as a rule ("new handlers go here") leaves an
  agent to infer the rule from precedent, which breaks the first time the
  precedent is inconsistent.
- A "Never" list without a stated reason invites a future session to treat
  it as an arbitrary restriction and quietly violate it once it seems
  inconvenient — "never skip the tenant_id filter" survives review much
  better paired with "queries are partitioned by tenant_id; an unfiltered
  query scans every partition and can leak cross-tenant data."
- Letting the doc go stale (commands change, a service gets renamed) is
  worse than not having one — an agent trusts a documented command more
  than it verifies one, so a stale command produces a confidently wrong
  action, not a cautious question.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "The code is self-documenting, we don't need this" | Self-documenting code shows *what* exists, not *where new things should go* or *which command runs the tests* — those are conventions, not derivable from reading one file. |
| "Everyone on the team already knows this" | The document isn't for the team that wrote it — it's for the next session (human or agent) that has none of that shared context. |
| "We'll update it when it's out of date" | An agent trusts a stale doc as if it were current; "we'll update it later" means it's wrong for however long "later" takes. |

## Real-world grounding

GitLab's publicly published Handbook — one of the most cited examples of
an organization documenting its own specific conventions rather than
generic best practice — is explicit that its value comes from stating
*this company's* actual process (how a specific type of decision gets
made, which exact tool is used for what) rather than general management
advice available in any textbook; the same principle scoped down to one
codebase is what separates a useful project guidelines doc from a
restatement of `skills/coding-standards/`.

## Verification

- [ ] The document states this project's specific facts (envelope shape,
      partition column, exact commands), not generic advice
- [ ] The directory contract reads as a rule for new files, not just a
      snapshot of existing ones
- [ ] Every command is copy-pasteable with its actual flags
- [ ] The Always/Ask first/Never list gives a reason for each Never, not
      just the restriction
- [ ] The testing threshold, if any, states what's actually enforced in
      CI, not an aspirational number
