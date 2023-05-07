package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	app := testNewApp(t)
	err := app.Run([]string{"gptx", "version"})
	assert.NoError(t, err)

	ret := app.Writer.(*bytes.Buffer).String()
	assert.Equal(t, strings.TrimPrefix(`
Version:       0.0.0
Commit Hash:   unknown
`, "\n"), ret)
}
