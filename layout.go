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

// Insets is space inside a panel's edges, a Measure on each side; a nil
// Measure is no space. Set it with panel.SetBorder.
func Insets(ui *UI, top, left, bottom, right Measure) unison.Border {
	return insets{ui, top, left, bottom, right}
}

type insets struct {
	ui                       *UI
	top, left, bottom, right Measure
}

func (i insets) Insets() geom.Insets {
	return geom.Insets{Top: orZero(i.ui, i.top), Left: orZero(i.ui, i.left), Bottom: orZero(i.ui, i.bottom), Right: orZero(i.ui, i.right)}
}

func (i insets) Draw(*unison.Canvas, geom.Rect) {}

// Px is a Measure for a one-off design-pixel value, scaled with the
// interface size. Reach for a named Size… measure first.
func Px(design int) Measure {
	return func(m *tokens.Interface) int { return m.Px(design) }
}

// Width holds a panel to a width, as a Measure, and lets it take the height
// it needs at that width.
func Width(ui *UI, w Measure, child unison.Paneler) *unison.Panel {
	return Sized(ui, w, nil, child)
}

// Height holds a panel to a height, as a Measure, and lets it take the width
// its container gives it.
func Height(ui *UI, h Measure, child unison.Paneler) *unison.Panel {
	p := Sized(ui, nil, h, child)
	p.SetLayoutData(nil)
	return p
}

// Sized holds a panel to a width and a height, as Measures; a nil height
// lets it take the height it needs at that width.
func Sized(ui *UI, w, h Measure, child unison.Paneler) *unison.Panel {
	p := unison.NewPanel()
	p.AddChild(child)
	p.SetLayout(&fixedSize{ui: ui, w: w, h: h})
	// A column or row that fills its children would stretch this past the
	// width it holds, so it asks to be left at its own size.
	p.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Start, VAlign: align.Middle})
	return p
}

// fixedSize is the layout Width and Sized use: it reports the measured size
// and lays its one child out to fill it. A panel's layout, when it has one,
// is what unison asks for its size, so a fixed size has to be a layout's
// answer.
type fixedSize struct {
	ui   *UI
	w, h Measure
}

func (f *fixedSize) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	if f.w == nil {
		// Only the height is held: the width is whatever the container gives.
		h := float32(f.h.Of(f.ui))
		var w float32
		for _, c := range target.Children() {
			_, pref, _ := c.Sizes(geom.NewSize(hint.Width, h))
			w = max(w, pref.Width)
		}
		return geom.NewSize(0, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
	}
	width := float32(f.w.Of(f.ui))
	var h float32
	if f.h != nil {
		h = float32(f.h.Of(f.ui))
	} else {
		for _, c := range target.Children() {
			_, pref, _ := c.Sizes(geom.NewSize(width, 0))
			h = max(h, pref.Height)
		}
	}
	size := geom.NewSize(width, h)
	if b := target.Border(); b != nil {
		size = size.Add(b.Insets().Size())
	}
	return size, size, size
}

func (f *fixedSize) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	for _, c := range target.Children() {
		// A child that asks to be centred, as Centred asks, keeps its own
		// size in the middle; any other fills.
		box := r
		if d, ok := c.LayoutData().(*unison.FlexLayoutData); ok {
			_, pref, _ := c.Sizes(geom.NewSize(r.Width, 0))
			if d.HAlign == align.Middle {
				box.Width = min(pref.Width, r.Width)
				box.X = r.X + (r.Width-box.Width)/2
			}
			if d.VAlign == align.Middle {
				box.Height = min(pref.Height, r.Height)
				box.Y = r.Y + (r.Height-box.Height)/2
			}
		}
		c.SetFrameRect(box)
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

// FullWidth makes a panel take the width the row, column or frame it sits in
// has to spare, centred on the height.
func FullWidth[T unison.Paneler](p T) T {
	p.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Middle, HGrab: true})
	return p
}

// At puts a panel at a point inside a panel the width of the container and
// the given height.
func At(ui *UI, x, y, height Measure, child unison.Paneler) *unison.Panel {
	p := unison.NewPanel()
	p.AddChild(child)
	p.SetLayout(atPoint{ui, x, y, height})
	return p
}

type atPoint struct {
	ui      *UI
	x, y, h Measure
}

func (a atPoint) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	h := orZero(a.ui, a.h)
	return geom.NewSize(0, h), geom.NewSize(hint.Width, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (a atPoint) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	for _, c := range target.Children() {
		_, p, _ := c.Sizes(geom.Size{})
		c.SetFrameRect(geom.NewRect(r.X+orZero(a.ui, a.x), r.Y+orZero(a.ui, a.y), p.Width, p.Height))
	}
}

// Left keeps a panel at its own width at the start of the row or column it
// sits in, where a column would otherwise stretch it across.
func Left[T unison.Paneler](p T) T {
	p.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Start, VAlign: align.Middle})
	return p
}

// Centred puts a panel in the middle of the room it is given.
func Centred[T unison.Paneler](p T) T {
	p.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Middle, VAlign: align.Middle, HGrab: true, VGrab: true})
	return p
}

// showOnly makes a panel's children the parts that are not hidden, in the
// order given, changing nothing when they already are. unison's flex layout
// gives a hidden child its room and its gap, so an optional part is taken
// out rather than hidden.
func showOnly(p *unison.Panel, parts ...unison.Paneler) {
	want := make([]*unison.Panel, 0, len(parts))
	for _, c := range parts {
		if !c.AsPanel().Hidden {
			want = append(want, c.AsPanel())
		}
	}
	have := p.Children()
	same := len(have) == len(want)
	for i := 0; same && i < len(want); i++ {
		same = have[i] == want[i]
	}
	if same {
		return
	}
	p.RemoveAllChildren()
	for _, c := range want {
		p.AddChild(c)
	}
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
