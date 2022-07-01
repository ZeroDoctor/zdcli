package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/zdcli/cmd"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/util"
	zdgoutil "github.com/zerodoctor/zdgo-util"
	"github.com/zerodoctor/zdvault"
)

func StartLua(cmd string, cfg *config.Config) {
	info := command.Info{
		Command: cfg.ShellCmd, // TODO: allow user to set lua endpoint
		Args:    []string{cfg.LuaCmd + " build-app.lua " + cmd},
		Dir:     cfg.RootScriptDir, // TODO: allow user to set lua direcoty
		Ctx:     context.Background(),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}

	err := command.Exec(&info)
	if err != nil {
		logger.Errorf("failed command [error=%s]", err.Error())
	}
}

func RunUI(cfg *config.Config) {
	running := true
	for running {
		exit := StartTui(cfg)
		switch exit.Code {
		case comp.EXIT_EDT:
			cmd.EditLua(exit.Msg, cfg)
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
		return
	}

	if cfg.VaultTokens == nil {
		cfg.VaultTokens = make(map[string]string)
	}

	for k, t := range cfg.VaultTokens {
		zdvault.SetToken(k, t)
	}

	app := cli.NewApp()
	app.Commands = []*cli.Command{
		// main lua commands
		cmd.NewLuaCmd(cfg),
		cmd.RemoveLuaCmd(cfg),
		cmd.EditLuaCmd(cfg),
		cmd.ListLuaCmd(cfg),

		// one or more sub commands
		cmd.NewAlertCmd(),
		cmd.NewPasteCmd(),
		cmd.NewVaultCmd(cfg),

		// meta stuff
		SetupCmd(cfg),
		UICmd(cfg),
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

	go func() {
		err := app.RunContext(ctx, os.Args)
		if err != nil {
			logger.Fatalf("failed to run cli [error=%s]", err.Error())
		}
		cancel()
	}()

	zdgoutil.OnExitWithContext(ctx, func(s os.Signal, i ...interface{}) {
		cancel()
	})

	fmt.Println("Good Bye.")
}
