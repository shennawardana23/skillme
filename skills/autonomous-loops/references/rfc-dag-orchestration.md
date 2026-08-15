# RFC-driven multi-agent DAG orchestration

The most sophisticated pattern in `autonomous-loops`: decompose a written
spec/RFC into a dependency DAG of work units, run each unit through a
tiered quality pipeline in its own isolated workspace, and land units
through a merge queue that evicts and retries on conflict. This pattern
is illustrative — build the pieces that fit your own tooling rather than
assuming any specific named implementation is already present.

## Architecture overview

```
RFC / spec document
      |
      v
DECOMPOSITION (agent reads the RFC, proposes work units + dependency DAG)
      |
      v
For each DAG layer, in dependency order:
  - run each unit's quality pipeline in parallel, isolated workspaces
  - land completed units through a merge queue
  - evicted units re-enter the next layer with conflict context attached
```

## RFC decomposition

A decomposition pass over the RFC should produce, per work unit:

- a stable identifier and human-readable name
- which RFC sections it addresses
- its dependencies (other unit IDs it must land after)
- concrete acceptance criteria
- a complexity tier (see below)

Decomposition rules that keep merge risk low:

- Prefer fewer, more cohesive units over many tiny ones — more units
  means more merge coordination.
- Minimize file overlap between units; overlapping units must land
  sequentially, which serializes work that was supposed to be parallel.
- Keep implementation and its tests in the same unit — never split
  "implement X" and "test X" into separate units landing independently.
- Only declare a dependency where a real code dependency exists; a false
  dependency serializes work that could have run in parallel.

The DAG determines execution order — units with no unmet dependencies run
in parallel within a layer; a unit only starts once everything it depends
on has landed.

## Complexity tiers

Route pipeline depth by how much can go wrong, not by how big the diff
looks:

| Tier | Pipeline |
|---|---|
| Trivial | implement -> test |
| Small | implement -> test -> code review |
| Medium | research -> plan -> implement -> test -> review -> fix |
| Large | research -> plan -> implement -> test -> review -> fix -> final review |

This keeps expensive research/review stages off simple mechanical
changes while giving architecturally risky units the full pipeline.

## Author-bias elimination: separate context per stage

Each pipeline stage runs as its own agent invocation with its own context
window — research, plan, implement, test, and review are never the same
context. The critical property: **the reviewer never wrote the code it's
reviewing.** This is the same reasoning behind requiring an independent
reviewer for human-written code — a single reasoning process reviewing
its own output shares its own blind spots with the thing it's checking
(see `ai-regression-testing` for the same principle applied to
self-review of bug fixes).

## Merge queue with eviction

```
unit branch
  -> rebase onto main
       conflict? -> evict, capture conflict context (diffs, conflicting files)
  -> run build + tests
       fail? -> evict, capture test failure output
  -> pass -> land, delete branch
```

- Units with no file overlap can land speculatively in parallel.
- Units that do overlap land one at a time, rebasing after each landing.
- An evicted unit's next attempt gets the full conflict context (which
  files conflicted, what the diff looked like, what tests failed) fed
  back into its implementer stage — a blind retry with no new
  information just reproduces the same conflict.

This eviction protocol is conceptually the same discipline behind a
"merge queue" / "not rocket science rule" CI gate used in large open
source projects (e.g. the `bors` merge bot pattern): nothing lands on the
main line without passing the full check suite *after* rebasing onto the
current main line, so main never goes red from a stale merge.

## Isolation between units

Each unit should run in its own isolated workspace (a git worktree, a
container, or an equivalent), so parallel units can't clobber each
other's working directory state. All pipeline stages for the *same* unit
share that one workspace, so state built up in research/plan (context
files, partial notes) is available to the later implement/test/review
stages of that same unit.

## When this pattern is (and isn't) worth it

| Signal | Use the DAG pattern | Use a simpler pattern instead |
|---|---|---|
| Multiple interdependent work units | Yes | No |
| Real parallel-implementation need | Yes | No |
| Merge conflicts plausible between units | Yes | No — sequential is fine |
| Single-file change | No | Yes — sequential pipeline |
| Quick iteration on one thing | No | Yes — persistent session loop |

The coordination overhead (DAG planning, per-unit isolated workspaces, a
merge queue with eviction) only pays for itself when the alternative is
genuinely serializing work that could otherwise run in parallel, or when
merge conflicts between concurrent units are a real risk rather than a
theoretical one.
