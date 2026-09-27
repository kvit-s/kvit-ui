package kvitui

// Accept answers the dialog as its confirming button does, closing it and
// running OnAccept: what Return in a dialog's text field does, where the
// field holds the keyboard rather than the button.
func (d *Dialog) Accept() { d.accept() }
