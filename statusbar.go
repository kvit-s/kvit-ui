package kvitui

import (
	"slices"
	"strings"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// StatusFact is one fact in a status bar group: something to press that
// opens what it names.
type StatusFact struct {
	Text string
	// Symbol is an optional meaning name drawn before the text.
	Symbol string
	// Explanation is the one sentence the words do not say: what state the
	// thing is in, what pressing it opens.
	Explanation string
}

// StatusGroup is a named group of facts, such as "Waiting on you".
type StatusGroup struct {
	// Label names what the group is about; "" for a group that needs none.
	Label string
	Facts []StatusFact
}

// StatusBar is the strip along the bottom of a window: on the left what the
// application is doing now, which changes, and on the right the standing
// facts about what is open, which do not. Mixing the two is what makes a
// status bar unreadable.
//
// The right takes facts in three shapes: plain strings separated by dots,
// with nothing to press; named groups of facts, each fact a link that opens
// what it names; and whole controls at the end, for an action that is not a
// fact at all. A bar is one line in a window of any width, so its groups
// will eventually not fit. It shows the groups it has room for, left to
// right, and keeps the rest behind a link saying how many facts are not on
// the bar, which opens a menu of them: nothing is dropped, and nothing is
// silent about it.
type StatusBar struct {
	unison.Panel
	ui *UI
	// Activity says what is happening now; "" leaves the left blank, which is
	// the resting state. A bar that always has something to say trains the
	// reader to stop looking at it.
	Activity string
	// Facts are the standing facts, at the right, separated by dots.
	Facts []string
	// Groups are the facts that do something.
	Groups []StatusGroup
	// GroupsFirst puts the groups at the left, before the activity, for a bar
	// whose left end is the list of what is waiting rather than a sentence
	// about what is running. The plain facts and the controls stay at the
	// right either way.
	GroupsFirst bool
	// OnFact runs when a fact is pressed, on the bar or in the menu, with the
	// index of its group and its index in the group.
	OnFact func(group, fact int)

	activity *Label
	facts    []*Label // each fact after the first follows a dot label
	groups   []*unison.Panel
	overflow *Link
	controls *unison.Panel
	shown    int // how many groups are on the bar

	built       bool
	builtFirst  bool
	builtFacts  []string
	builtGroups []StatusGroup
}

// NewStatusBar returns a bar with controls at its right end.
func NewStatusBar(ui *UI, controls ...unison.Paneler) *StatusBar {
	s := &StatusBar{ui: ui}
	s.Self = s
	s.activity = NewLabel(ui, "")
	s.activity.Role, s.activity.Ink = RoleCaption, InkTextMuted
	s.overflow = NewLink(ui, "")
	s.overflow.Role, s.overflow.Symbol = RoleCaption, "more"
	s.overflow.OnActivate = s.openOverflow
	s.controls = Row(ui, SizeSpaceTight, controls...)
	s.SetLayout(syncing{Layout: statusLayout{s}, sync: s.sync})
	s.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		t := ui.Theme.Tokens()
		p := painterFor(gc, ui)
		r := s.ContentRect(true)
		p.fill(r, t.FooterBackground)
		p.fill(geom.NewRect(r.X, r.Y, r.Width, float32(ui.Interface.Hairline())), t.Border)
	}
	s.sync()
	return s
}

// sync rebuilds the fact and group panels when the caller has changed them.
func (s *StatusBar) sync() {
	s.activity.Text = s.Activity
	changed := !s.built || s.builtFirst != s.GroupsFirst
	if !s.built || !slices.Equal(s.builtFacts, s.Facts) {
		s.builtFacts = slices.Clone(s.Facts)
		s.facts = nil
		for i, f := range s.Facts {
			if i > 0 {
				dot := NewLabel(s.ui, "·")
				dot.Role, dot.Ink = RoleCaption, InkTextFaint
				s.facts = append(s.facts, dot)
			}
			l := NewLabel(s.ui, f)
			l.Role, l.Ink, l.Tabular = RoleCaption, InkTextMuted, true
			s.facts = append(s.facts, l)
		}
		changed = true
	}
	if !s.built || !groupsEqual(s.builtGroups, s.Groups) {
		s.builtGroups = cloneGroups(s.Groups)
		s.groups = nil
		for gi, g := range s.Groups {
			var parts []unison.Paneler
			if g.Label != "" {
				l := NewLabel(s.ui, g.Label)
				l.Role, l.Ink = RoleCaption, InkTextFaint
				parts = append(parts, l)
			}
			for fi, f := range g.Facts {
				link := NewLink(s.ui, f.Text)
				link.Role, link.Symbol, link.Explanation = RoleCaption, f.Symbol, f.Explanation
				link.OnActivate = func() { s.activate(gi, fi) }
				parts = append(parts, link)
			}
			s.groups = append(s.groups, Row(s.ui, SizeSpaceNear, parts...))
		}
		changed = true
	}
	if changed {
		s.built, s.builtFirst = true, s.GroupsFirst
		s.rebuild()
	}
}

func (s *StatusBar) activate(group, fact int) {
	if s.OnFact != nil {
		s.OnFact(group, fact)
	}
}

// rebuild puts the parts back as children in the order they are read.
func (s *StatusBar) rebuild() {
	s.RemoveAllChildren()
	var order []unison.Paneler
	groups := func() {
		for _, g := range s.groups {
			order = append(order, g)
		}
		order = append(order, s.overflow)
	}
	if s.GroupsFirst {
		groups()
		order = append(order, s.activity)
	} else {
		order = append(order, s.activity)
		groups()
	}
	for _, f := range s.facts {
		order = append(order, f)
	}
	order = append(order, s.controls)
	for _, p := range order {
		s.AddChild(p)
	}
}

// hidden are the facts not on the bar, with the group and fact each came
// from.
func (s *StatusBar) hidden() (items []MenuItem, groups []string) {
	for gi := s.shown; gi < len(s.Groups); gi++ {
		g := s.Groups[gi]
		if g.Label != "" {
			groups = append(groups, g.Label)
		}
		for fi, f := range g.Facts {
			label := f.Text
			// The group's name goes with the fact: "2 behind" out of the bar
			// has lost which group it was in.
			if g.Label != "" {
				label = g.Label + " — " + f.Text
			}
			items = append(items, MenuItem{Text: PlainMenuText(label), OnSelect: func() { s.activate(gi, fi) }})
		}
	}
	return items, groups
}

func (s *StatusBar) openOverflow() {
	items, _ := s.hidden()
	if len(items) > 0 {
		s.ui.ShowMenu(s.overflow, "", items)
	}
}

// Shown is how many groups are on the bar; the rest are behind the overflow
// link.
func (s *StatusBar) Shown() int { return s.shown }

// statusLayout lays the bar out as a row, a space apart and a space in from
// each end, with the activity (or, when there is none, empty room) taking
// what is left. It first works out how many groups fit.
type statusLayout struct{ s *StatusBar }

func (l statusLayout) pref(p unison.Paneler) geom.Size {
	_, size, _ := p.AsPanel().Sizes(geom.Size{})
	return size
}

// groupWidth is a group's natural width with its leading space.
func (l statusLayout) groupWidth(i int) float32 {
	w := l.pref(l.s.groups[i]).Width
	if i > 0 {
		w += float32(l.s.ui.Interface.Space())
	}
	return w
}

func (l statusLayout) factsWidth() float32 {
	var w float32
	for i, f := range l.s.facts {
		if i > 0 {
			w += float32(l.s.ui.Interface.Space())
		}
		w += l.pref(f).Width
	}
	return w
}

// fit decides how many groups the bar shows at a width: the ones that fit,
// in order and never with a gap, in what is left after the margins, the
// activity's floor, the facts and the controls. It is tried twice, because
// the link saying what is hidden only exists when something is, and takes
// room of its own once it does.
func (l statusLayout) fit(width float32) {
	s, m := l.s, l.s.ui.Interface
	space := float32(m.Space())
	taken := 2 * space
	if s.Activity != "" {
		// About eight words at the default size: below this the activity
		// is cut short rather than pushing a fact off the bar.
		taken += float32(m.Px(96)) + space
	}
	if len(s.facts) > 0 {
		taken += l.factsWidth() + space
	}
	if cw := l.pref(s.controls).Width; cw > 0 {
		taken += cw + space
	}
	fitting := func(room float32) int {
		used := float32(0)
		for i := range s.groups {
			if used+l.groupWidth(i) > room {
				return i
			}
			used += l.groupWidth(i)
		}
		return len(s.groups)
	}
	s.shown = fitting(width - taken)
	if s.shown < len(s.groups) {
		s.overflow.Hidden = false
		s.updateOverflow()
		s.shown = fitting(width - taken - l.overflowWidth())
	}
	s.updateOverflow()
	for i, g := range s.groups {
		g.Hidden = i >= s.shown
	}
}

func (s *StatusBar) updateOverflow() {
	items, groups := s.hidden()
	s.overflow.Hidden = len(items) == 0
	s.overflow.Text = s.ui.Number(len(items)) + " more"
	// Which groups they came from, for a reader who hears the control rather
	// than sees the bar.
	s.overflow.Accessibility.Description = strings.Join(groups, ", ")
}

func (l statusLayout) overflowWidth() float32 {
	w := l.pref(l.s.overflow).Width
	if l.s.shown > 0 {
		w += float32(l.s.ui.Interface.Space())
	}
	return w
}

func (l statusLayout) height() float32 {
	s, m := l.s, l.s.ui.Interface
	h := float32(m.StatusBarHeight())
	tight := float32(m.SpaceTight())
	h = max(h, l.pref(s.controls).Height+tight)
	for _, g := range s.groups {
		h = max(h, l.pref(g).Height+tight)
	}
	return h
}

func (l statusLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	h := l.height()
	w := hint.Width
	if w <= 0 {
		w = float32(l.s.ui.Interface.Px(600))
	}
	return geom.NewSize(0, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
}

// slot is one item across the bar: a panel, or none for the empty room, its
// width, and whether it takes what is left.
type slot struct {
	panels []*unison.Panel
	width  float32
	fill   bool
}

func (l statusLayout) PerformLayout(target *unison.Panel) {
	s, m := l.s, l.s.ui.Interface
	r := target.ContentRect(false)
	l.fit(r.Width)
	space := float32(m.Space())
	var groups slot
	if len(s.groups) > 0 {
		for i := 0; i < s.shown; i++ {
			groups.panels = append(groups.panels, s.groups[i].AsPanel())
			groups.width += l.groupWidth(i)
		}
		if !s.overflow.Hidden {
			groups.panels = append(groups.panels, s.overflow.AsPanel())
			groups.width += l.overflowWidth()
		}
	}
	// The empty room is always an item, so the activity is two spaces from
	// the groups at either end, and it is what fills when there is no
	// activity.
	room := slot{fill: s.Activity == ""}
	activity := slot{panels: []*unison.Panel{s.activity.AsPanel()}, fill: true}
	var slots []slot
	if s.GroupsFirst {
		if len(s.groups) > 0 {
			slots = append(slots, groups)
		}
		slots = append(slots, room)
		if s.Activity != "" {
			slots = append(slots, activity)
		}
	} else {
		if s.Activity != "" {
			slots = append(slots, activity)
		}
		slots = append(slots, room)
		if len(s.groups) > 0 {
			slots = append(slots, groups)
		}
	}
	s.activity.Hidden = s.Activity == ""
	if len(s.facts) > 0 {
		var facts slot
		for _, f := range s.facts {
			facts.panels = append(facts.panels, f.AsPanel())
		}
		facts.width = l.factsWidth()
		slots = append(slots, facts)
	}
	// The controls' place is always there, as in , which keeps a space
	// before the right margin whether or not it holds anything.
	slots = append(slots, slot{panels: []*unison.Panel{s.controls}, width: l.pref(s.controls).Width})
	fixed := space * float32(len(slots)-1)
	fills := 0
	for _, sl := range slots {
		if sl.fill {
			fills++
		} else {
			fixed += sl.width
		}
	}
	spare := max(0, r.Width-2*space-fixed)
	x := r.X + space
	for _, sl := range slots {
		w := sl.width
		if sl.fill {
			w = spare / float32(max(1, fills))
		}
		l.place(sl, x, w, r)
		x += w + space
	}
}

// place sets the panels of a slot side by side from x, a space apart, each
// centred on the bar's height; the one panel of a filling slot takes the
// slot's width. The space between groups, and before the overflow link, is
// the leading space each carries in .
func (l statusLayout) place(sl slot, x, width float32, r geom.Rect) {
	space := float32(l.s.ui.Interface.Space())
	for _, p := range sl.panels {
		size := l.pref(p)
		w := size.Width
		if sl.fill {
			w = width
		}
		p.SetFrameRect(geom.NewRect(x, r.Y+(r.Height-size.Height)/2, w, size.Height))
		x += w + space
	}
}

func groupsEqual(a, b []StatusGroup) bool {
	return slices.EqualFunc(a, b, func(x, y StatusGroup) bool {
		return x.Label == y.Label && slices.Equal(x.Facts, y.Facts)
	})
}

func cloneGroups(g []StatusGroup) []StatusGroup {
	out := make([]StatusGroup, len(g))
	for i, x := range g {
		out[i] = StatusGroup{Label: x.Label, Facts: slices.Clone(x.Facts)}
	}
	return out
}
