package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// ChoiceButton is a button with a choice of value fused to its right side,
// drawn as one control: [Measure | All ▾]. The button runs the action, and
// the choice shows the value it runs with, so the reader sees what a press
// will do before pressing. Pressing the choice opens the list of values;
// choosing one changes the value and runs nothing.
//
// It is for an action repeated with a setting changed less often than the
// action is run, such as how many items Measure takes. When each value is
// an action of its own, chosen and run in one press, a SplitButton of
// variants is right instead (CMD-17).
//
// The two parts are two stops in the Tab order. The choice is a Select: Up
// and Down change the value without opening the list, and typing the start
// of a value's words chooses it.
type ChoiceButton struct {
	unison.Panel
	ui *UI
	// Command, when set, is what the button runs, as a Button's is: its
	// words, symbol, explanation and whether it can act come from the
	// command, and its Run reads Current for the value. A command that
	// cannot act disables both parts.
	Command *Command
	// Text is what the button says, in sentence case: kvit-ui writes it in
	// title case on the macOS form (WRD-6). The values are shown as written.
	Text string
	// Form is how both parts look: ordinary, primary, quiet or flat.
	Form ButtonForm
	// Danger draws both parts in the danger colour.
	Danger bool
	// Symbol is an optional meaning name drawn before the button's words.
	Symbol string
	// Busy says the action is running: the button's words change to
	// BusyText, as a Button's do.
	Busy bool
	// BusyText is what the busy button says: an -ing phrase naming the
	// work, such as "Measuring…" (WRD-25).
	BusyText string
	// Explanation is one sentence the words cannot say: why the button is
	// disabled, or what pressing it does. It is the button's tooltip and
	// accessible description.
	Explanation string
	// ChoiceLabel names the choice for a screen reader, saying what its
	// values are, such as "How many".
	ChoiceLabel string
	// Options are the values to choose from.
	Options []Option
	// Current is the Value of the chosen option.
	Current string
	// OnChoose runs after a value is chosen, with the value.
	OnChoose func(value string)
	// OnPress runs when the button is pressed, with the value chosen. A
	// Command runs in its place when one is set.
	OnPress func(value string)

	main   *Button
	choice *Select
}

// NewChoiceButton returns an ordinary button with a choice of values fused
// to it, the first value chosen.
func NewChoiceButton(ui *UI, label, choiceLabel string, options ...Option) *ChoiceButton {
	c := &ChoiceButton{ui: ui, Text: label, ChoiceLabel: choiceLabel, Options: options}
	c.Self = c
	if len(options) > 0 {
		c.Current = options[0].Value
	}
	c.main = NewButton(ui, label)
	c.main.fusedRight = true
	c.main.OnClick = func() {
		if c.OnPress != nil {
			c.OnPress(c.Current)
		}
	}
	c.choice = NewSelect(ui, choiceLabel)
	c.choice.fusedTo = c.main
	c.choice.OnChoose = func(value string) {
		c.Current = value
		if c.OnChoose != nil {
			c.OnChoose(value)
		}
	}
	c.AddChild(c.main)
	c.AddChild(c.choice)
	c.SetLayout(syncing{Layout: choiceButtonLayout{c}, sync: c.sync})
	// The pair draws before its parts, so each part draws the fields as they
	// are now, such as a Current the app has just set.
	c.DrawCallback = func(*unison.Canvas, geom.Rect) { c.sync() }
	return c
}

// sync copies the fields into the two parts. The choice is enabled while
// the button is: the value means nothing apart from the action, and one
// half of a shape greyed reads as a broken control.
func (c *ChoiceButton) sync() {
	b, s := c.main, c.choice
	b.Command, b.Text, b.Form, b.Danger, b.Symbol = c.Command, c.Text, c.Form, c.Danger, c.Symbol
	b.Busy, b.BusyText, b.Explanation = c.Busy, c.BusyText, c.Explanation
	s.Label, s.Options, s.Current = c.ChoiceLabel, c.Options, c.Current
	if c.Command == nil {
		b.SetEnabled(c.Enabled())
	}
	b.sync()
	s.SetEnabled(b.Enabled())
}

type choiceButtonLayout struct{ c *ChoiceButton }

func (l choiceButtonLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	_, bp, _ := l.c.main.Sizes(geom.Size{})
	_, sp, _ := l.c.choice.Sizes(geom.Size{})
	size := geom.NewSize(bp.Width+sp.Width, max(bp.Height, sp.Height))
	return size, size, size
}

func (l choiceButtonLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	_, sp, _ := l.c.choice.Sizes(geom.Size{})
	l.c.main.SetFrameRect(geom.NewRect(r.X, r.Y, r.Width-sp.Width, r.Height))
	l.c.choice.SetFrameRect(geom.NewRect(r.Right()-sp.Width, r.Y, sp.Width, r.Height))
}

// ProvideAccessibility names the pair as a group, by the button's words; a
// screen reader then reads the button and the choice with its value.
func (c *ChoiceButton) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	c.sync()
	n.Name = c.main.Words()
}
