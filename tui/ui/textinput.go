package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TIOption func(*TextInput)

func WithTIPassword() TIOption {
	return func(ti *TextInput) {
		ti.Input.EchoMode = textinput.EchoPassword
	}
}

type TextInput struct {
	WasCancel bool
	Input     textinput.Model
	err       error
}

func NewTextInput(options ...TIOption) *TextInput {
	input := textinput.New()
	ti := &TextInput{Input: input}

	for _, opt := range options {
		opt(ti)
	}

	return ti
}

func (i *TextInput) Init() tea.Cmd { return textinput.Blink }

func (i *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			i.WasCancel = true
		case tea.KeyEnter:
			return i, tea.Quit
		}
	case error:
		i.err = msg
		return i, nil
	}

	var cmd tea.Cmd
	i.Input, cmd = i.Input.Update(msg)

	return i, cmd
}

func (i *TextInput) View() string {
	return fmt.Sprintf("%s", i.Input.View())
}

func (i *TextInput) Focus() tea.Cmd { return i.Input.Focus() }
func (i *TextInput) Blur()          { i.Input.Blur() }
