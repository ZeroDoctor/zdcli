package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelWithFocus interface {
	Focus() tea.Cmd
	Blur()
	SetModel(tea.Model)
	tea.Model
}

type Models struct {
	focusIndex int
	models     []ModelWithFocus
}

func NewModels(models ...ModelWithFocus) *Models {
	return &Models{
		focusIndex: 0,
		models:     models,
	}
}

func (m *Models) Init() tea.Cmd {
	return textinput.Blink // TODO: rethink this
}

func (m *Models) Update(msg tea.Msg) (tea.Msg, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			switch msg.Type {
			case tea.KeyEnter:
				if m.focusIndex == len(m.models) {
					return m, tea.Quit
				}
				m.focusIndex++
			case tea.KeyUp:
				m.focusIndex--
			case tea.KeyDown:

				m.focusIndex++
			}

			if m.focusIndex > len(m.models) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.models)
			}

			cmds := make([]tea.Cmd, len(m.models))
			for i := 0; i < len(m.models); i++ {
				if i == m.focusIndex {
					cmds[i] = m.models[i].Focus()
					continue
				}

				m.models[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmds := make([]tea.Cmd, len(m.models))

	for i := range m.models {
		var model tea.Model
		model, cmds[i] = m.models[i].Update(msg)
		m.models[i].SetModel(model)
	}

	return m, tea.Batch(cmds...)
}

func (m *Models) View() string {
	var b strings.Builder

	for i := range m.models {
		b.WriteString(m.models[i].View())
		if i < len(m.models)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
