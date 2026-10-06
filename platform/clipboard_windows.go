package platform

import (
	"bytes"
	"time"
	"unsafe"

	"github.com/richardwilkes/toolbox/v2/xruntime"
	"github.com/richardwilkes/unison"
	"golang.org/x/sys/windows"
)

var (
	procOpenClipboard              = user32.NewProc("OpenClipboard")
	procCloseClipboard             = user32.NewProc("CloseClipboard")
	procIsClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable")
	procGetClipboardData           = user32.NewProc("GetClipboardData")
	procRegisterClipboardFormat    = user32.NewProc("RegisterClipboardFormatW")
	procGlobalLock                 = kernel32.NewProc("GlobalLock")
	procGlobalUnlock               = kernel32.NewProc("GlobalUnlock")
	procGlobalSize                 = kernel32.NewProc("GlobalSize")
)

const (
	cfDIB   = 8
	cfDIBV5 = 17
)

// pngFormatNames are the names under which programs register a PNG on the
// clipboard: "PNG" by browsers, the Snipping Tool and Office, "image/png"
// by others.
var pngFormatNames = []string{"PNG", "image/png"}

// windowsClipboardPicture reads the clipboard's PNG where a program put one,
// since it keeps the picture's transparency, and its bitmap otherwise. It
// skips a headless session, whose clipboard is its own in memory, so a test
// never reads the desktop's clipboard.
func windowsClipboardPicture() ([]byte, bool) {
	if unison.ActiveHeadlessScreen() != nil || !openClipboard() {
		return nil, false
	}
	defer procCloseClipboard.Call()
	for _, name := range pngFormatNames {
		if data := clipboardBytes(registeredFormat(name)); bytes.HasPrefix(data, pngSignature) {
			return data, true
		}
	}
	for _, format := range []uintptr{cfDIBV5, cfDIB} {
		if picture, ok := dibPNG(clipboardBytes(format)); ok {
			return picture, true
		}
	}
	return nil, false
}

// openClipboard opens the clipboard for reading, trying again for a moment
// while another program has it open.
func openClipboard() bool {
	for range 10 {
		if r, _, _ := procOpenClipboard.Call(0); r != 0 {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return false
}

// registeredFormat is the number of the clipboard format registered under
// name, which RegisterClipboardFormatW answers whether or not a program
// registered it first; 0 when it cannot.
func registeredFormat(name string) uintptr {
	text, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0
	}
	r, _, _ := procRegisterClipboardFormat.Call(uintptr(unsafe.Pointer(text)))
	return r
}

// clipboardBytes copies what the open clipboard holds in format; nil when
// it holds nothing in it.
func clipboardBytes(format uintptr) []byte {
	if format == 0 {
		return nil
	}
	if r, _, _ := procIsClipboardFormatAvailable.Call(format); r == 0 {
		return nil
	}
	handle, _, _ := procGetClipboardData.Call(format)
	if handle == 0 {
		return nil
	}
	locked, _, _ := procGlobalLock.Call(handle)
	if locked == 0 {
		return nil
	}
	defer procGlobalUnlock.Call(handle)
	size, _, _ := procGlobalSize.Call(handle)
	if size == 0 {
		return nil
	}
	return bytes.Clone(unsafe.Slice(xruntime.PtrFromUintptr[byte](locked), size))
}
