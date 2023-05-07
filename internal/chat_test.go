package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChatService_DisableOutputAnimation(t *testing.T) {
	app := testNewApp(t)
	r, err := getRepository(app)
	assert.NoError(t, err)
	c, err := r.NewChatService(app.Writer)
	assert.NoError(t, err)
	c.DisableOutputAnimation()

	assert.Equal(t, false, c.Writer.UseAnimation)
}
