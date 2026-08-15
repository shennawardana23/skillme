# Instinct file structure (worked example)

An illustrative on-disk layout for implementing the instinct model
described in the main skill. Adapt paths and formats to your own tooling
— the structural properties (project isolation, separate observation vs.
instinct vs. evolved-artifact directories) are what matters, not these
exact names.

```
<learning-root>/
├── projects.json              # registry: project hash -> name/path/remote
├── observations.jsonl          # global observations (fallback, no project detected)
├── instincts/
│   ├── personal/               # global auto-learned instincts
│   └── inherited/              # global imported instincts
├── evolved/
│   ├── agents/                 # global generated agent definitions
│   ├── skills/                 # global generated skills
│   └── commands/                # global generated commands
└── projects/
    ├── a1b2c3d4e5f6/            # one project, identified by a stable hash
    │   ├── project.json         # id / name / root path / remote, mirrored
    │   ├── observations.jsonl   # this project's raw tool-call observations
    │   ├── instincts/
    │   │   ├── personal/        # project-specific auto-learned instincts
    │   │   └── inherited/       # project-specific imported instincts
    │   └── evolved/
    │       ├── skills/
    │       ├── commands/
    │       └── agents/
    └── f6e5d4c3b2a1/            # another project, fully isolated
        └── ...
```

## Project identification

In priority order:

1. An explicit project-directory environment variable, if your harness
   sets one.
2. The git remote URL, hashed to a stable ID — this makes the same repo
   checked out on two different machines resolve to the same project ID.
3. The repository's top-level path, as a fallback when there's no remote
   (a local-only repo) — machine-specific, but still better than pooling
   everything globally.
4. If no project can be detected at all (not in a git repo), fall back to
   the global scope rather than failing.

## Observation record shape

Each line of an observations log is one tool-call event:

```json
{"ts": "2026-08-15T09:12:03Z", "project_id": "a1b2c3d4e5f6", "tool": "Edit", "outcome": "user_corrected", "note": "reverted panic to explicit error return"}
```

Keep these terse and local — they're the raw material an observer pass
reads to propose or reinforce instincts, not something meant for direct
human reading or export.

## Minimal command surface

Whatever you build to manage this system needs, at minimum, a way to:

- list current instincts (project-scoped + global) with their confidence
- cluster related instincts into an evolved skill/command/agent
- export instincts (filterable by scope/domain) for sharing
- import instincts from someone else's export
- promote a project-scoped instinct to global scope, individually or in
  bulk against the "2+ projects, confidence >= 0.8" rule from the main
  skill
- list known projects and their instinct counts

These map to whatever interface fits your environment — a CLI, a set of
slash commands, or a background job — the operations are what matter,
not a specific implementation.
