package kvitui_test

import (
	"slices"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// On the Windows form Alt and a control's access key press it: the letter
// the app chose with "&", or the first free one; OK, Cancel and Close take
// none; the macOS form has none (KEY-15).
func TestAccessKeysOperateControls(t *testing.T) {
	pressed := map[string]int{}
	button := func(ui *kvitui.UI, words string) *kvitui.Button {
		b := kvitui.NewButton(ui, words)
		b.OnClick = func() { pressed[words]++ }
		return b
	}
	var archive *kvitui.Check
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		archive = kvitui.NewCheck(ui, "Show &archived")
		return []unison.Paneler{button(ui, "Import"), button(ui, "Export"), button(ui, "&Reconcile"), button(ui, "Cancel"), archive}
	})
	var keys []kvitui.AccessKey
	screen.Do(func() { keys = w.AccessKeys() })
	got := map[string]rune{}
	for _, k := range keys {
		got[k.Words] = k.Key
	}
	if got["Import"] != 'i' || got["Export"] != 'e' || got["Reconcile"] != 'r' || got["Show archived"] != 'a' {
		t.Errorf("the access keys are %v", got)
	}
	if _, ok := got["Cancel"]; ok {
		t.Error("Cancel was given an access key")
	}
	screen.KeyPress(unison.KeyR, mod.Option)
	screen.KeyPress(unison.KeyA, mod.Option)
	if pressed["&Reconcile"] != 1 || !archive.Checked {
		t.Errorf("Alt+R pressed Reconcile %d times and Alt+A left the box %v", pressed["&Reconcile"], archive.Checked)
	}
	if n := screen.AccessibilityNodeFor(archive); n == nil || n.Name != "Show archived" {
		t.Errorf("the box's name kept its marker: %+v", n)
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.KeyPress(unison.KeyI, mod.Option)
	if pressed["Import"] != 0 {
		t.Error("the macOS form pressed a button for Option-I")
	}
}

// The access keys are underlined while Alt is held, or at all times when
// Windows' setting asks for it; never on the macOS form (KEY-15).
func TestAccessKeysShowWhenTheSystemAsks(t *testing.T) {
	var button *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewButton(ui, "Import")
		return []unison.Paneler{kvitui.Left(button)}
	})
	shown := func() bool {
		var on bool
		screen.Do(func() { on = w.ShowsAccessKeys() })
		return on
	}
	box := func() geom.Rect {
		var r geom.Rect
		screen.Do(func() { r = button.RectToRoot(button.ContentRect(true)) })
		return r
	}
	drawn := func() []byte {
		screen.Sync()
		img := screen.Capture()
		r := box()
		var out []byte
		for y := int(r.Y); y < int(r.Bottom()); y++ {
			for x := int(r.X); x < int(r.Right()); x++ {
				c := img.NRGBAAt(x, y)
				out = append(out, c.R, c.G, c.B)
			}
		}
		return out
	}
	if shown() {
		t.Fatal("the access keys show before Alt is held")
	}
	plain := drawn()
	screen.Do(func() { ui.UnderlineAccessKeys = true; w.MarkForRedraw() })
	if !shown() {
		t.Fatal("the access keys do not show with the system's setting on")
	}
	if slices.Equal(drawn(), plain) {
		t.Error("the button looks the same with its access key underlined")
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	if shown() {
		t.Error("the macOS form shows access keys")
	}
}

// A row of tabs, a segmented control and a radio group are each one stop in
// the Tab order, and the arrow keys move inside them; a radio group stops at
// its ends on the Windows form, where Ctrl and an arrow move without
// choosing, and goes round on the macOS form (KEY-9).
func TestArrowKeysMoveInsideGroups(t *testing.T) {
	var tabs []*kvitui.Tab
	var period *kvitui.Segmented
	var radio *kvitui.RadioGroup
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		var row []unison.Paneler
		for i, name := range []string{"Accounts", "Budgets", "Reports"} {
			tab := kvitui.NewTab(ui, name)
			tab.Selected = i == 0
			tab.OnClick = func() {
				for _, o := range tabs {
					o.Selected = o == tab
				}
			}
			tabs = append(tabs, tab)
			row = append(row, tab)
		}
		period = kvitui.NewSegmented(ui, "Period", kvitui.Option{Value: "m", Label: "Month"}, kvitui.Option{Value: "y", Label: "Year"})
		radio = kvitui.NewRadioGroup(ui, "Theme", kvitui.RadioOption{Value: "l", Label: "Light"}, kvitui.RadioOption{Value: "d", Label: "Dark"})
		radio.Current = "d"
		return []unison.Paneler{kvitui.Row(ui, kvitui.Px(0), row...), kvitui.Left(period), radio}
	})
	focused := func() string {
		if n := focusedNode(screen, w); n != nil {
			return n.Name
		}
		return ""
	}
	screen.Do(func() { tabs[0].Focus() })
	screen.KeyPress(unison.KeyRight, 0)
	if !tabs[1].Selected || focused() != "Budgets" {
		t.Errorf("Right on a tab left %q focused", focused())
	}
	screen.KeyPress(unison.KeyTab, 0)
	if focused() != "Month" {
		t.Errorf("Tab from the tabs went to %q, not the segmented control", focused())
	}
	screen.KeyPress(unison.KeyRight, 0)
	if period.Current != "y" {
		t.Error("Right on a segment did not choose the next")
	}
	screen.KeyPress(unison.KeyTab, 0)
	if focused() != "Dark" {
		t.Errorf("Tab went to %q, not the chosen radio option", focused())
	}
	screen.KeyPress(unison.KeyDown, 0)
	if radio.Current != "d" {
		t.Error("on the Windows form Down on the last option went round")
	}
	// Ctrl and an arrow move without choosing, and Space chooses.
	screen.KeyPress(unison.KeyUp, mod.Control)
	if radio.Current != "d" || focused() != "Light" {
		t.Errorf("Ctrl+Up chose %q and left %q focused", radio.Current, focused())
	}
	screen.KeyPress(unison.KeyTab, mod.Shift)
	screen.KeyPress(unison.KeyTab, 0)
	if focused() != "Dark" {
		t.Errorf("Tab back into the group after Ctrl+Up went to %q, not the chosen option", focused())
	}
	screen.KeyPress(unison.KeyUp, mod.Control)
	screen.KeyPress(unison.KeySpace, 0)
	if radio.Current != "l" {
		t.Errorf("Space after Ctrl+Up chose %q", radio.Current)
	}
	screen.KeyPress(unison.KeyDown, 0)
	if radio.Current != "d" || focused() != "Dark" {
		t.Errorf("Down chose %q and left %q focused", radio.Current, focused())
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.KeyPress(unison.KeyDown, 0)
	if radio.Current != "l" {
		t.Error("on the macOS form Down on the last option did not go round")
	}
	// Ctrl+Tab chooses the next tab from anywhere in the window.
	screen.KeyPress(unison.KeyTab, mod.Control)
	if !tabs[2].Selected {
		t.Error("Ctrl+Tab did not choose the next tab")
	}
}

// Find's keys put the keyboard in the window's search field with its text
// selected (KEY-13); Escape closes a type-ahead's list and abandons a dialog
// being moved (KEY-7); Space pages a region (KEY-11).
func TestFindEscapeAndSpace(t *testing.T) {
	var search *kvitui.SearchField
	var other *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		search = kvitui.NewSearchField(ui)
		search.Placeholder = "Search accounts"
		other = kvitui.NewButton(ui, "Import")
		return []unison.Paneler{kvitui.Left(search), other}
	})
	screen.Do(func() {
		search.SetText("bank")
		other.Focus()
	})
	screen.KeyPress(unison.KeyE, mod.Control)
	if n := focusedNode(screen, w); n == nil || n.Role != role.TextField {
		t.Fatalf("Ctrl+E left the keyboard on %v", n)
	}
	screen.Do(func() {
		if s, e := search.Edit().Selection(); s != 0 || e != 4 {
			t.Errorf("Find left the selection at %d to %d, not the whole text", s, e)
		}
	})
	// A dialog dragged by its title row goes back where it was on Escape.
	var d *kvitui.Dialog
	screen.Do(func() {
		d = kvitui.NewDialog(ui, "Rename account")
		d.ConfirmText = "Rename"
		d.Open(w)
	})
	screen.Sync()
	var before geom.Rect
	screen.Do(func() { before = d.FrameRect() })
	m := ui.Interface
	grab := screen.PanelPoint(d, geom.NewPoint(float32(m.ViewMargin()), float32(m.ViewMargin()+m.ControlHeight()/2)))
	screen.MouseDown(grab, unison.ButtonLeft, 0)
	screen.MouseMove(geom.NewPoint(grab.X-80, grab.Y-60), 0)
	screen.Sync()
	screen.KeyPress(unison.KeyEscape, 0)
	screen.MouseUp(geom.NewPoint(grab.X-80, grab.Y-60), unison.ButtonLeft, 0)
	screen.Sync()
	var after geom.Rect
	screen.Do(func() { after = d.FrameRect() })
	if after != before || len(w.Popups()) != 1 {
		t.Errorf("Escape during the move left the dialog at %v, not %v, with %d popups", after, before, len(w.Popups()))
	}
}

// Typing the start of a name jumps to it in a tree and a select (KEY-9).
func TestTypingJumpsToARow(t *testing.T) {
	var tree *kvitui.Tree
	var currency *kvitui.Select
	var chosen []int
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		tree = kvitui.NewTree(ui, "Accounts",
			kvitui.TreeNode{Label: "Checking"}, kvitui.TreeNode{Label: "Credit card"}, kvitui.TreeNode{Label: "Savings"})
		tree.OnChoose = func(path []int) { chosen = path }
		currency = kvitui.NewSelect(ui, "Currency",
			kvitui.Option{Value: "eur", Label: "Euro"}, kvitui.Option{Value: "gbp", Label: "Pound"}, kvitui.Option{Value: "usd", Label: "US dollar"})
		return []unison.Paneler{kvitui.Height(ui, kvitui.Px(120), tree), kvitui.Left(currency)}
	})
	screen.Do(func() { tree.RequestFocus() })
	screen.Type("s")
	if !slices.Equal(chosen, []int{2}) {
		t.Errorf("typing s chose %v, not Savings", chosen)
	}
	screen.Type("c")
	screen.Type("c")
	if !slices.Equal(chosen, []int{1}) {
		t.Errorf("typing c twice chose %v, not the second row starting with c", chosen)
	}
	screen.Do(func() { currency.Focus() })
	screen.Type("u")
	if currency.Current != "usd" {
		t.Errorf("typing u in the select chose %q", currency.Current)
	}
}

// On the macOS form with full keyboard access off, Tab goes only to text
// fields and lists (KEY-10).
func TestTabFollowsTheMacsKeyboardSetting(t *testing.T) {
	var name, memo *kvitui.Field
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		name, memo = kvitui.NewField(ui), kvitui.NewField(ui)
		name.Label, memo.Label = "Name", "Memo"
		return []unison.Paneler{kvitui.Left(name), kvitui.NewButton(ui, "Import"), kvitui.Left(memo)}
	})
	screen.Do(func() {
		ui.SetForm(kvitui.FormMac)
		ui.FullKeyboardAccess = false
		name.Focus()
	})
	screen.KeyPress(unison.KeyTab, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Memo" {
		t.Errorf("Tab without full keyboard access went to %v, not the next field", n)
	}
	screen.Do(func() { ui.FullKeyboardAccess = true })
	screen.KeyPress(unison.KeyTab, mod.Shift)
	if n := focusedNode(screen, w); n == nil || n.Name != "Import" {
		t.Errorf("Shift+Tab with full keyboard access went to %v, not the button", n)
	}
}
