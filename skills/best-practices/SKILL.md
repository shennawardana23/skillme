---
name: best-practices
description: Apply modern web development best practices for security, compatibility, and code quality. Use when asked to "apply best practices", "security audit", "modernize code", "code quality review", or "check for vulnerabilities".
license: Apache-2.0
metadata:
  category: "web-quality"
  version: "0.1.0"
---

# Best practices

Modern web development standards based on Lighthouse best practices audits. Covers security, browser compatibility, and code quality patterns.

## Security

**HTTPS everywhere:**
```html
<!-- ❌ Mixed content -->
<img src="http://example.com/image.jpg">

<!-- ✅ HTTPS only -->
<img src="https://example.com/image.jpg">
```
Avoid protocol-relative URLs (`//example.com/...`) — they hide the actual scheme from reviewers. Send `Strict-Transport-Security: max-age=31536000; includeSubDomains; preload`.

**Input sanitization:**
```javascript
// ❌ XSS vulnerable
element.innerHTML = userInput;
document.write(userInput);

// ✅ Safe text content
element.textContent = userInput;

// ✅ If HTML needed, sanitize
import DOMPurify from 'dompurify';
element.innerHTML = DOMPurify.sanitize(userInput);
```

**Secure cookies (server-side):**
```
Set-Cookie: session=abc123; Secure; HttpOnly; SameSite=Strict; Path=/
```

Beyond these, a real security posture needs: a CSP with nonces (`script-src`, `frame-ancestors`, `base-uri`, `form-action`); Trusted Types (`require-trusted-types-for 'script'`) to stop raw strings reaching `innerHTML`/`eval` even under a strict CSP; Subresource Integrity hashes on every third-party `<script>`/`<link>` you don't control; the standard security-header set (HSTS, `X-Content-Type-Options`, `Referrer-Policy`, `Permissions-Policy` — **never** `X-XSS-Protection`, it's deprecated and was itself a vulnerability); and `npm audit` run regularly, watching specifically for prototype-pollution via recursive merges of untrusted input (`_.merge`, `$.extend`). Full syntax and code for all of these is in [references/security-hardening.md](references/security-hardening.md) — load it when implementing or auditing any of them.

## Browser compatibility

Ship an HTML5 doctype, `<meta charset="UTF-8">` as the first element in `<head>`, and a responsive viewport meta tag. Detect features (`'IntersectionObserver' in window`, CSS `@supports`), never sniff the user agent — UA strings are brittle and spoofable. Prefer bundling polyfills at build time (Babel/SWC + `core-js`) over runtime checks, and never load a polyfill from a third-party CDN you don't control — the 2024 `polyfill.io` supply-chain attack served malware to ~100k sites via exactly that pattern. Full markup and polyfill-loading code: [references/compatibility-and-legacy-apis.md](references/compatibility-and-legacy-apis.md).

## Deprecated APIs

Avoid `document.write` (blocks parsing — use dynamic `createElement`/`appendChild` instead) and synchronous XHR (blocks the main thread — use `fetch`). Attach `touchstart`/`wheel` listeners with `{ passive: true }` so the browser doesn't wait on a possible `preventDefault()` before scrolling. Replacement code: [references/compatibility-and-legacy-apis.md](references/compatibility-and-legacy-apis.md).

## Console & errors

```javascript
// ❌ Errors in production
console.log('Debug info'); // Remove in production
throw new Error('Unhandled'); // Catch all errors

// ✅ Proper error handling
try {
  riskyOperation();
} catch (error) {
  errorTracker.captureException(error);
  showErrorMessage('Something went wrong. Please try again.');
}
```

Wrap React trees in an error boundary (`static getDerivedStateFromError` + `componentDidCatch`, reporting to your error tracker and rendering a fallback), and always register global handlers for uncaught errors and unhandled promise rejections:
```javascript
window.addEventListener('error', (event) => errorTracker.captureException(event.error));
window.addEventListener('unhandledrejection', (event) => errorTracker.captureException(event.reason));
```

## Source maps

Never ship `devtool: 'source-map'` to production — it exposes your original source. Use `'hidden-source-map'` (uploaded to your error tracker but not referenced in the bundle) or disable maps entirely in production builds, and strip `sourcesContent` from anything you do upload — the map otherwise embeds your full unminified source. For Vite, use `sourcemap: 'hidden'`. Full config: [references/compatibility-and-legacy-apis.md](references/compatibility-and-legacy-apis.md).

## Performance best practices

```javascript
// ❌ Blocking script / CSS import
<script src="heavy-library.js"></script>
@import url('other-styles.css');

// ✅ Deferred script, parallel-loading link tags
<script defer src="heavy-library.js"></script>
<link rel="stylesheet" href="styles.css">
```

Use event delegation (one listener on a container, matched via `e.target.matches(...)`) instead of attaching a handler to every list item. Always remove listeners you add — prefer `AbortController` (`{ signal }`) so cleanup is a single `controller.abort()` call instead of matching every `addEventListener`/`removeEventListener` pair by hand.

## Code quality

Use semantic HTML5 elements (`<header>`, `<nav>`, `<main>`, `<article>`) instead of `<div class="...">` soup — screen readers, SEO, and future maintainers all rely on it. Keep IDs unique, keep list children as `<li>`, and never nest interactive elements (`<button>` inside `<a>`). When sizing `<img>`, use the image's *actual* aspect ratio for `width`/`height` (prevents layout shift) and `object-fit` only when you deliberately need to crop.

## Permissions & privacy

Request permissions (geolocation, camera, mic) in response to a user action, with an explanation of why — never on page load, which reads as hostile and gets auto-denied. Restrict unused powerful features via `Permissions-Policy`.

## Audit checklist

The full security/compatibility/code-quality/UX checklist and the tool list to run it with (`npm audit`, SecurityHeaders.com, W3C Validator, Lighthouse, Observatory) live in [references/audit-checklist.md](references/audit-checklist.md) — load it when actually running or reporting on an audit.

## Gotchas

- `X-XSS-Protection` is deprecated and was itself exploitable — a strict CSP + Trusted Types replaces it, don't send it.
- A CSP alone doesn't stop DOM-XSS via `innerHTML`/`eval` — you need Trusted Types (or manual sanitization) on top of it.
- `Object.assign(target, userInput)` looks safe but isn't, when `target` can reach `Object.prototype` and `userInput` carries `__proto__` — use a null-prototype object or `structuredClone` for untrusted data.
- Loading a polyfill or library from a CDN you don't control is a supply-chain risk even over HTTPS — pin it with SRI or self-host.
- Source maps with `sourcesContent` intact leak your full unminified source if anyone obtains the `.map` file.

## References

- [MDN Web Security](https://developer.mozilla.org/en-US/docs/Web/Security)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Web Quality Audit](../web-quality-audit/SKILL.md)
