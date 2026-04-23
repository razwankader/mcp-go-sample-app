# MCP Go Sample App

A comprehensive Go application that showcases the full capabilities of the Model Context Protocol (MCP), including tools, resources, prompts, sampling, logging, notifications, and root management, with the help of a CLI interface, document management, and AI-powered chat capabilities.

## Features

- **MCP Tools**: Complete tool registration, calling, and result handling
- **MCP Resources**: Static and templated resource management
- **MCP Prompts**: Reusable prompt templates with argument support
- **MCP Sampling**: Server-side AI model integration through sampling callbacks
- **MCP Server Logging and Notifications**: Progress tracking through notifications and structured logging with configurable levels
- **MCP Roots**: File system root management for workspace directories
- **In-memory Document Management**: Document operations demonstrating MCP tools, resources, and sampling
- **CLI Interface**: Modern command-line interface using Cobra with multiple commands
- **AI Chat**: Interactive chat with Claude AI integration

## Quick Start

### Prerequisites

- Go 1.21 or later
- Anthropic API key

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd mcp-go-sample-app
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configuration**
   ```bash
   cp config-example.yaml config.yaml
   ```
   #### Update config.yaml
   
   ```yaml
   anthropic:
      api_key: "your-anthropic-api-key-here"

   claude:
      model: "claude-3-5-sonnet-20241022"
      max_tokens: 8000
      itinerary_max_tokens: 4000
   ```

4. **Build the application**
   ```bash
   go build
   ```

## Usage

### Interactive Chat
```bash
./mcp-go-sample-app chat
```

### Interactive Chat (with MCP roots feature)
```bash
./mcp-go-sample-app chat -d /<local-machine-base-path>/Documents -d /<local-machine-base-path>/Desktop
```
Roots are a way to grant MCP servers access to specific files and folders on your local machine. This means the MCP server currently has access only to the **Documents** and **Desktop** folders on your local machine.

### List All Available In-Memory Documents
```bash
./mcp-go-sample-app docs
```

### Document Summarization (MCP Sampling feature)
```bash
./mcp-go-sample-app summarize --doc report.pdf
```
### Help
```bash
./mcp-go-sample-app --help
./mcp-go-sample-app chat --help
./mcp-go-sample-app summarize --help
```

## MCP Feature Examples

The in-memory document management showcases various MCP capabilities:

### MCP Tools Example
Claude automatically calls MCP tools based on your prompts. Here are common scenarios:

Read Tool - Accessing Document Content
```bash
./mcp-go-sample-app chat
> read the content of @report.pdf
# Claude automatically calls the read MCP tool to fetch report.pdf content
```

Edit Tool - Modifying Document Content
```bash
./mcp-go-sample-app chat
> replace the text "Go MCP SDKs" with "Python MCP SDKs" in report.pdf  
# Claude automatically calls the edit MCP tool to update the document
```

### MCP Resources Example  
```bash
# List all documents (MCP Resource)
./mcp-go-sample-app docs
# Accesses MCP resource at docs://documents to list available documents

```
### MCP Prompts Example
```bash
# Interactive chat mode
./mcp-go-sample-app chat

> /format report.pdf 
# Executes MCP prompt template to reformat document in Markdown
```

### MCP Sampling Example
```bash
# The MCP server uses sampling to call Claude AI via MCP client for document summarization
./mcp-go-sample-app summarize --doc report.pdf
# Demonstrates MCP sampling: server-side AI integration through MCP callbacks
```

## MCP Topics Covered

### 1. MCP Transport
- **Stdio Transport**: Communication between client and server via standard input/output
- **Command Transport**: Spawning MCP server subprocess from client

### 2. MCP Tools
- **Tool Registration**: Registering tools on the server side
- **Tool Calling**: Client-side tool invocation with parameters
- **Tool Results**: Handling structured tool responses

### 3. MCP Resources
- **Resource Registration**: Static and templated resource URIs
- **Resource Reading**: Accessing document content via MCP resources
- **File System Roots**: Managing workspace directories

### 4. MCP Prompts
- **Prompt Templates**: Reusable prompt definitions with arguments
- **Prompt Execution**: Client-side prompt invocation

### 5. MCP Sampling
- **Sampling Callbacks**: Server-side AI model integration
- **Message Translation**: Converting MCP messages to provider-specific formats
- **Response Handling**: Processing AI responses back to MCP format

### 6. MCP Notifications
- **Progress Notifications**: Real-time progress updates for long operations
- **Logging Notifications**: Structured logging with configurable levels
- **Progress Callbacks**: Client-side progress tracking

### 7. MCP Lifecycle
- **Client Initialization**: Setting up MCP clients with handlers
- **Session Management**: Connection lifecycle and cleanup
- **Error Handling**: Robust error propagation and recovery

## Testing MCP Server Independently

You can test the MCP server independently using the MCP Inspector tool.

### Run Inspector (separately)
```bash
npx @modelcontextprotocol/inspector
```

A dashboard will be loaded with URL http://localhost:6274/

### In Left Panel
- Choose **TransportType** = STDIO
- **Command** = go
- **Arguments** = run ./mcp-server/main.go
- Press **Connect**

You will see all the available Tools, Resources and Prompts are loaded. You can test those with proper input parameters.


## Project Structure

```
mcp-go-sample-app/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command definition
│   ├── chat.go            # Interactive chat command
│   ├── summarize.go       # Document summarization command
│   └── docs.go            # Document listing command
├── chat/                   # Chat interface components
│   ├── chat.go            # Core chat logic
│   ├── clchat.go          # CLI chat implementation
│   ├── cli.go             # CLI application
│   └── tools.go           # MCP tool bridging
├── claude/                 # Claude AI integration
│   └── claude.go          # Claude API wrapper
├── config/                 # Configuration management
│   └── config.go          # Config loading
├── mcp-client/            # MCP client implementation
│   └── client.go          # MCP client with sampling
├── mcp-server/            # MCP server implementation
│   └── main.go            # MCP server with tools/resources
├── config.yaml            # Application configuration
├── config-example.yaml    # Configuration template
├── main.go                # Application entry point
├── go.mod                 # Go module definition
└── README.md              # This file
```


### Key Components

- **CLI Layer**: Command-line interface using Cobra framework
- **Chat Layer**: Interactive chat interface with document context
- **MCP Layer**: Protocol implementation for client-server communication
- **AI Integration**: Claude API wrapper with message handling


## Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io/) - The protocol this project demonstrates
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Official Go implementation of MCP
- [Anthropic Claude](https://www.anthropic.com/claude) - AI model used for chat and summarization
- [Cobra](https://github.com/spf13/cobra) - CLI framework

