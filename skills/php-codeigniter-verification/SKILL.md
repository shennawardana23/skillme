---
name: php-codeigniter-verification
description: Use when asked to verify a CodeIgniter 4 project before a PR, after a major refactor or dependency upgrade, or before deploying to staging/production. Runs a phased env-check -> lint -> test -> security -> migration -> route/filter -> deploy-readiness pipeline via the spark CLI, stopping at the first failing gate rather than reporting every phase's output at once.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# PHP CodeIgniter Verification Loop

A gated pipeline, not a checklist to run in full regardless of outcome. Each
phase exists to catch a class of failure cheaply before a later, more
expensive phase wastes time on a codebase that was already broken. Stop at
the first failing gate, report it, and do not proceed to later phases.

## Phase 1: Environment and Composer

```bash
php -v
composer --version
php spark --version
composer validate
```

Confirm `.env` sets `CI_ENVIRONMENT` correctly for the target
(`production`, `testing`, `development`) and that `app/Config/App.php`'s
`$baseURL` is an explicit HTTPS URL for anything beyond local development
(see `php-codeigniter-security`). **Why first**: a wrong `CI_ENVIRONMENT`
changes which error/display settings and which `Config\Database` group
CI4 assumes — every later phase's results are unreliable if this is wrong.

## Phase 2: Lint and Static Analysis

```bash
vendor/bin/php-cs-fixer fix --dry-run --diff
vendor/bin/phpstan analyse
```

CI4's official coding-standard package is `CodeIgniter/coding-standard`
(a PHP-CS-Fixer ruleset), not Laravel Pint — don't assume Pint is present
just because a project is PHP; confirm which formatter this specific
project actually installed before running it.

## Phase 3: Tests

```bash
php spark test
```

See `php-codeigniter-tdd` for what a healthy CI4 test suite looks like —
in particular, confirm the `tests` database group (Phase 1's environment
check should have already caught a missing or misconfigured one) before
trusting this phase's results at all.

## Phase 4: Security and Dependency Audit

```bash
composer audit
```

Run after tests, not before — a vulnerable dependency doesn't block
verifying application behavior, but it does block release. Pair with a
manual pass over `php-codeigniter-security`'s checklist (CSRF mode,
`esc()` contexts, `$allowedFields`) since `composer audit` only flags known
CVEs in locked dependency versions, not this project's own code.

## Phase 5: Database and Migrations

```bash
php spark migrate --namespace App
php spark migrate:status
```

Confirm every migration has a working `down()` (CI4's migration base class
requires one, but an empty/no-op `down()` "works" without actually
reversing anything) and that destructive migrations (dropped columns,
dropped tables) have an explicit backup plan before running against a
target with real data.

## Phase 6: Routes and Filters

```bash
php spark routes
```

`spark routes` prints every registered route with its resolved filter
chain — use it to confirm a route group that's supposed to be
auth/tenant-gated actually shows the expected filter alias attached.
Per `php-codeigniter-patterns`' Gotchas, a `Filters.php` alias typo doesn't
throw; the route still resolves with the filter silently missing. This
phase is the cheapest way to catch that class of bug before it reaches a
feature test or, worse, production — grep the `spark routes` output for
the expected filter alias next to every route that should have it, don't
just skim for "does this look roughly right."

## Phase 7: Build and Deploy Readiness

```bash
php spark cache:clear
```

Confirm `writable/` (CI4's cache/logs/session storage directory — the CI4
equivalent of Laravel's `storage/`) is writable on the target host, and
that any environment-specific `Config\*` overrides (an
`app/Config/Production/App.php`-style override, if the project uses
environment-scoped config) resolve to the values expected for this
deploy target.

## Minimal flow (fastest useful gate)

```bash
php -v && composer --version && php spark --version
composer validate
vendor/bin/php-cs-fixer fix --dry-run --diff
vendor/bin/phpstan analyse
php spark test
composer audit
php spark migrate:status
php spark routes
```

## Gotchas

- **CI4 uses `spark`, not `artisan`** — a verification script copy-pasted
  from a Laravel project's CI config (`php artisan migrate`,
  `php artisan test`) will fail with a "command not found"-shaped error in
  a CI4 project, not a meaningful test failure; confirm the CLI entrypoint
  before assuming a copied pipeline works unchanged.
- **`php spark routes` output reflects `Config\Filters` exactly as
  configured** — including a typo'd alias that silently does nothing. The
  command will show the route as *not* having the intended filter attached;
  it will not warn that an alias doesn't exist, because from the router's
  perspective an unmatched or typo'd filter reference for a route group is
  simply absent, not an error.
- **`composer audit` says nothing about this project's own code** — a
  clean audit result is not a substitute for a manual pass over
  `php-codeigniter-security`'s checklist.
- **A missing or empty `down()` migration "passes" `migrate:status`** —
  status reporting only confirms which migrations have run, not that
  rollback actually works; a genuinely reversible migration needs to be
  verified by actually running `migrate:rollback` in a disposable
  environment, not inferred from status output alone.
- **`CI_ENVIRONMENT=production` disables CodeIgniter's detailed error
  display by default** — a verification run against a misconfigured
  `production` environment can mask a real Phase 1-level problem behind a
  generic error page instead of the actual exception; run Phase 1's
  environment check against the *actual* target environment value, not
  assume `development` locally is representative.

## Real-world grounding

CI4's `php spark routes` command exists specifically because filter wiring
in `Config\Filters.php` is declarative and string-keyed (aliases mapped to
route patterns) rather than attached directly at the route-definition call
site the way Laravel's `->middleware()` chain is — that indirection is
exactly what makes a typo'd alias invisible without a command that resolves
and prints the final filter chain per route. Treat "I read Config/Filters.php
and it looks right" as insufficient verification; the resolved output from
`spark routes` is the only thing that reflects what actually runs.

## Verification

- [ ] `CI_ENVIRONMENT` and `Config\App::$baseURL` match the actual target
      environment
- [ ] `vendor/bin/php-cs-fixer` (or this project's actual configured
      formatter) and `phpstan` both pass
- [ ] `php spark test` passes, against a confirmed-separate `tests`
      database group
- [ ] `composer audit` is clean, or flagged vulnerabilities are
      consciously accepted with a reason
- [ ] Every migration has a real, working `down()`
- [ ] `php spark routes` output shows the expected filter alias attached
      to every route group that requires auth/tenant-scoping
- [ ] `writable/` is confirmed writable on the target host
