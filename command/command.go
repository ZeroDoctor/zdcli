package command

import (
	"bufio"
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

	CombindOutErr bool

	OutBuffer string
	OutFunc   func(msg []byte) (int, error)

	InFunc func(io.WriteCloser) error

	ErrBuffer string
	ErrFunc   func(msg []byte) (int, error)
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

func readWithFunc(fn func(msg []byte) (int, error), reader io.ReadCloser, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if scanner.Err() != nil {
			errChan <- scanner.Err()
		}
		fn(scanner.Bytes())
	}

	if scanner.Err() != nil {
		errChan <- scanner.Err()
	}
}

func Exec(info *Info) error {
	var err error
	hasFunc := info.ErrFunc != nil || info.OutFunc != nil || info.InFunc != nil

	if info.Dir == "" {
		info.Dir = "."
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
		go readWithFunc(info.ErrFunc, stderr, &wg, errChan)
	}

	if info.OutFunc != nil {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdout pipe %s", err.Error())
		}
		wg.Add(1)
		go readWithFunc(info.OutFunc, stdout, &wg, errChan)
	}

	stdInChan := make(chan bool, 1)
	if info.InFunc != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdin pipe %s", err.Error())
		}

		go func() {
			defer stdin.Close()
			select {
			case <-stdInChan:
			default:
				err := info.InFunc(stdin)
				if err != nil {
					errChan <- err
					return
				}

			}
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

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s", err.Error())
	}

	var errs []error
	if hasFunc {
		go func() {
			wg.Wait()
			time.Sleep(100 * time.Millisecond)
			close(doneChan)
		}()

	loop:
		for {
			select {
			case <-doneChan:
				break loop
			case err := <-errChan:
				errs = append(errs, err)
			}
		}
		close(errChan)
	}

	if errs != nil {
		cmd.Wait()

		info.ErrBuffer = berr.String()
		info.OutBuffer = bout.String()

		close(stdInChan)

		return combindErrs(errs)
	}

	err = cmd.Wait()

	info.ErrBuffer = berr.String()
	info.OutBuffer = bout.String()

	close(stdInChan)

	return combindErrs(errs)
}

func combindErrs(errs []error) error {
	var errStr strings.Builder

	for _, err := range errs {
		errStr.WriteString("[" + err.Error() + "]")
	}

	return fmt.Errorf("%s", errStr.String())
}
