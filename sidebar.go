package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// Sidebar is the list of places down the left edge. Collapsed is the state a
// window puts it in when it is narrow: the rail, where every item is its
// symbol alone. Every SidebarItem in it takes the sidebar's state rather than
// being told one by one, which is what stops one item in a rail still
// drawing its label.
type Sidebar struct {
	unison.Panel
	ui *UI
	// Collapsed draws the sidebar as its rail.
	Collapsed bool
}

// NewSidebar returns an expanded sidebar holding items, top to bottom.
func NewSidebar(ui *UI, items ...unison.Paneler) *Sidebar {
	s := &Sidebar{ui: ui}
	s.Self = s
	for _, it := range items {
		if it.AsPanel().LayoutData() == nil {
			it.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		}
		s.AddChild(it)
	}
	column := &spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: Px(0)}
	s.SetLayout(syncing{Layout: &sidebarLayout{s: s, column: column}, sync: s.sync})
	return s
}

func (s *Sidebar) sync() {
	for _, c := range s.Children() {
		if item, ok := c.Self.(*SidebarItem); ok {
			item.Collapsed = s.Collapsed
		}
	}
}

// sidebarLayout stacks the items and asks for the sidebar's width, or the
// rail's when collapsed.
type sidebarLayout struct {
	s      *Sidebar
	column unison.Layout
}

func (l *sidebarLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	minSize, prefSize, maxSize = l.column.LayoutSizes(target, hint)
	m := l.s.ui.Interface
	prefSize.Width = float32(m.SidebarWidth())
	if l.s.Collapsed {
		prefSize.Width = float32(m.RailWidth())
	}
	minSize.Width = min(minSize.Width, prefSize.Width)
	return minSize, prefSize, maxSize
}

func (l *sidebarLayout) PerformLayout(target *unison.Panel) { l.column.PerformLayout(target) }

// ProvideAccessibility names the sidebar as the application's navigation.
func (s *Sidebar) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.List
	}
	n.Name = "Navigation"
}

// SidebarItem is one place in the sidebar. The selected item is marked by a
// bar down its leading edge as well as by a tint: in the high-contrast theme
// the tint alone is a two-percent difference in lightness, invisible in a
// screenshot. The symbol is required because of the rail: collapsed, the
// label is gone and the symbol is the whole item, and an item with only
// words would vanish.
type SidebarItem struct {
	control
	// Text names the place.
	Text string
	// Symbol is the meaning name of the item's symbol.
	Symbol string
	// Selected marks the place the window is showing.
	Selected bool
	// Collapsed draws the item as part of the rail; its Sidebar sets it.
	Collapsed bool
	// Count is how many things are waiting there, drawn as a badge; below
	// zero draws none, and zero draws no badge but is still announced.
	Count int
	// Counted is the word for one of the things counted, so a screen reader
	// hears "Review, 1 decision"; "" says "items".
	Counted string
	// CountedPlural is the word for several; "" adds an s to Counted.
	CountedPlural string
	// CountInRail keeps the badge once the sidebar collapses, on the
	// symbol's upper corner. Off by default: a rail is a strip of symbols,
	// and an application whose rail has to keep a backlog in front of the
	// reader turns it on.
	CountInRail bool
	// CountMax is the largest count the badge writes out, 99 unless set.
	// Raise it where the place the item points at states the true figure,
	// or the badge's "99+" and that figure disagree about the same thing.
	CountMax int
	// OnClick runs when the item is pressed.
	OnClick func()

	badge *Badge
}

// NewSidebarItem returns an unselected item with no count.
func NewSidebarItem(ui *UI, label, symbol string) *SidebarItem {
	it := &SidebarItem{Text: label, Symbol: symbol, Count: -1, CountMax: 99}
	it.Self = it
	it.initControl(ui, func() {
		if it.OnClick != nil {
			it.OnClick()
		}
	}, unison.KeySpace)
	it.ringRadius = func() float32 { return 0 }
	it.badge = NewBadge(ui, 0)
	it.badge.Tone = BadgeNeutral
	it.AddChild(it.badge)
	it.SetLayout(syncing{Layout: sidebarItemLayout{it}, sync: it.sync})
	it.DrawCallback = it.draw
	// In the rail the label is gone, and the tooltip is the only thing that
	// says where the item goes.
	it.tip = func() string {
		if it.Collapsed {
			return it.Text
		}
		return ""
	}
	return it
}

// sync hands the count to the badge while it is to be drawn, and none
// otherwise.
func (it *SidebarItem) sync() {
	b := it.badge
	b.Count = 0
	if it.Count > 0 && (!it.Collapsed || it.CountInRail) {
		b.Count = it.Count
	}
	b.Max, b.Counted, b.CountedPlural = it.CountMax, it.Counted, it.CountedPlural
}

// badgeRight is the space between the badge and the item's right edge.
func (it *SidebarItem) badgeRight() float32 {
	if it.Collapsed {
		return float32(it.ui.Interface.SpaceTight())
	}
	return float32(it.ui.Interface.Space())
}

type sidebarItemLayout struct{ it *SidebarItem }

func (l sidebarItemLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.it.ui.Interface
	h := float32(m.RowHeightSlim())
	return geom.NewSize(float32(m.RailWidth()), h), geom.NewSize(float32(m.SidebarWidth()), h), geom.NewSize(unison.DefaultMaxSize, h)
}

// PerformLayout puts the badge at the right, centred on the row, or on the
// symbol's upper corner in the rail.
func (l sidebarItemLayout) PerformLayout(target *unison.Panel) {
	it := l.it
	r := target.ContentRect(false)
	_, bp, _ := it.badge.Sizes(geom.Size{})
	y := (r.Height - bp.Height) / 2
	if it.Collapsed {
		y -= float32(it.ui.Interface.SpaceNear())
	}
	it.badge.SetFrameRect(geom.NewRect(r.Right()-it.badgeRight()-bp.Width, r.Y+y, bp.Width, bp.Height))
}

func (it *SidebarItem) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := it.ui, it.ui.Theme.Tokens(), it.ui.Interface
	p := painterFor(gc, ui)
	r := it.ContentRect(false)
	switch {
	case it.Selected:
		p.fill(r, t.SelectionTint)
		p.fill(geom.NewRect(r.X, r.Y, float32(m.SpaceTight()), r.Height), t.Accent)
	case it.hovered:
		p.fill(r, t.HoverTint)
	}
	s := float32(m.IconSizeSmall())
	x := r.X + float32(m.SpaceLoose())
	if it.Collapsed {
		x = r.X + (r.Width-s)/2
	}
	ink := t.TextMuted
	if it.Selected {
		ink = t.Accent
	}
	if g, ok := icons.Glyph(it.Symbol); ok {
		drawGlyph(gc, ui, g, s, ink, geom.NewRect(x, r.Y+(r.Height-s)/2, s, s))
	}
	if it.Collapsed {
		return
	}
	weight, ink := text.Regular, t.TextSecondary
	if it.Selected {
		weight, ink = text.Bold, t.TextPrimary
	}
	left := x + s + float32(m.Space())
	// The label stops short of where a badge goes whether or not one is
	// drawn, so a label does not reach further on an item with no count.
	_, bp, _ := it.badge.Sizes(geom.Size{})
	right := r.Right() - it.badgeRight() - max(float32(m.PillHeight()), bp.Width) - float32(m.SpaceNear())
	l := ui.Fonts.Layout([]text.Span{{Text: it.Text, Style: ui.Chrome(ui.Size(RoleBody), weight, ink)}},
		text.Options{MaxWidth: max(0, right-left), Elide: true})
	_, h := l.Size()
	l.Draw(gc, left, r.Y+(r.Height-h)/2)
}

// Name is what a screen reader is told: the label, and the count with its
// noun when there is one, "Review, 11 decisions".
func (it *SidebarItem) Name() string {
	if it.Count < 0 {
		return it.Text
	}
	if it.Counted == "" {
		return it.Text + ", " + it.ui.CountPhrase(it.Count, "item", "")
	}
	return it.Text + ", " + it.ui.CountPhrase(it.Count, it.Counted, it.CountedPlural)
}

// ProvideAccessibility describes the item as one place in the list, and
// says whether it is the selected one.
func (it *SidebarItem) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.ListItem
	}
	n.Name = it.Name()
	n.Selectable = true
	n.Selected = it.Selected
	if it.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction presses the item for a screen reader.
func (it *SidebarItem) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !it.Enabled() {
		return false
	}
	it.fire()
	return true
}
