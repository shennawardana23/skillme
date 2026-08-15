---
name: api-and-interface-design
description: Guides internal interface and module-boundary design — Go interfaces, package APIs, and cross-team contracts inside a codebase, as distinct from the wire-level REST/RPC contract shape covered by the api-design skill. Use when defining an interface between packages, deciding where an interface should live, choosing between accepting an interface or a concrete type, evolving a public function signature without breaking callers, or reviewing a Go package's exported API for coupling and future flexibility.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# API and Interface Design

This skill covers **in-process boundaries**: interfaces between packages,
modules, or teams inside a single codebase or binary. For the shape of a
**wire-level contract** (REST/RPC request/response bodies, error envelopes,
pagination, versioning of an HTTP or RPC API), use the `api-design` skill
instead — the two are complementary: a service's HTTP handler is itself a
consumer of the internal interfaces this skill is about.

## Core rule: define the interface where it's consumed, not where it's implemented

The single most common Go interface mistake is declaring an interface next
to its implementation ("I built a `Store`, so I export a `Store` interface
next to it"). This forces every consumer to depend on the producer's
package just to get the interface type, and it tempts the producer into a
single fat interface covering everything the type does.

```go
// BAD: interface lives in the producer package, sized to the whole type.
package store

type Store interface {
    GetReservation(ctx context.Context, id string) (*Reservation, error)
    SaveReservation(ctx context.Context, r *Reservation) error
    GetGuest(ctx context.Context, id string) (*Guest, error)
    SaveGuest(ctx context.Context, g *Guest) error
    // ... ten more methods, because that's everything *SQLStore does
}

// GOOD: consumer package declares only the interface it needs.
package billing

// ReservationGetter is the only dependency billing has on storage.
type ReservationGetter interface {
    GetReservation(ctx context.Context, id string) (*Reservation, error)
}

func NewInvoicer(store ReservationGetter) *Invoicer { ... }
```

The producer (`store.SQLStore`) never needs to declare an interface at
all — it just needs to satisfy whatever its consumers declare, which Go's
structural typing does implicitly. This is the direct analogue of Go's
own guidance ("accept interfaces, return structs"): a constructor should
return a concrete type so callers get full access to it, and take
interfaces as parameters so the function only asks for what it uses.

## Core rules

1. **Small interfaces, defined by consumers.** One or two methods is the
   norm (`io.Reader`, `sort.Interface`). A ten-method interface is a sign
   the consumer actually depends on a concrete implementation and should
   just take one, or the interface should be split per-consumer.
2. **Accept interfaces, return structs.** Functions and constructors take
   the narrowest interface they need as a parameter, and return a
   concrete type. Returning an interface hides the concrete type's own
   methods from every caller, including ones that could safely use them.
3. **Extend via new, optional parameters — never widen an existing
   signature's meaning.** Adding a field to a struct passed by value, or
   a new method to an interface nobody outside this package implements,
   is safe. Adding a required parameter to an exported function, or a
   method to an interface implemented by external packages, breaks every
   implementer — use the functional-options pattern instead:
   ```go
   type Option func(*clientConfig)

   func WithTimeout(d time.Duration) Option {
       return func(c *clientConfig) { c.timeout = d }
   }

   // NewClient(url) still works; new callers add options without
   // breaking anyone who already calls NewClient(url).
   func NewClient(url string, opts ...Option) *Client { ... }
   ```
4. **One error-handling contract per package boundary.** Callers must be
   able to use `errors.Is`/`errors.As` against sentinel or typed errors
   you export (`var ErrNotFound = errors.New(...)`), not string-match
   `err.Error()`. Wrap with `%w`, never `%v`, so the chain survives.
5. **The One-Version Rule.** Don't make consumers choose between two
   incompatible versions of the same interface at once (a `v1.Store` and
   a `v2.Store` implemented by different types with overlapping callers).
   Prefer extending the existing interface additively, or migrating all
   consumers to the new one in one coordinated change — see
   `deprecation-and-migration` for the process when the old interface
   must eventually go.
6. **Validate at the boundary where untrusted input enters** (HTTP
   handler, CLI flag parsing, config file load) — internal function
   calls receiving already-validated types should not re-validate; see
   `api-design` for the full boundary-validation rule at the wire level.

## Discriminated unions without a language feature

Go has no sum types, but the same "make illegal states unrepresentable"
goal is achievable with an interface plus a private marker method, giving
consumers exhaustive-style handling via a type switch:

```go
type ReservationState interface{ isReservationState() }

type Pending struct{}
type Confirmed struct{ ConfirmedAt time.Time }
type Cancelled struct{ Reason string }

func (Pending) isReservationState()   {}
func (Confirmed) isReservationState() {}
func (Cancelled) isReservationState() {}

func describe(s ReservationState) string {
    switch v := s.(type) {
    case Pending:
        return "pending"
    case Confirmed:
        return "confirmed at " + v.ConfirmedAt.String()
    case Cancelled:
        return "cancelled: " + v.Reason
    default:
        panic(fmt.Sprintf("unhandled state %T", v)) // catches a future added case
    }
}
```

Because `isReservationState` is unexported, no package outside this one
can add a new implementation — the switch above really is exhaustive
against every case that exists today.

## Gotchas

- A one-method interface named after the method (`Reader`, `Closer`,
  `ReservationGetter`) reads clearly at the call site; a interface named
  after the producer type (`StoreInterface`) is a sign it was designed
  producer-first and is probably too big.
- `interface{}` (or `any`) as a parameter type isn't "flexible design," it
  deletes the compiler's ability to catch a caller passing the wrong
  thing — prefer a small concrete interface or a generic type parameter.
- Adding a method to an interface you export and that external packages
  implement is a breaking change even though nothing in your own repo
  fails to compile — search for `var _ YourInterface = (*T)(nil)` style
  assertions and any known external implementers before changing it.
- Table-driven tests that construct a fake by hand-implementing a small
  interface are cheap; a fake implementing a twelve-method interface is
  itself a maintenance burden that argues for shrinking the interface.
- Returning an interface from a constructor (`func New() Reader`) instead
  of the concrete type blocks every caller from later using
  `errors.As(err, &myConcreteErr)` style type assertions or calling a
  method that isn't on the interface, even if they need to.

## Real-world grounding

The Go 1 compatibility promise (in effect since Go 1.0, 2012) is the
industry's most disciplined public demonstration of additive-only
interface evolution: the standard library has added methods, types, and
functions for over a decade while guaranteeing "existing programs will
continue to run correctly, unchanged, over the lifetime of that version" —
achieved specifically by never changing an existing exported signature's
meaning and by growing interfaces only where it could prove no external
implementer would break (an audit only the Go team could realistically do,
which is why in your own code the safer default is to add a new
interface or new options rather than widen an old one).

## Verification

- [ ] Interfaces are declared in the consuming package, sized to what that consumer calls
- [ ] Constructors return concrete types; parameters accept the narrowest interface needed
- [ ] No exported function signature or externally-implemented interface changed in place — extension used functional options or a new type instead
- [ ] Errors are wrapped with `%w` and exposed as sentinels/typed errors consumers can `errors.Is`/`errors.As` against
- [ ] No two incompatible versions of the same interface are live across consumers at once (One-Version Rule)
- [ ] Wire-level contract concerns (HTTP shape, pagination, error envelope) are handled per the `api-design` skill, not duplicated here
