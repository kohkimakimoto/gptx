package internal

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/urfave/cli/v2"
	"io"
	"strings"
)

var ChatCommand = &cli.Command{
	Name:      "chat",
	Aliases:   []string{"c"},
	Usage:     "Chat with GPT",
	ArgsUsage: "[prompt...]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:               "no-loading",
			Usage:              "Disable loading animation",
			DisableDefaultText: true,
		},
		&cli.BoolFlag{
			Name:               "no-animation",
			Usage:              "Disable typing animation",
			DisableDefaultText: true,
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Specify a `name` for the conversation",
		},
		&cli.StringFlag{
			Name:    "label",
			Aliases: []string{"l"},
			Usage:   "Specify a `label` for the conversation",
		},
		&cli.StringFlag{
			Name:    "resume",
			Aliases: []string{"r"},
			Usage:   "Resume a `conversation`. You can specify a conversation id or name",
		},
		&cli.BoolFlag{
			Name:               "editor",
			Aliases:            []string{"e"},
			Usage:              "Open $EDITOR to make a prompt.",
			DisableDefaultText: true,
		},
		&cli.StringFlag{
			Name:  "model",
			Usage: "Specify a `model` working with ChatGPT",
		},
		&cli.Float64Flag{
			Name:  "temperature",
			Usage: "Specify a temperature",
			Value: 1,
		},
		&cli.Float64Flag{
			Name:  "top-p",
			Usage: "Specify a top_p",
			Value: 1,
		},
		&cli.StringSliceFlag{
			Name:    "hook",
			Aliases: []string{"H"},
			Usage:   "Specify a `hook`s",
		},
		&cli.StringSliceFlag{
			Name:    "env",
			Aliases: []string{"E"},
			Usage:   "Specify a environment variable for hooks. For example: -E FOO=BAR -E BAZ=QUX",
		},
		&cli.BoolFlag{
			Name:               "interactive",
			Aliases:            []string{"i"},
			Usage:              "Run in interactive mode",
			DisableDefaultText: true,
		},
		&cli.BoolFlag{
			Name:               "no-cache",
			Usage:              "Disable cache",
			DisableDefaultText: true,
		},
		&cli.BoolFlag{
			Name:               "on-memory",
			Aliases:            []string{"m"},
			Usage:              "Disable the ability to save or load the conversation to and from the disk, which also includes the behavior of the --no-cache option.",
			DisableDefaultText: true,
		},
	},
	Action: chatAction,
}

var chatAction = repositoryAwareAction(func(c *cli.Context, r *Repository) error {
	prompt := strings.Join(c.Args().Slice(), " ")
	noLoading := c.Bool("no-loading")
	noAnimation := c.Bool("no-animation")
	name := c.String("name")
	label := c.String("label")
	resume := c.String("resume")
	editor := c.Bool("editor")
	model := c.String("model")
	temperature := float32(c.Float64("temperature"))
	topP := float32(c.Float64("top-p"))
	hookNames := c.StringSlice("hook")
	interactive := c.Bool("interactive")
	noCache := c.Bool("no-cache")
	onMemory := c.Bool("on-memory")
	hooksEnv := c.StringSlice("env")

	if interactive && isPipe(c.App.Reader) {
		return fmt.Errorf("interactive mode is not supported with pipe")
	}

	if isPipe(c.App.Reader) {
		b, err := io.ReadAll(c.App.Reader)
		if err != nil {
			return err
		}
		stdinPrompt := string(b)
		if stdinPrompt != "" {
			prompt = stdinPrompt + "\n" + prompt
		}
	}

	// validate prompt
	if prompt == "" && !editor && !interactive {
		return fmt.Errorf("prompt is required")
	}

	if editor {
		_prompt, err := getPromptFromEditor(prompt)
		if err != nil {
			return err
		}
		prompt = _prompt
	}

	if interactive && len(prompt) > 0 {
		return fmt.Errorf("prompt is not supported in interactive mode")
	}

	// initialize chat service
	sv, err := r.NewChatService(c.App.Writer)
	if err != nil {
		return err
	}
	if noAnimation {
		sv.DisableOutputAnimation()
	}
	sv.NoLoading = noLoading

	if model == "" {
		model = r.Config.Model
	}
	sv.Model = model
	sv.Temperature = temperature
	sv.TopP = topP
	sv.NoCache = noCache
	sv.OnMemory = onMemory
	sv.HooksEnv = hooksEnv

	if err := sv.InitConversation(resume, name, label); err != nil {
		return err
	}

	if err := sv.LoadHooks(hookNames); err != nil {
		return err
	}

	if interactive {
		// REPL mode
		return doREPL(c, r, sv)
	} else {
		if err := sv.Chat(prompt); err != nil {
			if !isErrCancel(err) {
				return err
			}
		}
	}

	return nil
})

func doREPL(c *cli.Context, r *Repository, sv *ChatService) error {
	l, err := readline.NewEx(&readline.Config{
		Prompt:       "> ",
		AutoComplete: completer,
		HistoryFile:  r.PathResolver.HistoryFilePath(),
	})
	if err != nil {
		return err
	}
	defer l.Close()

	fmt.Fprintf(c.App.Writer, "Welcome to gptx v%s. \nType \".help\" for more information.\n", Version)
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, ".") {
			switch line {
			case ".help":
				_, _ = fmt.Fprintln(c.App.Writer, strings.TrimSpace(`
Available commands:
  .editor   Enter editor mode to input multi-line text
  .exit:    Exit REPL
  .quit:    Alias for .exit
  .help:    Show this help

Press Ctrl+C to abort current expression, Ctrl+D to exit the REPL
`))
			case ".exit", ".quit":
				goto exit
			case ".editor":
				_, _ = fmt.Fprintln(c.App.Writer, "Entering editor mode (Ctrl+D to finish, Ctrl+C to cancel)")
				line = doREPLEditor(l)
				l.SetPrompt("> ")
				goto chat
			default:
				_, _ = fmt.Fprintf(c.App.Writer, "Unknown command: %s\n", line)
			}
			continue
		}

	chat:
		if line != "" {
			if err := sv.Chat(line); err != nil {
				if !isErrCancel(err) {
					_, _ = fmt.Fprintf(c.App.ErrWriter, "%s\n", err.Error())
				}
			}
		}
	}
exit:
	return nil
}

func doREPLEditor(l *readline.Instance) string {
	l.SetPrompt(">> ")
	var ml string
	for {
		if line, err := l.Readline(); err == nil {
			if ml == "" {
				ml = line
			} else {
				ml = ml + "\n" + line
			}
		} else {
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					ml = ""
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
		}
	}
	return ml
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem(".editor"),
	readline.PcItem(".exit"),
	readline.PcItem(".quit"),
	readline.PcItem(".help"),
)
