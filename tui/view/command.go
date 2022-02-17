package view

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/zerodoctor/zdcli/tui"
)

type Command struct {
	g       *gocui.Gui
	msgChan chan interface{}

	running bool
	tui.CommandState
}

func NewCommand(g *gocui.Gui, state tui.CommandState) *Command {
	c := &Command{
		g:            g,
		msgChan:      make(chan interface{}, 100),
		CommandState: state,
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
		buf := v.Buffer()
		if !c.running {
			c.running = true
			go func(buf string) {
				c.StdIn = make(chan string, 100)

				c.VM.SendView("screen", Data{
					Type: "msg", Msg: fmt.Sprintf("[zd] starting fork [cmd=%s]\n", buf),
				})

				if err := c.Fork(c.VM, c.StdIn, buf); err != nil {
					c.VM.SendView("screen", Data{
						Type: "msg", Msg: "[zd] " + err.Error(),
					})
				}

				close(c.StdIn)
				c.running = false
				c.VM.SendView("screen", Data{
					Type: "msg", Msg: "[zd] fork exited\n\n",
				})
			}(buf)
		} else {
			c.StdIn <- buf
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
