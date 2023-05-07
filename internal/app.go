package internal

import (
	"github.com/urfave/cli/v2"
)

func Run(args []string) error {
	return NewApp(NewRepository(NewPathResolver(getAppHomeDir()))).Run(args)
}

func NewApp(r *Repository) *cli.App {
	app := cli.NewApp()
	app.Name = "gptx"
	app.Version = Version
	app.Usage = "An extensible command-line utility powered by ChatGPT"
	app.Copyright = "Copyright (c) 2023 Kohki Makimoto"
	app.Commands = []*cli.Command{
		ChatCommand,
		CleanCommand,
		ConfigCommand,
		DeleteCommand,
		InitCommand,
		InspectCommand,
		ListCommand,
		RenameCommand,
		VersionCommand,
	}

	// setup logic to resolve custom subcommands
	defaultHelpAction := app.Action
	app.Action = func(c *cli.Context) error {
		args := c.Args()
		if args.Present() {
			return trySubcommand(c)
		}
		return defaultHelpAction(c)
	}

	// setup repository
	app.Metadata = map[string]interface{}{
		"repository": r,
	}

	return app
}

func getRepository(app *cli.App) (*Repository, error) {
	r := app.Metadata["repository"].(*Repository)
	if err := r.Init(); err != nil {
		return nil, err
	}
	return r, nil
}
