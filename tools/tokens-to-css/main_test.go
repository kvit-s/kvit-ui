package main

import (
	"os"
	"path/filepath"
	"testing"
)

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
