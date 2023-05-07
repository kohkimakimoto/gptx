package internal

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

var trySubcommand = repositoryAwareAction(func(c *cli.Context, r *Repository) error {
	if err := updatePathEnv(r.PathResolver); err != nil {
		return err
	}

	args := c.Args()
	subCmdName := args.First()
	subCmd := fmt.Sprintf("gptx-%s", subCmdName)
	if strings.HasPrefix(subCmdName, "gptx-hook-") {
		return fmt.Errorf("the command name '%s' is reserved for hooks", subCmdName)
	}
	subCmdPath, err := exec.LookPath(subCmd)
	if err != nil {
		return fmt.Errorf("the command '%s' was not found: %w", subCmdName, err)
	}

	// run custom subcommand
	cmd := exec.Command(subCmdPath, args.Tail()...)
	cmd.Stdin = c.App.Reader
	cmd.Stdout = c.App.Writer
	cmd.Stderr = c.App.ErrWriter
	return cmd.Run()
})
