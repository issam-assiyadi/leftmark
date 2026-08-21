package scrollview

import "github.com/jroimartin/gocui"

const (
	borderHorizontal  = '─'
	borderVertical    = '│'
	cornerTopLeft     = '╭'
	cornerTopRight    = '╮'
	cornerBottomLeft  = '╰'
	cornerBottomRight = '╯'
)

// drawRoundedFrame draws a border with rounded corners around
// (x0,y0)-(x1,y1). gocui hard-codes square corners for any view with
// Frame set to true, so views that want rounded corners keep Frame
// false and get their border drawn here instead, replicating gocui's
// own edge/corner/title logic (see gocui's drawFrameEdges,
// drawFrameCorners and drawTitle) with rounded corner glyphs.
func drawRoundedFrame(g *gocui.Gui, v *gocui.View, x0, y0, x1, y1 int, title string) error {
	maxX, maxY := g.Size()

	fgColor, bgColor := frameColors(g, v)

	for x := x0 + 1; x < x1 && x < maxX; x++ {
		if x < 0 {
			continue
		}
		if y0 > -1 && y0 < maxY {
			if err := g.SetRune(x, y0, borderHorizontal, fgColor, bgColor); err != nil {
				return err
			}
		}
		if y1 > -1 && y1 < maxY {
			if err := g.SetRune(x, y1, borderHorizontal, fgColor, bgColor); err != nil {
				return err
			}
		}
	}
	for y := y0 + 1; y < y1 && y < maxY; y++ {
		if y < 0 {
			continue
		}
		if x0 > -1 && x0 < maxX {
			if err := g.SetRune(x0, y, borderVertical, fgColor, bgColor); err != nil {
				return err
			}
		}
		if x1 > -1 && x1 < maxX {
			if err := g.SetRune(x1, y, borderVertical, fgColor, bgColor); err != nil {
				return err
			}
		}
	}

	corners := []struct {
		x, y int
		ch   rune
	}{
		{x0, y0, cornerTopLeft},
		{x1, y0, cornerTopRight},
		{x0, y1, cornerBottomLeft},
		{x1, y1, cornerBottomRight},
	}
	for _, c := range corners {
		if c.x >= 0 && c.y >= 0 && c.x < maxX && c.y < maxY {
			if err := g.SetRune(c.x, c.y, c.ch, fgColor, bgColor); err != nil {
				return err
			}
		}
	}

	if title != "" && y0 >= 0 && y0 < maxY {
		titleFgColor := fgColor
		if g.Highlight && g.CurrentView() == v {
			titleFgColor |= gocui.AttrBold
		}
		for i, ch := range title {
			x := x0 + i + 2
			if x < 0 {
				continue
			} else if x > x1-2 || x >= maxX {
				break
			}
			if err := g.SetRune(x, y0, ch, titleFgColor, bgColor); err != nil {
				return err
			}
		}
	}

	return nil
}

// frameColors returns the color pair a bordered view (or anything
// visually tied to it, such as its scrollbar) should be drawn with:
// the selection colors while it holds focus, the default frame
// colors otherwise.
func frameColors(g *gocui.Gui, v *gocui.View) (gocui.Attribute, gocui.Attribute) {
	if g.Highlight && g.CurrentView() == v {
		return g.SelFgColor, g.SelBgColor
	}
	return g.FgColor, g.BgColor
}

// FocusFgColor returns the foreground color this pane's border and
// scrollbar are currently drawn with, so other elements tied to this
// pane (such as its selected row) can share the same focus/unfocused
// accent instead of hardcoding their own colors.
func (v *View) FocusFgColor(g *gocui.Gui) gocui.Attribute {
	if v == nil || g == nil {
		return gocui.ColorDefault
	}
	wrapper, err := g.View(v.wrapperViewName)
	if err != nil {
		return gocui.ColorDefault
	}
	fgColor, _ := frameColors(g, wrapper)
	return fgColor
}
