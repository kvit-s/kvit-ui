package uitest

import (
	"github.com/richardwilkes/unison/accessibility"
)

// Responses lists how a window changed between two readings of what a screen
// reader is told about it, leaving out what the act of pressing or typing
// changes by itself: where the keyboard focus is. A node added or removed, a
// name, value, description or state that changed, text typed or deleted,
// and a part that moved all count. An empty list means nothing a reader can
// see or hear changed.
func Responses(before, after *accessibility.Tree) []accessibility.Event {
	if before == nil || after == nil {
		return nil
	}
	var out []accessibility.Event
	for _, e := range accessibility.Diff(before, after) {
		switch e.Kind {
		case accessibility.FocusChanged, accessibility.TextSelectionChanged,
			accessibility.WindowActivated, accessibility.WindowDeactivated:
			continue
		}
		out = append(out, e)
	}
	return out
}

// CheckResponds does act, which should be what a reader does (press a
// button, choose a menu line, press a key), lets the window settle, and fails
// the test when nothing a reader can see or hear changed: every action shows
// its result where the reader is looking (agents-ux-mining/rubric.md, R16).
// Copy that changes nothing on screen, a Send that leaves the message where it
// was, and a Run again that keeps offering itself fail it. What is compared is
// what a screen reader is told about the window, as Responses reads it,
// rather than its picture, since pressing a control also moves the pointer
// onto it and changes how it is drawn. what names the action in the failure,
// as in "pressing Copy".
func (s *Session) CheckResponds(what string, act func()) {
	s.t.Helper()
	before := s.Tree()
	act()
	s.Sync()
	if len(Responses(before, s.Tree())) == 0 {
		s.t.Errorf("%s changed nothing a reader can see or hear; an action shows its result where the reader is looking", what)
	}
}
