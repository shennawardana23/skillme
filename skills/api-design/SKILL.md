---
name: api-design
description: Guides REST and RPC API design across Go, TypeScript, and PHP backends. Use when designing new endpoints, defining request/response contracts, planning breaking or additive schema changes, or reviewing an existing API for consistency in error handling, pagination, or naming.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# API and Interface Design

Design interfaces that are hard to misuse: consistent error shape, additive
changes only once a contract ships, and validation only at the boundary
where untrusted input enters.

## Hyrum's Law

> With enough users of an API, every observable behavior becomes a de facto
> contract — regardless of what the documentation promises.

Undocumented quirks, error text, timing, and field ordering all become
things someone depends on once observed. Design implication: don't expose
what you're not willing to keep stable, and plan for deprecation at design
time, not as an afterthought.

## Core rules

1. **Contract first.** Define the request/response types before
   implementing the handler.
2. **One error shape, everywhere.**
   ```json
   { "error": { "code": "VALIDATION_ERROR", "message": "email is required", "details": {} } }
   ```
   Never mix throwing, returning `null`, and returning `{error}` across
   different endpoints in the same API — the caller can't predict behavior.
3. **Validate at the boundary only.** HTTP handlers, form submissions, and
   third-party API responses are untrusted; internal function calls that
   already received validated types are not — re-validating there is
   wasted code that obscures where the real boundary is.
4. **Additive changes only, once shipped.** New optional fields are safe.
   Changing a field's type, removing a field, or changing status-code
   semantics for an existing endpoint is a breaking change — ship it as a
   new version or a new endpoint, never as a silent modification.
5. **Pagination on every list endpoint**, from the first version — adding
   it later is itself a breaking change to response shape.

## Gotchas

- A `PATCH` that silently ignores unknown fields is easy to build and easy
  to misuse — a client sending a typo'd field name gets no error and no
  effect, and won't notice until data looks wrong days later. Reject
  unknown fields at the boundary instead of ignoring them.
- REST resource URLs with verbs (`/api/createTask`) are a naming smell that
  correlates with inconsistent HTTP-method semantics elsewhere in the same
  API — grep for this pattern as a quick signal to review further.
- A third-party API response is untrusted data even when the third party is
  a well-known vendor: validate its shape before using it in any decision
  or rendering path.

## Real-world grounding

Twitter's 2018–2023 API changes — including the abrupt shutdown of the free
tier and repeated breaking changes to v1.1 without a viable migration
path — are a widely cited case study in what NOT to do: third-party
integrations built against observable (if undocumented) behavior broke with
no deprecation window, destroying an entire ecosystem of tools built on
the platform. Stripe's API, by contrast, is commonly cited for the opposite
pattern: it pins each integration to the API version active when that
integration's account was created via a version header, so Stripe can ship
breaking changes to new integrations while existing ones keep working
indefinitely on their pinned version.

## Verification

- [ ] Every endpoint has typed request/response schemas
- [ ] Every error response follows one consistent shape
- [ ] Validation happens at the boundary, not scattered through internal calls
- [ ] List endpoints are paginated from day one
- [ ] New fields are additive and optional; no in-place breaking changes
- [ ] Naming is consistent (plural nouns, camelCase fields, `is`/`has` booleans)

See `references/rest-and-typescript-patterns.md` for resource design,
discriminated unions, and branded-ID patterns.
