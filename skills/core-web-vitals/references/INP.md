# INP optimization reference

Detailed code patterns, the full optimization checklist, and the
debugging script for Interaction to Next Paint. Load this when actually
fixing an INP issue, not for a general Core Web Vitals overview.

## Common INP issues

**1. Long tasks blocking main thread**
```javascript
// ❌ Long synchronous task
function processLargeArray(items) {
  items.forEach(item => expensiveOperation(item));
}

// ✅ Break into chunks and yield to the scheduler. scheduler.yield() is the
//    recommended modern API — its continuation is queued at a boosted
//    priority so the rest of your work resumes ahead of unrelated tasks,
//    while still letting the browser handle pending input first.
async function processLargeArray(items) {
  const CHUNK_SIZE = 100;
  for (let i = 0; i < items.length; i += CHUNK_SIZE) {
    items.slice(i, i + CHUNK_SIZE).forEach(expensiveOperation);

    if ('scheduler' in window && 'yield' in scheduler) {
      await scheduler.yield();
    } else {
      // Fallback for browsers without scheduler.yield (Safari, older Firefox).
      // setTimeout(0) yields but loses priority — your continuation may run
      // after unrelated tasks the browser picked up in between.
      await new Promise(r => setTimeout(r, 0));
    }
  }
}
```

**2. Heavy event handlers**
```javascript
// ❌ All work in handler
button.addEventListener('click', () => {
  // Heavy computation
  const result = calculateComplexThing();
  // DOM updates
  updateUI(result);
  // Analytics
  trackEvent('click');
});

// ✅ Prioritize visual feedback, then yield before doing the heavy work
button.addEventListener('click', async () => {
  // 1. Immediate visual feedback (cheap DOM update)
  button.classList.add('loading');

  // 2. Yield so the browser can paint the loading state before we block
  if ('scheduler' in window && 'yield' in scheduler) {
    await scheduler.yield();
  }

  // 3. Now do the heavy work — the user already saw the click register
  const result = calculateComplexThing();
  updateUI(result);

  // 4. Lowest-priority work last, when the main thread is idle
  if ('requestIdleCallback' in window) {
    requestIdleCallback(() => trackEvent('click'));
  } else {
    setTimeout(() => trackEvent('click'), 0);
  }
});
```

**3. Third-party scripts**
```javascript
// ❌ Eagerly loaded, blocks interactions
<script src="https://heavy-widget.com/widget.js"></script>

// ✅ Lazy loaded on interaction or visibility
const loadWidget = () => {
  import('https://heavy-widget.com/widget.js')
    .then(widget => widget.init());
};
button.addEventListener('click', loadWidget, { once: true });
```

**4. Excessive re-renders (React/Vue)**
```javascript
// ❌ Re-renders entire tree
function App() {
  const [count, setCount] = useState(0);
  return (
    <div>
      <Counter count={count} />
      <ExpensiveComponent /> {/* Re-renders on every count change */}
    </div>
  );
}

// ✅ Memoized expensive components
const MemoizedExpensive = React.memo(ExpensiveComponent);

function App() {
  const [count, setCount] = useState(0);
  return (
    <div>
      <Counter count={count} />
      <MemoizedExpensive />
    </div>
  );
}
```

## INP optimization checklist

```markdown
- [ ] No tasks > 50ms on main thread
- [ ] Event handlers complete quickly (< 100ms)
- [ ] Visual feedback provided immediately
- [ ] Heavy work deferred with requestIdleCallback
- [ ] Third-party scripts don't block interactions
- [ ] Debounced input handlers where appropriate
- [ ] Web Workers for CPU-intensive operations
```

## INP debugging

```javascript
// Identify slow interactions. durationThreshold: 40 matches what the
// web-vitals library uses — 16 (one frame) fires on nearly every interaction
// and drowns the console; 40 surfaces interactions that are starting to feel
// sluggish without spamming.
new PerformanceObserver((list) => {
  for (const entry of list.getEntries()) {
    if (entry.duration > 200) {
      console.warn('Slow interaction:', {
        type: entry.name,
        duration: entry.duration,
        processingStart: entry.processingStart,
        processingEnd: entry.processingEnd,
        target: entry.target
      });
    }
  }
}).observe({ type: 'event', buffered: true, durationThreshold: 40 });
```

For field debugging across real users, prefer the `web-vitals/attribution` build of the [web-vitals library](https://github.com/GoogleChrome/web-vitals) — `onINP()` from that build attaches a `LoAF` (Long Animation Frame) breakdown identifying the longest script and the input/processing/presentation phase that ate the budget.
