package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenameCommand(t *testing.T) {
	t.Run("rename missing argument(s)", func(t *testing.T) {
		app := testNewApp(t)
		err := app.Run([]string{"gptx", "rename"})
		assert.Error(t, err)
		assert.Equal(t, "missing argument(s)", err.Error())
	})

	t.Run("rename", func(t *testing.T) {
		app := testNewApp(t)
		r, err := getRepository(app)
		assert.NoError(t, err)

		co := NewConversation()
		co.Name = "conversation001"

		s, err := r.StoreManager.Open()
		assert.NoError(t, err)
		err = s.CreateConversation(co)
		assert.NoError(t, err)

		err = app.Run([]string{"gptx", "rename", "1", "conversation001-updated"})
		assert.NoError(t, err)

		s, err = r.StoreManager.Open()
		assert.NoError(t, err)
		co2, err := s.GetConversationById(1)
		assert.NoError(t, err)
		assert.Equal(t, "conversation001-updated", co2.Name)

		co3, err := s.GetConversationByName("conversation001-updated")
		assert.NoError(t, err)
		assert.Equal(t, "conversation001-updated", co3.Name)

		_, err = s.GetConversationByName("conversation001")
		assert.Error(t, err)
		assert.IsType(t, &ConversationNotFoundError{}, err)

	})
}
