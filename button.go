package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/role"
)

// ButtonForm is how a button with words on it looks.
type ButtonForm int

const (
	// ButtonOrdinary is outlined: every action that is not the one the
	// screen is for.
	ButtonOrdinary ButtonForm = iota
	// ButtonPrimary is filled with the accent: the action the screen is for,
	// one per region of the screen. A screen with three has told the reader
	// nothing about which to press.
	ButtonPrimary
	// ButtonQuiet has the chip ground and a hairline of the faint border, for
	// a dense strip where a row of full outlines would read as a fence. It is
	// drawn as pressable at rest rather than only once the pointer is on it,
	// as kvit-cash's copy of the library draws it.
	ButtonQuiet
	// ButtonFlat has no ground and no outline at rest, only the hover tint
	// under the pointer and the selection tint while pressed or checked, for
	// a toolbar of many buttons told apart by their words, as Kvit's
	// toolbar draws them (Toolbar, BarBackground).
	ButtonFlat
)

// Button is a button with words on it. Danger is separate from the form,
// because a destructive action can take any of the three: a primary Delete
// in a confirmation, an ordinary one in a settings row, a quiet one in a
// strip. The words on a filled button take the colour that contrasts with
// its own fill, so a filled danger button stays readable.
type Button struct {
	control
	// Text is what the button says.
	Text string
	// Form is ordinary, primary or quiet.
	Form ButtonForm
	// Danger draws the button in the danger colour.
	Danger bool
	// Symbol is an optional meaning name drawn before the words.
	Symbol string
	// Busy says the button is doing something: it stays enabled so its name
	// can say what, and the words change to BusyText rather than a spinner
	// appearing, since a spinner says "wait" without saying for what.
	Busy bool
	// BusyText is what a busy button says; "Working…" unless set.
	BusyText string
	// Explanation is one sentence the words cannot say: why the button is
	// disabled, or what pressing it opens. It is the tooltip and the
	// accessible description, and a disabled button still shows it.
	Explanation string
	// Checkable makes a press toggle Checked.
	Checkable bool
	// Checked draws a button that turns a mode on as on: tinted and outlined
	// in the accent, two marks rather than one, since a tint alone is a few
	// percent of lightness in the high-contrast theme.
	Checked bool
	// OnClick runs when the button is pressed.
	OnClick func()
}

// NewButton returns an ordinary button.
func NewButton(ui *UI, label string) *Button {
	b := &Button{Text: label, BusyText: "Working…"}
	b.Self = b
	b.initControl(ui, b.click, unison.KeySpace)
	b.SetSizer(b.sizes)
	b.DrawCallback = b.draw
	b.tip = func() string { return b.Explanation }
	return b
}

func (b *Button) click() {
	if b.Checkable {
		b.Checked = !b.Checked
		b.MarkForRedraw()
	}
	if b.OnClick != nil {
		b.OnClick()
	}
}

// words is what the button says now.
func (b *Button) words() string {
	if b.Busy {
		return b.BusyText
	}
	return b.Text
}

func (b *Button) fill(t tokens.Tokens) palette.Color {
	if b.Danger {
		return t.Danger
	}
	return t.Accent
}

func (b *Button) foreground(t tokens.Tokens) palette.Color {
	switch {
	case !b.Enabled():
		return t.TextDisabled
	case b.Form == ButtonPrimary:
		return tokens.LabelOn(b.fill(t))
	case b.Danger:
		return t.Danger
	}
	return t.TextPrimary
}

func (b *Button) label() *text.Layout {
	ui := b.ui
	st := ui.Chrome(ui.Size(RoleBody), text.Regular, b.foreground(ui.Theme.Tokens()))
	return ui.Fonts.Layout([]text.Span{{Text: b.words(), Style: st}}, text.Options{})
}

// content is the width and height of the symbol and the words, a near space
// apart.
func (b *Button) content() (width, height float32, l *text.Layout) {
	m := b.ui.Interface
	l = b.label()
	width, height = l.Size()
	if b.Symbol != "" {
		s := float32(m.IconSizeSmall())
		width += s + float32(m.SpaceNear())
		height = max(height, s)
	}
	return width, height, l
}

func (b *Button) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := b.ui.Interface
	w, h, _ := b.content()
	size := geom.NewSize(w+2*float32(m.SpaceLoose()), max(float32(m.ControlHeight()), h))
	return size, size, size
}

func (b *Button) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := b.ui, b.ui.Theme.Tokens(), b.ui.Interface
	p := painterFor(gc, ui)
	r := b.ContentRect(false)
	radius := float32(m.RadiusControl())
	primary := b.Form == ButtonPrimary
	var ground *palette.Color
	set := func(c palette.Color) { ground = &c }
	switch {
	case !b.Enabled():
		if primary {
			set(t.ChipBackground)
		}
	case primary && b.pressed:
		set(b.fill(t).Darker(1.15))
	case primary:
		set(b.fill(t))
	case b.pressed, b.Checked:
		set(t.SelectionTint)
	case b.hovered:
		set(t.HoverTint)
	case b.Form == ButtonQuiet:
		set(t.ChipBackground)
	}
	if ground != nil {
		p.round(r, radius, *ground)
	}
	if !primary && (b.Form != ButtonFlat || b.Checked) {
		edge := t.BorderStrong
		switch {
		case !b.Enabled():
			edge = t.Border
		case b.Checked:
			edge = t.Accent
		case b.Danger:
			edge = t.Danger
		case b.Form == ButtonQuiet:
			edge = t.Border
		}
		p.outline(r, radius, float32(m.Hairline()), edge)
	}
	w, _, l := b.content()
	x := r.X + (r.Width-w)/2
	if b.Symbol != "" {
		s := float32(m.IconSizeSmall())
		if g, ok := icons.Glyph(b.Symbol); ok {
			drawGlyph(gc, ui, g, s, b.foreground(t), geom.NewRect(x, r.Y+(r.Height-s)/2, s, s))
		}
		x += s + float32(m.SpaceNear())
	}
	_, h := l.Size()
	l.Draw(gc, x, r.Y+(r.Height-h)/2)
}

// ProvideAccessibility describes the button by its words, what it is doing
// when busy, its explanation, and its checked state when it has one.
func (b *Button) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	n := bl.Node()
	if n.Role == role.Auto {
		n.Role = role.Button
	}
	n.Name = b.words()
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
func (b *Button) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !b.Enabled() {
		return false
	}
	b.click()
	return true
}
