package kvitui

import (
	"slices"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// shownIn are a panel's children that are not hidden, in order.
func shownIn(target *unison.Panel) []*unison.Panel {
	var out []*unison.Panel
	for _, c := range target.Children() {
		if !c.Hidden {
			out = append(out, c)
		}
	}
	return out
}

// insetsOf are the insets of a panel's border, none without one.
func insetsOf(target *unison.Panel) geom.Insets {
	if b := target.Border(); b != nil {
		return b.Insets()
	}
	return geom.Insets{}
}

// Stack is a panel whose children are laid out by a StackLayout: top to
// bottom, a gap apart, each as wide as the stack. Unlike a Column, a hidden
// child takes no room and no gap, so a part comes and goes by its Hidden
// field alone.
func Stack(ui *UI, gap Measure, children ...unison.Paneler) *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(NewStackLayout(ui, gap))
	for _, c := range children {
		p.AddChild(c)
	}
	return p
}

// StackLayout lays the children that are not hidden out top to bottom, a gap
// apart. A hidden child gets no room and no gap, and its frame is left as it
// was. Every other child is asked for its size at the stack's content width
// and gets that whole width and the height it asks for there; the stack
// neither stretches nor shrinks a child's height.
//
// A child says anything else about itself with StackData. The stack does not
// read unison's flex layout data, which Width, Sized and Centred set on the
// panels they make: such a panel in a stack gets the stack's whole width.
//
// LayoutSizes and PerformLayout are exported, so a layout of an application's
// own can call them and add to them.
type StackLayout struct {
	ui *UI
	// Gap is the room between one child and the next.
	Gap Measure
	// MaxWidth holds each child to a width, at the stack's start, except a
	// child whose data says Wide; nil is no limit. With it set, a child's
	// Own and End are not consulted.
	MaxWidth Measure
	// SkipEmpty also leaves out a shown child that asks for no height, with
	// no gap for it, and gives every child left out, hidden ones included, a
	// frame of no height just below the last child placed. A child left at
	// the frame it had lies over what is laid out there now and takes its
	// presses.
	SkipEmpty bool
}

// NewStackLayout returns a stack's layout with a gap between the children.
func NewStackLayout(ui *UI, gap Measure) *StackLayout { return &StackLayout{ui: ui, Gap: gap} }

// StackData is a stack child's layout data, set with SetLayoutData, by value
// or by pointer.
type StackData struct {
	// Own keeps the child at its own width, at the stack's start.
	Own bool
	// End, with Own, puts the child at the stack's end instead.
	End bool
	// Wide exempts the child from the stack's MaxWidth.
	Wide bool
	// Left and Right are room kept clear before and after the child: it is
	// measured at, and framed to, the stack's width less both. Nil is none.
	Left, Right Measure
}

// OwnWidth keeps a panel at its own width at the start of the stack it sits
// in.
func OwnWidth[T unison.Paneler](p T) T {
	p.AsPanel().SetLayoutData(StackData{Own: true})
	return p
}

// OwnWidthEnd keeps a panel at its own width at the end of the stack it sits
// in.
func OwnWidthEnd[T unison.Paneler](p T) T {
	p.AsPanel().SetLayoutData(StackData{Own: true, End: true})
	return p
}

func stackDataOf(c *unison.Panel) StackData {
	switch d := c.LayoutData().(type) {
	case StackData:
		return d
	case *StackData:
		if d != nil {
			return *d
		}
	}
	return StackData{}
}

// stackPart is one child a stack places and the box it gets, from the
// top-left corner of the stack's content.
type stackPart struct {
	panel *unison.Panel
	box   geom.Rect
}

// arrange is where each child goes in a content width, 0 for none offered,
// and the height and the widest request of the lot. Laying out for real also
// narrows a child that keeps its own width.
func (s *StackLayout) arrange(target *unison.Panel, width float32, real bool) (parts []stackPart, widest, height float32) {
	gap := orZero(s.ui, s.Gap)
	for _, c := range shownIn(target) {
		d := stackDataOf(c)
		left, right := orZero(s.ui, d.Left), orZero(s.ui, d.Right)
		x, w := left, width
		if left != 0 || right != 0 {
			w = max(0, width-left-right)
		}
		if s.MaxWidth != nil && !d.Wide && (w > 0 || real) {
			w = min(w, orZero(s.ui, s.MaxWidth))
		}
		_, p, _ := c.Sizes(geom.NewSize(w, 0))
		if s.SkipEmpty && p.Height <= 0 {
			continue
		}
		if real && s.MaxWidth == nil && d.Own {
			own := min(p.Width, w)
			if d.End {
				x += w - own
			}
			w = own
		}
		if len(parts) > 0 {
			height += gap
		}
		parts = append(parts, stackPart{panel: c, box: geom.NewRect(x, height, w, p.Height)})
		widest = max(widest, left+p.Width+right)
		height += p.Height
	}
	return parts, widest, height
}

// LayoutSizes implements unison.Layout.
func (s *StackLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	in := insetsOf(target)
	width := hint.Width
	if width > 0 {
		width = max(0, width-in.Width())
	}
	_, w, h := s.arrange(target, width, false)
	size := geom.NewSize(w+in.Width(), h+in.Height())
	return geom.NewSize(in.Width(), size.Height), size, geom.NewSize(unison.DefaultMaxSize, size.Height)
}

// PerformLayout implements unison.Layout.
func (s *StackLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	parts, _, height := s.arrange(target, r.Width, true)
	placed := make([]*unison.Panel, 0, len(parts))
	for _, p := range parts {
		box := p.box
		box.X += r.X
		box.Y += r.Y
		p.panel.SetFrameRect(box)
		placed = append(placed, p.panel)
	}
	if !s.SkipEmpty {
		return
	}
	for _, c := range target.Children() {
		if !slices.Contains(placed, c) {
			c.SetFrameRect(geom.NewRect(r.X, r.Y+height, r.Width, 0))
		}
	}
}
