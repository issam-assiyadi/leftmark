package scrollview

import "github.com/awesome-gocui/gocui"

var roundedFrameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}

func frameColors(g *gocui.Gui, v *gocui.View) (gocui.Attribute, gocui.Attribute) {
	if g.Highlight && g.CurrentView() == v {
		return g.SelFrameColor, g.SelBgColor
	}
	return g.FrameColor, g.BgColor
}

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
