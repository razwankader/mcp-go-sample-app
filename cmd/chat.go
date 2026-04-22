package cmd

import (
	"context"
	"log"
	"mcp-go-sample-app/chat"
	"mcp-go-sample-app/claude"
	"mcp-go-sample-app/config"
	mcpclient "mcp-go-sample-app/mcp-client"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

var (
	chatCmd = &cobra.Command{
		Use:   "chat",
		Short: "Start an interactive chat session with MCP and Claude",
		Long:  "Start an interactive chat session with MCP and Claude. This command connects to the MCP server, lists available prompts and documents, and allows you to interact with them using a CLI interface.",
		Run:   initChat,
	}
)

func init() {
	rootCmd.AddCommand(chatCmd)
}

func initChat(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	cfg := config.Get()

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
