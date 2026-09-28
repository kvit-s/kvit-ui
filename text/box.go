package text

import (
	"math"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"
)

// Box keeps room in a line for something the caller draws over the text,
// such as a typeset formula. A span whose style has a Box with a Width above
// zero is laid out as one unbreakable piece that wide, whatever its text (one
// U+FFFC is usual), and its text is not drawn; its background is, so a
// selection or a find match over it still shows. The line holding it keeps
// at least Ascent above its baseline and Descent below it: a line of a fixed
// Pitch grows by what the box needs beyond the pitch, and keeps the pitch
// when the box fits in it.
//
// The caller finds where a box was put with CaretAt, at the box's first rune,
// and LineBaseline.
type Box struct {
	Width, Ascent, Descent float32
}

func (b Box) on() bool { return b.Width > 0 }

func toFixed(v float32) fixed.Int26_6 { return fixed.Int26_6(math.Round(float64(v) * 64)) }

// boxOutput is the run a box span is laid out as: one glyph covering its runes
// [start, start+n) of the paragraph, the box's width across, and the box's
// ascent and descent as the line bounds. It has no font, so it is never drawn.
func boxOutput(b Box, size float32, start, n int) shaping.Output {
	w := toFixed(b.Width)
	g := shaping.Glyph{Width: w, Advance: w, XAdvance: w, ClusterIndex: start, RuneCount: n, GlyphCount: 1}
	return shaping.Output{
		Advance:    w,
		Size:       toFixed(size),
		Glyphs:     []shaping.Glyph{g},
		LineBounds: shaping.Bounds{Ascent: toFixed(b.Ascent), Descent: -toFixed(b.Descent)},
		Direction:  di.DirectionLTR,
		Runes:      shaping.Range{Offset: start, Count: n},
	}
}

// boxedLineBox is lineBox for a line whose text needs ascent and descent and
// whose boxes need boxAscent and boxDescent.
func (l *Layout) boxedLineBox(ascent, descent, boxAscent, boxDescent, lineMult, y float32, opt Options) (height, baseline float32) {
	if opt.Pitch <= 0 {
		return l.lineBox(max(ascent, boxAscent), max(descent, boxDescent), lineMult, y, opt)
	}
	height, baseline = l.lineBox(ascent, descent, lineMult, y, opt)
	if boxAscent <= baseline && boxDescent <= height-baseline {
		return height, baseline
	}
	// The baseline moves down to a whole pixel under the box's ascent, and the
	// line keeps below it the larger of the pitch's own room and the box's
	// descent.
	below := max(height-baseline, boxDescent)
	baseline = max(baseline, float32(math.Ceil(float64(y+boxAscent)))-y)
	return baseline + below, baseline
}

// LineBaseline is line i's baseline, from the layout's top.
func (l *Layout) LineBaseline(i int) float32 { return l.lines[i].baseline }
