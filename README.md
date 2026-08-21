# skillme

[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Skill Eval](https://github.com/shennawardana23/skillme/actions/workflows/skill-eval.yml/badge.svg)](https://github.com/shennawardana23/skillme/actions/workflows/skill-eval.yml)
[![Test Plugin Installation](https://github.com/shennawardana23/skillme/actions/workflows/test-plugin-install.yml/badge.svg)](https://github.com/shennawardana23/skillme/actions/workflows/test-plugin-install.yml)
[![Agent Skills spec](https://img.shields.io/badge/spec-agentskills.io-blue)](https://agentskills.io/specification)

**A testable knowledge base of real engineering practice — not a collection of advice nobody has checked.**

skillme is 132 skills spanning backend, frontend, databases, web quality,
domain modeling, QA, product management, and agentic engineering itself
— each one grounded in an official doc, a real codebase, or a documented
incident rather than general impression, and each one backed by its own
deterministic test suite so a claim survives being checked against a real
model's output instead of being trusted because it reads well. See
[What's covered](#whats-covered) for the actual topic list and
[Coverage](#coverage) for exactly how many have been checked against a
live model run, not just schema-validated, as of the last full sweep.

It's packaged as a Claude Code plugin so any AI coding assistant that
supports the [Agent Skills](https://agentskills.io/specification) format
can load it, but the content itself is meant to read like something a
senior engineer would hand a new hire — not tool-specific instructions,
and not tied to any one company or codebase.

---

## Table of contents

- [skillme](#skillme)
  - [Table of contents](#table-of-contents)
  - [Quick Start](#quick-start)
  - [Why this exists](#why-this-exists)
  - [Coverage](#coverage)
  - [What's covered](#whats-covered)
  - [Specialist review skills](#specialist-review-skills)
  - [How a skill gets proven, not just written](#how-a-skill-gets-proven-not-just-written)
  - [Anatomy of a skill](#anatomy-of-a-skill)
  - [Requirements](#requirements)
  - [Install](#install)
    - [Claude Code plugin (managed bundle)](#claude-code-plugin-managed-bundle)
    - [Any agent, as editable files (via skills.sh)](#any-agent-as-editable-files-via-skillssh)
    - [Installing just one skill, not the whole catalog](#installing-just-one-skill-not-the-whole-catalog)
  - [Project Structure](#project-structure)
  - [Evaluating a skill locally](#evaluating-a-skill-locally)
    - [Provider/model fallback](#providermodel-fallback)
  - [Documentation](#documentation)
  - [Contributing](#contributing)
  - [License](#license)

---

## Quick Start

Fastest path in, any agent — no clone, no registration:

```bash
npx skills@latest add shennawardana23/skillme
```

<details>
<summary><b>Claude Code (recommended)</b></summary>

**Marketplace install:**

```bash
claude plugin marketplace add shennawardana23/skillme
claude plugin install skillme@skillme
```

**Local / development — clone first, point Claude Code at the path:**

```bash
git clone https://github.com/shennawardana23/skillme.git
claude --plugin-dir /path/to/skillme
```

</details>

<details>
<summary><b>Just one skill, not the whole catalog</b></summary>

```bash
npx skills@latest add shennawardana23/skillme --skill go-service-idioms
```

</details>

Full comparison of the two install paths — managed plugin bundle vs.
editable files via `skills.sh` — is in [Install](#install) below.

<details>
<summary><b>Prove a skill actually works, don't just install it</b></summary>

```bash
go build -o smeval ./cmd/smeval
./smeval run skills/go-service-idioms -benchmark
```

Then open, for that run's `smeval-workspace/runs/go-service-idioms/iteration-1/`:

- `report.html` — pass/fail per case, with the quoted assertion evidence
- `feedback.json` — a stub for a human to record what an assertion can't
  catch (see [`TESTING.md`](TESTING.md#reviewing-results-with-a-human))
- `benchmark.json` — with-skill vs. without-skill, side by side; see
  [`examples/benchmark.sample.json`](examples/benchmark.sample.json) for
  what real output looks like and how to read it

</details>

---

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

---

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

---

## What's covered

Every name below is a real, linked skill — not a topic label. Each area
lists a handful out of the full 132; run `scripts/list-skills.sh` for the
complete, current list.

| Area | Skills you'd reach for |
| --- | --- |
| Go | [`go-service-idioms`](skills/go-service-idioms/SKILL.md), [`golang-testing`](skills/golang-testing/SKILL.md), [`gin-lambda-api-service-patterns`](skills/gin-lambda-api-service-patterns/SKILL.md), [`genkit-go-flows`](skills/genkit-go-flows/SKILL.md), [`jwt-tenant-scoped-authorization`](skills/jwt-tenant-scoped-authorization/SKILL.md) |
| PHP | [`laravel-patterns`](skills/laravel-patterns/SKILL.md), [`laravel-security`](skills/laravel-security/SKILL.md), [`php-codeigniter-patterns`](skills/php-codeigniter-patterns/SKILL.md), [`php-codeigniter-legacy-patterns`](skills/php-codeigniter-legacy-patterns/SKILL.md), [`php-codeigniter-security`](skills/php-codeigniter-security/SKILL.md) |
| Databases | [`postgres-patterns`](skills/postgres-patterns/SKILL.md), [`postgres-hotel-partitioning`](skills/postgres-hotel-partitioning/SKILL.md), [`mysql-patterns`](skills/mysql-patterns/SKILL.md), [`database-migrations`](skills/database-migrations/SKILL.md) |
| Frontend | [`vue-nuxt-frontend-patterns`](skills/vue-nuxt-frontend-patterns/SKILL.md), [`frontend-ui-engineering`](skills/frontend-ui-engineering/SKILL.md), [`frontend-patterns`](skills/frontend-patterns/SKILL.md) |
| Web quality | [`web-quality-audit`](skills/web-quality-audit/SKILL.md), [`core-web-vitals`](skills/core-web-vitals/SKILL.md), [`accessibility`](skills/accessibility/SKILL.md), [`seo`](skills/seo/SKILL.md), [`best-practices`](skills/best-practices/SKILL.md) |
| Hospitality domain | [`hotel-rate-and-inventory-modeling`](skills/hotel-rate-and-inventory-modeling/SKILL.md), [`channel-manager-ota-integration`](skills/channel-manager-ota-integration/SKILL.md) |
| Agentic engineering | [`multi-agent-orchestration`](skills/multi-agent-orchestration/SKILL.md), [`agent-drift-and-re-anchoring`](skills/agent-drift-and-re-anchoring/SKILL.md), [`skill-catalog-authoring`](skills/skill-catalog-authoring/SKILL.md), [`skill-inspector`](skills/skill-inspector/SKILL.md), [`mcp-server-patterns`](skills/mcp-server-patterns/SKILL.md) |
| QA & testing | [`qa-test-strategy-design`](skills/qa-test-strategy-design/SKILL.md), [`test-driven-development`](skills/test-driven-development/SKILL.md), [`e2e-testing`](skills/e2e-testing/SKILL.md), [`ai-regression-testing`](skills/ai-regression-testing/SKILL.md) |
| Process & delivery | [`code-review-and-quality`](skills/code-review-and-quality/SKILL.md), [`ci-cd-and-automation`](skills/ci-cd-and-automation/SKILL.md), [`incident-response-and-postmortems`](skills/incident-response-and-postmortems/SKILL.md), [`release-management-and-rollback-planning`](skills/release-management-and-rollback-planning/SKILL.md) |
| Product & leadership | [`writing-product-requirements`](skills/writing-product-requirements/SKILL.md), [`feature-prioritization-frameworks`](skills/feature-prioritization-frameworks/SKILL.md), [`one-on-ones-and-feedback-frameworks`](skills/one-on-ones-and-feedback-frameworks/SKILL.md) |
| Policy & philosophy | [`error-handling-philosophy`](skills/error-handling-philosophy/SKILL.md), [`dependency-and-license-policy`](skills/dependency-and-license-policy/SKILL.md), [`backward-compatibility-and-deprecation-policy`](skills/backward-compatibility-and-deprecation-policy/SKILL.md) |

---

## Specialist review skills

Six skills in the catalog exist specifically to review something —
another skill, a diff, a config, a page — rather than to teach a pattern.
They intentionally divide the same general territory by *what's being
reviewed*, not by vague overlapping scope:

| Skill | What it does | Use when |
| --- | --- | --- |
| [`skill-inspector`](skills/skill-inspector/SKILL.md) | Runs NVIDIA SkillSpector plus a manual source read on a third-party agent skill, then returns an APPROVE / CAUTION / REJECT verdict | You found a skill in someone else's repo and want to know if it's safe before installing it |
| [`security-review`](skills/security-review/SKILL.md) | Reads application source (Go, TS, PHP, infra config) for auth gaps, unvalidated input, file-upload holes, leaked secrets, unsafe third-party calls | You're reviewing a PR or an existing codebase for security bugs |
| [`security-and-hardening`](skills/security-and-hardening/SKILL.md) | Threat-models trust boundaries and designs OWASP Top 10 prevention into a feature | You're designing or building a feature, before the vulnerability exists to find |
| [`security-scan`](skills/security-scan/SKILL.md) | Grades your own agent harness config — `CLAUDE.md`, `.claude/settings.json`, MCP configs, hooks — A through F for prompt-injection exposure | You want to know how exposed your own Claude Code setup is, not your application |
| [`code-review-and-quality`](skills/code-review-and-quality/SKILL.md) | Sizes a diff, applies severity labels, and can help design a team's review process | You're reviewing a pull request, or defining how your team should |
| [`web-quality-audit`](skills/web-quality-audit/SKILL.md) | Runs performance, accessibility, SEO, and best-practices checks and returns one prioritized report | You need a single pass across all four web-quality dimensions instead of four separate audits |

These aren't a separate "agent" layer sitting outside the catalog — each
is an ordinary skill with its own `SKILL.md` and `evals/evals.json`, like
every other entry. `security-review`/`security-and-hardening`/
`security-scan` in particular read as near-duplicates by name alone; each
one's own description states the exact dividing line against its
siblings, which is the pattern to follow if you add a skill that sounds
close to an existing one.

---

## How a skill gets proven, not just written

```
skills/<name>/SKILL.md + evals/evals.json      ← prompt + assertions, written by hand
              │
              ▼
        smeval validate                        ← schema check, free, no model call
              │
              ▼
        smeval run                             ← spawns a fresh, isolated `claude` CLI
              │                                    process (only this skill visible)
              ▼
        grading.json + report.html + feedback.json
        quoted evidence + human review — never a bare pass/fail
```

No claim in this catalog ships on "seems right." `smeval` — this
repository's own Go eval runner, no third-party dependency — validates the
schema for free, then actually runs each case against the real `claude`
CLI in a freshly-isolated workspace (the skill under test is the *only*
skill visible) and grades the response with concrete, quoted evidence.
See [How it works](#evaluating-a-skill-locally) below to run this
yourself, and [`TESTING.md`](TESTING.md) for the full methodology,
including a real workspace-isolation bug this exact loop caught and fixed
in itself.

---

## Anatomy of a skill

```
skills/go-service-idioms/
├── SKILL.md              ┌─ frontmatter ─────────────────────────┐
│                         │ name: go-service-idioms               │
│                         │ description: Use when asked to "write │
│                         │   a Go function"... (trigger phrases) │
│                         │ license: Apache-2.0                   │
│                         │ metadata: { version, category }       │
│                         └───────────────────────────────────────┘
│                         Body: the actual guidance, written for the
│                         agent executing it — imperative, grounded in
│                         a real source, ending in Gotchas +
│                         Real-world grounding sections.
├── references/           Optional detail, loaded only when the body
│                         actually points to it — keeps the always-in-
│                         context body short (progressive disclosure).
└── evals/evals.json      This skill's own test suite — see above.
```

Every claim of the form "do X because Y" in a `SKILL.md` body should
trace to something real: an official doc, a documented incident, a
verified library/API. `skill-catalog-authoring`
(`skills/skill-catalog-authoring/SKILL.md`) is the actual, enforced
specification for all of this — required fields, naming rules, the
optional `skill-docs/<category>/<name>.md` human-facing page convention —
this section is the map, not the source of truth.

---

## Requirements

- Go 1.25+ (to build `smeval`)
- [`uv`](https://docs.astral.sh/uv/) (`uvx`) — runs the official
  `agentskills` spec validator without a separate install
- The [`claude` CLI](https://github.com/anthropics/claude-code), logged in
  or with `ANTHROPIC_API_KEY` set — only needed to actually *run* an eval
  against a live model; schema validation works without it

---

## Install

Two ways in, two philosophies. Pick one — installing both leaves you with
every skill twice.

| | Claude Code plugin | `skills.sh` |
| --- | --- | --- |
| Gets you | The whole catalog, one managed bundle | Just the skill(s) you pick |
| Editable? | No — read-only, reinstall to update | Yes — plain files in your project |
| Works with | Claude Code only | Claude Code, Codex, Cursor, others |
| Registration needed | Yes, this repo's own marketplace | No |

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

---

## Project Structure

```text
skillme/
├── skills/<name>/
│   ├── SKILL.md              # the skill itself — written for the agent
│   ├── references/           # optional detail, loaded only when needed
│   └── evals/evals.json      # prompt + assertions, this catalog's test suite
├── skill-docs/<category>/<name>.md   # optional, human-facing "why/when" page
├── examples/benchmark.sample.json    # real -benchmark output, annotated in TESTING.md
├── cmd/smeval/                # the eval runner CLI
├── internal/{evalspec,engine,grading,harness,report}/   # smeval's implementation
├── .github/workflows/
│   ├── skill-eval.yml            # CI: build/vet/test, then schema-validate every skill
│   └── test-plugin-install.yml   # CI: claude plugin validate + a real marketplace-add/install
├── CLAUDE.md + .claude/{commands,rules}/   # eval-authoring rules + /eval-skill, /new-skill
└── .claude-plugin/{plugin.json,marketplace.json}
```

---

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

---

## Documentation

| Doc | Covers |
| --- | --- |
| This README | Install, what's covered, coverage honesty, repository layout |
| [`TESTING.md`](TESTING.md) | The full eval methodology — schema vs. live, reading a failure's evidence, `-benchmark`, human review (`feedback.json`), real bugs this catalog's own testing loop has found and fixed in itself |
| [`CONTRIBUTING.md`](CONTRIBUTING.md) | Skill/eval-authoring conventions, the pre-PR checklist, what CI does and doesn't check |
| `skills/skill-catalog-authoring/SKILL.md` | The enforced spec for a skill's directory layout and frontmatter — the source of truth [Anatomy of a skill](#anatomy-of-a-skill) above summarizes |
| `skill-docs/<category>/<name>.md` | Optional, per-skill human-facing "why/when would I reach for this" pages — not every skill has one; see `skill-catalog-authoring/references/skill-docs-template.md` for the template |
| [`CLAUDE.md`](CLAUDE.md) + [`.claude/rules/skill-authoring.md`](.claude/rules/skill-authoring.md) | The eval-authoring rules Claude Code loads automatically in this repo |
| [`.claude/commands/eval-skill.md`](.claude/commands/eval-skill.md), [`.claude/commands/new-skill.md`](.claude/commands/new-skill.md) | `/eval-skill <name>` and `/new-skill <name>` — run this catalog's own eval loop or scaffold a new entry without leaving the chat |

---

## Contributing

New skills, fixes to existing ones, and improvements to `smeval` itself
are welcome — see [`CONTRIBUTING.md`](CONTRIBUTING.md) for the required
directory layout and frontmatter conventions, the eval-authoring mistakes
this catalog has actually made (and how to avoid repeating them), the
pre-PR checklist, and what CI does and doesn't check.

---

## License

Apache-2.0 — see [LICENSE](LICENSE).
