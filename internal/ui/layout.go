package ui

import "github.com/awesome-gocui/gocui"

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

	if err := a.Categories.Layout(g, 0, 0, sidebarWidth-1, maxY-1); err != nil {
		return err
	}

	contentLeft := sidebarWidth
	contentTop := 0
	contentRight := maxX - 1
	contentBottom := maxY - 1
	if contentLeft < contentRight && contentTop < contentBottom {
		if err := a.Content.Layout(g, contentLeft, contentTop, contentRight, contentBottom); err != nil {
			return err
		}
	}

	return a.Preview.Layout(g, maxX, maxY)
}
