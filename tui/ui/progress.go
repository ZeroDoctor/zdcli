package ui

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const (
	PADDING   = 2
	MAX_WIDTH = 80
)

const (
	GRAY_COL  = "#D7D7D7"
	GREEN_COL = "#ADFDAD"
	BLUE_COL  = "#AADDFF"
	RED_COL   = "#FF77AA"
)

var (
	defStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(GRAY_COL)).Render
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(RED_COL)).Render
)

var EOF_ProgErr error = errors.New("EOF")

type DefTick float64
type DefProgErr error

type Progress struct {
	progress progress.Model
	err      error
}

func NewProgress(work func() (*tea.Program, error)) *Progress {
	prog := &Progress{
		progress: progress.New(
			progress.WithScaledGradient(GREEN_COL, BLUE_COL),
			progress.WithColorProfile(termenv.ANSI256),
		),
	}

	go func() {
		var p *tea.Program
		var err error

		if p, err = work(); err != nil {
			p.Send(DefProgErr(err))
			return
		}

		p.Send(DefProgErr(EOF_ProgErr))
	}()

	return prog
}

func (p *Progress) Init() tea.Cmd { return nil }

func (p *Progress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	tea.Printf("update was called...[msg_type=%T]\n", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return p, tea.Quit

		}

	case tea.WindowSizeMsg:
		p.progress.Width = msg.Width - PADDING*2 - 4
		if p.progress.Width > MAX_WIDTH {
			p.progress.Width = MAX_WIDTH
		}

		return p, nil

	case DefProgErr:
		if msg == EOF_ProgErr {
			var cmds []tea.Cmd

			cmds = append(cmds,
				tea.Sequentially(
					tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg { // pause a bit before quiting
						return nil
					}),
					tea.Quit,
				),
			)

			cmds = append(cmds, p.progress.IncrPercent(float64(1.0)))

			return p, tea.Batch(cmds...)
		}

		p.err = msg

		return p, tea.Sequentially(
			tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg { // pause a bit before quiting
				return nil
			}),
			tea.Quit,
		)

	case DefTick:
		var cmds []tea.Cmd

		if msg >= 1.0 {
			cmds = append(cmds,
				tea.Sequentially(
					tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg { // pause a bit before quiting
						return nil
					}),
					tea.Quit,
				),
			)
		}

		cmds = append(cmds, p.progress.IncrPercent(float64(msg)))

		return p, tea.Batch(cmds...)

	case progress.FrameMsg:
		pm, cmd := p.progress.Update(msg)
		p.progress = pm.(progress.Model)

		return p, cmd
	}

	return p, nil
}

func (p *Progress) View() string {
	fn := func() string { return defStyle("working...") }
	if p.err != nil {
		fn = func() string { return errStyle(p.err.Error()) }
	}

	pad := strings.Repeat(" ", PADDING)
	return pad + fn() + "\n" +
		pad + p.progress.View() + "\n"
}
