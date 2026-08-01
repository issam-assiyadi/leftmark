package app

import (
	"fmt"

	"github.com/issam-assiyadi/leftmark/internal/markdown"
	"github.com/jroimartin/gocui"
)

func (a *App) Render(g *gocui.Gui) error {
	// render the sidebar content
	sv, err := g.View(a.SidebarView)
	if err != nil {
		return err
	}
	sv.Clear()

	for i, p := range a.Pages {
		prefix := "  "
		if i == a.CurrentPage {
			prefix = "> "
		}

		_, _ = fmt.Fprintf(sv, "%s%s\n", prefix, p.Title)
	}

	return a.Content.Render(g, func(v *gocui.View, contentWidth int) error {
		if len(a.Pages) == 0 {
			_, _ = fmt.Fprintln(v, "No pages")
			return nil
		}

		if contentWidth <= 0 {
			contentWidth = 80
		}

		md := a.Pages[a.CurrentPage].Content
		// FIXME: I should investigate more in that, it seems like we do this in a wrong way.
		rendered, err := markdown.Render(md, contentWidth)
		if err != nil {
			_, _ = fmt.Fprint(v, md)
			return nil
		}

		_, _ = fmt.Fprint(v, rendered)
		return nil
	})
}
