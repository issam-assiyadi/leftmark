// Package ui is the gocui terminal dashboard - one client of
// application.Service among others (see internal/cliadapter), with no
// knowledge of its own reachable from outside this binary.
package ui

import (
	"github.com/issam-assiyadi/leftmark/application"
	"github.com/issam-assiyadi/leftmark/domain"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/scrollview"
)

var categoryOrder = []domain.Kind{domain.KindTODO, domain.KindFIXME, domain.KindNOTE, domain.KindQUESTION}

type App struct {
	Service *application.Service

	Items           []domain.Item
	itemsByCategory map[domain.Kind][]domain.Item

	CategorySelected int
	ItemSelected     int

	Categories *scrollview.View
	Content    *scrollview.View

	focused string

	pendingEdit *domain.Item
}

func New(svc *application.Service) (*App, error) {
	a := &App{
		Service: svc,
		Categories: scrollview.New(scrollview.Config{
			BaseName:      "categories",
			Title:         " [1] Categories ",
			HideScrollbar: true,
		}),
		Content: scrollview.New(scrollview.Config{
			BaseName:       "content",
			Title:          " [2] Items ",
			ScrollbarWidth: 3,
		}),
	}
	a.focused = a.Categories.WrapperName()

	items, err := svc.Scan()
	if err != nil {
		return nil, err
	}
	a.setItems(items)

	return a, nil
}

func (a *App) setItems(items []domain.Item) {
	a.Items = items

	a.itemsByCategory = make(map[domain.Kind][]domain.Item, len(categoryOrder))
	for _, item := range items {
		a.itemsByCategory[item.Kind] = append(a.itemsByCategory[item.Kind], item)
	}

	if a.CategorySelected < 0 {
		a.CategorySelected = 0
	}
	if a.CategorySelected >= len(categoryOrder) {
		a.CategorySelected = len(categoryOrder) - 1
	}

	n := len(a.currentCategoryItems())
	if a.ItemSelected >= n {
		a.ItemSelected = n - 1
	}
	if a.ItemSelected < 0 {
		a.ItemSelected = 0
	}
}

func (a *App) currentCategory() domain.Kind {
	return categoryOrder[a.CategorySelected]
}

func (a *App) currentCategoryItems() []domain.Item {
	return a.itemsByCategory[a.currentCategory()]
}

func (a *App) selectedItem() (domain.Item, bool) {
	items := a.currentCategoryItems()
	if a.ItemSelected < 0 || a.ItemSelected >= len(items) {
		return domain.Item{}, false
	}
	return items[a.ItemSelected], true
}
