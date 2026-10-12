package main

// The gallery's icon. It is icons/kvit-ui-gallery.svg, drawn by
// tools/appicon at the sizes the window title asks for; the Windows
// executable carries it as its resource (winres/kvit-ui-gallery.ico,
// rsrc_windows_amd64.syso).

import (
	"embed"
	"path"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

//go:embed icons/*.png
var iconFiles embed.FS

// iconScale draws one pixel of an icon picture to one pixel of the screen.
const iconScale = 1

// setTitleIcons gives the windows made after it the gallery's icon: unison
// gives a window the system's generic icon rather than the executable's.
// It runs once unison has started, since unison decodes PNGs only then.
func setTitleIcons() {
	entries, _ := iconFiles.ReadDir("icons")
	var images []*unison.Image
	for _, e := range entries {
		data, err := iconFiles.ReadFile(path.Join("icons", e.Name()))
		if err != nil {
			continue
		}
		if img, err := unison.NewImageFromBytes(data, geom.NewPoint(iconScale, iconScale)); err == nil {
			images = append(images, img)
		}
	}
	if len(images) > 0 {
		unison.DefaultTitleIcons = images
	}
}
