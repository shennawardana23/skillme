---
description: Anti-duplication guardrail plus eval-authoring rules for adding or changing a skill
paths:
  - "skills/**"
---

# Adding or changing a skill

This catalog already has 132 skills, so a new idea more often overlaps
an existing skill than fills a real gap. Before creating a new
`skills/<name>/` directory or significantly reworking one:

- Run the pre-flight checks in
  [CONTRIBUTING.md](../../CONTRIBUTING.md#before-proposing-a-new-skill):
  search the catalog, check open PRs (`gh pr list --state open`), and
  justify the gap.
- Prefer extending an existing skill over adding a near-duplicate.
  `security-review` / `security-and-hardening` / `security-scan` already
  read as near-duplicates by name alone — each one's own `description`
  states the exact dividing line against its siblings, which is the
  pattern to follow, not more overlapping skills.
- Keep `SKILL.md` within the shape in the README's
  [Anatomy of a skill](../../README.md#anatomy-of-a-skill) section and
  the full spec in `skills/skill-catalog-authoring/SKILL.md`. Never
  duplicate content between skills — reference the other skill instead.

## Eval-authoring rules

- Grade against `outputs/response.md`, never a `.log` file — a log can
  contain tool-call text that coincidentally matches a substring
  assertion even when the actual response doesn't.
- Assertions must survive rephrasing: prefer `contains_any` over one
  literal string, and never assert on markdown formatting (backticks,
  hyphens, commas) a correct answer could render differently.
- `not_contains` has to ban the *wrong* thing specifically — a blunt ban
  on a whole word can fail a correct answer that uses it validly.
- `smeval run` must pass before opening a PR — `smeval validate` alone
  only proves the JSON parses, not that the skill works.

Full checklist and the real mistakes this catalog has made:
`CONTRIBUTING.md`. This rule points to it rather than restating it.
