---
name: documentation-lookup
description: Fetch current library and framework documentation via the Context7 MCP server instead of relying on training data. Use for setup or configuration questions, API references, code examples that depend on a library's behavior, or whenever the user names a specific framework or library (React, Next.js, Prisma, Supabase, etc.).
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Documentation Lookup (Context7)

Training data goes stale the moment a library ships a new version. When a question depends on a library's actual current behavior, fetch it live via Context7 rather than answering from memory.

## Core Concepts

- **Context7**: an MCP server exposing live, indexed library documentation — use it in place of training data for anything version- or API-specific.
- **resolve-library-id**: returns Context7-compatible library IDs (format `/org/project` or `/org/project/version`) from a library name and query.
- **query-docs**: fetches documentation and code snippets for a given library ID and question. Always call `resolve-library-id` first — never call `query-docs` with a guessed library ID.

## When to Use

- Setup or configuration questions ("how do I configure middleware in this framework?")
- Requests for code that depends on a specific library's API ("write a query using this ORM")
- API or reference questions ("what are this library's auth methods?")
- The user names a specific framework or library

## How It Works

### Step 1: Resolve the Library ID

Call `resolve-library-id` with the library name taken from the question and the user's full question as the query (the full question improves relevance ranking). A valid Context7 library ID is required before calling `query-docs` — never skip this step.

### Step 2: Select the Best Match

From the results, choose using: name match (prefer the exact or closest match to what was asked), benchmark score (higher indicates better documentation quality), source reputation (prefer high/medium reputation), and version (if the user specified one, prefer a version-specific library ID when listed).

### Step 3: Fetch the Documentation

Call `query-docs` with the selected library ID and the user's specific question — be specific to get relevant snippets rather than a generic dump.

Limit: no more than 3 total calls to `resolve-library-id`/`query-docs` combined per question. If the answer is still unclear after 3 calls, state the uncertainty explicitly and use the best information gathered rather than guessing further.

### Step 4: Use the Documentation

Answer using the fetched, current information; include relevant code examples from the docs; cite the library or version when it matters ("in this framework's current major version...").

## Best Practices

- Use the user's full question as the query where possible — it improves relevance ranking more than a keyword fragment.
- When a version is mentioned, prefer the version-specific library ID from the resolve step.
- Prefer official or primary packages over community forks when multiple matches exist.
- Redact API keys, passwords, tokens, and any other secret-shaped text from the query before sending it to Context7 — treat the user's question as potentially containing a secret before it leaves the session.

## Gotchas

- A library name resolving to multiple candidate IDs is common for popular libraries with several maintained forks or rewrites — the benchmark score and source-reputation signals exist specifically to break this tie; don't default to the first result.
- Calling `query-docs` without a version-specific ID when the user asked about an older major version will return current-version documentation that may directly contradict the behavior they're actually running — check the version match before trusting the answer.
- Hitting the 3-call limit without a clear answer is a signal to say so, not a signal to keep guessing across more calls — an honest "documentation is unclear on this point" beats a confident answer built on a partial fetch.

## Real-world grounding

Context7 is a real, publicly available MCP server (built by Upstash) that indexes official documentation and code examples for a large and growing set of libraries and frameworks specifically so that LLM coding assistants can retrieve current, version-aware documentation instead of relying on a frozen training-data snapshot — it's the concrete tool this skill's two-step resolve/query workflow is built around, not a generic description of "look things up."

## Verification

- [ ] `resolve-library-id` was called before any `query-docs` call
- [ ] The selected library ID matches the user's specified version, when one was given
- [ ] No more than 3 total tool calls were made before either answering or stating uncertainty
- [ ] Any secret-shaped text in the user's question was redacted before querying
