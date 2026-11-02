package kvitui

import (
	"log/slog"

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
	// drawn as pressable at rest rather than only once the pointer is on it.
	ButtonQuiet
	// ButtonFlat has no ground and no outline at rest, only the hover tint
	// under the pointer and the selection tint while pressed or checked, for
	// a toolbar of many buttons told apart by their words.
	ButtonFlat
)

// Button is a button with words on it. Danger is separate from the form,
// because a destructive action can take any of the three: a primary Delete
// in a confirmation, an ordinary one in a settings row, a quiet one in a
// strip. The words on a filled button take the colour that contrasts with
// its own fill, so a filled danger button stays readable.
type Button struct {
	control
	// Command, when set, is what the button runs: its words, symbol,
	// explanation, whether it can act and whether it is on all come from
	// the command, read each time the button is drawn.
	Command *Command
	// Text is what the button says, in sentence case: kvit-ui writes it in
	// title case on the macOS form (WRD-6).
	Text string
	// Asks marks a button that asks for more before acting, such as one
	// opening a dialog: "…" is added to its words (WRD-4).
	Asks bool
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
	// BusyText is what a busy button says: an -ing phrase naming the work,
	// "Importing statement…", never "Working…" (WRD-25). A busy button
	// without one keeps its words and is warned about once.
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

	// trailing is a symbol drawn after the words, as a menu button's
	// down-arrow is.
	trailing string
	// fusedRight says a part is fused to the button's right side, as a
	// ChoiceButton's choice is: the corners there are square and the
	// outline open, and the part draws the line between them.
	fusedRight bool
	// warnedBusy is set once a busy button without BusyText was warned about.
	warnedBusy bool
}

// NewButton returns an ordinary button.
func NewButton(ui *UI, label string) *Button {
	b := &Button{Text: label}
	b.Self = b
	b.initControl(ui, b.click, unison.KeySpace)
	b.pressKeys = buttonKeys
	b.SetSizer(b.sizes)
	b.DrawCallback = b.draw
	b.tip = func() string {
		if b.Command != nil {
			b.sync()
			return b.Command.tooltip(ui, plainWords(b.label()))
		}
		return b.Explanation
	}
	return b
}

// NewCommandButton returns an ordinary button that runs a command.
func NewCommandButton(ui *UI, c *Command) *Button {
	b := NewButton(ui, "")
	b.Command = c
	return b
}

func (b *Button) click() {
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

// sync reads the command's state: whether it can act and whether it is on.
func (b *Button) sync() {
	if c := b.Command; c != nil {
		b.SetEnabled(c.CanRun())
		b.Checked = c.IsChecked()
	}
}

// label is the button's words at rest, in the form's case, with the
// ellipsis of a button that asks, and the access-key marker where the app
// wrote one.
func (b *Button) label() string {
	if c := b.Command; c != nil {
		return c.Label(b.ui)
	}
	t := b.Text
	if b.Asks {
		t += "…"
	}
	return b.ui.CommandCase(t)
}

// Words are what the button says now: its words in the form's case with the
// ellipsis of a button that asks, or its busy words.
func (b *Button) Words() string {
	if b.Busy {
		if b.BusyText != "" {
			return b.BusyText
		}
		if !b.warnedBusy {
			b.warnedBusy = true
			slog.Warn("kvitui: a busy button has no BusyText saying what it is doing", "button", b.Text)
		}
	}
	return plainWords(b.label())
}

// accessWords are the button's words for its access key: none while busy,
// or for OK, Cancel and Close, which Enter and Escape press.
func (b *Button) accessWords() string {
	if b.Busy || noAccessKey[plainWords(b.label())] {
		return ""
	}
	return b.label()
}

// accessPress presses the button, as its access key does.
func (b *Button) accessPress() {
	if b.Enabled() {
		b.fire()
	}
}

// symbol is the meaning name of the symbol before the words, or "".
func (b *Button) symbol() string {
	if b.Symbol == "" && b.Command != nil {
		return b.Command.SymbolName()
	}
	return b.Symbol
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

func (b *Button) layout() *text.Layout {
	ui := b.ui
	st := ui.Chrome(ui.Size(RoleBody), text.Regular, b.foreground(ui.Theme.Tokens()))
	return ui.Fonts.Layout(ui.accessSpans(b.AsPanel(), b.Words(), st), text.Options{})
}

// content is the width and height of the symbol and the words, a near space
// apart.
func (b *Button) content() (width, height float32, l *text.Layout) {
	m := b.ui.Interface
	b.sync()
	l = b.layout()
	width, height = l.Size()
	if b.symbol() != "" {
		s := float32(m.IconSizeSmall())
		width += s + float32(m.SpaceNear())
		height = max(height, s)
	}
	if b.trailing != "" {
		width += float32(m.IconSizeSmall() + m.SpaceNear())
	}
	return width, height, l
}

func (b *Button) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := b.ui.Interface
	w, h, _ := b.content()
	size := geom.NewSize(w+2*float32(m.SpaceLoose()), max(float32(m.ControlHeight()), h))
	return size, size, size
}

// ground is the fill under a part of the button's form, or false for none:
// the button itself, or a part fused to it, pressed or under the pointer, in
// the button's form, danger and enabled state.
func (b *Button) ground(t tokens.Tokens, pressed, hovered, checked bool) (palette.Color, bool) {
	primary := b.Form == ButtonPrimary
	switch {
	case !b.Enabled():
		return t.ChipBackground, primary
	case primary && pressed:
		return b.fill(t).Darker(1.15), true
	case primary:
		return b.fill(t), true
	case pressed, checked:
		return t.SelectionTint, true
	case hovered:
		return t.HoverTint, true
	case b.Form == ButtonQuiet:
		return t.ChipBackground, true
	}
	return palette.Color{}, false
}

// edge is the colour of the button's outline, or false for none: a primary
// button has its fill instead, and a flat one has none unless it is on.
func (b *Button) edge(t tokens.Tokens) (palette.Color, bool) {
	if b.Form == ButtonPrimary || (b.Form == ButtonFlat && !b.Checked) {
		return palette.Color{}, false
	}
	switch {
	case !b.Enabled():
		return t.Border, true
	case b.Checked:
		return t.Accent, true
	case b.Danger:
		return t.Danger, true
	case b.Form == ButtonQuiet:
		return t.Border, true
	}
	return t.BorderStrong, true
}

// shape is the box the button's ground and outline are drawn in: its own,
// reaching past a fused right side by the corner radius, so the panel's clip
// squares the corners there and cuts the outline open.
func (b *Button) shape(r geom.Rect, radius float32) geom.Rect {
	if b.fusedRight {
		r.Width += radius
	}
	return r
}

// fusedArrow is the size of the down-arrow on a part fused to a button's
// right side: smaller than a symbol beside words, since it only marks the
// part as opening a list.
var fusedArrow = Px(10)

// drawFusedPart draws the ground and outline of a part fused to the
// button's right side, as a ChoiceButton's choice and a SplitButton's arrow
// are, whose box is r: the button's ground for the part's own press and
// hover, the button's outline around both, open on the left, and the line
// between the two along the part's left edge.
func (b *Button) drawFusedPart(gc *unison.Canvas, r geom.Rect, pressed, hovered bool) {
	// unison draws a panel's children last first, so the part is drawn
	// before the button and reads the button's command itself.
	b.sync()
	ui, t, m := b.ui, b.ui.Theme.Tokens(), b.ui.Interface
	p := painterFor(gc, ui)
	radius := float32(m.RadiusControl())
	hairline := float32(m.Hairline())
	// Drawn past the left side by the corner radius, where the panel's clip
	// squares the corners and cuts the outline open.
	shape := r
	shape.X -= radius
	shape.Width += radius
	if ground, ok := b.ground(t, pressed, hovered, false); ok {
		p.round(shape, radius, ground)
	}
	line, share := t.Border, float32(1)
	if edge, ok := b.edge(t); ok {
		line = edge
		p.outline(shape, radius, hairline, line)
	} else if b.Form == ButtonPrimary && b.Enabled() {
		// A filled pair has no outline: the line between the parts is the
		// ink of the words on the fill, at a fifth of its strength, which
		// parts the two without cutting across the fill. The high-contrast
		// theme keeps it at full strength, since its edges are there to be
		// seen.
		line = b.foreground(t)
		if ui.Theme.ResolvedTheme() != tokens.HighContrast {
			share = 0.2
		}
	}
	p.fillTint(geom.NewRect(r.X, r.Y, hairline, r.Height), line, share)
}

func (b *Button) draw(gc *unison.Canvas, _ geom.Rect) {
	b.sync()
	ui, t, m := b.ui, b.ui.Theme.Tokens(), b.ui.Interface
	p := painterFor(gc, ui)
	r := b.ContentRect(false)
	radius := float32(m.RadiusControl())
	shape := b.shape(r, radius)
	if ground, ok := b.ground(t, b.pressed, b.hovered, b.Checked); ok {
		p.round(shape, radius, ground)
	}
	if edge, ok := b.edge(t); ok {
		p.outline(shape, radius, float32(m.Hairline()), edge)
	}
	w, _, l := b.content()
	x := r.X + (r.Width-w)/2
	if symbol := b.symbol(); symbol != "" {
		s := float32(m.IconSizeSmall())
		if g, ok := icons.Glyph(symbol); ok {
			drawGlyph(gc, ui, g, s, b.foreground(t), geom.NewRect(x, r.Y+(r.Height-s)/2, s, s))
		}
		x += s + float32(m.SpaceNear())
	}
	lw, h := l.Size()
	l.Draw(gc, x, r.Y+(r.Height-h)/2)
	if b.trailing != "" {
		s := float32(m.IconSizeSmall())
		if g, ok := icons.Glyph(b.trailing); ok {
			drawGlyph(gc, ui, g, s, b.foreground(t), geom.NewRect(x+lw+float32(m.SpaceNear()), r.Y+(r.Height-s)/2, s, s))
		}
	}
}

// ProvideAccessibility describes the button by its words, what it is doing
// when busy, its explanation, and its checked state when it has one.
func (b *Button) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	n := bl.Node()
	if n.Role == role.Auto {
		n.Role = role.Button
	}
	b.sync()
	n.Name = b.Words()
	n.Description = b.Explanation
	if b.Command != nil && b.Explanation == "" {
		n.Description = b.Command.Explanation
	}
	if b.Checkable || b.Checked || (b.Command != nil && b.Command.Checked != nil) {
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
