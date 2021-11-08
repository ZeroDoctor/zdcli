package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"context"
)

// TODO: add std input func to command info

type Info struct {
	Command string
	Args    []string
	Dir     string
	Ctx     context.Context
	Verbose bool

	OutBufSize int
	OutBuffer  string
	OutFunc    func(msg []byte) (int, error)

	ErrBufSize int
	ErrBuffer  string
	ErrFunc    func(msg []byte) (int, error)
}

func (c Info) parseCommand() (string, []string) {
	split := strings.Split(c.Command, " ")
	return split[0], append(split[1:], c.Args...)
}

func readWithFunc(bufferSize int, fn func(msg []byte) (int, error), reader io.ReadCloser, wg *sync.WaitGroup, errChan chan error) {
	wg.Add(1)
	defer wg.Done()
	defer reader.Close()
	if fn == nil {
		return
	}

	buffer := make([]byte, bufferSize)
	for {
		_, err := reader.Read(buffer)
		if err != nil {
			fmt.Printf("closing stream: [error=%s]\n", err)
			return
		}
		if len(buffer) > 0 {
			_, err := fn(buffer)
			if err != nil {
				errChan <- err
				return
			}
		}
		buffer = make([]byte, bufferSize)
	}
}

func Exec(info *Info) error {
	hasFunc := info.OutFunc != nil || info.ErrFunc != nil

	if info.Dir == "" {
		info.Dir = "."
	}
	if info.OutBufSize <= 0 {
		info.OutBufSize = 1
	}
	if info.ErrBufSize <= 0 {
		info.ErrBufSize = 1
	}
	if info.Ctx == nil {
		info.Ctx = context.Background()
	}

	command, args := info.parseCommand()
	cmd := exec.CommandContext(info.Ctx, command, args...)
	cmd.Dir = info.Dir
	if info.Verbose {
		fmt.Printf("exec:\n\t[command=%s]\n\t[args=%v]\n\t[dir=%s]\n", command, args, info.Dir)
	}

	var errChan chan error
	var doneChan chan bool
	if hasFunc {
		errChan = make(chan error, 2)
		doneChan = make(chan bool, 1)
	}

	var wg sync.WaitGroup
	if info.OutFunc != nil {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout pipe %s", err.Error())
		}
		go readWithFunc(info.OutBufSize, info.OutFunc, stdout, &wg, errChan)
	}

	if info.ErrFunc != nil {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to get stderr pipe %s", err.Error())
		}
		go readWithFunc(info.ErrBufSize, info.ErrFunc, stderr, &wg, errChan)
	}

	var bout bytes.Buffer
	if info.OutFunc == nil {
		cmd.Stdout = &bout
	}

	var berr bytes.Buffer
	if info.ErrFunc == nil {
		cmd.Stderr = &berr
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s", err.Error())
	}

	if hasFunc {
		go func() {
			wg.Wait()
			close(doneChan)
		}()

	loop:
		for {
			select {
			case <-doneChan:
				break loop
			case newErr := <-errChan:
				if err != nil {
					err = fmt.Errorf("[%v] [%v]", err, newErr)
					continue
				}
				err = newErr
			}
		}
		close(errChan)
	}
	err = cmd.Wait()

	info.OutBuffer = bout.String()
	info.ErrBuffer = berr.String()

	return err
}
