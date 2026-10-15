# 0003. Distribute as a Claude Code plugin and through the `npx skills` CLI

## Status
Accepted (recorded retroactively)

## Date
2026-08-15 (plugin manifests in the initial release, commit `79b9b5b`; install CI added 2026-08-21, `7877a19`)

## Context
Users want either the whole catalog managed for them or a few skills they can edit. Claude Code has a
plugin marketplace mechanism; other agents read skill folders from their own paths. Evidence:
`.claude-plugin/plugin.json` (`"skills": "./skills"`), `.claude-plugin/marketplace.json`,
`.github/workflows/test-plugin-install.yml`, and the install section of the README.

## Decision
We ship two paths and support both:
1. A Claude Code plugin with this repo as its own marketplace
   (`claude plugin marketplace add shennawardana23/skillme`, `claude plugin install skillme@skillme`),
   tested by a real install in CI.
2. The open `npx skills` CLI (skills.sh), which reads `skills/<name>/SKILL.md` from GitHub and writes
   chosen skills into any of about 70 agents' folders, with no registration on our side.

## Alternatives considered
- **Per-tool plugin packages** (Codex, Cursor, Copilot marketplaces): native updates in each tool, but
  each needs its own manifest and testing. Not done yet; a future ADR can add them.
- **An npm package**: familiar to JavaScript users, but adds packaging and release work for plain
  Markdown. Rejected.
- **Clone-only**: simplest, but every user maintains copies by hand. Kept as the manual fallback, not
  the main path.

## Consequences
- Claude Code users get one managed bundle with updates; everyone else gets editable files.
- Installing both paths into one agent loads every skill twice; the README warns about this.
- `npx skills` serves whatever is on GitHub `main`, so local or gitignored skills never appear there.
- `npx skills` lists every `SKILL.md` it finds in the repo, so the slides workspace's five authoring
  skills show up beside the catalog (see ADR-0005).
- Cursor can import a repo only through a `.cursor-plugin/marketplace.json`, and Codex and Copilot have
  their own plugin marketplaces; we provide none of those manifests, so those users rely on `npx skills`
  or manual copies.
