package command

import (
	"bytes"
	"fmt"
	"io"
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

// func readWithFunc(fn func(msg []byte) (int, error), reader io.ReadCloser, wg *sync.WaitGroup, errChan chan error) {
// 	defer wg.Done()
// 	defer reader.Close()
//
// 	scanner := bufio.NewScanner(reader)
// 	scanner.Split(bufio.ScanLines)
// 	for scanner.Scan() {
// 		if scanner.Err() != nil {
// 			break
// 		}
// 		fn([]byte(scanner.Text() + "\n"))
// 	}
//
// 	if scanner.Err() != nil {
// 		errChan <- scanner.Err()
// 	}
// }

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

	var stdin io.WriteCloser
	var errChan chan error
	if info.InFunc != nil {
		errChan = make(chan error, 2)

		stdin, err = cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdin pipe %s", err.Error())
		}

		go func() {
			defer stdin.Close()
			_, err := info.InFunc(stdin)
			if err != nil {
				errChan <- err
				return
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
	if err = cmd.Wait(); err != nil {
		errs = append(errs, err)
	}

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
