package main

import (
	"context"
	"io"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/tui"
)

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

func update(g *gocui.Gui, vm *tui.ViewManager) {
	go clock(vm)
}

func ExecCommand(vm *tui.ViewManager, StdInChan chan string, com string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info := command.Info{
		Command: com,
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			return len(msg), vm.SendView("screen", tui.NewData("msg", string(msg)))
		},
		OutFunc: func(msg []byte) (int, error) {
			return len(msg), vm.SendView("screen", tui.NewData("msg", string(msg)))
		},
		InFunc: func(w io.WriteCloser) (int, error) {
			vm.SendView("screen", tui.NewData("msg", "starting look for input\n"))
			for {
				in := <-StdInChan
				vm.SendView("screen", tui.NewData("msg", "got back "+in+"\n"))
				return w.Write([]byte(in))
			}
		},
	}

	err := command.Exec(&info)
	if err != nil {
		vm.SendView("screen", tui.NewData("msg", err.Error()+"\n"))
	}

	vm.SendView("header", tui.NewData("msg", "Done - "+com))
	return nil
}
