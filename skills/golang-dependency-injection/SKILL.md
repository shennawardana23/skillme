---
name: golang-dependency-injection
description: This skill should be used when the user asks to "wire up dependencies in Go without a framework", "should I use wire, dig, fx, or samber/do", "generate a Go dependency injection container", "add a lifecycle hook to an fx app", "why does my dig container fail to invoke", or is choosing between or implementing google/wire, uber-go/dig, uber-go/fx, or samber/do for constructing a Go application's dependency graph. go-service-idioms's default is manual constructor functions wired by hand in main() — reach for this skill once that manual wiring gets unwieldy enough to justify a tool, and to pick the right one.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Go Dependency Injection

Manual constructor wiring in `main()` — call each constructor in
dependency order, pass the results down — is the correct default for
most Go services; it's explicit, has zero magic, and any Go developer can
read it without knowing a DI library. Reach for a tool only once that
wiring genuinely becomes hard to follow (dozens of components, several
optional/environment-specific variants) — and then pick based on the
table below, not on familiarity alone.

## Picking one

| Tool | Mechanism | Errors surface | Reach for it when |
|---|---|---|---|
| `google/wire` | Compile-time code generation — `wire.Build` inside a `//go:build wireinject` injector function produces a real `wire_gen.go` | Build time (missing/ambiguous provider is a codegen failure, not a runtime panic) | You want DI with zero runtime reflection cost and errors caught before the binary even runs |
| `uber-go/dig` | Runtime reflection — a `*dig.Container` resolves constructor arguments by type | Runtime, at `Invoke`/`Provide` call time | You want a lightweight container without wire's code-generation step, and runtime resolution cost is acceptable |
| `uber-go/fx` | Built on `dig`, adds application lifecycle (`OnStart`/`OnStop` hooks) and module composition | Runtime, at `fx.New`/app start | You're building a long-running service that needs coordinated startup/shutdown order across many components, not just object construction |
| `samber/do` (v2) | Runtime, generics-based (`Provide[T]`, `Invoke[T]`) instead of reflection over `interface{}` | Runtime, at `Invoke`/`MustInvoke` call time | You want dig/fx-like runtime DI but with compile-time type safety on the calling side via generics, and optionally its health-check/shutdown hooks |

Wire and dig/fx/do split along one axis: **compile-time (wire) vs.
runtime (dig, fx, do)**. Dig, fx, and do split along a second axis: dig
is a bare container, fx adds lifecycle management on top of dig, and
`samber/do` trades dig's `interface{}`-based reflection for Go generics
so `Invoke[MyService](injector)` is type-checked at compile time even
though resolution itself still happens at runtime.

## google/wire

```go
//go:build wireinject

func InitializeApp() (*App, error) {
    wire.Build(NewConfig, NewDB, NewStore, NewApp)
    return nil, nil // never actually runs; wire replaces this file
}
```

`wire.Build(providers...)` takes function values (each contributing its
first return type to the graph) and `wire.NewSet` groups reusable sets of
them. Running `wire` generates `wire_gen.go` with the real
`InitializeApp` implementation — the checked-in generated file, not the
injector template, is what actually ships. `wire.Bind(new(Iface),
new(*Impl))` maps an interface a constructor depends on to the concrete
type that satisfies it.

## uber-go/dig and uber-go/fx

```go
c := dig.New()
c.Provide(NewConfig)
c.Provide(NewDB)
c.Provide(NewStore)
if err := c.Invoke(func(s *Store) { s.Run() }); err != nil {
    log.Fatal(err)
}
```

`fx` wraps this in an `*fx.App` and adds `fx.Lifecycle`: a constructor
that needs to start/stop background work registers a `fx.Hook` via
`lifecycle.Append(fx.Hook{OnStart: ..., OnStop: ...})` instead of
managing goroutines by hand in `main()`:

```go
fx.New(
    fx.Provide(NewConfig, NewDB, NewServer),
    fx.Invoke(func(lc fx.Lifecycle, srv *Server) {
        lc.Append(fx.Hook{
            OnStart: func(ctx context.Context) error { return srv.Start(ctx) },
            OnStop:  func(ctx context.Context) error { return srv.Stop(ctx) },
        })
    }),
).Run()
```

`dig.In`/`dig.Out` structs let a constructor take or return several
values as named struct fields instead of a long parameter/return list;
`dig.Name`/`dig.Group` disambiguate or collect multiple values of the
same type.

## samber/do (v2)

```go
injector := do.New()
do.Provide(injector, func(i do.Injector) (*DB, error) { return NewDB(cfg) })
do.Provide(injector, func(i do.Injector) (*Store, error) {
    db := do.MustInvoke[*DB](i)
    return NewStore(db), nil
})
store := do.MustInvoke[*Store](injector)
```

`do.Provide[T]` registers a lazy constructor for `T`; `do.Invoke[T]`/
`do.MustInvoke[T]` resolve it, instantiating on first use. `do.Shutdown[T]`
and `do.HealthCheck[T]` give per-service lifecycle hooks similar in spirit
to fx's, but addressed by type via generics rather than fx's
lifecycle-object pattern.

## Gotchas

- Wire's injector function body is never actually executed — it exists
  only as a template `wire` reads to generate `wire_gen.go`. Editing the
  generated file by hand instead of re-running `wire` after changing the
  injector is a common mistake that silently drifts the two out of sync.
- A `dig`/`fx`/`do` missing-provider error only surfaces the first time
  that specific type is resolved — a container can build successfully
  and only fail deep into a request path the first time a rarely-used
  constructor is invoked, unlike wire's build-time failure for the same
  mistake.
- `fx.Invoke` functions run eagerly at app start in the order they're
  registered — a constructor with expensive or side-effecting logic
  should go in `fx.Provide` (lazy, resolved on demand) unless it
  genuinely needs to run unconditionally at startup.
- Value groups in `dig` (`dig.Group`) and named values both disambiguate
  "multiple things of the same type," but for different reasons —
  groups collect an unordered set (e.g., all registered HTTP
  middlewares), named values pick one specific instance among several
  (e.g., a "primary" vs. "replica" `*sql.DB`). Using named values where a
  group was meant loses the rest of the set silently.
- None of these four tools replace the "accept interfaces, return
  concrete structs" guidance in `go-service-idioms` — they automate
  *wiring* constructors together, not designing what those constructors'
  signatures should look like.

## Real-world grounding

`wire.Build`/`wire.NewSet`/`wire.Bind`, `dig.Container.Provide`/`Invoke`,
`fx.Lifecycle`/`fx.Hook`, and `do.Provide[T]`/`do.MustInvoke[T]` are
documented directly on pkg.go.dev under `github.com/google/wire`,
`go.uber.org/dig`, `go.uber.org/fx`, and `github.com/samber/do/v2`
respectively.

## Cross-references

- For the manual-wiring default and general package-layout guidance
  these tools automate around, see `go-service-idioms`.
- For lifecycle-hook goroutines (`OnStart`/`OnStop`) that themselves need
  correct cancellation, see `golang-concurrency` and `golang-context`.
