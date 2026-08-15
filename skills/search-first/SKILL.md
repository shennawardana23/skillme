---
name: search-first
description: Search for an existing solution — in this repo, as an installed MCP server, as an existing skill, or as a maintained package — before writing custom code for a new utility, dependency, or integration. Use before creating a new helper/abstraction, before adding a dependency, or whenever a request to "add X functionality" is about to turn into new code that already exists somewhere else.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Search First

Check whether the problem is already solved before writing code to solve it. This is the only research-family skill in this catalog that ends in a coding decision rather than a report — `research`, `deep-research`, and `market-research` all produce answers or documents; this one produces a choice between adopting, extending, composing, or building.

## When to Activate

- Starting a new feature that likely has an existing solution
- Adding a dependency or integration
- About to write a new utility, helper, or abstraction
- The request is "add X functionality" and the next step would otherwise be writing custom code

## Workflow

```
1. NEED ANALYSIS       — define what functionality is needed, and any language/framework constraint
2. PARALLEL SEARCH     — repo grep, installed MCP servers, skills/ catalog, package registry, GitHub
3. EVALUATE            — score candidates: functionality fit, maintenance activity, community size, docs, license, dependency weight
4. DECIDE              — Adopt as-is / Extend or wrap / Compose 2-3 small packages / Build custom, but informed by the search
5. IMPLEMENT           — install the package, configure the MCP server, or write the minimal custom code
```

## Decision Matrix

| Signal | Action |
|--------|--------|
| Exact match, well-maintained, permissive license | **Adopt** — install and use directly |
| Partial match, good foundation | **Extend** — install + write a thin wrapper |
| Multiple weak matches | **Compose** — combine 2-3 small packages |
| Nothing suitable found | **Build** — write custom, but informed by what the search ruled out |

## Quick Mode (inline check)

Before writing a utility or adding functionality, run through:
0. Does this already exist in the repo? Grep through relevant modules/tests first.
1. Is this a common problem? Check the Go standard library and the package registry (`pkg.go.dev`, npm, PyPI as relevant).
2. Is there an MCP server for this already configured?
3. Is there a skill for this? Check `skills/*/SKILL.md` in this repo.
4. Is there a maintained GitHub implementation or template before writing net-new code?

## Full Mode (agent)

For non-trivial functionality, delegate the search itself: launch a research agent to search the package registry, MCP servers, this repo's skills, and GitHub, and return a structured comparison with a recommendation — rather than doing the search inline while already committed to a coding mindset.

## Search Shortcuts by Category

### Development Tooling (Go-first)
- Linting → `golangci-lint`, `staticcheck` (other ecosystems: `eslint`, `ruff`)
- Formatting → `gofmt` / `goimports` (other ecosystems: `prettier`, `black`)
- Testing → `go test` with table-driven tests (other ecosystems: `jest`, `pytest`)
- HTTP client → standard library `net/http` first; add a client package only for genuine gaps (retries, circuit breaking)

### AI/LLM Integration
- Claude SDK / API usage → check `documentation-lookup` (Context7 MCP) for current docs before writing integration code
- Document processing → check for maintained Go or Python libraries before hand-rolling a parser

### Data & APIs
- Validation → Go: hand-rolled validation or a struct-tag validator library; TS: `zod`; Python: `pydantic`
- Database access → check for an existing MCP server before writing a new client wrapper

### Content & Publishing
- Markdown processing → `goldmark` (Go), `remark`/`unified` (JS)
- Image optimization → check for an existing pipeline before adding a new dependency

## Anti-Patterns

- **Jumping to code**: writing a utility without checking whether one already exists in the repo or a package registry
- **Ignoring MCP**: not checking whether an installed MCP server already provides the capability
- **Over-customizing**: wrapping a library so heavily it loses the benefit of using it at all
- **Dependency bloat**: installing a large package for one small feature it happens to contain

## Gotchas

- "Adopt" is not risk-free — check the license and maintenance activity even for an exact functional match; an abandoned or copyleft-licensed package can cost more later than writing 20 lines now.
- A partial match that requires heavy wrapping is often a sign to re-run the search with different terms rather than committing to "Extend" — the right package may exist under different naming.
- Searching only the package registry and skipping the repo grep step produces false negatives: the functionality may already exist internally under a name that doesn't match the search terms used externally.

## Real-world grounding

The 2016 npm "left-pad" incident — a widely documented case where an 11-line package (`left-pad`) was unpublished and broke builds across a large portion of the JavaScript ecosystem that had transitively depended on it — is the standard cautionary tale for the "Compose" and "Build" branches of this decision matrix: pulling in a tiny, thinly-maintained dependency to avoid writing a few lines of code carries its own real fragility risk, which the maintenance-activity and dependency-weight scoring in the Evaluate step exists to catch. It sits alongside the general "don't reinvent the wheel" engineering norm that motivates the Adopt/Extend path in the first place — both failure modes are real, which is why this skill scores candidates rather than defaulting to either extreme.

## Verification

- [ ] The repo was grepped for existing internal solutions before searching externally
- [ ] Installed MCP servers and this repo's `skills/` catalog were checked
- [ ] Candidates were scored on functionality, maintenance, license, and dependency weight — not chosen on the first result
- [ ] The final decision (Adopt/Extend/Compose/Build) is stated explicitly with its rationale
