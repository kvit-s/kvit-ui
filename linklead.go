package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// SetLead draws p before the link's words in place of its symbol: a mark the
// application draws its own way, such as one that moves while something
// runs. It is part of the link: a press on it follows the link, and the
// pointer resting on it rests on the link, which the link arranges by taking
// over p's pointer callbacks. Each link needs a panel of its own. nil takes
// the lead away.
func (l *Link) SetLead(p unison.Paneler) {
	// A panel with a layout ignores its sizer, so the layout gives the
	// sizer's answer and places the lead.
	if p == nil {
		l.SetLayout(nil)
	} else {
		l.SetLayout(linkLayout{l})
	}
	l.swapLead(l.lead, p)
	l.lead = p
}

// Lead is the panel SetLead gave, nil for none.
func (l *Link) Lead() unison.Paneler { return l.lead }

func (l *Link) leadSize() geom.Size {
	_, p, _ := l.lead.AsPanel().Sizes(geom.Size{})
	return p
}

// linkLayout is a link with a lead: its sizes are the link's own, and the
// lead sits where the symbol would, centred on the link's height.
type linkLayout struct{ l *Link }

func (y linkLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	return y.l.sizes(hint)
}

func (y linkLayout) PerformLayout(*unison.Panel) {
	if y.l.lead == nil {
		return
	}
	size := y.l.leadSize()
	b := y.l.ContentRect(false)
	y.l.lead.AsPanel().SetFrameRect(geom.NewRect(b.X, b.Y+(b.Height-size.Height)/2, size.Width, size.Height))
}
