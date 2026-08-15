# Test-Driven Development — Patterns and Anti-Patterns

## Test state, not interactions

Assert on the outcome of an operation, not on which internal methods were
called. Interaction-based tests break under refactors that don't change
behavior.

```go
// Good: asserts the outcome
func TestListTasks_SortedByCreatedAtDesc(t *testing.T) {
	tasks := listTasks(t, sortByCreatedAtDesc)
	if !tasks[0].CreatedAt.After(tasks[1].CreatedAt) {
		t.Fatal("expected newest first")
	}
}
```

## DAMP over DRY in tests

Production code favors DRY. Tests favor DAMP (Descriptive And Meaningful
Phrases) — each test should read as a self-contained specification, even at
the cost of some duplication across tests. A shared helper that hides the
actual input/output of a test forces the reader to trace through
indirection to understand what's being verified.

## Real implementation > fake > stub > mock

Preference order, most to least trustworthy:

1. **Real implementation** — highest confidence, catches real integration bugs.
2. **Fake** — an in-memory stand-in with real behavior (e.g., an in-memory
   key-value store implementing the same interface as the real one).
3. **Stub** — returns canned data, no real behavior.
4. **Mock (interaction-verifying)** — asserts on calls made; use sparingly,
   only when the real dependency is slow, non-deterministic, or has
   uncontrollable side effects (a payment gateway, an email provider).

Over-mocking is the single most common cause of a test suite that's green
while production is broken.

## Anti-pattern table

| Anti-pattern | Problem | Fix |
|---|---|---|
| Testing implementation details | Breaks on behavior-preserving refactors | Assert on inputs/outputs only |
| Flaky, order-dependent tests | Erodes trust in the whole suite | Deterministic assertions, isolated state per test |
| Snapshot abuse | Large snapshots nobody actually reviews | Use narrowly, review every diff |
| Mocking the function under test | Passes unconditionally | Mock only external dependencies |
| No reproduction test for a bug fix | Can't prove the fix works, no regression guard | Prove-It Pattern: write the failing test first |

## Arrange-Act-Assert

```typescript
it('marks overdue tasks after the deadline passes', () => {
  // Arrange
  const task = createTask({ deadline: new Date('2025-01-01') });
  // Act
  const result = checkOverdue(task, new Date('2025-01-02'));
  // Assert
  expect(result.isOverdue).toBe(true);
});
```
