package mcpclient

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPClient struct {
	client  *mcp.Client
	session *mcp.ClientSession
}

func New() *MCPClient {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "sample-mcp-client",
		Version: "v1.0.0",
	}, nil)

	return &MCPClient{
		client: client,
	}
}

func (c *MCPClient) Connect(ctx context.Context, transport *mcp.CommandTransport) error {
	session, err := c.client.Connect(ctx, transport, nil)
	if err != nil {
		return err
	}
	c.session = session
	return nil
}

func (c *MCPClient) ListTools(ctx context.Context) (*mcp.ListToolsResult, error) {
	return c.session.ListTools(ctx, &mcp.ListToolsParams{})
}

func (c *MCPClient) CallTool(ctx context.Context, name string, arguments map[string]any) (*mcp.CallToolResult, error) {
	return c.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: arguments,
	})
}

func (c *MCPClient) ReadResource(ctx context.Context, uri string) (*mcp.ReadResourceResult, error) {
	return c.session.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: uri,
	})
}

func (c *MCPClient) ListPrompts(ctx context.Context) (*mcp.ListPromptsResult, error) {
	return c.session.ListPrompts(ctx, &mcp.ListPromptsParams{})
}

func (c *MCPClient) GetPrompt(ctx context.Context, name string, arguments map[string]string) (*mcp.GetPromptResult, error) {
	return c.session.GetPrompt(ctx, &mcp.GetPromptParams{
		Name:      name,
		Arguments: arguments,
	})
}

func (c *MCPClient) Close() error {
	if c.session != nil {
		return c.session.Close()
	}
	return nil
}
