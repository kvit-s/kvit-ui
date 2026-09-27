package main

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/kvit-s/kvit-ui/tokens"
)

// TestShots writes the whole screenshot set headlessly, one image per page,
// theme and shot size, and checks each is at least the window's size and is
// not blank. KVIT_SHOTS keeps the images in a directory of your choosing.
func TestShots(t *testing.T) {
	dir := os.Getenv("KVIT_SHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	written, err := writeShots(dir, os.Getenv("KVIT_QT_SHOTS"))
	if err != nil {
		t.Fatal(err)
	}
	if want := len(pages) * len(tokens.BuiltInThemes()) * len(shotSizes); len(written) != want {
		t.Fatalf("wrote %d screenshots, want %d", len(written), want)
	}
	for _, name := range written {
		f, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		b := img.Bounds()
		if b.Dx() < windowWidth || b.Dy() < windowHeight {
			t.Errorf("%s is %d × %d, smaller than the window", name, b.Dx(), b.Dy())
		}
		// Blank would be one colour throughout; count a sample of distinct ones.
		seen := map[uint32]bool{}
		for y := b.Min.Y; y < b.Max.Y; y += 7 {
			for x := b.Min.X; x < b.Max.X; x += 7 {
				r, g, bl, _ := img.At(x, y).RGBA()
				seen[r>>8<<16|g>>8<<8|bl>>8] = true
			}
		}
		if len(seen) < 50 {
			t.Errorf("%s has only %d colours: nothing was drawn", name, len(seen))
		}
	}
}
