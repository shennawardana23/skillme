# smeval evals.json — Schema Reference

`smeval` is this repository's own eval runner (`cmd/smeval`) — no
third-party dependency. It follows Anthropic's documented Agent Skills
evaluation methodology (https://agentskills.io/skill-creation/evaluating-skills):
a prompt and assertions per case, graded with concrete evidence, results
written into a workspace of `iteration-N/<case>/{with_skill,without_skill}`
directories. Where that methodology leaves grading strategy open ("give the
outputs to an LLM, or use a script for mechanical checks"), `smeval`
implements the mechanical/script side as a small built-in check DSL so
grading stays free and deterministic — see `internal/grading`.

## evals/evals.json

```json
{
  "skill_name": "go-service-idioms",
  "evals": [
    {
      "id": "error-wrapping",
      "prompt": "The exact task given to the agent under test.",
      "expected_output": "Human-readable description of what success looks like.",
      "assertions": [
        {
          "text": "Human-readable statement of what this checks",
          "check": { "contains_all": ["%w", "fmt.Errorf"] }
        }
      ],
      "timeout_seconds": 180
    }
  ]
}
```

Rules: `skill_name` must match the containing skill's directory name.
`evals[].id` must be unique within the file (it doubles as the case's
workspace directory name). Every assertion needs at least one check
condition — `smeval validate` rejects an assertion with none.

## Check conditions

| Field | Checks against | Meaning |
|---|---|---|
| `contains_all` | final response text | every string must appear |
| `contains_any` | final response text | at least one string must appear |
| `not_contains` | final response text | none of these strings may appear |
| `matches_any` | final response text | at least one Go regexp must match |
| `not_matches` | final response text | none of these Go regexps may match |
| `files_exist` | case workspace | every path (relative to the workspace) must exist |
| `file_contains` | case workspace | `{"path": "...", "contains": "..."}` — the file at `path` must contain the substring |

A `Check` may combine multiple fields; all set fields must pass for the
assertion to pass. Use `files_exist`/`file_contains` only for cases that
ask the agent to write real files (see the `scaffold-new-skill` eval in
`skills/skill-catalog-authoring/evals/evals.json` for a worked example) —
most cases in this catalog grade the conversational response text
and never touch `files_*`.

## Running

```bash
go build -o smeval ./cmd/smeval

smeval validate skills/<skill-name>              # free, no model calls
smeval run      skills/<skill-name>               # with_skill only (default)
smeval run      skills/<skill-name> -benchmark    # also runs without_skill, writes benchmark.json
```

Flags for `run` (may appear before or after the skill directory argument):

| Flag | Default | Meaning |
|---|---|---|
| `-primary-model` | `sonnet` | Model for the first attempt |
| `-fallback-model` | `opus` | Model for smeval's own outer retry, and passed to claude's native `--fallback-model` |
| `-timeout` | `3m0s` | Per-attempt timeout |
| `-benchmark` | `false` | Also run each case without the skill installed |
| `-output-dir` | `<skill-dir>-workspace` | Workspace root |
| `-include` | (none) | Only run eval IDs containing this substring |

## How grading works (no LLM judge, by design)

`smeval` never spends a second model call grading a case — every assertion
is a deterministic string/regex/file check evaluated in Go
(`internal/grading`). This mirrors the "prefer a verification script over
LLM judgment for mechanically-checkable assertions" principle from the
reference methodology. If a future case genuinely needs semantic judgment
(tone, holistic quality) that cannot be reduced to a check, that is a
signal to extend the `Check` struct deliberately — not to bolt on an
LLM-judge escape hatch by default.

## Provider/model fallback (how it actually works)

Two layers, and they are not redundant:

1. **Native, in-process** — `smeval` passes `--fallback-model` straight
   through to the `claude` CLI, which already retries within the same
   invocation when the primary model is overloaded or unavailable. This is
   verified behavior, not assumed: an unrecognized `--model` name still
   returns `is_error:false` when `--fallback-model` is set.
2. **Outer, cross-process** (`internal/engine.Run`) — engaged only when
   the whole invocation fails in a way the native layer cannot see: the
   process exits non-zero, the timeout is exceeded, or stdout is not
   parseable JSON. On any of these, `smeval` retries once, in a fresh
   process, using the fallback model as the new primary. A run that
   completes with `is_error:false` is never retried for content reasons —
   failing assertions are a grading outcome, not an engine failure.

See `internal/engine/engine_test.go` for the test that proves the outer
layer engages using a fake `claude` binary, deterministically, without
depending on a real provider actually failing.

## Workspace and harness isolation

Each case runs in its own throwaway plugin directory (`internal/harness`)
containing only the skill under test — never the other skills in this
catalog, and never that skill's own `evals/` directory (which would leak
grading criteria into the model's context). Each `with_skill`/`without_skill`
run also gets its own workspace directory
(`iteration-N/<case>/<configuration>/workspace/`) so a case that writes
files does so in isolation, never into the repository working tree.

## CI wiring

`.github/workflows/skill-eval.yml` builds `smeval` from source, installs
the `claude` CLI, discovers every `skills/*/evals/evals.json` dynamically,
and runs one job per skill. Adding a skill following the layout above
requires no workflow changes.
