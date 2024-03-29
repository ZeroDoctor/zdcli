package cmdstate

import (
	"context"
	"fmt"
	"time"

	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdtui/data"
)

type ForkState struct {
	vm     data.IViewManager
	stdin  chan string
	state  *data.Stack
	cancel func()
	cfg    *config.Config
}

func NewForkState(vm data.IViewManager, state *data.Stack, cmd string, cfg *config.Config) *ForkState {
	fork := &ForkState{
		vm:    vm,
		stdin: make(chan string, 5),
		state: state,
		cfg:   cfg,
	}

	go fork.Start(cmd)
	return fork
}

func (fs *ForkState) Start(cmd string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fs.cancel = cancel

	info := command.Info{
		Command: fs.cfg.ShellCmd,
		Args:    []string{cmd},
		Ctx:     ctx,

		ErrFunc: func(msg []byte) (int, error) {
			return len(msg), fs.vm.SendView("screen", NewData("msg", string(msg)))

		},
		OutFunc: func(msg []byte) (int, error) {
			return len(msg), fs.vm.SendView("screen", NewData("msg", string(msg)))

		},
		InFunc: func(ctx context.Context) (string, error) {
			time.Sleep(100 * time.Millisecond)

			select {
			case in := <-fs.stdin:
				fs.vm.SendView("screen", NewData("msg", in+"\n"))
				return in, nil
			default:
			}

			return "", command.ErrStdInNone
		},
	}

	fs.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] starting fork [cmd=%s]\n", cmd)))

	err := command.Exec(&info)
	if err != nil {
		fs.vm.SendView("screen", NewData("msg", err.Error()+"\n"))
	}
	close(fs.stdin)

	fs.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] closing fork [cmd=%s]\n", cmd)))
	fs.vm.SendView("header", NewData("msg", "Done - "+cmd))

	fs.state.Pop()
	return nil
}

func (fs *ForkState) Exec(cmd string) error {
	fs.stdin <- cmd
	return nil
}

func (fs *ForkState) Stop() error {
	if fs.cancel != nil {
		fs.cancel()
		return nil
	}

	return ErrCommandNotRunning
}
