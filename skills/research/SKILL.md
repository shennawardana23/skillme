---
name: research
description: Run a structured, codebase-first technical research workflow to answer a specific engineering question — check internal docs and source before searching externally, verify claims against actual code, and stop within a bounded number of iterations. Use for narrow, verifiable technical questions like "does this library version support X" or "which error-handling pattern does this codebase use", not for open-ended topic research.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Research

Answer a specific technical question with evidence, not a guess. This is the narrowest of the research-family skills: `market-research` handles business decisions and `deep-research` handles open-ended web topics, but this one is for a single, verifiable engineering question about this codebase or its dependencies — a version behavior, an API signature, an existing convention.

## When to Activate

- A specific, answerable technical question needs resolving before writing code
- Deciding between two library APIs or two internal patterns
- Confirming version-specific behavior before relying on it
- Verifying a claim about "how this codebase does X" before extending it

## The Research Question

Before searching, write the question in one sentence. A good question:
- Is specific: "Does this project's ORM version support batch upserts?" not "How does the ORM work?"
- Has a verifiable answer: yes/no, a version number, a code pattern
- Scopes the investigation: "in this codebase" or "in the Go ecosystem" — not both at once

## Process

### 1. Check what's already known

Before searching externally, check, in order:
- This repo's `docs/` directory and any ADRs (see `documentation-and-adrs`)
- Existing skills under `skills/` — a skill may already encode the answer
- `go.mod` (or the relevant manifest) for the currently pinned dependency version
- Comments in the relevant source files
- Any project-wide contract or conventions file

You may already have the answer without a single external search.

### 2. Search precisely

Use specific terms, not concepts:
- Bad: "how to handle errors in Go"
- Good: "`errors.As` vs `errors.Is` for wrapped sentinel errors"

Check primary sources first: official documentation (pkg.go.dev for Go, the library's own docs site), the actual source in the module cache (Go: `$(go env GOMODCACHE)`, or `go doc <package>`), and GitHub issues/release notes for version-specific behavior.

### 3. Verify claims with code

Never trust documentation alone for version-specific behavior:
- Find the actual source in the module cache or vendor directory
- Confirm the API signature matches what the docs describe for the version actually pinned
- Read the upstream package's own tests for real usage examples — tests don't drift out of sync with behavior the way prose docs do

### 4. Record findings

```
Finding: <one sentence summary>
Evidence: <source URL or file path + line>
Confidence: HIGH | MEDIUM | LOW
```

A LOW confidence finding needs follow-up verification before anyone acts on it — don't round it up to HIGH just because it's the only lead you found.

### 5. Make a decision

If a decision is needed after research:
- State the recommendation in one sentence
- List the top 2 trade-offs
- Note what evidence would change the recommendation

## Good vs Bad Research

| Good | Bad |
|---|---|
| Primary sources (official docs, source code) | Blog posts without cited sources |
| Version-specific evidence, checked against the pinned version | Undated or old forum answers |
| "Evidence X proves Y" | "It seems like Y" |
| Acknowledges uncertainty explicitly | Presents a guess as a fact |

## Time Limit

If a research question is not resolved within 3 iterations of searching, stop and:
1. State what is known vs. unknown
2. Propose a safe default action
3. Flag what needs expert (human) input

## Gotchas

- Documentation on a package's official site frequently describes the *latest* version's behavior, not the version actually pinned in this project's manifest — always cross-check the doc's stated version against `go.mod` (or equivalent) before relying on it.
- A GitHub issue marked "closed" is not proof the described behavior was fixed in a released version — check whether the fix shipped in a tagged release, and which one.
- Two credible sources disagreeing is itself a finding worth recording (LOW confidence, with both sources cited) rather than picking whichever one you read first.

## Real-world grounding

Go's own documentation culture is the direct model for the "verify against source, not prose" rule in this skill: `pkg.go.dev` and `go doc` render documentation generated directly from source-code comments, and the Go team's own convention is that doc comments live next to the code they describe so they can't silently drift out of sync the way a separate wiki or blog post can. When behavior is genuinely ambiguous, reading the standard library's own source (available locally via `go doc -src` or in the Go toolchain's source tree) is the standard, well-established way experienced Go engineers resolve it rather than trusting secondary write-ups.

## Verification

- [ ] The research question was written as one specific, verifiable sentence before searching
- [ ] Internal sources (docs/, ADRs, go.mod, existing skills) were checked before external search
- [ ] At least one claim was verified against actual source code, not documentation alone
- [ ] Findings are recorded with an explicit confidence level
- [ ] If unresolved after 3 iterations, a safe default action was proposed instead of continuing indefinitely
