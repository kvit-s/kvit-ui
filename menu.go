package kvitui

import (
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// MenuItem is one line of a menu opened with ShowMenu.
type MenuItem struct {
	// Text is what the line says. An "&" before a letter makes it the line's
	// access key, as in the Qt library's menus: "Open &Folder…" is drawn
	// with the F underlined and chosen by typing F while the menu is open.
	// "&&" is an ampersand. macOS, which has no access keys, draws the words
	// without the marker.
	Text string
	// Key is the shortcut shown at the right of the line, which is how a menu
	// teaches that there is a faster way; the zero value shows none.
	Key unison.KeyBinding
	// Symbol is an optional meaning name drawn before the words.
	Symbol string
	// Explanation is one sentence saying what choosing the line does, for a
	// line whose words cannot say it: shown as its tooltip and told to a
	// screen reader, never in place of the words.
	Explanation string
	// Danger draws a destructive line in the danger colour. It is separated
	// by more than colour: the caller puts a Separator before it, since a red
	// line among black ones is the commonest colour-only distinction there is
	// and misreading it costs the most.
	Danger bool
	// Checked marks the entry in use, with a tick in place of its symbol and
	// its words in bold, so it still stands out in grey.
	Checked bool
	// Disabled draws the line and does nothing when chosen.
	Disabled bool
	// Separator makes the line a divider instead, as goes before a
	// destructive entry at the bottom of a menu.
	Separator bool
	// OnSelect runs when the line is chosen.
	OnSelect func()
	// Items make the line open a submenu of its own instead of acting: beside
	// the line, when the pointer rests on it or on Right, Return or Space, and
	// closed again by Left or Escape. The line draws a chevron where a
	// shortcut would go.
	Items []MenuItem
}

// AccessText splits a menu line's words at the access key: "Open &Folder…"
// is "Open Folder…" with the key F at rune 5. key is 0 and at -1 for words
// without one.
func AccessText(s string) (plain string, key rune, at int) {
	at = -1
	var b []rune
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if rs[i] == '&' && i+1 < len(rs) {
			i++
			if rs[i] != '&' && key == 0 {
				key, at = unicode.ToLower(rs[i]), len(b)
			}
		}
		b = append(b, rs[i])
	}
	return string(b), key, at
}

// PlainMenuText escapes the ampersands of words that are not a menu line's
// own, such as a name read from disk, so none is taken for an access key.
func PlainMenuText(s string) string { return strings.ReplaceAll(s, "&", "&&") }

// menuIDs numbers the lines of native menus, clear of the ids unison and an
// application's menu bar use.
var menuIDs = 0x40000000

// ShowMenu opens a menu of items under anchor, and returns the function that
// closes it. On Windows and Linux it is the Kvit menu, drawn in the window's
// popup layer as the Qt library's is: symbols, a destructive line in the
// danger colour, explanations as tooltips, and the arrow keys, Home, End,
// Return and Escape. On macOS it is the system's own menu, which draws only
// the words, the tick and the shortcut.
func (u *UI) ShowMenu(anchor unison.Paneler, title string, items []MenuItem) (closeMenu func()) {
	return u.ShowMenuAt(anchor, anchor.AsPanel().ContentRect(true), title, items)
}

// ShowMenuAt opens a menu of items with its top left corner at the bottom
// left of one part of owner, given in owner's own coordinates, such as one
// column of a table's header or the point a right-click landed on, and
// returns the function that closes it. A menu with no room below opens
// above, and one with no room to the right is pulled back inside the window.
func (u *UI) ShowMenuAt(owner unison.Paneler, part geom.Rect, title string, items []MenuItem) (closeMenu func()) {
	if runtime.GOOS == "darwin" || u.windowOf(owner) == nil {
		return u.nativeMenuAt(owner, part, title, items)
	}
	m := newMenu(u, title, items)
	m.open(owner, part)
	return m.close
}

// nativeMenuAt opens unison's menu: the system's own on macOS, and drawn by
// unison outside a Kvit window elsewhere.
func (u *UI) nativeMenuAt(owner unison.Paneler, part geom.Rect, title string, items []MenuItem) func() {
	f := unison.DefaultMenuFactory()
	m := nativeMenu(f, title, items)
	// No line is lit until the pointer or an arrow key reaches one, as the Qt
	// menu opens.
	m.Popup(owner.AsPanel().RectToRoot(part), -1)
	return m.Dispose
}

// nativeMenu builds unison's menu of items, a submenu for a line with items.
func nativeMenu(f unison.MenuFactory, title string, items []MenuItem) unison.Menu {
	menuIDs++
	m := f.NewMenu(menuIDs, title, nil)
	for _, it := range items {
		if it.Separator {
			m.InsertSeparator(-1, true)
			continue
		}
		words, _, _ := AccessText(it.Text)
		if len(it.Items) > 0 {
			m.InsertMenu(-1, nativeMenu(f, words, it.Items))
			continue
		}
		menuIDs++
		it := it
		mi := f.NewItem(menuIDs, words, it.Key,
			func(unison.MenuItem) bool { return !it.Disabled },
			func(unison.MenuItem) {
				if it.OnSelect != nil {
					it.OnSelect()
				}
			})
		if it.Checked {
			mi.SetCheckState(check.On)
		}
		m.InsertItem(-1, mi)
	}
	return m
}

// menu is the Kvit menu: a list of commands in the window's popup layer. It
// takes the keyboard focus while open and gives it back when it closes.
type menu struct {
	unison.Panel
	ui       *UI
	title    string
	items    []MenuItem
	lit      int // the lit line, -1 for none
	pressed  bool
	hide     func()
	previous *unison.Panel
	tip      partTip
	parent   *menu // the menu this one opened from, for a submenu
	sub      *menu // the submenu open beside a line, or nil
	subLine  int   // the line the submenu belongs to
	hoverGen int   // counts pointer moves, so a late submenu opening can tell it is stale
}

func newMenu(ui *UI, title string, items []MenuItem) *menu {
	m := &menu{ui: ui, title: title, items: items, lit: -1}
	m.Self = m
	m.tip = partTip{ui: ui, owner: m}
	m.SetFocusable(true)
	m.SetSizer(m.sizes)
	m.DrawCallback = m.draw
	m.KeyDownCallback = m.keyDown
	m.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool { m.hover(m.lineAt(where)); return true }
	m.MouseDragCallback = func(where geom.Point, _ int, _ mod.Modifiers) bool { m.hover(m.lineAt(where)); return true }
	// The first move over the menu arrives as the pointer entering it.
	m.MouseEnterCallback = m.MouseMoveCallback
	m.MouseExitCallback = func() bool {
		// The pointer leaving for the open submenu leaves its line lit.
		if m.sub == nil {
			m.light(-1)
		}
		return true
	}
	m.MouseDownCallback = func(where geom.Point, _, _ int, _ mod.Modifiers) bool {
		m.pressed = true
		m.light(m.lineAt(where))
		return true
	}
	m.MouseUpCallback = func(where geom.Point, _ int, _ mod.Modifiers) bool {
		if m.pressed {
			m.pressed = false
			if i := m.lineAt(where); i >= 0 {
				if len(m.items[i].Items) > 0 {
					m.openSub(i, false)
				} else {
					m.choose(i)
				}
			}
		}
		return true
	}
	m.FrameChangeCallback = func() {
		if w := m.Window(); w != nil {
			ui.watchWindow(w)
		}
	}
	return m
}

// frame is the width of the edge and the padding inside it.
func (m *menu) frame() float32 {
	i := m.ui.Interface
	return float32(i.Hairline() + i.SpaceSnug())
}

func (m *menu) lineHeight(it MenuItem) float32 {
	if it.Separator {
		return float32(m.ui.Interface.Hairline())
	}
	return float32(m.ui.Interface.RowHeightSlim())
}

// words lays out a line's words, cut short to width when width is above 0.
func (m *menu) words(it MenuItem, width float32) *text.Layout {
	ui, t := m.ui, m.ui.Theme.Tokens()
	ink := t.TextPrimary
	switch {
	case it.Disabled:
		ink = t.TextDisabled
	case it.Danger:
		ink = t.Danger
	}
	weight := text.Regular
	if it.Checked {
		weight = text.Bold
	}
	st := ui.Chrome(ui.Size(RoleBody), weight, ink)
	words, _, at := AccessText(it.Text)
	spans := []text.Span{{Text: words, Style: st}}
	if at >= 0 {
		rs := []rune(words)
		key := st
		key.Underline = true
		spans = []text.Span{{Text: string(rs[:at]), Style: st}, {Text: string(rs[at : at+1]), Style: key}, {Text: string(rs[at+1:]), Style: st}}
	}
	return ui.Fonts.Layout(spans, text.Options{MaxWidth: width, Elide: width > 0})
}

func (m *menu) shortcut(it MenuItem) *text.Layout {
	if it.Key.KeyCode == 0 {
		return nil
	}
	ui := m.ui
	return ui.Fonts.Layout([]text.Span{{Text: it.Key.String(), Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, ui.Theme.Tokens().TextFaint)}}, text.Options{})
}

// sizes: 200 design pixels wide, or wider for a long line, and as tall as its
// lines.
func (m *menu) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	i := m.ui.Interface
	w := float32(i.Px(200)) - 2*m.frame()
	h := float32(0)
	for _, it := range m.items {
		h += m.lineHeight(it)
		if it.Separator {
			continue
		}
		lw, _ := m.words(it, 0).Size()
		need := lw + float32(i.Px(72))
		if k := m.shortcut(it); k != nil {
			kw, _ := k.Size()
			need += kw
		} else if len(it.Items) > 0 {
			need += float32(i.IconSizeSmall())
		}
		w = max(w, float32(i.Px(180)), need)
	}
	size := geom.NewSize(w+2*m.frame(), h+2*m.frame())
	return size, size, size
}

// lineBox is where a line is.
func (m *menu) lineBox(i int) geom.Rect {
	r := m.ContentRect(false)
	y := r.Y + m.frame()
	for j := range i {
		y += m.lineHeight(m.items[j])
	}
	return geom.NewRect(r.X+m.frame(), y, r.Width-2*m.frame(), m.lineHeight(m.items[i]))
}

// lineAt is the line that can be chosen under a point, or -1.
func (m *menu) lineAt(where geom.Point) int {
	for i, it := range m.items {
		if !it.Separator && !it.Disabled && where.In(m.lineBox(i)) {
			return i
		}
	}
	return -1
}

// light lights a line, and shows its explanation beside it.
func (m *menu) light(i int) {
	if i == m.lit {
		return
	}
	m.lit = i
	m.MarkForRedraw()
	if i < 0 || m.items[i].Explanation == "" {
		m.tip.clear()
		return
	}
	m.tip.tooltip(i, m.items[i].Explanation, func() geom.Rect { return m.lineBox(i) })
}

// hover lights the line under the pointer, and after a moment's rest opens
// its submenu, or closes the submenu of another line: the pause lets the
// pointer cross a neighbouring line on its way into a submenu.
func (m *menu) hover(i int) {
	if i < 0 && m.sub != nil {
		return
	}
	m.light(i)
	m.hoverGen++
	gen := m.hoverGen
	unison.InvokeTaskAfter(func() {
		if gen != m.hoverGen || m.hide == nil || m.lit != i || i < 0 {
			return
		}
		switch {
		case len(m.items[i].Items) > 0 && (m.sub == nil || m.subLine != i):
			m.openSub(i, false)
		case len(m.items[i].Items) == 0 && m.sub != nil:
			m.sub.close()
		}
	}, 200*time.Millisecond)
}

// openSub opens a line's submenu beside the line, taking the keyboard, with
// its first line lit when first is set, as Right and Return open it.
func (m *menu) openSub(i int, first bool) {
	if m.sub != nil {
		if m.subLine == i {
			if first {
				m.sub.light(m.sub.step(-1, 1))
			}
			return
		}
		m.sub.close()
	}
	w := m.ui.windowOf(m)
	if w == nil {
		return
	}
	it := m.items[i]
	subTitle, _, _ := AccessText(it.Text)
	sub := newMenu(m.ui, subTitle, it.Items)
	sub.parent, sub.previous = m, m.AsPanel()
	m.sub, m.subLine = sub, i
	m.light(i)
	line := func() geom.Rect { return partIn(m, m.lineBox(i)) }
	outer := func() geom.Rect { return anchorIn(m) }
	sub.hide = w.Show(&Popup{Panel: sub, Anchor: m, OnEscape: sub.close, OnPressOutside: sub.close,
		Place: func(bounds geom.Rect, size geom.Size) geom.Rect {
			// Beside the menu, its first line level with the line it came from,
			// or on the other side where there is no room.
			a, o := line(), outer()
			x, y := o.Right(), a.Y-sub.frame()
			if x+size.Width > bounds.Right() {
				x = max(bounds.X, o.X-size.Width)
			}
			if y+size.Height > bounds.Bottom() {
				y = max(bounds.Y, bounds.Bottom()-size.Height)
			}
			return geom.NewRect(x, y, size.Width, size.Height)
		}})
	if first {
		sub.light(sub.step(-1, 1))
	}
	unison.InvokeTask(func() {
		if sub.hide != nil {
			sub.RequestFocus()
		}
	})
	m.MarkForRedraw()
}

// root is the menu a chain of submenus opened from.
func (m *menu) root() *menu {
	for m.parent != nil {
		m = m.parent
	}
	return m
}

// holdsFocus reports whether the menu or a submenu of it holds the focus.
func (m *menu) holdsFocus() bool {
	for q := m; q != nil; q = q.sub {
		if q.Focused() {
			return true
		}
	}
	return false
}

func (m *menu) open(owner unison.Paneler, part geom.Rect) {
	w := m.ui.windowOf(owner)
	if w == nil {
		return
	}
	m.previous = w.CurrentFocus()
	below := func() geom.Rect { return partIn(owner, part) }
	m.hide = w.Show(&Popup{Panel: m, Anchor: owner, OnEscape: m.close, OnPressOutside: m.close,
		Place: func(bounds geom.Rect, size geom.Size) geom.Rect {
			a := below()
			x, y := a.X, a.Bottom()
			if x+size.Width > bounds.Right() {
				x = max(bounds.X, bounds.Right()-size.Width)
			}
			if y+size.Height > bounds.Bottom() {
				y = max(bounds.Y, a.Y-size.Height)
			}
			return geom.NewRect(x, y, size.Width, size.Height)
		}})
	unison.InvokeTask(func() {
		if m.hide != nil {
			m.RequestFocus()
		}
	})
}

// close takes the menu and any submenu of it away, and gives the focus back
// to where it was: a submenu's goes back to the menu it opened from.
func (m *menu) close() {
	if m.hide == nil {
		return
	}
	focused := m.holdsFocus()
	if m.sub != nil {
		m.sub.close()
	}
	m.tip.clear()
	m.hide()
	m.hide = nil
	if p := m.parent; p != nil && p.sub == m {
		p.sub = nil
		p.MarkForRedraw()
	}
	if focused && m.previous != nil && m.previous.Window() != nil {
		m.previous.RequestFocus()
	}
}

// choose closes the menu and runs the line, in that order, so a line that
// opens something else opens it over a closed menu.
func (m *menu) choose(i int) {
	it := m.items[i]
	if it.Separator || it.Disabled {
		return
	}
	if len(it.Items) > 0 {
		m.openSub(i, true)
		return
	}
	// A line in a submenu closes the whole chain.
	m.root().close()
	if it.OnSelect != nil {
		it.OnSelect()
	}
}

// step moves the light to the next line that can be chosen, wrapping.
func (m *menu) step(from, by int) int {
	n := len(m.items)
	for range n {
		from = (from + by + n) % n
		if it := m.items[from]; !it.Separator && !it.Disabled {
			return from
		}
	}
	return -1
}

func (m *menu) keyDown(key unison.KeyCode, _ mod.Modifiers, _ bool) bool {
	switch key {
	case unison.KeyDown:
		m.light(m.step(max(-1, m.lit), 1))
	case unison.KeyUp:
		start := m.lit
		if start < 0 {
			start = len(m.items)
		}
		m.light(m.step(start, -1))
	case unison.KeyHome:
		m.light(m.step(-1, 1))
	case unison.KeyEnd:
		m.light(m.step(len(m.items), -1))
	case unison.KeyReturn, unison.KeyNumPadEnter, unison.KeySpace:
		if m.lit >= 0 {
			m.choose(m.lit)
		}
	case unison.KeyRight:
		if m.lit < 0 || len(m.items[m.lit].Items) == 0 {
			return false
		}
		m.openSub(m.lit, true)
	case unison.KeyLeft:
		if m.parent == nil {
			return false
		}
		m.close()
	case unison.KeyTab:
		m.root().close()
	default:
		return m.accessKey(key)
	}
	return true
}

// accessKey chooses the line whose access key was typed, or with several
// lines sharing it, lights the next of them.
func (m *menu) accessKey(key unison.KeyCode) bool {
	var typed rune
	switch {
	case key >= unison.KeyA && key <= unison.KeyZ:
		typed = rune('a' + int(key-unison.KeyA))
	case key >= unison.Key0 && key <= unison.Key9:
		typed = rune('0' + int(key-unison.Key0))
	default:
		return false
	}
	var hits []int
	for i, it := range m.items {
		if _, k, _ := AccessText(it.Text); k == typed && !it.Disabled && !it.Separator {
			hits = append(hits, i)
		}
	}
	switch len(hits) {
	case 0:
		return false
	case 1:
		m.light(hits[0])
		if len(m.items[hits[0]].Items) > 0 {
			m.openSub(hits[0], true)
		} else {
			m.choose(hits[0])
		}
		return true
	}
	next := hits[0]
	for _, h := range hits {
		if h > m.lit {
			next = h
			break
		}
	}
	m.light(next)
	return true
}

func (m *menu) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, i := m.ui, m.ui.Theme.Tokens(), m.ui.Interface
	p := painterFor(gc, ui)
	r := m.ContentRect(false)
	radius := float32(i.RadiusControl())
	p.round(r, radius, t.PopupBackground)
	p.outline(r, radius, float32(i.Hairline()), t.BorderStrong)
	snug := float32(i.SpaceSnug())
	for j, it := range m.items {
		box := m.lineBox(j)
		if it.Separator {
			p.fill(geom.NewRect(box.X+snug, box.Y, max(0, box.Width-2*snug), box.Height), t.Border)
			continue
		}
		if j == m.lit {
			p.fill(box, t.HoverTint)
		}
		symbol, ink := it.Symbol, t.TextMuted
		switch {
		case it.Checked:
			// The entry in use is marked as such whatever else it would show,
			// since the state is what the reader is scanning for.
			symbol, ink = "check", t.Accent
		case it.Danger:
			ink = t.Danger
		}
		if symbol != "" {
			s := float32(i.IconSizeSmall())
			if g, ok := icons.Glyph(symbol); ok {
				drawGlyph(gc, ui, g, s, ink, geom.NewRect(box.X+float32(i.Space()), box.Y+(box.Height-s)/2, s, s))
			}
		}
		right := box.Right() - float32(i.Space())
		if k := m.shortcut(it); k != nil {
			kw, kh := k.Size()
			k.Draw(gc, right-kw, box.Y+(box.Height-kh)/2)
			right -= kw + float32(i.SpaceLoose())
		} else if len(it.Items) > 0 {
			s := float32(i.IconSizeSmall())
			if g, ok := icons.Glyph("chevron-right"); ok {
				drawGlyph(gc, ui, g, s, t.TextMuted, geom.NewRect(right-s, box.Y+(box.Height-s)/2, s, s))
			}
			right -= s + float32(i.SpaceLoose())
		}
		x := box.X + float32(i.Px(30))
		l := m.words(it, max(1, right-x))
		_, h := l.Size()
		l.Draw(gc, x, box.Y+(box.Height-h)/2)
	}
}

// ProvideAccessibility describes the menu and its lines, the focus on the lit
// one.
func (m *menu) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	n.Role = role.Menu
	n.Name = m.title
	for i, it := range m.items {
		if it.Separator {
			continue
		}
		box := m.lineBox(i)
		id := b.AddVirtualChild(i, func(v *accessibility.Node) {
			v.Role = role.MenuItem
			v.Name, _, _ = AccessText(it.Text)
			v.Description = joinWords(it.Explanation, shortcutWords(it.Key))
			v.Bounds = box
			v.Disabled = it.Disabled
			v.Focusable = !it.Disabled
			if it.Checked {
				v.HasCheck, v.Checked = true, check.On
			}
			if !it.Disabled {
				v.Actions = v.Actions.With(accessibility.Press, accessibility.Focus)
			}
			if len(it.Items) > 0 {
				v.Expandable, v.Expanded = true, m.sub != nil && m.subLine == i
				if !it.Disabled {
					v.Actions = v.Actions.With(accessibility.Expand, accessibility.Collapse)
				}
			}
		})
		if i == m.lit {
			b.FocusChild(id)
		}
	}
}

func shortcutWords(k unison.KeyBinding) string {
	if k.KeyCode == 0 {
		return ""
	}
	return k.String()
}

// PerformAccessibilityAction chooses a line, or lights it, for a screen
// reader.
func (m *menu) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	i, ok := req.Key.(int)
	if !ok || i < 0 || i >= len(m.items) || m.items[i].Separator || m.items[i].Disabled {
		return false
	}
	switch req.Action {
	case accessibility.Press, accessibility.Expand:
		m.choose(i)
	case accessibility.Collapse:
		if m.sub != nil && m.subLine == i {
			m.sub.close()
		}
	case accessibility.Focus:
		m.light(i)
	default:
		return false
	}
	return true
}

// contextMenuKey is the client-data key under which a panel keeps what opens
// its context menu.
const contextMenuKey = "kvitui.contextMenu"

// contextOpener opens a panel's context menu for a point in its own
// coordinates, and reports whether it opened one.
type contextOpener func(where geom.Point) bool

// SetContextMenu gives a panel a context menu: what menu returns opens at
// the pointer on a right-click on the panel or anything inside it, and on the
// Menu key or Shift+F10 while the keyboard is in it, under the panel holding
// the focus, or where that panel's ContextMenuAnchor says. menu is given the
// point in the panel's coordinates; no items opens nothing.
func (u *UI) SetContextMenu(p unison.Paneler, menu func(where geom.Point) (title string, items []MenuItem)) {
	panel := p.AsPanel()
	panel.ClientData()[contextMenuKey] = contextOpener(func(where geom.Point) bool {
		title, items := menu(where)
		if len(items) == 0 {
			return false
		}
		u.ShowMenuAt(panel, geom.NewRect(where.X, where.Y, 0, 0), title, items)
		return true
	})
}

// setContextOpener gives a panel a context menu it opens itself, as a
// table's header opens a column's menu under the column.
func setContextOpener(p unison.Paneler, open contextOpener) {
	p.AsPanel().ClientData()[contextMenuKey] = open
}

// isContextMenuKey reports whether a key asks for a context menu: the Menu
// key, or Shift+F10, as on every desktop.
func isContextMenuKey(key unison.KeyCode, mods mod.Modifiers) bool {
	switch key {
	case unison.KeyMenu:
		return mods&mod.NonSticky == 0
	case unison.KeyF10:
		return mods&mod.NonSticky == mod.Shift
	}
	return false
}

// openContextMenu opens the context menu of p, or of the nearest panel above
// it that has one, for a point in p's coordinates. A panel with unison's own
// context menu, and nothing Kvit above it, is left to unison.
func openContextMenu(p *unison.Panel, where geom.Point) bool {
	for q := p; q != nil; q = q.Parent() {
		if open, ok := q.ClientData()[contextMenuKey].(contextOpener); ok {
			return open(q.PointFromRoot(p.PointToRoot(where)))
		}
		if q.ContextMenuCallback != nil {
			return false
		}
	}
	return false
}

// contextKeys opens the context menu of the panel holding the focus, for the
// Menu key and Shift+F10.
func (w *Window) contextKeys(key unison.KeyCode, mods mod.Modifiers) bool {
	if !isContextMenuKey(key, mods) {
		return false
	}
	focus := w.CurrentFocus()
	if focus == nil {
		return false
	}
	at := geom.NewPoint(0, focus.ContentRect(true).Height)
	if a, ok := focus.Self.(interface{ ContextMenuAnchor() geom.Point }); ok {
		at = a.ContextMenuAnchor()
	}
	return openContextMenu(focus, at)
}

// contextPress opens the context menu of what a right-click landed on, given
// in the window content's coordinates, giving it the focus first where it
// can hold it.
func (w *Window) contextPress(where geom.Point) bool {
	content := w.Content()
	target := content.PanelAt(where)
	if target == nil {
		return false
	}
	for q := target; q != nil; q = q.Parent() {
		if _, ok := q.ClientData()[contextMenuKey].(contextOpener); ok {
			if q.Focusable() && !q.Focused() {
				q.RequestFocus()
			}
			break
		}
	}
	return openContextMenu(target, target.PointFromRoot(content.PointToRoot(where)))
}

// applyMenuTheme draws unison's menus in the Kvit colours and type: the popup
// ground with a strong hairline edge, the hover tint on the line under the
// pointer, lines a slim row tall, and the shortcut in the small role.
func (u *UI) applyMenuTheme() {
	t, m := u.Theme.Tokens(), u.Interface
	title := unisonFont(u, m.FontFamily(), u.Size(RoleBody))
	th := &unison.DefaultMenuItemTheme
	th.TitleFont = title
	th.KeyFont = unisonFont(u, m.FontFamily(), u.Size(RoleSmall))
	th.BackgroundColor, th.OnBackgroundColor = Color(t.PopupBackground), Color(t.TextPrimary)
	th.SelectionColor, th.OnSelectionColor = Color(t.HoverTint), Color(t.TextPrimary)
	v := max(0, (float32(m.RowHeightSlim())-title.LineHeight())/2)
	space := float32(m.Space())
	th.ItemBorder = unison.NewEmptyBorder(geom.Insets{Top: v, Bottom: v, Left: space, Right: space})
	th.SeparatorBorder = unison.NewEmptyBorder(geom.NewVerticalInsets(float32(m.SpaceSnug())))
	th.KeyGap = float32(m.SpaceLoose())
	radius := float32(m.RadiusControl())
	unison.DefaultMenuTheme.MenuBorder = unison.NewCompoundBorder(
		unison.NewLineBorder(Color(t.BorderStrong), geom.NewSize(radius, radius), geom.NewUniformInsets(float32(m.Hairline())), false),
		unison.NewEmptyBorder(geom.NewUniformInsets(float32(m.SpaceSnug()))))
}
