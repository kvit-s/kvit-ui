package main

import (
	"context"
	"fmt"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/unison"
)

// TestShots writes the whole screenshot set headlessly, one image per page,
// theme and shot size, and checks each is at least the window's size and is
// not blank. KVIT_SHOTS keeps the images in a directory of your choosing.
//
// It is also the check that every page opens in all four themes with no
// warning: anything logged at warning level or above while the set is
// written fails it, a panic in a component's drawing included, since unison
// logs those. The one warning expected is the icon page's own, which shows
// the mark an unknown symbol name draws.
func TestShots(t *testing.T) {
	dir := os.Getenv("KVIT_SHOTS")
	if dir == "" {
		dir = t.TempDir()
	}
	logged := newWarnings()
	previous := slog.Default()
	slog.SetDefault(slog.New(logged))
	written, err := writeShots(dir, os.Getenv("KVIT_QT_SHOTS"))
	slog.SetDefault(previous)
	if err != nil {
		t.Fatal(err)
	}
	expected := 0
	for _, w := range logged.lines() {
		if strings.Contains(w, "no symbol with this name") && strings.Contains(w, "not-a-symbol") {
			expected++
			continue
		}
		t.Errorf("warned while writing the set: %s", w)
	}
	if expected == 0 {
		t.Error("the icon page's warning about an unknown symbol was not seen, so warnings are not being caught")
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

// The sidebar's filter keeps the pages whose name or group has the words in
// it, and says how many it kept.
func TestTheSidebarFilters(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	var g *gallery
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: windowWidth, Height: windowHeight},
		unison.StartupFinishedCallback(func() { g, _ = newGallery(ui, "KvitIcon", nil) }))
	if err != nil {
		t.Fatal(err)
	}
	defer screen.Stop()
	screen.Sync()
	screen.Do(func() { g.search.SetText("row") })
	screen.Sync()
	screen.Do(func() {
		if _, ok := g.rows["KvitSlimRow"]; !ok {
			t.Error("filtering for row dropped KvitSlimRow")
		}
		if _, ok := g.rows["KvitLabel"]; ok {
			t.Error("filtering for row kept KvitLabel")
		}
		if g.search.Matches != len(g.rows) {
			t.Errorf("the filter says %d matches and shows %d", g.search.Matches, len(g.rows))
		}
	})
	// A page row opens its page.
	screen.Click(screen.PanelCenter(g.rows["KvitSlimRow"]))
	if g.page != "KvitSlimRow" {
		t.Errorf("clicking the KvitSlimRow row opened %q", g.page)
	}
}

// TestEveryControlSaysWhatItIs opens every page and checks that everything a
// reader can reach or press tells a screen reader what kind of thing it is
// and what it is called: a control with no name is read as "button" and
// nothing else.
func TestEveryControlSaysWhatItIs(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	ui.Theme.SetReducedMotion(true)
	var g *gallery
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: windowWidth, Height: windowHeight},
		unison.StartupFinishedCallback(func() { g, _ = newGallery(ui, "Foundations", nil) }))
	if err != nil {
		t.Fatal(err)
	}
	defer screen.Stop()
	screen.EnableAccessibility()
	screen.Sync()
	for _, page := range pageNames() {
		screen.Do(func() { g.setPage(page) })
		screen.Sync()
		for _, problem := range uitest.Unnamed(screen.AccessibilityTree(g.wnd.Window)) {
			t.Errorf("%s: %s", page, problem)
		}
		screen.KeyPress(unison.KeyEscape, 0)
	}
}

// warnings keeps what is logged at warning level or above, through the
// default logger and every logger derived from it.
type warnings struct {
	kept  *kept
	attrs []slog.Attr
}

type kept struct {
	mu    sync.Mutex
	lines []string
}

func newWarnings() *warnings { return &warnings{kept: &kept{}} }

func (w *warnings) lines() []string {
	w.kept.mu.Lock()
	defer w.kept.mu.Unlock()
	return append([]string(nil), w.kept.lines...)
}

func (w *warnings) Enabled(_ context.Context, l slog.Level) bool { return l >= slog.LevelWarn }

func (w *warnings) Handle(_ context.Context, r slog.Record) error {
	line := r.Message
	for _, a := range w.attrs {
		line += " " + a.String()
	}
	r.Attrs(func(a slog.Attr) bool {
		line += " " + a.String()
		return true
	})
	w.kept.mu.Lock()
	w.kept.lines = append(w.kept.lines, line)
	w.kept.mu.Unlock()
	return nil
}

func (w *warnings) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &warnings{kept: w.kept, attrs: append(append([]slog.Attr(nil), w.attrs...), attrs...)}
}

func (w *warnings) WithGroup(string) slog.Handler { return w }

// The skill's catalogue is what --catalog writes, so it describes what the
// repository contains.
func TestTheCatalogueIsCurrent(t *testing.T) {
	root, err := moduleRoot()
	if err != nil {
		t.Fatal(err)
	}
	want, err := renderCatalog(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, catalogPath))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("%s is not what --catalog writes; run go run ./cmd/kvit-ui-gallery --catalog %s", catalogPath, catalogPath)
	}
}
