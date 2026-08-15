# Style Presets Reference

Read this before generating any slide HTML: the mandatory viewport-fitting
CSS base, density limits, preset catalog, and CSS gotchas all live here.

Favor abstract shapes over illustrations unless the user explicitly asks for
illustrated content.

## The viewport-fit rule

One slide equals exactly one viewport height. If content doesn't fit, split
it into another slide — never scroll inside a slide, and never shrink text
past a readable size to make it fit.

### Density limits

| Slide type | Maximum content |
|---|---|
| Title | 1 heading + 1 subtitle + optional tagline |
| Content | 1 heading + 4-6 bullets, or 2 short paragraphs |
| Feature grid | 6 cards max |
| Code | 8-10 lines max |
| Quote | 1 quote + attribution |
| Image | 1 image, ideally under 60vh |

## Mandatory base CSS

Include this block in every generated deck, then theme on top of it —
don't reinvent the viewport-fitting mechanics per deck.

```css
/* VIEWPORT FITTING: MANDATORY BASE */
html, body { height: 100%; overflow-x: hidden; }
html { scroll-snap-type: y mandatory; scroll-behavior: smooth; }

.slide {
    width: 100vw;
    height: 100vh;
    height: 100dvh;
    overflow: hidden;
    scroll-snap-align: start;
    display: flex;
    flex-direction: column;
    position: relative;
}

.slide-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    max-height: 100%;
    overflow: hidden;
    padding: var(--slide-padding);
}

:root {
    --title-size: clamp(1.5rem, 5vw, 4rem);
    --h2-size: clamp(1.25rem, 3.5vw, 2.5rem);
    --h3-size: clamp(1rem, 2.5vw, 1.75rem);
    --body-size: clamp(0.75rem, 1.5vw, 1.125rem);
    --small-size: clamp(0.65rem, 1vw, 0.875rem);

    --slide-padding: clamp(1rem, 4vw, 4rem);
    --content-gap: clamp(0.5rem, 2vw, 2rem);
    --element-gap: clamp(0.25rem, 1vw, 1rem);
}

.card, .container, .content-box {
    max-width: min(90vw, 1000px);
    max-height: min(80vh, 700px);
}

.grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(100%, 250px), 1fr));
    gap: clamp(0.5rem, 1.5vw, 1rem);
}

img, .image-container {
    max-width: 100%;
    max-height: min(50vh, 400px);
    object-fit: contain;
}

/* Short-viewport breakpoints: scale type/spacing down, hide chrome */
@media (max-height: 700px) {
    :root {
        --slide-padding: clamp(0.75rem, 3vw, 2rem);
        --content-gap: clamp(0.4rem, 1.5vw, 1rem);
        --title-size: clamp(1.25rem, 4.5vw, 2.5rem);
        --h2-size: clamp(1rem, 3vw, 1.75rem);
    }
}

@media (max-height: 600px) {
    :root {
        --slide-padding: clamp(0.5rem, 2.5vw, 1.5rem);
        --title-size: clamp(1.1rem, 4vw, 2rem);
        --body-size: clamp(0.7rem, 1.2vw, 0.95rem);
    }
    .nav-dots, .keyboard-hint, .decorative { display: none; }
}

@media (max-height: 500px) {
    :root {
        --slide-padding: clamp(0.4rem, 2vw, 1rem);
        --title-size: clamp(1rem, 3.5vw, 1.5rem);
        --h2-size: clamp(0.9rem, 2.5vw, 1.25rem);
        --body-size: clamp(0.65rem, 1vw, 0.85rem);
    }
}

@media (max-width: 600px) {
    :root { --title-size: clamp(1.25rem, 7vw, 2.5rem); }
    .grid { grid-template-columns: 1fr; }
}

@media (prefers-reduced-motion: reduce) {
    *, *::before, *::after {
        animation-duration: 0.01ms !important;
        transition-duration: 0.2s !important;
    }
    html { scroll-behavior: auto; }
}
```

### Viewport checklist

- Every `.slide` sets `height: 100vh`, `height: 100dvh`, `overflow: hidden`
- Every font size and spacing value uses `clamp()` or a viewport unit
- Images carry a `max-height` constraint
- Grids adapt via `auto-fit` + `minmax()`
- Short-height breakpoints exist at 700px, 600px, and 500px
- Anything that feels cramped gets split into another slide, not shrunk

## Mood to preset mapping

| Mood | Good presets |
|---|---|
| Impressed / confident | Bold Signal, Electric Studio, Dark Botanical |
| Excited / energized | Creative Voltage, Neon Cyber, Split Pastel |
| Calm / focused | Notebook Tabs, Paper & Ink, Swiss Modern |
| Inspired / moved | Dark Botanical, Vintage Editorial, Pastel Geometry |

## Preset catalog

Each preset: vibe, best-for context, font pairing, palette, and the one
signature visual move that makes it recognizable.

1. **Bold Signal** — confident, keynote-ready. Pitch decks, launches.
   Archivo Black + Space Grotesk. Charcoal base, hot-orange focal card,
   crisp white text. Signature: oversized section numbers on a high-contrast
   card over a dark field.
2. **Electric Studio** — clean, agency-polished. Client presentations,
   strategy reviews. Manrope only. Black/white/saturated cobalt. Signature:
   two-panel split, sharp editorial alignment.
3. **Creative Voltage** — energetic, retro-modern. Creative studios, brand
   storytelling. Syne + Space Mono. Electric blue, neon yellow, deep navy.
   Signature: halftone textures, badges, punchy contrast.
4. **Dark Botanical** — elegant, atmospheric. Luxury brands, premium
   narratives. Cormorant + IBM Plex Sans. Near-black, warm ivory, blush,
   gold, terracotta. Signature: blurred abstract circles, fine rules,
   restrained motion.
5. **Notebook Tabs** — editorial, tactile. Reports, structured reviews.
   Bodoni Moda + DM Sans. Cream paper on charcoal, pastel tabs. Signature:
   paper-sheet card with colored side tabs and binder details.
6. **Pastel Geometry** — approachable, friendly. Product overviews,
   onboarding. Plus Jakarta Sans only. Pale blue field, cream card, soft
   pink/mint/lavender accents. Signature: vertical pills, rounded cards,
   soft shadows.
7. **Split Pastel** — playful, creative. Agency intros, workshops. Outfit
   only. Peach + lavender split, mint badges. Signature: split backdrop,
   rounded tags, light grid overlay.
8. **Vintage Editorial** — witty, magazine-inspired. Personal brands,
   opinionated talks. Fraunces + Work Sans. Cream, charcoal, dusty warm
   accents. Signature: bordered callouts, punchy serif headlines.
9. **Neon Cyber** — futuristic, kinetic. AI/infra/dev-tools talks. Clash
   Display + Satoshi. Midnight navy, cyan, magenta. Signature: glow,
   particles, data-radar energy.
10. **Terminal Green** — developer-focused, hacker-clean. API/CLI/engineering
    demos. JetBrains Mono only. GitHub-dark palette + terminal green.
    Signature: scan lines, command-line framing.
11. **Swiss Modern** — minimal, data-forward. Corporate, product strategy.
    Archivo + Nunito. White, black, signal red. Signature: visible grids,
    asymmetry, geometric discipline.
12. **Paper & Ink** — literary, story-driven. Essays, manifesto decks.
    Cormorant Garamond + Source Serif 4. Warm cream, charcoal, crimson.
    Signature: pull quotes, drop caps, elegant rules.

If the user already knows the preset name, use it directly — skip preview
generation.

## Animation feel mapping

| Feeling | Motion direction |
|---|---|
| Dramatic / cinematic | Slow fades, parallax, large scale-ins |
| Techy / futuristic | Glow, particles, grid motion, scramble text |
| Playful / friendly | Springy easing, rounded shapes, floating motion |
| Professional / corporate | Subtle 200-300ms transitions, clean cuts |
| Calm / minimal | Very restrained movement, whitespace-first |
| Editorial / magazine | Strong hierarchy, staggered text/image interplay |

## CSS gotcha: negated functions

Never negate a `clamp()`/`min()`/`max()` call directly — browsers silently
ignore the whole declaration instead of erroring:

```css
/* Silently ignored — do not write this */
right: -clamp(28px, 3.5vw, 44px);
margin-left: -min(10vw, 100px);

/* Correct */
right: calc(-1 * clamp(28px, 3.5vw, 44px));
margin-left: calc(-1 * min(10vw, 100px));
```

## Validation sizes

Desktop: 1920×1080, 1440×900, 1280×720. Tablet: 1024×768, 768×1024. Mobile:
375×667, 414×896. Landscape phone: 667×375, 896×414.

## Anti-patterns

Purple-on-white startup templates; Inter/Roboto/Arial as the visual voice
unless the user explicitly wants utilitarian neutrality; bullet walls, tiny
type, or scrolling code blocks; decorative illustrations where abstract
geometry would carry the point better.
