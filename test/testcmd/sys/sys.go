package sys

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Process struct {
	cmd  string
	args []string
	attr *os.ProcAttr

	*os.Process
}

func NewProcess(cmd string, args []string, dir string) (*Process, error) {
	if cmd == "" {
		return nil, errors.New("command cannot be empty")
	}

	var err error

	split := strings.Split(cmd, " ")
	if len(split) > 1 {
		args = append(split[:], args...)
	}

	cmd, err = exec.LookPath(split[0])
	if err != nil {
		return nil, err
	}

	return &Process{
		cmd:  cmd,
		args: args,
		attr: &os.ProcAttr{
			Dir:   dir,
			Env:   os.Environ(),
			Files: []*os.File{nil, nil, nil},
		},
	}, nil
}

func (p *Process) SetInWriter(f *os.File) { p.attr.Files[0] = f }

func (p *Process) SetOutReader(f *os.File) { p.attr.Files[1] = f }

func (p *Process) SetErrReader(f *os.File) { p.attr.Files[2] = f }

func (p *Process) Start() error {
	var err error
	fmt.Printf("[cmd=%s] [args=%v]\n", p.cmd, p.args)
	p.Process, err = os.StartProcess(p.cmd, p.args, p.attr)
	return err
}
