package inter

import (
	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/tui/comp"
)

type IViewManager interface {
	SendView(string, interface{}) error
	AddView(*gocui.Gui, IView) error
	RemoveView(*gocui.Gui, string) error
	NextView(*gocui.Gui, *gocui.View) error
	SetCurrentViewOnTop(*gocui.Gui, string) (*gocui.View, error)
	Quit(*gocui.Gui, *gocui.View) error
	G() *gocui.Gui
	SetExitMsg(comp.ExitMessage)
}
