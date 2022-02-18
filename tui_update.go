package main

import (
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/tui"
)

func update(g *gocui.Gui, vm *tui.ViewManager) {
	go clock(vm)
}

func clock(vm *tui.ViewManager) {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

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
