# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Install dependencies
go mod tidy

# Run the chat CLI (auto-spawns MCP server subprocess)
ANTHROPIC_API_KEY=<key> go run .

# Run the MCP server directly (stdio transport, meant to be spawned by main)
go run ./mcp-server
```

Configuration lives in `config.yaml`. `ANTHROPIC_API_KEY` must be set as an environment variable (or under `anthropic.api_key` in `config.yaml`).

No test suite or linter config exists yet.

## Architecture

An MCP (Model Context Protocol) sample app combined with an interactive Claude chat CLI. `main.go` spawns the MCP server as a subprocess and then starts a REPL backed by the Anthropic SDK.

**Transport**: `main.go` spawns `mcp-server` as a subprocess via `CommandTransport` (`go run ./mcp-server`). All MCP communication is message-passing over stdio — not HTTP.

**Packages**:
- `main.go` — entry point; loads config, connects MCP client, starts the `CliApp` REPL
- `chat/chat.go` — `Chat` struct; implements the agentic loop (`RunMessages`) that drives tool-use cycles until `stop_reason != tool_use`
- `chat/claude.go` — `Claude` struct; thin wrapper around the Anthropic SDK client
- `chat/clchat.go` — `CliChat` extends `Chat`; handles doc listing, `@mention` resource injection, and `/command` prompt dispatch
- `chat/cli.go` — `CliApp`; stdin REPL that lists available docs and prompts on startup
- `chat/tools.go` — bridges MCP tools to Anthropic's tool format; routes tool call results back to the correct MCP client
- `config/config.go` — uses viper to load `config.yaml` and overlay `ANTHROPIC_API_KEY` from the environment
- `mcp-client/client.go` — thin wrapper around the MCP SDK client
- `mcp-server/main.go` — MCP server; registers tools, resources, and prompts against an in-memory document map

**Chat features**:
- `@docname` in a query injects the document's content as context before sending to Claude
- `/format <docname>` dispatches the MCP `format` prompt; the agentic loop then applies it via the `edit` tool
- The agentic loop continues calling MCP tools (`read`, `edit`) until Claude stops with `stop_reason != tool_use`

**Server capabilities**:
- Tools: `read` (fetch a doc by name), `edit` (update doc content)
- Resources: `docs://documents` (list all), `docs://documents/{name}` (get one)
- Prompts: `format` (instructs Claude to reformat a doc using the `edit` tool)

**Shared state**: The `docs` map in `mcp-server/main.go` is module-level. Tools and resource handlers both reference it directly, so `edit` tool mutations are immediately visible to subsequent resource reads.

**Configuration** (`config.yaml`): `claude.model`, `claude.max_tokens`, `claude.itinerary_max_tokens`

**SDKs**: `github.com/modelcontextprotocol/go-sdk v1.4.1`, `github.com/anthropics/anthropic-sdk-go`
