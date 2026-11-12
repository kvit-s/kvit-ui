package kvitui

import (
	"weak"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
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

	// Height, when set, is the row's height in place of the Form's. It is
	// read on every layout, so a measure that reads the row's state, such as
	// one taller with a second line, is followed.
	Height Measure
	// FitContent makes the row as tall as its tallest shown child at the
	// width offered less the margins' width, plus the margins' height.
	// Height wins over it.
	FitContent bool
	// Margins, when set, replace the row's own margins on all four sides;
	// the chevron's room is still added on the right while the row opens
	// something. With it set, a release counts anywhere in the row's box,
	// margins included.
	Margins unison.Border
	// Place, when set, places the row's children itself, given the box
	// inside the margins, in place of the row's own line. The row then asks
	// for the width it is offered.
	Place func(content geom.Rect)
	// Leading, Title, Subtitle and Trailing are slots for a row that needs
	// no arithmetic of its own. When any is set the row holds those panels
	// and nothing else: Leading at its own size at the left, Title taking
	// the width left with Subtitle under it, and Trailing at its own width
	// at the right, a space apart, each centred on the row's height. A
	// hidden slot takes no room and no gap.
	Leading, Title, Subtitle, Trailing unison.Paneler

	// Bar draws an accent bar this wide down the row's left edge while it is
	// Selected; nil draws none.
	Bar Measure
	// Radius rounds the hover, selection and cursor tints and the ring; nil
	// leaves them square.
	Radius Measure
	// RingAtEdge draws the ring just inside the row's edge, where it is
	// otherwise a ring's width further in.
	RingAtEdge bool
	// RingOnAnyFocus shows the ring however the focus arrived, where it
	// otherwise shows only for focus that came from the keyboard.
	RingOnAnyFocus bool
	// NoHoverTint leaves the row untinted under the pointer.
	NoHoverTint bool
	// NoHoverTintDisabled leaves a disabled row untinted under the pointer,
	// where it is otherwise tinted as an enabled one is.
	NoHoverTintDisabled bool

	// HoverWhole keeps the row hovered while the pointer is on anything
	// inside it. unison tells only the innermost panel that the pointer
	// arrived, so a row full of labels otherwise loses its tint over its own
	// words.
	HoverWhole bool
	// OnHover runs after the hover changed, once the event that changed it
	// is over, and not at all when it changed back within that event.
	OnHover func(over bool)
	// PressOn is when a press activates the row: on the button's release
	// inside the row (the default), or as the button goes down.
	PressOn RowPress
	// AnyButton lets every button do what the left one does.
	AnyButton bool
	// OnDoubleClick runs as the left button goes down for the second time;
	// the release still activates the row.
	OnDoubleClick func()
	// OnPress runs as a button the row accepts goes down on it, after the
	// focus is asked for, with where in the row and which button.
	OnPress func(where geom.Point, button int)
	// OnDrag runs as the pointer moves with the button held, with where it
	// is in the row's coordinates.
	OnDrag func(where geom.Point)
	// OnRelease runs first when the button is let go; true says the release
	// was used, as the end of a drag is, and the row does not activate.
	OnRelease func() (handled bool)

	// OnKey is offered every key that reaches the row before the row's own
	// handling: repeats, keys with modifiers, keys on a row that is not
	// Interactive and keys that come up from a control inside the row among
	// them. True says the key was used.
	OnKey func(key unison.KeyCode, mods mod.Modifiers) bool
	// OnFocus runs when the row gains the focus, with whether a press gave
	// it, which is the opposite of what KeyboardFocus then answers; OnBlur
	// when it loses it.
	OnFocus func(byPointer bool)
	OnBlur  func()

	// Description is the row's accessible description, on a row that shows
	// no chevron; one that does keeps OpensLabel.
	Description string
	// Level is the row's one-based depth in a tree, 0 for a row in a flat
	// list; Expandable and Expanded say whether it can open and has.
	Level                int
	Expandable, Expanded bool

	hoverSaid, hoverPending bool
	watched                 *unison.Window // the window whose pointer the row follows for HoverWhole
}

// RowPress is when a press activates a list row.
type RowPress int

const (
	// PressOnRelease activates the row when the button is let go inside it.
	PressOnRelease RowPress = iota
	// PressOnDown activates the row as the button goes down, and the
	// release does nothing.
	PressOnDown
)

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
		if r.OnKey != nil && r.OnKey(key, mods) {
			return true
		}
		if !r.Interactive || !r.Focused() {
			return false
		}
		return keyDown(key, mods, repeat)
	}
	// A row nobody presses leaves a press to what it sits in, such as a card
	// that opens something.
	r.MouseDownCallback = r.mouseDown
	r.MouseDragCallback = r.mouseDrag
	r.MouseUpCallback = r.mouseUp
	enter, exit := r.MouseEnterCallback, r.MouseExitCallback
	r.MouseEnterCallback = func(where geom.Point, mods mod.Modifiers) bool {
		if r.HoverWhole {
			// The window's pointer watch says whether the row is hovered.
			r.pointerAt(where)
			return false
		}
		enter(where, mods)
		r.settleHover()
		return false
	}
	r.MouseExitCallback = func() bool {
		if r.HoverWhole {
			if !r.KeyboardFocus() {
				r.hideTip()
			}
			return false
		}
		exit()
		r.settleHover()
		return false
	}
	gained, lost := r.GainedFocusCallback, r.LostFocusCallback
	r.GainedFocusCallback = func() {
		gained()
		if r.OnFocus != nil {
			r.OnFocus(!r.keyboardFocus)
		}
	}
	r.LostFocusCallback = func() {
		lost()
		if r.OnBlur != nil {
			r.OnBlur()
		}
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
		if slots := r.slots(); len(slots) > 0 {
			holdOnly(r.AsPanel(), slots)
		}
		flex.Columns = max(1, len(r.Children()))
		r.watchPointer()
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

// mouseDown presses the row: the left button, or any with AnyButton. It
// gives the row the focus, and activates it at once where PressOn says so.
func (r *ListRow) mouseDown(where geom.Point, button, clicks int, _ mod.Modifiers) bool {
	if !r.Interactive || !r.Enabled() || (button != unison.ButtonLeft && !r.AnyButton) {
		return false
	}
	r.pressed = r.PressOn == PressOnRelease
	r.hideTip()
	// The focus arrives after this returns; the marker says where it came
	// from until then. A click also ends any keyboard focus shown.
	r.pointerFocus = true
	r.keyboardFocus = false
	r.RequestFocus()
	r.MarkForRedraw()
	if r.OnPress != nil {
		r.OnPress(where, button)
	}
	if clicks == 2 && button == unison.ButtonLeft && r.OnDoubleClick != nil {
		r.OnDoubleClick()
	}
	if r.PressOn == PressOnDown {
		r.fire()
	}
	return true
}

// releaseBox is the part of the row a release has to be in to activate it.
func (r *ListRow) releaseBox() geom.Rect { return r.ContentRect(r.Margins != nil) }

func (r *ListRow) mouseDrag(where geom.Point, _ int, _ mod.Modifiers) bool {
	if r.PressOn == PressOnRelease {
		if inside := where.In(r.releaseBox()); inside != r.pressed {
			r.pressed = inside
			r.MarkForRedraw()
		}
	}
	if r.OnDrag != nil {
		r.OnDrag(where)
	}
	return true
}

func (r *ListRow) mouseUp(where geom.Point, _ int, _ mod.Modifiers) bool {
	if !r.Interactive {
		return false
	}
	was := r.pressed
	r.pressed = false
	r.pointerFocus = false
	r.MarkForRedraw()
	if r.OnRelease != nil && r.OnRelease() {
		return true
	}
	if was && r.PressOn == PressOnRelease && r.Enabled() && where.In(r.releaseBox()) {
		r.fire()
	}
	return true
}

// watchPointer follows the pointer over the row's whole window while
// HoverWhole is set, since unison tells only the innermost panel under the
// pointer that it has arrived.
func (r *ListRow) watchPointer() {
	w := r.Window()
	if !r.HoverWhole || w == nil || w == r.watched {
		return
	}
	r.watched = w
	// The window keeps the listener, which must not keep the row.
	row := weak.Make(r)
	r.ui.OnPointer(w, r.AsPanel(), func(where geom.Point, inside bool) {
		r := row.Value()
		if r == nil || !r.HoverWhole {
			return
		}
		if over := PointerOn(r.AsPanel(), where, inside); over != r.hovered {
			r.hovered = over
			r.MarkForRedraw()
			r.settleHover()
		}
	})
}

// settleHover tells OnHover that the hover changed, once the event that
// changed it is over. The pointer is reported as gone and back each time it
// crosses from one panel to the next, so the row is left and entered again
// as the pointer goes from its padding to its words; told at once, a tooltip
// that follows the hover would close and open again.
func (r *ListRow) settleHover() {
	if r.OnHover == nil || r.hoverPending {
		return
	}
	r.hoverPending = true
	unison.InvokeTask(func() {
		r.hoverPending = false
		if r.hovered == r.hoverSaid {
			return
		}
		r.hoverSaid = r.hovered
		if r.OnHover != nil {
			r.OnHover(r.hovered)
		}
	})
}

// slots are the row's slots that are set, in the order they are read.
func (r *ListRow) slots() []*unison.Panel {
	var out []*unison.Panel
	for _, s := range []unison.Paneler{r.Leading, r.Title, r.Subtitle, r.Trailing} {
		if s != nil {
			out = append(out, s.AsPanel())
		}
	}
	return out
}

// holdOnly makes a panel's children exactly these, changing nothing when
// they already are.
func holdOnly(p *unison.Panel, want []*unison.Panel) {
	have := p.Children()
	same := len(have) == len(want)
	for i := 0; same && i < len(want); i++ {
		same = have[i] == want[i]
	}
	if same {
		return
	}
	p.RemoveAllChildren()
	for _, c := range want {
		p.AddChild(c)
	}
}

// slotShown is a slot that is set and not hidden, or nil.
func slotShown(s unison.Paneler) *unison.Panel {
	if s == nil || s.AsPanel().Hidden {
		return nil
	}
	return s.AsPanel()
}

// arrangeSlots places the slots in a box, when set is true, and returns the
// width they ask for between them and the height of the tallest part.
func (r *ListRow) arrangeSlots(box geom.Rect, set bool) (width, height float32) {
	gap := float32(r.ui.Interface.Space())
	lead, title, sub, trail := slotShown(r.Leading), slotShown(r.Title), slotShown(r.Subtitle), slotShown(r.Trailing)
	left, right := box.X, box.Right()
	parts := 0
	if lead != nil {
		_, p, _ := lead.Sizes(geom.Size{})
		h := p.Height
		if set {
			h = min(h, box.Height)
			lead.SetFrameRect(geom.NewRect(left, box.Y+(box.Height-h)/2, p.Width, h))
		}
		left += p.Width + gap
		width, height = width+p.Width, max(height, p.Height)
		parts++
	}
	if trail != nil {
		_, p, _ := trail.Sizes(geom.Size{})
		h := p.Height
		if set {
			h = min(h, box.Height)
			trail.SetFrameRect(geom.NewRect(right-p.Width, box.Y+(box.Height-h)/2, p.Width, h))
		}
		right -= p.Width + gap
		width, height = width+p.Width, max(height, p.Height)
		parts++
	}
	if title != nil || sub != nil {
		room := max(0, right-left)
		var words, tall float32
		var heights [2]float32
		for i, c := range []*unison.Panel{title, sub} {
			if c == nil {
				continue
			}
			_, own, _ := c.Sizes(geom.Size{})
			words = max(words, own.Width)
			heights[i] = own.Height
			if box.Width > 0 {
				_, at, _ := c.Sizes(geom.NewSize(room, 0))
				heights[i] = at.Height
			}
			tall += heights[i]
		}
		if set {
			y := box.Y + max(0, (box.Height-tall)/2)
			for i, c := range []*unison.Panel{title, sub} {
				if c == nil {
					continue
				}
				c.SetFrameRect(geom.NewRect(left, y, room, heights[i]))
				y += heights[i]
			}
		}
		width, height = width+words, max(height, tall)
		parts++
	}
	return width + gap*float32(max(0, parts-1)), height
}

// contentHeight is the height of the row's tallest shown child at a width.
func (r *ListRow) contentHeight(width float32) float32 {
	if len(r.slots()) > 0 {
		_, h := r.arrangeSlots(geom.NewRect(0, 0, width, 0), false)
		return h
	}
	var h float32
	for _, c := range r.Children() {
		if c.Hidden {
			continue
		}
		_, p, _ := c.Sizes(geom.NewSize(width, 0))
		h = max(h, p.Height)
	}
	return h
}

// height is the row's height for its form.
func (r *ListRow) height() float32 {
	if r.Height != nil {
		return float32(r.Height.Of(r.ui))
	}
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
	if r.Margins != nil {
		in = r.Margins.Insets()
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
	r := l.r
	in := r.contentInsets()
	h := r.height()
	if r.Height == nil && r.FitContent {
		width := hint.Width
		if width > 0 {
			width = max(0, width-in.Width())
		}
		h = r.contentHeight(width) + in.Height()
	}
	switch {
	case r.Place != nil:
		return geom.NewSize(0, h), geom.NewSize(hint.Width, h), geom.NewSize(unison.DefaultMaxSize, h)
	case len(r.slots()) > 0:
		w, _ := r.arrangeSlots(geom.Rect{}, false)
		return geom.NewSize(in.Width(), h), geom.NewSize(w+in.Width(), h), geom.NewSize(unison.DefaultMaxSize, h)
	}
	minSize, prefSize, _ = l.Layout.LayoutSizes(target, hint)
	return geom.NewSize(minSize.Width, h), geom.NewSize(prefSize.Width, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l listRowLayout) PerformLayout(target *unison.Panel) {
	r := l.r
	switch {
	case r.Place != nil:
		r.Place(target.ContentRect(false))
	case len(r.slots()) > 0:
		r.arrangeSlots(target.ContentRect(false), true)
	default:
		l.Layout.PerformLayout(target)
	}
}

func (r *ListRow) drawGround(gc *unison.Canvas, _ geom.Rect) {
	t, m := r.ui.Theme.Tokens(), r.ui.Interface
	p := painterFor(gc, r.ui)
	b := r.ContentRect(true)
	fill := p.fill
	if r.Radius != nil {
		radius := float32(r.Radius.Of(r.ui))
		fill = func(box geom.Rect, c palette.Color) { p.round(box, radius, c) }
	}
	switch {
	case r.Selected:
		fill(b, t.SelectionTint)
	case r.Current:
		fill(b, t.FocusTint)
	case r.Interactive && r.hovered && !r.NoHoverTint && (r.Enabled() || !r.NoHoverTintDisabled):
		fill(b, t.HoverTint)
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
	if r.Bar != nil && r.Selected {
		bar := geom.NewRect(b.X, b.Y, float32(r.Bar.Of(r.ui)), b.Height)
		painterFor(gc, r.ui).fill(bar, t.Accent)
	}
	// Outlined as well as tinted: a tint alone is a two-percent difference in
	// lightness in the high-contrast theme, and invisible in grey.
	if onIt || (r.RingOnAnyFocus && r.Focused()) {
		w := float32(m.FocusRingWidth())
		inset := w + w/2
		if r.RingAtEdge {
			inset = w / 2
		}
		ring := b.Inset(geom.NewUniformInsets(inset))
		paint := Color(t.FocusRing).Paint(gc, ring, paintstyle.Stroke)
		paint.SetStrokeWidth(w)
		if r.Radius != nil {
			radius := float32(r.Radius.Of(r.ui))
			gc.DrawRoundedRect(ring, geom.NewSize(radius, radius), paint)
		} else {
			gc.DrawRect(ring, paint)
		}
	}
}

// ProvideAccessibility describes the row by what it does: a list item when
// it can be pressed, its line of text when it has a label, and otherwise a
// grouping of what it holds. A row holding a control is a grouping too, label
// or not, since a screen reader is told nothing of what is inside a line of
// text, and the control would be gone from what it can reach.
func (r *ListRow) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		switch {
		case r.Interactive:
			n.Role = role.ListItem
		case r.Label != "" && !holdsControl(r.AsPanel()):
			n.Role = role.Label
		default:
			n.Role = role.Group
		}
	}
	n.Name = r.Label
	if r.opens() {
		n.Description = r.OpensLabel
	} else if r.Description != "" {
		n.Description = r.Description
	}
	if r.Level > 0 {
		n.Level = r.Level
	}
	if r.Expandable {
		n.Expandable, n.Expanded = true, r.Expanded
	}
	n.Selectable = r.Interactive || r.Selected
	n.Selected = r.Selected
	if r.Interactive && r.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// holdsControl reports whether a control sits anywhere inside a panel. A Kvit
// control counts while disabled, when it cannot take the focus, so a row does
// not change what it is as its control comes and goes.
func holdsControl(p *unison.Panel) bool {
	for _, c := range p.Children() {
		if _, kvit := c.Self.(ringed); kvit || c.Focusable() || holdsControl(c) {
			return true
		}
	}
	return false
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
