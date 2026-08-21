# skillme

[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Skill Eval](https://github.com/shennawardana23/skillme/actions/workflows/skill-eval.yml/badge.svg)](https://github.com/shennawardana23/skillme/actions/workflows/skill-eval.yml)

A written-down, testable knowledge base of real engineering practice — Go,
PHP (Laravel, CodeIgniter, legacy CodeIgniter), PostgreSQL/MySQL, Vue/Nuxt,
web quality (Core Web Vitals, accessibility, SEO), hotel-domain modeling,
QA, product management, agentic engineering, and policy — plus
`skill-inspector`, which reviews a *third-party* agent skill for safety
before you install it. 132 topics, drawn from official docs, real
codebases, and documented incidents rather than general impression — see
[Coverage](#coverage) below for exactly how many have been checked against
a live model run, not just schema-validated, as of the last full sweep.

It's packaged as a Claude Code plugin so any AI coding assistant that
supports the [Agent Skills](https://agentskills.io/specification) format
can load it, but the content itself is meant to read like something a
senior engineer would hand a new hire — not tool-specific instructions,
and not tied to any one company or codebase.

## Why this exists

Most "prompt libraries" are collections of good-sounding advice nobody has
checked. Every topic here ships with a small, deterministic test suite so
a claim either survives being checked against a real model's output or
gets fixed — not trusted because it reads well. That checking is done by
[`smeval`](#evaluating-a-skill-locally), this repository's own Go eval
runner, following the same evaluation methodology
[Anthropic documents for Agent Skills](https://agentskills.io/skill-creation/evaluating-skills):
a prompt plus assertions per case, graded with concrete, quoted evidence,
not a vague pass/fail.

## Coverage

"Schema-valid" and "actually works against a real model" are different
claims — see [`TESTING.md`](TESTING.md). This repo tracks the second one
honestly rather than asserting it catalog-wide:

- Every skill passes schema validation (`agentskills validate` +
  `smeval validate`) — enforced in CI on every push.
- A rolling live-test sweep (`scripts/test-all-skills.sh`) checks each
  skill against a real model run. Its accumulating, always-current result
  is `smeval-workspace/test-all-results/summary.txt` — read that file for
  the exact count as of right now, not a number here that will drift out
  of date.
- The original 125-skill catalog completed a full live-test sweep: every
  partial failure found was individually diagnosed against the model's
  actual saved output before being touched. Every one turned out to be an
  eval-authoring bug (a grading assertion too narrow to match a
  correct-but-differently-phrased answer, a markdown backtick breaking a
  literal substring match, a blunt negative check banning a phrase the
  model used correctly) — not an actual skill-content defect — and was
  fixed and reverified against the model's saved output before being
  counted as passing. A couple of failures turned out to be ordinary
  single-run model variance (confirmed clean on an immediate re-run)
  rather than anything requiring a fix. One failure was a real bug in
  `smeval` itself, not the skill or the eval — a case could write files
  into this actual repo instead of its isolated workspace; see
  `TESTING.md`'s "A fixed workspace-escape bug" for the full story, now
  fixed and verified.
- `skill-inspector` and the 6 web-quality skills were added after that
  sweep and are live-verified individually (their own eval cases run and
  passing), but haven't yet gone through a full-catalog sweep pass of
  their own alongside everything else. Read
  `smeval-workspace/test-all-results/summary.txt` for the exact,
  currently-accumulated count rather than trusting a number here — don't
  read "132 skills" as "132 independently proven."

This is the same rigor bar this catalog's own eval-quality bug reports
apply to any other skill's claims — including a real, previously-hidden
test-isolation bug in `smeval` itself, caught only because a benchmark run
showed a suspiciously "no difference" result and was investigated rather
than accepted (see `internal/engine/engine.go`'s `--setting-sources`
comment for the full story).

## What's covered

| Area | Examples |
| --- | --- |
| Go | idioms, testing, service patterns, Gin+Lambda deployment, Genkit |
| PHP | Laravel, CodeIgniter 4, legacy CodeIgniter 2/3, security, TDD |
| Databases | PostgreSQL patterns, partitioning, MySQL/MariaDB, migrations |
| Frontend | Vue/Nuxt, frontend engineering patterns |
| Web quality | Core Web Vitals (LCP/INP/CLS), accessibility (WCAG 2.2), SEO, security/best-practices audits |
| Hospitality domain | rate/inventory modeling, channel manager/OTA integration |
| Agentic engineering | multi-agent orchestration, agent drift, eval design, third-party skill safety review (`skill-inspector`) |
| Process | QA strategy, TDD, code review, security review, CI/CD |
| Product & leadership | PR-FAQs, prioritization, feedback, incident response |
| Policy & philosophy | error handling, dependency policy, deprecation |

Run `scripts/list-skills.sh` for the full, current, exact list — the table
above is illustrative, not exhaustive.

## Requirements

- Go 1.25+ (to build `smeval`)
- [`uv`](https://docs.astral.sh/uv/) (`uvx`) — runs the official
  `agentskills` spec validator without a separate install
- The [`claude` CLI](https://github.com/anthropics/claude-code), logged in
  or with `ANTHROPIC_API_KEY` set — only needed to actually *run* an eval
  against a live model; schema validation works without it

## Install

Two ways in, two philosophies. The Claude Code plugin installs the whole
catalog as a managed, read-only bundle — you `git pull` this repo and
reinstall to pick up changes, rather than editing the skills in place.
The `skills.sh` installer copies individual skill files straight into your
own project, editable, across whichever coding agent you use. Pick one —
installing both leaves you with every skill twice.

### Claude Code plugin (managed bundle)

```bash
# Straight from GitHub — no clone needed:
claude plugin marketplace add shennawardana23/skillme
claude plugin install skillme@skillme

# Already have it cloned locally? Point at the path instead:
claude plugin marketplace add /path/to/skillme
claude plugin install skillme@skillme

# Or, for a quick local session with no marketplace registration:
claude --plugin-dir /path/to/skillme
```

`skillme` isn't (yet) in Anthropic's official plugin marketplace, so
`claude plugin install skillme` on its own won't resolve it — `marketplace
add` above is what registers this repo's own `.claude-plugin/marketplace.json`
(`source: "./"`) so `install` has something to find. To pick up later
changes: `claude plugin update skillme` (restart the session to apply).

### Any agent, as editable files (via skills.sh)

```bash
npx skills@latest add shennawardana23/skillme
```

[`skills.sh`](https://www.skills.sh/) works against any repo with a
`skills/<name>/SKILL.md` layout by convention — no registration needed on
this repo's end. It writes the skills you pick into your own project as
ordinary, editable files, for whichever agent(s) you tell it to target
(Claude Code, Codex, Cursor, and others). Nothing updates behind your
back; pull this repo's later changes yourself with `npx skills@latest
update` when you want them.

### Installing just one skill, not the whole catalog

The Claude Code plugin route above is all-or-nothing — a plugin installs
as one managed bundle, there's no per-skill flag. `skills.sh` is the one
that supports picking a single skill out of this catalog:

```bash
npx skills@latest add shennawardana23/skillme --skill go-service-idioms

# more than one, still without the rest of the catalog:
npx skills@latest add shennawardana23/skillme --skill go-service-idioms security-review

# not sure of the exact name yet:
npx skills@latest add shennawardana23/skillme --list
```

`--list` prints every skill name in this repo without installing
anything — run `scripts/list-skills.sh` locally for the same list, or
browse `skills/` directly. Everything installed this way lands as plain,
editable files in your own project; update just that skill later with
`npx skills@latest update go-service-idioms`.

## Repository layout

```text
skillme/
├── skills/<name>/
│   ├── SKILL.md              # the skill itself — written for the agent
│   ├── references/           # optional detail, loaded only when needed
│   └── evals/evals.json      # prompt + assertions, this catalog's test suite
├── skill-docs/<category>/<name>.md   # optional, human-facing "why/when" page
├── cmd/smeval/                # the eval runner CLI
├── internal/{evalspec,engine,grading,harness,report}/   # smeval's implementation
├── .github/workflows/skill-eval.yml   # CI: build/vet/test, then validate + run every skill
└── .claude-plugin/{plugin.json,marketplace.json}
```

## Evaluating a skill locally

```bash
# Official Agent Skills spec compliance (frontmatter, naming) — free
uvx --from skills-ref agentskills validate skills/go-service-idioms

# This catalog's own eval schema and cases
go build -o smeval ./cmd/smeval

smeval validate skills/go-service-idioms            # free, no model calls — schema only
smeval run      skills/go-service-idioms            # runs cases against the local `claude` CLI
smeval run      skills/go-service-idioms -benchmark # also runs without the skill, for comparison
```

`smeval validate` only proves the JSON is well-formed — it catches schema
errors after editing `evals.json`, nothing more. `smeval run` is the one
that actually proves a skill works: it executes each case against the
local `claude` CLI (headless, `--output-format json`) and grades the real
output. Results land in `smeval-workspace/runs/<name>/iteration-N/`, including a
styled `report.html`. See [`TESTING.md`](TESTING.md) for how to read a
failure's evidence, iterate on a single case, and interpret a
`-benchmark` comparison.

### Provider/model fallback

`smeval run` accepts `-primary-model` (default `sonnet`) and
`-fallback-model` (default `opus`). Fallback happens in two layers:
`claude`'s own native `--fallback-model` handles a model being overloaded
or unavailable within one invocation; `smeval`'s own outer retry (in
`internal/engine`) catches whole-process failures the native layer can't
see (non-zero exit, timeout, unparseable output) and retries once in a
fresh process against the fallback model. Both are exercised by real tests
in `internal/engine/engine_test.go`.

## Contributing

New skills and improvements to existing ones are welcome. Start with the
`skill-catalog-authoring` skill (`skills/skill-catalog-authoring/SKILL.md`)
for the required directory layout, frontmatter conventions, and
`evals.json` schema — it also documents the optional `skill-docs/` page
convention (`skill-catalog-authoring/references/skill-docs-template.md`).

Before opening a PR:

1. `go build ./... && go vet ./... && gofmt -l .`
2. `uvx --from skills-ref agentskills validate skills/<name>` and
   `smeval validate skills/<name>` for anything you added or changed
3. `smeval run skills/<name>` — **do this yourself, it's not a CI gate.**
   A passing schema is not the same as a working skill; CI only checks
   spec compliance and schema (no API key, no live model calls, on
   purpose — see below). Actually running the eval against a real model
   is the step that catches real bugs, and it's on you to do it before
   opening the PR.

CI (`.github/workflows/skill-eval.yml`) discovers any skill directory
under `skills/` that contains an `evals/evals.json` automatically — no
workflow edits are needed to pick up a new skill. It deliberately stops at
`agentskills validate` + `smeval validate` (schema/spec only): re-running
every skill's live eval on every push would mean real API cost and
significant runtime at catalog scale for even a one-line change to a
single skill.

## License

Apache-2.0 — see [LICENSE](LICENSE).
