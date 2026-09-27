// Package icons names every symbol the Kvit apps draw. The symbols are the
// glyphs of the Phosphor icon font, embedded here, and a call site asks for
// one by what it means ("chevron-right") rather than by what it looks like
// ("caret-right"), so a component does not change if the drawing moves to a
// different glyph. A Phosphor name also works, as the escape hatch for a
// symbol that has no meaning name yet.
package icons

import (
	_ "embed"
	"slices"
	"sort"
)

// Font is the Phosphor icon font (regular weight), MIT-licensed; see
// Phosphor-LICENSE.txt beside it.
//
//go:embed Phosphor.ttf
var Font []byte

// FontFamily is the family name the Phosphor font declares.
const FontFamily = "Phosphor"

// Glyph returns the character that draws name, a meaning name first and then
// a Phosphor name, and false when name is neither. A caller turns false into
// something visible rather than a blank box, so a typo cannot ship unseen.
func Glyph(name string) (rune, bool) {
	for _, m := range meanings {
		if m[0] == name {
			name = m[1]
			break
		}
	}
	r, ok := glyphs[name]
	return r, ok
}

// MeaningNames lists the names a call site should be written in, in the Qt
// library's order.
func MeaningNames() []string {
	out := make([]string, len(meanings))
	for i, m := range meanings {
		out[i] = m[0]
	}
	return out
}

// GlyphNames lists every Phosphor name the font carries, sorted.
func GlyphNames() []string {
	out := make([]string, 0, len(glyphs))
	for n := range glyphs {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// Draws returns the Phosphor name a meaning name draws, or name itself when
// it is not a meaning name.
func Draws(name string) string {
	if i := slices.IndexFunc(meanings, func(m [2]string) bool { return m[0] == name }); i >= 0 {
		return meanings[i][1]
	}
	return name
}
