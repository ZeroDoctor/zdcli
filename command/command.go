package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"context"
)

// TODO: add std input func to command info

type Info struct {
	Command string
	Args    []string
	Dir     string
	Ctx     context.Context

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

func removeEmpty(buffer []byte) []byte {
	first := -1

	var result []byte
	for i, b := range buffer {
		if b != 0 && first == -1 {
			first = i
		}

		if b == 0 && first != -1 {
			result = append(result, buffer[first:i]...)
			first = -1
		}
	}

	if first != -1 {
		result = append(result, buffer[first:]...)
	}

	return result
}

func readWithFunc(bufferSize int, fn func(msg []byte) (int, error), reader io.ReadCloser, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	defer reader.Close()

	buffer := make([]byte, bufferSize)
	for {
		_, err := reader.Read(buffer)
		if err != nil {
			errChan <- fmt.Errorf("closing stream: [error=%s]", err)
			return
		}
		buffer = removeEmpty(buffer)
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
	hasFunc := info.ErrFunc != nil || info.OutFunc != nil

	if info.Dir == "" {
		info.Dir = "."
	}
	if info.ErrBufSize <= 0 {
		info.ErrBufSize = 1
	}
	if info.OutBufSize <= 0 {
		info.OutBufSize = 1
	}
	if info.Ctx == nil {
		info.Ctx = context.Background()
	}

	command, args := info.parseCommand()
	cmd := exec.CommandContext(info.Ctx, command, args...)
	cmd.Dir = info.Dir

	var errChan chan error
	var doneChan chan bool
	if hasFunc {
		errChan = make(chan error, 2)
		doneChan = make(chan bool, 1)
	}

	var wg sync.WaitGroup
	if info.ErrFunc != nil {
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to get stderr pipe %s", err.Error())
		}
		wg.Add(1)
		go readWithFunc(info.ErrBufSize, info.ErrFunc, stderr, &wg, errChan)
	}

	if info.OutFunc != nil {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout pipe %s", err.Error())
		}
		wg.Add(1)
		go readWithFunc(info.OutBufSize, info.OutFunc, stdout, &wg, errChan)
	}

	var berr bytes.Buffer
	if info.ErrFunc == nil {
		cmd.Stderr = &berr
	}

	var bout bytes.Buffer
	if info.OutFunc == nil {
		cmd.Stdout = &bout
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s", err.Error())
	}

	if hasFunc {
		go func() {
			wg.Wait()
			time.Sleep(500 * time.Millisecond)
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
	if err != nil {
		cmd.Wait()
		info.ErrBuffer = berr.String()
		info.OutBuffer = bout.String()

		return err
	}

	err = cmd.Wait()

	info.ErrBuffer = berr.String()
	info.OutBuffer = bout.String()

	return err
}
