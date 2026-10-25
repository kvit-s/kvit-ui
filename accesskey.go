package kvitui

import (
	"slices"
	"strings"
	"unicode"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// On the Windows form, which Linux takes, every control with words has an
// access key: Alt and a letter of its words, unique in the window, which
// operates the control at once, or moves the keyboard to it when several
// share the letter (KEY-15, decided Q-11). The letter is underlined while
// Alt is held. An app chooses the letter by writing "&" before it in the
// words, "&Import"; a control without one gets the first free letter of its
// words, the first letters of the words first. OK, Cancel and Close take
// none (Enter and Escape serve them), nor do links, tabs, symbol buttons,
// toolbar items and static text. The macOS form has no access keys, and
// draws the words without the marker.

// accessKeyed is a control with words an access key can operate.
type accessKeyed interface {
	// accessWords are the control's words as written, with an "&" before
	// the letter the app chose, or "" for a control that takes no key now.
	accessWords() string
	// accessPress is what Alt and the key do: press, toggle or choose.
	accessPress()
	AsPanel() *unison.Panel
}

// AccessKey is the access key a control has in a window.
type AccessKey struct {
	// Panel is the control.
	Panel *unison.Panel
	// Words are its words as drawn.
	Words string
	// Key is the letter, in lower case; 0 for a control left without one.
	Key rune
	// Chosen says the app chose the letter with "&".
	Chosen bool
	// at is the letter's place in Words, in runes.
	at int
}

// noAccessKey are the words of buttons that take no access key, since Enter
// and Escape press them.
var noAccessKey = map[string]bool{"OK": true, "Cancel": true, "Close": true}

// AccessKeys are the access keys of the controls in the window's keyboard
// scope, the topmost dialog when one is open, on the Windows form; none on
// the macOS form. A key the menu bar's titles use is not given to a control.
func (w *Window) AccessKeys() []AccessKey {
	if w.ui.form != FormWindows {
		return nil
	}
	scope, modal := w.accessScope()
	used := map[rune]bool{}
	if !modal && w.menuBar != nil && !w.menuBar.isNative() {
		for _, m := range w.menuBar.Menus() {
			if _, k, _ := AccessText(m.Title); k != 0 {
				used[k] = true
			}
		}
	}
	var keys []AccessKey
	var walk func(p *unison.Panel)
	walk = func(p *unison.Panel) {
		if p.Hidden {
			return
		}
		if _, ok := p.Self.(*Toolbar); ok {
			return
		}
		if a, ok := p.Self.(accessKeyed); ok {
			if words := a.accessWords(); words != "" {
				plain, k, at := AccessText(words)
				keys = append(keys, AccessKey{Panel: p, Words: plain, Key: k, Chosen: k != 0, at: at})
			}
		}
		for _, c := range p.Children() {
			walk(c)
		}
	}
	walk(scope)
	for _, k := range keys {
		if k.Key != 0 {
			used[k.Key] = true
		}
	}
	for i, k := range keys {
		if k.Key == 0 {
			keys[i].Key, keys[i].at = freeLetter(k.Words, used)
			if keys[i].Key != 0 {
				used[keys[i].Key] = true
			}
		}
	}
	return keys
}

// accessScope is where access keys reach: the topmost modal popup, or the
// whole window.
func (w *Window) accessScope() (scope *unison.Panel, modal bool) {
	for i := len(w.popups) - 1; i >= 0; i-- {
		if w.popups[i].Modal {
			return w.popups[i].Panel.AsPanel(), true
		}
	}
	return w.Content(), false
}

// freeLetter is the first letter of words not yet used: the first letters
// of the words first, then any letter or digit in order. When every letter
// is taken the first letter is shared, which Windows allows: Alt and a
// shared letter moves the keyboard between the controls that have it.
func freeLetter(words string, used map[rune]bool) (rune, int) {
	if r, at := unusedLetter(words, used); r != 0 {
		return r, at
	}
	for i, r := range []rune(words) {
		if unicode.IsLetter(r) && r < unicode.MaxASCII {
			return unicode.ToLower(r), i
		}
	}
	return 0, -1
}

// unusedLetter is freeLetter without the sharing.
func unusedLetter(words string, used map[rune]bool) (rune, int) {
	rs := []rune(words)
	pick := func(i int) (rune, bool) {
		r := unicode.ToLower(rs[i])
		ok := (unicode.IsLetter(r) || unicode.IsDigit(r)) && r < unicode.MaxASCII && !used[r]
		return r, ok
	}
	for i := range rs {
		if (i == 0 || rs[i-1] == ' ') && unicode.IsLetter(rs[i]) {
			if r, ok := pick(i); ok {
				return r, i
			}
		}
	}
	for i := range rs {
		if r, ok := pick(i); ok {
			return r, i
		}
	}
	return 0, -1
}

// accessKeyPress operates the control whose access key Alt and a letter are,
// or moves the keyboard to the next of several sharing it.
func (w *Window) accessKeyPress(key unison.KeyCode, mods mod.Modifiers) bool {
	if w.ui.form != FormWindows || mods&mod.NonSticky != mod.Option {
		return false
	}
	typed, ok := typedLetter(key)
	if !ok {
		return false
	}
	var hits []*unison.Panel
	for _, k := range w.AccessKeys() {
		if k.Key == typed && k.Panel.Enabled() {
			hits = append(hits, k.Panel)
		}
	}
	switch len(hits) {
	case 0:
		return false
	case 1:
		if a, ok := hits[0].Self.(accessKeyed); ok {
			if c, ok := hits[0].Self.(interface{ Focus() }); ok {
				c.Focus()
			}
			a.accessPress()
		}
		return true
	}
	next := 0
	if at := slices.Index(hits, w.CurrentFocus()); at >= 0 {
		next = (at + 1) % len(hits)
	}
	if c, ok := hits[next].Self.(interface{ Focus() }); ok {
		c.Focus()
	} else {
		hits[next].RequestFocus()
	}
	return true
}

// accessSpans are a control's words as drawn, with its access key
// underlined while Alt is held in its window on the Windows form, or always
// where the system asks for that.
func (u *UI) accessSpans(p *unison.Panel, words string, st text.Style) []text.Span {
	plain := []text.Span{{Text: words, Style: st}}
	w := u.windowOf(p)
	if w == nil || !w.ShowsAccessKeys() {
		return plain
	}
	for _, k := range w.accessKeysNow() {
		if k.Panel != p || k.at < 0 {
			continue
		}
		rs := []rune(words)
		if k.at >= len(rs) {
			return plain
		}
		key := st
		key.Underline = true
		return []text.Span{{Text: string(rs[:k.at]), Style: st}, {Text: string(rs[k.at : k.at+1]), Style: key}, {Text: string(rs[k.at+1:]), Style: st}}
	}
	return plain
}

// ShowsAccessKeys reports whether the window underlines its access keys
// now: while Alt is held, or at all times where the system says so, on the
// Windows form.
func (w *Window) ShowsAccessKeys() bool {
	return w.ui.form == FormWindows && (w.altHeld || w.ui.UnderlineAccessKeys)
}

// accessKeysNow are the access keys worked out when Alt went down, or at
// the start of the frame while they are always shown, so a drawing does
// not work them out for every control.
func (w *Window) accessKeysNow() []AccessKey {
	if w.accessCache == nil {
		w.accessCache = w.AccessKeys()
	}
	return w.accessCache
}

// setAltHeld records whether Alt is held, which shows the access keys, and
// redraws the window when that changes.
func (w *Window) setAltHeld(held bool) {
	if held == w.altHeld {
		return
	}
	w.altHeld = held
	w.accessCache = nil
	w.MarkForRedraw()
}

// plainWords are words without their access-key marker.
func plainWords(s string) string {
	if !strings.Contains(s, "&") {
		return s
	}
	plain, _, _ := AccessText(s)
	return plain
}
