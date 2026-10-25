package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// SettingsPane is one pane of an app's settings: its name, the symbol drawn
// over the name, and what it holds, a column of SettingRows.
type SettingsPane struct {
	// Name is a noun for what the pane covers, "General" or "Accounts":
	// the pane's button, and the window's title while it is shown.
	Name string
	// Symbol is the meaning name of the symbol over the name.
	Symbol string
	// Content is what the pane holds.
	Content unison.Paneler
}

// SettingsPanes are an app's settings in panes, as the macOS form has them
// (WIN-16): the panes' buttons in a row across the top, each a symbol over
// its name, and the pane shown under them. Each change takes effect as it is
// made, so there is no OK, Cancel or Apply (WIN-23). A SettingsWindow holds
// it; on the Windows form, which Linux takes, the same panes make a page in
// the app's window instead (NewSettingsPage).
type SettingsPanes struct {
	unison.Panel
	ui *UI
	// Panes are the panes, in the order their buttons show.
	Panes []SettingsPane
	// Current is the pane shown.
	Current int
	// OnShow runs when the reader shows another pane, so the app can open
	// the settings on it the next time.
	OnShow func(pane int)

	shownHook func() // what holds the panes does when another is shown
	row       *unison.Panel
	buttons   []*settingsPaneButton
	rule      *Divider
	shown     unison.Paneler
}

// NewSettingsPanes returns the panes with the first one shown.
func NewSettingsPanes(ui *UI, panes ...SettingsPane) *SettingsPanes {
	s := &SettingsPanes{ui: ui, Panes: panes}
	s.Self = s
	s.row = unison.NewPanel()
	s.rule = NewDivider(ui)
	s.SetLayout(syncing{Layout: settingsPanesLayout{s}, sync: s.sync})
	s.sync()
	return s
}

// Show shows a pane.
func (s *SettingsPanes) Show(pane int) {
	if pane < 0 || pane >= len(s.Panes) || pane == s.Current {
		return
	}
	s.Current = pane
	s.sync()
	s.MarkForLayoutAndRedraw()
	if s.shownHook != nil {
		s.shownHook()
	}
	if s.OnShow != nil {
		s.OnShow(pane)
	}
}

// Title is the name of the pane shown, in the form's case, as the window
// holding the panes is titled.
func (s *SettingsPanes) Title() string {
	if s.Current < 0 || s.Current >= len(s.Panes) {
		return ""
	}
	return s.ui.CommandCase(s.Panes[s.Current].Name)
}

// sync makes a button for each pane and puts the shown pane under the row.
// Only the pane is swapped, so the button the keyboard is on stays in the
// window and keeps it.
func (s *SettingsPanes) sync() {
	s.Current = max(0, min(s.Current, len(s.Panes)-1))
	if len(s.buttons) != len(s.Panes) {
		for len(s.buttons) < len(s.Panes) {
			s.buttons = append(s.buttons, newSettingsPaneButton(s, len(s.buttons)))
		}
		s.buttons = s.buttons[:len(s.Panes)]
		kids := make([]unison.Paneler, len(s.buttons))
		for i, b := range s.buttons {
			kids[i] = b
		}
		s.row.SetLayout(&spaced{FlexLayout: unison.FlexLayout{Columns: len(kids), EqualColumns: true}, ui: s.ui, gap: SizeSpaceTight})
		showOnly(s.row, kids...)
	}
	if s.row.Parent() != s.AsPanel() {
		showOnly(s.AsPanel(), s.row, s.rule)
	}
	var want unison.Paneler
	if len(s.Panes) > 0 {
		want = s.Panes[s.Current].Content
	}
	if s.shown != nil && (want == nil || s.shown.AsPanel() != want.AsPanel()) {
		s.shown.AsPanel().RemoveFromParent()
		s.shown = nil
	}
	if want != nil && want.AsPanel().Parent() != s.AsPanel() {
		s.AddChild(want)
	}
	s.shown = want
}

type settingsPanesLayout struct{ s *SettingsPanes }

// arrange places the row of buttons centred at the top, the rule under it
// across the whole width, and the pane under the rule inside a view margin,
// and returns the height they take.
func (l settingsPanesLayout) arrange(r geom.Rect, set bool) float32 {
	s, m := l.s, l.s.ui.Interface
	margin, gap := float32(m.ViewMargin()), float32(m.Space())
	_, rp, _ := s.row.Sizes(geom.Size{})
	y := r.Y + float32(m.SpaceNear())
	if set {
		s.row.SetFrameRect(geom.NewRect(r.X+(r.Width-rp.Width)/2, y, rp.Width, rp.Height))
	}
	y += rp.Height + float32(m.SpaceNear())
	_, dp, _ := s.rule.Sizes(geom.Size{})
	if set {
		s.rule.SetFrameRect(geom.NewRect(r.X, y, r.Width, dp.Height))
	}
	y += dp.Height
	if s.shown != nil {
		inner := max(0, r.Width-2*margin)
		_, cp, _ := s.shown.AsPanel().Sizes(geom.NewSize(inner, 0))
		y += gap
		if set {
			s.shown.AsPanel().SetFrameRect(geom.NewRect(r.X+margin, y, inner, cp.Height))
		}
		y += cp.Height + margin
	}
	return y - r.Y
}

func (l settingsPanesLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	s, m := l.s, l.s.ui.Interface
	_, rp, _ := s.row.Sizes(geom.Size{})
	w := max(float32(m.Px(560)), rp.Width+2*float32(m.ViewMargin()))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := l.arrange(geom.NewRect(0, 0, w, 0), false)
	return geom.NewSize(rp.Width, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l settingsPanesLayout) PerformLayout(target *unison.Panel) {
	l.arrange(target.ContentRect(true), true)
}

// ProvideAccessibility leaves the holder out: the row is a tab list, and the
// pane shown is what it controls.
func (s *SettingsPanes) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	b.Node().Role = role.None
}

// settingsPaneButton is one pane's button: its symbol over its name, the
// pane shown marked by a tinted ground, the accent on its symbol and its
// name in bold, so the mark does not rest on colour alone. The row is one
// stop in the Tab order, and Left and Right show the pane before or after
// (KEY-9).
type settingsPaneButton struct {
	control
	s    *SettingsPanes
	pane int
}

func newSettingsPaneButton(s *SettingsPanes, pane int) *settingsPaneButton {
	b := &settingsPaneButton{s: s, pane: pane}
	b.Self = b
	b.initControl(s.ui, func() { s.Show(b.pane) }, unison.KeySpace, unison.KeyReturn)
	keys := b.KeyDownCallback
	b.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if mods&mod.NonSticky == 0 {
			if next, moved := roveStep(siblings[*settingsPaneButton](b.AsPanel()), b, key, false, false); moved {
				if next != nil {
					next.SetFocusable(true)
					next.Focus()
					next.fire()
				}
				return true
			}
		}
		return keys(key, mods, repeat)
	}
	b.SetSizer(b.sizes)
	b.DrawCallback = b.draw
	return b
}

func (b *settingsPaneButton) shown() bool { return b.pane == b.s.Current }

func (b *settingsPaneButton) name() *text.Layout {
	ui, tk := b.ui, b.ui.Theme.Tokens()
	weight, ink := text.Regular, tk.TextSecondary
	if b.shown() {
		weight, ink = text.Bold, tk.TextPrimary
	}
	words := ""
	if b.pane < len(b.s.Panes) {
		words = ui.CommandCase(b.s.Panes[b.pane].Name)
	}
	return ui.Fonts.Layout([]text.Span{{Text: words, Style: ui.Chrome(ui.Size(RoleSmall), weight, ink)}}, text.Options{})
}

func (b *settingsPaneButton) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	oneStop(siblings[*settingsPaneButton](b.AsPanel()), func(o *settingsPaneButton) bool { return o.shown() })
	m := b.ui.Interface
	nw, nh := b.name().Size()
	pad := float32(m.SpaceNear())
	w := max(nw, float32(m.IconSize())) + 2*float32(m.Space())
	h := pad + float32(m.Px(24)) + float32(m.SpaceTight()) + nh + pad
	size := geom.NewSize(max(w, float32(m.Px(72))), h)
	return size, size, size
}

func (b *settingsPaneButton) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := b.ui, b.ui.Theme.Tokens(), b.ui.Interface
	p := painterFor(gc, ui)
	r := b.ContentRect(false)
	radius := float32(m.RadiusControl())
	switch {
	case b.shown():
		p.round(r, radius, tk.SelectionTint)
	case b.hovered:
		p.round(r, radius, tk.HoverTint)
	}
	pad := float32(m.SpaceNear())
	ink := tk.TextSecondary
	if b.shown() {
		ink = tk.Accent
	}
	if b.pane < len(b.s.Panes) {
		if g, ok := icons.Glyph(b.s.Panes[b.pane].Symbol); ok {
			drawGlyph(gc, ui, g, float32(m.Px(24)), ink, geom.NewRect(r.X, r.Y+pad, r.Width, float32(m.Px(24))))
		}
	}
	name := b.name()
	nw, _ := name.Size()
	name.Draw(gc, r.X+(r.Width-nw)/2, r.Y+pad+float32(m.Px(24))+float32(m.SpaceTight()))
}

// ProvideAccessibility announces the button as a tab, selected while its
// pane is shown.
func (b *settingsPaneButton) ProvideAccessibility(bl *unison.AccessibilityBuilder) {
	n := bl.Node()
	if n.Role == role.Auto {
		n.Role = role.Tab
	}
	if b.pane < len(b.s.Panes) {
		n.Name = b.ui.CommandCase(b.s.Panes[b.pane].Name)
	}
	n.Selectable = true
	n.Selected = b.shown()
	if b.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// NewSettingsPage returns an app's settings as the Windows form has them,
// which Linux takes (WIN-16): a page for the app's window, filling it as one
// scrolling column with each pane a section under its name, in order, and
// About last, at most a readable width wide. Each change takes effect as it
// is made, so there is no OK, Cancel or Apply (WIN-23). The same panes make
// the SettingsWindow of the macOS form, so an app writes its settings once;
// a pane's content is in one place at a time, the page or the window.
func NewSettingsPage(ui *UI, about *About, panes ...SettingsPane) *Region {
	var sections []unison.Paneler
	for _, pane := range panes {
		parts := []unison.Paneler{FullWidth(NewSectionHeading(ui, pane.Name))}
		if pane.Content != nil {
			parts = append(parts, FullWidth(pane.Content))
		}
		section := Column(ui, SizeSpaceNear, parts...)
		section.Accessibility.Role, section.Accessibility.Name = role.Group, pane.Name
		sections = append(sections, FullWidth(section))
	}
	if about != nil {
		sections = append(sections, FullWidth(NewAboutSection(ui, *about)))
	}
	column := Column(ui, SizeViewMargin, sections...)
	page := unison.NewPanel()
	page.AddChild(column)
	page.SetLayout(settingsPageLayout{ui: ui, column: column})
	return NewRegion(ui, page)
}

// settingsPageLayout keeps the settings column at the leading edge, as wide
// as the page up to a readable width, so a setting's control stays near its
// name in a wide window.
type settingsPageLayout struct {
	ui     *UI
	column *unison.Panel
}

func (l settingsPageLayout) width(available float32) float32 {
	most := float32(l.ui.Interface.Px(720))
	if available <= 0 {
		return most
	}
	return min(available, most)
}

func (l settingsPageLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	w := l.width(hint.Width)
	cmin, cp, _ := l.column.Sizes(geom.NewSize(w, 0))
	return geom.NewSize(cmin.Width, cp.Height), geom.NewSize(w, cp.Height), geom.NewSize(unison.DefaultMaxSize, cp.Height)
}

func (l settingsPageLayout) PerformLayout(target *unison.Panel) {
	r := target.ContentRect(false)
	w := l.width(r.Width)
	_, cp, _ := l.column.Sizes(geom.NewSize(w, 0))
	l.column.SetFrameRect(geom.NewRect(r.X, r.Y, w, cp.Height))
}

// SettingsWindow is the app's settings on the macOS form (WIN-16): a window
// of its own, opened by Settings… in the app menu, holding SettingsPanes.
// It is titled with the pane shown, takes the height of the pane shown, and
// cannot be resized; Command-W closes it. unison cannot take the minimise
// button off a window, which a Mac's settings window does not have. On the
// Windows form, which Linux takes, settings are a page in the app's window
// made of the same panes (Page), and Command runs the app's own way of
// showing it.
type SettingsWindow struct {
	ui *UI
	// Panes are what the window holds.
	Panes *SettingsPanes

	window *Window
	// top is where the window's content was last left, its top-left corner
	// on the screen, so it opens there again; nil until it has been open.
	top *geom.Point
}

// NewSettingsWindow returns the settings window for panes, not yet open.
func NewSettingsWindow(ui *UI, panes ...SettingsPane) *SettingsWindow {
	s := &SettingsWindow{ui: ui, Panes: NewSettingsPanes(ui, panes...)}
	s.Panes.shownHook = s.fit
	return s
}

// Open opens the window, or brings it to the front when it is open, with
// the pane last shown.
func (s *SettingsWindow) Open() error {
	if s.window != nil {
		s.window.ToFront()
		return nil
	}
	w, err := newWindow(s.ui, s.Panes.Title(), unison.NotResizableWindowOption())
	if err != nil {
		return err
	}
	s.window = w
	w.SetBody(s.Panes)
	w.AddCommands(&Command{Std: StdClose, Run: s.Close})
	willClose := w.WillCloseCallback
	w.WillCloseCallback = func() {
		willClose()
		at := w.ContentRect().Point
		s.top = &at
		if s.window == w {
			s.window = nil
		}
	}
	w.MinMaxContentSizeCallback = func() (minimum, maximum geom.Size) {
		size := s.contentSize()
		return size, size
	}
	s.place()
	w.ToFront()
	// The keyboard starts on the pane shown, as unison starts it on the
	// first control of a window it activates; without a focus, Tab would
	// go nowhere.
	unison.InvokeTask(func() {
		if s.window == w && w.CurrentFocus() == nil {
			w.FocusNext()
		}
	})
	return nil
}

// Close closes the window, if it is open.
func (s *SettingsWindow) Close() {
	if s.window != nil {
		s.window.Dispose()
		s.window = nil
	}
}

// Window is the window while it is open, or nil.
func (s *SettingsWindow) Window() *Window { return s.window }

// Page is the Windows form's settings page made of the window's panes,
// ending with about where it is given: what showPage, given to Command, puts
// in the app's window.
func (s *SettingsWindow) Page(about *About) *Region {
	return NewSettingsPage(s.ui, about, s.Panes.Panes...)
}

// Command is the Settings command: on the macOS form it opens the window;
// on the Windows form it runs showPage, which shows the settings page in
// the app's window.
func (s *SettingsWindow) Command(showPage func()) *Command {
	return &Command{Std: StdSettings, Run: func() {
		if s.ui.form == FormMac {
			_ = s.Open()
		} else if showPage != nil {
			showPage()
		}
	}}
}

// contentSize is the size the pane shown wants.
func (s *SettingsWindow) contentSize() geom.Size {
	_, p, _ := s.Panes.Sizes(geom.Size{})
	_, h, _ := s.Panes.Sizes(geom.NewSize(p.Width, 0))
	return geom.NewSize(p.Width, h.Height)
}

// place sizes the window for the pane shown and puts it where it was last
// left, or the first time centred across the screen a fifth of the way
// down, where a Mac opens a settings window, and then fits it on the screen.
func (s *SettingsWindow) place() {
	w := s.window
	s.fit()
	size := s.contentSize()
	at := geom.Point{}
	if s.top != nil {
		at = *s.top
	} else if d := unison.PrimaryDisplay(); d != nil {
		u := d.Usable
		at = geom.NewPoint(u.X+(u.Width-size.Width)/2, u.Y+u.Height/5)
	}
	w.SetContentRect(geom.Rect{Point: at, Size: size})
	w.FitToScreen()
}

// fit titles the window with the pane shown and gives it the pane's height,
// keeping its top edge where it is, as a Mac's settings window grows and
// shrinks downwards.
func (s *SettingsWindow) fit() {
	w := s.window
	if w == nil {
		return
	}
	w.SetTitle(s.Panes.Title())
	size := s.contentSize()
	r := w.ContentRect()
	w.SetContentRect(geom.NewRect(r.X, r.Y, size.Width, size.Height))
	w.Content().MarkForLayoutAndRedraw()
}
