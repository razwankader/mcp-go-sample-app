package chat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type CliApp struct {
	agent *CliChat
}

func NewCliApp(agent *CliChat) *CliApp {
	return &CliApp{agent: agent}
}

func (a *CliApp) initialize(ctx context.Context) {
	ids, err := a.agent.ListDocIDs(ctx)
	if err == nil && len(ids) > 0 {
		fmt.Printf("In-Memory Documents: %s\n\n", strings.Join(ids, ", "))
		fmt.Println("Tips:\n 1. Read or edit document content by mentioning document name with @<docname>\n\n 2. Run prompts with /format <docname>")
	}

	prompts, err := a.agent.ListPrompts(ctx)
	if err == nil && len(prompts) > 0 {
		var names []string
		for _, p := range prompts {
			names = append(names, "/"+p.Name)
		}
	}
	fmt.Println()
}

func (a *CliApp) Run(ctx context.Context) error {
	a.initialize(ctx)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println()
			return nil
		}
		if err != nil {
			return err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			return nil
		}

		response, err := a.agent.Run(ctx, input)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		fmt.Printf("\nResponse:\n%s\n\n", response)
	}
}
