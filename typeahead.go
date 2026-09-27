package kvitui

import (
	"strconv"
	"strings"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/pathop"
	"github.com/richardwilkes/unison/enums/role"
)

// Suggestion is one thing a TypeAhead can offer: a value and the words shown
// for it, which are the value unless set.
type Suggestion struct {
	Value, Label string
}

func (s Suggestion) words() string {
	if s.Label != "" {
		return s.Label
	}
	return s.Value
}

// TypeAhead is a field that offers matches as the reader types, for a list
// too long to read where the reader already knows roughly what they want; a
// Select is for a list they will read. Whether the reader may make something
// new, a value that is in no list yet, is the property that matters, and it
// is a required argument rather than a default: a tag picker should let the
// reader make a tag, while a category picker should not, since a typo would
// make a second category beside the right one.
//
// The list shows only while the field holds the keyboard, so a field filled
// by the application opens nothing nobody asked for. It opens below the
// field, or above it where there is no room below. Up and Down move its
// cursor and Return takes the match; the caret stays in the field throughout.
type TypeAhead struct {
	Field
	// Source is everything that could be matched.
	Source []Suggestion
	// AllowNew offers to make what was typed when nothing matches it exactly.
	AllowNew bool
	// MaximumSuggestions caps the matches shown; 8 unless set.
	MaximumSuggestions int
	// OnChoose runs with the value chosen, or the words typed when a new one
	// is made.
	OnChoose func(value string)

	trailing unison.Paneler
	cursor   int  // the list's cursor
	shut     bool // closed until the reader types again
	setting  bool // the application is setting the text
	list     *typeAheadList
	hide     func()
}

// NewTypeAhead returns an empty type-ahead over source. allowNew says whether
// the reader may make a value that is not in it.
func NewTypeAhead(ui *UI, label string, allowNew bool, source ...Suggestion) *TypeAhead {
	a := &TypeAhead{Source: source, AllowNew: allowNew, MaximumSuggestions: 8}
	a.ui = ui
	a.Self = a
	a.initField(ui, false)
	a.Label = label
	a.list = newTypeAheadList(a)
	near := func() float32 { return float32(ui.Interface.SpaceNear()) }
	a.padRight = func() float32 {
		if a.trailing != nil {
			return preferred(a.trailing).Width
		}
		return near()
	}
	a.place = func(box geom.Rect) {
		if a.trailing != nil {
			size := preferred(a.trailing)
			a.trailing.AsPanel().SetFrameRect(geom.NewRect(box.Right()-size.Width, box.Y+(box.Height-size.Height)/2, size.Width, size.Height))
		}
	}
	a.changed = func() {
		a.cursor, a.list.scrollY = 0, 0
		if !a.setting {
			a.shut = false
		}
		a.update()
	}
	a.describe = func() string {
		switch n := len(a.Matches()); n {
		case 0:
			return ""
		case 1:
			return "1 suggestion"
		default:
			return strconv.Itoa(n) + " suggestions"
		}
	}
	gained, lost := a.edit.GainedFocusCallback, a.edit.LostFocusCallback
	a.edit.GainedFocusCallback = func() { gained(); a.update() }
	a.edit.LostFocusCallback = func() { lost(); a.update() }
	keyDown := a.edit.KeyDownCallback
	a.edit.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if a.keyDown(key) {
			return true
		}
		return keyDown(key, mods, repeat)
	}
	return a
}

// SetTrailing puts one control inside the field at its right-hand end, such
// as the button that acts on the chosen value, and the field holds back the
// room it takes. A button after the field the width of its own label grows
// with the value and pushes whatever follows it along.
func (a *TypeAhead) SetTrailing(p unison.Paneler) {
	if a.trailing != nil {
		a.trailing.AsPanel().RemoveFromParent()
	}
	a.trailing = p
	if p != nil {
		a.AddChildAtIndex(p, 0)
	}
	a.MarkForLayoutAndRedraw()
}

// Matches are the suggestions whose words hold what is typed, ignoring case,
// in the order of the source, at most MaximumSuggestions of them.
func (a *TypeAhead) Matches() []Suggestion {
	needle := strings.ToLower(strings.TrimSpace(a.Text()))
	if needle == "" {
		return nil
	}
	limit := a.MaximumSuggestions
	if limit <= 0 {
		limit = 8
	}
	var found []Suggestion
	for _, s := range a.Source {
		if len(found) >= limit {
			break
		}
		if strings.Contains(strings.ToLower(s.words()), needle) {
			found = append(found, s)
		}
	}
	return found
}

// offersNew reports whether the list ends with a row making what was typed.
func (a *TypeAhead) offersNew() bool { return a.AllowNew && strings.TrimSpace(a.Text()) != "" }

// update shows the list while the field holds the keyboard and there is
// something to offer, and takes it away otherwise.
func (a *TypeAhead) update() {
	want := !a.shut && a.edit.Focused() && (len(a.Matches()) > 0 || a.offersNew())
	switch {
	case want && a.hide == nil:
		if w := a.ui.windowOf(a); w != nil {
			a.hide = w.Show(&Popup{Panel: a.list, Place: a.placeList, Anchor: a})
		}
	case !want && a.hide != nil:
		a.hide()
		a.hide = nil
	}
	a.list.MarkForLayoutAndRedraw()
	if w := a.Window(); w != nil {
		w.Content().MarkForLayoutRecursively()
	}
}

// placeList puts the list under the field, as wide as it, or above it where
// there is no room below: a menu near the bottom of a window opens upward on
// every desktop.
func (a *TypeAhead) placeList(bounds geom.Rect, size geom.Size) geom.Rect {
	f := partIn(a, a.edit.FrameRect())
	gap := float32(a.ui.Interface.Space())
	if f.Bottom()+size.Height+gap <= bounds.Bottom() {
		return geom.NewRect(f.X, f.Bottom(), f.Width, size.Height)
	}
	return geom.NewRect(f.X, f.Y-size.Height, f.Width, size.Height)
}

func (a *TypeAhead) keyDown(key unison.KeyCode) bool {
	n := len(a.Matches())
	switch key {
	case unison.KeyDown:
		if a.hide == nil {
			return false
		}
		a.cursor = min(n-1, a.cursor+1)
		a.list.reveal(a.cursor)
	case unison.KeyUp:
		if a.hide == nil {
			return false
		}
		a.cursor = max(0, a.cursor-1)
		a.list.reveal(a.cursor)
	case unison.KeyReturn, unison.KeyNumPadEnter:
		a.accept()
	default:
		return false
	}
	a.list.MarkForRedraw()
	return true
}

// accept takes the match under the cursor, or makes what was typed where
// that is allowed.
func (a *TypeAhead) accept() {
	matches := a.Matches()
	switch {
	case a.cursor >= 0 && a.cursor < len(matches):
		a.choose(matches[a.cursor])
	case a.offersNew():
		a.chooseNew()
	}
}

func (a *TypeAhead) choose(s Suggestion) {
	a.SetText(s.words())
	if a.OnChoose != nil {
		a.OnChoose(s.Value)
	}
}

func (a *TypeAhead) chooseNew() {
	typed := strings.TrimSpace(a.Text())
	a.shut = true
	a.update()
	if a.OnChoose != nil {
		a.OnChoose(typed)
	}
}

// SetText replaces the text and shuts the list until the reader types: a
// field filled by the application, as a record pane fills its fields, opens
// no list that nobody asked for.
func (a *TypeAhead) SetText(s string) {
	a.setting = true
	a.Field.SetText(s)
	a.setting = false
	a.shut = true
	a.update()
}

// typeAheadList is the list of matches under a type-ahead. It takes no
// focus: the caret has to stay in the field, or the next key goes nowhere.
type typeAheadList struct {
	unison.Panel
	a       *TypeAhead
	hovered int
	scrollY float32 // how far the matches are scrolled
	ease    wheelScroll
}

func newTypeAheadList(a *TypeAhead) *typeAheadList {
	l := &typeAheadList{a: a, hovered: -1}
	l.Self = l
	l.ease = wheelScroll{ui: a.ui,
		at:     func() (x, y float32) { return 0, l.scrollY },
		travel: func() (x, y float32) { return 0, l.travel() },
		moveTo: func(_, y float32) { l.scrollTo(y) }}
	l.MouseWheelCallback = func(_, delta geom.Point, _ mod.Modifiers) bool { return l.ease.wheel(delta) }
	l.SetSizer(l.sizes)
	l.DrawCallback = l.draw
	l.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool { l.hover(l.rowAt(where)); return false }
	l.MouseExitCallback = func() bool { l.hover(-1); return false }
	l.MouseDownCallback = func(where geom.Point, _, _ int, _ mod.Modifiers) bool { return true }
	l.MouseUpCallback = func(where geom.Point, _ int, _ mod.Modifiers) bool {
		i := l.rowAt(where)
		matches := a.Matches()
		switch {
		case i >= 0 && i < len(matches):
			a.choose(matches[i])
		case i == len(matches) && a.offersNew():
			a.chooseNew()
		}
		return true
	}
	return l
}

func (l *typeAheadList) rows() int {
	n := len(l.a.Matches())
	if l.a.offersNew() {
		n++
	}
	return n
}

func (l *typeAheadList) rowHeight() float32 { return float32(l.a.ui.Interface.RowHeightSlim()) }

// view is the height the matches are shown in: all of them up to 200 design
// pixels, scrolling past that.
func (l *typeAheadList) view() float32 {
	return min(float32(len(l.a.Matches()))*l.rowHeight(), float32(l.a.ui.Interface.Px(200)))
}

func (l *typeAheadList) travel() float32 {
	return max(0, float32(len(l.a.Matches()))*l.rowHeight()-l.view())
}

func (l *typeAheadList) scrollTo(y float32) {
	y = max(0, min(l.travel(), y))
	if y != l.scrollY {
		l.scrollY = y
		l.MarkForRedraw()
	}
}

// reveal scrolls the cursor's match into view.
func (l *typeAheadList) reveal(i int) {
	top := float32(i) * l.rowHeight()
	switch {
	case top < l.scrollY:
		l.scrollTo(top)
	case top+l.rowHeight() > l.scrollY+l.view():
		l.scrollTo(top + l.rowHeight() - l.view())
	}
}

func (l *typeAheadList) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.a.ui.Interface
	matches := l.view()
	extra := float32(0)
	if l.a.offersNew() {
		extra = l.rowHeight()
	}
	h := min(matches+extra+float32(m.Space()), float32(m.Px(240)))
	size := geom.NewSize(float32(m.Px(220)), h)
	return size, size, size
}

func (l *typeAheadList) rowAt(where geom.Point) int {
	r := l.ContentRect(false).Inset(geom.NewUniformInsets(float32(l.a.ui.Interface.SpaceSnug())))
	if !where.In(r) {
		return -1
	}
	if where.Y >= r.Y+l.view() {
		if l.a.offersNew() && where.Y < r.Y+l.view()+l.rowHeight() {
			return len(l.a.Matches())
		}
		return -1
	}
	return int((where.Y - r.Y + l.scrollY) / l.rowHeight())
}

func (l *typeAheadList) hover(i int) {
	if i != l.hovered {
		l.hovered = i
		l.MarkForRedraw()
	}
}

func (l *typeAheadList) draw(gc *unison.Canvas, _ geom.Rect) {
	a, ui := l.a, l.a.ui
	t, m := ui.Theme.Tokens(), ui.Interface
	p := painterFor(gc, ui)
	b := l.ContentRect(false)
	drawSurface(gc, ui, b)
	r := b.Inset(geom.NewUniformInsets(float32(m.SpaceSnug())))
	gc.Save()
	gc.ClipRect(r, pathop.Intersect, false)
	rowH := l.rowHeight()
	near := float32(m.SpaceNear())
	matches := a.Matches()
	gc.Save()
	gc.ClipRect(geom.NewRect(r.X, r.Y, r.Width, l.view()), pathop.Intersect, false)
	for i, s := range matches {
		box := geom.NewRect(r.X, r.Y+float32(i)*rowH-l.scrollY, r.Width, rowH)
		switch {
		case i == a.cursor:
			p.fill(box, t.FocusTint)
		case i == l.hovered:
			p.fill(box, t.HoverTint)
		}
		lay := ui.Fonts.Layout([]text.Span{{Text: s.words(), Style: ui.Chrome(ui.Size(RoleBody), text.Regular, t.TextPrimary)}},
			text.Options{MaxWidth: max(1, box.Width-2*near), Elide: true})
		_, h := lay.Size()
		lay.Draw(gc, box.X+near, box.Y+(rowH-h)/2)
	}
	gc.Restore()
	if a.offersNew() {
		// The row says what it will make, in quotes: a row reading "Create"
		// under a list of matches does not say which of them it means.
		box := geom.NewRect(r.X, r.Y+l.view(), r.Width, rowH)
		if l.hovered == len(matches) {
			p.fill(box, t.HoverTint)
		}
		s := float32(m.IconSizeSmall())
		if g, ok := icons.Glyph("plus"); ok {
			drawGlyph(gc, ui, g, s, t.Accent, geom.NewRect(box.X+near, box.Y+(rowH-s)/2, s, s))
		}
		lay := ui.Fonts.Layout([]text.Span{{Text: "Create \"" + strings.TrimSpace(a.Text()) + "\"", Style: ui.Chrome(ui.Size(RoleBody), text.Regular, t.Accent)}},
			text.Options{MaxWidth: max(1, box.Width-3*near-s), Elide: true})
		_, h := lay.Size()
		lay.Draw(gc, box.X+2*near+s, box.Y+(rowH-h)/2)
	}
	gc.Restore()
}

// ProvideAccessibility describes the list of matches.
func (l *typeAheadList) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	n.Role = role.List
	n.Name = "Suggestions"
}
