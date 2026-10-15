# 0002. Build our own eval runner with deterministic grading, no LLM judge

## Status
Accepted (recorded retroactively)

## Date
2026-08-15 (initial public release, commit `79b9b5b`)

## Context
The catalog's claim is that its skills are checked, not just written. That needs a way to run each skill
against a real model and grade the answer repeatably, cheaply enough to do across 100+ skills, and with
evidence a reviewer can read. Evidence: `cmd/smeval`, `internal/grading` ("computed deterministically
rather than by an LLM judge"), `skills/skill-catalog-authoring/references/smeval-schema.md`
("no LLM judge, by design").

## Decision
We maintain `smeval`, a Go CLI with no third-party eval dependency. It validates each skill's
`evals/evals.json`, runs every case through the `claude` CLI, and grades responses with deterministic
string, regex and file checks, quoting the evidence. Human judgment is captured separately in
`feedback.json`, enforced by `smeval feedback-check`.

## Alternatives considered
- **LLM-as-judge grading**: handles nuance, but costs a second model call per case, varies between runs
  and is hard to audit. Rejected for repeatability and cost.
- **An existing eval framework**: less code to own, but adds a dependency and does not fit the
  one-skill-per-workspace isolation we need (ADR-0004). Rejected.
- **Manual review only**: catches nuance, but does not scale to a whole catalog and leaves no
  machine-checkable record. Rejected as the only method; kept as the `feedback.json` step.

## Consequences
- Grading is free, fast and repeatable; CI can validate everything without a model.
- Assertions can be wrong in their own way. The 2026-08-17 sweep commit (`e1028b9`) fixed 41 evals whose
  failures were eval-authoring bugs (too-narrow matches, markdown breaking substring checks, over-broad
  negative checks); the authoring rules in `.claude/rules/skill-authoring.md` exist because of that.
- Deterministic checks cannot judge prose quality or "correct but misses the point", which is why the
  human review step is mandatory.
- We own the runner's bugs too: one let a case write into the real repo (fixed, see `TESTING.md`).
- Live runs use the `claude` CLI only, so results describe Claude models, not other providers.
