---
description: Validate and run a skill's eval suite with smeval, then summarize the result
---

Run the full eval loop for skill `$ARGUMENTS` and report back:

1. `go build -o smeval ./cmd/smeval` (skip if `./smeval` is already newer
   than `cmd/smeval/` and `internal/`).
2. `uvx --from skills-ref agentskills validate skills/$ARGUMENTS`
3. `./smeval validate skills/$ARGUMENTS`
4. `./smeval run skills/$ARGUMENTS -benchmark`
5. Report pass/fail per case with the quoted assertion evidence, the
   `report.html` path, and the `feedback.json` path — tell the user to
   open and fill in `feedback.json` before treating the run as done (see
   `TESTING.md`'s "Reviewing results with a human" section).

If `$ARGUMENTS` is empty, ask which skill (a directory under `skills/`)
to run before doing anything else.
