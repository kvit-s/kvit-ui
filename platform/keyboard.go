package platform

// FullKeyboardAccess reports whether the system lets Tab reach every
// control, rather than only text fields and lists. It is the Mac's Keyboard
// navigation setting; every other system lets Tab reach every control.
func FullKeyboardAccess() bool { return fullKeyboardAccess() }

// UnderlineAccessKeys reports whether the system asks for access keys to be
// underlined at all times, rather than only while Alt is held. It is
// Windows' "Underline access keys" setting (Settings, Accessibility,
// Keyboard); every other system answers false.
func UnderlineAccessKeys() bool { return underlineAccessKeys() }
