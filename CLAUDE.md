# CLAUDE.md

Guidance for Claude Code (claude.ai/code) when working with code in this repo.

## Commands

```bash
# Install dependencies
go mod tidy

# Run chat CLI (spawns MCP server subprocess)
ANTHROPIC_API_KEY=<key> go run .

# Run MCP server directly (stdio, spawned by main)
go run ./mcp-server
```

Config in `config.yaml`. `ANTHROPIC_API_KEY` must be set in env (or in `config.yaml` under `anthropic.api_key`).

No test suite or linter config yet.

## Architecture

MCP sample app + interactive Claude chat CLI. `main.go` spawns MCP server subprocess, then runs REPL backed by Anthropic SDK.

**Transport**: `main.go` spawns `mcp-server` via `CommandTransport`. All MCP comms over stdio — not HTTP.

**Packages**:
- `main.go` — entry; loads config, connects MCP client, starts `CliApp` REPL
- `chat/chat.go` — `Chat`: agentic loop (`RunMessages`), tool-use cycle until `stop_reason != tool_use`
- `chat/claude.go` — `Claude`: thin Anthropic SDK wrapper
- `chat/clchat.go` — `CliChat`: extends `Chat`; doc listing, `@mention` resource injection, `/command` prompt dispatch
- `chat/cli.go` — `CliApp`: stdin REPL; lists available docs/prompts on start
- `chat/tools.go` — MCP→Anthropic tool bridging; routes tool calls to correct MCP client
- `config/config.go` — viper loads `config.yaml` + `ANTHROPIC_API_KEY` env var
- `mcp-client/client.go` — thin MCP SDK client wrapper
- `mcp-server/main.go` — MCP server; tools, resources, prompts against in-memory doc map

**Chat features**:
- `@docname` in query → injects doc content as context before sending to Claude
- `/format <docname>` → dispatches MCP `format` prompt; agentic loop applies it
- Agentic loop: Claude calls MCP tools (`read`, `edit`) until `stop_reason != tool_use`

**Server capabilities**:
- Tools: `read` (fetch doc by name), `edit` (update doc content)
- Resources: `docs://documents` (list all), `docs://documents/{name}` (get one)
- Prompts: `format` (reformat doc via `edit` tool)

**Shared state**: `docs` map in `mcp-server/main.go` module-level — `edit` mutations immediately visible to subsequent reads.

**Config** (`config.yaml`): `claude.model`, `claude.max_tokens`, `claude.itinerary_max_tokens`

**SDKs**: `github.com/modelcontextprotocol/go-sdk v1.4.1`, `github.com/anthropics/anthropic-sdk-go`