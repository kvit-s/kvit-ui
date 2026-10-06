package platform

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/png"
	"math/bits"
)

// A picture on Windows' clipboard is always there as a device-independent
// bitmap (DIB), the formats CF_DIB and CF_DIBV5, whatever else it is there
// as, since Windows makes each of them from the other or from CF_BITMAP. A
// DIB is a BMP file without its 14-byte file header: a BITMAPINFOHEADER (40
// bytes) or a larger header up to BITMAPV5HEADER (124 bytes), the colour
// masks or the palette, then the rows of pixels. Go's BMP decoder reads only
// some of the layouts programs put there (no colour masks after a 40-byte
// header, and none but its own inside a larger one), and takes a 32-bit
// pixel's fourth byte as its transparency where a screenshot leaves that
// byte 0 throughout, so the DIB is read here.

const (
	biRGB            = 0
	biBitfields      = 3
	biAlphaBitfields = 6
	dibHeaderSize    = 40
	// dibMaxPixels bounds the size a DIB may claim, 256 megapixels, so a
	// damaged header cannot ask for gigabytes.
	dibMaxPixels = 1 << 28
)

var errDIB = errors.New("platform: not a bitmap this reads")

// dibPNG turns a DIB into PNG bytes; false for anything it cannot read.
func dibPNG(dib []byte) ([]byte, bool) {
	picture, err := decodeDIB(dib)
	if err != nil {
		return nil, false
	}
	var out bytes.Buffer
	if err = png.Encode(&out, picture); err != nil {
		return nil, false
	}
	return out.Bytes(), true
}

// decodeDIB reads a DIB of 1, 4 or 8 bits a pixel with a palette, of 24 bits,
// or of 16 or 32 bits with the default colour masks or its own. A 32-bit
// picture whose every pixel has transparency 0 is taken as opaque: that byte
// is unused in a screenshot.
func decodeDIB(dib []byte) (*image.NRGBA, error) {
	if len(dib) < dibHeaderSize {
		return nil, errDIB
	}
	le := binary.LittleEndian
	headerSize := int64(le.Uint32(dib[0:]))
	width := int64(int32(le.Uint32(dib[4:])))
	height := int64(int32(le.Uint32(dib[8:])))
	planes, depth := le.Uint16(dib[12:]), le.Uint16(dib[14:])
	compression, colorsUsed := le.Uint32(dib[16:]), int64(le.Uint32(dib[32:]))
	if headerSize < dibHeaderSize || headerSize > int64(len(dib)) || planes != 1 {
		return nil, errDIB
	}
	// Rows run from the bottom up, or from the top down when the height is
	// negative.
	topDown := height < 0
	if topDown {
		height = -height
	}
	if width <= 0 || height <= 0 || width*height > dibMaxPixels {
		return nil, errDIB
	}
	offset := headerSize
	// The masks of red, green, blue and transparency.
	var masks [4]uint32
	var palette [][4]byte
	switch {
	case compression == biRGB && (depth == 1 || depth == 4 || depth == 8):
		if colorsUsed == 0 || colorsUsed > 1<<depth {
			colorsUsed = 1 << depth
		}
		if offset+4*colorsUsed > int64(len(dib)) {
			return nil, errDIB
		}
		// An index past the colours used is black.
		palette = make([][4]byte, 1<<depth)
		for i := range palette {
			palette[i][3] = 0xff
		}
		for i := range colorsUsed {
			entry := dib[offset+4*i:]
			palette[i] = [4]byte{entry[2], entry[1], entry[0], 0xff}
		}
		offset += 4 * colorsUsed
	case compression == biRGB && (depth == 16 || depth == 24 || depth == 32):
		masks = [4]uint32{0xff0000, 0xff00, 0xff, 0xff000000}
		if depth == 16 {
			masks = [4]uint32{0x7c00, 0x03e0, 0x001f, 0}
		}
	case (compression == biBitfields || compression == biAlphaBitfields) && (depth == 16 || depth == 32):
		count := int64(3)
		if compression == biAlphaBitfields {
			count = 4
		}
		// The masks follow a 40-byte header and are inside a larger one,
		// which has room for the transparency mask from 56 bytes on.
		if headerSize == dibHeaderSize {
			offset += 4 * count
		} else if headerSize >= 56 {
			count = 4
		} else if headerSize < dibHeaderSize+4*count {
			return nil, errDIB
		}
		if dibHeaderSize+4*count > int64(len(dib)) {
			return nil, errDIB
		}
		for i := range count {
			masks[i] = le.Uint32(dib[dibHeaderSize+4*i:])
		}
	default:
		return nil, errDIB
	}
	// A bitmap of more than 8 bits a pixel may still carry a palette, for
	// drawing on a screen of fewer colours, which goes before the pixels.
	if palette == nil {
		offset += 4 * min(colorsUsed, int64(len(dib)))
	}
	stride := (width*int64(depth) + 31) / 32 * 4
	if offset+stride*height > int64(len(dib)) {
		return nil, errDIB
	}
	picture := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))
	transparent := true
	for y := range height {
		row := dib[offset+y*stride:]
		target := height - 1 - y
		if topDown {
			target = y
		}
		line := picture.Pix[target*int64(picture.Stride):]
		for x := range width {
			var pixel [4]byte
			switch depth {
			case 1, 4, 8:
				at := x * int64(depth)
				index := row[at/8] >> (8 - int64(depth) - at%8) & (1<<depth - 1)
				pixel = palette[index]
			case 24:
				pixel = [4]byte{row[3*x+2], row[3*x+1], row[3*x], 0xff}
			default:
				var value uint32
				if depth == 16 {
					value = uint32(le.Uint16(row[2*x:]))
				} else {
					value = le.Uint32(row[4*x:])
				}
				pixel = [4]byte{maskedChannel(value, masks[0]), maskedChannel(value, masks[1]),
					maskedChannel(value, masks[2]), 0xff}
				if masks[3] != 0 {
					pixel[3] = maskedChannel(value, masks[3])
				}
			}
			transparent = transparent && pixel[3] == 0
			copy(line[4*x:], pixel[:])
		}
	}
	if transparent {
		for i := 3; i < len(picture.Pix); i += 4 {
			picture.Pix[i] = 0xff
		}
	}
	return picture, nil
}

// maskedChannel is the bits of value under mask, scaled to 0 to 255.
func maskedChannel(value, mask uint32) byte {
	if mask == 0 {
		return 0
	}
	shift := bits.TrailingZeros32(mask)
	largest := uint64(mask >> shift)
	return byte((uint64((value&mask)>>shift)*255 + largest/2) / largest)
}
