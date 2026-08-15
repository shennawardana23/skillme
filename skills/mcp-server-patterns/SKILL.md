---
name: mcp-server-patterns
description: Guides building Model Context Protocol (MCP) servers — tools, resources, prompts, transport choice, and schema validation — biased toward the official Go SDK with TypeScript SDK notes where it differs. Use when implementing a new MCP server, adding a tool or resource to an existing one, choosing between stdio and HTTP transport, or debugging MCP registration and transport issues.
license: Apache-2.0
metadata:
  version: "0.1.0"
---

# MCP Server Patterns

The Model Context Protocol (MCP) lets an AI assistant call tools, read
resources, and use prompts exposed by your server. The spec and SDKs
evolve; before writing registration code, check Context7 (query "MCP")
or https://modelcontextprotocol.io for the current method signatures —
this skill covers concepts and structure that outlive any one SDK
version, not exact API calls.

## Core concepts

- **Tools** — actions the model can invoke (run a search, execute a
  command, mutate state). Always have a side effect budget in mind:
  read-only tools are safe to retry; mutating tools need idempotency.
- **Resources** — read-only data the model can fetch by URI (file
  contents, a database row, an API response). No side effects, ever — if
  it mutates anything, it's a tool, not a resource.
- **Prompts** — reusable, parameterized templates a client can surface to
  the user (e.g., in Claude Desktop's prompt picker).
- **Transport** — `stdio` for local clients spawning your server as a
  subprocess (Claude Desktop, Claude Code); **Streamable HTTP** for
  remote clients (a hosted server multiple users connect to). Legacy
  HTTP+SSE exists only for backward compatibility with pre-2025-03-26
  clients.

Keep tool/resource/prompt logic independent of transport — construct the
server's capabilities once, then plug in whichever transport the
entrypoint needs. This is the same boundary discipline as
`api-and-interface-design`: the transport is a delivery mechanism, not
part of the domain logic.

## Go server skeleton

The official Go SDK (`github.com/modelcontextprotocol/go-sdk`) is this
organization's default for new MCP servers, matching the Go-first
policy. Verify exact method names against the SDK's current release —
the shape below is illustrative of the structure, not a guaranteed
current signature:

```go
package main

import (
    "context"
    "log"

    "github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchInput struct {
    Query string `json:"query" jsonschema:"the search query"`
    Limit int    `json:"limit,omitempty" jsonschema:"max results, default 10"`
}

func searchHotels(ctx context.Context, req *mcp.CallToolRequest, in SearchInput) (*mcp.CallToolResult, any, error) {
    if in.Query == "" {
        return nil, nil, fmt.Errorf("query is required")
    }
    results, err := lookupHotels(ctx, in.Query, in.Limit)
    if err != nil {
        return nil, nil, fmt.Errorf("search hotels: %w", err)
    }
    return &mcp.CallToolResult{
        Content: []mcp.Content{&mcp.TextContent{Text: formatResults(results)}},
    }, nil, nil
}

func main() {
    server := mcp.NewServer(&mcp.Implementation{Name: "hotel-search", Version: "1.0.0"}, nil)
    mcp.AddTool(server, &mcp.Tool{Name: "search_hotels", Description: "Search hotels by name or city"}, searchHotels)

    if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
        log.Fatal(err)
    }
}
```

For a remote deployment, swap the transport at the entrypoint only — the
tool registration above is unchanged:

```go
handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil)
http.ListenAndServe(":8080", handler)
```

## TypeScript notes

The TypeScript SDK (`@modelcontextprotocol/sdk`) exposes the same three
primitives; recent versions favor `server.registerTool(name, config,
handler)` / `registerResource` / `registerPrompt` over older `tool()` /
`resource()` free functions from earlier releases — check the installed
package's version against current docs before copying an example, since
this API has changed shape more than once. Use **Zod** for input schema
validation, mirroring the boundary-validation rule in `api-design`: every
tool input is untrusted until validated, because it can originate from
model-generated arguments as easily as from a malicious prompt injection
upstream.

## Schema-first tool design

Define the input schema before the handler, and make it strict:

```go
type CreateReservationInput struct {
    HotelID  int    `json:"hotel_id" jsonschema:"required"`
    GuestID  string `json:"guest_id" jsonschema:"required"`
    CheckIn  string `json:"check_in" jsonschema:"required,format=date"`
    CheckOut string `json:"check_out" jsonschema:"required,format=date"`
}
```

A loose schema (everything optional, `any` fields) pushes validation
into the handler body where it's easy to skip a case; a strict schema
lets the SDK reject a malformed call before your handler code ever runs.

## Best practices

- **Idempotent tools where possible.** A model may retry a tool call
  after a timeout without knowing whether the first call succeeded —
  design mutating tools to accept a client-supplied idempotency key, or
  make the operation naturally idempotent (upsert, not insert).
- **Structured errors the model can act on**, not raw stack traces or
  driver error text: return a message like `"reservation not found for
  id X"` rather than a raw `sql.ErrNoRows` string — the model needs to
  decide whether to retry, ask the user, or try a different tool.
- **Document cost and rate limits in the tool description itself** for
  any tool that calls a metered or rate-limited external API — the model
  has no other way to know calling this tool 50 times in a row is
  expensive.
- **Pin the SDK version** in `go.mod`/`package.json` and read release
  notes before bumping — this protocol and its SDKs are under active
  revision.

## Gotchas

- A resource handler that has a side effect (writes a file, calls an
  external API with billing implications) violates the tool/resource
  distinction and will surprise a client that assumes resources are safe
  to fetch speculatively or repeatedly for context — make it a tool.
- `stdio` transport assumes your server's stdout is exclusively the MCP
  JSON-RPC channel — any stray `fmt.Println`/`console.log` debug output
  written to stdout instead of stderr corrupts the protocol stream and
  produces a confusing client-side parse error far from the actual bug.
- Registering a tool with a vague description ("does stuff with hotels")
  gives the model no signal about when to call it versus a similarly
  named tool — the description is the model's only interface contract,
  write it like API documentation, not a code comment.
- Streamable HTTP requires session/auth handling remote-server authors
  often skip in a local-first prototype — don't ship a stdio-only
  prototype's auth assumptions (none) straight to a multi-tenant remote
  deployment.

## Real-world grounding

The MCP specification itself replaced its original HTTP+SSE transport
with **Streamable HTTP** in the 2025-03-26 spec revision, specifically to
support stateless deployments and simplify resumability — official SDKs
maintained backward-compatible support for HTTP+SSE clients for a
transition window rather than breaking them outright, the same
advisory-deprecation pattern covered in `deprecation-and-migration`
applied to a wire protocol instead of a library API.

## Verification

- [ ] Tool/resource/prompt logic has no direct dependency on the chosen transport
- [ ] Every tool has a strict, explicit input schema — no free-form `any` payloads
- [ ] Mutating tools are idempotent or accept an idempotency key
- [ ] Tool descriptions are written as documentation for the model, not internal comments
- [ ] stdio servers write only protocol messages to stdout; all logging goes to stderr
- [ ] SDK version is pinned and current signatures were checked against Context7/official docs, not assumed from memory
