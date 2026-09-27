package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// Divider is a rule between two things: one hairline in the border colour.
// The border colour rather than the strong one, because a rule is decoration
// and not the edge of a control, and a separator drawn like the edge of
// something clickable reads as one.
type Divider struct {
	unison.Panel
	ui *UI
	// Vertical draws the rule top to bottom rather than across.
	Vertical bool
	// Inset stops the rule short of both ends, for a rule inside a padded
	// container where a full-length line would cut the padding in half.
	Inset Measure
}

// NewDivider returns a horizontal rule.
func NewDivider(ui *UI) *Divider {
	d := &Divider{ui: ui}
	d.Self = d
	d.SetSizer(d.sizes)
	d.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		painterFor(gc, ui).fill(d.rule(d.ContentRect(false)), ui.Theme.Tokens().Border)
	}
	return d
}

func (d *Divider) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	h := float32(d.ui.Interface.Hairline())
	if d.Vertical {
		return geom.NewSize(h, 0), geom.NewSize(h, 0), geom.NewSize(h, unison.DefaultMaxSize)
	}
	return geom.NewSize(0, h), geom.NewSize(0, h), geom.NewSize(unison.DefaultMaxSize, h)
}

// rule is the part of b the line covers.
func (d *Divider) rule(b geom.Rect) geom.Rect {
	var inset float32
	if d.Inset != nil {
		inset = float32(d.Inset.Of(d.ui))
	}
	if d.Vertical {
		return geom.NewRect(b.X, b.Y+inset, b.Width, max(0, b.Height-2*inset))
	}
	return geom.NewRect(b.X+inset, b.Y, max(0, b.Width-2*inset), b.Height)
}

// ProvideAccessibility says the rule is a separator.
func (d *Divider) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	if n := b.Node(); n.Role == role.Auto {
		n.Role = role.Separator
	}
}
