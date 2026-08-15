---
name: scaffold
description: Generates a complete, compilable project scaffold -- directory structure, entry point, config, Makefile, README -- for a new Go (or other language) service, following clean architecture and dependency injection instead of global state. Use when starting a brand-new service or CLI from nothing, not when adding a feature to an existing project.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Project Scaffold

Generate a new service's skeleton so it's immediately buildable, runnable,
and testable — not a pile of TODOs. Every file produced must compile and
every mandatory file must exist; a scaffold that needs the user to fill in
missing pieces before `go build` succeeds has failed at the one thing a
scaffold is for.

## When to use

Starting a genuinely new service, CLI, or module with no existing code to
extend. Not for adding a feature or file to a project that already
exists — that's ordinary implementation work, not scaffolding.

## Phase 1: Requirements

Extract, or ask for if missing:

1. **Project name and purpose** — one line.
2. **Tech stack** — language, framework, database, external services. This
   catalog defaults to Go/PostgreSQL unless told otherwise.
3. **Key features** the scaffold must include.
4. **Architecture pattern** — clean/hexagonal layering, or a flatter
   structure for something genuinely small.
5. **Constraints** — team conventions, existing infrastructure this must
   fit into.

If critical information is missing (especially the tech stack or
database), ask one focused question rather than guessing — a scaffold
built on a wrong assumption about the database gets thrown away, not
edited.

## Phase 2: Plan the layout before generating files

For a Go service:

```
cmd/<name>/main.go            → signal-aware entry point, graceful shutdown
internal/config/config.go     → env-var config with defaults
internal/domain/<entity>.go   → core types, no framework dependency
internal/service/<service>.go → business logic
internal/repository/<repo>.go → persistence, behind an interface
internal/handler/<handler>.go → HTTP/RPC boundary
go.mod
Makefile
README.md
.env.example
```

State the planned layout before generating file contents — a plan the user
can glance at and redirect is cheaper to correct than a full tree of
generated files that turns out to be the wrong shape.

## Phase 3: Generate

For every file:

- Complete, compilable content — no `// TODO: implement` in a file the
  scaffold claims is done.
- Correct `package` declaration and only the imports actually used.
- Doc comments on exported symbols (Go: `// Foo does X.` above every
  exported func/type).
- Dependency injection, not global state — a `New(cfg Config, db *sql.DB)
  *Service` constructor, not a package-level `var db *sql.DB` other files
  reach into.
- Errors returned and wrapped with context, not swallowed or panicked —
  `panic` is acceptable only in `main`'s own startup init, never inside
  request-handling or service code.
- `context.Context` threaded through anything that does I/O.

**Mandatory files** for a Go service scaffold:

1. `cmd/<name>/main.go` — handles `SIGINT`/`SIGTERM` and shuts down
   gracefully (drains in-flight requests) rather than exiting immediately.
2. `internal/config/config.go` — reads env vars with sane defaults, fails
   fast (returns an error, not a panic) if a required var is missing.
3. `Makefile` — at minimum `make build`, `make test`, `make run`,
   `make lint`.
4. `README.md` — setup, run, test, and deploy instructions a new
   contributor can follow without asking anyone.
5. `.env.example` — every env var the config reads, documented, with no
   real secret values.
6. `go.mod` with pinned dependency versions (or the language's equivalent
   manifest).

## Phase 4: Review before handing it back

- Every import is used, and every imported package is actually a
  dependency in the manifest — not a package that happens to exist in the
  module cache from another project.
- No `TODO` left in a file the scaffold presents as complete.
- No hardcoded secrets, and no hardcoded `localhost` addresses that should
  be config-driven instead.
- The entry point handles `SIGINT`/`SIGTERM`.
- `make build` (or the equivalent) actually succeeds against the
  generated tree — verify this, don't assume it from having written
  syntactically plausible code.

## Gotchas

- A scaffold with global `var db *sql.DB` "for convenience" bakes untested
  code into every file that touches it from day one — every consumer
  becomes hard to unit-test without a real database, which is the opposite
  of what a good scaffold should set up.
- Godoc comments generated for internal (unexported) helpers add noise
  without adding the API-contract value a doc comment provides on an
  exported symbol — reserve them for what's actually exported.
- A `.env.example` that's missing a var the config code actually reads is
  worse than no `.env.example` — a new contributor trusts it completely
  and then hits a confusing "missing required env var" failure the file
  promised wouldn't happen.
- Claiming "make build succeeds" without having run it is the single most
  common way a scaffold ships with a typo'd import path or an unused
  import that fails `go vet`.

## Real-world grounding

`golang-standards/project-layout` on GitHub — one of the most-starred Go
repositories, despite explicitly not being an official Go team
document — became the de facto reference for the `cmd/`, `internal/`,
`pkg/` split used above precisely because Go's toolchain itself enforces
`internal/`'s import boundary at compile time: code outside the module
cannot import an `internal/` package, which is why "keep business logic
in `internal/`, unless intentionally exposing a public API via `pkg/`" is
a load-bearing convention rather than a stylistic preference.

## Verification

- [ ] Every mandatory file exists (entry point, config, Makefile, README,
      env example, dependency manifest)
- [ ] The generated tree actually builds — verified, not assumed
- [ ] No TODO comments in files presented as complete
- [ ] No hardcoded secrets or localhost addresses
- [ ] The entry point shuts down gracefully on SIGINT/SIGTERM
- [ ] Business logic is injected via constructors, not global state
