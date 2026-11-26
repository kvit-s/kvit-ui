package kvitui

import (
	"slices"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// Sidebar is the list of places down the left edge. Collapsed is the state a
// window puts it in when it is narrow: the rail, where every item is its
// symbol alone. Every SidebarItem in it takes the sidebar's state rather than
// being told one by one, which is what stops one item in a rail still
// drawing its label.
//
// A SidebarHeading among the items names the group of places under it, and a
// SidebarItem holds the places inside it (SidebarItem.Add), which it shows
// while it is open. The sidebar's children are the rows shown: the places
// inside a closed item are not among them, and the rail shows the top level
// alone, without the headings.
//
// The places scroll when they do not fit above the foot, which keeps its
// place: with the wheel, with a bar over their right-hand edge that is there
// only while there is something out of sight, and to wherever the keyboard
// goes. The rail scrolls without the bar, having no room for one.
type Sidebar struct {
	unison.Panel
	ui *UI
	// Collapsed draws the sidebar as its rail.
	Collapsed bool

	entries []*unison.Panel   // the top level, as given: places and headings
	footer  []*unison.Panel   // the items held at the foot
	places  *unison.Panel     // holds the rows shown, and scrolls them
	column  unison.Layout     // stacks the rows in places
	bar     *unison.ScrollBar // over the right-hand edge of places
	shown   []*SidebarItem    // the items shown, the foot's after the places'
	scrollY float32
	ease    wheelScroll
}

// NewSidebar returns an expanded sidebar holding items, top to bottom.
func NewSidebar(ui *UI, items ...unison.Paneler) *Sidebar {
	s := &Sidebar{ui: ui}
	s.Self = s
	s.column = &spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: Px(0)}
	s.places = unison.NewPanel()
	s.places.SetLayout(sidebarPlaces{s})
	s.places.FocusChangeInHierarchyCallback = func(_, to *unison.Panel) { s.reveal(to) }
	for _, it := range items {
		fillSidebar(it.AsPanel())
		s.entries = append(s.entries, it.AsPanel())
		s.places.AddChild(it)
	}
	s.bar = NewScrollBar(ui, false)
	s.bar.ChangedCallback = func() { s.scrollTo(s.bar.Value()) }
	s.ease = wheelScroll{ui: ui,
		at:     func() (x, y float32) { return 0, s.scrollY },
		travel: func() (x, y float32) { return 0, s.travel() },
		moveTo: func(_, y float32) { s.scrollTo(y) }}
	s.MouseWheelCallback = func(_, delta geom.Point, _ mod.Modifiers) bool { return s.ease.wheel(delta) }
	// unison draws the first child on top: the bar is over the places.
	s.AddChild(s.bar)
	s.AddChild(s.places)
	s.SetLayout(syncing{Layout: sidebarLayout{s}, sync: s.sync})
	return s
}

// Rows are the rows shown above the foot, top to bottom: the places and the
// headings, without the places inside a closed item.
func (s *Sidebar) Rows() []*unison.Panel { return slices.Clone(s.places.Children()) }

// travel is how far the places can scroll.
func (s *Sidebar) travel() float32 {
	r := s.places.ContentRect(false)
	_, pref, _ := s.column.LayoutSizes(s.places, geom.NewSize(r.Width, 0))
	return max(0, pref.Height-r.Height)
}

// Position is how far down the places are scrolled.
func (s *Sidebar) Position() float32 { return s.scrollY }

// ScrollTo scrolls the places so y is at the top, kept within their ends.
func (s *Sidebar) ScrollTo(y float32) { s.scrollTo(y) }

func (s *Sidebar) scrollTo(y float32) {
	y = max(0, min(s.travel(), y))
	if y != s.scrollY {
		s.scrollY = y
		s.places.MarkForLayoutAndRedraw()
		s.MarkForLayoutAndRedraw()
	}
}

// reveal brings the row the keyboard has just reached into view.
func (s *Sidebar) reveal(to *unison.Panel) {
	if to == nil || to.Parent() != s.places {
		return
	}
	// Rows that have just come are placed before one of them is looked for.
	s.ValidateLayout()
	box := to.FrameRect()
	top, view := box.Y+s.scrollY, s.places.ContentRect(false).Height
	switch {
	case top < s.scrollY:
		s.scrollTo(top)
	case top+box.Height > s.scrollY+view:
		s.scrollTo(top + box.Height - view)
	}
}

// sidebarPlaces stacks the rows and moves them up by how far they are
// scrolled.
type sidebarPlaces struct{ s *Sidebar }

func (l sidebarPlaces) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	return l.s.column.LayoutSizes(target, hint)
}

func (l sidebarPlaces) PerformLayout(target *unison.Panel) {
	l.s.column.PerformLayout(target)
	l.s.scrollY = max(0, min(l.s.travel(), l.s.scrollY))
	for _, c := range target.Children() {
		r := c.FrameRect()
		r.Y -= l.s.scrollY
		c.SetFrameRect(toWholePixels(target, r))
	}
}

// AddFooter puts items at the foot of the sidebar, below the places, in
// order: Help, then Settings last, where the Windows form puts the settings
// command (WIN-16, CMD-3).
func (s *Sidebar) AddFooter(items ...unison.Paneler) {
	for _, it := range items {
		fillSidebar(it.AsPanel())
		s.footer = append(s.footer, it.AsPanel())
		s.AddChild(it)
	}
	s.MarkForLayoutAndRedraw()
}

// fillSidebar makes a row as wide as the sidebar unless it says otherwise.
func fillSidebar(p *unison.Panel) {
	if p.LayoutData() == nil {
		p.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
	}
}

// each calls f with every item of the sidebar, the places inside closed items
// included, top to bottom.
func (s *Sidebar) each(f func(*SidebarItem)) {
	var walk func(it *SidebarItem)
	walk = func(it *SidebarItem) {
		f(it)
		for _, c := range it.items {
			walk(c)
		}
	}
	for _, p := range slices.Concat(s.entries, s.footer) {
		if it, ok := p.Self.(*SidebarItem); ok {
			walk(it)
		}
	}
}

// Commands are the commands the sidebar's items run.
func (s *Sidebar) Commands() []*Command {
	var out []*Command
	s.each(func(it *SidebarItem) {
		if it.Command != nil {
			out = append(out, it.Command)
		}
	})
	return out
}

// sync makes the sidebar's children the rows shown, and tells each item the
// sidebar's state and how deep it is.
func (s *Sidebar) sync() {
	// A row added with AddChild joins the top level.
	known := map[*unison.Panel]bool{s.places: true, s.bar.AsPanel(): true}
	for _, p := range slices.Concat(s.entries, s.footer) {
		known[p] = true
	}
	nests := false
	s.each(func(it *SidebarItem) {
		known[it.AsPanel()] = true
		nests = nests || len(it.items) > 0
	})
	for _, c := range s.Children() {
		if !known[c] {
			fillSidebar(c)
			s.entries = append(s.entries, c)
		}
	}
	var rows []*unison.Panel
	s.shown = s.shown[:0]
	var show func(it *SidebarItem, depth int)
	show = func(it *SidebarItem, depth int) {
		it.Collapsed, it.depth, it.inTree, it.sidebar = s.Collapsed, depth, nests, s
		rows = append(rows, it.AsPanel())
		s.shown = append(s.shown, it)
		if s.Collapsed || !it.Expanded {
			return
		}
		for _, c := range it.items {
			if !c.Hidden {
				show(c, depth+1)
			}
		}
	}
	for _, p := range s.entries {
		switch row := p.Self.(type) {
		case *SidebarItem:
			if !p.Hidden {
				show(row, 0)
			}
		case *SidebarHeading:
			row.inTree = nests
			if !p.Hidden && !s.Collapsed {
				rows = append(rows, p)
			}
		default:
			if !p.Hidden {
				rows = append(rows, p)
			}
		}
	}
	showInPlace(s.places, rows)
	// The foot's items are rows of the sidebar itself, and one group with
	// the places for the keyboard.
	for _, p := range s.footer {
		if it, ok := p.Self.(*SidebarItem); ok && !p.Hidden {
			it.Collapsed, it.depth, it.inTree, it.sidebar = s.Collapsed, 0, nests, s
			s.shown = append(s.shown, it)
		}
	}
}

// showInPlace makes want the children of p, leaving a child that stays where
// it is, so a row that holds the keyboard keeps it while rows come and go
// around it.
func showInPlace(p *unison.Panel, want []*unison.Panel) {
	if slices.Equal(p.Children(), want) {
		return
	}
	for _, c := range slices.Clone(p.Children()) {
		if !slices.Contains(want, c) {
			p.RemoveChild(c)
		}
	}
	for i, w := range want {
		if have := p.Children(); i >= len(have) || have[i] != w {
			p.AddChildAtIndex(w, i)
		}
	}
}

// nests reports whether any place holds places inside it.
func (s *Sidebar) nests() bool {
	nests := false
	s.each(func(it *SidebarItem) { nests = nests || len(it.items) > 0 })
	return nests
}

// sidebarLayout gives the foot's items the foot, in order, the last at the
// bottom, and the places the rest. It asks for the sidebar's width, or the
// rail's when collapsed, and for the height of every row; given less, the
// places scroll.
type sidebarLayout struct{ s *Sidebar }

// foot is the height the foot's items take at a width.
func (l sidebarLayout) foot(width float32) float32 {
	var h float32
	for _, p := range l.s.footer {
		if !p.Hidden {
			_, pref, _ := p.Sizes(geom.NewSize(width, 0))
			h += pref.Height
		}
	}
	return h
}

func (l sidebarLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.s.ui.Interface
	_, places, _ := l.s.places.Sizes(geom.NewSize(hint.Width, 0))
	foot := l.foot(hint.Width)
	prefSize = geom.NewSize(float32(m.SidebarWidth()), places.Height+foot)
	if l.s.Collapsed {
		prefSize.Width = float32(m.RailWidth())
	}
	// It can be as short as its foot and one row, scrolling the rest.
	minSize = geom.NewSize(min(float32(m.RailWidth()), prefSize.Width), foot+min(places.Height, float32(m.RowHeightSlim())))
	return minSize, prefSize, geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l sidebarLayout) PerformLayout(target *unison.Panel) {
	s := l.s
	r := target.ContentRect(false)
	y := r.Bottom()
	for i := len(s.footer) - 1; i >= 0; i-- {
		p := s.footer[i]
		if p.Hidden {
			continue
		}
		_, pref, _ := p.Sizes(geom.NewSize(r.Width, 0))
		y -= pref.Height
		p.SetFrameRect(geom.NewRect(r.X, y, r.Width, pref.Height))
	}
	view := geom.NewRect(r.X, r.Y, r.Width, max(0, y-r.Y))
	s.places.SetFrameRect(view)
	s.places.NeedsLayout = true
	s.places.ValidateLayout()
	// The bar is there only while something is out of sight, and never in
	// the rail, which has no room for it.
	travel := s.travel()
	width := float32(s.ui.Interface.SpaceWide())
	s.bar.Hidden = travel <= 0 || s.Collapsed
	s.bar.SetFrameRect(geom.NewRect(view.Right()-width, view.Y, width, view.Height))
	s.bar.SetRange(s.scrollY, view.Height, view.Height+travel)
}

// ProvideAccessibility names the sidebar as the application's navigation: a
// list of places, or a tree where a place holds places inside it.
func (s *Sidebar) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.List
		if s.nests() {
			n.Role = role.Tree
		}
	}
	n.Name = "Navigation"
}

// SidebarHeading names the group of places under it in a sidebar: "Folders"
// above the folders, "Tags" above the tags. It is words alone and takes no
// press and no key; the rail leaves it out.
type SidebarHeading struct {
	unison.Panel
	ui *UI
	// Text names the group.
	Text string

	inTree bool // some place of its sidebar holds places; its Sidebar sets it
}

// NewSidebarHeading returns a heading for the places that follow it.
func NewSidebarHeading(ui *UI, label string) *SidebarHeading {
	h := &SidebarHeading{ui: ui, Text: label}
	h.Self = h
	h.SetSizer(func(geom.Size) (minSize, prefSize, maxSize geom.Size) {
		m := ui.Interface
		height := float32(m.RowHeightSlim())
		return geom.NewSize(float32(m.RailWidth()), height), geom.NewSize(float32(m.SidebarWidth()), height),
			geom.NewSize(unison.DefaultMaxSize, height)
	})
	h.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		m := ui.Interface
		r := h.ContentRect(false)
		// The words start where the places' symbols do.
		left := r.X + float32(m.SpaceLoose())
		if h.inTree {
			left += float32(m.SpaceNear())
		}
		l := ui.Fonts.Layout([]text.Span{{Text: h.Text, Style: ui.Chrome(ui.Size(RoleCaption), text.Semibold, ui.Theme.Tokens().TextMuted)}},
			text.Options{MaxWidth: max(0, r.Right()-float32(m.Space())-left), Elide: true})
		// The words stand on the foot of the row, nearer the places they name
		// than the group above.
		_, height := l.Size()
		l.Draw(gc, left, r.Bottom()-float32(m.SpaceSnug())-height)
	}
	return h
}

// ProvideAccessibility gives a screen reader the heading's words.
func (h *SidebarHeading) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Heading
	}
	n.Name = h.Text
}

// SidebarItem is one place in the sidebar. The selected item is marked by a
// bar down its leading edge as well as by a tint: in the high-contrast theme
// the tint alone is a two-percent difference in lightness, invisible in a
// screenshot. The symbol, or a lead in its place (SetLead), is required
// because of the rail: collapsed, the label is gone and the symbol is the
// whole item, and an item with only words would vanish.
//
// An item holds the places inside it (Add): a folder its folders. It then has
// a mark before its symbol that opens and closes it, and Right and Left open
// and close it from the keyboard, Right stepping into an open one and Left
// out to the item a place is inside. Each level stands a loose space further
// in. A closed item, and an item in the rail, takes the bar of the selected
// place when that place is inside it, so the reader sees where they are.
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
	// Command, when set, is what the item runs, and gives its words and
	// symbol where they are empty: the settings command at the foot of the
	// sidebar is one.
	Command *Command
	// Expanded shows the places inside the item (Add).
	Expanded bool
	// OnToggle runs after the reader opened or closed the item, with the new
	// state.
	OnToggle func(expanded bool)

	badge       *Badge
	lead        unison.Paneler // drawn in place of the symbol (SetLead)
	items       []*SidebarItem // the places inside it
	within      *SidebarItem   // the item it is inside, nil at the top level
	depth       int            // how many items it is inside; its Sidebar sets it
	inTree      bool           // some place of its sidebar holds places
	markPressed bool           // the pointer went down on the open-and-close mark
	sidebar     *Sidebar       // the sidebar showing it, which sets it
}

// group is the items the arrow keys move among and Tab stops at once: every
// item its sidebar shows, the foot's included, or the items beside it when
// it stands outside a sidebar.
func (it *SidebarItem) group() []*SidebarItem {
	if s := it.sidebar; s != nil && (it.Parent() == s.places || it.Parent() == s.AsPanel()) {
		return s.shown
	}
	return siblings[*SidebarItem](it.AsPanel())
}

// NewSidebarCommand returns an item that runs a command, named and drawn by
// it: the settings command at the foot of the sidebar is one (WIN-16).
func NewSidebarCommand(ui *UI, c *Command) *SidebarItem {
	name := c.Name
	if name == "" {
		name = ui.StdName(c.Std)
	}
	it := NewSidebarItem(ui, name, c.SymbolName())
	it.Command = c
	return it
}

// NewSidebarItem returns an unselected item with no count.
func NewSidebarItem(ui *UI, label, symbol string) *SidebarItem {
	it := &SidebarItem{Text: label, Symbol: symbol, Count: -1, CountMax: 99}
	it.Self = it
	it.initControl(ui, func() {
		if it.Command != nil {
			ui.runCommand(it.Command)
		}
		if it.OnClick != nil {
			it.OnClick()
		}
	}, unison.KeySpace, unison.KeyReturn, unison.KeyNumPadEnter)
	it.ringRadius = func() float32 { return 0 }
	// The items are one stop in the Tab order; Up and Down move the keyboard
	// among them without opening one, and Space or Enter opens it (KEY-9).
	keys := it.KeyDownCallback
	it.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if mods&mod.NonSticky == 0 && it.stepInOrOut(key) {
			return true
		}
		if mods&mod.NonSticky == 0 {
			if next, moved := roveStep(it.group(), it, key, true, false); moved {
				if next != nil {
					next.SetFocusable(true)
					next.Focus()
				}
				return true
			}
		}
		return keys(key, mods, repeat)
	}
	// A press on the open-and-close mark opens or closes the item and goes
	// nowhere; a second press in a row on the item does the same.
	down, up := it.MouseDownCallback, it.MouseUpCallback
	it.MouseDownCallback = func(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
		it.markPressed = button == unison.ButtonLeft && it.Enabled() && it.onMark(where)
		if it.markPressed {
			return true
		}
		if button == unison.ButtonLeft && clicks == 2 && it.opens() && it.Enabled() {
			it.toggle()
		}
		// Only the group's one stop takes the keyboard until another member
		// is pressed, which then has to be able to.
		if button == unison.ButtonLeft && it.Enabled() {
			it.SetFocusable(true)
		}
		return down(where, button, clicks, mods)
	}
	it.MouseUpCallback = func(where geom.Point, button int, mods mod.Modifiers) bool {
		if !it.markPressed {
			return up(where, button, mods)
		}
		it.markPressed = false
		if it.onMark(where) {
			it.toggle()
		}
		return true
	}
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

// SetLead draws p in place of the item's symbol: a picture of the place's own,
// such as a project's icon. It is part of the item: a press on it presses the
// item, the pointer resting on it rests on the item, and a screen reader hears
// the item's label and nothing of p. Centred on where the symbol stands, in the
// sidebar and in the rail, it may be larger than the symbol and the label stays
// where it is, so names with a picture and names with a symbol line up; one too
// large to stop short of the label is drawn smaller. Each item needs a panel of
// its own. nil gives the place back to the symbol.
func (it *SidebarItem) SetLead(p unison.Paneler) {
	it.swapLead(it.lead, p)
	it.lead = p
	if p != nil {
		silenceLead(p.AsPanel())
	}
}

// Lead is the panel SetLead gave, nil for none.
func (it *SidebarItem) Lead() unison.Paneler { return it.lead }

// Add puts places inside the item, after those it holds: a folder's folders.
// They are shown under it, a level further in, while the item is Expanded.
func (it *SidebarItem) Add(places ...*SidebarItem) {
	for _, p := range places {
		fillSidebar(p.AsPanel())
		p.within = it
		it.items = append(it.items, p)
	}
	it.rowsChanged()
}

// Clear takes out the places inside the item.
func (it *SidebarItem) Clear() {
	for _, p := range it.items {
		p.within = nil
	}
	it.items = nil
	it.rowsChanged()
}

// Items are the places inside the item.
func (it *SidebarItem) Items() []*SidebarItem { return slices.Clone(it.items) }

// Depth is how many items the item is inside: zero at the top level. Its
// sidebar works it out at layout.
func (it *SidebarItem) Depth() int { return it.depth }

// Within is the item this one is inside, nil at the sidebar's top level.
func (it *SidebarItem) Within() *SidebarItem { return it.within }

// SetExpanded opens or closes the item. The keyboard, when it is on a place
// inside an item that closes, moves to the item.
func (it *SidebarItem) SetExpanded(open bool) {
	if it.Expanded == open {
		return
	}
	it.Expanded = open
	if !open {
		for _, sib := range it.group() {
			if sib.Focused() && sib.inside(it) {
				it.SetFocusable(true)
				it.Focus()
			}
		}
	}
	it.rowsChanged()
}

// inside reports whether the item is somewhere inside other.
func (it *SidebarItem) inside(other *SidebarItem) bool {
	for up := it.within; up != nil; up = up.within {
		if up == other {
			return true
		}
	}
	return false
}

// rowsChanged has the sidebar work out its rows again, and everything above
// it its size, since a sidebar that scrolls is as tall as its rows.
func (it *SidebarItem) rowsChanged() {
	for p := it.Parent(); p != nil; p = p.Parent() {
		p.NeedsLayout = true
	}
	it.MarkForRedraw()
}

// opens reports whether the item shows a mark that opens and closes it.
func (it *SidebarItem) opens() bool { return len(it.items) > 0 && !it.Collapsed }

// toggle opens or closes the item for the reader.
func (it *SidebarItem) toggle() {
	it.SetExpanded(!it.Expanded)
	if it.OnToggle != nil {
		it.OnToggle(it.Expanded)
	}
}

// onMark reports whether a point is on the open-and-close mark: anywhere
// before the symbol.
func (it *SidebarItem) onMark(where geom.Point) bool {
	r := it.ContentRect(false)
	return it.opens() && where.In(r) && where.X < it.symbolBox(r).X
}

// stepInOrOut answers Right and Left on an item of a sidebar whose places
// hold places: Right opens the item, or steps to the first place inside an
// open one; Left closes it, or steps out to the item it is inside.
func (it *SidebarItem) stepInOrOut(key unison.KeyCode) bool {
	to := it
	switch {
	case it.Collapsed || !it.Enabled():
		return false
	case key == unison.KeyRight && it.opens() && !it.Expanded,
		key == unison.KeyLeft && it.opens() && it.Expanded:
		it.toggle()
		return true
	case key == unison.KeyRight && it.opens():
		to = it.items[0]
	case key == unison.KeyLeft && it.within != nil:
		to = it.within
	default:
		return false
	}
	if to.Parent() != nil && !to.Hidden {
		to.SetFocusable(true)
		to.Focus()
	}
	return true
}

// holdsSelected reports whether the selected place is inside the item and
// not shown: the item is closed, or the sidebar is its rail.
func (it *SidebarItem) holdsSelected() bool {
	if it.Expanded && !it.Collapsed {
		return false
	}
	var selected func(in *SidebarItem) bool
	selected = func(in *SidebarItem) bool {
		return slices.ContainsFunc(in.items, func(p *SidebarItem) bool { return p.Selected || selected(p) })
	}
	return selected(it)
}

// sync hands the count to the badge while it is to be drawn, and none
// otherwise.
func (it *SidebarItem) sync() {
	oneStop(it.group(), func(o *SidebarItem) bool { return o.Selected })
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

// symbolBox is where the symbol stands in the item's content rectangle r, and
// what a lead is centred on.
func (it *SidebarItem) symbolBox(r geom.Rect) geom.Rect {
	m := it.ui.Interface
	s := float32(m.IconSizeSmall())
	x := r.X + float32(m.SpaceLoose()+it.depth*m.SpaceLoose())
	if it.inTree {
		// Room for the open-and-close mark, on every row, so the symbols of
		// places that open and places that do not stay in one column.
		x += float32(m.SpaceNear())
	}
	if it.Collapsed {
		x = r.X + (r.Width-s)/2
	}
	return geom.NewRect(x, r.Y+(r.Height-s)/2, s, s)
}

// markBox is where the open-and-close mark stands: against the symbol, on
// the side the row starts.
func (it *SidebarItem) markBox(r geom.Rect) geom.Rect {
	s := float32(it.ui.Interface.SpaceWide())
	return geom.NewRect(it.symbolBox(r).X-s-float32(it.ui.Interface.SpaceTight()), r.Y+(r.Height-s)/2, s, s)
}

// placeLead centres the lead on the symbol's place, reaching at most to a tight
// space short of the label.
func (it *SidebarItem) placeLead(r geom.Rect) {
	if it.lead == nil {
		return
	}
	m := it.ui.Interface
	box := it.symbolBox(r)
	lead := it.lead.AsPanel()
	size := leadSize(lead, geom.NewSize(box.Width+2*float32(m.Space()-m.SpaceTight()), r.Height))
	at := geom.NewRect(box.CenterX()-size.Width/2, box.CenterY()-size.Height/2, size.Width, size.Height)
	lead.SetFrameRect(toWholePixels(it.AsPanel(), at))
}

// PerformLayout puts the lead where the symbol stands, and the badge at the
// right, centred on the row, or on the symbol's upper corner in the rail.
func (l sidebarItemLayout) PerformLayout(target *unison.Panel) {
	it := l.it
	r := target.ContentRect(false)
	_, bp, _ := it.badge.Sizes(geom.Size{})
	y := (r.Height - bp.Height) / 2
	if it.Collapsed {
		y -= float32(it.ui.Interface.SpaceNear())
	}
	it.badge.SetFrameRect(geom.NewRect(r.Right()-it.badgeRight()-bp.Width, r.Y+y, bp.Width, bp.Height))
	it.placeLead(r)
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
	if !it.Selected && it.holdsSelected() {
		p.fill(geom.NewRect(r.X, r.Y, float32(m.SpaceTight()), r.Height), t.Accent)
	}
	if it.opens() {
		name := "chevron-right"
		if it.Expanded {
			name = "chevron-down"
		}
		if g, ok := icons.Glyph(name); ok {
			mark := it.markBox(r)
			drawGlyph(gc, ui, g, mark.Width, t.TextMuted, mark)
		}
	}
	box := it.symbolBox(r)
	ink := t.TextMuted
	if it.Selected {
		ink = t.Accent
	}
	if g, ok := icons.Glyph(it.Symbol); ok && it.lead == nil {
		drawGlyph(gc, ui, g, box.Width, ink, box)
	}
	if it.Collapsed {
		return
	}
	weight, ink := text.Regular, t.TextSecondary
	if it.Selected {
		weight, ink = text.Bold, t.TextPrimary
	}
	left := box.Right() + float32(m.Space())
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
		if it.inTree {
			n.Role = role.Row
		}
	}
	n.Name = it.Name()
	n.Selectable = true
	n.Selected = it.Selected
	if it.inTree {
		n.Level = it.depth + 1
	}
	if it.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
	if it.opens() {
		n.Expandable, n.Expanded = true, it.Expanded
		if it.Enabled() {
			n.Actions = n.Actions.With(accessibility.Expand, accessibility.Collapse)
		}
	}
}

// PerformAccessibilityAction presses the item for a screen reader.
func (it *SidebarItem) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if !it.Enabled() {
		return false
	}
	switch req.Action {
	case accessibility.Press:
		it.fire()
	case accessibility.Expand, accessibility.Collapse:
		if !it.opens() {
			return false
		}
		if it.Expanded != (req.Action == accessibility.Expand) {
			it.toggle()
		}
	default:
		return false
	}
	return true
}
