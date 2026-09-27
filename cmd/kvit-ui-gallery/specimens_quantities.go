package main

// The Quantities group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

func figureForms(ui *kvitui.UI) unison.Paneler {
	unmeasured := kvitui.NewFigure(ui, "", "")
	unmeasured.Measured = false
	bounded := kvitui.NewFigure(ui, "12", "d")
	bounded.Bounded = true
	large := kvitui.NewFigure(ui, "94", "%")
	large.Role = kvitui.RoleDisplay
	return kvitui.Column(ui, kvitui.SizeSpaceNear,
		kvitui.NewFigure(ui, "1,284.50", "GBP"), kvitui.NewFigure(ui, "0", "d"), unmeasured, bounded, large)
}
