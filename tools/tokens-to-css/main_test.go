package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// body is a stylesheet without its opening comment, which says where it came
// from and so differs between the Go and the Qt tool.
func body(css string) string {
	if i := strings.Index(css, "/* ═══"); i >= 0 {
		return css[i:]
	}
	return css
}

// The stylesheet in ux/ is what the generator writes, so nobody edits it by
// hand.
func TestTheStylesheetIsGenerated(t *testing.T) {
	want, err := render()
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join("..", "..", "ux", "tokens.css"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Error("ux/tokens.css is not what tools/tokens-to-css writes; run go run ./tools/tokens-to-css > ux/tokens.css")
	}
}

// Every value is the Qt library's: the stylesheet kvit-cash's copy of the Qt
// tool writes is the same after its opening comment.
func TestTheStylesheetIsTheQtOne(t *testing.T) {
	qt, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), "kvit-cash", "third-party", "kvit-ui", "ux", "tokens.css"))
	if err != nil {
		t.Skip("kvit-cash's copy of the Qt library is not here:", err)
	}
	want, err := render()
	if err != nil {
		t.Fatal(err)
	}
	if body(want) != body(string(qt)) {
		t.Error("the stylesheet differs from the one kvit-cash's copy of the Qt tool writes")
	}
}
