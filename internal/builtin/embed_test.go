package builtin

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestInitLibexecFiles(t *testing.T) {
	dir := testTempDir(t)
	err := InitLibexecFiles(dir)
	assert.NoError(t, err)

	// check if the file exists
	assert.FileExists(t, filepath.Join(dir, "gptx-hook-shell"))
}

func testTempDir(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "gptx_test_")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Remove(tempDir)
	})
	return tempDir
}
