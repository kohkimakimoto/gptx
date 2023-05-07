package internal

import (
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"strings"
	"time"
)

var ListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List conversations",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "begin",
			Aliases: []string{"b"},
			Usage:   "Load the conversations starting from the given conversation `id`",
		},
		&cli.BoolFlag{
			Name:               "reverse",
			Aliases:            []string{"r"},
			Usage:              "Sort the conversations in descending id order",
			DisableDefaultText: true,
		},
		&cli.IntFlag{
			Name:        "limit",
			Aliases:     []string{"L"},
			Usage:       "Limit the `number` of displayed conversations",
			DefaultText: "0 (no limit)",
		},
		&cli.StringFlag{
			Name:    "label",
			Aliases: []string{"l"},
			Usage:   "Filter the conversations by `label`",
		},
		&cli.BoolFlag{
			Name:               "quiet",
			Aliases:            []string{"q"},
			Usage:              "Only display the conversation ids",
			DisableDefaultText: true,
		},
	},
	Action: listAction,
}

var listAction = repositoryAwareAction(func(c *cli.Context, r *Repository) (err error) {
	query := &ListConversationsQuery{
		Reverse: c.Bool("reverse"),
		Limit:   c.Int("limit"),
		Label:   c.String("label"),
	}

	quiet := c.Bool("quiet")

	begin, err := getUint64ValueFromStringFlag(c, "begin")
	if err != nil {
		return err
	}
	if begin != 0 {
		query.Begin = &begin
	}

	store, err := r.StoreManager.Open()
	if err != nil {
		return err
	}
	defer store.Close()

	list, err := store.ListConversations(query)
	if err != nil {
		return err
	}

	t := NewSimpleTableWriter(c.App.Writer)
	if !quiet {
		t.AppendHeader(table.Row{
			"ID",
			"PROMPT",
			"MESSAGES",
			"NAME",
			"LABEL",
			"HOOKS",
			"CREATED",
			"ELAPSED",
		})
	}
	for _, c := range list.Conversations {
		if !quiet {
			t.AppendRow([]interface{}{
				c.Id,
				truncateChars(strings.ReplaceAll(c.Prompt, "\n", " "), 50),
				len(c.Messages),
				c.Name,
				c.Label,
				strings.Join(c.Hooks, ", "),
				c.CreatedAt.Format(time.RFC3339),
				humanize.Time(c.CreatedAt),
			})
		} else {
			t.AppendRow([]interface{}{c.Id})
		}
	}
	t.Render()
	return nil
})
