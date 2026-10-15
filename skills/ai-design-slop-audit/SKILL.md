---
name: ai-design-slop-audit
description: This skill should be used when the user wants an evidence-backed review of an existing UI for generic "AI-generated" design defaults — asking to "audit this design", "why does this look AI-generated", "review this for design slop", or wants a removal-first cleanup plan without a speculative redesign. Point-in-time audit of a screenshot, rendered page, or design code producing a prioritized, cite-the-evidence findings table. Does not edit anything and is not a build-time guide — see frontend-ui-engineering's "Avoiding the generic AI aesthetic" section for the rules to follow while building instead of auditing after the fact.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "web-quality"
---

# AI Design Slop Audit

A diagnostic audit, not a redesign. Find the smallest set of removals or
corrections that would improve the interface while keeping its existing
direction. This complements `frontend-ui-engineering` the way
`code-reviewer` complements writing code: that skill says what to do
while building; this one grades what already shipped.

## Boundaries

- Don't guess whether AI made the design — judge the visible result.
- Don't assign a numeric "slop score."
- Don't reject a technique in isolation. A gradient, serif, dark theme,
  glass effect, card, or animation can be a deliberate, working choice.
- Don't prescribe a new font, palette, layout, or design system unless
  asked for one.
- Don't turn the audit into an implementation task.
- Mark anything outside the visible/provided evidence as unknown rather
  than assuming it.

## Gather evidence

Inspect whatever the request actually provides: screenshots or
recordings, the rendered page at relevant viewports, interaction/state
changes, and source (markup, tokens, styles, copy). Every finding cites a
concrete location — a component, a behavior, a line of copy. A generic
tendency with no evidence in the artifact isn't a finding.

## Classify each finding

- **Quality defect** — an established usability, accessibility, content,
  responsive, or runtime problem.
- **Slop pattern** — a repeated default, decorative layer, or
  generated-looking convention that carries no useful role here.

## What to look for

**Decorative stacking** — glow, gradient text, glass, borders, shadows,
grid backgrounds, particles, and beams combined so several layers do the
same decorative job or compete with the content.

**Template repetition** — the same card treatment on nearly every
content block; repeated icon-heading-description tiles with
interchangeable content; nested rounded containers that don't
communicate hierarchy; a generic landing-page section order unrelated to
the actual product's usage journey.

**Typography and copy clutter** — duplicate text, empty superlatives and
generated-sounding filler, decorative type that obscures hierarchy, long
centered paragraphs that are hard to scan.

**Motion theater** — motion that delays access, repeats mechanically,
moves a target the user is trying to click, or ignores
`prefers-reduced-motion`. Motion that communicates state, causality, or
spatial change is not slop — preserve it. (For a deeper check of
animation craft specifically, hand the diff to `animation-review`.)

**Fabricated proof** — invented or unverifiable metrics, testimonials,
logos, ratings, or activity feeds presented as real evidence.

**Established UI failures** — unclear or competing primary actions, low
contrast, clipping/overflow, broken responsive behavior, information
available only on hover, missing states for the observed flow, spacing
or color inconsistency that reads as accidental rather than expressive.

## Removal test

For every candidate:

1. State what job it does — information, state, action, hierarchy, or
   brand meaning.
2. Ask whether removing it improves clarity without losing that job.
3. If yes, recommend removal or consolidation.
4. If no, recommend the smallest correction using the system already in
   place.
5. Only suggest a replacement when deletion would create a real loss —
   default to subtraction.

## Prioritize and report

- **P0** — blocks task completion, a severe accessibility issue, or
  deceptive proof.
- **P1** — materially harms comprehension, trust, or navigation.
- **P2** — repeated slop or inconsistency that weakens hierarchy/identity.
- **P3** — minor polish, limited impact.

Return the 5–8 highest-impact findings by default; group repeated
instances into one systemic finding rather than listing each occurrence.

```md
## Verdict
One paragraph: the dominant problem, and what to remove first.

## Checked scope
Artifact, screen, state, and viewport actually inspected.

## Findings
| Priority | Class | Pattern | Evidence | Harm | Remove or fix |
|---|---|---|---|---|---|

## Unknowns
States or behavior that couldn't be verified from the given evidence.
```

Keep the report compact; omit empty sections. Lead with evidence, not
taste. Recommend removal before restyling, and end with the single
change that would produce the largest improvement.
