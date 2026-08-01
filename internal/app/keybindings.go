package app

import (
	"github.com/jroimartin/gocui"
)

func (a *App) BindKeys(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return a.Content.ScrollDown(g)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return a.Content.ScrollUp(g)
	}); err != nil {
		return err
	}

	return nil
}
