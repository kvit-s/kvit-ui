package kvitui

import "time"

// SetCaretRestAfter shortens, for a test, how long a field's caret blinks
// with the field left alone, and returns what puts it back. A test registers
// the return with t.Cleanup before it starts a headless screen, so it runs
// once the screen has stopped.
func SetCaretRestAfter(d time.Duration) (restore func()) {
	saved := caretRestAfter
	caretRestAfter = d
	return func() { caretRestAfter = saved }
}
