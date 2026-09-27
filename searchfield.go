package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
)

// SearchField is a field that filters something. It differs from Field in
// what it does rather than how it looks: a magnifier makes it findable
// without reading the placeholder, a clear button appears only when there is
// something to clear, and Escape clears rather than reverting, which is what
// a reader expects of a filter. It also says how many things the filter
// left, because filtering is the one interaction whose whole outcome happens
// somewhere else on the screen, and a screen reader user typing into a
// filter that says nothing gets no feedback at all.
type SearchField struct {
	Field
	// Matches is how many things the filter left; below zero says nothing,
	// for a filter whose result is not a countable list.
	Matches int
	// MatchedNoun is the word for one of them, "result" unless set.
	MatchedNoun string
	// MatchedNounPlural is the word for several; "" adds an s.
	MatchedNounPlural string

	clear *IconButton
}

// NewSearchField returns an empty filter field.
func NewSearchField(ui *UI) *SearchField {
	s := &SearchField{Matches: -1, MatchedNoun: "result"}
	s.ui = ui
	s.Self = s
	s.initField(ui, false)
	s.Placeholder = "Filter…"
	m := ui.Interface
	s.padLeft = func() float32 { return float32(m.ControlHeight()) }
	s.padRight = func() float32 {
		if s.Text() != "" {
			return float32(m.ControlHeight())
		}
		return float32(m.SpaceNear())
	}
	s.decorate = func(gc *unison.Canvas, box geom.Rect) {
		if g, ok := icons.Glyph("search"); ok {
			sz := float32(m.IconSizeSmall())
			drawGlyph(gc, ui, g, sz, ui.Theme.Tokens().TextFaint,
				geom.NewRect(box.X+float32(m.SpaceNear()), box.Y+(box.Height-sz)/2, sz, sz))
		}
	}
	// The clear button is a control's height square, the size of every
	// other control, over the field's right end.
	s.clear = NewIconButton(ui, "close", "Clear the filter")
	s.clear.OnClick = func() {
		s.SetText("")
		s.edit.RequestFocus()
	}
	s.AddChildAtIndex(s.clear, 0)
	s.place = func(box geom.Rect) {
		s.clear.Hidden = s.Text() == ""
		c := float32(m.ControlHeight())
		s.clear.SetFrameRect(geom.NewRect(box.Right()-c, box.Y+(box.Height-c)/2, c, c))
	}
	keyDown := s.edit.KeyDownCallback
	s.edit.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if key == unison.KeyEscape && s.Text() != "" {
			s.SetText("")
			return true
		}
		return keyDown(key, mods, repeat)
	}
	named := s.edit.Accessibility.Callback
	s.edit.Accessibility.Callback = func(n *accessibility.Node) {
		named(n)
		if n.Name == "" {
			n.Name = "Filter"
		}
		n.Description = s.MatchPhrase()
	}
	return s
}

// MatchPhrase is what a screen reader is told the filter left: "1 transaction",
// "250,000 transactions", or "" when the count is not given.
func (s *SearchField) MatchPhrase() string {
	return s.ui.CountPhrase(s.Matches, s.MatchedNoun, s.MatchedNounPlural)
}
