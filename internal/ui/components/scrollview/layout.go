package scrollview

import "github.com/awesome-gocui/gocui"

func (v *View) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v == nil || g == nil {
		return nil
	}
	if x0 >= x1 || y0 >= y1 {
		return nil
	}

	// overlaps (0, unused here) only matters for T-junction/cross frame
	// runes between views that share an exact edge — none of our panes
	// do, they're independent rectangles that may freely overlap.
	wrapper, err := g.SetView(v.wrapperViewName, x0, y0, x1, y1, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		wrapper.Frame = true
		wrapper.FrameRunes = roundedFrameRunes
		wrapper.Title = v.title
	}

	contentLeft := x0 + 1
	contentTop := y0
	contentBottom := y1 - 1

	contentRight := x1 - 1
	if !v.hideScrollbar {
		contentRight = x1 - v.scrollbarWidth
	}

	if contentLeft >= contentRight || contentTop >= contentBottom {
		return nil
	}

	contentView, err := g.SetView(v.contentViewName, contentLeft, contentTop, contentRight, contentBottom, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		contentView.Frame = false
		contentView.Wrap = false
		contentView.Autoscroll = false
	}

	if v.hideScrollbar {
		return nil
	}

	scrollLeft := contentRight
	scrollTop := y0
	scrollRight := x1
	scrollBottom := y1

	if scrollLeft >= scrollRight || scrollTop >= scrollBottom {
		return nil
	}

	scrollView, err := g.SetView(v.scrollViewName, scrollLeft, scrollTop, scrollRight, scrollBottom, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		scrollView.Frame = false
		scrollView.Wrap = false
		scrollView.BgColor = gocui.ColorDefault
	}
	scrollView.FgColor, _ = frameColors(g, wrapper)

	return nil
}

// Delete removes this pane's views from g, if present. Layout recreates
// them from scratch (fresh scroll origin, one-time init rerun) the next
// time it's called. Intended for a pane that comes and goes, such as a
// modal: gocui keeps drawing whatever views remain registered every
// frame regardless of whether Layout is still being called for them, so
// simply no longer laying a pane out does not hide it.
func (v *View) Delete(g *gocui.Gui) {
	if v == nil || g == nil {
		return
	}
	for _, name := range []string{v.wrapperViewName, v.contentViewName, v.scrollViewName} {
		_ = g.DeleteView(name)
	}
}
