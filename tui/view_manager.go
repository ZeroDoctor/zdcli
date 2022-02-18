package tui

import (
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
)

type exit int8

const (
	EXIT_SUC exit = iota
	EXIT_CMD
)

type ExitMsg struct {
	Code exit
	Msg  string
}

type View interface {
	Layout(*gocui.Gui) error
	PrintView()
	Display(string)
	Name() string
	Channel() chan interface{}
}

type ViewManager struct {
	currentView int
	views       []View
	wg          sync.WaitGroup
	g           *gocui.Gui

	shutdown chan bool
	exitMsg  ExitMsg
}

func NewViewManager(g *gocui.Gui, views []View, currentView int) *ViewManager {
	vm := &ViewManager{
		views:       views,
		currentView: currentView,
		shutdown:    make(chan bool),
		g:           g,

		exitMsg: ExitMsg{Code: EXIT_SUC}, // TODO: redo exit handling
	}

	for _, view := range views {
		vm.wg.Add(1)
		go func(view View, wg *sync.WaitGroup) {
			view.PrintView()
			wg.Done()
		}(view, &vm.wg)
	}

	return vm
}

func (vm *ViewManager) Layout(g *gocui.Gui) error {
	// TODO: handle view collisions
	for _, view := range vm.views {
		if err := view.Layout(g); err != nil {
			return err
		}
	}

	if _, err := g.SetCurrentView(vm.views[vm.currentView].Name()); err != nil {
		return err
	}

	return nil
}

func (vm *ViewManager) Wait()               { vm.wg.Wait() }
func (vm *ViewManager) Shutdown() chan bool { return vm.shutdown }

func (vm *ViewManager) SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	return g.SetViewOnTop(name)
}

func (vm ViewManager) SendView(viewname string, data interface{}) error {
	for _, view := range vm.views {
		if view.Name() == viewname {
			retryCount := 3
			count := 0

		loop:
			for count < retryCount {
				select {
				case view.Channel() <- data:
					break loop
				default:
					time.Sleep(100 * time.Millisecond)
					count++
				}
			}

			return nil
		}
	}

	return gocui.ErrUnknownView
}

func (vm *ViewManager) AddView(g *gocui.Gui, view View) error {
	vm.wg.Add(1)
	go func(wg *sync.WaitGroup) {
		view.PrintView()
		wg.Done()
	}(&vm.wg)

	if err := view.Layout(g); err != nil {
		return err
	}

	if _, err := vm.SetCurrentViewOnTop(g, view.Name()); err != nil {
		return err
	}

	vm.views = append(vm.views, view)

	return nil
}

// TODO: remove view

func (vm *ViewManager) RemoveView(g *gocui.Gui, name string) error {
	var view View
	var index int

	for i, v := range vm.views {
		if v.Name() == name {
			view = v
			index = i
		}
	}

	close(view.Channel())
	g.DeleteView(name)

	// remove view from slice
	vm.views[index] = vm.views[len(vm.views)-1]
	vm.views = vm.views[:len(vm.views)-1]

	return nil
}

// # for keybindings

func (vm *ViewManager) NextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (vm.currentView + 1) % len(vm.views)

	name := vm.views[nextIndex].Name()

	if _, err := vm.SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	vm.currentView = nextIndex

	return nil
}

func (vm *ViewManager) Quit(g *gocui.Gui, v *gocui.View) error {
	close(vm.shutdown)
	for _, view := range vm.views {
		if view.Channel() != nil {
			close(view.Channel())
		}
	}

	return gocui.ErrQuit
}

func (vm *ViewManager) ExitMsg() ExitMsg { return vm.exitMsg }
