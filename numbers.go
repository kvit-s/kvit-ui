package kvitui

import (
	"math"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Number writes a whole number with its digits grouped the way the reader's
// locale groups them: "250,000" in en-US, "250.000" in de-DE, and "250000"
// under the C locale. A six-figure count run together is read digit by digit.
func (u *UI) Number(n int) string {
	if u.Locale == language.Und {
		return strconv.Itoa(n)
	}
	return message.NewPrinter(u.Locale).Sprintf("%d", n)
}

// CountPhrase is a count with the noun for what it counts: "1 transaction",
// "250,000 transactions". plural is the noun for several; "" adds an s to
// singular, which is right for most English nouns and wrong for "entry", so
// a noun that pluralises some other way passes its own. A negative count says
// nothing at all, which is right for something that is not a countable list.
func (u *UI) CountPhrase(n int, singular, plural string) string {
	switch {
	case n < 0:
		return ""
	case n == 1:
		return u.Number(n) + " " + singular
	case plural != "":
		return u.Number(n) + " " + plural
	default:
		return u.Number(n) + " " + singular + "s"
	}
}

// toFixed writes a number with a fixed count of decimal places, rounding a
// half away from zero as JavaScript's toFixed does, where Go's formatting
// rounds it to even: 0.5 shown with no decimals is "1" here, where fmt
// writes "0".
func toFixed(v float64, places int) string {
	scale := math.Pow(10, float64(places))
	r := math.Round(v*scale) / scale
	if r == 0 {
		r = 0 // not "-0"
	}
	return strconv.FormatFloat(r, 'f', max(0, places), 64)
}

// DecimalSeparator is the character the reader's locale writes between the
// whole part of a number and its fraction: "." in en-US and under the C
// locale, "," in de-DE.
func (u *UI) DecimalSeparator() string {
	if u.Locale == language.Und {
		return "."
	}
	s := []rune(message.NewPrinter(u.Locale).Sprintf("%.1f", 1.5))
	if len(s) < 3 {
		return "."
	}
	return string(s[1])
}
