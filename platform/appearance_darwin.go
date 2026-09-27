package platform

// readPlatform reads the defaults NSWorkspace's accessibility preferences are
// backed by. A change made while the application runs is picked up by the
// next Refresh.
func readPlatform() (dark, high, reduced, available bool) {
	dark = commandOutput("defaults", "read", "-g", "AppleInterfaceStyle") == "Dark"
	high = commandOutput("defaults", "read", "com.apple.universalaccess", "increaseContrast") == "1"
	reduced = commandOutput("defaults", "read", "com.apple.Accessibility", "ReduceMotionEnabled") == "1"
	return dark, high, reduced, true
}
