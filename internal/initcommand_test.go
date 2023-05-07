package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCommand(t *testing.T) {
	t.Run("init", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "gptx_test_")
		if err != nil {
			t.Fatal(err)
		}
		// remove it because we want to test if it is created by init command
		_ = os.RemoveAll(dir)
		t.Cleanup(func() {
			_ = os.RemoveAll(dir)
		})

		_ = os.Setenv(AppHomeEnvKey, dir)

		app := testNewApp(t)
		err = app.Run([]string{"gptx", "init"})
		assert.NoError(t, err)
		assert.DirExists(t, dir)
		assert.FileExists(t, filepath.Join(dir, "config.toml"))
	})
}
