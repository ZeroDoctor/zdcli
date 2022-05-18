package cmdstate

import (
	"context"
	"fmt"
	"time"

	"github.com/zerodoctor/zdcli/command"
	"github.com/zerodoctor/zdcli/config"
	"github.com/zerodoctor/zdcli/tui/comp"
	"github.com/zerodoctor/zdcli/tui/inter"
)

type ForkState struct {
	vm     inter.IViewManager
	stdin  chan string
	state  *comp.Stack
	cancel func()
	cfg *config.Config
}

func NewForkState(vm inter.IViewManager, state *comp.Stack, cmd string, cfg *config.Config) *ForkState {
	fork := &ForkState{
		vm:    vm,
		stdin: make(chan string, 5),
		state: state,
		cfg: cfg,
	}

	go fork.Start(cmd)
	return fork
}

func (fs *ForkState) Start(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	fs.cancel = cancel

	info := command.Info{
		Command: fs.cfg.ShellCmd,
		Args: []string{cmd},
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
			case <-ctx.Done():
			case in := <-fs.stdin:
				fs.vm.SendView("screen", NewData("msg", in+"\n"))
				return in, nil
			default:
			}

			return "", command.ErrStdInNone
		},
	}

	fs.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] starting fork %s\n", cmd)))

	err := command.Exec(&info)
	if err != nil {
		fs.vm.SendView("screen", NewData("msg", err.Error()+"\n"))
	}
	close(fs.stdin)

	fs.vm.SendView("screen", NewData("msg", fmt.Sprintf("[zd] closing fork %s\n", cmd)))
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
