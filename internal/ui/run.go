package ui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

func (a *App) Run() error {
	// supportOverlaps=true: the preview modal is laid out on top of the
	// still-visible sidebar/tree, and this is what lets its border draw
	// correctly there instead of being redrawn over by them.
	g, err := gocui.NewGui(gocui.Output256, true)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Mouse = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan
	g.SelFrameColor = gocui.ColorCyan

	g.SetManagerFunc(func(g *gocui.Gui) error {
		// Render before layout: Render() is what updates each pane's
		// title (e.g. the "- <category>" suffix), and layout draws
		// the border/title using the view's *current* title. Drawing
		// first would render last frame's title, one step behind.
		if err := a.Render(g); err != nil {
			return err
		}
		return a.layout(g)
	})

	if err := a.BindKeys(g); err != nil {
		return err
	}

	// Views only come into existence once the manager func runs, which
	// gocui otherwise defers to the first MainLoop iteration. Run it here
	// so the categories view already exists by the time we focus it —
	// per-view keybindings never fire while g.currentView is nil.
	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("unable to set initial view:", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}
