package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewApp(t *testing.T) {
	// just run the command
	a := testNewApp(t)
	err := a.Run([]string{"gptx"})
	assert.NoError(t, err)
}
