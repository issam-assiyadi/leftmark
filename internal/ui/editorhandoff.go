package ui

import (
	"errors"
	"log"

	"github.com/jroimartin/gocui"
)

// errOpenEditor is returned by a keybinding handler to make gocui's
// MainLoop return immediately, giving Run a clean point to close the Gui,
// shell out to $EDITOR with full terminal control, and reopen a fresh Gui
// afterward. jroimartin/gocui v0.5.0 has no suspend/resume API, so this
// close/run/reopen loop is the mechanism, not a stopgap.
var errOpenEditor = errors.New("ui: open in editor")

func (a *App) Run() error {
	for {
		g, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			return err
		}

		g.SetManagerFunc(func(g *gocui.Gui) error {
			if err := a.layout(g); err != nil {
				return err
			}
			return a.Render(g)
		})

		if err := a.BindKeys(g); err != nil {
			g.Close()
			return err
		}

		if _, err := g.SetCurrentView(a.SidebarView); err != nil {
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
