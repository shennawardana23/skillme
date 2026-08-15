---
name: jwt-tenant-scoped-authorization
description: Guides multi-tenant authorization where a JWT claim carries a tenant/organization identifier that every endpoint must check against the resource being accessed. Use when adding a new endpoint to a multi-tenant API, reviewing JWT-based auth code, deciding what a JWT claim should carry, or investigating a cross-tenant data access report.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# JWT Tenant-Scoped Authorization

In a multi-tenant API, "the request has a valid JWT" and "the request is
allowed to access this specific resource" are two different checks — the
first proves who's asking, the second proves they're allowed to touch
*this* row. A service that only does the first (a global auth middleware
that validates the token) and relies on every individual endpoint to also
do the second is one missed endpoint away from a cross-tenant data leak,
and that gap doesn't show up in a request that only ever tests your own
tenant's data.

## The check that must repeat, and where it belongs

A middleware validates the JWT itself (signature, expiry) once, globally.
The tenant-scoping check — does the resource identifier in *this* request
belong to the tenant in the token's claim — cannot live only in
middleware, because it depends on the specific resource each endpoint
handles:

```go
func (h *RoomHandler) Get(c *gin.Context) {
    claim := auth.ClaimFromContext(c)
    requestedTenantID := c.Param("tenantId")

    if requestedTenantID != claim.TenantID {
        c.JSON(http.StatusForbidden, ErrorResponse{"tenant mismatch"})
        return
    }
    // ... proceed only after the explicit per-endpoint check
}
```

Treat this per-endpoint check as mandatory boilerplate, not optional
defense-in-depth — for a resource-scoped multi-tenant API, it is the only
thing preventing tenant A's valid token from reading or writing tenant B's
data, independent of whatever authorization the data layer itself does.

## Claim design

Keep a JWT claim minimal and purpose-built: a tenant identifier plus
standard claims (issuer, expiry, subject) is usually enough for
tenant-scoping; avoid embedding data that changes independently of the
token's lifetime (a role that gets revoked mid-session, a permission list
that grows) unless the service has a real token-revocation or short-lived-
token strategy to keep that embedded data from going stale before the
token expires.

## Gotchas

- **A global "is this JWT valid" middleware does not imply per-resource
  authorization** — the tenant-scoping check has to be repeated at each
  endpoint that accepts a resource/tenant identifier from the request,
  because only that endpoint knows which resource is actually being
  touched. Treat a new endpoint that's missing this check as a real
  cross-tenant vulnerability (IDOR), not a style nit.
- **Do not use an archived/unmaintained JWT library for new work.** Some
  early, widely-copied Go JWT libraries are no longer maintained; verify
  a JWT dependency is actively maintained before adding new code against
  it, and treat "upgrade this JWT library" as a fork/migration decision
  (different import path, possibly different API), not a routine version
  bump.
- **Comparing a static API key or secret with `!=`/`==` is not
  constant-time** and leaks timing information about how much of the
  secret matched — use a constant-time comparison
  (`subtle.ConstantTimeCompare`in Go's standard library, or the
  language's equivalent) for any secret/token comparison, not a plain
  equality check.
- **`Access-Control-Allow-Origin: *` combined with
  `Access-Control-Allow-Credentials: true` is invalid and rejected by
  browsers for credentialed requests** — if cross-origin requests with
  cookies/credentials stop working, check for this exact combination
  before assuming the issue is elsewhere; a wildcard origin cannot be
  paired with a credentials flag.
- **Returning HTTP 200 with the actual error in the response body for
  auth failures** (rather than a proper 401/403 status) is sometimes a
  deliberate choice to avoid tripping infrastructure-level alerting on
  4xx rates — if you find this pattern, confirm whether it's intentional
  before "fixing" it to return a standard status code, since doing so
  can silently break an ops team's existing alerting assumptions.
- **A manually-parsed `Authorization: Bearer <token>` header with no
  null-safety** (e.g. blindly indexing the split result) throws on any
  malformed or missing header rather than returning a clean 401 — guard
  the split/parse explicitly rather than assuming the header is always
  well-formed.

## Real-world grounding

This per-endpoint tenant-check gap is the textbook shape of an IDOR
(Insecure Direct Object Reference) vulnerability — OWASP's own
documentation of this class of bug describes exactly this pattern: an
application correctly authenticates the caller but fails to verify that
the specific object being requested belongs to that caller, letting an
authenticated user substitute a different identifier and access another
tenant's data. It is one of the most commonly reported vulnerability
classes in multi-tenant SaaS applications precisely because authentication
correctness gives a false sense that authorization is also handled.

## Verification

- [ ] Every endpoint that accepts a tenant/resource identifier from the
      request explicitly checks it against the caller's JWT claim, not
      just relying on global auth middleware
- [ ] The JWT library in use is actively maintained, not an archived fork
- [ ] Secret/token comparisons use a constant-time comparison, not `==`/`!=`
- [ ] CORS configuration never pairs a wildcard origin with
      credentials-allowed
- [ ] Any non-standard HTTP status behavior on auth failure (e.g. 200
      with an error body) is a documented, deliberate choice, not an
      oversight
- [ ] Manual header parsing (Authorization, API keys) guards against a
      missing or malformed value rather than assuming well-formed input
