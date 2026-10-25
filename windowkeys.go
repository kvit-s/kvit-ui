package kvitui

import (
	"slices"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// windowKeys are the keys a window answers itself when no command has
// them: Ctrl+Tab and Ctrl+Shift+Tab choose the next or previous tab of a row
// of tabs on every system, and Ctrl+Page Down and Ctrl+Page Up do on the
// Windows form (KEY-9, decided Q-18); Find's keys put the keyboard in the
// window's search field with its text selected (KEY-13, decided Q-20).
func (w *Window) windowKeys(key unison.KeyCode, mods mod.Modifiers) bool {
	if w.popupOpen() {
		return false
	}
	held := mods & mod.NonSticky
	switch {
	case key == unison.KeyTab && held == mod.Control:
		return w.stepTab(1)
	case key == unison.KeyTab && held == mod.Control|mod.Shift:
		return w.stepTab(-1)
	case w.ui.form == FormWindows && held == mod.Control && key == unison.KeyPageDown:
		return w.stepTab(1)
	case w.ui.form == FormWindows && held == mod.Control && key == unison.KeyPageUp:
		return w.stepTab(-1)
	}
	for _, s := range StdKeys(StdFind, w.ui.form) {
		if s.Matches(w.ui.form, key, mods) {
			if f := w.searchField(); f != nil {
				f.Focus()
				f.edit.SelectAll()
				return true
			}
		}
	}
	return false
}

// limitedTab moves Tab and Shift+Tab among the text fields and lists only,
// on the macOS form with full keyboard access off (KEY-10).
func (w *Window) limitedTab(key unison.KeyCode, mods mod.Modifiers) bool {
	if key != unison.KeyTab || mods&mod.NonSticky&^mod.Shift != 0 || w.ui.form != FormMac || w.ui.FullKeyboardAccess {
		return false
	}
	var stops []*unison.Panel
	var walk func(p *unison.Panel)
	walk = func(p *unison.Panel) {
		if p.Hidden || !p.Enabled() {
			return
		}
		if p.Focusable() && tabStopWithoutFullAccess(p) {
			stops = append(stops, p)
		}
		for _, c := range p.Children() {
			walk(c)
		}
	}
	scope, _ := w.accessScope()
	walk(scope)
	if len(stops) == 0 {
		return true
	}
	at := slices.Index(stops, w.CurrentFocus())
	next := 0
	switch {
	case at < 0 && mods.ShiftDown():
		next = len(stops) - 1
	case at >= 0 && mods.ShiftDown():
		next = (at - 1 + len(stops)) % len(stops)
	case at >= 0:
		next = (at + 1) % len(stops)
	}
	stops[next].RequestFocus()
	return true
}

// tabStopWithoutFullAccess reports whether a panel is a text field or a
// list, which Tab reaches on a Mac without full keyboard access.
func tabStopWithoutFullAccess(p *unison.Panel) bool {
	switch p.Self.(type) {
	case *unison.Field, *tableBody, *Tree, *dualSide:
		return true
	}
	return false
}

// stepTab chooses the tab after or before the chosen one, going round at the
// ends, in the row of tabs the keyboard is in, or else the first row of tabs
// in the window.
func (w *Window) stepTab(by int) bool {
	var row []*Tab
	if f := w.CurrentFocus(); f != nil {
		for p := f; p != nil && row == nil; p = p.Parent() {
			if t, ok := p.Self.(*Tab); ok {
				row = siblings[*Tab](t.AsPanel())
			}
		}
	}
	if row == nil {
		var walk func(p *unison.Panel) bool
		walk = func(p *unison.Panel) bool {
			if p.Hidden {
				return false
			}
			if t, ok := p.Self.(*Tab); ok {
				row = siblings[*Tab](t.AsPanel())
				return true
			}
			return slices.ContainsFunc(p.Children(), walk)
		}
		walk(w.stage)
		if row == nil {
			walk(w.Content())
		}
	}
	if len(row) < 2 {
		return false
	}
	at := slices.IndexFunc(row, func(t *Tab) bool { return t.Selected })
	n := len(row)
	for range n {
		at = (at + by + n) % n
		if row[at].Enabled() {
			row[at].fire()
			if f := w.CurrentFocus(); f != nil {
				if _, ok := f.Self.(*Tab); ok {
					row[at].Focus()
				}
			}
			return true
		}
	}
	return false
}

// searchField is the window's search field, the first shown in it, or nil.
func (w *Window) searchField() *SearchField {
	var found *SearchField
	var walk func(p *unison.Panel) bool
	walk = func(p *unison.Panel) bool {
		if p.Hidden || !p.Enabled() {
			return false
		}
		if s, ok := p.Self.(*SearchField); ok {
			found = s
			return true
		}
		return slices.ContainsFunc(p.Children(), walk)
	}
	walk(w.Content())
	return found
}

// startDrag says a drag that Escape abandons is under way, with what puts
// things back as they were before it started (KEY-7).
func (w *Window) startDrag(undo func()) { w.dragUndo = undo }

// endDrag says the drag is over.
func (w *Window) endDrag() { w.dragUndo = nil }

// dragKeys abandons a drag under way on Escape.
func (w *Window) dragKeys(key unison.KeyCode) bool {
	if key != unison.KeyEscape || w.dragUndo == nil {
		return false
	}
	undo := w.dragUndo
	w.dragUndo = nil
	undo()
	return true
}
