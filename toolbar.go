package kvitui

import (
	"slices"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// ToolGroup is where on a toolbar a command sits.
type ToolGroup int

const (
	// ToolLeading is the start of the bar, where its commands begin.
	ToolLeading ToolGroup = iota
	// ToolCentre is the middle, for navigation or a title.
	ToolCentre
	// ToolTrailing is the end, beside the More button: Help, Settings and
	// the view commands on the Windows form, and the one primary action on
	// the macOS form (CMD-18).
	ToolTrailing
)

// Toolbar is a strip of the commands for a window's main tasks, at the top of
// a main window, made from the app's commands (CMD-18). On the Windows form
// a window with few commands may have a toolbar with a More menu instead of
// a menu bar (CMD-1).
//
// Commands are added in order of importance. When the bar is too narrow,
// first the words beside the symbols go to the tooltips, and then the least
// important commands move into a More menu at the trailing end, which exists
// only then, so every command stays reachable at any width (CMD-19, CMD-22).
// The window a toolbar is shown in answers its commands' keys for as long as
// it is shown there, after the commands added to the window itself.
type Toolbar struct {
	unison.Panel
	ui *UI
	// Label names the toolbar for a screen reader: "Statement".
	Label string
	// Words shows each command's words beside its symbol, where there is
	// room, on the Windows form (CMD-22). A command without a symbol always
	// shows its words.
	Words bool

	items    []*toolItem
	more     *IconButton
	overflow []*toolItem // the items in the More menu now
	words    bool        // whether the words show now
}

// toolItem is one command or menu on a toolbar, and the two buttons that
// can show it: its symbol alone, and its symbol with its words.
type toolItem struct {
	group ToolGroup
	cmd   *Command
	menu  []MenuItem // a menu of its own, for an item that opens one
	icon  *IconButton
	word  *Button
}

// NewToolbar returns an empty toolbar.
func NewToolbar(ui *UI, label string) *Toolbar {
	t := &Toolbar{ui: ui, Label: label, Words: true}
	t.Self = t
	t.more = NewIconButton(ui, "more", "More")
	t.more.OnClick = t.openMore
	OpensPopup(t.more)
	t.SetLayout(toolbarLayout{t})
	return t
}

// Add puts commands in a group, in order of importance.
func (t *Toolbar) Add(group ToolGroup, cmds ...*Command) *Toolbar {
	for _, c := range cmds {
		it := &toolItem{group: group, cmd: c}
		it.icon = NewCommandIconButton(t.ui, c)
		it.word = NewCommandButton(t.ui, c)
		it.word.Form = ButtonFlat
		t.items = append(t.items, it)
	}
	t.rebuild()
	return t
}

// AddMenu puts a button that opens a menu in a group: name and symbol are
// the button's, in the More menu a submenu's.
func (t *Toolbar) AddMenu(group ToolGroup, name, symbol string, items ...MenuItem) *Toolbar {
	c := &Command{Name: name, Symbol: symbol}
	it := &toolItem{group: group, cmd: c, menu: items}
	c.Run = func() { t.openItemMenu(it) }
	it.icon = NewCommandIconButton(t.ui, c)
	it.word = NewCommandButton(t.ui, c)
	OpensPopup(it.icon)
	OpensPopup(it.word)
	it.word.Form = ButtonFlat
	t.items = append(t.items, it)
	t.rebuild()
	return t
}

// Commands are the toolbar's commands, and those in its menus, each once.
func (t *Toolbar) Commands() []*Command {
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
	for _, it := range t.items {
		if it.menu != nil {
			walk(it.menu)
		} else if !slices.Contains(out, it.cmd) {
			out = append(out, it.cmd)
		}
	}
	return out
}

// InMore are the commands in the More menu now, the least important last.
func (t *Toolbar) InMore() []*Command {
	var out []*Command
	for _, it := range t.overflow {
		out = append(out, it.cmd)
	}
	return out
}

// rebuild makes the buttons the bar's children, in reading order: the
// leading group, the centre, the trailing group, then More.
func (t *Toolbar) rebuild() {
	t.RemoveAllChildren()
	for _, g := range []ToolGroup{ToolLeading, ToolCentre, ToolTrailing} {
		for _, it := range t.items {
			if it.group == g {
				t.AddChild(it.icon)
				t.AddChild(it.word)
			}
		}
	}
	t.AddChild(t.more)
	t.MarkForLayoutAndRedraw()
}

// showsWords reports whether an item would show its words, given whether
// words are on.
func (t *Toolbar) showsWords(it *toolItem, words bool) bool {
	return it.cmd.SymbolName() == "" || (words && t.Words && !t.ui.Mac())
}

// itemWidth is the width an item takes.
func (t *Toolbar) itemWidth(it *toolItem, words bool) float32 {
	var p geom.Size
	if t.showsWords(it, words) {
		_, p, _ = it.word.Sizes(geom.Size{})
	} else {
		_, p, _ = it.icon.Sizes(geom.Size{})
	}
	return p.Width
}

// arrange works out what fits in a width: whether the words show, and which
// items go to the More menu, the least important first.
func (t *Toolbar) arrange(width float32) (words bool, overflow []*toolItem) {
	m := t.ui.Interface
	gap, groupGap := float32(m.SpaceNear()), float32(m.SpaceLoose())
	need := func(shown []*toolItem, words bool, more bool) float32 {
		total := float32(0)
		groups := map[ToolGroup]bool{}
		for _, it := range shown {
			total += t.itemWidth(it, words) + gap
			groups[it.group] = true
		}
		if more {
			_, p, _ := t.more.Sizes(geom.Size{})
			total += p.Width + gap
		}
		return total + float32(max(0, len(groups)-1))*groupGap
	}
	shown := slices.Clone(t.items)
	for _, words := range []bool{true, false} {
		if need(shown, words, false) <= width {
			return words, nil
		}
	}
	for len(shown) > 0 {
		overflow = append([]*toolItem{shown[len(shown)-1]}, overflow...)
		shown = shown[:len(shown)-1]
		if need(shown, false, true) <= width {
			break
		}
	}
	return false, overflow
}

// openMore opens the More menu: the commands that did not fit, in order of
// importance.
func (t *Toolbar) openMore() {
	var items []MenuItem
	for _, it := range t.overflow {
		if it.menu != nil {
			items = append(items, MenuItem{Text: it.cmd.Name, Symbol: it.cmd.SymbolName(), Items: it.menu})
			continue
		}
		items = append(items, MenuItem{Command: it.cmd})
	}
	t.ui.ShowMenu(t.more, "More", items)
}

// openItemMenu opens an item's own menu under its button.
func (t *Toolbar) openItemMenu(it *toolItem) {
	var anchor unison.Paneler = it.icon
	if !it.word.Hidden {
		anchor = it.word
	}
	t.ui.ShowMenu(anchor, it.cmd.Name, it.menu)
}

type toolbarLayout struct{ t *Toolbar }

func (l toolbarLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	t, m := l.t, l.t.ui.Interface
	h := float32(m.ControlHeight())
	w := float32(0)
	for _, it := range t.items {
		w += t.itemWidth(it, true) + float32(m.SpaceNear())
	}
	w += float32(2 * m.SpaceLoose())
	_, more, _ := t.more.Sizes(geom.Size{})
	return geom.NewSize(more.Width, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l toolbarLayout) PerformLayout(target *unison.Panel) {
	t, m := l.t, l.t.ui.Interface
	r := target.ContentRect(false)
	words, overflow := t.arrange(r.Width)
	t.words, t.overflow = words, overflow
	gap := float32(m.SpaceNear())
	place := func(p unison.Paneler, x, w float32) {
		_, pref, _ := p.AsPanel().Sizes(geom.Size{})
		p.AsPanel().SetFrameRect(geom.NewRect(x, r.Y+(r.Height-pref.Height)/2, w, pref.Height))
	}
	groupItems := func(g ToolGroup) []*toolItem {
		var out []*toolItem
		for _, it := range t.items {
			it.icon.Hidden, it.word.Hidden = true, true
			if it.group == g && !slices.Contains(overflow, it) {
				out = append(out, it)
			}
		}
		return out
	}
	width := func(items []*toolItem) float32 {
		w := float32(0)
		for _, it := range items {
			w += t.itemWidth(it, words) + gap
		}
		return max(0, w-gap)
	}
	lay := func(items []*toolItem, x float32) {
		for _, it := range items {
			b := unison.Paneler(it.icon)
			if t.showsWords(it, words) {
				b = it.word
			}
			b.AsPanel().Hidden = false
			w := t.itemWidth(it, words)
			place(b, x, w)
			x += w + gap
		}
	}
	leading, centre, trailing := groupItems(ToolLeading), groupItems(ToolCentre), groupItems(ToolTrailing)
	// Every item is hidden by groupItems; those shown are revealed by lay.
	lay(leading, r.X)
	lay(centre, r.X+(r.Width-width(centre))/2)
	right := r.Right()
	t.more.Hidden = len(overflow) == 0
	if !t.more.Hidden {
		_, p, _ := t.more.Sizes(geom.Size{})
		right -= p.Width
		place(t.more, right, p.Width)
		right -= gap
	}
	lay(trailing, right-width(trailing))
}

// ProvideAccessibility names the strip as a toolbar.
func (t *Toolbar) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Toolbar
	}
	n.Name = t.Label
}
