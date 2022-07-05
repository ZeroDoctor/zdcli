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

	if err := km.SetKey("", gocui.KeyEsc, gocui.ModNone, vm.Quit); err != nil {
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

	// cx, cy := v.Cursor()
	// xOff := 0
	// yOff := 1
	//
	// sx, _ := v.Size()
	//
	// lineLength := len(v.Buffer())
	// maxRow := int(lineLength / sx)
	// if maxRow > 0 {
	// 	currRow := int(cx / sx)
	// 	if currRow-1 > 0 {
	// 		yOff = 0
	// 		xOff = sx
	// 	}
	// }
	//
	// cx, cy = checkVertCursor(cx-xOff, cy-yOff)
	// if err := v.SetCursor(cx, cy); err != nil {
	// 	return err
	// }

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

	// xOff := 0
	// yOff := 1
	//
	// sx, sy := v.Size() // screen width and height
	// cx, cy := v.Cursor()
	//
	// lineLength := len(v.Buffer())
	// maxRow := int(lineLength / sx)
	// if maxRow > 0 {
	// 	currRow := int(cx / sx)
	// 	if currRow+1 < maxRow {
	// 		yOff = 0
	// 		xOff = sx
	// 	}
	// }
	//
	// if err := v.SetCursor(cx+xOff, cy+yOff); err != nil {
	// 	return err
	// }
	//
	// trueY := (cx / sx) + cy + 2 // add 2 here because the screen starts at 0 instead of 1
	// if trueY > sy {
	ox, oy := v.Origin()
	if err := v.SetOrigin(ox, oy+1); err != nil {
		return err
	}
	// }

	return nil
}

// func checkVertCursor(cx, cy int) (int, int) {
// 	if cx < 0 {
// 		cx = 0
// 	}
//
// 	if cy < 0 {
// 		cy = 0
// 	}
//
// 	return cx, cy
// }

// TODO: LeftScreen and RightScreen
