package main

// The pictures the specimens put in front of a line's words: a picture has its
// own colours, so they are read from files rather than written with the design
// values. The SVG and the square PNG are the gallery's icon
// (icons/kvit-ui-gallery.svg); pictures/wide.png is a PNG twice as wide as it
// is tall, to show a picture keeping its shape.

import (
	_ "embed" // for the pictures
	"sync"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/filtermode"
	"github.com/richardwilkes/unison/enums/mipmapmode"
)

//go:embed icons/kvit-ui-gallery.svg
var svgPictureFile string

//go:embed icons/kvit-ui-gallery-64.png
var pngPictureFile []byte

//go:embed pictures/wide.png
var widePictureFile []byte

// The pictures are decoded when first asked for, since unison decodes a PNG
// only once it has started, and once, since every shot builds its page again.
var (
	svgPicture = sync.OnceValue(func() unison.Drawable {
		svg, err := unison.NewSVGFromContentString(svgPictureFile)
		if err != nil {
			panic(err)
		}
		return &unison.DrawableSVG{SVG: svg, Size: svg.Size()}
	})
	pngPicture  = sync.OnceValue(func() unison.Drawable { return pngFrom(pngPictureFile) })
	widePicture = sync.OnceValue(func() unison.Drawable { return pngFrom(widePictureFile) })
)

func pngFrom(data []byte) unison.Drawable {
	img, err := unison.NewImageFromBytes(data, geom.NewPoint(iconScale, iconScale))
	if err != nil {
		panic(err)
	}
	return img
}

// pictureLead is a picture as a panel for SetLead, 16 design pixels square. It
// blends the picture's pixels, since unison's default sampling leaves a reduced
// PNG jagged.
func pictureLead(ui *kvitui.UI, picture unison.Drawable) unison.Paneler {
	p := unison.NewPanel()
	p.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		sampling := &unison.SamplingOptions{FilterMode: filtermode.Linear, MipMapMode: mipmapmode.Linear}
		picture.DrawInRect(gc, p.ContentRect(false), sampling, nil)
	}
	return kvitui.Sized(ui, kvitui.Px(16), kvitui.Px(16), p)
}
