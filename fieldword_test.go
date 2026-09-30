package kvitui

import (
	"testing"

	"github.com/richardwilkes/unison/enums/mod"
)

// A word delete takes the spaces and punctuation next to the caret, then the
// word beyond them, and never crosses a line break: right next to one it
// takes only the break.
func TestAWordDeleteStopsAtTheWordAndTheLine(t *testing.T) {
	r := []rune("one, two_2\nthree")
	for _, c := range []struct {
		pos, left, right int
	}{
		{10, 5, 11},  // after "two_2": back to "t", forward over the break only
		{5, 0, 10},   // before "two_2": back over ", " and "one"
		{11, 10, 16}, // after the break: back over the break only
		{0, 0, 3},
		{16, 11, 16},
	} {
		if got := fieldWordLeft(r, c.pos); got != c.left {
			t.Errorf("back from %d stops at %d, want %d", c.pos, got, c.left)
		}
		if got := fieldWordRight(r, c.pos); got != c.right {
			t.Errorf("forward from %d stops at %d, want %d", c.pos, got, c.right)
		}
	}
}

// Ctrl takes a word on Windows and Linux and Option on macOS; AltGr (Ctrl
// with Alt) and plain keys do not.
func TestTheWordDeleteKeyFollowsThePlatform(t *testing.T) {
	saved := optionDeletesWord
	defer func() { optionDeletesWord = saved }()
	optionDeletesWord = false
	if !deletesWord(mod.Control) || !deletesWord(mod.Control|mod.Shift) || deletesWord(mod.Control|mod.Option) ||
		deletesWord(mod.Option) || deletesWord(0) {
		t.Error("on Windows and Linux only Ctrl should take a word")
	}
	optionDeletesWord = true
	if !deletesWord(mod.Option) || deletesWord(mod.OSMenuCommand()) || deletesWord(0) {
		t.Error("on macOS only Option should take a word")
	}
}
