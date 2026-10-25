package kvitui_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/platform"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// A setting row names its switch, which says On or Off on the Windows form;
// Settings at the foot of the sidebar keeps it out of Tools (WIN-16, CTL-6).
func TestSettingsLiveInRowsAndTheSidebarsFoot(t *testing.T) {
	var sw *kvitui.Switch
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		sw = kvitui.NewSwitch(ui, "")
		return []unison.Paneler{kvitui.FullWidth(kvitui.NewSettingRow(ui, "Reduce motion", "Fewer moving parts.", sw))}
	})
	if n := screen.AccessibilityNodeFor(sw); n == nil || n.Name != "Reduce motion" {
		t.Errorf("the switch in a setting row is %+v", n)
	}
	settings := &kvitui.Command{Std: kvitui.StdSettings, Run: func() {}}
	var side *kvitui.Sidebar
	var gear *kvitui.SidebarItem
	screen.Do(func() {
		side = kvitui.NewSidebar(ui, kvitui.NewSidebarItem(ui, "Accounts", "wallet"))
		gear = kvitui.NewSidebarCommand(ui, settings)
		side.AddFooter(gear)
		w.SetSidebar(side)
		w.SetMenuBar(kvitui.NewMenuBar(ui).Add(settings, &kvitui.Command{Std: kvitui.StdQuit, Run: func() {}}))
	})
	screen.Sync()
	screen.Do(func() {
		if bottom, foot := side.ContentRect(false).Bottom(), gear.FrameRect().Bottom(); foot != bottom {
			t.Errorf("Settings ends at %.1f, not the sidebar's foot at %.1f", foot, bottom)
		}
		for _, m := range w.MenuBar().Menus() {
			if m.Title == "&Tools" {
				t.Error("the menu bar has a Tools menu for a Settings the sidebar holds")
			}
		}
		if off := w.OffSurface(); len(off) != 0 {
			t.Errorf("commands on no surface: %v", off)
		}
	})
}

// On the macOS form Settings opens a window of panes, titled with the pane
// shown and as tall as it, switched by the pointer or the arrow keys and
// closed by Command-W; on the Windows form the command shows the app's page
// instead (WIN-16).
func TestSettingsOpenAWindowOfPanesOnAMac(t *testing.T) {
	screen, ui, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler { return nil })
	rows := func(names ...string) unison.Paneler {
		var out []unison.Paneler
		for _, n := range names {
			out = append(out, kvitui.FullWidth(kvitui.NewSettingRow(ui, n, "", kvitui.NewSwitch(ui, ""))))
		}
		return kvitui.Column(ui, kvitui.SizeSpace, out...)
	}
	var settings *kvitui.SettingsWindow
	var shown []int
	page := 0
	screen.Do(func() {
		settings = kvitui.NewSettingsWindow(ui,
			kvitui.SettingsPane{Name: "General", Symbol: "settings", Content: rows("Reduce motion")},
			kvitui.SettingsPane{Name: "Accounts", Symbol: "bank", Content: rows("Hide closed accounts", "Group by bank", "Show balances")},
			kvitui.SettingsPane{Name: "Privacy", Symbol: "lock", Content: rows("Lock after a minute")})
		settings.Panes.OnShow = func(pane int) { shown = append(shown, pane) }
		cmd := settings.Command(func() { page++ })
		cmd.Run()
		if page != 1 || settings.Window() != nil {
			t.Errorf("on the Windows form Settings showed the page %d times and opened %v", page, settings.Window())
		}
		ui.SetForm(kvitui.FormMac)
		cmd.Run()
	})
	screen.Sync()
	var sw *kvitui.Window
	var height float32
	screen.Do(func() {
		sw = settings.Window()
		if sw == nil {
			t.Fatal("on the macOS form Settings did not open a window")
		}
		height = sw.ContentRect().Height
		if sw.Title() != "General" {
			t.Errorf("the window is titled %q, not after the pane shown", sw.Title())
		}
	})
	tabs := nodes(screen, sw, role.Tab)
	if len(tabs) != 3 || !tabs[0].Selected || tabs[1].Name != "Accounts" {
		t.Fatalf("the panes' buttons are %+v", tabs)
	}
	// The keyboard starts on the pane shown, and Right shows the next.
	if n := focusedNode(screen, sw); n == nil || n.Name != "General" {
		t.Errorf("the keyboard starts on %v, not the pane shown", n)
	}
	screen.KeyPress(unison.KeyRight, 0)
	screen.Sync()
	screen.Do(func() {
		if sw.Title() != "Accounts" || sw.ContentRect().Height <= height {
			t.Errorf("after Right the window is %q and %.0f tall, from %.0f", sw.Title(), sw.ContentRect().Height, height)
		}
	})
	// The nodes' bounds are in the window; the press is on the screen.
	var origin geom.Point
	screen.Do(func() { origin = sw.ContentRect().Point })
	screen.Click(tabs[2].Bounds.Center().Add(origin))
	screen.Sync()
	screen.Do(func() {
		if sw.Title() != "Privacy" || !slices.Equal(shown, []int{1, 2}) {
			t.Errorf("a press on Privacy showed %q, and the panes shown were %v", sw.Title(), shown)
		}
	})
	// Moved and closed, it opens again where it was left.
	left := geom.NewPoint(120, 90)
	screen.Do(func() {
		r := sw.ContentRect()
		sw.SetContentRect(geom.Rect{Point: left, Size: r.Size})
	})
	screen.KeyPress(unison.KeyW, mod.Command)
	screen.Do(func() {
		if settings.Window() != nil {
			t.Error("Command-W left the settings window open")
		}
		if err := settings.Open(); err != nil {
			t.Fatal(err)
		}
	})
	screen.Sync()
	screen.Do(func() {
		if at := settings.Window().ContentRect().Point; at != left {
			t.Errorf("the window opened again at %v, not where it was left at %v", at, left)
		}
		settings.Close()
	})
}

// On the Windows form Settings shows a page made of the same panes: one
// column of sections under their names, in order, About last, no wider than
// a readable width, with no OK, Cancel or Apply (WIN-16, WIN-23).
func TestSettingsMakeAPageOnWindows(t *testing.T) {
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler { return nil })
	var row *kvitui.SettingRow
	screen.Do(func() {
		row = kvitui.NewSettingRow(ui, "Reduce motion", "", kvitui.NewSwitch(ui, ""))
		settings := kvitui.NewSettingsWindow(ui,
			kvitui.SettingsPane{Name: "General", Symbol: "settings", Content: kvitui.Column(ui, kvitui.SizeSpace, kvitui.FullWidth(row))},
			kvitui.SettingsPane{Name: "Accounts", Symbol: "bank", Content: kvitui.Column(ui, kvitui.SizeSpace,
				kvitui.FullWidth(kvitui.NewSettingRow(ui, "Group by bank", "", kvitui.NewSwitch(ui, ""))))})
		about := &kvitui.About{Name: "Kvit Cash", Version: "2.4.1"}
		settings.Command(func() { w.SetBody(settings.Page(about)) }).Run()
		if settings.Window() != nil {
			t.Error("on the Windows form Settings opened a window")
		}
	})
	screen.Sync()
	var names []string
	for _, n := range nodes(screen, w, role.Heading) {
		names = append(names, n.Name)
	}
	if !slices.Equal(names, []string{"General", "Accounts", "About"}) {
		t.Errorf("the page's headings are %v", names)
	}
	for _, n := range nodes(screen, w, role.Button) {
		switch n.Name {
		case "OK", "Cancel", "Apply":
			t.Errorf("the settings page has %s", n.Name)
		}
	}
	screen.Do(func() {
		if got, most := row.FrameRect().Width, float32(ui.Interface.Px(720)); got > most || got < most/2 {
			t.Errorf("in a window 1500 wide a setting row is %.0f wide, not up to %.0f", got, most)
		}
	})
}

// A form row's label is above its field on the Windows form and beside it
// on the macOS form, names the field, marks it required, and takes the
// field's access key (LAY-6, CON-4, KEY-15).
func TestAFormRowLabelsItsField(t *testing.T) {
	var row *kvitui.FormRow
	var field *kvitui.Field
	screen, ui, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		field = kvitui.NewField(ui)
		row = kvitui.NewFormRow(ui, "&Payee", field)
		row.Required = true
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(300), row), kvitui.NewButton(ui, "Import")}
	})
	if n := screen.AccessibilityNodeFor(field.Edit()); n == nil || n.Name != "Payee" || !strings.Contains(n.Description, "Required") {
		t.Errorf("the field is %+v", n)
	}
	label := row.Children()[0]
	screen.Do(func() {
		if label.FrameRect().Bottom() > field.FrameRect().Y {
			t.Error("on the Windows form the label is not above the field")
		}
	})
	screen.KeyPress(unison.KeyP, mod.Option)
	if !field.Edit().Focused() {
		t.Error("Alt+P did not put the keyboard in the field")
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.Do(func() {
		if label.FrameRect().Right() > field.FrameRect().X || label.FrameRect().Y > field.FrameRect().Bottom() {
			t.Error("on the macOS form the label is not beside the field")
		}
	})
}

// A field masks a secret, reports text past its limit without cutting it,
// and shows a clear button while focused with text on the Windows form only
// (CTL-7).
func TestAFieldMasksLimitsAndClears(t *testing.T) {
	var secret, name *kvitui.Field
	screen, ui, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		secret, name = kvitui.NewField(ui), kvitui.NewField(ui)
		secret.Label, secret.Masked = "Password", true
		name.Label, name.MaxLength = "Code", 4
		return []unison.Paneler{kvitui.Left(secret), kvitui.Left(name)}
	})
	if n := screen.AccessibilityNodeFor(secret.Edit()); n == nil || !n.Protected {
		t.Errorf("a masked field is told as %+v", n)
	}
	screen.Click(screen.PanelCenter(name.Edit()))
	screen.Type("ABCDEF")
	if name.Text() != "ABCDEF" || !strings.Contains(name.Error, "2 characters over") {
		t.Errorf("past the limit the field holds %q and says %q", name.Text(), name.Error)
	}
	clearShown := func() bool {
		shown := false
		screen.Do(func() {
			for _, c := range name.Children() {
				if b, ok := c.Self.(*kvitui.IconButton); ok && b.Label == "Clear" && !b.Hidden {
					shown = true
				}
			}
		})
		return shown
	}
	if !clearShown() {
		t.Error("the Windows form shows no clear button in a focused field with text")
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.Sync()
	if clearShown() {
		t.Error("the macOS form shows a clear button")
	}
}

// A window's title names its document in the form's way and marks unsaved
// changes; AllowClose can keep it open; full screen fills the screen and
// gives the frame back (WIN-12, WIN-15, WIN-19, WIN-20).
func TestTheWindowTitlesAndKeepsItsDocument(t *testing.T) {
	screen, ui, w := windowSession(t, func(*kvitui.UI) []unison.Paneler { return nil })
	screen.Do(func() {
		ui.AppName = "Kvit Cash"
		w.SetDocument("May budget")
		w.SetEdited(true)
	})
	if got := w.Title(); got != "*May budget - Kvit Cash" {
		t.Errorf("the Windows form's title is %q", got)
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.Sync()
	if got := w.Title(); got != "May budget" {
		t.Errorf("the macOS form's title is %q", got)
	}
	asked := 0
	screen.Do(func() {
		w.AllowClose = func() bool { asked++; return false }
		w.AttemptClose()
	})
	screen.Sync()
	if asked != 1 || !w.IsValid() {
		t.Errorf("AllowClose was asked %d times and the window is open: %v", asked, w.IsValid())
	}
	var before geom.Rect
	screen.Do(func() {
		before = w.FrameRect()
		w.SetFullScreen(true)
	})
	screen.Do(func() {
		if d := w.Display(); d != nil && w.FrameRect() != d.Frame {
			t.Errorf("full screen is %v, not the screen's %v", w.FrameRect(), d.Frame)
		}
		w.SetFullScreen(false)
		if w.FrameRect() != before {
			t.Errorf("leaving full screen put the window at %v, not %v", w.FrameRect(), before)
		}
	})
}

// A window put back from the settings file is where it was left, fitted to
// the screen (WIN-14).
func TestAWindowIsPutBackWhereItWasLeft(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ui.json")
	ui, err := kvitui.New(kvitui.Options{SettingsPath: path, IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	var w *kvitui.Window
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 1500, Height: 1000},
		unison.StartupFinishedCallback(func() {
			w, _ = kvitui.NewWindow(ui, "ledger")
			w.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.Sync()
	screen.Do(func() {
		w.RememberFrame("window.ledger")
		w.SetFrameRect(geom.NewRect(40, 30, 1000, 700))
		if w.MovedCallback != nil {
			w.MovedCallback()
		}
		v, ok := ui.Settings.Value("window.ledger")
		if f, _ := v.([]any); !ok || len(f) != 4 {
			t.Fatalf("the frame kept is %v", v)
		}
		// Somewhere off the screen, as on a screen since unplugged.
		ui.Settings.SetValue("window.ledger", []float64{5000, 4000, 1000, 700})
		w.RememberFrame("window.ledger")
		if r, d := w.FrameRect(), w.Display(); d != nil && r != d.FitRectOnto(r) {
			t.Errorf("the window was put back at %v, off the screen's %v", r, d.Usable)
		}
	})
}

// A tip blocks nothing, is not shown again once closed, and comes back
// after ResetTips (HLP-4).
func TestATipBlocksNothingAndIsShownOnce(t *testing.T) {
	var target, other *kvitui.Button
	pressed := 0
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		target = kvitui.NewButton(ui, "Split")
		other = kvitui.NewButton(ui, "Merge")
		other.OnClick = func() { pressed++ }
		return []unison.Paneler{kvitui.Row(ui, kvitui.SizeSpace, target, other)}
	})
	tip := kvitui.NewTip(ui, "split", "Split a transaction", "Share one payment among several categories.")
	shown := false
	screen.Do(func() { shown = tip.ShowBeside(target) })
	screen.Sync()
	if !shown || len(w.Popups()) != 1 {
		t.Fatalf("the tip was not shown: %v", shown)
	}
	screen.Click(screen.PanelCenter(other))
	if pressed != 1 || !tip.Shown() {
		t.Errorf("with a tip shown a press elsewhere pressed %d times and closed it: %v", pressed, !tip.Shown())
	}
	for _, n := range nodes(screen, w, role.Button) {
		if n.Name == "Close" {
			screen.Click(n.Bounds.Center())
		}
	}
	screen.Do(func() { shown = tip.ShowBeside(target) })
	if shown || tip.Shown() {
		t.Error("a tip closed once was shown again")
	}
	screen.Do(func() {
		ui.ResetTips()
		shown = tip.ShowBeside(target)
	})
	if !shown {
		t.Error("ResetTips did not bring the tip back")
	}
}

// Back follows the reader's own history, is disabled with nowhere to go,
// and answers its standard key (LAY-10); the breadcrumb's "…" lists the
// folded places.
func TestBackFollowsTheHistory(t *testing.T) {
	var went []any
	var back *kvitui.IconButton
	var h *kvitui.History
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		h = kvitui.NewHistory(ui, "accounts", func(p any) { went = append(went, p) })
		back = kvitui.NewBackButton(ui, h)
		return []unison.Paneler{back}
	})
	screen.Do(func() { _ = w.AddCommands(h.Command()) })
	if n := screen.AccessibilityNodeFor(back); n == nil || !n.Disabled || n.Name != "Back" {
		t.Errorf("Back with no history is %+v", n)
	}
	screen.Do(func() {
		h.Visit("budgets")
		h.Visit("reports")
	})
	screen.KeyPress(unison.KeyLeft, mod.Option)
	if !slices.Equal(went, []any{"budgets"}) || h.Current() != "budgets" {
		t.Errorf("Alt+Left went to %v", went)
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.KeyPress(unison.KeyOpenBracket, mod.Command)
	if h.CanGoBack() {
		t.Error("Cmd-[ did not go back to the start")
	}
	var crumbs *kvitui.Breadcrumb
	screen.Do(func() {
		crumbs = kvitui.NewBreadcrumb(ui, kvitui.Crumb{Label: "Ledger"}, kvitui.Crumb{Label: "2026"}, kvitui.Crumb{Label: "May"},
			kvitui.Crumb{Label: "Week 2"}, kvitui.Crumb{Label: "Groceries"}, kvitui.Crumb{Label: "Receipt"})
		w.SetBody(crumbs)
	})
	screen.Sync()
	screen.Click(screen.PanelCenter(crumbs.Children()[2]))
	var lines []string
	for _, n := range nodes(screen, w, role.MenuItem) {
		lines = append(lines, n.Name)
	}
	if !slices.Equal(lines, []string{"2026", "May", "Week 2"}) {
		t.Errorf("the folded crumbs listed %v", lines)
	}
}

// Progress offers Stop and Pause, and says the time left (FDB-3, FDB-10).
func TestProgressCanBeStoppedAndPaused(t *testing.T) {
	var p *kvitui.Progress
	halted := 0
	var paused []bool
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		p = kvitui.NewProgress(ui, "Importing statement", 0.4)
		p.HaltText, p.OnHalt = "Stop", func() { halted++ }
		p.Pausable, p.OnPause = true, func(on bool) { paused = append(paused, on) }
		p.Remaining = 3 * time.Minute
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(400), p)}
	})
	if n := screen.AccessibilityNodeFor(p); n == nil || !strings.Contains(n.Description, "About 3 minutes left") {
		t.Errorf("the bar is told as %+v", n)
	}
	press := func(name string) {
		for _, n := range nodes(screen, w, role.Button) {
			if n.Name == name {
				screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: n.ID, Action: accessibility.Press})
				screen.Sync()
				return
			}
		}
		t.Errorf("no %q button", name)
	}
	press("Pause")
	press("Resume")
	press("Stop")
	if halted != 1 || !slices.Equal(paused, []bool{true, false}) {
		t.Errorf("Stop halted %d times and Pause gave %v", halted, paused)
	}
}

// The undo stack names Undo and Redo for the edit they reverse (FDB-7,
// WRD-11); a notification is not posted while a Kvit window is in front
// (MSG-10); a help link opens its topic in the form's way (HLP-6); a
// selectable label copies its text (CTL-18).
func TestUndoNotifyHelpAndCopy(t *testing.T) {
	var help *kvitui.HelpLink
	var address *kvitui.Label
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.HelpRoot = "help"
		help = kvitui.NewHelpLink(ui, "budgets#rollover", "How does a budget roll over?")
		address = kvitui.NewLabel(ui, "IBAN GB33 BUKB 2020 1555 5555 55")
		address.Selectable = true
		return []unison.Paneler{help, address}
	})
	stack := kvitui.NewUndoStack(ui)
	amount := 10
	stack.Add("Delete transaction", func() { amount = 10 }, func() { amount = 0 })
	amount = 0
	undo, redo := stack.Commands()
	if undo.Label(ui) != "Undo delete transaction" || redo.CanRun() {
		t.Errorf("Undo says %q and Redo can run: %v", undo.Label(ui), redo.CanRun())
	}
	undo.Do()
	if amount != 10 || redo.Label(ui) != "Redo delete transaction" || undo.CanRun() {
		t.Errorf("after Undo the amount is %d and Redo says %q", amount, redo.Label(ui))
	}
	tray := platform.NewTestTray(platform.Authorized)
	screen.Do(func() {
		if ui.Notify(tray, platform.Notification{Title: "Sync finished"}) {
			t.Error("a notification was posted while the window was in front")
		}
	})
	links := nodes(screen, w, role.Link)
	if len(links) != 1 || links[0].Name != "How does a budget roll over?" {
		t.Fatalf("the help link is %+v", links)
	}
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: links[0].ID, Action: accessibility.Press})
	screen.Sync()
	if urls := screen.OpenedURLs(); len(urls) != 1 || !strings.HasSuffix(urls[0], "/help/budgets.html#rollover") {
		t.Errorf("the help link opened %v", urls)
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.Sync()
	if n := nodes(screen, w, role.Link); len(n) != 0 {
		t.Error("the macOS form shows the help as a link")
	}
	screen.Click(screen.PanelCenter(address))
	screen.KeyPress(unison.KeyC, mod.Command)
	if got := unison.ClipboardGetText(); got != address.Text {
		t.Errorf("copying the label put %q on the clipboard", got)
	}
}

// With a popover open, a press on a button closes it without pressing the
// button, and a press on another popup's opener opens that one in the same
// press (PTR-13).
func TestAPressOutsideAPopoverPressesNothing(t *testing.T) {
	var opener, behind *kvitui.Button
	var pick *kvitui.Select
	pressed := 0
	var pop *kvitui.Popover
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		opener = kvitui.NewButton(ui, "Filter")
		behind = kvitui.NewButton(ui, "Import")
		behind.OnClick = func() { pressed++ }
		pick = kvitui.NewSelect(ui, "Account", kvitui.Option{Value: "a", Label: "Checking"}, kvitui.Option{Value: "b", Label: "Savings"})
		return []unison.Paneler{kvitui.Row(ui, kvitui.SizeColumnGap, opener, behind, pick)}
	})
	open := func() {
		screen.Do(func() {
			pop = kvitui.NewPopover(ui, "Filter", kvitui.NewLabel(ui, "Only this month"))
			pop.Open(opener, kvitui.PlaceBelow(ui, opener))
		})
		screen.Sync()
	}
	open()
	screen.Click(screen.PanelCenter(behind))
	if pop.Opened() || pressed != 0 {
		t.Errorf("a press on a button closed the popover: %v, and pressed the button %d times", !pop.Opened(), pressed)
	}
	open()
	screen.Click(screen.PanelCenter(pick))
	if pop.Opened() || len(nodes(screen, w, role.Menu)) != 1 {
		t.Error("a press on a select's opener did not close the popover and open the list")
	}
}
