package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// shotSizes are the interface sizes the set is written at: the default, and
// the two ends of the range.
var shotSizes = []int{tokens.DefaultInterfaceSize, tokens.MinInterfaceSize, tokens.MaxInterfaceSize}

// shotName names a screenshot: "<theme>-<page>.png"
// at the default size, "<theme>-<size>px-<page>.png" otherwise, and
// "mac-<theme>-<page>.png" for the macOS form.
func shotName(form kvitui.Form, theme string, size int, page string) string {
	prefix := ""
	if form == kvitui.FormMac {
		prefix = "mac-"
	}
	if size == tokens.DefaultInterfaceSize {
		return fmt.Sprintf("%s%s-%s.png", prefix, theme, page)
	}
	return fmt.Sprintf("%s%s-%dpx-%s.png", prefix, theme, size, page)
}

// shot is one screenshot in the set: a form, a theme and an interface size.
type shot struct {
	form  kvitui.Form
	theme string
	size  int
}

// shots are the screenshots taken of every page: the Windows form (also
// Linux's) in every theme at each shot size, and the macOS form once, in the
// light theme at the default size, since the form changes words, order and
// keys rather than colours or sizes.
func shots() []shot {
	var out []shot
	for _, theme := range tokens.BuiltInThemes() {
		for _, size := range shotSizes {
			out = append(out, shot{kvitui.FormWindows, theme, size})
		}
	}
	return append(out, shot{kvitui.FormMac, tokens.Light, tokens.DefaultInterfaceSize})
}

// shotScreenHeight is the height of the headless screen the shots are drawn
// on, in pixels: a screen rather than anything in a layout, tall enough for
// the longest page at the largest interface size to be drawn whole.
const shotScreenHeight = 16000

// writeShots draws every page in every theme at each shot size, headlessly,
// into dir. It returns the files written.
func writeShots(dir string) ([]string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		return nil, err
	}
	// Motion is stilled for the run: a shot is taken as soon as the page is
	// laid out, and anything that eases, such as the sidebar's width after a
	// size change, would be caught partway.
	ui.Theme.SetReducedMotion(true)
	var g *gallery
	var startErr error
	widest := tokens.NewInterface()
	widest.SetFontSize(tokens.MaxInterfaceSize)
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: float32(widest.WidthDrawn()), Height: shotScreenHeight},
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
	var written []string
	for _, page := range pageNames() {
		for _, sh := range shots() {
			// Make the window tall enough to show the whole page without
			// scrolling, then capture it.
			screen.Do(func() {
				ui.SetForm(sh.form)
				// Built again for every shot, so a specimen that opens a
				// menu once it is shown opens it for each one.
				g.setPage(page)
				ui.Theme.SetThemeID(sh.theme)
				ui.Interface.SetFontSize(sh.size)
				// The status bar says which shot is being written.
				g.activity = fmt.Sprintf("Writing screenshots: %s / %s / %d px / %s", sh.form, sh.theme, sh.size, page)
				g.sync()
				g.fit()
			})
			screen.Sync()
			screen.Do(g.fit)
			screen.Sync()
			img := screen.CaptureWindow(g.wnd.Window)
			name := shotName(sh.form, sh.theme, sh.size, page)
			if err := savePNG(filepath.Join(dir, name), img); err != nil {
				return written, err
			}
			written = append(written, name)
			// Escape closes a menu a specimen opened.
			screen.KeyPress(unison.KeyEscape, 0)
		}
	}
	for _, e := range screen.Errors() {
		return written, e
	}
	return written, nil
}

// fit sizes the window so the page shows in full: as wide as the design
// width at the current interface size, which is the width the gallery's
// window takes and so where its pages wrap, and at least 960 tall, taller
// when the page needs it.
func (g *gallery) fit() {
	m := g.ui.Interface
	width := float32(m.WidthDrawn())
	need := float32(m.HeaderHeight()+m.StatusBarHeight()) + g.pageHeight(width-float32(m.SidebarWidth()))
	if bar := g.wnd.MenuBar(); bar != nil {
		_, p, _ := bar.Sizes(geom.Size{})
		need += p.Height
	}
	h := max(float32(windowHeight), need)
	g.wnd.SetContentRect(geom.NewRect(0, 0, width, h))
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
