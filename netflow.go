package kvitui

import (
	"math"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/pathop"
	"github.com/richardwilkes/unison/enums/role"
)

// NetFlow is two opposed series over the same periods, on one baseline and
// one scale, with what they come to drawn as a line: what came in is drawn up
// from the baseline and what went out is drawn down, so which was larger in a
// period is a glance at where the two bars of a column meet. Two Sparks, one
// over the other, cannot say that however they are scaled, since each is
// drawn up from its own bottom edge and reading one against the other is
// comparing two bars an inch apart with no common line. Each column names its
// period underneath, and the names thin out to whatever spacing they fit in,
// counted back from the newest column so the end of the series is always
// named. It is kvit-cash's addition to the Qt library (its Cash flow widget).
type NetFlow struct {
	unison.Panel
	ui *UI
	// Ins and Outs are what came in and went out over each period, indexed
	// together; NotMeasured is a period nobody measured, drawn as a tick on
	// the baseline rather than as a zero.
	Ins, Outs []float64
	// Periods name the columns, one each, in the same order.
	Periods []string
	// Maximum is the scale both directions are drawn against; 0 takes the
	// largest value present.
	Maximum float64
	// RisingLabel, FallingLabel and NetLabel name the parts in a key above
	// the plot; all empty draws no key.
	RisingLabel, FallingLabel, NetLabel string
	// RisingInk and FallingInk are the bars' colours: the first two
	// categorical colours unless set, never a meaning colour by default,
	// since success and danger mean finished and stalled everywhere else. A
	// caller whose directions are good and bad news says so. NetInk is the
	// line's; the primary text colour unless set.
	RisingInk, FallingInk, NetInk Ink
	// Label names the chart for a screen reader.
	Label string
	// PlotHeight is the plot's height without the key and the names; 64
	// design pixels unless set.
	PlotHeight Measure
}

// NewNetFlow returns a chart of what came in and went out over periods.
func NewNetFlow(ui *UI, ins, outs []float64, periods ...string) *NetFlow {
	n := &NetFlow{ui: ui, Ins: ins, Outs: outs, Periods: periods, PlotHeight: Px(64)}
	n.Self = n
	n.SetSizer(func(hint geom.Size) (geom.Size, geom.Size, geom.Size) {
		w := hint.Width
		if w <= 0 {
			w = float32(ui.Interface.Px(240))
		}
		h := n.height(w)
		return geom.NewSize(0, h), geom.NewSize(float32(ui.Interface.Px(240)), h), geom.NewSize(unison.DefaultMaxSize, h)
	})
	n.DrawCallback = n.draw
	return n
}

func (n *NetFlow) columns() int { return max(len(n.Ins), len(n.Outs)) }

func (n *NetFlow) scale() float64 {
	if n.Maximum > 0 {
		return n.Maximum
	}
	top := 0.0
	for _, series := range [][]float64{n.Ins, n.Outs} {
		for _, v := range series {
			if measured(v) && v > top {
				top = v
			}
		}
	}
	if top > 0 {
		return top
	}
	return 1
}

// share is a measured value as a share of the scale, held to it, or -1.
func (n *NetFlow) share(series []float64, i int) float64 {
	if i >= len(series) || !measured(series[i]) {
		return -1
	}
	return min(1, max(0, series[i])/n.scale())
}

func (n *NetFlow) inks() (rising, falling, net Ink) {
	rising, falling, net = n.RisingInk, n.FallingInk, n.NetInk
	if rising == nil {
		rising = InkCategorical(0)
	}
	if falling == nil {
		falling = InkCategorical(1)
	}
	if net == nil {
		net = InkTextPrimary
	}
	return rising, falling, net
}

// keyPart is one entry of the key: a block for a bar or a rule for the line,
// and its word.
type keyPart struct {
	words *text.Layout
	ink   Ink
	line  bool
	box   geom.Rect // where it goes, relative to the chart's top left
}

// key lays the key out in a width, wrapping rather than running off the end,
// and says how tall it is.
func (n *NetFlow) key(width float32) ([]keyPart, float32) {
	ui, m := n.ui, n.ui.Interface
	rising, falling, net := n.inks()
	var parts []keyPart
	for _, p := range []struct {
		words string
		ink   Ink
		line  bool
	}{{n.RisingLabel, rising, false}, {n.FallingLabel, falling, false}, {n.NetLabel, net, true}} {
		if p.words == "" {
			continue
		}
		l := ui.Fonts.Layout([]text.Span{{Text: p.words, Style: ui.Chrome(ui.Size(RoleCaption), text.Regular, ui.Theme.Tokens().TextMuted)}}, text.Options{})
		parts = append(parts, keyPart{words: l, ink: p.ink, line: p.line})
	}
	if len(parts) == 0 {
		return nil, 0
	}
	x, y, line := float32(0), float32(0), float32(0)
	for i := range parts {
		w, h := parts[i].words.Size()
		swatch := float32(m.SpaceNear())
		if parts[i].line {
			swatch = float32(m.SpaceLoose())
		}
		total := swatch + float32(m.SpaceSnug()) + w
		if x > 0 && x+total > width {
			x, y = 0, y+line+float32(m.Space())
			line = 0
		}
		parts[i].box = geom.NewRect(x, y, total, h)
		line = max(line, h)
		x += total + float32(m.Space())
	}
	return parts, y + line
}

func (n *NetFlow) captionHeight() float32 {
	ui := n.ui
	_, h := ui.Fonts.Layout([]text.Span{{Text: "0", Style: ui.Chrome(ui.Size(RoleCaption), text.Regular, ui.Theme.Tokens().TextFaint)}}, text.Options{}).Size()
	return h
}

func (n *NetFlow) height(width float32) float32 {
	m := n.ui.Interface
	h := float32(n.PlotHeight.Of(n.ui))
	if _, kh := n.key(width); kh > 0 {
		h += kh + float32(m.SpaceSnug())
	}
	if len(n.Periods) > 0 {
		h += float32(m.SpaceSnug()) + n.captionHeight()
	}
	return h
}

// slot is where a column starts and how wide it is. The boundaries are
// rounded rather than the widths, so the gaps come out even.
func (n *NetFlow) slot(width float32, i int) (x, w float32) {
	snug := float64(n.ui.Interface.SpaceSnug())
	count := float64(max(1, n.columns()))
	start := func(i int) float32 { return float32(math.Round(float64(i) * (float64(width) + snug) / count)) }
	return start(i), max(1, start(i+1)-start(i)-float32(snug))
}

func (n *NetFlow) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := n.ui, n.ui.Theme.Tokens(), n.ui.Interface
	p := painterFor(gc, ui)
	r := n.ContentRect(false)
	rising, falling, net := n.inks()
	hair := float32(m.Hairline())
	y := r.Y
	parts, kh := n.key(r.Width)
	for _, part := range parts {
		b := part.box
		_, h := part.words.Size()
		if part.line {
			p.fill(geom.NewRect(r.X+b.X, r.Y+b.Y+(h-2*hair)/2, float32(m.SpaceLoose()), 2*hair), part.ink.Of(ui))
			part.words.Draw(gc, r.X+b.X+float32(m.SpaceLoose()+m.SpaceSnug()), r.Y+b.Y)
		} else {
			s := float32(m.SpaceNear())
			p.round(geom.NewRect(r.X+b.X, r.Y+b.Y+(h-s)/2, s, s), float32(m.RadiusBar()), part.ink.Of(ui))
			part.words.Draw(gc, r.X+b.X+s+float32(m.SpaceSnug()), r.Y+b.Y)
		}
	}
	if kh > 0 {
		y += kh + float32(m.SpaceSnug())
	}
	plotH := float32(n.PlotHeight.Of(ui))
	plot := geom.NewRect(r.X, y, r.Width, plotH)
	gc.Save()
	gc.ClipRect(plot, pathop.Intersect, false)
	half := max(1, float32(math.Floor(float64(plotH-hair)/2)))
	middle := plot.Y + half
	p.fill(geom.NewRect(plot.X, middle, plot.Width, hair), t.Border)
	type point struct {
		x, y float32
		ok   bool
	}
	points := make([]point, n.columns())
	for i := range n.columns() {
		sx, sw := n.slot(plot.Width, i)
		bar := min(sw, float32(m.Px(20)))
		bx := plot.X + sx + float32(math.Round(float64(sw-bar)/2))
		up, down := n.share(n.Ins, i), n.share(n.Outs, i)
		if up >= 0 {
			h := max(hair, float32(math.Round(float64(half)*up)))
			p.round(geom.NewRect(bx, middle-h, bar, h), float32(m.RadiusBar()), rising.Of(ui))
		}
		if down >= 0 {
			h := max(hair, float32(math.Round(float64(half)*down)))
			p.round(geom.NewRect(bx, middle+hair, bar, h), float32(m.RadiusBar()), falling.Of(ui))
		}
		if up < 0 && down < 0 {
			// A period nobody measured: a muted tick on the baseline, an
			// absence rather than a period where nothing happened.
			p.fill(geom.NewRect(bx, middle, bar, hair), t.TextFaint)
			continue
		}
		points[i] = point{plot.X + sx + sw/2, middle - half*float32(max(0, up)-max(0, down)), true}
	}
	// What the two come to, a segment between each pair of measured columns
	// and a dot at each, broken where a column has no net.
	ink := net.Of(ui)
	paint := Color(ink).Paint(gc, plot, paintstyle.Stroke)
	paint.SetStrokeWidth(2 * hair)
	for i := 0; i+1 < len(points); i++ {
		if points[i].ok && points[i+1].ok {
			gc.DrawLine(geom.NewPoint(points[i].x, points[i].y), geom.NewPoint(points[i+1].x, points[i+1].y), paint)
		}
	}
	dot := 3 * hair
	for _, pt := range points {
		if pt.ok {
			p.round(geom.NewRect(pt.x-dot/2, pt.y-dot/2, dot, dot), dot/2, ink)
		}
	}
	gc.Restore()
	if len(n.Periods) == 0 {
		return
	}
	// The names, at whatever spacing the longest fits in.
	st := ui.Chrome(ui.Size(RoleCaption), text.Regular, t.TextFaint)
	longest := float32(0)
	for _, name := range n.Periods {
		w, _ := ui.Fonts.Layout([]text.Span{{Text: name, Style: st}}, text.Options{}).Size()
		longest = max(longest, w)
	}
	cols := n.columns()
	step := 1
	if cols > 0 {
		pitch := max(1, r.Width/float32(cols))
		step = max(1, int(math.Ceil(float64((longest+float32(m.SpaceSnug()))/pitch))))
	}
	ny := plot.Bottom() + float32(m.SpaceSnug())
	for i := range cols {
		if i >= len(n.Periods) || (cols-1-i)%step != 0 {
			continue
		}
		l := ui.Fonts.Layout([]text.Span{{Text: n.Periods[i], Style: st}}, text.Options{})
		w, _ := l.Size()
		sx, sw := n.slot(r.Width, i)
		l.Draw(gc, r.X+float32(math.Round(float64(sx+sw/2-w/2))), ny)
	}
}

// ProvideAccessibility names the chart and says how many periods it holds
// and which way each direction is drawn.
func (n *NetFlow) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	node := b.Node()
	if node.Role == role.Auto {
		node.Role = role.Image
	}
	node.Name = n.Label
	count := n.columns()
	periods := n.ui.Number(count) + " periods"
	if count == 1 {
		periods = "1 period"
	}
	node.Description = periods + ", what came in above the line and what went out below it"
}
