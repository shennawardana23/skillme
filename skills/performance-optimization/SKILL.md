---
name: performance-optimization
description: Guides measurement-first performance optimization for Go backend services (pprof profiling, allocation reduction, N+1 avoidance) and frontend Core Web Vitals for the Vue/Nuxt stack. Use when a performance requirement exists, a regression is suspected, profiling reveals a bottleneck, or someone proposes an optimization with no measurement backing it.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Performance Optimization

Measure before optimizing. Performance work without measurement is
guessing, and guessing produces complexity that doesn't fix what's
actually slow. The workflow is always: **measure → identify → fix →
verify → guard against regression** — never skip straight to "fix."

## When NOT to use this skill

Don't optimize before you have evidence of a problem. Premature
optimization adds complexity that costs more than the performance it
buys — see the grounding note below. If there's no measurement, the
first step is always to get one, not to guess at a fix.

## Backend: profile before you fix

Go ships production-grade profiling in the standard library — use it
instead of guessing which function is slow:

```go
import _ "net/http/pprof"

// exposes /debug/pprof/ on the given mux — gate behind internal-only
// network access or auth, never expose on a public listener
go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()
```

```bash
# CPU profile: capture 30s of a live process, view where time actually goes
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap profile: find what's allocating/retained
go tool pprof http://localhost:6060/debug/pprof/heap

# Inside the pprof shell
(pprof) top10        # highest cumulative time/allocations
(pprof) list FuncName # annotated source for one function
(pprof) web           # call graph (requires graphviz)
```

For micro-level "is this change actually faster," use Go's benchmark
harness rather than eyeballing wall-clock time:

```go
func BenchmarkSerialize(b *testing.B) {
    for i := 0; i < b.N; i++ {
        serialize(payload)
    }
}
```

```bash
go test -bench=. -benchmem -run=^$ ./...
```

`-benchmem` matters as much as timing: an optimization that halves time
but doubles allocations often loses under real GC pressure — measure
both.

## Where to start measuring, by symptom

```
What is slow?
├── A specific endpoint          --> CPU profile that endpoint under load
├── Memory growing over time     --> heap profile, look for retained large slices/maps
├── Intermittent slowness        --> check goroutine profile for lock contention, GC pause metrics
├── All endpoints, uniformly     --> check connection pool exhaustion, CPU/memory limits, GC pressure
└── A specific DB-backed endpoint --> check query plan (EXPLAIN ANALYZE) before touching app code
```

## Common backend anti-patterns

**N+1 queries** — see `backend-patterns` for the full pattern; it's the
single most common backend performance bug and the first thing to check
before profiling anything else.

**Unbounded fetch:**

```go
// BAD: no limit — grows without bound as the table grows
rows, _ := db.QueryContext(ctx, "SELECT * FROM reservations")

// GOOD: paginated, bounded
rows, _ := db.QueryContext(ctx, "SELECT * FROM reservations WHERE id > $1 ORDER BY id LIMIT 100", lastID)
```

**Allocation churn in a hot path** — reuse buffers instead of allocating
per request:

```go
// BAD: new buffer every call, in a function called per-request
func encode(v any) []byte {
    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(v)
    return buf.Bytes()
}

// GOOD: pooled buffers for a known hot path
var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func encode(v any) []byte {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() { buf.Reset(); bufPool.Put(buf) }()
    json.NewEncoder(buf).Encode(v)
    out := make([]byte, buf.Len())
    copy(out, buf.Bytes())
    return out
}
```

Only reach for `sync.Pool` once a heap profile actually shows this
allocation as significant — pooling every small allocation
speculatively adds complexity and GC-unfriendly retained memory for no
measured benefit.

**Missing caching for expensive, rarely-changing reads** — see the
cache-aside pattern in `backend-patterns`.

## Frontend: Core Web Vitals (Vue/Nuxt)

| Metric | Good | Needs improvement | Poor |
|---|---|---|---|
| LCP (Largest Contentful Paint) | ≤ 2.5s | ≤ 4.0s | > 4.0s |
| INP (Interaction to Next Paint) | ≤ 200ms | ≤ 500ms | > 500ms |
| CLS (Cumulative Layout Shift) | ≤ 0.1 | ≤ 0.25 | > 0.25 |

Measure with both a synthetic tool (Lighthouse, Chrome DevTools
Performance tab — reproducible, good for CI regression gates) and
real-user monitoring (the `web-vitals` library reporting from actual
visitors — required to confirm a fix helped real users, not just the lab
environment). For Nuxt-specific SSR data-fetching and hydration
performance issues, see `vue-nuxt-frontend-patterns` — this skill covers
the measurement and Core Web Vitals targets; that one covers the
Nuxt-specific fetching and reactivity patterns that cause the
regressions.

## Performance budgets, enforced in CI

```
JS bundle (initial load):  < 200KB gzipped
API response time:         < 200ms p95
Lighthouse Performance:    >= 90
```

```bash
npx lhci autorun            # Lighthouse CI regression gate
go test -bench=. -benchmem  # Go benchmark regression gate (compare against a saved baseline)
```

A budget with no enforcement is a suggestion; wire the check into
`ci-cd-and-automation`'s pipeline so a regression fails the build instead
of shipping and getting caught by a user report three weeks later.

## Gotchas

- `go tool pprof` on a CPU profile shows *time spent*, not *call count* —
  a function called once that's slow and a function called a million
  times that's individually fast can show the same "top10" entry; check
  both profile type and call count before concluding which to optimize.
- Benchmarking a function in isolation can hide the real bottleneck if
  the real system is I/O-bound (waiting on the network/DB) rather than
  CPU-bound — a CPU profile of a slow endpoint often shows most time in
  `runtime.gopark` (blocked), which means the fix is elsewhere (the
  query, the downstream call), not in the Go code itself.
- `React.memo`/`useMemo`-equivalent memoization patterns (and Vue's
  `computed`) overused are their own performance cost — added
  indirection and comparison overhead on values that were cheap to
  recompute in the first place. Measure before wrapping everything.
- Exposing `net/http/pprof` on a public-facing listener leaks internal
  memory layout and can be used for a denial-of-service via repeated
  profile requests — always bind it to a localhost-only or
  internal-network-only listener.

## Real-world grounding

"Premature optimization is the root of all evil" is Donald Knuth's most
quoted line from his 1974 paper *Structured Programming with go to
Statements* — the full context is explicitly about spending disciplined
effort on the ~3% of code that measurement shows matters, and treating
everything else as not worth the complexity cost of hand-tuning; it is
not, as often misquoted, an argument against ever thinking about
performance, but the earliest well-known articulation of the
measure-first discipline this skill is built around.

## Verification

- [ ] A profile, benchmark, or production metric — not intuition — identified the specific bottleneck
- [ ] The fix targets that measured bottleneck, not a generic "best practice" change
- [ ] Before/after numbers exist for the same measurement
- [ ] No N+1 queries or unbounded fetches were introduced by the fix itself
- [ ] A regression guard exists (CI benchmark baseline, Lighthouse budget, or equivalent) so the fix doesn't silently erode later
- [ ] `pprof` or equivalent debug endpoints are not exposed on a public listener
