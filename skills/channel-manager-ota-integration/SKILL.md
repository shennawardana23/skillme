---
name: channel-manager-ota-integration
description: Guides integrating a hotel PMS with online travel agencies (OTAs) via a channel manager. Use when implementing rate/availability push, reservation pull, reconciling a channel-manager sync failure, or reviewing code that treats an OTA's view of availability as always in sync with the PMS.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "hotel-domain"
---

# Channel Manager / OTA Integration

A channel manager syncs rate and availability from a PMS out to online
travel agencies (Booking.com, Agoda, Expedia, and similar) and pulls
reservations back in. The defining constraint this skill exists for:
that sync is inherently eventually-consistent, not transactional — code
that assumes the OTA's view and the PMS's view are always identical will
misdiagnose a normal, transient state as a bug, or worse, will miss a
genuine desync that's silently overselling rooms.

## Treat sync lag as normal, not an error

A push (PMS → OTA: updated rate or availability) and a pull (OTA → PMS:
new reservation, modification, cancellation) each have their own latency
and failure surface. "OTA shows sold out while the PMS shows available"
or the reverse is a normal, expected transient state during the sync
window — code and alerting should distinguish "still catching up" from
"actually desynced," rather than treating every observed mismatch as an
incident.

## A rejected push is a silent oversell risk

If a rate/availability push to an OTA is rejected (malformed payload,
OTA-side validation failure, a transient API error) and that rejection
isn't retried and surfaced, the OTA continues selling against its last
successfully-synced state while the PMS believes the update already took
effect. This produces silent overselling from the OTA side with no
PMS-side signal that anything is wrong — a push must be treated as
failed-until-confirmed, with retry and alerting on sustained failure, not
fire-and-forget.

## Rate-plan mapping isn't 1:1 across OTAs

Each OTA has its own rate-plan model and its own mapping conventions
(cancellation policy taxonomy, meal-plan inclusion codes, promotion
eligibility rules) that don't automatically correspond to the PMS's
internal rate plans. Assuming a PMS rate plan maps cleanly to "the
equivalent" OTA rate plan without an explicit, maintained mapping table
per OTA produces either rejected pushes or, worse, a booking that arrives
with terms the PMS doesn't actually have a matching internal rate plan
for.

## Reservation events, not reservation state

OTA reservation modifications and cancellations commonly arrive as new
events (a "modify" message, a "cancel" message) rather than as an
in-place update to a previously-pulled reservation record. Treating every
incoming event as an idempotent upsert keyed by the OTA's own reservation
identifier — rather than assuming each pull represents the full current
state — avoids duplicate-booking bugs when an event is redelivered or
arrives out of order.

## Gotchas

- **A rate/availability mismatch between PMS and OTA during the normal
  sync window is not automatically a bug** — alerting that fires on every
  observed mismatch instead of on a mismatch that persists past the
  expected sync latency will generate constant false-positive noise and
  train people to ignore it.
- **A rejected push is a silent oversell risk if not retried and
  surfaced** — fire-and-forget push code looks correct until the first
  OTA-side validation failure, at which point the OTA keeps selling
  against stale data with nothing in the PMS flagging that the update
  never landed.
- **Rate-plan mapping is per-OTA, not universal** — a mapping table
  assumed to be "the same shape" across OTAs breaks the moment a second
  OTA is integrated with a different cancellation-policy taxonomy.
- **Reservation events must be handled idempotently by the OTA's own
  reservation ID** — treating a redelivered or out-of-order modification
  event as a fresh reservation instead of an upsert against the existing
  one creates duplicate bookings that are hard to reconcile after the
  fact.
- **A cancellation event arriving after a modification event (out of
  order) needs explicit handling** — processing events strictly in
  arrival order without a sequence/timestamp check can apply a stale
  cancellation on top of a newer modification, or vice versa.

## Real-world grounding

The push/pull, eventually-consistent nature of PMS-to-OTA channel
management is a standard, industry-wide integration shape — every
major channel manager (SiteMinder, RateGain, and similar) operates on
this same push-availability/pull-reservation model with the same class of
sync-window and mapping problems, which is why "is this actually desynced
or just catching up" is one of the first diagnostic questions any
hospitality engineering team learns to ask before escalating a
channel-sync alert.

## Verification

- [ ] Alerting distinguishes a mismatch within the expected sync window
      from one that has persisted past it
- [ ] Every rate/availability push has retry-on-failure and surfaces a
      sustained failure, rather than being fire-and-forget
- [ ] Rate-plan mapping is maintained per-OTA, not assumed universal
      across channels
- [ ] Incoming reservation events are handled as idempotent upserts keyed
      by the OTA's reservation ID, not assumed to always be a fresh
      record
- [ ] Out-of-order event arrival (a cancellation after a later
      modification, or vice versa) is handled explicitly, not just
      processed in arrival order
