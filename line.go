package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// Line is a panel whose children are laid out by a LineLayout: side by side,
// a gap apart, centred on the line's height. Unlike a Row, a hidden child
// takes no room and no gap.
func Line(ui *UI, gap Measure, children ...unison.Paneler) *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(NewLineLayout(ui, gap))
	for _, c := range children {
		p.AddChild(c)
	}
	return p
}

// LineLayout lays the children that are not hidden out side by side, a gap
// apart. A child marked with Grow takes its share of the width the others
// leave, one marked with FixedWidth is held to that width, and any other
// keeps its own width, measured with no width offered. Each child is then
// asked for its height at the width it got, cut to the line's height when it
// is taller, and centred on the line's content height.
//
// The line asks for the height of its tallest child at the width that child
// would get, and for the sum of its children's own widths. It does not read
// unison's flex layout data.
//
// LayoutSizes and PerformLayout are exported, so a layout of an application's
// own can call them and add to them.
type LineLayout struct {
	ui *UI
	// Gap is the room between one child and the next.
	Gap Measure
	// GrowLast makes the last shown child grow when no child is marked with
	// Grow.
	GrowLast bool
	// ClipFixed cuts a child that does not grow to the room left in the
	// line, and measures it again at that width.
	ClipFixed bool
	// FillsHint makes the width asked for the sum of the widths the children
	// would get in the width offered, so a line with a grower asks for the
	// whole width it is offered.
	FillsHint bool
	// InsetsUncounted leaves the panel's insets out of the sizes reported,
	// and uses the width offered as it comes. Children are still placed
	// inside the insets.
	InsetsUncounted bool
	// Top centres the children on the height of the tallest of them,
	// measured from the top of the content, and does not cut them to the
	// content's height.
	Top bool
	// NoShrink makes the least size the line reports the size it asks for,
	// and lets a child keep its own height in a shorter line.
	NoShrink bool
}

// NewLineLayout returns a line's layout with a gap between the children.
func NewLineLayout(ui *UI, gap Measure) *LineLayout { return &LineLayout{ui: ui, Gap: gap} }

// LineData is a line child's layout data, set with SetLayoutData, by value or
// by pointer.
type LineData struct {
	// Grow, above zero, gives the child this share of the width left over.
	Grow float32
	// Width, when set, holds the child to this width.
	Width Measure
}

// Grow makes a panel take the width the line it sits in has left over.
func Grow[T unison.Paneler](p T) T { return GrowShare(p, 1) }

// GrowShare makes a panel take a share of the width the line it sits in has
// left over, against the shares of the line's other growers.
func GrowShare[T unison.Paneler](p T, share float32) T {
	p.AsPanel().SetLayoutData(LineData{Grow: share})
	return p
}

// FixedWidth holds a panel to a width in the line it sits in.
func FixedWidth[T unison.Paneler](p T, w Measure) T {
	p.AsPanel().SetLayoutData(LineData{Width: w})
	return p
}

func lineDataOf(c *unison.Panel) LineData {
	switch d := c.LayoutData().(type) {
	case LineData:
		return d
	case *LineData:
		if d != nil {
			return *d
		}
	}
	return LineData{}
}

// widths are the shown children, the width each gets in a content width, 0
// for none offered, and the width each would take by itself.
func (l *LineLayout) widths(target *unison.Panel, width float32) (parts []*unison.Panel, got, own []float32) {
	parts = shownIn(target)
	gap := orZero(l.ui, l.Gap)
	got = make([]float32, len(parts))
	own = make([]float32, len(parts))
	shares := make([]float32, len(parts))
	used := gap * float32(max(0, len(parts)-1))
	var total float32
	for i, c := range parts {
		d := lineDataOf(c)
		switch {
		case d.Grow > 0:
			shares[i] = d.Grow
			total += d.Grow
			_, p, _ := c.Sizes(geom.Size{})
			own[i] = p.Width
		case d.Width != nil:
			own[i] = orZero(l.ui, d.Width)
		default:
			_, p, _ := c.Sizes(geom.Size{})
			own[i] = p.Width
		}
	}
	if total == 0 && l.GrowLast && len(parts) > 0 {
		shares[len(parts)-1], total = 1, 1
	}
	for i := range parts {
		if shares[i] == 0 {
			got[i] = own[i]
			used += own[i]
		}
	}
	x := float32(0)
	for i := range parts {
		switch {
		case shares[i] > 0 && width > 0:
			got[i] = max(0, width-used) * shares[i] / total
		case shares[i] > 0:
			got[i] = own[i]
		case l.ClipFixed && width > 0:
			got[i] = min(got[i], max(0, width-x))
		}
		x += got[i] + gap
	}
	return parts, got, own
}

// LayoutSizes implements unison.Layout.
func (l *LineLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	in := insetsOf(target)
	if l.InsetsUncounted {
		in = geom.Insets{}
	}
	width := hint.Width
	if width > 0 {
		width = max(0, width-in.Width())
	}
	parts, got, own := l.widths(target, width)
	var w, h float32
	for i, c := range parts {
		_, p, _ := c.Sizes(geom.NewSize(got[i], 0))
		h = max(h, p.Height)
		if l.FillsHint {
			w += got[i]
		} else {
			w += own[i]
		}
	}
	w += orZero(l.ui, l.Gap)*float32(max(0, len(parts)-1)) + in.Width()
	h += in.Height()
	minSize = geom.NewSize(in.Width(), h)
	if l.NoShrink {
		minSize.Width = w
	}
	return minSize, geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
}

// PerformLayout implements unison.Layout.
func (l *LineLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	parts, got, _ := l.widths(target, r.Width)
	gap := orZero(l.ui, l.Gap)
	heights := make([]float32, len(parts))
	var tallest float32
	for i, c := range parts {
		_, p, _ := c.Sizes(geom.NewSize(got[i], 0))
		heights[i] = p.Height
		tallest = max(tallest, p.Height)
	}
	x := r.X
	for i, c := range parts {
		h, within := heights[i], r.Height
		switch {
		case l.Top:
			within = tallest
		case !l.NoShrink:
			h = min(h, r.Height)
		}
		c.SetFrameRect(geom.NewRect(x, r.Y+(within-h)/2, got[i], h))
		x += got[i] + gap
	}
}
