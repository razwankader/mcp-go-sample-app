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
	docsCmd = &cobra.Command{
		Use:   "docs",
		Short: "List document IDs available in the MCP server",
		Long:  "List document IDs available in the MCP server. This command connects to the MCP server, reads the 'docs://documents' resource, and prints the list of document IDs.",
		Run:   docs,
	}
)

func init() {
	rootCmd.AddCommand(docsCmd)
}

func docs(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	client := mcpclient.New()
	transport := &mcp.CommandTransport{
		Command: exec.Command("go", "run", "./mcp-server"),
	}
	if err := client.Connect(ctx, transport); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	result, err := client.ReadResource(ctx, "docs://documents")
	if err != nil {
		log.Fatalf("Error calling resource: %v", err)
	}

	if len(result.Contents) == 0 {
		log.Fatalf("No documents found")
	}

	var ids []string
	if err := json.Unmarshal([]byte(result.Contents[0].Text), &ids); err != nil {
		log.Fatalf("Error while unmarshaling document IDs: %v", err)
	}

	if len(ids) > 0 {
		fmt.Printf("List of Documents: \n%s\n", strings.Join(ids, "\n"))
	}
}
