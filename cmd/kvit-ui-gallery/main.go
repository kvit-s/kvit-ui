// Command kvit-ui-gallery shows every kvit-ui component with its states, in
// each theme and at each interface size. For now it opens an empty window.
//
//	kvit-ui-gallery               browse
//	kvit-ui-gallery --smoke 3s    close after 3 s and report when the first frame was drawn
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/richardwilkes/unison"
)

var started = time.Now()

func main() {
	smoke := flag.Duration("smoke", 0, "close the window after this long, printing when its first frame was drawn")
	flag.Parse()
	unison.Start(unison.StartupFinishedCallback(func() {
		var wnd *unison.Window
		var err error
		wnd, err = newGalleryWindow(func() {
			if *smoke > 0 {
				r := wnd.ContentRect()
				s := wnd.BackingScale()
				fmt.Printf("first frame after %d ms, %s/%s, window %.0f x %.0f at scale %.2f\n",
					time.Since(started).Milliseconds(), runtime.GOOS, runtime.GOARCH, r.Width, r.Height, s.X)
				unison.InvokeTaskAfter(wnd.Dispose, *smoke)
			}
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}))
}
