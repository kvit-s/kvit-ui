package platform

import "unsafe"

// spiGetKeyboardCues is SPI_GETKEYBOARDCUES, which answers true when
// Windows underlines access keys without waiting for Alt.
const spiGetKeyboardCues = 0x100A

func fullKeyboardAccess() bool { return true }

func underlineAccessKeys() bool {
	var on int32
	r, _, _ := procSystemParametersInfo.Call(spiGetKeyboardCues, 0, uintptr(unsafe.Pointer(&on)), 0)
	return r != 0 && on != 0
}
