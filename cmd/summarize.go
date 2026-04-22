package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	mcpclient "mcp-go-sample-app/mcp-client"
	"os/exec"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

var (
	docName      string
	summarizeCmd = &cobra.Command{
		Use:   "summarize",
		Short: "Summarize a document using the summarize tool and mcp sampling feature",
		Long:  "Summarize a document using the summarize tool and mcp sampling feature. This command connects to the MCP server, calls the summarize tool with the specified document name, and prints the summarized content.",
		Run:   summarize,
	}
)

func init() {
	summarizeCmd.PersistentFlags().StringVar(&docName, "doc", "", "document name")
	summarizeCmd.MarkPersistentFlagRequired("doc")

	rootCmd.AddCommand(summarizeCmd)
}

func summarize(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	client := mcpclient.New()
	transport := &mcp.CommandTransport{
		Command: exec.Command("go", "run", "./mcp-server"),
	}
	if err := client.Connect(ctx, transport); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	toolName := "summarize"
	inputMap := make(map[string]any)
	inputMap["name"] = docName

	output, err := client.CallTool(ctx, toolName, inputMap)
	if err != nil {
		log.Fatalf("Error calling tool: %v", err)
	}

	var textParts []string
	for _, item := range output.Content {
		b, err := json.Marshal(item)
		if err != nil {
			continue
		}
		var tc struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}
		if err := json.Unmarshal(b, &tc); err != nil {
			continue
		}
		if tc.Text != "" {
			textParts = append(textParts, tc.Text)
		}
	}
	content := strings.Join(textParts, "\n")
	fmt.Printf("Output: %s", content)
}
