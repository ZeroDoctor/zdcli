package tui

import (
	"fmt"

	"github.com/zerodoctor/zdcli/tui/data"
)

type CommandManager struct {
	vm           *ViewManager
	stateManager data.Stack
}

func NewCommandManager(vm *ViewManager, state data.ICmdStateManager) *CommandManager {
	cm := &CommandManager{
		vm:           vm,
		stateManager: data.NewStack(),
	}
	state.SetStack(&cm.stateManager)
	cm.stateManager.Push(state)

	return cm
}

func (cm *CommandManager) Cmd(cmd string) {
	if cm.stateManager.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.stateManager.Peek().(data.ICmdState).Exec(cmd)
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s | %s]\n", err.Error(), cmd)))
	}
}

func (cm *CommandManager) Kill() {
	if cm.stateManager.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.stateManager.Peek().(data.ICmdState).Stop()
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s requesting kill]\n", err.Error())))
	}
}
