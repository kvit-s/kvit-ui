//go:build !linux

package platform

// windowsClipboardPicture is no picture: only a Linux program can be running
// under WSL, and elsewhere unison's clipboard is the system's own.
func windowsClipboardPicture() ([]byte, bool) { return nil, false }
