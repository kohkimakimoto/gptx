package internal

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
)

var ConfigCommand = &cli.Command{
	Name:  "config",
	Usage: "Display configuration",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:               "pretty",
			Aliases:            []string{"p"},
			Usage:              "Pretty print",
			DisableDefaultText: true,
		},
	},
	Action: configAction,
}

var configAction = repositoryAwareAction(func(c *cli.Context, r *Repository) error {
	var buf []byte
	if c.Bool("pretty") {
		_buf, err := json.MarshalIndent(r.Config, "", "  ")
		if err != nil {
			return err
		}
		buf = _buf
	} else {
		_buf, err := json.Marshal(r.Config)
		if err != nil {
			return err
		}
		buf = _buf
	}

	_, _ = fmt.Fprintln(c.App.Writer, string(buf))
	return nil
})
