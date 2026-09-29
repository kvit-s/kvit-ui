package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/role"
)

// Check is a checkbox, in three states. The third, partial, is what a parent
// row shows when some of its children are checked and some are not: drawn as
// unchecked it loses that anything under it is selected, and drawn as
// checked it says something false that the reader acts on.
type Check struct {
	control
	// Text is the box's label.
	Text string
	// Checked is the box's state.
	Checked bool
	// Partial says some of what the box stands for is checked; it is drawn
	// with a dash, and a press checks the whole.
	Partial bool
	// OnChange runs after a press, with the new state.
	OnChange func(checked bool)
	// Name is what a screen reader calls a box with no words beside it, such
	// as the one on each row of a table; the words unless set.
	Name string

	box *unison.Panel
}

// NewCheck returns an unchecked box with a label.
func NewCheck(ui *UI, label string) *Check {
	c := &Check{Text: label}
	c.Self = c
	c.initControl(ui, c.toggle, unison.KeySpace)
	c.box = unison.NewPanel()
	c.box.DrawCallback = c.drawBox
	c.AddChild(c.box)
	c.SetLayout(checkLayout{c})
	c.DrawCallback = c.drawLabel
	return c
}

func (c *Check) toggle() {
	c.Checked = c.Partial || !c.Checked
	c.Partial = false
	c.MarkForRedraw()
	if c.OnChange != nil {
		c.OnChange(c.Checked)
	}
}

// focusRing puts the ring around the box, not the label, as the check
// draws it.
func (c *Check) focusRing() (*unison.Panel, float32, bool) {
	return c.box, float32(c.ui.Interface.RadiusBar()), c.KeyboardFocus()
}

func (c *Check) label() *text.Layout {
	ui, t := c.ui, c.ui.Theme.Tokens()
	ink := t.TextPrimary
	if !c.Enabled() {
		ink = t.TextDisabled
	}
	return ui.Fonts.Layout([]text.Span{{Text: c.Text, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ink)}}, text.Options{})
}

// checkPadding is the space on each side of a check box, the padding
// Quick's own style gives its CheckBox, which KvitCheck keeps. It is six
// pixels at every interface size, as the style's is.
const checkPadding = 6

type checkLayout struct{ c *Check }

func (l checkLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.c.ui.Interface
	w, h := l.c.label().Size()
	size := geom.NewSize(float32(2*checkPadding+m.Px(16)+m.SpaceNear())+w, max(float32(m.ControlHeight()), h))
	return size, size, geom.NewSize(unison.DefaultMaxSize, size.Height)
}

func (l checkLayout) PerformLayout(target *unison.Panel) {
	s := float32(l.c.ui.Interface.Px(16))
	r := target.ContentRect(false)
	l.c.box.SetFrameRect(geom.NewRect(r.X+checkPadding, r.Y+(r.Height-s)/2, s, s))
}

func (c *Check) drawBox(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := c.ui, c.ui.Theme.Tokens(), c.ui.Interface
	p := painterFor(gc, ui)
	r := c.box.ContentRect(false)
	radius := float32(m.RadiusBar())
	on := c.Checked || c.Partial
	if on && c.Enabled() {
		p.round(r, radius, t.Accent)
	}
	edge := t.BorderStrong
	if !c.Enabled() {
		edge = t.Border
	}
	p.outline(r, radius, float32(m.Hairline()), edge)
	if on {
		symbol := "check"
		if c.Partial {
			symbol = "minus"
		}
		// A disabled box has no fill, so its mark takes the disabled text
		// colour rather than the colour made for the accent fill.
		ink := tokens.LabelOn(t.Accent)
		if !c.Enabled() {
			ink = t.TextDisabled
		}
		if g, ok := icons.Glyph(symbol); ok {
			drawGlyph(gc, ui, g, r.Width-float32(m.SpaceSnug()), ink, r)
		}
	}
}

func (c *Check) drawLabel(gc *unison.Canvas, _ geom.Rect) {
	m := c.ui.Interface
	r := c.ContentRect(false)
	l := c.label()
	_, h := l.Size()
	l.Draw(gc, r.X+float32(checkPadding+m.Px(16)+m.SpaceNear()), r.Y+(r.Height-h)/2)
}

// ProvideAccessibility describes the box as a check box in one of three
// states.
func (c *Check) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.CheckBox
	}
	n.Name = c.Text
	if c.Name != "" {
		n.Name = c.Name
	}
	n.HasCheck = true
	switch {
	case c.Partial:
		n.Checked = check.Mixed
	case c.Checked:
		n.Checked = check.On
	default:
		n.Checked = check.Off
	}
	if c.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction toggles the box for a screen reader.
func (c *Check) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !c.Enabled() {
		return false
	}
	c.fire()
	return true
}
