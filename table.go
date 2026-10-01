package kvitui

import (
	"math"
	"slices"
	"strings"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/pathop"
	"github.com/richardwilkes/unison/enums/role"
)

// TableColumn says what one column of a table model is: its title, its
// width at rest, the kind of cell that draws it, and whether it can order the
// rows. A view's own column order, widths and hidden set are the view's state
// rather than the model's, since two windows onto one model may show
// different columns.
type TableColumn struct {
	// Title names the column in its header.
	Title string
	// Width is the column's width at rest, in design pixels; 120 unless set.
	Width int
	// Kind is how its cells are drawn.
	Kind CellKind
	// Unsortable says the column cannot order the rows.
	Unsortable bool
}

// TableModel is what a Table shows. The models belong to the applications,
// because what is in a row is the application's business; the library ships
// the view. The model does the filtering and the sorting and reports the rows
// it now has, so the view never holds a row: it asks for the cells it draws.
type TableModel interface {
	// Rows is how many rows there are to show, which after a filter is fewer
	// than the model holds.
	Rows() int
	// Columns describes the columns, in the model's order.
	Columns() []TableColumn
	// Cell is what one cell holds.
	Cell(row, column int) CellValue
}

// TableRowNamer is a TableModel with its own sentence for a whole row, which
// a screen reader says as each cell's description. Without one, the row is
// named by what its text-bearing columns hold, in column order, which for a
// ledger is the date, the payee and the amount.
type TableRowNamer interface {
	RowName(row int) string
}

// Table is a dense, configurable table over a model, for a list of any size:
// it draws the cells it shows and asks the model for nothing else, so a
// quarter of a million rows cost what a screenful costs.
//
// The reader owns the columns: they sort, resize, move and hide them from the
// header, by pointer and by keyboard, and every one of those is also in the
// column menu, which is where a reader looks for what a control can do. The
// header is one stop in the tab order rather than one per column, with a
// cursor the arrow keys move, so a twelve-column table does not cost twelve
// presses of Tab to walk past. A row opens its record on one press, and a
// Check column draws a box per row with a box in its header for every row
// shown. The table sorts nothing and filters nothing itself: it says what the
// reader asked for, and the model, which is the application's, does it.
type Table struct {
	unison.Panel
	ui *UI
	// Model is what the table shows.
	Model TableModel
	// Label names the table for a screen reader; "Table" unless set.
	Label string
	// EmptyTitle and EmptyDetail are what the table says when it has no
	// rows, which is a thing to say rather than a blank grid to leave.
	EmptyTitle, EmptyDetail string
	// HiddenColumns are the model columns the reader has hidden.
	HiddenColumns []int
	// Order is the model columns in the order they are shown; nil is the
	// model's order. The reader changes it by moving columns.
	Order []int
	// Widths are the columns the reader has resized, by model column, in
	// design pixels.
	Widths map[int]int
	// SortColumn is the column the rows are ordered by, -1 for none, and
	// SortAscending its direction. The header sets both when the reader
	// sorts, and says so through OnSort.
	SortColumn    int
	SortAscending bool
	// HeaderChecked and HeaderPartial are what the box in a Check column's
	// header shows: every row shown, or some of them. The application owns
	// the answer, since what "every row" means changes with the filter.
	HeaderChecked, HeaderPartial bool
	// RowsOpen draws, in a strip of its own down the right-hand edge, the
	// mark that says a press on a row opens something.
	RowsOpen bool
	// OpensColumn is the column holding the row's name, drawn in the link
	// colour to say the row opens something; -1 for none. It is the
	// alternative to RowsOpen, and the two are not drawn together.
	OpensColumn int
	// OpensLabel says in words what pressing a row opens, on the strip's
	// tooltip.
	OpensLabel string
	// OnRowPressed runs on one press on a row, which is how a ledger row
	// opens its record; OnRowActivated on a double press and on Return.
	OnRowPressed, OnRowActivated func(row int)
	// OnSort runs when the reader asks for a different order. The table has
	// already moved its indicator; a handler that cannot sort puts
	// SortColumn and SortAscending back.
	OnSort func(column int, ascending bool)
	// OnHeaderToggled runs when the box in a Check column's header is
	// pressed, with the state asked for.
	OnHeaderToggled func(checked bool)
	// OnCellToggled runs when one row's box is pressed, with the state asked
	// for. The table ticks nothing itself: the model does, and the table
	// redraws from it, so a request that is refused leaves no tick behind.
	OnCellToggled func(row, column int, checked bool)

	header     *tableHeader
	body       *tableBody
	opens      *tableOpens
	vbar, hbar *unison.ScrollBar
	empty      *EmptyState
	stamps     *cellStamps
	scrollX    float32
	scrollY    float32
	curRow     int // the row the keyboard is on, -1 for none
	curCol     int // the model column the keyboard is on
	anchor     int // where a Shift press extends the selection from
	selected   map[int]bool
	hoveredRow int
	cursor     int // the header's cursor, a model column, -1 for none
	namedRow   int
	namedText  string
	ease       wheelScroll
}

// NewTable returns a table over a model.
func NewTable(ui *UI, model TableModel) *Table {
	t := &Table{ui: ui, Model: model, EmptyTitle: "Nothing here yet", SortColumn: -1, SortAscending: true,
		OpensColumn: -1, curRow: -1, anchor: -1, hoveredRow: -1, cursor: -1, namedRow: -1,
		selected: map[int]bool{}, Widths: map[int]int{}}
	t.Self = t
	t.stamps = newCellStamps(ui)
	t.header = newTableHeader(t)
	t.body = newTableBody(t)
	t.opens = newTableOpens(t)
	t.vbar = NewScrollBar(ui, false)
	t.hbar = NewScrollBar(ui, true)
	t.vbar.ChangedCallback = func() { t.ScrollTo(t.scrollX, t.vbar.Value()) }
	t.hbar.ChangedCallback = func() { t.ScrollTo(t.hbar.Value(), t.scrollY) }
	t.empty = NewEmptyState(ui, "")
	t.empty.Symbol = "list"
	for _, p := range []unison.Paneler{t.header, t.body, t.opens, t.vbar, t.hbar, t.empty} {
		t.AddChild(p)
	}
	t.ease = wheelScroll{ui: ui,
		at:     func() (x, y float32) { return t.scrollX, t.scrollY },
		travel: func() (x, y float32) { return t.travel() },
		moveTo: t.ScrollTo}
	t.Accessibility.Role = role.None
	t.SetLayout(syncing{Layout: tableLayout{t}, sync: t.sync})
	return t
}

// columns are the model's columns.
func (t *Table) columns() []TableColumn {
	if t.Model == nil {
		return nil
	}
	return t.Model.Columns()
}

// RowCount is how many rows the table is showing.
func (t *Table) RowCount() int {
	if t.Model == nil {
		return 0
	}
	return t.Model.Rows()
}

func (t *Table) sync() {
	n := len(t.columns())
	if !isOrder(t.Order, n) {
		t.Order = make([]int, n)
		for i := range t.Order {
			t.Order[i] = i
		}
	}
	rows := t.RowCount()
	if t.curRow >= rows {
		t.curRow = rows - 1
	}
	if t.hoveredRow >= rows {
		t.hoveredRow = -1
	}
	if t.cursor >= n || (t.cursor >= 0 && t.isHidden(t.cursor)) {
		t.cursor = t.step(t.cursor, 1)
	}
	t.empty.Hidden = rows > 0
	t.empty.Title, t.empty.Detail = t.EmptyTitle, t.EmptyDetail
	t.opens.Hidden = !t.RowsOpen
}

// isOrder reports whether order holds each of 0 to n-1 once.
func isOrder(order []int, n int) bool {
	if len(order) != n {
		return false
	}
	seen := make([]bool, n)
	for _, c := range order {
		if c < 0 || c >= n || seen[c] {
			return false
		}
		seen[c] = true
	}
	return true
}

func (t *Table) isHidden(column int) bool { return slices.Contains(t.HiddenColumns, column) }

// width is a column's width now: nothing when hidden, the reader's when
// resized, and the model's otherwise.
func (t *Table) width(column int) float32 {
	if t.isHidden(column) {
		return 0
	}
	m := t.ui.Interface
	if w, ok := t.Widths[column]; ok {
		return float32(m.Px(w))
	}
	cols := t.columns()
	if column >= 0 && column < len(cols) && cols[column].Width > 0 {
		return float32(m.Px(cols[column].Width))
	}
	return float32(m.Px(120))
}

// columnPlace is where a shown column is: its model column, its left edge
// before scrolling, and its width.
type columnPlace struct {
	column   int
	x, width float32
}

// places are the shown columns, left to right.
func (t *Table) places() []columnPlace {
	var out []columnPlace
	x := float32(0)
	for _, c := range t.Order {
		w := t.width(c)
		if w <= 0 {
			continue
		}
		out = append(out, columnPlace{column: c, x: x, width: w})
		x += w
	}
	return out
}

func (t *Table) placeOf(column int) (columnPlace, bool) {
	for _, p := range t.places() {
		if p.column == column {
			return p, true
		}
	}
	return columnPlace{}, false
}

func (t *Table) contentWidth() float32 {
	w := float32(0)
	for _, c := range t.Order {
		w += t.width(c)
	}
	return w
}

func (t *Table) rowHeight() float32 { return float32(t.ui.Interface.RowHeightSlim()) }

// travel is how far the rows can scroll, across and down.
func (t *Table) travel() (x, y float32) {
	b := t.body.ContentRect(false)
	return max(0, t.contentWidth()-b.Width), max(0, float32(t.RowCount())*t.rowHeight()-b.Height)
}

// ScrollTo scrolls the rows so x and y are at the view's top left, kept
// within their ends.
func (t *Table) ScrollTo(x, y float32) {
	tx, ty := t.travel()
	x, y = max(0, min(tx, x)), max(0, min(ty, y))
	if x == t.scrollX && y == t.scrollY {
		return
	}
	t.scrollX, t.scrollY = x, y
	t.syncBars()
	t.body.MarkForRedraw()
	t.header.MarkForRedraw()
	t.opens.MarkForRedraw()
	t.body.pointerMoved()
}

// Position is how far the rows are scrolled, across and down.
func (t *Table) Position() (x, y float32) { return t.scrollX, t.scrollY }

func (t *Table) syncBars() {
	b := t.body.ContentRect(false)
	t.vbar.SetRange(t.scrollY, b.Height, float32(t.RowCount())*t.rowHeight())
	t.hbar.SetRange(t.scrollX, b.Width, t.contentWidth())
}

// step is the next shown column from a model column, in the order shown,
// wrapping at both ends; from -1 with 1 is the first and with -1 the last.
func (t *Table) step(from, by int) int {
	n := len(t.Order)
	if n == 0 {
		return -1
	}
	at := slices.Index(t.Order, from)
	if at < 0 {
		at = -1
		if by < 0 {
			at = n
		}
	}
	for range n {
		at += by
		if at < 0 {
			at = n - 1
		} else if at >= n {
			at = 0
		}
		if t.width(t.Order[at]) > 0 {
			return t.Order[at]
		}
	}
	return from
}

type tableLayout struct{ t *Table }

func (l tableLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.t.ui.Interface
	return geom.NewSize(float32(m.Px(120)), float32(m.RowHeightSlim()*3)), geom.NewSize(float32(m.Px(600)), float32(m.Px(400))),
		geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l tableLayout) PerformLayout(target *unison.Panel) {
	t, m := l.t, l.t.ui.Interface
	b := target.ContentRect(false)
	head, strip := float32(m.RowHeightSlim()), float32(m.SpaceWide())
	opens := float32(0)
	if t.RowsOpen {
		opens = float32(m.IconSizeSmall() + 2*m.SpaceNear())
	}
	t.header.SetFrameRect(geom.NewRect(b.X, b.Y, b.Width, head))
	body := geom.NewRect(b.X, b.Y+head, max(0, b.Width-strip-opens), max(0, b.Height-head-strip))
	t.body.SetFrameRect(body)
	t.opens.SetFrameRect(geom.NewRect(body.Right(), body.Y, opens, body.Height))
	t.vbar.SetFrameRect(geom.NewRect(b.Right()-strip, body.Y, strip, body.Height))
	t.hbar.SetFrameRect(geom.NewRect(b.X, body.Bottom(), max(0, b.Width-strip), strip))
	_, p, _ := t.empty.Sizes(geom.NewSize(b.Width, 0))
	t.empty.SetFrameRect(geom.NewRect(b.X, b.Y+(b.Height-p.Height)/2, b.Width, p.Height))
	x, y := t.scrollX, t.scrollY
	t.scrollX, t.scrollY = -1, -1 // so ScrollTo applies the clamp and the bars
	t.ScrollTo(x, y)
}

// Refresh redraws the table after the model's values changed and its rows
// did not.
func (t *Table) Refresh() {
	t.namedRow = -1
	t.MarkForLayoutAndRedraw()
	t.body.MarkForRedraw()
	t.header.MarkForRedraw()
}

// Reset is what to call after the model's rows changed, by a filter or a
// reload: the selection and the keyboard's row are cleared, since the rows
// they named are gone, and the view is kept within the rows there are now.
func (t *Table) Reset() {
	clear(t.selected)
	t.curRow, t.anchor, t.hoveredRow = -1, -1, -1
	t.body.tip.clear()
	t.body.card.clear()
	t.Refresh()
}

// IsRowSelected reports whether a row is selected.
func (t *Table) IsRowSelected(row int) bool { return t.selected[row] }

// SelectedRows are the selected rows, in order.
func (t *Table) SelectedRows() []int {
	out := make([]int, 0, len(t.selected))
	for r := range t.selected {
		out = append(out, r)
	}
	slices.Sort(out)
	return out
}

// CurrentRow is the row the keyboard is on, or -1.
func (t *Table) CurrentRow() int { return t.curRow }

// FocusRow puts the keyboard on a row: the row is scrolled into view and the
// table takes the focus. A surface raised from a row uses it to hand the
// focus back to the row that raised it when it closes, which is the only way
// the reader keeps their place.
func (t *Table) FocusRow(row int) {
	if row < 0 || row >= t.RowCount() {
		return
	}
	t.curRow, t.anchor = row, row
	if t.curCol < 0 || t.isHidden(t.curCol) {
		t.curCol = t.step(-1, 1)
	}
	t.reveal(row, -1)
	t.body.asked = true
	t.body.RequestFocus()
	t.body.MarkForRedraw()
}

// reveal scrolls a row, and a column when it is not -1, into view.
func (t *Table) reveal(row, column int) {
	b := t.body.ContentRect(false)
	x, y := t.scrollX, t.scrollY
	if row >= 0 {
		top := float32(row) * t.rowHeight()
		switch {
		case top < y:
			y = top
		case top+t.rowHeight() > y+b.Height:
			y = top + t.rowHeight() - b.Height
		}
	}
	if p, ok := t.placeOf(column); ok {
		switch {
		case p.x < x:
			x = p.x
		case p.x+p.width > x+b.Width:
			x = p.x + p.width - b.Width
		}
	}
	t.ScrollTo(x, y)
}

// SortBy orders the rows by a column, or by nothing for -1, and says so.
func (t *Table) SortBy(column int, ascending bool) {
	t.SortColumn, t.SortAscending = column, ascending
	t.header.MarkForRedraw()
	if t.OnSort != nil {
		t.OnSort(column, ascending)
	}
}

// activate is what a press on a column's header, or Return on it, does: a
// Check column's box is pressed; a sortable column sorts, and the column
// already sorted on turns its order around.
func (t *Table) activate(column int) {
	cols := t.columns()
	if column < 0 || column >= len(cols) {
		return
	}
	if cols[column].Kind == CellCheck {
		if t.OnHeaderToggled != nil {
			t.OnHeaderToggled(!t.HeaderChecked)
		}
		return
	}
	if cols[column].Unsortable {
		return
	}
	t.SortBy(column, t.SortColumn != column || !t.SortAscending)
}

// HideColumn hides a column.
func (t *Table) HideColumn(column int) {
	if column < 0 || t.isHidden(column) {
		return
	}
	t.HiddenColumns = append(slices.Clone(t.HiddenColumns), column)
	t.Refresh()
}

// ResizeColumnBy widens a column by some pixels, or narrows it for a
// negative amount, never below 48 design pixels: a column narrowed to
// nothing is a column the reader cannot find again to widen.
func (t *Table) ResizeColumnBy(column int, pixels float32) {
	if column < 0 || t.isHidden(column) {
		return
	}
	m := t.ui.Interface
	w := max(float32(m.Px(48)), t.width(column)+pixels)
	t.Widths[column] = int(math.Round(float64(w) / m.Scale()))
	t.Refresh()
}

// MoveColumnBy moves a column one shown column left for -1 or right for 1.
func (t *Table) MoveColumnBy(column, by int) {
	next := t.neighbour(column, by)
	if next < 0 {
		return
	}
	from, to := slices.Index(t.Order, column), slices.Index(t.Order, next)
	t.Order = slices.Delete(slices.Clone(t.Order), from, from+1)
	t.Order = slices.Insert(t.Order, to, column)
	t.Refresh()
}

// neighbour is the shown column beside a column, left for -1 and right for
// 1, or -1 at the end.
func (t *Table) neighbour(column, by int) int {
	places := t.places()
	for i, p := range places {
		if p.column == column {
			if j := i + by; j >= 0 && j < len(places) {
				return places[j].column
			}
			return -1
		}
	}
	return -1
}

// checkColumn is the column that draws boxes, or -1.
func (t *Table) checkColumn() int {
	for i, c := range t.columns() {
		if c.Kind == CellCheck {
			return i
		}
	}
	return -1
}

// rowName is what a screen reader says about a whole row after each cell.
// One row's sentence is kept, since every cell of the row asks for it.
func (t *Table) rowName(row int) string {
	if n, ok := t.Model.(TableRowNamer); ok {
		return n.RowName(row)
	}
	if row == t.namedRow {
		return t.namedText
	}
	var parts []string
	for c, col := range t.columns() {
		if col.Kind == CellMarks || col.Kind == CellCheck {
			continue
		}
		v := t.Model.Cell(row, c)
		if s := joinWords(v.shown(), v.Unit); s != "" && !v.Unmeasured {
			parts = append(parts, s)
		}
	}
	t.namedRow, t.namedText = row, strings.Join(parts, ", ")
	return t.namedText
}

// CopySelection is the selected rows as tab-separated text, one row to a
// line, in the column order the reader is looking at, each value with its
// unit: a Money cell draws its currency beside the amount rather than inside
// it, and a copy of bare numbers from a workspace holding several currencies
// says nothing. Ctrl+C puts it on the clipboard.
func (t *Table) CopySelection() string {
	var lines []string
	cols := t.columns()
	for _, row := range t.SelectedRows() {
		var cells []string
		for _, p := range t.places() {
			if cols[p.column].Kind == CellCheck {
				continue
			}
			v := t.Model.Cell(row, p.column)
			cells = append(cells, joinWords(v.shown(), v.Unit))
		}
		lines = append(lines, strings.Join(cells, "\t"))
	}
	return strings.Join(lines, "\n")
}

// columnMenu opens the menu of a column: sorting, moving, resizing and
// hiding it, which are the things the columns belong to the reader for.
func (t *Table) columnMenu(column int) {
	cols := t.columns()
	if column < 0 || column >= len(cols) {
		return
	}
	t.cursor = column
	t.header.MarkForRedraw()
	sortable := !cols[column].Unsortable && cols[column].Kind != CellCheck
	m := t.ui.Interface
	items := []MenuItem{
		{Text: "Sort ascending", Disabled: !sortable, OnSelect: func() { t.SortBy(column, true) }},
		{Text: "Sort descending", Disabled: !sortable, OnSelect: func() { t.SortBy(column, false) }},
		{Text: "Clear the sort", Disabled: t.SortColumn < 0, OnSelect: func() { t.SortBy(-1, true) }},
		{Separator: true},
		{Text: "Move this column left", Disabled: t.neighbour(column, -1) < 0, OnSelect: func() { t.MoveColumnBy(column, -1) }},
		{Text: "Move this column right", Disabled: t.neighbour(column, 1) < 0, OnSelect: func() { t.MoveColumnBy(column, 1) }},
		{Text: "Widen this column", OnSelect: func() { t.ResizeColumnBy(column, float32(m.Px(32))) }},
		{Text: "Narrow this column", OnSelect: func() { t.ResizeColumnBy(column, -float32(m.Px(32))) }},
		{Separator: true},
		{Text: "Hide this column", OnSelect: func() { t.HideColumn(column) }},
		{Text: "Show every column", Disabled: len(t.HiddenColumns) == 0, OnSelect: func() {
			t.HiddenColumns = nil
			t.Refresh()
		}},
	}
	box := geom.NewRect(0, 0, float32(m.Px(120)), t.header.FrameRect().Height)
	if p, ok := t.placeOf(column); ok {
		box.X = max(0, min(t.header.FrameRect().Width-box.Width, p.x-t.scrollX))
	}
	t.ui.ShowMenuAt(t.header, box, cols[column].Title, items)
}

// focusLook tells focus given by the keyboard from focus given any other way,
// as a control does, for the parts of a table that draw a cursor.
type focusLook struct {
	keyboard bool // the focus came from the keyboard, or a key has been used since
	asked    bool // set while FocusRow is giving the focus
}

func (f *focusLook) gained(ui *UI) {
	f.keyboard = ui.keyTurn || f.asked
	f.asked = false
}

// tableHeader is the row of column titles: one stop in the tab order, with a
// cursor the arrow keys move between the columns.
type tableHeader struct {
	unison.Panel
	t       *Table
	look    focusLook
	hovered int // the column under the pointer, -1 for none
	over    geom.Point
	button  *IconButton // stamped for the column menu's button
	tip     partTip
	// A press, and the drag that may follow it.
	pressCol   int
	pressX     float32
	moved      bool
	resizing   int // the column being resized, -1 for none
	startWidth float32
}

func newTableHeader(t *Table) *tableHeader {
	h := &tableHeader{t: t, hovered: -1, pressCol: -1, resizing: -1}
	h.Self = h
	h.tip = partTip{ui: t.ui, owner: h}
	h.button = NewIconButton(t.ui, "more-vertical", "Column options")
	h.SetFocusable(true)
	h.DrawCallback = h.draw
	h.GainedFocusCallback = func() {
		h.look.gained(t.ui)
		if t.cursor < 0 {
			t.cursor = t.step(-1, 1)
		}
		h.MarkForRedraw()
	}
	h.LostFocusCallback = func() { h.MarkForRedraw() }
	h.KeyDownCallback = h.keyDown
	h.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool { h.pointerAt(where); return false }
	h.MouseEnterCallback = func(where geom.Point, _ mod.Modifiers) bool { h.pointerAt(where); return false }
	h.MouseExitCallback = func() bool {
		h.hovered = -1
		h.tip.clear()
		h.MarkForRedraw()
		return false
	}
	// A right-click, the Menu key and Shift+F10 open a column's menu under
	// the column: under the pointer for a click, under the cursor for a key.
	setContextOpener(h, func(where geom.Point) bool {
		col := h.columnAt(where)
		if col < 0 {
			return false
		}
		t.columnMenu(col)
		return true
	})
	h.MouseDownCallback = h.mouseDown
	h.MouseDragCallback = h.mouseDrag
	h.MouseUpCallback = h.mouseUp
	h.UpdateCursorCallback = func(where geom.Point) *unison.Cursor {
		if h.edgeAt(where) >= 0 || h.resizing >= 0 {
			return unison.ResizeHorizontalCursor()
		}
		return unison.ArrowCursor()
	}
	h.FrameChangeCallback = func() {
		if w := h.Window(); w != nil {
			t.ui.watchWindow(w)
		}
	}
	return h
}

// columnBox is a shown column's box in the header, scrolled with the rows.
func (h *tableHeader) columnBox(p columnPlace) geom.Rect {
	r := h.ContentRect(false)
	return geom.NewRect(r.X+p.x-h.t.scrollX, r.Y, p.width, r.Height)
}

// columnAt is the model column under a point, or -1.
func (h *tableHeader) columnAt(where geom.Point) int {
	for _, p := range h.t.places() {
		if b := h.columnBox(p); where.X >= b.X && where.X < b.Right() {
			return p.column
		}
	}
	return -1
}

// edgeAt is the column whose right edge is under a point, which a drag
// resizes, or -1.
func (h *tableHeader) edgeAt(where geom.Point) int {
	grab := float32(h.t.ui.Interface.Px(4))
	for _, p := range h.t.places() {
		if right := h.columnBox(p).Right(); where.X >= right-grab && where.X <= right+grab {
			return p.column
		}
	}
	return -1
}

// ContextMenuAnchor is where the Menu key opens a column's menu: in the
// column the cursor is on.
func (h *tableHeader) ContextMenuAnchor() geom.Point {
	if p, ok := h.t.placeOf(h.t.cursor); ok {
		b := h.columnBox(p)
		return geom.NewPoint(b.X+b.Width/2, b.Bottom())
	}
	return geom.Point{}
}

// buttonBox is where a column's menu button is.
func (h *tableHeader) buttonBox(p columnPlace) geom.Rect {
	m := h.t.ui.Interface
	b := h.columnBox(p)
	s := float32(m.ControlHeight())
	return geom.NewRect(b.Right()-float32(m.SpaceTight())-s, b.Y+(b.Height-s)/2, s, s)
}

// checkBox is where a Check column's box is.
func (h *tableHeader) checkBox(p columnPlace) geom.Rect {
	m := h.t.ui.Interface
	b := h.columnBox(p)
	size := preferred(h.t.stamps.check)
	return geom.NewRect(b.X+float32(m.SpaceTight()), b.Y+(b.Height-float32(m.ControlHeight()))/2, size.Width, float32(m.ControlHeight()))
}

// showsButton reports whether a column's menu button is drawn: under the
// pointer and under the keyboard cursor. A button on every header all the
// time turns a twelve-column header into a row of dots.
func (h *tableHeader) showsButton(column int) bool {
	return h.hovered == column || h.cursorOn(column)
}

func (h *tableHeader) cursorOn(column int) bool {
	return h.Focused() && h.look.keyboard && h.t.cursor == column
}

func (h *tableHeader) draw(gc *unison.Canvas, _ geom.Rect) {
	t, ui := h.t, h.t.ui
	tk, m := ui.Theme.Tokens(), ui.Interface
	p := painterFor(gc, ui)
	r := h.ContentRect(false)
	p.fill(r, tk.PanelBackground)
	cols := t.columns()
	for _, pl := range t.places() {
		b := h.columnBox(pl)
		if b.Right() <= r.X || b.X >= r.Right() {
			continue
		}
		gc.Save()
		gc.ClipRect(b, pathop.Intersect, false)
		col := cols[pl.column]
		if col.Kind == CellCheck {
			c := t.stamps.check
			c.Checked, c.Partial = t.HeaderChecked, t.HeaderPartial
			stampAt(gc, c, h.checkBox(pl))
		} else {
			right := b.Right() - float32(m.SpaceTight())
			if h.showsButton(pl.column) {
				bb := h.buttonBox(pl)
				h.button.hovered = h.over.In(bb) && h.hovered == pl.column
				stampAt(gc, h.button, bb)
				right = bb.X - float32(m.SpaceTight())
			}
			sorted := t.SortColumn == pl.column
			if sorted {
				s := float32(m.IconSizeSmall())
				name := "sort-ascending"
				if !t.SortAscending {
					name = "sort-descending"
				}
				if g, ok := icons.Glyph(name); ok {
					drawGlyph(gc, ui, g, s, tk.Accent, geom.NewRect(right-s, b.Y+(b.Height-s)/2, s, s))
				}
				right -= s
			}
			ink := tk.TextMuted
			if sorted {
				// The sorted column is named in the reader's own text colour,
				// a second channel beside the indicator's hue; its shape is
				// the third.
				ink = tk.TextPrimary
			}
			x := b.X + float32(m.SpaceNear())
			l := ui.Fonts.Layout([]text.Span{{Text: col.Title, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, ink)}},
				text.Options{MaxWidth: max(1, right-float32(m.SpaceTight())-x), Elide: true})
			_, lh := l.Size()
			l.Draw(gc, x, b.Y+(b.Height-lh)/2)
		}
		if h.cursorOn(pl.column) {
			w := float32(m.FocusRingWidth())
			ring := b.Inset(geom.NewUniformInsets(w + w/2))
			paint := Color(tk.FocusRing).Paint(gc, ring, paintstyle.Stroke)
			paint.SetStrokeWidth(w)
			radius := float32(m.RadiusControl())
			gc.DrawRoundedRect(ring, geom.NewSize(radius, radius), paint)
		}
		gc.Restore()
	}
	hair := float32(m.Hairline())
	p.fill(geom.NewRect(r.X, r.Bottom()-hair, r.Width, hair), tk.BorderStrong)
}

func (h *tableHeader) pointerAt(where geom.Point) {
	h.over = where
	col := h.columnAt(where)
	if col != h.hovered {
		h.hovered = col
	}
	h.MarkForRedraw()
	t := h.t
	p, ok := t.placeOf(col)
	cols := t.columns()
	switch {
	case !ok:
		h.tip.clear()
	case cols[col].Kind == CellCheck && where.In(h.checkBox(p)):
		h.tip.tooltip([2]int{col, 0}, "Select every row the current filter shows", func() geom.Rect { return h.checkBox(p) })
	case cols[col].Kind != CellCheck && where.In(h.buttonBox(p)):
		h.tip.tooltip([2]int{col, 1}, "Column options", func() geom.Rect { return h.buttonBox(p) })
	default:
		h.tip.clear()
	}
}

func (h *tableHeader) mouseDown(where geom.Point, button, _ int, _ mod.Modifiers) bool {
	t := h.t
	h.tip.clear()
	col := h.columnAt(where)
	if button != unison.ButtonLeft {
		return false
	}
	if e := h.edgeAt(where); e >= 0 {
		h.resizing, h.pressX, h.startWidth = e, where.X, t.width(e)
		return true
	}
	if p, ok := t.placeOf(col); ok && t.columns()[col].Kind != CellCheck && h.showsButton(col) && where.In(h.buttonBox(p)) {
		t.columnMenu(col)
		return true
	}
	h.pressCol, h.pressX, h.moved = col, where.X, false
	return true
}

func (h *tableHeader) mouseDrag(where geom.Point, _ int, _ mod.Modifiers) bool {
	t, m := h.t, h.t.ui.Interface
	if h.resizing >= 0 {
		w := max(float32(m.Px(48)), h.startWidth+where.X-h.pressX)
		t.Widths[h.resizing] = int(math.Round(float64(w) / m.Scale()))
		t.Refresh()
		return true
	}
	if h.pressCol < 0 {
		return true
	}
	if !h.moved && math.Abs(float64(where.X-h.pressX)) < float64(m.Px(8)) {
		return true
	}
	h.moved = true
	// The column follows the pointer one place at a time: it moves once the
	// pointer is past the middle of the column beside it.
	for _, by := range []int{-1, 1} {
		next := t.neighbour(h.pressCol, by)
		if next < 0 {
			continue
		}
		p, _ := t.placeOf(next)
		b := h.columnBox(p)
		if (by < 0 && where.X < b.X+b.Width/2) || (by > 0 && where.X > b.X+b.Width/2) {
			t.MoveColumnBy(h.pressCol, by)
		}
	}
	return true
}

func (h *tableHeader) mouseUp(where geom.Point, _ int, _ mod.Modifiers) bool {
	t := h.t
	if h.resizing >= 0 {
		h.resizing = -1
		return true
	}
	col := h.pressCol
	h.pressCol = -1
	if col >= 0 && !h.moved && h.columnAt(where) == col {
		t.cursor = col
		t.activate(col)
		h.MarkForRedraw()
	}
	return true
}

func (h *tableHeader) keyDown(key unison.KeyCode, mods mod.Modifiers, _ bool) bool {
	t, m := h.t, h.t.ui.Interface
	col := t.cursor
	h.look.keyboard = true
	defer h.MarkForRedraw()
	// Width first, because Ctrl+Left is a resize and a bare Left is a move.
	if mods.OSMenuCommandDown() {
		switch key {
		case unison.KeyLeft:
			t.ResizeColumnBy(col, -float32(m.Px(16)))
			return true
		case unison.KeyRight:
			t.ResizeColumnBy(col, float32(m.Px(16)))
			return true
		}
	}
	// Alt+Down is the chord every desktop uses to open a control's list; the
	// window opens the same menu for the Menu key and Shift+F10.
	if mods.OptionDown() && key == unison.KeyDown {
		t.columnMenu(col)
		return true
	}
	switch key {
	case unison.KeyLeft:
		t.cursor = t.step(col, -1)
	case unison.KeyRight:
		t.cursor = t.step(col, 1)
	case unison.KeyHome:
		t.cursor = t.step(-1, 1)
	case unison.KeyEnd:
		t.cursor = t.step(-1, -1)
	case unison.KeyReturn, unison.KeyNumPadEnter, unison.KeySpace:
		t.activate(col)
	case unison.KeyDown:
		// Out of the header and into the rows, where the reader was heading.
		t.body.asked = true
		t.body.RequestFocus()
	default:
		return false
	}
	if t.cursor >= 0 {
		t.reveal(-1, t.cursor)
	}
	return true
}

// selectShownName is the name of the box in a Check column's header, which
// selects every row shown.
const selectShownName = "Select the rows shown"

// headerName is what a screen reader says about a column's header: its
// title, and the order it is holding the rows in. A Check column's header is
// named by its title, as every other header is, and the box inside it by
// what the box does (selectShownName); a Check column with no title is named
// by its box, since a header the keyboard stops on needs a name.
func (t *Table) headerName(column int) string {
	title := t.columns()[column].Title
	switch {
	case t.columns()[column].Kind == CellCheck && title == "":
		return selectShownName
	case t.SortColumn != column:
		return title
	case t.SortAscending:
		return title + ", sorted ascending"
	}
	return title + ", sorted descending"
}

// ProvideAccessibility describes the header as a row of column headers, with
// the focus on the one the cursor is on. A Check column's header holds its
// box as a node of its own: the header is named by the column's title and
// the box by what it does, and the box, not the header, says whether it is
// ticked.
func (h *tableHeader) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	t := h.t
	n := b.Node()
	n.Role = role.TableHeader
	n.Name = "Column headers"
	cols := t.columns()
	for i, p := range t.places() {
		col := cols[p.column]
		box := h.columnBox(p)
		id := b.AddVirtualChild(p.column, func(v *accessibility.Node) {
			v.Role = role.ColumnHeader
			v.Name = t.headerName(p.column)
			v.ColumnIndex = i
			v.Bounds = box
			switch {
			case col.Kind == CellCheck:
				v.Description = "Enter selects or clears every row shown; the menu key opens this column's options"
			case col.Unsortable:
				v.Description = "This column does not sort; the menu key opens its options"
			default:
				v.Description = "Enter sorts this column; the menu key opens its options"
			}
			v.Focusable = true
			v.Actions = v.Actions.With(accessibility.Press, accessibility.ShowContextMenu, accessibility.Focus)
		})
		if col.Kind == CellCheck {
			checkBox := h.checkBox(p)
			b.AddVirtualChildOf(id, headerBoxKey(p.column), func(v *accessibility.Node) {
				v.Role = role.CheckBox
				v.Name = selectShownName
				v.Bounds = checkBox
				v.HasCheck = true
				v.Checked = check.Off
				switch {
				case t.HeaderPartial:
					v.Checked = check.Mixed
				case t.HeaderChecked:
					v.Checked = check.On
				}
				v.Actions = v.Actions.With(accessibility.Press)
			})
		}
		if p.column == t.cursor {
			b.FocusChild(id)
		}
	}
}

// headerBoxKey is the key of the box in a Check column's header among the
// header's virtual children, which key the headers themselves by column.
type headerBoxKey int

// PerformAccessibilityAction presses a column's header or the box in it,
// opens a column's menu, or puts the cursor on a column.
func (h *tableHeader) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if box, ok := req.Key.(headerBoxKey); ok {
		if req.Action != accessibility.Press {
			return false
		}
		h.t.activate(int(box))
		h.MarkForRedraw()
		return true
	}
	col, ok := req.Key.(int)
	if !ok || col < 0 || col >= len(h.t.columns()) {
		return false
	}
	switch req.Action {
	case accessibility.Press:
		h.t.activate(col)
	case accessibility.ShowContextMenu:
		h.t.columnMenu(col)
	case accessibility.Focus:
		h.t.cursor = col
		h.RequestFocus()
	default:
		return false
	}
	h.MarkForRedraw()
	return true
}

// tableOpens is the strip that says a press on a row opens something: one
// mark per drawn row, down the right-hand edge where the eye ends however far
// the columns are scrolled, and not a column, so not part of the arrangement
// the reader owns.
type tableOpens struct {
	unison.Panel
	t   *Table
	tip partTip
}

func newTableOpens(t *Table) *tableOpens {
	o := &tableOpens{t: t}
	o.Self = o
	o.tip = partTip{ui: t.ui, owner: o}
	o.DrawCallback = o.draw
	// The word behind the mark, raised by resting anywhere on the strip:
	// one tooltip for the strip rather than one per row.
	o.MouseEnterCallback = func(geom.Point, mod.Modifiers) bool {
		o.tip.tooltip("opens", t.OpensLabel, func() geom.Rect { return o.ContentRect(false) })
		return false
	}
	o.MouseExitCallback = func() bool { o.tip.clear(); return false }
	return o
}

func (o *tableOpens) draw(gc *unison.Canvas, _ geom.Rect) {
	t, ui := o.t, o.t.ui
	tk, m := ui.Theme.Tokens(), ui.Interface
	r := o.ContentRect(false)
	g, ok := icons.Glyph("chevron-right")
	if !ok {
		return
	}
	rowH := t.rowHeight()
	s := float32(m.IconSizeSmall())
	rows := t.RowCount()
	for row := int(t.scrollY / rowH); row < rows; row++ {
		y := r.Y + float32(row)*rowH - t.scrollY
		if y >= r.Bottom() {
			break
		}
		ink := tk.TextFaint
		if row == t.hoveredRow {
			ink = tk.TextSecondary
		}
		drawGlyph(gc, ui, g, s, ink, geom.NewRect(r.X+float32(math.Round(float64(r.Width-s)/2)), y+float32(math.Round(float64(rowH-s)/2)), s, s))
	}
}

// ProvideAccessibility leaves the strip out: whether a row opens is said by
// the row's press action.
func (o *tableOpens) ProvideAccessibility(b *unison.AccessibilityBuilder) { b.Node().Ignored = true }
