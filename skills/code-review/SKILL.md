---
name: code-review
description: Reviews code changes for security, correctness, performance, and maintainability, in Go, TypeScript, PHP, Vue, or any language. Use when given a PR URL, a diff, or a file path; before merging a change; or when asked to check for injection risks, missing edge cases, missing authorization checks, or error-handling gaps.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Code Review

Review changes in a fixed order — security first, since a maintainability
issue costs refactor time but a missed injection or auth gap costs a
breach. Work through every category even for a small diff; skipping
categories on "obviously small" changes is how small changes ship real bugs.

## Phase 1 — context

Before judging anything, establish: what language/framework conventions
apply, what problem the change solves, and its scope (feature, fix,
refactor, security, performance).

## Phase 2 — review order

1. **Security** — input validation, authN/authZ on every new endpoint or
   operation, injection risk (SQL, command, template), secrets in code or
   logs, path traversal in file operations.
2. **Correctness** — business logic matches the stated intent, boundary
   conditions (empty input, max values, concurrent access), error paths are
   handled and tested, no silent failures.
3. **Performance** — no new unbounded loops over unbounded data, no leaked
   resources (unclosed files/connections, growing slices/maps held forever),
   caching invalidation is correct if caching was touched.
4. **Maintainability** — functions stay focused, names carry intent without
   needing an explanatory comment, no duplicated logic, public APIs are
   documented.
5. **Tests** — new behavior is tested, error paths are tested, table-driven
   (Go) or parameterized (JS/TS/PHP) tests used for input variations.

## Gotchas

- For PostgreSQL tables partitioned by `hotel_id` in this organization's
  schemas, a missing `hotel_id` filter is a correctness AND a performance
  defect simultaneously — the query engine scans every partition instead of
  pruning to one. Flag this even if the query "looks" correct. See the
  `postgres-hotel-partitioning` skill for the specific patterns to check.
- A change that only touches test files or comments still needs the
  Security pass if it changes what a test asserts — a weakened assertion
  ("skip this check in CI") is a security regression disguised as test
  maintenance.
- Identical-looking conditional branches are the highest-value thing to
  read character-by-character, not skim: Apple's 2014 "goto fail" SSL bug
  was a single duplicated `goto fail;` line that caused certificate
  validation to always succeed, silently, for over a year — the kind of
  defect that is invisible on a fast skim and only surfaces on a careful
  line-by-line read of a short, unremarkable-looking function.

## Output format

```
## Risk Assessment: LOW | MEDIUM | HIGH | CRITICAL

## Critical (must fix before merge)
- [file:line] Issue description
  Fix: concrete suggestion

## High (should fix)
- ...

## Medium (recommended)
- ...

## What's good
- At least one genuine positive observation

## Score: N/10
```

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "It's a one-line change, skip the full pass" | One-line changes are exactly where "goto fail"-class bugs hide — line count doesn't correlate with risk. |
| "Tests pass so it's fine" | Passing tests only prove what they assert; a missing authz check with no test for it stays invisible. |
| "The partition filter is implicit through the join" | Postgres does not infer a `hotel_id` filter from a join alone — verify it appears explicitly in the query or every plan generated from it. |
