---
name: php-codeigniter-tdd
description: Drives CodeIgniter 4 test development — CIUnitTestCase, DatabaseTestTrait's migrate/refresh/seed lifecycle, the dedicated `tests` database group, and FeatureTestTrait for HTTP-level assertions. Use when writing or fixing a CodeIgniter 4 test, when a test touches the database, or when a test suite is unexpectedly slow or appears to be affecting the wrong database.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# PHP CodeIgniter TDD

Red-green-refactor for CodeIgniter 4 (CI4), with the framework-specific parts:
which trait wires up the database per test, why CI4 insists on a separate
`tests` database group, and how to assert against a full HTTP request/response
cycle instead of calling a controller method directly.

For the generic RED-GREEN-REFACTOR cycle, see `skills/test-driven-development/`.
This skill covers what's specific to CI4 on top of that cycle.

## Base class and traits

Every CI4 test extends `CodeIgniter\Test\CIUnitTestCase`. Add
`CodeIgniter\Test\DatabaseTestTrait` only when the test actually touches the
database — it adds real per-test cost (migration/refresh/seed), so a pure
unit test (a value object, a calculator) should not use it.

```php
use CodeIgniter\Test\CIUnitTestCase;
use CodeIgniter\Test\DatabaseTestTrait;

final class ReservationModelTest extends CIUnitTestCase
{
    use DatabaseTestTrait;

    protected $refresh = true;
    protected $seed    = 'ReservationSeeder';

    public function testCreateRejectsCheckOutBeforeCheckIn(): void
    {
        $model = new \App\Models\ReservationModel();

        $result = $model->insert([
            'guest_id'  => 1,
            'check_in'  => '2026-06-10',
            'check_out' => '2026-06-05',
        ]);

        $this->assertFalse($result);
        $this->assertArrayHasKey('check_out', $model->errors());
    }
}
```

## The `tests` database group — not optional, not decorative

`DatabaseTestTrait` requires a dedicated `tests` connection group defined in
`app/Config/Database.php`, separate from `default`. Per the CI4 user guide,
this exists specifically so database tests never touch "your other data" —
the trait's `$migrate`/`$refresh`/`$seed` behavior runs migrations, seeds, and
(with `$refresh = true`) rolls back or re-migrates the schema, against
whichever connection the `tests` group actually points to.

```php
// app/Config/Database.php
public array $tests = [
    'DSN'      => '',
    'hostname'  => 'localhost',
    'database'  => 'app_test',        // a real, separate database
    'username'  => 'root',
    'password'  => '',
    'DBDriver'  => 'MySQLi',
    'port'      => 3306,
];
```

The real, documented risk this guards against: if `tests` is left pointing
at the same database as `default` (copy-pasted config, or an env var that
happens to resolve to the same DSN in a shared dev environment), running the
suite migrates/refreshes/seeds against live development data — the trait
does not verify the `tests` group is actually distinct from `default` before
running. Confirm the `tests` group's database name is genuinely separate
before the first `DatabaseTestTrait` test ever runs, not after a suite
mysteriously wipes local data.

## Trait properties: migrate, migrateOnce, refresh, seed

- **`$migrate`** — run migrations for this test class (defaults true when
  the trait is used).
- **`$migrateOnce`** — run migrations only once for the whole run, not
  per-test; faster for a large migration set when tests don't need a fully
  clean schema between each one.
- **`$refresh`** — roll back and re-run migrations between tests for a
  clean slate; the safer default for tests that mutate data, at a real
  time cost.
- **`$seed`** — a seeder class name to populate baseline fixture data after
  migration.
- **`$namespace`** — where test-specific migrations live, defaulting to
  `tests/_support/Database/Migrations` rather than the app's own
  `app/Database/Migrations` — a test suite with its own schema fixtures
  doesn't have to touch the app's real migration set.

## Feature tests: HTTP-level, not controller-method-level

```php
use CodeIgniter\Test\CIUnitTestCase;
use CodeIgniter\Test\FeatureTestTrait;

final class ReservationsControllerTest extends CIUnitTestCase
{
    use FeatureTestTrait;

    public function testCreateReservationReturns201(): void
    {
        $result = $this->post('/api/reservations', [
            'guest_id'  => 1,
            'check_in'  => '2026-06-10',
            'check_out' => '2026-06-15',
        ]);

        $result->assertStatus(201);
        $result->assertJSONFragment(['status' => 'confirmed']);
    }
}
```

`FeatureTestTrait`'s `$this->get()`/`$this->post()`/`$this->call()` run a
request through the full router → filters → controller stack — this is what
proves a Filter alias is actually attached to the right route group (a
`Filters.php` typo, per `php-codeigniter-patterns`' Gotchas, fails silently
at the routing layer and a feature test hitting the real route is the
cheapest way to catch it). A test that calls the controller method directly
(`(new ReservationsController())->store()`) never exercises filters at all.

## Gotchas

- **A test class using `DatabaseTestTrait` against a `tests` group that
  isn't genuinely separate from `default` can migrate/refresh/seed real
  development data.** Verify the `tests` group's database name before
  trusting the suite is safe to run freely.
- **`DatabaseTestTrait` methods run in `setUp()`/`tearDown()`** — overriding
  either method in a test class without calling `parent::setUp()` /
  `parent::tearDown()` silently disables migration/refresh/seed, and the
  test may pass or fail for the wrong reason (stale schema, leftover rows
  from a previous test).
- **`$refresh = true` is the safer default but the slower one.** A suite
  that's unexpectedly slow is often using `$refresh` on every test class
  when `$migrateOnce` plus transaction-scoped test data would suffice —
  profile before assuming the framework itself is slow.
- **Feature tests via `FeatureTestTrait` still hit Filters and Config**,
  including CSRF (see `php-codeigniter-security`) — a feature test that
  POSTs without disabling or accounting for CSRF in the test environment
  fails for a reason unrelated to the code under test; CI4's test
  environment typically has `Config\Security::$csrfProtection` handled via
  `WithHeaders`/`skipCSRF` helpers on the test request rather than a
  disabled filter in production config.

## Real-world grounding

The CI4 user guide's explicit framing — "keep your other data safe" — for
why `DatabaseTestTrait` requires its own `tests` group is not a stylistic
preference; it directly addresses the same class of incident every
framework's test-database isolation exists to prevent: a test suite that
silently mutates or truncates a shared database because no one configured a
genuinely separate one. This is the single most consequential setup step in
this skill and the one most likely to be skipped when a project scaffolds
quickly from a tutorial rather than from the full user guide.

## Verification

- [ ] Every test touching the database uses `DatabaseTestTrait`, not a bare
      `CIUnitTestCase` making live queries
- [ ] `app/Config/Database.php`'s `tests` group points at a database
      genuinely separate from `default`
- [ ] Overridden `setUp()`/`tearDown()` call their `parent::` equivalents
- [ ] Tests that must prove a Filter/route wiring works use
      `FeatureTestTrait`'s HTTP-level calls, not a direct controller
      method call
- [ ] `$refresh`/`$migrateOnce` choice matches the test's actual isolation
      need, not copied unreflectively from another test class
