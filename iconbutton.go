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
	// Command, when set, is what the button runs: its symbol, name,
	// explanation, whether it can act and whether it is on all come from the
	// command, read each time the button is drawn, and on the Windows form
	// its tooltip shows the command's shortcut.
	Command *Command
	// Symbol is the meaning name of the symbol drawn.
	Symbol string
	// Label says what the button does: a verb first, in sentence case, as a
	// tooltip is written on both forms.
	Label string
	// Asks marks a button that asks for more before acting: "…" is added to
	// its tooltip (WRD-4).
	Asks bool
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
	b.pressKeys = buttonKeys
	b.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		s := float32(ui.Interface.ControlHeight())
		if b.Size != nil {
			s = float32(b.Size.Of(ui))
		}
		size := geom.NewSize(s, s)
		return size, size, size
	})
	b.DrawCallback = b.draw
	b.tip = func() string {
		if !b.TooltipEnabled || b.Name() == "" {
			return ""
		}
		if c := b.Command; c != nil {
			b.sync()
			say := c.tooltip(ui, plainWords(c.name(ui)))
			if b.Explanation != "" {
				say = joinLines(say, b.Explanation)
			}
			return say
		}
		label := b.Label
		if b.Asks {
			label += "…"
		}
		return joinLines(label, b.Explanation)
	}
	return b
}

// NewCommandIconButton returns a quiet icon button that runs a command.
func NewCommandIconButton(ui *UI, c *Command) *IconButton {
	b := NewIconButton(ui, "", "")
	b.Command = c
	return b
}

// Name is what the button is called: its label, or its command's name.
func (b *IconButton) Name() string {
	if b.Label == "" && b.Command != nil {
		n := b.Command.Name
		if n == "" {
			n = b.ui.StdName(b.Command.Std)
		}
		return plainWords(n)
	}
	return b.Label
}

// symbol is the meaning name of the symbol drawn.
func (b *IconButton) symbol() string {
	if b.Symbol == "" && b.Command != nil {
		return b.Command.SymbolName()
	}
	return b.Symbol
}

// sync reads the command's state: whether it can act and whether it is on.
func (b *IconButton) sync() {
	if c := b.Command; c != nil {
		b.SetEnabled(c.CanRun())
		b.Checked = c.IsChecked()
	}
}

func (b *IconButton) click() {
	if c := b.Command; c != nil {
		c.Do()
		return
	}
	if b.Checkable {
		b.Checked = !b.Checked
		b.MarkForRedraw()
	}
	if b.OnClick != nil {
		b.OnClick()
	}
}

func (b *IconButton) draw(gc *unison.Canvas, _ geom.Rect) {
	b.sync()
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
	if g, ok := icons.Glyph(b.symbol()); ok {
		drawGlyph(gc, ui, g, float32(size), ink, r)
	}
}

// ProvideAccessibility describes the button: its label as the name, its
// explanation as the description, and its checked state when it has one.
func (b *IconButton) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	b.sync()
	n := bl.Node()
	toggle := b.Checkable || b.Checked || (b.Command != nil && b.Command.Checked != nil)
	if n.Role == role.Auto {
		n.Role = role.Button
		if toggle {
			n.Role = role.ToggleButton
			n.Pressed = b.Checked
		}
	}
	n.Name = b.Name()
	n.Description = b.Explanation
	if b.Command != nil && b.Explanation == "" {
		n.Description = b.Command.Explanation
	}
	if toggle {
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
