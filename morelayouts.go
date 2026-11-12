package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// FillLayout gives every child the whole of the panel's content, one over
// another, and asks for nothing.
func FillLayout() unison.Layout { return fillLayout{} }

type fillLayout struct{}

func (fillLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	return geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (fillLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	for _, c := range target.Children() {
		c.SetFrameRect(r)
	}
}

// OwnSizeLayout gives a panel's first child the whole of its content and asks
// for what the child asks. A hidden child keeps its room, since its sizes are
// still asked, and a panel with no children asks for nothing.
type OwnSizeLayout struct {
	// CountInsets adds the panel's insets to what the child asks, and takes
	// their width off the width the child is asked at.
	CountInsets bool
	// HiddenGone makes a hidden child take no room.
	HiddenGone bool
}

// LayoutSizes implements unison.Layout.
func (o OwnSizeLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	if len(target.Children()) == 0 {
		return geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
	}
	c := target.Children()[0]
	if o.HiddenGone && c.Hidden {
		return geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, 0)
	}
	if !o.CountInsets {
		return c.Sizes(hint)
	}
	in := insetsOf(target)
	inner := hint
	if inner.Width > 0 {
		inner.Width = max(0, inner.Width-in.Width())
	}
	minSize, prefSize, maxSize = c.Sizes(inner)
	grow := geom.NewSize(in.Width(), in.Height())
	return minSize.Add(grow), prefSize.Add(grow), maxSize.Add(grow)
}

// PerformLayout implements unison.Layout.
func (o OwnSizeLayout) PerformLayout(target *unison.Panel) {
	if len(target.Children()) > 0 {
		target.Children()[0].SetFrameRect(target.ContentRect(false))
	}
}

// AtMost wraps a layout so the panel is never taller than a Measure: it is
// as tall as the wrapped layout asks at the width offered, up to the limit,
// and neither shorter nor taller than that. Widths pass through unchanged.
func AtMost(ui *UI, height Measure, l unison.Layout) unison.Layout {
	return &atMost{Layout: l, ui: ui, h: height}
}

type atMost struct {
	unison.Layout
	ui *UI
	h  Measure
}

func (a *atMost) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	minSize, prefSize, maxSize = a.Layout.LayoutSizes(target, geom.NewSize(hint.Width, 0))
	h := min(prefSize.Height, float32(a.h.Of(a.ui)))
	minSize.Height, prefSize.Height, maxSize.Height = h, h, h
	return minSize, prefSize, maxSize
}

// FlowLayout lays the children that are not hidden out left to right, a gap
// apart, starting a new line where the next does not fit. The same gap
// separates one line from the next. Each child keeps its own width, cut to
// the line's width, and its own height.
type FlowLayout struct {
	ui *UI
	// Gap is the room between children in a line and between lines.
	Gap Measure
	// Unasked is the width to assume when none is offered; nil assumes the
	// width of every child on one line.
	Unasked Measure
	// FillsHint makes the width asked for the width offered or assumed,
	// where it is otherwise the width of every child on one line.
	FillsHint bool
	// ItemWidth, when set, answers the width every child is measured at and
	// framed to, from the line's width, in place of the child's own.
	ItemWidth func(lineWidth float32) float32
}

// NewFlowLayout returns a wrapping layout with a gap between the children.
func NewFlowLayout(ui *UI, gap Measure) *FlowLayout { return &FlowLayout{ui: ui, Gap: gap} }

// place puts the children in lines no wider than width, 0 for no limit,
// from origin, and returns the height of the lot.
func (f *FlowLayout) place(target *unison.Panel, origin geom.Point, width float32, set bool) float32 {
	gap := orZero(f.ui, f.Gap)
	var x, y, lineHeight float32
	for _, c := range shownIn(target) {
		var w float32
		var p geom.Size
		if f.ItemWidth != nil {
			w = f.ItemWidth(width)
			_, p, _ = c.Sizes(geom.NewSize(w, 0))
		} else {
			_, p, _ = c.Sizes(geom.Size{})
			w = p.Width
			if width > 0 {
				w = min(w, width)
			}
		}
		if x > 0 && x+w > width && width > 0 {
			x = 0
			y += lineHeight + gap
			lineHeight = 0
		}
		if set {
			c.SetFrameRect(geom.NewRect(origin.X+x, origin.Y+y, w, p.Height))
		}
		x += w + gap
		lineHeight = max(lineHeight, p.Height)
	}
	return y + lineHeight
}

// LayoutSizes implements unison.Layout.
func (f *FlowLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	in := insetsOf(target)
	gap := orZero(f.ui, f.Gap)
	var oneLine float32
	for i, c := range shownIn(target) {
		_, p, _ := c.Sizes(geom.Size{})
		if i > 0 {
			oneLine += gap
		}
		oneLine += p.Width
	}
	width := hint.Width - in.Width()
	if hint.Width <= 0 {
		width = oneLine
		if f.Unasked != nil {
			width = float32(f.Unasked.Of(f.ui))
		}
	}
	h := f.place(target, geom.Point{}, width, false) + in.Height()
	w := oneLine
	if f.FillsHint {
		w = max(0, width)
	}
	return geom.NewSize(in.Width(), h), geom.NewSize(w+in.Width(), h), geom.NewSize(unison.DefaultMaxSize, h)
}

// PerformLayout implements unison.Layout.
func (f *FlowLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	f.place(target, r.Point, r.Width, true)
}

// Relayout asks for a panel and everything it sits in to be laid out again
// and drawn; with children, everything inside it too.
func Relayout(p unison.Paneler, children bool) {
	panel := p.AsPanel()
	if children {
		panel.MarkForLayoutRecursively()
	}
	panel.MarkForLayoutRecursivelyUpward()
	panel.MarkForRedraw()
}
