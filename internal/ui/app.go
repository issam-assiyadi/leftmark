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
	Filter   domain.Filter

	SidebarView string
	Content     *scrollview.View

	pending     *pendingResolve
	pendingEdit *domain.Item
}

// pendingResolve tracks a done-or-delete confirmation waiting on a y/n
// keypress; it blocks nothing else in the UI, it's just state read by
// Render and the y/n/esc keybindings.
type pendingResolve struct {
	item domain.Item
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

	if _, err := svc.Scan(); err != nil {
		return nil, err
	}
	a.refresh()

	return a, nil
}

// refresh re-reads the service's current snapshot through the active
// filter. Call after any Scan/PromoteToDoing/ResolveDone, all of which
// already update the service's own snapshot.
func (a *App) refresh() {
	a.Items = a.Service.List(a.Filter)
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
