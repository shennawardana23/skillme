# Communication Patterns and Escalation Protocols

Load this when actually writing a schedule notification, escalation, or
status update — not needed for day-to-day scheduling decisions.

## Communication Patterns

- **Daily schedule publication**: clear, structured, table format — the shop floor doesn't read paragraphs.
- **Schedule change notification**: urgent header, reason, specific affected jobs, new sequence/timing, effective time.
- **Disruption escalation**: lead with impact magnitude (constraint hours lost, orders at risk), then cause, then response, then the decision needed from management.
- **Overtime request**: quantify the business case explicitly — cost of overtime vs. at-risk revenue, plus union-rule compliance.
- **Customer delivery impact**: never surprise the customer — notify as soon as a delay is likely, with the new date, cause (without blaming internal teams), and recovery plan.
- **Maintenance coordination**: specific window requested, business justification, and the cost of deferring it.

## Escalation Protocols

| Trigger | Action | Timeline |
|---|---|---|
| Constraint work center down >30 min unplanned | Alert production manager + maintenance manager | Immediate |
| Plan adherence <80% for a shift | Root cause analysis with shift supervisor | Within 4 hours |
| Customer order projected to miss ship date | Notify sales and customer service with revised ETA | Within 2 hours of detection |
| Overtime exceeds weekly budget by >20% | Escalate to plant manager with cost-benefit analysis | Within 1 business day |
| Constraint OEE <65% for 3 consecutive shifts | Trigger focused improvement event | Within 1 week |
| Quality yield at constraint <93% | Joint review with quality engineering | Within 24 hours |

Escalation chain: Scheduler → Production Manager/Shift Superintendent (30 min for constraint issues) → Plant Manager (2 hours for customer-impacting issues) → VP Operations (same day for multi-customer impact or safety-related changes).
