package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	HookTypePreMessage  = "pre-message"
	HookTypePostMessage = "post-message"
	HookTypeFinish      = "finish"
)

type HookFactory struct{}

func (f *HookFactory) NewHook(name string) (*Hook, error) {
	h := &Hook{}
	h.Name = name

	commandPath, err := exec.LookPath(f.resolveHookCommandPath(name))
	if err != nil {
		return nil, err
	}
	h.CommandPath = commandPath

	return h, nil
}

func (f *HookFactory) resolveHookCommandPath(name string) string {
	if name == "" {
		return ""
	}

	if strings.Contains(name, string(os.PathSeparator)) {
		// path is specified
		// NOTE: Specifying path is not recommended. It is mainly for testing.
		return name
	} else {
		// just a name is specified
		if !strings.HasPrefix(name, "gptx-hook-") {
			name = fmt.Sprintf("gptx-hook-%s", name)
		}
	}

	return name
}

type Hook struct {
	Name        string
	CommandPath string
}

func (h *Hook) Command() *exec.Cmd {
	cmd := exec.Command(h.CommandPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
