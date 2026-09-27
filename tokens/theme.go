// Package tokens holds the design values every Kvit app draws with: the
// colours of four themes (Theme), the chrome's type scale and geometry, all
// derived from one interface size (Interface), and the document's text
// settings (Typography). It is the Go port of kvit-ui's src/tokens, and uses
// the same setting keys, so both versions read the same settings file.
//
// Nothing here draws. The root kvitui package applies these values to unison.
package tokens

import (
	"strings"

	"github.com/kvit-s/kvit-ui/palette"
)

// Tokens is one theme's colours. Every colour the library draws with is a
// field here; a component reads these rather than writing a colour of its own.
type Tokens struct {
	// Surfaces.
	WindowBackground, PanelBackground, ListBackground, FooterBackground,
	PopupBackground, ChipBackground, BannerBackground, CodePanelBackground palette.Color

	// Text. OnAccent is the one token no table sets: it is derived from the
	// accent by LabelOn, so a built-in accent and a user's custom one go
	// through the same rule.
	TextPrimary, TextSecondary, TextMuted, TextFaint, TextDisabled, BannerText, OnAccent palette.Color

	// Lines and glyphs. Border is decorative (the rule between two panels,
	// deliberately below 3:1); BorderStrong is where a control's edge is, and
	// is held to 3:1 against the surface behind it.
	Border, BorderStrong, QuoteBar, MutedGlyph palette.Color

	// Interactive tints.
	HoverTint, BlockHoverTint, FocusTint, FocusRing, SelectionTint,
	SelectionActiveTint, BlockSelectionTint palette.Color

	// Accent and meaning.
	Accent, Danger, DangerBright, Success, Warning, PinColor palette.Color

	// Inline styling.
	Marker, InlineCodeBackground, HighlightBackground, Link, SearchMatchBackground,
	SearchCurrentBackground, ChangedTextBackground, AddedTextBackground,
	RemovedTextBackground palette.Color

	// Code highlighting.
	CodeKeyword, CodeType, CodeString, CodeComment, CodeNumber palette.Color

	// The tip callout's own hue.
	CalloutTip palette.Color

	// Dashboard vocabulary: two effort axes, discovered scope and three signals.
	AxisAttention, AxisAttentionText, AxisAgent, AxisAgentText, ScopeDiscovered,
	SignalHard, SignalSoft, SignalHygiene, HatchAlt palette.Color

	// The three chart ramps: eight categorical steps, seven sequential, seven
	// diverging. Each is checked against this theme's surfaces and reserved
	// hues by the palette package's validators in the tests.
	CategoricalRamp, SequentialRamp, DivergingRamp []palette.Color
}

// The built-in themes, plus "system", which follows the desktop.
const (
	Light        = "light"
	Dark         = "dark"
	Sepia        = "sepia"
	HighContrast = "highContrast"
	System       = "system"
)

// AvailableThemes lists every theme id a setting may hold, "system" last.
func AvailableThemes() []string { return []string{Light, Dark, Sepia, HighContrast, System} }

// BuiltInThemes lists the four themes that have colour tables.
func BuiltInThemes() []string { return []string{Light, Dark, Sepia, HighContrast} }

// ReducedMotionSettings lists the values the reduced-motion setting accepts.
func ReducedMotionSettings() []string { return []string{"on", "off", "system"} }

// nearBlack is the darkest label a derived accent label may be: pure black on
// a mid accent is harsher than the rest of the interface.
var nearBlack = palette.Hex("#1a1a1a")

var white = palette.Hex("#ffffff")

// LabelOn is the label colour for text drawn on a fill: near-black or white,
// whichever contrasts more with it. OnAccent is computed with it, and a label
// on any other fill, a danger button or a tag's own colour, should ask it too.
func LabelOn(fill palette.Color) palette.Color {
	if palette.ContrastRatio(nearBlack, fill) >= palette.ContrastRatio(white, fill) {
		return nearBlack
	}
	return white
}

// TokensFor returns a built-in theme's own table, with OnAccent derived from
// its accent. An unknown id gives the light table.
func TokensFor(resolved string) Tokens {
	var t Tokens
	switch resolved {
	case Dark:
		t = darkTable
	case Sepia:
		t = sepiaTable
	case HighContrast:
		t = highContrastTable
	default:
		t = lightTable
	}
	t.OnAccent = LabelOn(t.Accent)
	return t
}

// DisplayName is the label a menu shows for a theme id.
func DisplayName(id string) string {
	if id == HighContrast {
		return "High contrast"
	}
	if id == "" {
		return id
	}
	return strings.ToUpper(id[:1]) + id[1:]
}

// ColorPalette is the user-data palette for folders and tags: content, not
// chrome, so the same in every theme.
func ColorPalette() []palette.Color {
	return hexes("#e05c5c", "#e0a04c", "#58a866", "#4a90d9", "#9068c8", "#d06ca8")
}

// ColorPaletteNames names ColorPalette's entries, in order, for screen readers.
func ColorPaletteNames() []string { return []string{"Red", "Amber", "Green", "Blue", "Purple", "Pink"} }

// HighlightPalette is the soft tints the highlight-colour picker offers. Text
// sits on these, so they are pale where ColorPalette is saturated.
func HighlightPalette() []palette.Color {
	return hexes("#fdf3a9", "#ffd9a8", "#c9ecc9", "#c9e4ff", "#f2ccf2")
}

// HighlightPaletteNames names HighlightPalette's entries, in order.
func HighlightPaletteNames() []string { return []string{"Yellow", "Peach", "Mint", "Sky", "Lilac"} }

// ColorName is what a swatch is announced as: a palette entry's own name, one
// of the two greys the text-colour picker adds, "Theme default" for the empty
// string, or a name built from the hex value for anything else. Colours are
// compared parsed, so "#E05C5C" and "#e05c5c" get the same name.
func ColorName(value string) string {
	if strings.TrimSpace(value) == "" {
		return "Theme default"
	}
	wanted, err := palette.ParseHex(value)
	if err != nil {
		return value
	}
	for i, c := range ColorPalette() {
		if c.Equal(wanted) {
			return ColorPaletteNames()[i]
		}
	}
	for i, c := range HighlightPalette() {
		if c.Equal(wanted) {
			return HighlightPaletteNames()[i]
		}
	}
	switch wanted.Hex() {
	case "#333333":
		return "Near black"
	case "#888888":
		return "Grey"
	}
	return "Custom colour " + wanted.Hex()
}

func hexes(hs ...string) []palette.Color {
	out := make([]palette.Color, len(hs))
	for i, h := range hs {
		out[i] = palette.Hex(h)
	}
	return out
}

// Categorical returns categorical step i of a table, wrapping past the eighth
// in either direction: a chart handed nine series has to draw the ninth as
// something, and a repeat is a smaller lie than an unannounced grey.
func (t Tokens) Categorical(i int) palette.Color {
	n := len(t.CategoricalRamp)
	if n == 0 {
		return t.TextMuted
	}
	return t.CategoricalRamp[((i%n)+n)%n]
}

// ReservedHues lists every colour in the table that already means something,
// which chart ramps keep away from. Surfaces, text and lines are absent: the
// contrast rule already holds a step away from a surface, and nobody mistakes
// a bar for the colour of the text beside it. The two axis text colours are
// absent because they sit so close to their bar hues that reserving both
// leaves no room for an eight-step ramp, and they are never drawn as a mark.
func (t Tokens) ReservedHues() []palette.Color {
	return []palette.Color{
		t.Accent, t.Danger, t.DangerBright, t.Success, t.Warning, t.Link,
		t.PinColor, t.CalloutTip, t.FocusRing, t.AxisAttention, t.AxisAgent,
		t.ScopeDiscovered, t.SignalHard, t.SignalSoft, t.SignalHygiene,
	}
}

// Surfaces lists the grounds chart marks are drawn on, the window's first.
func (t Tokens) Surfaces() []palette.Color {
	return []palette.Color{t.WindowBackground, t.PanelBackground, t.ListBackground, t.ChipBackground}
}

// QtTokenNames lists the tokens the Qt library's tables set, as the
// generator found them; every colour field of Tokens but OnAccent is among
// them.
func QtTokenNames() []string { return append([]string(nil), qtTokenNames...) }
