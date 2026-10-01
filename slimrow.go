package kvitui

import (
	"strings"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// SlimRow is the list row a reader sees most: one line holding a name, what
// it is, one phrase about it and one figure at the right. The order is fixed,
// and the fixed order is the point: a reader scanning down the left edge
// finds the name in the same place on every screen, and scanning down the
// right edge finds the number.
type SlimRow struct {
	ListRow
	// Name is what the row is about.
	Name string
	// Kind says what sort of thing it is, where the name does not.
	Kind string
	// Phrase says one thing about its state, muted, because it is context
	// rather than identity.
	Phrase string
	// Figure is the number at the right, as the caller formatted it.
	Figure string
	// Unit follows the figure.
	Unit string
	// Measured is false for a value nobody measured, drawn as "—". A zero is
	// a measurement.
	Measured bool
	// Symbol is an optional meaning name drawn before the name.
	Symbol string

	line   *unison.Panel
	figure *Figure
}

// NewSlimRow returns a slim row for a name.
func NewSlimRow(ui *UI, name string) *SlimRow {
	s := &SlimRow{Name: name, Measured: true}
	s.Rule = true
	s.Self = s
	s.figure = NewFigure(ui, "", "")
	s.line = unison.NewPanel()
	s.line.AddChild(s.figure)
	s.line.SetLayout(syncing{Layout: slimLineLayout{s}, sync: s.sync})
	s.line.DrawCallback = s.drawLine
	s.line.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
	s.initListRow(ui, s.line)
	return s
}

// ProvideAccessibility describes the row as the list row does, except that a
// row nobody presses is a named group rather than a line of text, so its
// figure is still read out as a part of it.
func (s *SlimRow) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	s.ListRow.ProvideAccessibility(b)
	if n := b.Node(); !s.Interactive && n.Role == role.Label {
		n.Role = role.Group
	}
}

func (s *SlimRow) sync() {
	var parts []string
	for _, p := range []string{s.Name, s.Kind, s.Phrase} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	s.Label = strings.Join(parts, ", ")
	s.figure.Value, s.figure.Unit, s.figure.Measured = s.Figure, s.Unit, s.Measured
	s.figure.Hidden = s.Figure == "" && s.Measured
}

// slimLine is where each part of the line goes.
type slimLine struct {
	icon, name, kind, phrase, figure geom.Rect
}

// texts are the name, kind and phrase laid out at their natural widths.
func (s *SlimRow) texts() (name, kind, phrase *text.Layout) {
	ui := s.ui
	lay := func(str string, r TypeRole, c Ink) *text.Layout {
		return ui.Fonts.Layout([]text.Span{{Text: str, Style: ui.Chrome(ui.Size(r), text.Regular, c.Of(ui))}}, text.Options{})
	}
	return lay(s.Name, RoleBody, InkTextPrimary), lay(s.Kind, RoleSmall, InkTextFaint), lay(s.Phrase, RoleSmall, InkTextMuted)
}

// arrange lays the line out across a width: margins a near space in, parts
// a space apart, the name at most 45% of the row and the kind at most a
// quarter, the phrase taking what is left, and the figure at its natural
// width at the end.
func (s *SlimRow) arrange(b geom.Rect) slimLine {
	m := s.ui.Interface
	row := s.ContentRect(true).Width
	gap, edge := float32(m.Space()), float32(m.SpaceNear())
	name, kind, _ := s.texts()
	var out slimLine
	x, right := b.X+edge, b.Right()-edge
	if s.Symbol != "" {
		sz := float32(m.IconSizeSmall())
		out.icon = geom.NewRect(x, b.Y+(b.Height-sz)/2, sz, sz)
		x += sz + gap
	}
	if !s.figure.Hidden {
		_, fp, _ := s.figure.Sizes(geom.Size{})
		out.figure = geom.NewRect(right-fp.Width, b.Y+(b.Height-fp.Height)/2, fp.Width, fp.Height)
		right = out.figure.X - gap
	}
	nw, _ := name.Size()
	nw = min(nw, row*0.45)
	var kw float32
	if s.Kind != "" {
		kw, _ = kind.Size()
		kw = min(kw, row*0.25)
	}
	// What does not fit comes out of the phrase first, then the kind, then
	// the name.
	fixed := nw
	if s.Kind != "" {
		fixed += gap + kw
	}
	pw := right - x - fixed - gap
	if pw < 0 {
		short := -pw
		pw = 0
		cut := min(kw, short)
		kw -= cut
		nw = max(0, nw-(short-cut))
	}
	out.name = geom.NewRect(x, b.Y, nw, b.Height)
	x += nw + gap
	if s.Kind != "" {
		out.kind = geom.NewRect(x, b.Y, kw, b.Height)
		x += kw + gap
	}
	out.phrase = geom.NewRect(x, b.Y, pw, b.Height)
	return out
}

type slimLineLayout struct{ s *SlimRow }

func (l slimLineLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	s := l.s
	m := s.ui.Interface
	name, kind, phrase := s.texts()
	gap := float32(m.Space())
	w := 2 * float32(m.SpaceNear())
	nw, h := name.Size()
	w += nw + gap
	if s.Symbol != "" {
		w += float32(m.IconSizeSmall()) + gap
	}
	if s.Kind != "" {
		kw, _ := kind.Size()
		w += kw + gap
	}
	pw, _ := phrase.Size()
	w += pw
	if !s.figure.Hidden {
		_, fp, _ := s.figure.Sizes(geom.Size{})
		w += gap + fp.Width
	}
	return geom.NewSize(0, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l slimLineLayout) PerformLayout(target *unison.Panel) {
	if !l.s.figure.Hidden {
		l.s.figure.SetFrameRect(l.s.arrange(target.ContentRect(false)).figure)
	}
}

func (s *SlimRow) drawLine(gc *unison.Canvas, _ geom.Rect) {
	ui := s.ui
	b := s.line.ContentRect(false)
	at := s.arrange(b)
	if s.Symbol != "" {
		if g, ok := icons.Glyph(s.Symbol); ok {
			drawGlyph(gc, ui, g, at.icon.Width, ui.Theme.Tokens().TextMuted, at.icon)
		}
	}
	put := func(str string, r TypeRole, c Ink, box geom.Rect) {
		if str == "" || box.Width <= 0 {
			return
		}
		l := ui.Fonts.Layout([]text.Span{{Text: str, Style: ui.Chrome(ui.Size(r), text.Regular, c.Of(ui))}},
			text.Options{MaxWidth: box.Width, Elide: true})
		_, h := l.Size()
		l.Draw(gc, box.X, box.Y+(box.Height-h)/2)
	}
	put(s.Name, RoleBody, InkTextPrimary, at.name)
	put(s.Kind, RoleSmall, InkTextFaint, at.kind)
	put(s.Phrase, RoleSmall, InkTextMuted, at.phrase)
}
