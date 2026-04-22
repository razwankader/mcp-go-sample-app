package cmd

import (
	"log"
	"mcp-go-sample-app/config"

	"github.com/spf13/cobra"
)

var (
	// rootCmd is the root command of backup service
	rootCmd = &cobra.Command{
		Use:   "mcp-go-sample-app",
		Short: "MCP Go Sample App demonstrates how to build an app using the Model Context Protocol (MCP) with a local CLI interface.",
		Long:  "MCP Go Sample App demonstrates how to build an app using the Model Context Protocol (MCP) with a local CLI interface.",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func initConfig() {
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}

	cfg := config.Get()
	if cfg.Claude.Model == "" {
		log.Fatalln("Error: CLAUDE_MODEL cannot be empty. Update config")
	}

	if cfg.Anthropic.APIKey == "" {
		log.Fatalln("Error: ANTHROPIC_API_KEY cannot be empty. Update config")
	}
}
