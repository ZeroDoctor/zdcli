package cmdstate

// TODO: finish creating alert state
// TODO: create options for alert to run background or not

import (
	"github.com/zerodoctor/zdcli/alert"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/tui/inter"
)

type AlertState struct {
	vm    inter.IViewManager
	state *comp.Stack

	alerts []*alert.Alert
}

func (a *AlertState) Exec(cmd string) error {
	return nil
}

func (a *AlertState) Stop() error {
	return nil
}
