---
name: golang-grpc
description: This skill should be used when the user asks to "write a gRPC service in Go", "add a gRPC interceptor", "return a gRPC status error", "stream data over gRPC in Go", "chain gRPC unary interceptors", or writes or reviews Go code using google.golang.org/grpc — server/client setup, interceptors, streaming, deadlines, or status codes. Distinct from golang-context (general context composition) and jwt-tenant-scoped-authorization (HTTP-oriented tenant auth) — this skill is specifically about gRPC's own service/interceptor/status-code mechanics.
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Go gRPC Service Patterns

gRPC-go's server and client are both built around interceptors (its
middleware mechanism) and a fixed vocabulary of status codes instead of
HTTP status codes — the two things worth getting right before writing
service logic.

## Server and client setup

```go
srv := grpc.NewServer(
    grpc.ChainUnaryInterceptor(loggingInterceptor, authInterceptor),
    grpc.ChainStreamInterceptor(loggingStreamInterceptor),
)
pb.RegisterGuestServiceServer(srv, &guestServer{})
```

Use `grpc.NewClient(target, opts...)` for new client code — `grpc.Dial`
and `grpc.DialContext` are deprecated in favor of it. `NewClient` does not
perform I/O immediately; the connection is established lazily on first
use unless you explicitly wait for it ready.

## Interceptors: gRPC's middleware

A unary interceptor wraps every unary RPC; a stream interceptor wraps
every streaming RPC. Chain them with `ChainUnaryInterceptor`/
`ChainStreamInterceptor` — **the first interceptor listed is outermost**,
so ordering matters (a request-logging interceptor typically goes first
so it wraps everything, including auth failures further in):

```go
func authInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler) (any, error) {
    ctx, err := authenticate(ctx)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid credentials")
    }
    return handler(ctx, req)
}
```

Every interceptor must call `handler(ctx, req)` (or its stream
equivalent) exactly once on the success path — an interceptor that
returns early without calling `handler` silently short-circuits the RPC
for reasons that are easy to lose track of once several interceptors are
chained.

## Status codes, not HTTP codes

Return errors via `status.Error(code, msg)` or `status.Errorf(code, fmt,
args...)`, using a `codes.Code` from `google.golang.org/grpc/codes` —
`codes.NotFound`, `codes.InvalidArgument`, `codes.Unauthenticated`,
`codes.PermissionDenied`, `codes.AlreadyExists`, `codes.DeadlineExceeded`.
A plain Go `error` returned from a handler becomes `codes.Unknown` on the
wire, which loses the caller's ability to branch on what actually went
wrong — always pick a specific code for expected failure conditions.

```go
if guest == nil {
    return nil, status.Errorf(codes.NotFound, "guest %s not found", id)
}
```

Attach structured details beyond a message with `status.WithDetails` when
the client needs machine-readable data (e.g., which field failed
validation), not just a human-readable string embedded in the message.

## Deadlines propagate automatically — don't re-invent them

A client-set deadline (`grpc.NewClient` request context with
`context.WithTimeout`) travels with the RPC and is available server-side
as the handler's `ctx`'s own deadline — there's no separate "gRPC
timeout" concept to configure beyond the context passed in. Check
`ctx.Err()` in long-running handlers and streaming loops exactly as
`golang-context` describes; gRPC gives you nothing extra here beyond a
context that's already correctly wired.

## Streaming

A `ServerStream`/`ClientStream` sends multiple messages over one RPC.
Every `Send`/`Recv` pair must handle `io.EOF` on `Recv` as "the other side
finished sending," not as an error to propagate:

```go
for {
    msg, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        return status.Errorf(codes.Internal, "recv: %v", err)
    }
    process(msg)
}
```

A server-streaming handler must keep checking `stream.Context().Err()`
(equivalent to a unary handler's `ctx.Err()`) inside long send loops so a
client that disconnects doesn't leave the server writing into the void.

## Gotchas

- `ChainUnaryInterceptor`'s ordering is outermost-first, not
  innermost-first — a common mistake is assuming the *last* interceptor
  in the list runs first; it's the opposite, matching how you'd nest them
  by hand.
- Returning a plain `errors.New("not found")` instead of
  `status.Error(codes.NotFound, ...)` makes the client see `codes.Unknown`
  regardless of the message text — clients that branch on `status.Code(err)`
  (the correct pattern) get nothing useful to branch on.
- `grpc.Dial`/`DialContext` are deprecated but still widely seen in
  existing code and examples — prefer `grpc.NewClient` in new code, and
  don't assume a tutorial using `Dial` reflects current guidance.
- Forgetting to check `err == io.EOF` on a stream `Recv()` loop and
  instead treating it as a real error turns every normal stream
  completion into a logged failure.
- An interceptor that panics without recovering takes down the whole gRPC
  server process for every in-flight RPC, not just the one that
  triggered it — install a recovery interceptor
  (`google.golang.org/grpc` has no built-in one; use
  `grpc_recovery.UnaryServerInterceptor` from
  `github.com/grpc-ecosystem/go-grpc-middleware` or a small hand-rolled
  equivalent) ahead of business-logic interceptors in the chain.

## Real-world grounding

`grpc.NewClient` superseding `Dial`/`DialContext`, `ChainUnaryInterceptor`'s
outermost-first ordering, and the `codes.Code` vocabulary are documented
directly on pkg.go.dev under `google.golang.org/grpc` and
`google.golang.org/grpc/codes`.

## Cross-references

- For the context/deadline mechanics an RPC's `ctx` follows once inside a
  handler, see `golang-context`.
- For HTTP-oriented tenant-scoped authorization, see
  `jwt-tenant-scoped-authorization` — gRPC auth is enforced via
  interceptors instead, following the pattern above.
- For general Go service structure (package layout, error wrapping)
  outside the gRPC-specific pieces, see `go-service-idioms`.
