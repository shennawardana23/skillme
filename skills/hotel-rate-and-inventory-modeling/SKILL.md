---
name: hotel-rate-and-inventory-modeling
description: Guides modeling hotel rate plans, room inventory, and availability/allotment for a multi-brand, multi-property hotel management company. Use when designing a rate plan or room-type schema, implementing availability or overbooking logic, modeling allotment across a channel manager, or reviewing a design that conflates rate and inventory concerns.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "hotel-domain"
---

# Hotel Rate and Inventory Modeling

Rate and inventory are two independent axes that a schema or service layer
must keep genuinely separate — a **room type** (what physically exists to
sell) and a **rate plan** (a price/terms combination for selling it) are
not the same entity and don't share a lifecycle. This is the schema
distinction most new hotel-domain code gets wrong, and everything else in
this skill follows from keeping it straight. For the database
partitioning/query mechanics once a schema is decided, see
`postgres-hotel-partitioning` and `erd-expert` — this skill covers the
domain modeling above that layer.

## Room type vs. rate plan: two independent axes

A room type (deluxe, suite, standard twin) is a physical/product concept:
how many exist, their capacity, their amenities. A rate plan (best
available rate, non-refundable, member rate, corporate rate) is a
commercial concept layered on top: a price, cancellation terms, and
inclusions, for a given room type and date range. The same room type
carries multiple simultaneous rate plans; modeling a "rate" as if it were
a property of the room type itself (rather than a many-to-many
relationship keyed by room type × rate plan × date) causes combinatorial
schema explosion the moment a new promotional rate is added, since it
looks like it needs a new room-type row instead of a new rate-plan row.

## Allotment vs. raw inventory

**Allotment** is the count of rooms released to a specific sales channel
(a specific OTA, a specific rate plan, a group block) — it is a
*subset assignment* against physical inventory, not the inventory count
itself. Tracking only raw room-type inventory and treating every channel
as competing for the same undifferentiated pool makes controlled
overbooking (a deliberate yield-management decision) indistinguishable
from an accidental double-booking. A design needs both: total physical
inventory, and a separate allotment ledger per channel/rate plan that sums
to no more than that total (with an explicit, bounded overbooking margin
where the business intends one).

## Rate cascading

Rate plans commonly derive from one another rather than existing
independently — a non-refundable rate is often "base rate minus a fixed
discount," a member rate "base rate minus a percentage." Modeling each
derived rate as an independently-stored row (rather than a computed
relationship referencing its base rate) causes silent drift: the base
rate changes, the derived rates don't, and the two fall out of the
relationship the business actually intends. Store the derivation
relationship explicitly (a base-rate reference plus a rule), not just the
resulting number.

## Multi-brand, multi-property isolation

A multi-brand hotel group needs isolation on two axes simultaneously:
property (each hotel's own inventory and rates) and brand (reporting,
loyalty program rules, and rate strategy that can span every property
under one brand). A tenancy model scoped only by property — even a
well-designed one, per `postgres-hotel-partitioning` — can still leak
cross-brand reporting incorrectly if brand isn't also modeled as a first-
class dimension rather than inferred from which properties happen to
exist under it at query time.

## Gotchas

- **Conflating room type and rate plan into one entity** is the single
  most common schema mistake here — it works fine until a second rate
  plan needs adding, at which point it looks like it requires duplicating
  room-type rows instead of adding a rate-plan row.
- **Overbooking is often intentional, not a bug** — a design that can't
  distinguish "we deliberately allotted more to this channel than
  physical rooms exist, within a managed margin" from "two channels sold
  the same physical room by accident" will either block legitimate yield
  management or fail to catch a real double-booking, depending on which
  way the bug leans.
- **Derived rates stored as independent values drift from their base
  rate** silently — a change to the base rate doesn't propagate, and the
  discount relationship the business thinks still holds quietly stops
  being true.
- **Brand-level isolation is easy to omit when property-level isolation
  already exists** — a schema that partitions cleanly by property can
  still produce wrong brand-level rollups if brand membership isn't
  modeled explicitly and instead inferred by joining through properties.
- **Allotment ledgers need to sum to a bounded total**, not just exist per
  channel — without an explicit reconciliation check, per-channel
  allotments can silently exceed physical inventory with no single query
  surfacing that they no longer add up.

## Real-world grounding

Rate plan hierarchies (a "book direct" or "member" rate derived from a
publicly-published base rate) and channel allotment are standard,
long-established hospitality revenue-management concepts, not
company-specific inventions — the same distinction between physical room
inventory and channel-level allotment underlies how every PMS/channel
manager integration in the industry works, which is why getting it wrong
at the schema level tends to surface later as a channel-sync
reconciliation problem (see `channel-manager-ota-integration`) rather than
as an obvious modeling bug up front.

## Verification

- [ ] Room type and rate plan are modeled as separate entities in a
      many-to-many relationship, not one collapsed concept
- [ ] Allotment is tracked per channel/rate plan, separate from raw
      physical room-type inventory, and reconciles to a bounded total
- [ ] Derived rate plans reference their base rate and derivation rule,
      rather than storing an independently-drifting number
- [ ] Brand is modeled as an explicit dimension, not inferred solely from
      which properties currently belong to it
- [ ] A deliberate, bounded overbooking margin is distinguishable in the
      schema from an unintended double-booking
