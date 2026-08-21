# Contributing to skillme

New skills, fixes to existing ones, and improvements to `smeval` itself
are all welcome. This doc is the map — it points at the specific
reference for each piece rather than repeating them, since the actual
conventions live in the skill and the eval methodology, not in prose
here that would drift out of sync with them.

## Before you start

- Go 1.25+ (to build `smeval`)
- [`uv`](https://docs.astral.sh/uv/) (`uvx`) — runs the official
  `agentskills` spec validator with no separate install
- The [`claude` CLI](https://github.com/anthropics/claude-code), logged
  in or with `ANTHROPIC_API_KEY` set — only needed to actually *run* a
  live eval; schema validation works without it

```bash
git clone https://github.com/shennawardana23/skillme
cd skillme
go build -o smeval ./cmd/smeval
go build ./... && go vet ./... && gofmt -l .   # confirm a clean baseline
```

## Before proposing a new skill

This catalog is already 132 skills, so a new idea more often overlaps an
existing one than fills a real gap. Before opening a PR for a new
`skills/<name>/` directory:

1. **Search the catalog.** Skim [What's covered](README.md#whats-covered)
   and `skills/` itself for something that already covers your idea,
   whole or in part.
2. **Check open PRs.** `gh pr list --state open` — don't add to a cluster
   of near-duplicate proposals on the same topic.
3. **Read the anatomy.** Confirm your idea fits the shape in the README's
   [Anatomy of a skill](README.md#anatomy-of-a-skill) section and the
   full spec in `skills/skill-catalog-authoring/SKILL.md` — an
   actionable, verifiable workflow, not general advice.
4. **State the gap in your PR description.** Say explicitly why this
   isn't covered by an existing skill or open PR. If it overlaps one,
   propose extending that skill instead of adding a new directory —
   `security-review` / `security-and-hardening` / `security-scan` already
   read as near-duplicates by name alone; each one's own `description`
   states the exact dividing line against its siblings, which is the
   pattern to follow instead of adding more overlap.

## Adding or editing a skill

Start with the `skill-catalog-authoring` skill
(`skills/skill-catalog-authoring/SKILL.md`) for the required directory
layout, frontmatter fields (`name`, `description`, `license`,
`metadata.version`, optional `metadata.category`), and the `evals.json`
schema. It also documents the optional `skill-docs/<category>/<name>.md`
human-facing page convention
(`skills/skill-catalog-authoring/references/skill-docs-template.md`) —
add one when the skill benefits from a "why would I reach for this"
explanation beyond what the agent-facing `SKILL.md` itself needs to say.

A few conventions this catalog holds to more strictly than the spec
requires:

- **Ground it in something real.** A skill's non-obvious claims should
  trace to an official doc, a real codebase pattern, or a documented
  incident — not general impression. If you can't point at a source, say
  so rather than asserting it as fact.
- **Generalize, don't attribute.** Skills read like something a senior
  engineer would hand a new hire, not a personal essay — no "we found
  that," no crediting an individual or an unrelated project by name for
  a pattern.
- **`name` must match the directory exactly** (lowercase,
  hyphen-separated, no underscores) — the spec validator and `smeval`
  both enforce this, and it's the single most common first-time
  validation failure.
- **Don't duplicate content between skills.** If two skills need the same
  explanation, reference the other skill by name instead of restating
  it — a duplicated paragraph drifts out of sync the first time only one
  copy gets fixed.
- **Reference material belongs in `references/`, not inline in `SKILL.md`**
  — keeps the always-in-context body short (progressive disclosure); see
  [Anatomy of a skill](README.md#anatomy-of-a-skill).

## Writing evals that actually test something

Read [`TESTING.md`](TESTING.md) before writing your first `evals.json` —
it documents the two-layer test model (schema vs. live), and, more
importantly, the specific eval-authoring mistakes this catalog has
actually made and fixed, so you don't repeat them:

- Grading the chat response when the model correctly wrote its real
  answer to a file (fix: point the assertion at the file with
  `files_exist`/`file_contains`, and tell the prompt exactly what to
  name it).
- A markdown backtick, a hyphen, or a comma breaking what should be a
  loose substring match (fix: `matches_any` with a tolerant regex, or a
  looser substring).
- Mixing `contains_any` and `matches_any` inside one `Check` object
  when you meant "either phrasing is fine" — within a single `Check`,
  every check type that's set must independently pass (AND, not OR).
  If two check types are meant to catch the *same* fact under different
  phrasings, they belong in one `matches_any` list of alternatives, not
  split across two check types.
- A blunt `not_contains` banning a phrase the model uses correctly while
  explaining what it did *not* do (fix: ban the specific accusatory
  phrasing that would appear if the model actually got it wrong, not the
  bare topic word).

## Before opening a PR

1. `go build ./... && go vet ./... && gofmt -l .` — for any change under
   `internal/` or `cmd/`
2. `uvx --from skills-ref agentskills validate skills/<name>` and
   `./smeval validate skills/<name>` for any skill you added or changed
3. `./smeval run skills/<name>` — **do this yourself; it is not a CI
   gate.** A passing schema is not the same as a working skill, and CI
   deliberately does not run live evals (see below) — actually running
   the eval against a real model is the step that catches real bugs, and
   it's on you to do it before opening the PR. If it fails, read
   `TESTING.md`'s "Diagnosing a failure" section before assuming the
   skill's content is wrong — check whether it's actually one of the
   eval-authoring patterns above.

## What CI does and doesn't check

`.github/workflows/skill-eval.yml` discovers any skill directory under
`skills/` that contains an `evals/evals.json` automatically — no
workflow edits needed to pick up a new skill. It runs schema/spec
validation only (`agentskills validate` + `smeval validate`), for free,
on every push. It deliberately does **not** run live evals: re-running
every skill's live model call on every push would mean real API cost and
significant runtime at this catalog's size, for even a one-line change
to a single skill. That's why step 3 above is on you, not the CI gate.

## Reporting a bug in `smeval` itself, not a skill

If you hit something that looks like the harness's fault rather than a
skill's content — a case's output landing somewhere other than its
isolated workspace, a check type behaving unexpectedly, a fallback that
didn't engage when it should have — read `internal/engine/engine.go` and
`TESTING.md`'s incident writeups first; two real bugs of exactly this
shape (an environment-isolation leak, a workspace-escape bug) have
already been found and fixed this way, both by treating a suspicious
result as worth investigating rather than accepting it.
