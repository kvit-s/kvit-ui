package kvitui

import "github.com/richardwilkes/unison"

// WindowOf is the Kvit window a panel is shown in, or nil when it is in no
// window or in one that is not a Kvit window. A component that shows a popup
// on behalf of a panel it did not build, such as an editor's / menu, finds
// the window's popup layer with it.
func (u *UI) WindowOf(p unison.Paneler) *Window { return u.windowOf(p) }
