package kvitui

import "github.com/richardwilkes/unison"

// KeyTurn reports whether a window is handling a key press now, so that
// focus gained now was moved by a key. It is false for a press of the
// pointer and for the focus a window hands its first control when it opens.
func (u *UI) KeyTurn() bool { return u.keyTurn }

// WatchKeys makes a window report its key presses to KeyTurn, and draw the
// focus ring of whichever Kvit control holds its keyboard focus. A Kvit
// window, and any window a Kvit control is laid out in, is watched already;
// this is for a panel of an application's own in any other unison window.
// Calling it again for the same window does nothing.
//
// The watch wraps the window's KeyDownCallback and calls the one it
// replaced. An application that sets its own afterwards has to call the one
// it replaces.
func (u *UI) WatchKeys(w *unison.Window) {
	if w != nil {
		u.watchWindow(w)
	}
}

// KeyFocus tracks, for a panel that takes the focus itself, whether the
// focus it holds came from the keyboard, by the rule Kvit's controls use: a
// key moved it there, or the component asked for it as an arrow key does,
// and a press of the pointer did not. The panel shows its focus ring while
// Shows answers true.
//
// Call Pressed from the panel's MouseDownCallback before it asks for the
// focus, Asked before the component gives the panel the focus itself, Gained
// from its GainedFocusCallback and Lost from its LostFocusCallback.
type KeyFocus struct {
	pressed, asked, keyboard bool
}

// Pressed says a press is about to give the panel the focus. A press on a
// panel that holds the focus already ends any keyboard focus shown.
func (f *KeyFocus) Pressed() {
	f.pressed = true
	f.keyboard = false
}

// Asked says the component is about to give the panel the focus itself, as
// an arrow key in a list does.
func (f *KeyFocus) Asked() { f.asked = true }

// Gained records where the focus just gained came from.
func (f *KeyFocus) Gained(ui *UI) {
	f.keyboard = !f.pressed && (f.asked || ui.keyTurn)
	f.pressed, f.asked = false, false
}

// Lost records that the focus has gone.
func (f *KeyFocus) Lost() { f.pressed, f.asked, f.keyboard = false, false, false }

// Shows reports whether p holds the focus and it came from the keyboard.
func (f *KeyFocus) Shows(p *unison.Panel) bool { return f.keyboard && p != nil && p.Focused() }
