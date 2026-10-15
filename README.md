# skillme ░░░▒▒▒░░░▒▒▒

[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Skill Eval](https://github.com/shennawardana23/skillme/actions/workflows/skill-eval.yml/badge.svg)](https://github.com/shennawardana23/skillme/actions/workflows/skill-eval.yml)
[![Test Plugin Installation](https://github.com/shennawardana23/skillme/actions/workflows/test-plugin-install.yml/badge.svg)](https://github.com/shennawardana23/skillme/actions/workflows/test-plugin-install.yml)
[![Agent Skills spec](https://img.shields.io/badge/spec-agentskills.io-blue)](https://agentskills.io/specification)

**A tested catalog of engineering skills for AI coding agents.** 143 skills
(Go, PHP, databases, security, testing, web quality, agentic engineering,
product and more), each one a plain `SKILL.md` folder in the open
[Agent Skills](https://agentskills.io) format, so it works in Claude Code,
Codex, Cursor, GitHub Copilot, Gemini CLI, OpenCode and any other tool that
reads that format. What sets it apart: every skill ships its own eval suite,
and `smeval`, this repo's own runner, checks the skill against a real model
instead of trusting that it reads well.

| | |
| --- | --- |
| **Skills** | 143, in 16 areas ([full catalog](#the-catalog)) |
| **Format** | Agent Skills (`skills/<name>/SKILL.md`), progressive disclosure |
| **Works with** | Any Agent Skills client ([install per tool](#install-it-in-your-ai-tool)); plain Markdown for everything else |
| **Quality gate** | Per-skill `evals/evals.json`, graded by [`smeval`](#how-a-skill-is-proven) with quoted evidence, no LLM judge |
| **License** | Apache-2.0 |

**Jump to:** [Quick start](#quick-start) ·
[Install per tool](#install-it-in-your-ai-tool) ·
[Use a skill](#use-a-skill) ·
[Catalog](#the-catalog) ·
[How a skill is proven](#how-a-skill-is-proven) ·
[For contributors](#for-contributors) ·
[Decisions (ADRs)](#architecture-decisions) ·
[Troubleshooting](#troubleshooting)

---

## Who it is for

- **Developers using an AI coding agent** who want it to apply a vetted
  pattern (a Go context rule, a Postgres indexing rule, a security review
  checklist) instead of re-explaining house rules every session.
- **Teams** that want the same engineering guidance in every repo and every
  tool, from one source.
- **People building their own skill catalog**: the eval runner, isolation,
  CI split and catalog-wide checks here are working code you can read or fork.

**Not what this is**: a general-purpose prompt library or a creative-writing
collection. Every entry is an engineering pattern with a runnable test behind
it, graded on whether it is true and specific, not on how persuasive it reads.
It also does not replace the framework or library docs a skill is based on:
each skill ends with a "Real-world grounding" section naming its primary source.

---

## Why this exists

Most prompt libraries are good-sounding advice nobody has checked. Every topic
here ships with a small, deterministic test suite, so a claim either survives
being checked against a real model's output or gets fixed. The checking is done
by [`smeval`](#evaluating-a-skill-locally), this repository's own Go eval runner,
following the method [documented for Agent Skills](https://agentskills.io/skill-creation/evaluating-skills):
a prompt plus assertions per case, graded with concrete, quoted evidence rather
than a bare pass or fail. The content is written the way a senior engineer
would brief a new hire: not tool-specific, and not tied to one company.

---

## Quick start

One command, any supported agent, no clone:

```bash
npx skills@latest add shennawardana23/skillme              # pick skills and agents interactively
npx skills@latest add shennawardana23/skillme --list       # list skill names, install nothing
npx skills@latest add shennawardana23/skillme --skill go-service-idioms -a claude-code
```

Using Claude Code? Install the whole catalog as a managed plugin instead:

```bash
claude plugin marketplace add shennawardana23/skillme
claude plugin install skillme@skillme
```

Then ask your agent for a task the skill covers ("review this Go handler",
"design the rate plan tables"). It loads the matching skill on its own. See
[Use a skill](#use-a-skill) to call one by name.

> Pick **one** install method per agent. Installing the plugin and the
> `npx skills` copy into the same agent loads every skill twice.

---

## Install it in your AI tool

### Two ways, two trade-offs

| | Claude Code plugin | `npx skills` (skills.sh) | Copy or link the folders |
| --- | --- | --- | --- |
| You get | The whole catalog, one bundle | Only the skills you pick | Only the skills you pick |
| Files | Managed, read-only | Plain files in your project or home | Plain files |
| Updates | `claude plugin update skillme` | `npx skills@latest update` | `git pull` (if you linked) |
| Tools | Claude Code | About 70 agents | Any tool that reads a skills folder |

### Per tool

Folder paths below come from each tool's own documentation (checked October
2026). `npx skills ... -a <id>` writes to the same folders for you.

<details>
<summary><b>Claude Code</b> (plugin, recommended)</summary>

```bash
# Straight from GitHub
claude plugin marketplace add shennawardana23/skillme
claude plugin install skillme@skillme

# From a local clone (development), no marketplace needed
git clone https://github.com/shennawardana23/skillme.git
claude --plugin-dir ./skillme
```

Single skills without the plugin: `npx skills@latest add shennawardana23/skillme --skill <name> -a claude-code`,
or copy `skills/<name>/` into `~/.claude/skills/` (all projects) or `.claude/skills/` (one repo).
Plugin skills are called as `/skillme:<name>`. Restart the session after `claude plugin update skillme`.
</details>

<details>
<summary><b>OpenAI Codex</b> (CLI, IDE, app)</summary>

Codex reads `.agents/skills/` (from the working directory up to the repo root) and `~/.agents/skills/`.

```bash
npx skills@latest add shennawardana23/skillme -a codex        # this project
npx skills@latest add shennawardana23/skillme -a codex -g     # all projects
```

Restart Codex if a new skill does not appear. Call one explicitly with `$<name>`.
</details>

<details>
<summary><b>GitHub Copilot</b> (VS Code and JetBrains agent mode, Copilot CLI, coding agent, code review)</summary>

Copilot reads project skills from `.github/skills/`, `.claude/skills/` or `.agents/skills/`,
and personal skills from `~/.copilot/skills/` or `~/.agents/skills/`.

```bash
npx skills@latest add shennawardana23/skillme -a github-copilot
```

Commit `.github/skills/` (or `.agents/skills/`) to give the Copilot coding agent the same skills.
</details>

<details>
<summary><b>Cursor</b></summary>

Cursor reads `.agents/skills/` and `.cursor/skills/` in the project, `~/.agents/skills/` and
`~/.cursor/skills/` for the user, and also loads `.claude/skills/` and `.codex/skills/` for compatibility.

```bash
npx skills@latest add shennawardana23/skillme -a cursor
```

Only `~/.cursor/skills/` syncs to Cursor Cloud Agents.
</details>

<details>
<summary><b>Gemini CLI</b></summary>

Gemini CLI reads `.gemini/skills/` or `.agents/skills/` in the workspace and `~/.gemini/skills/` or `~/.agents/skills/` for the user.

```bash
gemini skills install https://github.com/shennawardana23/skillme.git --path skills             # whole catalog
gemini skills install https://github.com/shennawardana23/skillme.git --path skills/go-service-idioms --scope workspace
```

</details>

<details>
<summary><b>OpenCode, Windsurf, Kiro, Roo Code, Goose, Cline, Amp, Qwen Code, Trae and others</b></summary>

Use the agent id with `npx skills`:

```bash
npx skills@latest add shennawardana23/skillme -a opencode      # .agents/skills/ or ~/.config/opencode/skills/
npx skills@latest add shennawardana23/skillme -a windsurf      # .windsurf/skills/
npx skills@latest add shennawardana23/skillme -a kiro-cli      # .kiro/skills/
npx skills@latest add shennawardana23/skillme -a roo           # .roo/skills/
npx skills@latest add shennawardana23/skillme -a goose         # .goose/skills/
```

The full list of agent ids and folders is in the [skills CLI README](https://github.com/vercel-labs/skills).
OpenCode also reads `.claude/skills/`, so a Claude Code project install works there too.
</details>

<details>
<summary><b>Manual install for any tool that reads a skills folder</b></summary>

Most clients now read `~/.agents/skills/` (Codex, Cursor, Copilot, Gemini CLI, OpenCode).
Link the skills you want, so `git pull` keeps them current:

```bash
git clone https://github.com/shennawardana23/skillme.git ~/src/skillme
mkdir -p ~/.agents/skills
ln -s ~/src/skillme/skills/go-service-idioms ~/.agents/skills/go-service-idioms
# Claude Code reads ~/.claude/skills/ instead:
ln -s ~/src/skillme/skills/go-service-idioms ~/.claude/skills/go-service-idioms
```

The folder name must equal the `name` in the skill's frontmatter. Copy instead of link on systems without symlinks.
</details>

<details>
<summary><b>Claude apps and the Claude API</b></summary>

Claude (web, desktop, mobile) and the Claude API both support Agent Skills. Upload or attach a
skill folder as described in [Anthropic's Agent Skills docs](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview).
</details>

<details>
<summary><b>Tools without skills support</b> (ChatGPT chat, Gemini app, local models in Ollama or LM Studio, raw APIs)</summary>

A skill is plain Markdown, so you can use it as instructions:

1. Open `skills/<name>/SKILL.md` and paste the body (everything after the frontmatter) into the
   system prompt, custom instructions, project instructions or GPT/Gem instructions.
2. If the body points to a file under `references/`, paste that file only when the task needs it.
   That keeps the prompt small, which is what the skill format does for you automatically elsewhere.
3. Keep it to the few skills the task needs. You lose automatic selection, so you choose.

```bash
awk 'f>=2{print} /^---$/{f++}' skills/go-service-idioms/SKILL.md | pbcopy   # body only (macOS; xclip or clip.exe elsewhere)
```

</details>

---

## Use a skill

Skills load in three steps ([progressive disclosure](https://agentskills.io/specification#progressive-disclosure)):
at startup the agent reads only each skill's `name` and `description` (about 100 tokens each); when your
request matches a description, it reads that `SKILL.md`; it opens files under `references/` only when the
body points to them. So 143 installed skills cost little context until one is needed.

- **Automatically:** describe the task in plain words. Each description lists the phrases it reacts to,
  for example `go-service-idioms` reacts to "write a Go function" and "review this Go handler".
- **By name:** Claude Code plugin `/skillme:<name>`, Claude Code personal or project install `/<name>`,
  Codex `$<name>`, other tools: name the skill in your message ("use the security-review skill").
- **Not sure which one?** The [`orchestrator`](skills/orchestrator/SKILL.md) skill routes a request to the
  right skill, and [`using-agent-skills`](skills/using-agent-skills/SKILL.md) explains how to work with skills in general.

### Specialist review skills

Six skills exist to review something rather than to teach a pattern. They split the work by *what is
being reviewed*, not by vague overlapping scope:

| Skill | What it does | Use when |
| --- | --- | --- |
| [`skill-inspector`](skills/skill-inspector/SKILL.md) | Runs NVIDIA SkillSpector plus a manual source read on a third-party agent skill, then returns APPROVE / CAUTION / REJECT | You found a skill in someone else's repo and want to know if it is safe before installing it |
| [`security-review`](skills/security-review/SKILL.md) | Reads application source (Go, TS, PHP, infra config) for auth gaps, unvalidated input, upload holes, leaked secrets, unsafe third-party calls | Reviewing a PR or an existing codebase for security bugs |
| [`security-and-hardening`](skills/security-and-hardening/SKILL.md) | Threat-models trust boundaries and designs OWASP Top 10 prevention into a feature | Designing or building a feature, before the vulnerability exists |
| [`security-scan`](skills/security-scan/SKILL.md) | Grades your agent harness (`CLAUDE.md`, `.claude/settings.json`, MCP configs, hooks) A to F for prompt-injection exposure | Checking your own setup, not your application |
| [`code-review-and-quality`](skills/code-review-and-quality/SKILL.md) | Sizes a diff, applies severity labels, helps design a team's review process | Reviewing a pull request, or defining how your team reviews |
| [`web-quality-audit`](skills/web-quality-audit/SKILL.md) | Runs performance, accessibility, SEO and best-practices checks into one prioritized report | One pass across all four web-quality areas |

`security-review`, `security-and-hardening` and `security-scan` sound alike; each description states the
exact line against its siblings, which is the pattern to follow when you add a skill that sounds close to
an existing one.

---

## The catalog

<!-- generated: 143 skills in 16 areas -->
| Area | Skills | Count |
| --- | --- | ---: |
| **Go** | [`go-expert`](skills/go-expert/SKILL.md), [`go-service-idioms`](skills/go-service-idioms/SKILL.md), [`golang-patterns`](skills/golang-patterns/SKILL.md), [`golang-testing`](skills/golang-testing/SKILL.md), [`golang-concurrency`](skills/golang-concurrency/SKILL.md), [`golang-context`](skills/golang-context/SKILL.md), [`golang-database`](skills/golang-database/SKILL.md), [`golang-dependency-injection`](skills/golang-dependency-injection/SKILL.md), [`golang-grpc`](skills/golang-grpc/SKILL.md), [`gin-lambda-api-service-patterns`](skills/gin-lambda-api-service-patterns/SKILL.md), [`genkit-go-flows`](skills/genkit-go-flows/SKILL.md), [`adk-go-agent-builder`](skills/adk-go-agent-builder/SKILL.md), [`jwt-tenant-scoped-authorization`](skills/jwt-tenant-scoped-authorization/SKILL.md), [`layered-logging-and-alerting-patterns`](skills/layered-logging-and-alerting-patterns/SKILL.md), [`raw-sql-querybuilder-testing-patterns`](skills/raw-sql-querybuilder-testing-patterns/SKILL.md), [`scaffold`](skills/scaffold/SKILL.md) | 16 |
| **PHP (Laravel, CodeIgniter)** | [`laravel-patterns`](skills/laravel-patterns/SKILL.md), [`laravel-security`](skills/laravel-security/SKILL.md), [`laravel-tdd`](skills/laravel-tdd/SKILL.md), [`laravel-verification`](skills/laravel-verification/SKILL.md), [`php-codeigniter-patterns`](skills/php-codeigniter-patterns/SKILL.md), [`php-codeigniter-legacy-patterns`](skills/php-codeigniter-legacy-patterns/SKILL.md), [`php-codeigniter-legacy-auth-session`](skills/php-codeigniter-legacy-auth-session/SKILL.md), [`php-codeigniter-security`](skills/php-codeigniter-security/SKILL.md), [`php-codeigniter-tdd`](skills/php-codeigniter-tdd/SKILL.md), [`php-codeigniter-verification`](skills/php-codeigniter-verification/SKILL.md) | 10 |
| **Databases and data modeling** | [`postgres-patterns`](skills/postgres-patterns/SKILL.md), [`postgres-hotel-partitioning`](skills/postgres-hotel-partitioning/SKILL.md), [`mysql-patterns`](skills/mysql-patterns/SKILL.md), [`database-migrations`](skills/database-migrations/SKILL.md), [`erd-expert`](skills/erd-expert/SKILL.md), [`content-hash-cache-pattern`](skills/content-hash-cache-pattern/SKILL.md) | 6 |
| **Backend, APIs and operations** | [`api-design`](skills/api-design/SKILL.md), [`api-and-interface-design`](skills/api-and-interface-design/SKILL.md), [`backend-patterns`](skills/backend-patterns/SKILL.md), [`twelve-factor-app-principles`](skills/twelve-factor-app-principles/SKILL.md), [`idempotency-and-retry-safety`](skills/idempotency-and-retry-safety/SKILL.md), [`observability-and-instrumentation`](skills/observability-and-instrumentation/SKILL.md), [`docker-patterns`](skills/docker-patterns/SKILL.md), [`deployment-patterns`](skills/deployment-patterns/SKILL.md), [`mcp-server-patterns`](skills/mcp-server-patterns/SKILL.md) | 9 |
| **Frontend and mobile** | [`frontend-patterns`](skills/frontend-patterns/SKILL.md), [`frontend-ui-engineering`](skills/frontend-ui-engineering/SKILL.md), [`vue-nuxt-frontend-patterns`](skills/vue-nuxt-frontend-patterns/SKILL.md), [`compose-multiplatform-patterns`](skills/compose-multiplatform-patterns/SKILL.md), [`browser-testing-with-devtools`](skills/browser-testing-with-devtools/SKILL.md), [`frontend-slides`](skills/frontend-slides/SKILL.md) | 6 |
| **Web quality, performance and design** | [`web-quality-audit`](skills/web-quality-audit/SKILL.md), [`core-web-vitals`](skills/core-web-vitals/SKILL.md), [`performance`](skills/performance/SKILL.md), [`performance-optimization`](skills/performance-optimization/SKILL.md), [`accessibility`](skills/accessibility/SKILL.md), [`seo`](skills/seo/SKILL.md), [`best-practices`](skills/best-practices/SKILL.md), [`ai-design-slop-audit`](skills/ai-design-slop-audit/SKILL.md), [`animation-opportunity-scan`](skills/animation-opportunity-scan/SKILL.md), [`animation-improvement-planning`](skills/animation-improvement-planning/SKILL.md), [`animation-review`](skills/animation-review/SKILL.md) | 11 |
| **Security and privacy** | [`security-review`](skills/security-review/SKILL.md), [`security-and-hardening`](skills/security-and-hardening/SKILL.md), [`security-scan`](skills/security-scan/SKILL.md), [`security-analysis`](skills/security-analysis/SKILL.md), [`skill-inspector`](skills/skill-inspector/SKILL.md), [`data-retention-and-privacy-by-design`](skills/data-retention-and-privacy-by-design/SKILL.md), [`dependency-and-license-policy`](skills/dependency-and-license-policy/SKILL.md) | 7 |
| **Testing and QA** | [`qa-test-strategy-design`](skills/qa-test-strategy-design/SKILL.md), [`test-driven-development`](skills/test-driven-development/SKILL.md), [`tdd-workflow`](skills/tdd-workflow/SKILL.md), [`test-generation`](skills/test-generation/SKILL.md), [`e2e-testing`](skills/e2e-testing/SKILL.md), [`exploratory-testing-techniques`](skills/exploratory-testing-techniques/SKILL.md), [`ai-regression-testing`](skills/ai-regression-testing/SKILL.md), [`eval-harness`](skills/eval-harness/SKILL.md), [`bug-triage-and-severity-classification`](skills/bug-triage-and-severity-classification/SKILL.md) | 9 |
| **Code review and quality** | [`code-review`](skills/code-review/SKILL.md), [`code-review-and-quality`](skills/code-review-and-quality/SKILL.md), [`code-reviewer`](skills/code-reviewer/SKILL.md), [`diff-analysis`](skills/diff-analysis/SKILL.md), [`code-simplification`](skills/code-simplification/SKILL.md), [`coding-standards`](skills/coding-standards/SKILL.md), [`me-on-review`](skills/me-on-review/SKILL.md), [`my-own`](skills/my-own/SKILL.md), [`doubt-driven-development`](skills/doubt-driven-development/SKILL.md) | 9 |
| **Debugging and incidents** | [`debug`](skills/debug/SKILL.md), [`debugging-and-error-recovery`](skills/debugging-and-error-recovery/SKILL.md), [`incident-response-and-postmortems`](skills/incident-response-and-postmortems/SKILL.md), [`error-handling-philosophy`](skills/error-handling-philosophy/SKILL.md) | 4 |
| **Delivery and release** | [`ci-cd-and-automation`](skills/ci-cd-and-automation/SKILL.md), [`git-workflow-and-versioning`](skills/git-workflow-and-versioning/SKILL.md), [`release-management-and-rollback-planning`](skills/release-management-and-rollback-planning/SKILL.md), [`shipping-and-launch`](skills/shipping-and-launch/SKILL.md), [`deprecation-and-migration`](skills/deprecation-and-migration/SKILL.md), [`backward-compatibility-and-deprecation-policy`](skills/backward-compatibility-and-deprecation-policy/SKILL.md) | 6 |
| **Planning, specs and documentation** | [`spec-driven-development`](skills/spec-driven-development/SKILL.md), [`planning`](skills/planning/SKILL.md), [`planning-and-task-breakdown`](skills/planning-and-task-breakdown/SKILL.md), [`lifecycle-planner`](skills/lifecycle-planner/SKILL.md), [`incremental-implementation`](skills/incremental-implementation/SKILL.md), [`idea-refine`](skills/idea-refine/SKILL.md), [`interview-me`](skills/interview-me/SKILL.md), [`defining-done-and-acceptance-criteria`](skills/defining-done-and-acceptance-criteria/SKILL.md), [`source-driven-development`](skills/source-driven-development/SKILL.md), [`search-first`](skills/search-first/SKILL.md), [`documentation-and-adrs`](skills/documentation-and-adrs/SKILL.md), [`technical-decision-records`](skills/technical-decision-records/SKILL.md), [`documentation-lookup`](skills/documentation-lookup/SKILL.md), [`project-guidelines-example`](skills/project-guidelines-example/SKILL.md) | 14 |
| **Research** | [`research`](skills/research/SKILL.md), [`deep-research`](skills/deep-research/SKILL.md), [`market-research`](skills/market-research/SKILL.md) | 3 |
| **Agentic engineering** | [`agentic-engineering`](skills/agentic-engineering/SKILL.md), [`ai-first-engineering`](skills/ai-first-engineering/SKILL.md), [`agent-harness-construction`](skills/agent-harness-construction/SKILL.md), [`agent-drift-and-re-anchoring`](skills/agent-drift-and-re-anchoring/SKILL.md), [`multi-agent-orchestration`](skills/multi-agent-orchestration/SKILL.md), [`orchestrator`](skills/orchestrator/SKILL.md), [`team-builder`](skills/team-builder/SKILL.md), [`autonomous-loops`](skills/autonomous-loops/SKILL.md), [`continuous-agent-loop`](skills/continuous-agent-loop/SKILL.md), [`continuous-learning`](skills/continuous-learning/SKILL.md), [`continuous-learning-v2`](skills/continuous-learning-v2/SKILL.md), [`context-engineering`](skills/context-engineering/SKILL.md), [`strategic-compact`](skills/strategic-compact/SKILL.md), [`prompt-optimizer`](skills/prompt-optimizer/SKILL.md), [`using-agent-skills`](skills/using-agent-skills/SKILL.md), [`skill-catalog-authoring`](skills/skill-catalog-authoring/SKILL.md), [`data-scraper-agent`](skills/data-scraper-agent/SKILL.md) | 17 |
| **Product, leadership and communication** | [`writing-product-requirements`](skills/writing-product-requirements/SKILL.md), [`feature-prioritization-frameworks`](skills/feature-prioritization-frameworks/SKILL.md), [`metrics-and-north-star-frameworks`](skills/metrics-and-north-star-frameworks/SKILL.md), [`experiment-design-and-ab-testing`](skills/experiment-design-and-ab-testing/SKILL.md), [`user-research-synthesis`](skills/user-research-synthesis/SKILL.md), [`roadmapping-and-tradeoff-communication`](skills/roadmapping-and-tradeoff-communication/SKILL.md), [`stakeholder-communication-and-status-updates`](skills/stakeholder-communication-and-status-updates/SKILL.md), [`one-on-ones-and-feedback-frameworks`](skills/one-on-ones-and-feedback-frameworks/SKILL.md), [`engineering-onboarding-design`](skills/engineering-onboarding-design/SKILL.md), [`engineering-values-and-tradeoff-principles`](skills/engineering-values-and-tradeoff-principles/SKILL.md), [`code-ownership-and-boundaries`](skills/code-ownership-and-boundaries/SKILL.md), [`content-engine`](skills/content-engine/SKILL.md) | 12 |
| **Industry domains** | [`hotel-rate-and-inventory-modeling`](skills/hotel-rate-and-inventory-modeling/SKILL.md), [`channel-manager-ota-integration`](skills/channel-manager-ota-integration/SKILL.md), [`customs-trade-compliance`](skills/customs-trade-compliance/SKILL.md), [`production-scheduling`](skills/production-scheduling/SKILL.md) | 4 |

Run `make list` (or `scripts/list-skills.sh`) for the live list with eval case counts. Seventeen skills
also have a human-facing "why and when" page under [`skill-docs/`](skill-docs/).

> `npx skills add shennawardana23/skillme --list` shows 148 entries: these 143 plus five slide-authoring
> skills (`create-slide`, `slide-authoring`, `apply-comments`, `create-theme`, `current-slide`) that belong
> to the [slides workspace](#slides-workspace) under `.claude/skills/` and `.agents/skills/`. Skip them
> unless you author decks with open-slide.

---

## How a skill is proven

```
skills/<name>/SKILL.md + evals/evals.json      <- prompt + assertions, written by hand
              |
              v
        smeval validate                        <- schema check, free, no model call
              |
              v
        smeval run                             <- fresh `claude` CLI process, isolated workspace,
              |                                   only this skill installed
              v
        grading.json + report.html + feedback.json
        quoted evidence per assertion + a human-review stub, never a bare pass or fail
```

- **Deterministic grading, no LLM judge.** Each assertion is a string, regex or file check, so a result is
  repeatable and free to grade ([ADR-0002](docs/adr/0002-own-deterministic-eval-runner.md)).
- **Isolation.** Each case runs in its own workspace where the skill under test is the only skill visible,
  so another skill cannot cover for a weak one ([ADR-0004](docs/adr/0004-isolate-each-eval-case.md)).
- **With and without.** `-benchmark` runs every case again without the skill and writes `benchmark.json`,
  so you can see whether the skill changes the answer at all.
- **A human still reads it.** Every run writes `feedback.json`; `smeval feedback-check` fails until someone
  has reviewed every case.

### Coverage, stated honestly

"Schema-valid" and "works against a real model" are different claims ([`TESTING.md`](TESTING.md)).

- **Every skill** passes `agentskills validate`, `smeval validate` and the static security scan on every
  push (CI).
- **Live runs.** The rolling sweep (`scripts/test-all-skills.sh`, `make test-all`) records results in
  `smeval-workspace/test-all-results/summary.txt` on the machine that ran it (gitignored). The last full
  record, written 17 August 2026, covers 125 skills and 362 cases: 304 cases passed fully, 77 skills passed
  every case and 48 passed some. Re-run it for a current number instead of trusting this one.
- **What the failures turned out to be.** When the original 125-skill sweep was diagnosed against the
  model's saved output, the partial failures were eval-authoring bugs (an assertion too narrow for a
  correct answer phrased differently, a markdown backtick breaking a literal match, a negative check
  banning a phrase used correctly), a few single-run variances, and one real bug in `smeval` itself: a case
  could write into this repo instead of its workspace ([`TESTING.md`](TESTING.md#a-fixed-workspace-escape-bug-the-isolated-case-wrote-to-the-real-repo)).
  A known, still-open confound is [account-level instruction bleed](TESTING.md#a-known-unfixed-confound-account-level-instruction-bleed).
- **Not yet in a full sweep:** the 18 skills added after it, including `accessibility`, `seo`,
  `web-quality-audit`, the `golang-*` additions, `skill-inspector` and the PR #5 skills. Each has its own
  eval suite. Do not read "143 skills" as "143 independently proven".

---

## For contributors

### Requirements

| Tool | Needed for |
| --- | --- |
| Go 1.25+ | Building `smeval` (`go.mod` pins `go 1.25`) |
| [`uv`](https://docs.astral.sh/uv/) | `uvx --from skills-ref agentskills validate`, the official spec validator, without a separate install |
| [`claude` CLI](https://github.com/anthropics/claude-code), logged in or with `ANTHROPIC_API_KEY` | Live evals only; validation and scans work without it |
| Node.js 18+ and npm | Only the [slides workspace](#slides-workspace) |

### Evaluating a skill locally

```bash
# Official Agent Skills spec compliance (frontmatter, naming): free
uvx --from skills-ref agentskills validate skills/go-service-idioms

# This catalog's own eval schema and cases
go build -o smeval ./cmd/smeval
./smeval validate skills/go-service-idioms              # free, schema only
./smeval run      skills/go-service-idioms              # runs cases against the local `claude` CLI
./smeval run      skills/go-service-idioms -benchmark   # also runs without the skill, for comparison
```

`smeval validate` only proves the JSON is well formed. `smeval run` proves the skill works: it runs each
case through the local `claude` CLI (headless, `--output-format json`) and grades the real output. Results
land in `smeval-workspace/runs/<name>/iteration-N/`:

- `report.html`: pass or fail per case, with the quoted assertion evidence
- `feedback.json`: where a human records what an assertion cannot catch ([guide](TESTING.md#reviewing-results-with-a-human))
- `benchmark.json`: with-skill versus without-skill; [`examples/benchmark.sample.json`](examples/benchmark.sample.json) shows real output

**Provider and model fallback.** `smeval run` takes `-primary-model` (default `sonnet`) and `-fallback-model`
(default `opus`). Fallback works in two layers: the `claude` CLI's own `--fallback-model` handles an
overloaded or unavailable model within one call, and `smeval`'s outer retry (`internal/engine`) catches
whole-process failures (non-zero exit, timeout, unparseable output) and retries once in a fresh process on
the fallback model. Both are covered by tests in `internal/engine/engine_test.go`. Other flags: `-timeout`,
`-include <case-substring>`, `-output-dir`.

### Makefile

`make help` lists every target. The ones you will use:

| Target | What it does |
| --- | --- |
| `make ci` | `build` + `vet` + `unit-test` + `fmt-check`, the same steps as CI's `test` job |
| `make validate SKILL=<name>` / `make validate-all` | Spec and schema validation, free |
| `make security-scan SKILL=<name>` / `security-scan-all` | Static risk scan, free |
| `make eval SKILL=<name>` | Live eval for one skill (`INCLUDE=<case-substring>` narrows to one case) |
| `make eval-benchmark SKILL=<name>` / `make benchmark-report SKILL=<name>` | Live eval with a without-skill baseline, and its comparison |
| `make test-all` / `test-all-fresh` / `test-all-only SKILL=<name>` | Live-test the whole catalog, resumable |
| `make list` | Every skill with its eval case count |
| `make report SKILL=<name>` / `make cat SKILL=<name> CASE=<id>` | Open the latest `report.html` / print one case's response and grading |
| `make feedback SKILL=<name>` / `make feedback-check SKILL=<name>` | Open the review stub / fail while entries are unreviewed |
| `make trigger-validate` / `make trigger-run` / `make catalog-report` | Catalog-wide trigger-accuracy cases |
| `make similarity-check THRESHOLD=0.NN` | Advisory overlap check across all skill descriptions |
| `make clean` | Remove the `smeval` binary and local eval workspaces (both gitignored) |

### Catalog-wide tooling

Some questions cannot be answered by testing one skill alone. All of these are free (no model calls)
except the first.

**Does the agent pick the right skill, and does a new one steal triggers from an old one?**
`harness.BuildCatalog` installs the whole catalog into one harness, and
`smeval trigger-run catalog-evals/trigger-accuracy.json` grades cases against it. This is where triggering
accuracy and cross-skill coexistence get checked, two of Anthropic's
[enterprise evaluation dimensions](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/enterprise).
Real model calls, real cost; see [`catalog-evals/README.md`](catalog-evals/README.md).

**Is anything in a skill unsafe** (a bundled script, a live network call, a hardcoded credential)?
`smeval security-scan skills/<name>` checks against Anthropic's enterprise risk-tier indicators. It is
deliberately narrow where static analysis cannot be trusted: skills that teach about prompt injection
(`security-review`, `skill-inspector`) use the same words as an attack. `internal/riskscan` documents which
indicators fail a build and which are advisory. Runs in CI on every push.

**Does a new skill idea duplicate an existing one?** `smeval similarity-check` scores every pair of skill
descriptions by TF-IDF cosine similarity. It is lexical, not semantic: it catches shared vocabulary, not
paraphrase. Always advisory; `internal/similarity` explains the trade-off.

**Did anyone read the model's response, or just trust the green check?** Every run writes a `feedback.json`
stub. An audit found every one empty, so `smeval feedback-check` now fails until each case is reviewed. It is
step 4 of the pre-PR checklist in `CONTRIBUTING.md`.

**Context-dependent cases.** A case can set `"fixture": "<name>"` to seed files from the skill's
`evals/fixtures/<name>/` into its workspace before the prompt runs. `debug`'s `traces-bug-across-files` case is
a worked example: a real two-file Go bug the model has to find. `smeval validate` catches a missing fixture
before a live run wastes a call on it.

### Add or change a skill

1. **Search first.** An overlapping skill is more likely than a real gap (`make list`, `make similarity-check`).
   The [authoring rule](.claude/rules/skill-authoring.md) and [`CONTRIBUTING.md`](CONTRIBUTING.md#before-proposing-a-new-skill)
   have the pre-flight checks.
2. **Scaffold** with `/new-skill <name>` in Claude Code, or copy an existing skill. The enforced spec is
   [`skill-catalog-authoring`](skills/skill-catalog-authoring/SKILL.md).
3. **Write `evals/evals.json`** with assertions that survive rephrasing (`contains_any`, not one exact
   sentence), and `not_contains` that bans the wrong thing, not a whole word.
4. **Check:** `make validate`, `make security-scan`, then `make eval` (or `/eval-skill <name>`) until it
   passes, then review `feedback.json`.
5. **Open a PR.** CI validates and scans every skill; live evals are run by people, because they cost money.

### Anatomy of a skill

```
skills/go-service-idioms/
├── SKILL.md              ┌─ frontmatter ─────────────────────────┐
│                         │ name: go-service-idioms (= folder)    │
│                         │ description: what it does + when, with│
│                         │   the trigger phrases                 │
│                         │ license: Apache-2.0                   │
│                         │ metadata: { version, category }       │
│                         └───────────────────────────────────────┘
│                         Body: instructions for the agent, imperative,
│                         grounded in a real source, ending in Gotchas
│                         and Real-world grounding.
├── references/           Optional detail, read only when the body points
│                         to it (progressive disclosure).
└── evals/
    ├── evals.json        This skill's test cases.
    └── fixtures/<name>/  Optional files seeded into a case's workspace.
```

Every "do X because Y" in a body should trace to something real: an official doc, a documented incident or a
verified API. [`skill-catalog-authoring`](skills/skill-catalog-authoring/SKILL.md) is the source of truth;
this section is the map.

### Repository layout

```text
skillme/
├── skills/<name>/                     the catalog (143 skills)
├── skill-docs/<area>/<name>.md        optional human-facing "why and when" pages
├── catalog-evals/                     catalog-wide trigger-accuracy and coexistence cases
├── cmd/smeval/                        the eval runner CLI
├── internal/{engine,evalspec,grading,harness,report,riskscan,similarity}/   smeval's implementation
├── scripts/                           list-skills.sh, test-all-skills.sh (resumable full sweep)
├── examples/benchmark.sample.json     real -benchmark output, annotated in TESTING.md
├── docs/                              project overview and architecture decisions (docs/adr/)
├── Makefile                           make help: wraps smeval and the scripts
├── .claude-plugin/                    plugin.json + marketplace.json (Claude Code distribution)
├── CLAUDE.md, .claude/{rules,commands}/   authoring rule loaded by Claude Code; /eval-skill, /new-skill
├── .github/workflows/                 skill-eval.yml, test-plugin-install.yml
├── slides/, themes/, assets/, package.json, open-slide.config.ts   the slides workspace (below)
├── AGENTS.md                          agent rules for the slides workspace
└── .claude/skills/, .agents/skills/   open-slide's deck-authoring skills, managed by @open-slide/core
```

### What CI checks

| Workflow | Checks |
| --- | --- |
| `skill-eval.yml` | `go build`, `go vet`, `go test`, gofmt; for every skill: `agentskills validate`, `smeval validate`, `smeval security-scan` |
| `test-plugin-install.yml` | `claude plugin validate .`, then a real `marketplace add` and `plugin install` |

CI never calls a model. Live evals are a deliberate, manual step.

---

## Slides workspace

This repo also hosts an [open-slide](https://www.npmjs.com/package/@open-slide/core) workspace for decks
about the project, such as `slides/skillme-overview/`. The Go code and the catalog do not depend on it.

```bash
npm install
npm run dev        # dev server with hot reload
npm run build      # static bundle in dist/ (netlify.toml and vercel.json rewrite every path to index.html)
npm run preview
```

Each deck is `slides/<id>/index.tsx` on a fixed 1920 x 1080 canvas. Deck-authoring rules are in
[`AGENTS.md`](AGENTS.md) and the skills under `.claude/skills/`; update those skills with
`npm run sync:skills`, not by hand.

---

## Architecture decisions

The reasons behind the big choices are recorded as ADRs in [`docs/adr/`](docs/adr/README.md):

| ADR | Decision |
| --- | --- |
| [0001](docs/adr/0001-agent-skills-format.md) | Use the open Agent Skills format as the only content format |
| [0002](docs/adr/0002-own-deterministic-eval-runner.md) | Build our own eval runner with deterministic grading, no LLM judge |
| [0003](docs/adr/0003-distribution-plugin-and-skills-cli.md) | Distribute as a Claude Code plugin and through the `npx skills` CLI |
| [0004](docs/adr/0004-isolate-each-eval-case.md) | Run every eval case in an isolated workspace with one skill |
| [0005](docs/adr/0005-host-slides-workspace-in-repo.md) | Host the open-slide deck workspace in this repo |

---

## Documentation map

| Document | What it covers |
| --- | --- |
| [`TESTING.md`](TESTING.md) | The eval method: schema versus live, reading a failure's evidence, `-benchmark`, human review, real bugs the loop found in itself |
| [`CONTRIBUTING.md`](CONTRIBUTING.md) | Authoring conventions, the eval mistakes this catalog has made, the pre-PR checklist, what CI does not check |
| [`skill-catalog-authoring`](skills/skill-catalog-authoring/SKILL.md) | The enforced spec for skill layout and frontmatter |
| `skill-docs/<area>/<name>.md` | Optional per-skill human pages; template in `skills/skill-catalog-authoring/references/skill-docs-template.md` |
| [`catalog-evals/README.md`](catalog-evals/README.md) | Catalog-wide trigger-accuracy and coexistence testing |
| [`docs/skillme-project-overview.md`](docs/skillme-project-overview.md) | Leadership summary and comparison with similar catalogs |
| [`docs/adr/`](docs/adr/README.md) | Architecture decision records |
| [`CLAUDE.md`](CLAUDE.md), [`.claude/rules/skill-authoring.md`](.claude/rules/skill-authoring.md) | Rules Claude Code loads in this repo |
| [`.claude/commands/eval-skill.md`](.claude/commands/eval-skill.md), [`new-skill.md`](.claude/commands/new-skill.md) | `/eval-skill <name>` and `/new-skill <name>` from inside Claude Code |

---

## Troubleshooting

| Problem | Fix |
| --- | --- |
| Every skill appears twice | You installed the plugin and an `npx skills` copy into the same agent. Remove one. |
| `marketplace add` fails with an SSH error | `git config --global url."https://github.com/".insteadOf "git@github.com:"`, then retry. |
| A new or updated skill does not show up | Restart the agent session. Plugin users run `claude plugin update skillme` first. |
| A skill does not trigger | Call it by name (see [Use a skill](#use-a-skill)); check its folder name equals its frontmatter `name`. |
| Symlinks fail (Windows, some sandboxes) | `npx skills@latest add ... --copy`, or copy the folders instead of linking. |
| A skill you made locally is missing from `npx skills` | `npx skills` reads GitHub, not your disk. Push it, and check it is not gitignored. |
| `git add .` fails on `smeval-workspace/` | It holds nested eval repos and is gitignored on purpose. Stage paths explicitly. |

---

## Security

Skills carry instructions and sometimes scripts that an agent will act on, so treat a third-party skill like
a dependency. Before installing one from elsewhere, run [`skill-inspector`](skills/skill-inspector/SKILL.md)
or `smeval security-scan <dir>`. Report a problem in this catalog through a GitHub issue; for anything
sensitive, contact the maintainer directly instead of filing it publicly.

## Contributing

New skills, fixes and `smeval` improvements are welcome. Start with [`CONTRIBUTING.md`](CONTRIBUTING.md):
it has the layout and frontmatter rules, the eval-writing mistakes this catalog has made, the pre-PR
checklist, and what CI does and does not check.

## License

[Apache-2.0](LICENSE). Each skill declares `license: Apache-2.0`, except four that omit the field
(`adk-go-agent-builder`, `go-service-idioms`, `postgres-hotel-partitioning`, `skill-catalog-authoring`)
and fall under the repository license.
