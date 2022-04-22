package cmdstate

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/tui/inter"
	"github.com/zerodoctor/zdcli/tui/ui"
	"github.com/zerodoctor/zdcli/util"
)

var (
	ErrCommandNotRunning error = errors.New("command not running")
	ErrUnknownCommand    error = errors.New("unknown command")
)

func NewData(t, m string) comp.Data {
	return comp.Data{Type: t, Msg: m}
}

type State struct {
	vm inter.IViewManager

	state *comp.Stack
}

func NewState(vm inter.IViewManager, state *comp.Stack) *State {
	return &State{vm: vm, state: state}
}

func (s *State) Exec(cmd string) error {
	split := strings.Split(cmd, " ")
	cmd = strings.Join(split[1:], " ")

	switch split[0] {
	case "exec":
		if len(split) > 2 && (split[1] == "--tty" || split[1] == "-t") {
			cmd = strings.Join(split[2:], " ")

			s.vm.SetExitMsg(comp.ExitMessage{
				Code: comp.EXIT_CMD,
				Msg:  cmd,
			})

			s.vm.G().UpdateAsync(func(g *gocui.Gui) error {
				return s.vm.Quit(s.vm.G(), nil)
			})

			return nil
		}

		s.state.Push(NewForkState(s.vm, s.state, cmd))
		return nil

	case "lua":
		if len(split) > 2 && (split[1] == "--tty" || split[1] == "-t") {
			cmd = strings.Join(split[2:], " ")

			s.vm.SetExitMsg(comp.ExitMessage{
				Code: comp.EXIT_LUA,
				Msg:  cmd,
			})

			s.vm.G().UpdateAsync(func(g *gocui.Gui) error {
				return s.vm.Quit(s.vm.G(), nil)
			})

			return nil
		}

		s.state.Push(NewLuaState(s.vm, s.state, cmd))
		return nil

	case "edit":
		s.vm.SetExitMsg(comp.ExitMessage{
			Code: comp.EXIT_EDT,
			Msg:  cmd,
		})

		s.vm.G().UpdateAsync(func(*gocui.Gui) error {
			return s.vm.Quit(s.vm.G(), nil)
		})

		return nil

	case "go":
		return nil

	case "ls":
		path := util.EXEC_PATH + "/lua/scripts"
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
