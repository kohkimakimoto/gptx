package internal

import (
	"github.com/urfave/cli/v2"
)

var VersionCommand = &cli.Command{
	Name:   "version",
	Usage:  "Print the version information",
	Action: versionAction,
}

var versionAction = repositoryAwareAction(func(c *cli.Context, r *Repository) (err error) {
	t := NewSimpleTableWriter(c.App.Writer)
	t.AppendRow([]interface{}{"Version:", Version})
	t.AppendRow([]interface{}{"Commit Hash:", CommitHash})
	t.Render()
	return nil
})
