package main

// The specimens of the settings, forms, tips, history, help and undo.

import (
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

func settingRows(ui *kvitui.UI) unison.Paneler {
	motion := kvitui.NewSwitch(ui, "")
	motion.Checked = true
	theme := kvitui.NewSelect(ui, "", kvitui.Option{Value: "system", Label: "Follow the system"},
		kvitui.Option{Value: "light", Label: "Light"}, kvitui.Option{Value: "dark", Label: "Dark"})
	size := kvitui.NewStepper(ui, "", 10, 24)
	size.Value, size.Unit = 13, "px"
	return kvitui.Width(ui, kvitui.Px(520), kvitui.Column(ui, kvitui.Px(0),
		kvitui.FullWidth(kvitui.NewSettingRow(ui, "Reduce motion", "Panels open without sliding.", motion)),
		kvitui.FullWidth(kvitui.NewSettingRow(ui, "Theme", "", theme)),
		kvitui.FullWidth(kvitui.NewSettingRow(ui, "Interface size", "Every part of the window grows with it.", size))))
}

func settingsAbout(ui *kvitui.UI) unison.Paneler {
	// On Windows and Linux the settings page ends with About; on a Mac About
	// is in the app menu, which opens the same section in a window.
	return kvitui.Width(ui, kvitui.Px(520), kvitui.NewAboutSection(ui, kvitui.About{
		Name: "Kvit Cash", Version: "2.4.1", Details: []string{"Build 2026-10-24"},
		Links: []kvitui.AboutLink{{Text: "Report a problem"}, {Text: "Third-party notices"}, {Text: "Terms of use"}},
	}))
}

// settingsPanes are three panes of an app's settings, each a column of
// setting rows.
func settingsPanes(ui *kvitui.UI) []kvitui.SettingsPane {
	row := func(label, description string) unison.Paneler {
		return kvitui.FullWidth(kvitui.NewSettingRow(ui, label, description, kvitui.NewSwitch(ui, "")))
	}
	return []kvitui.SettingsPane{
		{Name: "General", Symbol: "settings", Content: kvitui.Column(ui, kvitui.Px(0),
			row("Reduce motion", "Panels open without sliding."), row("Open last vault", ""))},
		{Name: "Accounts", Symbol: "bank", Content: kvitui.Column(ui, kvitui.Px(0),
			row("Hide closed accounts", ""), row("Group by bank", ""), row("Show balances", "In the sidebar, beside each account."))},
		{Name: "Privacy", Symbol: "lock", Content: kvitui.Column(ui, kvitui.Px(0),
			row("Lock when idle", "After five minutes without input."))},
	}
}

func settingsPanesOnPage(ui *kvitui.UI) unison.Paneler {
	// The panes as a Mac's settings window holds them, drawn on the page in
	// a ruled frame where the window would be.
	frame := kvitui.NewPanel(ui)
	frame.RuleTop, frame.RuleBottom, frame.RuleLeft, frame.RuleRight = true, true, true, true
	frame.SetLayout(&unison.FlexLayout{Columns: 1})
	frame.AddChild(kvitui.FullWidth(kvitui.NewSettingsPanes(ui, settingsPanes(ui)...)))
	return kvitui.Width(ui, kvitui.Px(560), frame)
}

func settingsPageOnPage(ui *kvitui.UI) unison.Paneler {
	// On Windows and Linux the same panes make a page in the app's window,
	// one scrolling column ending with About, shown here in a ruled frame
	// where the window's body would be.
	page := kvitui.NewSettingsPage(ui, &kvitui.About{Name: "Kvit Cash", Version: "2.4.1"}, settingsPanes(ui)...)
	frame := kvitui.NewPanel(ui)
	frame.RuleTop, frame.RuleBottom, frame.RuleLeft, frame.RuleRight = true, true, true, true
	frame.SetLayout(&unison.FlexLayout{Columns: 1, HAlign: align.Fill, VAlign: align.Fill})
	frame.AddChild(page)
	page.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
	return kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(380), frame))
}

func settingsWindowOpen(ui *kvitui.UI) unison.Paneler {
	// Settings… in a Mac's app menu runs the window's Command, which opens
	// the window there; on Windows and Linux the same command shows the
	// app's settings page. This button opens the window on either.
	window := kvitui.NewSettingsWindow(ui, settingsPanes(ui)...)
	open := kvitui.NewButton(ui, "Open settings window")
	open.OnClick = func() { _ = window.Open() }
	return kvitui.Left(open)
}

func sidebarFooter(ui *kvitui.UI) unison.Paneler {
	// Settings is the last item of the sidebar, at its foot, below Help; a
	// window whose sidebar holds it has no Settings in a Tools menu.
	accounts := kvitui.NewSidebarItem(ui, "Accounts", "wallet")
	accounts.Selected = true
	side := kvitui.NewSidebar(ui, accounts, kvitui.NewSidebarItem(ui, "Budgets", "chart"))
	run := func() {}
	side.AddFooter(kvitui.NewSidebarCommand(ui, &kvitui.Command{Std: kvitui.StdHelp, Run: run}),
		kvitui.NewSidebarCommand(ui, &kvitui.Command{Std: kvitui.StdSettings, Run: run}))
	return kvitui.Height(ui, kvitui.Px(220), side)
}

func formRows(ui *kvitui.UI) unison.Paneler {
	payee := kvitui.NewField(ui)
	payee.SetText("Harlow Market")
	payeeRow := kvitui.NewFormRow(ui, "&Payee", payee)
	payeeRow.Required = true
	memo := kvitui.NewFormRow(ui, "&Memo", kvitui.NewField(ui))
	return kvitui.Width(ui, kvitui.Px(320), kvitui.Column(ui, kvitui.SizeSpaceLoose, kvitui.FullWidth(payeeRow), kvitui.FullWidth(memo)))
}

func formRowFields(ui *kvitui.UI) unison.Paneler {
	// A secret is masked; a field past its limit keeps what was typed and
	// says by how much it is over, rather than cutting it.
	password := kvitui.NewField(ui)
	password.Masked = true
	password.SetText("correct horse")
	code := kvitui.NewField(ui)
	code.MaxLength = 6
	code.SetText("BUKB20201")
	return kvitui.Width(ui, kvitui.Px(320), kvitui.Column(ui, kvitui.SizeSpaceLoose,
		kvitui.FullWidth(kvitui.NewFormRow(ui, "Pass&word", password)),
		kvitui.FullWidth(kvitui.NewFormRow(ui, "Sort &code", code))))
}

func tipBeside(ui *kvitui.UI) unison.Paneler {
	split := kvitui.NewButton(ui, "Split")
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(170), kvitui.Column(ui, kvitui.Px(0),
		kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(120), unison.NewPanel())), kvitui.Left(split))))
	tip := kvitui.NewTip(ui, "gallery.split", "Split a transaction", "Share one payment among several categories.")
	tip.Action = "Try it"
	whenShown(stage, func() {
		// The gallery shows the tip every time; an app shows it once.
		ui.ResetTips()
		tip.ShowBeside(split)
	})
	return stage
}

func historyBack(ui *kvitui.UI) unison.Paneler {
	var tabs []*kvitui.Tab
	var h *kvitui.History
	show := func(place any) {
		for _, t := range tabs {
			t.Selected = t.Text == place
			t.MarkForRedraw()
		}
	}
	h = kvitui.NewHistory(ui, "Accounts", show)
	for _, name := range []string{"Accounts", "Budgets", "Reports"} {
		t := kvitui.NewTab(ui, name)
		t.Selected = name == "Accounts"
		t.OnClick = func() {
			h.Visit(name)
			show(name)
		}
		tabs = append(tabs, t)
	}
	row := []unison.Paneler{kvitui.NewBackButton(ui, h)}
	for _, t := range tabs {
		row = append(row, t)
	}
	return kvitui.Row(ui, kvitui.SizeSpace, row...)
}

func helpLinkTopic(ui *kvitui.UI) unison.Paneler {
	// A link worded as the question its topic answers on Windows and Linux;
	// the round question-mark button on a Mac.
	return kvitui.NewHelpLink(ui, "budgets#rollover", "How does a budget roll over?")
}

func undoStackMenu(ui *kvitui.UI) unison.Paneler {
	stack := kvitui.NewUndoStack(ui)
	stack.Add("Delete transaction", func() {}, func() {})
	undo, redo := stack.Commands()
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(80), unison.NewPanel()))
	whenShown(stage, func() {
		ui.ShowMenuAt(stage, geom.Rect{}, "Edit", []kvitui.MenuItem{{Command: undo}, {Command: redo}})
	})
	return stage
}

func progressHalt(ui *kvitui.UI) unison.Paneler {
	p := kvitui.NewProgress(ui, "Importing statement", 0.4)
	p.HaltText, p.OnHalt = "Stop", func() {}
	p.Pausable = true
	p.Remaining = 3 * time.Minute
	return kvitui.Width(ui, kvitui.Px(480), p)
}

func labelSelectable(ui *kvitui.UI) unison.Paneler {
	// Content the reader may want again can be selected and copied; a
	// control's label cannot.
	iban := kvitui.NewLabel(ui, "GB33 BUKB 2020 1555 5555 55")
	iban.Mono, iban.Selectable = true, true
	return iban
}

func windowDocument(ui *kvitui.UI) unison.Paneler {
	// Source only, as every window is. The title is the document and the
	// app on Windows and Linux, with an asterisk while there are unsaved
	// changes, and the document alone on a Mac; AllowClose asks before the
	// window closes; RememberFrame puts it back where it was left.
	ui.AppName = "Kvit Cash"
	window, err := kvitui.NewWindow(ui, "")
	if err != nil {
		return nil
	}
	window.SetDocument("May budget")
	window.SetEdited(true)
	window.AllowClose = func() bool { return true }
	window.RememberFrame("window.budget")
	window.SetBody(kvitui.NewRegion(ui, kvitui.NewSlimRow(ui, "Groceries")))
	window.ToFront()
	return nil
}
