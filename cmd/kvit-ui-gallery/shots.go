package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// shotSizes are the interface sizes the set is written at: the default, and
// the two ends of the range, as the Qt gallery's size variants are.
var shotSizes = []int{tokens.DefaultInterfaceSize, tokens.MinInterfaceSize, tokens.MaxInterfaceSize}

// shotName names a screenshot as the Qt gallery does: "<theme>-<page>.png"
// at the default size, "<theme>-<size>px-<page>.png" otherwise.
func shotName(theme string, size int, page string) string {
	if size == tokens.DefaultInterfaceSize {
		return fmt.Sprintf("%s-%s.png", theme, page)
	}
	return fmt.Sprintf("%s-%dpx-%s.png", theme, size, page)
}

// writeShots draws every page in every theme at each shot size, headlessly,
// into dir. With compare set to the Qt gallery's shot directory, it also
// writes compare/<name>, the Go image above the Qt one, for every name both
// have. It returns the files written.
func writeShots(dir, compare string) ([]string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		return nil, err
	}
	var g *gallery
	var startErr error
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: windowWidth, Height: 16000},
		unison.StartupFinishedCallback(func() {
			g, startErr = newGallery(ui, "Foundations", nil)
		}))
	if err != nil {
		return nil, err
	}
	defer screen.Stop()
	screen.Sync()
	if startErr != nil {
		return nil, startErr
	}
	names := make([]string, 0, len(pages))
	for n := range pages {
		names = append(names, n)
	}
	sort.Strings(names)
	var written []string
	for _, page := range names {
		for _, theme := range tokens.BuiltInThemes() {
			for _, size := range shotSizes {
				// Make the window tall enough to show the whole page without
				// scrolling, then capture it.
				screen.Do(func() {
					g.page = page
					ui.Theme.SetThemeID(theme)
					ui.Interface.SetFontSize(size)
					g.fit()
				})
				screen.Sync()
				screen.Do(g.fit)
				screen.Sync()
				img := screen.CaptureWindow(g.wnd)
				name := shotName(theme, size, page)
				if err := savePNG(filepath.Join(dir, name), img); err != nil {
					return written, err
				}
				written = append(written, name)
				if compare != "" {
					if err := writeComparison(dir, compare, name, img); err != nil {
						return written, err
					}
				}
			}
		}
	}
	for _, e := range screen.Errors() {
		return written, e
	}
	return written, nil
}

// fit sizes the window so the page shows in full: at least the gallery's
// usual 1440 × 960, taller when the page needs it.
func (g *gallery) fit() {
	m := g.ui.Interface
	width := float32(windowWidth - m.SidebarWidth())
	need := float32(m.HeaderHeight()+m.StatusBarHeight()) + g.pageHeight(width) + float32(m.Space())
	h := max(float32(windowHeight), need)
	g.wnd.SetContentRect(geom.NewRect(0, 0, windowWidth, h))
	g.wnd.Content().MarkForLayoutRecursively()
	g.wnd.ValidateLayout()
	g.wnd.MarkForRedraw()
}

func savePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// writeComparison stacks the Go image above the Qt one of the same name.
func writeComparison(dir, qtDir, name string, goImg image.Image) error {
	qf, err := os.Open(filepath.Join(qtDir, name))
	if err != nil {
		return nil // no Qt reference for this page
	}
	defer qf.Close()
	qtImg, err := png.Decode(qf)
	if err != nil {
		return err
	}
	gb, qb := goImg.Bounds(), qtImg.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, max(gb.Dx(), qb.Dx()), gb.Dy()+qb.Dy()+8))
	draw.Draw(out, image.Rect(0, 0, gb.Dx(), gb.Dy()), goImg, gb.Min, draw.Src)
	draw.Draw(out, image.Rect(0, gb.Dy()+8, qb.Dx(), gb.Dy()+8+qb.Dy()), qtImg, qb.Min, draw.Src)
	if err := os.MkdirAll(filepath.Join(dir, "compare"), 0o755); err != nil {
		return err
	}
	return savePNG(filepath.Join(dir, "compare", name), out)
}
