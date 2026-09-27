package kvitui

import (
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// BeforeAfter is one value as it stands and the value something proposes to
// replace it with: the shape every screen asking the reader to approve a
// change is built from, of which kvit-cash alone has four (the import
// preview, the quarantine group, the diff before a recipe update, the review
// of an agent's proposal). The old value is on the left in the muted colour
// and the new one on the right in the ordinary colour, with an arrow between,
// so position, weight of colour and the arrow's direction all say which is
// which. Both are Figures, so a side nobody measured is an em dash: a record
// being added has no before, which is not a before of zero. Where the two are
// the same, the arrow and the second value go, since an arrow between equal
// values sends the reader looking for a difference.
type BeforeAfter struct {
	unison.Panel
	ui *UI
	// Before and After are the values, already formatted.
	Before, After string
	// Unit is drawn after each value.
	Unit string
	// BeforeMeasured and AfterMeasured false draw an em dash on that side.
	BeforeMeasured, AfterMeasured bool
	// Label says what the pair is a value of; "" for none, where a heading
	// already says it.
	Label string
	// Role is both values' type role, so the two are always the same size.
	Role TypeRole

	name          *Label
	before, after *Figure
	arrow         *Icon
}

// NewBeforeAfter returns a pair of measured values in the body role.
func NewBeforeAfter(ui *UI, label, before, after string) *BeforeAfter {
	b := &BeforeAfter{ui: ui, Label: label, Before: before, After: after, BeforeMeasured: true, AfterMeasured: true, Role: RoleBody}
	b.Self = b
	b.name = NewLabel(ui, "")
	b.name.Role, b.name.Ink = RoleSmall, InkTextMuted
	b.before = NewFigure(ui, "", "")
	b.before.Ink = InkTextMuted
	b.after = NewFigure(ui, "", "")
	b.arrow = NewIcon(ui, "arrow-right")
	b.arrow.Ink, b.arrow.Size = InkTextFaint, SizeIconSizeSmall
	for _, p := range []unison.Paneler{b.name, b.before, b.arrow, b.after} {
		p.AsPanel().SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
	}
	b.SetLayout(syncing{Layout: &spaced{FlexLayout: unison.FlexLayout{Columns: 4, VAlign: align.Middle}, ui: ui, gap: SizeSpaceNear, horizontal: true},
		sync: b.sync})
	return b
}

// Unchanged reports whether the two values are the same.
func (b *BeforeAfter) Unchanged() bool {
	return b.BeforeMeasured == b.AfterMeasured && b.Before == b.After
}

func (b *BeforeAfter) sync() {
	b.name.Text = b.Label
	b.name.Hidden = b.Label == ""
	b.before.Value, b.before.Unit, b.before.Measured, b.before.Role = b.Before, b.Unit, b.BeforeMeasured, b.Role
	b.after.Value, b.after.Unit, b.after.Measured, b.after.Role = b.After, b.Unit, b.AfterMeasured, b.Role
	b.arrow.Hidden, b.after.Hidden = b.Unchanged(), b.Unchanged()
	showOnly(b.AsPanel(), b.name, b.before, b.arrow, b.after)
	if l, ok := b.Layout().(syncing); ok {
		if s, ok := l.Layout.(*spaced); ok {
			s.Columns = max(1, len(b.Children()))
		}
	}
}

// ProvideAccessibility reads the pair as a sentence: "Amount, was 42.00 GBP,
// now 44.50 GBP", or "Date, unchanged, was 2026-08-14".
func (b *BeforeAfter) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	n := bl.Node()
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	was := "was not set"
	if b.BeforeMeasured {
		was = "was " + joinWords(b.Before, b.Unit)
	}
	now := "now not set"
	if b.AfterMeasured {
		now = "now " + joinWords(b.After, b.Unit)
	}
	switch {
	case b.Unchanged() && b.Label != "":
		n.Name = b.Label + ", unchanged, " + was
	case b.Unchanged():
		n.Name = "unchanged, " + was
	case b.Label != "":
		n.Name = b.Label + ", " + was + ", " + now
	default:
		n.Name = was + ", " + now
	}
}
