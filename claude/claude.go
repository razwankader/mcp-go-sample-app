package claude

import (
	"context"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Claude struct {
	client anthropic.Client
	model  string
}

func NewClaude(apiKey string, model string) *Claude {
	return &Claude{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:  model,
	}
}

func (c *Claude) Chat(ctx context.Context, messages []anthropic.MessageParam, system string, tools []anthropic.ToolUnionParam) (*anthropic.Message, error) {
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: 8000,
		Messages:  messages,
	}
	if system != "" {
		params.System = []anthropic.TextBlockParam{{
			Text:         system,
			CacheControl: anthropic.NewCacheControlEphemeralParam(),
		}}
	}
	if len(tools) > 0 {
		params.Tools = tools
	}
	return c.client.Messages.New(ctx, params)
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
