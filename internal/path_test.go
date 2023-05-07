package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestGetAppHomeDir(t *testing.T) {
	t.Run("GPTX_HOME is not set", func(t *testing.T) {
		originalHome := os.Getenv("GPTX_HOME")
		t.Cleanup(func() {
			os.Setenv("GPTX_HOME", originalHome)
		})
		os.Setenv("GPTX_HOME", "")
		assert.Equal(t, filepath.Join(os.Getenv("HOME"), ".gptx"), getAppHomeDir())
	})

	t.Run("GPTX_HOME is set", func(t *testing.T) {
		originalHome := os.Getenv("GPTX_HOME")
		t.Cleanup(func() {
			os.Setenv("GPTX_HOME", originalHome)
		})
		os.Setenv("GPTX_HOME", "/tmp/to/.gptx")
		assert.Equal(t, "/tmp/to/.gptx", getAppHomeDir())
	})
}

func TestPathResolver_ConfigFilePath(t *testing.T) {
	pr := NewPathResolver("/tmp/gptx")
	assert.Equal(t, filepath.Join("/tmp/gptx", "config.toml"), pr.ConfigFilePath())
}

func TestPathResolver_DBFilePath(t *testing.T) {
	pr := NewPathResolver("/tmp/gptx")
	assert.Equal(t, filepath.Join("/tmp/gptx", "gptx.db"), pr.DBFilePath())
}

func TestPathResolver_CacheDBFilePath(t *testing.T) {
	pr := NewPathResolver("/tmp/gptx")
	assert.Equal(t, filepath.Join("/tmp/gptx", "cache.db"), pr.CacheDBFilePath())
}

func TestPathResolver_HistoryFilePath(t *testing.T) {
	pr := NewPathResolver("/tmp/gptx")
	assert.Equal(t, filepath.Join("/tmp/gptx", "history.txt"), pr.HistoryFilePath())
}

func TestPathResolver_LibExecDir(t *testing.T) {
	pr := NewPathResolver("/tmp/gptx")
	assert.Equal(t, filepath.Join("/tmp/gptx", "libexec"), pr.LibExecDir())
}

func TestPathResolver_LibExecFilePath(t *testing.T) {
	pr := NewPathResolver("/tmp/gptx")
	assert.Equal(t, filepath.Join("/tmp/gptx", "libexec", "gptx-hook-shell"), pr.LibExecFilePath("gptx-hook-shell"))
}

func TestUpdatePathEnv(t *testing.T) {
	t.Run("PATH is not set", func(t *testing.T) {
		originalPath := os.Getenv("PATH")
		t.Cleanup(func() {
			os.Setenv("PATH", originalPath)
		})
		os.Setenv("PATH", "")
		pr := NewPathResolver("/tmp/gptx")
		assert.NoError(t, updatePathEnv(pr))
		assert.Equal(t, pr.LibExecDir(), os.Getenv("PATH"))
	})

	t.Run("PATH is set", func(t *testing.T) {
		originalPath := os.Getenv("PATH")
		t.Cleanup(func() {
			os.Setenv("PATH", originalPath)
		})
		os.Setenv("PATH", "/usr/local/bin")
		pr := NewPathResolver("/tmp/gptx")
		assert.NoError(t, updatePathEnv(pr))
		assert.Equal(t, pr.LibExecDir()+":/usr/local/bin", os.Getenv("PATH"))
	})

	t.Run("PATH is set and already updated", func(t *testing.T) {
		originalPath := os.Getenv("PATH")
		t.Cleanup(func() {
			os.Setenv("PATH", originalPath)
		})
		os.Setenv("PATH", "path/to/libexec:/usr/local/bin")
		pr := NewPathResolver("/tmp/gptx")
		pr.pathEnvUpdated = true
		assert.NoError(t, updatePathEnv(pr))
		assert.Equal(t, "path/to/libexec:/usr/local/bin", os.Getenv("PATH"))
	})
}
