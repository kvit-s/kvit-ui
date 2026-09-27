package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/check"
)

// MenuItem is one line of a menu opened with ShowMenu.
type MenuItem struct {
	// Text is what the line says.
	Text string
	// Key is the shortcut shown at the right of the line, which is how a menu
	// teaches that there is a faster way; the zero value shows none.
	Key unison.KeyBinding
	// Checked marks the entry in use.
	Checked bool
	// Disabled draws the line and does nothing when chosen.
	Disabled bool
	// Separator makes the line a divider instead, as goes before a
	// destructive entry at the bottom of a menu.
	Separator bool
	// OnSelect runs when the line is chosen.
	OnSelect func()
}

// menuIDs numbers the lines of menus opened with ShowMenu, clear of the ids
// unison and an application's menu bar use.
var menuIDs = 0x40000000

// ShowMenu opens a menu of items under anchor. It is unison's menu: drawn in
// the window on Windows and Linux in the Kvit colours and type, and a native
// menu on macOS, and in both it answers the arrow keys, Return and Escape.
func (u *UI) ShowMenu(anchor unison.Paneler, title string, items []MenuItem) {
	f := unison.DefaultMenuFactory()
	menuIDs++
	m := f.NewMenu(menuIDs, title, nil)
	for _, it := range items {
		if it.Separator {
			m.InsertSeparator(-1, true)
			continue
		}
		menuIDs++
		it := it
		mi := f.NewItem(menuIDs, it.Text, it.Key,
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
	p := anchor.AsPanel()
	m.Popup(p.RectToRoot(p.ContentRect(true)), 0)
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
