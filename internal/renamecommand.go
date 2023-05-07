package internal

import (
	"errors"
	"github.com/urfave/cli/v2"
)

var RenameCommand = &cli.Command{
	Name:      "rename",
	Usage:     "Rename a conversation",
	ArgsUsage: `[conversation] [new name]`,
	Action:    renameAction,
}

var renameAction = repositoryAwareAction(func(c *cli.Context, r *Repository) error {
	if c.NArg() != 2 {
		return errors.New("missing argument(s)")
	}

	src := c.Args().Get(0)
	newName := c.Args().Get(1)

	store, err := r.StoreManager.Open()
	if err != nil {
		return err
	}
	defer store.Close()

	return store.RenameConversationByKey(NewConversationKey(src), newName)
})
