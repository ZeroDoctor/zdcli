package tui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

type KeyManager struct {
	g *gocui.Gui
}

func NewKeyManager(g *gocui.Gui, vm *ViewManager) *KeyManager {
	km := &KeyManager{g: g}
	if err := km.SetKey("", gocui.KeyCtrlC, gocui.ModNone, vm.Quit); err != nil {
		log.Fatal(err)
	}
	if err := km.SetKey("", gocui.KeyTab, gocui.ModNone, vm.NextView); err != nil {
		log.Fatal(err)
	}
	return km
}

func (km *KeyManager) SetKey(viewname string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) error {
	return km.g.SetKeybinding(viewname, key, mod, handler)
}

// TODO: improve up and down movement in context of text wrapping!

func UpScreen(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	cx, cy := v.Cursor()
	if cy > 0 {
		if err := v.SetCursor(cx, cy-1); err != nil {
			return err
		}
	}

	ox, oy := v.Origin()
	if oy > 0 {
		var w int
		if w, _ = HasWrapped(cy, v); w > 0 {
			w -= 1
		}
		if err := v.SetOrigin(ox, oy-1-w); err != nil {
			return err
		}
	}

	return nil
}

func DownScreen(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	lines := v.ViewLinesHeight()
	cx, cy := v.Cursor()
	if cy < lines-1 {
		if err := v.SetCursor(cx, cy+1); err != nil {
			return err
		}

		_, y := v.Size()
		if cy+1 > y-1 && cy+1 < lines-1 {
			var w int
			if w, _ = HasWrapped(cy, v); w > 0 {
				w -= 1
			}
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1+w); err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO: LeftScreen and RightScreen

// # helper function

func HasWrapped(y int, v *gocui.View) (int, error) {
	_, maxY := v.Size()
	line, err := v.Line(y)
	return (len(line) / (maxY + 1)), err
}
