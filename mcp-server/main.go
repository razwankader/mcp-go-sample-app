package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var docs = map[string]string{
	"deposition.md":   "This deposition covers the testimony of Angela Smith, P.E.",
	"report.pdf":      "Several third party Go MCP SDKs inspired the development and design of this official SDK, and continue to be viable alternatives, notably mcp-go, originally authored by Ed Zynda. We are grateful to Ed as well as the other contributors to mcp-go, and to authors and contributors of other SDKs such as mcp-golang and go-mcp. Thanks to their work, there is a thriving ecosystem of Go MCP clients and servers.",
	"financials.docx": "These financials outline the project's budget and expenditures.",
	"outlook.pdf":     "This document presents the projected future performance of the system.",
	"plan.md":         "The plan outlines the steps for the project's implementation.",
	"spec.txt":        "These specifications define the technical requirements for the equipment.",
}

type readDocInput struct {
	Name string `json:"name" jsonschema:"the name of the document to read"`
}

type readDocOutput struct {
	Content string `json:"content" jsonschema:"the content of the document"`
}

// readDocument is a tool that reads the content of a document given its name.
func readDocument(ctx context.Context, req *mcp.CallToolRequest, input readDocInput) (*mcp.CallToolResult, readDocOutput, error) {
	if content, exists := docs[input.Name]; exists {
		return nil, readDocOutput{Content: content}, nil
	}
	return nil, readDocOutput{}, fmt.Errorf("document not found: %s", input.Name)
}

type editDocInput struct {
	Name      string `json:"name" jsonschema:"the name of the document to edit"`
	OldString string `json:"old_string" jsonschema:"the existing content of the document"`
	NewString string `json:"new_string" jsonschema:"the new content of the document"`
}

type editDocOutput struct {
	Success bool `json:"success" jsonschema:"whether the document was successfully edited"`
}

// editDocument is a tool that edits the content of a document given its name, the existing content, and the new content.
func editDocument(ctx context.Context, req *mcp.CallToolRequest, input editDocInput) (*mcp.CallToolResult, editDocOutput, error) {
	if content, exists := docs[input.Name]; exists {
		docs[input.Name] = strings.ReplaceAll(content, input.OldString, input.NewString)
		return nil, editDocOutput{Success: true}, nil
	}

	return nil, editDocOutput{Success: false}, fmt.Errorf("document not found: %s", input.Name)
}

// listDocuments is a resource that lists all available documents.
func listDocuments(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	var documentNames []string
	for name := range docs {
		documentNames = append(documentNames, name)
	}

	content, err := json.Marshal(documentNames)
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{
			URI:      req.Params.URI,
			MIMEType: "application/json",
			Text:     string(content),
		}},
	}, nil
}

func getDocumentContents(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	docName := strings.TrimPrefix(req.Params.URI, "docs://documents/")

	if content, exists := docs[docName]; exists {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:      req.Params.URI,
				MIMEType: "text/plain",
				Text:     content,
			}},
		}, nil
	}

	return nil, fmt.Errorf("document not found: %s", docName)
}

func formatDoccumentPrompt(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	docName, ok := req.Params.Arguments["doc_name"]
	if !ok {
		return nil, fmt.Errorf("missing or invalid doc_name parameter")
	}

	prompt := fmt.Sprintf(`Your goal is to reformat a document to be written with markdown syntax.

	The id of the document you need to reformat is:
	<document_id>
	%s
	</document_id>

	Add in headers, bullet points, tables, etc as necessary. Feel free to add in structure.
	Use the 'edit_document' tool to edit the document. After the document has been reformatted...
	`, docName)

	return &mcp.GetPromptResult{
		Messages: []*mcp.PromptMessage{{
			Role:    mcp.Role("user"),
			Content: &mcp.TextContent{Text: prompt},
		}},
	}, nil
}

type summarizeDocInput struct {
	Name string `json:"name" jsonschema:"the name of the document to summerize"`
}

type summarizeDocOutput struct {
	SummerizeContent string `json:"content" jsonschema:"the summerize content of the document"`
}

// summarizeDocumentContent is a tool that summarizes the content of a document given its name.
// It uses Sampling feature to call the model through MCP client
func summarizeDocumentContent(ctx context.Context, req *mcp.CallToolRequest, input summarizeDocInput) (*mcp.CallToolResult, summarizeDocOutput, error) {
	req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
		Message:  "Summarize Progress",
		Progress: 0.3,
		Total:    1.0,
	})
	req.Session.Log(ctx, &mcp.LoggingMessageParams{
		Data:  fmt.Sprintf("Fetching %s content to summarize", input.Name),
		Level: "info",
	})

	content, exists := docs[input.Name]
	if !exists {
		return nil, summarizeDocOutput{}, fmt.Errorf("document not found: %s", input.Name)
	}

	time.Sleep(5 * time.Second) // simulate a long running process of summarization

	prompt := fmt.Sprintf(`Please summarize the following text: %s`, content)

	req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
		Message:  "Summarize Progress",
		Progress: 0.6,
		Total:    1.0,
	})
	req.Session.Log(ctx, &mcp.LoggingMessageParams{
		Data:  "Sending prompt to mcp client for calling Claude through sampling feature",
		Level: "info",
	})

	result, err := req.Session.CreateMessage(ctx, &mcp.CreateMessageParams{
		Messages: []*mcp.SamplingMessage{
			{
				Role: mcp.Role("user"),
				Content: &mcp.TextContent{
					Text: prompt,
				},
			},
		},
		MaxTokens:    4000,
		SystemPrompt: "You are a helpful research assistant.",
	})
	if err != nil {
		return nil, summarizeDocOutput{}, err
	}

	time.Sleep(5 * time.Second) // simulate a long running process of summarization

	var summaryText string
	textContent, ok := result.Content.(*mcp.TextContent)
	if ok {
		summaryText = textContent.Text
	} else {
		return nil, summarizeDocOutput{}, fmt.Errorf("unexpected content type: %T", result.Content)
	}

	req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
		Message:  "Summarize Progress",
		Progress: 1.0,
		Total:    1.0,
	})
	req.Session.Log(ctx, &mcp.LoggingMessageParams{
		Data:  fmt.Sprintf("Received response from mcp client: \n %s", summaryText),
		Level: "info",
	})

	return nil, summarizeDocOutput{SummerizeContent: summaryText}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sample-mcp-server",
		Version: "v1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "read",
		Description: "Reads the content of a document given its name",
	}, readDocument)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "edit",
		Description: "Edits the content of a document given its name, old string and new string",
	}, editDocument)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "summarize",
		Description: "Summarize the content of a document given its name through the sampling feature",
	}, summarizeDocumentContent)

	server.AddResource(&mcp.Resource{
		Name:        "list",
		Description: "Lists all available documents",
		URI:         "docs://documents",
		MIMEType:    "application/json",
	}, listDocuments)

	server.AddResourceTemplate(&mcp.ResourceTemplate{
		Name:        "get",
		Description: "Gets the content of a document given its name",
		URITemplate: "docs://documents/{doc_name}",
		MIMEType:    "application/json",
	}, getDocumentContents)

	server.AddPrompt(&mcp.Prompt{
		Name:        "format",
		Description: "Rewrites the contents of the document in Markdown format",
		Arguments: []*mcp.PromptArgument{{
			Name:        "doc_name",
			Description: "The name of the document to reformat",
			Required:    true,
		}},
	}, formatDoccumentPrompt)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
