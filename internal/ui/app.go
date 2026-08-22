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

	// collapsed tracks which directory/file paths (see row.path) are
	// collapsed in the content pane's tree. Absent key = expanded, so new
	// paths default to expanded and collapse choices survive rescan/setItems
	// without any seeding or reset logic.
	collapsed map[string]bool

	CategorySelected int
	RowSelected      int

	Categories *scrollview.View
	Content    *scrollview.View

	focused string

	pendingEdit *domain.Item
}

func New(svc *application.Service) (*App, error) {
	a := &App{
		Service:   svc,
		collapsed: make(map[string]bool),
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

	n := len(a.visibleRows())
	if a.RowSelected >= n {
		a.RowSelected = n - 1
	}
	if a.RowSelected < 0 {
		a.RowSelected = 0
	}
}

func (a *App) visibleRows() []row {
	return buildRows(a.currentCategoryItems(), a.collapsed)
}

func (a *App) selectedRow() (row, bool) {
	rows := a.visibleRows()
	if a.RowSelected < 0 || a.RowSelected >= len(rows) {
		return row{}, false
	}
	return rows[a.RowSelected], true
}

func (a *App) currentCategory() domain.Kind {
	return categoryOrder[a.CategorySelected]
}

func (a *App) currentCategoryItems() []domain.Item {
	return a.itemsByCategory[a.currentCategory()]
}

func (a *App) selectedItem() (domain.Item, bool) {
	r, ok := a.selectedRow()
	if !ok || r.kind != rowKindItem {
		return domain.Item{}, false
	}
	return r.item, true
}
