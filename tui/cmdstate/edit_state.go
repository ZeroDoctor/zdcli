package cmdstate

import (
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/tui/inter"
)

type EditState struct {
	vm     inter.IViewManager
	state  *comp.Stack
	cancel func()
}
