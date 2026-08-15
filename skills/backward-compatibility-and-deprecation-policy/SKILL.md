---
name: backward-compatibility-and-deprecation-policy
description: Use when the user asks to "deprecate this endpoint/field", "set a sunset date", "plan a breaking change rollout", "write a deprecation notice", "how long should we support the old version", or is retiring an API version, a config format, a CLI flag, or an internal library's old interface. Guides the organizational process and timeline for deprecating something already shipped — not the contract-design question itself (see the api-design skill for that).
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Backward Compatibility and Deprecation Policy

This skill is about the *process* of retiring something that already
shipped — the timeline, the communication, the enforcement mechanism. The
question of how to design a contract so it doesn't need this process as
often (additive changes, versioning, error-shape consistency) is covered
by the `api-design` skill; read that one when you're designing the
contract itself, this one when you're retiring something already in use.

## Hyrum's Law sets the baseline expectation

> "With a sufficient number of users of an API, it does not matter what
> you promise in the contract: all observable behaviors of your system
> will be depended on by somebody." (Hyrum Wright, popularized via the
> *Software Engineering at Google* book)

This means: **the deprecation process has to assume real callers depend on
behavior you never intended to guarantee**, not just the documented
contract. A deprecation plan built only around "we announced it in the
changelog" underestimates how many callers won't see that changelog before
their integration breaks.

## Procedure: deprecating something already in production use

1. **Determine actual usage before announcing anything.** Instrument the
   old path (a log line, a metric counter, a response header) so you know
   who's still calling it — announcing a sunset date for something with
   unknown callers is how outages happen. If you can't measure usage,
   that's the first gap to close, not a reason to skip measurement.
2. **Publish the deprecation with three concrete pieces of information**:
   what's being replaced, what to migrate to, and a specific sunset
   date — "this will be removed eventually" is not a deprecation notice,
   it's a warning with no actionable deadline.
3. **Mark it as deprecated in the artifact itself**, not only in
   documentation: a `Deprecation` HTTP response header, a compiler
   deprecation annotation (`// Deprecated:` Go doc comment,
   `@deprecated` in other ecosystems), a CLI flag that prints a warning
   when used. Documentation-only deprecation is invisible to someone who
   never reads docs but does read a warning printed in front of them.
4. **Give a grace period sized to your actual caller population and their
   release cadence** — an internal API consumed only by services your own
   team deploys weekly needs a much shorter window than a public API
   consumed by third parties who batch their own dependency upgrades
   quarterly or yearly.
5. **Track migration progress actively**, don't just wait for the sunset
   date — if usage of the old path hasn't dropped as the date approaches,
   that's a signal to reach out to remaining callers directly (for
   internal APIs) or extend the deadline with clear justification (for
   public ones), not to remove it on schedule regardless.
6. **Remove on schedule once usage is at or near zero**, with a final
   direct check immediately before removal — don't let "we announced it
   months ago" substitute for confirming actual current usage is gone.

## Choosing a grace period

There's no universal number — set it from these inputs, in order of
weight: (1) how many external/uncontrolled callers exist and how often
they release, (2) whether the deprecated thing has any known usage at all
(zero known usage can justify near-immediate removal, with a fallback
plan), (3) contractual/SLA commitments already made to specific
customers, (4) the cost of maintaining the deprecated path in parallel
with its replacement. Write down which of these drove the specific number
chosen — a deprecation timeline without a stated rationale looks arbitrary
to affected callers and invites requests to change it.

## Gotchas

- A "deprecated" label with no sunset date is not a deprecation — it's
  permanent dead weight that nobody has authorization to actually remove,
  because no date was ever set as the forcing function.
- Removing a deprecated field/endpoint exactly on its announced date
  without checking current usage first has caused real outages when usage
  didn't actually reach zero — the date is a target for the removal
  *attempt*, not a substitute for the final usage check.
- Internal deprecations inside a single company are not exempt from
  Hyrum's Law — another team's code depending on an internal API's
  undocumented quirk (error message text, field ordering, timing) is just
  as real a dependency as an external customer's, and just as capable of
  breaking silently.
- A deprecation notice that only appears in a changelog or release notes
  reaches engineers who read changelogs — it does not reach the engineer
  three years from now who copy-pasted a code sample from an old internal
  wiki page. In-artifact warnings (compiler/linter/runtime) reach both.
- Version-pinning schemes (e.g., pinning callers to the API version active
  when their integration was created) let you deprecate old *behavior* for
  new integrations immediately while giving existing integrations a
  runway measured in "until they choose to move," not a fixed calendar
  date — a different, sometimes better-fitting model than a universal
  sunset date.

## Real-world grounding

Hyrum's Law is documented in *Software Engineering at Google* (Winters,
Manshreck, Tannenbaum, O'Reilly, 2020) as an observation the authors
attribute to Google engineer Hyrum Wright, and is cited there specifically
in the context of why API deprecation at Google's scale requires
usage-measurement and gradual migration processes rather than a
documentation-only notice — the same organizational discipline this skill
describes.

## Verification

- [ ] Actual current usage of the deprecated thing was measured before
      announcing a sunset date
- [ ] The deprecation notice states what's replacing it and gives a
      specific date, not an open-ended "eventually"
- [ ] The deprecation is marked in the artifact itself (header, doc
      comment, runtime warning), not only in external documentation
- [ ] Migration progress is tracked before the sunset date, not assumed
- [ ] Usage is re-checked immediately before actual removal, regardless of
      how long ago it was announced
