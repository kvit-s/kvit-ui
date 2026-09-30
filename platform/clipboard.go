package platform

// WindowsClipboardPicture answers the picture on Windows' clipboard as PNG
// bytes when the program runs under Windows Subsystem for Linux (WSL), and
// false anywhere else or when Windows' clipboard holds no picture.
//
// Under WSL the program's window is an X11 window, and WSLg's bridge between
// Windows' clipboard and X11 offers X11 programs text only, so a picture
// copied in Windows never appears on unison's clipboard. PowerShell is asked
// for it instead, as a Windows program WSL runs. That takes about 0.65 s and
// the caller waits for it, so ask only after unison's clipboard has shown no
// picture.
func WindowsClipboardPicture() ([]byte, bool) { return windowsClipboardPicture() }
