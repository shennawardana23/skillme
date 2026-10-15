# 0004. Run every eval case in an isolated workspace with one skill

## Status
Accepted (recorded retroactively)

## Date
2026-08-15 (`internal/harness` in the initial release; subprocess isolated from the tester's global
`CLAUDE.md` in `1083fde`; workspace-escape fix on 2026-08-17 in `e1028b9`)

## Context
A skill's eval should measure that skill. If the model under test can also see other skills, the
tester's personal instructions, or this repo's files, a pass might come from somewhere else, and a case
could change the repo it is testing.

## Decision
`smeval run` installs only the skill under test into a throwaway plugin directory (never its own
`evals/`), gives each case and each configuration (with and without the skill) its own workspace, and
starts the `claude` CLI so it does not load the tester's global settings. Catalog-wide questions
(does the right skill get picked) are tested separately by `smeval trigger-run`, which installs the
whole catalog on purpose.

## Alternatives considered
- **Run cases inside the repo checkout**: simple, but a case can read or write real files and see every
  skill. Rejected; a real case did write into the repo before the fix.
- **Run with the whole catalog installed**: closer to real use, but a strong neighbour can hide a weak
  skill. Kept only for the trigger-accuracy suite.

## Consequences
- A pass means this skill, in this prompt, produced this evidence.
- Each run creates many workspaces, including nested git repos, under `smeval-workspace/`; it is
  gitignored and `git add .` must not include it.
- One confound is still open: account-level instruction bleed from the logged-in `claude` account
  (`TESTING.md`). Isolation covers the filesystem and settings, not the account.
