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
		if a.pending != nil {
			_, _ = fmt.Fprintf(v, "Mark %q done.\n\nDelete the tagged comment from the source?\n\n  [y] delete it\n  [n] keep it, just mark done\n  [esc] cancel\n", a.pending.item.Text)
			return nil
		}

		item, ok := a.selectedItem()
		if !ok {
			_, _ = fmt.Fprintln(v, "No item selected")
			return nil
		}

		_, _ = fmt.Fprintf(v, "%s  %s\n", item.Kind, item.Status)
		_, _ = fmt.Fprintf(v, "%s:%d\n\n", item.File, item.Line)
		_, _ = fmt.Fprintln(v, item.Text)
		if item.ID != "" && !item.Located {
			_, _ = fmt.Fprintln(v, "\nnot found in the current tree (branch switch, or deleted by hand)")
		}
		return nil
	})
}

func renderRow(item domain.Item) string {
	status := "[" + string(item.Status) + "]"
	if item.ID != "" && !item.Located {
		status = "[missing]"
	}
	return fmt.Sprintf("%-9s %-9s %s:%d %s", item.Kind, status, item.File, item.Line, item.Text)
}
