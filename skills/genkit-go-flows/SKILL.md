---
name: genkit-go-flows
description: Guides building AI features with Genkit for Go - flows, prompts, structured generation, tool calling, and middleware. Use when defining a Genkit flow, writing a prompt, calling genkit.Generate/GenerateData, defining a tool, debugging Genkit output with traces, or asked how Genkit differs from Google's Agent Development Kit (ADK).
license: Apache-2.0
metadata:
  version: "0.1.0"
  category: "go"
---

# Genkit for Go: Flows, Prompts, and Tools

Genkit (`github.com/firebase/genkit/go`) is Google/Firebase's application
framework for building AI features — flows, prompt management, structured
generation, tool calling, RAG, and evals — with first-class Go support
alongside JS/TS. It is not the same tool as `adk-go-agent-builder`'s Agent
Development Kit: Genkit is the broader app-framework layer (any AI feature,
not just autonomous agents), while ADK is specifically for building
orchestrated, tool-using agents. A service can reasonably use both — Genkit
for its prompt/generation plumbing and observability, ADK for the parts
that are genuinely agentic.

## Defining a flow

A flow is Genkit's typed, traceable unit of AI-feature logic — register it
against a `*genkit.Genkit` instance rather than calling generation
functions ad hoc, so it shows up in traces and can be run standalone via
the CLI:

```go
g := genkit.Init(ctx)

flow := genkit.DefineFlow(g, "summarizeReservation", func(ctx context.Context, input ReservationInput) (string, error) {
    resp, err := genkit.Generate(ctx, g,
        ai.WithPrompt("Summarize this reservation in one sentence: %+v", input),
    )
    if err != nil {
        return "", fmt.Errorf("generate summary: %w", err)
    }
    return resp.Text(), nil
})
```

Pass `g` (the `*Genkit` instance) explicitly into every call that needs
it — it is Genkit's central action registry (flows, prompts, tools all
register against it). Storing it as a package-level global instead of
threading it through breaks the documented pattern and makes a service
harder to test with an isolated registry per test.

## Running with traces (the dev loop that actually works)

```bash
genkit start -- go run .              # traced, with Dev UI
genkit flow:run summarizeReservation '{"hotelId":"h1"}' -- go run .   # traced, non-interactive
```

`go run .` on its own does not capture Genkit dev traces — only running
under `genkit start` (or the equivalent `genkit flow:run` for a single,
self-terminating invocation) wires up the trace collector. Debugging a
flow by running the binary directly is debugging blind: no trace, no Dev
UI inspection of the actual prompt/response the model saw.

`genkit start` is a persistent, blocking server process — never use it as
a step in a non-interactive script or CI job. `genkit flow:run <name>
'<json-input>' -- go run .` runs one flow invocation with tracing and
terminates on its own, which is what a script or CI step actually needs.

## Structured output

```go
type Availability struct {
    Available bool   `json:"available"`
    Reason    string `json:"reason" jsonschema:"description=Why availability is or isn't there, in plain language"`
}

resp, err := genkit.Generate(ctx, g,
    ai.WithPrompt("Is room type %s available on %s?", roomType, date),
    ai.WithOutputType(Availability{}),
)
```

`jsonschema:"description=..."` struct tags materially affect output
quality — the model is guided by those field descriptions when producing
structured output, not just by the field names and types. A struct with
untagged fields tends to produce lower-quality structured results than the
same struct with real descriptions on each field.

## Tool calling

```go
lookupRoomRate := genkit.DefineTool(g, "lookupRoomRate",
    "Looks up the current rate for a room type and date",
    func(ctx context.Context, input RateLookupInput) (RateLookupOutput, error) {
        return rateService.Lookup(ctx, input.RoomType, input.Date)
    },
)

resp, err := genkit.Generate(ctx, g,
    ai.WithPrompt(userQuery),
    ai.WithTools(lookupRoomRate),
)
```

## Gotchas

- **The correct module is `github.com/firebase/genkit/go`.** Some
  installed documentation/examples reference an older `github.com/genkit-ai/genkit/go`
  import path — that path is stale; verify the actual current module via
  `pkg.go.dev` or `go get` before trusting an example's import block
  verbatim, since propagating a stale import produces a build that won't
  resolve.
- **`go run .` alone gives no trace visibility** — always develop against
  `genkit start -- go run .` (interactive) or `genkit flow:run <name>
  '<json>' -- go run .` (scriptable) instead of the bare binary when
  debugging flow behavior.
- **Never put `genkit start` in a CI/script step** — it's a blocking
  server, not a one-shot command; use `genkit flow:run` for a
  self-terminating, scriptable invocation instead.
- **Structured output quality depends on `jsonschema` description tags**,
  not just field names/types — add real descriptions to output-type
  fields rather than assuming the schema alone is enough guidance for the
  model.
- **The agent-oriented experimental API is gated and will panic if the
  gate is missing** — anything under Genkit's experimental agent surface
  requires `genkit.Init(ctx, genkit.WithExperimental())`; omitting that
  option and calling the experimental API panics rather than failing
  gracefully. Don't reach for the experimental agent surface for a
  simple, stateless flow — it exists for the genuinely agentic case.
- **Default model IDs are not stable across Genkit releases** — the
  official docs describe defaults as subject to frequent change; pin an
  explicit model ID for production code rather than relying on whatever
  a given Genkit version's default happens to be, and re-check it after
  any Genkit upgrade.

## Real-world grounding

Genkit for Go is maintained as part of Google/Firebase's multi-language
Genkit project (JS/TS first, Go and Python following), documented at
genkit.dev and distributed as `github.com/firebase/genkit/go` on
pkg.go.dev — verify any import path or API shape against those sources
directly rather than an older tutorial, since the project has moved
module paths at least once already (the `genkit-ai/genkit` →
`firebase/genkit` transition noted above).

## Verification

- [ ] Flows are registered via `genkit.DefineFlow(g, ...)` with `g` passed
      explicitly, not held as a package global
- [ ] Development/debugging runs under `genkit start` or `genkit flow:run`,
      never the bare binary, when trace visibility matters
- [ ] No `genkit start` invocation exists inside a non-interactive script
      or CI step
- [ ] Structured output types have `jsonschema` description tags on their
      fields
- [ ] Any experimental agent-surface usage explicitly enables
      `genkit.WithExperimental()` first
- [ ] The model ID used in production is explicit, not left as an
      unpinned default
