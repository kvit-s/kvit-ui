package kvitui

import (
	"math"
	"time"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// Pane is the side pane: detail about the one thing the reader has
// selected, beside the list they selected it from. It is the same width
// wherever it appears: a pane that is 320 pixels on one screen and 420 on
// another makes the list beside it reflow as the reader moves between them,
// which reads as the application losing its place. Put it over a view with
// WithPane; it slides in and out rather than appearing, following the
// reduced-motion setting.
type Pane struct {
	Panel
	// Title names what the pane shows.
	Title string
	// CloseLabel is what the close control is called, on hover and to a
	// screen reader: "Close record" for a pane holding one record, since the
	// reader is closing the record rather than the furniture it arrived in.
	CloseLabel string
	// Closable draws the close control; true unless set, and off for a pane
	// the reader cannot dismiss (kvit-cash).
	Closable bool
	// OnClose runs when the close control is pressed.
	OnClose func()

	close  *IconButton
	body   *unison.Panel
	open   bool
	shown  float32 // how far the pane has slid in, 0 to 1
	moving bool
}

// NewPane returns an open pane with a title, holding content.
func NewPane(ui *UI, title string, content ...unison.Paneler) *Pane {
	p := &Pane{Title: title, CloseLabel: "Close the pane", Closable: true, open: true, shown: 1}
	p.ui = ui
	p.Self = p
	p.DrawCallback = p.drawPane
	p.DrawOverCallback = p.drawRules
	p.RuleLeft = true
	p.close = NewIconButton(ui, "close", "")
	p.close.OnClick = func() {
		if p.OnClose != nil {
			p.OnClose()
		}
	}
	p.body = unison.NewPanel()
	p.body.SetLayout(&unison.FlexLayout{Columns: 1})
	for _, c := range content {
		if c.AsPanel().LayoutData() == nil {
			c.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
		}
		p.body.AddChild(c)
	}
	p.AddChild(p.close)
	p.AddChild(p.body)
	p.SetLayout(syncing{Layout: paneLayout{p}, sync: func() {
		p.close.Label = p.CloseLabel
		p.close.Hidden = !p.Closable
	}})
	return p
}

// Open reports whether the pane is open.
func (p *Pane) Open() bool { return p.open }

// SetOpen slides the pane in or out over 160 ms, easing out; with motion
// reduced it moves at once.
func (p *Pane) SetOpen(open bool) {
	if p.open == open {
		return
	}
	p.open = open
	target := float32(0)
	if open {
		target = 1
	}
	scale := p.ui.Theme.MotionScale()
	relayout := func() {
		if parent := p.Parent(); parent != nil {
			parent.MarkForLayoutAndRedraw()
		}
	}
	if scale == 0 || p.Window() == nil {
		p.shown = target
		relayout()
		return
	}
	if p.moving {
		return
	}
	p.moving = true
	from, start := p.shown, time.Now()
	var step func()
	step = func() {
		to := float32(0)
		if p.open {
			to = 1
		}
		if to != target {
			target, from, start = to, p.shown, time.Now()
		}
		t := min(1, float64(time.Since(start))/(160*scale*float64(time.Millisecond)))
		p.shown = from + (target-from)*float32(1-math.Pow(1-t, 3))
		if t >= 1 {
			p.shown, p.moving = target, false
		}
		relayout()
		if p.moving {
			unison.InvokeTaskAfter(step, 16*time.Millisecond)
		}
	}
	unison.InvokeTaskAfter(step, 16*time.Millisecond)
}

func (p *Pane) head() float32 { return float32(p.ui.Interface.BreadcrumbHeight()) }

func (p *Pane) drawPane(gc *unison.Canvas, r geom.Rect) {
	p.draw(gc, r)
	ui, t, m := p.ui, p.ui.Theme.Tokens(), p.ui.Interface
	b := p.ContentRect(true)
	pt := painterFor(gc, ui)
	pt.fill(geom.NewRect(b.X, b.Y+p.head(), b.Width, float32(m.Hairline())), t.Border)
	left := b.X + float32(m.ViewMargin())
	right := b.Right() - float32(m.SpaceNear())
	if p.Closable {
		right = p.close.FrameRect().X - float32(m.Space())
	}
	l := ui.Fonts.Layout([]text.Span{{Text: p.Title, Style: ui.Chrome(ui.Size(RoleStrong), text.Bold, t.TextPrimary)}},
		text.Options{MaxWidth: max(0, right-left), Elide: true})
	_, h := l.Size()
	l.Draw(gc, left, b.Y+(p.head()-h)/2)
}

// paneLayout is the head, a breadcrumb's height, with the close control at
// its right, a rule under it, and the body below.
type paneLayout struct{ p *Pane }

func (l paneLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.p.ui.Interface
	w := float32(m.PaneWidth())
	_, bp, _ := l.p.body.Sizes(geom.NewSize(w, 0))
	h := l.p.head() + float32(m.Hairline()) + bp.Height
	return geom.NewSize(w, l.p.head()), geom.NewSize(w, h), geom.NewSize(w, unison.DefaultMaxSize)
}

func (l paneLayout) PerformLayout(target *unison.Panel) {
	p, m := l.p, l.p.ui.Interface
	b := target.ContentRect(false)
	_, cp, _ := p.close.Sizes(geom.Size{})
	p.close.SetFrameRect(geom.NewRect(b.Right()-float32(m.SpaceNear())-cp.Width, b.Y+(p.head()-cp.Height)/2, cp.Width, cp.Height))
	top := b.Y + p.head() + float32(m.Hairline())
	p.body.SetFrameRect(geom.NewRect(b.X, top, b.Width, max(0, b.Bottom()-top)))
}

// ProvideAccessibility names the pane by its title.
func (p *Pane) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name = p.Title
}

// WithPane lays a view out with a pane over its right edge, slid in as far as
// the pane is open. The view keeps its whole width either way, so nothing in
// it moves when the pane comes and goes. under may be nil.
func WithPane(ui *UI, under unison.Paneler, pane *Pane) *unison.Panel {
	c := unison.NewPanel()
	c.AddChild(pane)
	if under != nil {
		c.AddChild(under)
	}
	c.SetLayout(withPaneLayout{under, pane})
	return c
}

type withPaneLayout struct {
	under unison.Paneler
	pane  *Pane
}

func (l withPaneLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	_, pp, _ := l.pane.Sizes(geom.Size{})
	prefSize = pp
	if l.under != nil {
		_, up, _ := l.under.AsPanel().Sizes(hint)
		prefSize = geom.NewSize(max(up.Width, pp.Width), max(up.Height, pp.Height))
	}
	return geom.Size{}, prefSize, geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l withPaneLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	if l.under != nil {
		l.under.AsPanel().SetFrameRect(r)
	}
	w := float32(l.pane.ui.Interface.PaneWidth())
	l.pane.Hidden = l.pane.shown == 0
	l.pane.SetFrameRect(geom.NewRect(r.Right()-w*l.pane.shown, r.Y, w, r.Height))
}
