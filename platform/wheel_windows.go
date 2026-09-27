package platform

import "unsafe"

const spiGetWheelScrollLines = 0x0068

// readWheelScrollLines asks Windows, whose mouse settings hold the number.
// A setting of "one screen at a time" reads as a very large number, which is
// passed on as it is.
func readWheelScrollLines() int {
	var lines uint32
	if r, _, _ := procSystemParametersInfo.Call(spiGetWheelScrollLines, 0, uintptr(unsafe.Pointer(&lines)), 0); r == 0 {
		return 0
	}
	return int(lines)
}
