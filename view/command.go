package view

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

type Command struct {
	g       *gocui.Gui
	msgChan chan interface{}

	fn func(v *gocui.View) error
}

func NewCommand(g *gocui.Gui, fn func(*gocui.View) error) *Command {
	c := &Command{
		g:       g,
		fn:      fn,
		msgChan: make(chan interface{}, 100),
	}
	return c
}

func (c Command) Name() string               { return "command" }
func (c *Command) Channel() chan interface{} { return c.msgChan }
func (c *Command) Send(msg Data)             { c.msgChan <- msg }

func (c *Command) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(c.Name(), 0, (maxY-(maxY/15))-1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = c.Name()
		v.Wrap = true
		v.Editable = true
		v.Editor = c
	}
	return nil
}

func (c *Command) PrintView() {
	for msg := range c.msgChan {
		var str string
		m := msg.(Data)

		switch m.Type {
		}

		c.Display(str)
	}
}

func (c *Command) Display(msg string) {
	c.g.UpdateAsync(func(g *gocui.Gui) error {
		v, err := g.View(c.Name())
		if err != nil {
			return err
		}

		v.Clear()
		fmt.Fprint(v, msg)
		return nil
	})
}

func (c *Command) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// TODO: Ctrl-Backspace
	// TODO: Ctrl-Arrow_Keys
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyTab:

	case key == gocui.KeyEnter:
		err := c.fn(v)
		if err != nil {
			log.Fatal(err)
		}

		v.SetCursor(0, 0)
		v.Clear()
	case key == gocui.KeyArrowDown:
	case key == gocui.KeyArrowUp:
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0)
	}
}
