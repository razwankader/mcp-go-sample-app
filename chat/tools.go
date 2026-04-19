package chat

import (
	"context"
	"encoding/json"
	"fmt"

	mcpclient "mcp-go-sample-app/mcp-client"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func getAllTools(ctx context.Context, clients map[string]*mcpclient.MCPClient) ([]anthropic.ToolUnionParam, error) {
	var tools []anthropic.ToolUnionParam
	for _, client := range clients {
		result, err := client.ListTools(ctx)
		if err != nil {
			return nil, err
		}
		for _, t := range result.Tools {
			tools = append(tools, mcpToolToAnthropicTool(t))
		}
	}
	return tools, nil
}

func mcpToolToAnthropicTool(t *mcp.Tool) anthropic.ToolUnionParam {
	// Marshal the MCP schema and extract properties for Anthropic's format
	schemaJSON, _ := json.Marshal(t.InputSchema)
	var schema struct {
		Properties map[string]any `json:"properties"`
	}
	json.Unmarshal(schemaJSON, &schema)
	props := schema.Properties
	if props == nil {
		props = map[string]any{}
	}

	tool := anthropic.ToolParam{
		Name:        t.Name,
		Description: anthropic.String(t.Description),
		InputSchema: anthropic.ToolInputSchemaParam{
			Properties: props,
		},
	}
	return anthropic.ToolUnionParam{OfTool: &tool}
}

func findClientWithTool(ctx context.Context, clients map[string]*mcpclient.MCPClient, toolName string) *mcpclient.MCPClient {
	for _, client := range clients {
		result, err := client.ListTools(ctx)
		if err != nil {
			continue
		}
		for _, t := range result.Tools {
			if t.Name == toolName {
				return client
			}
		}
	}
	return nil
}

func executeToolRequests(ctx context.Context, clients map[string]*mcpclient.MCPClient, msg *anthropic.Message) ([]anthropic.ContentBlockParamUnion, error) {
	var results []anthropic.ContentBlockParamUnion

	for _, block := range msg.Content {
		variant, ok := block.AsAny().(anthropic.ToolUseBlock)
		if !ok {
			continue
		}

		toolUseID := block.ID
		toolName := variant.Name
		toolInputRaw := variant.JSON.Input.Raw()

		client := findClientWithTool(ctx, clients, toolName)
		if client == nil {
			results = append(results, anthropic.NewToolResultBlock(toolUseID, "Could not find that tool", true))
			continue
		}

		var inputMap map[string]any
		if err := json.Unmarshal([]byte(toolInputRaw), &inputMap); err != nil {
			results = append(results, anthropic.NewToolResultBlock(toolUseID, "Invalid tool input", true))
			continue
		}

		output, err := client.CallTool(ctx, toolName, inputMap)
		if err != nil {
			results = append(results, anthropic.NewToolResultBlock(toolUseID, fmt.Sprintf("Error: %v", err), true))
			continue
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
		content, _ := json.Marshal(textParts)
		results = append(results, anthropic.NewToolResultBlock(toolUseID, string(content), output.IsError))
	}

	return results, nil
}
