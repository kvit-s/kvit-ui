package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// Select is a choice from a list too long to lay out: a font family, a
// currency, a folder. For a few mutually exclusive options worth seeing at
// once, Segmented is right: a three-item list hides two of the answers behind
// a press. The list opens as unison's menu, as the owner chose for menus, with
// the chosen option ticked; Up and Down change the choice without opening it.
type Select struct {
	control
	// Label names the choice for a screen reader.
	Label string
	// Options are the choices.
	Options []Option
	// Current is the Value of the chosen option.
	Current string
	// OnChoose runs after an option is chosen, with its value.
	OnChoose func(value string)
}

// NewSelect returns a choice of options, the first chosen.
func NewSelect(ui *UI, label string, options ...Option) *Select {
	s := &Select{Label: label, Options: options}
	s.Self = s
	if len(options) > 0 {
		s.Current = options[0].Value
	}
	s.initControl(ui, s.open, unison.KeySpace, unison.KeyReturn, unison.KeyNumPadEnter)
	keyDown := s.KeyDownCallback
	s.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		switch key {
		case unison.KeyUp:
			return s.step(-1)
		case unison.KeyDown:
			if mods.OptionDown() {
				s.open()
				return true
			}
			return s.step(1)
		}
		return keyDown(key, mods, repeat)
	}
	s.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		m := ui.Interface
		h := float32(m.ControlHeight())
		return geom.NewSize(float32(m.Px(60)), h), geom.NewSize(float32(m.Px(180)), h), geom.NewSize(unison.DefaultMaxSize, h)
	})
	s.DrawCallback = s.draw
	return s
}

// index is the position of the chosen option, or -1.
func (s *Select) index() int {
	for i, o := range s.Options {
		if o.Value == s.Current {
			return i
		}
	}
	return -1
}

// Choose makes the option with a value the chosen one.
func (s *Select) Choose(value string) {
	s.Current = value
	s.MarkForRedraw()
	if s.OnChoose != nil {
		s.OnChoose(value)
	}
}

func (s *Select) step(by int) bool {
	if len(s.Options) == 0 || !s.Enabled() {
		return false
	}
	i := max(0, min(len(s.Options)-1, s.index()+by))
	if s.Options[i].Value != s.Current {
		s.Choose(s.Options[i].Value)
	}
	return true
}

func (s *Select) open() {
	var items []MenuItem
	for _, o := range s.Options {
		items = append(items, MenuItem{Text: o.Label, Checked: o.Value == s.Current, OnSelect: func() { s.Choose(o.Value) }})
	}
	s.ui.ShowMenu(s, s.Label, items)
}

// shown is the chosen option's label.
func (s *Select) shown() string {
	if i := s.index(); i >= 0 {
		return s.Options[i].Label
	}
	return ""
}

func (s *Select) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := s.ui, s.ui.Theme.Tokens(), s.ui.Interface
	p := painterFor(gc, ui)
	r := s.ContentRect(false)
	radius := float32(m.RadiusControl())
	ground, ink := t.PopupBackground, t.TextPrimary
	if !s.Enabled() {
		ground, ink = t.ChipBackground, t.TextDisabled
	}
	p.round(r, radius, ground)
	edge := t.BorderStrong
	if s.Focused() {
		edge = t.FocusRing
	}
	p.outline(r, radius, float32(m.Hairline()), edge)
	// The indicator at the small icon size: at the full size it overpowers
	// the words beside it.
	c := float32(m.IconSizeSmall())
	near := float32(m.SpaceNear())
	if g, ok := icons.Glyph("chevron-down"); ok {
		drawGlyph(gc, ui, g, c, t.TextMuted, geom.NewRect(r.Right()-near-c, r.Y+(r.Height-c)/2, c, c))
	}
	room := r.Width - near - c - float32(m.SpaceLoose()) - near
	l := ui.Fonts.Layout([]text.Span{{Text: s.shown(), Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ink)}},
		text.Options{MaxWidth: max(0, room), Elide: true})
	_, h := l.Size()
	l.Draw(gc, r.X+near, r.Y+(r.Height-h)/2)
}

// ProvideAccessibility describes the select as a combo box with its choice
// as its value.
func (s *Select) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.ComboBox
	}
	n.Name, n.Value = s.Label, s.shown()
	n.Expandable = true
	if s.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction opens the list for a screen reader.
func (s *Select) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !s.Enabled() {
		return false
	}
	s.open()
	return true
}
