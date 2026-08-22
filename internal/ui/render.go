package ui

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/mattn/go-runewidth"

	"github.com/issam-assiyadi/leftmark/domain"
)

const rowStyleReset = "\x1b[0m"

// fileRowStyle marks a file row bold so it stands out from the directory
// rows above it and the item rows nested under it, without assuming
// anything about the terminal's color scheme.
const fileRowStyle = "\x1b[1m"

// basicColorSGR maps gocui's 8 basic color constants to their standard
// foreground SGR code, the encoding selectedRowStyle needs to write a raw
// ANSI escape into a view's content. gocui.Attribute packs color,
// validity, and style bits together rather than storing that code
// directly, so this is a lookup rather than arithmetic on the value —
// and only needs to cover the handful of colors this app actually
// assigns to g.SelFgColor.
var basicColorSGR = map[gocui.Attribute]int{
	gocui.ColorBlack:   30,
	gocui.ColorRed:     31,
	gocui.ColorGreen:   32,
	gocui.ColorYellow:  33,
	gocui.ColorBlue:    34,
	gocui.ColorMagenta: 35,
	gocui.ColorCyan:    36,
	gocui.ColorWhite:   37,
}

func selectedRowStyle(fgColor gocui.Attribute) string {
	code, ok := basicColorSGR[fgColor]
	if !ok {
		code = 39
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
	if err := a.Content.Render(g, focus, func(v *gocui.View, contentWidth int) error {
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
	}); err != nil {
		return err
	}

	markerStyle := selectedRowStyle(a.Preview.FocusFgColor(g))
	gutterWidth := previewGutterWidth(len(a.previewLines))
	focusLine := a.Preview.FocusLine()
	return a.Preview.Render(g, func(v *gocui.View, contentWidth int) error {
		for i, line := range a.previewLines {
			if i != focusLine {
				_, _ = fmt.Fprintln(v, line)
				continue
			}
			marker := fmt.Sprintf("%*d │ %s", gutterWidth, i+1, a.previewFocusText)
			marker = runewidth.FillRight(marker, contentWidth)
			_, _ = fmt.Fprintf(v, "%s%s%s\n", markerStyle, marker, rowStyleReset)
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
