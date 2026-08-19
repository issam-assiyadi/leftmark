// Package ui is the gocui terminal dashboard - one client of
// application.Service among others (see internal/cliadapter), with no
// knowledge of its own reachable from outside this binary.
package ui

import (
	"github.com/issam-assiyadi/leftmark/application"
	"github.com/issam-assiyadi/leftmark/domain"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/scrollview"
)

type App struct {
	Service *application.Service

	Items    []domain.Item
	Selected int

	SidebarView string
	Content     *scrollview.View

	pendingEdit *domain.Item
}

func New(svc *application.Service) (*App, error) {
	a := &App{
		Service:     svc,
		SidebarView: "sidebar",
		Content: scrollview.New(scrollview.Config{
			BaseName:       "content",
			Title:          "Detail",
			ScrollbarWidth: 3,
		}),
	}

	items, err := svc.Scan()
	if err != nil {
		return nil, err
	}
	a.setItems(items)

	return a, nil
}

// setItems replaces the displayed items and clamps the current selection to
// stay within bounds.
func (a *App) setItems(items []domain.Item) {
	a.Items = items
	if a.Selected >= len(a.Items) {
		a.Selected = len(a.Items) - 1
	}
	if a.Selected < 0 {
		a.Selected = 0
	}
}

func (a *App) selectedItem() (domain.Item, bool) {
	if a.Selected < 0 || a.Selected >= len(a.Items) {
		return domain.Item{}, false
	}
	return a.Items[a.Selected], true
}
