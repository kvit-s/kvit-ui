package kvitui

import (
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

// spaced is a flex layout whose gap is a Measure, applied on every layout,
// so the gap follows the interface size rather than keeping the value it had
// when the panel was built.
type spaced struct {
	unison.FlexLayout
	ui         *UI
	gap        Measure
	horizontal bool
}

func (s *spaced) apply() {
	g := float32(s.gap.Of(s.ui))
	if s.horizontal {
		s.HSpacing = g
	} else {
		s.VSpacing = g
	}
}

// LayoutSizes implements unison.Layout.
func (s *spaced) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	s.apply()
	return s.FlexLayout.LayoutSizes(target, hint)
}

// PerformLayout implements unison.Layout.
func (s *spaced) PerformLayout(target *unison.Panel) {
	s.apply()
	s.FlexLayout.PerformLayout(target)
}

// Row lays panels out side by side, a gap apart, centred on one line.
func Row(ui *UI, gap Measure, children ...unison.Paneler) *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(&spaced{FlexLayout: unison.FlexLayout{Columns: max(1, len(children)), VAlign: align.Middle}, ui: ui, gap: gap, horizontal: true})
	for _, c := range children {
		if c.AsPanel().LayoutData() == nil {
			c.AsPanel().SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
		}
		p.AddChild(c)
	}
	return p
}

// Column stacks panels top to bottom, a gap apart, each as wide as the
// column.
func Column(ui *UI, gap Measure, children ...unison.Paneler) *unison.Panel {
	p := unison.NewPanel()
	p.SetLayout(&spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: gap})
	for _, c := range children {
		if c.AsPanel().LayoutData() == nil {
			c.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		}
		p.AddChild(c)
	}
	return p
}

// padding is a border that is only space, a Measure wide on every side.
type padding struct {
	ui *UI
	m  Measure
}

func (p padding) Insets() geom.Insets                { return geom.NewUniformInsets(float32(p.m.Of(p.ui))) }
func (p padding) Draw(_ *unison.Canvas, _ geom.Rect) {}

// Padding is space on every side of a panel's content, a Measure wide, that
// follows the interface size. Set it with panel.SetBorder.
func Padding(ui *UI, m Measure) unison.Border { return padding{ui, m} }

// Px is a Measure for a one-off design-pixel value, scaled with the
// interface size. Reach for a named Size… measure first.
func Px(design int) Measure {
	return func(m *tokens.Interface) int { return m.Px(design) }
}

// Width holds a panel to a width, as a Measure, and lets it take the height
// it needs at that width.
func Width(ui *UI, w Measure, child unison.Paneler) *unison.Panel {
	p := unison.NewPanel()
	p.AddChild(child)
	p.SetLayout(&fixedWidth{ui: ui, w: w})
	return p
}

// fixedWidth is the layout Width uses: it reports the measured width and
// lays its one child out to fill it. A panel's layout, when it has one, is
// what unison asks for its size, so the width has to be a layout's answer.
type fixedWidth struct {
	ui *UI
	w  Measure
}

func (f *fixedWidth) LayoutSizes(target *unison.Panel, _ geom.Size) (minSize, prefSize, maxSize geom.Size) {
	width := float32(f.w.Of(f.ui))
	var h float32
	for _, c := range target.Children() {
		_, pref, _ := c.Sizes(geom.NewSize(width, 0))
		h = max(h, pref.Height)
	}
	size := geom.NewSize(width, h)
	return size, size, size
}

func (f *fixedWidth) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	for _, c := range target.Children() {
		c.SetFrameRect(r)
	}
}

// AtLeast wraps a layout so the panel is never shorter than a Measure.
func AtLeast(ui *UI, height Measure, l unison.Layout) unison.Layout {
	return &atLeast{Layout: l, ui: ui, h: height}
}

type atLeast struct {
	unison.Layout
	ui *UI
	h  Measure
}

func (a *atLeast) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	minSize, prefSize, maxSize = a.Layout.LayoutSizes(target, hint)
	h := float32(a.h.Of(a.ui))
	minSize.Height = max(minSize.Height, h)
	prefSize.Height = max(prefSize.Height, h)
	maxSize.Height = max(maxSize.Height, prefSize.Height)
	return minSize, prefSize, maxSize
}

// FullWidth makes a panel as wide as the column or frame it sits in, which
// is what `width: parent.width` says in QML.
func FullWidth[T unison.Paneler](p T) T {
	p.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
	return p
}

// syncing runs sync before every size question and every layout, so a
// component whose child panels follow its fields picks up a field changed
// since the last layout.
type syncing struct {
	unison.Layout
	sync func()
}

func (s syncing) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	s.sync()
	return s.Layout.LayoutSizes(target, hint)
}

func (s syncing) PerformLayout(target *unison.Panel) {
	s.sync()
	s.Layout.PerformLayout(target)
}
