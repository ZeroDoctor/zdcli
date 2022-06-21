package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInput struct {
	Input textinput.Model
	err   error
}

func NewTextInput() *TextInput {
	input := textinput.New()
	input.Focus()
	return &TextInput{Input: input}
}

func (i *TextInput) Init() tea.Cmd { return textinput.Blink }

func (i *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
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
