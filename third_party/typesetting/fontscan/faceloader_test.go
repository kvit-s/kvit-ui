package fontscan

import (
	"testing"

	"github.com/go-text/typesetting/font"
)

// A loader set with SetFaceLoader is used for each system font the map
// chooses, once per font, and nil restores loading from disk.
func TestSetFaceLoader(t *testing.T) {
	fm := NewFontMap(nil)
	if err := fm.UseSystemFonts(t.TempDir()); err != nil {
		t.Skipf("no system fonts: %v", err)
	}
	var calls []Location
	fm.SetFaceLoader(func(loc Location) (*font.Face, error) {
		calls = append(calls, loc)
		fp := Footprint{Location: loc}
		return fp.loadFromDisk()
	})
	fm.SetQuery(Query{Families: []string{"sans-serif"}})
	if face := fm.ResolveFace('a'); face == nil {
		t.Fatal("no face resolved")
	}
	if len(calls) == 0 {
		t.Fatal("the loader was not used")
	}
	n := len(calls)
	fm.ResolveFace('b')
	if len(calls) != n {
		t.Errorf("a face already loaded was loaded again: %d calls, then %d", n, len(calls))
	}
	fm.SetFaceLoader(nil)
	fm.SetQuery(Query{Families: []string{"monospace"}})
	fm.ResolveFace('a')
	if len(calls) != n {
		t.Error("the loader was used after being removed")
	}
}
