package kvitui

import (
	"math"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// MoveTo puts the dialog's top-left corner x across and y down from the
// top-left corner of the window's content: where it opens, or where it goes
// at once when it is open. The window keeps the whole dialog inside it,
// however far out x and y are. Px of what OnMove reported opens the dialog
// where the reader left it last time.
func (d *Dialog) MoveTo(x, y Measure) {
	d.x, d.y = x, y
	if p := d.Parent(); d.hide != nil && p != nil {
		p.MarkForLayoutAndRedraw()
	}
}

// Position is where the dialog's top-left corner is put, in design pixels
// from the top-left corner of the window's content, and false while it is
// put in the middle, as it is until MoveTo or the reader moves it.
func (d *Dialog) Position() (x, y int, placed bool) {
	if d.x == nil || d.y == nil {
		return 0, 0, false
	}
	return d.design(float32(d.x.Of(d.ui))), d.design(float32(d.y.Of(d.ui))), true
}

// design is a length in pixels in design pixels.
func (d *Dialog) design(pixels float32) int {
	return int(math.Round(float64(pixels) / d.ui.Interface.Scale()))
}

// place is where the open dialog goes in the window's content: the width a
// dialog is and the height it needs at that width, in the middle, or with its
// top-left corner where MoveTo or the reader put it, kept inside.
func (d *Dialog) place(bounds geom.Rect, _ geom.Size) geom.Rect {
	w := float32(d.ui.Interface.Px(420))
	_, pref, _ := d.Sizes(geom.NewSize(w, 0))
	size := geom.NewSize(w, pref.Height)
	if d.x == nil || d.y == nil {
		return PlaceCentred(nil)(bounds, size)
	}
	at := geom.NewPoint(bounds.X+float32(d.x.Of(d.ui)), bounds.Y+float32(d.y.Of(d.ui)))
	return keepInside(geom.Rect{Point: at, Size: size}, bounds)
}

// keepInside moves a box as little as it takes to be inside bounds, or to
// their top-left corner when it is larger, which keeps a dialog's title and
// close control in reach.
func keepInside(box, bounds geom.Rect) geom.Rect {
	box.X = max(bounds.X, min(box.X, bounds.Right()-box.Width))
	box.Y = max(bounds.Y, min(box.Y, bounds.Bottom()-box.Height))
	return box
}

// titleRow is the band a drag moves the dialog by: the title's row with the
// margin above it and the one below, down to where the detail starts. The
// close control in it takes its own presses.
func (d *Dialog) titleRow() geom.Rect {
	m := d.ui.Interface
	r := d.ContentRect(true)
	return geom.NewRect(r.X, r.Y, r.Width, float32(2*m.ViewMargin()+m.ControlHeight()))
}

// inContent is a point in the dialog in its window's content coordinates,
// which do not move as the dialog does.
func (d *Dialog) inContent(where geom.Point) geom.Point {
	pt := d.PointToRoot(where)
	if w := d.Window(); w != nil {
		return w.Content().PointFromRoot(pt)
	}
	return pt
}

func (d *Dialog) mouseDown(where geom.Point, button, _ int, _ mod.Modifiers) bool {
	if d.hide == nil || button != unison.ButtonLeft || !where.In(d.titleRow()) {
		return false
	}
	d.dragging, d.moved = true, false
	d.grab, d.grabAt = d.inContent(where), d.FrameRect().Point
	return true
}

// mouseDrag keeps the point of the dialog that was pressed under the pointer,
// as far as the window allows: past an edge the dialog stops, and it follows
// again once the pointer is back to where that point would be.
func (d *Dialog) mouseDrag(where geom.Point, _ int, _ mod.Modifiers) bool {
	p := d.Parent()
	if !d.dragging || p == nil {
		return false
	}
	bounds := p.ContentRect(false)
	to := d.grabAt.Add(d.inContent(where).Sub(d.grab))
	box := keepInside(geom.Rect{Point: to, Size: d.FrameRect().Size}, bounds)
	x, y := d.design(box.X-bounds.X), d.design(box.Y-bounds.Y)
	if now := d.FrameRect().Point; d.design(now.X-bounds.X) == x && d.design(now.Y-bounds.Y) == y {
		return true
	}
	d.moved = true
	d.MoveTo(Px(x), Px(y))
	return true
}

func (d *Dialog) mouseUp(geom.Point, int, mod.Modifiers) bool {
	if !d.dragging {
		return false
	}
	d.dragging = false
	if d.moved && d.OnMove != nil {
		x, y, _ := d.Position()
		d.OnMove(x, y)
	}
	return true
}
