package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputForm struct {
	focusIndex int
	Inputs     []*TextInput
	err        error
	WasCancel  bool
}

func NewTextInputForm(inputs ...*TextInput) *TextInputForm {
	return &TextInputForm{Inputs: inputs}
}

func (t *TextInputForm) Init() tea.Cmd {
	return textinput.Blink
}

func (t *TextInputForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if t.focusIndex == len(t.Inputs) {
				break
			}

			placeHolder := t.Inputs[t.focusIndex].Input.Placeholder
			t.Inputs[t.focusIndex].Input.SetValue(placeHolder)

			t.focusIndex++
			if t.focusIndex > len(t.Inputs) {
				t.focusIndex = len(t.Inputs)
			}

			cmds := make([]tea.Cmd, len(t.Inputs))
			for i := range t.Inputs {
				if i == t.focusIndex {
					cmds[i] = t.Inputs[i].Focus()
					continue
				}

				t.Inputs[i].Blur()
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			t.WasCancel = true
			return t, tea.Quit
		case tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			switch msg.Type {
			case tea.KeyEnter:
				t.focusIndex++
				if t.focusIndex >= len(t.Inputs) {
					return t, tea.Quit
				}
			case tea.KeyUp:
				t.focusIndex--
			case tea.KeyDown:
				t.focusIndex++
			}

			if t.focusIndex > len(t.Inputs) {
				t.focusIndex = 0
			} else if t.focusIndex < 0 {
				t.focusIndex = len(t.Inputs)
			}

			cmds := make([]tea.Cmd, len(t.Inputs))
			for i := range t.Inputs {
				if i == t.focusIndex {
					cmds[i] = t.Inputs[i].Focus()
					continue
				}

				t.Inputs[i].Blur()
			}
		}
	}

	cmds := make([]tea.Cmd, len(t.Inputs))

	for i := range cmds {
		t.Inputs[i].Input, cmds[i] = t.Inputs[i].Input.Update(msg)
	}

	return t, tea.Batch(cmds...)
}

func (t *TextInputForm) View() string {
	var b strings.Builder

	for i := range t.Inputs {
		b.WriteString(t.Inputs[i].View())
		if i < len(t.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
