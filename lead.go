package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/toolbox/v2/xreflect"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/filtermode"
	"github.com/richardwilkes/unison/enums/mipmapmode"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// disabledPictureOpacity is how much of a picture shows on a line that cannot
// act: a picture has its own colours, so it cannot take the disabled text
// colour as a symbol would.
const disabledPictureOpacity = 0.4

// swapLead puts p where old was: a panel the control draws before its words,
// which the pointer treats as part of the control. It is the control's first
// child, so whatever else the control holds, such as a badge, draws over it.
func (c *control) swapLead(old, p unison.Paneler) {
	if old != nil && old.AsPanel().Parent() == c.AsPanel() {
		old.AsPanel().RemoveFromParent()
	}
	if p != nil {
		c.AddChildAtIndex(p, 0)
		c.passPointer(p.AsPanel())
	}
	c.MarkForLayoutAndRedraw()
}

// passPointer makes p and every panel inside it, such as the one Sized holds,
// tell the control where the pointer is. unison tells only the panel under the
// pointer. A press needs nothing here: unison offers it to the panel's parents
// until one takes it, and the control does.
func (c *control) passPointer(p *unison.Panel) {
	here := func(where geom.Point) geom.Point { return c.PointFromRoot(p.PointToRoot(where)) }
	p.MouseEnterCallback = func(where geom.Point, mods mod.Modifiers) bool {
		return c.MouseEnterCallback(here(where), mods)
	}
	p.MouseMoveCallback = func(where geom.Point, mods mod.Modifiers) bool {
		return c.MouseMoveCallback(here(where), mods)
	}
	p.MouseExitCallback = func() bool { return c.MouseExitCallback() }
	p.UpdateCursorCallback = c.UpdateCursorCallback
	for _, child := range p.Children() {
		c.passPointer(child)
	}
}

// silenceLead takes a lead and everything inside it out of what a screen
// reader is told, so the control's name is all that is heard of it.
func silenceLead(p *unison.Panel) {
	p.Accessibility.Role = role.None
	for _, child := range p.Children() {
		silenceLead(child)
	}
}

// leadSize is the size a lead asks for, made smaller with its shape kept where
// it would not fit room.
func leadSize(lead *unison.Panel, room geom.Size) geom.Size {
	_, pref, _ := lead.Sizes(geom.Size{})
	if pref.Width <= room.Width && pref.Height <= room.Height {
		return pref
	}
	return fitted(geom.Rect{Size: room}, pref).Size
}

// fitted is the largest rectangle of natural's shape that fits box, centred in
// it, or box itself for a natural size that has no shape.
func fitted(box geom.Rect, natural geom.Size) geom.Rect {
	if natural.Width <= 0 || natural.Height <= 0 {
		return box
	}
	scale := min(box.Width/natural.Width, box.Height/natural.Height)
	w, h := natural.Width*scale, natural.Height*scale
	return geom.NewRect(box.X+(box.Width-w)/2, box.Y+(box.Height-h)/2, w, h)
}

// pictureSampling blends a picture's pixels when it is drawn at another size,
// which unison's default sampling does not, leaving a reduced picture jagged.
var pictureSampling = &unison.SamplingOptions{FilterMode: filtermode.Linear, MipMapMode: mipmapmode.Linear}

// drawPicture draws a picture fitted inside box with its shape kept, on whole
// pixels of the window, and dimmed when dim is set. A picture is drawn with no
// paint, since a paint would replace an SVG's own colours with one.
func drawPicture(gc *unison.Canvas, owner *unison.Panel, picture unison.Drawable, box geom.Rect, dim bool) {
	if dim {
		gc.SaveWithOpacity(disabledPictureOpacity)
		defer gc.Restore()
	}
	picture.DrawInRect(gc, toWholePixels(owner, fitted(box, picture.LogicalSize())), pictureSampling, nil)
}

// hasPicture reports whether a picture is there to draw. A loader's nil image
// in the interface counts as none.
func hasPicture(picture unison.Drawable) bool { return !xreflect.IsNil(picture) }
