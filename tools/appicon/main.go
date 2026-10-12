// Command appicon draws the gallery's icon,
// cmd/kvit-ui-gallery/icons/kvit-ui-gallery.svg, at the sizes the program
// carries it: the PNGs the window title takes
// (cmd/kvit-ui-gallery/icons/kvit-ui-gallery-<size>.png, which identity.go
// embeds) and the Windows icon (cmd/kvit-ui-gallery/winres/kvit-ui-gallery.ico,
// which winres/make-syso.sh compiles into the executable). The pictures are
// committed, so no build needs a renderer: rerun this after editing the SVG,
// then make-syso.sh.
//
//	go run ./tools/appicon
//
// The SVG is drawn by unison's own raster canvas, so the pictures need
// nothing installed. The .ico holds a PNG at each size, which Windows reads
// since Vista. Kvit Works draws its icon with the same tool.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

const (
	source = "cmd/kvit-ui-gallery/icons/kvit-ui-gallery.svg"
	// pngs is where the window's sizes go, as kvit-ui-gallery-<size>.png.
	pngs = "cmd/kvit-ui-gallery/icons"
	ico  = "cmd/kvit-ui-gallery/winres/kvit-ui-gallery.ico"
	// ppi draws one pixel of the picture to one point of the drawing.
	ppi = 72
)

// windowSizes are the sizes the window title asks for; icoSizes those
// Windows Explorer and the taskbar do.
var (
	windowSizes = []int{16, 24, 32, 48, 64}
	icoSizes    = []int{16, 24, 32, 48, 64, 128, 256}
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "appicon:", err)
		os.Exit(1)
	}
}

func run() error {
	data, err := os.ReadFile(filepath.FromSlash(source))
	if err != nil {
		return err
	}
	svg, err := unison.NewSVGFromContentString(string(data))
	if err != nil {
		return err
	}
	drawn := map[int][]byte{}
	for _, size := range icoSizes {
		if drawn[size], err = draw(svg, size); err != nil {
			return err
		}
	}
	for _, size := range windowSizes {
		file := filepath.FromSlash(fmt.Sprintf("%s/kvit-ui-gallery-%d.png", pngs, size))
		if err := os.WriteFile(file, drawn[size], 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", filepath.ToSlash(file))
	}
	if err := os.WriteFile(filepath.FromSlash(ico), pack(drawn), 0o644); err != nil {
		return err
	}
	fmt.Println("wrote", ico)
	return nil
}

// draw is the SVG as a PNG size pixels square.
func draw(svg *unison.SVG, size int) ([]byte, error) {
	img, err := unison.NewImageFromDrawing(size, size, ppi, func(c *unison.Canvas) {
		svg.DrawInRect(c, geom.NewRect(0, 0, float32(size), float32(size)), nil, nil)
	})
	if err != nil {
		return nil, err
	}
	return img.ToPNG(9)
}

// pack is an .ico of the PNGs drawn at icoSizes: its header, a directory
// entry for each size, then the PNGs themselves.
func pack(drawn map[int][]byte) []byte {
	var out bytes.Buffer
	write := func(v any) { _ = binary.Write(&out, binary.LittleEndian, v) }
	write([3]uint16{0, 1, uint16(len(icoSizes))})
	offset := 6 + 16*len(icoSizes)
	for _, size := range icoSizes {
		// Width and height are a byte each, and 0 stands for 256.
		dim := uint8(size % 256)
		write([4]uint8{dim, dim, 0, 0})
		write([2]uint16{1, 32})
		write([2]uint32{uint32(len(drawn[size])), uint32(offset)})
		offset += len(drawn[size])
	}
	for _, size := range icoSizes {
		out.Write(drawn[size])
	}
	return out.Bytes()
}
