package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/role"
)

// IconButtonForm is how an icon button looks at rest.
type IconButtonForm int

const (
	// Quiet has no resting ground and draws one on hover, which is what a
	// dense toolbar wants.
	Quiet IconButtonForm = iota
	// Ordinary keeps its outline the whole time, so its edges are findable.
	Ordinary
)

// IconButton is a button whose whole label is a symbol. Label is required:
// it is what a screen reader announces and what the tooltip shows, so the
// two cannot disagree. Explanation is a second sentence for what the name
// cannot hold, such as why the button is disabled; it follows the label in
// the tooltip and becomes the accessible description.
type IconButton struct {
	control
	// Symbol is the meaning name of the symbol drawn.
	Symbol string
	// Label says what the button does.
	Label string
	// Explanation says more than the label can; optional.
	Explanation string
	// Form is Quiet or Ordinary.
	Form IconButtonForm
	// Dense draws the symbol at the small icon size, the size every symbol
	// beside words is drawn at, rather than the full one.
	Dense bool
	// Checkable makes a press toggle Checked.
	Checkable bool
	// Checked draws the button as on: an accent edge as well as a tint, so
	// the state does not rest on colour alone.
	Checked bool
	// TooltipEnabled shows the tooltip; true unless turned off by a control
	// that already shows the label in a surface of its own.
	TooltipEnabled bool
	// Size is the button's square side; a control's height unless set.
	Size Measure
	// OnClick runs when the button is pressed.
	OnClick func()
}

// NewIconButton returns a quiet icon button.
func NewIconButton(ui *UI, symbol, label string) *IconButton {
	b := &IconButton{Symbol: symbol, Label: label, TooltipEnabled: true}
	b.Self = b
	b.initControl(ui, b.click, unison.KeySpace)
	b.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		s := float32(ui.Interface.ControlHeight())
		if b.Size != nil {
			s = float32(b.Size.Of(ui))
		}
		size := geom.NewSize(s, s)
		return size, size, size
	})
	b.DrawCallback = b.draw
	b.UpdateTooltipCallback = func(geom.Point, geom.Rect) geom.Rect {
		if b.TooltipEnabled && b.Label != "" {
			b.Tooltip = newTooltip(ui, b.Label, b.Explanation)
		} else {
			b.Tooltip = nil
		}
		return b.RectToRoot(b.ContentRect(false))
	}
	return b
}

func (b *IconButton) click() {
	if b.Checkable {
		b.Checked = !b.Checked
		b.MarkForRedraw()
	}
	if b.OnClick != nil {
		b.OnClick()
	}
}

func (b *IconButton) draw(gc *unison.Canvas, _ geom.Rect) {
	ui := b.ui
	t := ui.Theme.Tokens()
	m := ui.Interface
	p := painterFor(gc, ui)
	r := b.ContentRect(false)
	radius := float32(m.RadiusControl())
	if b.Enabled() {
		switch {
		case b.Checked, b.pressed:
			p.round(r, radius, t.SelectionTint)
		case b.hovered:
			p.round(r, radius, t.HoverTint)
		}
	}
	if b.Checked {
		p.outline(r, radius, float32(m.Hairline()), t.Accent)
	} else if b.Form == Ordinary {
		p.outline(r, radius, float32(m.Hairline()), t.BorderStrong)
	}
	ink := t.TextSecondary
	switch {
	case b.Checked:
		ink = t.Accent
	case !b.Enabled():
		ink = t.TextDisabled
	}
	size := m.IconSize()
	if b.Dense {
		size = m.IconSizeSmall()
	}
	if g, ok := icons.Glyph(b.Symbol); ok {
		drawGlyph(gc, ui, g, float32(size), ink, r)
	}
}

// ProvideAccessibility describes the button: its label as the name, its
// explanation as the description, and its checked state when it has one.
func (b *IconButton) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	n := bl.Node()
	if n.Role == role.Auto {
		n.Role = role.Button
		if b.Checkable || b.Checked {
			n.Role = role.ToggleButton
			n.Pressed = b.Checked
		}
	}
	n.Name = b.Label
	n.Description = b.Explanation
	if b.Checkable || b.Checked {
		n.HasCheck = true
		n.Checked = check.Off
		if b.Checked {
			n.Checked = check.On
		}
	}
	if b.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction presses the button for a screen reader.
func (b *IconButton) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !b.Enabled() {
		return false
	}
	b.click()
	return true
}
