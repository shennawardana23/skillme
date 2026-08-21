---
description: Scaffold a new skill directory following this repo's spec
---

Scaffold a new skill named `$ARGUMENTS` under `skills/`:

1. Read `skills/skill-catalog-authoring/SKILL.md` first — it's the
   enforced spec for directory layout, required frontmatter, and naming.
2. Create `skills/$ARGUMENTS/SKILL.md` with frontmatter (`name`,
   a `description` written as trigger phrasing, `license: Apache-2.0`,
   `metadata: { version: "0.1.0", category }`) and a body grounded in a
   real doc, codebase, or incident — no invented claims.
3. Create `skills/$ARGUMENTS/evals/evals.json` with real cases, following
   `CONTRIBUTING.md`'s eval-authoring guidance (grade the response, not
   the log; assertions that survive rephrasing).
4. Run `/eval-skill $ARGUMENTS` and get it fully passing before opening a PR.

If `$ARGUMENTS` is empty, ask for the skill name first.
