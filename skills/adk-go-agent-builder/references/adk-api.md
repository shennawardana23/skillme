# ADK Go — Extended API Reference

Source of truth: https://github.com/google/adk-go (Apache-2.0). Verify
against the actual module (`go doc google.golang.org/adk/v2/...` after
`go get google.golang.org/adk/v2`) before relying on any signature below for
a production change — this reference is a map to the upstream examples, not
a frozen copy of them, and the upstream repo is under active development.

## Upstream example directory map

`google/adk-go/examples/`:

| Directory | Demonstrates |
|---|---|
| `quickstart/` | Minimal single agent with a built-in tool, run through `launcher/full` |
| `tools/multipletools/` | Working around the Gemini-tool + custom-tool conflict via agent-as-tool |
| `tools/loadartifacts/`, `tools/loadmemory/` | Artifact and memory-backed tools |
| `multiagent/single_turn/` | Basic agent-as-tool composition |
| `multiagent/collaboration/` | Multiple specialist agents collaborating |
| `multiagent/task_sub_agent/` | Sub-agent transfer (hierarchical hand-off) |
| `workflowagents/sequential/`, `sequentialCode/` | `sequentialagent` — run sub-agents in a fixed order |
| `workflowagents/parallel/` | `parallelagent` — run sub-agents concurrently, merge results |
| `workflowagents/loop/` | `loopagent` — repeat a sub-agent until a stop condition |
| `a2a/` | Agent-to-Agent protocol interop (`a2aproject/a2a-go`) |
| `mcp/` | Using MCP servers as ADK tools (`modelcontextprotocol/go-sdk`) |
| `bidi/` | Bidirectional/streaming sessions |
| `rest/`, `web/` | Serving an agent over REST / a web UI |
| `vertexai/` | Vertex AI-backed model instead of the direct Gemini API |
| `agentengine/` | Deploying to Google's managed Agent Engine |
| `telemetry/` | Tracing/observability hooks |
| `toolconfirmation/` | Human-in-the-loop tool-call confirmation |
| `skills/` | ADK's own notion of a "skill" for an agent (distinct from a Claude Code Skill) |

## Workflow agents (`agent/workflowagents/*`)

Use these when the orchestration is a fixed structural shape known ahead of
time (always run A then B, always run these three in parallel and merge, or
retry a step until a condition holds) — reach for plain `SubAgents`
transfer or `agenttool` composition instead when the routing decision is
something the LLM should make dynamically per conversation.

- `agent/workflowagents/sequentialagent` — runs a fixed list of sub-agents
  in order, passing state forward.
- `agent/workflowagents/parallelagent` — runs a fixed list of sub-agents
  concurrently and merges their results.
- `agent/workflowagents/loopagent` — reruns a sub-agent (or sequence) until
  a stop condition is met or a max-iteration bound is hit; always set the
  bound explicitly to avoid an unbounded loop against a paid model API.

## Model providers beyond direct Gemini

`model/gemini.NewModel` calls the Gemini API directly with an API key. For
Vertex AI-backed deployments (organization service-account auth, regional
routing, quota pooled at the project level rather than per API key), see
`examples/vertexai/` instead — the `model` interface `llmagent.Config.Model`
accepts is provider-agnostic, so swapping `gemini.NewModel` for a Vertex AI
model constructor does not otherwise change how `llmagent.New` or
`functiontool.New` are used.

## Checking installed version and API surface locally

Before writing non-trivial ADK Go code, confirm the exact API surface of the
version actually pinned in the target module rather than assuming this
reference is current:

```bash
go get google.golang.org/adk/v2@latest
go doc google.golang.org/adk/v2/agent/llmagent
go doc google.golang.org/adk/v2/tool/functiontool
```
