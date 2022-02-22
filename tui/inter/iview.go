package inter

import "github.com/awesome-gocui/gocui"

type IView interface {
	Layout(*gocui.Gui) error
	PrintView()
	Display(string)
	Name() string
	Channel() chan interface{}
}
