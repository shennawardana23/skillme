---
name: git-workflow-and-versioning
description: Guides git commit discipline, trunk-based branching, and semantic versioning/changelogs for released code, including Go modules' major-version import path rule. Use when committing, branching, resolving conflicts, cutting a release, choosing a semantic version bump, tagging a release, or writing a changelog entry.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Git Workflow and Versioning

Commits are save points, branches are sandboxes, and history is
documentation. A version number is the separate promise you make to
*consumers* of your code — the moment anything outside your own working
copy depends on a piece of code, "latest on main" stops being a
sufficient answer to "what am I running, and is it safe to upgrade?"

## Trunk-based development (default)

Keep `main` always deployable. Work in feature branches that merge back
within 1–3 days — long-lived branches accumulate merge risk every day
they exist, and DORA/*Accelerate* research consistently correlates
trunk-based development with high-performing engineering teams. Prefer
a feature flag over a long-lived branch for anything incomplete (see
`ci-cd-and-automation`).

```
main ──●──●──●──●──●──●──●──  (always deployable)
        ╲      ╱  ╲    ╱
         ●──●─╱    ●──╱    ← short-lived feature branches (1-3 days)
```

## Commit discipline

**Commit each working increment**, not one giant commit at the end —
each commit is a point you can revert to if the next change breaks
something. **Each commit does one logical thing**; don't mix a refactor
with a feature, or formatting with behavior change, in the same commit.

```
# Good: atomic, self-contained commits
a1b2c3d feat: add reservation cancellation endpoint
d4e5f6g refactor: extract status-transition validation to its own function
h7i8j9k test: add table-driven tests for cancellation edge cases

# Bad: everything mixed together
x1y2z3a add cancellation, fix unrelated bug, reformat whole file, bump deps
```

**Messages explain why, not just what:**

```
feat: add email validation to registration endpoint

Prevents malformed addresses from reaching the database. Uses the
same validation approach as auth.ts for consistency.
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`. Target ~100–300
lines per commit/PR; split anything approaching 1000 lines.

## Branching

```
feature/<short-description>   feature/reservation-cancellation
fix/<short-description>       fix/duplicate-charge-on-retry
chore/<short-description>     chore/bump-go-version
refactor/<short-description>  refactor/reservation-service
```

Branch from `main`, delete after merge, keep it short-lived. For parallel
work by multiple agents or engineers, use `git worktree add
../repo-feature-a feature/a` — each worktree is an independent checkout
on its own branch, so nobody needs to stash or switch branches to work
on something else concurrently.

## Pre-commit hygiene

```bash
git diff --staged                                          # review before committing
git diff --staged | grep -iE "password|secret|api_key|token"  # scan for accidental secrets
go test ./... && golangci-lint run                          # tests + lint pass locally first
```

Automate with a pre-commit hook (`lefthook`, `pre-commit`, or Husky for
Node) rather than relying on everyone remembering to run it by hand.

## Using git for debugging

```bash
git bisect start
git bisect bad HEAD
git bisect good <known-good-commit>   # git checks out midpoints; test each

git blame path/to/file.go             # who last touched this line, and in which commit
git log --grep="reservation cancel"   # search commit messages for a keyword
```

`git bisect` only works well against a history of atomic commits that
each build and pass tests independently — a history of giant, mixed
commits makes every bisect step ambiguous about which part of the commit
actually caused the regression.

## Semantic versioning

For anything with consumers (a published library, an internal shared
package, a versioned API), version `MAJOR.MINOR.PATCH`:

```
MAJOR  breaking change — consumers must change their code to upgrade
MINOR  new functionality, backward-compatible — safe to upgrade
PATCH  bug fix, backward-compatible — safe to upgrade
```

The number is a promise; make the code match it. A "patch" that changes
behavior a consumer observably relied on is a major change wearing a
disguise — this is Hyrum's Law again, covered from the interface-design
side in `api-and-interface-design`. When genuinely unsure whether a
change is breaking, treat it as breaking; a surprise major release is
far cheaper than a broken consumer.

### Go modules: the major-version import path rule

Go modules enforce semantic versioning at the import-path level in a way
most ecosystems don't: **a v2 or later major version must live at a new
import path** (`module example.com/pkg/v2`), not just a new git tag —
this is what lets a v1 consumer and a v2 consumer coexist in the same
build's module graph without one silently resolving to the other's
breaking API:

```go
// go.mod for the v2 release
module github.com/archipelago/sentec-client/v2

go 1.23
```

```go
// consumer picks the major version explicitly via the import path
import sentecv2 "github.com/archipelago/sentec-client/v2"
```

A v1.x → v2.0.0 tag pushed **without** updating the module path to
`/v2` is a broken Go module from the consumer's perspective — `go get`
will not resolve it as intended.

### Tag the release; derive the version from the tag

```bash
git tag -a v1.4.0 -m "Release 1.4.0"
git push origin v1.4.0
```

A release is an immutable point in history. Derive the version shown in
build output, `go.mod`, or `package.json` from the tag at build/release
time rather than hand-editing a version string in a source file — that
keeps the artifact, the tag, and the changelog from ever disagreeing.

### Changelog written for humans, not `git log`

```markdown
## [1.4.0] - 2025-06-12
### Added
- Bulk reservation import via CSV
### Fixed
- Timezone drift in recurring rate-plan schedules
### Deprecated
- `GET /v1/reservations/all` — use paginated `GET /v1/reservations` (removal in 2.0)
```

Write the changelog entry in the same change that makes the change,
while the user-facing impact is still fresh — reconstructing it from
commit messages at release time loses exactly the impact framing a
changelog exists to provide. Breaking entries get a migration note (see
`deprecation-and-migration`); the release process itself is a
deployment/CI concern (`ci-cd-and-automation`, `deployment-patterns`) —
this section is only the versioning contract that feeds them.

## Gotchas

- `git commit --amend` on a commit that's already been pushed and
  reviewed rewrites history a collaborator may have already based work
  on — prefer a new commit once anyone else could plausibly have pulled
  the branch.
- Squash-merging destroys the atomic-commit narrative that made `git
  bisect` and `git blame` useful during review — decide the squash
  policy per-repo deliberately, not as a default nobody chose.
- A Go module major-version bump that forgets the `/v2` import-path
  change isn't caught by `go build` in the *publishing* repo — it only
  surfaces when a consumer tries to `go get` the new major version and
  gets confusing resolution behavior.
- A `.gitignore` added after secrets were already committed doesn't
  remove them from history — a leaked credential needs rotation, not
  just a `.gitignore` entry, once it has been pushed even once.

## Real-world grounding

DORA's *Accelerate* research (Forsgren, Humble, Kim) is the widely cited
industry source establishing the correlation between trunk-based
development, small/frequent changes, and elite software delivery
performance — the same research that underlies "faster is safer" as a
principle in `ci-cd-and-automation`: smaller, more frequent, atomic
changes are consistently associated with lower change-failure rates than
large infrequent ones, not higher.

## Verification

- [ ] Each commit does one logical thing and states why, not just what
- [ ] Feature branches are short-lived (1-3 days) and deleted after merge
- [ ] No secrets appear in `git diff --staged` before committing
- [ ] Version bumps match actual impact: breaking → major, additive → minor, fix → patch
- [ ] Go modules bumping to v2+ update the import path (`/v2`) to match, not just the tag
- [ ] Releases are tagged, with the version derived from the tag, not hand-edited out of sync
- [ ] Every user-facing release has a curated changelog entry written when the change was made
