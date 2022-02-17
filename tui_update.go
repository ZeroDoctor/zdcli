package main

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/tui"
	"github.com/zerodoctor/zdcli/tui/view"
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
			vm.SendView("header", view.NewData("clock", str))
		}
	}

}

func ExecCommand(vm *tui.ViewManager, StdInChan chan string, com string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	info := command.Info{
		Command: com,
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			return len(msg), vm.SendView("screen", view.NewData("msg", string(msg)))
		},
		OutFunc: func(msg []byte) (int, error) {
			return len(msg), vm.SendView("screen", view.NewData("msg", string(msg)))
		},
		InFunc: func(w io.WriteCloser, done <-chan struct{}) (int, error) {
			select {
			case <-done:
				return 0, errors.New("EOF")
			case in := <-StdInChan:
				if len(in) <= 0 {
					return 0, nil
				}

				n, err := w.Write([]byte(in + "\r\n"))
				if err != nil {
					vm.SendView("screen", view.NewData("msg", err.Error()+"\n"))
				}

				return n, err
			}
		},
	}

	err := command.Exec(&info)
	if err != nil {
		vm.SendView("screen", view.NewData("msg", err.Error()+"\n"))
	}

	vm.SendView("header", view.NewData("msg", "Done - "+com))
	return nil
}
