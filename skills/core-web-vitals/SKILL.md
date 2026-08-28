---
name: core-web-vitals
description: Optimize Core Web Vitals (LCP, INP, CLS) for better page experience and search ranking. Use when asked to "improve Core Web Vitals", "fix LCP", "reduce CLS", "optimize INP", "page experience optimization", or "fix layout shifts".
license: Apache-2.0
metadata:
  category: "web-quality"
  version: "0.1.0"
---

# Core Web Vitals optimization

Targeted optimization for the three Core Web Vitals metrics that affect Google Search ranking and user experience.

## The three metrics

| Metric | Measures | Good | Needs work | Poor |
|--------|----------|------|------------|------|
| **LCP** | Loading | ≤ 2.5s | 2.5s – 4s | > 4s |
| **INP** | Interactivity | ≤ 200ms | 200ms – 500ms | > 500ms |
| **CLS** | Visual Stability | ≤ 0.1 | 0.1 – 0.25 | > 0.25 |

Google measures at the **75th percentile** — 75% of page visits must meet "Good" thresholds.

---

## LCP: Largest Contentful Paint

LCP measures when the largest visible content element renders — usually a hero image, video poster, large background image, `<svg>` element, or large text block.

In rough order of impact, the common causes are: slow server response (TTFB > 800ms), the LCP image not preloaded with `fetchpriority="high"`, render-blocking CSS/JS in `<head>`, and client-side-rendered content that isn't present in the initial HTML (fix with SSR/SSG/streaming instead of fetching after mount). Beyond the landed page, the Speculation Rules API can prerender likely-next-page navigations so the LCP a user experiences on the *next* click is near-instant.

Full code patterns for each fix, the Speculation Rules setup and its caveats, the LCP optimization checklist, and the element-identification script are in [references/LCP.md](references/LCP.md) — load it when actually fixing an LCP issue.

---

## INP: Interaction to Next Paint

INP measures responsiveness across ALL interactions (clicks, taps, key presses) during a page visit. It reports the worst interaction (at 98th percentile for high-traffic pages).

Total INP = **Input Delay** + **Processing Time** + **Presentation Delay**

| Phase | Target | Optimization |
|-------|--------|--------------|
| Input Delay | < 50ms | Reduce main thread blocking |
| Processing | < 100ms | Optimize event handlers |
| Presentation | < 50ms | Minimize rendering work |

The recurring causes are long synchronous tasks blocking the main thread, event handlers that do heavy work before giving visual feedback, eagerly-loaded third-party scripts, and excessive component re-renders. The fix pattern for the first two is the same: give immediate visual feedback, then yield to the scheduler (`scheduler.yield()`, falling back to `setTimeout(0)`) before doing the expensive work.

Full code patterns, the INP optimization checklist, and the debugging script (including the `web-vitals/attribution` build for field data) are in [references/INP.md](references/INP.md) — load it when actually fixing an INP issue.

---

## CLS: Cumulative Layout Shift

CLS measures unexpected layout shifts. A shift occurs when a visible element changes position between frames without user interaction.

**CLS Formula:** `impact fraction × distance fraction`

The recurring causes are images/iframes/embeds without reserved dimensions, content injected above existing content, web fonts causing FOUT/FOIT reflow, and animations that transition layout properties (`height`/`width`) instead of `transform`/`opacity`.

Full code patterns, the CLS optimization checklist, and the layout-shift debugging script are in [references/CLS.md](references/CLS.md) — load it when actually fixing a CLS issue.

---

## Measurement tools

### Lab testing
- **Chrome DevTools** → Performance panel, Lighthouse
- **WebPageTest** → Detailed waterfall, filmstrip
- **Lighthouse CLI** → `npx lighthouse <url>`

### Field data (real users)
- **Chrome User Experience Report (CrUX)** → BigQuery or API
- **Search Console** → Core Web Vitals report
- **web-vitals library** → Send to your analytics

```javascript
import {onLCP, onINP, onCLS} from 'web-vitals';

function sendToAnalytics({name, value, rating}) {
  gtag('event', name, {
    event_category: 'Web Vitals',
    value: Math.round(name === 'CLS' ? value * 1000 : value),
    event_label: rating
  });
}

onLCP(sendToAnalytics);
onINP(sendToAnalytics);
onCLS(sendToAnalytics);
```

Lab tools (DevTools, Lighthouse, WebPageTest) show what *could* happen on one run; field data (CrUX, Search Console, the web-vitals library) shows what real users actually experienced at the 75th percentile — treat lab scores as diagnostic, not as the ranking signal itself.

## Framework quick fixes

Next.js, React, and Vue/Nuxt each have idiomatic one-liners for the LCP/INP/CLS fixes above (e.g. `next/image` with `priority`, `useTransition`, `NuxtImg` with `preload`). See [references/framework-quick-fixes.md](references/framework-quick-fixes.md).

## References

- [web.dev LCP](https://web.dev/articles/lcp)
- [web.dev INP](https://web.dev/articles/inp)
- [web.dev CLS](https://web.dev/articles/cls)
- [Performance skill](../performance/SKILL.md)
