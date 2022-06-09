package main

import (
	"context"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/cmd"
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
			cmd.StartEdit(exit.Msg, cfg)
			time.Sleep(100 * time.Millisecond)
			continue

		case comp.EXIT_LUA:
			cmd.StartLua(exit.Msg, cfg)
			time.Sleep(100 * time.Millisecond)
			continue

		case comp.EXIT_CMD:

		}

		running = false
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
		cmd.NewCmd(cfg),
		cmd.RemoveCmd(cfg),
		cmd.EditCmd(cfg),
		cmd.ListCmd(cfg),
		SetupCmd(cfg),
		UICmd(cfg),
		cmd.AlertCmd(),
		cmd.PasteCmd(),
	}

	app.Action = func(ctx *cli.Context) error {
		if ctx.Args().Len() <= 0 {
			RunUI(cfg)
			return nil
		}

		cmd.StartLua(strings.Join(ctx.Args().Slice(), " "), cfg)

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
