package main

// The Foundation group's specimens. Each function is one specimen, and the
// gallery shows its source as the specimen's code sample, so a sample is
// always code that compiles.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/unison"
)

func labelRoles(ui *kvitui.UI) unison.Paneler {
	role := func(r kvitui.TypeRole, s string) *kvitui.Label {
		l := kvitui.NewLabel(ui, s)
		l.Role = r
		return l
	}
	return kvitui.Column(ui, kvitui.SizeSpaceSnug,
		role(kvitui.RoleDisplay, "display — a page title"),
		role(kvitui.RoleHeadline, "headline — a pane title"),
		role(kvitui.RoleTitle, "title — a section heading"),
		role(kvitui.RoleStrong, "strong — a name"),
		role(kvitui.RoleBody, "body — row text and prose"),
		role(kvitui.RoleSmall, "small — chip labels and sub-lines"),
		role(kvitui.RoleCaption, "caption — kind tags and counts"),
	)
}

func labelMonoAndTabular(ui *kvitui.UI) unison.Paneler {
	id := kvitui.NewLabel(ui, "a1b2c3d4")
	id.Mono = true
	amount := kvitui.NewLabel(ui, "1,234.56")
	amount.Tabular = true
	return kvitui.Column(ui, kvitui.SizeSpaceSnug, id, amount)
}

func iconSizes(ui *kvitui.UI) unison.Paneler {
	search := kvitui.NewIcon(ui, "search")
	search.Size = kvitui.SizeIconSizeSmall
	trash := kvitui.NewIcon(ui, "trash")
	trash.Ink = kvitui.InkDanger
	success := kvitui.NewIcon(ui, "success")
	success.Ink = kvitui.InkSuccess
	return kvitui.Row(ui, kvitui.SizeSpace, search, kvitui.NewIcon(ui, "chevron-right"), trash, success)
}

func iconUnknown(ui *kvitui.UI) unison.Paneler {
	return kvitui.NewIcon(ui, "not-a-symbol")
}

func iconButtonForms(ui *kvitui.UI) unison.Paneler {
	remove := kvitui.NewIconButton(ui, "trash", "Delete")
	remove.Form = kvitui.Ordinary
	pin := kvitui.NewIconButton(ui, "pin", "Pin")
	pin.Checked = true
	duplicate := kvitui.NewIconButton(ui, "copy", "Duplicate")
	duplicate.Explanation = "Only a saved report can be duplicated."
	duplicate.SetEnabled(false)
	settings := kvitui.NewIconButton(ui, "settings", "Settings")
	settings.Focus()
	return kvitui.Row(ui, kvitui.SizeSpace,
		kvitui.NewIconButton(ui, "pencil", "Edit"), remove, pin, duplicate, settings)
}

func iconButtonSizes(ui *kvitui.UI) unison.Paneler {
	// A symbol standing on its own is drawn at the full icon size. Dense
	// draws it at the small one, the size every symbol beside words in this
	// library is drawn at: a strip of buttons across a header, or on a
	// heading bar, reads as larger than its row at the full size.
	strip := func(dense bool) *unison.Panel {
		search := kvitui.NewIconButton(ui, "search", "Search")
		add := kvitui.NewIconButton(ui, "plus", "New track")
		branch := kvitui.NewIconButton(ui, "git", "Switch branch")
		search.Dense, add.Dense, branch.Dense = dense, dense, dense
		return kvitui.Row(ui, kvitui.SizeSpace, search, add, branch)
	}
	return kvitui.Column(ui, kvitui.SizeSpace, strip(false), strip(true))
}

func iconButtonExplanation(ui *kvitui.UI) unison.Paneler {
	// Label names the button and is what a screen reader is told.
	// Explanation is the sentence beside it, for what the name cannot hold:
	// what pressing this changes, or why it cannot be pressed. Both go to
	// the tooltip and to the accessible description, in that order.
	rename := kvitui.NewIconButton(ui, "rename", "Rename")
	rename.Explanation = "Renames the track and its branch. Its history and folder stay unchanged."
	rename.Focus()
	archive := kvitui.NewIconButton(ui, "archive", "Archive")
	archive.Explanation = "The branch has work that has not been pushed."
	archive.SetEnabled(false)
	return kvitui.Row(ui, kvitui.SizeSpace, rename, archive)
}

func linkForms(ui *kvitui.UI) unison.Paneler {
	calendar := kvitui.NewLink(ui, "Open calendar")
	calendar.Symbol = "calendar"
	focused := kvitui.NewLink(ui, "Keyboard focus")
	focused.Focus()
	return kvitui.Row(ui, kvitui.SizeSpaceLoose, kvitui.NewLink(ui, "Privacy policy"), calendar, focused)
}

func linkElided(ui *kvitui.UI) unison.Paneler {
	return kvitui.Width(ui, kvitui.Px(120), kvitui.NewLink(ui, "A destination whose full name does not fit here"))
}

func labelCutShort(ui *kvitui.UI) unison.Paneler {
	path := "projects/ledger/exports/2026/quarterly-report-final.pdf"
	end := kvitui.NewLabel(ui, path)
	middle := kvitui.NewLabel(ui, path)
	// The end of a file's name is what tells it from its neighbours.
	middle.Middle = true
	two := kvitui.NewLabel(ui, "A wrapping label kept to two lines is cut short at the end of the second, so a long message takes the same room in every row of a list however much it says.")
	two.Wrap, two.MaxLines, two.Ink = true, 2, kvitui.InkTextSecondary
	return kvitui.Width(ui, kvitui.Px(260), kvitui.Stack(ui, kvitui.SizeSpaceSnug, end, middle, two))
}

func labelItalicAndAligned(ui *kvitui.UI) unison.Paneler {
	aside := kvitui.NewLabel(ui, "An aside, in italics")
	aside.Italic = true
	centred := kvitui.NewLabel(ui, "In the middle")
	centred.Align = text.AlignMiddle
	total := kvitui.NewLabel(ui, "At the end")
	total.Align = text.AlignEnd
	id := kvitui.NewLabel(ui, "a1b2c3d4, semibold")
	id.Mono, id.WeightedMono, id.Weight = true, true, text.Semibold
	return kvitui.Width(ui, kvitui.Px(260), kvitui.Stack(ui, kvitui.SizeSpaceSnug, aside, centred, total, id))
}

func labelForScreenReaders(ui *kvitui.UI) unison.Paneler {
	// A heading by role, a count with its meaning spelt out, and a line that
	// takes no room until it has something to say.
	heading := kvitui.NewLabel(ui, "Changes")
	heading.Role, heading.Heading = kvitui.RoleTitle, true
	count := kvitui.NewLabel(ui, "3 / 12")
	count.Name, count.Tabular = "3 of 12 files reviewed", true
	status := kvitui.NewLabel(ui, "")
	status.EmptyTakesNoRoom = true
	save := kvitui.NewButton(ui, "Save changes")
	save.OnClick = func() { status.SetText("Saved just now") }
	return kvitui.Stack(ui, kvitui.SizeSpaceSnug, heading, count, status, kvitui.OwnWidth(save))
}
