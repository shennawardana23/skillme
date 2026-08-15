---
name: php-codeigniter-legacy-auth-session
description: Guides authentication and session security for classic CodeIgniter 2.x/3.x applications - session/cookie design, CSRF, and object-level authorization, where the framework's built-in mechanisms are weaker or differently-shaped than CodeIgniter 4's. Use when reviewing login/session code, a cookie-setting call, or an endpoint that accepts an identifier (hotel_id, tenant_id, resource_id) from client input in a classic CodeIgniter app.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "php"
---

# Classic CodeIgniter Auth and Session Security

Classic CodeIgniter (2.x/3.x) apps often build their own authentication
and session layer on top of (or instead of) the framework's built-in
session library, and default to CSRF protection **off** rather than CI4's
default-on posture. Both facts mean auth/session security in a codebase
like this needs deliberate review rather than trusting framework defaults
the way a CI4 app's `php-codeigniter-security` skill would.

## Menu-level authorization is not object-level authorization

A common pattern: a controller checks whether the current user has
privilege to access a given *menu/feature* (a role check, a permission
lookup), then proceeds to act on a specific resource identifier taken
directly from request input, with no check that the identifier belongs to
something the user is actually allowed to touch:

```php
public function update_template()
{
    if (!$this->check_privilege('cms_template_edit')) {
        return $this->deny_access();
    }
    $hotel_id = $this->input->post('hotel_id'); // trusted with no ownership check
    $this->app_model->update_data($hotel_id, $this->input->post(), 'templates');
}
```

"Can this user use this feature" and "does this specific record belong to
something this user is allowed to modify" are two different checks — a
codebase that only does the first is vulnerable to an authenticated user
substituting a different resource identifier (a different hotel, tenant,
or account ID) and modifying data that isn't theirs. This is the same
class of bug `jwt-tenant-scoped-authorization` describes for token-based
APIs, showing up here in a session-based CMS instead.

## Session tokens must not be computable from user data

A session or auth cookie's value must be unpredictable given information
an attacker could know or guess (a user ID, a timestamp) — a session key
derived deterministically from the user's own ID (e.g. a hash of a fixed
prefix plus the user ID, with no random component) means anyone who knows
or can enumerate a user ID can compute that user's valid session key
without ever authenticating. A session token needs a genuinely random
component, not just any hash function applied to already-known data.

## Cookie hardening needs to be applied consistently

If a codebase has a "secure cookie" helper that sets `secure`/`httponly`/
`samesite` flags, confirm it's actually used for the primary auth cookie
itself, not just for secondary cookies — it's a common gap for the main
session/auth cookie to be set via a plain, unhardened `setcookie()` call
while a wrapper function with the correct flags exists and is used
elsewhere in the same codebase.

## Two disconnected "sessions" can coexist

A classic CodeIgniter app can have the framework's own native session
library configured and autoloaded (file- or database-backed) while the
actual authenticated-user state lives entirely in a separate store (Redis,
a custom table) keyed by a custom cookie — meaning `$this->session` and
"whether the user is actually logged in" are two unrelated systems. A
developer told to "check the session" needs to know which of the two
mechanisms actually governs auth state in a given codebase before
debugging in the wrong place.

## Gotchas

- **Menu/feature-level authorization checks do not imply object-level
  authorization** — a resource identifier taken from client input
  (`hotel_id`, `tenant_id`, any ownership-determining ID) needs its own
  explicit check against what the current user is allowed to touch,
  independent of whatever role/privilege check already ran.
- **A session key derived deterministically from a user ID (no random
  component) is computable by anyone who knows that ID** — treat this as
  a critical finding, not a style issue; a session token needs real
  randomness, not just any transformation of already-public data.
- **A "secure cookie" helper existing in a codebase doesn't mean it's
  used everywhere it should be** — check whether the primary auth cookie
  specifically goes through it, since it's a common gap for exactly the
  most important cookie to bypass the hardening wrapper that exists for
  others.
- **CSRF protection defaults to off in classic CodeIgniter** (unlike CI4's
  default-on posture) — check `$config['csrf_protection']` explicitly
  rather than assuming any CSRF protection exists unless it's been turned
  on deliberately.
- **A commented-out authentication check with a TODO is a real,
  currently-unauthenticated endpoint** — treat "auth check commented out
  pending a fix" as a live vulnerability to escalate immediately, not a
  known-and-accepted trade-off to leave alone because someone left a note
  about it.
- **Environment detection by literal hostname string match** (comparing
  `$_SERVER['SERVER_NAME']` against a hardcoded production domain string)
  is fragile — it silently falls through to whatever the "else" branch
  assumes (often looser local-dev defaults) for any hostname that doesn't
  exactly match, including a staging domain, a new environment, or a
  typo in the comparison string itself.
- **An in-process cache (e.g. APCu) fronting a remote secrets/config
  fetch can fail open to "call the remote service on every request"**
  with no visible error if the cache extension is disabled or unavailable
  on a given worker — this is a silent availability/performance cliff,
  not a hard failure that would get noticed quickly.

## Real-world grounding

The menu-level-vs-object-level authorization gap described here is the
same IDOR (Insecure Direct Object Reference) pattern OWASP documents for
API-based systems, appearing in a session-based CMS context instead of a
token-based API — the underlying mistake (authenticating and authorizing
*that a feature can be used* without separately authorizing *this specific
record*) is identical regardless of whether the surrounding system is a
modern JWT-based API or a classic server-rendered CMS. Predictable,
non-random session tokens are a long-documented session-management
weakness (OWASP's session management guidance is explicit that a session
identifier must be generated with a cryptographically strong random
source, not derived from user-known data) precisely because a
deterministic derivation collapses "have a valid session" down to "know
the input the derivation function used."

## Verification

- [ ] Every endpoint accepting a resource/ownership identifier from
      client input checks it against what the current user is actually
      allowed to touch, not just a feature-level privilege check
- [ ] Session/auth tokens have a genuinely random component, not a
      deterministic derivation from user-known data
- [ ] The primary auth cookie uses the same hardening (secure, httponly,
      samesite) as any other cookie a "secure cookie" helper protects
- [ ] `csrf_protection` is explicitly reviewed, not assumed on by default
- [ ] No authentication check is commented out in code reachable in
      production
- [ ] Environment detection does not rely solely on a literal hostname
      string match with an unreviewed fallback branch
