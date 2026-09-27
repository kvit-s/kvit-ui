package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/pathop"
)

// spiller is a component that draws past its own box, as a Qt item can: a
// trend's axis values centred on its top and bottom gridlines, a gauge's pace
// tick standing out of the bar. unison clips each panel's drawing to its box,
// so the part drawn outside it is drawn by the nearest panel above it that
// hosts spills, after everything in that panel. It is clipped where a Qt item
// would clip it: at the view of a scroll panel it is in, and at a panel that
// asks for it with clipSpills, as a disclosure's body does while it grows
// open. drawSpill draws in the component's own coordinates.
type spiller interface {
	drawSpill(gc *unison.Canvas)
}

// hostSpills makes a panel draw the spills of the components inside it, over
// its children and under anything the panel itself draws over them, such as
// a Kvit panel's rules. The window does this for its body and its sidebar, so
// a spill is drawn over the page it is on, and under the rail when the rail
// opens over the page and under a popup.
func (u *UI) hostSpills(host *unison.Panel) {
	if u.spillHosts[host] {
		return
	}
	if u.spillHosts == nil {
		u.spillHosts = map[*unison.Panel]bool{}
	}
	u.spillHosts[host] = true
	previous := host.DrawOverCallback
	host.DrawOverCallback = func(gc *unison.Canvas, rect geom.Rect) {
		u.drawSpills(gc, host, geom.Point{}, host.ContentRect(true))
		if previous != nil {
			previous(gc, rect)
		}
	}
}

// clipSpills makes a panel clip what the components inside it draw past their
// boxes to its own box.
func (u *UI) clipSpills(p *unison.Panel) {
	if u.spillClips == nil {
		u.spillClips = map[*unison.Panel]bool{}
	}
	u.spillClips[p] = true
}

// clipsSpills reports whether a panel clips the spills inside it: a view of a
// scroll panel, or one that asked to.
func (u *UI) clipsSpills(p *unison.Panel) bool {
	if u.spillClips[p] {
		return true
	}
	if parent := p.Parent(); parent != nil {
		_, scroll := parent.Self.(*unison.ScrollPanel)
		return scroll
	}
	return false
}

// drawSpills draws the spills inside p, whose origin is at origin in the
// host's coordinates and whose clipping containers, p among them, leave clip
// showing. A panel inside another host is left to that host.
func (u *UI) drawSpills(gc *unison.Canvas, p *unison.Panel, origin geom.Point, clip geom.Rect) {
	for _, c := range p.Children() {
		if c.Hidden || u.spillHosts[c] {
			continue
		}
		frame := c.FrameRect()
		at := origin.Add(frame.Point)
		if s, ok := c.Self.(spiller); ok && !clip.Empty() {
			gc.Save()
			gc.ClipRect(clip, pathop.Intersect, false)
			gc.Translate(at)
			s.drawSpill(gc)
			gc.Restore()
		}
		inner := clip
		if u.clipsSpills(c) {
			inner = clip.Intersect(geom.Rect{Point: at, Size: frame.Size})
		}
		if !inner.Empty() {
			u.drawSpills(gc, c, at, inner)
		}
	}
}

// spillHosted reports whether a panel above p draws p's spill. When none
// does, as outside a Kvit window, the component draws it itself, clipped to
// its box.
func (u *UI) spillHosted(p *unison.Panel) bool {
	for a := p.Parent(); a != nil; a = a.Parent() {
		if u.spillHosts[a] {
			return true
		}
	}
	return false
}
