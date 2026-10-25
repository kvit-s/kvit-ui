package platform

// MarkEdited shows or clears the system's own mark for a window whose
// document has unsaved changes, where the system has one: the dot in the
// close button on macOS. The window is found by its title, as unison keeps
// the system's window to itself, so every window with the title is marked.
// It reports whether one was found; on other systems it does nothing, and
// the title has the mark instead.
func MarkEdited(title string, edited bool) bool { return markEdited(title, edited) }
