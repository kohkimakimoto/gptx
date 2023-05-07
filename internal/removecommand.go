package internal

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strconv"
)

var DeleteCommand = &cli.Command{
	Name:      "remove",
	Aliases:   []string{"rm"},
	Usage:     "Remove one or more conversations",
	ArgsUsage: `[conversation_id...]`,
	Action:    deleteAction,
}

var deleteAction = repositoryAwareAction(func(c *cli.Context, r *Repository) (err error) {
	if c.NArg() == 0 {
		return fmt.Errorf("missing conversation id")
	}

	store, err := r.StoreManager.Open()
	if err != nil {
		return err
	}
	defer store.Close()

	for _, id := range c.Args().Slice() {
		uint64v, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid conversation id: %s", id)
		}
		if err := store.DeleteConversationById(uint64v); err != nil {
			return err
		}
	}
	return nil
})
