package kvitui_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// counted is a benchmark model that counts the cells it is asked for.
type counted struct {
	*kvitui.BenchmarkTableModel
	asked int
}

func (c *counted) Cell(row, column int) kvitui.CellValue {
	c.asked++
	return c.BenchmarkTableModel.Cell(row, column)
}

// focusedNode is the node the window's tree says holds the focus.
func focusedNode(screen *unison.HeadlessScreen, w *kvitui.Window) *accessibility.Node {
	tree := screen.AccessibilityTree(w.Window)
	if tree == nil {
		return nil
	}
	return tree.Nodes[tree.Focus]
}

// nodes are the window's nodes of a role.
func nodes(screen *unison.HeadlessScreen, w *kvitui.Window, r role.Enum) []*accessibility.Node {
	var out []*accessibility.Node
	for _, n := range screen.AccessibilityTree(w.Window).Nodes {
		if n.Role == r {
			out = append(out, n)
		}
	}
	return out
}

// What "sub-100 ms filtering" means: one filter change over 250,000 rows,
// to a filter nothing matches, which tests every row and keeps none, to one
// a reader would type, and back to none.
func TestFilteringAQuarterMillionRowsStaysUnderTheBudget(t *testing.T) {
	const budget = 100 * time.Millisecond
	rows := kvitui.NewBenchmarkTableModel(kvitui.BenchmarkRows)
	timed := func(filter string) time.Duration {
		start := time.Now()
		rows.SetFilter(filter)
		return time.Since(start)
	}
	none := timed("zzzzz")
	if rows.Rows() != 0 {
		t.Errorf("a filter nothing matches left %d rows", rows.Rows())
	}
	some := timed("Harlow")
	if n := rows.Rows(); n == 0 || n >= kvitui.BenchmarkRows {
		t.Errorf("a filter for one place left %d rows", n)
	}
	matched := rows.Rows()
	all := timed("")
	if rows.Rows() != kvitui.BenchmarkRows {
		t.Errorf("clearing the filter left %d rows", rows.Rows())
	}
	t.Logf("filtering 250,000 rows: %v to nothing, %v to %d matches, %v to clear (budget %v)", none, some, matched, all, budget)
	for _, d := range []time.Duration{none, some, all} {
		if d > budget {
			t.Errorf("a filter change took %v against a budget of %v", d, budget)
		}
	}
}

// Ticks are kept by the generated row, so narrowing the filter and widening
// it again finds them where they were, and the header's box says whether
// every row shown, some, or none is ticked.
func TestTicksSurviveAFilterChange(t *testing.T) {
	rows := kvitui.NewBenchmarkTableModel(1000)
	rows.SetChecked(3, true)
	rows.SetChecked(9, true)
	if !rows.SomeShownChecked() || rows.AllShownChecked() {
		t.Error("two ticks out of a thousand are not some")
	}
	rows.SetFilter("Transfer")
	rows.SetEveryShownChecked(true)
	if !rows.AllShownChecked() {
		t.Error("ticking every row shown did not tick them all")
	}
	rows.SetFilter("")
	if !rows.Cell(3, 12).Checked || !rows.Cell(9, 12).Checked || !rows.SomeShownChecked() || rows.AllShownChecked() {
		t.Error("the ticks did not survive the filter")
	}
}

// The claim a table of this size turns on: what a page costs to draw is what
// is on screen, not what is in the model. Each page of forty is timed on its
// own and the middle one has to fit in a 60 Hz frame; the worst is reported
// beside it rather than asserted on, since a page that waited on a busy
// machine says nothing about the table.
func TestATableDrawsAPageOfAQuarterMillionRowsInsideAFrame(t *testing.T) {
	const budget = 16 * time.Millisecond
	const pages = 40
	var table *kvitui.Table
	model := &counted{BenchmarkTableModel: kvitui.NewBenchmarkTableModel(kvitui.BenchmarkRows)}
	screen, ui, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		table = kvitui.NewTable(ui, model)
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(1200), kvitui.Height(ui, kvitui.Px(700), table))}
	})
	var times []time.Duration
	var asked []int
	screen.Do(func() {
		frame := table.FrameRect()
		page := frame.Height - float32(2*ui.Interface.RowHeightSlim())
		_, err := unison.NewImageFromDrawing(int(frame.Width), int(frame.Height), 72, func(gc *unison.Canvas) {
			for i := range pages {
				// The middle of the list, where a reader who has scrolled is.
				table.ScrollTo(0, float32(kvitui.BenchmarkRows/2)*float32(ui.Interface.RowHeightSlim())+float32(i)*page)
				model.asked = 0
				start := time.Now()
				table.Draw(gc, geom.Rect{Size: frame.Size})
				times = append(times, time.Since(start))
				asked = append(asked, model.asked)
			}
		})
		if err != nil {
			t.Error(err)
		}
	})
	sorted := slices.Clone(times)
	slices.Sort(sorted)
	median := sorted[len(sorted)/2]
	t.Logf("drawing a page of a 250,000-row table, %d pages: median %v, worst %v (budget %v); %d cells asked for a page",
		pages, median, sorted[len(sorted)-1], budget, slices.Max(asked))
	if median > budget {
		t.Errorf("the middle page took %v to draw against a budget of %v", median, budget)
	}
	// About 22 rows of 13 columns are on screen. A table that asked for
	// every row would ask for three million cells.
	if n := slices.Max(asked); n > 30*13 {
		t.Errorf("a page asked the model for %d cells", n)
	}
}

func TestTheHeaderSortsAndResizesFromTheKeyboard(t *testing.T) {
	var table *kvitui.Table
	var sorts []string
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		table = kvitui.NewTable(ui, kvitui.NewBenchmarkTableModel(400))
		table.OnSort = func(column int, ascending bool) {
			sorts = append(sorts, map[bool]string{true: "up", false: "down"}[ascending]+" "+string(rune('0'+column)))
		}
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(300), table))}
	})
	// The header is one stop in the tab order, with a cursor on a column.
	for range 3 {
		if n := focusedNode(screen, w); n != nil && n.Role == role.ColumnHeader {
			break
		}
		screen.KeyPress(unison.KeyTab, 0)
	}
	if n := focusedNode(screen, w); n == nil || n.Role != role.ColumnHeader || n.Name != "Date" {
		t.Fatalf("Tab did not reach the header's first column: %+v", n)
	}
	screen.KeyPress(unison.KeyReturn, 0)
	screen.KeyPress(unison.KeyReturn, 0)
	if !slices.Equal(sorts, []string{"up 0", "down 0"}) {
		t.Errorf("Return twice sorted %v", sorts)
	}
	if n := focusedNode(screen, w); n == nil || n.Name != "Date, sorted descending" {
		t.Errorf("the sorted column's header says %+v", n)
	}
	screen.KeyPress(unison.KeyRight, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Reference" {
		t.Errorf("Right moved the cursor to %+v", n)
	}
	screen.KeyPress(unison.KeyRight, mod.OSMenuCommand())
	screen.Do(func() {
		if got := table.Widths[1]; got != 126 {
			t.Errorf("Ctrl+Right made the column %d design pixels wide, want 126", got)
		}
	})
	// Down leaves the header for the rows.
	screen.KeyPress(unison.KeyDown, 0)
	screen.KeyPress(unison.KeyDown, 0)
	if n := focusedNode(screen, w); n == nil || n.Role != role.Cell {
		t.Errorf("Down from the header reached %+v", n)
	}
}

// titledChecks is the benchmark model with a title on its Check column, as
// an application's ledger has one.
type titledChecks struct{ *kvitui.BenchmarkTableModel }

func (m titledChecks) Columns() []kvitui.TableColumn {
	cols := slices.Clone(m.BenchmarkTableModel.Columns())
	cols[12].Title = "Selected"
	return cols
}

// A Check column's header is heard by the column's title, as every other
// header is, and the box inside it by what it does, with its tick; the box
// presses from a screen reader. An untitled Check column is named by its
// box, since the keyboard stops on its header.
func TestACheckColumnsHeaderIsNamedByItsTitleAndItsBoxByWhatItDoes(t *testing.T) {
	var table *kvitui.Table
	var toggled []bool
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		table = kvitui.NewTable(ui, titledChecks{kvitui.NewBenchmarkTableModel(40)})
		table.HiddenColumns = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		table.HeaderPartial = true
		table.OnHeaderToggled = func(checked bool) {
			toggled = append(toggled, checked)
			table.HeaderChecked, table.HeaderPartial = checked, false
		}
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(300), table))}
	})
	tree := screen.AccessibilityTree(w.Window)
	var header, box *accessibility.Node
	for _, n := range tree.Nodes {
		if n.Role == role.ColumnHeader && n.ColumnIndex == 1 {
			header = n
		}
	}
	if header == nil || header.Name != "Selected" || header.HasCheck {
		t.Fatalf("the Check column's header is %+v", header)
	}
	for _, id := range header.Children {
		if n := tree.Nodes[id]; n.Role == role.CheckBox {
			box = n
		}
	}
	if box == nil || box.Name != "Select the rows shown" || !box.HasCheck || box.Checked != check.Mixed {
		t.Fatalf("the box in the header is %+v", box)
	}
	if !screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: box.ID, Action: accessibility.Press}) ||
		!slices.Equal(toggled, []bool{true}) {
		t.Errorf("pressing the box asked for %v", toggled)
	}
	tree = screen.AccessibilityTree(w.Window)
	if n := tree.Nodes[box.ID]; n == nil || n.Checked != check.On {
		t.Errorf("after the press the box is %+v", n)
	}
	screen.Do(func() { table.Model = kvitui.NewBenchmarkTableModel(40); table.Refresh() })
	for _, n := range screen.AccessibilityTree(w.Window).Nodes {
		if n.Role == role.ColumnHeader && n.ColumnIndex == 1 && n.Name != "Select the rows shown" {
			t.Errorf("an untitled Check column's header is %+v", n)
		}
	}
}

func TestTheKeyboardWalksTheRowsAndTicksTheirBoxes(t *testing.T) {
	var table *kvitui.Table
	var toggled, pressed, activated []int
	rows := kvitui.NewBenchmarkTableModel(400)
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		table = kvitui.NewTable(ui, rows)
		table.HiddenColumns = []int{1, 2, 5, 8, 10, 11}
		table.OnCellToggled = func(row, column int, wanted bool) {
			toggled = append(toggled, row, column)
			rows.SetChecked(row, wanted)
			table.Refresh()
		}
		table.OnRowPressed = func(row int) { pressed = append(pressed, row) }
		table.OnRowActivated = func(row int) { activated = append(activated, row) }
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(300), table))}
	})
	screen.Do(func() { table.FocusRow(0) })
	screen.KeyPress(unison.KeyDown, 0)
	screen.Do(func() {
		if table.CurrentRow() != 1 || !slices.Equal(table.SelectedRows(), []int{1}) {
			t.Errorf("Down put the keyboard on row %d with %v selected", table.CurrentRow(), table.SelectedRows())
		}
	})
	// The focus is reported on the cell the keyboard is on, described by its
	// row, so a reader moving along a row knows which record it is.
	if n := focusedNode(screen, w); n == nil || n.Role != role.Cell || n.Name != "2020-01-02" ||
		!strings.Contains(n.Description, "Barrow & Co") {
		t.Errorf("the focused cell is %+v", n)
	}
	screen.KeyPress(unison.KeySpace, 0)
	if !slices.Equal(toggled, []int{1, 12}) || !rows.Cell(1, 12).Checked {
		t.Errorf("Space toggled %v", toggled)
	}
	screen.KeyPress(unison.KeyReturn, 0)
	if !slices.Equal(pressed, []int{1}) || !slices.Equal(activated, []int{1}) {
		t.Errorf("Return pressed %v and activated %v", pressed, activated)
	}
	screen.KeyPress(unison.KeyDown, mod.Shift)
	screen.Do(func() {
		if !slices.Equal(table.SelectedRows(), []int{1, 2}) {
			t.Errorf("Shift+Down selected %v", table.SelectedRows())
		}
		if got := table.CopySelection(); got != "2020-01-02\tBarrow & Co\tTransport\t-499.63\t9.77\tPending\n"+
			"2020-01-03\tColwyn & Co\tUtilities\t-499.26\t19.54\tDisputed" {
			t.Errorf("the copy is %q", got)
		}
	})
	// A box says it is a box and whether it is ticked.
	var box *accessibility.Node
	for _, n := range nodes(screen, w, role.CheckBox) {
		if n.RowIndex == 1 {
			box = n
		}
	}
	if box == nil || box.Name != "Select this row" || box.Checked != check.On {
		t.Errorf("row 2's box is %+v", box)
	}
}

// A press on a row's box ticks it and does not also open the row, which a
// press anywhere else on the row does.
func TestTickingABoxDoesNotAlsoOpenTheRow(t *testing.T) {
	var table *kvitui.Table
	var toggled, pressed []int
	screen, ui, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		table = kvitui.NewTable(ui, kvitui.NewBenchmarkTableModel(400))
		table.HiddenColumns = []int{1, 2, 5, 8, 10, 11}
		table.OnCellToggled = func(row, _ int, _ bool) { toggled = append(toggled, row) }
		table.OnRowPressed = func(row int) { pressed = append(pressed, row) }
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(300), table))}
	})
	m := ui.Interface
	// The box column starts after the six shown columns before it, 672
	// design pixels in; its box is a near space and the check's padding in.
	row := float32(m.RowHeightSlim() + m.RowHeightSlim()/2)
	box := geom.NewPoint(float32(m.Px(672)+m.SpaceNear()+6+m.Px(8)), row)
	screen.Click(screen.PanelPoint(table, box))
	if !slices.Equal(toggled, []int{0}) || len(pressed) != 0 {
		t.Errorf("a press on the box toggled %v and opened %v", toggled, pressed)
	}
	screen.Click(screen.PanelPoint(table, geom.NewPoint(float32(m.Px(40)), row)))
	screen.Do(func() {
		if !slices.Equal(pressed, []int{0}) || !table.IsRowSelected(0) {
			t.Errorf("a press on the row opened %v, selected %v", pressed, table.SelectedRows())
		}
	})
}

// A payee cut short, or shortened by the model, is shown whole under the
// pointer, and the amount behind a converted figure in a hover card.
func TestACellCutShortIsShownWholeUnderThePointer(t *testing.T) {
	var table *kvitui.Table
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		table = kvitui.NewTable(ui, kvitui.NewBenchmarkTableModel(400))
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(300), table))}
	})
	m := ui.Interface
	// The payee column starts after the date, reference and description.
	at := screen.PanelPoint(table, geom.NewPoint(float32(m.Px(96+110+220+40)), float32(m.RowHeightSlim()*3/2)))
	screen.MouseMove(at, 0)
	waitFor(t, screen, "the payee's tooltip", func() bool { return len(w.Popups()) == 1 })
	if n := screen.AccessibilityNodeFor(w.Popups()[0].Panel); n == nil || !strings.Contains(n.Name, "Ashford & Co (Holdings) Limited") {
		t.Errorf("the tooltip says %+v", n)
	}
	screen.MouseMove(screen.PanelPoint(table, geom.NewPoint(float32(m.Px(40)), float32(m.RowHeightSlim()*3/2))), 0)
	waitFor(t, screen, "the tooltip to go", func() bool { return len(w.Popups()) == 0 })
}

func TestTheEmptyStateFollowsTheRowCount(t *testing.T) {
	var table *kvitui.Table
	rows := kvitui.NewBenchmarkTableModel(1000)
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		table = kvitui.NewTable(ui, rows)
		table.EmptyTitle, table.EmptyDetail = "No transactions", "Nothing matches the current filter."
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(300), table))}
	})
	says := func() bool {
		for _, n := range screen.AccessibilityTree(w.Window).Nodes {
			if n.Name == "No transactions. Nothing matches the current filter." && !n.Ignored {
				return true
			}
		}
		return false
	}
	if says() {
		t.Error("a table with rows says it has none")
	}
	screen.Do(func() {
		rows.SetFilter("zzzzz")
		table.Reset()
	})
	screen.Sync()
	if !says() || table.RowCount() != 0 {
		t.Error("a table filtered to nothing does not say so")
	}
	screen.Do(func() {
		rows.SetFilter("")
		table.Reset()
	})
	screen.Sync()
	if says() {
		t.Error("the empty state stayed over rows that came back")
	}
}

func TestACellSaysWhatItHolds(t *testing.T) {
	var money, unmeasured, payee, marks *kvitui.Cell
	var picked *kvitui.Cell
	var asked []bool
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		money = kvitui.NewCell(ui, kvitui.CellMoney, kvitui.CellValue{Text: "-42.90", Unit: "GBP"})
		unmeasured = kvitui.NewCell(ui, kvitui.CellFigure, kvitui.CellValue{Unmeasured: true})
		payee = kvitui.NewCell(ui, kvitui.CellText, kvitui.CellValue{Text: "Harlow depot", FullText: "Harlow depot retainer"})
		marks = kvitui.NewCell(ui, kvitui.CellMarks, kvitui.CellValue{Marks: []kvitui.CellMark{
			kvitui.MarkFor(kvitui.ToneSuccess, "Settled"), kvitui.MarkFor(kvitui.ToneWarning, "Not reviewed")}})
		picked = kvitui.NewCell(ui, kvitui.CellCheck, kvitui.CellValue{Checked: true})
		picked.OnToggle = func(wanted bool) { asked = append(asked, wanted) }
		return []unison.Paneler{money, unmeasured, payee, marks, picked}
	})
	for _, c := range []struct {
		cell *kvitui.Cell
		name string
	}{{money, "-42.90 GBP"}, {unmeasured, "not measured"}, {payee, "Harlow depot retainer"}} {
		if n := screen.AccessibilityNodeFor(c.cell); n == nil || n.Role != role.Cell || n.Name != c.name {
			t.Errorf("a cell's node is %+v, want %q", n, c.name)
		}
	}
	if n := screen.AccessibilityNodeFor(marks); n == nil || n.Name != "2 marks" || len(n.Children) != 2 {
		t.Errorf("the marks cell's node is %+v", n)
	}
	// The box asks and does not tick itself.
	screen.Click(screen.PanelPoint(picked, geom.NewPoint(12, 15)))
	screen.Do(func() {
		if !slices.Equal(asked, []bool{false}) || !picked.Value.Checked {
			t.Errorf("a press asked %v and left the value %v", asked, picked.Value.Checked)
		}
	})
}

// What is behind a converted amount — the native amount, the rate and its
// date — is a hover card under the figure, shown at once and gone when the
// pointer leaves.
func TestAnAmountShowsWhatIsBehindItInAHoverCard(t *testing.T) {
	var fare *kvitui.Cell
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		fare = kvitui.NewCell(ui, kvitui.CellMoney, kvitui.CellValue{Text: "-8.12", Unit: "GBP",
			FullText: "-1,500 JPY\nat 184.7 JPY to the pound\nrate of 27 September 2026"})
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(200), fare)}
	})
	screen.MouseMove(screen.PanelCenter(fare), 0)
	waitFor(t, screen, "the hover card", func() bool { return len(w.Popups()) == 1 })
	screen.Do(func() {
		card := w.Popups()[0].Panel.AsPanel()
		if n := len(card.Children()); n != 3 {
			t.Errorf("the card has %d lines, want 3", n)
		}
		box := w.Content().RectFromRoot(fare.RectToRoot(fare.ContentRect(true)))
		if r := card.FrameRect(); r.Right() != box.Right() || r.Y != box.Bottom() {
			t.Errorf("the card is at %v under a cell at %v", r, box)
		}
	})
	screen.MouseMove(geom.NewPoint(1400, 900), 0)
	waitFor(t, screen, "the card to go", func() bool { return len(w.Popups()) == 0 })
}
