//go:build !linux && !windows

package platform

// windowsClipboardPicture is no picture: only a Linux program can be running
// under WSL, and on macOS unison's clipboard is the system's own.
func windowsClipboardPicture() ([]byte, bool) { return nil, false }
