package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestHookFactory_resolveHookCommandPath(t *testing.T) {
	cases := []struct {
		name     string
		expected string
	}{
		{
			name:     "gptx-hook-example",
			expected: "gptx-hook-example",
		},
		{
			name:     "example",
			expected: "gptx-hook-example",
		},
		{
			name:     "/path/to/custom/example",
			expected: "/path/to/custom/example",
		},
	}

	for _, c := range cases {
		f := &HookFactory{}
		assert.Equal(t, c.expected, f.resolveHookCommandPath(c.name))
	}
}

func TestHook_Command(t *testing.T) {
	app := testNewApp(t)
	r, err := getRepository(app)
	assert.NoError(t, err)
	err = updatePathEnv(r.PathResolver)
	assert.NoError(t, err)

	hookFile := filepath.Join(r.PathResolver.LibExecDir(), "gptx-hook-example")
	err = os.WriteFile(hookFile, []byte("#!/bin/sh\necho hello"), 0755)
	assert.NoError(t, err)

	f := &HookFactory{}
	hook, err := f.NewHook("example")
	assert.NoError(t, err)

	cmd := hook.Command()
	assert.Equal(t, hookFile, cmd.Path)
}
