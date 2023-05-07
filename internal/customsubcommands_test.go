package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCustomSubcommands(t *testing.T) {
	app := testNewApp(t)
	r, err := getRepository(app)
	assert.NoError(t, err)

	subcmd := filepath.Join(r.PathResolver.LibExecDir(), "gptx-helloworld")
	err = os.WriteFile(subcmd, []byte("#!/bin/sh\necho hello world"), 0755)
	assert.NoError(t, err)

	err = app.Run([]string{"gptx", "helloworld"})
	assert.NoError(t, err)

	out := app.Writer.(*bytes.Buffer).String()
	assert.Equal(t, "hello world", strings.TrimSpace(out))
}
