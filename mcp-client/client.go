package mcpclient

import (
	"context"
	"fmt"
	"log"
	"mcp-go-sample-app/claude"
	"mcp-go-sample-app/config"

	"github.com/anthropics/anthropic-sdk-go"
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
	}, &mcp.ClientOptions{
		CreateMessageHandler:        SamplingCallback,
		LoggingMessageHandler:       LoggingCallback,
		ProgressNotificationHandler: ProgressCallback,
	})

	return &MCPClient{
		client: client,
	}
}

func (c *MCPClient) Connect(ctx context.Context, transport *mcp.CommandTransport) error {
	session, err := c.client.Connect(ctx, transport, nil)
	if err != nil {
		return err
	}

	err = session.SetLoggingLevel(ctx, &mcp.SetLoggingLevelParams{
		Level: "info",
	})
	if err != nil {
		return err
	}

	c.session = session
	return nil
}

func (c *MCPClient) AddRoots(ctx context.Context, roots []*mcp.Root) {
	c.client.AddRoots(roots...)
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

func SamplingCallback(ctx context.Context, req *mcp.CreateMessageRequest) (*mcp.CreateMessageResult, error) {
	var claudeMessages []anthropic.MessageParam
	for _, m := range req.Params.Messages {
		textContent, ok := m.Content.(*mcp.TextContent)
		if !ok {
			return nil, fmt.Errorf("unsupported message content type: %T", m.Content)
		}

		if m.Role == mcp.Role("user") {
			claudeMessages = append(claudeMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(textContent.Text)))
		} else {
			claudeMessages = append(claudeMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(textContent.Text)))
		}
	}

	cfg := config.Get()

	claudeClient := claude.NewClaude(cfg.Anthropic.APIKey, cfg.Claude.Model)
	resp, err := claudeClient.Chat(ctx, claudeMessages, "", nil)
	if err != nil {
		return nil, fmt.Errorf("calling Claude: %w", err)
	}

	return &mcp.CreateMessageResult{
		Content: &mcp.TextContent{
			Text: claude.TextFromMessage(resp),
		},
	}, nil
}

func LoggingCallback(ctx context.Context, req *mcp.LoggingMessageRequest) {
	log.Printf("Log from server: %s\n", req.Params.Data)
}

func ProgressCallback(ctx context.Context, req *mcp.ProgressNotificationClientRequest) {
	percentage := (req.Params.Progress / req.Params.Total) * 100
	log.Printf("%s %.f%%\n", req.Params.Message, percentage)
}
