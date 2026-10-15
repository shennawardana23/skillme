---
name: me-on-review
description: This skill should be used when the user asks to "review this PR/repo as me", "compile a review of these changes", "give me a human-readable summary of this code review", or wants a whole pull request, branch, or repository's changes turned into one narrative report in their own voice rather than a per-file finding list. Aggregates and compiles findings across multiple files/commits into a single readable report with a merge/no-merge call, applying Archipelago/Sentec-specific lenses (PostgreSQL partition-by-hotel_id filtering, Go idioms, metric units, Vue over React, no fabricated metrics). Does not replace the mechanical per-file scanners — see code-review, code-reviewer, and security-review, whose findings this skill compiles and narrates rather than re-detects.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Me, On Review

A compilation skill, not a scanner. The mechanical work of finding
issues already belongs to `code-review` (per-diff correctness/security/
performance), `code-reviewer` (fixed-category severity scan), and
`security-review` (security-priority checklist). This skill's job starts
after those findings exist: turn a PR, a branch, or a whole repository's
change set into one report a human reads top to bottom in a few minutes,
in the reviewer's own voice, with a clear call at the end.

Use this for the *whole* unit of change — every file in a PR, every
commit on a branch, or a full-repo pass — not a single file or snippet.
For that, use one of the scanner skills directly.

## What "in my voice" means here

Layer these lenses on top of whatever the scanners found, every time:

- **PostgreSQL partition-by-hotel_id.** Any query against a table
  partitioned by `hotel_id` that doesn't filter on it explicitly is a
  correctness-and-performance finding, not a style nitpick — it lets the
  planner scan every partition. Flag it even if no scanner caught it;
  this is Archipelago-specific and generic reviewers don't know it.
- **Go and PostgreSQL by default.** A new service, script, or migration
  in a different stack needs a stated reason in the PR description. No
  reason stated is itself a finding.
- **Metric units.** Flag imperial units introduced in code, comments, API
  fields, or docs (miles, °F, lb) unless the change is integrating a
  third-party API that's contractually imperial.
- **Vue, not React**, for any new frontend surface in this org's repos.
  A new React component is a finding unless the PR states why the repo
  it's in is React-based.
- **No fabricated confidence.** Flag marketing language in comments,
  commit messages, or docs ("blazingly fast", "100% secure") and any
  invented metric/benchmark with no measurement behind it.
- **Undo-ability.** For anything destructive or hard to reverse (a
  migration, a config change, a deletion), check the PR shows how to roll
  it back. Missing rollback plan on a risky change is a finding, not a
  suggestion.

## Workflow

1. **Collect.** Run or read the output of the relevant scanner skill(s)
   for every changed file — don't re-derive their findings by hand.
2. **Apply the lenses above** across the whole change set, not per file —
   a hotel_id filtering gap in one file and a second one three files
   later are the same systemic finding, reported once.
3. **Group by theme, not by file.** A human reading this report should
   understand *what's wrong with the change* in a few groups (e.g.
   "Partition filtering," "Test coverage," "Rollback safety"), not scroll
   through a per-file list that repeats the same issue five times.
4. **Compile into narrative form** — see Output below. This is the
   differentiator: the scanners produce tables; this skill produces
   something a non-reviewing stakeholder (a lead, the PR author) can read
   without cross-referencing file names first.

## Output

```md
## Summary
2-4 sentences: what this change does, and the headline verdict.

## What's solid
Real things done well — specific, not filler praise. Skip this section
before inventing praise; "nothing stands out" is a valid, honest line.

## Findings, by theme
### <Theme, e.g. "Partition filtering">
Plain-language explanation, citing file:line for every instance,
followed by the fix.

## Recommendation
One of: Merge as-is / Merge after the required items below / Do not
merge until these are fixed — with the specific list that gates the
next step, separated from nice-to-haves.
```

- Cite `file:line` for every concrete claim; never generalize a finding
  found in one file across the whole PR without checking the others.
- State a fix for every finding — a report that only names a problem
  without a next step is incomplete.
- If the change set is clean, say so as a two-line summary. A short
  approval is a complete, correct output — don't pad it to look
  thorough.
