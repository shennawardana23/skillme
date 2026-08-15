---
name: technical-decision-records
description: This skill should be used when the user asks to "write an ADR", "document this architecture decision", "record why we chose X over Y", "create a decision record", or is capturing a significant, hard-to-reverse technical decision so future readers understand why it was made. Use for recording an already-made decision concisely — not for broad open-ended design exploration, which belongs in a design doc.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Technical Decision Records

An Architecture Decision Record (ADR), as originally proposed by Michael
Nygard in 2011, exists to answer one question six months from now: "why
did we do it this way?" It records a decision that was already made, not
an exploration of options still in flight — that distinction is what keeps
ADRs short.

## When to write one

Write an ADR when the decision is expensive to reverse, affects multiple
teams or services, or is the kind of choice someone will reasonably
question later without context (why Postgres over another datastore, why
this partitioning scheme, why this API versioning approach). Skip it for
decisions that are cheap to change or purely local to one function/file —
not every choice needs a permanent record, and an ADR backlog full of
trivial entries makes the important ones harder to find.

## The Nygard template

```
# NNNN. Title (short, describes the decision, not the problem)

## Status
Proposed | Accepted | Rejected | Deprecated | Superseded by ADR-00XX

## Context
What forces are at play — technical, business, team constraints — that
make this decision necessary. Written neutrally: state the forces, not
the conclusion.

## Decision
The decision that was made, stated actively: "We will use X."

## Consequences
What becomes easier or harder as a result — including the honest
downsides, not just the benefits. A consequences section with no
downsides is a sign the alternatives weren't seriously considered.
```

## Numbering, storage, and lifecycle

- Store under a predictable path (e.g. `docs/adr/NNNN-title.md`),
  sequentially numbered, so they're discoverable and referenceable by
  number in commit messages and PRs.
- Status moves forward, not backward: **Proposed → Accepted** (or
  **Rejected**), and later, if circumstances change, **Superseded by
  ADR-00XX** — a new ADR that supersedes an old one, rather than editing
  the old one's Decision section.
- **Never edit an Accepted ADR's Decision or Context to match reality
  after the fact.** If the decision changes, write a new ADR that
  supersedes it and link both directions. Editing history destroys the
  exact thing the record exists to preserve: what was actually decided,
  and why, at that time.
- Keep an index (a README listing number, title, status) — an ADR nobody
  can find is functionally the same as an ADR that was never written.

## Gotchas

- **An ADR is not a design doc.** A design doc explores multiple options
  broadly before a decision is made; an ADR records the decision that
  resulted, briefly, with just enough context to explain the "why." If the
  ADR is turning into a multi-page options-comparison document, that
  content probably belongs in a design doc that the ADR then references.
- **"We'll write the ADR later, once things settle" usually means never.**
  Context is freshest at decision time — write it within the same PR or
  shortly after, not as a retroactive cleanup task.
- **A Consequences section that's all upside is a red flag** that
  alternatives weren't genuinely weighed, or that the downsides are being
  hidden from future readers who'll need to know them.
- **Superseding, not editing, preserves the audit trail** — a team that
  edits old ADRs in place loses the ability to answer "what did we believe
  was true when we made this choice," which is often more valuable than
  the current state of belief.

## Real-world grounding

The ADR format originates from Michael Nygard's 2011 blog post
"Documenting Architecture Decisions," which proposed the lightweight
Context/Decision/Consequences structure specifically as a reaction against
heavyweight, rarely-read architecture documentation. It has since been
adopted widely across the industry (including being tracked as a
consistently "Adopt"-rated technique on ThoughtWorks' Technology Radar),
precisely because its brevity is the feature, not a shortcut.
