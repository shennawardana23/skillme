---
name: deployment-patterns
description: Guides production rollout strategy — rolling/blue-green/canary deployments, liveness/readiness/startup probes, startup config validation, and rollback planning. Use when choosing a deployment strategy, implementing health check endpoints or Kubernetes probes, planning a rollback path before a release, or reviewing a production-readiness checklist before a launch. For building the container image itself, see docker-patterns; for the pipeline that runs the deploy, see ci-cd-and-automation.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Deployment Patterns

This skill is about what happens once a built, tested artifact is ready
to go live: how traffic moves to it, how the platform knows it's healthy,
and how to get back to the last good state if it isn't. Building the
image is `docker-patterns`; the pipeline that builds, tests, and
triggers the deploy is `ci-cd-and-automation`.

## Choosing a rollout strategy

| Strategy | Mechanism | Rollback speed | Cost | Use when |
|---|---|---|---|---|
| **Rolling** (default) | Replace instances gradually; old and new run simultaneously | Redeploy previous image (not instant) | No extra infra | Standard deploys with backward-compatible changes |
| **Blue-green** | Two full environments; switch traffic atomically | Instant (switch back) | 2x infra during deploy | Critical services, near-zero tolerance for bad deploys |
| **Canary** | Small traffic % to new version first, then ramp | Fast (cut canary traffic) | Requires traffic-splitting infra | High-traffic services, risky changes |

Rolling deployment's hard requirement: **the old and new versions must
be able to run simultaneously against the same database and message
schema** — if a migration or message format isn't backward-compatible
with the version being replaced, a rolling deploy will have some
requests hit old code and some hit new code mid-rollout, and one of them
breaks. See `database-migrations`' expand/migrate/contract pattern for
how to make a schema change safe under this constraint.

## Health checks: three distinct signals, not one

```go
// Liveness: is the process alive at all? Should almost never fail.
// A failing liveness check triggers a restart — false positives are expensive.
func (h *Handler) Livez(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

// Readiness: can this instance serve traffic right now?
// Checks dependencies. A failing readiness check pulls the instance
// out of the load balancer WITHOUT restarting it.
func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
    if err := h.db.PingContext(r.Context()); err != nil {
        http.Error(w, "database unreachable", http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

```yaml
livenessProbe:
  httpGet: { path: /livez, port: 8080 }
  initialDelaySeconds: 10
  periodSeconds: 30
  failureThreshold: 3

readinessProbe:
  httpGet: { path: /readyz, port: 8080 }
  initialDelaySeconds: 5
  periodSeconds: 10
  failureThreshold: 2

startupProbe:                    # covers slow-starting apps (large cache warmup, migrations)
  httpGet: { path: /livez, port: 8080 }
  periodSeconds: 5
  failureThreshold: 30           # 30 * 5s = 150s max allowed startup time
```

Conflating liveness and readiness into one endpoint is the most common
mistake: a liveness check that also pings the database means a
transient database blip **restarts every instance simultaneously**
instead of just pulling them from rotation — turning a brief dependency
hiccup into a full outage during the restart storm.

## Startup config validation — fail fast, not fail-weird

Validate every required environment variable at process start, before
serving a single request:

```go
type Config struct {
    Port        int    `env:"PORT" envDefault:"8080"`
    DatabaseURL string `env:"DATABASE_URL,required"`
    JWTSecret   string `env:"JWT_SECRET,required"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    if len(cfg.JWTSecret) < 32 {
        return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
    }
    return cfg, nil
}

func main() {
    cfg, err := LoadConfig()
    if err != nil {
        log.Fatalf("config error: %v", err) // crash immediately, don't limp along
    }
    // ...
}
```

A service that starts successfully with a missing config value and
fails on the *first request that needs it* turns a deploy-time problem
into a production incident hours later — validate everything the
process needs at startup, exit non-zero immediately if anything is
missing or malformed.

## Rollback: plan it before you need it

```bash
kubectl rollout undo deployment/app          # Kubernetes: revert to previous ReplicaSet
kubectl rollout undo deployment/app --to-revision=3
```

Rollback checklist, verified **before** the release it protects, not
during the incident:

- [ ] Previous image/artifact is still available and tagged (not
      overwritten by the new build)
- [ ] Database migrations for this release are backward-compatible with
      the previous code version (expand/migrate/contract, never a bare
      destructive change in the same release as new code that needs it)
- [ ] Risky new behavior is behind a feature flag that can be disabled
      without a redeploy
- [ ] Rollback has actually been exercised in staging, not just
      documented

## Production readiness checklist

**Application**: structured logging with no PII (`observability-and-instrumentation`),
error handling covers edge cases, no hardcoded secrets.

**Infrastructure**: image is reproducible (pinned base image versions —
`docker-patterns`), resource limits set (CPU/memory), horizontal scaling
min/max configured, TLS on every endpoint.

**Monitoring**: RED metrics exported, alerts configured on symptoms with
runbook links (`observability-and-instrumentation`), uptime monitoring on
the readiness endpoint specifically (not liveness — liveness can be "up"
while the instance can't actually serve a real request).

**Operations**: rollback plan documented *and tested*, migration tested
against production-sized data (`database-migrations`), on-call rotation
defined for this service before it goes live.

## Gotchas

- A `readinessProbe` with too aggressive a `failureThreshold`/`periodSeconds`
  pulls an instance out of rotation on a single slow response — tune
  against the endpoint's actual p99, not a guess.
- `startupProbe` absence on a slow-starting service means the
  `livenessProbe`'s `initialDelaySeconds` has to cover worst-case startup
  time for every restart, forever — a dedicated `startupProbe` with a
  longer allowance decouples "how long can startup take" from "how
  quickly do we detect a genuinely hung process" once running.
- Rolling deployment with a schema change that isn't backward-compatible
  is a guaranteed partial-outage window — the old code path will run
  against the new schema (or vice versa) for however long the rollout
  takes, not an instant.
- Blue-green's "instant rollback" assumption breaks if the new version
  already wrote data in a format the old version can't read — the
  atomic-switch guarantee only covers the application layer, not shared
  mutable state.

## Real-world grounding

Knight Capital's 2012 production incident is the canonical case for
rollout-strategy risk: a new deployment was rolled out to 7 of 8
production servers, leaving one server running old code that
reinterpreted a repurposed configuration flag as a signal to activate a
long-dormant test feature — the mismatch between what was deployed where
went undetected, and the resulting erroneous trading activity cost the
firm roughly $440 million in the 45 minutes before it was manually
stopped. The lesson generalizes directly: a rollout strategy is only as
safe as the guarantee that every instance is verifiably running one
consistent, known version — "it's probably fully rolled out by now" is
not a verification step.

## Verification

- [ ] Rollout strategy matches the risk level of the change (rolling default; blue-green/canary for higher-risk releases)
- [ ] Liveness and readiness are separate endpoints with separate failure semantics
- [ ] All required config is validated at startup; the process exits non-zero on missing/invalid config rather than starting degraded
- [ ] A rollback path exists, is documented, and has been exercised in staging
- [ ] Any schema change in this release is backward-compatible with the version being replaced during rollout
- [ ] Alerts and uptime monitoring target the readiness endpoint, not just liveness
