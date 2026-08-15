---
name: documentation-and-adrs
description: Record architecture decisions and documentation that capture why, not just what — use when making a significant architectural decision, choosing between competing approaches, changing a public API, shipping user-facing behavior, or documenting context a future engineer or agent will need.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Documentation and ADRs

Code shows what was built. Documentation is the only place that explains why it was built this way and what alternatives were rejected — that context is what future humans and agents actually need and can't recover from the code alone.

## When to Use

- Making a significant architectural decision
- Choosing between competing approaches
- Adding or changing a public API
- Shipping a feature that changes user-facing behavior
- Onboarding a new engineer or agent to the project
- Finding yourself explaining the same rationale repeatedly

**Skip it for**: obvious code, comments that restate what the code already says, and throwaway prototypes.

## Architecture Decision Records (ADRs)

ADRs capture the reasoning behind a significant technical decision — the highest-value documentation available for its cost to write.

### When to Write One

Choosing a framework, library, or major dependency; designing a data model or schema; selecting an auth strategy; deciding between API architectures; choosing infrastructure or hosting; any decision that would be expensive to reverse.

### Template

Store ADRs under `docs/adr/` (or `docs/decisions/` — either is a common convention; use whichever this repo already has, and pick one consistently if neither exists yet), numbered sequentially:

```markdown
# ADR-001: Use PostgreSQL for primary database

## Status
Accepted | Superseded by ADR-XXX | Deprecated

## Date
2025-01-15

## Context
[Requirements and constraints that made a decision necessary]

## Decision
[What was decided, stated plainly]

## Alternatives Considered

### Option A
- Pros: ...
- Cons: ...
- Rejected: [specific reason]

## Consequences
[What this decision commits the team to, including new risks or new capabilities]
```

### Lifecycle

`PROPOSED → ACCEPTED → (SUPERSEDED or DEPRECATED)`. Never delete an old ADR — it's the historical record of why the *previous* decision was made, and a new decision should be written as a new ADR that explicitly references and supersedes it, not an edit to the old one.

## Inline Documentation

Comment the *why*, not the *what*:

```go
// BAD: restates the code
// increment retry count by 1
retryCount++

// GOOD: explains non-obvious intent
// Rate limit uses a sliding window — reset at the window boundary,
// not on a fixed schedule, to prevent burst attacks at window edges.
if time.Since(windowStart) > windowSize {
    count = 0
    windowStart = time.Now()
}
```

Don't comment self-explanatory code. Don't leave a `// TODO: add error handling` comment when you could just add the error handling now. Don't leave commented-out code — version control already has the history; delete it.

Document known gotchas directly at the point they'd bite someone:

```go
// InitializeTheme must be called before the first render. Calling it
// after hydration causes a flash of unstyled content because the theme
// context isn't available during SSR. See ADR-003 for the full rationale.
func InitializeTheme(theme Theme) { ... }
```

## API Documentation

For Go exported symbols, use standard doc comments starting with the identifier name so `go doc` and pkg.go.dev render them correctly:

```go
// CreateTask creates a new task for the given user.
//
// Returns ErrValidation if title is empty or exceeds 200 characters,
// and ErrUnauthenticated if the caller has no valid session.
func CreateTask(ctx context.Context, input CreateTaskInput) (*Task, error) {
    // ...
}
```

For REST APIs, an OpenAPI/Swagger spec is the equivalent contract — keep request/response schemas and error responses defined there, not only in prose.

## README Structure

Every project should cover: a one-paragraph description, a quick-start (clone, install, configure, run), a command table (build/test/lint/dev), a short architecture overview linking out to the relevant ADRs, and a contributing section.

## Changelog Maintenance

For shipped features, record Added/Fixed/Changed entries per release with enough specificity that a reader can tell what changed without reading the diff.

## Documentation for Agents

- **Rules files** (CLAUDE.md and equivalents) document conventions so an agent follows them instead of guessing.
- **Specs** should be kept current so an agent builds the right thing rather than an outdated one.
- **ADRs** let an agent understand *why* a past decision was made, which prevents it from silently re-deciding a settled question.
- **Inline gotchas** prevent an agent from falling into a known trap that isn't otherwise visible in the code's structure.

## Gotchas

- An ADR with only the winning option documented is worth much less than one that also names what was rejected and why — the rejected alternatives are what stop the same debate from being re-litigated in six months by someone (human or agent) who wasn't there the first time.
- "We'll document it once the API stabilizes" tends to mean it never gets documented — writing the doc first is itself a design review, since an API that's awkward to describe is usually awkward to use.
- A comment explaining *why* stays correct as the code around it changes; a comment explaining *what* silently drifts out of sync with the code the moment someone edits the logic without updating the comment — this is why only the former is worth writing.

## Real-world grounding

The ADR format used in this skill — Context / Decision / Alternatives Considered / Consequences, stored as small numbered files that are never deleted, only superseded — comes directly from Michael Nygard's 2011 blog post "Documenting Architecture Decisions," which introduced the lightweight ADR practice now widely adopted across the software industry specifically because it's cheap enough that teams actually keep doing it, unlike heavier architecture-documentation formats that get abandoned after the first few entries.

## Verification

- [ ] ADRs exist for the significant architectural decisions made in this project
- [ ] The README covers quick start, commands, and an architecture overview
- [ ] Exported/public API symbols carry parameter, return, and error documentation
- [ ] Known gotchas are documented inline at the point they'd bite someone
- [ ] No commented-out code remains
- [ ] Rules files are current with the project's actual conventions
