package uitest_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// found reports whether a problem of a rule mentioning a word was found.
func found(problems []uitest.Problem, rule, word string) bool {
	for _, p := range problems {
		if p.Rule == rule && strings.Contains(p.String(), word) {
			return true
		}
	}
	return false
}

// The label checks find each kind of problem planted on a screen, in the
// form it belongs to, and pass a screen written by the rules.
func TestTheLabelChecksFindWhatIsPlanted(t *testing.T) {
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		article := kvitui.NewButton(ui, "Delete the file")
		letter := kvitui.NewButton(ui, "B")
		colon := kvitui.NewButton(ui, "Export:")
		titled := kvitui.NewButton(ui, "Show All Columns")
		disabled := kvitui.NewButton(ui, "Archive")
		disabled.SetEnabled(false)
		busy := kvitui.NewButton(ui, "Save")
		busy.Busy, busy.BusyText = true, "Working…"
		search := kvitui.NewSearchField(ui)
		good := kvitui.NewButton(ui, "Export statement")
		return kvitui.Column(ui, kvitui.SizeSpace, article, letter, colon, titled, disabled, busy, search, good)
	})
	var problems []uitest.Problem
	s.Do(func() {
		d := kvitui.NewDialog(s.UI, "Error")
		d.ConfirmText = "Confirm"
		d.Open(s.Window)
	})
	s.Sync()
	s.Do(func() { problems = uitest.Labels(s.Window) })
	for _, want := range [][2]string{
		{"WRD-3", "Delete the file"}, {"CMD-22", `"B"`}, {"WRD-19", "Export:"}, {"WRD-6", "Show All Columns"},
		{"CTL-8", "Archive"}, {"WRD-25", "Working"}, {"CTL-15", "search field"}, {"WRD-7", "Confirm"}, {"MSG-4", "Error"},
	} {
		if !found(problems, want[0], want[1]) {
			t.Errorf("the checks missed %s on %s", want[0], want[1])
		}
	}
	for _, p := range problems {
		if strings.Contains(p.String(), "Export statement") {
			t.Errorf("a label written by the rules was reported: %s", p)
		}
	}
	// On the macOS form a sentence-case label written by hand is not title
	// case, which kvit-ui gives every button it draws, so only the dialog's
	// words and the planted labels are reported.
	s.SetForm(kvitui.FormMac)
	s.Do(func() { problems = uitest.Labels(s.Window) })
	if found(problems, "WRD-6", "Export Statement") {
		t.Error("on the macOS form a button kvit-ui wrote in title case was reported")
	}
}

// On the macOS form a button's words are in title case and a symbol
// button's name in sentence case, as Apple writes a tooltip, so a button and
// a symbol button for one command are not one command written two ways.
func TestTheAppChecksTakeTheMacFormsCases(t *testing.T) {
	s := uitest.Open(t, uitest.Options{Form: kvitui.FormMac}, func(ui *kvitui.UI) unison.Paneler {
		return kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewButton(ui, "Add agent"),
			kvitui.NewIconButton(ui, "plus", "Add agent"))
	})
	if problems := uitest.App(s.Collect("agents")); found(problems, "CON-1", "add agent") {
		t.Errorf("a button and a symbol button for one command were reported: %v", problems)
	}
}

// The whole-app checks find one command written two ways, one symbol with
// two meanings, one key for two commands, a standard command made by hand
// and a command on no surface.
func TestTheAppChecksFindWhatIsPlanted(t *testing.T) {
	run := func() {}
	var screens []uitest.Screen
	for i, build := range []func(ui *kvitui.UI) unison.Paneler{
		func(ui *kvitui.UI) unison.Paneler {
			capture := &kvitui.Command{Name: "Quick capture", Asks: true, Shortcut: kvitui.Keys(unison.KeyK, kvitui.ModPrimary), Run: run}
			archive := kvitui.NewIconButton(ui, "close", "Archive")
			closing := kvitui.NewIconButton(ui, "close", "Close")
			return kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewCommandButton(ui, capture), archive, closing)
		},
		func(ui *kvitui.UI) unison.Paneler {
			note := &kvitui.Command{Name: "Quick capture note", Shortcut: kvitui.Keys(unison.KeyK, kvitui.ModPrimary), Run: run}
			save := &kvitui.Command{Name: "Save", Run: run}
			other := kvitui.NewButton(ui, "Quick Capture")
			return kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewCommandButton(ui, note), kvitui.NewCommandButton(ui, save), other)
		},
	} {
		s := uitest.Open(t, uitest.Options{}, build)
		screens = append(screens, s.Collect([]string{"capture", "notes"}[i]))
		// One headless screen runs at a time.
		s.Screen.Stop()
	}
	problems := uitest.App(screens...)
	for _, want := range [][2]string{
		{"CON-1", "quick capture"}, {"CON-2", "close"}, {"KEY-5", "Ctrl+K"}, {"KEY-2", "Save"}, {"CMD-1", "Quick capture"},
	} {
		if !found(problems, want[0], want[1]) {
			t.Errorf("the checks missed %s on %s: %v", want[0], want[1], problems)
		}
	}
}

// A help link is worded as a whole question, never generically; the macOS
// form has one help button to a window; a generic Help button is reported;
// and the whole-app checks report a topic the app does not ship (HLP-6).
func TestTheHelpChecksFindWhatIsPlanted(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "budgets.html"), []byte(`<h2 id="rollover">Rolling over</h2>`), 0o644); err != nil {
		t.Fatal(err)
	}
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		ui.HelpRoot = root
		return kvitui.Column(ui, kvitui.SizeSpace,
			kvitui.NewHelpLink(ui, "budgets#rollover", "How does a budget roll over?"),
			kvitui.NewHelpLink(ui, "budgets#carry", "Learn more about budgets"),
			kvitui.NewHelpLink(ui, "goals", "What is a goal?!"),
			kvitui.NewButton(ui, "Help"))
	})
	var problems []uitest.Problem
	s.Do(func() { problems = uitest.Labels(s.Window) })
	for _, want := range [][2]string{
		{"HLP-6", "Learn more about budgets"}, {"HLP-6", "What is a goal?!"}, {"HLP-6", `button "Help"`},
	} {
		if !found(problems, want[0], want[1]) {
			t.Errorf("the checks missed %s on %s: %v", want[0], want[1], problems)
		}
	}
	if found(problems, "HLP-6", "How does a budget roll over?") || found(problems, "HLP-6", "help buttons") {
		t.Errorf("on the Windows form a good help link, or the count of help buttons, was reported: %v", problems)
	}
	s.SetForm(kvitui.FormMac)
	s.Do(func() { problems = uitest.Labels(s.Window) })
	if !found(problems, "HLP-6", "3 help buttons") {
		t.Errorf("on the macOS form three help buttons in one window were not reported: %v", problems)
	}
	app := uitest.App(s.Collect("budgets"))
	for _, want := range []string{"carry", "goals"} {
		if !found(app, "HLP-6", want) {
			t.Errorf("the app checks missed the topic %q that is not shipped: %v", want, app)
		}
	}
	if found(app, "HLP-6", "rollover") {
		t.Errorf("a shipped topic was reported: %v", app)
	}
}

// reversed lays its children out in the opposite order to the Tab order:
// the last child first, down a column or along a row.
type reversed struct{ row bool }

func (l reversed) LayoutSizes(target *unison.Panel, _ geom.Size) (minSize, prefSize, maxSize geom.Size) {
	var size geom.Size
	for _, c := range target.Children() {
		_, p, _ := c.Sizes(geom.Size{})
		if l.row {
			size.Width += p.Width
			size.Height = max(size.Height, p.Height)
		} else {
			size.Height += p.Height
			size.Width = max(size.Width, p.Width)
		}
	}
	return size, size, size
}

func (l reversed) PerformLayout(target *unison.Panel) {
	var at float32
	kids := target.Children()
	for i := len(kids) - 1; i >= 0; i-- {
		_, p, _ := kids[i].Sizes(geom.Size{})
		if l.row {
			kids[i].SetFrameRect(geom.NewRect(at, 0, p.Width, p.Height))
			at += p.Width
		} else {
			kids[i].SetFrameRect(geom.NewRect(0, at, p.Width, p.Height))
			at += p.Height
		}
	}
}

// Tab going back up a column or back along a row is reported; two columns
// read one after the other, and a dialog's help after its buttons, are not
// (KEY-9).
func TestTheTabOrderCheckFollowsTheReadingOrder(t *testing.T) {
	group := func(ui *kvitui.UI, row bool, words ...string) unison.Paneler {
		p := unison.NewPanel()
		p.SetLayout(reversed{row})
		for _, w := range words {
			p.AddChild(kvitui.NewButton(ui, w))
		}
		return kvitui.Left(p)
	}
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		columns := kvitui.Row(ui, kvitui.SizeSpace,
			kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewButton(ui, "Add account"), kvitui.NewButton(ui, "Import file")),
			kvitui.Column(ui, kvitui.SizeSpace, kvitui.NewButton(ui, "Export ledger"), kvitui.NewButton(ui, "Print report")))
		return kvitui.Column(ui, kvitui.SizeSpace, columns,
			group(ui, false, "Open budget", "Close budget"),
			group(ui, true, "Archive note", "Restore note"))
	})
	var problems []uitest.Problem
	s.Do(func() { problems = uitest.Labels(s.Window) })
	for _, want := range []string{`button "Close budget": Tab goes back up`, `button "Restore note": Tab goes back to the left`} {
		if !found(problems, "KEY-9", want) {
			t.Errorf("the check missed %s: %v", want, problems)
		}
	}
	for _, p := range problems {
		if p.Rule == "KEY-9" && (strings.Contains(p.Where, "Export ledger") || strings.Contains(p.Where, "Print report")) {
			t.Errorf("the second column was reported: %s", p)
		}
	}
	for _, form := range []kvitui.Form{kvitui.FormWindows, kvitui.FormMac} {
		s.SetForm(form)
		s.Do(func() {
			d := kvitui.NewDialog(s.UI, "Roll the budget over?")
			d.ConfirmText, d.OtherText = "Roll over", "Start empty"
			d.Help = kvitui.NewHelpLink(s.UI, "budgets#rollover", "How does a budget roll over?")
			d.Open(s.Window)
		})
		s.Sync()
		s.Do(func() { problems = uitest.Labels(s.Window) })
		if found(problems, "KEY-9", "") {
			t.Errorf("on the %s form a dialog with help was reported: %v", form, problems)
		}
		s.Screen.KeyPress(unison.KeyEscape, 0)
		s.Sync()
	}
}

// Apply is reported outside a properties dialog, and OK or Cancel on a page
// of settings; a properties dialog's own buttons are not (WIN-23).
func TestTheApplyChecksFindWhatIsPlanted(t *testing.T) {
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		settings := kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.SizeSpace,
			kvitui.FullWidth(kvitui.NewSettingRow(ui, "Reduce motion", "", kvitui.NewSwitch(ui, ""))),
			kvitui.Row(ui, kvitui.SizeSpace, kvitui.NewButton(ui, "OK"), kvitui.NewButton(ui, "Cancel"))))
		return kvitui.Column(ui, kvitui.SizeSpace, kvitui.Left(kvitui.NewButton(ui, "Apply")),
			kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(200), settings)))
	})
	var problems []uitest.Problem
	s.Do(func() { problems = uitest.Labels(s.Window) })
	for _, want := range []string{`button "Apply": Apply outside a properties dialog`, `button "OK": on a page of settings`, `button "Cancel": on a page of settings`} {
		if !found(problems, "WIN-23", want) {
			t.Errorf("the checks missed %s: %v", want, problems)
		}
	}
	s.Do(func() {
		d := kvitui.NewDialog(s.UI, "Checking properties")
		d.Kind = kvitui.DialogProperties
		d.Open(s.Window)
	})
	s.Sync()
	s.Do(func() { problems = uitest.Labels(s.Window) })
	n := 0
	for _, p := range problems {
		if p.Rule == "WIN-23" && strings.Contains(p.Where, "Apply") {
			n++
		}
	}
	// The page's own Apply is still there; the dialog's is not reported.
	if n != 1 {
		t.Errorf("with a properties dialog open, Apply was reported %d times: %v", n, problems)
	}
}

// Report fails what is not known, and a known problem that is gone.
func TestReportKeepsTheKnownListShrinking(t *testing.T) {
	problems := []uitest.Problem{{Rule: "WRD-3", Where: `button "Delete the file"`, What: "an article"}}
	rec := &recorder{}
	uitest.Report(rec, problems, `button "Delete the file"`, "gone")
	if len(rec.errors) != 1 || !strings.Contains(rec.errors[0], "gone") {
		t.Errorf("Report said %v", rec.errors)
	}
}

// recorder is a testing.TB that keeps what it is told.
type recorder struct {
	testing.TB
	errors []string
}

func (r *recorder) Helper() {}

func (r *recorder) Error(args ...any) { r.errors = append(r.errors, fmt.Sprint(args...)) }

func (r *recorder) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

// Names sharing a verb are one meaning of a symbol: several Add … buttons
// and a stepper's Increase … button on plus are not reported; two verbs
// outside such a table, as next period and show fields on one chevron, are
// (CON-2).
func TestSymbolUsesSharingAVerbAreOneMeaning(t *testing.T) {
	plus := func(ui *kvitui.UI, words string) *kvitui.Button {
		b := kvitui.NewButton(ui, words)
		b.Symbol = "plus"
		return b
	}
	var screens []uitest.Screen
	for i, build := range []func(ui *kvitui.UI) unison.Paneler{
		func(ui *kvitui.UI) unison.Paneler {
			return kvitui.Column(ui, kvitui.SizeSpace,
				plus(ui, "Add rule"), plus(ui, "Add condition"), plus(ui, "Add action"))
		},
		func(ui *kvitui.UI) unison.Paneler {
			return kvitui.Column(ui, kvitui.SizeSpace,
				kvitui.NewStepper(ui, "Interface size", 10, 12),
				kvitui.NewIconButton(ui, "chevron-down", "Next period"),
				kvitui.NewIconButton(ui, "chevron-down", "Show fields"))
		},
	} {
		s := uitest.Open(t, uitest.Options{}, build)
		screens = append(screens, s.Collect([]string{"rules", "ledger"}[i]))
		// One headless screen runs at a time.
		s.Screen.Stop()
	}
	problems := uitest.App(screens...)
	for _, symbol := range []string{"plus", "minus"} {
		if found(problems, "CON-2", symbol) {
			t.Errorf("uses sharing the adding or taking-away verb were reported: %v", problems)
		}
	}
	if !found(problems, "CON-2", "chevron-down") {
		t.Errorf("two verbs sharing one chevron were not reported: %v", problems)
	}
}

// A placeholder that is only a dash stands for a value that is not set,
// while words ending in punctuation are still reported (WRD-19).
func TestADashPlaceholderStandsForAnUnsetValue(t *testing.T) {
	field := func(ui *kvitui.UI, label, placeholder string) *kvitui.Field {
		f := kvitui.NewField(ui)
		f.Label, f.Placeholder = label, placeholder
		return f
	}
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		return kvitui.Column(ui, kvitui.SizeSpace,
			field(ui, "Payee", "—"), field(ui, "Notes", " – "),
			field(ui, "Code", "-"), field(ui, "Name", "Enter a name:"))
	})
	var problems []uitest.Problem
	s.Do(func() { problems = uitest.Labels(s.Window) })
	for _, label := range []string{"Payee", "Notes", "Code"} {
		if found(problems, "WRD-19", label) {
			t.Errorf("a dash standing for an unset value was reported: %v", problems)
		}
	}
	if !found(problems, "WRD-19", "Name") {
		t.Errorf("words ending in punctuation were not reported: %v", problems)
	}
}

// A compact select's words are the value in force: no WRD-1 in either form,
// and one calendar meaning whatever the period is (CON-2).
func TestACompactSelectReadsAsAValue(t *testing.T) {
	build := func(value string) func(ui *kvitui.UI) unison.Paneler {
		return func(ui *kvitui.UI) unison.Paneler {
			s := kvitui.NewSelect(ui, "Period",
				kvitui.Option{Value: "all", Label: "All history"},
				kvitui.Option{Value: "month", Label: "Month"})
			s.Symbol, s.Compact, s.Current = "calendar", true, value
			return kvitui.Column(ui, kvitui.SizeSpace, s)
		}
	}
	s := uitest.Open(t, uitest.Options{}, build("all"))
	var problems []uitest.Problem
	s.Do(func() { problems = uitest.Labels(s.Window) })
	s.SetForm(kvitui.FormMac)
	s.Do(func() { problems = append(problems, uitest.Labels(s.Window)...) })
	for _, p := range problems {
		if p.Rule == "WRD-1" && strings.Contains(p.String(), "All history") {
			t.Errorf("a compact select's value was read as a command: %s", p)
		}
	}
	// One headless screen runs at a time.
	s.Screen.Stop()
	var screens []uitest.Screen
	for i, value := range []string{"all", "month"} {
		r := uitest.Open(t, uitest.Options{}, build(value))
		screens = append(screens, r.Collect([]string{"ledger", "ledger-month"}[i]))
		// One headless screen runs at a time.
		r.Screen.Stop()
	}
	if app := uitest.App(screens...); found(app, "CON-2", "calendar") {
		t.Errorf("one calendar meaning period was reported: %v", app)
	}
}
