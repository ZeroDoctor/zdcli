package main

import (
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/zerodoctor/zdcli/logger"
	"github.com/zerodoctor/zdcli/tui/comp"
)

// TODO: Take a look at gocmd and urfave/cli

type EditCmd struct {
	Name string `arg:"" optional:"" help:"name of lua script to edit"`
}

func (e *EditCmd) Run() error {
	StartEdit(e.Name)
	return nil
}

type ListCmd struct{}

func (l *ListCmd) Run() error {
	StartLs()
	return nil
}

type UiCmd struct{}

func (u *UiCmd) Run() error {

	running := true
	for running {
		exit := StartTui()

		switch exit.Code {
		case comp.EXIT_EDT:
			StartEdit(exit.Msg)
			time.Sleep(100 * time.Millisecond)
			continue

		case comp.EXIT_CMD:

		case comp.EXIT_LUA:
			StartLua(exit.Msg)
			time.Sleep(100 * time.Millisecond)
			continue

		}

		running = false
	}

	return nil
}

type RunCmd struct {
	Script []string `arg:"" help:"script name and its arguments"`
}

func (r *RunCmd) Run() error {
	StartLua(strings.Join(r.Script, " "))
	return nil
}

var cli struct {
	Edit EditCmd `cmd:"" help:"edit a lua script"`
	List ListCmd `cmd:"" help:"list current lua script"`
	Ui   UiCmd   `cmd:"" help:"show a terminal user interface"`
	Run  RunCmd  `cmd:"" help:"run a lua script"`
}

func main() {
	logger.Init()
	ctx := kong.Parse(&cli)

	err := ctx.Run()
	if err == nil {
		return
	}

	if len(ctx.Command()) > 0 {
		StartEdit(ctx.Command())
	}

	logger.Print("Good Bye.")
}
