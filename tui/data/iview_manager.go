package data

import (
	"github.com/awesome-gocui/gocui"
)

type IViewManager interface {
	SendView(string, interface{}) error
	GetView(string) (IView, error)
	AddView(*gocui.Gui, IView) error
	RemoveView(*gocui.Gui, string) error
	NextView(*gocui.Gui, *gocui.View) error
	SetCurrentViewOnTop(*gocui.Gui, string) (*gocui.View, error)
	Quit(*gocui.Gui, *gocui.View) error
	G() *gocui.Gui
	SetExitMsg(ExitMessage)
}
