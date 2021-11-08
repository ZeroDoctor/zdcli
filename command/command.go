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
	command string
	args    []string
	dir     string
	ctx     context.Context
	verbose bool

	outBufSize int
	outBuffer  string
	outFunc    func(msg []byte) (int, error)

	errBufSize int
	errBuffer  string
	errFunc    func(msg []byte) (int, error)
}

func (c Info) parseCommand() (string, []string) {
	split := strings.Split(c.command, " ")
	return split[0], append(split[1:], c.args...)
}

func readWithFunc(bufferSize int, fn func(msg []byte) (int, error), reader io.ReadCloser, wg *sync.WaitGroup, errChan chan error) {
	wg.Add(1)
	defer wg.Done()
	if fn == nil {
		return
	}

	buffer := make([]byte, bufferSize)
	for {
		_, err := reader.Read(buffer)
		if err != nil {
			fmt.Printf("closing stream: [error=%s]\n", err)
			reader.Close()
			return
		}
		if len(buffer) > 0 {
			_, err := fn(buffer)
			if err != nil {
				errChan <- err
				return
			}
		}
	}
}

func Exec(info *Info) error {
	hasFunc := info.outFunc != nil || info.errFunc != nil

	if info.dir == "" {
		info.dir = "."
	}
	if info.outBufSize <= 0 {
		info.outBufSize = 1
	}
	if info.errBufSize <= 0 {
		info.errBufSize = 1
	}
	if info.ctx == nil {
		info.ctx = context.Background()
	}

	command, args := info.parseCommand()
	cmd := exec.CommandContext(info.ctx, command, args...)
	cmd.Dir = info.dir
	if info.verbose {
		fmt.Printf("exec:\n\t[command=%s]\n\t[args=%v]\n\t[dir=%s]\n", command, args, info.dir)
	}

	var errChan chan error
	var doneChan chan bool
	if hasFunc {
		errChan = make(chan error, 2)
		doneChan = make(chan bool, 1)
	}

	var wg sync.WaitGroup
	if info.outFunc != nil {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout pipe %s", err.Error())
		}
		go readWithFunc(info.outBufSize, info.outFunc, stdout, &wg, errChan)
	}

	if info.errFunc != nil {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to get stderr pipe %s", err.Error())
		}
		go readWithFunc(info.errBufSize, info.errFunc, stderr, &wg, errChan)
	}

	var bout bytes.Buffer
	if info.outFunc == nil {
		cmd.Stdout = &bout
	}

	var berr bytes.Buffer
	if info.errFunc == nil {
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

	info.outBuffer = bout.String()
	info.errBuffer = berr.String()

	return err
}
