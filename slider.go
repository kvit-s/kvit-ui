package kvitui

import (
	"math"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// Slider is a value chosen from a continuous range by dragging: right where
// the reader is choosing a feel rather than a number, such as an opacity or a
// threshold judged by its result, and wrong where the exact value matters,
// since a particular number is hard to hit and cannot be typed; a Stepper or
// a NumberField is for that. The value is always drawn beside it, because a
// slider whose position is its only output makes the reader estimate.
//
// The value stands to the right of the slider's box, as the Qt slider draws
// it, so a slider in a row wants room after it.
type Slider struct {
	control
	// Label names the value, for a screen reader.
	Label string
	// Unit follows the value.
	Unit string
	// Precision is how many decimal places the value is drawn with.
	Precision int
	// Value is the value, From and To the range; 0 to 1 unless set.
	Value, From, To float64
	// Step is how far an arrow key moves it; a tenth of the range unless set.
	Step float64
	// OnChange runs as the value moves.
	OnChange func(value float64)

	figure *Figure // stamped beside the slider
}

// NewSlider returns a slider over 0 to 1 at value.
func NewSlider(ui *UI, label string, value float64) *Slider {
	s := &Slider{Label: label, Value: value, To: 1}
	s.Self = s
	s.initControl(ui, nil)
	s.figure = NewFigure(ui, "", "")
	s.figure.Role = RoleSmall
	s.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		m := ui.Interface
		h := float32(m.ControlHeight())
		return geom.NewSize(float32(m.Px(60)), h), geom.NewSize(float32(m.Px(180)), h), geom.NewSize(unison.DefaultMaxSize, h)
	})
	s.DrawCallback = s.draw
	// The ring goes around the handle, as the Qt slider draws it, and the
	// slider draws it itself.
	s.noRing = true
	down := s.MouseDownCallback
	s.MouseDownCallback = func(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
		if !down(where, button, clicks, mods) {
			return false
		}
		s.setAt(where)
		return true
	}
	s.MouseDragCallback = func(where geom.Point, _ int, _ mod.Modifiers) bool {
		if s.pressed {
			s.setAt(where)
		}
		return true
	}
	s.MouseUpCallback = func(geom.Point, int, mod.Modifiers) bool {
		s.pressed, s.pointerFocus = false, false
		s.MarkForRedraw()
		return true
	}
	s.KeyDownCallback = s.keyDown
	return s
}

// track is the groove's extent, the slider's whole width: the Qt slider,
// unlike its check box, has no padding at its ends.
func (s *Slider) track() geom.Rect { return s.ContentRect(false) }

// position is where the value is along the range, 0 to 1.
func (s *Slider) position() float64 {
	if s.To == s.From {
		return 0
	}
	return max(0, min(1, (s.Value-s.From)/(s.To-s.From)))
}

func (s *Slider) handle() geom.Rect {
	a := s.track()
	d := float32(s.ui.Interface.Px(14))
	return geom.NewRect(a.X+float32(s.position())*(a.Width-d), a.Y+(a.Height-d)/2, d, d)
}

// set moves the value, kept within the range.
func (s *Slider) set(v float64) {
	lo, hi := min(s.From, s.To), max(s.From, s.To)
	v = max(lo, min(hi, v))
	if v == s.Value {
		return
	}
	s.Value = v
	s.MarkForRedraw()
	if s.OnChange != nil {
		s.OnChange(v)
	}
}

func (s *Slider) setAt(where geom.Point) {
	a := s.track()
	d := float32(s.ui.Interface.Px(14))
	if a.Width <= d {
		return
	}
	pos := float64((where.X - a.X - d/2) / (a.Width - d))
	s.set(s.From + max(0, min(1, pos))*(s.To-s.From))
}

func (s *Slider) step() float64 {
	if s.Step > 0 {
		return s.Step
	}
	return math.Abs(s.To-s.From) / 10
}

func (s *Slider) keyDown(key unison.KeyCode, _ mod.Modifiers, _ bool) bool {
	if !s.Enabled() {
		return false
	}
	switch key {
	case unison.KeyLeft, unison.KeyDown:
		s.set(s.Value - s.step())
	case unison.KeyRight, unison.KeyUp:
		s.set(s.Value + s.step())
	case unison.KeyHome:
		s.set(s.From)
	case unison.KeyEnd:
		s.set(s.To)
	default:
		return false
	}
	return true
}

func (s *Slider) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := s.ui, s.ui.Theme.Tokens(), s.ui.Interface
	p := painterFor(gc, ui)
	a := s.track()
	h := float32(m.BarHeight())
	groove := geom.NewRect(a.X, a.Y+(a.Height-h)/2, a.Width, h)
	radius := float32(m.RadiusBar())
	p.round(groove, radius, t.ChipBackground)
	fill := t.Accent
	if !s.Enabled() {
		fill = t.Border
	}
	p.round(geom.NewRect(groove.X, groove.Y, float32(s.position())*groove.Width, h), radius, fill)
	knob := s.handle()
	ground := t.PopupBackground
	if s.pressed {
		ground = t.HoverTint
	}
	p.round(knob, knob.Width/2, ground)
	p.outline(knob, knob.Width/2, float32(m.Hairline()), t.BorderStrong)
	if s.KeyboardFocus() {
		w := float32(m.FocusRingWidth())
		ring := knob.Inset(geom.NewUniformInsets(-w / 2))
		p.outline(ring, ring.Width/2, w, t.FocusRing)
	}
	if !s.ui.spillHosted(s.AsPanel()) {
		s.drawSpill(gc)
	}
}

// drawSpill draws the value a space to the right of the slider's box.
func (s *Slider) drawSpill(gc *unison.Canvas) {
	f := s.figure
	f.Value, f.Unit = toFixed(s.Value, s.Precision), s.Unit
	size := preferred(f)
	r := s.ContentRect(false)
	stampAt(gc, f, geom.NewRect(r.Right()+float32(s.ui.Interface.Space()), r.Y+(r.Height-size.Height)/2, size.Width, size.Height))
}

// ProvideAccessibility describes the slider by its name and value, with the
// steps a screen reader can take.
func (s *Slider) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Slider
	}
	n.Name = s.Label
	n.Value = joinWords(toFixed(s.Value, s.Precision), s.Unit)
	n.HasNumber = true
	n.Number, n.Min, n.Max, n.Step = s.Value, min(s.From, s.To), max(s.From, s.To), s.step()
	if s.Enabled() {
		n.Actions = n.Actions.With(accessibility.Increment, accessibility.Decrement, accessibility.SetValue)
	}
}

// PerformAccessibilityAction moves the value for a screen reader.
func (s *Slider) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if !s.Enabled() {
		return false
	}
	switch req.Action {
	case accessibility.Increment:
		s.set(s.Value + s.step())
	case accessibility.Decrement:
		s.set(s.Value - s.step())
	case accessibility.SetValue:
		s.set(req.Number)
	default:
		return false
	}
	return true
}
