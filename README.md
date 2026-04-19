# mcp-go-sample-app

An interactive chat CLI powered by Claude AI for document management using the Model Context Protocol (MCP). This sample application demonstrates MCP client-server communication with AI integration.

Built with the official MCP Go SDK: `github.com/modelcontextprotocol/go-sdk`.

## Features

- **Interactive Chat Interface**: Query documents and get AI-powered responses from Claude.
- **Document Management**: Read, edit, and format documents via MCP tools.
- **Mention Support**: Reference documents in queries using `@docname` syntax.
- **Prompt Commands**: Run predefined prompts like `/format docname` to reformat documents.
- **MCP Integration**: Full MCP client and server implementation with stdio transport.

## Prerequisites

- Go 1.23+ installed
- Anthropic API key (set as `ANTHROPIC_API_KEY` environment variable)

## Setup

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Set your Anthropic API key:
   ```bash
   export ANTHROPIC_API_KEY=your_api_key_here
   ```

3. (Optional) Configure settings in `config.yaml`:
   ```yaml
   anthropic:
     api_key: your_api_key_here
   claude:
     model: claude-haiku-4-5-20251001
     max_tokens: 2048
   ```

## Usage

Run the interactive chat application:

```bash
go run .
```

The app will:
- Start the MCP server automatically
- Display available documents and prompts
- Enter an interactive prompt for queries

### Example Interactions

- List documents: The app shows available docs on startup
- Query with document reference: `What is the status of the condenser tower? @report.pdf`
- Run a prompt: `/format deposition.md` (reformats the document using Claude)

### Available Prompts

- `/format <docname>`: Rewrites the document in Markdown format

## Architecture

- **Client**: Interactive CLI that integrates Claude AI with MCP tools
- **Server**: MCP server providing document read/edit tools and resources
- **Transport**: Stdio-based communication between client and server
- **State**: In-memory document storage (shared between client sessions)

## Project Structure

- `main.go`: CLI application entry point
- `chat/`: Claude integration and chat logic
- `config/`: Configuration management
- `mcp-client/`: MCP client wrapper
- `mcp-server/`: MCP server implementation
- `config.yaml`: Configuration file
