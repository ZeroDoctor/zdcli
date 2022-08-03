package tui

import "github.com/zerodoctor/zdcli/tui/data"

type ScreenManager struct {
	vm      *ViewManager
	screens []data.IView
}
