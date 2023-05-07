package internal

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/sashabaranov/go-openai"
)

var initialConfig = trimLeftSpaces(`
# This is a Gptx configuration file.

# OpenAI API Key. You can override this value by using the OPENAI_API_KEY environment variable.
openai_api_key = ""

# Default model for Chat API.
model = "gpt-3.5-turbo"

# Maximum number of cached responses.
max_cache_length = 100
`)

type Config struct {
	OpenAIAPIKey   string                 `toml:"openai_api_key"`   // OpenAI API Key
	Model          string                 `toml:"model"`            // Default setting for https://platform.openai.com/docs/api-reference/chat/create#chat/create-model
	MaxCacheLength int                    `toml:"max_cache_length"` // The maximum number of cached responses.
	m              map[string]interface{} `toml:"-"`                // This is an internal representation of Config for holding arbitrary keys.
}

func NewConfig() *Config {
	return &Config{
		OpenAIAPIKey:   "",
		Model:          openai.GPT3Dot5Turbo,
		MaxCacheLength: 100,
		m:              make(map[string]interface{}),
	}
}

func (c *Config) LoadFromFile(path string) error {
	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}
	if _, err := toml.DecodeFile(path, &c.m); err != nil {
		return nil
	}
	return nil
}

func (c *Config) MarshalJSON() ([]byte, error) {
	// copy map
	m := make(map[string]interface{}, len(c.m))
	for key, value := range c.m {
		m[key] = value
	}

	// override built-in keys
	m["openai_api_key"] = c.OpenAIAPIKey
	m["model"] = c.Model
	m["max_cache_length"] = c.MaxCacheLength

	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
