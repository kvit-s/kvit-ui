package kvitui

import (
	"slices"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/pathop"
)

// spiller is a component that draws past its own box, as a Qt item can: a
// trend's axis values centred on its top and bottom gridlines, a gauge's pace
// tick standing out of the bar. unison clips each panel's drawing to its box,
// so the part drawn outside it is drawn by the nearest panel above it that
// hosts spills, after everything in that panel. It is clipped where a Qt item
// would clip it: at the view of a scroll panel it is in, and at a spillClip,
// as a disclosure's body is while it grows open. drawSpill draws in the
// component's own coordinates.
type spiller interface {
	drawSpill(gc *unison.Canvas)
}

// hostSpills makes a panel of the window draw the spills of the components
// inside it, over its children and under anything the panel itself draws
// over them, such as a Kvit panel's rules. The window does this for its body
// and its sidebar, so a spill is drawn over the page it is on, and under the
// rail when the rail opens over the page and under a popup. The window keeps
// its hosts, so they go with it when it closes.
func (w *Window) hostSpills(host *unison.Panel) {
	if slices.Contains(w.spillHosts, host) {
		return
	}
	w.spillHosts = append(w.spillHosts, host)
	previous := host.DrawOverCallback
	host.DrawOverCallback = func(gc *unison.Canvas, rect geom.Rect) {
		drawSpills(gc, w.spillHosts, host, geom.Point{}, host.ContentRect(true))
		if previous != nil {
			previous(gc, rect)
		}
	}
}

// spillClip is a panel that clips what the components inside it draw past
// their boxes to its own box, as a disclosure's body does while it grows
// open. It says so by what it is, rather than by being listed somewhere that
// would keep it after it is gone.
type spillClip struct {
	unison.Panel
}

// newSpillClip returns an empty panel that clips the spills inside it.
func newSpillClip() *unison.Panel {
	p := &spillClip{}
	p.Self = p
	return p.AsPanel()
}

// clipsSpills reports whether a panel clips the spills inside it: a view of a
// scroll panel, or a spillClip.
func clipsSpills(p *unison.Panel) bool {
	if _, ok := p.Self.(*spillClip); ok {
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
// showing. A panel inside another of the window's hosts is left to that
// host.
func drawSpills(gc *unison.Canvas, hosts []*unison.Panel, p *unison.Panel, origin geom.Point, clip geom.Rect) {
	for _, c := range p.Children() {
		if c.Hidden || slices.Contains(hosts, c) {
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
		if clipsSpills(c) {
			inner = clip.Intersect(geom.Rect{Point: at, Size: frame.Size})
		}
		if !inner.Empty() {
			drawSpills(gc, hosts, c, at, inner)
		}
	}
}

// spillHosted reports whether a panel above p draws p's spill. When none
// does, as outside a Kvit window, the component draws it itself, clipped to
// its box.
func (u *UI) spillHosted(p *unison.Panel) bool {
	w := u.windowOf(p)
	if w == nil {
		return false
	}
	for a := p.Parent(); a != nil; a = a.Parent() {
		if slices.Contains(w.spillHosts, a) {
			return true
		}
	}
	return false
}
