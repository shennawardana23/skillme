# 0001. Use the open Agent Skills format as the only content format

## Status
Accepted (recorded retroactively)

## Date
2026-08-15 (initial public release, commit `79b9b5b`)

## Context
The catalog's content has to reach engineers in several AI coding tools, not one. Each tool had its own
way to load instructions (rules files, custom instructions, prompt snippets), and copying content into
each format would let the copies drift apart. The open [Agent Skills](https://agentskills.io/specification)
format defines one folder layout (`SKILL.md` with `name` and `description` frontmatter, optional
`references/`, `scripts/`, `assets/`) and loads it by progressive disclosure, so many skills can be
installed for little context cost.

## Decision
Every unit of content is an Agent Skills folder, `skills/<name>/SKILL.md`, validated against the official
spec (`agentskills validate`) in CI. No other content format is maintained.

## Alternatives considered
- **Tool-specific rules files** (`CLAUDE.md` fragments, `.cursor/rules`, Copilot instructions): native to
  each tool, but one copy per tool and no shared validation. Rejected for drift.
- **A custom format plus converters**: full control, but every consumer would need our tooling. Rejected.
- **One large guide document**: easy to read, but always fully in context and not selectable per task.
  Rejected for context cost.

## Consequences
- The same folders work in every client that reads the format; by October 2026 that includes Claude Code,
  Codex, Cursor, GitHub Copilot, Gemini CLI, OpenCode and more than 30 others (agentskills.io client list).
- Tools without skills support can still use a skill as plain Markdown instructions, without automatic
  selection.
- The spec's rules bind us: folder name must equal `name`, `description` is at most 1024 characters and
  must be valid YAML (an unquoted `: ` inside it breaks strict parsers), and the description is the only
  thing an agent sees before it decides to load a skill, so trigger quality depends on it.
- Spec changes upstream become our maintenance work.
