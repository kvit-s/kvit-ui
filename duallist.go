package kvitui

import (
	"slices"
	"strconv"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/role"
)

// DualList is two lists with items moving between them: what is available,
// and what is chosen and in what order, which is the shape of choosing which
// columns a table shows. The chosen side is ordered and the available side is
// not, because column order matters and which columns exist does not.
// Everything works from the keyboard: each list is one stop in the tab order
// whose cursor the arrow keys move and whose Return moves the item across,
// and the buttons between them move and reorder, each saying which way with
// an arrow and what it does in its name. A dual list whose only route is
// dragging between the panes is unusable without a mouse.
type DualList struct {
	unison.Panel
	ui *UI
	// Available and Chosen are the two lists; Chosen is in order.
	Available, Chosen []Option
	// AvailableLabel and ChosenLabel head the two lists.
	AvailableLabel, ChosenLabel string
	// OnChange runs after an item moves, with both lists as they now are.
	OnChange func(available, chosen []Option)

	left, right                      *dualSide
	show, hide, earlier, later       *IconButton
	leftHead, rightHead              *Label
	leftColumn, rightColumn, buttons *unison.Panel
}

// NewDualList returns the two lists.
func NewDualList(ui *UI, available, chosen []Option) *DualList {
	d := &DualList{ui: ui, Available: available, Chosen: chosen, AvailableLabel: "Available", ChosenLabel: "Shown, in order"}
	d.Self = d
	d.left = newDualSide(d, false)
	d.right = newDualSide(d, true)
	button := func(symbol, label string, act func()) *IconButton {
		b := NewIconButton(ui, symbol, label)
		b.Form = Ordinary
		b.OnClick = act
		return b
	}
	d.show = button("arrow-right", "Show the selected column", func() { d.moveAcross(false, d.left.cursor) })
	d.hide = button("arrow-left", "Hide the selected column", func() { d.moveAcross(true, d.right.cursor) })
	d.earlier = button("arrow-up", "Move the selected column earlier", func() { d.reorder(-1) })
	d.later = button("arrow-down", "Move the selected column later", func() { d.reorder(1) })
	d.leftHead, d.rightHead = NewLabel(ui, ""), NewLabel(ui, "")
	for _, l := range []*Label{d.leftHead, d.rightHead} {
		l.Role, l.Ink = RoleSmall, InkTextMuted
	}
	column := func(head *Label, side *dualSide) *unison.Panel {
		side.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
		c := Column(ui, SizeSpaceNear, head, side)
		c.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
		return c
	}
	d.leftColumn, d.rightColumn = column(d.leftHead, d.left), column(d.rightHead, d.right)
	d.buttons = Column(ui, SizeSpaceNear, d.show, d.hide, d.earlier, d.later)
	d.buttons.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
	d.AddChild(d.leftColumn)
	d.AddChild(d.buttons)
	d.AddChild(d.rightColumn)
	d.SetLayout(syncing{Layout: &spaced{FlexLayout: unison.FlexLayout{Columns: 3}, ui: ui, gap: SizeColumnGap, horizontal: true}, sync: d.sync})
	return d
}

func (d *DualList) sync() {
	d.leftHead.Text, d.rightHead.Text = d.AvailableLabel, d.ChosenLabel
	d.left.clamp()
	d.right.clamp()
	d.show.SetEnabled(d.left.cursor >= 0 && d.left.cursor < len(d.Available))
	d.hide.SetEnabled(d.right.cursor >= 0 && d.right.cursor < len(d.Chosen))
	d.earlier.SetEnabled(d.right.cursor > 0 && d.right.cursor < len(d.Chosen))
	d.later.SetEnabled(d.right.cursor >= 0 && d.right.cursor < len(d.Chosen)-1)
}

// changed redraws after a move and says so.
func (d *DualList) changed() {
	d.MarkForLayoutAndRedraw()
	d.left.MarkForRedraw()
	d.right.MarkForRedraw()
	if d.OnChange != nil {
		d.OnChange(d.Available, d.Chosen)
	}
}

// moveAcross moves an item from one list to the end of the other.
func (d *DualList) moveAcross(fromChosen bool, i int) {
	from, to := &d.Available, &d.Chosen
	if fromChosen {
		from, to = to, from
	}
	if i < 0 || i >= len(*from) {
		return
	}
	item := (*from)[i]
	*from = slices.Delete(slices.Clone(*from), i, i+1)
	*to = append(slices.Clone(*to), item)
	d.changed()
}

// reorder moves the chosen item under the cursor one place earlier or later.
func (d *DualList) reorder(by int) {
	at := d.right.cursor
	to := at + by
	if at < 0 || at >= len(d.Chosen) || to < 0 || to >= len(d.Chosen) {
		return
	}
	list := slices.Clone(d.Chosen)
	list[at], list[to] = list[to], list[at]
	d.Chosen = list
	d.right.cursor = to
	d.changed()
}

// dualSide is one of the two lists.
type dualSide struct {
	unison.Panel
	d       *DualList
	chosen  bool
	cursor  int
	hovered int
	scrollY float32
	ease    wheelScroll
}

func newDualSide(d *DualList, chosen bool) *dualSide {
	s := &dualSide{d: d, chosen: chosen, hovered: -1}
	s.Self = s
	s.SetFocusable(true)
	s.ease = wheelScroll{ui: d.ui,
		at:     func() (x, y float32) { return 0, s.scrollY },
		travel: func() (x, y float32) { return 0, s.travel() },
		moveTo: func(_, y float32) { s.scrollTo(y) }}
	s.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		m := d.ui.Interface
		return geom.NewSize(float32(m.Px(80)), s.rowHeight()), geom.NewSize(float32(m.Px(200)), float32(m.Px(240))),
			geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
	})
	s.DrawCallback = s.draw
	s.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool { s.hover(s.rowAt(where)); return false }
	s.MouseExitCallback = func() bool { s.hover(-1); return false }
	s.MouseDownCallback = func(where geom.Point, button, _ int, _ mod.Modifiers) bool {
		if button != unison.ButtonLeft {
			return false
		}
		if i := s.rowAt(where); i >= 0 {
			s.cursor = i
			// A press on an available item shows it; a press on a chosen one
			// puts the cursor there for the buttons.
			if !s.chosen {
				d.moveAcross(false, i)
			}
			d.MarkForLayoutAndRedraw()
			s.MarkForRedraw()
		}
		return true
	}
	s.MouseWheelCallback = func(_, delta geom.Point, _ mod.Modifiers) bool { return s.ease.wheel(delta) }
	s.KeyDownCallback = s.keyDown
	s.GainedFocusCallback = s.MarkForRedraw
	s.LostFocusCallback = s.MarkForRedraw
	return s
}

func (s *dualSide) items() []Option {
	if s.chosen {
		return s.d.Chosen
	}
	return s.d.Available
}

func (s *dualSide) rowHeight() float32 { return float32(s.d.ui.Interface.RowHeightSlim()) }

func (s *dualSide) clamp() {
	n := len(s.items())
	if s.cursor >= n {
		s.cursor = n - 1
	}
	if s.cursor < 0 && n > 0 {
		s.cursor = 0
	}
}

func (s *dualSide) travel() float32 {
	return max(0, float32(len(s.items()))*s.rowHeight()-s.ContentRect(false).Height)
}

func (s *dualSide) scrollTo(y float32) {
	y = max(0, min(s.travel(), y))
	if y != s.scrollY {
		s.scrollY = y
		s.MarkForRedraw()
	}
}

func (s *dualSide) rowAt(where geom.Point) int {
	r := s.ContentRect(false)
	if !where.In(r) {
		return -1
	}
	i := int((where.Y - r.Y + s.scrollY) / s.rowHeight())
	if i >= len(s.items()) {
		return -1
	}
	return i
}

func (s *dualSide) hover(i int) {
	if i != s.hovered {
		s.hovered = i
		s.MarkForRedraw()
	}
}

func (s *dualSide) keyDown(key unison.KeyCode, _ mod.Modifiers, _ bool) bool {
	n := len(s.items())
	if n == 0 {
		return false
	}
	switch key {
	case unison.KeyUp:
		s.cursor = max(0, s.cursor-1)
	case unison.KeyDown:
		s.cursor = min(n-1, s.cursor+1)
	case unison.KeyHome:
		s.cursor = 0
	case unison.KeyEnd:
		s.cursor = n - 1
	case unison.KeyReturn, unison.KeyNumPadEnter, unison.KeySpace:
		s.d.moveAcross(s.chosen, s.cursor)
	default:
		return false
	}
	top := float32(s.cursor) * s.rowHeight()
	view := s.ContentRect(false).Height
	switch {
	case top < s.scrollY:
		s.scrollTo(top)
	case top+s.rowHeight() > s.scrollY+view:
		s.scrollTo(top + s.rowHeight() - view)
	}
	s.d.MarkForLayoutAndRedraw()
	s.MarkForRedraw()
	return true
}

func (s *dualSide) draw(gc *unison.Canvas, _ geom.Rect) {
	ui := s.d.ui
	t, m := ui.Theme.Tokens(), ui.Interface
	p := painterFor(gc, ui)
	b := s.ContentRect(false)
	p.fill(b, t.PanelBackground)
	rowH := s.rowHeight()
	near := float32(m.SpaceNear())
	for i, o := range s.items() {
		y := b.Y + float32(i)*rowH - s.scrollY
		if y+rowH < b.Y {
			continue
		}
		if y >= b.Bottom() {
			break
		}
		box := geom.NewRect(b.X, y, b.Width, rowH)
		current := i == s.cursor
		switch {
		case current:
			p.fill(box, t.FocusTint)
		case i == s.hovered:
			p.fill(box, t.HoverTint)
		}
		x := box.X + near
		if s.chosen {
			st := ui.Chrome(ui.Size(RoleCaption), text.Regular, t.TextFaint)
			st.Tabular = true
			l := ui.Fonts.Layout([]text.Span{{Text: strconv.Itoa(i + 1), Style: st}}, text.Options{})
			w, h := l.Size()
			l.Draw(gc, x, y+(rowH-h)/2)
			x += w + near
		}
		l := ui.Fonts.Layout([]text.Span{{Text: o.Label, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, t.TextPrimary)}},
			text.Options{MaxWidth: max(1, box.Right()-near-x), Elide: true})
		_, h := l.Size()
		l.Draw(gc, x, y+(rowH-h)/2)
		// The cursor's row is outlined as well as tinted, as a list row
		// under a list's own cursor is.
		if current {
			w := float32(m.FocusRingWidth())
			ring := box.Inset(geom.NewUniformInsets(w + w/2))
			paint := Color(t.FocusRing).Paint(gc, ring, paintstyle.Stroke)
			paint.SetStrokeWidth(w)
			radius := float32(m.RadiusControl())
			gc.DrawRoundedRect(ring, geom.NewSize(radius, radius), paint)
		}
	}
	// The panel's rule on all four sides.
	hair := float32(m.Hairline())
	for _, r := range []geom.Rect{
		geom.NewRect(b.X, b.Y, b.Width, hair), geom.NewRect(b.X, b.Bottom()-hair, b.Width, hair),
		geom.NewRect(b.X, b.Y, hair, b.Height), geom.NewRect(b.Right()-hair, b.Y, hair, b.Height),
	} {
		p.fill(r, t.Border)
	}
}

// ProvideAccessibility describes a side as a list of its items, the focus
// on the one under the cursor.
func (s *dualSide) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	n.Role = role.List
	n.Name = s.d.AvailableLabel
	if s.chosen {
		n.Name = s.d.ChosenLabel
	}
	r := s.ContentRect(false)
	rowH := s.rowHeight()
	for i, o := range s.items() {
		id := b.AddVirtualChild(i, func(v *accessibility.Node) {
			v.Role = role.ListItem
			v.Name = o.Label
			v.RowIndex = i
			v.Bounds = geom.NewRect(r.X, r.Y+float32(i)*rowH-s.scrollY, r.Width, rowH)
			v.Selectable, v.Selected = true, i == s.cursor
			v.Focusable = true
			v.Actions = v.Actions.With(accessibility.Press, accessibility.Focus)
		})
		if i == s.cursor {
			b.FocusChild(id)
		}
	}
}

// PerformAccessibilityAction moves an item across, or puts the cursor on it,
// for a screen reader.
func (s *dualSide) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	i, ok := req.Key.(int)
	if !ok || i < 0 || i >= len(s.items()) {
		return false
	}
	switch req.Action {
	case accessibility.Press:
		s.d.moveAcross(s.chosen, i)
	case accessibility.Focus:
		s.cursor = i
		s.RequestFocus()
	default:
		return false
	}
	s.d.MarkForLayoutAndRedraw()
	return true
}
