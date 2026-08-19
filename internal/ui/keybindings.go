package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

func (a *App) BindKeys(g *gocui.Gui) error {
	bindings := []struct {
		key interface{}
		fn  func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyCtrlC, quit},
		{'q', quit},
		{gocui.KeyArrowDown, a.moveDown},
		{'j', a.moveDown},
		{gocui.KeyArrowUp, a.moveUp},
		{'k', a.moveUp},
		{gocui.KeyEnter, a.openSelectedInEditor},
		{'o', a.openSelectedInEditor},
		{'r', a.rescan},
	}

	for _, b := range bindings {
		if err := g.SetKeybinding("", b.key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (a *App) moveDown(g *gocui.Gui, v *gocui.View) error {
	if a.Selected < len(a.Items)-1 {
		a.Selected++
	}
	return nil
}

func (a *App) moveUp(g *gocui.Gui, v *gocui.View) error {
	if a.Selected > 0 {
		a.Selected--
	}
	return nil
}

func (a *App) openSelectedInEditor(g *gocui.Gui, v *gocui.View) error {
	item, ok := a.selectedItem()
	if !ok {
		return nil
	}
	a.pendingEdit = &item
	return errOpenEditor
}

func (a *App) rescan(g *gocui.Gui, v *gocui.View) error {
	items, err := a.Service.Scan()
	if err != nil {
		log.Println("scan:", err)
		return nil
	}
	a.setItems(items)
	return nil
}
