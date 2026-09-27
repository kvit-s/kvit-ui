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
	return p
}

func (p *Panel) draw(gc *unison.Canvas, _ geom.Rect) {
	t := p.ui.Theme.Tokens()
	pt := painterFor(gc, p.ui)
	b := p.ContentRect(true)
	pt.fill(b, t.PanelBackground)
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
