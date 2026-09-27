package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// Option is one choice in a Segmented control or a Select.
type Option struct {
	// Value is what the choice means to the caller.
	Value string
	// Label is what the segment says.
	Label string
}

// Segmented is one choice from two to five short options, all visible: a
// period (week, month, quarter, year) governing a whole screen, which should
// say what the screen is showing without being opened. Past about five
// options the segments grow too narrow for their words, and a select is the
// right control.
//
// The chosen segment is filled with the accent, the colour that means
// selection, with its label in the colour that contrasts with the fill and
// bold, and no outline of its own. This is kvit-cash's form of the control:
// an outlined segment raised out of the strip read as three buttons rather
// than as one choice of three, and a raised shape is the shape of a control
// being pressed.
type Segmented struct {
	unison.Panel
	ui *UI
	// Label names the choice for a screen reader, such as "Theme".
	Label string
	// Current is the Value of the chosen option.
	Current string
	// OnChoose runs after an option is chosen, with its value.
	OnChoose func(value string)

	segments []*segment
}

type segment struct {
	control
	s      *Segmented
	option Option
}

// NewSegmented returns a choice of options, the first chosen.
func NewSegmented(ui *UI, label string, options ...Option) *Segmented {
	s := &Segmented{ui: ui, Label: label}
	s.Self = s
	for _, o := range options {
		seg := &segment{s: s, option: o}
		seg.Self = seg
		seg.initControl(ui, func() { s.Choose(seg.option.Value) }, unison.KeySpace)
		seg.ringRadius = func() float32 { return float32(ui.Interface.RadiusBar()) }
		seg.SetSizer(seg.sizes)
		seg.DrawCallback = seg.draw
		seg.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Fill, VGrab: true})
		s.segments = append(s.segments, seg)
		s.AddChild(seg)
	}
	if len(options) > 0 {
		s.Current = options[0].Value
	}
	// The strip stands a hairline taller than a control on each side, so a
	// segment, the thing that is pressed, is the full control height: the
	// height of a tab.
	s.SetBorder(Padding(ui, SizeHairline))
	s.SetLayout(AtLeast(ui, SizeTabHeight, &unison.FlexLayout{Columns: max(1, len(options))}))
	s.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		t, m := ui.Theme.Tokens(), ui.Interface
		p := painterFor(gc, ui)
		r := s.ContentRect(true)
		radius := float32(m.RadiusControl())
		p.round(r, radius, t.ChipBackground)
		p.outline(r, radius, float32(m.Hairline()), t.BorderStrong)
	}
	return s
}

// Choose makes the option with a value the chosen one.
func (s *Segmented) Choose(value string) {
	s.Current = value
	s.MarkForRedraw()
	if s.OnChoose != nil {
		s.OnChoose(value)
	}
}

func (g *segment) label() *text.Layout {
	ui, t := g.ui, g.ui.Theme.Tokens()
	weight, ink := text.Regular, t.TextMuted
	if g.chosen() {
		weight, ink = text.Bold, tokens.LabelOn(t.Accent)
	}
	return ui.Fonts.Layout([]text.Span{{Text: g.option.Label, Style: ui.Chrome(ui.Size(RoleBody), weight, ink)}}, text.Options{})
}

func (g *segment) chosen() bool { return g.option.Value == g.s.Current }

func (g *segment) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := g.ui.Interface
	w, _ := g.label().Size()
	h := float32(m.TabHeight() - 2*m.Hairline())
	size := geom.NewSize(max(w+2*float32(m.SpaceLoose()), float32(m.Px(56))), h)
	return size, size, geom.NewSize(size.Width, unison.DefaultMaxSize)
}

func (g *segment) draw(gc *unison.Canvas, _ geom.Rect) {
	t, m := g.ui.Theme.Tokens(), g.ui.Interface
	r := g.ContentRect(false)
	radius := float32(m.RadiusBar())
	switch {
	case g.chosen():
		painterFor(gc, g.ui).round(r, radius, t.Accent)
	case g.hovered:
		painterFor(gc, g.ui).round(r, radius, t.HoverTint)
	}
	l := g.label()
	w, h := l.Size()
	l.Draw(gc, r.X+(r.Width-w)/2, r.Y+(r.Height-h)/2)
}

// ProvideAccessibility names the strip as a list of tabs.
func (s *Segmented) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.TabList
	}
	n.Name = s.Label
}

// ProvideAccessibility describes a segment as a tab that is or is not the
// chosen one.
func (g *segment) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Tab
	}
	n.Name = g.option.Label
	n.Selectable, n.Selected = true, g.chosen()
	if g.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction chooses the segment for a screen reader.
func (g *segment) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !g.Enabled() {
		return false
	}
	g.fire()
	return true
}
