package tui

import (
	"fmt"

	"github.com/zerodoctor/zdcli/tui/data"
)

type CommandManager struct {
	vm    *ViewManager
	stack data.Stack
}

func NewCommandManager(vm *ViewManager, state data.ICmdStateManager) *CommandManager {
	cm := &CommandManager{
		vm:    vm,
		stack: data.NewStack(),
	}
	state.SetStack(&cm.stack)
	cm.stack.Push(state)

	return cm
}

func (cm *CommandManager) Cmd(cmd string) {
	if cm.stack.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.stack.Peek().(data.ICmdState).Exec(cmd)
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s | %s]\n", err.Error(), cmd)))
	}
}

func (cm *CommandManager) Kill() {
	if cm.stack.Len() <= 0 {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=state slice is empty]\n")))
		return
	}

	err := cm.stack.Peek().(data.ICmdState).Stop()
	if err != nil {
		cm.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] [error=%s requesting kill]\n", err.Error())))
	}
}
