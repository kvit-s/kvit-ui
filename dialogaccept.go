package kvitui

// Accept answers the dialog as its confirming button does, closing it,
// unless ConfirmKeepsOpen, and running OnAccept: what Return in a dialog's
// text field does, where the field holds the keyboard rather than the
// button.
func (d *Dialog) Accept() { d.accept() }

// Close takes the dialog away without answering it, so neither OnAccept nor
// OnReject runs: for an application that has taken the answer another way,
// as one closes a dialog that stays open on confirming once the answer it
// checked is taken.
func (d *Dialog) Close() { d.finish() }
