package main

// The Marks group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

func badgeCounts(ui *kvitui.UI) unison.Paneler {
	neutral := kvitui.NewBadge(ui, 42)
	neutral.Tone = kvitui.BadgeNeutral
	danger := kvitui.NewBadge(ui, 1204)
	danger.Tone = kvitui.BadgeDanger
	return kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewBadge(ui, 3), neutral, danger, kvitui.NewBadge(ui, 0))
}

func badgeNouns(ui *kvitui.UI) unison.Paneler {
	// Same three pills, three different sentences for a screen reader:
	// "1 decision", "214 decisions", "214 entries". One already inflected
	// word would be right at one count and wrong at every other.
	one := kvitui.NewBadge(ui, 1)
	one.Counted = "decision"
	many := kvitui.NewBadge(ui, 214)
	many.Counted, many.Max = "decision", 999
	entries := kvitui.NewBadge(ui, 214)
	entries.Max, entries.Tone = 999, kvitui.BadgeNeutral
	entries.Counted, entries.CountedPlural = "entry", "entries"
	return kvitui.Row(ui, kvitui.SizeSpace, one, many, entries)
}
