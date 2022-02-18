package main

import (
	"errors"
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/tui"
	"github.com/zerodoctor/zdcli/tui/view"
)

func StartTui() {
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Fatal(err.Error())
	}

	g.Mouse = false
	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan

	vm := tui.NewViewManager(g, []tui.View{view.NewHeader(g), view.NewScreen(g)}, 1)
	cm := tui.NewCommandManager(vm)

	g.SetManagerFunc(vm.Layout)
	km := tui.NewKeyManager(g, vm)

	km.SetKey("screen", rune('j'), gocui.ModNone, tui.DownScreen)
	km.SetKey("screen", rune('k'), gocui.ModNone, tui.UpScreen)

	go update(g, vm)

	err = vm.AddView(g, view.NewCommand(g, cm))
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatal(err)
	}

	vm.Wait()
	g.Close()

}
