package kvitui

import (
	"math"
	"time"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/role"
)

// Switch is an option that is on or off and takes effect the moment it moves,
// where a Check is a value in a form that takes effect when the form is
// submitted: a settings row that applies at once is a switch, and a row of
// options above a Save button is checkboxes. The knob moves as well as the
// track filling, so the state is told by position as well as by colour.
type Switch struct {
	control
	// Text is what the switch turns on.
	Text string
	// Checked is whether it is on.
	Checked bool
	// OnChange runs after a press, with the new state.
	OnChange func(checked bool)

	track  *unison.Panel
	knob   float32 // where the knob is drawn, 0 off to 1 on
	moving bool
}

// NewSwitch returns a switch that is off.
func NewSwitch(ui *UI, label string) *Switch {
	s := &Switch{Text: label}
	s.Self = s
	s.initControl(ui, s.toggle, unison.KeySpace)
	s.track = unison.NewPanel()
	s.track.DrawCallback = s.drawTrack
	s.AddChild(s.track)
	s.SetLayout(syncing{Layout: switchLayout{s}, sync: s.sync})
	s.DrawCallback = s.drawLabel
	return s
}

func (s *Switch) toggle() {
	s.Checked = !s.Checked
	s.slide()
	if s.OnChange != nil {
		s.OnChange(s.Checked)
	}
}

// sync puts the knob where the state says when nothing is moving it, so a
// switch set by its owner is drawn set.
func (s *Switch) sync() {
	if !s.moving {
		s.knob = s.target()
	}
}

func (s *Switch) target() float32 {
	if s.Checked {
		return 1
	}
	return 0
}

// slide moves the knob to the state over 120 ms, easing out, or at once when
// motion is reduced.
func (s *Switch) slide() {
	scale := s.ui.Theme.MotionScale()
	if scale == 0 || s.Window() == nil {
		s.knob = s.target()
		s.MarkForRedraw()
		return
	}
	from, start := s.knob, time.Now()
	s.moving = true
	var step func()
	step = func() {
		t := min(1, float64(time.Since(start))/(120*scale*float64(time.Millisecond)))
		to := s.target()
		s.knob = from + (to-from)*float32(1-math.Pow(1-t, 3))
		s.MarkForRedraw()
		if t >= 1 {
			s.knob, s.moving = to, false
			return
		}
		unison.InvokeTaskAfter(step, 16*time.Millisecond)
	}
	unison.InvokeTaskAfter(step, 16*time.Millisecond)
}

// focusRing puts the ring around the track.
func (s *Switch) focusRing() (*unison.Panel, float32, bool) {
	return s.track, float32(s.ui.Interface.Px(16)) / 2, s.KeyboardFocus()
}

func (s *Switch) label() *text.Layout {
	ui, t := s.ui, s.ui.Theme.Tokens()
	ink := t.TextPrimary
	if !s.Enabled() {
		ink = t.TextDisabled
	}
	return ui.Fonts.Layout([]text.Span{{Text: s.Text, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ink)}}, text.Options{})
}

type switchLayout struct{ s *Switch }

func (l switchLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.s.ui.Interface
	w, h := l.s.label().Size()
	size := geom.NewSize(float32(2*checkPadding+m.Px(28)+m.Space())+w, max(float32(m.ControlHeight()), h))
	return size, size, geom.NewSize(unison.DefaultMaxSize, size.Height)
}

func (l switchLayout) PerformLayout(target *unison.Panel) {
	m := l.s.ui.Interface
	r := target.ContentRect(false)
	h := float32(m.Px(16))
	l.s.track.SetFrameRect(geom.NewRect(r.X+checkPadding, r.Y+(r.Height-h)/2, float32(m.Px(28)), h))
}

func (s *Switch) drawTrack(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := s.ui, s.ui.Theme.Tokens(), s.ui.Interface
	p := painterFor(gc, ui)
	r := s.track.ContentRect(false)
	on := s.Checked && s.Enabled()
	ground, edge := t.ChipBackground, t.BorderStrong
	if on {
		ground = t.Accent
	}
	if !s.Enabled() {
		edge = t.Border
	}
	p.round(r, r.Height/2, ground)
	p.outline(r, r.Height/2, float32(m.Hairline()), edge)
	tight := float32(m.SpaceTight())
	size := r.Height - float32(m.SpaceSnug())
	x := r.X + tight + s.knob*(r.Width-size-2*tight)
	ink := t.TextMuted
	if on {
		ink = tokens.LabelOn(t.Accent)
	}
	p.round(geom.NewRect(x, r.Y+tight, size, size), size/2, ink)
}

func (s *Switch) drawLabel(gc *unison.Canvas, _ geom.Rect) {
	m := s.ui.Interface
	r := s.ContentRect(false)
	l := s.label()
	_, h := l.Size()
	l.Draw(gc, r.X+float32(checkPadding+m.Px(28)+m.Space()), r.Y+(r.Height-h)/2)
}

// ProvideAccessibility describes the switch as a box that is ticked or not.
func (s *Switch) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.CheckBox
	}
	n.Name = s.Text
	n.HasCheck = true
	n.Checked = check.Off
	if s.Checked {
		n.Checked = check.On
	}
	if s.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction turns the switch over for a screen reader.
func (s *Switch) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !s.Enabled() {
		return false
	}
	s.fire()
	return true
}
