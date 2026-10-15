# 0005. Host the open-slide deck workspace in this repo

## Status
Accepted (recorded retroactively)

## Date
2026-10-15 (PR #5, commit `a9e657b`)

## Context
The project needs decks to present itself (for example `slides/skillme-overview/`). open-slide builds
decks as React components and ships its own agent rules and authoring skills. Putting the workspace
next to the catalog keeps the decks close to the content they describe.

## Decision
The open-slide workspace lives in this repo: `package.json`, `open-slide.config.ts`, `slides/`,
`themes/`, `assets/`, `AGENTS.md`, its authoring skills in `.claude/skills/` and `.agents/skills/`, and
static hosting configs (`netlify.toml`, `vercel.json`).

## Alternatives considered
- **A separate repo for decks**: keeps the catalog repo single-purpose and its agent rules clean, but
  decks drift from the content and need their own setup. Not chosen at the time.
- **Decks as static HTML** (for example with the `frontend-slides` skill): no Node toolchain, but no
  live editing workflow. Not chosen.

## Consequences
- A Node toolchain sits beside the Go one; it is optional for catalog work.
- PR #5 replaced the catalog `README.md` with open-slide's template README; it was restored and
  rewritten on 2026-10-15.
- The root `AGENTS.md` is open-slide's deck-authoring guide. Codex, Cursor, GitHub Copilot and other
  agents read `AGENTS.md`, so in this repo they are told they are "authoring slides", even when the task
  is catalog work. This needs a follow-up: either scope `AGENTS.md` to both purposes or move the deck
  workspace out.
- `npx skills add shennawardana23/skillme` lists the five open-slide skills beside the 143 catalog
  skills (148 in total).
