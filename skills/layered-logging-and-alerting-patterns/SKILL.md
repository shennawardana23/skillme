---
name: layered-logging-and-alerting-patterns
description: Guides services where a single logging call fans out to multiple side effects - error tracking (e.g. Sentry), a centralized log-shipping endpoint, and chat/webhook alerting - rather than just writing a log line. Use when adding a log statement to a service with this shape, reviewing why removing a "log statement" broke alerting, or designing observability for a service with more than one logging sink.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Layered Logging and Alerting

In some services, a single `logger.Error(...)` call is not "just a log
line" — it can simultaneously capture an event to an error tracker, ship
the log to a centralized collection endpoint, and (for certain severities
or subsystems) fire a chat/webhook alert. Treating a call site as if it
only writes local output is the source of most surprising regressions in
this shape of system — the fix here is understanding what a given log call
actually fans out to before changing or removing it.

## A log call can be a multi-sink side effect

```go
func Errorf(ctx context.Context, format string, args ...any) {
    msg := fmt.Sprintf(format, args...)
    logToLocalOutput(msg)
    captureToErrorTracker(ctx, msg)        // e.g. Sentry
    go sendToCentralizedLogEndpoint(msg)   // async, fire-and-forget by design
}
```

Before removing or "simplifying" a logging call in a codebase with this
shape, check what it actually does — an apparently-redundant log line can
be the only thing populating an external dashboard or triggering an
on-call alert; deleting it silently removes that signal with no local
indication anything changed.

## Async side effects need a flush guarantee on exit paths

If a log call ships data asynchronously (a background goroutine posting to
an external endpoint), a process that exits immediately after logging a
fatal error can terminate before that goroutine's request completes,
silently losing the log. A deliberate short delay before actually
terminating (after a `Fatal`/`Panic`-style log call specifically) is a
real, intentional pattern to give an in-flight async log delivery time to
finish — not dead code or an accidental performance issue to "optimize
away."

## Sampling and filtering are usually deliberate, not accidental

Error-tracking integrations commonly sample by environment (log everything
in development, sample a smaller percentage in production to control
volume/cost) and filter certain error classes before sending (a known-
noisy, low-value error message excluded via a `beforeSend`-style hook).
Both are typically deliberate operational decisions — verify the actual
current sampling rate and any filter rules before assuming an error "isn't
being tracked" is itself a bug, versus a documented filtering choice.

## Gotchas

- **A log call in this kind of system can have side effects well beyond
  writing a line** — capturing to an error tracker, shipping to a
  centralized endpoint, firing a chat alert — and "removing an unused log
  statement" can silently break one of those downstream effects with no
  local signal that anything changed.
- **A short sleep before process exit on a fatal/panic path can be a
  deliberate flush guarantee** for an async log-shipping goroutine, not
  dead code — removing it can cause fatal-path logs to be lost right when
  they matter most.
- **Environment-conditional sampling and message-based filtering in an
  error tracker are usually intentional** — check the actual configured
  rule before treating "this error isn't showing up in the dashboard" as
  a bug rather than a documented filter.
- **More than one logging API coexisting in one codebase (a legacy
  logging-library-compatible shim alongside a newer logger, or a tracer
  initialized but only partially wired into request paths) is often
  deliberate migration-in-progress state**, not accidental duplication —
  confirm which is the actual source of truth for a given code path before
  assuming either is dead code to remove.
- **A per-request "start transaction/span" call at the top of every
  handler is usually load-bearing boilerplate for tracing**, not
  copy-paste cruft — removing it from a new handler quietly drops that
  handler out of tracing/APM visibility with no error, just an absence of
  data later.

## Real-world grounding

Fanning a single log call out to multiple sinks (structured local output,
an error tracker, a metrics/alerting pipeline) is a standard pattern once a
service needs both human-searchable logs and automated incident response —
the risk it introduces (a "log cleanup" accidentally removing an alerting
trigger) is exactly why observability tooling vendors document explicit
"what does this call actually do" guidance for any wrapper logging function
rather than assuming a log call is side-effect-free by default.

## Verification

- [ ] Before removing or changing a log call, its actual downstream
      effects (error tracker capture, centralized shipping, alerting)
      have been checked, not assumed to be local-only
- [ ] Any deliberate delay before process exit on a fatal/panic path is
      recognized as a flush guarantee, not removed as an apparent
      inefficiency
- [ ] Sampling rates and filter rules in the error-tracking integration
      are checked before concluding an error "isn't being tracked" is a
      bug
- [ ] Coexisting logging APIs (a legacy shim, a newer logger, an
      under-wired tracer) are understood as migration state before either
      is assumed to be dead code
- [ ] New request handlers include the same tracing/span-start boilerplate
      existing handlers use, so they aren't silently invisible to tracing
