# Writing a skill's docs page

Every skill's `SKILL.md` is written for the agent executing it — terse,
imperative, spec-constrained. A `skill-docs/<category>/<skill-name>.md` page
is written for a different reader: a human deciding *whether and when* to
reach for this skill at all. It is not a copy of `SKILL.md` and not a second
copy of `references/` — it answers "why does this exist and when do I want
it," not "how do I execute it."

## Required section order

```markdown
## What it does

## When to reach for it

## Prerequisites          <!-- omit entirely if the skill has no real setup/state requirement -->

## <1-3 free-form sections in the skill's own vocabulary>

## Common questions

## It's working if

## Where it fits
```

- **What it does** — one or two sentences: the job, plus the one
  non-obvious constraint that makes this skill worth having rather than
  something an agent would just figure out unprompted (usually the same
  fact that anchors the skill's `## Gotchas` section).
- **When to reach for it** — the trigger condition, and the boundary
  against the nearest sibling skill if one exists (e.g. `debug` vs
  `debugging-and-error-recovery`, `postgres-patterns` vs
  `postgres-hotel-partitioning`). Name the sibling and the actual
  dividing line, not just "see also."
- **Prerequisites** — only include this section if the skill genuinely
  needs something set up first (a database group, a config file, a
  specific framework version). Most skills omit it.
- **Free-form middle** — 1-3 sections using the terminology this specific
  skill actually uses. Do not force every skill into the same subsection
  names; a database skill's middle looks nothing like a philosophy skill's.
- **Common questions** — real, documented points of confusion about this
  exact technology or practice, each traceable to the same class of primary
  source the skill's own `## Real-world grounding` section cites (an
  official doc's explicit warning, a well-known version-gated behavior, a
  commonly-reported mistake). **Never invent a question to fill space.** A
  skill grounded in concrete, version-specific, or historically-documented
  gotchas (a database engine, a framework) will have several real
  questions. A skill that's mostly a judgment framework (a process or
  philosophy skill) may have one or two, or none — in that case, shorten or
  omit the section rather than padding it with generic questions ("what if
  I don't know where to start?") that aren't actually sourced from anything.
- **It's working if** — a short bullet list a reader can check without
  opening `SKILL.md` — observable outcomes, not restated instructions.
- **Where it fits** — this skill's role relative to the rest of the
  catalog: is it a standalone reference, a step in a chain with named
  neighbors, a one-time setup skill, or an ongoing-maintenance skill? Name
  the actual neighboring skills by their real slug.

## Rules

- No H1 — the page's title is its filename/slug, not a heading inside it.
- No install/setup commands for the plugin itself (that's the README's job).
- No attribution to an individual as the "author" of the skill or the
  content it's grounded in — skillme's skills are catalog entries, not
  personal essays; write findings as the catalog's own stated fact
  ("CI4's docs state X"), not "we found that X."
- Link to a sibling skill by name (`postgres-hotel-partitioning`), not by a
  relative file path — the catalog's directory layout can change; a skill
  name is the stable reference.
- Category matches the skill's `metadata.category` frontmatter value (see
  the main `SKILL.md` conventions section) and determines the page's path:
  `skill-docs/<category>/<skill-name>.md`.

## The most important rule here

A "Common questions" section is evidence-sized, not padded to look
complete. A skill grounded in concrete, version-specific, or
historically-documented gotchas will have several real questions; a skill
that's mostly a judgment framework may have one, or none. Shorten or omit
the section rather than inventing a generic question to fill space.
