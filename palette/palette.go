// Package palette is the colour arithmetic the design values are checked
// with: sRGB colours, OKLab and OKLCH, perceptual distance, simulation of the
// three colour-vision deficiencies, WCAG contrast, and the rules a chart ramp
// has to meet. It is a port of kvit-ui's src/tokens/palette.cpp and gives the
// same answers.
package palette

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color is an opaque sRGB colour, each channel from 0 to 1.
type Color struct {
	R, G, B float64
}

// Hex parses "#rrggbb" or "#rgb". It panics on anything else, which suits the
// constant tables it is used for; ParseHex reports the error instead.
func Hex(s string) Color {
	c, err := ParseHex(s)
	if err != nil {
		panic(err)
	}
	return c
}

// ParseHex parses "#rrggbb" or "#rgb", in either case.
func ParseHex(s string) (Color, error) {
	t := strings.TrimSpace(s)
	if !strings.HasPrefix(t, "#") {
		return Color{}, fmt.Errorf("palette: %q is not a #rrggbb colour", s)
	}
	h := t[1:]
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) != 6 {
		return Color{}, fmt.Errorf("palette: %q is not a #rrggbb colour", s)
	}
	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return Color{}, fmt.Errorf("palette: %q is not a #rrggbb colour", s)
	}
	return RGB8(uint8(v>>16), uint8(v>>8), uint8(v)), nil
}

// RGB8 builds a colour from 8-bit channels.
func RGB8(r, g, b uint8) Color {
	return Color{float64(r) / 255, float64(g) / 255, float64(b) / 255}
}

// Shade is the black a floating surface darkens what is behind it with, at
// a share of its strength: the same in every theme, as the library's
// spotlight and floating view draw it.
var Shade = RGB8(0, 0, 0)

// RGBA8 returns the channels rounded to 8 bits.
func (c Color) RGBA8() (r, g, b uint8) {
	q := func(v float64) uint8 { return uint8(math.Round(math.Min(1, math.Max(0, v)) * 255)) }
	return q(c.R), q(c.G), q(c.B)
}

// Hex returns "#rrggbb" in lower case.
func (c Color) Hex() string {
	r, g, b := c.RGBA8()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// String implements fmt.Stringer.
func (c Color) String() string { return c.Hex() }

// Equal reports whether the two colours are the same at 8 bits per channel,
// which is the precision every colour in the tables is written in.
func (c Color) Equal(o Color) bool { return c.Hex() == o.Hex() }

// sRGB's transfer function and its inverse. Every piece of arithmetic below
// happens in linear light, because averaging, mixing and the colour-vision
// matrices are all linear operations.
func toLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func toEncoded(v float64) float64 {
	if v <= 0.0031308 {
		return v * 12.92
	}
	return 1.055*math.Pow(v, 1.0/2.4) - 0.055
}

type linear struct{ r, g, b float64 }

func linearOf(c Color) linear { return linear{toLinear(c.R), toLinear(c.G), toLinear(c.B)} }

func clamp01(v float64) float64 { return math.Min(1, math.Max(0, v)) }

func colorOf(l linear) Color {
	return Color{clamp01(toEncoded(clamp01(l.r))), clamp01(toEncoded(clamp01(l.g))), clamp01(toEncoded(clamp01(l.b)))}
}

type oklab struct{ l, a, b float64 }

// Björn Ottosson's OKLab.
func toOklab(c linear) oklab {
	l := 0.4122214708*c.r + 0.5363325363*c.g + 0.0514459929*c.b
	m := 0.2119034982*c.r + 0.6806995451*c.g + 0.1073969566*c.b
	s := 0.0883024619*c.r + 0.2817188376*c.g + 0.6299787005*c.b
	l_, m_, s_ := math.Cbrt(l), math.Cbrt(m), math.Cbrt(s)
	return oklab{
		0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_,
		1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_,
		0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_,
	}
}

func fromOklab(v oklab) linear {
	l_ := v.l + 0.3963377774*v.a + 0.2158037573*v.b
	m_ := v.l - 0.1055613458*v.a - 0.0638541728*v.b
	s_ := v.l - 0.0894841775*v.a - 1.2914855480*v.b
	l, m, s := l_*l_*l_, m_*m_*m_, s_*s_*s_
	return linear{
		+4.0767416621*l - 3.3077115913*m + 0.2309699292*s,
		-1.2684380046*l + 2.6097574011*m - 0.3413193965*s,
		-0.0041960863*l - 0.7034186147*m + 1.7076147010*s,
	}
}

// Oklch is a colour in OKLCH: lightness 0–1, chroma, hue in degrees.
type Oklch struct {
	L, C, H float64
}

// ToOklch converts a colour to OKLCH.
func ToOklch(c Color) Oklch {
	lab := toOklab(linearOf(c))
	hue := math.Atan2(lab.b, lab.a) * 180 / math.Pi
	if hue < 0 {
		hue += 360
	}
	return Oklch{lab.l, math.Hypot(lab.a, lab.b), hue}
}

// FromOklch converts back to sRGB, clamping anything outside the gamut.
func FromOklch(v Oklch) Color {
	r := v.H * math.Pi / 180
	return colorOf(fromOklab(oklab{v.L, v.C * math.Cos(r), v.C * math.Sin(r)}))
}

// InGamut reports whether an OKLCH value is inside sRGB, allowing a
// thousandth for rounding.
func InGamut(v Oklch) bool {
	r := v.H * math.Pi / 180
	l := fromOklab(oklab{v.L, v.C * math.Cos(r), v.C * math.Sin(r)})
	const slack = 0.001
	for _, ch := range []float64{l.r, l.g, l.b} {
		if ch < -slack || ch > 1+slack {
			return false
		}
	}
	return true
}

// DeltaE is the OKLab distance between two colours, times 100.
func DeltaE(a, b Color) float64 {
	x, y := toOklab(linearOf(a)), toOklab(linearOf(b))
	return 100 * math.Sqrt((x.l-y.l)*(x.l-y.l)+(x.a-y.a)*(x.a-y.a)+(x.b-y.b)*(x.b-y.b))
}

// Deficiency is one of the three dichromacies.
type Deficiency int

// The three deficiencies the ramps are checked against.
const (
	Protanopia   Deficiency = iota // no long-wavelength cone: red and green converge
	Deuteranopia                   // no medium-wavelength cone: the common one
	Tritanopia                     // no short-wavelength cone: blue and green converge
)

// Deficiencies lists all three, in the order the version checks them.
var Deficiencies = []Deficiency{Protanopia, Deuteranopia, Tritanopia}

func (d Deficiency) String() string {
	switch d {
	case Protanopia:
		return "protanopia"
	case Deuteranopia:
		return "deuteranopia"
	case Tritanopia:
		return "tritanopia"
	}
	return "unknown"
}

// The Machado, Oliveira and Fernandes (2009) matrices at full severity,
// applied in linear light. Rows are the red, green and blue outputs.
var deficiencyMatrix = map[Deficiency][9]float64{
	Protanopia: {
		0.152286, 1.052583, -0.204868,
		0.114503, 0.786281, 0.099216,
		-0.003882, -0.048116, 1.051998,
	},
	Deuteranopia: {
		0.367322, 0.860646, -0.227968,
		0.280085, 0.672501, 0.047413,
		-0.011820, 0.042940, 0.968881,
	},
	Tritanopia: {
		1.255528, -0.076749, -0.178779,
		-0.078411, 0.930809, 0.147602,
		0.004733, 0.691367, 0.303900,
	},
}

// Simulate returns the colour as a reader with the deficiency sees it.
func Simulate(c Color, d Deficiency) Color {
	m := deficiencyMatrix[d]
	l := linearOf(c)
	return colorOf(linear{
		m[0]*l.r + m[1]*l.g + m[2]*l.b,
		m[3]*l.r + m[4]*l.g + m[5]*l.b,
		m[6]*l.r + m[7]*l.g + m[8]*l.b,
	})
}

// RelativeLuminance is WCAG 2.1's relative luminance.
func RelativeLuminance(c Color) float64 {
	l := linearOf(c)
	return 0.2126*l.r + 0.7152*l.g + 0.0722*l.b
}

// ContrastRatio is WCAG 2.1's contrast ratio, from 1 (identical) to 21.
func ContrastRatio(a, b Color) float64 {
	la, lb := RelativeLuminance(a), RelativeLuminance(b)
	return (math.Max(la, lb) + 0.05) / (math.Min(la, lb) + 0.05)
}

// Darker divides the colour's HSV value by factor, as the colour::darker
// does: 1.15 is the pressed shade of a filled button.
func (c Color) Darker(factor float64) Color {
	if factor <= 0 {
		return c
	}
	hi := max(c.R, c.G, c.B)
	if hi == 0 {
		return c
	}
	// Scaling all three channels keeps the hue and saturation and divides
	// the value, which is the largest channel.
	scale := (hi / factor) / hi
	return Color{R: c.R * scale, G: c.G * scale, B: c.B * scale}
}
