//go:build !linux && !darwin && !windows

package platform

// newTrayBackend has no notification area to offer on this system.
func newTrayBackend(trayEvents) trayBackend { return noTray{} }
