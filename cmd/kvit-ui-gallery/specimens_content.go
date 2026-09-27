package main

// The Content group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

func dividerBothWays(ui *kvitui.UI) unison.Paneler {
	down := kvitui.NewDivider(ui)
	down.Vertical = true
	down.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Fill, VGrab: true})
	beside := kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewLabel(ui, "left"), down, kvitui.NewLabel(ui, "right"))
	beside.SetLayout(kvitui.AtLeast(ui, kvitui.Px(24), beside.Layout()))
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewDivider(ui), beside))
}

func panelWithRules(ui *kvitui.UI) unison.Paneler {
	p := kvitui.NewPanel(ui)
	p.RuleTop, p.RuleBottom = true, true
	label := kvitui.NewLabel(ui, "A panel")
	label.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Middle, VAlign: align.Middle, HGrab: true, VGrab: true})
	p.AddChild(label)
	p.SetLayout(kvitui.AtLeast(ui, kvitui.Px(60), &unison.FlexLayout{Columns: 1}))
	return kvitui.FullWidth(p)
}

func slimRowMeasured(ui *kvitui.UI) unison.Paneler {
	library := kvitui.NewSlimRow(ui, "kvit-ui")
	library.Symbol, library.Kind, library.Phrase = "folder", "library", "waiting on review"
	library.Figure, library.Unit = "3.5", "d"
	app := kvitui.NewSlimRow(ui, "kvit-cash")
	app.Symbol, app.Kind, app.Phrase = "folder", "application", "not started"
	app.Measured = false
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), library, app))
}

func rowPressableAndStatic(ui *kvitui.UI) unison.Paneler {
	// A row that opens something is reachable by the keyboard, and takes the
	// hover tint, the press and the chevron with it.
	opens := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, "opens the record")))
	opens.Interactive, opens.Label = true, "Groceries"
	// A field name beside its value declares nothing, draws no hover tint and
	// answers no press.
	static := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, "layout only")))
	static.Label = "Amount"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), opens, static))
}

func rowFromTheKeyboard(ui *kvitui.UI) unison.Paneler {
	// Tab reaches this row, the ring says which row the keyboard is on, and
	// Return, Enter or Space opens it. A screen reader's press action does
	// the same one thing.
	groceries := kvitui.NewListRow(ui, kvitui.NewLabel(ui, "keyboard focus — Return opens it"))
	groceries.Form, groceries.Interactive, groceries.Label = kvitui.RowSub, true, "Groceries"
	groceries.Focus()
	// A row is usually a container, and the key may be meant for something
	// inside it. The row answers only while the row itself has the keyboard,
	// so Space on this button presses the button and leaves the record shut.
	split := kvitui.NewButton(ui, "Split")
	split.Form = kvitui.ButtonQuiet
	rent := kvitui.NewListRow(ui, kvitui.FullWidth(kvitui.NewLabel(ui, "the button takes its own Space")), split)
	rent.Form, rent.Interactive, rent.Label = kvitui.RowSub, true, "Rent"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), groceries, rent))
}

func rowHeights(ui *kvitui.UI) unison.Paneler {
	row := func(form kvitui.RowForm, label, note string) *kvitui.ListRow {
		r := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, note)))
		r.Form, r.Label = form, label
		return r
	}
	slim := row(kvitui.RowSlim, "Slim", "slim — 30, selected")
	slim.Selected = true
	compact := row(kvitui.RowCompact, "Compact", "compact — 24, keyboard focus")
	compact.Current = true
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0),
		row(kvitui.RowFull, "Full", "full — 56"), row(kvitui.RowSub, "Sub", "sub — 48"), slim, compact))
}
