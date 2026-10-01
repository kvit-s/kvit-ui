package kvitui

import (
	"strconv"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// Stepper is a number with a minus and a plus beside it: an interface size, a
// column count, a number of days. It is for a small whole-number range a
// reader adjusts by one or two, which is where a slider is wrong: a slider
// makes a precise value hard to hit and does not say what the value is.
type Stepper struct {
	unison.Panel
	ui *UI
	// Label names what the number is, such as "Interface size".
	Label string
	// Unit follows the number, in the muted colour.
	Unit string
	// Value is the number.
	Value int
	// From and To are the ends of the range; Step is how far a press moves.
	From, To, Step int
	// OnChange runs after a press changes the value, with the new value.
	OnChange func(value int)
	// Follow, when set, gives the value from whatever owns it, read before
	// every layout, so the stepper shows the owner's value when it changes
	// elsewhere.
	Follow func() int

	minus, plus *IconButton
	figure      *Figure
}

// NewStepper returns a stepper over a range, at its lowest value.
func NewStepper(ui *UI, label string, from, to int) *Stepper {
	s := &Stepper{ui: ui, Label: label, From: from, To: to, Step: 1, Value: from}
	s.Self = s
	s.minus = NewIconButton(ui, "minus", "")
	s.plus = NewIconButton(ui, "plus", "")
	for _, b := range []*IconButton{s.minus, s.plus} {
		b.Form, b.Size = Ordinary, Px(22)
	}
	s.minus.OnClick = func() { s.apply(s.Value - s.Step) }
	s.plus.OnClick = func() { s.apply(s.Value + s.Step) }
	s.figure = NewFigure(ui, "", "")
	// A fixed width, so what is beside it does not shuffle sideways each time
	// the number gains or loses a digit.
	width := Width(ui, Px(34), s.figure)
	for _, p := range []unison.Paneler{s.minus, width, s.plus} {
		p.AsPanel().SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
		s.AddChild(p)
	}
	s.SetLayout(syncing{Layout: &spaced{FlexLayout: unison.FlexLayout{Columns: 3, VAlign: align.Middle}, ui: ui, gap: SizeSpaceSnug, horizontal: true}, sync: s.sync})
	return s
}

func (s *Stepper) sync() {
	if s.Follow != nil {
		s.Value = s.Follow()
	}
	s.figure.Value, s.figure.Unit = strconv.Itoa(s.Value), s.Unit
	s.minus.Label, s.plus.Label = "Decrease "+s.Label, "Increase "+s.Label
	s.minus.SetEnabled(s.Value > s.From)
	s.plus.SetEnabled(s.Value < s.To)
}

// apply moves the value within the range, and reports a change.
func (s *Stepper) apply(next int) {
	next = max(s.From, min(s.To, next))
	if next == s.Value {
		return
	}
	s.Value = next
	s.sync()
	s.MarkForLayoutAndRedraw()
	if s.OnChange != nil {
		s.OnChange(next)
	}
}

// ProvideAccessibility describes the stepper as a spin button with its
// value, range and step. unison reads a spin button as text and leaves its
// children out, so the minus and plus are offered as its decrement and
// increment actions instead.
func (s *Stepper) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.SpinButton
	}
	n.Name = s.Label
	n.Value = joinWords(strconv.Itoa(s.Value), s.Unit)
	n.HasNumber = true
	n.Number, n.Min, n.Max, n.Step = float64(s.Value), float64(s.From), float64(s.To), float64(s.Step)
	if s.Value < s.To {
		n.Actions = n.Actions.With(accessibility.Increment)
	}
	if s.Value > s.From {
		n.Actions = n.Actions.With(accessibility.Decrement)
	}
}

// PerformAccessibilityAction raises or lowers the value by a step.
func (s *Stepper) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	switch req.Action {
	case accessibility.Increment:
		s.apply(s.Value + s.Step)
	case accessibility.Decrement:
		s.apply(s.Value - s.Step)
	default:
		return false
	}
	return true
}
