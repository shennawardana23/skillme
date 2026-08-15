---
name: observability-and-instrumentation
description: Guides instrumenting Go (and TypeScript/PHP) services with structured logging, RED/USE metrics, and distributed tracing so production behavior is diagnosable from telemetry alone. Use when adding logging, metrics, or tracing to a new feature or service, reviewing a PR that adds retries/queues/external calls, setting up or reviewing alerting rules, or when a production issue took too long to diagnose because the data wasn't there.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Observability and Instrumentation

Code you can't observe is code you can't operate. Instrumentation is
written alongside the feature, not bolted on after the first incident —
by the time you need it in a live outage, it's too late to add.

## 1. Define "working" before instrumenting

Telemetry without a question behind it is noise. Write down 2–4 questions
an on-call engineer will ask about this feature before adding any signal:

```
FEATURE: reservation cancellation
QUESTIONS ON-CALL WILL ASK:
1. What fraction of cancellations succeed vs. fail?
2. When one fails, why? (payment refund error? already checked out? DB error?)
3. Is the downstream refund provider slower than usual right now?
→ every signal below answers one of these — nothing else gets added.
```

## 2. Pick the right signal

| Signal | Answers | Example |
|---|---|---|
| Structured log | "What happened in this one case?" | `event=reservation_cancel_failed reservation_id=... reason=already_checked_out` |
| Metric | "How often / how fast, in aggregate?" | p99 latency of the cancel endpoint |
| Trace | "Where did time go across services?" | One slow cancel, broken down by hop (DB → payment provider) |

Rule of thumb: metrics tell you **that** something is wrong, traces tell
you **where**, logs tell you **why**.

## 3. Structured logging (Go: `log/slog`)

Log events, not prose — every line is a structured record with a stable
event name:

```go
// BAD: string interpolation — unqueryable, inconsistent shape
log.Printf("reservation %s cancel failed for hotel %d: %v", id, hotelID, err)

// GOOD: structured, stable event name, machine-readable fields
logger.Error("reservation cancel failed",
    slog.String("event", "reservation_cancel_failed"),
    slog.Int64("hotel_id", hotelID),
    slog.String("reservation_id", id),
    slog.String("error", err.Error()),
    slog.String("request_id", requestID),
)
```

**Log levels, used consistently:**

| Level | Meaning | On-call action |
|---|---|---|
| `Error` | invariant broken, may need action | investigate |
| `Warn` | degraded but handled (retry succeeded, fallback used) | watch for trend |
| `Info` | significant business event (reservation confirmed) | none |
| `Debug` | diagnostic detail | off in production by default |

**Correlation IDs are mandatory.** Generate or accept a request ID at the
system boundary (middleware) and thread it through every log line, span,
and outbound call — otherwise a single request can't be reconstructed
from interleaved concurrent logs:

```go
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.Header.Get("X-Request-ID")
        if id == "" {
            id = uuid.NewString()
        }
        ctx := context.WithValue(r.Context(), requestIDKey{}, id)
        w.Header().Set("X-Request-ID", id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Never log secrets, tokens, or full PII** — allowlist fields explicitly;
never log a whole request body or a guest record wholesale. This is a
hard rule shared with `security-review`.

## 4. Metrics: RED and USE, never bare averages

For request-driven endpoints, instrument **RED**: Rate, Errors, Duration
(as a histogram). For resources (pools, queues, workers), use **USE**:
Utilization, Saturation, Errors.

```go
var cancelDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
    Name:    "reservation_cancel_duration_seconds",
    Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5},
}, []string{"status_class"}) // "2xx", "4xx", "5xx" -- never a raw status code or user ID
```

**Cardinality is the failure mode.** Every unique label combination is a
separate time series — labels must come from small, fixed sets:

```
OK as a label:     route="/reservations/:id"   status_class="5xx"   provider="stripe"
NEVER a label:     user_id, guest_id, hotel_id (if unbounded), raw error message, full URL
```

Track percentiles, never bare averages — an average hides the 1% of
requests having a terrible time. Read p50/p95/p99 from the histogram.

## 5. Distributed tracing

Use OpenTelemetry (`go.opentelemetry.io/otel`) — vendor-neutral, with
auto-instrumentation for `net/http`, gRPC, and common SQL drivers:

```go
tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
otel.SetTracerProvider(tp)

func (s *Service) Cancel(ctx context.Context, id string) error {
    ctx, span := tracer.Start(ctx, "ReservationService.Cancel")
    defer span.End()
    span.SetAttributes(attribute.String("reservation_id", id))
    // ...
}
```

Add manual spans around meaningful internal units of work (a refund
call, a DB transaction), and propagate context across every async
boundary — HTTP headers, queue message metadata — or the trace dies at
the gap. Sample low by default (head-based), and keep all error traces
if the backend supports tail sampling.

## 6. Alerting: symptoms, not causes

```
SYMPTOM (page-worthy):          CAUSE (dashboard, not a page):
error rate > 1% for 5 min       CPU at 85%
p99 latency > 2s                one pod restarted
queue age > 10 min              disk at 70%
```

Cause-based alerts fire when nothing is actually wrong for users, and
miss failure modes nobody predicted. Every alert needs: (1) to be
actionable — if the response is "ignore it, self-heals," delete it; (2) a
runbook link, even three lines; (3) a threshold justified by an SLO or
historical data, not a guess; (4) exactly two severities — page
(user-facing, act now) or ticket (degradation, act this week). A third
tier becomes noise everyone learns to ignore.

## 7. Verify the telemetry itself

Instrumentation is code — it can be wrong. Before calling it done:
force an error in staging and find it in logs by request ID; send test
traffic and confirm metric series appear with expected labels; follow
one request across services in the tracing UI with no broken spans; fire
each new alert once (temporarily lower the threshold) and confirm it
reaches the right channel.

## Gotchas

- A feature PR that adds retries, a queue, or a new external call with
  zero new telemetry is a red flag on its own — that's exactly the code
  path that fails in ways nothing else surfaces.
- Logging a `user_id` or `guest_id` as a **metric label** (not a log
  field) is a cardinality bomb that can take down the metrics backend
  itself, distinct from and much worse than an unqueryable log line.
- `slog`'s default handler in Go writes text, not JSON — explicitly
  configure `slog.NewJSONHandler` in production, or every downstream log
  aggregator has to re-parse unstructured text.
- An alert that fires daily and gets acknowledged without action is not
  monitoring, it's pager training people to ignore the pager — the fix
  is to retune or delete it, not to add a Slack mute.
- `console.log`/`fmt.Println` debug statements left in a request path
  are unstructured, unfilterable, and often the very thing that leaks a
  secret into logs by accident.

## Real-world grounding

The RED/USE framing and the "alert on symptoms, not causes" principle
trace back to Google's Site Reliability Engineering practice (the
publicly available SRE book, 2016) and its "four golden signals"
(latency, traffic, errors, saturation) — the core argument is identical
to what's above: paging a human for a cause metric (CPU, disk) that
hasn't yet produced user-visible harm trains the on-call rotation to
distrust and eventually ignore the pager, which is worse than having no
alert at all.

## Verification

- [ ] The on-call questions for this feature are written down; every signal maps to one
- [ ] All logs are structured (JSON in production), with a stable event name and a correlation/request ID on every line
- [ ] No secrets, tokens, or unredacted PII appear in any log line (spot-checked against actual output)
- [ ] RED metrics exist for every new endpoint and external dependency, with bounded label sets
- [ ] Latency is a histogram; p95/p99 are queryable, not just an average
- [ ] A single request can be followed end-to-end in the tracing UI without broken spans
- [ ] Every new alert is symptom-based, links a runbook, and was test-fired once before shipping
