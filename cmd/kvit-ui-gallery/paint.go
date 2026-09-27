package main

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

// painter draws one page. With a nil canvas it draws nothing and only
// measures, which is how a page finds its height before it is drawn.
type painter struct {
	gc *unison.Canvas
	ui *kvitui.UI
}

func (p painter) fill(r geom.Rect, c palette.Color) {
	if p.gc != nil {
		p.gc.DrawRect(r, kvitui.Color(c).Paint(p.gc, r, paintstyle.Fill))
	}
}

func (p painter) round(r geom.Rect, radius float32, c palette.Color) {
	if p.gc != nil {
		p.gc.DrawRoundedRect(r, geom.NewSize(radius, radius), kvitui.Color(c).Paint(p.gc, r, paintstyle.Fill))
	}
}

func (p painter) outline(r geom.Rect, radius, width float32, c palette.Color) {
	if p.gc != nil {
		paint := kvitui.Color(c).Paint(p.gc, r, paintstyle.Stroke)
		paint.SetStrokeWidth(width)
		inset := r.Inset(geom.NewUniformInsets(width / 2))
		p.gc.DrawRoundedRect(inset, geom.NewSize(radius, radius), paint)
	}
}

// text lays out spans within maxWidth (0 for no wrapping), draws them with
// their top-left at (x, y), and returns their size.
func (p painter) text(spans []text.Span, x, y, maxWidth float32) (w, h float32) {
	l := p.ui.Fonts.Layout(spans, text.Options{MaxWidth: maxWidth})
	if p.gc != nil {
		l.Draw(p.gc, x, y)
	}
	return l.Size()
}

func span(s string, st text.Style) []text.Span { return []text.Span{{Text: s, Style: st}} }
