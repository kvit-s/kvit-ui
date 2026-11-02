package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/role"
)

// MenuButton is a button that opens a menu of commands refining its own
// action, drawn with a down-arrow after its words so the menu is expected
// (CMD-17). It lists at least three commands: one or two are buttons of
// their own, and a view's main actions are on the screen rather than in a
// menu.
type MenuButton struct {
	Button
	// Items are the menu's lines.
	Items []MenuItem
}

// NewMenuButton returns an ordinary button that opens a menu.
func NewMenuButton(ui *UI, label string, items ...MenuItem) *MenuButton {
	b := &MenuButton{Items: items}
	b.Text = label
	b.Self = b
	b.trailing = "chevron-down"
	OpensPopup(b)
	b.initControl(ui, b.openMenu, unison.KeySpace, unison.KeyDown)
	b.pressKeys = func(ui *UI) []unison.KeyCode { return append(buttonKeys(ui), unison.KeyDown) }
	b.SetSizer(b.sizes)
	b.DrawCallback = b.draw
	b.tip = func() string { return b.Explanation }
	return b
}

func (b *MenuButton) openMenu() { b.ui.ShowMenu(b, b.label(), b.Items) }

// ProvideAccessibility describes the button as one that opens a menu.
func (b *MenuButton) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	b.Button.ProvideAccessibility(bl)
	n := bl.Node()
	n.Role = role.PopupButton
	n.Expandable = true
}

// PerformAccessibilityAction opens the menu for a screen reader.
func (b *MenuButton) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !b.Enabled() {
		return false
	}
	b.openMenu()
	return true
}

// SplitButton runs one of several variants of a command: its main part runs
// the variant it shows, and its arrow opens a menu of them all, the one the
// main part runs marked as the default (CMD-17). On the Windows form the
// main part runs the variant last chosen (decided Q-46); on the macOS form
// it runs the likeliest, the first, unless RememberLast is set for a
// variant people repeat. The two parts are drawn as one control, one
// outline around both and a line between them, as a ChoiceButton is.
//
// Made with NewToggleSplitButton it is a toggle split button instead, for a
// feature that is on or off and comes in styles, such as bullets: the main
// part turns the feature on and off and shows whether it is on, and the
// menu picks the style, the one in use marked, choosing one also turning the
// feature on (CMD-17).
type SplitButton struct {
	unison.Panel
	ui *UI
	// Variants are the commands the button offers, the likeliest first; a
	// toggle split button's styles.
	Variants []*Command
	// Current is the variant the main part runs, or the style in use.
	Current int
	// Toggle is a toggle split button's switch: a command whose Checked
	// says whether the feature is on and whose Run turns it on or off. nil
	// for a split button of variants.
	Toggle *Command
	// RememberLast makes the main part run the variant last chosen on the
	// macOS form too.
	RememberLast bool
	// MenuLabel names the arrow for a screen reader and its tooltip.
	MenuLabel string
	// Form is how both parts look: ordinary, primary, quiet or flat. A
	// toggle split button given the primary form is drawn ordinary, since
	// a filled button cannot show whether its feature is on.
	Form ButtonForm

	main  *Button
	arrow *IconButton
}

// NewSplitButton returns a split button offering variants.
func NewSplitButton(ui *UI, variants ...*Command) *SplitButton {
	s := &SplitButton{ui: ui, Variants: variants, MenuLabel: "More options"}
	s.Self = s
	s.main = NewButton(ui, "")
	s.arrow = NewIconButton(ui, "chevron-down", "")
	s.arrow.OnClick = s.openMenu
	s.main.fusedRight = true
	s.arrow.fusedTo = s.main
	OpensPopup(s.arrow)
	s.AddChild(s.main)
	s.AddChild(s.arrow)
	s.SetLayout(syncing{Layout: splitButtonLayout{s}, sync: s.sync})
	relayoutWhenMoved(s.AsPanel())
	// The pair draws before its parts, so each part draws the fields as they
	// are now, such as a Form the app has just set.
	s.DrawCallback = func(*unison.Canvas, geom.Rect) { s.sync() }
	return s
}

// sync copies the fields into the two parts.
func (s *SplitButton) sync() {
	s.Current = max(0, min(s.Current, len(s.Variants)-1))
	switch {
	case s.Toggle != nil:
		s.main.Command = s.Toggle
	case len(s.Variants) > 0:
		s.main.Command = s.Variants[s.Current]
	}
	s.arrow.Label = s.MenuLabel
	s.main.Form = s.Form
	if s.Toggle != nil && s.Form == ButtonPrimary {
		s.main.Form = ButtonOrdinary
	}
}

// NewToggleSplitButton returns a toggle split button: toggle turns the
// feature on and off, and styles are the ways it can look, the one in use
// at Current.
func NewToggleSplitButton(ui *UI, toggle *Command, styles ...*Command) *SplitButton {
	s := NewSplitButton(ui, styles...)
	s.Toggle = toggle
	s.MenuLabel = "Choose style"
	return s
}

// choose runs a variant chosen from the menu, and makes it the one the main
// part runs where the form says so. A style chosen is the one in use on
// both forms, and turns the feature on.
func (s *SplitButton) choose(i int) {
	if s.Toggle != nil || s.ui.form == FormWindows || s.RememberLast {
		s.Current = i
		s.MarkForLayoutAndRedraw()
	}
	s.ui.runCommand(s.Variants[i])
	if s.Toggle != nil && !s.Toggle.IsChecked() {
		s.ui.runCommand(s.Toggle)
	}
}

func (s *SplitButton) openMenu() {
	var items []MenuItem
	for i, c := range s.Variants {
		it := MenuItem{Command: c, Default: i == s.Current, OnSelect: func() { s.choose(i) }}
		if s.Toggle != nil {
			// Styles are options, one of which is in use, rather than
			// variants one of which is the default.
			it.Default, it.Radio, it.Checked = false, true, i == s.Current
		}
		items = append(items, it)
	}
	s.ui.ShowMenu(s, s.MenuLabel, items)
}

type splitButtonLayout struct{ s *SplitButton }

// The parts are whole pixels wide and placed on whole pixels of the window,
// as a ChoiceButton's are.
func (l splitButtonLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	_, mp, _ := l.s.main.Sizes(geom.Size{})
	_, ap, _ := l.s.arrow.Sizes(geom.Size{})
	size := geom.NewSize(wholeWidth(mp.Width)+wholeWidth(ap.Width), max(mp.Height, ap.Height))
	return size, size, size
}

func (l splitButtonLayout) PerformLayout(target *unison.Panel) {
	r := onWholePixels(target, target.ContentRect(false))
	_, ap, _ := l.s.arrow.Sizes(geom.Size{})
	w := wholeWidth(ap.Width)
	l.s.main.SetFrameRect(geom.NewRect(r.X, r.Y, r.Width-w, r.Height))
	l.s.arrow.SetFrameRect(geom.NewRect(r.Right()-w, r.Y, w, r.Height))
}

// ProvideAccessibility names the pair as a group.
func (s *SplitButton) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	switch {
	case s.Toggle != nil:
		n.Name = s.Toggle.Plain(s.ui)
	case len(s.Variants) > 0:
		n.Name = s.Variants[s.Current].Plain(s.ui)
	}
}
