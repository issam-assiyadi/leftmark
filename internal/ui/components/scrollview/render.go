package scrollview

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func (v *View) SetTitle(g *gocui.Gui, title string) error {
	if v == nil || g == nil {
		return nil
	}

	wrapper, err := g.View(v.wrapperViewName)
	if err != nil {
		if err == gocui.ErrUnknownView {
			return nil
		}
		return err
	}

	wrapper.Title = title
	return nil
}

func (v *View) Render(g *gocui.Gui, focus int, renderFn func(*gocui.View, int) error) error {
	if v == nil || g == nil {
		return nil
	}

	contentView, err := v.contentView(g)
	if err != nil || contentView == nil {
		return err
	}

	scrollView, err := v.scrollView(g)
	if err != nil || scrollView == nil {
		return err
	}

	contentView.Clear()
	scrollView.Clear()

	if err := scrollView.SetOrigin(0, 0); err != nil {
		return err
	}

	if renderFn != nil {
		if err := renderFn(contentView, v.renderWidth(contentView)); err != nil {
			return err
		}
	}

	if err := v.ensureVisible(contentView, focus); err != nil {
		return err
	}

	return v.renderScrollbar(contentView, scrollView)
}

func (v *View) ensureVisible(contentView *gocui.View, focus int) error {
	if contentView == nil || focus < 0 {
		return nil
	}

	_, height := contentView.Size()
	if height <= 0 {
		return nil
	}

	originX, originY := contentView.Origin()
	newOriginY := originY
	if focus < newOriginY {
		newOriginY = focus
	} else if focus >= newOriginY+height {
		newOriginY = focus - height + 1
	}

	if maxOrigin := v.maxOrigin(contentView); newOriginY > maxOrigin {
		newOriginY = maxOrigin
	}
	if newOriginY < 0 {
		newOriginY = 0
	}

	if newOriginY == originY {
		return nil
	}

	return contentView.SetOrigin(originX, newOriginY)
}

func (v *View) renderScrollbar(contentView, scrollView *gocui.View) error {
	if contentView == nil || scrollView == nil {
		return nil
	}

	lines := len(contentView.BufferLines())
	_, height := contentView.Size()
	scrollWidth, scrollHeight := scrollView.Size()
	if height <= 0 || scrollHeight <= 0 || scrollWidth <= 0 {
		return nil
	}
	if lines <= 0 {
		lines = height
	}

	maxOrigin := v.maxOrigin(contentView)
	originX, originY := contentView.Origin()
	if originY > maxOrigin {
		originY = maxOrigin
		if err := contentView.SetOrigin(originX, originY); err != nil {
			return err
		}
	}

	if lines <= height {
		return nil
	}

	thumbHeight := max(1, scrollHeight*height/lines)
	if thumbHeight > scrollHeight {
		thumbHeight = scrollHeight
	}

	thumbTop := 0
	if maxOrigin > 0 && scrollHeight > thumbHeight {
		thumbTop = originY * (scrollHeight - thumbHeight) / maxOrigin
	}

	track := repeatRune(' ', scrollWidth)
	thumb := repeatRune('█', scrollWidth)
	for i := 0; i < scrollHeight; i++ {
		line := track
		if i >= thumbTop && i < thumbTop+thumbHeight {
			line = thumb
		}
		_, _ = fmt.Fprintln(scrollView, line)
	}

	return nil
}
