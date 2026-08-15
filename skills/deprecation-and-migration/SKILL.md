---
name: deprecation-and-migration
description: Guides deprecating and removing old systems, APIs, or code, and migrating consumers to a replacement. Use when replacing a library, API, or internal system with a new one, sunsetting a feature, deciding whether to maintain or remove legacy code, or planning how a new system will eventually be retired.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Deprecation and Migration

Code is a liability, not an asset: every line has an ongoing cost in
tests, dependency updates, security patches, and onboarding overhead.
Deprecation is the discipline of removing code that no longer earns that
cost back. Most organizations are good at building; few are good at
removing — this skill is about the removing half.

## The deprecation decision

Answer these in order before deprecating anything:

1. **Does it still provide unique value?** If yes, maintain it; stop here.
2. **How many consumers depend on it?** Quantify with actual usage data
   (call metrics, import graph, dependency scan) — not a guess.
3. **Does a working replacement exist?** If not, build and prove the
   replacement in production first. Never deprecate without an
   alternative already available.
4. **What's the per-consumer migration cost?** Trivially automatable
   (codemod, `go fix`-style rewrite) vs. manual and high-effort changes
   the decision.
5. **What's the cost of *not* deprecating?** Security exposure, engineer
   time maintaining it, complexity tax on everything nearby.

## Compulsory vs. advisory

| Type | When | Mechanism |
|---|---|---|
| **Advisory** (default) | Old system is stable; migration optional | `// Deprecated:` comments, linter warnings, docs. Consumers migrate on their own schedule. |
| **Compulsory** | Security issue, unsustainable maintenance cost, or the old system blocks other work | Hard removal date, migration tooling provided, active tracking of remaining consumers |

Compulsory deprecation without migration tooling is just an announcement
that ships the work onto consumers — that's the "Churn Rule": whoever
owns the thing being deprecated owns migrating its users, or ships a
backward-compatible replacement that needs no migration at all.

## Go-specific mechanics

Mark deprecated exported identifiers with the standard doc-comment
convention `go vet` and IDEs recognize — this makes the deprecation
visible at every call site, not just in documentation nobody reads:

```go
// LegacyCharge processes a payment through the old gateway.
//
// Deprecated: use Charge instead; LegacyCharge will be removed after
// all callers migrate (tracked in JIRA-1234).
func LegacyCharge(amount int64) error { ... }
```

For a package whose exported API must change incompatibly, Go modules
require a new major version to live at a new import path
(`module example.com/pkg/v2`) rather than in place — this is the
mechanism that lets old and new major versions coexist in the module
graph without one silently breaking the other; see `git-workflow-and-versioning`
for how the version number itself should reflect the break.

## Migration patterns

**Strangler pattern** — shift traffic from old to new incrementally,
removing the old system only once it's fully idle:

```
Phase 1: new handles 0%,  old handles 100%
Phase 2: new handles 10%  (canary)
Phase 3: new handles 50%
Phase 4: new handles 100%, old idle
Phase 5: remove old
```

**Adapter pattern** — wrap the new implementation behind the old
interface so consumers don't need to change at all during the migration
window:

```go
// Old interface, delegating to the new implementation underneath.
type LegacyTaskService struct{ new *NewTaskService }

func (l *LegacyTaskService) GetTask(id int) (OldTask, error) {
    t, err := l.new.FindByID(context.Background(), strconv.Itoa(id))
    if err != nil {
        return OldTask{}, err
    }
    return toOldFormat(t), nil
}
```

**Feature-flag migration** — switch individual consumers over one at a
time, with an instant revert path if the new path misbehaves:

```go
func getTaskService(flags FeatureFlags, userID string) TaskService {
    if flags.Enabled("new-task-service", userID) {
        return newTaskService
    }
    return legacyTaskService
}
```

## The migration process

1. **Build the replacement** — proven in production, documented, covers
   every critical use case of the old system (not just the common path).
2. **Announce with a concrete plan** — status, replacement, migration
   guide with copy-pasteable steps, and either an advisory note or a hard
   removal date.
3. **Migrate incrementally, one consumer at a time** — identify
   touchpoints, switch, verify behavior matches, remove the old
   reference, confirm no regression. Track remaining consumers explicitly
   (a tracking issue, a dashboard of remaining call sites) rather than
   trusting memory.
4. **Remove only after verified zero usage** — check metrics/logs/import
   graph, not just "I think everyone migrated." Remove the code, its
   tests, its docs, its config, and the deprecation notice itself in the
   same change — a deprecation notice for code that no longer exists is
   its own kind of clutter.

## Zombie code

Code with active consumers but no owner, no recent commits, and no
one accountable for its failing tests or unpatched CVEs. It cannot stay
in limbo: either assign an owner and maintain it for real, or deprecate
it with the concrete plan above. Signs: 6+ months with no commits despite
active consumers; documentation referencing systems that no longer exist;
known vulnerabilities nobody has scheduled a fix for.

## Gotchas

- Hyrum's Law applies to removal, not just to APIs: with enough
  consumers, undocumented timing quirks and even bugs in the old system
  become depended-on behavior — "the replacement is functionally
  equivalent" is not sufficient; verify it against real consumer
  behavior, not just the documented contract.
- A `// Deprecated:` comment with no tracking issue and no removal
  target quietly becomes permanent — advisory deprecations that have sat
  for years with zero migration progress are a signal to either invest in
  migration tooling or accept the code is not actually going away.
- Deleting the old code before confirming zero callers via actual usage
  data (not team memory) is the single most common cause of migration
  incidents — a caller nobody remembered still existed breaks silently.
- Running two systems that do the same thing "just for a while" rarely
  stays temporary without an explicit forcing function (a removal date,
  a tracked issue) — indefinite parallel maintenance is strictly more
  expensive than either finishing the migration or reverting it.

## Real-world grounding

Python's 2-to-3 migration (Python 3 released 2008, Python 2 end-of-life
January 2020) is the widely cited cautionary case for underestimating
migration cost: the language broke backward compatibility without a
sufficiently long overlapping-support and tooling runway relative to the
size of the installed ecosystem, and the transition took over a decade
longer than planned industry-wide — the lesson generalizes directly:
compulsory deprecation without adequate migration tooling and a realistic
timeline for your actual consumer base, not an optimistic one, is what
turns a deprecation into a multi-year drag on the whole ecosystem around
it.

## Verification

- [ ] A working, production-proven replacement exists before deprecation is announced
- [ ] The deprecation type (advisory vs. compulsory) matches the actual urgency, not just convenience
- [ ] Go: deprecated exported identifiers carry a `// Deprecated:` comment; incompatible API changes get a new major-version import path
- [ ] Migration is tracked per-consumer, not assumed complete from memory
- [ ] Zero-usage is confirmed via metrics/logs/import graph before deletion
- [ ] Old code, its tests, its docs, its config, and its deprecation notice are all removed in the same change
