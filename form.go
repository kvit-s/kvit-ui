package kvitui

import "runtime"

// Form is which system's conventions a window follows where Apple's and
// Microsoft's guidelines differ: the order of a dialog's buttons, the case
// of labels, the command key and the shortcuts, where Settings and Exit or
// Quit go, access keys, the words for a toggle. Linux takes the Windows form
// (the owner's decision of 2026-10-23, docs/ui-guidelines Q-40).
//
// The form says what the app does, not what the system can do: the system's
// own menus on macOS, and which key deletes a word, follow the system the app
// runs on whatever the form, so a test can draw the macOS form on Linux.
type Form int

const (
	// FormWindows is Windows' conventions, used on Windows and Linux.
	FormWindows Form = iota
	// FormMac is the Mac's conventions.
	FormMac
)

// String is the form's name as the gallery shows it.
func (f Form) String() string {
	if f == FormMac {
		return "macOS"
	}
	return "Windows"
}

// SystemForm is the form of the system the app runs on.
func SystemForm() Form {
	if runtime.GOOS == "darwin" {
		return FormMac
	}
	return FormWindows
}

// Form is the form the app's windows follow now. Components read it at
// layout and draw time, as they read colours and sizes.
func (u *UI) Form() Form { return u.form }

// SetForm makes the app's windows follow a form, and lays them out and draws
// them again. An app leaves it at the system's; tests and the gallery set it
// to draw a screen in either form on any system.
func (u *UI) SetForm(f Form) {
	if f == u.form {
		return
	}
	u.form = f
	u.changed()
}

// Mac reports whether the app follows the macOS form.
func (u *UI) Mac() bool { return u.form == FormMac }
