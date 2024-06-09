package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"context"
)

var (
	ErrEndOfFile error = errors.New("EOF")
	ErrStdInNone error = errors.New("got nothing from stdin")
)

type Info struct {
	Command string
	Args    []string
	Env     []string
	Dir     string
	Ctx     context.Context

	CombindOutErr bool

	OutBuffer string
	OutFunc   func([]byte) (int, error)
	Stdout    *os.File

	InFunc func(context.Context) (string, error)
	InChan chan string
	Stdin  *os.File

	ErrBuffer string
	ErrFunc   func([]byte) (int, error)
	Stderr    *os.File
}

func (c Info) parseCommand() (string, []string) {
	split := strings.Split(c.Command, " ")
	return split[0], append(split[1:], c.Args...)
}

type Output struct {
	fn func(msg []byte) (int, error)
}

func (o *Output) Write(msg []byte) (int, error) { return o.fn(msg) }

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
	cmd.Env = info.Env
	// cmd.SysProcAttr = &syscall.SysProcAttr{}

	if info.ErrFunc != nil {
		out := &Output{
			fn: info.ErrFunc,
		}
		cmd.Stderr = out
	}

	if info.OutFunc != nil {
		out := &Output{
			fn: info.OutFunc,
		}
		cmd.Stdout = out
	}

	var wg sync.WaitGroup
	done := make(chan struct{}, 1)
	errChan := make(chan error, 4)

	if info.InFunc != nil {
		r, w, err := os.Pipe()
		if err != nil {
			return err
		}
		cmd.Stdin = r

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() {
				if err = w.Close(); err != nil {
					errChan <- err
				}
			}()

			for {
				select {
				case <-info.Ctx.Done():
					return
				case <-done:
					return
				default:
					in, err := info.InFunc(info.Ctx)
					if err == ErrStdInNone || len(in) <= 0 {
						continue
					}

					if err != nil {
						errChan <- err
						return
					}

					_, err = io.WriteString(w, in+LINEEND)
					if err != nil {
						errChan <- err
						return
					}
				}
			}
		}(&wg)
	}

	var berr bytes.Buffer
	if info.ErrFunc == nil {
		cmd.Stderr = &berr
	}

	var bout bytes.Buffer
	if info.OutFunc == nil {
		cmd.Stdout = &bout
	}

	if info.Stderr != nil {
		cmd.Stderr = info.Stderr
	}

	if info.Stdout != nil {
		cmd.Stdout = info.Stdout
	}

	if info.Stdin != nil {
		cmd.Stdin = info.Stdin
	}

	if err = cmd.Start(); err != nil {
		return fmt.Errorf("failed at start [error=%w]", err)
	}

	var errs []error
	if err = cmd.Wait(); err != nil {
		if info.Ctx.Err() != nil {
			errs = append(errs, fmt.Errorf("context [error=%w] [code=%s]", info.Ctx.Err(), err.Error()))
		} else {
			errs = append(errs, fmt.Errorf("failed at wait [error=%w]", err))
		}
	}

	close(done)
	wg.Wait()

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
