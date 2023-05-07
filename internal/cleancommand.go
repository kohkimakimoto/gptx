package internal

import (
	"github.com/urfave/cli/v2"
)

var CleanCommand = &cli.Command{
	Name:   "clean",
	Usage:  "Clean up the cache",
	Action: cleanAction,
}

var cleanAction = repositoryAwareAction(func(c *cli.Context, r *Repository) error {
	return r.CacheManager.Refresh()
})
