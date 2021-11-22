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
// figure out how to set correct cursor position with wrapped text

func UpScreen(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	return nil
}

func DownScreen(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	cx, cy := v.Cursor()

	line, err := v.Line(cy)
	if err != nil {
		return err
	}

	sx, sy := v.Size()

	if len(line) > sx && cx < sx {
		return v.SetCursor(len(line)-(len(line)-sx), cy)
	}

	if cx >= len(line) {
		cx = len(line) - 1
	}

	ox, oy := v.Origin()
	if err = v.SetCursor(cx, cy+1); err == nil && cy+1 >= sy {
		if err = v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	return nil
}

// TODO: LeftScreen and RightScreen
