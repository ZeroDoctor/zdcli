package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/cmdstate"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/tui"
	"github.com/zerodoctor/zdcli/tui/data"
	"github.com/zerodoctor/zdcli/tui/view"
)

func update(g *gocui.Gui, vm *tui.ViewManager, wg *sync.WaitGroup) {
	wg.Add(1)
	go clock(vm, wg)
}

func clock(vm *tui.ViewManager, wg *sync.WaitGroup) {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	defer wg.Done()

	for {
		select {
		case <-vm.Shutdown():
			return
		case <-tick.C:
			str := time.Now().Format("02/01/2006 15:04:05")
			vm.SendView("header", tui.NewData("clock", str))
		}
	}
}

func StartTui(cfg *config.Config) data.ExitMessage {
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	g.Mouse = false
	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan

	stack := data.NewStack()
	vm := tui.NewViewManager(g, []data.IView{view.NewHeader(g), view.NewScreen(g)}, 1)
	cm := tui.NewCommandManager(vm, cmdstate.NewState(vm, &stack, cfg))

	g.SetManagerFunc(vm.Layout)
	km := tui.NewKeyManager(g, vm)

	km.SetKey("screen", rune('j'), gocui.ModNone, tui.DownScreen)
	km.SetKey("screen", rune('k'), gocui.ModNone, tui.UpScreen)

	var wg sync.WaitGroup
	go update(g, vm, &wg)

	err = vm.AddView(g, view.NewCommand(g, cm))
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = vm.SetCurrentView(g, "command"); err != nil {
		log.Fatal(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatal(err)
	}

	wg.Wait()
	vm.Wait()
	g.Close()

	return vm.ExitMsg
}
