package tui

import "github.com/awesome-gocui/gocui"

type KeyManager struct {
	g *gocui.Gui
}

func NewKeyManager(g *gocui.Gui, vm *ViewManager) *KeyManager {
	km := &KeyManager{g: g}
	if err := km.SetKey("", gocui.KeyCtrlC, gocui.ModNone, vm.Quit); err != nil {
		panic(err)
	}
	if err := km.SetKey("", gocui.KeyTab, gocui.ModNone, vm.NextView); err != nil {
		panic(err)
	}
	return km
}

func (km *KeyManager) SetKey(viewname string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) error {
	return km.g.SetKeybinding(viewname, key, mod, handler)
}

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
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}

	return nil
}

func DownScreen(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		return err
	}

	_, rows := v.Size()
	if cy > rows {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	return nil
}

// TODO: LeftScreen and RightScreen
