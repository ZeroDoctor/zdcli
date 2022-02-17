package tui

import (
	"errors"
	"strings"
)

var ErrUnknownCommand error = errors.New("unknown command")

type ICommand interface {
	Exec(cmd string) error
}

type CommandManager struct {
	state Stack
}

func NewCommandManager() *CommandManager {
	cm := &CommandManager{
		state: NewStack(),
	}

	cm.state.Push(cm)

	return cm
}

type InitState struct{}

func NewInitState() *InitState {
	return &InitState{}
}

func (is *InitState) Exec(cmd string) error {
	split := strings.Split(cmd, " ")

	switch split[0] {
	case "app":
		return nil

	case "go":
		return nil

	}

	return ErrUnknownCommand
}

type ForkState struct {
	VM    *ViewManager
	Fork  func(*ViewManager, chan string, string) error
	StdIn chan string
}

func NewForkState() *ForkState {
	return &ForkState{}
}

func (fs *ForkState) Exec(cmd string) error {

}

/**

# command design

## callback

### fork

### goroutine

*/
