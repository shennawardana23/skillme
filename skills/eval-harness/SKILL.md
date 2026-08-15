---
name: eval-harness
description: Sets up eval-driven development (EDD) for an agent-assisted feature — define capability and regression evals before coding, grade them with code/model/human graders, and track reliability with pass@k and pass^k metrics. Use when defining pass/fail criteria for an AI coding task before implementation starts, measuring agent reliability across repeated attempts, or building a regression suite that guards prompt or agent changes.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Eval Harness

Eval-Driven Development treats evals as the unit tests of AI-assisted work: define expected behavior *before* implementation, run evals continuously during development, track regressions with every change, and use pass@k for reliability rather than a single pass/fail.

## Eval types

**Capability evals** — test whether the agent can do something new:
```
[CAPABILITY EVAL: feature-name]
Task: description of what the agent should accomplish
Success Criteria:
  - [ ] Criterion 1
  - [ ] Criterion 2
Expected Output: description of the expected result
```

**Regression evals** — confirm a change didn't break existing behavior:
```
[REGRESSION EVAL: feature-name]
Baseline: SHA or checkpoint name
Tests:
  - existing-test-1: PASS/FAIL
  - existing-test-2: PASS/FAIL
Result: X/Y passed (previously Y/Y)
```

## Grader types

1. **Code grader** — deterministic checks: `grep -q "expected pattern" file && echo PASS`; `go test ./... -run TestAuth`; `go build ./...`. Prefer this whenever the success criterion is mechanically checkable.
2. **Model grader (LLM-as-judge)** — for open-ended output: give the model a rubric ("does it solve the stated problem? is it well-structured? are edge cases handled?"), have it emit a 1-5 score with reasoning. Use only where a code grader genuinely can't express the criterion.
3. **Human grader** — flag for manual review when risk or ambiguity is too high to automate:
```
[HUMAN REVIEW REQUIRED]
Change: what changed
Reason: why human review is needed
Risk Level: LOW/MEDIUM/HIGH
```

Order of preference: code grader > model grader > human grader. Deterministic beats probabilistic; probabilistic beats un-reviewed.

## Metrics

- **pass@k** — "at least one success in k attempts." pass@1 is first-attempt reliability; pass@3 is success within 3 tries. Typical target: pass@3 > 90%.
- **pass^k** — "all k trials succeed," a higher bar for reliability. Use for critical paths where a single silent failure is unacceptable (auth, payments, data migrations): pass^3 = 100% means three consecutive clean runs, not one lucky one.

Recommended thresholds: capability evals pass@3 ≥ 0.90; regression evals pass^3 = 1.00 for release-critical paths.

## Workflow

**1. Define (before coding)** — write capability and regression evals, with explicit success criteria, before writing implementation code. This forces clear thinking about what "done" means and surfaces ambiguity while it's still cheap to resolve.

**2. Implement** — write code to pass the defined evals; don't let the evals drift to match whatever the implementation happened to produce.

**3. Evaluate** — run each capability eval, run the regression suite, record pass/fail per eval, compute pass@k across repeated attempts if reliability (not just correctness) is the question.

**4. Report:**
```
EVAL REPORT: feature-xyz
Capability Evals:
  create-user:     PASS (pass@1)
  hash-password:   PASS (pass@1)
  Overall:         2/2 passed
Regression Evals:
  login-flow:      PASS
  Overall:         1/1 passed
Metrics:
  pass@1: 67% (2/3 across repeated attempts)
  pass@3: 100%
Status: READY FOR REVIEW
```

## Storage

```
.claude/
  evals/
    feature-xyz.md      # eval definition
    feature-xyz.log     # run history
    baseline.json        # regression baselines
```

Version evals with the code they test — an eval file is a first-class artifact, not disposable scratch work.

## Anti-patterns

- **Overfitting prompts to known eval examples** — an agent (or its prompt) tuned until it passes the exact eval cases stops measuring general capability and starts measuring memorization of the eval set.
- **Measuring only happy-path outputs** — an eval suite with no adversarial or edge-case entries systematically overstates reliability.
- **Ignoring cost/latency drift while chasing pass rate** — a change that raises pass@3 from 90% to 95% while tripling token cost or latency is not an unambiguous win; report it alongside the pass rate, not instead of it.
- **Allowing flaky graders in release gates** — a model grader with high run-to-run variance blocking or passing releases non-deterministically is worse than no gate, because it trains the team to re-run until green instead of investigating.

## Gotchas

- pass@k is not the same statistic as "run it k times and see if any passed" when k is small relative to the true success rate — the standard unbiased estimator (as introduced for HumanEval) samples n ≥ k completions and computes an expectation, because naively running exactly k trials produces a biased, higher-variance estimate. If precision matters (comparing two prompts' reliability), sample more than k and use the unbiased formula rather than a single k-length run.
- A code grader checking for a string pattern (`grep -q "export function handleAuth"`) can pass on a refactor that renames the function's behavior but not its declaration line — a grader that only checks surface syntax is validating that code was written, not that it does what the eval intended; pair syntax checks with a behavioral check (an actual test run) wherever the criterion is really about behavior.
- pass^k for a "critical path" eval is meaningless if the k runs aren't independent — three consecutive runs against a warm cache, a memoized model response, or a test harness that leaks state between runs will report pass^3 = 100% while the underlying reliability is unmeasured. Confirm each trial starts from a clean, independent state before trusting the metric.

## Real-world grounding

pass@k as a metric originates from OpenAI's 2021 Codex paper ("Evaluating Large Language Models Trained on Code," Chen et al.), which introduced the unbiased sampling estimator specifically because naive k-trial pass/fail was too noisy to compare models reliably on HumanEval — the same estimator and the same "don't just run it k times" caveat apply directly to evaluating agentic coding tasks today.

## Verification

- [ ] Capability and regression evals were written before implementation started, not reverse-engineered from what the code does
- [ ] Every eval has an explicit success criterion, not just a description of the task
- [ ] Code graders are used wherever the criterion is mechanically checkable; model/human graders only where it isn't
- [ ] pass@k or pass^k is reported with the sample count it was computed from, not a bare percentage
- [ ] The eval suite includes at least one adversarial or edge-case entry, not only happy-path tasks
