package main

import (
	"context"
	"log"
	"mcp-go-sample-app/chat"
	"mcp-go-sample-app/claude"
	"mcp-go-sample-app/config"
	mcpclient "mcp-go-sample-app/mcp-client"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Claude.Model == "" {
		log.Fatal("Error: CLAUDE_MODEL cannot be empty. Update config")
	}

	if cfg.Anthropic.APIKey == "" {
		log.Fatal("Error: ANTHROPIC_API_KEY cannot be empty. Update config")
	}

	client := mcpclient.New()
	transport := &mcp.CommandTransport{
		Command: exec.Command("go", "run", "./mcp-server"),
	}
	if err := client.Connect(ctx, transport); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	clients := map[string]*mcpclient.MCPClient{"docs": client}
	claude := claude.NewClaude(cfg.Anthropic.APIKey, cfg.Claude.Model)

	cliChat := chat.NewCliChat(client, clients, claude)
	app := chat.NewCliApp(cliChat)

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
