---
name: laravel-verification
description: Use when asked to verify a Laravel project before a PR, after a major refactor or dependency upgrade, or before deploying to staging/production. Runs a phased env-check -> lint -> test -> security -> migration -> deploy-readiness -> queue pipeline, stopping at the first failing gate rather than reporting every phase's output at once.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Laravel Verification Loop

A gated pipeline, not a checklist to run in full regardless of outcome. Each
phase exists to catch a class of failure cheaply before a later, more
expensive phase wastes time on a codebase that was already broken. Stop at
the first failing gate, report it, and do not proceed to later phases —
running `php artisan test` against code that fails `composer validate` only
produces confusing secondary failures.

## Phase 1: Environment and Composer

```bash
php -v
composer --version
php artisan --version
composer validate
composer dump-autoload -o
```

Confirm `.env` exists, required keys are present, `APP_ENV` matches the
target (`production`, `staging`), and `APP_DEBUG=false` outside local/dev.
Using Sail: prefix commands with `./vendor/bin/sail`.

**Why first**: a broken autoloader or missing env key makes every later
phase fail for an unrelated reason. Stop here on any failure.

## Phase 2: Lint and Static Analysis

```bash
vendor/bin/pint --test
vendor/bin/phpstan analyse   # or: vendor/bin/psalm
```

**Why before tests**: static analysis catches type and null-safety errors
tests may not exercise, and it's near-instant compared to the full suite —
fail fast on the cheap check.

## Phase 3: Tests and Coverage

```bash
php artisan test
XDEBUG_MODE=coverage php artisan test --coverage   # CI only, see Gotchas
```

## Phase 4: Security and Dependency Audit

```bash
composer audit
```

Run after tests, not before: a vulnerable dependency doesn't block
verifying application behavior, but it does block release.

## Phase 5: Database and Migrations

```bash
php artisan migrate --pretend
php artisan migrate:status
```

Check migration filenames follow `Y_m_d_His_*`, that every migration with an
`up()` has a working `down()`, and that destructive migrations (dropped
columns, dropped tables) have an explicit backup plan — `--pretend` will not
catch a missing or broken rollback path (see Gotchas).

## Phase 6: Build and Deploy Readiness

```bash
php artisan optimize:clear
php artisan config:cache
php artisan route:cache
php artisan view:cache
```

Confirm these complete without error in production-equivalent config, and
that `storage/` and `bootstrap/cache/` are writable on the target host.
Run `config:cache` only after Phase 1's env check passes — see Gotchas for
why caching a bad config is worse than no cache.

## Phase 7: Queue and Scheduler

```bash
php artisan schedule:list
php artisan queue:failed
php artisan horizon:status       # if Horizon is installed
php artisan queue:monitor default --max=100   # backlog check, no processing
```

Staging only — dispatch a no-op job to a dedicated queue and confirm it
produces its expected side effect (log line, healthcheck row, metric):

```bash
php artisan tinker --execute="dispatch((new App\Jobs\QueueHealthcheck())->onQueue('healthcheck'))"
php artisan queue:work --once --queue=healthcheck
```

Never run the active dispatch step against a `production` queue connection.

## Minimal flow (fastest useful gate)

```bash
php -v && composer --version && php artisan --version
composer validate
vendor/bin/pint --test
vendor/bin/phpstan analyse
php artisan test
composer audit
php artisan migrate --pretend
php artisan config:cache
php artisan queue:failed
```

## Gotchas

- `php artisan config:cache` bakes every `config/*.php` value into a single
  cached file. Any code that calls `env()` **outside** a `config/` file
  (directly in a controller, service, or view) will silently read `null` in
  production once the cache is built, even though `.env` still has the
  value — the classic "works locally, breaks in prod" Laravel footgun. Grep
  for `env(` outside `config/` before caching.
- `php artisan migrate --pretend` prints the SQL Eloquent's schema builder
  *would* run — it does not execute and therefore does not validate raw
  `DB::statement(...)` calls inside a migration, since those aren't part of
  the schema builder's generated SQL. A migration that only contains raw
  statements will show as empty/no-op under `--pretend` and give false
  confidence.
- `XDEBUG_MODE=coverage` instruments every line and can slow the suite by
  several multiples — keep it CI/nightly-only, not part of the local
  inner-loop `php artisan test` run.
- `composer audit` only flags known CVEs in `composer.lock` — it says
  nothing about code you wrote, so a clean audit is not a substitute for
  Phase 2's static analysis or Phase 4's own review.
- A green Phase 3 with a red Phase 1 is impossible to trust: if `.env` was
  wrong, tests may have silently run against the wrong database or with
  debug mode masking errors.

## Real-world grounding

Laravel Pint (Phase 2) is Laravel's first-party wrapper around
[PHP-CS-Fixer](https://github.com/PHP-CS-Fixer/PHP-CS-Fixer), configured with
Laravel's own opinionated ruleset — that's why `pint --test` is the standard
format gate in Laravel CI pipelines rather than a bare `php-cs-fixer` call.
