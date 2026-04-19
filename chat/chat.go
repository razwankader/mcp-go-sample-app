package chat

import (
	"context"
	"fmt"
	mcpclient "mcp-go-sample-app/mcp-client"

	"github.com/anthropics/anthropic-sdk-go"
)

type Chat struct {
	Claude   *Claude
	Clients  map[string]*mcpclient.MCPClient
	Messages []anthropic.MessageParam
}

func NewChat(claude *Claude, clients map[string]*mcpclient.MCPClient) *Chat {
	return &Chat{
		Claude:  claude,
		Clients: clients,
	}
}

// RunMessages runs the agentic loop using the current Messages slice.
// Caller is responsible for appending the initial user message before calling.
func (c *Chat) RunMessages(ctx context.Context) (string, error) {
	for {
		tools, err := getAllTools(ctx, c.Clients)
		if err != nil {
			return "", fmt.Errorf("listing tools: %w", err)
		}

		resp, err := c.Claude.Chat(ctx, c.Messages, "", tools)
		if err != nil {
			return "", fmt.Errorf("calling Claude: %w", err)
		}

		c.Messages = append(c.Messages, resp.ToParam())

		if resp.StopReason != anthropic.StopReasonToolUse {
			return TextFromMessage(resp), nil
		}

		if text := TextFromMessage(resp); text != "" {
			fmt.Println(text)
		}

		toolResults, err := executeToolRequests(ctx, c.Clients, resp)
		if err != nil {
			return "", fmt.Errorf("executing tools: %w", err)
		}

		c.Messages = append(c.Messages, anthropic.NewUserMessage(toolResults...))
	}
}
