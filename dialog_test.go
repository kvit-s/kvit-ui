package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// isInDialog reports whether a node is inside the open dialog.
func isInDialog(screen *unison.HeadlessScreen, w *kvitui.Window, n *accessibility.Node) bool {
	for _, d := range nodes(screen, w, role.Dialog) {
		if n.Bounds.Intersects(d.Bounds) {
			return true
		}
	}
	return false
}

// A dialog's buttons are in the form's order, its close symbol is on the
// Windows form only, and an informing dialog has Close on the Windows form
// and OK on the macOS form (WIN-1, WIN-6, WRD-7).
func TestADialogFollowsTheForm(t *testing.T) {
	var opener *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		opener = kvitui.NewButton(ui, "Open")
		return []unison.Paneler{opener}
	})
	show := func(d func() *kvitui.Dialog) []string {
		screen.Do(func() { d().Open(w) })
		screen.Sync()
		var out []string
		for _, n := range nodes(screen, w, role.Button) {
			if isInDialog(screen, w, n) {
				out = append(out, n.Name)
			}
		}
		screen.KeyPress(unison.KeyEscape, 0)
		return out
	}
	ask := func() *kvitui.Dialog {
		d := kvitui.NewDialog(ui, "Save changes to the ledger?")
		d.ConfirmText, d.OtherText = "Save", "Don't save"
		return d
	}
	inform := func() *kvitui.Dialog { return kvitui.NewMessage(ui, "Import finished", "212 transactions were added.") }
	same := func(got []string, want ...string) bool {
		if len(got) != len(want) {
			return false
		}
		for i := range got {
			if got[i] != want[i] {
				return false
			}
		}
		return true
	}
	if got := show(ask); !same(got, "Close", "Save", "Don't save", "Cancel") {
		t.Errorf("on the Windows form the buttons are %v", got)
	}
	if got := show(inform); !same(got, "Close", "Close") {
		t.Errorf("on the Windows form a message has %v", got)
	}
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	if got := show(ask); !same(got, "Don't Save", "Cancel", "Save") {
		t.Errorf("on the macOS form the buttons are %v", got)
	}
	if got := show(inform); !same(got, "OK") {
		t.Errorf("on the macOS form a message has %v", got)
	}
}

// Return presses the default button: the action on the macOS form unless
// the dialog is destructive, and on the Windows form only where the app
// marks it, Enter pressing the focused button instead (KEY-6). Escape, and
// Command-period on the macOS form, cancel (KEY-7).
func TestReturnAndEscapeInADialogFollowTheForm(t *testing.T) {
	screen, ui, w := windowSession(t, func(*kvitui.UI) []unison.Paneler { return nil })
	var accepted, rejected int
	open := func(setup func(d *kvitui.Dialog)) {
		field := kvitui.NewField(ui)
		field.Label = "Name"
		screen.Do(func() {
			d := kvitui.NewDialog(ui, "Rename the ledger", field)
			d.ConfirmText = "Rename"
			d.OnAccept = func() { accepted++ }
			d.OnReject = func() { rejected++ }
			if setup != nil {
				setup(d)
			}
			d.Open(w)
		})
		screen.Sync()
	}
	reset := func() {
		accepted, rejected = 0, 0
		screen.Do(func() {
			for _, p := range w.Popups() {
				if d, ok := p.Panel.(*kvitui.Dialog); ok {
					d.Close()
				}
			}
		})
	}
	open(nil)
	if n := focusedNode(screen, w); n == nil || n.Name != "Name" {
		t.Errorf("the keyboard started on %v, not the dialog's first field", n)
	}
	screen.KeyPress(unison.KeyReturn, 0)
	if accepted != 0 {
		t.Error("on the Windows form Return in a field confirmed a dialog without a default")
	}
	screen.KeyPress(unison.KeyTab, 0)
	screen.KeyPress(unison.KeyTab, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Cancel" {
		t.Fatalf("two Tabs from the field went to %v, not Cancel", n)
	}
	screen.KeyPress(unison.KeyReturn, 0)
	if rejected != 1 {
		t.Error("on the Windows form Enter on the focused Cancel did not press it")
	}
	reset()
	open(func(d *kvitui.Dialog) { d.DefaultConfirm = true })
	screen.KeyPress(unison.KeyReturn, 0)
	if accepted != 1 {
		t.Error("on the Windows form Return in a field did not press the default button")
	}
	reset()
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	open(nil)
	screen.KeyPress(unison.KeyReturn, 0)
	if accepted != 1 {
		t.Error("on the macOS form Return did not press the default button")
	}
	reset()
	open(func(d *kvitui.Dialog) { d.Destructive = true })
	screen.KeyPress(unison.KeyReturn, 0)
	screen.KeyPress(unison.KeyPeriod, mod.Command)
	if accepted != 0 || rejected != 1 || len(w.Popups()) != 0 {
		t.Errorf("on the macOS form a destructive dialog took Return %d times, and Command-period cancelled %d times", accepted, rejected)
	}
}

// The keyboard stays in an open dialog, and goes back to where it was when
// the dialog closes (KEY-8).
func TestTheKeyboardStaysInADialogAndGoesBack(t *testing.T) {
	var opener, behind *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		opener = kvitui.NewButton(ui, "Rename")
		behind = kvitui.NewButton(ui, "Behind")
		return []unison.Paneler{opener, behind}
	})
	var d *kvitui.Dialog
	screen.Do(func() { opener.Focus() })
	screen.Sync()
	screen.Do(func() {
		d = kvitui.NewDialog(ui, "Rename the ledger")
		d.ConfirmText = "Rename"
		d.Open(w)
	})
	screen.Sync()
	for range 6 {
		screen.KeyPress(unison.KeyTab, 0)
		if n := focusedNode(screen, w); n == nil || !isInDialog(screen, w, n) {
			t.Fatalf("Tab left the dialog for %v", n)
		}
	}
	screen.KeyPress(unison.KeyEscape, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Rename" {
		t.Errorf("closing the dialog left the keyboard on %v, not its opener", n)
	}
}

// The close symbol and Escape ask OnDiscard before throwing away the
// reader's input, and Cancel does not (WIN-8); a second dialog is refused
// while one is open, and a message may open over it (WIN-7).
func TestADialogKeepsInputAndOpensOneAtATime(t *testing.T) {
	screen, ui, w := windowSession(t, func(*kvitui.UI) []unison.Paneler { return nil })
	asked, keep := 0, true
	var d *kvitui.Dialog
	screen.Do(func() {
		d = kvitui.NewDialog(ui, "New account")
		d.ConfirmText = "Create"
		d.OnDiscard = func() bool { asked++; return !keep }
		d.Open(w)
	})
	screen.Sync()
	screen.KeyPress(unison.KeyEscape, 0)
	if asked != 1 || len(w.Popups()) != 1 {
		t.Errorf("Escape asked %d times and left %d popups, want the dialog kept", asked, len(w.Popups()))
	}
	screen.Do(func() {
		second := kvitui.NewDialog(ui, "Import statement")
		second.ConfirmText = "Import"
		second.Open(w)
	})
	if len(w.Popups()) != 1 {
		t.Errorf("a second dialog opened over the first: %d popups", len(w.Popups()))
	}
	screen.Do(func() { kvitui.NewMessage(ui, "Account names must be unique", "").Open(w) })
	if len(w.Popups()) != 2 {
		t.Errorf("a message did not open over the dialog: %d popups", len(w.Popups()))
	}
	screen.KeyPress(unison.KeyEscape, 0)
	screen.Sync()
	cancel := ""
	for _, n := range nodes(screen, w, role.Button) {
		if n.Name == "Cancel" {
			cancel = n.Name
			screen.Click(n.Bounds.Center())
		}
	}
	if cancel == "" || asked != 1 || len(w.Popups()) != 0 {
		t.Errorf("Cancel asked %d times and left %d popups", asked-1, len(w.Popups()))
	}
}

// A dialog's help is the question-worded link at the bottom of its content
// on the Windows form, which Tab reaches after the buttons, and the round
// help button at the leading end of the button row on the macOS form
// (HLP-6, KEY-9).
func TestADialogPlacesItsHelpByForm(t *testing.T) {
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		return []unison.Paneler{kvitui.NewButton(ui, "Open")}
	})
	open := func() {
		screen.Do(func() {
			d := kvitui.NewDialog(ui, "Roll the budget over?")
			d.ConfirmText, d.OtherText = "Roll over", "Start empty"
			d.Help = kvitui.NewHelpLink(ui, "budgets#rollover", "How does a budget roll over?")
			d.Open(w)
		})
		screen.Sync()
	}
	open()
	var link *accessibility.Node
	var buttons []*accessibility.Node
	for _, n := range nodes(screen, w, role.Link) {
		if isInDialog(screen, w, n) && n.Name == "How does a budget roll over?" {
			link = n
		}
	}
	for _, n := range nodes(screen, w, role.Button) {
		if isInDialog(screen, w, n) && n.Name != "Close" {
			buttons = append(buttons, n)
		}
	}
	if link == nil || len(buttons) != 3 {
		t.Fatalf("on the Windows form the dialog has the link %v and the buttons %v", link, buttons)
	}
	for _, b := range buttons {
		if link.Bounds.Bottom() > b.Bounds.Y {
			t.Errorf("the help link reaches %v, below the top of %q at %v", link.Bounds.Bottom(), b.Name, b.Bounds.Y)
		}
	}
	// The keyboard starts on Cancel, the last button; Tab goes to the link.
	screen.KeyPress(unison.KeyTab, 0)
	if n := focusedNode(screen, w); n == nil || n.Role != role.Link {
		t.Errorf("Tab after the last button reached %v, not the help link", n)
	}
	screen.KeyPress(unison.KeyEscape, 0)

	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	open()
	var help *accessibility.Node
	buttons = nil
	for _, n := range nodes(screen, w, role.Button) {
		if !isInDialog(screen, w, n) {
			continue
		}
		if n.Name == "Help" {
			help = n
		} else {
			buttons = append(buttons, n)
		}
	}
	if help == nil || len(buttons) != 3 {
		t.Fatalf("on the macOS form the dialog has the help button %v and the buttons %v", help, buttons)
	}
	for _, b := range buttons {
		if help.Bounds.Right() > b.Bounds.X || help.Bounds.Y > b.Bounds.Bottom() || help.Bounds.Bottom() < b.Bounds.Y {
			t.Errorf("the help button at %v is not first in the row before %q at %v", help.Bounds, b.Name, b.Bounds)
		}
	}
}

// A properties dialog has OK, Cancel and Apply on the Windows form, Apply
// enabled once something changed and again after it is pressed, which
// leaves the dialog open; the macOS form leaves Apply out (WIN-23).
func TestApplyWaitsForAChange(t *testing.T) {
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		return []unison.Paneler{kvitui.NewButton(ui, "Open")}
	})
	var d *kvitui.Dialog
	applied := 0
	open := func() {
		screen.Do(func() {
			name := kvitui.NewField(ui)
			name.Label = "Name"
			d = kvitui.NewDialog(ui, "Checking properties", kvitui.FullWidth(name))
			d.Kind = kvitui.DialogProperties
			d.OnApply = func() { applied++ }
			name.OnChange = func(string) { d.MarkChanged() }
			d.Open(w)
		})
		screen.Sync()
	}
	apply := func() *accessibility.Node {
		for _, n := range nodes(screen, w, role.Button) {
			if n.Name == "Apply" && isInDialog(screen, w, n) {
				return n
			}
		}
		return nil
	}
	open()
	if got := d.ButtonWords(); len(got) != 3 || got[0] != "OK" || got[2] != "Apply" {
		t.Fatalf("on the Windows form a properties dialog has %v", got)
	}
	if n := apply(); n == nil || !n.Disabled {
		t.Fatalf("Apply before any change is %+v", n)
	}
	screen.Type("Joint")
	n := apply()
	if n == nil || n.Disabled {
		t.Fatalf("Apply after a change is %+v", n)
	}
	screen.Click(n.Bounds.Center())
	if applied != 1 || d.Window() == nil {
		t.Errorf("Apply ran OnApply %d times and left the dialog open: %v", applied, d.Window() != nil)
	}
	if n := apply(); n == nil || !n.Disabled {
		t.Errorf("Apply after it was pressed is %+v", n)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	open()
	if got := d.ButtonWords(); len(got) != 2 || got[1] != "OK" {
		t.Errorf("on the macOS form a properties dialog has %v", got)
	}
	screen.KeyPress(unison.KeyEscape, 0)
}

// A dialog hands out its buttons, so the application enables or presses them
// without finding them by their words.
func TestADialogHandsOutItsButtons(t *testing.T) {
	screen, ui, w := windowSession(t, func(*kvitui.UI) []unison.Paneler { return nil })
	var d *kvitui.Dialog
	screen.Do(func() {
		d = kvitui.NewDialog(ui, "Save changes to the ledger?")
		d.ConfirmText, d.OtherText = "Save", "Don't save"
		d.Open(w)
	})
	screen.Sync()
	var buttons []*kvitui.Button
	var words []string
	var confirm, other, cancel, apply *kvitui.Button
	screen.Do(func() {
		buttons = d.Buttons()
		words = d.ButtonWords()
		confirm, other, cancel, apply = d.ConfirmButton(), d.OtherButton(), d.CancelButton(), d.ApplyButton()
	})
	if len(buttons) != 3 || len(words) != 3 {
		t.Fatalf("an asking dialog shows %d buttons named %v", len(buttons), words)
	}
	for i, b := range buttons {
		if b.Words() != words[i] {
			t.Fatalf("Buttons()[%d] says %q, ButtonWords()[%d] says %q", i, b.Words(), i, words[i])
		}
	}
	if buttons[0] != confirm || buttons[1] != other || buttons[2] != cancel {
		t.Fatalf("on the Windows form the buttons are not the confirm, other and cancel buttons in order")
	}
	if apply != nil {
		t.Fatalf("an asking dialog has an apply button %q", apply.Words())
	}
	// The handle acts on the button: disabling the confirming button through
	// it disables the Save button the reader sees.
	screen.Do(func() { confirm.SetEnabled(false) })
	screen.Sync()
	disabled := false
	for _, n := range nodes(screen, w, role.Button) {
		if isInDialog(screen, w, n) && n.Name == "Save" {
			disabled = n.Disabled
		}
	}
	if !disabled {
		t.Error("disabling the confirming button did not disable the Save button")
	}
	screen.Do(func() { confirm.SetEnabled(true) })

	screen.Do(func() { ui.SetForm(kvitui.FormMac) })
	screen.Sync()
	screen.Do(func() {
		buttons = d.Buttons()
		confirm, other, cancel = d.ConfirmButton(), d.OtherButton(), d.CancelButton()
	})
	if len(buttons) != 3 || buttons[0] != other || buttons[1] != cancel || buttons[2] != confirm {
		t.Errorf("on the macOS form the buttons are not the other, cancel and confirm buttons in order")
	}
	// The form leaves Apply out on the macOS form, Done takes the action away
	// and an informing dialog's only button is its confirming one.
	screen.Do(func() {
		ui.SetForm(kvitui.FormWindows)
		p := kvitui.NewDialog(ui, "Checking properties")
		p.Kind = kvitui.DialogProperties
		if p.ApplyButton() == nil {
			t.Error("a properties dialog has no apply button on the Windows form")
		}
		ui.SetForm(kvitui.FormMac)
		if p.ApplyButton() != nil {
			t.Error("a properties dialog keeps its apply button on the macOS form")
		}
		ui.SetForm(kvitui.FormWindows)
		d.Done()
		if d.ConfirmButton() != nil || d.OtherButton() != nil || d.CancelButton() == nil {
			t.Error("once done, an asking dialog keeps no action and only its Close button")
		}
		m := kvitui.NewMessage(ui, "Import finished", "212 transactions were added.")
		ui.SetForm(kvitui.FormMac)
		if m.ConfirmButton() == nil || m.CancelButton() != nil || len(m.Buttons()) != 1 {
			t.Error("an informing dialog's only button is not its confirming button")
		}
	})
}
