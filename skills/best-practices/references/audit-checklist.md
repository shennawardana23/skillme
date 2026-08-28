# Audit checklist and tools reference

Full pass/fail checklist for running a best-practices audit, plus the
tool list to run it with. Load this when actually performing or
reporting on an audit — the day-to-day gotchas stay in `SKILL.md`.

## Audit checklist

### Security (critical)
- [ ] HTTPS enabled, no mixed content
- [ ] No vulnerable dependencies (`npm audit`)
- [ ] CSP headers configured (with `frame-ancestors`, `base-uri`, `form-action`)
- [ ] `require-trusted-types-for 'script'` enforced (or report-only during rollout)
- [ ] Third-party `<script>`/`<link rel="stylesheet">` pinned with SRI hashes
- [ ] Security headers present (HSTS, X-Content-Type-Options, Referrer-Policy)
- [ ] No exposed source maps (and `sourcesContent` stripped from uploaded ones)

### Compatibility
- [ ] Valid HTML5 doctype
- [ ] Charset declared first in head
- [ ] Viewport meta tag present
- [ ] No deprecated APIs used
- [ ] Passive event listeners for scroll/touch

### Code quality
- [ ] No console errors
- [ ] Valid HTML (no duplicate IDs)
- [ ] Semantic HTML elements used
- [ ] Proper error handling
- [ ] Memory cleanup in components

### UX
- [ ] No intrusive interstitials
- [ ] Permission requests in context
- [ ] Clear error messages
- [ ] Appropriate image aspect ratios

## Tools

| Tool | Purpose |
|------|---------|
| `npm audit` | Dependency vulnerabilities |
| [SecurityHeaders.com](https://securityheaders.com) | Header analysis |
| [W3C Validator](https://validator.w3.org) | HTML validation |
| Lighthouse | Best practices audit |
| [Observatory](https://observatory.mozilla.org) | Security scan |
