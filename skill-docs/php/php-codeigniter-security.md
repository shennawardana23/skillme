## What it does

Configures CodeIgniter 4's own security-relevant framework mechanisms —
CSRF mode, output escaping, HTTPS enforcement, and mass-assignment guards.
The defining constraint: CI4 ships a *working* default for each of these,
but "working" is not the same as "correct for this app" — the framework's
own docs explicitly call out that its default CSRF mode has a real gap for
apps that also use sessions, which is most apps.

## When to reach for it

Reach for this skill when reviewing or configuring `Config\Security`,
`Config\App`, or a Model's `$allowedFields` in a CodeIgniter 4 project, or
when a form/endpoint change needs a security review before merge. For
general injection/XSS/secrets review that applies regardless of framework,
use `security-review` instead — apply both to a CI4 app; this skill is
narrower and specific to CI4's own configuration surface.

## The default that needs a second look

CI4's CSRF protection defaults to `'cookie'` mode (Double Submit Cookie) —
it works without server-side session state, which is a reasonable default
for a framework that doesn't assume every app uses sessions. But the CI4
user guide states directly that cookie-based CSRF does not protect against
same-site attacks the way session-based (Synchronizer Token) protection
does. An app that already maintains sessions for authentication gains
nothing from the session-less default and loses that protection — this is
the single highest-value check this skill exists to prompt.

## Common questions

- **"Our app uses session-based login. Is the default `'cookie'` CSRF mode
  fine, or should we change it?"** Change it. The CI4 docs are explicit:
  cookie-based CSRF doesn't stop same-site attacks; session-based does. If
  the app already has server-side sessions for auth, there's no reason to
  keep the session-less default.
- **"We call `esc($value)` everywhere — is that enough?"** Only if every
  call's context matches where the value is actually written. `esc()`'s
  default context is `html`; a value interpolated into an inline `<script>`
  block or an unquoted attribute needs `esc($value, 'js')` or
  `esc($value, 'attr')` respectively — the default context does not
  neutralize every context's injection vector.
- **"We set `forceGlobalSecureRequests = true` — are we covered for
  HTTPS?"** Only if `Config\App::$baseURL` is also an explicit HTTPS URL.
  Leaving `$baseURL` empty falls back to auto-detecting it from request
  headers, which the framework doesn't recommend for production — the two
  settings work together, not independently.
- **"Is `$allowedFields` actually a security control, or just a data-shape
  thing?"** Both, but treat it as security first. An empty or missing
  `$allowedFields` on a Model that persists request data lets any
  client-supplied field reach any column — it's the same mechanism and the
  same risk class as Laravel's `$fillable`/`$guarded`, just under a
  different name.

## It's working if

- `Config\Security::$csrfProtection` is `'session'` for any app that also
  uses sessions for authentication
- Every `esc()` call's context argument matches its actual output location
- `Config\App::$baseURL` is an explicit HTTPS URL in production
- Every Model reachable from a request path declares a non-empty
  `$allowedFields`

## Where it fits

Pairs with `php-codeigniter-patterns` (which covers `$allowedFields` as an
architecture concern first) and `php-codeigniter-verification` (which
checks this skill's items as part of a pre-deploy gate). For the
non-CI4-specific security review — injection, secrets, generic auth/authz
— apply `security-review` alongside this one, not instead of it.
