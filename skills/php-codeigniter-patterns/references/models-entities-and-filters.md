# PHP CodeIgniter Patterns — Reference

## Entity `$datamap`: remapping column names without touching callers

```php
class Reservation extends \CodeIgniter\Entity\Entity
{
    // Maps property name `guestName` to the actual `guest_name` column,
    // so callers use $reservation->guestName regardless of the schema's
    // real column name.
    protected $datamap = [
        'guestName' => 'guest_name',
    ];

    protected $casts = [
        'is_cancelled' => 'boolean',
        'metadata'     => 'json-array',
        'status'       => 'enum[App\Enums\ReservationStatus]',
    ];
}
```

Supported built-in cast types include `integer`, `float`, `string`,
`boolean`, `array`, `json`, `json-array`, `csv`, `datetime`, `timestamp`, and
`enum[ClassName]`; a custom cast extends `CodeIgniter\Entity\Cast\BaseCast`
for anything not covered (e.g. a value object).

## Filter execution order (CI4 v4.5.0+)

`Config\Filters` has four applicable arrays, evaluated in this order for
"before" filters, and reverse for "after" filters:

```
before: required → globals → methods → filters → [controller runs] → route
after:  route → filters → globals → required
```

- **`$required`** (v4.5.0+) — filters that run on *every* request,
  intended for framework-level concerns (e.g. `forcehttps`, `pagecache`);
  not meant to be edited per-project in most apps.
- **`$globals`** — filters applied to all valid requests, with an `except`
  key to exclude specific URI patterns:
  ```php
  public array $globals = [
      'before' => ['csrf' => ['except' => ['api/webhooks/*']]],
      'after'  => ['toolbar'],
  ];
  ```
- **`$methods`** — keyed by HTTP method (`'post' => ['csrf']`), before-only.
- **`$filters`** — keyed by alias, with explicit URI patterns and optional
  arguments passed to the filter (`'group:admin,superadmin'` becomes
  `$arguments = ['admin', 'superadmin']` inside `before()`):
  ```php
  public array $filters = [
      'hotelscope' => ['before' => ['hotels/*']],
  ];
  ```

Nested `$routes->group()` calls **merge** their `filter` option with the
enclosing group's — a filter attached to an outer group still runs for
routes defined in an inner group, it does not need to be repeated.

## Query Builder join patterns (no relationship methods exist)

```php
// One-to-many: fetch parents, then children keyed by parent ID
$hotels = $this->hotelModel->whereIn('id', $hotelIds)->findAll();
$rooms  = $this->roomModel->whereIn('hotel_id', $hotelIds)->findAll();
$roomsByHotel = [];
foreach ($rooms as $room) {
    $roomsByHotel[$room->hotel_id][] = $room;
}

// Or a single JOIN when the result set is naturally flat
$rows = $this->db->table('reservations')
    ->select('reservations.*, hotels.name AS hotel_name')
    ->join('hotels', 'hotels.id = reservations.hotel_id')
    ->where('reservations.status', 'confirmed')
    ->get()
    ->getResult();
```

Neither pattern is "the CI4 way" over the other — pick per query based on
whether the parent/child shape or a flat row shape is more useful to the
caller, and be consistent within one feature area so reviewers aren't
guessing which pattern a new query should follow.
