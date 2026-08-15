---
name: e2e-testing
description: Guides Playwright end-to-end test design — Page Object Model, CI configuration, and flaky test diagnosis. Use when writing browser-based end-to-end tests, debugging a test that fails intermittently, setting up Playwright CI configuration, or reviewing E2E test code for race conditions and arbitrary waits.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# End-to-End Testing (Playwright)

E2E tests sit at the top of the test pyramid — expensive, slow, and the
first place flakiness shows up. Most flakiness traces back to one thing:
the test racing the page instead of waiting for a specific, observable
condition.

## Page Object Model

```typescript
export class ReservationsPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto("/reservations");
  }

  async search(query: string) {
    await this.page.locator('[data-testid="search-input"]').fill(query);
    await this.page.waitForResponse((r) => r.url().includes("/api/search"));
  }

  rows() {
    return this.page.locator('[data-testid="reservation-row"]');
  }
}
```

## The single highest-value fix: stop waiting on time, wait on state

```typescript
// Flaky: races the animation/network, passes locally, fails in CI
await page.waitForTimeout(2000);
await page.click('[data-testid="confirm"]');

// Stable: waits for the actual condition
await page.locator('[data-testid="confirm"]').waitFor({ state: "visible" });
await page.locator('[data-testid="confirm"]').click();
```

Playwright locators already auto-wait for actionability (attached,
visible, stable, enabled) before acting — `page.waitForTimeout` is almost
never the right tool and should be treated as a review finding, not a style
preference.

## Gotchas

- `page.click(selector)` on a raw selector string skips Playwright's
  auto-waiting in some older API patterns; `page.locator(selector).click()`
  is the form that gets the actionability checks. Prefer locators
  throughout, not selector strings passed directly to action methods.
- A test that passes locally but fails in CI is frequently a timing issue
  invisible on a fast local machine and a slow, resource-constrained CI
  runner — reproduce with `--repeat-each=10` before assuming it's an
  environment misconfiguration.
- `test.skip`/`test.fixme` used to silence a flaky test without an
  attached issue reference is a one-way door — the test quietly stops
  providing coverage and nobody notices it never runs again. Always attach
  a tracking reference when quarantining a test.
- CI-only flakiness is commonly a resource contention issue, not a logic
  bug: setting `workers: 1` for CI (serializing tests) and `retries: 2` is
  a legitimate stabilizing measure while investigating root cause, but
  should not be treated as the fix itself — a retried, silently-passing
  flaky test is still hiding a real race condition.

## Configuration baseline

```typescript
export default defineConfig({
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  use: {
    trace: "on-first-retry",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },
});
```

## Real-world grounding

Flaky E2E tests are a widely documented cross-industry problem, not a sign
of one team's carelessness — the common root causes converge across
codebases: arbitrary sleeps instead of condition waits, shared mutable
state between parallel test workers, and animation/transition timing. The
fix is never "add a retry and move on" as the permanent state — retries
mask the failure rate without addressing why the test races the
application in the first place; use them as a stabilizing measure while
you diagnose the actual race, not as the diagnosis's conclusion.

## Verification

- [ ] No `waitForTimeout` used as a substitute for waiting on a real condition
- [ ] Every action goes through a `page.locator(...)` call, not a bare selector
- [ ] Flaky/quarantined tests have an attached tracking reference
- [ ] CI config sets deterministic retries and worker count, separate from local dev config
- [ ] Critical user flows (not every flow) get E2E coverage — the rest belongs in the unit/integration layers
