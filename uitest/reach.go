package uitest

import (
	"fmt"
	"math"
	"strings"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// Reach lists what a reader cannot get to on a window's screen as it is laid
// out now (LAY-3): a control, a scroll view or a popup that the panels around
// it, or the window's own edge, cut off where nothing scrolls it into view.
// Inside a scroll view, a Region's or the Sidebar's, a control may lie outside
// what is shown, since scrolling brings it in; the scroll view itself must not
// be cut off. A popup that hangs from something in a scroll view that is out
// of sight, as a menu pinned beside a sample further down a page does, moves
// with it and comes into view as it is scrolled to, so it is read only while
// what it hangs from is in sight. A control cut off is reported once, not
// with every panel inside it. Call it on the interface's thread, as
// Session.Do does.
func Reach(w *kvitui.Window) []Problem {
	c := &reachCheck{form: w.UI().Form().String(), popups: map[*unison.Panel]bool{}, viewports: map[*unison.Panel]bool{},
		moving: map[*unison.Panel]bool{}}
	root := w.Content()
	window := root.RectToRoot(root.ContentRect(false))
	for _, p := range w.Popups() {
		if a := p.Anchor; a != nil && scrolled(a.AsPanel()) && !inside(a.AsPanel().RectToRoot(a.AsPanel().ContentRect(true)), window) {
			c.moving[p.Panel.AsPanel()] = true
			continue
		}
		c.popups[p.Panel.AsPanel()] = true
	}
	c.walk(root, edges{box: window, by: root, window: true})
	return c.out
}

// scrolled reports whether a panel is inside a scroll view: a Region's or
// another unison scroll panel's content, or a Sidebar's places.
func scrolled(p *unison.Panel) bool {
	for q := p; q != nil && q.Parent() != nil; q = q.Parent() {
		switch v := q.Parent().Self.(type) {
		case *unison.ScrollPanel:
			if q == v.ContentView() {
				return true
			}
		case *kvitui.Sidebar:
			if rows := v.Rows(); len(rows) > 0 && rows[0].Parent() == q {
				return true
			}
		}
	}
	return false
}

// edges are the bounds within which a panel can be seen: the box its
// ancestors leave, and the panel that last narrowed it.
type edges struct {
	box    geom.Rect
	by     *unison.Panel
	window bool // the window's own edge is what narrowed it
	open   bool // inside a scroll view: no bound but the panels below it
}

type reachCheck struct {
	form      string
	out       []Problem
	popups    map[*unison.Panel]bool
	viewports map[*unison.Panel]bool
	moving    map[*unison.Panel]bool // popups that come into view with what they hang from
}

// slack is how far, in logical pixels, a panel may pass its bounds before it
// counts as cut off, for the rounding of a layout to whole pixels.
const slack = 1

func (c *reachCheck) walk(p *unison.Panel, in edges) {
	if p.Hidden || c.moving[p] {
		return
	}
	box := p.RectToRoot(p.ContentRect(true))
	if box.Width <= 0 || box.Height <= 0 {
		// Folded away on purpose, as a closed disclosure's body is.
		return
	}
	switch v := p.Self.(type) {
	case *unison.ScrollPanel:
		c.viewports[v.ContentView()] = true
	case *kvitui.Sidebar:
		if rows := v.Rows(); len(rows) > 0 && rows[0].Parent() != nil {
			c.viewports[rows[0].Parent()] = true
		}
	}
	if c.reached(p) && !in.open && !inside(box, in.box) {
		c.out = append(c.out, Problem{Rule: "LAY-3", Where: reachName(p), What: cutWhere(box, in), Form: c.form})
		return
	}
	next := in
	switch {
	case c.viewports[p]:
		// What the scroll view holds may be anywhere: scrolling brings it in.
		next = edges{open: true}
	case in.open:
		next = edges{box: box, by: p}
	default:
		if cut := in.box.Intersect(box); cut != in.box {
			next = edges{box: cut, by: p}
		}
	}
	for _, child := range p.Children() {
		c.walk(child, next)
	}
}

// reached reports whether a panel is something a reader must be able to get
// to: a control that takes the keyboard or a press, a scroll view, or a popup.
func (c *reachCheck) reached(p *unison.Panel) bool {
	if c.popups[p] || p.Focusable() {
		return true
	}
	switch p.Self.(type) {
	case *unison.ScrollPanel, *kvitui.Region, *kvitui.Sidebar, *kvitui.Table, *kvitui.Tree, *kvitui.Dialog:
		return true
	}
	return false
}

// inside reports whether a box lies within bounds, give or take the slack.
func inside(box, bounds geom.Rect) bool {
	return box.X >= bounds.X-slack && box.Y >= bounds.Y-slack &&
		box.Right() <= bounds.Right()+slack && box.Bottom() <= bounds.Bottom()+slack
}

// cutWhere says how a box passes its bounds and what cuts it there.
func cutWhere(box geom.Rect, in edges) string {
	by := "the window's edge"
	if !in.window {
		by = "the " + panelKind(in.by) + " around it"
	}
	if box.Intersect(in.box).Empty() {
		return "out of sight, outside " + by + ", and nothing scrolls it into view"
	}
	var sides []string
	px := func(v float32) string { return fmt.Sprintf("%d px", int(math.Round(float64(v)))) }
	if d := box.Bottom() - in.box.Bottom(); d > slack {
		sides = append(sides, "its bottom "+px(d))
	}
	if d := in.box.Y - box.Y; d > slack {
		sides = append(sides, "its top "+px(d))
	}
	if d := box.Right() - in.box.Right(); d > slack {
		sides = append(sides, "its right "+px(d))
	}
	if d := in.box.X - box.X; d > slack {
		sides = append(sides, "its left "+px(d))
	}
	return strings.Join(sides, " and ") + " cut off by " + by + ", and nothing scrolls it into view"
}

// reachName names what is cut off in a problem: a control by its kind and
// words, as a Tab stop is named, and anything else by its kind.
func reachName(p *unison.Panel) string {
	switch v := p.Self.(type) {
	case *kvitui.Dialog:
		return "dialog " + quote(v.Title)
	case *kvitui.Region, *kvitui.Sidebar, *kvitui.Table, *kvitui.Tree, *unison.ScrollPanel:
		return panelKind(p)
	}
	return stopName(p)
}

// panelKind is the kind of a panel, the Kvit component it is or is part of:
// "Card", or "panel in Card" for a plain panel inside one.
func panelKind(p *unison.Panel) string {
	name := func(q *unison.Panel) (string, bool) {
		t := fmt.Sprintf("%T", q.Self)
		if strings.HasPrefix(t, "*kvitui.") {
			return strings.TrimPrefix(t, "*kvitui."), true
		}
		return t, false
	}
	if n, ok := name(p); ok {
		return n
	}
	for q := p.Parent(); q != nil; q = q.Parent() {
		if n, ok := name(q); ok {
			return "panel in " + n
		}
	}
	return "panel"
}

// CheckReach shrinks the window to the smallest size it allows, the
// interface's width and height floors (LAY-3), reads it there with Reach,
// and fails the test for every problem known does not list, as Report does.
// The window then has its size back. Call it on each screen of an app, after
// opening whatever the screen shows: a dialog, a menu, a pane.
func (s *Session) CheckReach(known ...string) {
	s.t.Helper()
	var was geom.Rect
	s.Do(func() {
		was = s.Window.ContentRect()
		m := s.UI.Interface
		s.Window.SetContentRect(geom.NewRect(was.X, was.Y, float32(m.WidthFloor()), float32(m.HeightFloor())))
	})
	s.Sync()
	var problems []Problem
	s.Do(func() {
		problems = Reach(s.Window)
		s.Window.SetContentRect(was)
	})
	s.Sync()
	Report(s.t, problems, known...)
}
