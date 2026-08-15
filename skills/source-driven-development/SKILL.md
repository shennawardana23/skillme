---
name: source-driven-development
description: Grounds framework-specific code decisions in official documentation instead of training data or memory. Use before writing framework- or library-specific code (forms, routing, data fetching, auth, state management), when the user asks for "current best practices" or "documented"/"verified" code, or before reviewing code that leans on framework conventions.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Source-Driven Development

Verify framework-specific patterns against official documentation before
writing them, rather than from memory. Training data goes stale, APIs
deprecate, and best practices evolve — a pattern that was correct in a
training snapshot can be the wrong one for the version actually installed.
Every framework-specific decision should trace back to a source the user
can check themselves.

## When to use

Building boilerplate or patterns that will be copied elsewhere in the
project; the user asks for "current," "documented," or "verified" code;
implementing something where the framework's recommended approach matters
(forms, routing, data fetching, state management, auth); reviewing code
that leans on framework conventions; any time about to write
framework-specific code from memory.

**Skip it for:** logic that's version-independent (loops, conditionals,
data structures), renames/typo fixes, or when the user explicitly wants
speed over verification.

## The process

```
DETECT ──→ FETCH ──→ IMPLEMENT ──→ CITE
```

### 1. Detect stack and versions

Read the dependency file — `package.json`, `composer.json`,
`go.mod`, `requirements.txt`/`pyproject.toml`, `Gemfile` — and state what
was found:

```
STACK DETECTED: React 19.1.0, Vite 6.2.0 (from package.json)
→ Fetching official docs for the relevant pattern.
```

If a version is missing or ambiguous, ask — don't guess. The version is
what determines which pattern is actually correct.

### 2. Fetch the specific documentation page

Fetch the page for the feature being implemented, not the framework
homepage and not a general web search.

| Priority | Source |
|---|---|
| 1 | Official documentation (react.dev, docs.djangoproject.com, pkg.go.dev) |
| 2 | Official blog / changelog |
| 3 | Web standards references (MDN, web.dev) |
| 4 | Browser/runtime compatibility tables (caniuse.com) |

Not authoritative, never cite as primary: Stack Overflow answers, blog
posts or tutorials, AI-generated summaries, or your own training data —
verifying that training data is the entire point of this skill.

```
BAD:  fetch the React homepage, or search "django auth best practices"
GOOD: fetch react.dev/reference/react/useActionState
GOOD: fetch docs.djangoproject.com/en/6.0/topics/auth/
```

When two official sources disagree (a migration guide contradicts the API
reference), surface the discrepancy to the user and verify which pattern
actually works against the detected version rather than picking one
silently.

### 3. Implement following the documented pattern

Use the signatures shown in the docs, not from memory. If the docs show a
newer approach, use it; if they deprecate a pattern, don't use it; if they
don't cover something, flag it as unverified rather than guessing.

When the docs and the existing codebase disagree, surface the conflict
instead of silently picking one:

```
CONFLICT: the codebase uses useState for form-loading state; React 19
docs recommend useActionState for this. (react.dev/reference/react/useActionState)
A) modern pattern (useActionState) — consistent with current docs
B) match existing code (useState) — consistent with the codebase
→ Which do you want?
```

### 4. Cite sources

Every framework-specific pattern gets a citation the user can verify.

```typescript
// React 19 form handling — react.dev/reference/react/useActionState#usage
const [state, formAction, isPending] = useActionState(submitOrder, initialState);
```

Full URLs, not shortened. Prefer deep anchor links
(`/useActionState#usage`) over the bare page — anchors tend to survive a
docs-site restructuring better than top-level URLs. Quote the relevant
passage when it justifies a non-obvious decision. If nothing verifiable
was found:

```
UNVERIFIED: no official documentation found for this pattern. Based on
training data and may be outdated — verify before production use.
```

A flagged gap is more useful than false confidence dressed up as a
citation.

## Gotchas

- A URL that resolves and looks plausible is not the same as a citation
  that was actually fetched and read this session — cite what was checked,
  not what seems likely to exist.
- Docs for the *latest* major version are not automatically correct for
  the version pinned in the dependency file; a migration guide's "new way"
  can require a version the project hasn't upgraded to yet.
- When official sources disagree with each other, defaulting to whichever
  was fetched first (rather than surfacing the conflict) silently commits
  the user to an unverified guess.
- "I recall this API" and "I fetched this API's current docs" produce
  identically confident-sounding prose — the difference only shows up in
  whether a citation is present, so the absence of a citation is itself a
  signal worth noticing during review.

## Real-world grounding

React's `useActionState` (introduced in React 19, previously shipped as
`useFormState` in a canary release) is a documented case where the
"recommended" hook for handling form-submission state changed name and
signature between versions in the framework's own official docs — code
written from a slightly earlier training snapshot would use the old name
and compile against nothing in a fresh install. Checking the fetched
version-pinned docs, rather than a remembered API name, is what catches
this before the user hits an import error.

## Verification

- [ ] Framework/library versions were identified from the dependency file
- [ ] Official documentation was fetched for the specific pattern being
      implemented, not the general homepage
- [ ] All citations are official docs, not blog posts, forums, or
      training-data recall
- [ ] Code follows the pattern shown in the version-pinned documentation
- [ ] Conflicts (docs vs. docs, or docs vs. existing code) were surfaced
      to the user, not silently resolved
- [ ] Anything that couldn't be verified is explicitly flagged as
      unverified
