package main

// The Feedback group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

func tooltipOnButton(ui *kvitui.UI) unison.Paneler {
	reconcile := kvitui.NewButton(ui, "Reconcile")
	// Shown here for good; a control shows its own on hover and on keyboard
	// focus.
	ui.ShowTooltip(reconcile, "Match these against the statement")
	return kvitui.FullWidth(kvitui.At(ui, kvitui.Px(60), kvitui.Px(40), kvitui.Px(80), reconcile))
}

func popoverWithForm(ui *kvitui.UI) unison.Paneler {
	settled := kvitui.NewCheck(ui, "Settled")
	settled.Checked = true
	filter := kvitui.NewPopover(ui, "Filter", settled, kvitui.NewCheck(ui, "Pending"), kvitui.NewCheck(ui, "Disputed"))
	filter.Width = kvitui.Px(240)
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(150), unison.NewPanel()))
	filter.Open(stage, kvitui.PlaceOver(stage))
	return stage
}

func hintOpen(ui *kvitui.UI) unison.Paneler {
	hint := kvitui.NewHint(ui, "About automatic matching",
		"Automatic matching compares the date, amount and reference. It never changes the imported statement.")
	hint.Open()
	return kvitui.FullWidth(kvitui.At(ui, nil, nil, kvitui.Px(130), hint))
}

func hoverCardRow(ui *kvitui.UI) unison.Paneler {
	name := kvitui.NewLabel(ui, "Payment to Harlow depot")
	name.Role = kvitui.RoleStrong
	when := kvitui.NewLabel(ui, "17 March 2026 at 09:14")
	when.Role, when.Ink = kvitui.RoleSmall, kvitui.InkTextMuted
	return kvitui.NewHoverCard(ui, name, when, kvitui.NewFigure(ui, "1,284.50", "GBP"))
}

func toastTones(ui *kvitui.UI) unison.Paneler {
	toast := func(words string, tone kvitui.Tone, action string) unison.Paneler {
		t := kvitui.NewToast(ui, words)
		t.Tone, t.Action = tone, action
		return kvitui.Left(t)
	}
	return kvitui.Column(ui, kvitui.SizeSpace,
		toast("Copied to the clipboard", kvitui.ToneInfo, ""),
		toast("Statement imported", kvitui.ToneSuccess, ""),
		toast("Two rows could not be matched", kvitui.ToneWarning, ""),
		toast("The file could not be read", kvitui.ToneDanger, ""),
		toast("40 transactions archived", kvitui.ToneSuccess, "Undo"))
}

func noticeForms(ui *kvitui.UI) unison.Paneler {
	licence := kvitui.NewNotice(ui, "Your licence expires in three days")
	licence.Tone, licence.Detail = kvitui.ToneWarning, "Renew before 30 March to keep syncing."
	licence.Action, licence.Dismissible = "Renew", true
	vault := kvitui.NewNotice(ui, "This vault could not be saved")
	vault.Tone, vault.Detail = kvitui.ToneDanger, "The disk is full. Nothing has been lost; the changes are still held."
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, licence, vault))
}

func dialogDestructive(ui *kvitui.UI) unison.Paneler {
	// Shown here without its modality so it sits on the page. A real one is
	// opened with Open, modal and in the middle of the window.
	confirm := kvitui.NewDialog(ui, "Delete 40 transactions?")
	confirm.Detail = "They will be removed from every report. This cannot be undone."
	confirm.ConfirmText, confirm.Destructive = "Delete them", true
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(210), kvitui.Centred(kvitui.Width(ui, kvitui.Px(420), confirm))))
}

func floatingRecord(ui *kvitui.UI) unison.Paneler {
	// Shown over a stand-in list, which is what a floating view is always
	// over: the card is centred on the region it floats in, and the region is
	// the list's, not the window's.
	row := func(name, figure string) unison.Paneler {
		r := kvitui.NewSlimRow(ui, name)
		r.Figure = figure
		return kvitui.FullWidth(r)
	}
	list := kvitui.Column(ui, kvitui.Px(0),
		row("Whole Foods", "-84.10"), row("Rent", "-1,850.00"), row("Payroll", "3,204.55"), row("Con Edison", "-96.22"))
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(260), list))
	detail := func(name, phrase string) unison.Paneler {
		r := kvitui.NewSlimRow(ui, name)
		r.Phrase = phrase
		return kvitui.FullWidth(r)
	}
	body := kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.SizeSpace,
		kvitui.FullWidth(kvitui.NewSectionHeading(ui, "Source record")),
		detail("Payee", "Whole Foods"), detail("Category", "Groceries")))
	record := kvitui.NewFloatingView(ui, "Whole Foods · 3 Sep", body)
	record.CloseLabel = "Close record"
	record.OnCloseRequested = record.Close
	whenShown(stage, func() { record.Open(stage) })
	return stage
}
