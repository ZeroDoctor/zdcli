package ui

import (
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	PADDING   = 2
	MAX_WIDTH = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type Tick float64

type Progress struct {
	progress progress.Model
	wg       sync.WaitGroup
}

func NewProgress(work func(p *Progress) any) *Progress {
	p := &Progress{}

	p.wg.Add(1)
	go func() {
		work(p)
		p.Update(Tick(1.0 - p.progress.Percent()))
		p.wg.Done()
	}()

	return p
}

func (p *Progress) Init() tea.Cmd { return nil }

func (p *Progress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.progress.Width = msg.Width - PADDING*2 - 4
		if p.progress.Width > MAX_WIDTH {
			p.progress.Width = MAX_WIDTH
		}

		return p, nil
	case Tick:
		if p.progress.Percent() >= 1.0 {
			return p, tea.Quit
		}
		cmd := p.progress.IncrPercent(float64(msg))

		return p, tea.Batch(cmd)
	case progress.FrameMsg:
		pm, cmd := p.progress.Update(msg)
		p.progress = pm.(progress.Model)

		return p, cmd
	}

	return p, nil
}

func (p *Progress) View() string {
	pad := strings.Repeat(" ", PADDING)
	return "\n" +
		pad + p.progress.View() + "\n\n" +
		pad + helpStyle("working...")
}

func (p *Progress) Wait() { p.wg.Wait() }
