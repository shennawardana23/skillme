# Testing a skill

Two layers, and they check different things — running the first does not
mean the second would pass:

1. **Schema validation** — free, no model calls. Proves the files are
   well-formed. Does not prove the skill actually works.
2. **Live eval** (`smeval run`) — calls the real `claude` CLI and grades
   its actual output. This is the one that catches real bugs.

Every batch of live-tested skills in this catalog's history has turned up
real bugs that schema validation missed — a grading assertion checking the
chat response when the model correctly wrote its answer to a file, a regex
Go's engine can't compile, an assertion phrased too narrowly to match a
correct-but-differently-worded answer. Schema-clean is not the same as
working. Run the live eval before trusting a skill, and before opening a
PR that touches one.

## 1. Schema validation

```bash
uvx --from skills-ref agentskills validate skills/<name>   # official Agent Skills spec
./smeval validate skills/<name>                              # this catalog's evals.json schema
```

Catches: frontmatter spec violations (bad `name`, missing `description`,
a stray top-level `version:` instead of `metadata.version`), a
`skill_name` in `evals.json` that doesn't match the directory, malformed
JSON, an assertion with no check condition. Catches nothing about whether
the skill's content is actually correct or whether the model follows it.

## 2. Live eval

```bash
go build -o smeval ./cmd/smeval   # once, or after pulling internal/ changes
./smeval run skills/<name>
```

Needs either an existing `claude` CLI login or `ANTHROPIC_API_KEY` set.
Each case in the skill's `evals/evals.json` gets sent to a real,
freshly-isolated Claude Code invocation (the skill under test is the
*only* skill visible — see `internal/harness`) and the response is graded
against that case's assertions.

### Reading the output

Console output while it runs:

```
⏳ safe-query-write: Binds hotel_id as an actual query parameter alongside status, not just...
   ✅ PASS safe-query-write: 3/3 assertions
```

or on failure:

```
   ❌ FAIL safe-query-write: 2/3 assertions
```

A styled version of the same summary lands at
`skills/<name>-workspace/iteration-N/report.html` — open it in a browser
rather than re-reading the terminal scrollback. `iteration-N`
auto-increments each run, so nothing overwrites a previous result.

### Diagnosing a failure

Don't stop at "it failed" — read the actual evidence before concluding
whether the skill is wrong or the eval case is wrong:

```bash
cat skills/<name>-workspace/iteration-1/<case-id>/with_skill/outputs/response.md   # what the model actually said
cat skills/<name>-workspace/iteration-1/<case-id>/with_skill/grading.json          # exactly which assertion failed, and why
```

`grading.json`'s `evidence` field quotes exactly what was/wasn't found —
e.g. `"missing required substring(s): PARTITION BY LIST"`. If the model's
answer looks substantively correct despite that, check whether it wrote
the real deliverable to a file in
`skills/<name>-workspace/iteration-1/<case-id>/with_skill/workspace/`
instead of the chat response — grading a case whose real output is a
file, but whose assertions only check the chat text, is the single most
common eval-authoring mistake in this catalog so far. Fix by pointing the
assertion at the file (`files_exist` / `file_contains`) instead of
widening the text match.

### Iterating on one case

```bash
./smeval run skills/<name> -include <case-id>
```

Runs only that case — much faster than re-running the whole skill while
fixing one broken assertion.

### Comparing against no skill at all

```bash
./smeval run skills/<name> -benchmark
cat skills/<name>-workspace/iteration-1/benchmark.json
```

Also runs every case with the skill absent entirely, and writes a
side-by-side comparison:

```json
{
  "with_skill":    { "pass_rate_mean": 1,    "time_seconds_mean": 9.7,  "tokens_mean": 86647 },
  "without_skill": { "pass_rate_mean": 0.89, "time_seconds_mean": 20.4, "tokens_mean": 72990 },
  "delta":         { "pass_rate_mean": 0.11, "time_seconds_mean": -10.6, "tokens_mean": 13657 }
}
```

If `delta.pass_rate_mean` is at or near zero, the skill isn't adding
measurable value for the cases you've written — either the model already
knows this without help, or the eval cases aren't actually probing the
thing the skill is supposed to teach. Both are worth investigating rather
than shipping a skill on faith.

**A single `-benchmark` run is one sample, not a verdict.** Real model
output varies run to run — a case that shows 3/3 without the skill once
and 2/3 the next run is ordinary variance, not a regression to chase.
Look for a consistent pattern across a couple of runs before concluding
anything from a single number.

## Provider/model fallback

`smeval run` accepts `-primary-model` (default `sonnet`) and
`-fallback-model` (default `opus`), exercised at two independent layers:
`claude`'s own native `--fallback-model` (recovers within one process
invocation) and `smeval`'s outer retry in `internal/engine` (a fresh
process attempt, only on a whole-process failure — timeout, non-zero
exit, unparseable output — never on a graded assertion failure). See
`internal/engine/engine_test.go` for both exercised against a fake
`claude` binary.

## What CI does and doesn't check

`.github/workflows/skill-eval.yml` runs schema validation
(`agentskills validate` + `smeval validate`) on every push/PR, for free.
It does **not** run the live eval — re-running every skill's live model
call on every push would mean real API cost and long runtime for even a
one-line change to a single skill, at this catalog's size. Running the
live eval before opening a PR is on you; see the Contributing section in
`README.md`.
