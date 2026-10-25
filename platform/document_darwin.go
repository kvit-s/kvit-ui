package platform

import "github.com/ebitengine/purego/objc"

// markEdited calls setDocumentEdited: on each of NSApp's windows with the
// title. AppKit is loaded by unison before any window exists.
func markEdited(title string, edited bool) bool {
	found := false
	withPool(func() {
		app := class("NSApplication")
		if app == 0 {
			return
		}
		windows := app.Send(sel("sharedApplication")).Send(sel("windows"))
		n := objc.Send[uint64](windows, sel("count"))
		for i := range n {
			win := windows.Send(sel("objectAtIndex:"), i)
			if goString(win.Send(sel("title"))) == title {
				win.Send(sel("setDocumentEdited:"), edited)
				found = true
			}
		}
	})
	return found
}
