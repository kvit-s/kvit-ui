//go:build !windows

package platform

// readWheelScrollLines has no desktop setting to read here.
func readWheelScrollLines() int { return 0 }
