package kvitui

import (
	"github.com/richardwilkes/unison"
)

// Command is one thing the app can do, made once and shown wherever it is
// offered: a menu line, a button, a symbol button on a toolbar, a shortcut.
// Every place draws its name, symbol and shortcut from the command, so a
// command has one name and one key everywhere (CON-1), and is disabled
// everywhere at once (CMD-9).
//
// A standard command, one with Std set, takes its name, symbol and keys
// from the platform's conventions unless they are set here: an app that has
// Save gets Ctrl+S on the Windows form and Cmd-S on the macOS form without
// choosing it (KEY-2). A window answers a command's keys once the command is
// added to it (Window.AddCommands).
type Command struct {
	// Name is what the command is called everywhere: a verb first, in
	// sentence case, "Export…" with the ellipsis left to Asks. kvit-ui gives
	// it the form's case. A standard command's name is the standard one
	// unless set.
	Name string
	// Std names a standard command, StdNone for the app's own.
	Std StdCommand
	// Symbol is the meaning name of its symbol, drawn on a symbol button and
	// before the words in a menu on the Windows form; a standard command has
	// its standard one unless set.
	Symbol string
	// Shortcut is the app's own key for the command; a standard command's
	// keys are the form's standard ones, and this adds one more. Give new
	// shortcuts only to frequent commands (KEY-5).
	Shortcut Shortcut
	// Asks marks a command that asks for more before it acts, a dialog, a
	// chooser or a confirmation: its name gains "…" on both forms (WRD-4).
	Asks bool
	// OpensWindow marks a command whose whole job is to show a window, such
	// as Settings: "…" on the macOS form only (WRD-4).
	OpensWindow bool
	// Explanation is one sentence the name cannot say, such as why the
	// command cannot act now: a tooltip's second line and the description a
	// screen reader is told.
	Explanation string
	// Enabled reports whether the command can act now; nil is always. A
	// command that cannot act stays where it is, disabled, and its key does
	// nothing (CMD-9, KEY-2). One too slow to check leaves this nil and
	// says why when it is run and cannot act.
	Enabled func() bool
	// Checked makes the command a switch, reporting whether it is on now: a
	// tick on its menu line and a checked button.
	Checked func() bool
	// TurnOn and TurnOff are a switch's words on the macOS form, where the
	// menu line says what choosing it does: "Show sidebar" while it is off,
	// "Hide sidebar" while it is on. The Windows form keeps Name and ticks
	// it (WRD-12). Empty keeps Name on both forms.
	TurnOn, TurnOff string
	// Run does it.
	Run func()
}

// StdCommand names one of the commands every desktop app is expected to
// have under the same name and key.
type StdCommand int

// The standard commands.
const (
	StdNone StdCommand = iota
	StdNew
	StdOpen
	StdSave
	StdSaveAs
	StdClose
	StdPrint
	StdUndo
	StdRedo
	StdCut
	StdCopy
	StdPaste
	StdDelete
	StdSelectAll
	StdFind
	StdFindNext
	StdReplace
	StdGoTo
	StdRename
	StdRefresh
	StdProperties
	StdFullScreen
	StdZoomIn
	StdZoomOut
	StdActualSize
	StdNewTab
	StdCloseTab
	StdReopenTab
	StdNextTab
	StdPreviousTab
	StdBack
	StdSettings
	StdHelp
	StdAbout
	StdQuit
	stdCount
)

// StdMenu names one of the standard menus of a menu bar.
type StdMenu int

// The standard menus, in the order a menu bar shows those its form has: the
// app menu and Window only on the macOS form, Tools only on the Windows
// form.
const (
	MenuNone StdMenu = iota
	MenuApp
	MenuFile
	MenuEdit
	MenuFormat
	MenuView
	MenuTools
	MenuWindow
	MenuHelp
)

// stdInfo is what a standard command is called, where it goes and its keys,
// on each form.
type stdInfo struct {
	name    string // sentence case
	macName string // where the macOS form words it differently
	symbol  string
	menu    StdMenu
	asks    bool // asks for more before acting, on both forms
	macAsks bool // asks on the macOS form only, as Settings… does
	win     []Shortcut
	mac     []Shortcut
	// focused is the unison command a focused control answers for it, such
	// as a text field's cut, which it does in preference to the app's.
	focused int
}

func k(key unison.KeyCode, mods Mods) Shortcut { return Shortcut{Key: key, Mods: mods} }

// stdCommands are the standard commands. The Windows keys are Microsoft's
// (KEY-2, with the owner's decisions of 2026-10-23: Redo on Ctrl+Y and
// Ctrl+Shift+Z, Find on Ctrl+F and Ctrl+E, Ctrl+Q never bound); the macOS
// keys are Apple's.
var stdCommands = [stdCount]stdInfo{
	StdNew:    {name: "New", symbol: "plus", menu: MenuFile, win: []Shortcut{k(unison.KeyN, ModPrimary)}, mac: []Shortcut{k(unison.KeyN, ModPrimary)}},
	StdOpen:   {name: "Open", symbol: "open", menu: MenuFile, asks: true, win: []Shortcut{k(unison.KeyO, ModPrimary)}, mac: []Shortcut{k(unison.KeyO, ModPrimary)}},
	StdSave:   {name: "Save", symbol: "save", menu: MenuFile, win: []Shortcut{k(unison.KeyS, ModPrimary)}, mac: []Shortcut{k(unison.KeyS, ModPrimary)}},
	StdSaveAs: {name: "Save as", menu: MenuFile, asks: true, mac: []Shortcut{k(unison.KeyS, ModPrimary|ModShift)}},
	StdClose:  {name: "Close", menu: MenuFile, win: []Shortcut{k(unison.KeyW, ModPrimary), k(unison.KeyF4, ModPrimary)}, mac: []Shortcut{k(unison.KeyW, ModPrimary)}},
	StdPrint:  {name: "Print", symbol: "print", menu: MenuFile, asks: true, win: []Shortcut{k(unison.KeyP, ModPrimary)}, mac: []Shortcut{k(unison.KeyP, ModPrimary)}},
	StdUndo:   {name: "Undo", symbol: "undo", menu: MenuEdit, win: []Shortcut{k(unison.KeyZ, ModPrimary)}, mac: []Shortcut{k(unison.KeyZ, ModPrimary)}},
	StdRedo:   {name: "Redo", symbol: "redo", menu: MenuEdit, win: []Shortcut{k(unison.KeyY, ModPrimary), k(unison.KeyZ, ModPrimary|ModShift)}, mac: []Shortcut{k(unison.KeyZ, ModPrimary|ModShift)}},
	StdCut:    {name: "Cut", symbol: "cut", menu: MenuEdit, win: []Shortcut{k(unison.KeyX, ModPrimary)}, mac: []Shortcut{k(unison.KeyX, ModPrimary)}, focused: unison.CutItemID},
	StdCopy:   {name: "Copy", symbol: "copy", menu: MenuEdit, win: []Shortcut{k(unison.KeyC, ModPrimary)}, mac: []Shortcut{k(unison.KeyC, ModPrimary)}, focused: unison.CopyItemID},
	StdPaste:  {name: "Paste", symbol: "paste", menu: MenuEdit, win: []Shortcut{k(unison.KeyV, ModPrimary)}, mac: []Shortcut{k(unison.KeyV, ModPrimary)}, focused: unison.PasteItemID},
	StdDelete: {name: "Delete", symbol: "trash", menu: MenuEdit, win: []Shortcut{k(unison.KeyDelete, 0)}, mac: []Shortcut{k(unison.KeyBackspace, 0)}, focused: unison.DeleteItemID},
	StdSelectAll: {name: "Select all", menu: MenuEdit, win: []Shortcut{k(unison.KeyA, ModPrimary)}, mac: []Shortcut{k(unison.KeyA, ModPrimary)},
		focused: unison.SelectAllItemID},
	StdFind:       {name: "Find", symbol: "search", menu: MenuEdit, win: []Shortcut{k(unison.KeyF, ModPrimary), k(unison.KeyE, ModPrimary)}, mac: []Shortcut{k(unison.KeyF, ModPrimary)}},
	StdFindNext:   {name: "Find next", menu: MenuEdit, win: []Shortcut{k(unison.KeyF3, 0)}, mac: []Shortcut{k(unison.KeyG, ModPrimary)}},
	StdReplace:    {name: "Replace", menu: MenuEdit, asks: true, win: []Shortcut{k(unison.KeyH, ModPrimary)}, mac: []Shortcut{k(unison.KeyF, ModPrimary|ModAlt)}},
	StdGoTo:       {name: "Go to", menu: MenuEdit, asks: true, win: []Shortcut{k(unison.KeyG, ModPrimary)}},
	StdRename:     {name: "Rename", symbol: "rename", menu: MenuEdit, win: []Shortcut{k(unison.KeyF2, 0)}},
	StdRefresh:    {name: "Refresh", symbol: "refresh", menu: MenuView, win: []Shortcut{k(unison.KeyF5, 0)}, mac: []Shortcut{k(unison.KeyR, ModPrimary)}},
	StdProperties: {name: "Properties", macName: "Get info", menu: MenuFile, win: []Shortcut{k(unison.KeyReturn, ModAlt)}, mac: []Shortcut{k(unison.KeyI, ModPrimary)}},
	StdFullScreen: {name: "Full screen", macName: "Enter full screen", symbol: "full-screen", menu: MenuView, win: []Shortcut{k(unison.KeyF11, 0)}, mac: []Shortcut{k(unison.KeyF, ModPrimary|ModControl)}},
	StdZoomIn:     {name: "Zoom in", symbol: "zoom-in", menu: MenuView, win: []Shortcut{k(unison.KeyEqual, ModPrimary)}, mac: []Shortcut{k(unison.KeyEqual, ModPrimary)}},
	StdZoomOut:    {name: "Zoom out", symbol: "zoom-out", menu: MenuView, win: []Shortcut{k(unison.KeyMinus, ModPrimary)}, mac: []Shortcut{k(unison.KeyMinus, ModPrimary)}},
	StdActualSize: {name: "Actual size", menu: MenuView, win: []Shortcut{k(unison.Key0, ModPrimary)}, mac: []Shortcut{k(unison.Key0, ModPrimary)}},
	StdNewTab:     {name: "New tab", menu: MenuFile, win: []Shortcut{k(unison.KeyT, ModPrimary)}, mac: []Shortcut{k(unison.KeyT, ModPrimary)}},
	StdCloseTab:   {name: "Close tab", menu: MenuFile, win: []Shortcut{k(unison.KeyW, ModPrimary), k(unison.KeyF4, ModPrimary)}, mac: []Shortcut{k(unison.KeyW, ModPrimary)}},
	StdReopenTab:  {name: "Reopen closed tab", menu: MenuFile, win: []Shortcut{k(unison.KeyT, ModPrimary|ModShift)}, mac: []Shortcut{k(unison.KeyT, ModPrimary|ModShift)}},
	StdNextTab: {name: "Next tab", menu: MenuView, win: []Shortcut{k(unison.KeyTab, ModControl), k(unison.KeyPageDown, ModControl)},
		mac: []Shortcut{k(unison.KeyTab, ModControl), k(unison.KeyCloseBracket, ModPrimary|ModShift)}},
	StdPreviousTab: {name: "Previous tab", menu: MenuView, win: []Shortcut{k(unison.KeyTab, ModControl|ModShift), k(unison.KeyPageUp, ModControl)},
		mac: []Shortcut{k(unison.KeyTab, ModControl|ModShift), k(unison.KeyOpenBracket, ModPrimary|ModShift)}},
	StdBack: {name: "Back", symbol: "arrow-left", menu: MenuView, win: []Shortcut{k(unison.KeyLeft, ModAlt)},
		mac: []Shortcut{k(unison.KeyOpenBracket, ModPrimary)}},
	StdSettings: {name: "Settings", symbol: "settings", menu: MenuTools, macAsks: true, mac: []Shortcut{k(unison.KeyComma, ModPrimary)}},
	StdHelp:     {name: "Help", symbol: "question", menu: MenuHelp, win: []Shortcut{k(unison.KeyF1, 0)}, mac: []Shortcut{k(unison.KeySlash, ModPrimary|ModShift)}},
	StdAbout:    {name: "About", menu: MenuHelp},
	StdQuit:     {name: "Exit", macName: "Quit", menu: MenuFile, mac: []Shortcut{k(unison.KeyQ, ModPrimary)}},
}

// stdInfoOf is a standard command's entry, or the zero entry for StdNone.
func stdInfoOf(s StdCommand) stdInfo {
	if s <= StdNone || s >= stdCount {
		return stdInfo{}
	}
	return stdCommands[s]
}

// StdName is a standard command's name in a form, without its ellipsis.
// Help, About and the macOS form's Quit name the app: "Kvit Notes help",
// "Quit Kvit Notes".
func (u *UI) StdName(s StdCommand) string {
	info := stdInfoOf(s)
	name := info.name
	if u.form == FormMac && info.macName != "" {
		name = info.macName
	}
	if u.AppName != "" {
		switch {
		case s == StdHelp:
			name = u.AppName + " help"
		case s == StdAbout:
			name = "About " + u.AppName
		case s == StdQuit && u.form == FormMac:
			name = "Quit " + u.AppName
		}
	}
	return name
}

// stdName is a standard command's name in a form, without the app's name.
func stdName(s StdCommand, f Form) string {
	info := stdInfoOf(s)
	if f == FormMac && info.macName != "" {
		return info.macName
	}
	return info.name
}

// StdKeys are a standard command's keys in a form.
func StdKeys(s StdCommand, f Form) []Shortcut {
	info := stdInfoOf(s)
	if f == FormMac {
		return info.mac
	}
	return info.win
}

// StdMenuOf is the menu a standard command goes in on a form: Settings and
// About go in the app menu on the macOS form, and on the Windows form
// Settings goes last in Tools, Exit last in File, and About in no menu,
// since it is a section of the Settings page there (CMD-3).
func StdMenuOf(s StdCommand, f Form) StdMenu {
	switch {
	case f == FormMac && (s == StdSettings || s == StdAbout || s == StdQuit):
		return MenuApp
	case f == FormWindows && s == StdAbout:
		return MenuNone
	}
	return stdInfoOf(s).menu
}

// name is the command's name in the app's form, before the form's case: a
// switch's words for its state on the macOS form, and the ellipsis when it
// asks for more.
func (c *Command) name(u *UI) string {
	n := c.Name
	if n == "" {
		n = u.StdName(c.Std)
	}
	if u.form == FormMac && c.Checked != nil && c.TurnOn != "" && c.TurnOff != "" {
		n = c.TurnOn
		if c.Checked() {
			n = c.TurnOff
		}
	}
	if c.asks(u.form) {
		n += "…"
	}
	return n
}

// asks reports whether the command's name takes an ellipsis in a form.
func (c *Command) asks(f Form) bool {
	info := stdInfoOf(c.Std)
	return c.Asks || info.asks || (f == FormMac && (c.OpensWindow || info.macAsks))
}

// Label is the command's name as it is shown in the app's form: its words
// for its state, its ellipsis, and the form's case for a command.
func (c *Command) Label(u *UI) string { return u.CommandCase(c.name(u)) }

// Plain is the command's name in the app's form without the ellipsis and
// the access-key marker, as a dialog it opens is titled.
func (c *Command) Plain(u *UI) string {
	n := c.Name
	if n == "" {
		n = u.StdName(c.Std)
	}
	return plainWords(u.CommandCase(n))
}

// SymbolName is the command's symbol: its own, or its standard one.
func (c *Command) SymbolName() string {
	if c.Symbol != "" {
		return c.Symbol
	}
	return stdInfoOf(c.Std).symbol
}

// Shortcuts are the keys that run the command in a form: its standard keys,
// then its own.
func (c *Command) Shortcuts(f Form) []Shortcut {
	keys := append([]Shortcut(nil), StdKeys(c.Std, f)...)
	if !c.Shortcut.IsZero() {
		keys = append(keys, c.Shortcut)
	}
	return keys
}

// FirstShortcut is the key a menu line and a tooltip show for the command
// in a form, or the zero shortcut.
func (c *Command) FirstShortcut(f Form) Shortcut {
	if keys := c.Shortcuts(f); len(keys) > 0 {
		return keys[0]
	}
	return Shortcut{}
}

// CanRun reports whether the command can act now.
func (c *Command) CanRun() bool { return c.Run != nil && (c.Enabled == nil || c.Enabled()) }

// IsChecked reports whether a switch is on now; false for a plain command.
func (c *Command) IsChecked() bool { return c.Checked != nil && c.Checked() }

// Do runs the command if it can act, and reports whether it did. Every
// window is drawn again afterwards, since what it did may have enabled or
// disabled other commands.
func (c *Command) Do() bool {
	if !c.CanRun() {
		return false
	}
	c.Run()
	for _, w := range unison.Windows() {
		w.MarkForRedraw()
	}
	return true
}

// Refresh draws every window again, so the controls and menu lines made from
// commands show whether each can act now. An app calls it when something
// changes what its commands can do other than a command running, such as
// data arriving in the background; a command that runs redraws by itself.
func (u *UI) Refresh() {
	for _, w := range unison.Windows() {
		w.MarkForRedraw()
	}
}

// tooltip is what a control running the command says while the pointer
// rests on it: the name, with the shortcut in parentheses on the Windows
// form ("Bold (Ctrl+B)", KEY-14, WRD-13), and the explanation on a line of
// its own.
func (c *Command) tooltip(u *UI, name string) string {
	if u.form == FormWindows {
		if s := c.FirstShortcut(u.form); !s.IsZero() {
			name += " (" + u.ShortcutText(s) + ")"
		}
	}
	return joinLines(name, c.Explanation)
}
