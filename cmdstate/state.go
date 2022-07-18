package cmdstate

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/tui/data"
	"github.com/zerodoctor/zdcli/tui/ui"
)

var (
	ErrCommandNotRunning error = errors.New("command not running")
	ErrUnknownCommand    error = errors.New("unknown command")
)

func NewData(t, m string) data.Data {
	return data.Data{Type: t, Msg: m}
}

type State struct {
	vm  data.IViewManager
	cfg *config.Config

	state *data.Stack
}

func NewState(vm data.IViewManager, state *data.Stack, cfg *config.Config) *State {
	return &State{vm: vm, state: state, cfg: cfg}
}

func (s *State) Exec(cmd string) error {
	split := strings.Split(cmd, " ")
	cmd = strings.Join(split[1:], " ")

	switch split[0] {
	case "exec":
		if len(split) > 2 && (split[1] == "--tty" || split[1] == "-t") {
			cmd = strings.Join(split[2:], " ")

			s.vm.SetExitMsg(data.ExitMessage{
				Code: data.EXIT_CMD,
				Msg:  cmd,
			})

			s.vm.G().UpdateAsync(func(g *gocui.Gui) error {
				return s.vm.Quit(s.vm.G(), nil)
			})

			return nil
		}

		s.state.Push(NewForkState(s.vm, s.state, cmd, s.cfg))
		return nil

	case "lua":
		if len(split) > 2 && (split[1] == "--tty" || split[1] == "-t") {
			cmd = strings.Join(split[2:], " ")

			s.vm.SetExitMsg(data.ExitMessage{
				Code: data.EXIT_LUA,
				Msg:  cmd,
			})

			s.vm.G().UpdateAsync(func(g *gocui.Gui) error {
				return s.vm.Quit(s.vm.G(), nil)
			})

			return nil
		}

		s.state.Push(NewLuaState(s.vm, s.state, cmd, s.cfg))
		return nil

	case "edit":
		s.vm.SetExitMsg(data.ExitMessage{
			Code: data.EXIT_EDT,
			Msg:  cmd,
		})

		s.vm.G().UpdateAsync(func(*gocui.Gui) error {
			return s.vm.Quit(s.vm.G(), nil)
		})

		return nil

	case "go":
		return nil

	case "ls":
		path := s.cfg.RootScriptDir + "/scripts"
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		var data [][]interface{}

		for _, file := range files {
			data = append(data, []interface{}{
				file.Mode(), file.Name(), file.Size(), file.ModTime(),
			})
		}

		screen, err := s.vm.GetView("screen")
		if err != nil {
			return err
		}

		table, err := ui.NewTable(
			[]string{"Mode", "Name", "Size", "Modify Time"},
			data,
			screen.Width(), 10*len(data),
		)
		if err != nil {
			return err
		}

		s.vm.SendView("screen", NewData("msg", table.View()+"\n"))

		return nil
	}

	return ErrUnknownCommand
}

func (s *State) Stop() error { return ErrCommandNotRunning }

func (s *State) Stack() *data.Stack { return s.state }
