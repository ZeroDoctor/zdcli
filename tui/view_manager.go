package tui

import (
	"sync"

	"github.com/awesome-gocui/gocui"
)

type View interface {
	Layout(*gocui.Gui) error
	PrintView(*gocui.Gui, *sync.WaitGroup)
	Display(*gocui.Gui, string)
	Name() string
	Channel() chan interface{}
}

type ViewManager struct {
	currentView int
	views       []View
	wg          sync.WaitGroup
}

func NewViewManager(g *gocui.Gui, views []View, currentView int) *ViewManager {
	vm := &ViewManager{
		views:       views,
		currentView: currentView,
	}

	for _, view := range views {
		vm.wg.Add(1)
		go func(view View, wg *sync.WaitGroup) {
			view.PrintView(g, wg)
			wg.Done()
		}(view, &vm.wg)
	}

	return vm
}

func (vm *ViewManager) Layout(g *gocui.Gui) error {
	for _, view := range vm.views {
		if err := view.Layout(g); err != nil {
			panic(err)
		}
	}

	return nil
}

func (vm *ViewManager) Wait() { vm.wg.Wait() }

func (vm *ViewManager) SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	return g.SetViewOnTop(name)
}

// # for keybindings

func (vm *ViewManager) NextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := vm.currentView % len(vm.views)

	name := vm.views[nextIndex].Name()

	if _, err := vm.SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	vm.currentView = nextIndex

	return nil
}

func (vm *ViewManager) Quit(g *gocui.Gui, v *gocui.View) error {
	for _, view := range vm.views {
		close(view.Channel())
	}
	return gocui.ErrQuit
}
