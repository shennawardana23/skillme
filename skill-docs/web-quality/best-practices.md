## What it does

Covers the modern web-development standards that don't fit cleanly into
performance, accessibility, or SEO — security headers and CSP, browser
compatibility, deprecated API avoidance, error handling, and source-map
hygiene. The thread running through the security section specifically:
each mitigation here targets a *specific, real, named* attack class
(Trusted Types for DOM-XSS sinks that a strict CSP alone doesn't stop,
Subresource Integrity for a compromised third-party CDN, `sourcesContent`
stripping for source maps that leak unminified original code) rather than
generic "harden everything" advice.

## When to reach for it

Reach for this for CSP/security-header questions, dependency
vulnerability checks, browser-compatibility patterns (feature detection,
polyfills), or general code-quality review that doesn't fit the other
three web-quality categories. For a full four-category audit, start at
`web-quality-audit` instead.

## Common questions

- **"We already have a strict CSP — do we still need Trusted Types?"**
  Yes if any code writes to `innerHTML`, calls `eval`, or touches another
  DOM-XSS sink. A strict CSP blocks loading untrusted *script files*, but
  it does nothing to stop a plain string from reaching those sinks —
  Trusted Types (Baseline across major browsers since early 2026) closes
  that specific gap by making the sinks reject raw strings and require a
  typed object from a named policy. Roll out with
  `Content-Security-Policy-Report-Only` first to find every sink usage
  before enforcing.
- **"Is Subresource Integrity actually necessary if we trust our CDN?"**
  The `polyfill.io` supply-chain attack (2024) is the concrete case for
  why: a previously-trusted CDN was compromised and used to serve malware
  to roughly 100,000 sites. SRI (`integrity="sha384-..."` on the
  `<script>`/`<link>` tag) makes the browser refuse to execute a file
  whose hash doesn't match, regardless of what the CDN currently serves.
- **"We strip `X-XSS-Protection` from our headers — is that a
  regression?"** No — the opposite. The legacy browser XSS auditor that
  header controlled was deprecated and removed (Chrome 78, Edge 17), and
  it had introduced its own vulnerabilities in some cases. Sending it
  does nothing useful in current browsers; a strict CSP plus Trusted
  Types is the actual current mitigation.

## It's working if

- CSP is enforced via response header (not just a meta tag), with
  `frame-ancestors`, `base-uri`, and `form-action` set
- Every third-party `<script>`/`<link rel="stylesheet">` from a CDN
  carries an `integrity` hash
- Production source maps either aren't shipped publicly or have
  `sourcesContent` stripped before upload to an error tracker
- `npm audit` (or equivalent) runs as part of the normal dependency
  workflow, not just occasionally by hand

## Where it fits

The fourth category under `web-quality-audit`, and the one most likely
to overlap with `security-review`/`security-and-hardening` for anything
that's really an application-security question rather than a
browser-facing best practice.
