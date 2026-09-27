package platform

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var procGetUserDefaultLocaleName = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetUserDefaultLocaleName")

// localeNameMaxLength is LOCALE_NAME_MAX_LENGTH, which counts the
// terminating zero.
const localeNameMaxLength = 85

// readLocale asks Windows for the user's region, which is already a BCP 47
// tag such as "en-GB".
func readLocale() string {
	buf := make([]uint16, localeNameMaxLength)
	if r, _, _ := procGetUserDefaultLocaleName.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf))); r == 0 {
		return ""
	}
	return windows.UTF16ToString(buf)
}
