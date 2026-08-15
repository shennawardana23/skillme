---
name: security-and-hardening
description: Hardens a feature against attack before or while it's built — threat-model trust boundaries with STRIDE, apply the Always/Ask-First/Never tiers, and prevent OWASP Top 10 classes by construction. Use when building anything that accepts user input, handles auth, sessions, file uploads, webhooks, payments, PII, or calls an LLM — as opposed to auditing code that already exists.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Security and Hardening

This is a build-time skill: it hardens a feature while it's being designed and written. For auditing code that already exists, use `security-review` (priority checklist) or `security-analysis` (pattern scan). Treat every external input as hostile, every secret as sacred, every authorization check as mandatory — security is a constraint on every line touching user data, auth, or external systems, not a phase at the end.

## Process: threat model first

Controls bolted on without a threat model are guesses. Before writing hardening code, spend five minutes thinking like an attacker:

1. **Map the trust boundaries** — where does untrusted data cross into the system? HTTP requests, form fields, file uploads, webhooks, third-party API responses, message queues, and **LLM output**. Every boundary is attack surface.
2. **Name the assets** — what's worth stealing or breaking? Credentials, PII, payment data, admin actions, money movement.
3. **Run STRIDE over each boundary** — a quick lens, not a ceremony:

| Threat | Ask | Typical mitigation |
|---|---|---|
| Spoofing | Can someone impersonate a user/service? | Authentication, signature verification |
| Tampering | Can data be altered in transit or at rest? | Integrity checks, parameterized queries, HTTPS |
| Repudiation | Can an action be denied later? | Audit logging of security events |
| Information disclosure | Can data leak? | Encryption, field allowlists, generic errors |
| Denial of service | Can it be overwhelmed? | Rate limiting, input size caps, timeouts |
| Elevation of privilege | Can a user gain rights they shouldn't? | Authorization checks, least privilege |

4. **Write abuse cases next to use cases** — for each feature, ask "how would I misuse this?" and make that the first test.

If you can't name the trust boundaries for a feature, you're not ready to secure it — this is OWASP A04: Insecure Design; most breaches begin in design, not code.

## The three-tier boundary system

**Always do (no exceptions):** validate all external input at the boundary; parameterize every database query; encode output via framework auto-escaping; use HTTPS everywhere; hash passwords with bcrypt/scrypt/argon2; set security headers (CSP, HSTS, X-Frame-Options); use httpOnly/secure/sameSite cookies for sessions; run a dependency audit before every release.

**Ask first (needs human approval):** new authentication flows or auth-logic changes; storing a new category of sensitive data (PII, payment info); new external service integrations; CORS changes; file upload handlers; rate-limit/throttle changes; granting elevated permissions.

**Never do:** commit secrets to version control; log sensitive data (passwords, tokens, full card numbers); trust client-side validation as a security boundary; disable security headers for convenience; use `eval()`/`innerHTML` with user data; store session tokens in client-accessible storage; expose stack traces to users.

## OWASP Top 10 classes, by construction

Prevent these classes rather than patching them later. See `references/security-patterns.md` for full code samples (parameterized queries, bcrypt + session config, CSP/helmet, per-resource authorization, SSRF host+IP pinning, zod schema validation, file-upload checks, rate limiting).

- **Injection** — parameterize; never concatenate untrusted input into SQL, shell, or a non-escaping template.
- **Broken authentication** — hash with bcrypt (cost ≥ 12); httpOnly/secure/sameSite session cookies with an expiry.
- **XSS** — rely on framework auto-escaping; sanitize with a library (e.g. DOMPurify) only when raw HTML rendering is unavoidable.
- **Broken access control** — check resource ownership on every mutating endpoint, not just that the caller is authenticated.
- **Security misconfiguration** — security headers and a restrictive CSP/CORS allowlist on by default, not opt-in.
- **Sensitive data exposure** — strip secret/internal fields before serializing a record for an API response.
- **SSRF** — any server-side fetch of a user-influenced URL needs a host allowlist, a check that *every* resolved IP is a public unicast address, and `redirect: 'error'`. The `169.254.169.254` cloud metadata address is the single most common SSRF target — the private/reserved-IP check catches it along with loopback and link-local ranges on both IPv4 and IPv6. Note the residual TOCTOU gap: `fetch` re-resolves DNS after your check, so a short-TTL DNS record can rebind between validation and connection — for high-risk surfaces, resolve once and connect to the pinned IP.

## Securing LLM/AI features

If the app calls an LLM — chatbot, summarizer, agent, RAG — it inherits the [OWASP Top 10 for LLM Applications](https://genai.owasp.org/llm-top-10/):

- **Treat all model output as untrusted input (LLM05).** Never pass it straight into `eval`, SQL, a shell, `innerHTML`, or a file path — validate and encode exactly as raw user input.
- **Assume prompts can be hijacked (LLM01).** Any untrusted text in the context window — a user message, a fetched web page, a PDF — can carry instructions. The system prompt is not a security boundary; enforce permissions in code.
- **Keep secrets and other users' data out of prompts (LLM02/LLM07).** Anything in context can be echoed back.
- **Constrain tool/agent permissions (LLM06).** Scope tools to the minimum, require confirmation for destructive actions, validate every tool argument.
- **Bound consumption (LLM10).** Cap tokens, request rate, and loop/recursion depth.
- **Isolate retrieval data (LLM08).** In RAG, partition embeddings per tenant and validate documents before indexing so poisoned content can't steer answers.

## Triaging dependency-audit results

```
npm audit / govulncheck reports a vulnerability
├── Severity critical/high
│   ├── Reachable in your code path? YES → fix immediately (update/patch/replace)
│   │                                 NO  → fix soon, not a hard blocker
│   └── Fix available? YES → update  |  NO → workaround, replace, or allowlist with a review date
├── Severity moderate → fix next release cycle if reachable in production; else backlog
└── Severity low → track during regular dependency updates
```

Also: commit the lockfile and install with `npm ci` (or Go's committed `go.sum`) in CI for reproducible builds; review new dependencies for maintenance and download counts before adding them (OWASP A06, LLM03: supply chain); watch for typosquats (`cross-env` vs `crossenv`); be wary of `postinstall` scripts in unfamiliar packages.

## Gotchas

- If a secret is ever committed, deleting the line or rewriting history is not enough — assume it's compromised the moment it reaches a remote. Rotate (revoke and reissue) first, purge history second.
- `fetch`'s SSRF guard has a TOCTOU gap: it re-resolves DNS after your allowlist check runs, so an attacker controlling a short-TTL DNS record can pass validation pointing at a public IP and then have the actual connection resolve to an internal one. Pin the resolved IP for high-risk surfaces instead of trusting a second resolution to match the first.
- CORS with `credentials: true` and an echoed `Origin` header (instead of a fixed allowlist) is equivalent to a wildcard for any site that can get a victim to make a request — a very common near-miss that "looks" locked down because it isn't a literal `*`.
- A framework's auto-escaping (React's JSX, Go's `html/template`) only protects the specific sink it targets — `dangerouslySetInnerHTML`, a raw `text/template`, or building a URL/attribute value by string concatenation each reopen XSS even inside an otherwise-safe framework.

## Real-world grounding

Capital One's 2019 breach came from an SSRF against a cloud instance-metadata endpoint that let the attacker retrieve temporary credentials and pull data from S3 — "the request came from our own infrastructure" was not itself a trust signal, and the fix class is exactly the allowlist-plus-IP-check pattern above. The Equifax breach in the same era traced to a known, already-patched Apache Struts CVE left unapplied for months — a reminder that the dependency-audit triage step above is not optional busywork.

## Verification

- [ ] Trust boundaries and STRIDE threats were named before writing hardening code, not after
- [ ] Dependency audit shows no unaddressed critical/high, reachable vulnerabilities
- [ ] No secrets in source or git history; any past exposure was rotated, not just deleted
- [ ] Every new endpoint checks authentication AND per-resource authorization
- [ ] Server-side URL fetches (if any) are allowlisted and IP-checked, not fetched raw
- [ ] Security headers present in the actual response (verify in DevTools, don't assume config took effect)
- [ ] LLM output (if used) is treated as untrusted before reaching eval/SQL/DOM/shell

See `references/security-patterns.md` for full code samples referenced above.
