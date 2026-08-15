# Observation schema examples

These are worked examples of the status/summary/next_actions/artifacts
shape from the main skill, applied to a few common Go-backed tool
categories. Adapt field names to your own harness; keep the shape.

## File edit tool

```json
{
  "status": "success",
  "summary": "applied 1 hunk to internal/booking/reservation.go",
  "next_actions": ["run go test ./internal/booking/... to confirm"],
  "artifacts": ["internal/booking/reservation.go"]
}
```

On failure (e.g. the file changed since it was last read):

```json
{
  "status": "error",
  "error_class": "stale_read",
  "summary": "file has changed on disk since your last read",
  "next_actions": ["re-read the file, then reapply the edit against current content"],
  "stop_condition": "if this repeats after a fresh read, stop and report — the file is being modified concurrently",
  "artifacts": ["internal/booking/reservation.go"]
}
```

## Database migration tool (micro-grained, high risk)

```json
{
  "status": "success",
  "summary": "applied migration 0042_add_notification_settings up on hotel_id-partitioned table",
  "next_actions": ["verify with a read against a sample hotel_id"],
  "artifacts": ["migrations/0042_add_notification_settings.sql"]
}
```

A migration tool should never silently no-op on "already applied" —
report it explicitly so the model doesn't assume it just ran:

```json
{
  "status": "warning",
  "summary": "migration 0042 was already applied; no changes made",
  "next_actions": ["confirm this is expected before continuing"],
  "artifacts": []
}
```

## Search/query tool

```json
{
  "status": "success",
  "summary": "3 matches for \"ReservationStatus\" across 2 files",
  "next_actions": ["read internal/booking/status.go:14 for the type definition"],
  "artifacts": ["internal/booking/status.go", "internal/booking/handlers.go"]
}
```

Zero matches is not an error — say so plainly instead of returning an
empty array with no explanation, which reads ambiguously as "the tool
call itself failed":

```json
{
  "status": "success",
  "summary": "no matches for \"ReservationStatus\" in the searched paths",
  "next_actions": ["try a broader search term or a different path scope"],
  "artifacts": []
}
```

## Rate-limited external API tool

```json
{
  "status": "error",
  "error_class": "rate_limited",
  "summary": "PMS API returned 429; retry-after 30s",
  "next_actions": ["wait at least 30s before retrying this exact call"],
  "stop_condition": "if rate-limited 3 times in a row, stop and report — do not busy-loop",
  "artifacts": []
}
```
