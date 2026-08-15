---
name: twelve-factor-app-principles
description: Use when the user asks to "check twelve-factor compliance", "review this service for 12-factor", "design a new microservice", "where should config live", "should this write to a local file", or is scaffolding a new backend service, containerizing an app, or debugging "works on my machine but not in prod" issues. Guides applying the publicly documented Twelve-Factor App methodology (codebase, config, backing services, statelessness, logs, dev/prod parity) to a service under review or in design.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Twelve-Factor App Principles

The Twelve-Factor App (12factor.net) is a methodology for building
services that can be deployed, scaled, and handed off between engineers
without hidden local state or environment-specific assumptions. This skill
paraphrases its own principles — for the canonical text, see the source
website rather than copying prose from here.

Rather than restate all twelve one by one, group them by the failure mode
each one prevents — this is more useful for review than memorizing a list.

## Group 1: One codebase, config and environment kept out of code

- **One codebase tracked in version control, many deploys from it.** Every
  running instance — local, staging, production — traces back to the same
  repository at different points in its history; if two environments are
  running code that lives in genuinely different repositories (a fork, a
  copy-pasted variant), you no longer have one deployable app, you have
  two apps pretending to be one. A service should never require its own
  divergent copy of the repo per environment.
- **Config in the environment, not in code.** Anything that varies between
  deploys (database URLs, API keys, feature flags, hostnames) belongs in
  environment variables or a secrets manager, never hardcoded or committed
  in a config file checked into version control. A quick test: could you
  make the repository public without leaking a secret or an
  environment-specific value? If not, something in "config" is misplaced.
- **Dev/prod parity.** Keep the gap between local, staging, and production
  as small as possible — same backing service *types* (don't run SQLite
  locally and Postgres in prod), same dependency versions, short lead time
  between a code change and its deploy. Most "works on my machine" bugs
  trace back to a parity gap here.

## Group 2: Backing services are attached resources, not internal state

- **Every external system a service talks to** — its database, a queue, a
  cache, a third-party API — **should be swappable purely by changing a
  config value**, with no code change required. A service shouldn't need
  a code change to point at a different Postgres instance — only a changed
  connection string. This makes it easy to point a local dev environment
  at a throwaway instance without touching the code that talks to it.
- **Explicit, declared dependencies.** Every library the service needs is
  declared in a manifest (`go.mod`, `package.json`) and pinned — never
  relying on something happening to be present on the host system.

## Group 3: Processes are stateless and disposable

- **Stateless processes.** Nothing that must survive a restart lives in
  process memory or on local disk — session state, uploaded files, and
  computed caches belong in a backing store (database, object storage,
  Redis) so any process can be killed and replaced without losing data.
- **Fast startup, graceful shutdown.** Processes should start in seconds
  and shut down cleanly on a termination signal, finishing in-flight work
  or returning it to a queue — this is what makes horizontal scaling and
  rolling deploys safe.
- **Concurrency via the process model.** Scale out by running more
  process instances (horizontally), not by making a single process more
  complex internally to handle more load.

## Group 4: Ports, build/release/run, logs as streams

- **Port binding / self-contained service.** The app exports its
  functionality by binding to a port itself (not requiring injection into
  a heavier external webserver container to become network-reachable) —
  this is what makes a service directly runnable and directly composable
  behind any router.
- **Strict separation of build, release, and run stages.** A build
  produces an immutable artifact; a release combines that artifact with a
  specific config; a run stage executes a release. You should never patch
  code directly in the running stage — every fix flows back through build.
- **Logs as an unbuffered event stream to stdout/stderr**, not written to
  or rotated by the app itself. Let the execution environment (container
  runtime, log aggregator) handle routing, storage, and rotation — an app
  that manages its own log files is coupling itself to a specific
  deployment environment.
- **Admin/one-off tasks (migrations, REPL, one-off scripts) run in an
  identical environment** to the long-running process (same code, same
  config), not from a different snapshot of the codebase that can drift out
  of sync.

## Procedure: reviewing a service for 12-factor compliance

1. Grep for hardcoded connection strings, API keys, or hostnames in source
   — anything that should be config but isn't.
2. Check whether the service writes anything to local disk that must
   survive a restart (uploaded files, session data, sqlite files used as
   real storage) — if so, that's state that needs to move to a backing
   service.
3. Check log output: does the app write straight to stdout, or does it
   manage its own log files/rotation?
4. Check whether a second instance of the process could run concurrently
   against the same backing services without corrupting shared state —
   if not, it isn't stateless yet.
5. Check whether `docker build` (or equivalent) plus a config change is
   enough to produce every environment's release, or whether some
   environments require hand-edited files post-build.

## Gotchas

- Config that's environment-*specific* but happens to be the same value
  across all environments today (e.g., a feature flag currently `true`
  everywhere) still belongs in config, not a hardcoded constant — "same
  value everywhere right now" is not the test; "could legitimately differ
  between environments" is.
- A service that writes structured JSON logs to stdout still satisfies the
  logs-as-streams factor even though it's not human-readable line-by-line —
  the requirement is about *where* logs go (a stream, unbuffered), not
  their format.
- Local caching in application memory (e.g., an in-process LRU cache) is a
  gray area: acceptable as a pure performance optimization the process can
  lose on restart without correctness impact, but a violation the moment
  any code path depends on that cache being warm or present for correct
  behavior.
- "Backing service" doesn't only mean a database — third-party APIs
  (payment gateways, email providers) and even another internal service
  count, and should be equally swappable via config/URL rather than
  hardcoded per environment.

## Real-world grounding

The Twelve-Factor App methodology was published by Heroku engineers
(originally largely attributed to Adam Wiggins) as a distillation of
patterns observed across many SaaS applications deployed on Heroku's
platform — it's publicly documented at https://12factor.net/ and remains
the most commonly cited reference for what makes a service
cloud-native/PaaS-portable, cited directly in Kubernetes and container
platform design discussions well beyond Heroku itself.

## Verification

- [ ] No secrets or environment-specific values are hardcoded in source
- [ ] The service is swappable to different backing service instances via
      config alone
- [ ] Nothing required for correctness lives only in process memory or
      local disk across a restart
- [ ] Logs go to stdout/stderr, not app-managed files
- [ ] A single build artifact plus config produces every environment's release
