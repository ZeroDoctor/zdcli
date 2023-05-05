package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/zerodoctor/go-logging"
	"github.com/zerodoctor/zdcli/cmd"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	zdgoutil "github.com/zerodoctor/zdgo-util"
	"github.com/zerodoctor/zdtui/data"
	"github.com/zerodoctor/zdvault"
)

func StartLua(cmd string, cfg *config.Config) {
	pwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("failed to get current working directory [error=%s]", err.Error())
		return
	}

	flags := fmt.Sprintf("--os_i %s --arch_i %s --pwd %s ", cfg.OS, cfg.Arch, pwd)
	info := command.Info{
		Command: cfg.ShellCmd,
		Args:    []string{cfg.LuaCmd + " build-app.lua " + flags + cmd},
		Dir:     cfg.RootLuaScriptDir,
		Ctx:     context.Background(),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}

	logger.Debugf("[cmd=%s %s]", info.Command, info.Args)

	err = command.Exec(&info)
	if err != nil {
		logger.Errorf("failed command [error=%s]", err.Error())
	}
}

func RunUI(cfg *config.Config) {
	running := true
	for running {
		exit := StartTui(cfg)
		switch exit.Code {
		case data.EXIT_EDT:
			cmd.EditLua(exit.Msg, cfg)
			time.Sleep(100 * time.Millisecond)
			continue

		case data.EXIT_LUA:
			StartLua(exit.Msg, cfg)
			time.Sleep(100 * time.Millisecond)
			continue

		case data.EXIT_CMD:

		}

		running = false
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

func SetupLogLevel() logging.Level {
	switch os.Getenv("ZDCLI_LOG_LEVEL") {
	case "0", "DEBUG":
		return logging.DEBUG
	case "1", "INFO":
		return logging.INFO
	case "2", "NOTICE":
		return logging.NOTICE
	case "3", "WARNING":
		return logging.WARNING
	case "4", "ERROR":
		return logging.ERROR
	case "5", "CRITICAL":
		return logging.CRITICAL
	}

	return logging.DEBUG
}

func main() {
	if err := godotenv.Load(util.EXEC_PATH + "/.env"); err != nil {
		fmt.Printf("[ERROR] env file not found [error=%s]\n", err.Error())
	}

	logLevel := SetupLogLevel()
	logger.Init(logLevel)

	cfg := &config.Config{}
	if err := cfg.Load(); err != nil {
		logger.Errorf("failed to save/load config [error=%s]", err.Error())
		return
	}

	logger.Debugf("loading config...(%s)\n%s", util.EXEC_PATH, cfg)

	if cfg.VaultTokens == nil {
		cfg.VaultTokens = make(map[string]string)
	}

	for k, t := range cfg.VaultTokens {
		zdvault.SetToken(k, t)
	}

	f, err := tea.LogToFile(util.EXEC_PATH+"/debug.log", "debug")
	if err != nil {
		logger.Fatalf("failed to create bubbletea log file [error=%s]", err.Error())
	}
	defer f.Close()

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
		cmd.NewSeaweedFS(cfg),
		cmd.NewSqliteCmd(cfg),

		// meta stuff
		cmd.NewSetupCmd(cfg),
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
