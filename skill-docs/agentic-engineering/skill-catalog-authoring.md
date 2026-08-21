## What it does

Governs how a skill enters or changes in this catalog: the required
directory layout, `SKILL.md` frontmatter, and the `evals.json` schema
`smeval` (this repo's own eval runner) reads. The non-obvious constraint:
a schema-valid skill and a working skill are different claims — schema
validation is free and catches malformed files, but only a live run
against a real model (`smeval run`) proves the guidance actually works.

## When to reach for it

Before adding a new `skills/<name>/` directory, before reworking an
existing skill's frontmatter or eval cases, or any time you're not sure
whether a change needs a schema check, a live run, or both. If your idea
overlaps an existing skill by name or scope, extend that skill instead of
adding a near-duplicate — `security-review` / `security-and-hardening` /
`security-scan` already read as near-duplicates by name alone, and each
one's own `description` states the exact dividing line against its
siblings; that's the pattern to follow, not more overlapping skills.

## The two-layer test loop

1. **Schema validation** — `uvx --from skills-ref agentskills validate
   skills/<name>` (official spec) and `./smeval validate skills/<name>`
   (this catalog's own `evals.json` schema). Free, no model call, catches
   malformed frontmatter and JSON — nothing about whether the content
   actually works.
2. **Live eval** — `./smeval run skills/<name>` (add `-benchmark` to also
   run every case with the skill absent, for a side-by-side comparison).
   Calls the real `claude` CLI and grades the actual response. This is
   the one that catches real bugs — see [`TESTING.md`](../../TESTING.md)
   for the actual mistakes this catalog has made and fixed this way.

`/eval-skill <name>` (`.claude/commands/eval-skill.md`) runs both steps
and summarizes the result without leaving the chat.

## Common questions

- **Schema validation passed — why did the live run still fail?** They
  check different things. A grading assertion can be checking the chat
  response when the model correctly wrote its real answer to a file, or
  be phrased too narrowly to match a correct-but-differently-worded
  answer — both are eval-authoring bugs, not proof the skill's content is
  wrong. Read the failure's `grading.json` evidence before concluding
  either way.
- **What's `feedback.json` for, if assertions already pass?** Assertion
  grading only checks what someone thought to write an assertion for.
  `feedback.json` is where a human records what a checklist can't —
  whether a headline is actually punchy, whether an answer is technically
  correct but misses the point. See "Reviewing results with a human" in
  [`TESTING.md`](../../TESTING.md).
- **`-benchmark`'s delta is near zero — is the skill broken?** Not
  necessarily broken, but not proven valuable either: either the model
  already knew this without help, or the eval cases aren't actually
  probing what the skill teaches. Worth investigating before shipping.

## It's working if

- `smeval validate` and the official `agentskills validate` both pass
- `smeval run` shows every case passing, with evidence in `report.html`
  that actually matches what the case's assertions claim to check
- `feedback.json` has been opened and filled in, not left as an
  unreviewed stub
- The skill's `name` matches its directory exactly, and its `description`
  states its trigger phrasing and — if a sibling skill's scope is close —
  the dividing line against it

## Where it fits

The skill most other contributions touch — `CONTRIBUTING.md` and this
catalog's `README.md` ("Anatomy of a skill," "How a skill gets proven")
both point back to it as the source of truth rather than restating its
rules. `TESTING.md` is the deep-dive on the live-eval half of this
skill's loop; this page is the "why and when," not a duplicate of either.
