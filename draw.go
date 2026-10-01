package kvitui

import (
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

// painter draws the shapes components are made of, in design colours.
type painter struct {
	gc *unison.Canvas
	ui *UI
}

func painterFor(gc *unison.Canvas, ui *UI) painter { return painter{gc, ui} }

func (p painter) fill(r geom.Rect, c palette.Color) {
	p.gc.DrawRect(r, Color(c).Paint(p.gc, r, paintstyle.Fill))
}

func (p painter) round(r geom.Rect, radius float32, c palette.Color) {
	p.gc.DrawRoundedRect(r, geom.NewSize(radius, radius), Color(c).Paint(p.gc, r, paintstyle.Fill))
}

// outline strokes a rounded rectangle whose outer edge is r.
func (p painter) outline(r geom.Rect, radius, width float32, c palette.Color) {
	paint := Color(c).Paint(p.gc, r, paintstyle.Stroke)
	paint.SetStrokeWidth(width)
	inner := max(0, radius-width/2)
	p.gc.DrawRoundedRect(r.Inset(geom.NewUniformInsets(width/2)), geom.NewSize(inner, inner), paint)
}

// roundTint fills a rounded rectangle with a colour at a share of its
// strength, over whatever is below it.
func (p painter) roundTint(r geom.Rect, radius float32, c palette.Color, share float32) {
	p.gc.DrawRoundedRect(r, geom.NewSize(radius, radius), Color(c).SetAlphaIntensity(share).Paint(p.gc, r, paintstyle.Fill))
}
