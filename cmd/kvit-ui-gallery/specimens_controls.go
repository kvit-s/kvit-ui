package main

// The Controls group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
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

func fieldStates(ui *kvitui.UI) unison.Paneler {
	payee := kvitui.NewField(ui)
	payee.Label, payee.Placeholder = "Payee", "Who was paid"
	reference := kvitui.NewField(ui)
	reference.Label = "Reference"
	reference.SetText("TX-00173404")
	amount := kvitui.NewField(ui)
	amount.Label, amount.Error = "Amount", "Not a number"
	amount.SetText("twelve")
	locked := kvitui.NewField(ui)
	locked.Label = "Locked"
	locked.SetText("read only")
	locked.SetEnabled(false)
	return kvitui.Width(ui, kvitui.Px(200), kvitui.Column(ui, kvitui.SizeSpaceLoose, payee, reference, amount, locked))
}

func searchFieldStates(ui *kvitui.UI) unison.Paneler {
	filtering := kvitui.NewSearchField(ui)
	filtering.SetText("harlow")
	filtering.Matches, filtering.MatchedNoun = 47, "transaction"
	// What a screen reader is told is a sentence: the digits grouped by the
	// reader's locale and the plural a word, so "250,000 entries" rather than
	// "250000 entry(s)". A noun that does not take an s says its own plural.
	year := kvitui.NewSearchField(ui)
	year.SetText("2026")
	year.Matches, year.MatchedNoun, year.MatchedNounPlural = 250000, "entry", "entries"
	return kvitui.Width(ui, kvitui.Px(220), kvitui.Column(ui, kvitui.SizeSpaceLoose, kvitui.NewSearchField(ui), filtering, year))
}

func stepperInterfaceSize(ui *kvitui.UI) unison.Paneler {
	size := kvitui.NewStepper(ui, "Interface size", tokens.MinInterfaceSize, tokens.MaxInterfaceSize)
	size.Unit = "px"
	size.Follow = ui.Interface.FontSize
	size.OnChange = ui.Interface.SetFontSize
	return size
}

func chipButtonTones(ui *kvitui.UI) unison.Paneler {
	chip := func(words string, tone kvitui.Tone) *kvitui.ChipButton {
		c := kvitui.NewChipButton(ui, words)
		c.Tone = tone
		return c
	}
	settled := chip("settled", kvitui.ToneSuccess)
	settled.Strong = true
	disputed := chip("disputed", kvitui.ToneDanger)
	disputed.Strong, disputed.Symbol = true, "warning"
	focused := chip("keyboard focus", kvitui.ToneNeutral)
	focused.Focus()
	return kvitui.Column(ui, kvitui.SizeSpace,
		kvitui.Left(kvitui.Row(ui, kvitui.SizeSpaceNear, chip("neutral", kvitui.ToneNeutral), chip("accent", kvitui.ToneAccent),
			chip("success", kvitui.ToneSuccess), chip("warning", kvitui.ToneWarning), chip("danger", kvitui.ToneDanger),
			chip("info", kvitui.ToneInfo))),
		kvitui.Left(kvitui.Row(ui, kvitui.SizeSpaceNear, settled, disputed, focused)))
}

func chipButtonLeadsTo(ui *kvitui.UI) unison.Paneler {
	// No chevron is drawn for a chip that acts. Whether pressing it discloses
	// a list under the row, leaves the view or opens another window belongs
	// to the destination, which is why the symbol is the caller's to name.
	changes := kvitui.NewChipButton(ui, "3 changes")
	changes.Symbol = "file"
	ahead := kvitui.NewChipButton(ui, "1 ahead")
	ahead.TrailingSymbol = "chevron-right"
	behind := kvitui.NewChipButton(ui, "2 behind")
	behind.TrailingSymbol = "chevron-right"
	hub := kvitui.NewChipButton(ui, "Open on the hub")
	hub.TrailingSymbol = "external"
	return kvitui.Row(ui, kvitui.SizeSpaceNear, changes, ahead, behind, hub)
}

func chipButtonUnavailable(ui *kvitui.UI) unison.Paneler {
	// A chip with a reason keeps its place in the tab order, still shows the
	// words in its tooltip and its accessible description, and does nothing
	// when pressed. The fill going away rather than changing hue is what a
	// reader who cannot separate the tones still sees.
	behind := kvitui.NewChipButton(ui, "2 behind")
	behind.UnavailableReason = "The other branch has not been fetched yet."
	inSync := kvitui.NewChipButton(ui, "in sync")
	inSync.Tone, inSync.UnavailableReason = kvitui.ToneSuccess, "There is nothing on either side to compare."
	return kvitui.Row(ui, kvitui.SizeSpaceNear, kvitui.NewChipButton(ui, "1 ahead"), behind, inSync)
}

func chipButtonElided(ui *kvitui.UI) unison.Paneler {
	// The label gives way and the symbols keep their size: a symbol at half
	// width is a smudge, and the trailing one says where pressing this goes.
	// The whole label is in the tooltip once it no longer fits.
	long := kvitui.NewChipButton(ui, "A fact whose whole phrase does not fit in this column")
	long.Symbol, long.TrailingSymbol = "warning", "chevron-right"
	short := kvitui.NewChipButton(ui, "A fact whose whole phrase does not fit in this column")
	return kvitui.Column(ui, kvitui.SizeSpaceSnug,
		kvitui.Width(ui, kvitui.Px(150), long), kvitui.Width(ui, kvitui.Px(90), short))
}

func chipButtonCurrent(ui *kvitui.UI) unison.Paneler {
	// One chip in a row of places to go is where the reader already is:
	// the selection tint under it, the accent on its edge, and bold words,
	// since the tints are nothing in a grayscale screenshot.
	changes := kvitui.NewChipButton(ui, "3 changes")
	changes.Symbol = "diff"
	ahead := kvitui.NewChipButton(ui, "1 ahead")
	ahead.Current, ahead.TrailingSymbol = true, "chevron-right"
	behind := kvitui.NewChipButton(ui, "2 behind")
	behind.TrailingSymbol = "chevron-right"
	return kvitui.Row(ui, kvitui.SizeSpaceNear, changes, ahead, behind)
}

func chipButtonExplanation(ui *kvitui.UI) unison.Paneler {
	// "1 ahead" is a count; what it opens is the one commit. Explanation is
	// that sentence; UnavailableReason is the other one, about a chip that
	// cannot be pressed at all.
	ahead := kvitui.NewChipButton(ui, "1 ahead")
	ahead.TrailingSymbol = "chevron-right"
	ahead.Explanation = "Opens the one commit this branch has and the remote does not."
	ahead.Focus()
	behind := kvitui.NewChipButton(ui, "2 behind")
	behind.UnavailableReason = "The remote has not been fetched yet."
	return kvitui.Row(ui, kvitui.SizeSpaceNear, ahead, behind)
}

func textAreaStates(ui *kvitui.UI) unison.Paneler {
	message := kvitui.NewTextArea(ui)
	message.Label, message.Placeholder = "Message", "What happened"
	source := kvitui.NewTextArea(ui)
	source.Label, source.Mono, source.ReadOnly = "Source", true, true
	source.SetText("# Notes\n\nThe file as it stands.")
	summary := kvitui.NewTextArea(ui)
	summary.Label, summary.Error = "Summary", "Say what changed"
	summary.SetText("...")
	return kvitui.Width(ui, kvitui.Px(280), kvitui.Column(ui, kvitui.SizeSpaceLoose, message, source, summary))
}

func textAreaPlain(ui *kvitui.UI) unison.Paneler {
	source := kvitui.NewTextArea(ui)
	source.Label, source.Plain, source.Mono, source.ReadOnly = "Source", true, true, true
	source.SetText("# Notes\n\nA document has no edges to find:\nthe pane is its edge.")
	return kvitui.Width(ui, kvitui.Px(280), source)
}

func textAreaUnderlay(ui *kvitui.UI) unison.Paneler {
	marked := kvitui.NewTextArea(ui)
	marked.Label, marked.Mono, marked.ReadOnly = "Reviewed source", true, true
	marked.SetText("one\ntwo\nthree")
	// A wash behind the second line, placed where the text says the fifth
	// character is.
	marked.Underlay = func(gc *unison.Canvas, at func(index int) geom.Rect) {
		box := at(4)
		box.Width = float32(ui.Interface.Px(40))
		gc.DrawRect(box, kvitui.Color(ui.Theme.Tokens().SelectionTint).Paint(gc, box, paintstyle.Fill))
	}
	return kvitui.Width(ui, kvitui.Px(280), marked)
}

func checkStates(ui *kvitui.UI) unison.Paneler {
	drafts := kvitui.NewCheck(ui, "Include drafts")
	drafts.Checked = true
	some := kvitui.NewCheck(ui, "Some of these")
	some.Partial = true
	unavailable := kvitui.NewCheck(ui, "Not available")
	unavailable.SetEnabled(false)
	return kvitui.Column(ui, kvitui.SizeSpaceNear,
		kvitui.Left(kvitui.NewCheck(ui, "Include archived")), kvitui.Left(drafts), kvitui.Left(some), kvitui.Left(unavailable))
}

func selectCurrency(ui *kvitui.UI) unison.Paneler {
	var options []kvitui.Option
	for _, code := range []string{"GBP", "EUR", "USD", "JPY", "CHF", "SEK"} {
		options = append(options, kvitui.Option{Value: code, Label: code})
	}
	return kvitui.NewSelect(ui, "Currency", options...)
}
