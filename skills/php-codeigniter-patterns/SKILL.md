---
name: php-codeigniter-patterns
description: Guides CodeIgniter 4 application architecture — Controllers, Models (Query Builder, not an ORM), Entities, the Validation library, Filters (CI4's middleware equivalent), and Services. Use when building a CI4 controller or route, working with a CI4 Model or Entity, wiring a Filter to a route group, or adding a custom Service.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# PHP CodeIgniter Patterns

This is the architecture layer for CodeIgniter 4 (CI4) applications: how a
request flows from Router → Filters → Controller → Model → Entity, and where
each concern belongs. It is not a testing guide (`php-codeigniter-tdd`), not
a security checklist (`php-codeigniter-security`), and not a pre-deploy
verification pipeline (`php-codeigniter-verification`) — apply those
alongside this one.

**The single biggest CI4-vs-Laravel difference**: CI4's `Model` class is
**not an ORM**. It has no relationship methods (`hasMany`, `belongsTo`), no
eager loading, and no Active Record-style object graph. A `Model` is a thin
wrapper that exposes Query Builder plus a handful of lifecycle conveniences
(mass-assignment guarding, timestamps, soft deletes, validation). Per the
[CI4 Model user guide](https://codeigniter4.github.io/userguide/models/model.html):
"Once you get the Query Builder instance, you can call methods of the Query
Builder. However, since Query Builder is not a Model, you cannot call
methods of the Model." Related rows are fetched with an explicit second
query or a manual `JOIN` in the Query Builder call — there is no `->with()`
to reach for, and no N+1-via-eager-loading gotcha to guard against, because
eager loading doesn't exist. The N+1 risk in CI4 is the mirror image: a
developer coming from Laravel/Eloquent instinctively expects a relationship
method to exist and is surprised when it doesn't.

## Layering

```
app/Controllers/  → routing + request/response shaping only
app/Filters/       → CI4's middleware equivalent (auth, rate limiting, CORS)
app/Models/        → Query Builder wrapper: allowedFields, validation, casts
app/Entities/      → typed row objects (optional; set via Model::$returnType)
app/Config/        → Services.php, Filters.php, Validation.php, Routes.php
```

## Models: Query Builder, not an ORM

```php
namespace App\Models;

use CodeIgniter\Model;

class ReservationModel extends Model
{
    protected $table            = 'reservations';
    protected $primaryKey       = 'id';
    protected $returnType       = \App\Entities\Reservation::class;
    protected $useTimestamps    = true;
    protected $useSoftDeletes   = true;

    // Fields the model will accept from insert()/update()/save() arrays.
    // Never include $primaryKey here.
    protected $allowedFields = ['guest_id', 'hotel_id', 'check_in', 'check_out', 'status'];

    protected $validationRules = [
        'guest_id'  => 'required|integer',
        'check_in'  => 'required|valid_date',
        'check_out' => 'required|valid_date',
    ];
}
```

`$allowedFields` is CI4's mass-assignment guard: "Any field names other than
these will be discarded" ([Model user
guide](https://codeigniter4.github.io/userguide/models/model.html)) — it is
the exact same class of protection as Laravel's `$fillable`, just under a
different name. `$validationRules` on the Model runs automatically just
before `insert()`/`update()`/`save()` persist data, so a Model that skips
`$validationRules` and relies solely on controller-level validation loses
protection for any other code path (a console command, a queued job) that
calls the Model directly.

Since CI4 v4.3.0, `update()` raises a `DatabaseException` if the generated
`UPDATE` statement would have no `WHERE` clause — a deliberate guard against
a Model call that accidentally updates every row in the table.

Reaching into the underlying Query Builder for anything beyond simple CRUD:

```php
$reservations = $this->reservationModel
    ->where('hotel_id', $hotelId)
    ->where('status', 'confirmed')
    ->orderBy('check_in', 'ASC')
    ->findAll();

// Manual join — there is no relationship method to call instead
$rows = $this->reservationModel
    ->select('reservations.*, guests.name AS guest_name')
    ->join('guests', 'guests.id = reservations.guest_id')
    ->where('reservations.hotel_id', $hotelId)
    ->findAll();
```

## Entities: typed row objects

Setting `Model::$returnType` to an Entity class makes `find()`/`findAll()`
return typed objects instead of arrays or `stdClass`:

```php
namespace App\Entities;

use CodeIgniter\Entity\Entity;

class Reservation extends Entity
{
    protected $casts = [
        'check_in'  => 'datetime',
        'check_out' => 'datetime',
        'metadata'  => 'json-array',
    ];

    // snake_case column `check_in` maps to CamelCase-prefixed accessor
    public function setCheckIn(string $value): static
    {
        $this->attributes['check_in'] = date('Y-m-d H:i:s', strtotime($value));
        return $this;
    }
}
```

Entities expose magic `__get()`/`__set()` that look for a `getX()`/`setX()`
method (snake_case column → PascalCase-prefixed method) before falling back
to raw attribute access — this is where per-field business logic (hashing,
normalizing, casting) belongs, not scattered across controllers.

## Validation: controller-level, and the `getVar()` gotcha

```php
public function update(int $id)
{
    if (! $this->validateData($this->request->getPost(), [
        'check_in'  => 'required|valid_date',
        'check_out' => 'required|valid_date|greater_than_field[check_in]',
    ])) {
        return $this->response->setJSON(['errors' => $this->validator->getErrors()])
            ->setStatusCode(422);
    }
    // ...
}
```

Reusable rule groups belong in `app/Config/Validation.php` as named public
properties, run via `$validation->run($data, 'groupName')` — this is the CI4
equivalent of a Laravel `FormRequest` class:

```php
// app/Config/Validation.php
public array $reservationCreate = [
    'guest_id'  => 'required|integer|is_not_unique[guests.id]',
    'check_in'  => 'required|valid_date',
    'check_out' => 'required|valid_date|greater_than_field[check_in]',
];
```

## Filters: CI4's middleware equivalent

```php
namespace App\Filters;

use CodeIgniter\Filters\FilterInterface;
use CodeIgniter\HTTP\RequestInterface;
use CodeIgniter\HTTP\ResponseInterface;

class HotelScopeFilter implements FilterInterface
{
    public function before(RequestInterface $request, $arguments = null)
    {
        if (! session()->get('hotel_id')) {
            return redirect()->to('/login');
        }
    }

    public function after(RequestInterface $request, ResponseInterface $response, $arguments = null)
    {
        // Filters can only modify/return the Response here, not stop execution.
    }
}
```

```php
// app/Config/Filters.php
public array $aliases = [
    'hotelscope' => \App\Filters\HotelScopeFilter::class,
];

public array $globals = [
    'before' => ['csrf'],
    'after'  => [],
];
```

```php
// app/Config/Routes.php — attach a filter to a route group
$routes->group('hotels/(:num)', ['filter' => 'hotelscope'], static function ($routes) {
    $routes->get('reservations', 'ReservationController::index');
});
```

The filter value must match an alias defined in `$aliases` — a typo here
fails silently at the routing layer (the route still matches, the filter
just never runs) rather than throwing, so a missing filter on a route group
is easy to miss in review.

## Services: the service locator

Custom cross-cutting collaborators are registered in
`app/Config/Services.php` (extends `CodeIgniter\Config\BaseService`) and
resolved through the `service()` helper or `\Config\Services::name()` — both
are equivalent. Services default to a **shared** instance (same object
returned on every call in one request); pass `false` for a fresh instance:

```php
// app/Config/Services.php
public static function pricingEngine(bool $getShared = true): PricingEngine
{
    if ($getShared) {
        return static::getSharedInstance('pricingEngine');
    }
    return new PricingEngine(config('Pricing'));
}
```

```php
$engine = service('pricingEngine');          // shared instance
$fresh  = service('pricingEngine', false);   // new instance
```

## Gotchas

- **`$this->validate()` reads `$_GET`, `$_POST`, and `$_COOKIE`, in that
  order** (via `$request->getVar()`) — a request with a cookie named the
  same as a validated field can have its cookie value validated instead of
  the POST body, and the CI4 docs explicitly recommend `validateData()`
  (validates an explicit array — typically `$this->request->getPost()`)
  for POST-only validation instead of the ambiguous `validate()` shorthand.
- **`$allowedFields` must never include the primary key** — the Model user
  guide states this explicitly; including it lets a client-supplied `id` in
  a `create`-style form silently retarget an `update()` call or corrupt
  autoincrement behavior on `insert()`.
- **There is no eager loading, because there is no ORM.** A developer
  porting Laravel muscle memory will look for `->with(['guest'])` and not
  find it — the correct CI4 pattern is either a second `findAll()` keyed by
  the foreign IDs, or a `JOIN` in the Query Builder call. Neither is wrong;
  picking one inconsistently across a codebase is a maintenance cost, not a
  functional bug.
- **A `Filters.php` alias typo fails silently.** The route still resolves
  and the controller still runs — the filter that was supposed to gate it
  simply never executes. Confirm a new filter alias with a test request,
  not just a code read.
- **`update()` without a `WHERE`-implying primary key throws since v4.3.0**
  — code written against an older CI4 version that relied on this being a
  silent full-table update will break on upgrade; treat the exception as
  the framework catching a bug, not a regression to work around.

## Real-world grounding

CI4's decision to keep `Model` as a Query Builder wrapper rather than
building a full ORM (Eloquent- or Doctrine-style) is a documented, explicit
design choice by the framework — the user guide's opening line on the topic
is functionally "Query Builder is not a Model, and a Model is not Query
Builder," drawing the boundary on purpose rather than as an omission. Teams
migrating a CI3 (pre-4.0) codebase, which had no formal Model class at all
and typically hand-rolled Query Builder calls directly in controllers, tend
to under-use `$allowedFields`/`$validationRules` because the old codebase
never had anywhere to put them — this is the most common defect this skill
catches in a CI3→CI4 migration review.

## Verification

- [ ] Every Model accepting request-derived data declares `$allowedFields`
      (never leaves it empty while still calling `insert()`/`update()`/`save()`
      with raw request arrays)
- [ ] `$allowedFields` never includes the primary key column
- [ ] Model-level `$validationRules` exist for any Model reachable from a
      non-HTTP path (console command, queued job), not just controller checks
- [ ] Controllers use `validateData()` (or an explicit rule group via
      `Config\Validation`) rather than bare `validate()` when the intent is
      "validate POST data only"
- [ ] Every route group requiring auth/tenant-scoping has a `filter` key
      whose value matches a real alias in `Config\Filters::$aliases`
- [ ] Cross-cutting collaborators are registered in `Config\Services.php`
      rather than `new`'d directly inside controllers

See `references/models-entities-and-filters.md` for Entity `$datamap`,
custom casts, filter execution order across `$required`/`$globals`/
`$methods`/`$filters`, and Query Builder join patterns.
