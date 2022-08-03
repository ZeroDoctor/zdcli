package data

import "github.com/awesome-gocui/gocui"

type ICmdState interface {
	Exec(cmd string) error
	Stop() error
}

type ICmdStateManager interface {
	Exec(cmd string) error
	Stop() error
	SetStack(*Stack)
}

type IView interface {
	Layout(*gocui.Gui) error
	Width() int
	Height() int
	PrintView()
	Display(string)
	Name() string
	Channel() chan interface{}
}

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
