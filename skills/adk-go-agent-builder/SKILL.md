---
name: adk-go-agent-builder
description: This skill should be used when the user asks to "build an agent with ADK Go", "create an ADK agent in Go", "add a tool to my Go agent", "use google.golang.org/adk", "compose sub-agents in ADK Go", or writes Go code against Google's Agent Development Kit for Go (google.golang.org/adk/v2). Provides verified API patterns for agents, models, function tools, and multi-agent composition.
metadata:
  version: "0.1.0"
---

# ADK Go Agent Builder

Google's Agent Development Kit has an official, actively maintained Go
module: `google.golang.org/adk/v2` (`go get google.golang.org/adk/v2`,
requires Go 1.26+). This is a distinct, real SDK — not a hypothetical or
Python-only port — with its own idioms that differ from the Python/Java/
Kotlin/TypeScript ADK variants. Ground every code sample in the packages
below rather than porting Python ADK syntax (`google.adk.agents.Agent`,
decorator-based tools) directly into Go; the Go API is call-based and
config-struct-driven.

## Core packages

| Package | Purpose |
|---|---|
| `google.golang.org/adk/v2/agent` | Core `agent.Agent` interface, `agent.Context`, `agent.NewSingleLoader` |
| `google.golang.org/adk/v2/agent/llmagent` | `llmagent.New(llmagent.Config{...})` — the standard LLM-backed agent |
| `google.golang.org/adk/v2/model/gemini` | `gemini.NewModel(ctx, modelName, *genai.ClientConfig)` |
| `google.golang.org/adk/v2/tool` | `tool.Tool` interface |
| `google.golang.org/adk/v2/tool/functiontool` | `functiontool.New(functiontool.Config{...}, handler)` — wrap a typed Go function as a tool |
| `google.golang.org/adk/v2/tool/agenttool` | `agenttool.New(subAgent, *agenttool.Config)` — expose one agent as a tool for another |
| `google.golang.org/adk/v2/tool/geminitool` | Built-in Gemini tools, e.g. `geminitool.GoogleSearch{}` |
| `google.golang.org/adk/v2/cmd/launcher` + `.../launcher/full` | CLI/server launcher that runs an agent |

## Minimal agent

```go
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/cmd/launcher/full"
	"google.golang.org/adk/v2/model/gemini"
)

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-flash-latest", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create model: %v", err)
	}

	weatherAgent, err := llmagent.New(llmagent.Config{
		Name:        "weather_agent",
		Model:       model,
		Description: "Answers questions about current weather.",
		Instruction: "Your SOLE purpose is to answer weather questions. Refuse anything else.",
	})
	if err != nil {
		log.Fatalf("create agent: %v", err)
	}

	l := full.NewLauncher()
	cfg := &launcher.Config{AgentLoader: agent.NewSingleLoader(weatherAgent)}
	if err := l.Execute(ctx, cfg, os.Args[1:]); err != nil {
		log.Fatalf("run: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
```

`Instruction` is the system-prompt-equivalent field — write it as a hard
behavioral constraint ("Your SOLE purpose is...", "Refuse anything else"),
not a vague topic hint, since the model uses it to decide what to refuse.

## Custom function tools

Wrap a typed Go function with `functiontool.New` instead of hand-rolling a
JSON-schema tool definition — ADK derives the schema from the `Input` /
`Output` struct's JSON tags:

```go
type ServerTimeInput struct{}

type ServerTimeOutput struct {
	RFC3339 string `json:"rfc3339"`
}

serverTimeTool, err := functiontool.New(functiontool.Config{
	Name:        "server_time",
	Description: "Returns the current server time in RFC3339 format.",
}, func(ctx agent.Context, in ServerTimeInput) (ServerTimeOutput, error) {
	return ServerTimeOutput{RFC3339: time.Now().UTC().Format(time.RFC3339)}, nil
})
if err != nil {
	log.Fatalf("create tool: %v", err)
}

clockAgent, err := llmagent.New(llmagent.Config{
	Name:        "clock_agent",
	Model:       model,
	Description: "Answers questions about the current time.",
	Instruction: "Use the server_time tool to answer time questions.",
	Tools:       []tool.Tool{serverTimeTool},
})
```

The handler signature is fixed: `func(ctx agent.Context, input In) (Out, error)`
— both `In` and `Out` must be JSON-serializable structs (even an empty
struct for no-input tools), and the returned `error` becomes a tool-call
failure the model can see and react to, not a Go-level panic.

## Composing multiple agents

Two distinct composition patterns exist — pick based on whether the parent
needs to keep control after the sub-agent responds:

**Agent-as-tool** (`agenttool.New`): the parent calls the sub-agent the same
way it calls any other tool, gets a result back, and keeps driving the
conversation. Use this when a specialist agent (e.g., one already using
`geminitool.GoogleSearch`, which cannot coexist with custom function tools
on the same agent — see below) needs to be composed with other tools:

```go
searchAgent, _ := llmagent.New(llmagent.Config{
	Name:        "search_agent",
	Model:       model,
	Description: "Performs Google Search.",
	Instruction: "You are a search specialist.",
	Tools:       []tool.Tool{geminitool.GoogleSearch{}},
})

rootAgent, _ := llmagent.New(llmagent.Config{
	Name:        "root_agent",
	Model:       model,
	Description: "Delegates search questions and answers everything else directly.",
	Instruction: "For anything requiring a web search, use the search_agent tool.",
	Tools:       []tool.Tool{agenttool.New(searchAgent, nil)},
})
```

**Sub-agent transfer** (`SubAgents` on the parent's config, see
`google.golang.org/adk/v2/agent/workflowagent` and
`.../agent/workflowagents/*` for sequential/parallel/loop orchestration):
the parent hands off the entire turn to a specialist and does not
necessarily regain control. Use this for hierarchical routing (a triage
agent that permanently hands off to the right specialist), not for a single
tool-like call-and-return.

`agenttool.Config{SkipSummarization: true}` skips the parent's LLM
summarization pass over the sub-agent's raw output — set it when the
sub-agent's output should be returned to the user verbatim rather than
re-narrated.

## Gemini tool + custom tool conflict

`geminitool.GoogleSearch{}` (and other built-in Gemini tools) cannot be
combined with custom `functiontool` tools on the same `llmagent` — this is a
genai API limitation, not an ADK bug. The workaround is exactly the
agent-as-tool pattern above: put `GoogleSearch` on its own agent, wrap that
agent with `agenttool.New`, and give the wrapped tool to the agent that also
holds custom function tools.

## Additional resources

For end-to-end runnable examples straight from the upstream repository
(`google/adk-go/examples/`), including `multiagent/`, `workflowagents/`, and
`tools/`, consult `references/adk-api.md`.
