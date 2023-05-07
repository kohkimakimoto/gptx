package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConfigCommand(t *testing.T) {
	t.Run("config", func(t *testing.T) {
		app := testNewApp(t)
		r, err := getRepository(app)
		if err != nil {
			t.Fatal(err)
		}
		r.Config = &Config{
			OpenAIAPIKey:   "sk-1234567890",
			Model:          "test_model",
			MaxCacheLength: 123,
		}
		err = app.Run([]string{"gptx", "config"})
		assert.NoError(t, err)

		ret := app.Writer.(*bytes.Buffer).String()
		assert.JSONEq(t, strings.TrimPrefix(`
{
  "openai_api_key": "sk-1234567890",
  "model": "test_model",
  "max_cache_length": 123
}
`, "\n"), ret)
	})

	t.Run("config pretty print", func(t *testing.T) {
		app := testNewApp(t)
		r, err := getRepository(app)
		if err != nil {
			t.Fatal(err)
		}
		r.Config = &Config{
			OpenAIAPIKey:   "sk-1234567890",
			Model:          "test_model",
			MaxCacheLength: 123,
		}
		err = app.Run([]string{"gptx", "config", "--pretty"})
		assert.NoError(t, err)

		ret := app.Writer.(*bytes.Buffer).String()
		// just check if the output is valid JSON (not pretty printed)
		assert.JSONEq(t, strings.TrimPrefix(`
{
  "openai_api_key": "sk-1234567890",
  "model": "test_model",
  "max_cache_length": 123
}
`, "\n"), ret)
	})

}
