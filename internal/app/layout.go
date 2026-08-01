package app

import (
	"github.com/jroimartin/gocui"
)

func (a *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if maxX <= 0 || maxY <= 0 {
		return nil
	}

	// gocui requires x0 < x1 and y0 < y1 for SetView.
	if maxX < 6 || maxY < 3 {
		return nil
	}

	sidebarWidth := max(maxX/4, 20)
	if sidebarWidth > maxX-2 {
		sidebarWidth = maxX - 2
	}
	if sidebarWidth < 2 {
		sidebarWidth = 2
	}

	sv, err := g.SetView(a.SidebarView, 0, 0, sidebarWidth-1, maxY-1)

	// err == gocui.ErrUnknownView => view doesn't exist
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		sv.Title = "Pages"
		sv.Highlight = true
		sv.SelBgColor = gocui.ColorGreen
		sv.Wrap = true
	}

	contentLeft := sidebarWidth
	contentTop := 0
	contentRight := maxX - 1
	contentBottom := maxY - 1
	if contentLeft >= contentRight || contentTop >= contentBottom {
		return nil
	}

	return a.Content.Layout(g, contentLeft, contentTop, contentRight, contentBottom)
}
