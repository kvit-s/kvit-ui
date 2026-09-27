package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// The generated file matches what the Qt source produces today. While the Qt
// library exists beside this one, a failure means the Qt source changed:
// rerun the generator (its package comment gives the command) and port
// whatever the change means. Skipped where ~/kvit-ui has no Qt source.
func TestGeneratedFileMatchesTheQtSource(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip(err)
	}
	src := filepath.Join(home, "kvit-ui", "src/qml/iconcatalog.cpp")
	if _, err := os.Stat(src); err != nil {
		t.Skip("no Qt source at " + src)
	}
	want, err := generate(src)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join("..", "..", "icons/catalog_gen.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("icons/catalog_gen.go is out of date with %s", src)
	}
}
