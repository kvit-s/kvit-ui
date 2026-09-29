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

// ChipButton is a chip's twin for a fact that opens something when pressed:
// "3 changes", "1 ahead". At rest it is drawn as a chip is, from the same
// tones, so a row mixing facts that act with facts that do not reads as one
// row of facts; what sets it apart is what a control does anyway: the ground
// changes under the pointer, the focus ring shows for the keyboard, and a
// screen reader is told it is a button.
//
// It has three states beside its tone. Current marks the chip whose
// destination is already open: the selection tint, an accent edge and bold
// words. Selectable makes it a quick filter that is on or off, as kvit-cash's
// copy of the library adds: a control's height, filled with the accent
// when Selected. And UnavailableReason says why it cannot be pressed; it
// stays in the tab order with its fill gone, so the one chip with something
// to explain is not the one a reader cannot reach.
type ChipButton struct {
	control
	// Text is the chip's words.
	Text string
	// Tone is what its colour means.
	Tone Tone
	// Strong fills the chip in its tone.
	Strong bool
	// Symbol is an optional meaning name before the words; TrailingSymbol one
	// after them, for a caller that knows where pressing leads.
	Symbol, TrailingSymbol string
	// Explanation says where pressing the chip goes, where the words do not.
	Explanation string
	// UnavailableReason says why the chip cannot be pressed; "" means it can.
	UnavailableReason string
	// Current marks the chip whose destination is already open.
	Current bool
	// Selectable makes the chip a choice that is on or off; Selected is
	// whether it is on.
	Selectable, Selected bool
	// OnActivate runs when an available chip is pressed.
	OnActivate func()
}

// NewChipButton returns a neutral chip button.
func NewChipButton(ui *UI, words string) *ChipButton {
	c := &ChipButton{Text: words}
	c.Self = c
	c.initControl(ui, func() {
		if c.UnavailableReason == "" && c.OnActivate != nil {
			c.OnActivate()
		}
	}, unison.KeySpace, unison.KeyReturn, unison.KeyNumPadEnter)
	c.ringRadius = func() float32 { return float32(ui.Interface.RadiusChip()) }
	c.SetSizer(c.sizes)
	c.DrawCallback = c.draw
	c.UpdateCursorCallback = func(geom.Point) *unison.Cursor {
		if c.UnavailableReason != "" {
			return unison.ArrowCursor()
		}
		return unison.PointingCursor()
	}
	c.tip = c.tooltip
	return c
}

func (c *ChipButton) chosen() bool { return c.Selectable && c.Selected }

func (c *ChipButton) lit() bool { return c.hovered || c.KeyboardFocus() }

// tooltip is whichever applies: why the chip cannot be pressed, the whole
// label where it was cut short, and where pressing it goes.
func (c *ChipButton) tooltip() string {
	cut := c.truncated()
	switch {
	case c.UnavailableReason != "":
		return c.UnavailableReason
	case cut && c.Explanation != "":
		return c.Text + "\n" + c.Explanation
	case cut:
		return c.Text
	}
	return c.Explanation
}

func (c *ChipButton) ink(tk tokens.Tokens) palette.Color {
	tone := c.Tone.color(tk)
	switch {
	case c.UnavailableReason != "":
		return tk.TextDisabled
	case c.chosen():
		return tokens.LabelOn(tk.Accent)
	case c.Current:
		return tk.TextPrimary
	case c.Strong:
		return tokens.LabelOn(tone)
	case c.Tone == ToneNeutral:
		return tk.TextSecondary
	}
	return tone
}

func (c *ChipButton) label(width float32) *text.Layout {
	ui, tk := c.ui, c.ui.Theme.Tokens()
	weight := text.Regular
	if c.chosen() || c.Current {
		weight = text.Bold
	}
	return ui.Fonts.Layout([]text.Span{{Text: c.Text, Style: ui.Chrome(ui.Size(RoleCaption), weight, c.ink(tk))}},
		text.Options{MaxWidth: width, Elide: width > 0})
}

// symbols is the width the two symbols take with their gaps.
func (c *ChipButton) symbols() float32 {
	m := c.ui.Interface
	var w float32
	for _, s := range []string{c.Symbol, c.TrailingSymbol} {
		if s != "" {
			w += float32(m.Caption() + m.SpaceSnug())
		}
	}
	return w
}

// truncated reports whether the words were cut short to fit.
func (c *ChipButton) truncated() bool {
	full, _ := c.label(0).Size()
	room := c.ContentRect(false).Width - 2*float32(c.ui.Interface.SpaceNear()) - c.symbols()
	return full > room+0.5
}

func (c *ChipButton) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := c.ui.Interface
	w, h := c.label(0).Size()
	height := float32(m.ChipHeight())
	if c.Selectable {
		// A choice pressed over and over is a control's height; a fact keeps
		// the mark's height, so rows of facts do not grow.
		height = float32(m.ControlHeight())
	}
	height = max(height, h)
	pad := 2 * float32(m.SpaceNear())
	return geom.NewSize(pad+c.symbols(), height), geom.NewSize(pad+c.symbols()+w, height), geom.NewSize(unison.DefaultMaxSize, height)
}

func (c *ChipButton) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := c.ui, c.ui.Theme.Tokens(), c.ui.Interface
	p := painterFor(gc, ui)
	r := c.ContentRect(false)
	radius := float32(m.RadiusChip())
	tone := c.Tone.color(tk)
	unavailable := c.UnavailableReason != ""
	press := func(col palette.Color) palette.Color {
		if c.pressed {
			return col.Darker(1.15)
		}
		return col
	}
	switch {
	case unavailable:
	case c.chosen():
		p.round(r, radius, press(tk.Accent))
	case c.Current:
		p.round(r, radius, tk.SelectionTint)
	case c.Strong:
		p.round(r, radius, press(tone))
	case c.Tone == ToneNeutral:
		ground := tk.ChipBackground
		if c.pressed || c.lit() {
			ground = tk.HoverTint
		}
		p.round(r, radius, ground)
	default:
		share := float32(0.16)
		if c.pressed || c.lit() {
			share = 0.28
		}
		p.roundTint(r, radius, tone, share)
	}
	if unavailable || c.Current || !(c.Strong || c.chosen()) {
		edge := tone
		switch {
		case unavailable:
			edge = tk.Border
		case c.Current:
			edge = tk.Accent
		case c.Tone == ToneNeutral && c.lit():
			edge = tk.BorderStrong
		case c.Tone == ToneNeutral:
			edge = tk.Border
		}
		p.outline(r, radius, float32(m.Hairline()), edge)
	}
	pad, gap, s := float32(m.SpaceNear()), float32(m.SpaceSnug()), float32(m.Caption())
	room := r.Width - 2*pad - c.symbols()
	l := c.label(room)
	lw, lh := l.Size()
	content := c.symbols() + lw
	x := r.X + (r.Width-content)/2
	glyph := func(name string) {
		if g, ok := icons.Glyph(name); ok {
			drawGlyph(gc, ui, g, s, c.ink(tk), geom.NewRect(x, r.Y+(r.Height-s)/2, s, s))
		}
		x += s + gap
	}
	if c.Symbol != "" {
		glyph(c.Symbol)
	}
	l.Draw(gc, x, r.Y+(r.Height-lh)/2)
	x += lw + gap
	if c.TrailingSymbol != "" {
		glyph(c.TrailingSymbol)
	}
}

// ProvideAccessibility describes the chip as a button: why it cannot be
// pressed, or where it goes; whether it is the current one; and, for a
// choice, whether it is on.
func (c *ChipButton) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Button
		if c.Selectable {
			n.Role = role.ToggleButton
			n.Pressed = c.Selected
		}
	}
	n.Name = c.Text
	n.Description = c.Explanation
	if c.UnavailableReason != "" {
		n.Description = c.UnavailableReason
	}
	n.Selected = c.Current
	if c.Selectable {
		n.HasCheck = true
		n.Checked = check.Off
		if c.Selected {
			n.Checked = check.On
		}
	}
	if c.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction presses the chip for a screen reader.
func (c *ChipButton) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !c.Enabled() {
		return false
	}
	c.fire()
	return true
}
