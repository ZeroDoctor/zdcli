package view

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

		out := msg
		cols, _ := v.Size()
		curr := v.LinesHeight()
		if curr <= 0 { // means this is the first line
			fmt.Fprint(v, out)
			return nil
		}

		line, err := v.Line(curr - 1)
		if err != nil {
			out = "[ERROR-1] failed to get line on screen view " + err.Error()
			fmt.Fprintln(v, out)
			return nil
		}

		// if current line has wrapped due to termainl resize (for example)
		// then create new line to write the current msg out
		if len(line) > cols {
			fmt.Fprint(v, "\n")

			cols, _ = v.Size()
			curr = v.LinesHeight()
			line, err = v.Line(curr - 1)
			if err != nil {
				out = "[ERROR-2] failed to get line on screen view " + err.Error()
				fmt.Fprintln(v, out)
				return nil
			}
		}

		// cols - (len(line) % cols) bascially means account for current line
		// in max number of columns left (cols - len(line))
		// and account for the line wrapping or (len(line) % cols)
		for i := len(msg); i > (cols - (len(line) % cols)); i -= cols {
			out = msg[:(cols-(len(line)%cols))] + "\n" + msg[(cols-(len(line)%cols)):]
		}

		// line := v.ViewLinesHeight()
		// _, rows := v.Size()
		// if line > rows {
		// 	ox, oy := v.Origin()
		// 	v.SetOrigin(ox, oy+1)
		// }

		fmt.Fprint(v, out)
		return nil
	})
}
