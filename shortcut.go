package kvitui

import (
	"strings"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// Mods are the modifier keys of a Shortcut, named for what they do on both
// forms rather than for one keyboard, so a shortcut is written once.
type Mods uint8

const (
	// ModPrimary is the command key: Ctrl on the Windows form, Command on
	// the macOS form (KEY-3).
	ModPrimary Mods = 1 << iota
	// ModShift is Shift, which adds a related variant: Ctrl+Shift+Z redoes
	// what Ctrl+Z undid.
	ModShift
	// ModAlt is Alt on the Windows form and Option on the macOS form. On the
	// Windows form Alt with a letter or a digit belongs to access keys and is
	// never a shortcut (KEY-15).
	ModAlt
	// ModControl is the Control key itself on the macOS form, beside
	// Command, as in Ctrl-Cmd-F for full screen. On the Windows form it is
	// Ctrl, the same key as ModPrimary.
	ModControl
)

// Shortcut is a key and its modifiers, written once for both forms. A
// shortcut on KeyEqual is the key with the plus sign on it: it answers with
// or without Shift and on the keypad's plus, and is written as a plus, as
// Ctrl+Plus zooms in on every keyboard. KeyMinus answers on the keypad's
// minus too.
type Shortcut struct {
	Key  unison.KeyCode
	Mods Mods
}

// Keys returns a shortcut: Keys(unison.KeyS, ModPrimary) is Ctrl+S on the
// Windows form and Cmd-S on the macOS form.
func Keys(key unison.KeyCode, mods Mods) Shortcut { return Shortcut{Key: key, Mods: mods} }

// IsZero reports whether the shortcut is no shortcut.
func (s Shortcut) IsZero() bool { return s.Key == 0 }

// modifiers is what unison reports while the shortcut's modifiers are held,
// in a form.
func (s Shortcut) modifiers(f Form) mod.Modifiers {
	var m mod.Modifiers
	if s.Mods&ModShift != 0 {
		m |= mod.Shift
	}
	if s.Mods&ModAlt != 0 {
		m |= mod.Option
	}
	if f == FormMac {
		if s.Mods&ModPrimary != 0 {
			m |= mod.Command
		}
		if s.Mods&ModControl != 0 {
			m |= mod.Control
		}
	} else if s.Mods&(ModPrimary|ModControl) != 0 {
		m |= mod.Control
	}
	return m
}

// Binding is the shortcut as unison's key binding in a form, for a menu
// unison draws.
func (s Shortcut) Binding(f Form) unison.KeyBinding {
	if s.IsZero() {
		return unison.KeyBinding{}
	}
	return unison.KeyBinding{KeyCode: s.Key, Modifiers: s.modifiers(f)}
}

// sameKey turns a keypad key into its twin on the main keyboard, and Return
// on the keypad into Return, since they mean the same (KEY-3).
func sameKey(key unison.KeyCode) unison.KeyCode {
	switch {
	case key >= unison.KeyNumPad0 && key <= unison.KeyNumPad9:
		return unison.Key0 + (key - unison.KeyNumPad0)
	}
	switch key {
	case unison.KeyNumPadEnter:
		return unison.KeyReturn
	case unison.KeyNumPadAdd:
		return unison.KeyEqual
	case unison.KeyNumPadSubtract:
		return unison.KeyMinus
	case unison.KeyNumPadDecimal:
		return unison.KeyPeriod
	case unison.KeyNumPadDivide:
		return unison.KeySlash
	}
	return key
}

// Matches reports whether a key press is the shortcut in a form. Caps Lock
// and Num Lock are ignored, and a keypad key is its twin.
func (s Shortcut) Matches(f Form, key unison.KeyCode, mods mod.Modifiers) bool {
	if s.IsZero() || sameKey(key) != sameKey(s.Key) {
		return false
	}
	held := mods & mod.NonSticky
	want := s.modifiers(f)
	if sameKey(s.Key) == unison.KeyEqual {
		// The plus is the shifted equals sign on most keyboards.
		held &^= mod.Shift
		want &^= mod.Shift
	}
	return held == want
}

// ShortcutText is how a shortcut is written in the app's form: "Ctrl+Shift+S"
// on the Windows form, "⇧⌘S" on the macOS form.
func (u *UI) ShortcutText(s Shortcut) string { return shortcutText(u.form, s) }

// ShortcutTextIn is how a shortcut is written in a form.
func ShortcutTextIn(f Form, s Shortcut) string { return shortcutText(f, s) }

func shortcutText(f Form, s Shortcut) string {
	if s.IsZero() {
		return ""
	}
	key, mods := keyName(f, s.Key), s.Mods
	if sameKey(s.Key) == unison.KeyEqual {
		// The plus is written as itself, without the Shift that types it.
		mods &^= ModShift
	}
	if f == FormMac && sameKey(s.Key) == unison.KeySlash && mods&ModShift != 0 {
		// Help's Cmd-? is written with the question mark Shift types.
		key, mods = "?", mods&^ModShift
	}
	if f == FormMac {
		// Apple's order: Control, Option, Shift, Command.
		var b strings.Builder
		if mods&ModControl != 0 {
			b.WriteString("⌃")
		}
		if mods&ModAlt != 0 {
			b.WriteString("⌥")
		}
		if mods&ModShift != 0 {
			b.WriteString("⇧")
		}
		if mods&ModPrimary != 0 {
			b.WriteString("⌘")
		}
		return b.String() + key
	}
	var parts []string
	if mods&(ModPrimary|ModControl) != 0 {
		parts = append(parts, "Ctrl")
	}
	if mods&ModAlt != 0 {
		parts = append(parts, "Alt")
	}
	if mods&ModShift != 0 {
		parts = append(parts, "Shift")
	}
	return strings.Join(append(parts, key), "+")
}

// keyName is how a key is written in a form: the names Windows' menus use,
// and the symbols on a Mac's keys.
func keyName(f Form, key unison.KeyCode) string {
	key = sameKey(key)
	mac := f == FormMac
	pick := func(windows, macOS string) string {
		if mac {
			return macOS
		}
		return windows
	}
	switch key {
	case unison.KeyEqual:
		return pick("Plus", "+")
	case unison.KeyMinus:
		return pick("Minus", "-")
	case unison.KeyReturn:
		return pick("Enter", "↩")
	case unison.KeyEscape:
		return pick("Esc", "⎋")
	case unison.KeyTab:
		return pick("Tab", "⇥")
	case unison.KeySpace:
		return "Space"
	case unison.KeyBackspace:
		return pick("Backspace", "⌫")
	case unison.KeyDelete:
		return pick("Del", "⌦")
	case unison.KeyInsert:
		return "Ins"
	case unison.KeyHome:
		return pick("Home", "↖")
	case unison.KeyEnd:
		return pick("End", "↘")
	case unison.KeyPageUp:
		return pick("PgUp", "⇞")
	case unison.KeyPageDown:
		return pick("PgDn", "⇟")
	case unison.KeyLeft:
		return pick("Left", "←")
	case unison.KeyRight:
		return pick("Right", "→")
	case unison.KeyUp:
		return pick("Up", "↑")
	case unison.KeyDown:
		return pick("Down", "↓")
	}
	return key.String()
}
