package kvitui

import (
	"github.com/kvit-s/kvit-ui/platform"
	"github.com/richardwilkes/unison"
)

// Notify posts a notification through the tray, unless one of the app's
// windows is in front: the reader is looking at the app, and a notification
// would only interrupt what they are doing (MSG-10). It reports whether it
// posted. A notification is for something that happened in the background,
// never the result of what the reader just did.
func (u *UI) Notify(t *platform.Tray, n platform.Notification) bool {
	if t == nil {
		return false
	}
	if w := unison.ActiveWindow(); w != nil && u.windows[w] != nil {
		return false
	}
	return t.Post(n)
}
