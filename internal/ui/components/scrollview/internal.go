package scrollview

import "github.com/jroimartin/gocui"

func (v *View) contentView(g *gocui.Gui) (*gocui.View, error) {
	if v == nil || g == nil {
		return nil, nil
	}

	view, err := g.View(v.contentViewName)
	if err != nil {
		if err == gocui.ErrUnknownView {
			return nil, nil
		}
		return nil, err
	}

	return view, nil
}

func (v *View) scrollView(g *gocui.Gui) (*gocui.View, error) {
	if v == nil || g == nil {
		return nil, nil
	}

	view, err := g.View(v.scrollViewName)
	if err != nil {
		if err == gocui.ErrUnknownView {
			return nil, nil
		}
		return nil, err
	}

	return view, nil
}

func (v *View) maxOrigin(contentView *gocui.View) int {
	if contentView == nil {
		return 0
	}

	_, height := contentView.Size()
	if height <= 0 {
		return 0
	}

	lines := len(contentView.BufferLines())
	maxOrigin := lines - height
	if maxOrigin < 0 {
		return 0
	}

	return maxOrigin
}

func (v *View) renderWidth(contentView *gocui.View) int {
	if contentView == nil {
		return 1
	}

	width, _ := contentView.Size()
	width -= v.scrollbarWidth
	if width < 1 {
		return 1
	}

	return width
}

func repeatRune(ch rune, count int) string {
	if count <= 0 {
		return ""
	}

	runes := make([]rune, count)
	for i := range runes {
		runes[i] = ch
	}

	return string(runes)
}
