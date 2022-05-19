package main

import (
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/util"
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

func main() {
	logger.Init()

	cfg := &config.Config{}
	if err := cfg.Load(); err != nil {
		cfg.ShellCmd = "bash -c -i"
		cfg.EditorCmd = "nvim"
		cfg.LuaCmd = "lua"
		cfg.RootScriptDir = util.EXEC_PATH + "/lua"
	}

	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "create a new lua script",
			Action: func(ctx *cli.Context) error {
				for _, arg := range ctx.Args().Slice() {
					CreateLua(arg, cfg)
				}

				return nil
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "remove a lua script or a directory",
			Action: func(ctx *cli.Context) error {
				for _, arg := range ctx.Args().Slice() {
					RemoveLua(arg, cfg)
				}

				return nil
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "edits a lua script",
			Action: func(ctx *cli.Context) error {
				StartEdit(ctx.Args().Get(0), cfg)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list current lua scripts",
			Action: func(ctx *cli.Context) error {
				StartLs(cfg)
				return nil
			},
		},
		{
			Name:  "setup",
			Usage: "setup lua, editor, and dir configs",
		},
		{
			Name:  "ui",
			Usage: "opens a custom terminal emulator",
			Action: func(ctx *cli.Context) error {
				RunUI(cfg)
				return nil
			},
		},
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

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalf("failed to run cli [error=%s]", err.Error())
	}

	logger.Print("Good Bye.")
}
