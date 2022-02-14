package tui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type Screen struct {
	g       *gocui.Gui
	msgChan chan interface{}
}

func NewScreen(g *gocui.Gui) *Screen {
	s := &Screen{
		g:       g,
		msgChan: make(chan interface{}, 10000),
	}
	return s
}

func (s Screen) Name() string               { return "screen" }
func (s *Screen) Channel() chan interface{} { return s.msgChan }
func (s *Screen) Send(msg Data)             { s.msgChan <- msg }

func (s *Screen) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(s.Name(), 0, (maxY/15)+1, maxX-1, (maxY-(maxY/15))-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = s.Name()
		v.Wrap = true
	}
	return nil
}

func (s *Screen) PrintView() {
	for msg := range s.msgChan {
		var str string
		m := msg.(Data)

		switch m.Type {
		case "msg":
			str = m.Msg.(string)
		}

		s.Display(str)
	}
}

func (s *Screen) Display(msg string) {
	s.g.UpdateAsync(func(g *gocui.Gui) error {
		v, err := g.View(s.Name())
		if err != nil {
			return err
		}

		fmt.Fprint(v, msg)

		ox, _ := v.Origin()
		_, sy := v.Size()

		y := (v.ViewLinesHeight() - sy) - 1
		if y < 0 {
			y = 0
		}

		err = v.SetOrigin(ox, y)
		if err != nil {
			return err
		}

		return nil
	})
}
