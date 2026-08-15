---
name: backend-patterns
description: Guides backend service architecture — layering, N+1 query prevention, caching, retries, rate limiting, and background job processing, in Go primarily and TypeScript/PHP where the pattern is language-agnostic. Use when structuring handler/service/repository layers, reviewing a service for N+1 queries, adding a cache in front of a slow read, implementing retry logic, or setting up a background job queue.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Backend Patterns

Structural patterns for server-side code, independent of any one
framework. For the shape of the request/response contract itself, see
`api-design`; for query- and index-level database patterns, see
`postgres-patterns`; for structured logs/metrics/traces, see
`observability-and-instrumentation`.

## Layering: handler → service → repository

Keep three concerns separate so each can be tested and changed
independently:

```go
// repository: only knows SQL, returns domain types or a wrapped error
type ReservationRepo struct{ db *sql.DB }

func (r *ReservationRepo) ByID(ctx context.Context, hotelID int64, id string) (*Reservation, error) {
    const q = `SELECT id, hotel_id, status FROM reservations WHERE hotel_id = $1 AND id = $2`
    var res Reservation
    if err := r.db.QueryRowContext(ctx, q, hotelID, id).Scan(&res.ID, &res.HotelID, &res.Status); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("reservation %s: %w", id, ErrNotFound)
        }
        return nil, fmt.Errorf("query reservation %s: %w", id, err)
    }
    return &res, nil
}

// service: business logic, orchestrates repositories, knows nothing about HTTP
type ReservationService struct{ repo *ReservationRepo }

func (s *ReservationService) Cancel(ctx context.Context, hotelID int64, id string) error {
    res, err := s.repo.ByID(ctx, hotelID, id)
    if err != nil {
        return err
    }
    if res.Status == StatusCheckedOut {
        return fmt.Errorf("cancel reservation %s: %w", id, ErrAlreadyCheckedOut)
    }
    return s.repo.SetStatus(ctx, hotelID, id, StatusCancelled)
}

// handler: only knows HTTP — decode, call service, map errors to status codes
func (h *Handler) CancelReservation(w http.ResponseWriter, r *http.Request) {
    hotelID, id := parseParams(r)
    if err := h.svc.Cancel(r.Context(), hotelID, id); err != nil {
        writeError(w, err) // maps ErrNotFound -> 404, ErrAlreadyCheckedOut -> 409, else 500
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
```

The service layer must not import `net/http`, and the repository layer
must not contain business rules (e.g. "can't cancel after checkout"
belongs in the service, not buried in a `WHERE status != 'CHECKED_OUT'`
clause the caller can't see). This mirrors the interface-boundary
guidance in `api-and-interface-design`: the service depends on a small
consumer-defined repository interface, not a concrete `*sql.DB`.

## N+1 query prevention

The most common backend performance bug: fetching a list, then issuing
one query per row to fetch related data.

```go
// BAD: N+1 — one query per reservation to fetch its guest
reservations, _ := repo.ListByHotel(ctx, hotelID)
for _, r := range reservations {
    guest, _ := guestRepo.ByID(ctx, hotelID, r.GuestID) // N queries
    r.Guest = guest
}

// GOOD: batch-fetch every guest in one query
reservations, _ := repo.ListByHotel(ctx, hotelID)
guestIDs := make([]string, len(reservations))
for i, r := range reservations {
    guestIDs[i] = r.GuestID
}
guests, _ := guestRepo.ByIDs(ctx, hotelID, guestIDs) // 1 query, IN ($1, $2, ...) or ANY($1)
guestByID := make(map[string]*Guest, len(guests))
for _, g := range guests {
    guestByID[g.ID] = g
}
for _, r := range reservations {
    r.Guest = guestByID[r.GuestID]
}
```

`ByIDs` still needs the `hotel_id` predicate for partition pruning — see
`postgres-hotel-partitioning`. In an ORM (Prisma, Eloquent, GORM), the
same fix is an explicit eager load (`Preload`, `with()`, `include`), not
a loop calling `Find` per row.

## Caching: cache-aside with explicit invalidation

```go
func (s *ConfigService) Get(ctx context.Context) (*Config, error) {
    if cached, ok := s.cache.Get("config"); ok {
        return cached.(*Config), nil
    }
    cfg, err := s.repo.LoadConfig(ctx)
    if err != nil {
        return nil, err
    }
    s.cache.SetWithTTL("config", cfg, 5*time.Minute)
    return cfg, nil
}
```

Two failure modes to design against explicitly:

- **Stale-after-write**: any write path that changes cached data must
  invalidate (or update) the cache entry in the same transaction/request,
  not rely on TTL expiry alone — otherwise readers see stale data for the
  full TTL window after a write.
- **Thundering herd**: when a hot key expires, many concurrent requests
  can all miss the cache and hit the database at once. Use
  single-flight (Go's `golang.org/x/sync/singleflight`) so only one
  request per key actually queries the backing store while the rest wait
  for its result.

## Retries with backoff — only for retryable failures

```go
func withRetry(ctx context.Context, maxAttempts int, fn func() error) error {
    var err error
    for attempt := 0; attempt < maxAttempts; attempt++ {
        if err = fn(); err == nil {
            return nil
        }
        if !isRetryable(err) {
            return err // don't retry a 400 or a validation error
        }
        backoff := time.Duration(1<<attempt) * 100 * time.Millisecond
        jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
        select {
        case <-time.After(backoff + jitter):
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return fmt.Errorf("after %d attempts: %w", maxAttempts, err)
}
```

Jitter matters at scale: without it, every client backing off from the
same outage retries in lockstep and re-creates the load spike each
interval. Only retry idempotent operations, or operations made idempotent
via an idempotency key — retrying a bare `POST /charge` can double-charge.

## Rate limiting

Prefer a token-bucket limiter (`golang.org/x/time/rate` in Go) over a
hand-rolled sliding window — it's O(1) per check and battle-tested:

```go
limiter := rate.NewLimiter(rate.Limit(100), 20) // 100 req/s, burst 20

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

A single process-wide limiter only protects one instance; behind a load
balancer with N instances, the effective limit is N times higher. For a
service-wide limit, back the limiter with Redis (`INCR` + `EXPIRE`, or a
Lua script for atomicity) instead of in-memory state.

## Background jobs: durable queue, not an in-memory slice

An in-memory job queue (a slice plus a goroutine) loses every queued job
on process restart or crash. For anything that must survive a restart,
use a durable queue: Postgres-backed (`SELECT ... FOR UPDATE SKIP LOCKED`,
see `postgres-patterns`), or a dedicated queue (SQS, Redis Streams, NATS
JetStream). Design every job handler to be safely re-run (idempotent) —
a worker can crash after doing the work but before acknowledging it, and
the job will be redelivered.

## Gotchas

- `SELECT *` in a repository method silently breaks the moment someone
  adds a column with a default the code doesn't expect, or drops one the
  code still scans into — select named columns explicitly.
- A service method that returns `(T, error)` and also logs the error
  internally causes it to be logged twice once the caller (correctly)
  logs it again at the boundary — log once, at the boundary, and wrap the
  error with context on the way up instead.
- A cache with no TTL and no invalidation path is not a cache, it's a
  slow memory leak that also serves stale data forever.
- Rate limiting only at the load balancer, and not per-user/per-key
  inside the service, means one abusive API key can still starve every
  other tenant sharing the aggregate limit.

## Real-world grounding

The repository/service/handler split described here is the backend
analogue of the layered architecture popularized under names like "ports
and adapters" (hexagonal architecture, Alistair Cockburn, 2005) and
"clean architecture" (Robert C. Martin) — both make the same core claim
this skill does: business logic (the service layer) should have zero
import-level dependency on delivery mechanism (HTTP) or storage
mechanism (SQL driver), so either can be swapped or unit-tested without
the other.

## Verification

- [ ] Service layer has no `net/http` (or framework) import; repository layer has no business rules
- [ ] Every list-then-fetch-related-data path is a single batched query, not a loop
- [ ] Every cache write path has an explicit invalidation or update on the corresponding write path
- [ ] Retries apply jitter and skip non-retryable errors; retried operations are idempotent
- [ ] Rate limiting is enforced per key/tenant, not only in aggregate
- [ ] Background jobs are durable (survive a process restart) and job handlers are idempotent
