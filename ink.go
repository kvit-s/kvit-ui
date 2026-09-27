package kvitui

import (
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/tokens"
)

// Ink is one of the theme's colours, chosen by what it is for rather than by
// its value, so a component given an Ink keeps following the theme as it
// changes. inks_gen.go has one for every token: InkTextMuted, InkDanger and
// the rest.
type Ink func(tokens.Tokens) palette.Color

// Of returns the ink's colour in the UI's current theme.
func (i Ink) Of(u *UI) palette.Color { return i(u.Theme.Tokens()) }

// Measure is one of the values Interface derives from the interface size,
// chosen by name (SizeSpace, SizeViewMargin and the rest in names_gen.go),
// so a gap or a padding given as a Measure follows the interface size as it
// changes rather than keeping the value it had when the panel was built.
type Measure func(*tokens.Interface) int

// Of returns the measure at the UI's current interface size.
func (m Measure) Of(u *UI) int { return m(u.Interface) }
