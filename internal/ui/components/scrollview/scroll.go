package scrollview

import "github.com/jroimartin/gocui"

func (v *View) ScrollDown(g *gocui.Gui) error {
	contentView, err := v.contentView(g)
	if err != nil || contentView == nil {
		return err
	}

	ox, oy := contentView.Origin()
	maxOrigin := v.maxOrigin(contentView)
	if oy >= maxOrigin {
		return nil
	}

	return contentView.SetOrigin(ox, oy+1)
}

func (v *View) ScrollUp(g *gocui.Gui) error {
	contentView, err := v.contentView(g)
	if err != nil || contentView == nil {
		return err
	}

	ox, oy := contentView.Origin()
	if oy <= 0 {
		return nil
	}

	return contentView.SetOrigin(ox, oy-1)
}

func (v *View) JumpTop(g *gocui.Gui) error {
	contentView, err := v.contentView(g)
	if err != nil || contentView == nil {
		return err
	}

	ox, _ := contentView.Origin()
	return contentView.SetOrigin(ox, 0)
}
