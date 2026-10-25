package kvitui

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// AddCommands makes the window answer the commands' keys: a key runs its
// command wherever the window has the focus, and does nothing, silently,
// while the command cannot act (KEY-2). A command is refused, with a
// warning, when one of its keys in either form is already another
// command's, a standard command's that it is not (KEY-5), or one the system
// or the access keys own (KEY-3, KEY-4); the error lists every refusal. A
// refused command is still added for its other keys, and still runs from
// its menu line and buttons.
//
// A key a text field uses, typing and moving the caret, goes to a focused
// text field rather than to a command. A standard command a focused control
// can carry out itself, such as Copy in a text field or Select all in a
// list, is done by the control.
func (w *Window) AddCommands(cmds ...*Command) error {
	var refused []string
	for _, c := range cmds {
		if c == nil || w.hasCommand(c) {
			continue
		}
		for _, f := range []Form{FormWindows, FormMac} {
			for _, s := range c.Shortcuts(f) {
				if why := w.refusal(c, f, s); why != "" {
					refused = append(refused, why)
				}
			}
		}
		w.commands = append(w.commands, c)
	}
	if len(refused) == 0 {
		return nil
	}
	for _, r := range refused {
		slog.Warn("kvitui: a shortcut was refused", "why", r)
	}
	return fmt.Errorf("kvitui: shortcuts refused: %s", strings.Join(refused, "; "))
}

// Commands are the commands the window answers the keys of, in the order
// they were added.
func (w *Window) Commands() []*Command { return append([]*Command(nil), w.commands...) }

func (w *Window) hasCommand(c *Command) bool {
	for _, have := range w.commands {
		if have == c {
			return true
		}
	}
	return false
}

// refusal says why a command may not have a key in a form, or "".
func (w *Window) refusal(c *Command, f Form, s Shortcut) string {
	name := c.Name
	if name == "" {
		name = stdName(c.Std, f)
	}
	say := func(why string) string {
		return fmt.Sprintf("%s on %s (%s form): %s", name, shortcutText(f, s), f, why)
	}
	if why := systemKey(f, s); why != "" {
		return say(why)
	}
	for _, other := range w.commands {
		for _, o := range other.Shortcuts(f) {
			if sameShortcut(o, s, f) {
				oname := other.Name
				if oname == "" {
					oname = stdName(other.Std, f)
				}
				return say("already " + oname + "'s")
			}
		}
	}
	for std := StdNone + 1; std < stdCount; std++ {
		if std == c.Std || !reservedStd(std) {
			continue
		}
		for _, o := range StdKeys(std, f) {
			if sameShortcut(o, s, f) {
				return say("the standard key for " + stdName(std, f))
			}
		}
	}
	return ""
}

// reservedStd reports whether a standard command's keys mean it in every
// app, so no other command may take them; the tab commands' keys are free
// in an app without document tabs (KEY-5), and the zoom keys in one that
// does not zoom.
func reservedStd(s StdCommand) bool {
	switch s {
	case StdNewTab, StdCloseTab, StdReopenTab, StdNextTab, StdPreviousTab, StdZoomIn, StdZoomOut, StdActualSize:
		return false
	}
	return true
}

// sameShortcut reports whether two shortcuts are the same keys in a form.
func sameShortcut(a, b Shortcut, f Form) bool {
	return sameKey(a.Key) == sameKey(b.Key) && a.modifiers(f) == b.modifiers(f)
}

// systemKey says why a key belongs to the system, the access keys or focus
// in a form, or "" when an app command may have it (KEY-3, KEY-4, KEY-5).
func systemKey(f Form, s Shortcut) string {
	key := sameKey(s.Key)
	letterOrDigit := (key >= unison.KeyA && key <= unison.KeyZ) || (key >= unison.Key0 && key <= unison.Key9)
	m := s.modifiers(f)
	switch {
	case key == unison.KeyTab && m&^mod.Shift == 0:
		return "Tab moves the focus"
	case key == unison.KeyMenu || (key == unison.KeyF10 && m == mod.Shift):
		return "it opens the context menu"
	}
	if f == FormWindows {
		switch {
		case m&(mod.Control|mod.Option) == mod.Control|mod.Option:
			return "Ctrl+Alt is AltGr on many keyboards, which types characters"
		case m&^mod.Shift == mod.Option && letterOrDigit:
			return "Alt with a letter or digit is an access key"
		case m == mod.Option && (key == unison.KeyF4 || key == unison.KeyTab || key == unison.KeyEscape || key == unison.KeySpace):
			return "the system's"
		case m&mod.Control != 0 && key == unison.KeyEscape:
			return "the system's"
		case m == 0 && key == unison.KeyF10:
			return "F10 goes to the menu bar"
		case m == mod.Control && key == unison.KeyQ:
			return "Ctrl+Q is never bound to an app command"
		case m&mod.Command != 0:
			return "the Windows key belongs to the system"
		}
		return ""
	}
	switch {
	case m == mod.Command && (key == unison.KeySpace || key == unison.KeyTab || key == unison.KeyH || key == unison.KeyM || key == unison.KeyBackQuote):
		return "the system's"
	case m == mod.Control && key == unison.KeySpace:
		return "it switches input sources"
	case m == mod.Shift|mod.Command && (key == unison.Key3 || key == unison.Key4 || key == unison.Key5):
		return "it takes a screenshot"
	case m == mod.Option|mod.Command && (key == unison.KeyEscape || key == unison.KeyD || key == unison.KeyH):
		return "the system's"
	case m == mod.Control|mod.Command && key == unison.KeyQ:
		return "it locks the screen"
	case m == mod.Control && key >= unison.KeyF1 && key <= unison.KeyF7:
		return "it belongs to keyboard access"
	case m == mod.Control && (key == unison.KeyLeft || key == unison.KeyRight || key == unison.KeyUp || key == unison.KeyDown):
		return "it switches spaces"
	}
	return ""
}

// commandFor is the command whose key a press is in the app's form, or nil:
// one added to the window first, then one on a toolbar shown in it.
func (w *Window) commandFor(key unison.KeyCode, mods mod.Modifiers) *Command {
	f := w.ui.form
	find := func(cmds []*Command) *Command {
		for _, c := range cmds {
			for _, s := range c.Shortcuts(f) {
				if s.Matches(f, key, mods) {
					return c
				}
			}
		}
		return nil
	}
	if c := find(w.commands); c != nil {
		return c
	}
	for _, t := range w.toolbars() {
		if c := find(t.Commands()); c != nil {
			return c
		}
	}
	return nil
}

// toolbars are the toolbars shown in the window.
func (w *Window) toolbars() []*Toolbar {
	var out []*Toolbar
	var walk func(p *unison.Panel)
	walk = func(p *unison.Panel) {
		if p.Hidden {
			return
		}
		if t, ok := p.Self.(*Toolbar); ok {
			out = append(out, t)
			return
		}
		for _, c := range p.Children() {
			walk(c)
		}
	}
	walk(w.Content())
	return out
}

// commandKeys runs the command a key press is the shortcut of, and reports
// whether the press is used up. While a menu, a dialog or a popover is open
// the keys are its own.
func (w *Window) commandKeys(key unison.KeyCode, mods mod.Modifiers) bool {
	c := w.commandFor(key, mods)
	if c == nil || w.popupOpen() {
		return false
	}
	focus := w.CurrentFocus()
	if focus != nil && isTextField(focus) && typingKey(key, mods) {
		return false
	}
	if runFocused(c, focus, w) {
		return true
	}
	// A command that cannot act uses up its key all the same, silently.
	c.Do()
	return true
}

// runFocused carries out a standard command in the focused control, when
// the control can: cut, copy and paste in a text field, select all in a
// list, undo and redo in a text field.
func runFocused(c *Command, focus *unison.Panel, src any) bool {
	if focus == nil {
		return false
	}
	if id := stdInfoOf(c.Std).focused; id != 0 && focus.CanPerformCmd(src, id) {
		focus.PerformCmd(src, id)
		return true
	}
	if (c.Std == StdUndo || c.Std == StdRedo) && isTextField(focus) {
		if m := unison.UndoManagerFor(focus); m != nil {
			if c.Std == StdUndo && m.CanUndo() {
				m.Undo()
				return true
			}
			if c.Std == StdRedo && m.CanRedo() {
				m.Redo()
				return true
			}
		}
	}
	return false
}

// focusedCan reports whether the focused control would carry out a
// standard command itself, for a menu line to show it as enabled.
func focusedCan(c *Command, focus *unison.Panel, src any) bool {
	if focus == nil {
		return false
	}
	if id := stdInfoOf(c.Std).focused; id != 0 && focus.CanPerformCmd(src, id) {
		return true
	}
	if (c.Std == StdUndo || c.Std == StdRedo) && isTextField(focus) {
		if m := unison.UndoManagerFor(focus); m != nil {
			return (c.Std == StdUndo && m.CanUndo()) || (c.Std == StdRedo && m.CanRedo())
		}
	}
	return false
}

// popupOpen reports whether a popup that takes the keyboard is shown.
func (w *Window) popupOpen() bool {
	for _, p := range w.popups {
		if !p.Passive && !p.Ambient {
			return true
		}
	}
	return false
}

// isTextField reports whether a panel is a text field the reader types in.
func isTextField(p *unison.Panel) bool {
	_, ok := p.Self.(*unison.Field)
	return ok
}

// typingKey reports whether a key press is one a text field uses itself:
// typing, with or without Shift, and moving or deleting with the arrows,
// Home, End, Page Up and Down, Backspace and Delete under any modifier.
// Function keys, Escape and Tab are not.
func typingKey(key unison.KeyCode, mods mod.Modifiers) bool {
	switch key {
	case unison.KeyLeft, unison.KeyRight, unison.KeyUp, unison.KeyDown, unison.KeyHome, unison.KeyEnd,
		unison.KeyPageUp, unison.KeyPageDown, unison.KeyBackspace, unison.KeyDelete:
		return true
	case unison.KeyEscape, unison.KeyTab:
		return false
	}
	if key >= unison.KeyF1 && key <= unison.KeyF25 {
		return false
	}
	return mods&mod.NonSticky&^mod.Shift == 0
}
