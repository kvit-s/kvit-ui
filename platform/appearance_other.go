//go:build !linux && !darwin && !windows

package platform

// readPlatform has nothing to read on this platform.
func readPlatform() (dark, high, reduced, available bool) { return false, false, false, false }
