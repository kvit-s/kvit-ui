package kvitui

import (
	"math"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// SplitView is two or more regions the reader resizes against each other.
// The handle between two regions is a strip a space wide with a hairline down
// its middle, so the target is comfortable and the rule is still thin, and it
// takes the keyboard focus and moves with the arrow keys: a split a reader can
// only set with a mouse is a split a keyboard user is stuck with.
type SplitView struct {
	unison.Panel
	ui *UI
	// Vertical stacks the regions top to bottom rather than side by side.
	Vertical bool
	// Fill is the region that takes the room the others leave; the last
	// unless set.
	Fill int

	panes   []unison.Paneler
	sizes   []Measure // each region's size along the split; nil takes its own
	handles []*splitHandle
}

// NewSplitView returns the regions side by side, the last taking the room
// the others leave.
func NewSplitView(ui *UI, panes ...unison.Paneler) *SplitView {
	s := &SplitView{ui: ui, Fill: len(panes) - 1, panes: panes, sizes: make([]Measure, len(panes))}
	s.Self = s
	for i, p := range panes {
		s.AddChild(p)
		if i > 0 {
			h := &splitHandle{s: s, before: i - 1}
			h.Self = h
			h.initControl(ui, nil)
			h.ringRadius = func() float32 { return 0 }
			h.MouseDownCallback = h.mouseDown
			h.MouseDragCallback = h.mouseDrag
			h.MouseUpCallback = func(geom.Point, int, mod.Modifiers) bool {
				h.pressed = false
				h.MarkForRedraw()
				return true
			}
			h.KeyDownCallback = h.keyDown
			h.UpdateCursorCallback = func(geom.Point) *unison.Cursor {
				if s.Vertical {
					return unison.ResizeVerticalCursor()
				}
				return unison.ResizeHorizontalCursor()
			}
			h.DrawCallback = h.draw
			s.handles = append(s.handles, h)
			s.AddChild(h)
		}
	}
	s.SetLayout(splitLayout{s})
	return s
}

// SetSize sets a region's size along the split, as the reader's dragging
// does.
func (s *SplitView) SetSize(pane int, size Measure) {
	if pane >= 0 && pane < len(s.sizes) {
		s.sizes[pane] = size
		s.MarkForLayoutAndRedraw()
	}
}

// along is a size's length along the split.
func (s *SplitView) along(size geom.Size) float32 {
	if s.Vertical {
		return size.Height
	}
	return size.Width
}

// natural is a region's size along the split before the fill takes the rest.
func (s *SplitView) natural(i int) float32 {
	if s.sizes[i] != nil {
		return float32(s.sizes[i].Of(s.ui))
	}
	return s.along(preferred(s.panes[i]))
}

type splitLayout struct{ s *SplitView }

func (l splitLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	s := l.s
	grip := float32(s.ui.Interface.Space()) * float32(len(s.handles))
	var along, across float32
	for i, p := range s.panes {
		along += s.natural(i)
		size := preferred(p)
		if s.Vertical {
			across = max(across, size.Width)
		} else {
			across = max(across, size.Height)
		}
	}
	pref := geom.NewSize(along+grip, across)
	if s.Vertical {
		pref = geom.NewSize(across, along+grip)
	}
	return geom.Size{}, pref, geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l splitLayout) PerformLayout(target *unison.Panel) {
	s := l.s
	r := target.ContentRect(false)
	grip := float32(s.ui.Interface.Space())
	total := s.along(r.Size) - grip*float32(len(s.handles))
	fixed := float32(0)
	for i := range s.panes {
		if i != s.Fill {
			fixed += s.natural(i)
		}
	}
	at := float32(0)
	place := func(p *unison.Panel, length float32) {
		if s.Vertical {
			p.SetFrameRect(geom.NewRect(r.X, r.Y+at, r.Width, length))
		} else {
			p.SetFrameRect(geom.NewRect(r.X+at, r.Y, length, r.Height))
		}
		at += length
	}
	for i, p := range s.panes {
		length := s.natural(i)
		if i == s.Fill {
			length = max(0, total-fixed)
		}
		place(p.AsPanel(), length)
		if i < len(s.handles) {
			place(s.handles[i].AsPanel(), grip)
		}
	}
}

// splitHandle is the strip between two regions.
type splitHandle struct {
	control
	s      *SplitView
	before int // the region before the handle
	// Where a drag started, in window coordinates, and the size the region
	// it resizes had then.
	grab, grabSize float32
}

// resized is the region a handle resizes: the one before it, unless that is
// the one filling, in which case the one after it.
func (h *splitHandle) resized() (pane int, sign float32) {
	if h.before == h.s.Fill {
		return h.before + 1, -1
	}
	return h.before, 1
}

// moveBy moves the handle some pixels along the split.
func (h *splitHandle) moveBy(pixels float32) {
	pane, _ := h.resized()
	h.resizeTo(h.s.natural(pane), pixels)
}

// resizeTo gives the region the handle resizes a size, moved some pixels
// from a starting size, kept within the room there is.
func (h *splitHandle) resizeTo(from, pixels float32) {
	s := h.s
	pane, sign := h.resized()
	total := s.along(s.ContentRect(false).Size) - float32(s.ui.Interface.Space())*float32(len(s.handles))
	room := total
	for i := range s.panes {
		if i != s.Fill && i != pane {
			room -= s.natural(i)
		}
	}
	size := max(0, min(room, from+sign*pixels))
	s.SetSize(pane, Px(int(math.Round(float64(size)/s.ui.Interface.Scale()))))
}

func (h *splitHandle) mouseDown(where geom.Point, button, _ int, _ mod.Modifiers) bool {
	if button != unison.ButtonLeft {
		return false
	}
	h.pressed = true
	pane, _ := h.resized()
	h.grab, h.grabSize = h.alongRoot(where), h.s.natural(pane)
	h.MarkForRedraw()
	return true
}

// alongRoot is how far along the split a point in the handle is, in window
// coordinates, which do not move as the handle does.
func (h *splitHandle) alongRoot(where geom.Point) float32 {
	pt := h.PointToRoot(where)
	return h.s.along(geom.NewSize(pt.X, pt.Y))
}

func (h *splitHandle) mouseDrag(where geom.Point, _ int, _ mod.Modifiers) bool {
	if h.pressed {
		h.resizeTo(h.grabSize, h.alongRoot(where)-h.grab)
	}
	return true
}

func (h *splitHandle) keyDown(key unison.KeyCode, _ mod.Modifiers, _ bool) bool {
	step := float32(h.ui.Interface.Px(16))
	switch key {
	case unison.KeyLeft, unison.KeyUp:
		h.moveBy(-step)
	case unison.KeyRight, unison.KeyDown:
		h.moveBy(step)
	default:
		return false
	}
	return true
}

func (h *splitHandle) draw(gc *unison.Canvas, _ geom.Rect) {
	t, m := h.ui.Theme.Tokens(), h.ui.Interface
	r := h.ContentRect(false)
	active := h.hovered || h.pressed || h.KeyboardFocus()
	w, ink := float32(m.Hairline()), t.Border
	if active {
		w, ink = float32(m.SpaceTight()), t.Accent
	}
	var rule geom.Rect
	if h.s.Vertical {
		rule = geom.NewRect(r.X, r.Y+float32(math.Round(float64(r.Height-w)/2)), r.Width, w)
	} else {
		rule = geom.NewRect(r.X+float32(math.Round(float64(r.Width-w)/2)), r.Y, w, r.Height)
	}
	painterFor(gc, h.ui).fill(rule, ink)
}

// ProvideAccessibility describes the handle as a separator a reader can move,
// with the size of the region it resizes as its value.
func (h *splitHandle) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	n.Role = role.Separator
	n.Name = "Resize the regions"
	pane, _ := h.resized()
	n.HasNumber = true
	n.Number = float64(h.s.natural(pane))
	n.Min, n.Max = 0, float64(h.s.along(h.s.ContentRect(false).Size))
	n.Focusable = true
	n.Actions = n.Actions.With(accessibility.Increment, accessibility.Decrement)
}

// PerformAccessibilityAction moves the handle for a screen reader.
func (h *splitHandle) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	step := float32(h.ui.Interface.Px(16))
	switch req.Action {
	case accessibility.Increment:
		h.moveBy(step)
	case accessibility.Decrement:
		h.moveBy(-step)
	default:
		return false
	}
	return true
}
