package internal

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConfig_LoadFromFile(t *testing.T) {
	t.Run("Load from file", func(t *testing.T) {
		tempFile := testTempFile(t, []byte(`
openai_api_key = "test-key"
model = "test-model"
max_cache_length = 123

# arbitrary keys
v1 = "bar"
v2 = 123
`))
		c := NewConfig()
		err := c.LoadFromFile(tempFile.Name())
		assert.NoError(t, err)

		assert.Equal(t, "test-key", c.OpenAIAPIKey)
		assert.Equal(t, "test-model", c.Model)
		assert.Equal(t, 123, c.MaxCacheLength)
		assert.Equal(t, "bar", c.m["v1"])
		assert.Equal(t, int64(123), c.m["v2"])
	})

	t.Run("No file", func(t *testing.T) {
		c := NewConfig()
		err := c.LoadFromFile("/tmp/file-not-found")
		assert.Error(t, err)
	})
}

func TestConfig_MarshalJSON(t *testing.T) {
	c := NewConfig()
	c.OpenAIAPIKey = "test-key"
	c.Model = "test-model"
	c.MaxCacheLength = 100
	c.m["v1"] = "bar"
	c.m["v2"] = 123

	buf, err := json.Marshal(c)
	assert.NoError(t, err)
	assert.JSONEq(t, strings.TrimPrefix(`
{
  "openai_api_key": "test-key",
  "model": "test-model",
  "max_cache_length": 100,
  "v1": "bar",
  "v2": 123
}`, "\n"), string(buf))

}
