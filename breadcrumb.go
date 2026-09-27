package kvitui

import (
	"slices"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// Crumb is one place in a breadcrumb trail.
type Crumb struct {
	Label string
	// ID is the caller's own name for the place, handed back when the crumb
	// is followed.
	ID string
}

// Breadcrumb says where the reader is, and gives the way back.
//
// The last crumb is the current place and is not a link: a last crumb that
// could be followed would be followed and do nothing, and a screen reader
// would announce a link that goes nowhere. A long trail is cut from the
// middle rather than the end, because the first crumb says which part of the
// application this is and the last says where the reader is, and what a
// reader can do without is the middle.
type Breadcrumb struct {
	unison.Panel
	ui *UI
	// Trail is the places from the root to the current one.
	Trail []Crumb
	// MaximumVisible is how many crumbs are drawn before the middle of the
	// trail is replaced by "…"; 4 unless set.
	MaximumVisible int
	// OnActivate runs when a crumb before the current one is followed.
	OnActivate func(Crumb)

	built    []Crumb
	builtMax int
}

// NewBreadcrumb returns a breadcrumb showing a trail, root first.
func NewBreadcrumb(ui *UI, trail ...Crumb) *Breadcrumb {
	b := &Breadcrumb{ui: ui, Trail: trail, MaximumVisible: 4}
	b.Self = b
	row := &spaced{FlexLayout: unison.FlexLayout{VAlign: align.Middle}, ui: ui, gap: SizeSpaceNear, horizontal: true}
	b.SetLayout(breadcrumbLayout{b, AtLeast(ui, SizeBreadcrumbHeight, row)})
	b.sync()
	return b
}

// Shown is the crumbs as drawn: the whole trail when it is short enough, and
// otherwise the first crumb, a crumb reading "…" with no ID that cannot be
// followed, and the last few.
func (b *Breadcrumb) Shown() []Crumb {
	most := b.MaximumVisible
	if most <= 0 {
		most = 4
	}
	if len(b.Trail) <= most {
		return b.Trail
	}
	// The "…" stands for several places, and picking one of them is what the
	// sidebar is for, so it is not a link.
	tail := b.Trail[len(b.Trail)-max(1, most-2):]
	return append([]Crumb{b.Trail[0], {Label: "…"}}, tail...)
}

// sync rebuilds the crumbs when the trail has changed since they were built.
func (b *Breadcrumb) sync() {
	if slices.Equal(b.built, b.Trail) && b.builtMax == b.MaximumVisible && b.built != nil {
		return
	}
	b.built, b.builtMax = slices.Clone(b.Trail), b.MaximumVisible
	if b.built == nil {
		b.built = []Crumb{}
	}
	b.RemoveAllChildren()
	shown := b.Shown()
	elided := len(shown) != len(b.Trail)
	for i, c := range shown {
		if i > 0 {
			chevron := NewIcon(b.ui, "chevron-right")
			chevron.Ink, chevron.Size = InkTextFaint, SizeCaption
			b.AddChild(chevron)
		}
		switch {
		case i == len(shown)-1:
			here := NewLabel(b.ui, c.Label)
			here.Weight = text.Bold
			b.AddChild(here)
		case elided && i == 1:
			gap := NewLabel(b.ui, c.Label)
			gap.Ink = InkTextFaint
			b.AddChild(gap)
		default:
			link := NewLink(b.ui, c.Label)
			link.OnActivate = func() {
				if b.OnActivate != nil {
					b.OnActivate(c)
				}
			}
			b.AddChild(link)
		}
	}
	for _, c := range b.Children() {
		c.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
	}
	if l, ok := b.Layout().(breadcrumbLayout); ok {
		l.row().Columns = len(b.Children())
	}
}

// ProvideAccessibility names the trail as a group; each crumb says itself.
func (b *Breadcrumb) ProvideAccessibility(ab *unison.AccessibilityBuilder) {
	n := ab.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name = "Breadcrumb"
}

// breadcrumbLayout rebuilds the crumbs if the trail changed, then lays them
// out in a row, a near space apart, at least a breadcrumb's height tall.
type breadcrumbLayout struct {
	b *Breadcrumb
	unison.Layout
}

func (l breadcrumbLayout) row() *spaced { return l.Layout.(*atLeast).Layout.(*spaced) }

func (l breadcrumbLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	l.b.sync()
	return l.Layout.LayoutSizes(target, hint)
}

func (l breadcrumbLayout) PerformLayout(target *unison.Panel) {
	l.b.sync()
	l.Layout.PerformLayout(target)
}
