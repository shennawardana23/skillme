---
name: coding-standards
description: Cross-language code-quality heuristics -- function size, nesting depth, magic numbers, comment intent, naming -- for reviewing or writing code in any language. Use when reviewing code for maintainability, setting up lint rules, or onboarding conventions. For Go-specific idioms use go-service-idioms, for API contracts use api-design, for Laravel use laravel-patterns, for Vue/Nuxt use vue-nuxt-frontend-patterns -- this skill is the language-agnostic layer underneath those.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Coding Standards

This is the language-agnostic layer: decision heuristics that apply
whether the file is Go, TypeScript, or PHP. For a specific language or
framework's idioms, use the dedicated skill instead —
`skills/go-service-idioms/`, `skills/api-design/`,
`skills/laravel-patterns/`, `skills/vue-nuxt-frontend-patterns/`. This
skill only covers what's genuinely cross-cutting: it skips KISS/DRY/YAGNI
definitions and other things any competent reviewer already knows, and
focuses on the concrete thresholds and rationalizations that actually
distinguish a maintainable change from one that technically works.

## When to use

Reviewing code for maintainability across languages; setting up lint
rules; establishing conventions a mixed-language team will actually follow
without a language-specific skill to lean on.

## Function size and nesting

A function whose body exceeds ~50 lines is usually doing more than one
thing — the fix is extracting named sub-steps, not shrinking whitespace:

```typescript
// One function, three responsibilities hidden inside it
function processMarketData() { /* validate, transform, persist, ~100 lines */ }

// Three named steps, each independently testable
function processMarketData() {
  const validated = validateData()
  const transformed = transformData(validated)
  return saveData(transformed)
}
```

Nesting past 3-4 levels is a comparable signal. Guard clauses (early
returns) collapse a pyramid of conditions into a flat list of exits,
which is easier to scan and to add a fifth condition to later:

```typescript
// 5 levels deep
if (user) { if (user.isAdmin) { if (market) { if (market.isActive) { if (hasPermission) { /* ... */ } } } } }

// Flat
if (!user) return
if (!user.isAdmin) return
if (!market) return
if (!market.isActive) return
if (!hasPermission) return
// ...
```

## Magic numbers and strings

An unexplained literal in a conditional or a timeout is a landmine for the
next reader, who has to guess whether `3` and `500` are related, arbitrary,
or load-bearing:

```typescript
if (retryCount > 3) { }
setTimeout(callback, 500)

// vs.
const MAX_RETRIES = 3
const DEBOUNCE_DELAY_MS = 500
if (retryCount > MAX_RETRIES) { }
```

The name doesn't just aid readability — it's the place a future change
gets made once instead of hunting down every literal `3` in the file to
check whether it's this one.

## Comments: why, not what

A comment restating what the next line already says adds a second thing
that can drift out of sync with the code; a comment explaining a
non-obvious *why* is the only kind that survives a refactor with its value
intact:

```typescript
// Bad — restates the code
// Increment counter by 1
count++

// Good — explains a non-obvious reason
// Exponential backoff avoids overwhelming the API during an outage
const delay = Math.min(1000 * Math.pow(2, retryCount), 30000)
```

## Naming

Names should describe what a thing is or does, not its type or a
placeholder: `isUserAuthenticated`, not `flag`; `fetchMarketData`, not
`market()` (a noun where a verb is expected reads as a getter that might
not do I/O, hiding the fact that it does). Each language in this catalog
has its own casing convention — `camelCase` for TS/JS, Go's exported
`PascalCase` — but never mix two casing conventions within one project;
that's a stronger tell of copy-pasted, unreviewed code than any single
bad name.

## Error handling

Every call that can fail needs an explicit decision at the call site: this
error propagates, or it's handled here — not silently swallowed. A caught
exception with an empty handler, or a Go error assigned to `_`, is a
decision made by omission that the next reader has to reverse-engineer.

```typescript
// Silent failure — the caller has no idea this can fail
async function fetchData(url) {
  const response = await fetch(url)
  return response.json()
}

// Explicit: what fails, and what happens when it does
async function fetchData(url: string) {
  const response = await fetch(url)
  if (!response.ok) throw new Error(`HTTP ${response.status}: ${response.statusText}`)
  return await response.json()
}
```

## Type safety

`any` (TypeScript) or an untyped `interface{}`/`any` parameter (Go) is not
neutral — it's a specific decision to move a class of bug from compile
time to runtime, for every future caller, not just the one being written
right now:

```typescript
// any defers the type check to runtime, for every caller
function getMarket(id: any): Promise<any> { }

// the compiler catches a wrong-typed argument at the call site
interface Market { id: string; name: string; status: 'active' | 'resolved' | 'closed' }
function getMarket(id: string): Promise<Market> { }
```

## Test naming (cross-language AAA)

A test name should describe the behavior under test, not restate that a
test exists — `test('works')` gives zero information when it fails; `test
falls back to substring search when Redis unavailable` names exactly what
broke. Structure the body Arrange-Act-Assert, in that order, so the shape
of the test itself communicates the behavior being verified.

## Gotchas

- "80% coverage" is a threshold a specific project chooses and enforces in
  CI, not a universal industry law — treating it as a fixed target that
  applies unconditionally to every project misses that the number should
  come from the project's own risk tolerance, not from a generic
  standards doc.
- 100% line coverage is compatible with zero real assertions — a test that
  calls a function and checks nothing threw exercises every line without
  verifying a single behavior. Coverage percentage measures what ran, not
  what was checked.
- Consistent formatting enforced by a linter/formatter is a solved problem
  in every language in this catalog (`gofmt`, `prettier`, `pint`) — a
  human debating tab-vs-space in review is spending review time on
  something a pre-commit hook should have already settled.
- A naming convention violated in exactly one file usually means that file
  was copy-pasted from a different project or generated by a different
  tool, and is worth a closer look at what else came along with it.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "This function is long but it's all one operation" | If it can't be summarized in one sentence without "and," it's more than one operation, regardless of how it reads to the person who just wrote it. |
| "I'll name it properly later" | Later rarely arrives before three other files start depending on the placeholder name, at which point renaming has a bigger blast radius. |
| "It's just a quick script, standards don't apply" | A "quick script" that gets copied into the next quick script inherits every shortcut taken in the first one — the second reader has no way to know which parts were deliberate. |

## Real-world grounding

Postel's robustness principle ("be conservative in what you send, liberal
in what you accept") — originally written for network protocol
implementers in RFC 761 (1980) — generalizes past networking to any
function boundary: validate and reject malformed input explicitly at the
boundary (liberal in what's accepted means handling more shapes of valid
input gracefully, not silently accepting invalid input), rather than
letting a bad value propagate several call frames deep before it surfaces
as a confusing failure far from its actual cause.

## Verification

- [ ] No function mixes more than one clear responsibility without being
      split into named steps
- [ ] No literal number/string in a conditional or timeout lacks a name
      explaining what it represents
- [ ] Comments explain why, not what — and none are stale relative to the
      code beside them
- [ ] No error is silently discarded; each either propagates or is
      handled with a visible reason
- [ ] No untyped escape hatch (`any`/untyped `interface{}`) was used where
      a concrete type was available
- [ ] Naming convention is consistent within the file/project, not mixed
