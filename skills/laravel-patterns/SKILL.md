---
name: laravel-patterns
description: Guides Laravel application architecture — controllers, Eloquent, service/action layers, queues, and API resources. Use when building Laravel controllers or routes, working with Eloquent models and relationships, designing API resources, or adding queued jobs, events, or caching to a Laravel app.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Laravel Patterns

Keep controllers thin — orchestration lives in services, single-purpose
logic lives in actions. Eloquent's convenience is also its biggest
footgun: mass assignment and N+1 queries both look correct until they hit
real data volume or an untrusted request body.

## Layering

```
app/Http/Controllers/  → routing + response shape only
app/Http/Requests/     → validation (FormRequest classes)
app/Actions/           → single-purpose use cases
app/Services/          → coordinating domain logic across multiple actions/models
app/Models/            → Eloquent models, casts, scopes, relationships
```

```php
final class CreateReservationAction
{
    public function __construct(private ReservationRepository $reservations) {}

    public function handle(CreateReservationData $data): Reservation
    {
        return $this->reservations->create($data);
    }
}

final class ReservationsController extends Controller
{
    public function __construct(private CreateReservationAction $createReservation) {}

    public function store(StoreReservationRequest $request): JsonResponse
    {
        $reservation = $this->createReservation->handle($request->toDto());
        return response()->json(['data' => ReservationResource::make($reservation)], 201);
    }
}
```

## Gotchas

- **Mass assignment is opt-in trust, not automatic safety.** `$fillable`
  only allowlists the fields Eloquent will accept from `create()`/`update()`
  arrays — a model with `$guarded = []` (guard nothing) accepts every field
  in the incoming array, including ones like `is_admin` or `hotel_id` that a
  request body should never be allowed to set directly. Always define
  `$fillable` explicitly for any model that accepts request-derived data;
  never set `$guarded = []` on such a model.
- **N+1 queries hide in plain view.** `Reservation::all()` followed by
  `$reservation->guest->name` inside a loop issues one query per
  reservation. Use `->with(['guest'])` (eager loading) whenever a
  relationship is accessed inside a loop over a collection — this is the
  single most common Eloquent performance defect, and it's invisible in a
  dev database with 10 rows.
- **Route-model binding without `scopeBindings()` allows cross-tenant
  access** on nested routes: `/accounts/{account}/projects/{project}`
  without scoped bindings will resolve `{project}` globally, letting a
  request supply a `project` ID belonging to a different `account` than the
  one in the URL. Use `Route::scopeBindings()` on any nested resource route.
- Queued job handlers must be idempotent — Laravel's queue driver can
  redeliver a job (worker crash mid-job, a retry after a transient failure),
  so a handler that isn't safe to run twice (e.g., "increment a counter"
  instead of "set to this value") will double-apply on redelivery.

## Real-world grounding

Mass-assignment vulnerabilities are not theoretical — they're the exact
mechanism behind a well-known class of Rails and Laravel incidents from the
2010s where a request body containing an unexpected field (`role: "admin"`,
`is_verified: true`) silently set a privileged column because the
framework's default behavior was to accept whatever fields were present.
This is precisely why Laravel requires either `$fillable` (allowlist) or
`$guarded` (denylist) to be declared explicitly on every model — treat a
model with neither defined, or with `$guarded = []`, as a review-blocking
finding.

## Verification

- [ ] Every model accepting request data declares an explicit `$fillable`
- [ ] Relationships accessed inside a loop are eager-loaded with `->with()`
- [ ] Nested resource routes use `Route::scopeBindings()`
- [ ] Queue job handlers are safe to run more than once
- [ ] Controllers contain no business logic — only request/response shaping

See `references/eloquent-and-api-patterns.md` for query scopes, custom
casts, form requests, and API resource pagination examples.
