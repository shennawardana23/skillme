# Hotel-Partitioned Tables — Extended Reference

## Recognizing the bug in review

The missing-filter bug is easy to miss in review because the query still
*works* — it returns correct rows, just slower and at far higher cost than
necessary, so it rarely surfaces until a table has enough hotels' worth of
data for a full scan to show up in latency graphs or database load.

Grep signal when reviewing a diff that touches a partitioned table: does the
new query's `WHERE` clause (or the join's `ON` clause) mention `hotel_id`?
If a new query string is added and `hotel_id` does not appear in it at all,
that is a near-certain bug — flag it even without running `EXPLAIN`.

```go
// Looks correct, passes tests against a single-hotel test database,
// and silently scans every partition once other hotels exist:
func (s *Store) ByStatus(ctx context.Context, status string) ([]Reservation, error) {
    const q = `SELECT * FROM reservations WHERE status = $1`
    ...
}
```

The fix is not just adding a `hotel_id` argument — it's removing any code
path that can call the partitioned table without one. Prefer deleting
`ByStatus` above entirely in favor of `ByHotelAndStatus`, rather than adding
a second overload that keeps the unsafe path reachable.

## Dynamic hotel sets

`hotel_id = ANY($1)` with a parameterized array is still prunable and safe
for "all reservations across these three hotels" queries:

```sql
SELECT * FROM reservations WHERE hotel_id = ANY($1) AND status = $2;
```

```go
rows, err := s.db.QueryContext(ctx, q, pq.Array(hotelIDs), status)
```

An empty `hotelIDs` slice bound to `ANY($1)` returns zero rows (not an
error) — no need for a separate empty-slice guard purely for correctness,
though guarding it can still save a round trip.

## Partial indexes per hot hotel

If one hotel is disproportionately larger or hotter than the rest, a
per-partition partial index can still help beyond the partition pruning
itself:

```sql
CREATE INDEX ON reservations_101 (status, checkin) WHERE status <> 'cancelled';
```

Partition-local indexes are independent of each other — an index added to
`reservations_101` does not automatically exist on `reservations_102`;
migrations that add an index to a partitioned parent must ensure the index
propagates to every existing partition (`CREATE INDEX ON reservations (...)`
against the parent does propagate automatically to all partitions in
PostgreSQL 11+; an index created directly against one partition does not).
