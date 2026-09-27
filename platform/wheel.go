package platform

// WheelScrollLines is how many lines one notch of a mouse wheel scrolls, as
// the desktop is set: 3 where it gives no answer, which is every desktop's
// default.
func WheelScrollLines() int {
	if n := readWheelScrollLines(); n > 0 {
		return n
	}
	return 3
}
