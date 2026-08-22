package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/mattn/go-runewidth"

	"github.com/issam-assiyadi/leftmark/domain"
)

const rowStyleReset = "\x1b[0m"

// fileRowStyle marks a file row bold so it stands out from the directory
// rows above it and the item rows nested under it, without assuming
// anything about the terminal's color scheme.
const fileRowStyle = "\x1b[1m"

// selectedRowStyle marks a row as selected with bold + reverse video,
// tinted with fgColor beforehand — the same accent frameColors gives
// this pane's border/scrollbar, so the highlight reads as "active"
// (accent-tinted) or "inactive" (theme-default) the same way they do.
// Reverse video does the contrast work, so the highlight still adapts
// to dark, light, or transparent terminal themes rather than assuming
// one. fgColor must come before bold/reverse: gocui's escape parser
// assigns fg/bg colors outright (rather than OR-ing them in), so
// setting it after bold/reverse would wipe those bits back out.
func selectedRowStyle(fgColor gocui.Attribute) string {
	code := 39
	if fgColor != gocui.ColorDefault {
		code = 30 + int(fgColor) - 1
	}
	return fmt.Sprintf("\x1b[%d;1;7m", code)
}

func (a *App) Render(g *gocui.Gui) error {
	categoryStyle := selectedRowStyle(a.Categories.FocusFgColor(g))
	if err := a.Categories.Render(g, a.CategorySelected, func(v *gocui.View, contentWidth int) error {
		rowWidth := contentWidth - len("  ")
		for i, kind := range categoryOrder {
			count := len(a.itemsByCategory[kind])
			row := "  " + formatCategoryRow(rowWidth, kind, count)
			if i == a.CategorySelected {
				_, _ = fmt.Fprintf(v, "%s%s%s\n", categoryStyle, row, rowStyleReset)
			} else {
				_, _ = fmt.Fprintln(v, row)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	rows := a.visibleRows()
	if err := a.Content.SetTitle(g, fmt.Sprintf(" [2] Items - %s ", a.currentCategory())); err != nil {
		return err
	}

	focus := a.RowSelected
	if len(rows) == 0 {
		focus = -1
	}

	itemStyle := selectedRowStyle(a.Content.FocusFgColor(g))
	return a.Content.Render(g, focus, func(v *gocui.View, contentWidth int) error {
		if len(rows) == 0 {
			_, _ = fmt.Fprintln(v, "No items in this category")
			return nil
		}
		for i, r := range rows {
			line := "  " + formatRow(r, contentWidth-2)
			switch {
			case i == a.RowSelected:
				line = runewidth.FillRight(line, contentWidth)
				_, _ = fmt.Fprintf(v, "%s%s%s\n", itemStyle, line, rowStyleReset)
			case r.kind == rowKindFile:
				_, _ = fmt.Fprintf(v, "%s%s%s\n", fileRowStyle, line, rowStyleReset)
			default:
				_, _ = fmt.Fprintln(v, line)
			}
		}
		return nil
	})
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
