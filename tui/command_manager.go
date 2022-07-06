package tui

import (
	"fmt"

	"github.com/zerodoctor/zdcli/tui/data"
)

type CommandManager struct {
	vm    *ViewManager
	state data.Stack
}

func NewCommandManager(vm *ViewManager, stateManager data.ICmdState) *CommandManager {
	cm := &CommandManager{
		vm:    vm,
		state: data.NewStack(),
	}

	cm.state.Push(stateManager)

	return cm
}

func (cm *CommandManager) Cmd(cmd string) {
	if cm.state.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.state.Peek().(data.ICmdState).Exec(cmd)
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s | %s]\n", err.Error(), cmd)))
	}
}

func (cm *CommandManager) Kill() {
	if cm.state.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.state.Peek().(data.ICmdState).Stop()
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s requesting kill]\n", err.Error())))
	}
}
