---
name: incident-response-and-postmortems
description: This skill should be used when the user asks to "write a postmortem", "run an incident retro", "do a blameless postmortem", "document this outage", or is drafting an incident timeline, root cause analysis, or action items after a production incident. Use for the incident-response process and the after-the-fact postmortem document — not for the live technical debugging itself (see the diagnose/debugging skills for that).
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Incident Response and Postmortems

A postmortem's job is to make the system (technical and organizational)
less likely to fail the same way twice — not to identify who to blame.
Blameless postmortem practice, as documented by Google SRE and by Etsy,
exists because blame-seeking postmortems make people hide information,
which makes the next incident worse, not better.

## Blameless language, concretely

Blameless doesn't mean no accountability — actions still get named owners
and deadlines. It means describing what happened in terms of the system
and process, not the person:

- Wrong: "Jane forgot to run the migration before deploying."
- Right: "The deploy pipeline did not enforce that the migration step ran
  before the application deploy step, so the deploy proceeded against a
  schema it wasn't compatible with."

The second version points at a fixable gap (add a pipeline gate); the
first version's only "fix" is "Jane, be more careful," which doesn't
survive Jane going on vacation or leaving the team.

## When to write one

Trigger a postmortem for: any user-facing outage above your SLO error
budget threshold, any data loss or data corruption incident, any security
incident, any incident requiring manual/emergency intervention (rollback,
hotfix, manual DB repair), or any incident where the on-call was paged
outside business hours. Err toward writing one — a short postmortem for a
near-miss is cheap; a missed pattern from skipping one is not.

## Postmortem structure

1. **Summary** — one paragraph: what broke, for how long, who/what was
   affected.
2. **Impact** — concrete numbers: duration, % of requests/users affected,
   which hotels/tenants (relevant in a multi-tenant, `hotel_id`-partitioned
   system — impact often needs to be scoped per-hotel, not just
   globally), revenue or SLA impact if known.
3. **Timeline** — UTC timestamps, one line per event, from first anomaly
   to full resolution: detection, escalation, diagnosis milestones,
   mitigation, resolution. Build this from logs/alerts/chat history, not
   memory — memory reorders events under stress.
4. **Root cause** — use "5 whys" to go past the first proximate cause.
   Stopping at "the deploy failed" is rarely the real cause; keep asking
   why until you hit a systemic gap (missing test, missing alert, missing
   pipeline gate, unclear ownership).
5. **Detection** — how was it noticed: alert, customer report, internal
   discovery? A postmortem where detection was "a customer told us" is
   itself a finding — that's a monitoring gap.
6. **Action items** — each one has a named owner and a due date, and is
   specific enough to be verifiably done or not done. "Communicate better"
   is not an action item; "add a pipeline gate that blocks deploy until
   the migration step reports success" is.
7. **Lessons learned** — what went well (fast detection, good runbook)
   alongside what didn't, so the good parts get reinforced too.

## Process

- Draft the timeline within 24-48 hours while memory and logs are fresh;
  don't wait for a "complete" root cause before starting the draft.
- Hold a review meeting with everyone involved plus at least one person
  who wasn't — an outside perspective catches assumptions insiders don't
  notice.
- Publish it somewhere durable and searchable, not just in the incident
  chat channel that will scroll away.
- Track action items to closure the same way any other committed work is
  tracked — an action-item backlog that never gets revisited is a
  postmortem practice in name only.

## Gotchas

- **A postmortem with zero concrete action items, or all "communicate
  better"/"be more careful" items, has not actually found the systemic
  cause** — go back to the 5 whys.
- **Don't wait for the "perfect" root cause to publish the timeline.**
  The timeline is valuable on its own even before root cause is fully
  understood; the two don't have to ship together.
- **Blameless doesn't mean anonymous or vague.** Naming what a specific
  system/pipeline/process did (not who) is precise and useful; refusing to
  name any specifics "to be safe" produces a postmortem nobody can act on.
- **Re-triage incident severity after the fact if the postmortem reveals
  broader impact than first reported** — see
  `bug-triage-and-severity-classification` for how severity classification
  should update on new information.
- **A recurring root cause across multiple postmortems is itself a
  finding** — track root-cause categories over time; if "missing input
  validation" keeps appearing, that's a signal for a systemic fix, not
  three unrelated incidents.

## Real-world grounding

Google's SRE book documents blameless postmortem culture as core SRE
practice (freely available: sre.google/sre-book/postmortem-culture/),
including the principle that postmortems focus on process and tooling
gaps rather than individual fault. John Allspaw's 2012 Etsy engineering
blog post, "Blameless PostMortems and a Just Culture," is the widely-cited
origin of applying aviation/healthcare "just culture" thinking to software
incident review, arguing that blame-seeking suppresses the information
needed to actually prevent recurrence.
