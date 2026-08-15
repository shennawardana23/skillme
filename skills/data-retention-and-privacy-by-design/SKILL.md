---
name: data-retention-and-privacy-by-design
description: Use when the user asks to "design a new data model that stores personal data", "add a retention policy", "figure out how long to keep guest data", "handle a data deletion request", "design for GDPR compliance", or is adding any field/table that captures a name, email, phone, ID document, or payment detail. Guides applying GDPR's data minimization principle and Ann Cavoukian's Privacy by Design framework at design time, not as a bolt-on compliance pass after launch.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Data Retention and Privacy by Design

Privacy decisions made at schema-design time are cheap. The same decisions
made after a data breach, a regulator inquiry, or a customer's deletion
request are expensive and sometimes impossible (data already replicated to
backups, analytics warehouses, and third-party processors). Design for
minimal collection and bounded retention from the first migration.

## Privacy by Design: 7 foundational principles (Cavoukian)

Ann Cavoukian's framework, adopted explicitly into GDPR (Article 25, "Data
protection by design and by default"), gives seven principles. The two
with the most day-to-day engineering weight:

1. **Proactive, not reactive.** Privacy risks are addressed before code
   ships, not patched in after an incident.
2. **Privacy as the default setting.** A user who does nothing gets the
   most private configuration automatically — opt-in for extra data use,
   never opt-out.
3. **Privacy embedded into design.** Not a separate compliance layer
   bolted onto a finished system, but a property of the architecture
   itself (schema shape, retention jobs, access controls).
4. **Full functionality — positive-sum.** Privacy protections shouldn't be
   framed as a tradeoff against product functionality; look for designs
   that deliver both (e.g., store a hashed lookup key instead of the raw
   value when only equality-matching is needed).
5. **End-to-end security.** Data is protected across its whole lifecycle —
   collection, transit, storage, and secure deletion — not just at rest.
6. **Visibility and transparency.** Data subjects and auditors can verify
   what's actually collected and done with data, not just read a privacy
   policy that describes intended behavior.
7. **Respect for user privacy.** Design choices center the individual's
   interests, not just the organization's convenience.

## GDPR's data minimization principle

GDPR Article 5(1)(c) requires personal data to be "adequate, relevant, and
limited to what is necessary" for the stated purpose. In practice this
means, for every field you're about to add to a schema:

- **Ask "do we need this specific field for a specific, already-identified
  purpose?"** — not "might this be useful for some future analytics
  question." Speculative collection ("let's capture it in case we need it
  later") is the single most common minimization violation.
- **Prefer derived/aggregated data over raw personal data** when the
  consuming feature only needs the aggregate (e.g., store "guest stayed 3
  nights" rather than a full itinerary, if only a stay-count feature
  consumes it).
- **Set a retention period at the point you add the field**, not later.
  "We'll decide retention when someone asks" means it never gets decided
  and the field accumulates indefinitely.

## Procedure: designing a feature that touches personal data

1. **Name the purpose** for each personal-data field in one sentence
   before writing the migration. If you can't name a specific purpose,
   don't collect the field yet.
2. **Set a retention period** and write it down next to the schema
   definition (a comment, a data dictionary entry, or an ADR) — e.g.,
   "guest email retained 3 years post-checkout for loyalty program, then
   purged." Every partitioned table (see `postgres-hotel-partitioning`
   skill) should have its retention decision documented at the same time
   as its partition key.
3. **Design the deletion path before the collection path.** If honoring a
   deletion request means a manual, multi-team, days-long process, that's
   a design defect discovered too late. A per-subject deletion job or
   cascading foreign-key delete should exist before the feature ships.
4. **Check downstream copies.** Personal data collected once often ends up
   in backups, analytics pipelines, logs, and third-party processors
   (payment providers, email senders). Deletion has to reach all of them,
   or note explicitly which are exempted (e.g., legally required financial
   records) and why.
5. **Default to the least-identifying representation** that still serves
   the purpose: a hashed/tokenized reference beats a raw value; an
   aggregate beats a raw event; a pseudonymous ID beats a real name, when
   the consuming code only needs to distinguish subjects, not identify them.

## Gotchas

- Logging a personal-data field "just for debugging" puts it in a system
  (log aggregation) that usually has a much longer and much less
  controlled retention period than the primary database — treat log
  statements as a data-collection decision, not a free action.
- A `deleted_at` soft-delete column does not satisfy a GDPR deletion
  request by itself — the row is still present and often still queryable
  or exportable until a hard-delete/purge job actually runs.
- Free-text fields (support ticket notes, special-request fields) are a
  common accidental home for personal data (a guest pastes their passport
  number into a "special requests" box) that no schema review will catch —
  minimization has to be paired with output/anonymization review for
  free-text fields feeding into analytics.
- "We might need it later" is not a purpose. If a future feature needs the
  field, add it when that feature is actually being built — collecting
  ahead of need is the textbook minimization violation regulators cite.
- Cross-border data transfer rules (e.g., GDPR's restrictions on
  transferring EU personal data outside the EEA without a valid transfer
  mechanism) interact with retention: a backup replicated to a
  non-compliant region is itself a separate compliance problem, independent
  of how long it's retained.

## Real-world grounding

Ann Cavoukian's "Privacy by Design: The 7 Foundational Principles" (2009,
Information and Privacy Commissioner of Ontario) is the publicly documented
source for the framework above, and is explicitly cited as the model for
GDPR Article 25's "data protection by design and by default" requirement —
this is a well-established, named lineage from a specific framework into a
specific binding regulation, not a general best-practice restatement.

## Verification

- [ ] Every new personal-data field has a one-sentence stated purpose
- [ ] A retention period is documented alongside the schema, not deferred
- [ ] A deletion path exists (automated purge or cascading delete), not
      just a soft-delete flag
- [ ] Downstream copies (logs, backups, analytics, third parties) are
      accounted for in the deletion design
- [ ] The least-identifying representation that serves the purpose was
      chosen over the raw/most-identifying one
