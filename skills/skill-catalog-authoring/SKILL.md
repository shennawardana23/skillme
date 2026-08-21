---
name: skill-catalog-authoring
description: This skill should be used when the user asks to "add a new skill to skillme", "add a skill to this catalog", "write an eval for this skill", "validate this skill's evals", "run smeval on this skill", or otherwise extends or maintains a skill inside this plugin's skills/ directory. Provides the required directory layout, SKILL.md conventions, and evals.json schema for this repository's own eval runner.
metadata:
  version: "0.1.0"
---

# Skill Catalog Authoring

This plugin (`skillme`) is a catalog of Claude Code skills, each shipped
with an eval suite graded by `smeval` — this repository's own eval runner
(`cmd/smeval`, plain Go, no third-party dependency) — so a skill's quality
is measurable and re-checked in CI on every change, rather than trusted on
first impression. `smeval` follows Anthropic's documented Agent Skills
evaluation methodology (https://agentskills.io/skill-creation/evaluating-skills):
a prompt and assertions per case, graded with concrete evidence. Follow this
process whenever adding a skill to `skills/` or modifying an existing one's
`evals/`.

## Directory layout for a new skill

```
skills/<skill-name>/
├── SKILL.md              # required
├── references/           # optional — detailed docs loaded only when needed
├── examples/              # optional — working code examples
└── evals/
    └── evals.json         # required — smeval's case file for this skill
```

Use kebab-case for `<skill-name>`. `evals.json`'s top-level `skill_name`
must match the directory name exactly — `smeval validate` rejects a
mismatch.

## SKILL.md conventions

Frontmatter follows the official spec at https://agentskills.io/specification
exactly — the fields below are the **only** ones a validator will accept;
anything else (including a bare top-level `version:`) is rejected:

| Field | Required | Constraint |
|---|---|---|
| `name` | Yes | 1–64 chars, lowercase unicode alphanumeric + hyphens, no leading/trailing/double hyphen, **must equal the directory name exactly** |
| `description` | Yes | 1–1024 chars. Third person ("This skill should be used when the user asks to..."), packed with concrete quoted trigger phrases, not a vague topic label |
| `license` | No | Short license name or reference to a bundled `LICENSE` |
| `compatibility` | No | 1–500 chars. Only include if there's a real environment requirement (e.g. `"Requires Go 1.26+ and google.golang.org/adk/v2"`) |
| `metadata` | No | String→string map for anything else — **this is where `version` goes** (`metadata: { version: "0.1.0" }`), never as a bare top-level field |
| `allowed-tools` | No | Space-separated pre-approved tools (experimental, client-dependent) |

- **Body** in imperative/infinitive form ("Wrap errors with %w", not "You
  should wrap errors"). Keep it under 500 lines / ~5,000 tokens per the spec
  — that's a ceiling to move detail out of into `references/`, not a floor
  to pad toward. Most skills in this catalog run far leaner (600–850 words).
- **Ground every factual/API claim against a primary source** before writing
  it — read the actual upstream repository, run `go doc`, or fetch the
  official docs rather than recalling API shapes from training data. A
  skill that teaches a wrong function signature is worse than no skill.
- **Gotchas over generic advice.** Per the official best-practices guide
  (https://agentskills.io/skill-creation/best-practices), the highest-value
  content is a `## Gotchas` section of concrete, environment-specific facts
  the agent would get wrong unprompted — not generic advice ("handle errors
  appropriately") it already knows.

## Writing evals/evals.json

```json
{
  "skill_name": "<skill-name>",
  "evals": [
    {
      "id": "case-one",
      "prompt": "The exact task given to the agent under test.",
      "expected_output": "Human-readable description of what success looks like.",
      "assertions": [
        {
          "text": "Human-readable statement of what this checks",
          "check": { "contains_all": ["some required substring"] }
        },
        {
          "text": "Does not do the anti-pattern this skill exists to prevent",
          "check": { "not_contains": ["an anti-pattern that must not appear"] }
        }
      ]
    }
  ]
}
```

Every assertion pairs a human-readable `text` with a deterministic `check`
— `contains_all` / `contains_any` / `not_contains` / `matches_any` (Go
regexp) / `not_matches`, plus `files_exist` / `file_contains` for cases
that ask the agent to actually write files to its workspace rather than
answer inline. Prefer these deterministic checks over adding an LLM-judge
grading path for anything mechanically checkable — see
`references/smeval-schema.md` for the full check reference and for how
`smeval` isolates each case's workspace and skill installation.

Write at minimum 2–3 cases per skill: the one obviously-correct-usage case,
and at least one case that specifically tries to catch the skill's most
likely failure mode (the bug the skill exists to prevent, or a plausible
wrong-but-plausible-sounding answer).

## Validating and running

```bash
# Official spec compliance (frontmatter, naming) — free, no model calls
uvx --from skills-ref agentskills validate skills/<skill-name>

# This catalog's own eval schema and cases — free, no model calls
go build -o smeval ./cmd/smeval
smeval validate skills/<skill-name>
smeval run      skills/<skill-name>
```

Run both validators on every new or edited skill before committing —
`agentskills validate` catches frontmatter/naming spec violations (it will
reject a bare `version:` field, a `name` that doesn't match the directory,
a `description` over 1024 chars); `smeval validate` catches this catalog's
own eval schema errors (bad JSON, a `skill_name` mismatch, an assertion
with no check conditions).

## CI wiring

`.github/workflows/skill-eval.yml` builds `smeval` from source, runs
`go build`/`go vet`/`go test` first, then discovers every
`skills/*/evals/evals.json` dynamically and runs one job per skill. Adding
a skill following the layout above requires no workflow changes — it is
picked up automatically the next time CI runs.

## Writing a skill's docs page

Alongside `SKILL.md` (for the agent), a skill can have a
`skill-docs/<category>/<skill-name>.md` page (for a human deciding whether
to reach for it at all) — a different document, not a copy. See
`references/skill-docs-template.md` for the required section order, the
evidence-sizing rule for its "Common questions" section, and the
`metadata.category` field that determines its path. Not every skill needs
one yet; write it when a skill is genuinely reachable by more than one
plausible sibling and a human needs the boundary spelled out, or when asked
to add docs coverage for a batch of skills.

## Additional resources

For the full `evals.json` check reference, the fallback-mechanism design,
and how workspace/harness isolation works, consult
`references/smeval-schema.md`. For the docs-page template, consult
`references/skill-docs-template.md`.
