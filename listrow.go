package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/role"
)

// RowForm is one of a list row's four heights. The height is chosen by what
// the row holds rather than by how many rows a view wants to fit: shrinking
// rows to fit more in is what made kvit-hub's earlier interface read as busy.
type RowForm int

const (
	// RowSlim is one line: a name, what it is, one phrase and one figure.
	RowSlim RowForm = iota
	// RowFull is a name over a description, with chips and figures beside.
	RowFull
	// RowSub is a line expanded under a row, indented under its parent.
	RowSub
	// RowCompact is a disclosure or a group heading.
	RowCompact
)

// ListRow is a row in a list: the ground, the rule under it, and what a row
// that can be pressed says about itself. What sits on it is the caller's;
// SlimRow is the arrangement nearly every list wants.
//
// Hover and keyboard focus are separate marks, because the row under the
// pointer and the row the keyboard is on are different rows. A row that does
// something when pressed says so before it is pressed, with a chevron at its
// trailing edge, and a row that does nothing says nothing: a list of records
// and a pane of field names beside their values are both built from this, and
// when every row tints under the pointer the tint stops meaning anything.
type ListRow struct {
	control
	// Form is the row's height.
	Form RowForm
	// Selected marks a chosen row.
	Selected bool
	// Current marks the row a list's own cursor is on, for a list that moves
	// a cursor with the arrow keys rather than the focus.
	Current bool
	// Rule draws a hairline under the row; true unless turned off. A sub
	// line and a compact row never draw one.
	Rule bool
	// Label is what a screen reader says for the row. A row made of several
	// separate texts reads as a jumble without it.
	Label string
	// Interactive makes the row something to press: it takes the focus, tints
	// under the pointer, draws the chevron, and answers the pointer, Return,
	// Enter, Space and a screen reader's press, all with OnActivate.
	Interactive bool
	// NoOpensMark leaves the chevron off an interactive row that ends in a
	// control of its own for the same purpose.
	NoOpensMark bool
	// OpensLabel says in words what pressing the row opens: the tooltip over
	// the chevron, and the row's accessible description.
	OpensLabel string
	// OnActivate runs when an interactive row is pressed.
	OnActivate func()
}

// NewListRow returns a slim row holding content, laid out left to right inside
// its margins, a space apart and centred on the row's height. A child that
// should take the width left over says so with its layout data.
func NewListRow(ui *UI, content ...unison.Paneler) *ListRow {
	r := &ListRow{Rule: true}
	r.Self = r
	r.initListRow(ui, content...)
	return r
}

func (r *ListRow) initListRow(ui *UI, content ...unison.Paneler) {
	r.initControl(ui, func() {
		if r.Interactive && r.OnActivate != nil {
			r.OnActivate()
		}
	}, unison.KeyReturn, unison.KeyNumPadEnter, unison.KeySpace)
	// A row is usually a container, and a key pressed on a button inside it
	// is the button's. The row answers keys only while it holds the focus
	// itself, so one Space does not both press the button and open the row.
	keyDown := r.KeyDownCallback
	r.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if !r.Interactive || !r.Focused() {
			return false
		}
		return keyDown(key, mods, repeat)
	}
	// A row nobody presses leaves a press to what it sits in, such as a card
	// that opens something.
	down, up := r.MouseDownCallback, r.MouseUpCallback
	r.MouseDownCallback = func(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
		return r.Interactive && down(where, button, clicks, mods)
	}
	r.MouseUpCallback = func(where geom.Point, button int, mods mod.Modifiers) bool {
		return r.Interactive && up(where, button, mods)
	}
	// The row draws its own ring, inside its edges, over what it holds.
	r.noRing = true
	for _, c := range content {
		if c.AsPanel().LayoutData() == nil {
			c.AsPanel().SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
		}
		r.AddChild(c)
	}
	r.SetBorder(listRowMargins{r})
	flex := &spaced{FlexLayout: unison.FlexLayout{VAlign: align.Middle}, ui: ui, gap: SizeSpace, horizontal: true}
	r.SetLayout(syncing{Layout: listRowLayout{r, flex}, sync: func() {
		r.SetFocusable(r.Interactive)
		flex.Columns = max(1, len(r.Children()))
	}})
	r.DrawCallback = r.drawGround
	r.DrawOverCallback = r.drawMarks
	r.tip = func() string {
		if r.opens() {
			return r.OpensLabel
		}
		return ""
	}
	r.tipArea = r.opensStrip
}

func (r *ListRow) opens() bool { return r.Interactive && !r.NoOpensMark }

// height is the row's height for its form.
func (r *ListRow) height() float32 {
	m := r.ui.Interface
	switch r.Form {
	case RowFull:
		return float32(m.RowHeight())
	case RowSub:
		return float32(m.RowHeightSub())
	case RowCompact:
		return float32(m.RowHeightCompact())
	}
	return float32(m.RowHeightSlim())
}

// listRowMargins are the margins the row's content sits inside: a sub line
// is indented under its parent, and an opening row keeps the chevron's room.
type listRowMargins struct{ r *ListRow }

func (b listRowMargins) Insets() geom.Insets { return b.r.contentInsets() }

func (b listRowMargins) Draw(*unison.Canvas, geom.Rect) {}

func (r *ListRow) contentInsets() geom.Insets {
	m := r.ui.Interface
	in := geom.Insets{Left: float32(m.SpaceTight()), Right: float32(m.SpaceTight())}
	if r.Form == RowSub {
		in.Left = float32(m.Px(34))
	}
	if r.opens() {
		in.Right += float32(m.IconSizeSmall() + m.SpaceTight())
	}
	return in
}

// opensStrip is the trailing strip the chevron is centred in; the whole
// strip raises the tooltip, since a small glyph is a poor thing to aim at.
func (r *ListRow) opensStrip() geom.Rect {
	m := r.ui.Interface
	b := r.ContentRect(true)
	w := float32(m.IconSizeSmall() + 2*m.SpaceTight())
	return geom.NewRect(b.Right()-w, b.Y, w, b.Height)
}

// listRowLayout is a row of the content, at the height of the row's form.
type listRowLayout struct {
	r *ListRow
	unison.Layout
}

func (l listRowLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	minSize, prefSize, _ = l.Layout.LayoutSizes(target, hint)
	h := l.r.height()
	return geom.NewSize(minSize.Width, h), geom.NewSize(prefSize.Width, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (r *ListRow) drawGround(gc *unison.Canvas, _ geom.Rect) {
	t, m := r.ui.Theme.Tokens(), r.ui.Interface
	p := painterFor(gc, r.ui)
	b := r.ContentRect(true)
	switch {
	case r.Selected:
		p.fill(b, t.SelectionTint)
	case r.Current:
		p.fill(b, t.FocusTint)
	case r.Interactive && r.hovered:
		p.fill(b, t.HoverTint)
	}
	// The rule belongs to the row rather than to the list, so a group that
	// ends partway down a list still closes.
	if r.Rule && r.Form != RowSub && r.Form != RowCompact {
		h := float32(m.Hairline())
		p.fill(geom.NewRect(b.X, b.Bottom()-h, b.Width, h), t.Border)
	}
}

// drawMarks draws, over the content, the chevron of a row that opens
// something and the ring of the row the keyboard is on.
func (r *ListRow) drawMarks(gc *unison.Canvas, _ geom.Rect) {
	t, m := r.ui.Theme.Tokens(), r.ui.Interface
	b := r.ContentRect(true)
	onIt := r.Current || r.KeyboardFocus()
	if r.opens() {
		ink := t.TextFaint
		if onIt {
			ink = t.TextSecondary
		}
		if g, ok := icons.Glyph("chevron-right"); ok {
			s := float32(m.IconSizeSmall())
			strip := r.opensStrip()
			drawGlyph(gc, r.ui, g, s, ink, geom.NewRect(strip.X+(strip.Width-s)/2, strip.Y+(strip.Height-s)/2, s, s))
		}
	}
	// Outlined as well as tinted: a tint alone is a two-percent difference in
	// lightness in the high-contrast theme, and invisible in grey.
	if onIt {
		w := float32(m.FocusRingWidth())
		ring := b.Inset(geom.NewUniformInsets(w + w/2))
		paint := Color(t.FocusRing).Paint(gc, ring, paintstyle.Stroke)
		paint.SetStrokeWidth(w)
		gc.DrawRect(ring, paint)
	}
}

// ProvideAccessibility describes the row by what it does: a list item when
// it can be pressed, its line of text when it has a label, and otherwise a
// grouping of what it holds.
func (r *ListRow) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		switch {
		case r.Interactive:
			n.Role = role.ListItem
		case r.Label != "":
			n.Role = role.Label
		default:
			n.Role = role.Group
		}
	}
	n.Name = r.Label
	if r.opens() {
		n.Description = r.OpensLabel
	}
	n.Selectable = r.Interactive || r.Selected
	n.Selected = r.Selected
	if r.Interactive && r.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction opens the row for a screen reader, which names
// the row it means rather than tabbing to it.
func (r *ListRow) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !r.Interactive || !r.Enabled() {
		return false
	}
	r.fire()
	return true
}
