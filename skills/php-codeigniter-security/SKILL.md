---
name: php-codeigniter-security
description: Use when configuring CodeIgniter 4-specific security mechanisms - CSRF protection mode (cookie vs session), output escaping with esc(), Config\App's baseURL and forceGlobalSecureRequests, or Model $allowedFields mass-assignment guards. For generic web-security review (injection, secrets, auth-vs-authz on any stack) use security-review instead; this skill is for getting CI4's own framework mechanisms configured correctly.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "php"
---

# PHP CodeIgniter Security Configuration

This skill is about *CodeIgniter 4's own framework mechanisms* — the config
properties and helpers specific to this framework that are easy to
misconfigure even when the surrounding application logic is otherwise sound.
It is not a general injection/XSS/secrets checklist — see `security-review`
for that; apply both when reviewing a CI4 app.

## CSRF: cookie vs. session protection

`Config\Security::$csrfProtection` selects one of two real, different
mechanisms:

```php
// app/Config/Security.php
public string $csrfProtection = 'session'; // or 'cookie' (the framework default)
```

- **`'cookie'`** — Double Submit Cookie pattern: a token is stored in a
  cookie and compared against a submitted form/header value. Works without
  server-side session state.
- **`'session'`** — Synchronizer Token Pattern: the token is stored
  server-side in the session and compared against the submitted value.

The CI4 user guide states this directly: **cookie-based CSRF protection does
not prevent same-site attacks** the way session-based protection does — if
the app already uses sessions for auth, use `'session'` for CSRF too, not the
framework's cookie default, which was designed for session-less use cases.

```php
// app/Config/Filters.php
public array $globals = [
    'before' => ['csrf' => ['except' => ['api/webhooks/*']]],
];
```

## Output escaping: `esc()`

```php
<?= esc($comment->body) ?>                    <!-- HTML context (default) -->
<script>const id = <?= esc($id, 'js') ?>;</script>
<a href="<?= esc($url, 'url') ?>">link</a>
<div data-value="<?= esc($value, 'attr') ?>">
```

`esc()`'s second argument selects the escaping context (`html`, `js`, `css`,
`url`, `attr`) — using the default `html` context for a value interpolated
into a `<script>` block or an `href` does not correctly neutralize that
context's injection vectors; the context argument must match where the
value is actually being written, not just "some HTML page."

## Mass assignment: `$allowedFields`

Covered in depth in `php-codeigniter-patterns` (it's an architecture
concern first), but it is exactly as much a security control as Laravel's
`$fillable` — an empty or missing `$allowedFields` on a Model that persists
request-derived data via `insert()`/`update()`/`save()` lets a client-supplied
field silently reach any column, including ones like `is_admin`, `hotel_id`,
or `role` that should never be client-settable. Treat a Model with no
`$allowedFields` declared, reachable from any request-handling path, as a
review-blocking finding.

## HTTPS enforcement and base URL

```php
// app/Config/App.php
public string $baseURL = 'https://app.example.com/';
public bool $forceGlobalSecureRequests = true;
```

`forceGlobalSecureRequests = true` redirects any HTTP request to HTTPS and
sets the `Strict-Transport-Security` header — but it has no effect if
`$baseURL` itself is left as `http://` or empty (CI4 falls back to
auto-detecting the base URL from request headers when `$baseURL` is empty,
which is explicitly discouraged in production since it trusts
client-supplied `Host`/`X-Forwarded-Host` headers for URL generation).
Setting an explicit HTTPS `$baseURL` is a prerequisite for
`forceGlobalSecureRequests` to be meaningful, not an independent setting.

## Session cookie hardening

```php
// app/Config/Session.php
public bool $cookieSecure = true;    // HTTPS-only cookie
public bool $cookieHTTPOnly = true;  // not readable from JavaScript (default)
public string $cookieSameSite = 'Lax'; // 'Strict' for high-risk flows
```

## Gotchas

- **Cookie-based CSRF does not stop same-site attacks** — the CI4 docs say
  this explicitly. An app using sessions for authentication should use
  `'session'` CSRF protection, not leave the framework's `'cookie'` default
  in place unreviewed.
- **`esc()`'s context argument must match the actual output location.**
  Escaping a value for `html` and then interpolating it into an inline
  `<script>` block or an unquoted HTML attribute does not neutralize
  that context's injection vector — the default context is not "safe
  everywhere."
- **`forceGlobalSecureRequests` without an explicit HTTPS `$baseURL`
  accomplishes less than it appears to** — confirm `$baseURL` is set
  explicitly in production config, not left to auto-detection.
- **An empty or absent `$allowedFields` is a silent mass-assignment hole**,
  identical in effect to Laravel's `$guarded = []` — see
  `php-codeigniter-patterns` for the full mechanism; flag it here as a
  security finding, not just an architecture nit.
- **CSRF token regeneration on every submission (`$regenerate = true`,
  the CI4 default) can break back/forward navigation, multiple open tabs,
  or concurrent AJAX requests** if the frontend doesn't refresh its stored
  token after each request — this is a usability/security tradeoff to
  make deliberately, not a bug to "fix" by disabling regeneration outright.

## Real-world grounding

CSRF's cookie-vs-session distinction in CI4 mirrors the broader, well-known
weakness of double-submit-cookie CSRF defenses generally: they protect
against a pure cross-site forged request but not against same-site
subdomain-level attacks, since a cookie is often readable/writable by
anything on the same registrable domain. This is exactly why session-based
(synchronizer token) CSRF protection is the stronger recommendation whenever
an app already maintains server-side session state for auth — it is not a
CI4-specific quirk, but CI4's documentation calling it out explicitly (rather
than leaving developers to discover it) is what makes the framework default
worth double-checking on every app rather than assuming it's already correct.

## Verification

- [ ] `Config\Security::$csrfProtection` is `'session'` for any app that
      also uses sessions for authentication
- [ ] Every `esc()` call's context argument matches where the value is
      actually rendered (`js`, `url`, `attr`, not a blanket default)
- [ ] Every Model reachable from request-handling code declares a non-empty
      `$allowedFields`
- [ ] `Config\App::$baseURL` is an explicit HTTPS URL in production, not
      left empty for auto-detection
- [ ] Session cookies set `cookieSecure = true` and an appropriate
      `cookieSameSite` value for the app's risk profile
