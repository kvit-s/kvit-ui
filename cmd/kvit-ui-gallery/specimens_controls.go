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

func buttonForms(ui *kvitui.UI) unison.Paneler {
	row := func(label, symbol string, busy bool) *unison.Panel {
		var buttons []unison.Paneler
		for _, form := range []kvitui.ButtonForm{kvitui.ButtonPrimary, kvitui.ButtonOrdinary, kvitui.ButtonQuiet} {
			b := kvitui.NewButton(ui, label)
			b.Form, b.Symbol, b.Busy = form, symbol, busy
			buttons = append(buttons, b)
		}
		return kvitui.Row(ui, kvitui.SizeSpace, buttons...)
	}
	return kvitui.Column(ui, kvitui.SizeSpace, row("Save", "", false), row("Schedule", "calendar", false), row("Save", "", true))
}

func buttonDanger(ui *kvitui.UI) unison.Paneler {
	primary := kvitui.NewButton(ui, "Delete")
	primary.Form, primary.Danger = kvitui.ButtonPrimary, true
	ordinary := kvitui.NewButton(ui, "Delete")
	ordinary.Danger, ordinary.Symbol = true, "trash"
	disabled := kvitui.NewButton(ui, "Disabled")
	disabled.SetEnabled(false)
	return kvitui.Row(ui, kvitui.SizeSpace, primary, ordinary, disabled)
}

func buttonChecked(ui *kvitui.UI) unison.Paneler {
	on := kvitui.NewButton(ui, "Select region")
	on.Checkable, on.Checked = true, true
	off := kvitui.NewButton(ui, "Select region")
	off.Checkable = true
	return kvitui.Row(ui, kvitui.SizeSpace, on, off)
}

func buttonExplanations(ui *kvitui.UI) unison.Paneler {
	// The reason a control is in the state it is in often lives somewhere the
	// reader cannot see. A disabled button with nothing to say about itself
	// is a grey rectangle and no account of it; Explanation is where that
	// sentence goes, shown on hover and announced as the accessible
	// description. It is read on a disabled button too, which is the case it
	// exists for.
	pull := kvitui.NewButton(ui, "Pull")
	pull.Explanation = "Brings the 2 commits on the remote into this branch."
	push := kvitui.NewButton(ui, "Push")
	push.SetEnabled(false)
	push.Explanation = "Nothing here has been committed yet."
	archive := kvitui.NewButton(ui, "Archive")
	archive.SetEnabled(false)
	archive.Danger = true
	archive.Explanation = "The branch has work that has not been pushed."
	return kvitui.Row(ui, kvitui.SizeSpace, pull, push, archive)
}
