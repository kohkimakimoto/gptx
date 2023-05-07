package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
)

var InspectCommand = &cli.Command{
	Name:      "inspect",
	Aliases:   []string{"i"},
	Usage:     "Display details information on one or more conversations",
	ArgsUsage: `[conversation...]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:               "pretty",
			Aliases:            []string{"p"},
			Usage:              "Pretty print",
			DisableDefaultText: true,
		},
	},
	Action: inspectAction,
}

var inspectAction = repositoryAwareAction(func(c *cli.Context, r *Repository) (err error) {
	if c.NArg() == 0 {
		return errors.New("missing conversation argument(s)")
	}

	store, err := r.StoreManager.Open()
	if err != nil {
		return err
	}
	defer store.Close()

	list := make([]*Conversation, 0, c.NArg())
	for _, key := range c.Args().Slice() {
		co, err := store.GetConversationByKey(NewConversationKey(key))
		if err != nil {
			return err
		}
		list = append(list, co)
	}

	pretty := c.Bool("pretty")
	var buf []byte
	if len(list) == 1 {
		if pretty {
			_buf, err := json.MarshalIndent(list[0], "", "  ")
			if err != nil {
				return err
			}
			buf = _buf
		} else {
			_buf, err := json.Marshal(list[0])
			if err != nil {
				return err
			}
			buf = _buf
		}
	} else {
		if pretty {
			_buf, err := json.MarshalIndent(list, "", "  ")
			if err != nil {
				return err
			}
			buf = _buf
		} else {
			_buf, err := json.Marshal(list)
			if err != nil {
				return err
			}
			buf = _buf
		}
	}
	_, _ = fmt.Fprintln(c.App.Writer, string(buf))
	return nil
})
