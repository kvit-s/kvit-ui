package kvitui

import (
	"slices"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// Popup is a surface shown above everything else in a Window: a tooltip, a
// popover, a dialog.
type Popup struct {
	// Panel is what is shown.
	Panel unison.Paneler
	// Place gives where it goes, in the window's content coordinates, from
	// the content's bounds and the panel's preferred size. It is asked on
	// every layout, so the popup follows its anchor when the window changes.
	Place func(bounds geom.Rect, size geom.Size) geom.Rect
	// Modal lays a scrim over everything below it that takes every press,
	// so the popup has to be answered first.
	Modal bool
	// OnEscape runs when Escape is pressed while this is the topmost popup
	// that is not passive; nil leaves Escape alone.
	OnEscape func()
	// Passive is a popup that only informs, such as a tooltip: it never has
	// the focus, so Escape passes it by and goes to the popup under it, or to
	// the window when there is none. In  a tooltip never takes the focus,
	// so Escape reached what was under it, such as the panel the control
	// with the tooltip is in.
	Passive bool
	// OnPressOutside runs when the pointer is pressed outside it; the press
	// still reaches whatever is under the pointer unless the popup is modal.
	OnPressOutside func()
	// Anchor, when set, is what the popup belongs to: once it is no longer in
	// the window, as when the page holding it is replaced, the popup goes too.
	Anchor unison.Paneler
}

// windowOf is the Kvit window a panel is shown in, or nil.
func (u *UI) windowOf(p unison.Paneler) *Window {
	if p == nil {
		return nil
	}
	if w := p.AsPanel().Window(); w != nil {
		return u.windows[w]
	}
	return nil
}

// Show puts a popup above everything in the window and returns the function
// that takes it away again. Popups shown later are above those shown before.
func (w *Window) Show(p *Popup) (hide func()) {
	w.popups = append(w.popups, p)
	w.rebuild()
	return func() {
		if i := slices.Index(w.popups, p); i >= 0 {
			w.popups = slices.Delete(w.popups, i, i+1)
			w.rebuild()
		}
	}
}

// Showing reports whether a popup is shown in the window.
func (w *Window) Showing(p *Popup) bool { return slices.Contains(w.popups, p) }

// top is the topmost popup, or nil.
func (w *Window) top() *Popup {
	if len(w.popups) == 0 {
		return nil
	}
	return w.popups[len(w.popups)-1]
}

// scrim is the layer behind a modal popup: it draws a faint shade over the
// window and takes every press, so nothing behind can be used until the
// popup is answered.
type scrim struct {
	unison.Panel
	ui *UI
}

func newScrim(ui *UI) *scrim {
	s := &scrim{ui: ui}
	s.Self = s
	s.MouseDownCallback = func(geom.Point, int, int, mod.Modifiers) bool { return true }
	s.MouseUpCallback = func(geom.Point, int, mod.Modifiers) bool { return true }
	s.MouseWheelCallback = func(geom.Point, geom.Point, mod.Modifiers) bool { return true }
	s.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		painterFor(gc, ui).roundTint(s.ContentRect(true), 0, ui.Theme.Tokens().TextPrimary, 0.18)
	}
	s.Accessibility.Role = role.None
	return s
}

// addPopups puts the popups first among the content's children, topmost
// first, which is the child unison draws last and asks first about a press,
// with a scrim behind the topmost modal one.
func (w *Window) addPopups(content *unison.Panel) {
	modal := -1
	for i, p := range w.popups {
		if p.Modal {
			modal = i
		}
	}
	at := 0
	for i := len(w.popups) - 1; i >= 0; i-- {
		content.AddChildAtIndex(w.popups[i].Panel, at)
		at++
		if i == modal {
			if w.scrim == nil {
				w.scrim = newScrim(w.ui)
			}
			content.AddChildAtIndex(w.scrim, at)
			at++
		}
	}
}

// layoutPopups places the scrim over the whole content and each popup where
// it asks to go.
func (w *Window) layoutPopups(r geom.Rect) {
	gone := slices.DeleteFunc(slices.Clone(w.popups), func(p *Popup) bool {
		return p.Anchor == nil || p.Anchor.AsPanel().Window() == w.Window
	})
	if len(gone) > 0 {
		w.popups = slices.DeleteFunc(w.popups, func(p *Popup) bool { return slices.Contains(gone, p) })
		for _, p := range gone {
			p.Panel.AsPanel().RemoveFromParent()
		}
	}
	if w.scrim != nil && w.scrim.Parent() != nil {
		w.scrim.SetFrameRect(r)
	}
	for _, p := range w.popups {
		_, size, _ := p.Panel.AsPanel().Sizes(geom.Size{})
		p.Panel.AsPanel().SetFrameRect(p.Place(r, size))
	}
}

// popupKeys closes the topmost popup on Escape, before any control sees the
// key. Passive popups are passed by. A popup that is not passive and has no
// OnEscape leaves the key to the focused control, as the typeahead's list
// does for its field, which closes the list.
func (w *Window) popupKeys(key unison.KeyCode) bool {
	if key != unison.KeyEscape {
		return false
	}
	for i := len(w.popups) - 1; i >= 0; i-- {
		p := w.popups[i]
		if p.Passive {
			continue
		}
		if p.OnEscape != nil {
			p.OnEscape()
			return true
		}
		return false
	}
	return false
}

// popupPress tells the popups a press landed outside them. It returns true to
// stop the press, which a modal popup does.
func (w *Window) popupPress(where geom.Point) bool {
	for i := len(w.popups) - 1; i >= 0; i-- {
		p := w.popups[i]
		inside := where.In(p.Panel.AsPanel().FrameRect())
		if inside {
			return false
		}
		if p.OnPressOutside != nil {
			p.OnPressOutside()
		}
		if p.Modal {
			return false // the scrim takes it
		}
	}
	return false
}

// PlaceBeside is where a tooltip goes, in the order kvit-cash's copy of the
// library settled on, taking the first that fits: beside the anchor on its
// trailing side, level with its top, since a row in a list has neighbours
// above and below and none across; inside the anchor's own band, at its
// trailing end, for a row as wide as its surface; then below it; then above.
func PlaceBeside(ui *UI, anchor unison.Paneler) func(bounds geom.Rect, size geom.Size) geom.Rect {
	return placeBeside(ui, func() geom.Rect { return anchorIn(anchor) })
}

// placeBeside is PlaceBeside for a box in the window's content coordinates,
// such as one cell of a table, asked on every layout.
func placeBeside(ui *UI, box func() geom.Rect) func(bounds geom.Rect, size geom.Size) geom.Rect {
	return func(bounds geom.Rect, size geom.Size) geom.Rect {
		a := box()
		m := ui.Interface
		gap, room := float32(m.SpaceTight()), float32(m.Space())
		switch {
		case a.Right()+gap+size.Width+room <= bounds.Right():
			y := a.Y
			if over := y + size.Height + room - bounds.Bottom(); over > 0 {
				y = max(bounds.Y, y-over)
			}
			return geom.NewRect(a.Right()+gap, y, size.Width, size.Height)
		case size.Height <= a.Height:
			return geom.NewRect(max(bounds.X, a.Right()-size.Width), a.Y+(a.Height-size.Height)/2, size.Width, size.Height)
		}
		x := a.X
		if over := x + size.Width + room - bounds.Right(); over > 0 {
			x -= over
		}
		if a.Bottom()+gap+size.Height+room <= bounds.Bottom() {
			return geom.NewRect(x, a.Bottom()+gap, size.Width, size.Height)
		}
		return geom.NewRect(x, a.Y-gap-size.Height, size.Width, size.Height)
	}
}

// PlaceBelow is under the anchor, a near space below it, starting at its
// left edge, or ending at its right edge where that would run off the window.
func PlaceBelow(ui *UI, anchor unison.Paneler) func(bounds geom.Rect, size geom.Size) geom.Rect {
	return func(bounds geom.Rect, size geom.Size) geom.Rect {
		a := anchorIn(anchor)
		x := a.X
		if x+size.Width > bounds.Right() {
			x = a.Right() - size.Width
		}
		return geom.NewRect(x, a.Bottom()+float32(ui.Interface.SpaceNear()), size.Width, size.Height)
	}
}

// PlaceOver is at the anchor's top-left corner.
func PlaceOver(anchor unison.Paneler) func(bounds geom.Rect, size geom.Size) geom.Rect {
	return func(_ geom.Rect, size geom.Size) geom.Rect {
		a := anchorIn(anchor)
		return geom.NewRect(a.X, a.Y, size.Width, size.Height)
	}
}

// PlaceCentred is in the middle of the anchor, or of the window when the
// anchor is nil.
func PlaceCentred(anchor unison.Paneler) func(bounds geom.Rect, size geom.Size) geom.Rect {
	return func(bounds geom.Rect, size geom.Size) geom.Rect {
		a := bounds
		if anchor != nil {
			a = anchorIn(anchor)
		}
		return geom.NewRect(a.X+(a.Width-size.Width)/2, a.Y+(a.Height-size.Height)/2, size.Width, size.Height)
	}
}

// anchorIn is an anchor's box in its window's content coordinates.
func anchorIn(anchor unison.Paneler) geom.Rect {
	return partIn(anchor, anchor.AsPanel().ContentRect(true))
}

// partIn is a box in a panel's own coordinates, such as one cell of a table,
// in its window's content coordinates.
func partIn(owner unison.Paneler, box geom.Rect) geom.Rect {
	p := owner.AsPanel()
	r := p.RectToRoot(box)
	if w := p.Window(); w != nil {
		return w.Content().RectFromRoot(r)
	}
	return r
}

// Popups are the popups shown, the topmost last.
func (w *Window) Popups() []*Popup { return slices.Clone(w.popups) }
