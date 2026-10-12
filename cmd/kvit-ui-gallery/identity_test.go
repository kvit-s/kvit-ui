package main

import (
	"slices"
	"testing"

	"github.com/richardwilkes/unison"
)

// The window title gets the gallery's icon at every size tools/appicon
// draws for it, each one decoding to a square picture of that size.
func TestTheWindowHasTheGalleryIcon(t *testing.T) {
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: windowWidth, Height: windowHeight})
	if err != nil {
		t.Fatal(err)
	}
	defer screen.Stop()
	screen.Sync()
	var sizes []int
	screen.Do(func() {
		setTitleIcons()
		for _, img := range unison.DefaultTitleIcons {
			s := img.Size()
			if s.Width != s.Height {
				t.Errorf("an icon picture is %.0f x %.0f, not square", s.Width, s.Height)
			}
			sizes = append(sizes, int(s.Width))
		}
		unison.DefaultTitleIcons = nil
	})
	slices.Sort(sizes)
	if want := []int{16, 24, 32, 48, 64}; !slices.Equal(sizes, want) {
		t.Errorf("the title icons come in sizes %v, want %v", sizes, want)
	}
}
