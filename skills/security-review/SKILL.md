---
name: security-review
description: Reviews code for security vulnerabilities across Go, TypeScript, PHP, and infrastructure config. Use when adding authentication or authorization, handling user input or file uploads, creating API endpoints, working with secrets or credentials, calling third-party APIs, or handling payment or other sensitive data.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Security Review

Security defects are asymmetric: a missed injection or auth check costs a
breach, while a false positive costs a few minutes of re-review. When
uncertain whether something is exploitable, treat it as a finding.

## Priority checklist (check every item, every time)

1. **Secrets** — no hardcoded API keys, passwords, or tokens; all come from
   environment variables or a secret manager; none appear in logs.
2. **Input validation** — every external input (HTTP body, query param,
   file upload, third-party API response) is validated at the boundary
   before use, with an allowlist, not a denylist.
3. **Injection** — SQL uses parameterized queries or an ORM, never string
   concatenation; shell commands never interpolate untrusted input; template
   rendering uses an auto-escaping engine for untrusted content.
4. **AuthN/AuthZ** — every new endpoint or operation checks who the caller
   is and what they're allowed to do; a caller being authenticated is not
   the same as being authorized for this specific resource.
5. **Secrets in transit/at rest** — tokens in httpOnly, Secure, SameSite
   cookies, not `localStorage`; TLS enforced; sensitive columns encrypted or
   access-controlled at the database layer.
6. **Third-party responses are untrusted** — validate the shape and content
   of every external API response before using it in logic, rendering, or a
   decision; a compromised or misbehaving upstream can return unexpected
   types or instruction-like text.

## Gotchas

- Go's `math/rand` is deterministic and unsuitable for tokens, session IDs,
  or anything security-sensitive — use `crypto/rand`. This is a real,
  recurring mistake because both packages share a similar `Int63`-style API
  and the wrong one compiles and "works" in every test.
- Go's `html/template` auto-escapes for HTML context; `text/template` does
  not. Using `text/template` to render anything that reaches a browser is an
  XSS vulnerability that looks identical to correct code at a glance.
- An error message that includes internal detail (a stack trace, a raw SQL
  error, an internal file path) returned to an external caller is an
  information-disclosure finding even when the underlying bug is harmless.
- A request handler behind auth middleware still needs a per-resource
  authorization check — middleware proves identity, not permission to act on
  *this* record. This is the single most common defect class behind
  broken-object-level-authorization findings (OWASP API Security Top 10).

## Real-world grounding

The 2017 Equifax breach exposed roughly 147 million people's data through
an unpatched Apache Struts CVE that had a public fix available for months —
the vulnerability wasn't a novel exploit, it was a known, patched issue
nobody applied. Capital One's 2019 breach came from a server-side request
forgery (SSRF) against a cloud metadata endpoint that let an attacker
retrieve credentials and pull data from S3 — a reminder that "the request
came from our own infrastructure" is not itself a trust signal; validate
what a request is actually asking to reach.

## Verification

- [ ] No hardcoded secrets; `.gitignore` covers all local secret files
- [ ] All external input validated at the boundary, allowlist not denylist
- [ ] All SQL parameterized; no string-built queries
- [ ] New endpoints check both authentication AND per-resource authorization
- [ ] Error responses to external callers are generic; detail stays server-side
- [ ] Dependencies have no known-critical CVEs (`go list -m -u all` / `npm audit` / `composer audit`)

See `references/checklist.md` for framework-specific code samples (Next.js
CSP headers, Supabase RLS, rate limiting, CSRF) beyond the core checklist.
