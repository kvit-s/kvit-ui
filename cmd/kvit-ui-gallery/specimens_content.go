package main

// The Content group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

func dividerBothWays(ui *kvitui.UI) unison.Paneler {
	down := kvitui.NewDivider(ui)
	down.Vertical = true
	down.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Fill, VGrab: true})
	beside := kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewLabel(ui, "left"), down, kvitui.NewLabel(ui, "right"))
	beside.SetLayout(kvitui.AtLeast(ui, kvitui.Px(24), beside.Layout()))
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewDivider(ui), beside))
}

func panelWithRules(ui *kvitui.UI) unison.Paneler {
	p := kvitui.NewPanel(ui)
	p.RuleTop, p.RuleBottom = true, true
	label := kvitui.NewLabel(ui, "A panel")
	label.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Middle, VAlign: align.Middle, HGrab: true, VGrab: true})
	p.AddChild(label)
	p.SetLayout(kvitui.AtLeast(ui, kvitui.Px(60), &unison.FlexLayout{Columns: 1}))
	return kvitui.FullWidth(p)
}

func slimRowMeasured(ui *kvitui.UI) unison.Paneler {
	library := kvitui.NewSlimRow(ui, "kvit-ui")
	library.Symbol, library.Kind, library.Phrase = "folder", "library", "waiting on review"
	library.Figure, library.Unit = "3.5", "d"
	app := kvitui.NewSlimRow(ui, "kvit-cash")
	app.Symbol, app.Kind, app.Phrase = "folder", "application", "not started"
	app.Measured = false
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), library, app))
}

func rowPressableAndStatic(ui *kvitui.UI) unison.Paneler {
	// A row that opens something is reachable by the keyboard, and takes the
	// hover tint, the press and the chevron with it.
	opens := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, "opens the record")))
	opens.Interactive, opens.Label = true, "Groceries"
	// A field name beside its value declares nothing, draws no hover tint and
	// answers no press.
	static := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, "layout only")))
	static.Label = "Amount"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), opens, static))
}

func rowFromTheKeyboard(ui *kvitui.UI) unison.Paneler {
	// Tab reaches this row, the ring says which row the keyboard is on, and
	// Return, Enter or Space opens it. A screen reader's press action does
	// the same one thing.
	groceries := kvitui.NewListRow(ui, kvitui.NewLabel(ui, "keyboard focus — Return opens it"))
	groceries.Form, groceries.Interactive, groceries.Label = kvitui.RowSub, true, "Groceries"
	groceries.Focus()
	// A row is usually a container, and the key may be meant for something
	// inside it. The row answers only while the row itself has the keyboard,
	// so Space on this button presses the button and leaves the record shut.
	split := kvitui.NewButton(ui, "Split")
	split.Form = kvitui.ButtonQuiet
	rent := kvitui.NewListRow(ui, kvitui.FullWidth(kvitui.NewLabel(ui, "the button takes its own Space")), split)
	rent.Form, rent.Interactive, rent.Label = kvitui.RowSub, true, "Rent"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), groceries, rent))
}

func rowHeights(ui *kvitui.UI) unison.Paneler {
	row := func(form kvitui.RowForm, label, note string) *kvitui.ListRow {
		r := kvitui.NewListRow(ui, kvitui.Centred(kvitui.NewLabel(ui, note)))
		r.Form, r.Label = form, label
		return r
	}
	slim := row(kvitui.RowSlim, "Slim", "slim — 30, selected")
	slim.Selected = true
	compact := row(kvitui.RowCompact, "Compact", "compact — 24, keyboard focus")
	compact.Current = true
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0),
		row(kvitui.RowFull, "Full", "full — 56"), row(kvitui.RowSub, "Sub", "sub — 48"), slim, compact))
}

func sectionHeadingForms(ui *kvitui.UI) unison.Paneler {
	mine := kvitui.NewSectionHeading(ui, "Waiting on me")
	mine.Counted, mine.Count = "project", 4
	mine.Action, mine.Collapsible = "Hand all to an agent", true
	// A group of one is still counted. The count is drawn wherever the
	// caller gives one, so it does not come and go as the group changes size.
	sams := kvitui.NewSectionHeading(ui, "Waiting on Sam")
	sams.Counted, sams.Count = "project", 1
	archived := kvitui.NewSectionHeading(ui, "Archived")
	archived.Kind = "closed last quarter"
	archived.Collapsible, archived.Expanded, archived.Strong = true, false, true
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, mine, sams, archived))
}

func sectionHeadingKeyboard(ui *kvitui.UI) unison.Paneler {
	// The ring says the heading has the keyboard; Return, Enter and Space
	// open and close it. The action beside it is a Link, which takes its own
	// keys, so tabbing on to it and pressing Space runs the action and leaves
	// the group where it was.
	mine := kvitui.NewSectionHeading(ui, "Waiting on me")
	mine.Counted, mine.Count = "project", 4
	mine.Action, mine.Collapsible = "Hand all to an agent", true
	mine.Focus()
	archived := kvitui.NewSectionHeading(ui, "Archived")
	archived.Kind = "closed last quarter"
	archived.Collapsible, archived.Expanded = true, false
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, mine, archived))
}

func sectionHeadingWrittenCount(ui *kvitui.UI) unison.Paneler {
	// Some counts are not numbers. A Changes heading reads "1 · +0 −0": one
	// changed file, and the lines added and removed across it, which the
	// service hands over already written and which no integer expresses.
	// CountText is drawn exactly as given, beside the name rather than at the
	// right end where a number goes.
	//
	// ActionSymbol draws the hoisted action as a symbol. The words stay in
	// Action and become the button's accessible name and its tooltip.
	changes := kvitui.NewSectionHeading(ui, "Changes")
	changes.CountText, changes.Collapsible = "1 · +0 −0", true
	changes.Action, changes.ActionSymbol = "Open the diff", "diff"
	agents := kvitui.NewSectionHeading(ui, "Agents")
	agents.CountText, agents.Collapsible = "2 running", true
	agents.Action, agents.ActionSymbol = "Start an agent", "plus"
	terminal := kvitui.NewSectionHeading(ui, "Terminal")
	terminal.Action, terminal.ActionSymbol = "Open a terminal", "terminal"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, changes, agents, terminal))
}

func sectionHeadingExplanation(ui *kvitui.UI) unison.Paneler {
	// ActionExplanation is the sentence beside the action's name, on the
	// action rather than on the heading: a heading is not something a reader
	// presses. Hover either action, or tab to it, to read it; it is the
	// accessible description too.
	agents := kvitui.NewSectionHeading(ui, "Agents")
	agents.CountText, agents.Collapsible = "2 running", true
	agents.Action, agents.ActionSymbol = "Start an agent", "plus"
	agents.ActionExplanation = "Creates an agent at this project’s root. An agent may change files. Nothing runs until you send its first message."
	changes := kvitui.NewSectionHeading(ui, "Changes")
	changes.CountText, changes.Collapsible = "1 · +0 −0", true
	changes.Action = "Review all"
	changes.ActionExplanation = "Opens every changed file as one diff."
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, agents, changes))
}

func cardForms(ui *kvitui.UI) unison.Paneler {
	interactive := kvitui.NewCard(ui, kvitui.NewLabel(ui, "Interactive"))
	interactive.Interactive = true
	selected := kvitui.NewCard(ui, kvitui.NewLabel(ui, "Selected"))
	selected.Selected = true
	return kvitui.Row(ui, kvitui.SizeColumnGap, kvitui.NewCard(ui, kvitui.NewLabel(ui, "A card")), interactive, selected)
}

func paneOpen(ui *kvitui.UI) unison.Paneler {
	pane := kvitui.NewPane(ui, "Transaction", kvitui.Centred(kvitui.NewLabel(ui, "Detail goes here")))
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(160), kvitui.WithPane(ui, nil, pane)))
}

func disclosureOpenAndClosed(ui *kvitui.UI) unison.Paneler {
	changed := kvitui.NewDisclosure(ui, "What changed", kvitui.NewLabel(ui, "Three files were rewritten."))
	changed.Count = 3
	changed.SetExpanded(true)
	unchanged := kvitui.NewDisclosure(ui, "What did not")
	unchanged.Count = 12
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, changed, unchanged))
}

func emptyStateWithAction(ui *kvitui.UI) unison.Paneler {
	empty := kvitui.NewEmptyState(ui, "No transactions yet")
	empty.Symbol = "wallet"
	empty.Detail = "Import a statement or add one by hand, and it will appear here."
	empty.Action = "Import a statement"
	return kvitui.FullWidth(empty)
}

func emptyStateDropTarget(ui *kvitui.UI) unison.Paneler {
	drop := kvitui.NewEmptyState(ui, "Drop a statement here")
	drop.Dashed, drop.Symbol = true, "file-arrow-down"
	drop.Detail = "CSV, OFX and QIF. The file is read on this machine and nothing is sent anywhere."
	drop.Action = "Choose a file"
	return kvitui.FullWidth(drop)
}

func emptyStateCompact(ui *kvitui.UI) unison.Paneler {
	// A column of sections, each of which may have nothing in it. The full
	// block is several times taller than the rows it stands in for, so a
	// stack of them uses the one-line form; the words are the same words.
	changes := kvitui.NewSectionHeading(ui, "Changes")
	changes.Count, changes.Counted = 0, "change"
	noChanges := kvitui.NewEmptyState(ui, "No changes")
	noChanges.Form = kvitui.EmptyCompact
	agents := kvitui.NewSectionHeading(ui, "Agents")
	agents.Count, agents.Counted = 0, "agent"
	noAgents := kvitui.NewEmptyState(ui, "No agents")
	noAgents.Form, noAgents.Symbol = kvitui.EmptyCompact, "robot"
	noAgents.Detail, noAgents.Action = "none started here yet", "Start one"
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceSnug, changes, noChanges, agents, noAgents))
}
