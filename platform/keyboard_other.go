//go:build !darwin && !windows

package platform

func fullKeyboardAccess() bool { return true }

func underlineAccessKeys() bool { return false }
