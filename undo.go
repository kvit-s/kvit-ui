package kvitui

// UndoStack is the app's edits that can be undone, each named for what it
// did (FDB-7, WRD-11): the Edit menu's lines read "Undo delete transaction"
// and "Redo delete transaction", and a toast's or a ConfirmInPlace's undo
// takes back the edit it reports. Every edit since the document was opened
// can be undone; Clear forgets them when another document opens.
type UndoStack struct {
	ui           *UI
	done, undone []undoEntry
	undo, redo   *Command
}

// undoEntry is one edit: its name and how to take it back and do it again.
type undoEntry struct {
	name       string
	undo, redo func()
}

// NewUndoStack returns an empty stack.
func NewUndoStack(ui *UI) *UndoStack {
	s := &UndoStack{ui: ui}
	s.undo = &Command{Std: StdUndo, Enabled: s.CanUndo, Run: s.Undo}
	s.redo = &Command{Std: StdRedo, Enabled: s.CanRedo, Run: s.Redo}
	s.rename()
	return s
}

// Add records an edit the reader just made, named in sentence case for what
// it did ("Delete transaction"), with how to take it back and do it again.
// It forgets the edits that were undone, which cannot be redone after a new
// one.
func (s *UndoStack) Add(name string, undo, redo func()) {
	s.done = append(s.done, undoEntry{name: name, undo: undo, redo: redo})
	s.undone = nil
	s.rename()
}

// CanUndo reports whether there is an edit to take back.
func (s *UndoStack) CanUndo() bool { return len(s.done) > 0 }

// CanRedo reports whether there is an edit taken back to do again.
func (s *UndoStack) CanRedo() bool { return len(s.undone) > 0 }

// UndoName is the name of the edit Undo takes back, "" for none.
func (s *UndoStack) UndoName() string {
	if len(s.done) == 0 {
		return ""
	}
	return s.done[len(s.done)-1].name
}

// RedoName is the name of the edit Redo does again, "" for none.
func (s *UndoStack) RedoName() string {
	if len(s.undone) == 0 {
		return ""
	}
	return s.undone[len(s.undone)-1].name
}

// Undo takes back the last edit.
func (s *UndoStack) Undo() {
	if len(s.done) == 0 {
		return
	}
	e := s.done[len(s.done)-1]
	s.done = s.done[:len(s.done)-1]
	if e.undo != nil {
		e.undo()
	}
	s.undone = append(s.undone, e)
	s.rename()
	s.ui.Refresh()
}

// Redo does again the last edit taken back.
func (s *UndoStack) Redo() {
	if len(s.undone) == 0 {
		return
	}
	e := s.undone[len(s.undone)-1]
	s.undone = s.undone[:len(s.undone)-1]
	if e.redo != nil {
		e.redo()
	}
	s.done = append(s.done, e)
	s.rename()
	s.ui.Refresh()
}

// Clear forgets every edit, as opening another document does.
func (s *UndoStack) Clear() {
	s.done, s.undone = nil, nil
	s.rename()
	s.ui.Refresh()
}

// Commands are the standard Undo and Redo commands for the stack, for the
// Edit menu: named for the edit they reverse, and done by a focused text
// field itself while it has an edit of its own to take back.
func (s *UndoStack) Commands() (undo, redo *Command) { return s.undo, s.redo }

// rename names the commands for the edits they reverse.
func (s *UndoStack) rename() {
	s.undo.Name, s.redo.Name = "Undo", "Redo"
	if n := s.UndoName(); n != "" {
		s.undo.Name = "Undo " + lowerFirst(n)
	}
	if n := s.RedoName(); n != "" {
		s.redo.Name = "Redo " + lowerFirst(n)
	}
}

// Toast returns a toast saying what the last edit did, with an Undo that
// takes it back: the confirmation of a reversible action, in place of a
// question before it.
func (s *UndoStack) Toast(words string) *Toast {
	t := NewToast(s.ui, words)
	t.Tone = ToneSuccess
	if s.CanUndo() {
		t.Action = "Undo"
		t.OnAction = s.Undo
	}
	return t
}

// Confirm returns a ConfirmInPlace saying what the last edit did, with an
// undo that takes it back.
func (s *UndoStack) Confirm(words string) *ConfirmInPlace {
	c := NewConfirmInPlace(s.ui, words)
	c.OnUndo = s.Undo
	return c
}
