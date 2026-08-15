---
name: bug-triage-and-severity-classification
description: This skill should be used when the user asks to "triage this bug", "what severity is this", "classify this incident/bug", "is this a Sev1 or Sev2", "what priority should this ticket be", or is deciding SLA/response urgency for a reported defect. Use for assigning severity and priority to bugs and incidents, distinguishing the two, and running a triage process — not for the technical debugging itself.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Bug Triage and Severity Classification

Severity and priority are two different axes, and conflating them is the
single most common triage mistake. Getting this right determines whether
the right bug gets fixed at 2am versus next sprint.

## Severity vs. Priority — they are orthogonal

- **Severity** = objective technical/business impact if the bug is not
  fixed. Does not change based on who reported it or how busy the team is.
- **Priority** = how soon it gets worked, given everything else in flight.
  Business context can shift this independently of severity.

A bug can be high severity, low priority (data corruption in a feature used
by one deprecated internal tool nobody touches this quarter) or low
severity, high priority (a cosmetic logo glitch on the page the CEO is
demoing to investors tomorrow). Both axes are legitimate; neither should be
inflated to force the other's outcome — inflating severity to jump the
queue burns the taxonomy's credibility for the next real Sev1.

## The Sev1-4 / P0-P3 taxonomy

This is the industry-common shape (used in variants by Google, PagerDuty,
Atlassian, and most on-call rotations):

| Severity | Definition | Example | Typical response SLA |
|---|---|---|---|
| Sev1 / P0 | Full outage, data loss, security breach, or revenue-blocking failure in production | Booking engine can't process any reservations; payment data exposed | Page on-call immediately, respond in minutes |
| Sev2 | Major functionality broken or badly degraded for a significant user subset, no workaround | Rate calculation wrong for one hotel brand; search returns errors for 20% of queries | Same business day, often within hours |
| Sev3 | Minor functionality broken, workaround exists, limited blast radius | A filter on a report page doesn't persist; one edge case in date formatting | Next sprint / few business days |
| Sev4 / P3 | Cosmetic, typo, or negligible impact | Misaligned button, wrong icon color | Backlog, best-effort |

Exact SLA numbers vary by org — the point of the taxonomy is that everyone
agrees on what the four buckets *mean*, so the SLA table can be applied
mechanically instead of debated per-ticket.

## Triage procedure

1. **Intake**: capture repro steps, environment, affected hotel(s)/tenant,
   and a screenshot/log if available. A bug without repro steps isn't
   ready for severity classification — get that first.
2. **Reproduce or corroborate**: confirm the reported behavior actually
   occurs (or find supporting evidence — logs, error rate graphs — if it
   isn't reproducible on demand).
3. **Classify severity** using the impact definition, not the reporter's
   framing. A VP calling something "critical" doesn't make it Sev1 if it
   affects zero production users.
4. **Assign priority** given current sprint/on-call load — this is where
   business judgment enters, severity stays fixed.
5. **Route** to an owning team/individual with the severity and priority
   attached, and the SLA clock starts.
6. **Re-triage on new information**: if reproduction reveals broader
   blast radius than first reported, re-score severity immediately rather
   than leaving the original (possibly wrong) label in place.

## Gotchas

- **Severity should not be adjusted to manage SLA pressure.** Downgrading
  a Sev2 to a Sev3 near a deadline to stop the SLA clock is gaming the
  taxonomy, not triage — it hides real risk from whoever reads the
  severity distribution later.
- **"Affects everyone a little" and "affects one customer a lot" are not
  automatically the same severity** — check the impact definition for
  which one your org's Sev2 actually means; teams often disagree here
  until it's written down explicitly.
- **A bug's severity can be worse than the immediate symptom.** A silently
  wrong number (e.g., a rate calculation off by a small amount) is often
  higher severity than a loud, obvious error — nobody works around silent
  data corruption because nobody sees it happening.
- **Don't skip reproduction.** Classifying severity from a vague report
  ("something's broken with search") produces an SLA commitment against a
  problem that isn't yet understood.

## Real-world grounding

The Sev1-4 naming convention and its accompanying response-time SLAs are
publicly documented in PagerDuty's incident response guidance (their
public incident-response documentation defines SEV1-SEV5 with concrete
response expectations), and the severity/priority distinction is standard
practice reflected across on-call and issue-tracking conventions
(Atlassian's Jira priority/severity fields, Google's internal incident
management practice referenced in the SRE Workbook). The specific numbers
differ between organizations; the two-axis structure does not.
