package internal

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

var InitCommand = &cli.Command{
	Name:  "init",
	Usage: "Init a gptx home directory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "directory",
			Aliases: []string{"d"},
			Usage:   "Specify a `directory` to init",
		},
	},
	Action: initAction,
}

func initAction(c *cli.Context) error {
	dir := c.String("directory")
	if dir == "" {
		dir = getAppHomeDir()
	}
	pr := NewPathResolver(dir)

	if _, err := os.Stat(pr.Dir); os.IsNotExist(err) {
		err = os.MkdirAll(pr.Dir, os.FileMode(0700))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("directory %s already exists", pr.Dir)
	}

	r := NewRepository(NewPathResolver(dir))
	if err := r.Init(); err != nil {
		return err
	}
	defer r.Close()

	_, _ = fmt.Fprintf(c.App.Writer, "Completed to init gptx home directory: %s\n", pr.Dir)
	return nil
}
