package builtin

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed gptx-hook-shell
var bShellHook []byte

var Files = map[string][]byte{
	"gptx-hook-shell": bShellHook,
}

func InitLibexecFiles(pathToLibexec string) error {
	for name, b := range Files {
		err := createExecutableFileIfNotExist(filepath.Join(pathToLibexec, name), b)
		if err != nil {
			return err
		}
	}
	return nil
}

// createExecutableFileIfNotExist creates an executable file if it does not exist.
func createExecutableFileIfNotExist(path string, b []byte) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	// Write the content to the file
	err = os.WriteFile(path, b, 0755)
	if err != nil {
		return err
	}

	return nil
}
