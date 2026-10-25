package kvitui_test

import (
	"slices"
	"strings"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// A shortcut is written the way each form writes it.
func TestShortcutsAreWrittenInEachFormsWay(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		s            kvitui.Shortcut
		windows, mac string
	}{
		{kvitui.Keys(unison.KeyS, kvitui.ModPrimary), "Ctrl+S", "⌘S"},
		{kvitui.Keys(unison.KeyZ, kvitui.ModPrimary|kvitui.ModShift), "Ctrl+Shift+Z", "⇧⌘Z"},
		{kvitui.Keys(unison.KeyEqual, kvitui.ModPrimary), "Ctrl+Plus", "⌘+"},
		{kvitui.Keys(unison.KeyF, kvitui.ModPrimary|kvitui.ModControl), "Ctrl+F", "⌃⌘F"},
		{kvitui.Keys(unison.KeySlash, kvitui.ModPrimary|kvitui.ModShift), "Ctrl+Shift+/", "⌘?"},
		{kvitui.Keys(unison.KeyF2, 0), "F2", "F2"},
		{kvitui.Keys(unison.KeyDelete, 0), "Del", "⌦"},
	}
	for _, c := range cases {
		ui.SetForm(kvitui.FormWindows)
		if got := ui.ShortcutText(c.s); got != c.windows {
			t.Errorf("Windows form: %q, want %q", got, c.windows)
		}
		ui.SetForm(kvitui.FormMac)
		if got := ui.ShortcutText(c.s); got != c.mac {
			t.Errorf("macOS form: %q, want %q", got, c.mac)
		}
	}
}

// An app that has Save gets Ctrl+S on the Windows form and Cmd-S on the
// macOS form without choosing a key, wherever the focus is in the window,
// typing in a field included; and a command that cannot act ignores its
// key without passing it on.
func TestAStandardCommandAnswersItsStandardKey(t *testing.T) {
	saved, passed := 0, 0
	canSave := true
	var field *kvitui.Field
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		field = kvitui.NewField(ui)
		field.Label = "Payee"
		return []unison.Paneler{kvitui.Left(field)}
	})
	save := &kvitui.Command{Std: kvitui.StdSave, Enabled: func() bool { return canSave }, Run: func() { saved++ }}
	screen.Do(func() {
		if err := w.AddCommands(save); err != nil {
			t.Error(err)
		}
		w.OnKeyDown = func(unison.KeyCode, mod.Modifiers, bool) bool { passed++; return false }
	})
	screen.Click(screen.PanelCenter(field.Edit()))
	screen.KeyPress(unison.KeyS, mod.Control)
	if saved != 1 {
		t.Errorf("Ctrl+S in a field saved %d times", saved)
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.KeyPress(unison.KeyS, mod.Control)
	screen.KeyPress(unison.KeyS, mod.Command)
	if saved != 2 {
		t.Errorf("on the macOS form Ctrl-S and Cmd-S saved %d times between them, want once", saved-1)
	}
	canSave = false
	passed = 0
	screen.KeyPress(unison.KeyS, mod.Command)
	if saved != 2 || passed != 0 {
		t.Errorf("Cmd-S with nothing to save saved %d times and passed the key on %d times", saved-2, passed)
	}
	if got := field.Text(); got != "" {
		t.Errorf("the shortcuts typed %q into the field", got)
	}
}

// A key a text field uses goes to the field, and a standard command the
// field carries out itself, such as Select all, is done by the field; the
// app's own command for it runs when the focus is elsewhere.
func TestATextFieldKeepsItsOwnKeys(t *testing.T) {
	deleted, selected := 0, 0
	var field *kvitui.Field
	var other *kvitui.Button
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		field = kvitui.NewField(ui)
		field.Label = "Payee"
		other = kvitui.NewButton(ui, "Import")
		return []unison.Paneler{kvitui.Left(field), other}
	})
	screen.Do(func() {
		err := w.AddCommands(
			&kvitui.Command{Std: kvitui.StdDelete, Run: func() { deleted++ }},
			&kvitui.Command{Std: kvitui.StdSelectAll, Run: func() { selected++ }})
		if err != nil {
			t.Error(err)
		}
		field.SetText("Acme")
	})
	screen.Click(screen.PanelCenter(field.Edit()))
	screen.KeyPress(unison.KeyHome, 0)
	screen.KeyPress(unison.KeyDelete, 0)
	screen.KeyPress(unison.KeyA, mod.Control)
	if deleted != 0 || selected != 0 || field.Text() != "cme" {
		t.Errorf("in a field Delete and Ctrl+A ran the app's commands %d and %d times and left %q", deleted, selected, field.Text())
	}
	screen.Do(func() { other.Focus() })
	screen.KeyPress(unison.KeyDelete, 0)
	screen.KeyPress(unison.KeyA, mod.Control)
	if deleted != 1 || selected != 1 {
		t.Errorf("on a button Delete and Ctrl+A ran the app's commands %d and %d times", deleted, selected)
	}
}

// A window refuses a second command on a key, another command on a
// standard key, and the keys the system and the access keys own.
func TestAWindowRefusesKeysThatAreTaken(t *testing.T) {
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler { return nil })
	refused := func(c *kvitui.Command) bool {
		var err error
		screen.Do(func() { err = w.AddCommands(c) })
		return err != nil
	}
	run := func() {}
	if refused(&kvitui.Command{Std: kvitui.StdSave, Run: run}) {
		t.Error("Save on its own key was refused")
	}
	for _, c := range []*kvitui.Command{
		{Name: "Sync", Shortcut: kvitui.Keys(unison.KeyS, kvitui.ModPrimary), Run: run},
		{Name: "Copy path", Shortcut: kvitui.Keys(unison.KeyC, kvitui.ModPrimary), Run: run},
		{Name: "Archive", Shortcut: kvitui.Keys(unison.KeyA, kvitui.ModAlt), Run: run},
		{Name: "Quick capture", Shortcut: kvitui.Keys(unison.KeyQ, kvitui.ModPrimary), Run: run},
		{Name: "Spotlight", Shortcut: kvitui.Keys(unison.KeySpace, kvitui.ModPrimary), Run: run},
		{Name: "Insert row", Shortcut: kvitui.Keys(unison.KeyI, kvitui.ModPrimary|kvitui.ModAlt), Run: run},
	} {
		if !refused(c) {
			t.Errorf("%s on its key was not refused", c.Name)
		}
	}
	if refused(&kvitui.Command{Name: "Quick capture", Shortcut: kvitui.Keys(unison.KeyK, kvitui.ModPrimary|kvitui.ModShift), Run: run}) {
		t.Error("Ctrl+Shift+K was refused")
	}
}

// A command made into a button, a symbol button and a menu line has one
// name, symbol and key in all three, in the form's case, with the shortcut
// in the symbol button's tooltip on the Windows form only, and is disabled
// in all three at once.
func TestACommandLooksTheSameEverywhere(t *testing.T) {
	can := true
	exp := &kvitui.Command{Name: "Export statement", Asks: true, Symbol: "external",
		Shortcut: kvitui.Keys(unison.KeyE, kvitui.ModPrimary|kvitui.ModShift),
		Enabled:  func() bool { return can }, Run: func() {}}
	var button *kvitui.Button
	var icon *kvitui.IconButton
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewCommandButton(ui, exp)
		icon = kvitui.NewCommandIconButton(ui, exp)
		return []unison.Paneler{button, icon}
	})
	name := func(p unison.Paneler) string { return screen.AccessibilityNodeFor(p).Name }
	if got := name(button); got != "Export statement…" {
		t.Errorf("the button says %q", got)
	}
	if got := name(icon); got != "Export statement" {
		t.Errorf("the symbol button is named %q", got)
	}
	screen.Do(func() { icon.Focus() })
	waitFor(t, screen, "the tooltip", func() bool { return len(w.Popups()) == 1 })
	tip := ""
	screen.Do(func() { tip = screen.AccessibilityNodeFor(w.Popups()[0].Panel).Name })
	if tip != "Export statement… (Ctrl+Shift+E)" {
		t.Errorf("the Windows form's tooltip is %q", tip)
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	if got := name(button); got != "Export Statement…" {
		t.Errorf("on the macOS form the button says %q", got)
	}
	screen.Do(func() {
		ui.SetForm(kvitui.FormWindows)
		ui.ShowMenu(button, "Statement", []kvitui.MenuItem{{Command: exp}})
	})
	screen.Sync()
	lines := nodes(screen, w, role.MenuItem)
	if len(lines) != 1 || lines[0].Name != "Export statement…" || !strings.Contains(lines[0].Description, "Ctrl+Shift+E") {
		t.Errorf("the menu line is %+v", lines)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	can = false
	screen.Do(func() { ui.Refresh() })
	screen.Sync()
	if !screen.AccessibilityNodeFor(button).Disabled || !screen.AccessibilityNodeFor(icon).Disabled {
		t.Error("a command that cannot act left its buttons enabled")
	}
}

// A context menu leaves out what cannot act, and on the macOS form shows
// no shortcuts but keeps Cut, Copy and Paste; a destructive line at the end
// gets a separator before it.
func TestAContextMenuShowsWhatApplies(t *testing.T) {
	var target *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		target = kvitui.NewButton(ui, "Row")
		return []unison.Paneler{target}
	})
	never := func() bool { return false }
	screen.Do(func() {
		ui.SetContextMenu(target, func(geom.Point) (string, []kvitui.MenuItem) {
			return "", []kvitui.MenuItem{
				{Command: &kvitui.Command{Std: kvitui.StdCopy, Run: func() {}}},
				{Command: &kvitui.Command{Std: kvitui.StdPaste, Enabled: never, Run: func() {}}},
				{Command: &kvitui.Command{Std: kvitui.StdRename, Enabled: never, Run: func() {}}},
				{Text: "Delete row", Danger: true},
			}
		})
	})
	open := func() []string {
		screen.Do(func() { target.Focus() })
		screen.KeyPress(unison.KeyMenu, 0)
		var out []string
		for _, n := range nodes(screen, w, role.MenuItem) {
			out = append(out, n.Name+"|"+n.Description)
		}
		screen.KeyPress(unison.KeyEscape, 0)
		return out
	}
	if got := open(); !slices.Equal(got, []string{"Copy|Ctrl+C", "Delete row|"}) {
		t.Errorf("on the Windows form the context menu is %v", got)
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	if got := open(); !slices.Equal(got, []string{"Copy|", "Paste|", "Delete Row|"}) {
		t.Errorf("on the macOS form the context menu is %v", got)
	}
}
