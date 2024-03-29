package cmdstate

import (
	"context"
	"fmt"
	"time"

	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdtui/data"
)

type LuaState struct {
	vm     data.IViewManager
	stdin  chan string
	state  *data.Stack
	cancel func()
	cfg    *config.Config
}

func NewLuaState(vm data.IViewManager, state *data.Stack, cmd string, cfg *config.Config) *LuaState {
	lua := &LuaState{
		vm:    vm,
		stdin: make(chan string, 5),
		state: state,
		cfg:   cfg,
	}

	go lua.Start(cmd)
	return lua
}

func (ls *LuaState) Start(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	ls.cancel = cancel

	flags := fmt.Sprintf("--os_i %s --arch_i %s ", ls.cfg.OS, ls.cfg.Arch)
	info := command.Info{
		Command: ls.cfg.ShellCmd,
		Args:    []string{ls.cfg.LuaCmd + " build-app.lua " + flags + cmd},
		Dir:     ls.cfg.RootLuaScriptDir,
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			return len(msg), ls.vm.SendView("screen", NewData("msg", string(msg)))

		},
		OutFunc: func(msg []byte) (int, error) {
			return len(msg), ls.vm.SendView("screen", NewData("msg", string(msg)))

		},
		InFunc: func(ctx context.Context) (string, error) {
			time.Sleep(100 * time.Millisecond)

			select {
			case in := <-ls.stdin:
				ls.vm.SendView("screen", NewData("msg", in+"\n"))
				return in, nil
			default:
			}

			return "", command.ErrStdInNone
		},
	}

	ls.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] starting script [cmd=%s]\n", cmd)))

	err := command.Exec(&info)
	if err != nil {
		ls.vm.SendView("screen", NewData("msg", err.Error()+"\n"))
	}
	close(ls.stdin)

	ls.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] closing script [cmd=%s]\n\n", cmd)))
	ls.vm.SendView("header", NewData("msg", "Done - "+cmd))

	ls.state.Pop()
	return nil
}

func (ls *LuaState) Exec(cmd string) error {
	ls.stdin <- cmd
	return nil
}

func (ls *LuaState) Stop() error {
	if ls.cancel != nil {
		ls.cancel()
		return nil
	}

	return ErrCommandNotRunning
}
