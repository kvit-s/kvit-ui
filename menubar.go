package kvitui

import (
	"log/slog"
	"runtime"
	"slices"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// MenuBar is a main window's menus, built from the app's commands, with the
// platform's standard menus in its standard order (CMD-1 to CMD-8). A
// standard command goes in its standard menu at its standard place without
// being told where: Save in File, Copy in Edit, Settings last in Tools on the
// Windows form and in the app menu on the macOS form, Exit last in File on
// the Windows form and Quit last in the app menu on the macOS form. The
// app's own commands go in a standard menu with AddTo or in a menu of their
// own with AddMenu.
//
// On the Windows form, which Linux takes, the bar is drawn by Kvit across
// the top of the window: File, Edit, View, the app's menus, Tools, Help,
// each with an access key. A press on a title opens its menu; Alt on its own
// or F10 puts the keyboard in the bar, and Alt with a title's letter opens
// that menu. On the macOS form on a Mac it is the system's menu bar: the app
// menu, File, Edit, Format, View, the app's menus, Window and Help. On the
// macOS form elsewhere, as in a test or the gallery, Kvit draws the macOS
// menus in the window, so the form can be seen on any system.
//
// Window.SetMenuBar puts the bar in a window and makes the window answer
// the keys of every command on it.
type MenuBar struct {
	unison.Panel
	ui    *UI
	std   []*Command
	extra map[StdMenu][]MenuItem
	own   []BarMenu

	window   *Window
	lit      int   // the title lit while the keyboard is in the bar, -1 for none
	openAt   int   // the title whose menu is open, -1 for none
	open     *menu // the menu open under a title
	previous *unison.Panel
	active   bool // the keyboard is in the bar
	altAlone bool // Alt is down with nothing else pressed since
	native   unison.Menu
}

// BarMenu is one menu of a menu bar: its title and its lines.
type BarMenu struct {
	// Title is the menu's name, with an access key: "&File".
	Title string
	// Items are its lines.
	Items []MenuItem
	// std is the standard menu it is, MenuNone for one of the app's.
	std StdMenu
}

// NewMenuBar returns an empty menu bar.
func NewMenuBar(ui *UI) *MenuBar {
	b := &MenuBar{ui: ui, extra: map[StdMenu][]MenuItem{}, lit: -1, openAt: -1}
	b.Self = b
	OpensPopup(b)
	b.SetSizer(b.sizes)
	b.DrawCallback = b.draw
	b.MouseDownCallback = func(where geom.Point, _, _ int, _ mod.Modifiers) bool {
		if i := b.titleAt(where); i >= 0 {
			if b.openAt == i {
				b.closeMenu(true)
			} else {
				b.openMenu(i, false)
			}
		}
		return true
	}
	b.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool {
		// With one menu open, the pointer moving onto another title opens
		// that one instead, as every menu bar does.
		if i := b.titleAt(where); i >= 0 && b.openAt >= 0 && i != b.openAt {
			b.openMenu(i, false)
		}
		return true
	}
	b.KeyDownCallback = b.keyDown
	b.LostFocusCallback = func() {
		if b.active && b.open == nil {
			b.leave(false)
		}
	}
	return b
}

// Add puts standard commands in their standard menus at their standard
// places in the app's form. A command of the app's own belongs in AddTo or
// AddMenu, and is left out with a warning.
func (b *MenuBar) Add(cmds ...*Command) *MenuBar {
	for _, c := range cmds {
		if c == nil {
			continue
		}
		if c.Std == StdNone {
			slog.Warn("kvitui: MenuBar.Add takes standard commands; give the app's own to AddTo or AddMenu", "command", c.Name)
			continue
		}
		b.std = append(b.std, c)
	}
	b.MarkForLayoutAndRedraw()
	return b
}

// AddTo puts the app's own lines in one of the standard menus, as a group
// of their own: in View before the standard lines (the window's parts and
// panes come first, CMD-6), in File before Exit, in Tools before Settings,
// and after the standard lines elsewhere.
func (b *MenuBar) AddTo(menu StdMenu, items ...MenuItem) *MenuBar {
	b.extra[menu] = append(b.extra[menu], items...)
	b.MarkForLayoutAndRedraw()
	return b
}

// AddMenu adds a menu of the app's own, after View and before Tools (or
// Window): title in sentence case, with its access key, "&Transaction".
func (b *MenuBar) AddMenu(title string, items ...MenuItem) *MenuBar {
	b.own = append(b.own, BarMenu{Title: title, Items: items})
	b.MarkForLayoutAndRedraw()
	return b
}

// Commands are every command on the bar, standard and the app's own, each
// once, in menu order.
func (b *MenuBar) Commands() []*Command {
	var out []*Command
	var walk func(items []MenuItem)
	walk = func(items []MenuItem) {
		for _, it := range items {
			if it.Command != nil && !slices.Contains(out, it.Command) {
				out = append(out, it.Command)
			}
			walk(it.Items)
		}
	}
	for _, m := range b.Menus() {
		walk(m.Items)
	}
	for _, c := range b.std {
		if !slices.Contains(out, c) {
			out = append(out, c)
		}
	}
	return out
}

// stdMenuLines are, per form and menu, the standard commands in their
// standard order, StdNone marking a separator.
var stdMenuLines = map[Form]map[StdMenu][]StdCommand{
	FormWindows: {
		MenuFile: {StdNew, StdNewTab, StdOpen, StdClose, StdCloseTab, StdReopenTab, StdNone, StdSave, StdSaveAs, StdNone,
			StdPrint, StdNone, StdProperties},
		MenuEdit: {StdUndo, StdRedo, StdNone, StdCut, StdCopy, StdPaste, StdNone, StdSelectAll, StdNone, StdDelete, StdRename,
			StdNone, StdFind, StdFindNext, StdReplace, StdGoTo},
		MenuView: {StdBack, StdNone, StdZoomIn, StdZoomOut, StdActualSize, StdNone, StdNextTab, StdPreviousTab, StdNone,
			StdFullScreen, StdRefresh},
		MenuHelp: {StdHelp},
	},
	FormMac: {
		MenuApp: {StdAbout},
		MenuFile: {StdNew, StdNewTab, StdOpen, StdNone, StdClose, StdCloseTab, StdReopenTab, StdSave, StdSaveAs, StdRename,
			StdProperties, StdNone, StdPrint},
		MenuEdit: {StdUndo, StdRedo, StdNone, StdCut, StdCopy, StdPaste, StdDelete, StdSelectAll, StdNone, StdFind,
			StdFindNext, StdReplace, StdGoTo},
		MenuView:   {StdBack, StdNone, StdRefresh, StdZoomIn, StdZoomOut, StdActualSize, StdNone, StdFullScreen},
		MenuWindow: {StdNextTab, StdPreviousTab},
		MenuHelp:   {StdHelp},
	},
}

// stdTitles are the standard menus' titles on the Windows form; the macOS
// form draws them without the access key, and names the app menu after the
// app.
var stdTitles = map[StdMenu]string{
	MenuFile: "&File", MenuEdit: "&Edit", MenuFormat: "F&ormat", MenuView: "&View", MenuTools: "&Tools",
	MenuWindow: "&Window", MenuHelp: "&Help",
}

// Menus are the bar's menus as the app's form arranges them, the empty ones
// left out.
func (b *MenuBar) Menus() []BarMenu {
	f := b.ui.form
	order := []StdMenu{MenuFile, MenuEdit, MenuFormat, MenuView, MenuNone, MenuTools, MenuHelp}
	if f == FormMac {
		order = []StdMenu{MenuApp, MenuFile, MenuEdit, MenuFormat, MenuView, MenuNone, MenuWindow, MenuHelp}
	}
	var out []BarMenu
	for _, sm := range order {
		if sm == MenuNone {
			out = append(out, b.own...)
			continue
		}
		if sm == MenuWindow && f != FormMac {
			continue
		}
		items := b.stdLines(sm)
		if len(items) == 0 {
			continue
		}
		title := stdTitles[sm]
		if sm == MenuApp {
			title = b.ui.AppName
			if title == "" {
				title = "App"
			}
		}
		out = append(out, BarMenu{Title: title, Items: items, std: sm})
	}
	return out
}

// stdLines are one standard menu's lines: the standard commands the app
// has in their order, with the app's own lines in their place.
func (b *MenuBar) stdLines(sm StdMenu) []MenuItem {
	f := b.ui.form
	has := func(s StdCommand) *Command {
		for _, c := range b.std {
			if c.Std == s && StdMenuOf(s, f) == sm {
				return c
			}
		}
		return nil
	}
	var std []MenuItem
	for _, s := range stdMenuLines[f][sm] {
		if s == StdNone {
			std = append(std, MenuItem{Separator: true})
		} else if c := has(s); c != nil {
			std = append(std, MenuItem{Command: c})
		}
	}
	// Commands placed in this menu that its standard order does not list,
	// such as Settings in Tools on the Windows form, unless the window's
	// sidebar holds it, which is where that form puts it (WIN-16).
	var last []MenuItem
	for _, c := range b.std {
		if StdMenuOf(c.Std, f) == sm && !slices.Contains(stdMenuLines[f][sm], c.Std) {
			if f == FormWindows && c.Std == StdSettings && b.window != nil && slices.Contains(b.window.sidebarCommands(), c) {
				continue
			}
			last = append(last, MenuItem{Command: c})
		}
	}
	extra := b.extra[sm]
	sep := []MenuItem{{Separator: true}}
	var items []MenuItem
	switch {
	case sm == MenuWindow:
		// The Window menu is always there on the macOS form (CMD-7).
		items = slices.Concat(b.macWindowLines(), sep, std, sep, extra)
	case sm == MenuView:
		items = slices.Concat(extra, sep, std, sep, last)
	case sm == MenuApp:
		// About, Settings…, the app's own, then the system's, then Quit.
		var settings, quit []MenuItem
		for _, it := range last {
			switch it.Command.Std {
			case StdSettings:
				settings = append(settings, it)
			case StdQuit:
				quit = append(quit, it)
			}
		}
		items = slices.Concat(std, sep, settings, sep, extra, sep, b.macSystemLines(), sep, quit)
	default:
		// Exit and Settings come last, after the app's own lines.
		items = slices.Concat(std, sep, extra, sep, last)
	}
	if !slices.ContainsFunc(items, func(it MenuItem) bool { return !it.Separator }) {
		return nil
	}
	return tidySeparators(items)
}

// macSystemLines are the app menu's lines the system owns, as the macOS
// form's preview shows them away from a Mac, where they cannot act.
func (b *MenuBar) macSystemLines() []MenuItem {
	name := b.ui.AppName
	if name == "" {
		name = "app"
	}
	return []MenuItem{
		{Text: "Hide " + name, Verbatim: true, Shortcut: Keys(unison.KeyH, ModPrimary), Disabled: true},
		{Text: "Hide Others", Verbatim: true, Shortcut: Keys(unison.KeyH, ModPrimary|ModAlt), Disabled: true},
		{Text: "Show All", Verbatim: true, Disabled: true},
	}
}

// macWindowLines are the Window menu's lines the system owns, as the macOS
// form's preview shows them away from a Mac, where they cannot act.
func (b *MenuBar) macWindowLines() []MenuItem {
	return []MenuItem{
		{Text: "Minimize", Verbatim: true, Shortcut: Keys(unison.KeyM, ModPrimary), Disabled: true},
		{Text: "Zoom", Verbatim: true, Disabled: true},
	}
}

// isNative reports whether the bar is the system's, which it is on the macOS
// form on a Mac.
func (b *MenuBar) isNative() bool { return runtime.GOOS == "darwin" && b.ui.form == FormMac }

// titleLayouts are the titles as drawn, without their access-key markers.
func (b *MenuBar) titleLayouts() []*text.Layout {
	ui, t := b.ui, b.ui.Theme.Tokens()
	var out []*text.Layout
	for _, m := range b.Menus() {
		words, _, at := AccessText(m.Title)
		weight := text.Regular
		if m.std == MenuApp {
			weight = text.Bold
		}
		st := ui.Chrome(ui.Size(RoleBody), weight, t.TextPrimary)
		spans := []text.Span{{Text: words, Style: st}}
		// The access keys are underlined while the keyboard is in the bar or
		// Alt is held, or always where the system asks (KEY-15); the macOS
		// form has none.
		if at >= 0 && !ui.Mac() && (b.active || b.altAlone || ui.UnderlineAccessKeys) {
			rs := []rune(words)
			key := st
			key.Underline = true
			spans = []text.Span{{Text: string(rs[:at]), Style: st}, {Text: string(rs[at : at+1]), Style: key}, {Text: string(rs[at+1:]), Style: st}}
		}
		out = append(out, ui.Fonts.Layout(spans, text.Options{}))
	}
	return out
}

// titleBoxes are where the titles are, in the bar's coordinates.
func (b *MenuBar) titleBoxes() []geom.Rect {
	m := b.ui.Interface
	r := b.ContentRect(false)
	pad := float32(m.Space())
	x := r.X + float32(m.SpaceNear())
	var out []geom.Rect
	for _, l := range b.titleLayouts() {
		w, _ := l.Size()
		out = append(out, geom.NewRect(x, r.Y, w+2*pad, r.Height))
		x += w + 2*pad
	}
	return out
}

func (b *MenuBar) titleAt(where geom.Point) int {
	for i, box := range b.titleBoxes() {
		if where.In(box) {
			return i
		}
	}
	return -1
}

func (b *MenuBar) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := b.ui.Interface
	if b.isNative() {
		return geom.Size{}, geom.Size{}, geom.Size{}
	}
	w := float32(0)
	for _, box := range b.titleBoxes() {
		w = max(w, box.Right())
	}
	h := float32(m.RowHeightCompact())
	return geom.NewSize(w, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (b *MenuBar) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := b.ui, b.ui.Theme.Tokens(), b.ui.Interface
	p := painterFor(gc, ui)
	r := b.ContentRect(false)
	p.fill(r, t.WindowBackground)
	p.fill(geom.NewRect(r.X, r.Bottom()-float32(m.Hairline()), r.Width, float32(m.Hairline())), t.Border)
	boxes := b.titleBoxes()
	for i, l := range b.titleLayouts() {
		box := boxes[i]
		switch {
		case i == b.openAt:
			p.fill(box, t.SelectionTint)
		case i == b.lit:
			p.fill(box, t.HoverTint)
		}
		w, h := l.Size()
		l.Draw(gc, box.X+(box.Width-w)/2, box.Y+(box.Height-h)/2)
	}
}

// OpenMenu opens the menu under the title at index i, as a press on the
// title does.
func (b *MenuBar) OpenMenu(i int) { b.openMenu(i, false) }

// openMenu opens the menu under a title, with its first line lit when
// first is set, as opening it from the keyboard does.
func (b *MenuBar) openMenu(i int, first bool) {
	menus := b.Menus()
	if i < 0 || i >= len(menus) {
		return
	}
	w := b.ui.windowOf(b)
	if w == nil {
		return
	}
	if b.open != nil {
		// Moving from one menu to the next keeps the keyboard in the bar.
		prev := b.open
		b.open = nil
		prev.onSide = nil
		prev.close()
	}
	if b.previous == nil {
		b.previous = w.CurrentFocus()
	}
	focus := b.previous
	items := b.ui.menuLines(menus[i].Items, false, focus)
	title, _, _ := AccessText(menus[i].Title)
	mn := newMenu(b.ui, title, items)
	mn.previous = b.previous
	mn.onSide = func(dir int) {
		n := len(b.Menus())
		b.openMenu((i+dir+n)%n, true)
	}
	mn.onClose = func() {
		if b.open == mn {
			b.open, b.openAt = nil, -1
			b.leave(true)
		}
	}
	b.open, b.openAt, b.lit = mn, i, i
	box := b.titleBoxes()[i]
	mn.open(b, geom.NewRect(box.X, box.Y, box.Width, box.Height))
	if first {
		mn.light(mn.step(-1, 1))
	}
	b.MarkForRedraw()
}

// closeMenu closes the open menu, keeping the keyboard in the bar when keep
// is set.
func (b *MenuBar) closeMenu(keep bool) {
	if b.open == nil {
		return
	}
	mn := b.open
	b.open, b.openAt = nil, -1
	mn.onClose = nil
	mn.close()
	if keep && b.active {
		b.RequestFocus()
	} else {
		b.leave(true)
	}
	b.MarkForRedraw()
}

// enter puts the keyboard in the bar, its first title lit, as Alt or F10
// does.
func (b *MenuBar) enter() {
	if len(b.Menus()) == 0 || b.isNative() {
		return
	}
	w := b.ui.windowOf(b)
	if w == nil {
		return
	}
	if !b.active {
		b.previous = w.CurrentFocus()
	}
	b.active, b.lit = true, max(0, b.lit)
	b.SetFocusable(true)
	b.RequestFocus()
	b.MarkForRedraw()
}

// leave takes the keyboard out of the bar, back to where it was when
// giveBack is set.
func (b *MenuBar) leave(giveBack bool) {
	previous := b.previous
	b.active, b.lit, b.previous = false, -1, nil
	b.SetFocusable(false)
	b.MarkForRedraw()
	if giveBack && previous != nil && previous.Window() != nil && b.open == nil {
		previous.RequestFocus()
	}
}

func (b *MenuBar) keyDown(key unison.KeyCode, mods mod.Modifiers, _ bool) bool {
	if !b.active {
		return false
	}
	n := len(b.Menus())
	switch key {
	case unison.KeyLeft:
		b.lit = (b.lit - 1 + n) % n
	case unison.KeyRight:
		b.lit = (b.lit + 1) % n
	case unison.KeyDown, unison.KeyReturn, unison.KeyNumPadEnter, unison.KeySpace:
		b.openMenu(b.lit, true)
		return true
	case unison.KeyEscape, unison.KeyF10, unison.KeyTab:
		b.leave(true)
		return true
	default:
		if mods&(mod.Control|mod.Command) == 0 {
			if i := b.accessIndex(key); i >= 0 {
				b.openMenu(i, true)
				return true
			}
		}
		return false
	}
	b.MarkForRedraw()
	return true
}

// accessIndex is the menu whose title has the access key typed, or -1.
func (b *MenuBar) accessIndex(key unison.KeyCode) int {
	typed, ok := typedLetter(key)
	if !ok || b.ui.Mac() {
		return -1
	}
	for i, m := range b.Menus() {
		if _, k, _ := AccessText(m.Title); k == typed {
			return i
		}
	}
	return -1
}

// typedLetter is the lower-case letter or digit a key types.
func typedLetter(key unison.KeyCode) (rune, bool) {
	switch {
	case key >= unison.KeyA && key <= unison.KeyZ:
		return rune('a' + int(key-unison.KeyA)), true
	case key >= unison.Key0 && key <= unison.Key9:
		return rune('0' + int(key-unison.Key0)), true
	}
	return 0, false
}

// windowKeyDown is offered every key press in the bar's window first: Alt
// with a title's letter opens its menu, F10 puts the keyboard in the bar,
// and Alt pressed on its own is noted, to put the keyboard in the bar when
// it is let go with nothing pressed in between.
func (b *MenuBar) windowKeyDown(key unison.KeyCode, mods mod.Modifiers) bool {
	if b.isNative() || b.ui.Mac() {
		return false
	}
	if key == unison.KeyLOption || key == unison.KeyROption {
		b.altAlone = mods&mod.NonSticky&^mod.Option == 0
		b.MarkForRedraw()
		return false
	}
	if b.altAlone {
		b.altAlone = false
		b.MarkForRedraw()
	}
	if key == unison.KeyF10 && mods&mod.NonSticky == 0 {
		if b.active {
			b.leave(true)
		} else {
			b.enter()
		}
		return true
	}
	if mods&mod.NonSticky == mod.Option {
		if i := b.accessIndex(key); i >= 0 {
			b.enter()
			b.openMenu(i, true)
			return true
		}
	}
	return false
}

// windowKeyUp puts the keyboard in the bar when Alt is let go after being
// pressed on its own, and takes it out again on a second Alt.
func (b *MenuBar) windowKeyUp(key unison.KeyCode) {
	if (key == unison.KeyLOption || key == unison.KeyROption) && b.altAlone {
		b.altAlone = false
		if b.active && b.open == nil {
			b.leave(true)
		} else if !b.active {
			b.enter()
		}
		b.MarkForRedraw()
	}
}

// ProvideAccessibility describes the bar and its menus' titles.
func (b *MenuBar) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	n := bl.Node()
	n.Role = role.MenuBar
	n.Name = "Menu bar"
	if b.isNative() {
		n.Ignored = true
		return
	}
	boxes := b.titleBoxes()
	for i, m := range b.Menus() {
		box := boxes[i]
		id := bl.AddVirtualChild(i, func(v *accessibility.Node) {
			v.Role = role.MenuItem
			v.Name, _, _ = AccessText(m.Title)
			v.Bounds = box
			v.Focusable = true
			v.Expandable, v.Expanded = true, b.openAt == i
			v.Actions = v.Actions.With(accessibility.Press, accessibility.Expand, accessibility.Collapse)
		})
		if b.active && i == b.lit {
			bl.FocusChild(id)
		}
	}
}

// PerformAccessibilityAction opens or closes a menu for a screen reader.
func (b *MenuBar) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	i, ok := req.Key.(int)
	if !ok || i < 0 || i >= len(b.Menus()) {
		return false
	}
	switch req.Action {
	case accessibility.Press, accessibility.Expand:
		b.openMenu(i, true)
	case accessibility.Collapse:
		b.closeMenu(false)
	default:
		return false
	}
	return true
}

// installNative puts the bar's menus in the system's menu bar for a window,
// which is what the macOS form is on a Mac. unison's app menu calls the
// settings command "Preferences…", so its line is renamed Settings….
func (b *MenuBar) installNative(w *Window) {
	f := unison.DefaultMenuFactory()
	about, settings := b.stdCommand(StdAbout), b.stdCommand(StdSettings)
	handler := func(c *Command) func(unison.MenuItem) {
		if c == nil {
			return nil
		}
		return func(unison.MenuItem) { b.ui.runCommand(c) }
	}
	f.BarForWindow(w.Window, func(bar unison.Menu) {
		app := unison.NewAppMenu(f, handler(about), handler(settings), nil)
		if mi := app.Item(unison.PreferencesItemID); mi != nil {
			mi.SetTitle(b.ui.CommandCase("Settings…"))
		}
		at := 3
		for _, it := range b.ui.menuLines(b.extra[MenuApp], false, nil) {
			if it.Separator {
				app.InsertSeparator(at, false)
			} else {
				app.InsertItem(at, nativeItem(f, FormMac, it))
			}
			at++
		}
		bar.InsertMenu(-1, app)
		for _, m := range b.Menus() {
			lines := func() []MenuItem {
				var focus *unison.Panel
				if aw := unison.ActiveWindow(); aw != nil {
					focus = aw.CurrentFocus()
				}
				return b.ui.menuLines(m.Items, false, focus)
			}
			switch m.std {
			case MenuApp:
				continue
			case MenuWindow:
				win := unison.NewWindowMenu(f, nil)
				fillNative(win, FormMac, lines())
				bar.InsertMenu(-1, win)
				continue
			}
			title, _, _ := AccessText(m.Title)
			bar.InsertMenu(-1, nativeMenuOf(f, FormMac, b.ui.CommandCase(title), lines))
		}
		if bar.Menu(unison.WindowMenuID) == nil {
			// The Window menu is always there on a Mac (CMD-7), before Help.
			at := -1
			if n := bar.Count(); n > 0 && b.stdHas(MenuHelp) {
				at = n - 1
			}
			bar.InsertMenu(at, unison.NewWindowMenu(f, nil))
		}
	})
}

// stdHas reports whether the bar shows a standard menu.
func (b *MenuBar) stdHas(sm StdMenu) bool {
	return slices.ContainsFunc(b.Menus(), func(m BarMenu) bool { return m.std == sm })
}

// stdCommand is the bar's standard command of a kind, or nil.
func (b *MenuBar) stdCommand(s StdCommand) *Command {
	for _, c := range b.std {
		if c.Std == s {
			return c
		}
	}
	return nil
}

// SetMenuBar puts a menu bar in the window, across the top on the Windows
// form and in the system's menu bar on the macOS form on a Mac, and makes the
// window answer the keys of every command on it.
func (w *Window) SetMenuBar(b *MenuBar) {
	if w.menuBar != nil {
		w.menuBar.window = nil
	}
	w.menuBar = b
	if b != nil {
		b.window = w
		// A refused key is reported by AddCommands; the command still runs
		// from its line.
		_ = w.AddCommands(b.Commands()...)
		if b.isNative() {
			b.installNative(w)
		}
	}
	w.rebuild()
}

// MenuBar is the window's menu bar, or nil.
func (w *Window) MenuBar() *MenuBar { return w.menuBar }

// sidebarCommands are the commands the items of the window's sidebar run.
func (w *Window) sidebarCommands() []*Command {
	if w.sidebar == nil {
		return nil
	}
	if s, ok := w.sidebar.AsPanel().Self.(*Sidebar); ok {
		return s.Commands()
	}
	return nil
}

// OffSurface are the commands the window runs, by their keys or from a
// button, that are on none of its command surfaces: not on its menu bar, not
// in its sidebar, and not on a toolbar in it or in that toolbar's More menu. Every command
// belongs on a visible surface (CMD-1); one reached only by a key, a button
// or a context menu cannot be found.
func (w *Window) OffSurface() []*Command {
	on := w.sidebarCommands()
	if w.menuBar != nil {
		on = append(on, w.menuBar.Commands()...)
	}
	var runs []*Command
	add := func(c *Command) {
		if c != nil && !slices.Contains(runs, c) {
			runs = append(runs, c)
		}
	}
	for _, c := range w.commands {
		add(c)
	}
	var walk func(p *unison.Panel)
	walk = func(p *unison.Panel) {
		switch v := p.Self.(type) {
		case *Toolbar:
			on = append(on, v.Commands()...)
			return
		case *Button:
			add(v.Command)
		case *IconButton:
			add(v.Command)
		}
		for _, c := range p.Children() {
			walk(c)
		}
	}
	walk(w.Content())
	var out []*Command
	for _, c := range runs {
		if !slices.Contains(on, c) {
			out = append(out, c)
		}
	}
	return out
}
