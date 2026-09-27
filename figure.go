package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// Figure is a measured value, with three rules that are the reason it is a
// component rather than a label with a number in it:
//
//   - tabular digits, so a column of figures lines up and a value that
//     changes does not shift what is beside it;
//   - the unit smaller and muted, because the unit is the same on every row
//     and the number is not, and setting them alike makes the reader re-read
//     the unit to find the number;
//   - an em dash where nothing was measured, rather than a zero: a balance
//     nobody has computed is not a balance of zero, and drawing it as 0
//     states something false.
type Figure struct {
	unison.Panel
	ui *UI
	// Value is the number as the caller formatted it.
	Value string
	// Unit follows the value; "" for none.
	Unit string
	// Measured is false for a value nobody measured, which draws "—".
	Measured bool
	// Bounded marks an upper bound rather than a measurement, drawn "~12".
	Bounded bool
	// Role is the value's type role; RoleBody unless set.
	Role TypeRole
	// Ink is the value's colour; InkTextPrimary unless set.
	Ink Ink
}

// NewFigure returns a measured figure in the body role.
func NewFigure(ui *UI, value, unit string) *Figure {
	f := &Figure{ui: ui, Value: value, Unit: unit, Measured: true, Role: RoleBody}
	f.Self = f
	f.SetSizer(f.sizes)
	f.DrawCallback = f.draw
	return f
}

func (f *Figure) layouts() (value, unit *text.Layout) {
	ui, t := f.ui, f.ui.Theme.Tokens()
	ink := f.Ink
	if ink == nil {
		ink = InkTextPrimary
	}
	shown, c := f.Value, ink.Of(ui)
	switch {
	case !f.Measured:
		// One glyph wide at every size, read as an absence rather than a
		// minus sign, and in need of no translation.
		shown, c = "—", t.TextFaint
	case f.Bounded:
		shown = "~" + f.Value
	}
	st := ui.Chrome(ui.Size(f.Role), text.Regular, c)
	st.Tabular = true
	value = ui.Fonts.Layout([]text.Span{{Text: shown, Style: st}}, text.Options{})
	if f.Unit != "" && f.Measured {
		r := RoleSmall
		if f.Role == RoleCaption {
			r = RoleCaption
		}
		unit = ui.Fonts.Layout([]text.Span{{Text: f.Unit, Style: ui.Chrome(ui.Size(r), text.Regular, t.TextMuted)}}, text.Options{})
	}
	return value, unit
}

func (f *Figure) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	value, unit := f.layouts()
	w, h := value.Size()
	if unit != nil {
		uw, uh := unit.Size()
		w += float32(f.ui.Interface.SpaceSnug()) + uw
		h = max(h, uh)
	}
	size := geom.NewSize(w, h)
	return size, size, size
}

func (f *Figure) draw(gc *unison.Canvas, _ geom.Rect) {
	b := f.ContentRect(false)
	value, unit := f.layouts()
	w, h := value.Size()
	value.Draw(gc, b.X, b.Y+(b.Height-h)/2)
	if unit != nil {
		_, uh := unit.Size()
		unit.Draw(gc, b.X+w+float32(f.ui.Interface.SpaceSnug()), b.Y+(b.Height-uh)/2)
	}
}

// Phrase is what a screen reader is told: "12 h", "at most 12 h", or "not
// measured".
func (f *Figure) Phrase() string {
	switch {
	case !f.Measured:
		return "not measured"
	case f.Bounded:
		return joinWords("at most", f.Value, f.Unit)
	}
	return joinWords(f.Value, f.Unit)
}

// ProvideAccessibility reads the figure as text.
func (f *Figure) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	n.Name = f.Phrase()
}

// joinWords joins the non-empty words with spaces.
func joinWords(words ...string) string {
	out := ""
	for _, w := range words {
		if w == "" {
			continue
		}
		if out != "" {
			out += " "
		}
		out += w
	}
	return out
}
