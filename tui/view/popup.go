package view

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type Popup struct {
	g       *gocui.Gui
	name    string
	msgChan chan interface{}
}

func NewPopup(g *gocui.Gui, name string) *Popup {
	p := &Popup{
		g:       g,
		name:    name + "_pop",
		msgChan: make(chan interface{}, 10),
	}

	return p
}

func (p Popup) Width() int  { return 0 }
func (p Popup) Height() int { return 0 }

func (p *Popup) Layout(*gocui.Gui) error {
	return nil
}

func (p *Popup) PrintView() {

}

func (p *Popup) Display(msg string) {
	p.g.UpdateAsync(func(g *gocui.Gui) error {
		v, err := g.View(p.Name())
		if err != nil {
			return err
		}

		v.Clear()
		fmt.Fprint(v, msg)
		return nil
	})
}

func (p *Popup) Name() string { return p.name }

func (p *Popup) Channel() chan interface{} { return p.msgChan }
