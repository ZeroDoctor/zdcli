package main

import (
	"context"
	"errors"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/alert"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/util"
	zdgoutil "github.com/zerodoctor/zdgo-util"
)

func RunUI(cfg *config.Config) {
	running := true
	for running {
		exit := StartTui(cfg)

		switch exit.Code {
		case comp.EXIT_EDT:
			StartEdit(exit.Msg, cfg)
			time.Sleep(100 * time.Millisecond)
			continue

		case comp.EXIT_LUA:
			StartLua(exit.Msg, cfg)
			time.Sleep(100 * time.Millisecond)
			continue

		case comp.EXIT_CMD:

		}

		running = false
	}
}

func PasteCmd() *cli.Command {
	return &cli.Command{
		Name:  "paste",
		Usage: "common commands to interact with pastebin.com. May need to login via this cli before use.",
		Subcommands: []*cli.Command{
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "upload files to pastebin.com",
				Action: func(ctx *cli.Context) error {
					PasteBinUpload(ctx.Args().Slice())

					return nil
				},
			},
			{
				Name:  "update",
				Usage: "update files in pastebin.com while keep the same pastebin key",
				Action: func(ctx *cli.Context) error {
					PasteBinUpdate(ctx.Args().Slice())

					return nil
				},
			},
		},
	}
}

func NewCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "create a new lua script",
		Action: func(ctx *cli.Context) error {
			for _, arg := range ctx.Args().Slice() {
				CreateLua(arg, cfg)
			}

			return nil
		},
	}
}

func RemoveCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "remove a lua script or a directory",
		Action: func(ctx *cli.Context) error {
			for _, arg := range ctx.Args().Slice() {
				RemoveLua(arg, cfg)
			}

			return nil
		},
	}
}

func EditCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "edit",
		Aliases: []string{"e"},
		Usage:   "edits a lua script",
		Action: func(ctx *cli.Context) error {
			StartEdit(ctx.Args().Get(0), cfg)
			return nil
		},
	}
}

func ListCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "list current lua scripts",
		Action: func(ctx *cli.Context) error {
			StartLs(cfg)
			return nil
		},
	}
}

func SetupCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "setup lua, editor, and dir configs",
	}
}

func UICmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "ui",
		Usage: "opens a custom terminal emulator",
		Action: func(ctx *cli.Context) error {
			RunUI(cfg)
			return nil
		},
	}
}

func AlertCmd() *cli.Command {
	return &cli.Command{
		Name:  "alert",
		Usage: "notifies user when an event happens",
		Subcommands: []*cli.Command{
			{
				Name:    "endpoint",
				Aliases: []string{"e"},
				Usage:   "create an alert when endpoint fails or returns status code not between 200-299",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "route",
						Usage:    "an endpoint i.e. https://google.com",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "message",
						Usage:    "a message to display when route/endpoint fails",
						Required: true,
					},
					&cli.IntFlag{
						Name:    "check-duration",
						Aliases: []string{"c"},
						Usage:   "period (in seconds) alert checks endpoint",
					},
				},

				Action: func(ctx *cli.Context) error {
					c, cancel := context.WithCancel(ctx.Context)
					defer cancel()

					checkDur := 5 * time.Second
					if sec := ctx.Int("check-duration"); sec > 0 {
						checkDur = time.Duration(sec) * time.Second
					}

					a := alert.WatchEndpoint(
						c,
						ctx.String("route"),
						ctx.String("message"),
						checkDur,
					)
					a.Wait()

					return nil
				},
			},
		},

		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelp(ctx)
			return errors.New("must provide additional subcommand(s)")
		},
	}
}

func main() {
	logger.Init()

	if err := godotenv.Load(util.EXEC_PATH + "/.env"); err != nil {
		logger.Info("env file not found [error=%s]", err.Error())
	}

	cfg := &config.Config{}
	if err := cfg.Load(); err != nil {
		logger.Errorf("failed to save/load config [error=%s]", err.Error())
	}

	app := cli.NewApp()
	app.Commands = []*cli.Command{
		NewCmd(cfg),
		RemoveCmd(cfg),
		EditCmd(cfg),
		ListCmd(cfg),
		SetupCmd(cfg),
		UICmd(cfg),
		AlertCmd(),
		PasteCmd(),
	}

	app.Action = func(ctx *cli.Context) error {
		if ctx.Args().Len() <= 0 {
			RunUI(cfg)
			return nil
		}

		StartLua(strings.Join(ctx.Args().Slice(), " "), cfg)

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go zdgoutil.OnExit(func(s os.Signal, i ...interface{}) {
		cancel()
	})

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		logger.Fatalf("failed to run cli [error=%s]", err.Error())
	}

	logger.Print("Good Bye.")
}
