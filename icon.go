package kvitui

import (
	"log/slog"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// Icon is a symbol asked for by what it means rather than by what it looks
// like (icons.MeaningNames lists the names), so the drawing can change
// without touching a call site. An unrecognised name draws a box with "?" in
// the danger colour and logs a warning naming it, so a typo shows on the
// first render instead of shipping as an empty space.
type Icon struct {
	unison.Panel
	ui *UI
	// Name is the symbol's meaning name.
	Name string
	// Ink is the colour; InkTextPrimary unless set.
	Ink Ink
	// Size is the glyph's size; SizeIconSize unless set.
	Size Measure
	// Label, when set, is what a screen reader is told the symbol means. An
	// icon without one is decoration and is hidden from screen readers.
	Label  string
	warned string
}

// NewIcon returns an icon at the interface's icon size.
func NewIcon(ui *UI, name string) *Icon {
	i := &Icon{ui: ui, Name: name}
	i.Self = i
	i.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		s := float32(i.size())
		size := geom.NewSize(s, s)
		return size, size, size
	})
	i.DrawCallback = i.draw
	return i
}

func (i *Icon) size() int {
	if i.Size != nil {
		return i.Size.Of(i.ui)
	}
	return i.ui.Interface.IconSize()
}

func (i *Icon) draw(gc *unison.Canvas, _ geom.Rect) {
	b := i.ContentRect(false)
	s := min(b.Width, b.Height, float32(i.size()))
	g, ok := icons.Glyph(i.Name)
	if !ok {
		i.drawUnknown(gc, b, s)
		return
	}
	ink := i.Ink
	if ink == nil {
		ink = InkTextPrimary
	}
	drawGlyph(gc, i.ui, g, s, ink.Of(i.ui), b)
}

// drawGlyph draws one symbol centred in r; the em box is square and the
// glyph fills it.
func drawGlyph(gc *unison.Canvas, ui *UI, g rune, size float32, c palette.Color, r geom.Rect) {
	st := text.Style{Family: icons.FontFamily, Size: size, Color: TextColor(c)}
	l := ui.Fonts.Layout([]text.Span{{Text: string(g), Style: st}}, text.Options{})
	w, h := l.Size()
	l.Draw(gc, r.X+(r.Width-w)/2, r.Y+(r.Height-h)/2)
}

func (i *Icon) drawUnknown(gc *unison.Canvas, b geom.Rect, s float32) {
	if i.warned != i.Name && i.Name != "" {
		slog.Warn("kvitui: no symbol with this name", "name", i.Name)
		i.warned = i.Name
	}
	t := i.ui.Theme.Tokens()
	m := i.ui.Interface
	box := geom.NewRect(b.X+(b.Width-s)/2, b.Y+(b.Height-s)/2, s, s)
	painterFor(gc, i.ui).outline(box, 0, float32(m.Hairline()), t.Danger)
	l := i.ui.Fonts.Layout([]text.Span{{Text: "?", Style: text.Style{Family: m.FontFamily(), Size: s * 0.7, Color: TextColor(t.Danger)}}}, text.Options{})
	w, h := l.Size()
	l.Draw(gc, box.X+(box.Width-w)/2, box.Y+(box.Height-h)/2)
}

// Recognized reports whether the name is a known symbol.
func (i *Icon) Recognized() bool {
	_, ok := icons.Glyph(i.Name)
	return ok
}

// ProvideAccessibility names a labelled icon as an image and hides a
// decorative one.
func (i *Icon) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if i.Label == "" {
		// Decoration: screen readers skip it.
		n.Ignored = true
		return
	}
	if n.Role == role.Auto {
		n.Role = role.Image
	}
	if n.Name == "" {
		n.Name = i.Label
	}
}
