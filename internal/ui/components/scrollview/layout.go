package scrollview

import "github.com/jroimartin/gocui"

func (v *View) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v == nil || g == nil {
		return nil
	}
	if x0 >= x1 || y0 >= y1 {
		return nil
	}

	wrapper, err := g.SetView(v.wrapperViewName, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		wrapper.Frame = false
	}
	if err := drawRoundedFrame(g, wrapper, x0, y0, x1, y1, v.title); err != nil {
		return err
	}

	contentLeft := x0 + 1
	contentTop := y0
	contentBottom := y1 - 1

	contentRight := x1 - 1
	if !v.hideScrollbar {
		contentRight = x1 - v.scrollbarWidth - 1
	}

	if contentLeft >= contentRight || contentTop >= contentBottom {
		return nil
	}

	contentView, err := g.SetView(v.contentViewName, contentLeft, contentTop, contentRight, contentBottom)
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

	scrollLeft := x1 - v.scrollbarWidth + 1
	scrollTop := y0
	scrollRight := x1
	scrollBottom := y1

	if scrollLeft >= scrollRight || scrollTop >= scrollBottom {
		return nil
	}

	scrollView, err := g.SetView(v.scrollViewName, scrollLeft, scrollTop, scrollRight, scrollBottom)
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
