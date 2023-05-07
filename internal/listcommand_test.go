package internal

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListCommand(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		app := testNewApp(t)
		r, err := getRepository(app)
		assert.NoError(t, err)

		s, err := r.StoreManager.Open()
		assert.NoError(t, err)

		for i := 1; i <= 30; i++ {
			co := NewConversation()
			co.Name = fmt.Sprintf("test-conversation-%d", i)
			err := s.CreateConversation(co)
			assert.NoError(t, err)
		}

		err = app.Run([]string{"gptx", "list"})
		assert.NoError(t, err)
		out := app.Writer.(*bytes.Buffer).String()
		// t.Log(out)
		// just check the header line
		assert.Regexp(t, `^ID\s+PROMPT\s+MESSAGES\s+NAME\s+LABEL\s+HOOKS\s+CREATED\s+ELAPSED`, out)
	})
}
