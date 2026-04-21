package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"mcp-go-sample-app/claude"
	mcpclient "mcp-go-sample-app/mcp-client"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CliChat struct {
	Chat
	docClient *mcpclient.MCPClient
}

func NewCliChat(docClient *mcpclient.MCPClient, clients map[string]*mcpclient.MCPClient, claude *claude.Claude) *CliChat {
	return &CliChat{
		Chat:      Chat{Claude: claude, Clients: clients},
		docClient: docClient,
	}
}

func (cc *CliChat) ListPrompts(ctx context.Context) ([]*mcp.Prompt, error) {
	result, err := cc.docClient.ListPrompts(ctx)
	if err != nil {
		return nil, err
	}
	return result.Prompts, nil
}

func (cc *CliChat) ListDocIDs(ctx context.Context) ([]string, error) {
	result, err := cc.docClient.ReadResource(ctx, "docs://documents")
	if err != nil {
		return nil, err
	}
	if len(result.Contents) == 0 {
		return nil, nil
	}
	var ids []string
	if err := json.Unmarshal([]byte(result.Contents[0].Text), &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (cc *CliChat) GetDocContent(ctx context.Context, docID string) (string, error) {
	result, err := cc.docClient.ReadResource(ctx, "docs://documents/"+docID)
	if err != nil {
		return "", err
	}
	if len(result.Contents) == 0 {
		return "", nil
	}
	return result.Contents[0].Text, nil
}

func (cc *CliChat) Run(ctx context.Context, query string) (string, error) {
	if err := cc.processQuery(ctx, query); err != nil {
		return "", err
	}
	return cc.Chat.RunMessages(ctx)
}

func (cc *CliChat) processQuery(ctx context.Context, query string) error {
	if strings.HasPrefix(query, "/") {
		return cc.processCommand(ctx, query)
	}
	return cc.processWithResources(ctx, query)
}

func (cc *CliChat) processCommand(ctx context.Context, query string) error {
	words := strings.Fields(query)
	if len(words) < 2 {
		return fmt.Errorf("usage: /<command> <doc_name>")
	}
	command := strings.TrimPrefix(words[0], "/")
	docID := words[1]

	result, err := cc.docClient.GetPrompt(ctx, command, map[string]string{"doc_name": docID})
	if err != nil {
		return fmt.Errorf("getting prompt %q: %w", command, err)
	}

	for _, msg := range result.Messages {
		cc.Messages = append(cc.Messages, convertPromptMessage(msg))
	}
	return nil
}

func (cc *CliChat) processWithResources(ctx context.Context, query string) error {
	mentions := extractMentions(query)

	var docContexts []string
	for _, mention := range mentions {
		content, err := cc.GetDocContent(ctx, mention)
		if err != nil {
			continue
		}
		docContexts = append(docContexts, fmt.Sprintf("<document id=%q>\n%s\n</document>", mention, content))
	}

	ctx_ := strings.Join(docContexts, "\n")
	prompt := buildQueryPrompt(query, ctx_)
	cc.Messages = append(cc.Messages, anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)))
	return nil
}

func extractMentions(query string) []string {
	var mentions []string
	for _, word := range strings.Fields(query) {
		if strings.HasPrefix(word, "@") {
			mentions = append(mentions, strings.TrimPrefix(word, "@"))
		}
	}
	return mentions
}

func buildQueryPrompt(query, context string) string {
	return fmt.Sprintf(`The user has a question:
<query>
%s
</query>

The following context may be useful in answering their question:
<context>
%s
</context>

Note the user's query might contain references to documents like "@report.docx". The "@" is only
included as a way of mentioning the doc. The actual name of the document would be "report.docx".
If the document content is included in this prompt, you don't need to use an additional tool to read the document.
Answer the user's question directly and concisely. Start with the exact information they need.
Don't refer to or mention the provided context in any way - just use it to inform your answer.`, query, context)
}

func convertPromptMessage(msg *mcp.PromptMessage) anthropic.MessageParam {
	role := string(msg.Role)

	text := extractTextFromContent(msg.Content)

	if role == "user" {
		return anthropic.NewUserMessage(anthropic.NewTextBlock(text))
	}
	return anthropic.NewAssistantMessage(anthropic.NewTextBlock(text))
}

func extractTextFromContent(content any) string {
	if tc, ok := content.(*mcp.TextContent); ok {
		return tc.Text
	}
	b, err := json.Marshal(content)
	if err != nil {
		return ""
	}
	var m struct {
		Text string `json:"text"`
	}
	json.Unmarshal(b, &m)
	return m.Text
}
