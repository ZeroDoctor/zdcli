package tui

import (
	"fmt"

	"github.com/zerodoctor/zdcli/tui/cmdstate"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/tui/inter"
)

type CommandManager struct {
	vm    *ViewManager
	state comp.Stack
}

func NewCommandManager(vm *ViewManager) *CommandManager {
	cm := &CommandManager{
		vm:    vm,
		state: comp.NewStack(),
	}

	cm.state.Push(cmdstate.NewState(vm, &cm.state))

	return cm
}

func (cm *CommandManager) Cmd(cmd string) {
	if cm.state.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.state.Peek().(inter.ICommandState).Exec(cmd)
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s | %s]\n", err.Error(), cmd)))
	}
}

func (cm *CommandManager) Kill() {
	if cm.state.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.state.Peek().(inter.ICommandState).Stop()
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s requesting kill]\n", err.Error())))
	}
}
