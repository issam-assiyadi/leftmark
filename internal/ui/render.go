package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/issam-assiyadi/leftmark/domain"
)

func (a *App) Render(g *gocui.Gui) error {
	sv, err := g.View(a.SidebarView)
	if err != nil {
		return err
	}
	sv.Clear()

	if len(a.Items) == 0 {
		_, _ = fmt.Fprintln(sv, "No items")
	}
	for i, item := range a.Items {
		prefix := "  "
		if i == a.Selected {
			prefix = "> "
		}
		_, _ = fmt.Fprintln(sv, prefix+renderRow(item))
	}

	return a.Content.Render(g, func(v *gocui.View, contentWidth int) error {
		item, ok := a.selectedItem()
		if !ok {
			_, _ = fmt.Fprintln(v, "No item selected")
			return nil
		}

		_, _ = fmt.Fprintf(v, "%s\n", item.Kind)
		_, _ = fmt.Fprintf(v, "%s:%d\n\n", item.File, item.Line)
		_, _ = fmt.Fprintln(v, item.Text)
		return nil
	})
}

func renderRow(item domain.Item) string {
	return fmt.Sprintf("%-9s %s:%d %s", item.Kind, item.File, item.Line, item.Text)
}
