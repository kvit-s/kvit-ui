package main

// The Flow group's specimens.

import (
	"fmt"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

func scrollBarBesideColumn(ui *kvitui.UI) unison.Paneler {
	rows := make([]unison.Paneler, 10)
	for i := range rows {
		rows[i] = kvitui.NewSlimRow(ui, fmt.Sprintf("Row %d", i+1))
	}
	region := kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.Px(0), rows...))
	return kvitui.Sized(ui, kvitui.Px(480), kvitui.Px(120), region)
}

func segmentedPeriod(ui *kvitui.UI) unison.Paneler {
	period := kvitui.NewSegmented(ui, "Period",
		kvitui.Option{Value: "week", Label: "Week"}, kvitui.Option{Value: "month", Label: "Month"},
		kvitui.Option{Value: "quarter", Label: "Quarter"}, kvitui.Option{Value: "year", Label: "Year"})
	period.Current = "month"
	whose := kvitui.NewSegmented(ui, "",
		kvitui.Option{Value: "All", Label: "All"}, kvitui.Option{Value: "Mine", Label: "Mine"})
	return kvitui.Column(ui, kvitui.SizeSpace, kvitui.Left(period), kvitui.Left(whose))
}
