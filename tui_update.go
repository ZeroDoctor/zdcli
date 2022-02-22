package main

import (
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/tui"
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
