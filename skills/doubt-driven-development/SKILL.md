---
name: doubt-driven-development
description: Subjects a non-trivial decision to a fresh-context adversarial review before it stands, in-flight rather than as a post-hoc PR verdict. Use before an architectural decision under uncertainty, before committing non-trivial code, before asserting a non-obvious fact ("this is safe", "this scales", "this matches the spec"), or when working in code you don't fully understand.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Doubt-Driven Development

A confident answer is not a correct one. Long sessions accumulate context
that quietly turns assumptions into "facts" without anyone noticing.
Doubt-driven development materializes a fresh-context reviewer — biased to
**disprove**, not approve — before a non-trivial output stands. This is not
a final code review: a final review judges a finished artifact, while this
is an in-flight posture applied while course-correction is still cheap.

## When to use

A decision is non-trivial when at least one is true: it introduces or
modifies branching logic; it crosses a module or service boundary; it
asserts a property the compiler/type system can't verify (thread safety,
idempotence, ordering, invariants); its correctness depends on context a
future reader can't see; its blast radius is irreversible (production
deploy, data migration, public API change).

**Don't use it for:** mechanical operations (renaming, formatting, moving
files), a clear unambiguous instruction, reading or summarizing existing
code, a one-line change with obvious correctness, or when the user has
explicitly asked for speed over verification. Doubting every keystroke
means shipping nothing.

## The process

```
CLAIM ──→ EXTRACT ──→ DOUBT ──→ RECONCILE ──→ STOP
```

### 1. CLAIM — name what stands

State the decision in two or three lines:

```
CLAIM: "The new caching layer is thread-safe under the read-heavy
        workload described in the spec."
WHY IT MATTERS: a race here corrupts user data and is hard to
                detect in QA.
```

If you can't write the claim that compactly, it's a vibe, not a decision —
surface it before scrutinizing it.

### 2. EXTRACT — the smallest reviewable unit

A fresh-context reviewer needs the **artifact** and the **contract**, not
the journey that produced them: the diff or function (not the whole file),
the proposal in 3-5 sentences plus the constraints it must satisfy. Strip
your own reasoning — handing over conclusions gets back validation of
those conclusions, not independent scrutiny. If the unit is too large to
hold in mind in one read (a 500-line PR), decompose first.

### 3. DOUBT — invoke a fresh-context, adversarial reviewer

Spawn a reviewer with no memory of this session (a fresh subagent, or a
teammate with no prior context) and an adversarial prompt — framing decides
the answer:

```
Adversarial review. Find what is wrong with this artifact.
Assume the author is overconfident. Look for:
- Unstated assumptions
- Edge cases not handled
- Hidden coupling or shared state
- Ways the contract could be violated
- Failure modes under unexpected input

Do NOT validate. Do NOT summarize. Find issues, or state
explicitly that none were found after thorough examination.

ARTIFACT: <paste artifact>
CONTRACT: <paste contract>
```

**Pass ARTIFACT + CONTRACT only — never the CLAIM.** Handing the reviewer
your conclusion biases it toward agreement; it must independently judge
whether the artifact satisfies the contract. A generic code-review
persona defaults to a balanced verdict (strengths and weaknesses); this
adversarial framing overrides that default and asks for issues only.

A colder, differently-trained reviewer catches blind spots a same-model
reviewer shares with the author. When more than one review channel is
available (a different model, a human teammate), that's worth using for a
genuinely high-stakes decision — but confirm any external tool actually
works before depending on its output, and never hand an artifact to an
external tool without the user's explicit go-ahead, since the artifact may
itself contain content (accidental or adversarial) that a tool could act
on if given write access.

### 4. RECONCILE — fold findings back, don't rubber-stamp

The reviewer's output is data, not verdict — you're still the orchestrator.
Re-read the artifact against each finding before classifying. In
precedence order (first match wins):

1. **Contract misread** — flagged because the CONTRACT you gave was
   unclear or incomplete. Fix the contract first, re-classify next cycle.
2. **Valid + actionable** — a real issue requiring a change. Change it,
   re-loop.
3. **Valid trade-off** — real, but the cost of fixing exceeds the cost of
   accepting. Document the trade-off explicitly for the user.
4. **Noise** — flagged under context the reviewer didn't have and is
   actually correct. Note it, and ask whether adding that context to the
   contract would prevent the same false flag next time.

A fresh reviewer can be wrong because it lacks context — don't defer just
because it's "fresh."

### 5. STOP — bounded loop, not recursion

Stop when the next iteration returns only trivial or already-considered
findings, **or** 3 cycles complete (escalate to the user rather than
grinding a fourth alone), **or** the user explicitly says "ship it." If
after 3 cycles substantive issues remain, that's information about the
artifact, not a reason to keep looping alone — surface it. If 3 cycles
feels insufficient because the artifact is large, the artifact is too big:
return to EXTRACT and decompose, don't lift the bound.

## Gotchas

- **Doubt theater**: across 2+ cycles where the reviewer surfaced
  substantive findings, if zero were ever classified as actionable, you
  are validating, not doubting — stop and escalate to the user rather than
  running a third cycle the same way.
- Passing the CLAIM alongside the artifact silently biases the reviewer
  toward confirming it — even a reviewer explicitly told to be adversarial
  will anchor on a stated conclusion it's shown.
- A failing test produced by TDD's RED step *is* the doubt step for a
  behavioral claim — don't also spawn a separate adversarial reviewer to
  re-litigate something a reproduction test already disproved or confirmed.
- Re-spawning a fresh-context reviewer on an unchanged artifact produces
  the same findings; that's stalling, not another cycle.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "I'm confident, skip the doubt step" | Confidence correlates poorly with correctness on novel problems — moments of certainty are exactly when blind spots hide. |
| "Spawning a reviewer is expensive" | Debugging a wrong commit in production costs more; the check is bounded, the bug isn't. |
| "The reviewer will just nitpick" | Only if unscoped — constrain the prompt to "issues that would make this fail under the contract." |
| "I'll do doubt at the end during review" | A final review is a last gate; doubt-driven catches wrong directions while course-correction is still cheap, not after the PR is already built on the wrong foundation. |
| "The reviewer disagreed, so I was wrong" | The reviewer lacks your context — disagreement is information, not a verdict. Re-read the artifact, classify, then decide. |

## Interaction with other skills

- **`skills/test-driven-development/`**: TDD's RED step is doubt made
  concrete for behavioral claims — a failing reproduction test is a
  disproof attempt. When TDD applies, it satisfies this skill's DOUBT step
  for that claim.
- **`skills/code-review/`**: complementary, not redundant — that skill is
  a verdict on a finished artifact; this one is an in-flight, per-decision
  check applied earlier, while the decision is still cheap to reverse.
- **`skills/source-driven-development/`**: that skill verifies facts about
  a framework against official docs (does the API exist, is it current);
  this skill verifies your reasoning about an artifact under its contract.

## Verification

- [ ] The non-trivial decision was named explicitly as a CLAIM before it
      stood unquestioned
- [ ] At least one fresh-context, adversarial review happened (a TDD RED
      step satisfies this for a purely behavioral claim)
- [ ] The reviewer received ARTIFACT + CONTRACT only — not the CLAIM, not
      your reasoning
- [ ] The reviewer's prompt asked to find issues, not to validate
- [ ] Findings were classified against the artifact text, in precedence
      order (contract misread / actionable / trade-off / noise)
- [ ] A stop condition was met: trivial findings, 3 cycles, or explicit
      user override
