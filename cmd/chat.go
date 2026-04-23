package cmd

import (
	"context"
	"log"
	"mcp-go-sample-app/chat"
	"mcp-go-sample-app/claude"
	"mcp-go-sample-app/config"
	mcpclient "mcp-go-sample-app/mcp-client"
	"net/url"
	"os/exec"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

var (
	rootPaths []string
	chatCmd   = &cobra.Command{
		Use:   "chat",
		Short: "Start an interactive chat session with MCP and Claude",
		Long:  "Start an interactive chat session with MCP and Claude. This command connects to the MCP server, lists available prompts and documents, and allows you to interact with them using a CLI interface.",
		Run:   initChat,
	}
)

func init() {
	chatCmd.PersistentFlags().StringSliceVarP(
		&rootPaths,
		"dir",
		"d",
		[]string{},
		"List of valid directory paths",
	)

	rootCmd.AddCommand(chatCmd)
}

func initChat(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	cfg := config.Get()

	client := mcpclient.New()
	transport := &mcp.CommandTransport{
		Command: exec.Command("go", "run", "./mcp-server"),
	}

	client.AddRoots(ctx, buildRoots(rootPaths))

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

func buildRoots(rootPaths []string) []*mcp.Root {
	var roots []*mcp.Root

	for _, path := range rootPaths {
		p, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		u := &url.URL{
			Scheme: "file",
			Path:   filepath.ToSlash(p),
		}

		name := filepath.Base(p)
		if name == "" {
			name = "Root"
		}

		roots = append(roots, &mcp.Root{
			URI:  u.String(),
			Name: name,
		})
	}

	return roots
}
