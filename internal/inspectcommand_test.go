package internal

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInspectCommand(t *testing.T) {
	t.Run("inspect missing conversation argument(s)", func(t *testing.T) {
		app := testNewApp(t)
		err := app.Run([]string{"gptx", "inspect"})
		assert.Error(t, err, "missing conversation argument(s)")
	})

	t.Run("inspect", func(t *testing.T) {
		app := testNewApp(t)
		r, err := getRepository(app)
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "conversation001"

		s, err := r.StoreManager.Open()
		assert.NoError(t, err)
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		err = app.Run([]string{"gptx", "inspect", "1"})
		assert.NoError(t, err)

		buf, _ := json.Marshal(co)
		out := app.Writer.(*bytes.Buffer).String()
		assert.JSONEq(t, string(buf), out)
	})

	t.Run("inspect pretty print", func(t *testing.T) {
		app := testNewApp(t)
		r, err := getRepository(app)
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "conversation001"

		s, err := r.StoreManager.Open()
		assert.NoError(t, err)
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		err = app.Run([]string{"gptx", "inspect", "--pretty", "1"})
		assert.NoError(t, err)

		buf, _ := json.MarshalIndent(co, "", "  ")
		out := app.Writer.(*bytes.Buffer).String()
		assert.JSONEq(t, string(buf), out)
	})
}
