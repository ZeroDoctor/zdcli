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
	"github.com/zerodoctor/zdcli/cmd/vault"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/util"
	zdgoutil "github.com/zerodoctor/zdgo-util"
	"github.com/zerodoctor/zdtui/data"
)

// RunScript routes `zd <name> <args...>` to lua, python, or bash based on where
// the script named by the first token lives. Lua is tried first and, when no
// script matches, falls back to lua so its "script not found" error is unchanged.
func RunScript(cmd string, cfg *config.Config) error {
	name := strings.Fields(cmd)
	if len(name) == 0 {
		return StartLua(cmd, cfg)
	}
	rel := strings.ReplaceAll(name[0], ".", "/")

	switch {
	case zdgoutil.FileExists(cfg.RootLuaScriptDir + "/" + rel + ".lua"):
		return StartLua(cmd, cfg)
	case zdgoutil.FileExists(cfg.RootPythonScriptDir + "/" + rel + ".py"):
		return StartPython(cmd, cfg)
	case zdgoutil.FileExists(cfg.RootBashScriptDir + "/" + rel + ".sh"):
		return StartBash(cmd, cfg)
	}

	return StartLua(cmd, cfg)
}

// StartPython runs `python main.py -s <name> -f <func> -a <json>` in the python
// script dir. Usage: `zd <name> [func] [json-kwargs]` (func defaults to main).
func StartPython(cmd string, cfg *config.Config) error {
	fields := strings.Fields(cmd)
	name := fields[0]
	fn := "main"
	if len(fields) > 1 {
		fn = fields[1]
	}
	args := "{}"
	if len(fields) > 2 {
		args = strings.Join(fields[2:], " ")
	}

	info := command.Info{
		Command: cfg.ShellCmd,
		Args:    []string{fmt.Sprintf("%s main.py -s %s -f %s -a '%s'", cfg.PythonCmd, name, fn, args)},
		Dir:     cfg.RootPythonScriptDir,
		Ctx:     context.Background(),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}

	logger.Debugf("[cmd=%s %s]", info.Command, info.Args)

	if err := command.Exec(&info); err != nil {
		return fmt.Errorf("failed command [error=%s]", err.Error())
	}

	return nil
}

// StartBash runs `bash <name>.sh <args...>` from the bash script dir.
func StartBash(cmd string, cfg *config.Config) error {
	fields := strings.Fields(cmd)
	rel := strings.ReplaceAll(fields[0], ".", "/")
	rest := strings.Join(fields[1:], " ")

	info := command.Info{
		Command: cfg.ShellCmd,
		Args:    []string{strings.TrimSpace(fmt.Sprintf("bash %s.sh %s", rel, rest))},
		Dir:     cfg.RootBashScriptDir,
		Ctx:     context.Background(),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}

	logger.Debugf("[cmd=%s %s]", info.Command, info.Args)

	if err := command.Exec(&info); err != nil {
		return fmt.Errorf("failed command [error=%s]", err.Error())
	}

	return nil
}

func StartLua(cmd string, cfg *config.Config) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory [error=%s]", err.Error())
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
		return fmt.Errorf("failed command [error=%s]", err.Error())
	}

	return nil
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
		os.Exit(1)
	}

	logLevel := SetupLogLevel()
	logger.Init(logLevel)

	cfg := &config.Config{}
	if err := cfg.Load(); err != nil {
		logger.Fatalf("failed to save/load config [error=%s]", err.Error())
	}

	logger.Debugf("loading config...(%s)\n%s", util.EXEC_PATH, cfg)

	if cfg.VaultTokens == nil {
		cfg.VaultTokens = make(map[string]string)
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
		vault.NewVaultCmd(cfg),
		cmd.NewSeaweedFS(cfg),
		cmd.NewSqliteCmd(cfg),
		cmd.NewFileServerCmd(),

		// meta stuff
		cmd.NewSetupCmd(cfg),
		{
			Name:  "version",
			Usage: "current version of cli",
			Action: func(ctx *cli.Context) error {
				fmt.Println(VERSION)
				return nil
			},
		},
		UICmd(cfg),
	}

	app.Action = func(ctx *cli.Context) error {
		if ctx.Args().Len() <= 0 {
			RunUI(cfg)
			return nil
		}

		return RunScript(strings.Join(ctx.Args().Slice(), " "), cfg)
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
}
