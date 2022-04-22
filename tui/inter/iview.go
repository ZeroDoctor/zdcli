package inter

import "github.com/awesome-gocui/gocui"

type IView interface {
	Layout(*gocui.Gui) error
	Width() int
	Height() int
	PrintView()
	Display(string)
	Name() string
	Channel() chan interface{}
}
