package platform

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// dib builds a DIB with a header of headerSize bytes, the masks inside it
// when it is larger than 40 bytes, then tail (masks after a 40-byte header,
// or a palette) and the pixel rows.
func dib(headerSize int, width, height int32, depth uint16, compression uint32, colorsUsed uint32,
	masks []uint32, tail, pixels []byte) []byte {
	header := make([]byte, headerSize)
	le := binary.LittleEndian
	le.PutUint32(header[0:], uint32(headerSize))
	le.PutUint32(header[4:], uint32(width))
	le.PutUint32(header[8:], uint32(height))
	le.PutUint16(header[12:], 1)
	le.PutUint16(header[14:], depth)
	le.PutUint32(header[16:], compression)
	le.PutUint32(header[32:], colorsUsed)
	for i, mask := range masks {
		le.PutUint32(header[40+4*i:], mask)
	}
	return append(append(header, tail...), pixels...)
}

func words(values ...uint32) []byte {
	out := make([]byte, 4*len(values))
	for i, v := range values {
		binary.LittleEndian.PutUint32(out[4*i:], v)
	}
	return out
}

func wantPixels(t *testing.T, got image.Image, want [][]color.NRGBA) {
	t.Helper()
	if size := got.Bounds().Size(); size.Y != len(want) || size.X != len(want[0]) {
		t.Fatalf("size %v, want %dx%d", size, len(want[0]), len(want))
	}
	for y, row := range want {
		for x, w := range row {
			if c := color.NRGBAModel.Convert(got.At(x, y)).(color.NRGBA); c != w {
				t.Errorf("pixel %d,%d is %v, want %v", x, y, c, w)
			}
		}
	}
}

var (
	red   = color.NRGBA{0xff, 0, 0, 0xff}
	green = color.NRGBA{0, 0xff, 0, 0xff}
	blue  = color.NRGBA{0, 0, 0xff, 0xff}
	white = color.NRGBA{0xff, 0xff, 0xff, 0xff}
)

// A 24-bit bitmap's rows run from the bottom up, each padded to four bytes,
// and its pixels are blue, green, red.
func TestA24BitBitmapIsReadBottomUp(t *testing.T) {
	pixels := []byte{
		0xff, 0, 0, 0xff, 0xff, 0xff, 0, 0, // bottom row: blue, white, padding
		0, 0, 0xff, 0, 0xff, 0, 0, 0, // top row: red, green, padding
	}
	got, err := decodeDIB(dib(40, 2, 2, 24, biRGB, 0, nil, nil, pixels))
	if err != nil {
		t.Fatal(err)
	}
	wantPixels(t, got, [][]color.NRGBA{{red, green}, {blue, white}})
}

// A screenshot is 32 bits a pixel with the colour masks after a 40-byte
// header and the fourth byte 0 throughout, which is an opaque picture.
func TestAScreenshotsUnusedFourthByteIsOpaque(t *testing.T) {
	masks := words(0xff0000, 0xff00, 0xff)
	pixels := words(0x00ff0000, 0x0000ff00)
	png, ok := dibPNG(dib(40, 2, 1, 32, biBitfields, 0, nil, masks, pixels))
	if !ok {
		t.Fatal("not read")
	}
	got, err := pngDecode(png)
	if err != nil {
		t.Fatal(err)
	}
	wantPixels(t, got, [][]color.NRGBA{{red, green}})
}

// A version 5 header holds its masks, transparency's among them, and a
// negative height runs the rows from the top down.
func TestAV5BitmapKeepsItsTransparency(t *testing.T) {
	masks := []uint32{0xff0000, 0xff00, 0xff, 0xff000000}
	pixels := words(0x80ff0000, 0x000000ff, 0xff00ff00, 0xffffffff)
	got, err := decodeDIB(dib(124, 2, -2, 32, biBitfields, 0, masks, nil, pixels))
	if err != nil {
		t.Fatal(err)
	}
	wantPixels(t, got, [][]color.NRGBA{{{0xff, 0, 0, 0x80}, {0, 0, 0xff, 0}}, {green, white}})
}

// A 32-bit bitmap without masks takes the fourth byte as transparency once
// any pixel has some.
func TestA32BitBitmapWithAlphaKeepsIt(t *testing.T) {
	got, err := decodeDIB(dib(40, 2, 1, 32, biRGB, 0, nil, nil, words(0x40ff0000, 0x00000000)))
	if err != nil {
		t.Fatal(err)
	}
	wantPixels(t, got, [][]color.NRGBA{{{0xff, 0, 0, 0x40}, {}}})
}

// An 8-bit bitmap names a colour of its palette for each pixel, and a 1-bit
// one packs eight pixels into a byte from the highest bit down.
func TestPalettedBitmapsAreRead(t *testing.T) {
	palette := words(0x00ff0000, 0x000000ff)
	got, err := decodeDIB(dib(40, 3, 1, 8, biRGB, 2, nil, palette, []byte{1, 0, 1, 0}))
	if err != nil {
		t.Fatal(err)
	}
	wantPixels(t, got, [][]color.NRGBA{{blue, red, blue}})

	got, err = decodeDIB(dib(40, 3, 1, 1, biRGB, 0, nil, palette, []byte{0b10100000, 0, 0, 0}))
	if err != nil {
		t.Fatal(err)
	}
	wantPixels(t, got, [][]color.NRGBA{{blue, red, blue}})
}

// A 16-bit bitmap without masks is five bits of each colour.
func TestA16BitBitmapIsFiveBitsAColour(t *testing.T) {
	pixels := []byte{0x00, 0x7c, 0x1f, 0x00} // 0x7c00 red, 0x001f blue
	got, err := decodeDIB(dib(40, 2, 1, 16, biRGB, 0, nil, nil, pixels))
	if err != nil {
		t.Fatal(err)
	}
	wantPixels(t, got, [][]color.NRGBA{{red, blue}})
}

// Too short, a header longer than the data, rows cut off, run-length
// compression and an impossible size are no picture.
func TestADamagedOrUnknownBitmapIsNoPicture(t *testing.T) {
	whole := dib(40, 2, 2, 24, biRGB, 0, nil, nil, make([]byte, 16))
	huge := dib(40, 1<<20, 1<<20, 24, biRGB, 0, nil, nil, nil)
	header := dib(40, 1, 1, 24, biRGB, 0, nil, nil, make([]byte, 4))
	binary.LittleEndian.PutUint32(header[0:], 200)
	for name, data := range map[string][]byte{
		"nothing":       nil,
		"short":         whole[:30],
		"header":        header,
		"rows cut off":  whole[:len(whole)-1],
		"run-length":    dib(40, 2, 2, 8, 1, 0, nil, make([]byte, 1024), make([]byte, 8)),
		"huge":          huge,
		"no width":      dib(40, 0, 2, 24, biRGB, 0, nil, nil, make([]byte, 16)),
		"alpha in a V2": dib(52, 1, 1, 32, biAlphaBitfields, 0, []uint32{0xff0000, 0xff00, 0xff}, nil, make([]byte, 4)),
	} {
		if _, ok := dibPNG(data); ok {
			t.Errorf("%s was read as a picture", name)
		}
	}
}

func pngDecode(data []byte) (image.Image, error) { return png.Decode(bytes.NewReader(data)) }
