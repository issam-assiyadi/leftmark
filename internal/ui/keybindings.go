package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

func (a *App) BindKeys(g *gocui.Gui) error {
	global := []struct {
		key interface{}
		fn  func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyCtrlC, quit},
		{'q', quit},
		{'r', a.rescan},
		{'1', a.focusCategories},
		{'2', a.focusItems},
	}
	for _, b := range global {
		if err := g.SetKeybinding("", b.key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}

	for _, viewname := range []string{a.Categories.WrapperName(), a.Categories.ContentViewName()} {
		if err := g.SetKeybinding(viewname, gocui.MouseLeft, gocui.ModNone, a.focusCategories); err != nil {
			return err
		}
	}

	categoryKeys := []struct {
		key interface{}
		fn  func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyArrowDown, a.categoryMoveDown},
		{'j', a.categoryMoveDown},
		{gocui.KeyArrowUp, a.categoryMoveUp},
		{'k', a.categoryMoveUp},
	}
	for _, b := range categoryKeys {
		if err := g.SetKeybinding(a.Categories.WrapperName(), b.key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}

	itemKeys := []struct {
		key interface{}
		fn  func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyArrowDown, a.rowMoveDown},
		{'j', a.rowMoveDown},
		{gocui.KeyArrowUp, a.rowMoveUp},
		{'k', a.rowMoveUp},
		{gocui.KeyEnter, a.activateSelectedRow},
		{'o', a.activateSelectedRow},
		{gocui.KeySpace, a.activateSelectedRow},
	}
	for _, viewname := range []string{a.Content.WrapperName(), a.Content.ContentViewName(), a.Content.ScrollViewName()} {
		if err := g.SetKeybinding(viewname, gocui.MouseLeft, gocui.ModNone, a.focusItems); err != nil {
			return err
		}
	}
	for _, b := range itemKeys {
		if err := g.SetKeybinding(a.Content.WrapperName(), b.key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (a *App) focusCategories(g *gocui.Gui, v *gocui.View) error {
	a.focused = a.Categories.WrapperName()
	_, err := g.SetCurrentView(a.focused)
	return err
}

func (a *App) focusItems(g *gocui.Gui, v *gocui.View) error {
	a.focused = a.Content.WrapperName()
	_, err := g.SetCurrentView(a.focused)
	return err
}

func (a *App) categoryMoveDown(g *gocui.Gui, v *gocui.View) error {
	if a.CategorySelected < len(categoryOrder)-1 {
		a.CategorySelected++
		a.RowSelected = 0
	}
	return nil
}

func (a *App) categoryMoveUp(g *gocui.Gui, v *gocui.View) error {
	if a.CategorySelected > 0 {
		a.CategorySelected--
		a.RowSelected = 0
	}
	return nil
}

func (a *App) rowMoveDown(g *gocui.Gui, v *gocui.View) error {
	if a.RowSelected < len(a.visibleRows())-1 {
		a.RowSelected++
	}
	return nil
}

func (a *App) rowMoveUp(g *gocui.Gui, v *gocui.View) error {
	if a.RowSelected > 0 {
		a.RowSelected--
	}
	return nil
}

func (a *App) activateSelectedRow(g *gocui.Gui, v *gocui.View) error {
	r, ok := a.selectedRow()
	if !ok {
		return nil
	}
	if r.kind == rowKindItem {
		return a.openSelectedInEditor(g, v)
	}
	a.collapsed[r.path] = !a.collapsed[r.path]
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
