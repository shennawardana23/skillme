---
name: laravel-tdd
description: Drives Laravel feature development with Pest or PHPUnit — factories, RefreshDatabase, fakes for queues/mail/notifications, Sanctum auth, and Inertia assertions. Use when writing or fixing a Laravel controller, Eloquent model, policy, job, or notification, when the project uses Pest/PHPUnit, or when asked to test an Inertia page's props.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Laravel TDD

Red-green-refactor for Laravel, with the framework-specific parts that
differ from a generic TDD loop: which database trait to reach for, which
`Facade::fake()` to use for a given side effect, and how to assert an
Inertia response instead of raw JSON.

For the generic RED-GREEN-REFACTOR cycle and the Prove-It bug-fix pattern,
see `skills/test-driven-development/`. This skill covers what's specific to
Laravel on top of that cycle.

## Test layers

- **Unit** — pure PHP: value objects, calculators, services with no
  framework dependency.
- **Feature** — HTTP endpoints, auth, validation, policies, response shape.
- **Integration** — database + queue + external boundaries together.

Pick the layer by what you need to prove, not by habit: a validation rule
belongs in a Feature test that hits the route (it proves the `FormRequest`
is wired in), not a Unit test that instantiates the rule class directly.

## Framework choice

Default to **Pest** for new test files when the project has it installed;
use **PHPUnit** only if the project already standardizes on it. Don't mix
styles within one test file.

```php
// Pest
test('owner can create project', function () {
    $user = User::factory()->create();

    $response = actingAs($user)->postJson('/api/projects', ['name' => 'New Project']);

    $response->assertCreated();
    assertDatabaseHas('projects', ['name' => 'New Project']);
});
```

```php
// PHPUnit — same behavior, class-based
final class ProjectControllerTest extends TestCase
{
    use RefreshDatabase;

    public function test_owner_can_create_project(): void
    {
        $user = User::factory()->create();

        $response = $this->actingAs($user)->postJson('/api/projects', ['name' => 'New Project']);

        $response->assertCreated();
        $this->assertDatabaseHas('projects', ['name' => 'New Project']);
    }
}
```

## Database trait: which one and why

- **`RefreshDatabase`** — the default for feature/integration tests. On a
  connection with transaction support it migrates once per test run (a
  static flag) and wraps each test in a rolled-back transaction; on
  `:memory:` SQLite (no cross-connection transaction visibility) it
  re-migrates before every test, which is slower.
- **`DatabaseTransactions`** — use when the schema is already migrated
  (e.g. a persistent test database) and you only need per-test rollback,
  skipping the migration-freshness check `RefreshDatabase` does.
- **`DatabaseMigrations`** — full migrate-and-drop every test. Reserve for
  tests that specifically verify migration behavior; it's the slowest
  option and unnecessary for ordinary feature tests.

## Factories and states

Define named states for the edge cases you'll actually test, not just the
happy-path row:

```php
$admin = User::factory()->state(['role' => 'admin'])->create();
$archived = Project::factory()->archived()->create();
```

## Faking side effects

Don't let a test send a real email, dispatch a real job, or hit a real
external API. Match the fake to the side effect:

```php
Queue::fake();
dispatch(new SendOrderConfirmation($order->id));
Queue::assertPushed(SendOrderConfirmation::class);

Notification::fake();
$user->notify(new InvoiceReady($invoice));
Notification::assertSentTo($user, InvoiceReady::class);

Http::fake();
// ... code that calls an external API ...
Http::assertSent(fn ($request) => $request->url() === 'https://api.example.com/charge');
```

`Bus::fake()` for jobs dispatched via `Bus::dispatch`/`dispatchNow`,
`Mail::fake()` for `Mail::send`/`Mailable::send`, `Event::fake()` for
domain events — pick the facade that matches how the code actually
triggers the side effect, not the closest-sounding one.

## Auth in tests

```php
use Laravel\Sanctum\Sanctum;

Sanctum::actingAs($user);
$response = $this->getJson('/api/projects');
$response->assertOk();
```

For policy checks directly (without a full HTTP round trip):

```php
$this->assertTrue(Gate::forUser($user)->allows('update', $project));
$this->assertFalse(Gate::forUser($otherUser)->allows('update', $project));
```

## Inertia responses

Assert on the component name and props, not on raw JSON — Inertia
responses aren't meant to be consumed as a plain API contract:

```php
$response = $this->actingAs($user)->get('/dashboard');

$response->assertInertia(fn (AssertableInertia $page) => $page
    ->component('Dashboard')
    ->where('user.id', $user->id)
    ->has('projects')
);
```

## Coverage and commands

- `php artisan test`, `vendor/bin/pest`, or `vendor/bin/phpunit`.
- Coverage needs `pcov` installed or `XDEBUG_MODE=coverage` set — without
  one of these, `--coverage` silently reports 0% rather than failing loudly.
- Set `DB_CONNECTION=sqlite` / `DB_DATABASE=:memory:` in the test
  environment (`phpunit.xml` or `.env.testing`) for fast, isolated runs
  that can't touch dev/prod data.

## Gotchas

- `RefreshDatabase` on `:memory:` SQLite re-migrates per test — if a suite
  is unexpectedly slow, check whether it's on `:memory:` when it could use
  a file-backed SQLite or a real Postgres/MySQL test database with
  transactions instead.
- Forgetting `Queue::fake()`/`Mail::fake()` doesn't just make a test slow —
  it can dispatch a real email or job during a test run against
  whatever queue/mail driver the test environment happens to have
  configured. Fake before triggering the side effect, not after.
- `assertDatabaseHas` after a mocked/faked repository call proves nothing —
  if the write path is faked, the row was never inserted; assert against
  the fake's own assertion methods (`Queue::assertPushed`, etc.) instead.
- Factory `state()` closures that reference `$this` or another factory's
  result at *definition* time (rather than inside the closure) capture a
  stale value shared across every invocation — keep dynamic data inside
  the closure body.

## Real-world grounding

Laravel's own `RefreshDatabase` trait exists because the older
`DatabaseMigrations` trait (still available) reruns every migration for
every test — for a schema with hundreds of migrations this turns a test
suite from seconds into minutes. The framework's documented recommendation
to default to `RefreshDatabase` and reserve `DatabaseMigrations` for
migration-specific tests is a direct response to that measured cost, not a
stylistic preference.

## Verification

- [ ] A failing test existed before the implementation change (RED)
- [ ] The correct database trait was chosen for the test's actual needs
- [ ] External side effects (queue, mail, notification, HTTP) are faked,
      not real
- [ ] Inertia responses are asserted via `assertInertia`, not raw JSON
- [ ] `php artisan test` (or `vendor/bin/pest`/`phpunit`) passes locally
