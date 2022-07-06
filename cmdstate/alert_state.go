package cmdstate

// TODO: finish creating alert state
// TODO: create options for alert to run background or not

import (
	"github.com/zerodoctor/zdcli/alert"
	"github.com/zerodoctor/zdcli/tui/data"
)

type AlertState struct {
	vm    data.IViewManager
	state *data.Stack

	alerts []*alert.Alert
}

func (a *AlertState) Exec(cmd string) error {
	return nil
}

func (a *AlertState) Stop() error {
	return nil
}
