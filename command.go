package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/tui"
	"github.com/zerodoctor/zdcli/view"
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
			vm.SendView("header", view.NewData("clock", str))
		}
	}

}

func update(g *gocui.Gui, vm *tui.ViewManager) {
	go clock(vm)
}

func ExecCommand(vm *tui.ViewManager) func(*gocui.View) error {
	return func(v *gocui.View) error {
		com := v.Buffer()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		info := command.Info{
			Command: com,
			Ctx:     ctx,

			ErrFunc: func(msg []byte) (int, error) {
				vm.SendView("screen", view.NewData("msg", string(msg)))
				return len(msg), nil
			},
			OutFunc: func(msg []byte) (int, error) {
				vm.SendView("screen", view.NewData("msg", string(msg)))
				return len(msg), nil
			},
			InFunc: func(wc io.WriteCloser) error {
				var line string
				_, err := fmt.Scanln(&line)
				if err != nil {
					return err
				}
				wc.Write([]byte(line))

				return nil
			},
		}

		err := command.Exec(&info)
		if err != nil {
			vm.SendView("screen", view.NewData("msg", err.Error()+"\n"))
		}

		vm.SendView("header", view.NewData("msg", "Done - "+com))
		return nil
	}
}
