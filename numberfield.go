package kvitui

import (
	"math"
	"strconv"
	"strings"

	"github.com/kvit-s/kvit-ui/text"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

// NumberField is a number the reader types, right-aligned, as a Figure is,
// because a column of these is compared down the units digit. It checks what
// is typed as it is typed and says what is wrong through the field's Error,
// rather than refusing keys: a field that ignores a key gives no reason, and
// the commonest cause is a decimal separator the reader's locale writes
// differently from the one the field expects. So both are accepted: "3,50"
// typed on a German keyboard is 3.5.
//
// The text is rounded to the field's decimal places when the reader leaves
// the field, not as they type, since reformatting under the caret moves it.
type NumberField struct {
	Field
	// Minimum and Maximum bound the value; no bound unless set.
	Minimum, Maximum float64
	// Decimals is how many decimal places the value has; 0 for a whole
	// number.
	Decimals int
}

// NewNumberField returns an empty number field with no bounds.
func NewNumberField(ui *UI) *NumberField {
	n := &NumberField{Minimum: math.Inf(-1), Maximum: math.Inf(1)}
	n.ui = ui
	n.Self = n
	n.initField(ui, false)
	n.initNumber()
	return n
}

func (n *NumberField) initNumber() {
	n.edit.HAlign = align.End
	n.changed = n.check
	lost := n.edit.LostFocusCallback
	n.edit.LostFocusCallback = func() {
		lost()
		if v, ok := n.Value(); ok && n.Valid() {
			n.SetText(n.format(v))
		}
	}
}

// SetText replaces the text and checks it.
func (n *NumberField) SetText(s string) {
	n.Field.SetText(s)
	n.check()
}

// normalised is the text with the locale's decimal separator made a point.
func (n *NumberField) normalised() string {
	s := strings.TrimSpace(n.Text())
	if sep := n.ui.DecimalSeparator(); sep != "." {
		s = strings.Replace(s, sep, ".", 1)
	}
	return s
}

// Value is the number typed, and false while the text is not a number.
func (n *NumberField) Value() (float64, bool) {
	s := n.normalised()
	lower := strings.ToLower(s)
	if s == "" || strings.Contains(lower, "inf") || strings.Contains(lower, "nan") || strings.Contains(lower, "x") {
		return math.NaN(), false
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN(), false
	}
	return v, true
}

// Valid reports whether the text is a number within the bounds.
func (n *NumberField) Valid() bool {
	v, ok := n.Value()
	return ok && v >= n.Minimum && v <= n.Maximum
}

// format writes a value at the field's decimal places with the reader's
// decimal separator.
func (n *NumberField) format(v float64) string {
	s := toFixed(v, n.Decimals)
	if sep := n.ui.DecimalSeparator(); sep != "." {
		s = strings.Replace(s, ".", sep, 1)
	}
	return s
}

// check says what is wrong with the text, if anything.
func (n *NumberField) check() {
	was := n.Error
	v, ok := n.Value()
	switch {
	case strings.TrimSpace(n.Text()) == "":
		n.Error = ""
	case !ok:
		n.Error = "Not a number"
	case v < n.Minimum:
		n.Error = "Must be at least " + num(n.Minimum)
	case v > n.Maximum:
		n.Error = "Must be at most " + num(n.Maximum)
	default:
		n.Error = ""
	}
	if n.Error != was {
		n.MarkForLayoutAndRedraw()
	}
}

// MoneyField is an amount of money: a number field that works in minor units,
// takes its decimal places from the currency, and shows the currency inside
// its right-hand end. The reader types 12.34 and MinorUnits is 1234, a whole
// number, which is what a ledger stores; money in floating point is where
// 0.1 + 0.2 is not 0.3 and a balance drifts by a penny nobody can find. The
// yen has no minor digits and the Bahraini dinar three, and a field that
// assumes two divides or multiplies those by a hundred. The currency is beside
// the amount rather than in a column heading, because a ledger can hold more
// than one.
type MoneyField struct {
	NumberField
	// Currency is the ISO 4217 code shown in the field; "" shows none, for a
	// field where the currency is already clear.
	Currency string
}

// NewMoneyField returns an empty amount in a currency with minorDigits
// decimal places: 2 for most, 0 for the yen, 3 for the dinar.
func NewMoneyField(ui *UI, currency string, minorDigits int) *MoneyField {
	f := &MoneyField{Currency: currency}
	f.Minimum, f.Maximum, f.Decimals = math.Inf(-1), math.Inf(1), minorDigits
	f.ui = ui
	f.Self = f
	f.initField(ui, false)
	f.initNumber()
	near := func() float32 { return float32(ui.Interface.SpaceNear()) }
	f.padRight = func() float32 {
		if f.Currency == "" {
			return near()
		}
		w, _ := f.marker().Size()
		return w + 2*near()
	}
	f.decorate = func(gc *unison.Canvas, box geom.Rect) {
		if f.Currency == "" {
			return
		}
		l := f.marker()
		w, h := l.Size()
		l.Draw(gc, box.Right()-near()-w, box.Y+(box.Height-h)/2)
	}
	f.spoken = func() string {
		if f.Currency == "" {
			return f.Label
		}
		return f.Label + " in " + f.Currency
	}
	return f
}

// MinorUnits is the amount as a whole number of minor units — pence, cents,
// sen — and false while the text is not a valid amount.
func (f *MoneyField) MinorUnits() (int64, bool) {
	v, ok := f.Value()
	if !ok || !f.Valid() {
		return 0, false
	}
	return int64(math.Round(v * math.Pow(10, float64(f.Decimals)))), true
}

// marker is the currency code as it is drawn: small and muted.
func (f *MoneyField) marker() *text.Layout {
	ui := f.ui
	return ui.Fonts.Layout([]text.Span{{Text: f.Currency, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, ui.Theme.Tokens().TextMuted)}}, text.Options{})
}
