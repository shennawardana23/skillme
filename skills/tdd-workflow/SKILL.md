---
name: tdd-workflow
description: Runs the full TDD workflow for a TypeScript/JavaScript feature — user journeys to test cases, unit tests (Jest/Vitest + Testing Library), API/integration tests, Playwright E2E, and a coverage gate wired into CI. Use when building a Next.js/React/Node feature end to end, not just a single unit test, or when asked to set up a coverage threshold or CI test step.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# TDD Workflow (TypeScript/JavaScript)

The RED-GREEN-REFACTOR loop tells you how to write one test. This skill is
one level up: how a whole feature's test surface is planned across the
test pyramid, wired into CI, and gated on coverage, in a Jest/Vitest +
Playwright stack. For the loop itself and the Prove-It bug-fix pattern,
see `skills/test-driven-development/` — that skill applies inside every
step below.

Use `skills/laravel-tdd/` instead for PHP/Laravel; use this skill for
Next.js/React/Node projects with Jest, Vitest, Testing Library, or
Playwright.

## From user journey to test cases

Start from the observable behavior, not the implementation:

```
As a user, I want to search for markets semantically,
so that I can find relevant markets even without exact keywords.
```

```typescript
describe('Semantic Search', () => {
  it('returns relevant markets for query', async () => {})
  it('handles empty query gracefully', async () => {})
  it('falls back to substring search when Redis unavailable', async () => {})
  it('sorts results by similarity score', async () => {})
})
```

Write the cases before any implementation exists; the shape of the
`describe` block is the spec for what "search" means here — including the
fallback and the empty-input path, not just the primary path.

## Test pyramid for this stack

- **Unit** (Jest/Vitest + Testing Library) — components and pure
  functions in isolation. Test user-visible behavior (`screen.getByText`),
  not internal state (`component.state.count`) — internal state can be
  correct while the rendered output is wrong, and refactors that preserve
  behavior will still break a state-shape assertion.
- **Integration** — API route handlers against a real or in-memory
  database boundary, external services mocked.
- **E2E (Playwright)** — a small number of critical user flows through a
  real browser. Prefer semantic selectors (`button:has-text("Submit")`,
  `[data-testid=...]`) over CSS class selectors, which break on any style
  refactor unrelated to the behavior under test.

```typescript
// Integration: API route
describe('GET /api/markets', () => {
  it('returns markets successfully', async () => {
    const response = await GET(new NextRequest('http://localhost/api/markets'))
    const data = await response.json()
    expect(response.status).toBe(200)
    expect(data.success).toBe(true)
  })

  it('returns 400 for an invalid query parameter', async () => {
    const response = await GET(new NextRequest('http://localhost/api/markets?limit=invalid'))
    expect(response.status).toBe(400)
  })
})
```

```typescript
// E2E: Playwright, semantic selectors
test('user can search and filter markets', async ({ page }) => {
  await page.goto('/markets')
  await page.fill('input[placeholder="Search markets"]', 'election')
  await expect(page.locator('[data-testid="market-card"]')).toHaveCount(5, { timeout: 5000 })
  await page.click('button:has-text("Active")')
  await expect(page.locator('[data-testid="market-card"]')).toHaveCount(3)
})
```

## Mocking external services

Mock at the module boundary, not inside the function under test:

```typescript
jest.mock('@/lib/redis', () => ({
  searchMarketsByVector: jest.fn(() => Promise.resolve([{ slug: 'test-market', similarity_score: 0.95 }])),
  checkRedisHealth: jest.fn(() => Promise.resolve({ connected: true })),
}))
```

## Coverage gate

A coverage threshold in CI turns "we should have good coverage" into an
enforced build failure:

```json
{
  "jest": {
    "coverageThresholds": {
      "global": { "branches": 80, "functions": 80, "lines": 80, "statements": 80 }
    }
  }
}
```

```yaml
# CI
- name: Run Tests
  run: npm test -- --coverage
```

Pick the threshold your project will actually enforce — 80% is a common
starting point, not a universal requirement — and treat a drop below it
as a build failure, not a warning to fix later.

## Test isolation

Each test sets up its own data; tests must not depend on execution order:

```typescript
// Each test is independent
test('creates user', () => {
  const user = createTestUser()
})
test('updates user', () => {
  const user = createTestUser() // not the user from the previous test
})
```

## Gotchas

- A coverage number can hit 100% of lines while asserting nothing — a test
  that calls a function and checks no exception was thrown exercises the
  line without verifying the behavior. Coverage measures what ran, not
  what was checked.
- `page.waitForTimeout(600)` in a Playwright test to wait out a debounce
  is a flake source under CI load; prefer waiting on the resulting DOM
  state (`await expect(locator).toHaveCount(...)`) which retries until it
  matches or times out, rather than a fixed sleep.
- Mocking `@/lib/redis` at the module level means a real Redis outage in
  production won't be caught by this test suite — integration or E2E tests
  against the fallback path are what actually verify the degraded-mode
  behavior, not the unit test with the module mocked out.
- CSS class selectors (`.css-class-xyz`) in E2E tests couple the test to
  styling; a CSS refactor with zero behavior change turns green E2E tests
  red for the wrong reason.

## Real-world grounding

Google's internal test-size taxonomy (small/medium/large, publicly
described in the "Software Engineering at Google" book) maps closely onto
this pyramid: small tests run in a single process with no I/O and are
expected to be fast and numerous; large tests may spin up real
dependencies and are deliberately kept few, because their slower, flakier
nature makes them expensive to run on every commit. The E2E tier here
plays that "large" role — reserved for flows whose failure would matter
enough to justify the cost.

## Verification

- [ ] Test cases exist for the primary path, an edge case (empty/invalid
      input), and at least one failure/fallback path
- [ ] Unit tests assert on user-visible output, not internal state
- [ ] E2E tests use semantic selectors, not CSS classes
- [ ] A coverage threshold is configured and enforced in CI, not just
      reported
- [ ] No test depends on another test's side effects or run order
