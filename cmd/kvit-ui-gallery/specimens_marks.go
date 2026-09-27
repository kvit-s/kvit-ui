package main

// The Marks group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
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

func chipTones(ui *kvitui.UI) unison.Paneler {
	chip := func(word string, tone kvitui.Tone) unison.Paneler {
		c := kvitui.NewChip(ui, word)
		c.Tone = tone
		return c
	}
	settled := kvitui.NewChip(ui, "settled")
	settled.Tone, settled.Strong = kvitui.ToneSuccess, true
	disputed := kvitui.NewChip(ui, "disputed")
	disputed.Tone, disputed.Strong, disputed.Symbol = kvitui.ToneDanger, true, "warning"
	return kvitui.Column(ui, kvitui.SizeSpace,
		kvitui.Left(kvitui.Row(ui, kvitui.SizeSpaceNear, chip("neutral", kvitui.ToneNeutral), chip("accent", kvitui.ToneAccent),
			chip("success", kvitui.ToneSuccess), chip("warning", kvitui.ToneWarning), chip("danger", kvitui.ToneDanger),
			chip("info", kvitui.ToneInfo))),
		kvitui.Left(kvitui.Row(ui, kvitui.SizeSpaceNear, settled, disputed)))
}

func chipExplanation(ui *kvitui.UI) unison.Paneler {
	// A chip that names a kind of thing is finished at its word. A chip that
	// names a state somebody has to act on is not, and Explanation is where
	// the rest of it goes: the tooltip and the accessible description.
	stale := kvitui.NewChip(ui, "stale")
	stale.Tone = kvitui.ToneWarning
	stale.Explanation = "The note changed since this was staged. Update redoes it against the note as it stands now."
	readOnly := kvitui.NewChip(ui, "read only")
	readOnly.Explanation = "Another session holds the write lease on this copy, so this turn only reads."
	return kvitui.Row(ui, kvitui.SizeSpaceNear, stale, readOnly)
}

func tagForms(ui *kvitui.UI) unison.Paneler {
	colours := tokens.ColorPalette()
	groceries := kvitui.NewTag(ui, "groceries")
	groceries.SetTint(colours[0])
	transport := kvitui.NewTag(ui, "transport")
	transport.SetTint(colours[2])
	removable := kvitui.NewTag(ui, "removable")
	removable.SetTint(colours[4])
	removable.Removable = true
	return kvitui.Row(ui, kvitui.SizeSpaceNear, groceries, transport, kvitui.NewTag(ui, "untinted"), removable)
}

func slugGrounds(ui *kvitui.UI) unison.Paneler {
	bare := kvitui.NewSlug(ui, "9f2c1ab4e77d0031")
	bare.Ground = false
	return kvitui.Column(ui, kvitui.SizeSpaceNear, kvitui.Left(kvitui.NewSlug(ui, "TX-00173404")), kvitui.Left(bare))
}

func dotLevels(ui *kvitui.UI) unison.Paneler {
	dot := func(ink kvitui.Ink, shape kvitui.Shape, hollow bool, label string) unison.Paneler {
		d := kvitui.NewDot(ui)
		d.Ink, d.Shape, d.Hollow, d.Label = ink, shape, hollow, label
		return d
	}
	return kvitui.Row(ui, kvitui.SizeSpace,
		dot(kvitui.InkSuccess, kvitui.ShapeCircle, false, "healthy"),
		dot(kvitui.InkWarning, kvitui.ShapeSquare, false, "slipping"),
		dot(kvitui.InkDanger, kvitui.ShapeDiamond, false, "stalled"),
		dot(kvitui.InkTextMuted, kvitui.ShapeCircle, true, "not measured"))
}

func signalCounts(ui *kvitui.UI) unison.Paneler {
	// At one the mark is the whole statement. Past one the count goes inside
	// it, and the mark grows sideways to hold the digits rather than growing
	// taller and pushing the row apart.
	signal := func(count int, ink kvitui.Ink, shape kvitui.Shape, hollow bool, label string) unison.Paneler {
		s := kvitui.NewSignal(ui, label)
		s.Count, s.Ink, s.Shape, s.Hollow = count, ink, shape, hollow
		return s
	}
	return kvitui.Column(ui, kvitui.SizeSpace,
		kvitui.Left(kvitui.Row(ui, kvitui.SizeSpace,
			signal(1, kvitui.InkAccent, kvitui.ShapeSquare, false, "1 needs you"),
			signal(1, kvitui.InkSuccess, kvitui.ShapeCircle, false, "1 running"),
			signal(1, kvitui.InkDanger, kvitui.ShapeSquare, true, "1 failed"))),
		kvitui.Left(kvitui.Row(ui, kvitui.SizeSpace,
			signal(4, kvitui.InkAccent, kvitui.ShapeSquare, false, "4 need you"),
			signal(12, kvitui.InkSuccess, kvitui.ShapeCircle, false, "12 running"),
			signal(240, kvitui.InkDanger, kvitui.ShapeSquare, true, "240 failed"))))
}

func signalBesideWords(ui *kvitui.UI) unison.Paneler {
	// Where these are drawn: at the head of a row, in front of what the row
	// is about. The mark is a caption's height plus a margin, which lets it
	// sit in a line of text without setting the line's height.
	scoutSignal := kvitui.NewSignal(ui, "2 need you")
	scoutSignal.Count = 2
	scoutName := kvitui.NewLabel(ui, "Dialog Scout")
	scoutName.Role = kvitui.RoleSmall
	scout := kvitui.NewListRow(ui, scoutSignal, scoutName)
	scout.Label = "Dialog Scout, 2 need you"

	nightlySignal := kvitui.NewSignal(ui, "1 running")
	nightlySignal.Ink, nightlySignal.Shape = kvitui.InkSuccess, kvitui.ShapeCircle
	nightlyName := kvitui.NewLabel(ui, "Nightly check")
	nightlyName.Role = kvitui.RoleSmall
	nightly := kvitui.NewListRow(ui, nightlySignal, nightlyName)
	nightly.Label = "Nightly check, running"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), scout, nightly))
}

func pipCounts(ui *kvitui.UI) unison.Paneler {
	days := kvitui.NewPip(ui, 3, 5)
	days.Label = "3 of 5 days recorded"
	checks := kvitui.NewPip(ui, 0, 4)
	checks.Ink, checks.Label = kvitui.InkWarning, "no checks passed"
	return kvitui.Column(ui, kvitui.SizeSpace, kvitui.Left(days), kvitui.Left(checks))
}
