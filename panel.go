package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// Panel is a region of the window with its own ground: a sidebar, a toolbar
// strip, the area behind a group of controls. A panel is part of the
// window's furniture and its edges are the window's edges, so it has no
// corner radius; a card is an object sitting on a surface, and has one.
type Panel struct {
	unison.Panel
	ui *UI
	// RuleTop, RuleBottom, RuleLeft and RuleRight draw a hairline along that
	// edge. A panel between two others usually wants one on the side it
	// abuts and none anywhere else.
	RuleTop, RuleBottom, RuleLeft, RuleRight bool
}

// NewPanel returns a panel with no rules.
func NewPanel(ui *UI) *Panel {
	p := &Panel{ui: ui}
	p.Self = p
	p.DrawCallback = p.draw
	p.DrawOverCallback = p.drawRules
	return p
}

// draw fills the panel's ground, under its children.
func (p *Panel) draw(gc *unison.Canvas, _ geom.Rect) {
	painterFor(gc, p.ui).fill(p.ContentRect(true), p.ui.Theme.Tokens().PanelBackground)
}

// drawRules draws the rules over the panel's children. They are the panel's
// own edges and the content sits inside them, but anything opaque a child
// draws at an edge would otherwise paint the rule out. The sidebar is where
// it showed in kvit-cash: a region's scroll bar is the panel's ground drawn
// right up against the rule, so the line down the sidebar's side went
// missing for the height of the scrolling account list.
func (p *Panel) drawRules(gc *unison.Canvas, _ geom.Rect) {
	t := p.ui.Theme.Tokens()
	pt := painterFor(gc, p.ui)
	b := p.ContentRect(true)
	h := float32(p.ui.Interface.Hairline())
	if p.RuleTop {
		pt.fill(geom.NewRect(b.X, b.Y, b.Width, h), t.Border)
	}
	if p.RuleBottom {
		pt.fill(geom.NewRect(b.X, b.Bottom()-h, b.Width, h), t.Border)
	}
	if p.RuleLeft {
		pt.fill(geom.NewRect(b.X, b.Y, h, b.Height), t.Border)
	}
	if p.RuleRight {
		pt.fill(geom.NewRect(b.Right()-h, b.Y, h, b.Height), t.Border)
	}
}
