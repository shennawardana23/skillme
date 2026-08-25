# Catalog-level evals

Every file under `skills/<name>/evals/` tests one skill **in isolation** —
`smeval run` installs only that skill via `harness.Build`, so a case can
prove a skill gives correct guidance, but it cannot prove Claude actually
picks that skill over a sibling when the whole catalog is available, or
that adding a new skill doesn't steal triggers from an existing one.
Those are two of the five evaluation dimensions Anthropic's own
[enterprise Agent Skills guidance](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/enterprise)
lists — *triggering accuracy* and *coexistence* — and this catalog had
zero automated coverage of either until this directory existed.

`trigger-accuracy.json` holds catalog-wide cases: `harness.BuildCatalog`
installs **every** skill under `skills/` at once (still excluding each
one's `evals/`, for the same leak-prevention reason `harness.Build`
excludes it for a single skill), then each case's prompt is graded
against the whole catalog being available — not just the one skill it's
"about."

## Running

```bash
go build -o smeval ./cmd/smeval
smeval trigger-validate catalog-evals/trigger-accuracy.json   # schema only, free
smeval trigger-run      catalog-evals/trigger-accuracy.json   # real model calls, real cost
smeval trigger-run      catalog-evals/trigger-accuracy.json -include <case-id>
```

Same flags as `smeval run` (`-primary-model`, `-fallback-model`,
`-timeout`, `-output-dir`, `-include`), plus `-skills-dir` to point at a
different catalog root. `-benchmark` is not supported here — there is no
single skill to strip out for a without_skill baseline when the whole
point is testing which skill (if any) gets picked among all of them.

## Writing a case

Same `evals.json` schema as a per-skill suite (`skill_name` is just a
label here, not a real `skills/<name>` directory — `trigger-validate`
skips the directory-name-match check `smeval validate` enforces). Ground
every case in a *real*, already-documented dividing line between two or
more actual skills in this catalog — e.g. a pair called out in a skill's
own `description` field, or in the README's
["Specialist review skills"](../README.md#specialist-review-skills)
table — rather than an invented ambiguity.

**A known sharp edge, found while writing the first cases:** a prompt
that hands over the failing test *and* the exact one-line fix needed
gives a capable model nothing left to reproduce or isolate, so it
reasonably answers directly without echoing a skill's structured output
template (e.g. `debug`'s Root Cause/Evidence/Fix/Prevention format) —
that's sensible behavior per Anthropic's own "Claude is already very
smart, don't add unneeded ceremony" guidance, not a sign the skill failed
to trigger. Testing whether a skill's structured *process* actually
engages needs a bug that genuinely requires investigation — ideally
across a small multi-file fixture the model has to search, not a
single self-contained snippet.

**Resolved:** `evalspec.Eval` now has an optional `fixture` field (a
directory name under the *skill's own* `evals/fixtures/`, not this
directory — fixtures are per-skill, not catalog-wide) — `smeval run`
copies it into the case's isolated workspace before the prompt runs, via
`harness.CopyTree`, for any command whose `runOptions.fixturesRoot` is
set (currently `run`; `trigger-run` doesn't support fixtures yet and
errors clearly rather than silently running a context-dependent case
greenfield). `debug/evals/evals.json`'s `traces-bug-across-files` case is
the worked example this note used to only propose: a real two-package Go
fixture where the panic surfaces in `pkg/server/server.go` but the actual
bug (`Timeout` never set) is in `pkg/config/config.go` — confirmed to
really reproduce (`go run .` panics at the exact line the case's prompt
quotes) before being wired up. Run live, the model used `go run .` to
reproduce, `grep` to isolate across files, and wrote a genuine
Root-Cause/Evidence/Fix/Prevention response with a real regression test —
not a lucky guess, an actual cross-file investigation. `smeval validate`
also checks a referenced fixture directory actually exists, so a typo'd
or deleted fixture fails the free schema check, not a wasted live run.

This also gives the catalog its greenfield-feasible vs. context-dependent
split, per-*case* rather than per-skill (a skill can mix both): a case
with no `fixture` is greenfield — testable in a completely empty
container, true of the overwhelming majority of existing cases — and one
with `fixture` set is context-dependent. `smeval validate` prints each
case's classification (`· <id>: greenfield` or
`· <id>: context-dependent (fixture: <name>)`).
