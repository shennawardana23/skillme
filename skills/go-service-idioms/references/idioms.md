# Go Service Idioms — Extended Reference

## Error boundary pattern

Wrap errors as they cross package boundaries, but stop wrapping and start
handling at the outermost layer that can actually do something about the
failure — log it, map it to an HTTP status, retry it, or send it to a queue's
dead-letter path.

```go
// internal/store/reservations.go
func (s *Store) Reservation(ctx context.Context, id string) (*Reservation, error) {
    var r Reservation
    err := s.db.GetContext(ctx, &r, reservationByIDQuery, id)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("reservation %s: %w", id, ErrNotFound)
    }
    if err != nil {
        return nil, fmt.Errorf("query reservation %s: %w", id, err)
    }
    return &r, nil
}

// internal/api/reservations.go — the boundary: decide what the caller sees
func (h *Handler) getReservation(w http.ResponseWriter, r *http.Request) {
    res, err := h.store.Reservation(r.Context(), reservationID(r))
    switch {
    case errors.Is(err, store.ErrNotFound):
        http.Error(w, "reservation not found", http.StatusNotFound)
        return
    case err != nil:
        h.logger.Error("get reservation failed", "err", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }
    writeJSON(w, res)
}
```

The store layer never logs — it only wraps and returns. The handler is the
single place that logs or maps to a status code. This avoids the common
anti-pattern of the same error being logged three times as it bubbles up
through three call sites.

## Goroutine leak: a concrete failure scenario

```go
// BUG: if the caller times out or stops reading, this goroutine blocks
// forever on the unbuffered channel send and is never collected.
func stream(ctx context.Context) <-chan Event {
    ch := make(chan Event)
    go func() {
        for _, e := range fetchAll() {
            ch <- e // no ctx check, no select — leaks on cancellation
        }
        close(ch)
    }()
    return ch
}
```

Fix: select on `ctx.Done()` alongside the send, or use a buffered channel
sized to the known bound and document why that bound is safe:

```go
func stream(ctx context.Context) <-chan Event {
    ch := make(chan Event)
    go func() {
        defer close(ch)
        for _, e := range fetchAll() {
            select {
            case ch <- e:
            case <-ctx.Done():
                return
            }
        }
    }()
    return ch
}
```

## Worker pool with bounded concurrency

```go
func processAll(ctx context.Context, items []Item, workers int) error {
    g, ctx := errgroup.WithContext(ctx)
    sem := make(chan struct{}, workers)
    for _, item := range items {
        item := item
        select {
        case sem <- struct{}{}:
        case <-ctx.Done():
            return ctx.Err()
        }
        g.Go(func() error {
            defer func() { <-sem }()
            return process(ctx, item)
        })
    }
    return g.Wait()
}
```

Use this over an unbounded `for range items { go ... }` whenever `items` can
be large or comes from external/user input — unbounded goroutine fan-out is
a resource-exhaustion vector, not just a style nit.

## Interface size

Define consumer-side interfaces at the point of use, sized to exactly what
that function needs:

```go
// internal/billing — this package only needs to read a reservation.
type reservationReader interface {
    Reservation(ctx context.Context, id string) (*Reservation, error)
}

func ChargeFor(ctx context.Context, r reservationReader, id string) error {
    res, err := r.Reservation(ctx, id)
    ...
}
```

`*store.Store` satisfies `reservationReader` without either package
importing the other's interface — `billing` depends on nothing but the
method set it uses, which keeps the dependency graph acyclic and makes
`ChargeFor` trivially testable with a hand-written fake.
