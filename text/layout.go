package text

import (
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/go-text/typesetting/di"
	gfont "github.com/go-text/typesetting/font"
	ot "github.com/go-text/typesetting/font/opentype"
	"github.com/go-text/typesetting/fontscan"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"
)

// Weights, on the usual 100–900 scale.
const (
	Regular  = 400
	Medium   = 500
	Semibold = 600
	Bold     = 700
)

// Color is a colour as 8-bit red, green, blue and alpha. A zero Color is
// transparent, which is how an optional colour says "none".
type Color struct{ R, G, B, A uint8 }

// Style is how a span of text is drawn.
type Style struct {
	Family     string  // "" and SansSerif are the desktop's face, Monospace its fixed-pitch one
	Size       float32 // em size in pixels
	Weight     int     // 100–900; 0 means Regular
	Italic     bool
	Color      Color
	Background Color // drawn behind the span's glyphs on its line; zero for none
	Underline  bool
	Strike     bool
	// Tabular draws every digit the same width, so a column of figures lines
	// up and a changing value does not shift the text beside it.
	Tabular bool
}

// Span is text in one style.
type Span struct {
	Text  string
	Style Style
}

// Options shape a whole layout.
type Options struct {
	// MaxWidth wraps lines at this width; 0 never wraps, though line breaks in
	// the text still start new lines.
	MaxWidth float32
	// LineHeight multiplies each line's natural height; 0 means 1. The extra
	// space is shared above and below the text.
	LineHeight float32
	// KeepTrailingSpace keeps the width of spaces at the end of a wrapped
	// line, which an editor wants so the caret can sit after them.
	KeepTrailingSpace bool
	// Elide keeps the text to one line no wider than MaxWidth, cutting it
	// short with "…" when it does not fit. A label that grows past its column
	// pushes whatever is beside it off the screen; one cut short is legible.
	Elide bool
	// Align places each line in the width: at its start (the default), in
	// its middle, or at its end. The width is MaxWidth, or the widest line
	// when there is none.
	Align Alignment
}

// Alignment is where a line sits across the width of its layout.
type Alignment int

// The three alignments.
const (
	AlignStart Alignment = iota
	AlignMiddle
	AlignEnd
)

// Layout is text shaped and broken into lines, ready to draw and to answer
// caret and hit-test questions. Offsets are rune indexes into the text.
type Layout struct {
	runes  []rune
	starts []int // rune offset where each span starts
	styles []Style
	lines  []line
	width  float32
	height float32
	fonts  *Fonts
}

type line struct {
	runs          []placedRun // in visual order
	start, end    int         // rune range, end exclusive
	hardEnd       bool        // a line break or the end of the text ends the line, not wrapping
	top, baseline float32
	height        float32
	width         float32
	stops         []stop // caret positions, one per rune boundary in the line
}

type placedRun struct {
	out   shaping.Output // offsets relative to its paragraph
	base  int            // the paragraph's first rune in the whole text
	x     float32
	width float32
	style int // index into Layout.styles
}

type stop struct {
	index int
	x     float32
}

func toF(v fixed.Int26_6) float32 { return float32(v) / 64 }

func aspectOf(s Style) gfont.Aspect {
	a := gfont.Aspect{Weight: gfont.Weight(s.Weight)}
	if s.Weight == 0 {
		a.Weight = gfont.WeightNormal
	}
	if s.Italic {
		a.Style = gfont.StyleItalic
	} else {
		a.Style = gfont.StyleNormal
	}
	return a
}

// Layout shapes the spans and breaks them into lines.
func (f *Fonts) Layout(spans []Span, opt Options) *Layout {
	var l *Layout
	if opt.Elide && opt.MaxWidth > 0 {
		l = f.elided(spans, opt)
	} else {
		l = f.layout(spans, opt)
	}
	l.align(opt)
	return l
}

// align moves each line across the layout's width, carrying its caret stops
// with it so hit tests agree with what is drawn.
func (l *Layout) align(opt Options) {
	if opt.Align == AlignStart {
		return
	}
	box := opt.MaxWidth
	if box <= 0 {
		box = l.width
	}
	for i := range l.lines {
		ln := &l.lines[i]
		off := box - ln.width
		if opt.Align == AlignMiddle {
			off /= 2
		}
		if off <= 0 {
			continue
		}
		for j := range ln.runs {
			ln.runs[j].x += off
		}
		for j := range ln.stops {
			ln.stops[j].x += off
		}
	}
}

func (f *Fonts) layout(spans []Span, opt Options) *Layout {
	f.mu.Lock()
	defer f.mu.Unlock()
	l := &Layout{fonts: f}
	for _, sp := range spans {
		l.starts = append(l.starts, len(l.runes))
		l.styles = append(l.styles, sp.Style)
		l.runes = append(l.runes, []rune(sp.Text)...)
	}
	if len(spans) == 0 {
		l.starts, l.styles = []int{0}, []Style{{Size: 14}}
	}
	lineMult := opt.LineHeight
	if lineMult <= 0 {
		lineMult = 1
	}

	// Each paragraph, the text between line breaks, is shaped and wrapped on
	// its own, as go-text's wrapper expects; offsets are shifted back to
	// positions in the whole text afterwards.
	y := float32(0)
	ps := 0
	for {
		pe := ps
		for pe < len(l.runes) && l.runes[pe] != '\n' {
			pe++
		}
		y = l.layoutParagraph(ps, pe, opt, lineMult, y)
		if pe >= len(l.runes) {
			break
		}
		ps = pe + 1
	}
	l.height = y
	return l
}

// lineBox is a line's height and the distance from its top to its baseline.
// The natural height is rounded up to a whole pixel before the multiplier,
// and the baseline to the nearest, as Qt's native text rendering does, so
// lines stack at whole pixels and a block of text is as tall here as in the
// Qt version.
func lineBox(ascent, descent, lineMult float32) (height, baseline float32) {
	natural := ascent + descent
	height = float32(math.Ceil(math.Ceil(float64(natural)) * float64(lineMult)))
	baseline = float32(math.Round(float64((height-natural)/2 + ascent)))
	return height, baseline
}

// layoutParagraph shapes and wraps runes [ps, pe), adds its lines from y
// down, and returns the y below them.
func (l *Layout) layoutParagraph(ps, pe int, opt Options, lineMult, y float32) float32 {
	f := l.fonts
	para := l.runes[ps:pe]
	// Shape each span's part of the paragraph with the span's font query; the
	// segmenter splits it further by direction, script and whichever font
	// covers each character.
	var outs []shaping.Output
	var seg shaping.Segmenter
	var shaper shaping.HarfbuzzShaper
	for i, st := range l.styles {
		start, end := l.starts[i], len(l.runes)
		if i+1 < len(l.starts) {
			end = l.starts[i+1]
		}
		start, end = max(start, ps), min(end, pe)
		if end <= start {
			continue
		}
		// A character followed by an emoji presentation mark (U+FE0F, or the
		// keycap U+20E3) belongs to the emoji font with its marks, even when the
		// span's font has the character itself: "1️⃣" and "❤️" split across two
		// fonts would not join.
		for a := start; a < end; {
			emoji := emojiContext(l.runes, a)
			b := a + 1
			for b < end && emojiContext(l.runes, b) == emoji {
				b++
			}
			families := queryFamilies(st.Family)
			if emoji {
				families = emojiFamilies
			}
			f.fm.SetQuery(fontscan.Query{Families: families, Aspect: aspectOf(st)})
			in := shaping.Input{
				Text:      para,
				RunStart:  a - ps,
				RunEnd:    b - ps,
				Direction: di.DirectionLTR,
				Size:      fixed.Int26_6(st.Size * 64),
				Language:  language.NewLanguage("en"),
			}
			if st.Tabular {
				in.FontFeatures = []shaping.FontFeature{{Tag: tnum, Value: 1}}
			}
			for _, part := range seg.Split(in, scriptFontmap{f.fm}) {
				out := shaper.Shape(part)
				markInk(&out, para)
				outs = append(outs, out)
			}
			a = b
		}
	}
	if len(outs) == 0 {
		// An empty paragraph still takes a line, at the metrics of the style in
		// force there, so a caret has somewhere to be.
		st := l.styles[l.styleAt(ps)]
		ascent, descent := l.emptyMetrics(st)
		h, baseline := lineBox(ascent, descent, lineMult)
		l.lines = append(l.lines, line{start: ps, end: ps, hardEnd: true, top: y, height: h,
			baseline: y + baseline, stops: []stop{{ps, 0}}})
		return y + h
	}
	maxWidth := math.MaxInt32 >> 7
	if opt.MaxWidth > 0 {
		maxWidth = int(math.Floor(float64(opt.MaxWidth)))
	}
	var wrapper shaping.LineWrapper
	wrapped, _ := wrapper.WrapParagraph(shaping.WrapConfig{
		Direction:                     di.DirectionLTR,
		DisableTrailingWhitespaceTrim: opt.KeepTrailingSpace,
	}, maxWidth, para, shaping.NewSliceIterator(outs))
	for n, wl := range wrapped {
		ln := line{start: math.MaxInt, end: 0, hardEnd: n == len(wrapped)-1}
		var ascent, descent float32
		for _, o := range wl {
			ascent = max(ascent, toF(o.LineBounds.Ascent))
			descent = max(descent, -toF(o.LineBounds.Descent))
			ln.start = min(ln.start, ps+o.Runes.Offset)
			ln.end = max(ln.end, ps+o.Runes.Offset+o.Runes.Count)
		}
		runs := append(shaping.Line(nil), wl...)
		sort.SliceStable(runs, func(a, b int) bool { return runs[a].VisualIndex < runs[b].VisualIndex })
		x := float32(0)
		for _, o := range runs {
			w := toF(o.Advance)
			ln.runs = append(ln.runs, placedRun{out: o, base: ps, x: x, width: w, style: l.styleAt(ps + o.Runes.Offset)})
			x += w
		}
		h, baseline := lineBox(ascent, descent, lineMult)
		ln.height = h
		ln.top = y
		ln.baseline = y + baseline
		ln.width = x
		ln.stops = l.stopsFor(&ln)
		l.lines = append(l.lines, ln)
		l.width = max(l.width, x)
		y += ln.height
	}
	return y
}

// markInk gives every glyph that is not whitespace a nominal ink width.
// Shaping faces are parsed without their outlines (see strippedFace), so
// every glyph reports zero ink width, and go-text's line wrapper treats a
// glyph with zero ink width at the end of a line as a space: it drops its
// advance and lets it hang past the wrap width. The text package reads no
// ink sizes itself, so the width only has to be non-zero for real glyphs.
func markInk(out *shaping.Output, text []rune) {
	for i := range out.Glyphs {
		g := &out.Glyphs[i]
		if g.Width != 0 || g.ClusterIndex >= len(text) || unicode.IsSpace(text[g.ClusterIndex]) {
			continue
		}
		g.Width = g.XAdvance
	}
}

var tnum = ot.MustNewTag("tnum")

// elided lays the spans out on one line no wider than opt.MaxWidth: as they
// are when they fit, otherwise the longest start that fits with "…" after it,
// the ellipsis in the style of the last character kept.
func (f *Fonts) elided(spans []Span, opt Options) *Layout {
	one := opt
	one.Elide, one.MaxWidth = false, 0
	var flat []Span
	for _, sp := range spans {
		// One line: line breaks become spaces.
		flat = append(flat, Span{Text: strings.ReplaceAll(sp.Text, "\n", " "), Style: sp.Style})
	}
	full := f.Layout(flat, one)
	if full.width <= opt.MaxWidth || full.Len() == 0 {
		return full
	}
	// The longest prefix whose width plus the ellipsis's fits.
	lastStyle := func(n int) Style { return full.styles[full.styleAt(max(0, n-1))] }
	ellipsisWidth := func(st Style) float32 {
		w, _ := f.Layout([]Span{{Text: "…", Style: st}}, one).Size()
		return w
	}
	keep := 0
	for _, s := range full.lines[0].stops {
		if s.index > keep && s.x+ellipsisWidth(lastStyle(s.index)) <= opt.MaxWidth {
			keep = s.index
		}
	}
	// Do not end on a space before the ellipsis.
	for keep > 0 && unicode.IsSpace(full.runes[keep-1]) {
		keep--
	}
	var out []Span
	for i, sp := range flat {
		start := full.starts[i]
		rs := []rune(sp.Text)
		if start >= keep {
			break
		}
		if start+len(rs) > keep {
			rs = rs[:keep-start]
		}
		out = append(out, Span{Text: string(rs), Style: sp.Style})
	}
	out = append(out, Span{Text: "…", Style: lastStyle(keep)})
	return f.Layout(out, one)
}

// emojiFamilies are the colour emoji fonts of the three platforms.
var emojiFamilies = []string{"Noto Color Emoji", "Apple Color Emoji", "Segoe UI Emoji", "emoji"}

// emojiContext reports whether rune i is an emoji presentation mark or is
// followed by one.
func emojiContext(rs []rune, i int) bool {
	isMark := func(r rune) bool { return r == 0xFE0F || r == 0x20E3 }
	if isMark(rs[i]) {
		return true
	}
	return i+1 < len(rs) && isMark(rs[i+1])
}

// emptyMetrics is the ascent and descent a style's font gives an empty line.
func (l *Layout) emptyMetrics(st Style) (ascent, descent float32) {
	f := l.fonts
	f.fm.SetQuery(fontscan.Query{Families: familiesFor(st.Family), Aspect: aspectOf(st)})
	ascent, descent = st.Size*0.8, st.Size*0.2
	if face := f.fm.ResolveFace(' '); face != nil {
		if ext, ok := face.FontHExtents(); ok {
			scale := st.Size / float32(face.Upem())
			ascent, descent = ext.Ascender*scale, -ext.Descender*scale
		}
	}
	return ascent, descent
}

func (l *Layout) styleAt(offset int) int {
	i := sort.Search(len(l.starts), func(i int) bool { return l.starts[i] > offset })
	return max(0, i-1)
}

// stopsFor lists the caret position of every rune boundary on a line. A
// cluster covering several runes (a ligature) has its width shared evenly; a
// right-to-left run places its boundaries from the right.
func (l *Layout) stopsFor(ln *line) []stop {
	pos := map[int]float32{}
	for _, r := range ln.runs {
		rtl := r.out.Direction.Progression() == di.TowardTopLeft
		x := r.x
		gs := r.out.Glyphs
		for i := 0; i < len(gs); {
			g := gs[i]
			n := max(1, g.GlyphCount)
			var w float32
			for k := i; k < i+n && k < len(gs); k++ {
				w += toF(gs[k].XAdvance)
			}
			runes := max(1, g.RuneCount)
			for k := 0; k <= runes; k++ {
				frac := w * float32(k) / float32(runes)
				idx := r.base + g.ClusterIndex + k
				if rtl {
					pos[idx] = x + w - frac
				} else {
					pos[idx] = x + frac
				}
			}
			x += w
			i += n
		}
	}
	out := make([]stop, 0, len(pos))
	for i, x := range pos {
		if i >= ln.start && i <= ln.end {
			out = append(out, stop{i, x})
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].index < out[b].index })
	if len(out) == 0 {
		out = []stop{{ln.start, 0}}
	}
	return out
}

// Size is the layout's width (its widest line) and height.
func (l *Layout) Size() (width, height float32) { return l.width, l.height }

// LineCount is the number of lines.
func (l *Layout) LineCount() int { return len(l.lines) }

// Len is the number of runes laid out.
func (l *Layout) Len() int { return len(l.runes) }

// Baseline is the first line's baseline, from the top.
func (l *Layout) Baseline() float32 { return l.lines[0].baseline }

// LineBounds returns line i's rune range and its top and height.
func (l *Layout) LineBounds(i int) (start, end int, top, height float32) {
	ln := l.lines[i]
	return ln.start, ln.end, ln.top, ln.height
}

func (l *Layout) lineFor(index int) int {
	for i, ln := range l.lines {
		if index < ln.end || (index == ln.end && ln.hardEnd) || i == len(l.lines)-1 {
			return i
		}
	}
	return 0
}

// CaretAt is where a caret before rune index goes: its x, and the top and
// height of its line.
func (l *Layout) CaretAt(index int) (x, top, height float32) {
	index = min(max(index, 0), len(l.runes))
	ln := l.lines[l.lineFor(index)]
	x = ln.width
	for _, s := range ln.stops {
		if s.index == index {
			x = s.x
			break
		}
		if s.index > index {
			x = s.x
			break
		}
	}
	return x, ln.top, ln.height
}

// IndexAt is the rune boundary nearest to a point, for placing a caret or
// starting a selection where the pointer is.
func (l *Layout) IndexAt(x, y float32) int {
	li := len(l.lines) - 1
	for i, ln := range l.lines {
		if y < ln.top+ln.height {
			li = i
			break
		}
	}
	ln := l.lines[li]
	best, bestD := ln.start, float32(math.MaxFloat32)
	for _, s := range ln.stops {
		// A wrapped line's end belongs to the next line.
		if s.index == ln.end && !ln.hardEnd {
			continue
		}
		if d := abs(s.x - x); d < bestD {
			best, bestD = s.index, d
		}
	}
	return best
}

func abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}
