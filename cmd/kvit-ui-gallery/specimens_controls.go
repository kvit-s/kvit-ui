package main

// The Controls group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

func tabForms(ui *kvitui.UI) unison.Paneler {
	all := kvitui.NewTab(ui, "All")
	all.Selected, all.Count = true, 1284
	uncategorised := kvitui.NewTab(ui, "Uncategorised")
	uncategorised.Count = 47
	return kvitui.Row(ui, kvitui.Px(0), all, uncategorised, kvitui.NewTab(ui, "Disputed"))
}

func tabExplanations(ui *kvitui.UI) unison.Paneler {
	changes := kvitui.NewTab(ui, "Changes")
	changes.Selected = true
	changes.Explanation = "What this version changed."
	source := kvitui.NewTab(ui, "Source")
	source.Explanation = "The file’s own text."
	history := kvitui.NewTab(ui, "History")
	history.Explanation = "Every version and who wrote it."
	return kvitui.Row(ui, kvitui.Px(0), changes, source, history)
}
