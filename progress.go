package kvitui

import (
	"math"
	"strconv"
	"time"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/pathop"
	"github.com/richardwilkes/unison/enums/role"
)

// Progress is how far through something the application is. Which of two
// states it is in is the most useful thing on the bar: with a known total the
// fill says how much is done, and without one a band moves to say only that
// work is happening, since an unknown total drawn as a bar creeping towards
// the end is what makes people distrust progress bars. It carries a label
// saying what is happening: "47%" says how far, and not how far through what.
type Progress struct {
	unison.Panel
	ui *UI
	// Value against Maximum, 1 unless set.
	Value, Maximum float64
	// Determinate says the total is known; false draws the moving band.
	Determinate bool
	// Label says what is happening.
	Label string
	// ShowPercent draws the share done beside the label; off unless set,
	// since the label saying what is happening is what the reader needs
	// (FDB-10).
	ShowPercent bool
	// Remaining is how long the work is expected to take still, said as
	// "About 3 minutes left" beside the label; 0 says nothing (FDB-10).
	Remaining time.Duration
	// HaltText is the words of the button that stops the work: Cancel while
	// stopping leaves no trace, Stop once what is done is kept (FDB-3,
	// WRD-8); "" draws no button. OnHalt runs when it is pressed.
	HaltText string
	OnHalt   func()
	// Pausable draws Pause beside the bar, Resume while Paused, for work
	// that cancelling would throw away (FDB-3); OnPause runs with the new
	// state.
	Pausable, Paused bool
	OnPause          func(paused bool)

	halt, pause *Button

	// The fill slides to a new value, but only while the slide can finish;
	// see draw.
	width    float32 // the fill's width as drawn
	wanted   float32 // where the last step asked it to be, -1 before any
	from, to float32
	started  time.Time
	sliding  bool
	band     time.Time // when the moving band started
}

// NewProgress returns a determinate bar at value of 1, labelled.
func NewProgress(ui *UI, label string, value float64) *Progress {
	p := &Progress{ui: ui, Label: label, Value: value, Maximum: 1, Determinate: true, wanted: -1}
	p.Self = p
	p.halt = NewButton(ui, "")
	p.halt.Form = ButtonQuiet
	p.halt.OnClick = func() {
		if p.OnHalt != nil {
			p.OnHalt()
		}
	}
	p.pause = NewButton(ui, "")
	p.pause.Form = ButtonQuiet
	p.pause.OnClick = func() {
		p.Paused = !p.Paused
		p.MarkForLayoutAndRedraw()
		if p.OnPause != nil {
			p.OnPause(p.Paused)
		}
	}
	p.SetLayout(syncing{Layout: progressLayout{p}, sync: p.sync})
	p.DrawCallback = p.draw
	return p
}

// sync shows the buttons the fields ask for.
func (p *Progress) sync() {
	p.halt.Text = p.HaltText
	p.pause.Text = "Pause"
	if p.Paused {
		p.pause.Text = "Resume"
	}
	var parts []unison.Paneler
	if p.Pausable {
		parts = append(parts, p.pause)
	}
	if p.HaltText != "" {
		parts = append(parts, p.halt)
	}
	showOnly(p.AsPanel(), parts...)
}

// buttonsWidth is the width the buttons take at the bar's end, with the gap
// before them.
func (p *Progress) buttonsWidth() float32 {
	w := float32(0)
	for _, c := range p.Children() {
		_, pref, _ := c.Sizes(geom.Size{})
		w += pref.Width + float32(p.ui.Interface.Space())
	}
	return w
}

type progressLayout struct{ p *Progress }

func (l progressLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.p.ui.Interface
	h := float32(m.RowHeightCompact())
	for _, c := range l.p.Children() {
		_, pref, _ := c.Sizes(geom.Size{})
		h = max(h, pref.Height)
	}
	extra := l.p.buttonsWidth()
	return geom.NewSize(float32(m.Px(60))+extra, h), geom.NewSize(float32(m.Px(200))+extra, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l progressLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	x := r.Right()
	kids := l.p.Children()
	for i := len(kids) - 1; i >= 0; i-- {
		_, pref, _ := kids[i].Sizes(geom.Size{})
		x -= pref.Width
		kids[i].SetFrameRect(geom.NewRect(x, r.Y+(r.Height-pref.Height)/2, pref.Width, pref.Height))
		x -= float32(l.p.ui.Interface.Space())
	}
}

// remainingWords say how long is left: "About 3 minutes left", or "" when
// the time is not known.
func (p *Progress) remainingWords() string {
	d := p.Remaining
	switch {
	case d <= 0:
		return ""
	case d < time.Minute:
		return "About " + p.ui.CountPhrase(max(1, int(math.Round(d.Seconds()))), "second", "") + " left"
	case d < time.Hour:
		return "About " + p.ui.CountPhrase(int(math.Round(d.Minutes())), "minute", "") + " left"
	}
	return "About " + p.ui.CountPhrase(int(math.Round(d.Hours())), "hour", "") + " left"
}

func (p *Progress) fraction() float64 {
	if p.Maximum <= 0 {
		return 0
	}
	return max(0, min(1, p.Value/p.Maximum))
}

func (p *Progress) percent() string {
	return strconv.Itoa(int(math.Round(p.fraction()*100))) + "%"
}

func (p *Progress) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := p.ui, p.ui.Theme.Tokens(), p.ui.Interface
	pt := painterFor(gc, ui)
	r := p.ContentRect(false)
	r.Width = max(0, r.Width-p.buttonsWidth())
	gap := float32(m.Space())
	right := r.Right()
	if left := p.remainingWords(); left != "" {
		l := ui.Fonts.Layout([]text.Span{{Text: left, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, t.TextMuted)}}, text.Options{})
		w, _ := l.Size()
		l.Draw(gc, right-w, r.Y)
		right -= w + gap
	}
	if p.Determinate && p.ShowPercent {
		st := ui.Chrome(ui.Size(RoleSmall), text.Regular, t.TextMuted)
		st.Tabular = true
		l := ui.Fonts.Layout([]text.Span{{Text: p.percent(), Style: st}}, text.Options{})
		w, _ := l.Size()
		l.Draw(gc, right-w, r.Y)
		right -= w + gap
	}
	l := ui.Fonts.Layout([]text.Span{{Text: p.Label, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, t.TextSecondary)}},
		text.Options{MaxWidth: max(1, right-r.X), Elide: true})
	l.Draw(gc, r.X, r.Y)

	h := float32(m.BarHeight())
	track := geom.NewRect(r.X, r.Bottom()-h, r.Width, h)
	radius := float32(m.RadiusBar())
	pt.round(track, radius, t.ChipBackground)
	gc.Save()
	gc.ClipRect(track, pathop.Intersect, false)
	scale := ui.Theme.MotionScale()
	if p.Determinate {
		pt.round(geom.NewRect(track.X, track.Y, p.fill(track.Width, scale), h), radius, t.Accent)
	} else if scale <= 0 {
		// With motion stilled the band does not move: the whole track takes
		// a muted fill, so the reader still sees that something is happening
		// rather than an empty bar that looks stalled.
		pt.roundTint(track, radius, t.Accent, 0.4)
	} else {
		if p.band.IsZero() {
			p.band = time.Now()
		}
		const loop = 1200 * time.Millisecond
		share := float32(time.Since(p.band)%loop) / float32(loop)
		w := track.Width * 0.3
		pt.round(geom.NewRect(track.X-w+share*(track.Width+w), track.Y, w, h), radius, t.Accent)
		p.nextFrame()
	}
	gc.Restore()
	if !p.Determinate {
		return
	}
	p.band = time.Time{}
}

// fill is the fill's width now. A new value slides the fill to it over
// 160 ms, which reads as one bar moving rather than a bar redrawn, but only
// when the bar is where the step before asked it to be. Long work reports its
// progress from the thread that also draws, so between two reports a slide
// gets little time and the shortfall compounds, until a label reading 42
// per cent can stand over a bar drawn at about 5. A step whose slide
// never arrived jumps instead, so the bar is never more than one report
// behind. The first value is drawn rather than slid to.
func (p *Progress) fill(width float32, scale float64) float32 {
	wanted := float32(math.Round(float64(width) * p.fraction()))
	if wanted != p.wanted {
		stalled := math.Abs(float64(p.width-p.wanted)) > 0.5
		p.wanted = wanted
		if stalled || scale <= 0 || p.Window() == nil {
			p.width, p.sliding = wanted, false
		} else {
			p.from, p.to, p.started, p.sliding = p.width, wanted, time.Now(), true
		}
	}
	if p.sliding {
		s := min(1, float64(time.Since(p.started))/(160*scale*float64(time.Millisecond)))
		p.width = p.from + (p.to-p.from)*float32(s)
		if s >= 1 {
			p.width, p.sliding = p.to, false
		} else {
			p.nextFrame()
		}
	}
	return p.width
}

// nextFrame asks for the bar to be drawn again a frame from now. It is asked
// only from a draw, so a bar that is not drawn, being hidden or out of its
// window, stops asking.
func (p *Progress) nextFrame() {
	unison.InvokeTaskAfter(p.MarkForRedraw, 16*time.Millisecond)
}

// ProvideAccessibility reads the bar as its label and how far it is, or that
// it is under way.
func (p *Progress) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.ProgressBar
	}
	n.Name = p.Label
	if p.Determinate {
		n.Description = joinWords(strconv.Itoa(int(math.Round(p.fraction()*100)))+" percent", p.remainingWords())
		n.HasNumber = true
		n.Number, n.Min, n.Max = p.fraction()*100, 0, 100
		return
	}
	n.Description = "in progress"
	n.Busy = true
}
