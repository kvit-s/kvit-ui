package main

// The Data group's specimens.

import (
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

func barForms(ui *kvitui.UI) unison.Paneler {
	attention := kvitui.NewBar(ui, 62, 100)
	attention.Label, attention.Unit = "Attention", "h"
	agent := kvitui.NewBar(ui, 38, 100)
	agent.Wide, agent.Bounded, agent.Ink = true, true, kvitui.InkAxisAgent
	agent.Label, agent.Unit = "Agent", "h"
	unrecorded := kvitui.NewBar(ui, 0, 100)
	unrecorded.Measured, unrecorded.Label = false, "Unrecorded"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose, attention, agent, unrecorded))
}

func stackedBarBreakdown(ui *kvitui.UI) unison.Paneler {
	month := kvitui.NewStackedBar(ui,
		kvitui.Segment{Value: 420, Label: "rent", Ink: kvitui.InkCategorical(0)},
		kvitui.Segment{Value: 180, Label: "groceries", Ink: kvitui.InkCategorical(1)},
		kvitui.Segment{Value: 95, Label: "transport", Ink: kvitui.InkCategorical(2)},
		kvitui.Segment{Value: 60, Label: "everything else", Ink: kvitui.InkCategorical(3)})
	month.Wide, month.Label = true, "Where the month went"
	return kvitui.FullWidth(month)
}

func sparkWithHole(ui *kvitui.UI) unison.Paneler {
	gap := kvitui.NotMeasured
	notes := kvitui.NewSpark(ui, 3, 5, 8, 6, gap, gap, 9, 12, 7, 4, 6, 11)
	notes.Label = "Notes written"
	return kvitui.Width(ui, kvitui.Px(160), notes)
}

func trendAndEmpty(ui *kvitui.UI) unison.Paneler {
	balance := kvitui.NewTrend(ui, 400, 620, 580, 900, kvitui.NotMeasured, 1400, 1250, 1700)
	balance.Label, balance.Unit, balance.MaximumY = "Balance", "GBP", 2000
	savings := kvitui.NewTrend(ui)
	savings.Label = "Savings"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose,
		balance, kvitui.Height(ui, kvitui.Px(90), savings)))
}

func trendTwoSeries(ui *kvitui.UI) unison.Paneler {
	gap := kvitui.NotMeasured
	worth := kvitui.NewTrend(ui, 400, 620, 580, 900, 1150, 1400, 1250, 1700)
	worth.Label, worth.SecondLabel, worth.Unit, worth.MaximumY = "Assets", "Liabilities", "GBP", 2000
	worth.Second = []float64{gap, gap, 300, 340, 320, 290, 260, 240}
	// Where the second account's history begins, drawn over the plot with
	// the same arithmetic the lines use.
	worth.Overlay = func(gc *unison.Canvas, plot geom.Rect) {
		x := plot.X + plot.Width*2/7
		mark := geom.NewRect(x, plot.Y, float32(ui.Interface.Hairline()), plot.Height)
		gc.DrawRect(mark, kvitui.Color(kvitui.InkMarker.Of(ui)).Paint(gc, mark, paintstyle.Fill))
		note := ui.Fonts.Layout(kvitui.Caption(ui, "history starts here", kvitui.InkTextFaint), kvitui.NoOptions)
		_, h := note.Size()
		note.Draw(gc, x+float32(ui.Interface.SpaceSnug()), plot.Bottom()-h)
	}
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(150), worth))
}

func distributionRows(ui *kvitui.UI) unison.Paneler {
	card := kvitui.NewDistribution(ui)
	card.Label, card.Unit, card.ScaleMaximum = "Time to settle, card", "d", 30
	card.Minimum, card.LowerQuartile, card.Median, card.UpperQuartile, card.Maximum = 1, 2, 3, 5, 21
	transfer := kvitui.NewDistribution(ui)
	transfer.Label, transfer.Unit, transfer.ScaleMaximum = "Time to settle, transfer", "d", 30
	transfer.Minimum, transfer.LowerQuartile, transfer.Median, transfer.UpperQuartile, transfer.Maximum = 1, 1, 2, 2, 4
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose, card, transfer))
}

func gaugePace(ui *kvitui.UI) unison.Paneler {
	gauge := func(value float64, label string) unison.Paneler {
		g := kvitui.NewGauge(ui, value, 600)
		g.Pace, g.Label, g.Unit = 0.6, label, "GBP"
		return g
	}
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose,
		gauge(240, "Groceries"), gauge(480, "Leisure"), gauge(720, "Transport")))
}

func deltaDirections(ui *kvitui.UI) unison.Paneler {
	delta := func(change float64, unit string, precision int, good kvitui.Direction) *kvitui.Delta {
		d := kvitui.NewDelta(ui, change)
		d.Unit, d.Precision, d.Good = unit, precision, good
		return d
	}
	unmeasured := kvitui.NewDelta(ui, 0)
	unmeasured.Measured = false
	return kvitui.Row(ui, kvitui.SizeSpaceLoose,
		delta(12.4, "%", 1, kvitui.Up), delta(-3.2, "%", 1, kvitui.Up), delta(18, "GBP", 0, kvitui.Down),
		delta(0, "", 0, kvitui.Neither), unmeasured)
}

func statTileRow(ui *kvitui.UI) unison.Paneler {
	balance := kvitui.NewStatTile(ui, "Balance", "4,182.30", "GBP")
	balance.HasChange, balance.Change, balance.Good = true, 240.10, kvitui.Up
	balance.History = []float64{3200, 3400, 3390, 3800, 4000, 4182}
	spent := kvitui.NewStatTile(ui, "Spent this month", "812.40", "GBP")
	spent.HasChange, spent.Change, spent.Good = true, 96.20, kvitui.Down
	spent.Caption = "12 days remaining"
	uncategorised := kvitui.NewStatTile(ui, "Uncategorised", "", "")
	uncategorised.Measured = false
	return kvitui.Row(ui, kvitui.SizeColumnGap, balance, spent, uncategorised)
}

func figureBlockRow(ui *kvitui.UI) unison.Paneler {
	reconciled := kvitui.NewFigureBlock(ui, "", "", "Reconciled")
	reconciled.Figure.Measured = false
	return kvitui.Row(ui, kvitui.Px(40),
		kvitui.NewFigureBlock(ui, "1,284", "", "Transactions"), kvitui.NewFigureBlock(ui, "97", "%", "Categorised"), reconciled)
}

func cellKinds(ui *kvitui.UI) unison.Paneler {
	cell := func(kind kvitui.CellKind, v kvitui.CellValue) unison.Paneler {
		return kvitui.FullWidth(kvitui.NewCell(ui, kind, v))
	}
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceTight,
		cell(kvitui.CellText, kvitui.CellValue{Text: "Payment to Harlow depot"}),
		cell(kvitui.CellDate, kvitui.CellValue{Date: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)}),
		cell(kvitui.CellSlug, kvitui.CellValue{Text: "TX-00173404"}),
		cell(kvitui.CellFigure, kvitui.CellValue{Text: "1284.5", Unit: "GBP"}),
		// Money takes an amount that has already been formatted. The cell
		// does no arithmetic: an amount is a count of the minor unit of its
		// own currency, and how many minor digits it has is the currency's.
		cell(kvitui.CellMoney, kvitui.CellValue{Text: "-42.90", Unit: "GBP"}),
		cell(kvitui.CellChip, kvitui.CellValue{Text: "Settled", Tone: kvitui.ToneSuccess}),
		cell(kvitui.CellFigure, kvitui.CellValue{Unmeasured: true})))
}

func cellStatesAndPick(ui *kvitui.UI) unison.Paneler {
	// Every state a row is in at once, one dot each. Each mark has a tone,
	// the shape that says the same thing without colour, and the word that
	// is both the tooltip and what a screen reader says.
	marks := kvitui.NewCell(ui, kvitui.CellMarks, kvitui.CellValue{Marks: []kvitui.CellMark{
		{Tone: kvitui.ToneSuccess, Shape: kvitui.ShapeCircle, Label: "Settled"},
		{Tone: kvitui.ToneWarning, Shape: kvitui.ShapeDiamond, Label: "Not reviewed"},
		{Tone: kvitui.ToneNeutral, Shape: kvitui.ShapeSquare, Label: "Has an attachment"},
	}})
	// The box does not tick itself. Whatever owns the selection does, and
	// the cell redraws from it, so a request that is refused leaves no tick.
	picked := kvitui.NewCell(ui, kvitui.CellCheck, kvitui.CellValue{Checked: true})
	picked.OnToggle = func(wanted bool) {
		picked.Value.Checked = wanted
		picked.MarkForLayoutAndRedraw()
	}
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceTight, kvitui.FullWidth(marks), kvitui.FullWidth(picked)))
}

func cellCutShort(ui *kvitui.UI) unison.Paneler {
	// A value wider than its column. The model gives the whole of it as
	// FullText, and the cell shows it under the pointer and under the
	// keyboard cursor, which the view says is here through Current.
	return kvitui.Width(ui, kvitui.Px(150), kvitui.NewCell(ui, kvitui.CellText, kvitui.CellValue{
		Text:     "Harlow depot retainer",
		FullText: "Harlow depot — quarterly maintenance retainer",
	}))
}

func tableQuarterMillion(ui *kvitui.UI) unison.Paneler {
	rows := kvitui.NewBenchmarkTableModel(kvitui.BenchmarkRows)
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(260), kvitui.NewTable(ui, rows)))
}

func tableSortingAndBoxes(ui *kvitui.UI) unison.Paneler {
	rows := kvitui.NewBenchmarkTableModel(400)
	table := kvitui.NewTable(ui, rows)
	// Six of the thirteen columns, so the boxes at the right-hand end are on
	// screen without scrolling.
	table.HiddenColumns = []int{1, 2, 5, 8, 10, 11}
	// The table draws the indicator and says what was asked for; putting the
	// rows in that order is the application's, because the model is. Nothing
	// here sorts, so the arrow stays where it was set.
	table.SortColumn, table.SortAscending = 0, false
	// The column of boxes, and the box in its header for every row shown.
	// The ticks are the model's, kept by the underlying row, so narrowing
	// the filter does not lose them.
	ticked := func() {
		table.HeaderChecked, table.HeaderPartial = rows.AllShownChecked(), rows.SomeShownChecked()
		table.Refresh()
	}
	table.OnHeaderToggled = func(wanted bool) { rows.SetEveryShownChecked(wanted); ticked() }
	table.OnCellToggled = func(row, _ int, wanted bool) { rows.SetChecked(row, wanted); ticked() }
	// The mark that says a press on a row opens something, drawn at rest in
	// a strip of its own rather than in a column, so a reader can tell which
	// lists lead anywhere without sweeping the pointer across them.
	table.RowsOpen = true
	table.OnRowPressed = func(row int) { rows.SetChecked(row, true); ticked() }
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(260), table))
}

func tableEmpty(ui *kvitui.UI) unison.Paneler {
	table := kvitui.NewTable(ui, kvitui.NewBenchmarkTableModel(0))
	table.EmptyTitle, table.EmptyDetail = "No transactions", "Nothing matches the current filter."
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(200), table))
}
