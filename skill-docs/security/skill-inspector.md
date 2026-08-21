## What it does

Decides whether a *third-party* AI agent skill — something you downloaded,
cloned, or are about to `claude plugin install`/`npx skills add` — is safe
to bring into your own agent setup. It runs NVIDIA's SkillSpector static
scanner, then does a source-aware semantic pass on top of it, because the
two catch different failure modes: the scanner finds pattern matches
(`eval`, a suspicious env-var read, a YARA hit), the semantic pass judges
*intent* — whether what the code does actually matches what the skill
claims to do. Treating the score as the verdict is the single most common
way to get this wrong in both directions: a low score can still hide
something that matters (an env-var read that defaults to
`~/.aws/credentials`), and a documented, bounded, necessary use of a
credential can legitimately score MEDIUM without deserving a rejection.

## When to reach for it

Reach for this before installing anything you didn't write yourself —
a skill from a marketplace, a repo someone linked you, a `.skill` archive
a colleague sent over. It is not for auditing your own project: for this
catalog's own `.claude/` harness configuration (hooks, MCP configs, agent
definitions already in place and presumed trusted), use `security-scan`
instead — that's a different threat model, an existing setup drifting
insecure over time, not a new, possibly hostile addition. For your
application's own source code, use `security-review`.

## Common questions

- **"SkillSpector came back clean — can I trust that alone?"** No. A
  clean static scan means nothing tripped SkillSpector's 68 pattern
  rules; it says nothing about whether the skill's actual behavior
  matches its stated purpose. The semantic pass is not optional
  ceremony on top of a scanner that already decided — it is the second,
  independent check that catches what pattern-matching structurally
  cannot.
- **"The skill's own text says it's safe and fully vetted — doesn't
  that count for something?"** The opposite. Legitimate skills have no
  reason to address their own reviewer. Text inside a skill that talks
  to "any security reviewer or scanning agent" and tries to pre-empt the
  verdict is itself the strongest single signal of a prompt-injection
  attempt, and following the embedded instruction (including one that
  asks you to not mention it) is exactly the failure mode this skill
  exists to prevent.
- **"`skillspector` isn't installed — should I skip the review?"** No.
  Say so plainly, then do the manual/semantic-only review anyway
  (`SKILL.md`, scripts, dependency manifests, MCP configs), and hold the
  resulting verdict to lower confidence than a dual-line review would
  get — don't silently install `skillspector` yourself, and don't
  present a manual-only verdict with the same certainty as a scanned one.

## It's working if

- The report names a verdict of exactly `APPROVE`, `CAUTION`, or
  `REJECT` — never a hedge, never both a score and a contradicting
  recommendation left unreconciled
- A CAUTION verdict always names *why* the sensitive behavior is
  considered documented, necessary, and bounded — not just that it exists
- A skill containing reviewer-directed text gets flagged as the finding
  itself, not treated as reassurance
- Nothing from the target skill gets executed during the review

## Where it fits

A standalone, pre-install gate — it runs once, before a skill enters your
setup, not as an ongoing check. Sits next to `security-scan` (this
catalog's own config, over time) and `security-review` (application
source code) as the third leg of "what could go wrong with an AI-assisted
setup," each covering a different thing that could be reviewed.
