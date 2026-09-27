package main

// The Flow group's specimens.

import (
	"fmt"
	"math"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
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

// whenShown runs f once, the first time a panel is laid out in a window.
func whenShown(p *unison.Panel, f func()) {
	done := false
	p.FrameChangeCallback = func() {
		if !done && p.Window() != nil {
			done = true
			unison.InvokeTask(f)
		}
	}
}

// command is a shortcut with the platform's command key.
func command(key unison.KeyCode) unison.KeyBinding {
	return unison.KeyBinding{KeyCode: key, Modifiers: mod.OSMenuCommand()}
}

func menuOpened(ui *kvitui.UI) unison.Paneler {
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(170), unison.NewPanel()))
	// Opened at the stage's top left once it is in a window, as the Qt page
	// opens its menu.
	whenShown(stage, func() {
		ui.ShowMenuAt(stage, geom.Rect{}, "", []kvitui.MenuItem{
			{Text: "Open", Symbol: "file", Key: command(unison.KeyO)},
			{Text: "Duplicate", Symbol: "copy", Key: command(unison.KeyD)},
			{Text: "Archive", Symbol: "archive"},
			{Separator: true},
			{Text: "Delete", Symbol: "trash", Danger: true},
		})
	})
	return stage
}

func menuItemForms(ui *kvitui.UI) unison.Paneler {
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(110), unison.NewPanel()))
	whenShown(stage, func() {
		ui.ShowMenuAt(stage, geom.Rect{}, "", []kvitui.MenuItem{
			{Text: "Reconcile", Symbol: "check", Key: command(unison.KeyR)},
			{Text: "Split", Symbol: "split", Disabled: true},
			{Text: "Delete", Symbol: "trash", Danger: true},
		})
	})
	return stage
}

func switchForms(ui *kvitui.UI) unison.Paneler {
	follow := kvitui.NewSwitch(ui, "Follow the system theme")
	follow.Checked = true
	cellular := kvitui.NewSwitch(ui, "Sync over cellular")
	cellular.SetEnabled(false)
	return kvitui.Column(ui, kvitui.SizeSpaceNear,
		kvitui.Left(follow), kvitui.Left(kvitui.NewSwitch(ui, "Reduce motion")), kvitui.Left(cellular))
}

func radioAppearance(ui *kvitui.UI) unison.Paneler {
	appearance := kvitui.NewRadioGroup(ui, "Appearance",
		kvitui.RadioOption{Value: "system", Label: "Follow the system",
			Detail: "Light or dark, whichever the desktop is set to."},
		kvitui.RadioOption{Value: "light", Label: "Always light"},
		kvitui.RadioOption{Value: "dark", Label: "Always dark"})
	appearance.Current = "system"
	return kvitui.FullWidth(appearance)
}

func progressForms(ui *kvitui.UI) unison.Paneler {
	reconciling := kvitui.NewProgress(ui, "Reconciling", 0)
	reconciling.Determinate = false
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose,
		kvitui.FullWidth(kvitui.NewProgress(ui, "Importing statement", 0.47)), kvitui.FullWidth(reconciling)))
}

func sliderWithValue(ui *kvitui.UI) unison.Paneler {
	opacity := kvitui.NewSlider(ui, "Opacity", 0.6)
	opacity.Unit = "%"
	return kvitui.Width(ui, kvitui.Px(200), opacity)
}

func splitTwoPanes(ui *kvitui.UI) unison.Paneler {
	pane := func(words string) unison.Paneler {
		p := kvitui.NewPanel(ui)
		p.SetLayout(&unison.FlexLayout{Columns: 1})
		p.AddChild(kvitui.Centred(kvitui.NewLabel(ui, words)))
		return p
	}
	split := kvitui.NewSplitView(ui, pane("left"), pane("right"))
	split.SetSize(0, kvitui.Px(160))
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(120), split))
}

func numberFieldForms(ui *kvitui.UI) unison.Paneler {
	field := func(label, text string, minimum, maximum float64, decimals int) unison.Paneler {
		f := kvitui.NewNumberField(ui)
		f.Label, f.Minimum, f.Maximum, f.Decimals = label, minimum, maximum, decimals
		f.SetText(text)
		return kvitui.Left(f)
	}
	unbounded := math.Inf(1)
	return kvitui.Column(ui, kvitui.SizeSpaceLoose,
		field("Days", "14", 1, 365, 0),
		field("Rate", "4.25", -unbounded, unbounded, 2),
		field("Days", "999", 1, 365, 0))
}

func moneyFieldCurrencies(ui *kvitui.UI) unison.Paneler {
	amount := func(text, currency string, minorDigits int) unison.Paneler {
		f := kvitui.NewMoneyField(ui, currency, minorDigits)
		f.Label = "Amount"
		f.SetText(text)
		return kvitui.Left(f)
	}
	return kvitui.Column(ui, kvitui.SizeSpaceLoose,
		amount("1284.50", "GBP", 2), amount("4200", "JPY", 0), amount("18.750", "BHD", 3))
}

func confirmAfterBulkEdit(ui *kvitui.UI) unison.Paneler {
	done := kvitui.NewConfirmInPlace(ui, "Recategorised as Groceries")
	done.Shown, done.Affected = true, 40
	return kvitui.FullWidth(done)
}

func timelineHistory(ui *kvitui.UI) unison.Paneler {
	return kvitui.FullWidth(kvitui.NewTimeline(ui, "Account history",
		kvitui.TimelineEntry{When: "14:02", What: "Statement imported", Who: "agent",
			Detail: "412 transactions, 8 unmatched", Tone: kvitui.ToneSuccess},
		kvitui.TimelineEntry{When: "11:20", What: "Two rows disputed", Who: "you", Tone: kvitui.ToneWarning},
		kvitui.TimelineEntry{When: "Yesterday", What: "Account opened", Who: "you"}))
}

func spotlightOnButton(ui *kvitui.UI) unison.Paneler {
	target := kvitui.NewButton(ui, "Import a statement")
	target.Form = kvitui.ButtonPrimary
	stage := kvitui.FullWidth(kvitui.At(ui, kvitui.Px(40), kvitui.Px(20), kvitui.Px(160), target))
	spot := kvitui.NewSpotlight(ui, target, "Start here", "Import a statement and the dashboard fills itself in.")
	spot.OnDismiss = spot.Close
	whenShown(stage, func() { spot.Open(stage) })
	return stage
}

func beforeAfterForms(ui *kvitui.UI) unison.Paneler {
	amount := kvitui.NewBeforeAfter(ui, "Amount", "42.00", "44.50")
	amount.Unit = "GBP"
	category := kvitui.NewBeforeAfter(ui, "Category", "", "Groceries")
	category.BeforeMeasured = false
	payee := kvitui.NewBeforeAfter(ui, "Payee", "TESCO 4471", "")
	payee.AfterMeasured = false
	date := kvitui.NewBeforeAfter(ui, "Date", "2026-08-14", "2026-08-14")
	balance := kvitui.NewBeforeAfter(ui, "", "1,284.50", "1,301.75")
	balance.Unit, balance.Role = "GBP", kvitui.RoleStrong
	return kvitui.Column(ui, kvitui.SizeSpaceNear,
		kvitui.Left(amount), kvitui.Left(category), kvitui.Left(payee), kvitui.Left(date), kvitui.Left(balance))
}

func treeAccounts(ui *kvitui.UI) unison.Paneler {
	leaf := func(label string) kvitui.TreeNode { return kvitui.TreeNode{Label: label} }
	accounts := kvitui.NewTree(ui, "Accounts",
		kvitui.TreeNode{Label: "Everyday", Children: []kvitui.TreeNode{leaf("Checking"), leaf("Joint checking"), leaf("Cash")}},
		kvitui.TreeNode{Label: "Savings", Children: []kvitui.TreeNode{
			leaf("Emergency fund"),
			{Label: "Certificates", Children: []kvitui.TreeNode{leaf("18 months"), leaf("3 years")}}}},
		leaf("Credit card"))
	// The top level open, and everything under it closed.
	accounts.ExpandTo(1)
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(250), accounts))
}

func typeAheadCategory(ui *kvitui.UI) unison.Paneler {
	var categories []kvitui.Suggestion
	for _, c := range []string{"Groceries", "Transport", "Utilities", "Rent", "Leisure", "Health", "Savings", "Income"} {
		categories = append(categories, kvitui.Suggestion{Value: c})
	}
	// A category picker that will not invent categories: a typo would make
	// a second category beside the right one.
	category := kvitui.NewTypeAhead(ui, "Category", false, categories...)
	category.Placeholder = "Start typing"
	return kvitui.Width(ui, kvitui.Px(240), category)
}

func dualListColumns(ui *kvitui.UI) unison.Paneler {
	columns := func(names ...string) []kvitui.Option {
		var out []kvitui.Option
		for _, n := range names {
			out = append(out, kvitui.Option{Value: n, Label: n})
		}
		return out
	}
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(220), kvitui.NewDualList(ui,
		columns("Payee", "Account", "Tags", "Note"), columns("Date", "Description", "Amount", "Balance"))))
}

func menuItemSubmenu(ui *kvitui.UI) unison.Paneler {
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(80), unison.NewPanel()))
	// A line with items of its own opens them beside it, and draws a
	// chevron where a shortcut would go.
	whenShown(stage, func() {
		ui.ShowMenuAt(stage, geom.Rect{}, "", []kvitui.MenuItem{
			{Text: "Turn into", Symbol: "rename", Items: []kvitui.MenuItem{{Text: "Heading"}, {Text: "Quote"}}},
			{Text: "Copy as", Symbol: "copy", Items: []kvitui.MenuItem{{Text: "Markdown"}, {Text: "Plain text"}}},
		})
	})
	return stage
}
