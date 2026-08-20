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

	categoryKeys := []struct {
		key interface{}
		fn  func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyArrowDown, a.categoryMoveDown},
		{'j', a.categoryMoveDown},
		{gocui.KeyArrowUp, a.categoryMoveUp},
		{'k', a.categoryMoveUp},
		{gocui.MouseLeft, a.focusCategories},
	}
	for _, b := range categoryKeys {
		if err := g.SetKeybinding(a.CategoriesView, b.key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}

	itemKeys := []struct {
		key interface{}
		fn  func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyArrowDown, a.itemMoveDown},
		{'j', a.itemMoveDown},
		{gocui.KeyArrowUp, a.itemMoveUp},
		{'k', a.itemMoveUp},
		{gocui.KeyEnter, a.openSelectedInEditor},
		{'o', a.openSelectedInEditor},
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
	a.focused = a.CategoriesView
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
		a.ItemSelected = 0
	}
	return nil
}

func (a *App) categoryMoveUp(g *gocui.Gui, v *gocui.View) error {
	if a.CategorySelected > 0 {
		a.CategorySelected--
		a.ItemSelected = 0
	}
	return nil
}

func (a *App) itemMoveDown(g *gocui.Gui, v *gocui.View) error {
	if a.ItemSelected < len(a.currentCategoryItems())-1 {
		a.ItemSelected++
	}
	return nil
}

func (a *App) itemMoveUp(g *gocui.Gui, v *gocui.View) error {
	if a.ItemSelected > 0 {
		a.ItemSelected--
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
