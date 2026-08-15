---
name: laravel-security
description: Use when configuring Laravel-specific auth/session/CSRF/authorization plumbing - Sanctum SPA authentication, session cookie hardening, policy/gate wiring, mass-assignment guards, signed URLs, or encrypted casts. For generic web-security review (injection, secrets, auth-vs-authz on any stack) use security-review instead; this skill is for getting Laravel's own framework mechanisms configured correctly.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Laravel Security Configuration

This skill is about *Laravel's own framework mechanisms* — the config keys,
middleware, and Eloquent guards that are specific to this framework and easy
to misconfigure even when the surrounding logic is otherwise secure. It is
not a general injection/XSS/secrets checklist — see `security-review` for
that; apply both when reviewing a Laravel app.

## Session and cookie hardening

Laravel's session config (`config/session.php`) drives cookie behavior for
every stateful request, including Sanctum's SPA mode.

```env
SESSION_SECURE_COOKIE=true      # cookie only sent over HTTPS
SESSION_HTTP_ONLY=true          # not readable from JavaScript (default true)
SESSION_SAME_SITE=lax           # use 'strict' for high-risk flows (banking, admin)
```

Regenerate the session on every privilege change, not just login — Laravel
does this automatically on `Auth::login()` via `session()->regenerate()`
inside the default `AuthenticatesUsers` flow, but a custom login controller
that calls `Auth::login($user)` directly without going through the framework
trait must call `$request->session()->regenerate()` itself, or session
fixation remains possible across the login boundary.

## Sanctum SPA authentication

Sanctum's SPA mode (cookie-based, not token-based) requires the frontend's
domain to be explicitly listed as *stateful* — anything not listed falls
back to token auth and won't carry CSRF/session cookies correctly:

```php
// config/sanctum.php
'stateful' => explode(',', env(
    'SANCTUM_STATEFUL_DOMAINS',
    'localhost,localhost:3000,127.0.0.1,127.0.0.1:8000,::1'
)),
```

The SPA and the API must share the same top-level domain (or be configured
with `SESSION_DOMAIN` covering both) — Sanctum's cookie approach cannot
authenticate a genuinely cross-site SPA; a fully separate domain needs
token-based Sanctum instead, not the SPA flow.

Client must hit `/sanctum/csrf-cookie` before the first stateful request:

```php
Route::middleware('auth:sanctum')->get('/me', fn (Request $r) => $r->user());
```

## Policies and gates: route- and method-level enforcement

Middleware answers "who is this," policies answer "can they do this to
*this record*." Wire both — a route behind `auth:sanctum` alone has no
per-resource check.

```php
// Controller-level, inside the action:
$this->authorize('update', $project);

// Route-level, before the action runs at all:
Route::put('/projects/{project}', [ProjectController::class, 'update'])
    ->middleware(['auth:sanctum', 'can:update,project']);
```

Register the policy in `AuthServiceProvider` (or let auto-discovery match
`ProjectPolicy` to `Project` by naming convention) — a policy class that
exists but isn't registered and doesn't match the naming convention silently
never runs, and `$this->authorize()` will throw `AuthorizationException`
only if the ability itself is undefined, not if the policy method logic is
wrong.

## Mass assignment guards

```php
final class Project extends Model
{
    protected $fillable = ['name', 'description', 'owner_id'];
    // never: protected $guarded = [];
}
```

`$guarded = []` combined with `Model::create($request->all())` accepts every
field the client sends, including ones that were never meant to be
client-settable (`is_admin`, `hotel_id`, `role`). Prefer `$fillable` as an
explicit allowlist over `$guarded` as a denylist — a denylist must be kept
in sync with every new column added later, an allowlist fails safe by
default. For anything touching authorization-relevant columns, bypass mass
assignment entirely and set fields individually after validation.

## Signed URLs for tamper-proof temporary links

Use signed routes instead of exposing a resource ID and trusting nobody
guesses the neighbor:

```php
$url = URL::temporarySignedRoute(
    'downloads.invoice', now()->addMinutes(15), ['invoice' => $invoice->id]
);

Route::get('/invoices/{invoice}/download', [InvoiceController::class, 'download'])
    ->name('downloads.invoice')
    ->middleware('signed');
```

The `signed` middleware validates the signature and expiry but does **not**
check authorization — a signed URL proves the link wasn't tampered with, not
that the current requester is allowed to view *this* invoice. Still call
`$this->authorize('view', $invoice)` inside the controller action.

## Encrypted attributes at rest

```php
protected $casts = ['api_token' => 'encrypted'];
```

Encryption uses `APP_KEY` — rotating `APP_KEY` without a migration step
makes every previously encrypted column unreadable. Never rotate `APP_KEY`
in place without first decrypting-and-re-encrypting existing rows under the
new key, or write access continues but existing encrypted reads throw
`DecryptException`.

## CORS and Sanctum credentials

```php
// config/cors.php
'allowed_origins' => ['https://app.example.com'],  // never '*' with credentials
'supports_credentials' => true,
```

`supports_credentials => true` with a wildcard origin is rejected by
browsers outright (the Fetch spec forbids `*` alongside credentialed
requests) — if cross-origin cookies stop working after adding Sanctum, check
this first before suspecting the session config.

## Gotchas

- `php artisan config:cache` in production means changes to `.env` after
  caching (including `SESSION_SECURE_COOKIE`, `SANCTUM_STATEFUL_DOMAINS`)
  have no effect until the cache is rebuilt — a security setting can look
  "changed" in `.env` while the running app still enforces the old value.
- A custom login flow that skips Laravel's `AuthenticatesUsers` trait and
  calls `Auth::login($user)` directly must manually regenerate the session
  — this is the most common way Laravel session-fixation protection gets
  silently dropped.
- `Route::middleware('can:update,project')` resolves `project` from the
  route-model-bound parameter *before* the policy runs; if the model isn't
  found the framework 404s before authorization is even evaluated — don't
  mistake a 404 for "authorization is working."
- `'guarded' => []` is sometimes added "temporarily" to unblock a form and
  never removed — grep for it in review; it silently reopens mass
  assignment on every attribute added afterward, including ones nobody
  intended to expose.

## Real-world grounding

CVE-2018-15133 was a deserialization RCE in Laravel's
`Illuminate\Broadcasting\PendingBroadcast` reachable when an attacker knew
(or brute-forced) the application's `APP_KEY` — it's the canonical reason
`APP_KEY` is treated as a secret on par with a database password, not a
throwaway config default, and why leaking `.env` (e.g. via a misconfigured
debug page) is a full-application compromise, not just an info leak.
