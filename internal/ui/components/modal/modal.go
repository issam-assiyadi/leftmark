// Package modal is a floating overlay pane built on scrollview.View. It
// opens on top of whatever panes are already on screen, closes back to
// whichever one was focused before it, and scrolls like any other pane.
// Features share one Modal by supplying their own title and render
// callback; open/close/focus bookkeeping is the same regardless of what's
// inside.
package modal

import (
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/scrollview"
)

type Config struct {
	BaseName       string
	ScrollbarWidth int
}

type Modal struct {
	view            *scrollview.View
	open            bool
	title           string
	focusLine       int
	focusPending    bool
	previousFocused string
}

func New(cfg Config) *Modal {
	return &Modal{
		view: scrollview.New(scrollview.Config{
			BaseName:       cfg.BaseName,
			ScrollbarWidth: cfg.ScrollbarWidth,
		}),
	}
}

func (m *Modal) WrapperName() string { return m.view.WrapperName() }

func (m *Modal) IsOpen() bool { return m.open }

func (m *Modal) FocusLine() int { return m.focusLine }

func (m *Modal) FocusFgColor(g *gocui.Gui) gocui.Attribute { return m.view.FocusFgColor(g) }

// Open switches the modal to its open state and returns the view name
// that should now hold focus. currentFocused is remembered so Close can
// restore it. focusLine scrolls that line into view once, on the first
// render after opening; pass a negative value to open at the top.
func (m *Modal) Open(currentFocused, title string, focusLine int) string {
	m.title = title
	m.focusLine = focusLine
	m.focusPending = focusLine >= 0
	m.previousFocused = currentFocused
	m.open = true
	return m.view.WrapperName()
}

// Close switches the modal to its closed state and returns the view name
// to restore focus to.
func (m *Modal) Close(g *gocui.Gui) string {
	m.open = false
	m.view.Delete(g)
	return m.previousFocused
}

func (m *Modal) ScrollDown(g *gocui.Gui) error { return m.view.ScrollDown(g) }
func (m *Modal) ScrollUp(g *gocui.Gui) error   { return m.view.ScrollUp(g) }
func (m *Modal) JumpTop(g *gocui.Gui) error    { return m.view.JumpTop(g) }
func (m *Modal) JumpBottom(g *gocui.Gui) error { return m.view.JumpBottom(g) }

// BindKeys registers the modal's standard controls on its own view, so
// they only fire while it's focused: q/Esc call onClose, and the arrow
// keys, j/k, PgUp/PgDn, Home, and End scroll it.
func (m *Modal) BindKeys(g *gocui.Gui, onClose func(*gocui.Gui, *gocui.View) error) error {
	keys := []struct {
		key interface{}
		fn  func(*gocui.Gui, *gocui.View) error
	}{
		{'q', onClose},
		{gocui.KeyEsc, onClose},
		{gocui.KeyArrowDown, func(g *gocui.Gui, v *gocui.View) error { return m.ScrollDown(g) }},
		{'j', func(g *gocui.Gui, v *gocui.View) error { return m.ScrollDown(g) }},
		{gocui.KeyArrowUp, func(g *gocui.Gui, v *gocui.View) error { return m.ScrollUp(g) }},
		{'k', func(g *gocui.Gui, v *gocui.View) error { return m.ScrollUp(g) }},
		{gocui.KeyPgdn, func(g *gocui.Gui, v *gocui.View) error { return m.ScrollDown(g) }},
		{gocui.KeyPgup, func(g *gocui.Gui, v *gocui.View) error { return m.ScrollUp(g) }},
		{gocui.KeyHome, func(g *gocui.Gui, v *gocui.View) error { return m.JumpTop(g) }},
		{gocui.KeyEnd, func(g *gocui.Gui, v *gocui.View) error { return m.JumpBottom(g) }},
	}
	for _, k := range keys {
		if err := g.SetKeybinding(m.WrapperName(), k.key, gocui.ModNone, k.fn); err != nil {
			return err
		}
	}
	return nil
}

// Layout lays the modal out as a centered overlay covering ~90% of the
// terminal, clamped to a usable minimum. It's a no-op while closed.
func (m *Modal) Layout(g *gocui.Gui, maxX, maxY int) error {
	if !m.open {
		return nil
	}
	x0, y0, x1, y1 := bounds(maxX, maxY)
	return m.view.Layout(g, x0, y0, x1, y1)
}

func bounds(maxX, maxY int) (x0, y0, x1, y1 int) {
	width := maxX * 9 / 10
	if width < 20 {
		width = min(20, maxX)
	}
	height := maxY * 9 / 10
	if height < 6 {
		height = min(6, maxY)
	}

	x0 = (maxX - width) / 2
	y0 = (maxY - height) / 2
	x1 = x0 + width - 1
	y1 = y0 + height - 1
	if x1 >= maxX {
		x1 = maxX - 1
	}
	if y1 >= maxY {
		y1 = maxY - 1
	}
	return x0, y0, x1, y1
}

// Render draws the modal's title and delegates content to renderFn. It's
// a no-op while closed.
func (m *Modal) Render(g *gocui.Gui, renderFn func(*gocui.View, int) error) error {
	if !m.open {
		return nil
	}
	if err := m.view.SetTitle(g, m.title); err != nil {
		return err
	}
	focus := -1
	if m.focusPending {
		focus = m.focusLine
		m.focusPending = false
	}
	return m.view.Render(g, focus, renderFn)
}
