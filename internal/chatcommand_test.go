package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestChatCommand(t *testing.T) {
	t.Run("chat", func(t *testing.T) {
		// simple chat request
		app := testNewApp(t)
		r, err := getRepository(app)
		assert.NoError(t, err)

		// override config
		r.Config.OpenAIAPIKey = "sk-dummykey..."
		// override http client
		r.ClientConfig.HTTPClient = testHttpClient(t, func(req *http.Request) *http.Response {
			// see also: https://platform.openai.com/docs/api-reference/chat/create
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(bytes.NewBuffer([]byte(strings.TrimPrefix(`
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "\n\nHello there, how may I assist you today?"
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}
`, "\n")))),
				Header: make(http.Header),
			}
		})

		err = app.Run([]string{"gptx", "chat", "Hello!"})
		assert.NoError(t, err)
		out := app.Writer.(*bytes.Buffer).String()
		assert.Contains(t, out, "Hello there, how may I assist you today?")
	})
}

// TODO: add more tests
