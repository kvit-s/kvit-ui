package kvitui

import (
	"math"
	"strconv"
	"strings"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/pathop"
	"github.com/richardwilkes/unison/enums/role"
	"github.com/richardwilkes/unison/enums/strokecap"
	"github.com/richardwilkes/unison/enums/strokejoin"
)

// NotMeasured marks a period, in a series of values, that nobody measured.
// It is drawn as an absence rather than as a zero, which would be a claim.
var NotMeasured = math.NaN()

// measured reports whether a value in a series was measured.
func measured(v float64) bool { return !math.IsNaN(v) }

// InkCategorical is categorical colour i of the current theme, for the
// segments of a breakdown.
func InkCategorical(i int) Ink {
	return func(t tokens.Tokens) palette.Color { return t.Categorical(i) }
}

// Caption is words in the caption role in a colour, as spans for the text
// layer, for drawing text over a chart.
func Caption(ui *UI, words string, ink Ink) []text.Span {
	return []text.Span{{Text: words, Style: ui.Chrome(ui.Size(RoleCaption), text.Regular, ink.Of(ui))}}
}

// NoOptions lays text out on one line at its natural width.
var NoOptions = text.Options{}

// num writes a number as briefly as it can be written exactly.
func num(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

// fixedSizer is a sizer of a design size that follows the interface size.
func fixedSizer(ui *UI, w, h Measure) func(geom.Size) (geom.Size, geom.Size, geom.Size) {
	return func(hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
		width, height := orZero(ui, w), orZero(ui, h)
		return geom.NewSize(0, height), geom.NewSize(width, height), geom.NewSize(unison.DefaultMaxSize, height)
	}
}

// Bar is one quantity against a stated scale. A bar with no value is drawn
// as its track and a tick where it would start, not as a zero-width fill,
// which cannot be told from a measured zero. A value that is an upper bound
// is hatched, which is geometry and survives grayscale. And the scale is
// stated rather than taken from the largest value around it, so two
// screenshots a day apart can be compared.
type Bar struct {
	unison.Panel
	ui *UI
	// Value against Maximum, the scale.
	Value, Maximum float64
	// Measured is false for a value nobody measured.
	Measured bool
	// Bounded marks the value as "at most this much".
	Bounded bool
	// Ink is the fill; the accent unless set.
	Ink Ink
	// Wide draws the taller form.
	Wide bool
	// Label and Unit say what the bar is, for a screen reader.
	Label, Unit string
}

// NewBar returns a measured bar against a scale.
func NewBar(ui *UI, value, maximum float64) *Bar {
	b := &Bar{ui: ui, Value: value, Maximum: maximum, Measured: true}
	b.Self = b
	b.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		h := SizeBarHeight
		if b.Wide {
			h = SizeBarHeightWide
		}
		return fixedSizer(ui, Px(120), h)(geom.Size{})
	})
	b.DrawCallback = b.draw
	return b
}

func (b *Bar) fraction() float32 {
	if b.Maximum <= 0 {
		return 0
	}
	return float32(max(0, min(1, b.Value/b.Maximum)))
}

func (b *Bar) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := b.ui, b.ui.Theme.Tokens(), b.ui.Interface
	p := painterFor(gc, ui)
	r := b.ContentRect(false)
	radius := float32(m.RadiusBar())
	p.round(r, radius, t.ChipBackground)
	if !b.Measured {
		// An absence drawn as an absence: a tick at the origin.
		p.round(geom.NewRect(r.X, r.Y, float32(m.SpaceTight()), r.Height), radius, t.TextFaint)
		return
	}
	ink := b.Ink
	if ink == nil {
		ink = InkAccent
	}
	fill := geom.NewRect(r.X, r.Y, float32(math.Round(float64(r.Width*b.fraction()))), r.Height)
	p.round(fill, radius, ink.Of(ui))
	if b.Bounded && fill.Width > 0 {
		drawHatch(gc, ui, fill)
	}
}

// drawHatch draws thin diagonals over a box in the theme's hatch colour.
func drawHatch(gc *unison.Canvas, ui *UI, box geom.Rect) {
	m := ui.Interface
	gc.Save()
	gc.ClipRect(box, pathop.Intersect, false)
	paint := Color(ui.Theme.Tokens().HatchAlt).SetAlphaIntensity(0.7).Paint(gc, box, paintstyle.Stroke)
	paint.SetStrokeWidth(float32(m.Hairline()))
	step := float32(m.SpaceSnug())
	for x := box.X - box.Height; x < box.Right(); x += step {
		gc.DrawLine(geom.NewPoint(x, box.Bottom()), geom.NewPoint(x+box.Height, box.Y), paint)
	}
	gc.Restore()
}

// ProvideAccessibility reads the bar as its value against its scale.
func (b *Bar) ProvideAccessibility(ab *unison.AccessibilityBuilder) {
	n := ab.Node()
	if n.Role == role.Auto {
		n.Role = role.ProgressBar
	}
	n.Name = b.Label
	switch {
	case !b.Measured:
		n.Description = "not measured"
	case b.Bounded:
		n.Description = joinWords("at most", num(b.Value), "of", num(b.Maximum), b.Unit)
	default:
		n.Description = joinWords(num(b.Value), "of", num(b.Maximum), b.Unit)
	}
	if b.Measured {
		n.HasNumber = true
		n.Number, n.Min, n.Max = b.Value, 0, b.Maximum
	}
}

// Segment is one part of a StackedBar.
type Segment struct {
	Value float64
	// Ink is the segment's colour; categorical colour i unless set.
	Ink Ink
	// Label says what the segment is, for a screen reader.
	Label string
}

// StackedBar is several quantities adding up to one total, in one bar, a gap
// apart: two fills that touch read as one fill with a colour change in it.
// A segment too small to see is drawn at the smallest visible width rather
// than dropped, because a bar whose parts do not add up to it misstates its
// own arithmetic.
type StackedBar struct {
	unison.Panel
	ui *UI
	// Segments are the parts, left to right.
	Segments []Segment
	// Maximum is the scale; the sum of the segments unless set, which is
	// right for a bar that is a breakdown of its own total.
	Maximum float64
	// Wide draws the taller form.
	Wide bool
	// Label names the bar for a screen reader.
	Label string
}

// NewStackedBar returns a bar of segments.
func NewStackedBar(ui *UI, segments ...Segment) *StackedBar {
	s := &StackedBar{ui: ui, Segments: segments}
	s.Self = s
	s.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		h := SizeBarHeight
		if s.Wide {
			h = SizeBarHeightWide
		}
		return fixedSizer(ui, Px(160), h)(geom.Size{})
	})
	s.DrawCallback = s.draw
	return s
}

func (s *StackedBar) scale() float64 {
	if s.Maximum > 0 {
		return s.Maximum
	}
	var sum float64
	for _, g := range s.Segments {
		sum += g.Value
	}
	return sum
}

func (s *StackedBar) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := s.ui, s.ui.Theme.Tokens(), s.ui.Interface
	p := painterFor(gc, ui)
	r := s.ContentRect(false)
	radius := float32(m.RadiusBar())
	p.round(r, radius, t.ChipBackground)
	scale := s.scale()
	if scale <= 0 {
		return
	}
	gap := float32(m.SpaceTight())
	x := r.X
	for i, g := range s.Segments {
		if g.Value <= 0 {
			continue
		}
		w := max(gap, float32(math.Round(float64(r.Width)*g.Value/scale)))
		ink := g.Ink
		if ink == nil {
			ink = InkCategorical(i)
		}
		p.round(geom.NewRect(x, r.Y, w, r.Height), radius, ink.Of(ui))
		x += w + gap
	}
}

// ProvideAccessibility reads the bar as its parts, "12 food, 5 rent".
func (s *StackedBar) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Image
	}
	n.Name = s.Label
	var parts []string
	for _, g := range s.Segments {
		parts = append(parts, joinWords(num(g.Value), g.Label))
	}
	n.Description = strings.Join(parts, ", ")
}

// Spark is a series as a row of slim bars rising from one baseline: the
// recent shape of something, small enough to sit in a row. A period nobody
// measured is a hairline on the baseline, not a missing bar.
type Spark struct {
	unison.Panel
	ui *UI
	// Values are the periods, oldest first; NotMeasured for a gap.
	Values []float64
	// Ink is the bars' colour; the accent unless set.
	Ink Ink
	// Maximum is the scale; the largest value unless set.
	Maximum float64
	// Label names the series for a screen reader.
	Label string
}

// NewSpark returns a spark of values.
func NewSpark(ui *UI, values ...float64) *Spark {
	s := &Spark{ui: ui, Values: values}
	s.Self = s
	s.SetSizer(fixedSizer(ui, Px(80), SizeRowHeightCompact))
	s.DrawCallback = s.draw
	return s
}

func (s *Spark) scale() float64 {
	if s.Maximum > 0 {
		return s.Maximum
	}
	top := 0.0
	for _, v := range s.Values {
		if measured(v) && v > top {
			top = v
		}
	}
	if top > 0 {
		return top
	}
	return 1
}

func (s *Spark) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := s.ui, s.ui.Theme.Tokens(), s.ui.Interface
	p := painterFor(gc, ui)
	r := s.ContentRect(false)
	n := max(1, len(s.Values))
	gap := float32(m.Hairline())
	start := func(i int) float32 { return float32(math.Round(float64(float32(i) * (r.Width + gap) / float32(n)))) }
	ink := s.Ink
	if ink == nil {
		ink = InkAccent
	}
	for i, v := range s.Values {
		x := r.X + start(i)
		w := max(1, start(i+1)-start(i)-gap)
		if !measured(v) {
			p.fill(geom.NewRect(x, r.Bottom()-gap, w, gap), t.TextFaint)
			continue
		}
		h := max(gap, float32(math.Round(float64(r.Height)*v/s.scale())))
		p.round(geom.NewRect(x, r.Bottom()-h, w, h), float32(m.RadiusBar()), ink.Of(ui))
	}
}

// ProvideAccessibility names the series and how many periods it holds.
func (s *Spark) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Image
	}
	n.Name = s.Label
	n.Description = s.ui.CountPhrase(len(s.Values), "period", "")
}

// Distribution is a spread of values: the whole range as a thin line, the
// middle half as a box, and the median as a mark, against a stated scale.
type Distribution struct {
	unison.Panel
	ui *UI
	// The five numbers.
	Minimum, LowerQuartile, Median, UpperQuartile, Maximum float64
	// ScaleMinimum and ScaleMaximum are the scale's ends.
	ScaleMinimum, ScaleMaximum float64
	// Measured is false for a spread nobody measured, drawn as "—".
	Measured bool
	// Ink is the box's colour; the accent unless set.
	Ink Ink
	// Label and Unit say what it is, for a screen reader.
	Label, Unit string
}

// NewDistribution returns a measured spread on a scale from 0 to 1.
func NewDistribution(ui *UI) *Distribution {
	d := &Distribution{ui: ui, ScaleMaximum: 1, Measured: true}
	d.Self = d
	d.SetSizer(fixedSizer(ui, Px(160), SizeBarHeightWide))
	d.DrawCallback = d.draw
	return d
}

func (d *Distribution) at(r geom.Rect, v float64) float32 {
	span := d.ScaleMaximum - d.ScaleMinimum
	if span <= 0 {
		return r.X
	}
	return r.X + float32(math.Round(float64(r.Width)*max(0, min(1, (v-d.ScaleMinimum)/span))))
}

func (d *Distribution) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := d.ui, d.ui.Theme.Tokens(), d.ui.Interface
	p := painterFor(gc, ui)
	r := d.ContentRect(false)
	if !d.Measured {
		l := ui.Fonts.Layout([]text.Span{{Text: "—", Style: ui.Chrome(ui.Size(RoleCaption), text.Regular, t.TextFaint)}}, text.Options{})
		w, h := l.Size()
		l.Draw(gc, r.X+(r.Width-w)/2, r.Y+(r.Height-h)/2)
		return
	}
	hair := float32(m.Hairline())
	lo, hi := d.at(r, d.Minimum), d.at(r, d.Maximum)
	p.fill(geom.NewRect(lo, r.Y+(r.Height-hair)/2, max(hair, hi-lo), hair), t.TextFaint)
	ink := d.Ink
	if ink == nil {
		ink = InkAccent
	}
	q1, q3 := d.at(r, d.LowerQuartile), d.at(r, d.UpperQuartile)
	box := geom.NewRect(q1, r.Y, max(float32(m.SpaceTight()), q3-q1), r.Height)
	radius := float32(m.RadiusBar())
	p.roundTint(box, radius, ink.Of(ui), 0.35)
	p.outline(box, radius, hair, ink.Of(ui))
	p.fill(geom.NewRect(d.at(r, d.Median), r.Y, float32(m.SpaceTight()), r.Height), t.TextPrimary)
}

// ProvideAccessibility reads the spread as its five numbers.
func (d *Distribution) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Image
	}
	n.Name = d.Label
	if !d.Measured {
		n.Description = "not measured"
		return
	}
	n.Description = joinWords("median", num(d.Median), d.Unit) + ", middle half " + num(d.LowerQuartile) +
		" to " + num(d.UpperQuartile) + ", range " + num(d.Minimum) + " to " + num(d.Maximum)
}

// Gauge is spending against an allowance: green while under, amber when
// ahead of the pace the period is at, red past the allowance with a bright
// end cap, so being over is a shape as well as a colour.
type Gauge struct {
	unison.Panel
	ui *UI
	// Value against Allowance.
	Value, Allowance float64
	// Pace is how far through the period it is, 0 to 1, drawn as a tick;
	// below zero draws none.
	Pace float64
	// Measured is false for a value nobody measured.
	Measured bool
	// Label and Unit say what it is, for a screen reader.
	Label, Unit string
}

// NewGauge returns a gauge of a value against an allowance, with no pace.
func NewGauge(ui *UI, value, allowance float64) *Gauge {
	g := &Gauge{ui: ui, Value: value, Allowance: allowance, Pace: -1, Measured: true}
	g.Self = g
	g.SetSizer(fixedSizer(ui, Px(160), SizeBarHeightWide))
	g.DrawCallback = g.draw
	return g
}

func (g *Gauge) over() bool { return g.Allowance > 0 && g.Value > g.Allowance }

func (g *Gauge) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := g.ui, g.ui.Theme.Tokens(), g.ui.Interface
	p := painterFor(gc, ui)
	r := g.ContentRect(false)
	radius := float32(m.RadiusBar())
	p.round(r, radius, t.ChipBackground)
	if !g.Measured {
		p.round(geom.NewRect(r.X, r.Y, float32(m.SpaceTight()), r.Height), radius, t.TextFaint)
		return
	}
	fraction := 0.0
	if g.Allowance > 0 {
		fraction = max(0, min(1, g.Value/g.Allowance))
	}
	tone := t.Success
	switch {
	case g.over():
		tone = t.Danger
	case g.Pace >= 0 && fraction > g.Pace+0.05:
		tone = t.Warning
	}
	p.round(geom.NewRect(r.X, r.Y, float32(math.Round(float64(r.Width)*fraction)), r.Height), radius, tone)
	if g.over() {
		snug := float32(m.SpaceSnug())
		p.fill(geom.NewRect(r.Right()-snug, r.Y, snug, r.Height), t.DangerBright)
	}
	if !g.ui.spillHosted(g.AsPanel()) {
		g.drawSpill(gc)
	}
}

// drawSpill draws the pace tick, which stands a tight space above the bar and
// a snug space further down than its top, past the gauge's box.
func (g *Gauge) drawSpill(gc *unison.Canvas) {
	if !g.Measured || g.Pace < 0 || g.Pace > 1 {
		return
	}
	m := g.ui.Interface
	r := g.ContentRect(false)
	x := r.X + float32(math.Round(float64(r.Width)*g.Pace))
	tight := float32(m.SpaceTight())
	painterFor(gc, g.ui).fill(geom.NewRect(x, r.Y-tight, float32(m.Hairline()), r.Height+float32(m.SpaceSnug())), g.ui.Theme.Tokens().TextPrimary)
}

// ProvideAccessibility reads the gauge as its value, allowance and how far
// over it is.
func (g *Gauge) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.ProgressBar
	}
	n.Name = g.Label
	switch {
	case !g.Measured:
		n.Description = "not measured"
	case g.over():
		n.Description = joinWords(num(g.Value), "of", num(g.Allowance), g.Unit) + ", over by " + num(g.Value-g.Allowance)
	default:
		n.Description = joinWords(num(g.Value), "of", num(g.Allowance), g.Unit)
	}
}

// Direction is which way a change is good.
type Direction int

// The directions.
const (
	// Neither says a change is neither good nor bad.
	Neither Direction = iota
	Up
	Down
)

// Delta is how much something changed and in which direction, carried by an
// arrow as well as a colour: green-up and red-down is the most common place a
// screen rests a meaning on hue, and the pair a deuteranope cannot separate.
// Whether up is good is the caller's to say, since a rise in spending and a
// rise in savings are the same arrow in opposite colours.
type Delta struct {
	unison.Panel
	ui *UI
	// Change is the change.
	Change float64
	// Unit follows it.
	Unit string
	// Measured is false for a change nobody measured.
	Measured bool
	// Good is the direction that is good; Neither draws the muted colour.
	Good Direction
	// Precision is how many decimal places to draw.
	Precision int

	figure *Figure
}

// NewDelta returns a measured change.
func NewDelta(ui *UI, change float64) *Delta {
	d := &Delta{ui: ui, Change: change, Measured: true}
	d.Self = d
	d.figure = NewFigure(ui, "", "")
	d.figure.Role = RoleSmall
	d.figure.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
	d.AddChild(d.figure)
	d.SetLayout(syncing{Layout: deltaLayout{d}, sync: d.sync})
	d.DrawCallback = d.draw
	return d
}

func (d *Delta) tone(t tokens.Tokens) palette.Color {
	if !d.Measured || d.Change == 0 || d.Good == Neither {
		return t.TextMuted
	}
	if (d.Change > 0) == (d.Good == Up) {
		return t.Success
	}
	return t.Danger
}

func (d *Delta) sync() {
	d.figure.Measured = d.Measured
	if d.Change == 0 {
		d.figure.Value, d.figure.Unit = "no change", ""
	} else {
		d.figure.Value, d.figure.Unit = strconv.FormatFloat(math.Abs(d.Change), 'f', d.Precision, 64), d.Unit
	}
	d.figure.Ink = func(t tokens.Tokens) palette.Color { return d.tone(t) }
}

func (d *Delta) arrowShown() bool { return d.Measured && d.Change != 0 }

type deltaLayout struct{ d *Delta }

func (l deltaLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.d.ui.Interface
	_, fp, _ := l.d.figure.Sizes(geom.Size{})
	w := fp.Width
	if l.d.arrowShown() {
		w += float32(m.Caption() + m.SpaceSnug())
	}
	size := geom.NewSize(w, max(fp.Height, float32(m.Caption())))
	return size, size, size
}

func (l deltaLayout) PerformLayout(target *unison.Panel) {
	m := l.d.ui.Interface
	r := target.ContentRect(false)
	_, fp, _ := l.d.figure.Sizes(geom.Size{})
	x := r.X
	if l.d.arrowShown() {
		x += float32(m.Caption() + m.SpaceSnug())
	}
	l.d.figure.SetFrameRect(geom.NewRect(x, r.Y+(r.Height-fp.Height)/2, fp.Width, fp.Height))
}

func (d *Delta) draw(gc *unison.Canvas, _ geom.Rect) {
	if !d.arrowShown() {
		return
	}
	symbol := "arrow-down"
	if d.Change > 0 {
		symbol = "arrow-up"
	}
	s := float32(d.ui.Interface.Caption())
	r := d.ContentRect(false)
	if g, ok := icons.Glyph(symbol); ok {
		drawGlyph(gc, d.ui, g, s, d.tone(d.ui.Theme.Tokens()), geom.NewRect(r.X, r.Y+(r.Height-s)/2, s, s))
	}
}

// ProvideAccessibility says the change in words: "up 12 %", "unchanged".
func (d *Delta) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	abs := strconv.FormatFloat(math.Abs(d.Change), 'f', d.Precision, 64)
	switch {
	case !d.Measured:
		n.Name = "not measured"
	case d.Change == 0:
		n.Name = "unchanged"
	case d.Change > 0:
		n.Name = joinWords("up", abs, d.Unit)
	default:
		n.Name = joinWords("down", abs, d.Unit)
	}
}

// FigureBlock is a figure with its name under it: one number a reader is
// meant to take away. The number is above and larger, because a row of these
// is read across the numbers.
type FigureBlock struct {
	unison.Panel
	ui     *UI
	Figure *Figure
	name   *Label
	// Label names the figure.
	Label string
}

// NewFigureBlock returns a figure in the headline role with its name.
func NewFigureBlock(ui *UI, value, unit, label string) *FigureBlock {
	f := &FigureBlock{ui: ui, Label: label}
	f.Self = f
	f.Figure = NewFigure(ui, value, unit)
	f.Figure.Role = RoleHeadline
	f.name = NewLabel(ui, label)
	f.name.Role, f.name.Ink = RoleSmall, InkTextMuted
	for _, p := range []unison.Paneler{f.Figure, f.name} {
		p.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Start})
		f.AddChild(p)
	}
	f.SetLayout(syncing{Layout: &spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: SizeSpaceTight},
		sync: func() { f.name.Text = f.Label }})
	return f
}

// ProvideAccessibility reads the block as "Balance: 1,284 GBP".
func (f *FigureBlock) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	if f.Figure.Measured {
		n.Name = f.Label + ": " + joinWords(f.Figure.Value, f.Figure.Unit)
	} else {
		n.Name = f.Label + ": not measured"
	}
}

// StatTile is a card holding one figure, what it is, how it changed and its
// recent shape: the shape of most dashboard widgets. The parts are optional
// and their order fixed, because a dashboard of tiles is read one part at a
// time down a column, and a tile that orders them differently breaks that
// scan for every tile beside it.
type StatTile struct {
	Card
	// Label names the figure; Value and Unit are the figure.
	Label, Value, Unit string
	// Measured is false for a figure nobody measured.
	Measured bool
	// HasChange draws Change with the Good direction.
	HasChange bool
	Change    float64
	Good      Direction
	// History is the recent shape; NotMeasured for a gap.
	History []float64
	// Caption is a line under it all.
	Caption string

	name    *Label
	figure  *Figure
	delta   *Delta
	spark   *Spark
	caption *Label
	line    *unison.Panel // the figure and the change
}

// NewStatTile returns a tile for a labelled figure.
func NewStatTile(ui *UI, label, value, unit string) *StatTile {
	s := &StatTile{Label: label, Value: value, Unit: unit, Measured: true}
	s.name = NewLabel(ui, "")
	s.name.Role, s.name.Ink = RoleSmall, InkTextMuted
	s.figure = NewFigure(ui, "", "")
	s.figure.Role = RoleDisplay
	s.delta = NewDelta(ui, 0)
	s.spark = NewSpark(ui)
	s.caption = NewLabel(ui, "")
	s.caption.Role, s.caption.Ink = RoleCaption, InkTextFaint
	room := unison.NewPanel()
	room.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
	s.line = FullWidth(Row(ui, SizeSpace, s.figure, room, s.delta))
	s.Padding = SizeSpaceLoose
	s.Self = s
	s.initCard(ui, s.name, s.line, s.spark, s.caption)
	s.SetLayout(syncing{Layout: widthAtLeast{AtLeast(ui, Px(96), &spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: SizeSpaceNear}), ui, Px(200)}, sync: s.sync})
	return s
}

func (s *StatTile) sync() {
	s.SetFocusable(s.Interactive)
	s.name.Text = s.Label
	s.figure.Value, s.figure.Unit, s.figure.Measured = s.Value, s.Unit, s.Measured
	s.delta.Hidden = !s.HasChange
	s.delta.Change, s.delta.Good, s.delta.Measured = s.Change, s.Good, s.Measured
	s.spark.Values, s.spark.Label = s.History, s.Label+" over time"
	s.spark.Hidden = len(s.History) == 0
	s.caption.Text = s.Caption
	s.caption.Hidden = s.Caption == ""
	showOnly(s.AsPanel(), s.name, s.line, s.spark, s.caption)
}

// ProvideAccessibility reads the tile as "Spending: 1,284 GBP", holding its
// parts.
func (s *StatTile) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	s.Card.ProvideAccessibility(b)
	n := b.Node()
	n.Name = s.Label + ": " + joinWords(s.Value, s.Unit)
	if !s.Measured {
		n.Name = s.Label + ": not measured"
	}
}

// widthAtLeast wraps a layout so the panel prefers at least a width.
type widthAtLeast struct {
	unison.Layout
	ui *UI
	w  Measure
}

func (l widthAtLeast) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	minSize, prefSize, maxSize = l.Layout.LayoutSizes(target, hint)
	prefSize.Width = max(prefSize.Width, orZero(l.ui, l.w))
	return minSize, prefSize, maxSize
}

// Trend is a line over periods against a stated scale, with gridlines and
// their values down the left, and an optional second series drawn dashed so
// the two do not rest on hue alone. Under the pointer it shows the value of
// the period there. With no periods it says so rather than drawing an empty
// axis.
type Trend struct {
	unison.Panel
	ui *UI
	// Points are the values, oldest first; NotMeasured for a gap.
	Points []float64
	// MinimumY and MaximumY are the scale's ends.
	MinimumY, MaximumY float64
	// Ink is the line's colour; the accent unless set.
	Ink Ink
	// Label and Unit say what the line is.
	Label, Unit string
	// Gridlines is how many horizontal lines, at least two; 3 unless set.
	Gridlines int
	// Second is an optional second series, dashed, in SecondInk (the
	// second categorical colour unless set), named SecondLabel.
	Second      []float64
	SecondInk   Ink
	SecondLabel string
	// Overlay draws over the plot after the lines, given the plot's box, so
	// an annotation lines up with the same arithmetic the lines use.
	Overlay func(gc *unison.Canvas, plot geom.Rect)

	hovered int // the period under the pointer, or -1
	empty   *EmptyState
}

// NewTrend returns a line of points on a scale from 0 to 1.
func NewTrend(ui *UI, points ...float64) *Trend {
	t := &Trend{ui: ui, Points: points, MaximumY: 1, Gridlines: 3, hovered: -1}
	t.Self = t
	t.empty = NewEmptyState(ui, "Nothing recorded yet")
	t.AddChild(t.empty)
	t.SetLayout(syncing{Layout: trendLayout{t}, sync: func() {
		t.empty.Hidden = t.periods() > 0
		name := t.Label
		if name == "" {
			name = "The series"
		}
		t.empty.Detail = name + " will appear here once there is something to plot."
	}})
	t.DrawCallback = t.draw
	t.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool { t.hover(where); return false }
	t.MouseEnterCallback = func(where geom.Point, _ mod.Modifiers) bool { t.hover(where); return false }
	t.MouseExitCallback = func() bool {
		t.hovered = -1
		t.MarkForRedraw()
		return false
	}
	return t
}

func (t *Trend) periods() int { return max(len(t.Points), len(t.Second)) }

// plot is the area the lines are drawn in: right of the axis values, and
// under the key when there are two series.
func (t *Trend) plot() geom.Rect {
	r := t.ContentRect(false)
	m := t.ui.Interface
	axis := float32(m.Px(36))
	top := float32(0)
	if len(t.Second) > 0 {
		_, h := t.caption("Hg").Size()
		top = h + float32(m.SpaceSnug())
	}
	return geom.NewRect(r.X+axis, r.Y+top, max(1, r.Width-axis), max(1, r.Height-top))
}

func (t *Trend) caption(s string) *text.Layout {
	ui := t.ui
	st := ui.Chrome(ui.Size(RoleCaption), text.Regular, ui.Theme.Tokens().TextMuted)
	return ui.Fonts.Layout([]text.Span{{Text: s, Style: st}}, text.Options{})
}

func (t *Trend) hover(where geom.Point) {
	p := t.plot()
	n := t.periods()
	h := -1
	if n > 0 {
		h = int(math.Round(float64((where.X - p.X) / p.Width * float32(n-1))))
		h = max(0, min(n-1, h))
	}
	if h != t.hovered {
		t.hovered = h
		t.MarkForRedraw()
	}
}

type trendLayout struct{ t *Trend }

func (l trendLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	return fixedSizer(l.t.ui, Px(320), Px(120))(geom.Size{})
}

func (l trendLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	_, p, _ := l.t.empty.Sizes(geom.NewSize(r.Width, 0))
	l.t.empty.SetFrameRect(geom.NewRect(r.X, r.Y+(r.Height-p.Height)/2, r.Width, p.Height))
}

func (t *Trend) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := t.ui, t.ui.Theme.Tokens(), t.ui.Interface
	pp := painterFor(gc, ui)
	r := t.ContentRect(false)
	plot := t.plot()
	lines := max(2, t.Gridlines)
	for i := range lines {
		if y := t.gridline(plot, i, lines); y < r.Bottom() {
			pp.fill(geom.NewRect(plot.X, y, plot.Width, float32(m.Hairline())), tk.Border)
		}
	}
	if !ui.spillHosted(t.AsPanel()) {
		t.drawSpill(gc)
	}
	n := t.periods()
	if n == 0 {
		return
	}
	span := t.MaximumY - t.MinimumY
	toY := func(v float64) float32 {
		if span <= 0 {
			return plot.Bottom()
		}
		return plot.Bottom() - float32((v-t.MinimumY)/span)*plot.Height
	}
	toX := func(i int) float32 {
		if n == 1 {
			return plot.X + plot.Width/2
		}
		return plot.X + float32(i)/float32(n-1)*plot.Width
	}
	stroke := func(series []float64, c palette.Color, dashed bool) {
		paint := Color(c).Paint(gc, plot, paintstyle.Stroke)
		paint.SetStrokeWidth(max(1, float32(m.Px(2))))
		paint.SetStrokeCap(strokecap.Round)
		paint.SetStrokeJoin(strokejoin.Round)
		if dashed {
			paint.SetPathEffect(unison.NewDashPathEffect([]float32{float32(m.Px(5)), float32(m.Px(4))}, 0))
		}
		path := unison.NewPath()
		drawing := false
		for i := range n {
			if i >= len(series) || !measured(series[i]) {
				drawing = false
				continue
			}
			pt := geom.NewPoint(toX(i), toY(series[i]))
			if drawing {
				path.LineTo(pt)
			} else {
				path.MoveTo(pt)
				drawing = true
			}
		}
		gc.DrawPath(path, paint)
	}
	secondInk := t.SecondInk
	if secondInk == nil {
		secondInk = InkCategorical(1)
	}
	ink := t.Ink
	if ink == nil {
		ink = InkAccent
	}
	if len(t.Second) > 0 {
		stroke(t.Second, secondInk.Of(ui), true)
		t.drawKey(gc, r, ink.Of(ui), secondInk.Of(ui))
	}
	stroke(t.Points, ink.Of(ui), false)
	if t.Overlay != nil {
		gc.Save()
		t.Overlay(gc, plot)
		gc.Restore()
	}
	if t.hovered >= 0 {
		x := toX(t.hovered)
		pp.fill(geom.NewRect(x, plot.Y, float32(m.Hairline()), plot.Height), tk.TextMuted)
		t.drawReadout(gc, plot, x, ink.Of(ui), secondInk.Of(ui))
	}
}

// gridline is the height of gridline i of lines, the first at the top.
func (t *Trend) gridline(plot geom.Rect, i, lines int) float32 {
	return plot.Y + float32(math.Round(float64(plot.Height*float32(i)/float32(lines-1))))
}

// drawSpill draws each gridline's value down the left, centred on its line,
// so the top and bottom values stand half a line past the trend's box, and
// the bottom gridline, which lies on the line below the box, as the Qt trend
// places it. Values are rounded half away from zero, as the Qt trend's
// toFixed rounds them.
func (t *Trend) drawSpill(gc *unison.Canvas) {
	ui, m := t.ui, t.ui.Interface
	r, plot := t.ContentRect(false), t.plot()
	lines := max(2, t.Gridlines)
	st := ui.Chrome(ui.Size(RoleCaption), text.Regular, ui.Theme.Tokens().TextFaint)
	st.Tabular = true
	for i := range lines {
		if y := t.gridline(plot, i, lines); y >= r.Bottom() {
			painterFor(gc, ui).fill(geom.NewRect(plot.X, y, plot.Width, float32(m.Hairline())), ui.Theme.Tokens().Border)
		}
		value := t.MaximumY - (t.MaximumY-t.MinimumY)*float64(i)/float64(lines-1)
		l := ui.Fonts.Layout([]text.Span{{Text: toFixed(value, 0), Style: st}},
			text.Options{MaxWidth: plot.X - r.X - float32(m.SpaceSnug()), Align: text.AlignEnd})
		_, h := l.Size()
		l.Draw(gc, r.X, t.gridline(plot, i, lines)-h/2)
	}
}

// valueAt is a series' value at the hovered period, in words.
func (t *Trend) valueAt(series []float64) string {
	if t.hovered < 0 || t.hovered >= len(series) || !measured(series[t.hovered]) {
		return "not measured"
	}
	return joinWords(num(series[t.hovered]), t.Unit)
}

func (t *Trend) drawReadout(gc *unison.Canvas, plot geom.Rect, x float32, first, second palette.Color) {
	ui, tk, m := t.ui, t.ui.Theme.Tokens(), t.ui.Interface
	lines := []struct {
		words string
		c     palette.Color
	}{{t.valueAt(t.Points), tk.TextPrimary}}
	if len(t.Second) > 0 {
		lines[0] = struct {
			words string
			c     palette.Color
		}{t.Label + ": " + t.valueAt(t.Points), first}
		lines = append(lines, struct {
			words string
			c     palette.Color
		}{t.SecondLabel + ": " + t.valueAt(t.Second), second})
	}
	y := plot.Y
	for _, ln := range lines {
		st := ui.Chrome(ui.Size(RoleCaption), text.Regular, ln.c)
		st.Tabular = true
		l := ui.Fonts.Layout([]text.Span{{Text: ln.words, Style: st}}, text.Options{})
		w, h := l.Size()
		lx := min(plot.Right()-w, x+float32(m.SpaceSnug()))
		l.Draw(gc, lx, y)
		y += h + float32(m.SpaceTight())
	}
}

// drawKey names the two series above the plot, each beside a sample of its
// line: solid for the first, dashed for the second.
func (t *Trend) drawKey(gc *unison.Canvas, r geom.Rect, first, second palette.Color) {
	ui, m := t.ui, t.ui.Interface
	p := painterFor(gc, ui)
	x := t.plot().X
	thick := max(1, float32(m.Px(2)))
	for i, e := range []struct {
		words string
		c     palette.Color
	}{{t.Label, first}, {t.SecondLabel, second}} {
		l := t.caption(e.words)
		w, h := l.Size()
		mid := r.Y + h/2
		if i == 0 {
			p.round(geom.NewRect(x, mid-thick/2, float32(m.Px(15)), thick), float32(m.RadiusBar()), e.c)
			x += float32(m.Px(15))
		} else {
			for j := range 2 {
				p.round(geom.NewRect(x+float32(j*m.Px(9)), mid-thick/2, float32(m.Px(6)), thick), float32(m.RadiusBar()), e.c)
			}
			x += float32(m.Px(15))
		}
		x += float32(m.SpaceSnug())
		l.Draw(gc, x, r.Y)
		x += w + float32(m.SpaceLoose())
	}
}

// ProvideAccessibility names the line, its periods and its scale.
func (t *Trend) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Image
	}
	n.Name = t.Label
	scale := num(t.MinimumY) + " to " + joinWords(num(t.MaximumY), t.Unit)
	if len(t.Second) > 0 {
		n.Description = t.ui.CountPhrase(t.periods(), "period", "") + ", " + scale + ", two series: " + t.Label + " and " + t.SecondLabel
		return
	}
	n.Description = t.ui.CountPhrase(len(t.Points), "point", "") + ", " + scale
}
