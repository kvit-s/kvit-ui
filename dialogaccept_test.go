package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/unison"
)

// A dialog that stays open on confirming runs OnAccept and is still there,
// as an answer the application refuses leaves it; Close takes it away
// without answering it, so neither OnAccept nor OnReject runs again, and
// closing a closed dialog does nothing. A dialog that closes on confirming,
// as every dialog does unless asked, is gone before its OnAccept runs.
func TestADialogThatStaysOpenOnConfirmingIsClosedByTheApplication(t *testing.T) {
	var dialog *kvitui.Dialog
	accepted, rejected := 0, 0
	openWhenAccepted := false
	screen, ui, w := windowSession(t, func(*kvitui.UI) []unison.Paneler { return nil })
	open := func(keep bool) {
		screen.Do(func() {
			dialog = kvitui.NewDialog(ui, "Rename the ledger")
			dialog.ConfirmText, dialog.ConfirmKeepsOpen = "Rename", keep
			// Renaming is the likely and safe answer, so it is the default,
			// which holds the keyboard.
			dialog.DefaultConfirm = true
			dialog.OnAccept = func() {
				accepted++
				openWhenAccepted = len(w.Popups()) != 0
			}
			dialog.OnReject = func() { rejected++ }
			dialog.Open(w)
		})
		screen.Sync()
	}

	open(true)
	// The confirming button holds the keyboard, so Space presses it.
	screen.KeyPress(unison.KeySpace, 0)
	screen.Do(dialog.Accept)
	if accepted != 2 || !openWhenAccepted || len(w.Popups()) != 1 {
		t.Fatalf("confirmed twice, the dialog was accepted %d times and left %d popups", accepted, len(w.Popups()))
	}
	screen.Do(dialog.Close)
	screen.Sync()
	if len(w.Popups()) != 0 || accepted != 2 || rejected != 0 {
		t.Errorf("closed, the dialog left %d popups, accepted %d times and rejected %d", len(w.Popups()), accepted,
			rejected)
	}
	screen.Do(dialog.Close)
	screen.Do(func() { dialog.Open(w) })
	screen.Sync()
	if len(w.Popups()) != 1 {
		t.Errorf("a closed dialog opened again shows %d popups", len(w.Popups()))
	}
	screen.KeyPress(unison.KeyEscape, 0)

	accepted = 0
	open(false)
	screen.KeyPress(unison.KeySpace, 0)
	if accepted != 1 || openWhenAccepted || len(w.Popups()) != 0 {
		t.Errorf("a dialog that closes on confirming was accepted %d times, open then: %v, leaving %d popups",
			accepted, openWhenAccepted, len(w.Popups()))
	}
}
