package kvitui

import (
	"strconv"
	"strings"
	"time"

	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// CellKind is how the cells of a column are drawn. The kind belongs to the
// column rather than to each value, which is what makes a column readable:
// every cell in it is aligned the same way, set in the same face, and says
// the same thing about a missing value. A cell that decided for itself would
// give a column where some rows are right-aligned and some are not,
// depending on whether that row's value happened to look like a number.
type CellKind int

// The kinds.
const (
	// CellText is ordinary text, left aligned, cut short with an ellipsis.
	CellText CellKind = iota
	// CellFigure is a measured value: tabular numerals, right aligned, an em
	// dash where nothing was measured.
	CellFigure
	// CellChip is a small labelled mark in the value's Tone.
	CellChip
	// CellSlug is a monospace identifier, without the ground a slug has
	// elsewhere: a box on every row reads as a column of boxes.
	CellSlug
	// CellDate is a date.
	CellDate
	// CellMoney is an amount already formatted by the caller, right aligned,
	// in the danger colour where its sign says it is money leaving.
	CellMoney
	// CellMarks is every state the row is in at once, one dot each.
	CellMarks
	// CellCheck is a box the reader ticks.
	CellCheck
)

var cellKindNames = [...]string{"Text", "Figure", "Chip", "Slug", "Date", "Money", "Marks", "Check"}

// String is the kind's name: "Text", "Figure", "Chip", "Slug", "Date",
// "Money", "Marks" or "Check".
func (k CellKind) String() string {
	if k < 0 || int(k) >= len(cellKindNames) {
		return cellKindNames[CellText]
	}
	return cellKindNames[k]
}

// CellMark is one state a row is in, drawn as one dot in a Marks cell. Label
// is both the word on hover and the name a screen reader says, so the two
// cannot drift apart. Shape is the second channel beside the hue, and Hollow
// doubles what the shapes can say, since a row can be in more states than
// there are shapes.
type CellMark struct {
	Tone   Tone
	Shape  Shape
	Hollow bool
	Label  string
}

// MarkFor is a mark in a tone with the shape that tone takes: a diamond for
// danger, a square for warning, a circle for the rest.
func MarkFor(tone Tone, label string) CellMark {
	shape := ShapeCircle
	switch tone {
	case ToneDanger:
		shape = ShapeDiamond
	case ToneWarning:
		shape = ShapeSquare
	}
	return CellMark{Tone: tone, Shape: shape, Label: label}
}

// CellValue is what one cell holds. A cell formats nothing: it is given the
// string it draws. Money is where that matters most, since an amount is a
// count of the minor unit of its own currency, and how many minor digits it
// has belongs to the currency — none for the yen, three for the dinar — which
// a cell printing two decimal places would get wrong.
type CellValue struct {
	// Text is the value as it is drawn.
	Text string
	// Date is a Date cell's date where Text is empty; it is drawn as
	// 2006-01-02.
	Date time.Time
	// Unmeasured says there is no value, which a Figure or Money cell draws
	// as an em dash: a balance nobody has computed is not a balance of zero.
	Unmeasured bool
	// Tone is a Chip cell's tone.
	Tone Tone
	// Marks are a Marks cell's states.
	Marks []CellMark
	// Checked is whether a Check cell's box is ticked.
	Checked bool
	// Unit is what a Figure or Money cell draws beside the value: a currency
	// code or a unit of measure, kept apart from the digits so a column
	// lines up down its decimal separator.
	Unit string
	// FullText is the whole value where Text is a shortened one; "" says
	// Text is all there is. For a Money cell it is what is behind the
	// amount, one fact to a line — the native amount, the rate, the date of
	// the rate — shown in a hover card.
	FullText string
}

// shown is the value as it is drawn.
func (v CellValue) shown() string {
	if v.Text == "" && !v.Date.IsZero() {
		return v.Date.Format("2006-01-02")
	}
	return v.Text
}

// negative reports whether an amount is money leaving, read from the sign
// already written into it, after any direction marks a right-to-left locale
// writes before it: the cell has a formatted string rather than a number.
func (v CellValue) negative() bool {
	s := strings.TrimLeft(v.Text, "‎‏؜")
	return strings.HasPrefix(s, "-") || strings.HasPrefix(s, "−")
}

// cellName is what a screen reader says for a cell: the whole value rather
// than the shortened one, since a screen reader that says a cut-short payee
// says a different payee.
func cellName(kind CellKind, v CellValue) string {
	switch kind {
	case CellMarks:
		switch len(v.Marks) {
		case 0:
			return ""
		case 1:
			return "1 mark"
		}
		return strconv.Itoa(len(v.Marks)) + " marks"
	case CellCheck:
		return ""
	}
	if v.Unmeasured {
		return "not measured"
	}
	if v.FullText != "" && kind != CellMoney {
		return v.FullText
	}
	return joinWords(v.shown(), v.Unit)
}

// cellLook is how a cell is drawn besides its value.
type cellLook struct {
	// selected draws the text in the primary colour.
	selected bool
	// leads draws the text in the link colour: the row's name, which the
	// reader presses to open the row.
	leads bool
}

// cellParts are where the parts of a drawn cell are, in the coordinates it
// was drawn in, for the pointer.
type cellParts struct {
	marks     []geom.Rect // each mark's box, wider than its dot
	check     geom.Rect   // the box and the room beside it that takes a press
	truncated bool        // the text was cut short
}

// cellStamps draw cells. Each is a component drawn in place without being
// added to the panel tree — configured for one cell, placed and drawn — which
// is how a table draws a page of cells without a panel per cell, and still
// draws a figure, a chip or a box exactly as the component draws it
// elsewhere.
type cellStamps struct {
	ui     *UI
	figure *Figure
	chip   *Chip
	slug   *Slug
	dot    *Dot
	check  *Check
}

func newCellStamps(ui *UI) *cellStamps {
	s := &cellStamps{ui: ui, figure: NewFigure(ui, "", ""), chip: NewChip(ui, ""), slug: NewSlug(ui, ""),
		dot: NewDot(ui), check: NewCheck(ui, "")}
	s.slug.Ground = false
	return s
}

// stampAt draws a detached panel at a box.
func stampAt(gc *unison.Canvas, p unison.Paneler, box geom.Rect) {
	panel := p.AsPanel()
	panel.SetFrameRect(box)
	panel.ValidateLayout()
	gc.Save()
	gc.Translate(box.Point)
	panel.Draw(gc, geom.Rect{Size: box.Size})
	gc.Restore()
}

// preferred is a panel's preferred size.
func preferred(p unison.Paneler) geom.Size {
	_, pref, _ := p.AsPanel().Sizes(geom.Size{})
	return pref
}

// cellText is the layout of a Text or Date cell's words in a width.
func (s *cellStamps) cellText(kind CellKind, v CellValue, look cellLook, width float32) *text.Layout {
	ui, t := s.ui, s.ui.Theme.Tokens()
	ink := t.TextSecondary
	switch {
	case kind == CellDate:
	case look.leads:
		ink = t.Link
	case look.selected:
		ink = t.TextPrimary
	}
	st := ui.Chrome(ui.Size(RoleBody), text.Regular, ink)
	st.Tabular = kind == CellDate
	return ui.Fonts.Layout([]text.Span{{Text: v.shown(), Style: st}}, text.Options{MaxWidth: width, Elide: width > 0})
}

// draw draws a cell's value in box, a near space in from each side, and
// says where its parts went. A nil canvas only works out where they go.
func (s *cellStamps) draw(gc *unison.Canvas, box geom.Rect, kind CellKind, v CellValue, look cellLook) cellParts {
	m := s.ui.Interface
	near := float32(m.SpaceNear())
	c := geom.NewRect(box.X+near, box.Y, max(0, box.Width-2*near), box.Height)
	var parts cellParts
	middle := func(size geom.Size, x float32) geom.Rect {
		return geom.NewRect(x, c.Y+(c.Height-size.Height)/2, size.Width, size.Height)
	}
	switch kind {
	case CellFigure, CellMoney:
		f := s.figure
		f.Value, f.Unit, f.Measured, f.Role, f.Ink = v.Text, v.Unit, !v.Unmeasured, RoleBody, nil
		if kind == CellMoney && v.negative() {
			f.Ink = InkDanger
		}
		size := preferred(f)
		if gc != nil {
			stampAt(gc, f, middle(size, c.Right()-size.Width))
		}
	case CellChip:
		s.chip.Text, s.chip.Tone = v.shown(), v.Tone
		if gc != nil {
			stampAt(gc, s.chip, middle(preferred(s.chip), c.X))
		}
	case CellSlug:
		s.slug.Text = v.shown()
		if gc != nil {
			size := preferred(s.slug)
			size.Width = min(size.Width, c.Width)
			stampAt(gc, s.slug, middle(size, c.X))
		}
	case CellMarks:
		loose := float32(m.SpaceLoose())
		x := c.X
		for _, mark := range v.Marks {
			item := geom.NewRect(x, c.Y+(c.Height-loose)/2, loose, loose)
			parts.marks = append(parts.marks, item)
			if gc != nil {
				d := s.dot
				d.Ink = func(tk tokens.Tokens) palette.Color { return mark.Tone.color(tk) }
				d.Shape, d.Hollow = mark.Shape, mark.Hollow
				size := preferred(d)
				stampAt(gc, d, geom.NewRect(item.X+(loose-size.Width)/2, item.Y+(loose-size.Height)/2, size.Width, size.Height))
			}
			x += loose + near
		}
	case CellCheck:
		s.check.Checked, s.check.Partial = v.Checked, false
		size := preferred(s.check)
		parts.check = middle(size, c.X)
		if gc != nil {
			stampAt(gc, s.check, parts.check)
		}
	default:
		l := s.cellText(kind, v, look, c.Width)
		parts.truncated = l.Text() != v.shown()
		if gc != nil {
			_, h := l.Size()
			l.Draw(gc, c.X, c.Y+(c.Height-h)/2)
		}
	}
	return parts
}

// disclosed is what a Text or Date cell shows on hover and under the
// keyboard cursor: the full value where it differs from what is drawn, or
// what is drawn where the column cut it short; "" for nothing.
func disclosed(kind CellKind, v CellValue, parts cellParts) string {
	if kind != CellText && kind != CellDate {
		return ""
	}
	if v.FullText != "" && v.FullText != v.shown() {
		return v.FullText
	}
	if parts.truncated {
		return v.shown()
	}
	return ""
}

// moneyCard is the hover card behind an amount: its FullText, one fact to a
// line, the first in the body role and the rest smaller and quieter.
func moneyCard(ui *UI, v CellValue) unison.Paneler {
	var lines []unison.Paneler
	for i, line := range strings.Split(v.FullText, "\n") {
		l := NewLabel(ui, line)
		l.Tabular = true
		if i > 0 {
			l.Role, l.Ink = RoleSmall, InkTextSecondary
		}
		lines = append(lines, l)
	}
	card := NewHoverCard(ui, lines...)
	card.SetLayout(&spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: SizeSpaceSnug})
	return card
}

// Cell is one cell of a table, drawn according to the kind of value its
// column holds, for a caller laying cells out itself; a Table draws its own.
// The value is disclosed where the cell cuts it short, under the pointer and
// under the keyboard cursor, because a reader who never touches the pointer
// has the same column to read. A cell is not a control and takes no focus of
// its own, so whatever owns the cursor says where it is through Current.
type Cell struct {
	unison.Panel
	ui *UI
	// Kind is how the cell is drawn.
	Kind CellKind
	// Value is what it holds.
	Value CellValue
	// Leads draws the text in the link colour: the row's name, which the
	// reader presses to open the row.
	Leads bool
	// Selected draws the text in the primary colour.
	Selected bool
	// Current says the keyboard cursor is on the cell, which discloses a
	// value cut short as the pointer does.
	Current bool
	// OnToggle asks for a Check cell's box to change. The cell does not tick
	// itself: whatever owns the selection does, and sets Value.Checked, so a
	// request that is refused leaves no tick behind.
	OnToggle func(checked bool)

	stamps   *cellStamps
	box      *Check // a Check cell's box, a control of its own
	pointer  geom.Point
	pointed  bool
	tip      partTip
	card     partTip
	wasFocus bool
}

// NewCell returns a cell of a kind holding a value.
func NewCell(ui *UI, kind CellKind, v CellValue) *Cell {
	c := &Cell{ui: ui, Kind: kind, Value: v, stamps: newCellStamps(ui)}
	c.Self = c
	c.tip = partTip{ui: ui, owner: c}
	c.card = partTip{ui: ui, owner: c}
	c.box = NewCheck(ui, "")
	c.box.Name = "Select this row"
	c.box.OnChange = func(wanted bool) {
		c.box.Checked = c.Value.Checked
		if c.OnToggle != nil {
			c.OnToggle(wanted)
		}
	}
	c.AddChild(c.box)
	c.SetLayout(syncing{Layout: cellLayout{c}, sync: c.sync})
	c.DrawCallback = c.draw
	c.MouseEnterCallback = func(where geom.Point, _ mod.Modifiers) bool { c.point(where, true); return false }
	c.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool { c.point(where, true); return false }
	c.MouseExitCallback = func() bool { c.point(geom.Point{}, false); return false }
	return c
}

func (c *Cell) sync() {
	c.box.Hidden = c.Kind != CellCheck
	c.box.Checked = c.Value.Checked
	if c.Current != c.wasFocus {
		c.wasFocus = c.Current
		c.disclose()
	}
}

type cellLayout struct{ c *Cell }

func (l cellLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	c, m := l.c, l.c.ui.Interface
	h := float32(m.RowHeightSlim())
	w := c.natural() + 2*float32(m.SpaceNear())
	return geom.NewSize(0, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l cellLayout) PerformLayout(target *unison.Panel) {
	c := l.c
	if c.Kind == CellCheck {
		c.box.SetFrameRect(c.stamps.draw(nil, target.ContentRect(false), CellCheck, c.Value, c.look()).check)
	}
}

// natural is the width the cell's content wants.
func (c *Cell) natural() float32 {
	s := c.stamps
	switch c.Kind {
	case CellFigure, CellMoney:
		s.figure.Value, s.figure.Unit, s.figure.Measured = c.Value.Text, c.Value.Unit, !c.Value.Unmeasured
		return preferred(s.figure).Width
	case CellChip:
		s.chip.Text = c.Value.shown()
		return preferred(s.chip).Width
	case CellSlug:
		s.slug.Text = c.Value.shown()
		return preferred(s.slug).Width
	case CellMarks:
		m := c.ui.Interface
		n := float32(len(c.Value.Marks))
		return max(0, n*float32(m.SpaceLoose())+(n-1)*float32(m.SpaceNear()))
	case CellCheck:
		return preferred(c.box).Width
	}
	w, _ := s.cellText(c.Kind, c.Value, c.look(), 0).Size()
	return w
}

func (c *Cell) look() cellLook { return cellLook{selected: c.Selected, leads: c.Leads} }

func (c *Cell) draw(gc *unison.Canvas, _ geom.Rect) {
	if c.Kind == CellCheck {
		return // the box draws itself
	}
	c.stamps.draw(gc, c.ContentRect(false), c.Kind, c.Value, c.look())
}

// point follows the pointer over the cell's parts.
func (c *Cell) point(where geom.Point, inside bool) {
	c.pointer, c.pointed = where, inside
	c.disclose()
}

// disclose shows what the part under the pointer, or the whole cell under
// the keyboard cursor, has to say: a mark's word, a value cut short, what is
// behind an amount.
func (c *Cell) disclose() {
	box := c.ContentRect(false)
	parts := c.stamps.draw(nil, box, c.Kind, c.Value, c.look())
	whole := func() geom.Rect { return c.ContentRect(false) }
	on := c.pointed || c.Current
	switch {
	case c.Kind == CellMarks && c.pointed:
		for i, r := range parts.marks {
			if c.pointer.In(r) {
				r := r
				c.tip.tooltip(i, c.Value.Marks[i].Label, func() geom.Rect { return r })
				return
			}
		}
		c.tip.clear()
	case c.Kind == CellMoney && on && c.Value.FullText != "":
		v := c.Value
		c.card.card("money", func() unison.Paneler { return moneyCard(c.ui, v) }, whole)
	case on:
		c.tip.tooltip("text", disclosed(c.Kind, c.Value, parts), whole)
	default:
		c.tip.clear()
		c.card.clear()
	}
}

// ProvideAccessibility reads the cell as its whole value. A Marks cell is a
// group of its marks, each named on its own, so a reader is told where one
// state ends and the next begins; a Check cell leaves its name to the box,
// which says it is a box and whether it is ticked.
func (c *Cell) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	n.Name = cellName(c.Kind, c.Value)
	switch c.Kind {
	case CellCheck:
		n.Role = role.Group
		n.Name = ""
	case CellMarks:
		n.Role = role.Group
		parts := c.stamps.draw(nil, c.ContentRect(false), c.Kind, c.Value, c.look())
		for i, mark := range c.Value.Marks {
			if mark.Label == "" {
				continue
			}
			b.AddVirtualChild(i, func(v *accessibility.Node) {
				v.Role, v.Name, v.Bounds = role.Image, mark.Label, parts.marks[i]
			})
		}
	default:
		if n.Role == role.Auto {
			n.Role = role.Cell
		}
	}
}
