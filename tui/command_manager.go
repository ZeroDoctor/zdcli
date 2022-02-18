package tui

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/command"
)

var ErrUnknownCommand error = errors.New("unknown command")

type ICommandState interface {
	Exec(cmd string) error
}

type CommandManager struct {
	vm    *ViewManager
	state Stack
}

func NewCommandManager(vm *ViewManager) *CommandManager {
	cm := &CommandManager{
		vm:    vm,
		state: NewStack(),
	}

	cm.state.Push(NewState(vm, &cm.state))

	return cm
}

func (cm *CommandManager) Cmd(cmd string) {
	if cm.state.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.state.Peek().(ICommandState).Exec(cmd)
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s | %s]\n", err.Error(), cmd)))
	}
}

type State struct {
	vm *ViewManager

	state *Stack
}

func NewState(vm *ViewManager, state *Stack) *State {
	return &State{vm: vm, state: state}
}

func (s *State) Exec(cmd string) error {
	split := strings.Split(cmd, " ")

	switch split[0] {
	case "exec":
		cmd = strings.Join(split[1:], " ")

		s.state.Push(NewForkState(s.vm, s.state, cmd))
		return nil

	case "lua":
		cmd = strings.Join(split[2:], " ")

		if len(split) > 2 && split[1] == "--tty" {
			s.vm.exitMsg = ExitMsg{
				Code: EXIT_CMD,
				Msg:  cmd,
			}

			s.vm.g.UpdateAsync(func(g *gocui.Gui) error {
				return s.vm.Quit(s.vm.g, nil)
			})

			return nil
		}

		s.state.Push(NewLuaState(s.vm, s.state, cmd))
		return nil

	case "go":
		return nil

	}

	return ErrUnknownCommand
}

type ForkState struct {
	vm    *ViewManager
	stdin chan string
	state *Stack
}

func NewForkState(vm *ViewManager, state *Stack, cmd string) *ForkState {
	fork := &ForkState{
		vm:    vm,
		stdin: make(chan string, 5),
		state: state,
	}

	go fork.Start(cmd)
	return fork
}

func (fs *ForkState) Start(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	info := command.Info{
		Command: cmd,
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			return len(msg), fs.vm.SendView("screen", NewData("msg", string(msg)))
		},
		OutFunc: func(msg []byte) (int, error) {
			return len(msg), fs.vm.SendView("screen", NewData("msg", string(msg)))
		},
		InFunc: func(w io.WriteCloser, done <-chan struct{}) (int, error) {
			select {
			case <-done:
				return 0, command.ErrEndOfFile
			case in := <-fs.stdin:
				if len(in) <= 0 {
					return 0, nil
				}

				fs.vm.SendView("screen", NewData("msg", in+"\n"))
				return io.WriteString(w, in+"\r\n")
			}
		},
	}

	fs.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] starting fork %s\n", cmd)))
	err := command.Exec(&info)
	if err != nil {
		fs.vm.SendView("screen", NewData("msg", err.Error()+"\n"))
	}
	close(fs.stdin)
	fs.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] closing fork %s\n", cmd)))

	fs.vm.SendView("header", NewData("msg", "Done - "+cmd))

	fs.state.Pop()
	return nil
}

func (fs *ForkState) Exec(cmd string) error {
	fs.stdin <- cmd
	return nil
}

// TODO: create Lua state
type LuaState struct {
	vm    *ViewManager
	stdin chan string
	state *Stack
}

func NewLuaState(vm *ViewManager, state *Stack, cmd string) *LuaState {
	lua := &LuaState{
		vm:    vm,
		stdin: make(chan string, 5),
		state: state,
	}

	go lua.Start(cmd)
	return lua
}

func (ls *LuaState) Start(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	info := command.Info{
		Command: "lua build-app.lua " + cmd, // TODO: allow user to set lua endpoint
		Dir:     "./lua/",                   // TODO: allow user to set lua direcoty
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			return len(msg), ls.vm.SendView("screen", NewData("msg", string(msg)))
		},
		OutFunc: func(msg []byte) (int, error) {
			return len(msg), ls.vm.SendView("screen", NewData("msg", string(msg)))
		},
		InFunc: func(w io.WriteCloser, done <-chan struct{}) (int, error) {
			select {
			case <-done:
				return 0, command.ErrEndOfFile
			case in := <-ls.stdin:
				ls.vm.SendView("screen", NewData("msg", in+"\n"))
				return io.WriteString(w, in+"\r\n")
			}
		},
	}

	ls.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] starting script %s\n", cmd)))
	err := command.Exec(&info)
	if err != nil {
		ls.vm.SendView("screen", NewData("msg", err.Error()+"\n"))
	}
	close(ls.stdin)
	ls.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] closing script %s\n\n", cmd)))

	ls.vm.SendView("header", NewData("msg", "Done - "+cmd))

	ls.state.Pop()
	return nil
}

func (ls *LuaState) Exec(cmd string) error {
	ls.stdin <- cmd
	return nil
}

// TODO: create GO state
type GoState struct{}
