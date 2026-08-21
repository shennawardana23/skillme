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

```text
⏳ safe-query-write: Binds hotel_id as an actual query parameter alongside status, not just...
   ✅ PASS safe-query-write: 3/3 assertions
```

or on failure:

```text
   ❌ FAIL safe-query-write: 2/3 assertions
```

A styled version of the same summary lands at
`smeval-workspace/runs/<name>/iteration-N/report.html` — open it in a browser
rather than re-reading the terminal scrollback. `iteration-N`
auto-increments each run, so nothing overwrites a previous result.

### Diagnosing a failure

Don't stop at "it failed" — read the actual evidence before concluding
whether the skill is wrong or the eval case is wrong:

```bash
cat smeval-workspace/runs/<name>/iteration-1/<case-id>/with_skill/outputs/response.md   # what the model actually said
cat smeval-workspace/runs/<name>/iteration-1/<case-id>/with_skill/grading.json          # exactly which assertion failed, and why
```

`grading.json`'s `evidence` field quotes exactly what was/wasn't found —
e.g. `"missing required substring(s): PARTITION BY LIST"`. If the model's
answer looks substantively correct despite that, check whether it wrote
the real deliverable to a file in
`smeval-workspace/runs/<name>/iteration-1/<case-id>/with_skill/workspace/`
instead of the chat response — grading a case whose real output is a
file, but whose assertions only check the chat text, is the single most
common eval-authoring mistake in this catalog so far. Fix by pointing the
assertion at the file (`files_exist` / `file_contains`) instead of
widening the text match.

### Re-verifying an assertion fix without a new API call

If the model's answer was substantively correct and only the assertion's
`contains_any`/`matches_any` list needed broadening (the common case — see
above), you don't need to re-run the case to prove the fix works. The
model's output is still sitting in `response.md`; only the grading logic
changed. Re-grade the saved response against the edited `evals.json`
directly, using `internal/grading.Grade` from a throwaway `main.go` (see
git history for `cmd/regrade-scratch` in past commits for the ~40-line
pattern: `evalspec.Load` the skill's `evals.json`, read the saved
`response.md`, call `grading.Grade`, print `AssertionResults`). This is
free, instant, and — because it grades the *exact same* model output the
original failure was diagnosed against — a stronger proof than a fresh
`-include` run, which would introduce new model-output variance on top of
the assertion change you're trying to isolate. Delete the throwaway tool
when done; it's not part of the shipped catalog.

### Iterating on one case

```bash
./smeval run skills/<name> -include <case-id>
```

Runs only that case — much faster than re-running the whole skill while
fixing one broken assertion.

### Comparing against no skill at all

```bash
./smeval run skills/<name> -benchmark
cat smeval-workspace/runs/<name>/iteration-1/benchmark.json
```

Also runs every case with the skill absent entirely, and writes a
side-by-side comparison — [`examples/benchmark.sample.json`](examples/benchmark.sample.json)
is real output from a run against this catalog, reproduced here for
reference:

```json
{
  "with_skill": {
    "pass_rate_mean": 1,
    "time_seconds_mean": 36.71,
    "tokens_mean": 92870.66666666667
  },
  "without_skill": {
    "pass_rate_mean": 0.8333333333333334,
    "time_seconds_mean": 37.186,
    "tokens_mean": 63351.666666666664
  },
  "delta": {
    "pass_rate_mean": 0.16666666666666663,
    "time_seconds_mean": -0.4759999999999991,
    "tokens_mean": 29519.000000000007
  }
}
```

Reading it: the skill won the one case the unskilled run missed
(`delta.pass_rate_mean` > 0, the signal that actually matters), took
about the same wall-clock time either way, and spent ~29.5k more tokens
doing it — the skill's own guidance costs tokens to read and follow, so a
positive `tokens_mean` delta alongside a positive `pass_rate_mean` delta
is the expected, healthy shape. A positive token delta with a *zero or
negative* pass-rate delta is the shape worth investigating — you paid
tokens for a skill that didn't change the outcome.

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

## Reviewing results with a human

Every `smeval run` scaffolds `iteration-N/feedback.json` — one empty-string
entry per case ID, e.g.:

```json
{
  "generate-prfaq-from-vague-idea": "",
  "critique-prfaq-missing-customer-quote-and-faq": "",
  "recommend-lighter-prd-over-prfaq": ""
}
```

This exists because assertion grading only checks what someone thought to
write an assertion for — [the spec's own framing](https://agentskills.io/skill-creation/evaluating-skills#reviewing-results-with-a-human)
for exactly this file. Some qualities are real but not decomposable into a
pass/fail check: whether a press release's headline is actually punchy,
whether a report's structure reads well, whether an output is technically
correct but misses the point. A human (or a separate qualitative-review
pass) opens each case's `outputs/response.md` (or the real deliverable
under `outputs/workspace/` if the case writes a file), and fills in that
case's entry — specific and actionable, not "looks bad."

Leaving an entry as `""` means "reviewed, it held up" — not "not yet
looked at." Don't skip the review because assertions all passed; a real
example from this catalog: `writing-product-requirements`'s
`generate-prfaq-from-vague-idea` case passed 3/3 assertions (has a
headline, has a customer quote, has both FAQ sections) on a run that
still had a genuine craft problem invisible to any of those three checks
— a 20-word press-release headline (real PR-FAQ headlines run ~10-12
words) and two raw `[Open question]` placeholder tokens left inline in
customer-facing FAQ copy instead of a provisional answer with the caveat
in prose. Both are exactly the "technically correct, still worth
tightening" gap assertions structurally can't express — caught only by
actually reading the output.

`feedback.json` is never overwritten if it already exists in that
iteration directory, so filling it in is safe to do at any point before
moving to the next iteration.

## Testing every skill in the catalog

```bash
scripts/test-all-skills.sh              # test every skill, resuming from prior results
scripts/test-all-skills.sh -fresh       # discard prior results, test everything again
scripts/test-all-skills.sh -only=<name> # just one skill
```

This is slow and makes real API calls — one skill at a time, real model
calls per case, at catalog scale that's real time and real cost. It's
resumable: results accumulate in
`smeval-workspace/test-all-results/summary.txt`, and a skill already
recorded there is skipped on the next invocation, so an interrupted run
just needs to be re-invoked, not restarted from scratch. The script prints
a "Needs attention" list of every skill that errored or partially failed
when it finishes (or is interrupted and re-run to completion).

An `ERROR` result (as opposed to a partial pass) means the engine itself
failed — a timeout, a non-zero exit, unparseable output — not that an
assertion failed. Check
`smeval-workspace/runs/<name>/iteration-N/<case-id>/with_skill/outputs/engine-error.txt`
for the reason before assuming anything about the skill's content; a
timeout on a case with an unusually long, thorough response may just need
a longer `timeout_seconds` in that case's `evals.json` entry, not a
content fix.

## A fixed workspace-escape bug: the isolated case wrote to the real repo

Found during this catalog's own live-test sweep, on
`skill-catalog-authoring`'s `scaffold-new-skill` case, whose prompt asks
the model to write real files "under `skills/http-retry-policy/`". The
case's own isolated workspace came back empty, but a real, untracked
`skills/http-retry-policy/` directory appeared in this actual repository
— matching the case's expected content exactly. The model didn't
misbehave; the isolation had a hole.

Root cause: each case's workspace (`smeval-workspace/`-adjacent
`smeval-workspace/runs/<name>/iteration-N/<case-id>/.../workspace/`) is a real
directory *inside this repo's own git working tree* (gitignored, but
still inside the tree — there is only one `.git`, at this repo's root).
`internal/engine` sets the subprocess's OS-level working directory
correctly via `cmd.Dir`, but a relative path that doesn't already exist
under the workspace and *does* exist at the repo's real root — like
`skills/http-retry-policy/`, which looks exactly like a normal path in
this catalog — can still resolve there if anything in the write path
trusts a discovered project/git root over the literal process cwd. Only
one case in the whole catalog has a prompt shaped like this (checked via
`grep -rln 'under skills/\|to skills/\|in skills/' skills/*/evals/evals.json`),
so the blast radius was narrow, but the mechanism wasn't.

Fixed two ways in the same commit, both defensible independently, proven
together by re-running the case live and confirming the real repo's
`skills/http-retry-policy/` did **not** reappear while the file correctly
landed inside the isolated workspace (4/4 assertions passed):

1. `internal/engine` now overrides the subprocess's inherited `PWD`
   environment variable to match `cmd.Dir` — `cmd.Dir` changes the actual
   OS working directory but Go does not touch the inherited env to
   match, and some tools trust `process.env.PWD` over `process.cwd()`.
2. `cmd/smeval/main.go` now runs `git init -q` inside every case
   workspace right after creating it, so the workspace has its own `.git`
   and nothing walking up looking for a project root can walk past it
   into the real repo.

If you ever see a live-run case write a file somewhere in this actual
repo instead of its workspace, this is the failure mode to suspect first
— check `git status --porcelain` for anything untracked right after a
suspicious run, before concluding the skill's content is wrong.

## A known, unfixed confound: account-level instruction bleed

`--setting-sources project,local` (see above) excludes the *local*
`user` settings tier — the tester machine's own `~/.claude/CLAUDE.md` and
project-independent hooks. It does **not** exclude anything tied to the
Claude account the `claude` CLI is authenticated as. If that account has
an organization-level instruction configured (e.g. a standing preference
to ask clarifying questions before starting an open-ended task), the
isolated subprocess can still inherit it — because it isn't a local
settings file `--setting-sources` has any control over.

Observed several times in this catalog's live-test sweep: an open-ended
"design this" / "write this from scratch" prompt got a page of
clarifying questions back instead of a direct answer, once with the
model explicitly citing "your org's preference to nail down context
early" in its own response — direct evidence of the mechanism, not a
guess. Two of those cases produced a normal, direct answer on a single
immediate re-run — ordinary soft-bias variance, no fix needed.

One case (`writing-product-requirements`'s `generate-prfaq-from-vague-idea`)
reproduced on *two* consecutive re-runs, meaning it isn't always a
one-shot fluke — for some prompt shapes the bias is sticky enough that a
plain re-run won't clear it. Do not "fix" a case like this by broadening
its assertions to accept clarifying questions as a pass — that hides a
real methodological gap instead of documenting it. What *is* legitimate,
and already an established pattern in this catalog (see
`test-driven-development`'s `prove-it-bug-fix-pattern` prompt), is adding
an explicit "don't ask clarifying questions — draft with reasonable
assumptions, flag them as open questions in the output itself" line to
the case's own prompt, so the eval reliably tests what the skill teaches
instead of testing whether this particular reflex fires this run. That
line closed `generate-prfaq-from-vague-idea` immediately, and revealed
what the model does once it stops deflecting is worth checking too —
that one also then wrote its answer to a file, the same file-vs-response
mistake described above, which needed its own fix on top.

If you hit this, re-run once before concluding anything; if it
reproduces, add the explicit anti-deflection line to that case's own
prompt rather than to the skill or to a shared default — this is a
per-case escape hatch, not a catalog-wide policy change.

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
