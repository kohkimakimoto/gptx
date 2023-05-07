package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveCommand(t *testing.T) {
	t.Run("remove missing conversation id", func(t *testing.T) {
		app := testNewApp(t)
		err := app.Run([]string{"gptx", "remove"})
		assert.Error(t, err)
		assert.Equal(t, "missing conversation id", err.Error())
	})

	t.Run("remove", func(t *testing.T) {
		app := testNewApp(t)
		r, err := getRepository(app)
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "conversation001"

		s, err := r.StoreManager.Open()
		assert.NoError(t, err)
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		err = app.Run([]string{"gptx", "remove", "1"})
		assert.NoError(t, err)

		s, err = r.StoreManager.Open()
		assert.NoError(t, err)
		_, err = s.GetConversationById(1)
		assert.Error(t, err)
		assert.IsType(t, &ConversationNotFoundError{}, err)
	})

}
