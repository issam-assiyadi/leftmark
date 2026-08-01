package ui

import (
	"log"

	"github.com/jroimartin/gocui"

	"github.com/issam-assiyadi/leftmark/domain"
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
		{'d', a.promoteSelected},
		{'x', a.beginResolve},
		{'y', a.confirmResolve(true)},
		{'n', a.confirmResolve(false)},
		{gocui.KeyEsc, a.cancelResolve},
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
	if !ok || !item.Located {
		return nil
	}
	a.pendingEdit = &item
	return errOpenEditor
}

// promoteSelected tags and tracks an open item as doing. No confirmation -
// unlike resolving to done, there's no destructive choice to make here.
func (a *App) promoteSelected(g *gocui.Gui, v *gocui.View) error {
	item, ok := a.selectedItem()
	if !ok || item.Status != domain.StatusOpen {
		return nil
	}
	if _, err := a.Service.PromoteToDoing(item.File, item.Line); err != nil {
		log.Println("promote:", err)
		return nil
	}
	a.refresh()
	return nil
}

func (a *App) beginResolve(g *gocui.Gui, v *gocui.View) error {
	item, ok := a.selectedItem()
	if !ok || item.Status != domain.StatusDoing {
		return nil
	}
	a.pending = &pendingResolve{item: item}
	return nil
}

func (a *App) confirmResolve(deleteComment bool) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if a.pending == nil {
			return nil
		}
		id := a.pending.item.ID
		a.pending = nil

		if _, err := a.Service.ResolveDone(id, !deleteComment); err != nil {
			log.Println("resolve:", err)
			return nil
		}
		a.refresh()
		return nil
	}
}

func (a *App) cancelResolve(g *gocui.Gui, v *gocui.View) error {
	a.pending = nil
	return nil
}

func (a *App) rescan(g *gocui.Gui, v *gocui.View) error {
	if _, err := a.Service.Scan(); err != nil {
		log.Println("scan:", err)
		return nil
	}
	a.refresh()
	return nil
}
