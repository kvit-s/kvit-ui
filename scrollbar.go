package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// NewScrollBar returns a scroll bar that occupies a strip of its own rather
// than floating over what it scrolls. A bar drawn over the content hides the
// right-hand column of a table and the last character of every cut-short
// label, and content that moves aside when a bar appears is worse, so the
// strip is kept whether or not the bar is showing. The bar is hidden when
// everything fits: a full-length handle says "there is more" when there is
// not.
//
// It is unison's scroll bar, so dragging, clicking the track and what a
// screen reader is told are unison's, drawn in the Kvit colours: the panel
// ground, and a handle in the border colour that turns the strong border
// colour under the pointer or while dragged.
func NewScrollBar(ui *UI, horizontal bool) *unison.ScrollBar {
	b := unison.NewScrollBar(horizontal)
	b.ThumbIndent = 0
	var hovered, dragging bool
	fit := func() {
		m := ui.Interface
		b.MinimumThickness = float32(m.SpaceWide())
		// A handle sized in proportion to 250,000 rows is a fraction of a
		// pixel, and cannot be grabbed.
		b.MinimumThumb = float32(m.Px(24))
	}
	fit()
	b.SetBorder(scrollBarInset{ui, horizontal})
	b.SetSizer(func(hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
		fit()
		return b.DefaultSizes(hint)
	})
	b.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		fit()
		if b.Extent() >= b.Max()*0.999 {
			return
		}
		t := ui.Theme.Tokens()
		p := painterFor(gc, ui)
		p.fill(b.ContentRect(true), t.PanelBackground)
		thumb := b.Thumb()
		if thumb.Empty() {
			return
		}
		ink := t.Border
		if hovered || dragging {
			ink = t.BorderStrong
		}
		p.round(thumb, min(thumb.Width, thumb.Height)/2, ink)
	}
	track := func(where geom.Point) {
		if over := where.In(b.Thumb()); over != hovered {
			hovered = over
			b.MarkForRedraw()
		}
	}
	b.MouseMoveCallback = func(where geom.Point, mods mod.Modifiers) bool {
		track(where)
		return b.DefaultMouseMove(where, mods)
	}
	b.MouseEnterCallback = func(where geom.Point, mods mod.Modifiers) bool {
		track(where)
		return b.DefaultMouseEnter(where, mods)
	}
	b.MouseExitCallback = func() bool {
		hovered = false
		b.MarkForRedraw()
		return b.DefaultMouseExit()
	}
	b.MouseDownCallback = func(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
		if b.Extent() >= b.Max()*0.999 {
			return false
		}
		dragging = true
		return b.DefaultMouseDown(where, button, clicks, mods)
	}
	b.MouseUpCallback = func(where geom.Point, button int, mods mod.Modifiers) bool {
		dragging = false
		track(where)
		return b.DefaultMouseUp(where, button, mods)
	}
	return b
}

// scrollBarInset keeps the handle a tight space in from both long edges of
// its strip.
type scrollBarInset struct {
	ui         *UI
	horizontal bool
}

func (s scrollBarInset) Insets() geom.Insets {
	w := float32(s.ui.Interface.SpaceTight())
	if s.horizontal {
		return geom.Insets{Top: w, Bottom: w}
	}
	return geom.Insets{Left: w, Right: w}
}

func (s scrollBarInset) Draw(*unison.Canvas, geom.Rect) {}
