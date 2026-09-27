package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/unison"
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
	if want := len(pageNames()) * len(tokens.BuiltInThemes()) * len(shotSizes); len(written) != want {
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
		// As wide as the design width at the shot's interface size, as the
		// Qt set is, and at least the gallery's usual height.
		m := tokens.NewInterface()
		for _, size := range shotSizes {
			if strings.Contains(name, fmt.Sprintf("-%dpx-", size)) {
				m.SetFontSize(size)
			}
		}
		b := img.Bounds()
		if b.Dx() != m.WidthDrawn() || b.Dy() < windowHeight {
			t.Errorf("%s is %d × %d, want %d wide and at least %d tall", name, b.Dx(), b.Dy(), m.WidthDrawn(), windowHeight)
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

// A page built at one interface size and shown at another is laid out as if
// it had been built at the second: every gap, padding and size follows the
// interface size, none keeps the value it had when the page was built.
func TestPagesFollowTheInterfaceSize(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	var g *gallery
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: windowWidth, Height: windowHeight},
		unison.StartupFinishedCallback(func() { g, _ = newGallery(ui, "Foundations", nil) }))
	if err != nil {
		t.Fatal(err)
	}
	defer screen.Stop()
	screen.Sync()
	width := float32(windowWidth - ui.Interface.SidebarWidth())
	for _, name := range pageNames() {
		var builtAt12, builtAt24 float32
		screen.Do(func() {
			ui.Interface.SetFontSize(12)
			g.setPage(name)
			builtAt12 = g.pageHeight(width)
			ui.Interface.SetFontSize(24)
			g.setPage(name)
			ui.Interface.SetFontSize(12)
			builtAt24 = g.pageHeight(width)
		})
		if builtAt12 != builtAt24 {
			t.Errorf("%s at 12 px is %.1f tall when built at 12 and %.1f when built at 24", name, builtAt12, builtAt24)
		}
	}
}

// A sample the gallery only shows is still run, as the Qt gallery test runs
// the window's sample, so it cannot drift from what compiles and works.
func TestTheSourceOnlySamplesRun(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	var ran int
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 1600, Height: 1000},
		unison.StartupFinishedCallback(func() {
			for _, e := range catalog {
				for _, s := range e.specimens {
					if s.sourceOnly {
						before := len(unison.Windows())
						s.build(ui)
						if len(unison.Windows()) <= before {
							t.Errorf("%s: %q opened no window", e.name, s.caption)
						}
						ran++
					}
				}
			}
		}))
	if err != nil {
		t.Fatal(err)
	}
	defer screen.Stop()
	screen.Sync()
	if ran == 0 {
		t.Error("no source-only sample was run")
	}
	for _, e := range screen.Errors() {
		t.Error(e)
	}
}
