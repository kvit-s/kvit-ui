package kvitui

import (
	"strconv"
	"strings"
	"time"
)

// BenchmarkRows is how many rows a BenchmarkTableModel generates unless told
// otherwise: the size of kvit-cash's transaction browser.
const BenchmarkRows = 250000

var (
	benchVerbs      = []string{"Payment", "Transfer", "Refund", "Deposit", "Withdrawal", "Fee", "Interest", "Adjustment"}
	benchPlaces     = []string{"Ashford", "Barrow", "Colwyn", "Dunmore", "Elmsley", "Fairholt", "Grantham", "Harlow", "Ivybridge", "Jarrow", "Kelmscott", "Lyndhurst"}
	benchCategories = []string{"Groceries", "Transport", "Utilities", "Rent", "Leisure", "Health", "Savings", "Income"}
	benchAccounts   = []string{"Current", "Savings", "Card", "Joint"}
	benchStatuses   = []string{"Settled", "Pending", "Disputed"}
	benchStart      = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	// Twelve columns of value and one of boxes. The claim the model is for is
	// about the twelve; the thirteenth is here so the Check kind is exercised
	// by the one model the library ships.
	benchColumns = []TableColumn{
		{Title: "Date", Width: 96, Kind: CellDate},
		{Title: "Reference", Width: 110, Kind: CellSlug},
		{Title: "Description", Width: 220},
		{Title: "Payee", Width: 150},
		{Title: "Category", Width: 110, Kind: CellChip},
		{Title: "Account", Width: 110},
		{Title: "Amount", Width: 96, Kind: CellMoney},
		{Title: "Balance", Width: 110, Kind: CellFigure},
		{Title: "Currency", Width: 64},
		{Title: "Status", Width: 110, Kind: CellMarks},
		{Title: "Tags", Width: 110, Kind: CellChip},
		{Title: "Note", Width: 200},
		{Title: "", Width: 40, Kind: CellCheck, Unsortable: true},
	}
)

// The columns the model treats specially.
const (
	benchPayee  = 3
	benchStatus = 9
	benchSelect = 12
)

// BenchmarkTableModel is a ledger of generated transactions, thirteen
// columns wide, that holds no row: every cell is worked out from its row
// number when it is asked for. It is the model the gallery shows a quarter of
// a million rows with, and the one the check that a Table scrolls smoothly and
// filters in under 100 ms at that size is run against.
//
// The ticks in its column of boxes are kept by generated row rather than by
// the row shown, so a reader who ticks forty rows, narrows the filter and
// widens it again finds the forty still ticked.
type BenchmarkTableModel struct {
	total        int
	filter       string
	matched      []int32
	checked      map[int]bool
	checkedShown int // how many of the rows shown are ticked, kept as they change
}

// NewBenchmarkTableModel returns a model of total generated rows.
func NewBenchmarkTableModel(total int) *BenchmarkTableModel {
	return &BenchmarkTableModel{total: max(0, total), checked: map[int]bool{}}
}

// Total is how many rows the model generates, whatever the filter shows.
func (b *BenchmarkTableModel) Total() int { return b.total }

// Rows is how many rows the filter shows.
func (b *BenchmarkTableModel) Rows() int {
	if b.filter == "" {
		return b.total
	}
	return len(b.matched)
}

// Columns are the thirteen columns.
func (b *BenchmarkTableModel) Columns() []TableColumn { return benchColumns }

// source is the generated row behind a row shown, or -1.
func (b *BenchmarkTableModel) source(row int) int {
	if b.filter == "" {
		if row < 0 || row >= b.total {
			return -1
		}
		return row
	}
	if row < 0 || row >= len(b.matched) {
		return -1
	}
	return int(b.matched[row])
}

func benchAmount(row int) float64 { return float64(row*37%100000)/100 - 500 }

func benchDescription(row int) string {
	return benchVerbs[row%len(benchVerbs)] + " to " + benchPlaces[(row/7)%len(benchPlaces)] + " depot"
}

func benchMarks(row int) []CellMark {
	var marks []CellMark
	switch row % len(benchStatuses) {
	case 0:
		marks = append(marks, CellMark{Tone: ToneSuccess, Shape: ShapeCircle, Label: "Settled"})
	case 1:
		marks = append(marks, CellMark{Tone: ToneWarning, Shape: ShapeDiamond, Label: "Pending"})
	default:
		marks = append(marks, CellMark{Tone: ToneDanger, Shape: ShapeSquare, Label: "Disputed"})
	}
	if row%7 == 0 {
		marks = append(marks, CellMark{Tone: ToneNeutral, Shape: ShapeSquare, Label: "Not reviewed"})
	}
	return marks
}

// Cell works one cell out from its row number.
func (b *BenchmarkTableModel) Cell(row, column int) CellValue {
	r := b.source(row)
	if r < 0 {
		return CellValue{Unmeasured: true}
	}
	switch column {
	case 0:
		return CellValue{Date: benchStart.AddDate(0, 0, r%2200)}
	case 1:
		return CellValue{Text: "TX-" + pad8(r)}
	case 2:
		return CellValue{Text: benchDescription(r)}
	case benchPayee:
		place := benchPlaces[r%len(benchPlaces)]
		return CellValue{Text: place + " & Co", FullText: place + " & Co (Holdings) Limited"}
	case 4:
		return CellValue{Text: benchCategories[r%len(benchCategories)]}
	case 5:
		return CellValue{Text: benchAccounts[r%len(benchAccounts)]}
	case 6:
		return CellValue{Text: strconv.FormatFloat(benchAmount(r), 'f', 2, 64)}
	case 7:
		if r%40 == 0 {
			return CellValue{Unmeasured: true}
		}
		return CellValue{Text: strconv.FormatFloat(float64(r*977%4000000)/100, 'f', -1, 64)}
	case 8:
		return CellValue{Text: "GBP"}
	case benchStatus:
		return CellValue{Text: benchStatuses[r%len(benchStatuses)], Marks: benchMarks(r)}
	case 10:
		return CellValue{Text: "#" + strings.ToLower(benchCategories[(r/3)%len(benchCategories)])}
	case 11:
		if r%11 == 0 {
			return CellValue{Text: "Checked against statement"}
		}
		return CellValue{}
	case benchSelect:
		return CellValue{Checked: b.checked[r]}
	}
	return CellValue{}
}

func pad8(n int) string {
	s := strconv.Itoa(n)
	if len(s) < 8 {
		s = strings.Repeat("0", 8-len(s)) + s
	}
	return s
}

// Filter is the text the descriptions shown contain.
func (b *BenchmarkTableModel) Filter() string { return b.filter }

// SetFilter shows only the rows whose description contains text, ignoring
// case; "" shows them all. It tests every generated row, which is the
// operation the 100 ms budget is about.
func (b *BenchmarkTableModel) SetFilter(text string) {
	b.filter = text
	b.matched = b.matched[:0]
	if text != "" {
		want := strings.ToLower(text)
		buf := make([]byte, 0, 64)
		for r := range b.total {
			buf = append(buf[:0], benchVerbs[r%len(benchVerbs)]...)
			buf = append(buf, " to "...)
			buf = append(buf, benchPlaces[(r/7)%len(benchPlaces)]...)
			buf = append(buf, " depot"...)
			if containsFold(buf, want) {
				b.matched = append(b.matched, int32(r))
			}
		}
	}
	b.recountChecked()
}

// containsFold reports whether s holds want, which is already lower case,
// ignoring the case of ASCII letters in s.
func containsFold(s []byte, want string) bool {
	n := len(want)
	for i := 0; i+n <= len(s); i++ {
		match := true
		for j := range n {
			c := s[i+j]
			if 'A' <= c && c <= 'Z' {
				c += 'a' - 'A'
			}
			if c != want[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// SetChecked ticks one row shown, or clears it.
func (b *BenchmarkTableModel) SetChecked(row int, checked bool) {
	r := b.source(row)
	if r < 0 || checked == b.checked[r] {
		return
	}
	if checked {
		b.checked[r] = true
		b.checkedShown++
	} else {
		delete(b.checked, r)
		b.checkedShown--
	}
}

// SetEveryShownChecked ticks every row the filter shows, or clears them.
func (b *BenchmarkTableModel) SetEveryShownChecked(checked bool) {
	rows := b.Rows()
	for row := range rows {
		if checked {
			b.checked[b.source(row)] = true
		} else {
			delete(b.checked, b.source(row))
		}
	}
	b.checkedShown = 0
	if checked {
		b.checkedShown = rows
	}
}

// AllShownChecked reports whether every row shown is ticked, and
// SomeShownChecked whether some are and not all: the three states of the box
// in the column's header.
func (b *BenchmarkTableModel) AllShownChecked() bool {
	return b.Rows() > 0 && b.checkedShown == b.Rows()
}

// SomeShownChecked reports whether some rows shown are ticked and not all.
func (b *BenchmarkTableModel) SomeShownChecked() bool {
	return b.checkedShown > 0 && !b.AllShownChecked()
}

// recountChecked counts the ticked rows the filter shows, once when the
// filter changes rather than whenever the header asks, since the header asks
// on every redraw.
func (b *BenchmarkTableModel) recountChecked() {
	b.checkedShown = 0
	if len(b.checked) == 0 {
		return
	}
	for row := range b.Rows() {
		if b.checked[b.source(row)] {
			b.checkedShown++
		}
	}
}
