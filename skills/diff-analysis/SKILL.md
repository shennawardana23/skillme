---
name: diff-analysis
description: Analyzes a git diff or PR across three tracks — complexity, style/convention, and security impact — run together, then synthesized into a single risk-graded report with required actions cited to file:line. Use when given a raw diff, PR number, or branch to triage before merge, especially when the ask is "how risky is this change" rather than a full line-by-line review.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Diff Analysis

A diff-shaped review: it only looks at *what changed*, not the whole file, and produces one risk verdict. Run all three tracks over the same diff, then synthesize — don't review them one after another as separate passes, since a single line often matters to more than one track (a deleted nil-check is both a Track A complexity signal and a Track C security signal).

## Track A: complexity

For each changed file:
1. **Cyclomatic complexity change** — did the diff add nested conditionals, new loops, new goroutines?
2. **Function length** — is any new or modified function longer than ~40 lines?
3. **Removed safety guards** — were any validation checks, nil-checks, or early-returns deleted?
4. **New external dependencies** — does the diff add new imports? Are they necessary, or does an existing dependency already cover this?

## Track B: style and convention

1. **Naming** — do new identifiers follow the language's naming convention?
2. **Error handling** — are new errors returned/wrapped properly? Any new `_ =` suppressions?
3. **Documentation** — are new exported symbols documented?
4. **Test additions** — does the diff include tests for the new behavior?
5. **Dead code** — does the diff add unreachable code or unused variables?

## Track C: security impact

1. **Removed auth checks** — did the diff delete or comment out authentication/authorization code?
2. **New SQL/command strings** — does the diff introduce a new string-built query or shell command?
3. **Weakened validation** — did any input validation become less strict (a `must` check turned into a warning, a bound loosened)?
4. **New exposed endpoints** — are new routes/handlers added without the same security middleware as their neighbors?
5. **Dependency updates** — if a lockfile/`go.mod` changed, are the new or updated dependencies from a trusted source with no known critical CVEs?

## Synthesis

```
## Risk Level: LOW | MEDIUM | HIGH | CRITICAL

### Risk justification
<One paragraph explaining the overall rating, referencing which track(s) drove it>

## Key Changes
- <What the diff actually does, in plain language>

## Required Actions (must complete before merge)
1. <Specific action with file:line reference>

## Recommended Improvements (should do)
1. ...

## Nice-to-have
1. ...

## Merge verdict: APPROVE | REQUEST_CHANGES | BLOCK
```

Always cite a specific file:line. A finding without a location is not actionable — the author has to re-find what you're describing, which defeats the point of triaging the diff instead of the whole file.

## Gotchas

- A diff can look small and still be high-risk: deleting a single `if err != nil { return err }` line is a one-line diff that removes a safety guard (Track A) and may also remove the only check standing between untrusted input and a downstream sink (Track C) — line count and risk are uncorrelated, weight the *kind* of line over its count.
- Context lines matter as much as `+`/`-` lines for Track C: an added endpoint can look secure in isolation but be missing the auth middleware every sibling route in the same file already has — compare against the file's own pattern, not an abstract standard.
- A diff-only view can miss whether a "new" function is actually new or a near-duplicate of an existing one elsewhere in the codebase; if Track A flags a large new function, a quick check for an existing near-equivalent belongs in Recommended Improvements even though it requires looking outside the diff.
- Dependency-lockfile diffs are easy to skim past because they're large and mechanical — but a single altered line in `go.sum`/`package-lock.json` can represent a full major-version jump or a supply-chain substitution; treat any lockfile change as requiring the same scrutiny as a `go.mod`/`package.json` line, not less.

## Real-world grounding

Apple's 2014 "goto fail" SSL bug was a single duplicated `goto fail;` line — a diff of essentially zero complexity by any line-count or nesting metric, yet it silently made certificate validation always succeed for over a year. It's the canonical case for why a diff-analysis process needs a dedicated Track C pass even on changes Track A would score as trivially low-risk: complexity metrics and security risk are different axes, and a diff can score low on one while being critical on the other.

## Verification

- [ ] All three tracks were run against the same diff, not just the track that seemed most relevant
- [ ] Every Required Action cites a specific file:line
- [ ] The risk level is justified by naming which track(s) drove it, not asserted alone
- [ ] Any removed check (auth, validation, nil-guard) was traced to see whether it has an equivalent replacement elsewhere in the diff, or is simply gone
