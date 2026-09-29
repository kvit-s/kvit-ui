// Command kvit-ui-gallery shows every kvit-ui component with its states, in
// each theme and at each interface size, as the library's gallery does.
// Pages not built yet are listed greyed.
//
//	kvit-ui-gallery                                   browse
//	kvit-ui-gallery --page Foundations --theme dark   one page, one theme
//	kvit-ui-gallery --interface-size 16               at a chosen size
//	kvit-ui-gallery --shots DIR                       write the screenshot set and exit
//	kvit-ui-gallery --smoke 3s                        close after 3 s, printing when the first frame was drawn
//
// In the window, Ctrl+1 to Ctrl+4 choose a theme and Ctrl+plus, Ctrl+minus
// and Ctrl+0 change the interface size.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"slices"
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/unison"
)

var started = time.Now()

func main() {
	page := flag.String("page", "Foundations", "the page to open")
	theme := flag.String("theme", "", "light, dark, sepia, highContrast or system (default: the saved choice)")
	size := flag.Int("interface-size", 0, "the interface size in pixels, 10 to 24 (default: the saved choice)")
	shots := flag.String("shots", "", "write the screenshot set into this directory and exit")
	smoke := flag.Duration("smoke", 0, "close the window after this long, printing when its first frame was drawn")
	heapProfile := flag.String("heap-profile", "", "with --smoke: write a heap profile to this file halfway through")
	catalogFile := flag.String("catalog", "", "write the vocabulary skill's component catalogue to this file and exit, from a checkout of the library ("+catalogPath+" there)")
	flag.Parse()

	if *catalogFile != "" {
		if err := writeCatalog(*catalogFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	if *shots != "" {
		written, err := writeShots(*shots)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("wrote %d screenshots to %s\n", len(written), *shots)
		return
	}

	ui, err := kvitui.New(kvitui.Options{SettingsPath: kvitui.DefaultSettingsPath("kvit-ui-gallery")})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *theme != "" {
		if !slices.Contains(tokens.AvailableThemes(), *theme) {
			fmt.Fprintf(os.Stderr, "unknown theme %q\n", *theme)
			os.Exit(2)
		}
		ui.Theme.SetThemeID(*theme)
	}
	if *size != 0 {
		ui.Interface.SetFontSize(*size)
	}
	unison.Start(
		unison.ThemeChangedCallback(ui.Appearance.Refresh),
		unison.StartupFinishedCallback(func() {
			var g *gallery
			var err error
			g, err = newGallery(ui, *page, func() {
				if *smoke > 0 {
					r := g.wnd.ContentRect()
					fmt.Printf("first frame after %d ms, %s/%s, window %.0f x %.0f at scale %.2f\n",
						time.Since(started).Milliseconds(), runtime.GOOS, runtime.GOARCH, r.Width, r.Height, g.wnd.BackingScale().X)
					unison.InvokeTaskAfter(func() {
						// What Go itself holds, to set against the process's memory.
						runtime.GC()
						var ms runtime.MemStats
						runtime.ReadMemStats(&ms)
						fmt.Printf("Go heap in use %d MB, Go memory from the system %d MB\n", ms.HeapInuse>>20, ms.Sys>>20)
						if *heapProfile != "" {
							if f, err := os.Create(*heapProfile); err == nil {
								_ = pprof.WriteHeapProfile(f)
								f.Close()
							}
						}
					}, *smoke/2)
					unison.InvokeTaskAfter(g.wnd.Dispose, *smoke)
				}
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		}))
}
