package ui

import (
	"errors"
	"log"

	"github.com/jroimartin/gocui"
)

var errOpenEditor = errors.New("ui: open in editor")

func (a *App) Run() error {
	for {
		g, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			return err
		}

		g.Mouse = true
		g.Highlight = true
		g.SelFgColor = gocui.ColorCyan

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
			g.Close()
			return err
		}

		// Views only come into existence once the manager func runs, which
		// gocui otherwise defers to the first MainLoop iteration. Run it here
		// so the categories view already exists by the time we focus it —
		// per-view keybindings never fire while g.currentView is nil.
		if err := a.layout(g); err != nil {
			g.Close()
			return err
		}
		if _, err := g.SetCurrentView(a.focused); err != nil {
			log.Println("unable to set initial view:", err)
		}

		runErr := g.MainLoop()
		g.Close()

		if runErr == nil || errors.Is(runErr, gocui.ErrQuit) {
			return nil
		}

		if errors.Is(runErr, errOpenEditor) {
			item := a.pendingEdit
			a.pendingEdit = nil
			if item != nil {
				if err := a.Service.OpenInEditor(item.File, item.Line); err != nil {
					log.Println("open in editor:", err)
				}
			}
			continue
		}

		return runErr
	}
}
