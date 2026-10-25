package platform

import "strconv"

// fullKeyboardAccess reads AppleKeyboardUIMode, whose second bit is the
// Keyboard navigation setting in System Settings.
func fullKeyboardAccess() bool {
	v, err := strconv.Atoi(commandOutput("defaults", "read", "-g", "AppleKeyboardUIMode"))
	return err == nil && v&2 != 0
}

// underlineAccessKeys is false: the Mac has no access keys.
func underlineAccessKeys() bool { return false }
