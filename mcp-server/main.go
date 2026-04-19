package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var docs = map[string]string{
	"deposition.md":   "This deposition covers the testimony of Angela Smith, P.E.",
	"report.pdf":      "The report details the state of a 20m condenser tower.",
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
