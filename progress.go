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
	// ShowPercent draws the share done beside the label.
	ShowPercent bool

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
	p := &Progress{ui: ui, Label: label, Value: value, Maximum: 1, Determinate: true, ShowPercent: true, wanted: -1}
	p.Self = p
	p.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		m := ui.Interface
		h := float32(m.RowHeightCompact())
		return geom.NewSize(float32(m.Px(60)), h), geom.NewSize(float32(m.Px(200)), h), geom.NewSize(unison.DefaultMaxSize, h)
	})
	p.DrawCallback = p.draw
	return p
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
	gap := float32(m.Space())
	right := r.Right()
	if p.Determinate && p.ShowPercent {
		st := ui.Chrome(ui.Size(RoleSmall), text.Regular, t.TextMuted)
		st.Tabular = true
		l := ui.Fonts.Layout([]text.Span{{Text: p.percent(), Style: st}}, text.Options{})
		w, _ := l.Size()
		l.Draw(gc, r.Right()-w, r.Y)
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
// gets little time and the shortfall compounds: the Qt library measured a
// label reading 42 per cent over a bar drawing about 5. A step whose slide
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
		n.Description = strconv.Itoa(int(math.Round(p.fraction()*100))) + " percent"
		n.HasNumber = true
		n.Number, n.Min, n.Max = p.fraction()*100, 0, 100
		return
	}
	n.Description = "in progress"
	n.Busy = true
}
