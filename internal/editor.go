package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func getPromptFromEditor(defaultPrompt string) (string, error) {
	// Create a temporary file.
	file, err := os.CreateTemp("", "gptx-prompt-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())

	// Write the default prompt to the file.
	if _, err := file.WriteString(defaultPrompt); err != nil {
		return "", err
	}

	// Open the file in the editor.
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Read the text from the file.
	data, err := os.ReadFile(file.Name())
	if err != nil {
		return "", err
	}
	output := string(data)

	if output == "" {
		return "", fmt.Errorf("couldn't get valid PROMPT from $EDITOR")
	}
	return output, nil
}
