package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"

	"github.com/issam-assiyadi/leftmark/domain"
)

func (a *App) Render(g *gocui.Gui) error {
	cv, err := g.View(a.CategoriesView)
	if err != nil {
		return err
	}
	cv.Clear()

	cvWidth, _ := cv.Size()
	rowWidth := cvWidth - len("  ") // prefix is always 2 columns, even when unselected
	for i, kind := range categoryOrder {
		prefix := "  "
		if i == a.CategorySelected {
			prefix = "> "
		}
		count := len(a.itemsByCategory[kind])
		_, _ = fmt.Fprintln(cv, prefix+formatCategoryRow(rowWidth, kind, count))
	}

	items := a.currentCategoryItems()
	if err := a.Content.SetTitle(g, fmt.Sprintf(" [2] Items - %s ", a.currentCategory())); err != nil {
		return err
	}

	return a.Content.Render(g, func(v *gocui.View, contentWidth int) error {
		if len(items) == 0 {
			_, _ = fmt.Fprintln(v, "No items in this category")
			return nil
		}
		for i, item := range items {
			prefix := "  "
			if i == a.ItemSelected {
				prefix = "> "
			}
			_, _ = fmt.Fprintln(v, prefix+renderItemRow(item))
		}
		return nil
	})
}

func renderItemRow(item domain.Item) string {
	return fmt.Sprintf("%s:%d %s", item.File, item.Line, item.Text)
}

// categoryMetaRightPad keeps the item count from sitting flush against the
// view's right border.
const categoryMetaRightPad = 1

// formatCategoryRow lays out a category name against its item count so the
// count always lands near the row's right edge, truncating the name with
// "..." rather than letting a long name push the count off-screen or wrap it
// onto its own line.
func formatCategoryRow(width int, kind domain.Kind, count int) string {
	meta := fmt.Sprintf("%d", count)
	title := string(kind)

	metaWidth := len(meta) + categoryMetaRightPad
	avail := width - metaWidth - 1 // 1 column gap between title and meta
	if avail < 0 {
		avail = 0
	}
	if len(title) > avail {
		if avail <= 3 {
			title = title[:avail]
		} else {
			title = title[:avail-3] + "..."
		}
	}

	gap := width - len(title) - metaWidth
	if gap < 1 {
		gap = 1
	}
	return title + strings.Repeat(" ", gap) + meta + strings.Repeat(" ", categoryMetaRightPad)
}
