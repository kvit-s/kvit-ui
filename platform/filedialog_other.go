//go:build !linux && !darwin && !windows

package platform

// openNative is unison's own dialog on a system without a native one here.
func (d FileDialog) openNative() []string { return d.openWithUnison() }
