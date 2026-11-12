package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// Mover lets the reader move a panel shown in a window's popup layer, such as
// a modal Popover with a title line of its own, by pressing in the band along
// its top and dragging. The point pressed stays under the pointer as far as
// the popup's bounds allow. Once moved, the panel keeps its place, measured
// from the corner of the bounds, when its content changes height or the
// window is resized, and is kept whole inside the bounds.
//
// The popup's Place function passes what it would answer through Place. A
// surface that should open in the middle every time calls Reset when it
// opens; one that should open where the reader left it does not. A Dialog
// has a drag of its own and needs none of this.
type Mover struct {
	// EscapeAbandons makes Escape during a drag put the panel back where the
	// drag began, in a Kvit window. Without it Escape does nothing to the
	// drag.
	EscapeAbandons bool

	ui       *UI
	panel    *unison.Panel
	band     func() float32
	bounds   geom.Rect
	at       geom.Point // where the panel was moved to, from the corner of the bounds
	moved    bool
	dragging bool
	grab     geom.Point // the point pressed, in the panel's coordinates
}

// NewMover wraps p's pointer callbacks, so set any of p's own before calling
// it; band is the height of the band a drag starts in, from the panel's top,
// in the panel's coordinates. A press anywhere else, and with any button but
// the left, goes to the callback that was there before, and a control in the
// band takes its own presses.
func NewMover(ui *UI, p unison.Paneler, band func() float32) *Mover {
	m := &Mover{ui: ui, panel: p.AsPanel(), band: band}
	panel := m.panel
	down, drag, up := panel.MouseDownCallback, panel.MouseDragCallback, panel.MouseUpCallback
	panel.MouseDownCallback = func(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
		if button != unison.ButtonLeft || where.Y < 0 || where.Y >= m.band() {
			return down != nil && down(where, button, clicks, mods)
		}
		m.dragging, m.grab = true, where
		if w := ui.windowOf(panel); w != nil && m.EscapeAbandons {
			at, moved := m.at, m.moved
			w.startDrag(func() {
				m.dragging, m.at, m.moved = false, at, moved
				if parent := panel.Parent(); parent != nil {
					parent.MarkForLayoutAndRedraw()
				}
			})
		}
		return true
	}
	panel.MouseDragCallback = func(where geom.Point, button int, mods mod.Modifiers) bool {
		if !m.dragging {
			return drag != nil && drag(where, button, mods)
		}
		m.drag(where)
		return true
	}
	panel.MouseUpCallback = func(where geom.Point, button int, mods mod.Modifiers) bool {
		if w := ui.windowOf(panel); w != nil && m.EscapeAbandons {
			w.endDrag()
		}
		if !m.dragging {
			return up != nil && up(where, button, mods)
		}
		m.dragging = false
		return true
	}
	return m
}

func (m *Mover) drag(where geom.Point) {
	f := m.panel.FrameRect()
	m.at = f.Point.Add(where.Sub(m.grab)).Sub(m.bounds.Point)
	m.moved = true
	f.Point = m.inside(m.bounds, f.Size)
	m.at = f.Point.Sub(m.bounds.Point)
	m.panel.SetFrameRect(f)
	m.panel.MarkForRedraw()
}

// inside is where the moved panel's top-left corner goes for it to be whole
// inside bounds, or at their top-left corner when it is larger.
func (m *Mover) inside(bounds geom.Rect, size geom.Size) geom.Point {
	x := min(bounds.X+m.at.X, bounds.Right()-size.Width)
	y := min(bounds.Y+m.at.Y, bounds.Bottom()-size.Height)
	return geom.Point{X: max(x, bounds.X), Y: max(y, bounds.Y)}
}

// Place is for the popup's Place function: bounds as that function is given
// them, and where the panel would go unmoved. It answers unmoved until the
// reader moves the panel, and afterwards a box of unmoved's size where they
// moved it, kept whole inside bounds.
func (m *Mover) Place(bounds, unmoved geom.Rect) geom.Rect {
	m.bounds = bounds
	if !m.moved {
		return unmoved
	}
	unmoved.Point = m.inside(bounds, unmoved.Size)
	return unmoved
}

// Moved reports whether the reader has moved the panel since the last Reset.
func (m *Mover) Moved() bool { return m.moved }

// Reset forgets where the reader moved the panel, so Place answers unmoved
// again, as for a surface that opens in the middle every time.
func (m *Mover) Reset() {
	m.moved, m.at, m.dragging = false, geom.Point{}, false
	if parent := m.panel.Parent(); parent != nil {
		parent.MarkForLayoutAndRedraw()
	}
}

// DialogButtonOrder is a surface's own buttons in the order the form reads
// them: on the Windows form, which Linux takes, the actions and then the
// safe button; on the macOS form the safe button and then the actions. A
// Dialog orders its own buttons itself.
func DialogButtonOrder(ui *UI, safe unison.Paneler, actions ...unison.Paneler) []unison.Paneler {
	out := make([]unison.Paneler, 0, len(actions)+1)
	if ui.Form() == FormMac {
		return append(append(out, safe), actions...)
	}
	return append(append(out, actions...), safe)
}
