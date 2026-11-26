package kvitui

import (
	"time"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// caretRestAfter is how long a field's caret blinks with nothing typed, no
// key or press and no caret move, before it stops blinking and stays shown.
// It is GTK's own (gtk-cursor-blink-timeout) and the time kvit-notes' editor
// waits, so every caret in a window stops alike. A variable, so a test can
// shorten it.
var caretRestAfter = 10 * time.Second

// caretHold is how long unison is told to keep a resting caret shown: longer
// than any window stays open.
const caretHold = 100 * 365 * 24 * time.Hour

// caretRest stops a field's caret blinking once the field has been left
// alone, since unison redraws the whole window for every blink.
//
// unison's field blinks its caret every BlinkRate for as long as it holds the
// keyboard. A blink flips the caret only once a moment the field keeps to
// itself has passed, and every change of the selection sets that moment one
// BlinkRate ahead and shows the caret. The field offers no other way to set
// either, so the caret is made to rest by moving it away and straight back
// with BlinkRate set to caretHold for the two moves (holdCaret): the caret is
// shown, and later blinks find the moment still ahead and neither flip it nor
// redraw. The same two moves with the ordinary BlinkRate let it blink again.
//
// unison's blink timer goes on running every BlinkRate while the caret rests,
// each time finding nothing to do. Stopping it would take a BlinkRate as long
// as caretHold for the blink that follows, and the field starts no new timer
// while one is waiting, so the caret could not blink again when asked to.
type caretRest struct {
	// last is when a key, a press or a caret move last reached the field.
	last time.Time
	// resting is set while the caret is held shown.
	resting bool
	// queued is set while a check of how long the field has been left alone
	// is waiting to run.
	queued bool
	// seen is the focus and selection the field was last drawn with.
	seen caretSeen
}

// caretSeen is what a field was drawn with that says the caret moved or
// started to blink: live is set while the field is enabled and holds the
// keyboard of the active window, which is when unison blinks its caret.
type caretSeen struct {
	live       bool
	start, end int
}

// caretInput is a key or a press reaching the field: the caret blinks again
// if it was resting, and stops again once the field is left alone.
func (f *Field) caretInput() {
	f.rest.last = time.Now()
	if f.rest.resting {
		f.rest.resting = false
		holdCaret(f.edit, f.edit.BlinkRate)
	}
	f.queueCaretCheck(caretRestAfter)
}

// caretDrawn notices, as the field is drawn, a caret move or the caret
// starting to blink that no key or press of the field's own brought: a move
// made by the application or a screen reader, the keyboard arriving by Tab,
// the window becoming the active one again, the field being enabled. Every
// such change redraws the field, so drawing is where it is seen.
func (f *Field) caretDrawn() {
	e := f.edit
	start, end := e.Selection()
	now := caretSeen{live: e.Focused() && e.Enabled(), start: start, end: end}
	if now == f.rest.seen {
		return
	}
	f.rest.seen = now
	if !now.live {
		return
	}
	f.rest.last = time.Now()
	if f.rest.resting {
		f.rest.resting = false
		// Not while drawing: the moves mark the field for redrawing.
		unison.InvokeTask(func() {
			if !f.rest.resting {
				holdCaret(e, e.BlinkRate)
			}
		})
	}
	f.queueCaretCheck(caretRestAfter)
}

// queueCaretCheck checks after a while whether the field has been left alone
// long enough, unless a check is already waiting.
func (f *Field) queueCaretCheck(after time.Duration) {
	if f.rest.queued {
		return
	}
	f.rest.queued = true
	unison.InvokeTaskAfter(f.checkCaret, after)
}

// checkCaret makes the caret rest once the field holding the keyboard has
// been left alone for caretRestAfter, and checks again later when it has not
// been. A field without the keyboard is checked again when it next has it.
func (f *Field) checkCaret() {
	f.rest.queued = false
	e := f.edit
	if f.rest.resting || !e.Focused() || !e.Enabled() {
		return
	}
	if left := caretRestAfter - time.Since(f.rest.last); left > 0 {
		f.queueCaretCheck(left)
		return
	}
	f.rest.resting = true
	holdCaret(e, caretHold)
}

// holdCaret shows a field's caret and keeps it shown for hold, through the
// one route unison's field offers: a change of its selection. The selection
// is moved away and straight back, the text and selection ending as they
// were, with the field's scrolling and that of the panels around it put back
// so a view the person has scrolled away from the caret stays where it is.
// An empty field has no other place for the caret, so it is given a space
// and has it taken away again; neither is reported as a change.
func holdCaret(e *unison.Field, hold time.Duration) {
	rate, offset, scroll := e.BlinkRate, e.ScrollOffset(), e.ScrollRectIntoViewCallback
	e.BlinkRate = hold
	e.ScrollRectIntoViewCallback = func(geom.Rect) bool { return true }
	was := e.GetFieldState()
	away := *was
	switch {
	case was.Text == "":
		away.Text, away.SelectionStart, away.SelectionEnd, away.SelectionAnchor = " ", 1, 1, 1
	case was.SelectionStart != was.SelectionEnd:
		away.SelectionEnd, away.SelectionAnchor = was.SelectionStart, was.SelectionStart
	case was.SelectionStart > 0:
		away.SelectionStart, away.SelectionEnd, away.SelectionAnchor = was.SelectionStart-1, was.SelectionStart-1, was.SelectionStart-1
	default:
		away.SelectionStart, away.SelectionEnd, away.SelectionAnchor = 1, 1, 1
	}
	e.ApplyFieldState(&away)
	e.ApplyFieldState(was)
	e.ScrollRectIntoViewCallback = scroll
	e.SetScrollOffset(offset)
	e.BlinkRate = rate
}
