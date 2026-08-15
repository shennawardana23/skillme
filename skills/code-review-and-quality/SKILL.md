---
name: code-review-and-quality
description: Runs code review as a change-management process — right-sizing the diff, labeling findings by required-vs-optional severity, applying structural remedies for complexity, and resolving reviewer/author disagreements. Use when a change needs splitting before review, when review comments need severity labels, when a "should I approve this" call is being made, or when setting up a team's review process rather than auditing a single file.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Code Review and Quality

This is a process skill, not a line-by-line checklist — see `code-reviewer` for the mechanical category scan and `security-review` for the security-priority walkthrough. This skill governs the surrounding decisions: how big should this change be, how should feedback be labeled, when do you approve versus block, and how do disagreements get resolved. It is grounded in Google's public Code Review Developer Guide and Engineering Practices documentation, adapted to be model-agnostic.

## The approval standard

Approve a change when it definitely improves overall code health, even if it isn't perfect. Perfect code doesn't exist. Don't block a change because it isn't exactly how you'd have written it — if it improves the codebase and follows project conventions, approve it. This standard exists specifically to stop reviews from becoming a venue for imposing personal style.

## Change sizing — decide before reviewing content

```
~100 lines changed   → Good, reviewable in one sitting
~300 lines changed   → Acceptable if it is a single logical change
~1000 lines changed  → Too large — ask the author to split it
```

Watch the *file's total size* too, not just the diff: a change that pushes an already-large file past roughly 1000 total lines is a signal to extract helpers or modules first, then add the new logic — decompose, then add, not the reverse.

**What counts as "one change":** a self-contained modification addressing one thing, with its own tests, that leaves the system working. Refactoring and feature work in the same diff is two changes — split them (small renames can stay at reviewer discretion).

**Splitting strategies:**

| Strategy | How | When |
|---|---|---|
| Stack | Submit a small change, base the next one on it | Sequential dependencies |
| By file group | Separate diffs for groups needing different reviewers | Cross-cutting concerns |
| Horizontal | Land shared code/stubs first, then consumers | Layered architecture |
| Vertical | Break into full-stack slices of the feature | Feature work |

Complete file deletions and pure automated refactors are the exception — they can be large if the reviewer only needs to verify intent.

## Severity labeling

Label every comment so the author can triage without guessing:

| Prefix | Meaning | Author action |
|---|---|---|
| *(no prefix)* | Required change | Must address before merge |
| **Critical:** | Blocks merge | Security vulnerability, data loss, broken functionality |
| **Nit:** | Minor, optional | May ignore — formatting, style preference |
| **Optional:**/**Consider:** | Suggestion | Worth thinking about, not required |
| **FYI** | Informational | No action needed |

Order findings by leverage: correctness and security first, then structural regressions and missed simplifications, then everything else. If there's one structural problem and ten nits, the structural problem *is* the review — don't bury it.

## Structural remedies

When flagging complexity, propose the specific move, not just "this is complex":

- Replace a chain of conditionals with a typed model or explicit dispatcher.
- Collapse duplicate branches into one clearer flow.
- Separate orchestration from business logic.
- Move feature-specific logic out of a shared module into the package that owns the concept.
- Reuse the canonical helper instead of a bespoke near-duplicate.
- Make a type boundary explicit (kill a gratuitous `any`/cast/silent fallback) so downstream branching disappears.
- Delete a pass-through wrapper that adds indirection without clarifying the API.

The test for "did this refactor help": count the concepts a reader must hold to follow the change. If a "cleaner" version leaves that count unchanged, it relocated complexity rather than reducing it — prefer the version where whole branches, modes, or layers disappear.

## Review process

1. **Understand the context** — what is this trying to accomplish, what spec/task does it implement, what behavior changes?
2. **Review the tests first** — they reveal intent and coverage faster than the implementation does. Do tests exist, do they test behavior (not implementation details), are edge cases covered, would they catch a regression?
3. **Review the implementation** against correctness, readability, architecture, security, and performance (delegate the mechanical per-category pass to `code-reviewer` if a deep scan is needed).
4. **Verify the verification** — what tests were run, did the build pass, was it tested manually, are there before/after screenshots for UI changes?

## Change descriptions

Every change needs a description that stands alone in history. First line: short, imperative, standalone ("Delete the FizzBuzz RPC", not "Deleting the FizzBuzz RPC"). Body: what and why, with context not visible in the diff itself — link bugs, benchmarks, design docs. Anti-patterns to reject: "Fix bug", "Fix build", "Moving code from A to B", "Phase 1".

## Handling disagreements

Resolve disputes in this order: (1) technical facts and data override opinions, (2) the style guide is the absolute authority on style, (3) software design is judged on engineering principles, not personal preference, (4) codebase consistency is acceptable only if it doesn't degrade overall health. Don't accept "I'll clean it up later" — require it before submission, or require the author to file a self-assigned follow-up bug.

## Dead code hygiene

After a refactor, list orphaned code explicitly and ask before deleting:

```
DEAD CODE IDENTIFIED:
- formatLegacyDate() in src/utils/date.ts — replaced by formatDate()
- LEGACY_API_URL constant in src/config.ts — no remaining references
→ Safe to remove these?
```

## Dependency discipline

Before approving a new dependency: does the existing stack already solve this? How large is it? Is it actively maintained? Any known vulnerabilities (`npm audit` / `go list -m -u all`)? Is the license compatible? Prefer the standard library and existing utilities — every dependency is a liability the reviewer inherits too.

## Review speed

A typical change should get multiple review rounds in a single day; respond within one business day as the ceiling, not the target. Prioritize a fast individual response over a fast final approval — quick feedback reduces author frustration even across several rounds. If a change is too large to review responsibly, ask the author to split it rather than rubber-stamping it.

## Gotchas

- "The tests pass" is not evidence of a good review — tests only prove what they assert. A missing per-resource authorization check with no test written against it stays invisible through a full green CI run.
- A refactor that moves code around without shrinking the number of concepts a reader must track is not a simplification — it is complexity in a new location. Judge the result, not the intent.
- Small diffs are exactly where oversized-file and bolted-on-conditional problems hide, because a 5-line addition to a 900-line file "looks" too small to warrant asking whether the file should be split first.
- "I'll clean it up later" is a specific, well-documented failure mode — deferred cleanup overwhelmingly does not happen. Treat it as a request for a follow-up bug with an owner, not a promise to hold the reviewer to.

## Real-world grounding

Google's public Code Review Developer Guide — the source this skill is grounded in — explicitly states the standard is "improves overall code health" rather than "is perfect," precisely because reviewers who hold changes to a personal-perfection bar create review bottlenecks that outweigh the marginal quality gained; the same guide is the origin of the CL-size guidance (~100 lines ideal, ~1000 too large) used here.

## Verification

- [ ] All Critical issues are resolved
- [ ] All required (no-prefix) changes are resolved or explicitly deferred with justification
- [ ] The change is one logical unit, sized for single-sitting review (or a documented exception)
- [ ] Tests pass and the build succeeds
- [ ] The change description stands alone and explains why, not just what
- [ ] Dead code identified during the change was surfaced, not silently deleted or silently kept
