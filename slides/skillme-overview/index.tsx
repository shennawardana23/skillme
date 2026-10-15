/* Hallmark-adapted stamp (open-slide has no tokens.css/log.json — recorded here instead)
 * macrostructure: mixed (hero / quote-hero / terminal-card / dark-bullets / card-grid /
 *   big-number / compare ×2 / closing-hero)
 * theme: custom — Natural Linen coordinating palette, cool half (sage / slate / olive-charcoal)
 * tone: editorial, muted, restrained
 * pre-emit self-critique: Philosophy 4 · Hierarchy 4 · Execution 4 · Specificity 4 · Restraint 5 · Variety 4
 * — v2: fixed the v1 sin (eyebrow+heading+bullets repeated as one template on every page).
 *   Every page below has a different shape; background swatch varies deliberately per page,
 *   using the SAME reference palette as depth rather than one flat color everywhere.
 * — v3: added an Install page borrowing the terminal-card shape from
 *   slides/getting-started (its Init/Prompt pages) — same window-chrome + typed-command
 *   structure, recolored into this deck's own palette (dark charcoal terminal, sage
 *   prompt/check marks) rather than copying that deck's dark techy hues wholesale.
 */
import type { CSSProperties, ReactNode } from 'react';
import type { DesignSystem, Page, SlideMeta, SlideTransition } from '@open-slide/core';

// Palette: the cooler half of the Sherwin Williams "Natural Linen" coordinating
// grid — cream, sage, slate blue-gray, olive-gray, dark olive-charcoal. Used as
// a small family of surfaces (not one flat background), the way the reference
// image itself reads as a *combination*, not a single swatch.
export const design: DesignSystem = {
  palette: { bg: '#E7E2D3', text: '#33342B', accent: '#4F5A62' },
  fonts: {
    display: 'Georgia, "Times New Roman", serif',
    body: '-apple-system, BlinkMacSystemFont, "Inter", system-ui, sans-serif',
  },
  typeScale: { hero: 176, body: 40 },
  radius: 10,
};

// Full swatch family (same source image), for surfaces the DesignSystem shape
// doesn't cover.
const cream = design.palette.bg; // lightest
const charcoal = design.palette.text; // darkest, doubles as a page bg on p3
const slate = design.palette.accent;
const sage = '#9CA69A';
const oliveGray = '#585F4E';
const lightBlueGray = '#AEB7BB';
const onDark = '#EDE9DD'; // text color used on the two dark-bg pages
const onDarkMuted = '#B7BEB2';
const cardBg = '#DED6C0';

const PAD = 130;

// RISE — quiet, 6px Y + opacity. One DNA for the whole deck.
const EASE_OUT = 'cubic-bezier(0, 0, 0.2, 1)';
const EASE_IN = 'cubic-bezier(0.4, 0, 1, 1)';
export const transition: SlideTransition = {
  duration: 200,
  exit: {
    duration: 140,
    easing: EASE_IN,
    keyframes: [
      { opacity: 1, transform: 'translateY(0)' },
      { opacity: 0, transform: 'translateY(-4px)' },
    ],
  },
  enter: {
    duration: 200,
    delay: 80,
    easing: EASE_OUT,
    keyframes: [
      { opacity: 0, transform: 'translateY(6px)' },
      { opacity: 1, transform: 'translateY(0)' },
    ],
  },
};

const base = {
  width: '100%',
  height: '100%',
  fontFamily: 'var(--osd-font-body)',
} as const;

// Small swatch strip — the literal reference palette as a design element,
// not decoration bolted on: it's the visual citation for "this deck's colors
// came from here," reused on the cover and the close for a bookend rhyme.
const SwatchStrip = () => (
  <div style={{ display: 'flex', gap: 10 }}>
    {[cream, sage, oliveGray, lightBlueGray, slate, charcoal].map((c) => (
      <div key={c} style={{ width: 22, height: 22, borderRadius: '50%', background: c }} />
    ))}
  </div>
);

// --- Page 1: Cover — hero wordmark + swatch citation -----------------------

const Cover: Page = () => (
  <div
    style={{
      ...base,
      background: cream,
      color: charcoal,
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      padding: `0 ${PAD}px`,
    }}
  >
    <SwatchStrip />
    <h1
      style={{
        fontFamily: 'var(--osd-font-display)',
        fontSize: 'var(--osd-size-hero)',
        fontWeight: 700,
        margin: '36px 0 28px',
        lineHeight: 1,
      }}
    >
      skillme
    </h1>
    <p style={{ fontSize: 32, lineHeight: 1.5, color: oliveGray, maxWidth: 1300, fontWeight: 400 }}>
      A testable knowledge base of real engineering practice — not a
      collection of advice nobody has checked.
    </p>
  </div>
);

// --- Page 2: Why — a single large statement, not a bullet list -------------

const WhyThisExists: Page = () => (
  <div
    style={{
      ...base,
      background: sage,
      color: charcoal,
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      padding: `0 ${PAD}px`,
    }}
  >
    <div style={{ fontSize: 26, letterSpacing: '0.18em', textTransform: 'uppercase', fontWeight: 600, opacity: 0.65 }}>
      Why
    </div>
    <h2
      style={{
        fontFamily: 'var(--osd-font-display)',
        fontSize: 88,
        fontWeight: 700,
        lineHeight: 1.15,
        margin: '28px 0 0',
        maxWidth: 1560,
      }}
    >
      Most "prompt libraries" are good-sounding advice nobody has checked.
    </h2>
    <div style={{ fontSize: 34, lineHeight: 1.5, marginTop: 40, color: oliveGray, maxWidth: 1300 }}>
      Every one of skillme's 137+ skills ships with a deterministic test
      suite, graded on real model output — one rigor bar, not an exception.
    </div>
  </div>
);

// --- Pages 3–4: terminal-card pair, structurally mirroring getting-started's
// Init (p4) + Prompt (p5) — window-chrome terminal, then a two-pane agent-CLI
// split — recolored into this deck's palette and filled with skillme's own
// actual install command and a real skill trigger (code-review), not the
// other deck's slide-authoring content. ---------------------------------

const mono = 'ui-monospace, SFMono-Regular, "SF Mono", Menlo, monospace';

// Typewriter + stream-in animations for the two terminal pages — same idea as
// getting-started's own type/stream effect (command text reveals character by
// character, results fade in staggered after), scoped under an so- prefix so
// it can't collide with any other slide's injected styles.
const terminalStyles = `
  @keyframes so-type {
    from { clip-path: inset(0 100% 0 0); }
    to   { clip-path: inset(0 0 0 0); }
  }
  @keyframes so-fadeIn { from { opacity: 0; } to { opacity: 1; } }
  @keyframes so-blink { 0%, 49% { opacity: 1; } 50%, 100% { opacity: 0; } }
  .so-stream { opacity: 0; animation: so-fadeIn 0.4s ease forwards; }
  .so-caret {
    display: inline-block;
    width: 12px;
    height: 26px;
    background: currentColor;
    animation: so-blink 1s steps(1) infinite;
  }
`;
const TerminalStyles = () => <style>{terminalStyles}</style>;

// A line of command text that reveals left-to-right, one character-step at a
// time — clip-path is relative to the span's own box, so this works
// regardless of the container's width (a plain animated `width` would reveal
// almost instantly, since the target is the container, not the text).
const TypedLine = ({ text, color, delay = 0 }: { text: string; color: string; delay?: number }) => (
  <span
    style={{
      display: 'inline-block',
      whiteSpace: 'nowrap',
      color,
      clipPath: 'inset(0 100% 0 0)',
      animation: `so-type ${Math.max(text.length * 0.035, 0.6)}s steps(${text.length}, end) forwards`,
      animationDelay: `${delay}s`,
    }}
  >
    {text}
  </span>
);

const TrafficLights = () => (
  <div style={{ display: 'flex', gap: 10 }}>
    {['#ff5f56', '#ffbd2e', '#27c93f'].map((c) => (
      <span
        key={c}
        style={{
          width: 14,
          height: 14,
          borderRadius: '50%',
          background: c,
          boxShadow: 'inset 0 0 0 1px rgba(0,0,0,0.25)',
        }}
      />
    ))}
  </div>
);

const TerminalWindow = ({
  title,
  children,
  style,
}: {
  title: string;
  children: ReactNode;
  style?: CSSProperties;
}) => (
  <div
    style={{
      background: charcoal,
      borderRadius: 'var(--osd-radius)',
      overflow: 'hidden',
      display: 'flex',
      flexDirection: 'column',
      boxShadow: '0 30px 60px -24px rgba(0,0,0,0.35)',
      ...style,
    }}
  >
    <div style={{ height: 52, padding: '0 20px', display: 'flex', alignItems: 'center', gap: 16, background: oliveGray, flexShrink: 0 }}>
      <TrafficLights />
      <div style={{ flex: 1, textAlign: 'center', fontFamily: mono, fontSize: 20, color: onDarkMuted }}>
        {title}
      </div>
      <div style={{ minWidth: 40 }} />
    </div>
    {children}
  </div>
);

const TerminalLine = ({ cmd, result, delay = 0 }: { cmd?: string; result?: string; delay?: number }) =>
  cmd !== undefined ? (
    <div className="so-stream" style={{ display: 'flex', gap: 16, minHeight: 38, animationDelay: `${delay}s` }}>
      <span style={{ color: sage }}>$</span>
      <TypedLine text={cmd} color={onDark} delay={delay} />
    </div>
  ) : (
    <div className="so-stream" style={{ display: 'flex', gap: 16, minHeight: 38, animationDelay: `${delay}s` }}>
      <span style={{ color: sage }}>✔</span>
      <TypedLine text={result ?? ''} color={onDarkMuted} delay={delay} />
    </div>
  );

const PromptCaret = ({ delay = 0 }: { delay?: number }) => (
  <div className="so-stream" style={{ display: 'flex', gap: 16, alignItems: 'center', animationDelay: `${delay}s` }}>
    <span style={{ color: sage }}>$</span>
    <span className="so-caret" style={{ color: onDark }} />
  </div>
);

// Page 3 — mirrors getting-started's Init page (p4): eyebrow + heading +
// subtitle + one terminal card, filled with skillme's real install path.
const InstallTerminal: Page = () => (
  <div
    style={{
      ...base,
      background: cream,
      color: charcoal,
      padding: PAD,
      display: 'flex',
      flexDirection: 'column',
      gap: 48,
    }}
  >
    <TerminalStyles />
    <div>
      <div style={{ fontSize: 26, letterSpacing: '0.18em', textTransform: 'uppercase', fontWeight: 600, color: oliveGray }}>
        01 / Install
      </div>
      <h2 style={{ fontFamily: 'var(--osd-font-display)', fontSize: 84, fontWeight: 700, margin: '20px 0 0', lineHeight: 1.15 }}>
        Two ways in, one repo
      </h2>
      <p style={{ fontSize: 32, lineHeight: 1.5, color: oliveGray, margin: '20px 0 0', maxWidth: 1300 }}>
        A managed plugin bundle for the whole team, or editable files via
        skills.sh — pick one, never both.
      </p>
    </div>
    <TerminalWindow title="~/skillme — zsh" style={{ flex: 1 }}>
      <div
        style={{
          flex: 1,
          padding: '32px 44px',
          fontFamily: mono,
          fontSize: 26,
          lineHeight: 1.5,
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          gap: 14,
        }}
      >
        <TerminalLine cmd="claude plugin marketplace add shennawardana23/skillme" delay={0.2} />
        <TerminalLine result="Marketplace registered" delay={2.2} />
        <TerminalLine cmd="claude plugin install skillme@skillme" delay={3.1} />
        <TerminalLine result="skillme@skillme installed — 137 skills active" delay={4.6} />
        <PromptCaret delay={6.4} />
      </div>
    </TerminalWindow>
  </div>
);

const AgentLine = ({
  speaker,
  text,
  delay = 0,
}: {
  speaker: 'user' | 'assistant' | 'tool';
  text: string;
  delay?: number;
}) => {
  const label = speaker === 'user' ? 'you' : speaker === 'assistant' ? 'agent' : 'tool';
  const color = speaker === 'user' ? onDark : speaker === 'assistant' ? sage : onDarkMuted;
  return (
    <div className="so-stream" style={{ display: 'flex', gap: 14, marginBottom: 20, animationDelay: `${delay}s` }}>
      <span style={{ fontFamily: mono, fontSize: 18, color: onDarkMuted, width: 60, flexShrink: 0, paddingTop: 3 }}>
        {label}
      </span>
      <div style={{ fontFamily: mono, fontSize: 23, lineHeight: 1.5 }}>
        <TypedLine text={text} color={color} delay={delay} />
      </div>
    </div>
  );
};

// Page 4 — mirrors getting-started's Prompt page (p5): eyebrow + heading +
// two-pane grid (agent CLI on the left, its real result on the right). The
// right pane is skillme's own actual differentiator — a passing smeval run —
// not a copy of the other deck's live-slide-preview pane.
const PromptTrigger: Page = () => (
  <div
    style={{
      ...base,
      background: cream,
      color: charcoal,
      padding: PAD,
      display: 'flex',
      flexDirection: 'column',
      gap: 40,
    }}
  >
    <TerminalStyles />
    <div>
      <div style={{ fontSize: 26, letterSpacing: '0.18em', textTransform: 'uppercase', fontWeight: 600, color: oliveGray }}>
        02 / Prompt the agent
      </div>
      <h2 style={{ fontFamily: 'var(--osd-font-display)', fontSize: 84, fontWeight: 700, margin: '20px 0 0', lineHeight: 1.15 }}>
        Ask. Watch the right skill trigger.
      </h2>
    </div>
    <div style={{ flex: 1, display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 40, minHeight: 0 }}>
      <TerminalWindow title="claude · ~/hotel-api">
        <div style={{ flex: 1, padding: '28px 36px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <AgentLine speaker="user" delay={0.2} text="review this PR for missing authorization checks" />
          <AgentLine speaker="assistant" delay={2.0} text="Using code-review to check this diff…" />
          <AgentLine speaker="tool" delay={3.5} text="read diff.patch" />
          <AgentLine speaker="tool" delay={4.2} text="skill loaded → code-review" />
          <PromptCaret delay={5.3} />
        </div>
      </TerminalWindow>
      <TerminalWindow title="smeval run code-review">
        <div
          style={{
            flex: 1,
            padding: '32px 44px',
            fontFamily: mono,
            fontSize: 24,
            lineHeight: 1.6,
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            gap: 14,
          }}
        >
          <TerminalLine cmd="smeval run code-review" delay={0.3} />
          <TerminalLine result="missing-auth-check ... PASS" delay={1.3} />
          <TerminalLine result="sql-injection-detection ... PASS" delay={2.4} />
          <TerminalLine result="2/2 cases passed" delay={3.7} />
        </div>
      </TerminalWindow>
    </div>
  </div>
);

// --- Page 5: Proven — dark surface, heading + bullets (earns its shape) ----

const ProvenNotWritten: Page = () => (
  <div style={{ ...base, background: charcoal, color: onDark, padding: PAD, display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
    <div style={{ fontSize: 26, letterSpacing: '0.18em', textTransform: 'uppercase', fontWeight: 600, color: lightBlueGray }}>
      The eval layer
    </div>
    <h2 style={{ fontFamily: 'var(--osd-font-display)', fontSize: 84, fontWeight: 700, margin: '28px 0 0', lineHeight: 1.15 }}>
      Proven, not just written
    </h2>
    <ul style={{ margin: '56px 0 0', paddingLeft: 44 }}>
      <li style={{ fontSize: 40, lineHeight: 1.6, marginBottom: 28, color: onDark }}>
        smeval — this repo's own zero-dependency Go eval runner
      </li>
      <li style={{ fontSize: 40, lineHeight: 1.6, marginBottom: 28, color: onDark }}>
        Each skill tested in isolation, graded with quoted evidence
      </li>
      <li style={{ fontSize: 40, lineHeight: 1.6, marginBottom: 0, color: onDarkMuted }}>
        A claim survives being checked, or it gets fixed
      </li>
    </ul>
  </div>
);

// --- Page 6: Catalog-wide checks (explicit card instances, not a map) ------

const ToolCard = ({ name, blurb, accent }: { name: string; blurb: string; accent: string }) => (
  <div
    style={{
      flex: 1,
      background: cardBg,
      borderRadius: 'var(--osd-radius)',
      padding: '32px 28px',
      display: 'flex',
      flexDirection: 'column',
      gap: 14,
      borderTop: `6px solid ${accent}`,
    }}
  >
    <div style={{ fontSize: 32, fontWeight: 700, fontFamily: 'var(--osd-font-display)', color: charcoal }}>{name}</div>
    <div style={{ fontSize: 24, lineHeight: 1.5, color: oliveGray }}>{blurb}</div>
  </div>
);

const CatalogWideChecks: Page = () => (
  <div style={{ ...base, background: cream, color: charcoal, padding: PAD, display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
    <div style={{ fontSize: 26, letterSpacing: '0.18em', textTransform: 'uppercase', fontWeight: 600, color: oliveGray }}>
      At 137+ skills
    </div>
    <h2 style={{ fontFamily: 'var(--osd-font-display)', fontSize: 84, fontWeight: 700, margin: '28px 0 56px', lineHeight: 1.15 }}>
      Catalog-wide checks
    </h2>
    <div style={{ display: 'flex', gap: 36 }}>
      <ToolCard name="trigger-run" blurb="Does the right skill win among all of them?" accent={sage} />
      <ToolCard name="security-scan" blurb="Static risk scan, free, no model calls" accent={oliveGray} />
      <ToolCard name="similarity-check" blurb="Flags duplicate or overlapping skills" accent={lightBlueGray} />
      <ToolCard name="feedback-check" blurb="Closes the human-review gap, enforced" accent={slate} />
    </div>
  </div>
);

// --- Page 7: Big number — text as the whole design --------------------------

const BigNumber: Page = () => (
  <div
    style={{
      ...base,
      background: slate,
      color: onDark,
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'flex-start',
      padding: `0 ${PAD}px`,
    }}
  >
    <div style={{ fontSize: 28, letterSpacing: '0.18em', textTransform: 'uppercase', fontWeight: 600, color: lightBlueGray }}>
      Token-efficient by enforcement
    </div>
    <div
      style={{
        fontFamily: 'var(--osd-font-display)',
        fontSize: 340,
        fontWeight: 700,
        lineHeight: 1,
        margin: '24px 0 32px',
      }}
    >
      5,484
    </div>
    <p style={{ fontSize: 36, lineHeight: 1.5, color: onDarkMuted, maxWidth: 1300, margin: 0 }}>
      words already relocated from always-loaded SKILL.md bodies into
      on-demand references/ this session — verified line by line, audited,
      not just policy.
    </p>
  </div>
);

// --- Page 8: Context-dependent testing — two-column before/after -----------

const CompareColumn = ({
  label,
  body,
  tone,
}: {
  label: string;
  body: string;
  tone: 'muted' | 'strong';
}) => (
  <div style={{ flex: 1 }}>
    <div
      style={{
        fontSize: 24,
        letterSpacing: '0.14em',
        textTransform: 'uppercase',
        fontWeight: 700,
        color: tone === 'strong' ? slate : oliveGray,
        marginBottom: 20,
      }}
    >
      {label}
    </div>
    <div style={{ fontSize: 34, lineHeight: 1.55, color: tone === 'strong' ? charcoal : oliveGray, maxWidth: 680 }}>
      {body}
    </div>
  </div>
);

const ContextDependentTesting: Page = () => (
  <div style={{ ...base, background: cream, color: charcoal, padding: PAD, display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
    <h2 style={{ fontFamily: 'var(--osd-font-display)', fontSize: 76, fontWeight: 700, margin: '0 0 56px', lineHeight: 1.15, maxWidth: 1500 }}>
      Tests real investigation, not just prompts
    </h2>
    <div style={{ display: 'flex', gap: 64, alignItems: 'flex-start' }}>
      <CompareColumn
        label="A described bug"
        body="Hands the model everything in the prompt text — nothing left to actually go find."
        tone="muted"
      />
      <div style={{ width: 2, alignSelf: 'stretch', background: lightBlueGray }} />
      <CompareColumn
        label="A fixture-backed case"
        body="Seeds real files into an isolated workspace — the model reproduces, isolates, and fixes a genuine cross-file Go bug."
        tone="strong"
      />
    </div>
  </div>
);

// --- Page 9: Free CI / costed eval — asymmetric split, not a grid -----------

const FreeCiCostedEval: Page = () => (
  <div style={{ ...base, background: cream, color: charcoal, padding: PAD, display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
    <h2 style={{ fontFamily: 'var(--osd-font-display)', fontSize: 76, fontWeight: 700, margin: 0, lineHeight: 1.15, maxWidth: 1400 }}>
      Free CI, costed eval — on purpose
    </h2>
    <div style={{ display: 'flex', gap: 0, marginTop: 56, borderRadius: 'var(--osd-radius)', overflow: 'hidden' }}>
      <div style={{ flex: 1.1, background: oliveGray, color: onDark, padding: '36px 40px' }}>
        <div style={{ fontSize: 26, letterSpacing: '0.1em', textTransform: 'uppercase', fontWeight: 700, marginBottom: 16 }}>
          Free, every push
        </div>
        <div style={{ fontSize: 32, lineHeight: 1.5 }}>Schema validation + static security scan gate every push</div>
      </div>
      <div style={{ flex: 1, background: cardBg, color: charcoal, padding: '36px 40px' }}>
        <div style={{ fontSize: 26, letterSpacing: '0.1em', textTransform: 'uppercase', fontWeight: 700, color: oliveGray, marginBottom: 16 }}>
          Costed, by hand
        </div>
        <div style={{ fontSize: 32, lineHeight: 1.5, color: oliveGray }}>Live model evals stay a deliberate, human-run step</div>
      </div>
    </div>
  </div>
);

// --- Page 10: Closing — bookend hero + swatch citation ----------------------

const Closing: Page = () => (
  <div
    style={{
      ...base,
      background: cream,
      color: charcoal,
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      padding: `0 ${PAD}px`,
    }}
  >
    <div style={{ fontSize: 26, letterSpacing: '0.18em', textTransform: 'uppercase', fontWeight: 600, color: oliveGray }}>
      Open standard
    </div>
    <h2
      style={{
        fontFamily: 'var(--osd-font-display)',
        fontSize: 92,
        fontWeight: 700,
        margin: '28px 0 40px',
        lineHeight: 1.15,
        maxWidth: 1500,
      }}
    >
      Built to be checked, not just believed
    </h2>
    <p style={{ fontSize: 34, lineHeight: 1.5, color: oliveGray, maxWidth: 1300, margin: '0 0 48px' }}>
      agentskills.io spec-compliant, validated in CI. Apache-2.0 — fork the
      eval harness, not just the idea.
    </p>
    <SwatchStrip />
    <div style={{ fontSize: 24, color: oliveGray, letterSpacing: '0.04em', marginTop: 24 }}>
      github.com/shennawardana23/skillme
    </div>
  </div>
);

export const meta: SlideMeta = {
  title: 'skillme — Key Features',
  createdAt: '2026-08-30T08:50:31.899Z',
};

export default [
  Cover,
  WhyThisExists,
  InstallTerminal,
  PromptTrigger,
  ProvenNotWritten,
  CatalogWideChecks,
  BigNumber,
  ContextDependentTesting,
  FreeCiCostedEval,
  Closing,
] satisfies Page[];
