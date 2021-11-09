package view

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type Header struct {
	g       *gocui.Gui
	msgChan chan interface{}

	msg     string
	permMsg string
}

func NewHeader(g *gocui.Gui) *Header {
	h := &Header{
		g:       g,
		msgChan: make(chan interface{}, 25),
	}
	return h
}

func (h Header) Name() string               { return "header" }
func (h *Header) Channel() chan interface{} { return h.msgChan }
func (h *Header) Send(msg Data)             { h.msgChan <- msg }

func (h *Header) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(h.Name(), 0, 0, maxX-1, (maxY / 15), 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = h.Name()
		v.Wrap = false
	}
	return nil
}

func (h *Header) PrintView() {
	for msg := range h.msgChan {
		var str string
		m := msg.(Data)

		switch m.Type {
		case "clock":
			h.permMsg = m.Msg.(string) + "|"
		case "msg":
			h.msg = m.Msg.(string) + "|"
		case "temp":
			str = m.Msg.(string)
		}

		h.Display(h.permMsg + h.msg + str)
	}
}

func (h *Header) Display(msg string) {
	h.g.UpdateAsync(func(g *gocui.Gui) error {
		v, err := g.View(h.Name())
		if err != nil {
			return err
		}

		v.Clear()
		fmt.Fprint(v, msg)
		return nil
	})
}
