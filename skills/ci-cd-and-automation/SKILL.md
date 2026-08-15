---
name: ci-cd-and-automation
description: Guides CI/CD pipeline design — quality gates, GitHub Actions for Go/TypeScript/PHP, feature-flag-decoupled deploys, secrets scoping, and pipeline optimization. Use when setting up or modifying a build/test/deploy pipeline, adding an automated quality gate, configuring CI secrets, or diagnosing a slow or flaky pipeline. For the rollout strategy the pipeline triggers, see deployment-patterns; for the image it builds, see docker-patterns.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# CI/CD and Automation

CI/CD is the enforcement mechanism for every other engineering
discipline in this catalog — it's what makes "tests pass," "lint is
clean," and "the migration was reviewed" apply to every change, every
time, instead of depending on a human remembering to check.

**Shift left**: catch problems as early in the pipeline as possible. A
lint failure costs a minute; the same defect reaching production costs
hours of incident response — order pipeline stages cheapest-and-fastest
first so a change fails fast on a trivial problem instead of waiting
20 minutes for a full integration suite to discover it.

**No gate can be skipped.** If lint fails, fix the lint issue — don't
disable the rule to unblock a merge. If a test is flaky, fix the
flakiness — don't re-run until it passes.

## The quality gate pipeline (Go)

```
PR opened
  → gofmt / golangci-lint
  → go vet
  → go build ./...
  → go test ./... -race -cover
  → govulncheck ./...
  → (main only) build & push image → deploy staging → smoke test → deploy production
```

```yaml
# .github/workflows/ci.yml
name: CI
on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: "1.23", cache: true }
      - uses: golangci/golangci-lint-action@v6

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: "1.23", cache: true }
      - run: go build ./...
      - run: go test ./... -race -cover
      - run: go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Run `lint` and `test` as **separate parallel jobs**, not sequential steps
in one job — a lint failure and a test failure are independent facts
about the same change; running them in parallel means the PR author
sees both results in the time of the slower one, not the sum of both.

## Database-backed integration tests

```yaml
  integration:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_DB: testdb
          POSTGRES_USER: ci_user
          POSTGRES_PASSWORD: ${{ secrets.CI_DB_PASSWORD }}
        ports: ["5432:5432"]
        options: >-
          --health-cmd pg_isready --health-interval 10s --health-timeout 5s --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: "1.23", cache: true }
      - name: Run migrations
        run: migrate -path migrations -database "$DATABASE_URL" up
        env:
          DATABASE_URL: postgresql://ci_user:${{ secrets.CI_DB_PASSWORD }}@localhost:5432/testdb
      - run: go test ./... -tags=integration
        env:
          DATABASE_URL: postgresql://ci_user:${{ secrets.CI_DB_PASSWORD }}@localhost:5432/testdb
```

Even a CI-only, throwaway test database's credentials belong in GitHub
Secrets, not hardcoded in the workflow file — this builds the habit that
prevents a real credential from ever being hardcoded "just this once."

## Feeding CI failures back into the development loop

```
CI fails
    → copy the exact failure output
    → feed it back: "CI failed with this error: [paste]. Fix and verify
      locally (go test ./...) before pushing again."
    → fix pushed → CI re-runs
```

| Failure type | Typical fix |
|---|---|
| Lint | Run the linter's `--fix` locally, commit the result |
| Type/compile error | Read the exact file:line, fix the type |
| Test failure | Follow `test-driven-development`'s red-green loop — reproduce, then fix |
| Flaky test | Investigate root cause (`superpowers:systematic-debugging`-style), never just re-run to green |

## Feature flags decouple deploy from release

```go
func handleCheckout(flags FeatureFlags, userID string) http.HandlerFunc {
    if flags.Enabled("new-checkout-flow", userID) {
        return newCheckoutHandler
    }
    return legacyCheckoutHandler
}
```

Flags let you: ship code to `main` and production without turning it on,
roll back by flipping a flag instead of redeploying, canary a feature to
1% of users before 100%, and A/B test behavior directly. **Flags are not
permanent** — set a removal target when a flag is created; a flag still
live a year later with no plan to remove it is technical debt
(`deprecation-and-migration` covers cleaning these up).

## Build Cop

Designate a rotating owner responsible for keeping CI green. When `main`'s
build breaks, the Build Cop's job — not necessarily the author of the
breaking change — is to fix forward or revert immediately. Without this,
a broken build lingers while everyone assumes someone else is handling
it, and every other PR queues up behind it.

## Pipeline optimization, in order of impact

```
Slow pipeline (>10 min)?
├── Cache dependencies       → actions/setup-go's cache: true, or actions/cache for go build cache
├── Parallelize jobs         → lint/test/build as separate jobs, not sequential steps
├── Path-filter unrelated CI → skip integration/e2e jobs for docs-only changes
├── Shard the test suite     → matrix build across packages/test files
└── Move slow tests off the critical path → run nightly/scheduled instead of per-PR
```

```yaml
jobs:
  test:
    strategy:
      matrix:
        pkg: ["./internal/reservation/...", "./internal/billing/...", "./internal/guest/..."]
    steps:
      - run: go test ${{ matrix.pkg }} -race
```

## Secrets scoping

```
.env.example        → committed (template only, no real values)
.env                → never committed
CI secrets           → GitHub Secrets, scoped to the environment (never prod creds in a PR-triggered job)
Production secrets   → deployment platform's own secret store / vault
```

CI should never hold production secrets in a workflow triggered by an
untrusted PR (e.g. `pull_request` from a fork) — use `pull_request_target`
only with extreme care, or gate production-secret-bearing jobs to
`push` on protected branches.

## Gotchas

- `npm audit`/`govulncheck`/equivalent scanners catching a vulnerability
  in a transitive dependency with no available fix yet is a real
  operational decision (accept known risk with a tracked follow-up vs.
  block the pipeline) — don't silently downgrade the audit level to make
  it pass without that decision being made and recorded.
- A workflow trigger on `pull_request_target` runs with the base
  repository's secrets even for a fork's PR — using it to run untrusted
  PR code (rather than just to comment on a PR) is a supply-chain risk,
  not a convenience.
- Caching `go build`/`node_modules` across workflow runs can serve a
  stale cache after a dependency version bump if the cache key doesn't
  include the lockfile hash — key the cache on `go.sum`/`package-lock.json`
  hash, not just the branch name.
- A pipeline with no path filtering runs the full E2E suite on a
  documentation-only PR — harmless but slow, and it trains contributors
  to ignore CI wait times rather than investigate whether the pipeline
  itself needs the optimization steps above.

## Real-world grounding

The 2021 Codecov supply-chain compromise is a concrete, well-documented
case for scoping CI secrets carefully: attackers modified Codecov's
bash-uploader script (widely referenced directly from CI pipelines
across thousands of organizations) to exfiltrate environment variables
— including CI secrets — from every pipeline that ran it, for roughly
two months before detection. The generalizable lesson: a third-party
script pulled into CI at run time has the same access as your own
pipeline code, so any secret exposed to a CI job should be scoped to
only what that specific job needs, not to the full production secret
set "for convenience."

## Verification

- [ ] Lint, test, and build run as separate parallel jobs, not sequential steps
- [ ] No quality gate is skipped or silenced to force a merge
- [ ] CI secrets are scoped separately from production secrets, and never exposed to untrusted PR-triggered jobs
- [ ] Failing CI feeds back a specific, actionable error, not just a red X
- [ ] Feature flags used for risky changes have a tracked removal target
- [ ] Pipeline runtime is actively managed (caching, parallelism, path filters) rather than left to grow unchecked
