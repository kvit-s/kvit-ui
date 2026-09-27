package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/pathop"
	"github.com/richardwilkes/unison/enums/role"
)

// tableBody is the rows: it draws the ones in view and nothing else.
type tableBody struct {
	unison.Panel
	t    *Table
	look focusLook
	// asked is set while the table is giving the rows the focus itself.
	asked   bool
	pointer geom.Point
	pointed bool
	tip     partTip
	card    partTip
	// A press, and the drag that may follow it.
	pressRow    int
	pressClicks int
}

// A cell part's key, for the surface shown for it.
type cellKey struct{ row, column, part int }

func newTableBody(t *Table) *tableBody {
	b := &tableBody{t: t, pressRow: -1}
	b.Self = b
	b.tip = partTip{ui: t.ui, owner: b}
	b.card = partTip{ui: t.ui, owner: b}
	b.SetFocusable(true)
	b.DrawCallback = b.draw
	b.GainedFocusCallback = func() {
		b.look.asked = b.asked
		b.asked = false
		b.look.gained(t.ui)
		b.MarkForRedraw()
		b.disclose()
	}
	b.LostFocusCallback = func() {
		b.MarkForRedraw()
		b.disclose()
	}
	b.KeyDownCallback = b.keyDown
	b.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool {
		b.pointer, b.pointed = where, true
		b.pointerMoved()
		return false
	}
	b.MouseEnterCallback = b.MouseMoveCallback
	b.MouseExitCallback = func() bool {
		b.pointed = false
		b.pointerMoved()
		return false
	}
	b.MouseDownCallback = b.mouseDown
	b.MouseDragCallback = b.mouseDrag
	b.MouseUpCallback = b.mouseUp
	b.MouseWheelCallback = func(_, delta geom.Point, _ mod.Modifiers) bool { return t.ease.wheel(delta) }
	b.UpdateCursorCallback = func(where geom.Point) *unison.Cursor {
		// A row that opens something says so before it is pressed, with the
		// hand a link has everywhere else.
		if row, _ := b.cellAt(where); row >= 0 {
			return unison.PointingCursor()
		}
		return unison.ArrowCursor()
	}
	b.FrameChangeCallback = func() {
		if w := b.Window(); w != nil {
			t.ui.watchWindow(w)
		}
	}
	return b
}

// keyboardHeld reports whether the keyboard is in the rows, which is when
// the row it is on is ringed and the cell it is on discloses what it cuts
// short.
func (b *tableBody) keyboardHeld() bool { return b.Focused() && b.look.keyboard }

// rowTop is where a row starts, scrolled.
func (b *tableBody) rowTop(row int) float32 {
	t := b.t
	return b.ContentRect(false).Y + float32(row)*t.rowHeight() - t.scrollY
}

// cellBox is a shown column's box on a row, scrolled.
func (b *tableBody) cellBox(row int, p columnPlace) geom.Rect {
	r := b.ContentRect(false)
	return geom.NewRect(r.X+p.x-b.t.scrollX, b.rowTop(row), p.width, b.t.rowHeight())
}

// cellAt is the row and shown column under a point, or -1 for the row.
func (b *tableBody) cellAt(where geom.Point) (int, columnPlace) {
	t := b.t
	r := b.ContentRect(false)
	if !where.In(r) {
		return -1, columnPlace{}
	}
	row := int((where.Y - r.Y + t.scrollY) / t.rowHeight())
	if row < 0 || row >= t.RowCount() {
		return -1, columnPlace{}
	}
	x := where.X - r.X + t.scrollX
	for _, p := range t.places() {
		if x >= p.x && x < p.x+p.width {
			return row, p
		}
	}
	return row, columnPlace{column: -1}
}

func (b *tableBody) cellLook(row, column int) cellLook {
	return cellLook{selected: b.t.selected[row], leads: column == b.t.OpensColumn}
}

func (b *tableBody) draw(gc *unison.Canvas, _ geom.Rect) {
	t, ui := b.t, b.t.ui
	tk, m := ui.Theme.Tokens(), ui.Interface
	p := painterFor(gc, ui)
	r := b.ContentRect(false)
	rows := t.RowCount()
	if rows == 0 {
		return
	}
	cols := t.columns()
	places := t.places()
	rowH := t.rowHeight()
	hair, ring := float32(m.Hairline()), float32(m.FocusRingWidth())
	for row := int(t.scrollY / rowH); row < rows; row++ {
		y := b.rowTop(row)
		if y >= r.Bottom() {
			break
		}
		selected, hovered := t.selected[row], row == t.hoveredRow
		for _, pl := range places {
			box := b.cellBox(row, pl)
			if box.Right() <= r.X || box.X >= r.Right() {
				continue
			}
			switch {
			case selected:
				p.fill(box, tk.SelectionTint)
			case row == t.curRow && pl.column == t.curCol:
				p.fill(box, tk.FocusTint)
			case hovered:
				p.fill(box, tk.HoverTint)
			}
			gc.Save()
			gc.ClipRect(box, pathop.Intersect, false)
			t.stamps.draw(gc, box, cols[pl.column].Kind, t.Model.Cell(row, pl.column), b.cellLook(row, pl.column))
			gc.Restore()
			p.fill(geom.NewRect(box.X, box.Bottom()-hair, box.Width, hair), tk.Border)
			if hovered {
				// The row underlined as a link is, which is the second half
				// of saying it opens something; the tint alone says only
				// that the pointer is here.
				p.fill(geom.NewRect(box.X, box.Bottom()-ring, box.Width, ring), tk.Link)
			}
		}
	}
	// The row the keyboard is on, ringed across the whole row: a shape rather
	// than a tint, so it survives a grayscale reading, and a different thing
	// from the selection.
	if b.keyboardHeld() && t.curRow >= 0 && t.curRow < rows {
		box := geom.NewRect(r.X, b.rowTop(t.curRow), r.Width, rowH).Inset(geom.NewUniformInsets(ring / 2))
		paint := Color(tk.FocusRing).Paint(gc, box, paintstyle.Stroke)
		paint.SetStrokeWidth(ring)
		radius := float32(m.RadiusControl())
		gc.DrawRoundedRect(box, geom.NewSize(radius, radius), paint)
	}
}

// pointerMoved follows the pointer over the rows: the row under it is
// hovered, and the part of the cell under it says what it has to say.
func (b *tableBody) pointerMoved() {
	t := b.t
	row := -1
	if b.pointed {
		row, _ = b.cellAt(b.pointer)
	}
	if row != t.hoveredRow {
		t.hoveredRow = row
		b.MarkForRedraw()
		t.opens.MarkForRedraw()
	}
	b.disclose()
}

// disclose shows what the part under the pointer has to say, or else the cell
// under the keyboard cursor: a mark's word, a value cut short, what is behind
// an amount.
func (b *tableBody) disclose() {
	t := b.t
	cols := t.columns()
	if b.pointed {
		if row, p := b.cellAt(b.pointer); row >= 0 && p.column >= 0 {
			if b.discloseCell(row, p, cols[p.column].Kind, true) {
				return
			}
		}
	}
	if b.keyboardHeld() && t.curRow >= 0 {
		if p, ok := t.placeOf(t.curCol); ok && b.discloseCell(t.curRow, p, cols[p.column].Kind, false) {
			return
		}
	}
	b.tip.clear()
	b.card.clear()
}

// discloseCell shows what one cell has to say, and reports whether it had
// anything.
func (b *tableBody) discloseCell(row int, p columnPlace, kind CellKind, pointer bool) bool {
	t := b.t
	v := t.Model.Cell(row, p.column)
	box := b.cellBox(row, p)
	parts := t.stamps.draw(nil, box, kind, v, b.cellLook(row, p.column))
	at := func() geom.Rect { return b.cellBox(row, p) }
	switch {
	case kind == CellMarks && pointer:
		for i, r := range parts.marks {
			if b.pointer.In(r) && v.Marks[i].Label != "" {
				dx, dy := r.X-box.X, r.Y-box.Y
				b.card.clear()
				b.tip.tooltip(cellKey{row, p.column, i}, v.Marks[i].Label, func() geom.Rect {
					c := at()
					return geom.NewRect(c.X+dx, c.Y+dy, r.Width, r.Height)
				})
				return true
			}
		}
	case kind == CellMoney && v.FullText != "":
		b.tip.clear()
		b.card.card(cellKey{row, p.column, -1}, func() unison.Paneler { return moneyCard(t.ui, v) }, at)
		return true
	default:
		if words := disclosed(kind, v, parts); words != "" {
			b.card.clear()
			b.tip.tooltip(cellKey{row, p.column, -1}, words, at)
			return true
		}
	}
	return false
}

// selectOnly makes one row the whole selection.
func (t *Table) selectOnly(row int) {
	clear(t.selected)
	t.selected[row] = true
}

// selectRange makes the rows from the anchor to row the selection.
func (t *Table) selectRange(row int) {
	clear(t.selected)
	from, to := min(t.anchor, row), max(t.anchor, row)
	if t.anchor < 0 {
		from = row
	}
	for r := from; r <= to; r++ {
		t.selected[r] = true
	}
}

func (b *tableBody) mouseDown(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
	t := b.t
	b.look.keyboard = false
	b.RequestFocus()
	b.tip.clear()
	b.card.clear()
	if button != unison.ButtonLeft {
		return false
	}
	row, p := b.cellAt(where)
	if row < 0 || p.column < 0 {
		return true
	}
	cols := t.columns()
	if cols[p.column].Kind == CellCheck {
		v := t.Model.Cell(row, p.column)
		// The box takes the press, so ticking a row does not also open it.
		if parts := t.stamps.draw(nil, b.cellBox(row, p), CellCheck, v, cellLook{}); where.In(parts.check) {
			if t.OnCellToggled != nil {
				t.OnCellToggled(row, p.column, !v.Checked)
			}
			return true
		}
	}
	switch {
	case mods.DiscontiguousSelectionDown():
		t.selected[row] = !t.selected[row]
		if !t.selected[row] {
			delete(t.selected, row)
		}
		t.anchor = row
	case mods.ShiftDown() && t.anchor >= 0:
		t.selectRange(row)
	default:
		t.selectOnly(row)
		t.anchor = row
	}
	t.curRow, t.curCol = row, p.column
	b.pressRow, b.pressClicks = row, clicks
	b.MarkForRedraw()
	return true
}

func (b *tableBody) mouseDrag(where geom.Point, _ int, _ mod.Modifiers) bool {
	t := b.t
	if b.pressRow < 0 {
		return true
	}
	// A drag across the rows selects them: the reader's task in a list this
	// size is usually "these forty".
	r := b.ContentRect(false)
	row := int((where.Y - r.Y + t.scrollY) / t.rowHeight())
	row = max(0, min(t.RowCount()-1, row))
	if row != t.curRow {
		t.selectRange(row)
		t.curRow = row
		t.reveal(row, -1)
		b.MarkForRedraw()
	}
	return true
}

func (b *tableBody) mouseUp(where geom.Point, _ int, _ mod.Modifiers) bool {
	t := b.t
	row := b.pressRow
	b.pressRow = -1
	if at, _ := b.cellAt(where); row < 0 || at != row {
		return true
	}
	// One press opens the record; the double press is kept, because a reader
	// who double-clicks a row means the same thing twice.
	if b.pressClicks >= 2 {
		if t.OnRowActivated != nil {
			t.OnRowActivated(row)
		}
	} else if t.OnRowPressed != nil {
		t.OnRowPressed(row)
	}
	return true
}

func (b *tableBody) keyDown(key unison.KeyCode, mods mod.Modifiers, _ bool) bool {
	t := b.t
	rows := t.RowCount()
	if mods.OSMenuCommandDown() {
		switch key {
		case unison.KeyC:
			if s := t.CopySelection(); s != "" {
				unison.ClipboardSetText(s)
			}
			return true
		case unison.KeyA:
			for r := range rows {
				t.selected[r] = true
			}
			b.MarkForRedraw()
			return true
		}
	}
	if rows == 0 {
		return false
	}
	page := max(1, int(b.ContentRect(false).Height/t.rowHeight())-1)
	row := t.curRow
	switch key {
	case unison.KeyUp:
		row--
	case unison.KeyDown:
		row++
	case unison.KeyPageUp:
		row -= page
	case unison.KeyPageDown:
		row += page
	case unison.KeyHome:
		row = 0
	case unison.KeyEnd:
		row = rows - 1
	case unison.KeyLeft, unison.KeyRight:
		by := 1
		if key == unison.KeyLeft {
			by = -1
		}
		if next := t.neighbour(t.curCol, by); next >= 0 {
			t.curCol = next
		} else if _, ok := t.placeOf(t.curCol); !ok {
			t.curCol = t.step(-1, 1)
		}
		b.look.keyboard = true
		t.reveal(-1, t.curCol)
		b.MarkForRedraw()
		b.disclose()
		return true
	case unison.KeyReturn, unison.KeyNumPadEnter:
		if t.curRow < 0 {
			return false
		}
		if t.OnRowPressed != nil {
			t.OnRowPressed(t.curRow)
		}
		if t.OnRowActivated != nil {
			t.OnRowActivated(t.curRow)
		}
		return true
	case unison.KeySpace:
		box := t.checkColumn()
		if t.curRow < 0 || box < 0 {
			return false
		}
		if t.OnCellToggled != nil {
			t.OnCellToggled(t.curRow, box, !t.Model.Cell(t.curRow, box).Checked)
		}
		return true
	default:
		return false
	}
	if t.curRow < 0 && (key == unison.KeyUp || key == unison.KeyPageUp) {
		row = 0
	}
	row = max(0, min(rows-1, row))
	if _, ok := t.placeOf(t.curCol); !ok {
		t.curCol = t.step(-1, 1)
	}
	t.curRow = row
	if mods.ShiftDown() && t.anchor >= 0 {
		t.selectRange(row)
	} else {
		t.selectOnly(row)
		t.anchor = row
	}
	b.look.keyboard = true
	t.reveal(row, -1)
	b.MarkForRedraw()
	b.disclose()
	return true
}

// ProvideAccessibility describes the rows as a table: each row in view, and
// the keyboard's row, with its cells, the focus reported on the cell the
// keyboard is on. A cell's description is its row's name, since a reader
// moving along a row of twelve columns who is told "42.10" is told nothing
// about which record that belongs to.
func (b *tableBody) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	t := b.t
	n := bl.Node()
	n.Role = role.Table
	n.Name = t.Label
	if n.Name == "" {
		n.Name = "Table"
	}
	n.Multiselectable = true
	rows := t.RowCount()
	n.RowCount = rows
	places := t.places()
	n.ColumnCount = len(places)
	n.Controls = append(n.Controls, bl.IDFor(t.header))
	n.Actions = n.Actions.Without(accessibility.Press)
	if rows == 0 {
		return
	}
	cols := t.columns()
	r := b.ContentRect(false)
	rowH := t.rowHeight()
	first := int(t.scrollY / rowH)
	last := min(rows-1, int((t.scrollY+r.Height)/rowH))
	described := make([]int, 0, last-first+2)
	for row := first; row <= last; row++ {
		described = append(described, row)
	}
	if t.curRow >= 0 && (t.curRow < first || t.curRow > last) {
		described = append(described, t.curRow)
	}
	var focus accessibility.NodeID
	for _, row := range described {
		name := t.rowName(row)
		rowID := bl.AddVirtualChild(rowKey(row), func(v *accessibility.Node) {
			v.Role = role.Row
			v.Name = name
			v.RowIndex = row
			v.Bounds = geom.NewRect(r.X, b.rowTop(row), r.Width, rowH)
			v.Selectable, v.Selected = true, t.selected[row]
			v.Focusable = true
			v.Actions = v.Actions.With(accessibility.Select, accessibility.AddToSelection,
				accessibility.RemoveFromSelection, accessibility.ScrollIntoView, accessibility.Focus)
			if t.OnRowPressed != nil || t.OnRowActivated != nil {
				v.Actions = v.Actions.With(accessibility.Press)
			}
		})
		if row == t.curRow {
			focus = rowID
		}
		for i, p := range places {
			kind := cols[p.column].Kind
			val := t.Model.Cell(row, p.column)
			box := b.cellBox(row, p)
			cellID := bl.AddVirtualChildOf(rowID, cellAxKey{row, p.column}, func(v *accessibility.Node) {
				v.ColumnIndex = i
				v.RowIndex = row
				v.Bounds = box
				v.Focusable = true
				v.Actions = v.Actions.With(accessibility.Focus, accessibility.ScrollIntoView)
				switch kind {
				case CellCheck:
					// A box says it is a box and whether it is ticked; a
					// name on a cell around it would be said twice.
					v.Role = role.CheckBox
					v.Name = "Select this row"
					v.Description = name
					v.HasCheck = true
					v.Checked = check.Off
					if val.Checked {
						v.Checked = check.On
					}
					v.Actions = v.Actions.With(accessibility.Press)
				case CellMarks:
					// The marks name themselves, one each, so a reader is
					// told where one state ends and the next begins.
					v.Role = role.Group
					v.Name = cellName(kind, val)
					v.Description = name
				default:
					v.Role = role.Cell
					v.Name = cellName(kind, val)
					v.Description = name
				}
			})
			if kind == CellMarks {
				parts := t.stamps.draw(nil, box, kind, val, cellLook{})
				for j, mark := range val.Marks {
					if mark.Label == "" {
						continue
					}
					bl.AddVirtualChildOf(cellID, markKey{row, p.column, j}, func(v *accessibility.Node) {
						v.Role, v.Name, v.Bounds = role.Image, mark.Label, parts.marks[j]
					})
				}
			}
			if row == t.curRow && p.column == t.curCol && b.keyboardHeld() {
				focus = cellID
			}
		}
	}
	if focus != 0 {
		bl.FocusChild(focus)
	}
}

// rowKey is a row's key among the table's virtual children.
type rowKey int

// cellAxKey is one cell's key among the table's virtual children.
type cellAxKey struct{ row, column int }

// markKey is one mark's key among the table's virtual children.
type markKey struct{ row, column, mark int }

// PerformAccessibilityAction does for a screen reader what a press, a click
// with a modifier or an arrow key does.
func (b *tableBody) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	t := b.t
	switch key := req.Key.(type) {
	case rowKey:
		row := int(key)
		if row < 0 || row >= t.RowCount() {
			return false
		}
		switch req.Action {
		case accessibility.Select:
			t.selectOnly(row)
			t.anchor = row
		case accessibility.AddToSelection:
			t.selected[row] = true
		case accessibility.RemoveFromSelection:
			delete(t.selected, row)
		case accessibility.Focus:
			t.FocusRow(row)
		case accessibility.ScrollIntoView:
			t.reveal(row, -1)
		case accessibility.Press:
			if t.OnRowPressed != nil {
				t.OnRowPressed(row)
			}
			if t.OnRowActivated != nil {
				t.OnRowActivated(row)
			}
		default:
			return false
		}
	case cellAxKey:
		row := key.row
		if row < 0 || row >= t.RowCount() || key.column < 0 || key.column >= len(t.columns()) {
			return false
		}
		switch req.Action {
		case accessibility.Press:
			if t.columns()[key.column].Kind != CellCheck || t.OnCellToggled == nil {
				return false
			}
			t.OnCellToggled(row, key.column, !t.Model.Cell(row, key.column).Checked)
		case accessibility.Focus:
			t.curCol = key.column
			t.FocusRow(row)
		case accessibility.ScrollIntoView:
			t.reveal(row, key.column)
		default:
			return false
		}
	default:
		return false
	}
	b.MarkForRedraw()
	return true
}
