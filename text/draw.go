package text

import (
	gfont "github.com/go-text/typesetting/font"
	cfont "github.com/richardwilkes/canvas/font"
	"github.com/richardwilkes/canvas/textblob"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

// Unison converts a colour for drawing.
func (c Color) Unison() unison.Color {
	return unison.ARGB(float32(c.A)/255, int(c.R), int(c.G), int(c.B))
}

type fontKey struct {
	loc  gfont.FontID
	size float32
}

// font returns the canvas font for a face at a size, from a cache.
func (f *Fonts) font(face *gfont.Face, size float32) *cfont.Font {
	if f.sized == nil {
		f.sized = map[fontKey]*cfont.Font{}
	}
	key := fontKey{f.fm.FontLocation(face.Font), size}
	if cf, ok := f.sized[key]; ok {
		return cf
	}
	tf, err := f.typeface(face)
	if err != nil {
		return nil
	}
	cf := cfont.NewFont(tf, size, 1, 0)
	cf.SetSubpixel(true)
	cf.SetForceAutoHinting(true)
	cf.SetHinting(cfont.HintingNone)
	f.sized[key] = cf
	return cf
}

// Draw draws the layout with its top-left corner at (x, y).
func (l *Layout) Draw(gc *unison.Canvas, x, y float32) {
	l.fonts.mu.Lock()
	defer l.fonts.mu.Unlock()
	for _, ln := range l.lines {
		for _, r := range ln.runs {
			if bg := l.styles[r.style].Background; bg.A > 0 {
				rect := geom.NewRect(x+r.x, y+ln.top, r.width, ln.height)
				gc.DrawRect(rect, bg.Unison().Paint(gc, rect, paintstyle.Fill))
			}
		}
		for _, r := range ln.runs {
			st := l.styles[r.style]
			origin := geom.NewPoint(x+r.x, y+ln.baseline-st.Rise)
			paint := st.Color.Unison().Paint(gc, geom.Rect{}, paintstyle.Fill)
			if len(r.out.Glyphs) > 0 && !st.Box.on() {
				if cf := l.fonts.font(r.out.Face, toF(r.out.Size)); cf != nil {
					b := textblob.NewBuilder()
					buf := b.AllocRunPos(cf, len(r.out.Glyphs), nil)
					var cx float32
					for i, g := range r.out.Glyphs {
						buf.Glyphs[i] = uint16(g.GlyphID)
						buf.Pos[2*i] = cx + toF(g.XOffset)
						buf.Pos[2*i+1] = -toF(g.YOffset)
						cx += toF(g.XAdvance)
					}
					if blob := b.Make(); blob != nil {
						gc.DrawTextBlob(blob, origin, paint)
					}
				}
			}
			if (st.Underline || st.Strike) && !st.Box.on() {
				thick := max(1, st.Size/14)
				line := st.Color.Unison().Paint(gc, geom.Rect{}, paintstyle.Stroke)
				line.SetStrokeWidth(thick)
				if st.Underline {
					uy := origin.Y + max(thick, st.Size*0.1)
					gc.DrawLine(geom.NewPoint(origin.X, uy), geom.NewPoint(origin.X+r.width, uy), line)
				}
				if st.Strike {
					sy := origin.Y - st.Size*0.28
					gc.DrawLine(geom.NewPoint(origin.X, sy), geom.NewPoint(origin.X+r.width, sy), line)
				}
			}
		}
	}
}
