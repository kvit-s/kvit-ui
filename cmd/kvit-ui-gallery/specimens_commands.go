package main

// The Commands group's specimens.

import (
	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

func commandEverywhere(ui *kvitui.UI) unison.Paneler {
	export := &kvitui.Command{Name: "Export statement", Asks: true, Symbol: "external",
		Shortcut: kvitui.Keys(unison.KeyE, kvitui.ModPrimary|kvitui.ModShift), Run: func() {}}
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(60), unison.NewPanel()))
	whenShown(stage, func() {
		ui.ShowMenuAt(stage, geom.Rect{}, "Statement", []kvitui.MenuItem{{Command: export}})
	})
	return kvitui.Column(ui, kvitui.SizeSpace,
		kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewCommandButton(ui, export), kvitui.NewCommandIconButton(ui, export)), stage)
}

func commandStandard(ui *kvitui.UI) unison.Paneler {
	run := func() {}
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(150), unison.NewPanel()))
	whenShown(stage, func() {
		ui.ShowMenuAt(stage, geom.Rect{}, "Edit", []kvitui.MenuItem{
			{Command: &kvitui.Command{Std: kvitui.StdUndo, Run: run}},
			{Command: &kvitui.Command{Std: kvitui.StdRedo, Run: run}},
			{Separator: true},
			{Command: &kvitui.Command{Std: kvitui.StdSelectAll, Run: run}},
			{Command: &kvitui.Command{Std: kvitui.StdFind, Run: run}},
		})
	})
	return stage
}

// windowSurfaces are the commands of an imagined statement window.
func windowSurfaces(ui *kvitui.UI) (*kvitui.MenuBar, *kvitui.Toolbar) {
	run := func() {}
	save := &kvitui.Command{Std: kvitui.StdSave, Run: run}
	undo := &kvitui.Command{Std: kvitui.StdUndo, Run: run}
	find := &kvitui.Command{Std: kvitui.StdFind, Run: run}
	settings := &kvitui.Command{Std: kvitui.StdSettings, Run: run}
	help := &kvitui.Command{Std: kvitui.StdHelp, Run: run}
	quit := &kvitui.Command{Std: kvitui.StdQuit, Run: run}
	reconcile := &kvitui.Command{Name: "Reconcile", Symbol: "check", Run: run}
	importStatement := &kvitui.Command{Name: "Import", Symbol: "arrow-down", Asks: true, Run: run}
	shown := true
	sidebar := &kvitui.Command{Name: "Sidebar", TurnOn: "Show sidebar", TurnOff: "Hide sidebar",
		Checked: func() bool { return shown }, Run: func() { shown = !shown }}
	bar := kvitui.NewMenuBar(ui).
		Add(save, undo, find, settings, help, quit).
		AddTo(kvitui.MenuView, kvitui.MenuItem{Command: sidebar}).
		AddMenu("&Statement", kvitui.MenuItem{Command: importStatement}, kvitui.MenuItem{Command: reconcile})
	tools := kvitui.NewToolbar(ui, "Statement").
		Add(kvitui.ToolLeading, importStatement, reconcile).
		Add(kvitui.ToolTrailing, find, settings)
	return bar, tools
}

func menuBarWindow(ui *kvitui.UI) unison.Paneler {
	bar, tools := windowSurfaces(ui)
	body := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(40), unison.NewPanel()))
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpace, kvitui.FullWidth(bar), kvitui.FullWidth(tools), body))
}

func menuBarOpened(ui *kvitui.UI) unison.Paneler {
	bar, _ := windowSurfaces(ui)
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(150), unison.NewPanel()))
	// The View menu, opened once the bar is in a window: the window's own
	// parts first, then the standard lines.
	whenShown(stage, func() {
		for i, m := range bar.Menus() {
			if m.Title == "&View" || m.Title == "View" {
				bar.OpenMenu(i)
			}
		}
	})
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.Px(0), kvitui.FullWidth(bar), stage))
}

func toolbarWideAndNarrow(ui *kvitui.UI) unison.Paneler {
	_, wide := windowSurfaces(ui)
	_, narrow := windowSurfaces(ui)
	return kvitui.FullWidth(kvitui.Column(ui, kvitui.SizeSpaceLoose,
		kvitui.FullWidth(wide), kvitui.Width(ui, kvitui.Px(150), narrow)))
}

func menuButtonRefinements(ui *kvitui.UI) unison.Paneler {
	run := func() {}
	sort := kvitui.NewMenuButton(ui, "Sort",
		kvitui.MenuItem{Text: "By date", Radio: true, Checked: true, OnSelect: run},
		kvitui.MenuItem{Text: "By payee", Radio: true, OnSelect: run},
		kvitui.MenuItem{Text: "By amount", Radio: true, OnSelect: run})
	return sort
}

func splitButtonVariants(ui *kvitui.UI) unison.Paneler {
	run := func() {}
	split := kvitui.NewSplitButton(ui,
		&kvitui.Command{Name: "Export as CSV", Run: run},
		&kvitui.Command{Name: "Export as PDF", Run: run},
		&kvitui.Command{Name: "Export as OFX", Run: run})
	split.MenuLabel = "Choose format"
	return split
}

func toggleSplitButtonStyles(ui *kvitui.UI) unison.Paneler {
	// The main part turns the bullets on and off, and shows them on here;
	// the menu picks how they look, and choosing a style turns them on.
	on := true
	bullets := &kvitui.Command{Name: "Bullets", Symbol: "list-bullets",
		Checked: func() bool { return on }, Run: func() { on = !on }}
	run := func() {}
	return kvitui.NewToggleSplitButton(ui, bullets,
		&kvitui.Command{Name: "Dot bullets", Run: run},
		&kvitui.Command{Name: "Dash bullets", Run: run},
		&kvitui.Command{Name: "Arrow bullets", Run: run})
}

func menuItemStates(ui *kvitui.UI) unison.Paneler {
	stage := kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(170), unison.NewPanel()))
	whenShown(stage, func() {
		ui.ShowMenuAt(stage, geom.Rect{}, "View", []kvitui.MenuItem{
			{Text: "Status bar", TurnOn: "Show status bar", TurnOff: "Hide status bar", Checked: true},
			{Separator: true},
			{Text: "Compact rows", Radio: true},
			{Text: "Comfortable rows", Radio: true, Checked: true},
			{Separator: true},
			{Text: "Bold", Mixed: true},
		})
	})
	return stage
}
