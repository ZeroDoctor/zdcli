package cmdstate

import (
	"context"
	"fmt"
	"time"

	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/tui/inter"
	"github.com/zerodoctor/zdcli/util"
)

type LuaState struct {
	vm     inter.IViewManager
	stdin  chan string
	state  *comp.Stack
	cancel func()
}

func NewLuaState(vm inter.IViewManager, state *comp.Stack, cmd string) *LuaState {
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
	ls.cancel = cancel

	info := command.Info{
		Command: "lua build-app.lua " + cmd, // TODO: allow user to set lua endpoint
		Dir:     util.EXEC_PATH + "/lua/",   // TODO: allow user to set lua direcoty
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
			case <-ctx.Done():
			case in := <-ls.stdin:
				ls.vm.SendView("screen", NewData("msg", in+"\n"))
				return in, nil
			default:
			}

			return "", command.ErrStdInNone
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

func (ls *LuaState) Stop() error {
	if ls.cancel != nil {
		ls.cancel()
		return nil
	}

	return ErrCommandNotRunning
}
