---
name: shipping-and-launch
description: Prepare a production launch so it's reversible, observable, and incremental — pre-launch checklist, feature-flag strategy, staged rollout with rollback thresholds, and post-launch monitoring. Use when preparing to deploy to production, releasing a significant user-facing change, migrating data or infrastructure, or opening a beta program.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Shipping and Launch

The goal isn't just to deploy — it's to deploy so that a bad outcome is cheap to detect and cheap to undo. Every launch should be reversible, observable, and incremental.

## When to Use

- Deploying a feature to production for the first time
- Releasing a significant change to users
- Migrating data or infrastructure
- Opening a beta or early-access program
- Any deployment that carries risk — which is effectively all of them

## Pre-Launch Checklist

**Code quality**: all tests pass; build succeeds with no warnings; lint and type checks pass; code reviewed; no unresolved TODOs blocking launch; no debug logging left in; error handling covers expected failure modes.

**Security**: no secrets in code or version control; dependency audit shows no critical/high vulnerabilities; input validation on user-facing endpoints; auth/authz checks in place; security headers configured; rate limiting on auth endpoints; CORS scoped to specific origins, not a wildcard.

**Performance**: Core Web Vitals within "good" thresholds if user-facing web; no N+1 queries in critical paths; images optimized; bundle size within budget; database queries have appropriate indexes; caching configured for static assets and repeated queries.

**Accessibility**: keyboard navigation works for interactive elements; a screen reader can convey structure and content; color contrast meets WCAG 2.1 AA; focus management is correct for modals/dynamic content; error messages are descriptive and associated with the relevant field.

**Infrastructure**: environment variables set in production; database migrations applied or ready; DNS/SSL configured; CDN configured for static assets; logging and error reporting configured; a health-check endpoint exists and responds.

**Documentation**: README reflects any new setup requirement; API docs current; ADRs written for architectural decisions made (see `documentation-and-adrs`); changelog updated.

## Feature Flag Strategy

Ship behind a flag to decouple deployment from release — the code reaches production inert, and gets turned on separately.

Lifecycle: deploy with the flag off → enable for team/internal beta → gradual rollout (5% → 25% → 50% → 100%) → monitor at every stage → clean up the flag and dead code path once fully rolled out.

Rules: every flag has an owner and an expiration date; clean up within roughly 2 weeks of full rollout; don't nest flags (creates exponential combinations to reason about); test both flag states in CI, not just the "on" path.

## Staged Rollout

```
1. Deploy to staging — full test suite, manual smoke test of critical flows
2. Deploy to production, flag OFF — verify health check, check error monitoring
3. Enable for team (flag ON internally) — 24-hour monitoring window
4. Canary (flag ON for ~5% of users) — 24-48 hour window, compare canary vs. baseline
5. Gradual increase (25% → 50% → 100%) — same monitoring at each step, ability to roll back to the previous percentage at any point
6. Full rollout — monitor for 1 week, then clean up the flag
```

### Rollout Decision Thresholds

| Metric | Advance | Hold and investigate | Roll back |
|---|---|---|---|
| Error rate | Within 10% of baseline | 10-100% above baseline | >2x baseline |
| P95 latency | Within 20% of baseline | 20-50% above baseline | >50% above baseline |
| Client JS errors | No new error types | New errors, <0.1% of sessions | New errors, >0.1% of sessions |
| Business metrics | Neutral or positive | Decline <5% (may be noise) | Decline >5% |

Roll back immediately, without waiting for the next scheduled check, if: error rate more than doubles baseline, P95 latency rises more than 50%, user-reported issues spike, data integrity issues appear, or a security vulnerability is discovered.

## Monitoring and Observability

Track application metrics (error rate overall and per-endpoint, response time percentiles, request volume, active users, key business metrics), infrastructure metrics (CPU/memory, DB connection pool usage, disk space, network latency, queue depth), and client metrics (Core Web Vitals, JS errors, client-perceived API error rate, page load time). Set up error reporting (an error boundary on the client, structured error middleware on the server) before launch, not after the first incident — and never expose internal error details to the end user in the response.

### Post-Launch Verification (first hour)

Health check returns 200; error-monitoring dashboard shows no new error types; latency dashboard shows no regression; the critical user flow was tested manually; logs are flowing and readable; the rollback mechanism was verified ready (dry run if possible).

## Rollback Strategy

Every deployment needs a written rollback plan *before* it happens, not improvised during an incident:

```markdown
## Rollback Plan for [Feature/Release]

### Trigger Conditions
- Error rate > 2x baseline
- P95 latency > [X]ms
- User reports of [specific issue]

### Rollback Steps
1. Disable the feature flag, OR deploy the previous version
2. Verify: health check, error monitoring
3. Communicate the rollback to the team

### Database Considerations
- Does the migration have a rollback path?
- What happens to data written by the new feature — preserved or cleaned up?

### Time to Rollback
- Feature flag: <1 minute
- Redeploy previous version: <5 minutes
- Database rollback: <15 minutes
```

## Gotchas

- "It works in staging" is not evidence it will work in production — production has different data shape, traffic patterns, and edge cases that staging rarely reproduces; the first hour of production monitoring is not optional just because staging passed.
- A feature flag with no owner or expiration date tends to become permanent — the dead code path and the flag check both linger indefinitely, and eventually nobody remembers which behavior is real.
- Rolling back is not an admission of failure — shipping a broken feature and leaving it broken while debugging in production is the actual failure; the rollback plan existing and being used is the system working as designed.

## Real-world grounding

Knight Capital's 2012 trading-software deployment incident — a widely documented case where a botched production deployment (an old feature flag/code path unexpectedly reactivated during a partial rollout across a fleet of servers) caused a major, rapid financial loss within about 45 minutes of market open — is the canonical real-world cautionary tale for exactly the failure modes this skill guards against: no working rollback plan, insufficient monitoring to catch the problem fast, and a rollout that wasn't staged or reversible enough to contain the blast radius before real damage was done.

## Verification

Before deploying:
- [ ] Pre-launch checklist completed across all sections
- [ ] Feature flag configured, with an owner and expiration date, if applicable
- [ ] A written rollback plan exists with trigger conditions and step-by-step actions
- [ ] Monitoring dashboards are set up and someone is watching them
- [ ] The team is notified of the deployment window

After deploying:
- [ ] Health check returns 200
- [ ] Error rate and latency are within the "advance" threshold, not just "not obviously broken"
- [ ] The critical user flow was verified manually
- [ ] Logs are flowing
- [ ] The rollback path was verified ready, not just assumed to work
