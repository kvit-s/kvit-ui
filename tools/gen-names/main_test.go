package main

import (
	"bytes"
	"os"
	"testing"
)

// names_gen.go matches Tokens and Interface; rerun the generator if not.
func TestNamesAreUpToDate(t *testing.T) {
	want, err := generate()
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("../../names_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Error("names_gen.go is out of date: go run ./tools/gen-names > names_gen.go")
	}
}
