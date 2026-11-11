package kvitui

import (
	"runtime"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// Select is a choice from a list too long to lay out: a font family, a
// currency, a folder. For a few mutually exclusive options worth seeing at
// once, Segmented is right: a three-item list hides two of the answers behind
// a press. The list opens as unison's menu, as the owner chose for menus, with
// the chosen option ticked; Up and Down change the choice without opening it.
// An option's Symbol or Picture is drawn before its words on its line in the
// list, and never in the closed select, which in every form draws the chosen
// option's words alone.
//
// Compact draws the same choice as a quiet button: the optional Symbol before
// the chosen option's words, no field outline. It is the control for a filter
// that states its value, such as a ledger's period reading All history, where
// the full-width field would take the room the list it filters needs. Its
// words are read as a value, never as a command.
//
// Plain draws the same choice as words in a row of a list, with a small arrow
// after them and a ground only under the pointer. It is the control for a
// value repeated down a list, where each row states one answer that reads with
// the row's name as a sentence and few rows get changed: there, Segmented
// repeats every answer on every row, outlined fields are heavier than what
// they say, and compact selects read as a column of buttons.
type Select struct {
	control
	// Label names the choice for a screen reader.
	Label string
	// Options are the choices.
	Options []Option
	// Current is the Value of the chosen option.
	Current string
	// Placeholder is what the select says while no option is chosen.
	Placeholder string
	// Symbol is an optional meaning name drawn before the chosen option's
	// words in the compact form, naming what the choice is about, such as
	// calendar for a period.
	Symbol string
	// Compact draws the select as a quiet button with the chosen option's
	// words, for a filter stating its value, rather than as a full-width
	// field.
	Compact bool
	// Plain draws the select as the chosen option's words and a small arrow,
	// with no ground or border until the pointer is on it, for a value
	// repeated down a list. It takes the width its layout gives it, its words
	// stand PlainSelectInset in from its left edge, and words too long for the
	// width are cut, with the arrow still shown. It wins over Compact.
	Plain bool
	// OnChoose runs after an option is chosen, with its value.
	OnChoose func(value string)

	jump       typeSelect // choosing an option by typing its words
	systemMenu bool       // the system's own menu is open, holding the program
	// fusedTo is the button whose right side the select is fused to, as a
	// ChoiceButton's choice is: drawn in the button's form with a
	// down-arrow, its left corners square, and the line between the two
	// along its left edge. nil for a select of its own.
	fusedTo *Button
}

// PlainSelectInset is how far a plain Select's words stand in from its left
// edge, which gives the ground under the pointer room around them. A list that
// names the column above its selects insets the name by this much, so the name
// stands over the words.
var PlainSelectInset = SizeSpaceNear

// NewSelect returns a choice of options, the first chosen.
func NewSelect(ui *UI, label string, options ...Option) *Select {
	s := &Select{Label: label, Options: options}
	s.Self = s
	OpensPopup(s)
	if len(options) > 0 {
		s.Current = options[0].Value
	}
	s.initControl(ui, s.open, unison.KeySpace, unison.KeyReturn, unison.KeyNumPadEnter)
	keyDown := s.KeyDownCallback
	s.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		switch key {
		case unison.KeyUp:
			return s.step(-1)
		case unison.KeyDown:
			if mods.OptionDown() {
				s.open()
				return true
			}
			return s.step(1)
		}
		return keyDown(key, mods, repeat)
	}
	// Typing the start of an option's words chooses it, open or closed
	// (KEY-9).
	s.RuneTypedCallback = func(ch rune) bool {
		if !s.jump.typing(ch) {
			return false
		}
		if i := s.jump.next(ch, s.index(), len(s.Options), func(i int) string { return s.Options[i].Label }); i >= 0 && s.Options[i].Value != s.Current {
			s.Choose(s.Options[i].Value)
		}
		return true
	}
	// Words cut short show whole where the pointer rests on them (PTR-12).
	s.tip = func() string {
		if s.Plain && s.plainWordsCut() && !s.menuOpen() {
			return s.shown()
		}
		return ""
	}
	s.SetSizer(s.sizes)
	s.DrawCallback = s.draw
	return s
}

// sizes is the select's size: a full-width field in the field form, its
// natural width as a quiet button in the compact form, the width its layout
// gives it in the plain form, read at layout time so it follows the interface
// size.
func (s *Select) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := s.ui.Interface
	h := float32(m.ControlHeight())
	if s.fusedTo != nil {
		size := geom.NewSize(s.fusedWidth(), h)
		return size, size, size
	}
	if s.Plain {
		pref := s.plainWidth()
		return geom.NewSize(min(float32(m.Px(60)), pref), h), geom.NewSize(pref, h), geom.NewSize(unison.DefaultMaxSize, h)
	}
	if !s.Compact {
		return geom.NewSize(float32(m.Px(60)), h), geom.NewSize(float32(m.Px(180)), h), geom.NewSize(unison.DefaultMaxSize, h)
	}
	w, _ := s.compactContent()
	size := geom.NewSize(w+2*float32(m.SpaceLoose()), max(h, float32(m.Px(22))))
	return size, size, size
}

// index is the position of the chosen option, or -1.
func (s *Select) index() int {
	for i, o := range s.Options {
		if o.Value == s.Current {
			return i
		}
	}
	return -1
}

// Choose makes the option with a value the chosen one.
func (s *Select) Choose(value string) {
	s.Current = value
	s.MarkForRedraw()
	if s.OnChoose != nil {
		s.OnChoose(value)
	}
}

func (s *Select) step(by int) bool {
	if len(s.Options) == 0 || !s.Enabled() {
		return false
	}
	i := max(0, min(len(s.Options)-1, s.index()+by))
	if s.Options[i].Value != s.Current {
		s.Choose(s.Options[i].Value)
	}
	return true
}

func (s *Select) open() {
	var items []MenuItem
	for _, o := range s.Options {
		items = append(items, MenuItem{Text: PlainMenuText(o.Label), Verbatim: true, Symbol: o.Symbol, Picture: o.Picture,
			Radio: true, Checked: o.Value == s.Current, OnSelect: func() { s.Choose(o.Value) }})
	}
	if s.Plain && runtime.GOOS == "darwin" {
		// The system's menu holds the program until it closes, so the ground
		// that marks the row has to be drawn before it opens.
		s.systemMenu = true
		s.MarkForRedraw()
		if w := s.Window(); w != nil {
			w.FlushDrawing()
		}
		defer func() {
			s.systemMenu = false
			s.MarkForRedraw()
		}()
	}
	s.ui.ShowMenu(s, s.Label, items)
}

// menuOpen reports whether the select's list is open now.
func (s *Select) menuOpen() bool {
	if s.systemMenu {
		return true
	}
	w := s.ui.windowOf(s)
	return w != nil && w.menuFrom(s)
}

// shown is the chosen option's label.
func (s *Select) shown() string {
	if i := s.index(); i >= 0 {
		return s.Options[i].Label
	}
	return s.Placeholder
}

func (s *Select) draw(gc *unison.Canvas, _ geom.Rect) {
	if s.fusedTo != nil {
		s.drawFused(gc)
		return
	}
	if s.Plain {
		s.drawPlain(gc)
		return
	}
	if s.Compact {
		s.drawCompact(gc)
		return
	}
	ui, t, m := s.ui, s.ui.Theme.Tokens(), s.ui.Interface
	p := painterFor(gc, ui)
	r := s.ContentRect(false)
	radius := float32(m.RadiusControl())
	ground, ink := t.PopupBackground, t.TextPrimary
	if !s.Enabled() {
		ground, ink = t.ChipBackground, t.TextDisabled
	}
	p.round(r, radius, ground)
	edge := t.BorderStrong
	if s.Focused() {
		edge = t.FocusRing
	}
	p.outline(r, radius, float32(m.Hairline()), edge)
	// The indicator at the small icon size: at the full size it overpowers
	// the words beside it.
	c := float32(m.IconSizeSmall())
	near := float32(m.SpaceNear())
	if g, ok := icons.Glyph("chevron-down"); ok {
		drawGlyph(gc, ui, g, c, t.TextMuted, geom.NewRect(r.Right()-near-c, r.Y+(r.Height-c)/2, c, c))
	}
	room := r.Width - near - c - float32(m.SpaceLoose()) - near
	l := ui.Fonts.Layout([]text.Span{{Text: s.shown(), Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ink)}},
		text.Options{MaxWidth: max(0, room), Elide: true})
	_, h := l.Size()
	l.Draw(gc, r.X+near, r.Y+(r.Height-h)/2)
}

// compactContent is the width and height of the compact form's symbol and
// words, read at layout and draw time so both agree.
func (s *Select) compactContent() (width, height float32) {
	m := s.ui.Interface
	ink := s.ui.Theme.Tokens().TextPrimary
	if !s.Enabled() {
		ink = s.ui.Theme.Tokens().TextDisabled
	}
	l := s.ui.Fonts.Layout([]text.Span{{Text: s.shown(), Style: s.ui.Chrome(s.ui.Size(RoleBody), text.Regular, ink)}},
		text.Options{})
	width, height = l.Size()
	if s.Symbol != "" {
		c := float32(m.IconSizeSmall())
		width += c + float32(m.SpaceNear())
		height = max(height, c)
	}
	return width, height
}

// drawCompact draws the select as a quiet button: the chip ground, the faint
// border rather than the field's strong one, the optional symbol before the
// chosen option's words. Hover and press answer as a button's do.
func (s *Select) drawCompact(gc *unison.Canvas) {
	ui, t, m := s.ui, s.ui.Theme.Tokens(), s.ui.Interface
	p := painterFor(gc, ui)
	r := s.ContentRect(false)
	radius := float32(m.RadiusControl())
	var ground palette.Color
	switch {
	case !s.Enabled():
		ground = t.ChipBackground
	case s.pressed:
		ground = t.SelectionTint
	case s.hovered:
		ground = t.HoverTint
	default:
		ground = t.ChipBackground
	}
	p.round(r, radius, ground)
	edge := t.Border
	if s.Focused() {
		edge = t.FocusRing
	}
	p.outline(r, radius, float32(m.Hairline()), edge)
	ink := t.TextPrimary
	if !s.Enabled() {
		ink = t.TextDisabled
	}
	x := r.X + float32(m.SpaceLoose())
	if s.Symbol != "" {
		c := float32(m.IconSizeSmall())
		if g, ok := icons.Glyph(s.Symbol); ok {
			drawGlyph(gc, ui, g, c, ink, geom.NewRect(x, r.Y+(r.Height-c)/2, c, c))
		}
		x += c + float32(m.SpaceNear())
	}
	l := ui.Fonts.Layout([]text.Span{{Text: s.shown(), Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ink)}},
		text.Options{})
	_, h := l.Size()
	l.Draw(gc, x, r.Y+(r.Height-h)/2)
}

// plainRoom is the plain form's room around its words: the inset before them,
// the gap between them and the arrow, the arrow's size and the space after it.
func (s *Select) plainRoom() (before, gap, arrow, after float32) {
	m := s.ui.Interface
	inset := float32(PlainSelectInset.Of(s.ui))
	return inset, float32(m.SpaceNear()), float32(m.IconSizeSmall()), inset
}

// plainLayout is words as the plain form draws them.
func (s *Select) plainLayout(words string, ink palette.Color, opt text.Options) *text.Layout {
	return s.ui.Fonts.Layout([]text.Span{{Text: words, Style: s.ui.Chrome(s.ui.Size(RoleBody), text.Regular, ink)}}, opt)
}

// plainWidth is the plain form's preferred width: the widest option's words,
// so the select keeps its width whatever is chosen, and the room around them.
func (s *Select) plainWidth() float32 {
	ink := s.ui.Theme.Tokens().TextPrimary
	widest := s.widestWords(func(words string) *text.Layout { return s.plainLayout(words, ink, text.Options{}) })
	before, gap, arrow, after := s.plainRoom()
	return before + wholeWidth(widest) + gap + arrow + after
}

// plainWordsRoom is the width the plain form has for its words now.
func (s *Select) plainWordsRoom() float32 {
	before, gap, arrow, after := s.plainRoom()
	return s.ContentRect(false).Width - before - gap - arrow - after
}

// plainWordsCut reports whether the chosen option's words are too long for the
// plain form's width now.
func (s *Select) plainWordsCut() bool {
	w, _ := s.plainLayout(s.shown(), s.ui.Theme.Tokens().TextPrimary, text.Options{}).Size()
	return w > s.plainWordsRoom()
}

// drawPlain draws the select as words in a row of a list: the chosen option's
// words and, directly after them, a small arrow that says a press opens a
// choice. There is no ground and no border at rest, so the control repeated
// down a column weighs no more than the words. Under the pointer the whole
// control tints, and while pressed or with its list open it takes the
// selection's tint, which marks the row the open list belongs to.
func (s *Select) drawPlain(gc *unison.Canvas) {
	ui, t, m := s.ui, s.ui.Theme.Tokens(), s.ui.Interface
	r := s.ContentRect(false)
	ink, mark := t.TextPrimary, t.TextMuted
	switch {
	case !s.Enabled():
		ink, mark = t.TextDisabled, t.TextDisabled
	case s.pressed || s.menuOpen():
		painterFor(gc, ui).round(r, float32(m.RadiusControl()), t.SelectionTint)
	case s.hovered:
		painterFor(gc, ui).round(r, float32(m.RadiusControl()), t.HoverTint)
	}
	before, gap, arrow, after := s.plainRoom()
	// At least 1, since a width of 0 would be no limit at all.
	l := s.plainLayout(s.shown(), ink, text.Options{MaxWidth: max(1, s.plainWordsRoom()), Elide: true})
	w, h := l.Size()
	x := r.X + before
	l.Draw(gc, x, r.Y+(r.Height-h)/2)
	if g, ok := icons.Glyph("chevron-down"); ok {
		// A whole number of pixels after the words, so the arrow draws the same
		// on every row, and never past the room kept for it.
		at := min(x+wholeWidth(w)+gap, r.Right()-after-arrow)
		drawGlyph(gc, ui, g, arrow, mark, geom.NewRect(at, r.Y+(r.Height-arrow)/2, arrow, arrow))
	}
}

// widestWords is the width of the widest words the select can show: each
// option's, and the placeholder's, as layout lays them out.
func (s *Select) widestWords(layout func(words string) *text.Layout) float32 {
	widest, _ := layout(s.Placeholder).Size()
	for _, o := range s.Options {
		w, _ := layout(o.Label).Size()
		widest = max(widest, w)
	}
	return widest
}

// valueLayout is the words of one option as the fused form draws them, in
// the button's chrome and ink.
func (s *Select) valueLayout(words string) *text.Layout {
	ui := s.ui
	ink := s.fusedTo.foreground(ui.Theme.Tokens())
	return ui.Fonts.Layout([]text.Span{{Text: words, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ink)}}, text.Options{})
}

// fusedRoom is the fused form's room around its words: the padding before
// them, the least gap between them and the down-arrow, the arrow's size and
// the padding after it. The padding is tighter than a button's, since the
// arrow only marks the part as a list and the value is what the part says.
func (s *Select) fusedRoom() (before, gap, arrow, after float32) {
	m := s.ui.Interface
	return float32(m.Space()), float32(m.SpaceSnug()), float32(fusedArrow.Of(s.ui)), float32(m.SpaceSnug())
}

// fusedWidth is the fused form's width: the widest option's words, so the
// pair keeps its width whatever is chosen, and the room around them.
func (s *Select) fusedWidth() float32 {
	before, gap, arrow, after := s.fusedRoom()
	return before + s.widestWords(s.valueLayout) + gap + arrow + after
}

// drawFused draws the select as the part fused to a button's right side,
// with the button's ground and outline, the chosen option's words, and a
// down-arrow at the right, in the button's ink. The words are centred in
// the room the widest option takes, as a button's are in its box, so a short
// value such as All does not sit against the line with a gap after it.
func (s *Select) drawFused(gc *unison.Canvas) {
	b := s.fusedTo
	r := s.ContentRect(false)
	b.drawFusedPart(gc, r, s.pressed, s.hovered)
	before, gap, arrow, after := s.fusedRoom()
	l := s.valueLayout(s.shown())
	w, h := l.Size()
	room := r.Width - before - gap - arrow - after
	l.Draw(gc, r.X+before+max(0, room-w)/2, r.Y+(r.Height-h)/2)
	if g, ok := icons.Glyph("chevron-down"); ok {
		drawGlyph(gc, s.ui, g, arrow, b.foreground(s.ui.Theme.Tokens()), geom.NewRect(r.Right()-after-arrow, r.Y+(r.Height-arrow)/2, arrow, arrow))
	}
}

// ProvideAccessibility describes the select as a combo box with its choice
// as its value.
func (s *Select) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.ComboBox
	}
	n.Name, n.Value = s.Label, s.shown()
	n.Expandable = true
	if s.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction opens the list for a screen reader.
func (s *Select) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !s.Enabled() {
		return false
	}
	s.open()
	return true
}
