package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// RadioOption is one choice in a RadioGroup, with an optional sentence under
// it saying what it means.
type RadioOption struct {
	Value  string
	Label  string
	Detail string
}

// RadioGroup is one choice from a handful, all visible, each with room for a
// sentence saying what it means: "Follow the system" beside "Always dark"
// beside "Always light", each with a line under it. A Segmented control is
// for two to four short options whose words explain themselves, and a Select
// for a list too long to lay out. The arrow keys move within the group and
// choose, and Tab leaves it, which is what a screen reader user expects of a
// radio group and what a column of separate controls does not do.
type RadioGroup struct {
	unison.Panel
	ui *UI
	// Label names the group for a screen reader.
	Label string
	// Current is the chosen option's value.
	Current string
	// OnChoose runs when the reader chooses an option.
	OnChoose func(value string)

	options []*radioOption
}

type radioOption struct {
	control
	g      *RadioGroup
	option RadioOption
	circle *unison.Panel
}

// NewRadioGroup returns a group of options, none chosen.
func NewRadioGroup(ui *UI, label string, options ...RadioOption) *RadioGroup {
	g := &RadioGroup{ui: ui, Label: label}
	g.Self = g
	for _, o := range options {
		opt := &radioOption{g: g, option: o}
		opt.Self = opt
		opt.initControl(ui, func() { g.Choose(opt.option.Value) }, unison.KeySpace)
		opt.ringRadius = func() float32 { return float32(ui.Interface.Px(16)) / 2 }
		keys := opt.KeyDownCallback
		opt.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
			if g.arrow(opt, key) {
				return true
			}
			return keys(key, mods, repeat)
		}
		opt.circle = unison.NewPanel()
		opt.circle.DrawCallback = opt.drawCircle
		opt.AddChild(opt.circle)
		opt.SetLayout(radioLayout{opt})
		opt.DrawCallback = opt.draw
		opt.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		g.options = append(g.options, opt)
		g.AddChild(opt)
	}
	g.SetLayout(syncing{Layout: &spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: SizeSpaceNear}, sync: g.sync})
	return g
}

// sync makes the chosen option, or the first when none is, the group's one
// stop in the tab order.
func (g *RadioGroup) sync() {
	stop := 0
	for i, o := range g.options {
		if o.option.Value == g.Current {
			stop = i
		}
	}
	for i, o := range g.options {
		o.SetFocusable(i == stop)
	}
}

// Choose makes the option with a value the chosen one.
func (g *RadioGroup) Choose(value string) {
	g.Current = value
	g.sync()
	g.MarkForRedraw()
	if g.OnChoose != nil {
		g.OnChoose(value)
	}
}

// arrow moves to the option before or after, choosing it.
func (g *RadioGroup) arrow(from *radioOption, key unison.KeyCode) bool {
	by := 0
	switch key {
	case unison.KeyUp, unison.KeyLeft:
		by = -1
	case unison.KeyDown, unison.KeyRight:
		by = 1
	default:
		return false
	}
	at := 0
	for i, o := range g.options {
		if o == from {
			at = i
		}
	}
	n := len(g.options)
	for range n {
		at = (at + by + n) % n
		if o := g.options[at]; o.Enabled() {
			g.Choose(o.option.Value)
			o.Focus()
			return true
		}
	}
	return true
}

// focusRing puts the ring around the circle, as the Qt radio button draws it.
func (o *radioOption) focusRing() (*unison.Panel, float32, bool) {
	return o.circle, float32(o.ui.Interface.Px(16)) / 2, o.KeyboardFocus()
}

// indent is where the words start: past the circle and a space.
func (o *radioOption) indent() float32 {
	m := o.ui.Interface
	return float32(m.Px(16) + m.Space())
}

func (o *radioOption) layouts(width float32) (label, detail *text.Layout) {
	ui, t := o.ui, o.ui.Theme.Tokens()
	ink := t.TextPrimary
	if !o.Enabled() {
		ink = t.TextDisabled
	}
	label = ui.Fonts.Layout([]text.Span{{Text: o.option.Label, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ink)}},
		text.Options{MaxWidth: width, Elide: width > 0})
	if o.option.Detail != "" {
		// The detail is a sentence, so it wraps where the label is cut short:
		// cutting the end off an explanation loses the half that explains.
		detail = ui.Fonts.Layout([]text.Span{{Text: o.option.Detail, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, t.TextMuted)}},
			text.Options{MaxWidth: width})
	}
	return label, detail
}

type radioLayout struct{ o *radioOption }

func (l radioLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	o := l.o
	width := float32(0)
	if hint.Width > 0 {
		width = max(1, hint.Width-o.indent())
	}
	label, detail := o.layouts(width)
	w, h := label.Size()
	if detail != nil {
		dw, dh := detail.Size()
		w = max(w, dw)
		h += float32(o.ui.Interface.SpaceTight()) + dh
	}
	h = max(h, float32(o.ui.Interface.Px(16)))
	return geom.NewSize(o.indent(), h), geom.NewSize(o.indent()+w, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l radioLayout) PerformLayout(target *unison.Panel) {
	o := l.o
	r := target.ContentRect(false)
	label, _ := o.layouts(max(1, r.Width-o.indent()))
	_, lh := label.Size()
	s := float32(o.ui.Interface.Px(16))
	// Centred on the first line rather than on the option, which with a
	// line of detail under it is two lines tall.
	o.circle.SetFrameRect(geom.NewRect(r.X, r.Y+(lh-s)/2, s, s))
}

func (o *radioOption) drawCircle(gc *unison.Canvas, _ geom.Rect) {
	t, m := o.ui.Theme.Tokens(), o.ui.Interface
	p := painterFor(gc, o.ui)
	r := o.circle.ContentRect(false)
	edge := t.BorderStrong
	if !o.Enabled() {
		edge = t.Border
	}
	p.outline(r, r.Width/2, float32(m.Hairline()), edge)
	if o.option.Value == o.g.Current {
		dot := r.Inset(geom.NewUniformInsets(float32(m.SpaceNear()) / 2))
		p.round(dot, dot.Width/2, t.Accent)
	}
}

func (o *radioOption) draw(gc *unison.Canvas, _ geom.Rect) {
	r := o.ContentRect(false)
	x := r.X + o.indent()
	label, detail := o.layouts(max(1, r.Width-o.indent()))
	label.Draw(gc, x, r.Y)
	if detail != nil {
		_, lh := label.Size()
		detail.Draw(gc, x, r.Y+lh+float32(o.ui.Interface.SpaceTight()))
	}
}

// ProvideAccessibility names the group.
func (g *RadioGroup) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name = g.Label
}

// ProvideAccessibility reads an option as its label and its detail, and
// whether it is the one chosen.
func (o *radioOption) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.RadioButton
	}
	n.Name = o.option.Label
	if o.option.Detail != "" {
		n.Name = o.option.Label + ". " + o.option.Detail
	}
	n.HasCheck = true
	n.Checked = check.Off
	if o.option.Value == o.g.Current {
		n.Checked = check.On
	}
	if o.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction chooses the option for a screen reader.
func (o *radioOption) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !o.Enabled() {
		return false
	}
	o.fire()
	return true
}
