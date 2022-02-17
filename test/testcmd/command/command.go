package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"context"
)

// TODO: add std input func to command info

type Info struct {
	Command string
	Args    []string
	Dir     string
	Ctx     context.Context

	CombindOutErr bool

	OutBuffer string
	OutFunc   func(msg []byte) (int, error)

	InFunc func(io.WriteCloser) (int, error)

	ErrBuffer string
	ErrFunc   func(msg []byte) (int, error)
}

func (c Info) parseCommand() (string, []string) {
	split := strings.Split(c.Command, " ")
	return split[0], append(split[1:], c.Args...)
}

type InOut struct {
	fn func(msg []byte) (int, error)
}

func (o *InOut) Write(msg []byte) (int, error) { return o.fn(msg) }
func (o *InOut) Read(msg []byte) (int, error)  { return o.fn(msg) }

func Exec(info *Info) error {
	var err error

	if info.Dir == "" {
		info.Dir = "."
	}

	if info.Ctx == nil {
		info.Ctx = context.Background()
	}

	command, args := info.parseCommand()
	cmd := exec.CommandContext(info.Ctx, command, args...)
	cmd.Dir = info.Dir

	if info.ErrFunc != nil {
		out := &InOut{
			fn: info.ErrFunc,
		}
		cmd.Stderr = out
	}

	if info.OutFunc != nil {
		out := &InOut{
			fn: info.OutFunc,
		}
		cmd.Stdout = out
	}

	done := make(chan bool, 1)
	var errChan chan error
	if info.InFunc != nil {

		_, w := io.Pipe()
		cmd.Stdin = os.Stdin

		go func() {
			select {
			case <-done:
				if err = w.Close(); err != nil {
					fmt.Println("[error]", err.Error())
				}
				return
			default:
			}

			// _, err := info.InFunc(w)
			// if err != nil {
			// 	errChan <- err
			// 	return
			// }

			// err = w.Close()
			// if err != nil {
			// 	fmt.Println("[error]", err.Error())
			// }
		}()
	}

	var berr bytes.Buffer
	if info.ErrFunc == nil {
		cmd.Stderr = &berr
	}

	var bout bytes.Buffer
	if info.OutFunc == nil {
		cmd.Stdout = &bout
	}

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("failed at start [error=%w]", err)
	}

	var errs []error
	fmt.Println("[parent] waiting...")
	if err = cmd.Wait(); err != nil {

		if info.Ctx.Err() != nil {
			errs = append(errs, fmt.Errorf("context [error=%w]", info.Ctx.Err()))
		}

		errs = append(errs, fmt.Errorf("failed at wait [error=%w]", err))
	}
	fmt.Println("[parent] done waiting...")

	close(done)

	info.ErrBuffer = berr.String()
	info.OutBuffer = bout.String()

	for len(errChan) > 0 {
		errs = append(errs, <-errChan)
	}

	return combindErrs(errs)
}

func combindErrs(errs []error) error {
	if errs == nil {
		return nil
	}

	var errStr strings.Builder

	for _, err := range errs {
		errStr.WriteString("[" + err.Error() + "]")
	}

	return fmt.Errorf("%s", errStr.String())
}
