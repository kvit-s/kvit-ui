package kvitui

import (
	"runtime"
	"unicode"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// Deleting a word in a field or text area. unison's field deletes one
// character whatever the modifiers, so the word keys are taken before it
// sees them. A word is letters, digits and "_", as for the field's own
// Option+arrow moves.

// optionDeletesWord is whether Option, rather than Ctrl, makes Backspace and
// Delete take a word, as on macOS. The tests set it.
var optionDeletesWord = runtime.GOOS == "darwin"

// deletesWord reports whether Backspace or Delete with mods takes a word:
// Ctrl on Windows and Linux, Option on macOS. AltGr, which arrives as Ctrl
// with Alt, does not.
func deletesWord(mods mod.Modifiers) bool {
	ctrl, alt := mods.OSMenuCommandDown(), mods.OptionDown()
	if optionDeletesWord {
		return alt && !ctrl
	}
	return ctrl && !alt
}

// deleteWord deletes back to the start of the word before the caret, or
// forward to the end of the word after it, as one undo step. A selection
// is deleted as it is.
func deleteWord(e *unison.Field, forward bool) {
	start, end := e.Selection()
	if start == end {
		r := []rune(e.Text())
		if forward {
			end = fieldWordRight(r, end)
		} else {
			start = fieldWordLeft(r, start)
		}
		if start == end {
			return
		}
		e.SetSelection(start, end)
	}
	e.Delete()
}

// fieldWordLeft is where a word delete back from pos stops: past the spaces
// and punctuation before pos, then past the word before those. A line break
// is never crossed; right after one, only the break goes.
func fieldWordLeft(r []rune, pos int) int {
	pos = min(max(pos, 0), len(r))
	if pos > 0 && r[pos-1] == '\n' {
		return pos - 1
	}
	for pos > 0 && r[pos-1] != '\n' && !isFieldWordRune(r[pos-1]) {
		pos--
	}
	for pos > 0 && isFieldWordRune(r[pos-1]) {
		pos--
	}
	return pos
}

// fieldWordRight is fieldWordLeft forward.
func fieldWordRight(r []rune, pos int) int {
	pos = min(max(pos, 0), len(r))
	if pos < len(r) && r[pos] == '\n' {
		return pos + 1
	}
	for pos < len(r) && r[pos] != '\n' && !isFieldWordRune(r[pos]) {
		pos++
	}
	for pos < len(r) && isFieldWordRune(r[pos]) {
		pos++
	}
	return pos
}

func isFieldWordRune(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' }
