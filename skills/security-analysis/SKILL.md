---
name: security-analysis
description: Runs a two-pass vulnerability scan on a snippet, file, or diff — a static pattern pass across fixed categories (injection, secrets, auth, crypto, deserialization, path traversal, resource exhaustion), then an LLM-reasoning pass for business-logic and race-condition flaws the patterns can't catch. Use for a dedicated security audit that must produce a per-finding severity, exact location, and concrete remediation.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Security Analysis

Two passes, run in order. The static pass is mechanical and exhaustive — report every match, even ones that turn out to be false positives on inspection, and only discard a match after confirming why it's safe. The reasoning pass catches what patterns structurally cannot: this is not a re-run of `security-review`'s priority checklist or `security-and-hardening`'s proactive threat-modeling — it's a scan workflow that ends in a severity-counted findings report.

## Pass 1: static pattern scan

Check every category below against the target. Treat this like a grep sweep — you're looking for *shapes*, not judging exploitability yet.

**Injection**
- SQL: string concatenation or `fmt.Sprintf`/`+` building a query with an external value; raw input placed directly in a `WHERE`/`ORDER BY` clause
- Command: `exec.Command`/`os/exec`/`subprocess`/`child_process.exec` where any argument traces back to external input
- Template: untrusted input rendered through a non-escaping template engine (Go `text/template` instead of `html/template`)
- LDAP: unsanitized input interpolated into an LDAP filter string

**Credential and secret exposure**
- Hardcoded passwords, API keys, tokens, connection strings in source
- Secrets appearing in log statements, even at DEBUG level
- Secrets or internal detail echoed back in error responses
- Environment variables logged directly instead of redacted

**Authentication and authorization**
- Endpoints reachable with no authentication check at all
- Authentication present but no check that the caller is authorized for *this* resource (ownership/role check)
- JWT validation skipped, or `alg: none`/unverified-signature acceptance
- Session tokens with low entropy or no expiry

**Cryptographic weaknesses**
- MD5 or SHA1 for anything security-sensitive (password hashing, MACs)
- ECB mode block ciphers
- Hardcoded IV/nonce (reused nonce breaks stream-cipher and GCM confidentiality)
- `math/rand` (or language equivalent) used for tokens, session IDs, or anything security-sensitive

**Unsafe deserialization**
- `json.Unmarshal` into `interface{}`/`any` with no schema validation downstream
- `encoding/gob`, pickle, or similar binary deserialization of untrusted input
- YAML/TOML parsing of user-supplied content without type constraints (YAML in particular supports type coercion and anchors that can be abused)

**Path traversal**
- File operations building a path from user input without `filepath.Clean` + an allowlist-prefix check
- A file server exposing a directory with user-controlled subpaths

**Resource exhaustion**
- No rate limiting on an expensive endpoint (auth, search, export)
- Request body read with no size cap (`io.LimitedReader`/`bodyLimit`)
- Outbound HTTP calls with no timeout

## Pass 2: reasoning pass

After the static sweep, reason about what patterns can't catch:

- **Business logic vulnerabilities** — e.g. a discount code endpoint that trusts a client-supplied price, a workflow that can be replayed out of order to reach an invalid state.
- **Race conditions in auth or state checks** — check-then-act sequences (check balance, then debit) with no lock or atomic operation between them (TOCTOU).
- **Trust boundary violations** — data that crossed from an untrusted zone (third-party API response, uploaded file, another tenant's data) and is now used without being re-validated at the new boundary.

## Output format

```
## Security Risk Level: LOW | MEDIUM | HIGH | CRITICAL

## Findings

### [CRITICAL] Type — Brief title
- Location: file:line
- Description: what the vulnerability is and how it's exploited
- Remediation: exact fix, with a code example

### [HIGH] Type — Brief title
...

## Summary
Total: N critical, N high, N medium, N low
```

Never report a finding without a concrete remediation and a code example — "this could be a problem" without a fix wastes the reader's time re-deriving what you already know.

## Gotchas

- A parameterized-looking query can still be vulnerable if the *table or column name* — not the value — is built from user input; placeholders only protect values, never identifiers. Scan identifier interpolation as its own SQL-injection subcase.
- `json.Unmarshal(data, &v)` into a concrete struct is not automatically safe just because it's not `interface{}` — a struct field of type `interface{}`, or a custom `UnmarshalJSON` that does its own type-switching, reopens the same hole one level down.
- A JWT library that "validates the signature" can still accept `alg: none` or an attacker-chosen algorithm if the verification call doesn't pin the expected algorithm explicitly — this was the mechanism behind a well-known class of 2015-era JWT library vulnerabilities and still needs an explicit allowlisted-algorithm check today.
- Rate limiting applied only at a reverse proxy or API gateway layer misses any endpoint reachable by an internal service-to-service call — check whether the limiter is enforced at the same boundary the finding's threat model assumes.

## Real-world grounding

The 2017 Equifax breach traced to an unpatched Apache Struts CVE with a public fix available for months — a known, patched vulnerability nobody applied, not a novel exploit. It's a reminder that a static pattern scan and dependency-CVE check together catch a large share of real incidents; the reasoning pass exists for the remainder — the business-logic and race-condition classes that show up in bug bounty reports far more often than classic injection does today.

## Verification

- [ ] Every category in Pass 1 was checked, not just the ones that seemed likely to apply
- [ ] Every finding includes file:line, a description of exploitability, and a concrete remediation
- [ ] The reasoning pass considered at least business logic and race conditions, not just static patterns
- [ ] The summary count matches the findings actually listed above it
