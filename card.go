package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// Card is a bounded block of content sitting on a surface: a dashboard
// widget, a summary, a group of related figures. It has a corner radius
// where a panel has none, because it is an object on a surface rather than a
// part of the window.
//
// Its outline is the ordinary border colour unless the card can be pressed,
// in which case it is a control and takes the strong border colour, the
// hover tint, and, as kvit-cash's copy of the Qt library adds, the chevron a
// pressable row carries, in its top corner. A press on something inside the
// card that acts on a press of its own is that thing's, and not also the
// card's. Unlike the Qt card, a pressable card also takes the keyboard focus
// and opens on Return, Enter or Space.
type Card struct {
	control
	// Interactive makes the card something to press.
	Interactive bool
	// NoOpensMark leaves the chevron off a pressable card whose ways in are
	// the rows inside it, so each row carries its own.
	NoOpensMark bool
	// OpensLabel says in words what pressing the card opens: the tooltip on
	// the chevron, and the card's accessible description.
	OpensLabel string
	// Selected marks a chosen card.
	Selected bool
	// Padding is the space inside the edge; the loose space unless set.
	Padding Measure
	// OnActivate runs when a pressable card is pressed.
	OnActivate func()
}

// NewCard returns a card holding content, stacked top to bottom.
func NewCard(ui *UI, content ...unison.Paneler) *Card {
	c := &Card{Padding: SizeSpaceLoose}
	c.Self = c
	c.initControl(ui, func() {
		if c.Interactive && c.OnActivate != nil {
			c.OnActivate()
		}
	}, unison.KeyReturn, unison.KeyNumPadEnter, unison.KeySpace)
	down, up := c.MouseDownCallback, c.MouseUpCallback
	c.MouseDownCallback = func(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
		return c.Interactive && down(where, button, clicks, mods)
	}
	c.MouseUpCallback = func(where geom.Point, button int, mods mod.Modifiers) bool {
		return c.Interactive && up(where, button, mods)
	}
	c.ringRadius = func() float32 { return float32(ui.Interface.RadiusCard()) }
	for _, p := range content {
		if p.AsPanel().LayoutData() == nil {
			p.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		}
		c.AddChild(p)
	}
	c.SetBorder(cardInsets{c})
	c.SetLayout(syncing{Layout: &unison.FlexLayout{Columns: 1}, sync: func() { c.SetFocusable(c.Interactive) }})
	c.DrawCallback = c.draw
	c.UpdateTooltipCallback = func(where geom.Point, _ geom.Rect) geom.Rect {
		c.Tooltip = nil
		strip := c.opensStrip()
		if c.opens() && c.OpensLabel != "" && where.In(strip) {
			c.Tooltip = newTooltip(ui, c.OpensLabel, "")
		}
		return c.RectToRoot(strip)
	}
	return c
}

func (c *Card) opens() bool { return c.Interactive && !c.NoOpensMark }

// opensStrip is where the chevron sits: the top corner, on the card's first
// line rather than halfway down, since a card is as tall as its content and
// a mark in the middle belongs to whatever it happens to be beside.
func (c *Card) opensStrip() geom.Rect {
	m := c.ui.Interface
	b := c.ContentRect(true)
	pad := orZero(c.ui, c.Padding)
	w := float32(m.IconSizeSmall() + 2*m.SpaceTight())
	return geom.NewRect(b.Right()-pad-w, b.Y+pad, w, float32(m.RowHeightCompact()))
}

// cardInsets are the padding, and the chevron's room at the right of a card
// that opens something.
type cardInsets struct{ c *Card }

func (b cardInsets) Insets() geom.Insets {
	c := b.c
	pad := orZero(c.ui, c.Padding)
	in := geom.NewUniformInsets(pad)
	if c.opens() {
		in.Right += float32(c.ui.Interface.IconSizeSmall() + c.ui.Interface.SpaceNear())
	}
	return in
}

func (b cardInsets) Draw(*unison.Canvas, geom.Rect) {}

func (c *Card) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := c.ui, c.ui.Theme.Tokens(), c.ui.Interface
	p := painterFor(gc, ui)
	r := c.ContentRect(true)
	radius := float32(m.RadiusCard())
	ground := t.ListBackground
	switch {
	case c.Selected:
		ground = t.SelectionTint
	case c.Interactive && c.hovered:
		ground = t.HoverTint
	}
	p.round(r, radius, ground)
	edge := t.Border
	if c.Interactive {
		edge = t.BorderStrong
	}
	p.outline(r, radius, float32(m.Hairline()), edge)
	if c.opens() {
		ink := t.TextFaint
		if c.hovered {
			ink = t.TextSecondary
		}
		if g, ok := icons.Glyph("chevron-right"); ok {
			s := float32(m.IconSizeSmall())
			strip := c.opensStrip()
			drawGlyph(gc, ui, g, s, ink, geom.NewRect(strip.X+(strip.Width-s)/2, strip.Y+(strip.Height-s)/2, s, s))
		}
	}
}

// ProvideAccessibility describes a pressable card as a group that can be
// pressed, with what it opens as its description.
func (c *Card) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	if c.opens() {
		n.Description = c.OpensLabel
	}
	n.Selectable, n.Selected = c.Interactive || c.Selected, c.Selected
	if c.Interactive && c.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction presses the card for a screen reader.
func (c *Card) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !c.Interactive || !c.Enabled() {
		return false
	}
	c.fire()
	return true
}
