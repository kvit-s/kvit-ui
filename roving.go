package kvitui

import (
	"github.com/richardwilkes/unison"
)

// A group of controls that acts as one, a row of tabs, a segmented control,
// a radio group, the sidebar, is one stop in the Tab order, and the arrow
// keys move inside it (KEY-9): Tab lands on the chosen member, or the first
// when none is chosen, and leaves the group again.

// siblings are the members of a group: the controls of one kind among a
// panel's parent's children, in order.
func siblings[T unison.Paneler](p *unison.Panel) []T {
	var out []T
	parent := p.Parent()
	if parent == nil {
		return out
	}
	for _, c := range parent.Children() {
		if m, ok := c.Self.(T); ok && !c.Hidden {
			out = append(out, m)
		}
	}
	return out
}

// roveStep is the member a key moves to from one: the arrow keys along the
// group's way to the one before or after, Home and End to the ends, round
// from the last to the first where wrap is set. A disabled member is passed
// by. It reports false for a key that does not move.
func roveStep[T unison.Paneler](members []T, from T, key unison.KeyCode, vertical, wrap bool) (T, bool) {
	var none T
	at := -1
	for i, m := range members {
		if m.AsPanel() == from.AsPanel() {
			at = i
		}
	}
	n := len(members)
	// A member alone is no group: its keys go on to what holds it.
	if at < 0 || n < 2 {
		return none, false
	}
	back, forward := unison.KeyLeft, unison.KeyRight
	if vertical {
		back, forward = unison.KeyUp, unison.KeyDown
	}
	by := 0
	switch key {
	case back:
		by = -1
	case forward:
		by = 1
	case unison.KeyHome:
		at, by = -1, 1
	case unison.KeyEnd:
		at, by = n, -1
	default:
		return none, false
	}
	for range n {
		next := at + by
		if next < 0 || next >= n {
			if !wrap || key == unison.KeyHome || key == unison.KeyEnd {
				return none, true
			}
			next = (next + n) % n
		}
		at = next
		if members[at].AsPanel().Enabled() {
			return members[at], true
		}
	}
	return none, true
}

// oneStop makes the chosen member of a group, or the first when none is
// chosen, the group's one stop in the Tab order. A member holding the
// keyboard, which the arrows moved to, keeps it.
func oneStop[T unison.Paneler](members []T, chosen func(T) bool) {
	stop := 0
	for i, m := range members {
		if chosen(m) {
			stop = i
			break
		}
	}
	for i, m := range members {
		m.AsPanel().SetFocusable(i == stop || m.AsPanel().Focused())
	}
}
