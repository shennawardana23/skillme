---
name: golang-patterns
description: Use when designing Go struct/interface APIs or optimizing allocation-heavy hot paths - functional options for constructors, defining narrow consumer-side interfaces, making zero values useful, and reducing allocations with preallocated slices, strings.Builder, and sync.Pool. For error handling, context propagation, or basic package layout use go-service-idioms instead; this skill is about API shape and allocation efficiency, not service plumbing.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# Go Struct/Interface Design and Allocation Efficiency

This skill covers two things `go-service-idioms` doesn't: how to shape Go
structs and interfaces so an API is hard to misuse, and how to cut
allocations on a hot path once a benchmark (see `golang-testing`) has
identified one. Apply the design guidance when writing new types and
constructors; apply the allocation guidance only where profiling data
justifies it — these are optimizations, not defaults to reach for
everywhere.

## Make the zero value useful

Design a type so `var t T` is immediately usable without a constructor call,
the way `bytes.Buffer` and `sync.Mutex` work.

```go
// Good: zero value works — count starts at 0, mutex starts unlocked
type Counter struct {
    mu    sync.Mutex
    count int
}

// Bad: zero value panics — nil map has no entries and can't be written to
type BadCounter struct {
    counts map[string]int
}
```

If a type genuinely cannot have a useful zero value (it needs a required
dependency), say so with an unexported field and a constructor, not a
runtime nil-check scattered across every method.

## Accept interfaces, return structs

A function should accept the narrowest interface it needs and return a
concrete type, not an interface — returning an interface hides information
the caller may need (which concrete fields exist, which extra methods are
available) for no benefit to the callee.

```go
// Good
func ProcessData(r io.Reader) (*Result, error) { ... }

// Bad — caller loses access to anything beyond io.Reader's own methods,
// and gains nothing since ProcessData already had the concrete type
func ProcessData(r io.Reader) (io.Reader, error) { ... }
```

## Define interfaces where they're consumed, not where they're implemented

The consumer package should declare the small interface it needs; the
implementing package stays unaware such an interface exists.

```go
// package service — defines exactly what it needs, nothing more
type UserStore interface {
    GetUser(id string) (*User, error)
}

type Service struct{ store UserStore }

// package postgres — implements GetUser (and much more) with no import
// of, or awareness of, the service package's interface
type Store struct{ db *sql.DB }
func (s *Store) GetUser(id string) (*User, error) { ... }
```

This keeps interfaces small and consumer-specific instead of one large
`Repository` interface every implementation must satisfy in full, and it
avoids an import cycle between the consumer and provider packages.

## Functional options for constructors with optional config

Use functional options once a constructor has more than two or three
optional parameters, instead of a config struct with exported zero-value
ambiguity (was `Timeout: 0` explicit or just unset?) or a long positional
argument list.

```go
type Server struct {
    addr    string
    timeout time.Duration
    logger  *slog.Logger
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{addr: addr, timeout: 30 * time.Second, logger: slog.Default()}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

## Embedding for method promotion, not for shortcuts

Embed a type when the outer type genuinely *is-a* superset of the inner
type's behavior (a `Server` that logs is reasonably a `*Logger` plus
networking). Don't embed purely to save writing forwarding methods for a
relationship that isn't really "is-a" — it leaks the embedded type's full
method set into the outer type's API, including methods callers shouldn't
rely on.

```go
type Logger struct{ prefix string }
func (l *Logger) Log(msg string) { fmt.Printf("[%s] %s\n", l.prefix, msg) }

type Server struct {
    *Logger
    addr string
}
// s := &Server{Logger: &Logger{prefix: "SERVER"}, addr: ":8080"}
// s.Log("starting") // promoted from *Logger
```

## Allocation efficiency (apply only where a benchmark justifies it)

**Preallocate slices with a known or estimable size** — `append` on a nil
slice reallocates and copies repeatedly as it grows:

```go
// Bad: repeated reallocation as the slice grows
var results []Result
for _, item := range items {
    results = append(results, process(item))
}

// Good: one allocation up front
results := make([]Result, 0, len(items))
```

**Use `strings.Builder` (or `strings.Join`) instead of `+=` in a loop** —
each `+=` allocates a new string, making the loop O(n²) in total bytes
copied:

```go
var sb strings.Builder
for i, p := range parts {
    if i > 0 {
        sb.WriteString(",")
    }
    sb.WriteString(p)
}
// or, when there's no per-item logic: strings.Join(parts, ",")
```

**Use `sync.Pool` for short-lived objects allocated on a hot path** —
scratch buffers, encoder state — reused across calls instead of allocated
and garbage-collected every time. Only worth the complexity when a
benchmark (`golang-testing`) shows GC pressure from that allocation site:

```go
var bufferPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func ProcessRequest(data []byte) []byte {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() { buf.Reset(); bufferPool.Put(buf) }()
    buf.Write(data)
    return buf.Bytes()
}
```

A `sync.Pool` is not a cache with guaranteed retention — the runtime is free
to evict pooled items at any GC cycle, so never rely on it for anything
beyond amortizing allocation cost.

## Recommended linter baseline

```yaml
# .golangci.yml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - unconvert
    - unparam
```

Run `gofmt -l .` / `goimports -l .` in CI to fail on unformatted code rather
than just relying on editor integration, since a missed local `gofmt` run is
otherwise invisible until review.

## Gotchas

- Loop-variable capture in a goroutine or closure (`for _, id := range ids {
  go func() { use(id) }() }`) needed an explicit `id := id` shadow before Go
  1.22 — as of 1.22 each iteration gets its own variable and the shadow is a
  no-op, not a bug. Check the module's `go` directive in `go.mod` before
  assuming which behavior applies to a given codebase.
- A large embedded interface's full method set becomes part of the outer
  type's public API whether or not that was intended — embedding for
  convenience is an API design decision, not just an implementation detail.
- `sync.Pool`'s `New` function must return a pointer type (or something
  cheap to box) — pooling non-pointer values defeats the purpose, since
  `Get()`/`Put()` traffic in `any` and a value type still allocates on the
  interface conversion.

## Real-world grounding

The standard library's own `fmt` package uses `sync.Pool` internally to
reuse the `pp` printer struct across `Printf`/`Sprintf` calls — the same
pattern recommended above, applied by the stdlib to one of the hottest call
paths in a typical Go program.
