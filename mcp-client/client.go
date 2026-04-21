package mcpclient

import (
	"context"
	"fmt"
	"mcp-go-sample-app/config"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
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
		CreateMessageHandler: SamplingCallback,
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

func SamplingCallback(ctx context.Context, req *mcp.CreateMessageRequest) (*mcp.CreateMessageResult, error) {
	fmt.Println("balllllllll")

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

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(cfg.Claude.Model),
		MaxTokens: 8000,
		Messages:  claudeMessages,
	}

	claudeClient := anthropic.NewClient(option.WithAPIKey(cfg.Anthropic.APIKey))
	resp, err := claudeClient.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("calling Claude: %w", err)
	}

	text := TextFromMessage(resp)

	return &mcp.CreateMessageResult{
		Content: &mcp.TextContent{
			Text: text,
		},
	}, nil
}

func TextFromMessage(msg *anthropic.Message) string {
	var parts []string
	for _, block := range msg.Content {
		if tb, ok := block.AsAny().(anthropic.TextBlock); ok {
			parts = append(parts, tb.Text)
		}
	}
	return strings.Join(parts, "\n")
}
