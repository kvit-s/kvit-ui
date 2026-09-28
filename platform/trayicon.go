package platform

// The tray icon's picture at the size a system draws it: the app gives it
// at a few sizes, and each backend asks for the one it needs.

import (
	"image"

	"golang.org/x/image/draw"
)

// iconAt returns the picture size pixels square: the smallest of sizes at
// least that large, or the largest when none is, scaled to fit. It is nil
// when there is no picture.
func iconAt(sizes []image.Image, size int) *image.NRGBA {
	var best image.Image
	for _, img := range sizes {
		if img == nil {
			continue
		}
		side := min(img.Bounds().Dx(), img.Bounds().Dy())
		switch {
		case best == nil:
			best = img
		case side >= size && (bestSide(best) < size || side < bestSide(best)):
			best = img
		case bestSide(best) < size && side > bestSide(best):
			best = img
		}
	}
	if best == nil || size <= 0 {
		return nil
	}
	out := image.NewNRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(out, out.Bounds(), best, best.Bounds(), draw.Src, nil)
	return out
}

func bestSide(img image.Image) int { return min(img.Bounds().Dx(), img.Bounds().Dy()) }
