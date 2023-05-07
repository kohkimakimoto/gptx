package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanCommand(t *testing.T) {
	t.Run("clean", func(t *testing.T) {
		// just run the command
		app := testNewApp(t)
		err := app.Run([]string{"gptx", "clean"})
		assert.NoError(t, err)
	})
}
