package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"

	"github.com/issam-assiyadi/leftmark/domain"
)

func (a *App) Render(g *gocui.Gui) error {
	if err := a.Categories.Render(g, a.CategorySelected, func(v *gocui.View, contentWidth int) error {
		rowWidth := contentWidth - len("  ")
		for i, kind := range categoryOrder {
			prefix := "  "
			if i == a.CategorySelected {
				prefix = "> "
			}
			count := len(a.itemsByCategory[kind])
			_, _ = fmt.Fprintln(v, prefix+formatCategoryRow(rowWidth, kind, count))
		}
		return nil
	}); err != nil {
		return err
	}

	items := a.currentCategoryItems()
	if err := a.Content.SetTitle(g, fmt.Sprintf(" [2] Items - %s ", a.currentCategory())); err != nil {
		return err
	}

	focus := a.ItemSelected
	if len(items) == 0 {
		focus = -1
	}

	return a.Content.Render(g, focus, func(v *gocui.View, contentWidth int) error {
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

const categoryMetaRightPad = 1

func formatCategoryRow(width int, kind domain.Kind, count int) string {
	meta := fmt.Sprintf("%d", count)
	title := string(kind)

	metaWidth := len(meta) + categoryMetaRightPad
	avail := width - metaWidth - 1
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
