package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Anthropic AnthropicConfig
	Claude    ClaudeConfig
}

type AnthropicConfig struct {
	APIKey string
}

type ClaudeConfig struct {
	Model              string
	MaxTokens          int64
	ItineraryMaxTokens int64
}

// Load reads config.yaml (from the working directory) and overlays
// environment variables. ANTHROPIC_API_KEY must be set in the environment.
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config.yaml: %w", err)
	}

	v.SetEnvPrefix("")
	v.AutomaticEnv()
	v.BindEnv("anthropic.api_key", "ANTHROPIC_API_KEY")

	return &Config{
		Anthropic: AnthropicConfig{
			APIKey: v.GetString("anthropic.api_key"),
		},
		Claude: ClaudeConfig{
			Model:              v.GetString("claude.model"),
			MaxTokens:          v.GetInt64("claude.max_tokens"),
			ItineraryMaxTokens: v.GetInt64("claude.itinerary_max_tokens"),
		},
	}, nil
}
