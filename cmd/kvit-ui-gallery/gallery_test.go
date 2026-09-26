package main

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardwilkes/unison"
)

// TestGalleryOpens opens the gallery headlessly, checks that it draws, and
// saves what it drew.
func TestGalleryOpens(t *testing.T) {
	var wnd *unison.Window
	drawn := false
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 1100, Height: 800},
		unison.StartupFinishedCallback(func() {
			var err error
			if wnd, err = newGalleryWindow(func() { drawn = true }); err != nil {
				t.Error(err)
			}
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.Sync()
	var title string
	screen.Do(func() { title = wnd.Title() })
	if title != "kvit-ui gallery" {
		t.Errorf("window title %q", title)
	}
	if !drawn {
		t.Error("the window never drew")
	}
	f, err := os.Create(filepath.Join(t.TempDir(), "gallery.png"))
	if err != nil {
		t.Fatal(err)
	}
	if err = png.Encode(f, screen.Capture()); err != nil {
		t.Error(err)
	}
	if err = f.Close(); err != nil {
		t.Error(err)
	}
	for _, e := range screen.Errors() {
		t.Error(e)
	}
}
