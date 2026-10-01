package icons_test

import (
	"bytes"
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/kvit-s/kvit-ui/icons"
)

func TestEveryMeaningNameResolves(t *testing.T) {
	names := icons.MeaningNames()
	if len(names) <= 40 {
		t.Fatalf("only %d meaning names", len(names))
	}
	for _, n := range names {
		if _, ok := icons.Glyph(n); !ok {
			t.Errorf("the meaning name %q resolves to nothing", n)
		}
	}
	// Two pairs share a drawing on purpose; any other pair is a mistake.
	allowed := map[string]bool{"agent and robot": true, "robot and agent": true, "rename and pencil": true, "pencil and rename": true}
	byGlyph := map[rune]string{}
	for _, n := range names {
		g, _ := icons.Glyph(n)
		if prev, ok := byGlyph[g]; ok && !allowed[prev+" and "+n] {
			t.Errorf("%s and %s draw the same glyph", prev, n)
		}
		byGlyph[g] = n
	}
}

func TestTheFontHasAGlyphForEveryMeaningName(t *testing.T) {
	face, err := font.ParseTTF(bytes.NewReader(icons.Font))
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range icons.MeaningNames() {
		g, _ := icons.Glyph(n)
		if _, ok := face.NominalGlyph(g); !ok {
			t.Errorf("the font has no glyph for %s (U+%04X)", n, g)
		}
	}
}

func TestPhosphorNamesResolveDirectly(t *testing.T) {
	a, _ := icons.Glyph("caret-right")
	b, _ := icons.Glyph("chevron-right")
	if a != b {
		t.Error("caret-right and chevron-right should be the same glyph")
	}
	if _, ok := icons.Glyph("acorn"); !ok {
		t.Error("a Phosphor name with no meaning should still resolve")
	}
}

func TestAnUnknownNameResolvesToNothing(t *testing.T) {
	for _, n := range []string{"no-such-symbol", "", "chevron_right", "Chevron-Right"} {
		if _, ok := icons.Glyph(n); ok {
			t.Errorf("%q resolved to a glyph", n)
		}
	}
}

func TestTheEditorsLiteralsAllHaveSomewhereToGo(t *testing.T) {
	for _, n := range []string{"chevron-right", "chevron-left", "close", "plus", "sidebar", "list", "search", "note", "tag",
		"messages-square", "message-square", "send", "chevron-down", "zoom-in", "zoom-out", "pencil", "check", "archive", "rotate-ccw"} {
		if _, ok := icons.Glyph(n); !ok {
			t.Errorf("%q has no glyph", n)
		}
	}
}

// The catalogue lists every glyph the font carries, so no symbol in the font
// is unreachable by name.
func TestTheCatalogueIsTheWholeFont(t *testing.T) {
	face, err := font.ParseTTF(bytes.NewReader(icons.Font))
	if err != nil {
		t.Fatal(err)
	}
	names := icons.GlyphNames()
	missing := 0
	for _, n := range names {
		g, _ := icons.Glyph(n)
		if _, ok := face.NominalGlyph(g); !ok {
			missing++
		}
	}
	if missing > 0 {
		t.Errorf("%d catalogued glyphs are not in the font", missing)
	}
	if len(names) < 1500 {
		t.Errorf("only %d glyphs catalogued; the font has about 1,530", len(names))
	}
}
