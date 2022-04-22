package main

import (
	"errors"
	"log"
	"sync"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/tui"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/tui/inter"
	"github.com/zerodoctor/zdcli/tui/view"
)

func StartTui() comp.ExitMessage {
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	g.Mouse = false
	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan

	vm := tui.NewViewManager(g, []inter.IView{view.NewHeader(g), view.NewScreen(g)}, 1)
	cm := tui.NewCommandManager(vm)

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

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatal(err)
	}

	wg.Wait()
	vm.Wait()
	g.Close()

	return vm.ExitMsg
}
