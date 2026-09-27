package kvitui

import (
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// BadgeTone is the colour of a badge's pill.
type BadgeTone int

const (
	// BadgeAccent is the default: something the reader has to look at.
	BadgeAccent BadgeTone = iota
	// BadgeNeutral is a count that asks for nothing.
	BadgeNeutral
	// BadgeDanger is a count of things that are wrong.
	BadgeDanger
)

// Badge is a count attached to something else: unread items, pending
// changes, matches. The number is small on the screen and possibly large in
// value, so past Max it says "99+" rather than growing wide enough to move
// what it is attached to, and at zero it is not drawn at all, because a
// badge reading 0 says "look here" about nothing. The number is written with
// the reader's digit grouping, drawn and announced alike.
type Badge struct {
	unison.Panel
	ui *UI
	// Count is the number; zero or less hides the badge.
	Count int
	// Max is the largest count written out; 99 unless set.
	Max int
	// Counted is the word for one of the things counted, so a screen reader
	// hears "1 decision" or "11 decisions"; "" says "items".
	Counted string
	// CountedPlural is the word for several; "" adds an s to Counted. Two
	// words rather than one already inflected, because one word is right at
	// one count and wrong at every other.
	CountedPlural string
	// Tone is the pill's colour.
	Tone BadgeTone
}

// NewBadge returns an accent badge showing a count.
func NewBadge(ui *UI, count int) *Badge {
	b := &Badge{ui: ui, Count: count, Max: 99}
	b.Self = b
	b.SetSizer(b.sizes)
	b.DrawCallback = b.draw
	return b
}

// Shows reports whether the badge is drawn, which it is for a count above
// zero.
func (b *Badge) Shows() bool { return b.Count > 0 }

// Text is what the pill says: the count, or the cap and a plus past it.
func (b *Badge) Text() string {
	most := b.Max
	if most <= 0 {
		most = 99
	}
	if b.Count > most {
		return b.ui.Number(most) + "+"
	}
	return b.ui.Number(b.Count)
}

// Phrase is what a screen reader is told: the whole count and its noun,
// "214 decisions", never capped, since a screen reader has room for the
// number and is the one reader who cannot look it up somewhere else.
func (b *Badge) Phrase() string {
	if b.Counted == "" {
		return b.ui.CountPhrase(b.Count, "item", "")
	}
	return b.ui.CountPhrase(b.Count, b.Counted, b.CountedPlural)
}

func (b *Badge) tone(t tokens.Tokens) palette.Color {
	switch b.Tone {
	case BadgeDanger:
		return t.Danger
	case BadgeNeutral:
		return t.TextMuted
	}
	return t.Accent
}

func (b *Badge) layout() *text.Layout {
	st := b.ui.Chrome(b.ui.Size(RoleCaption), text.Regular, tokens.LabelOn(b.tone(b.ui.Theme.Tokens())))
	st.Tabular = true
	return b.ui.Fonts.Layout([]text.Span{{Text: b.Text(), Style: st}}, text.Options{})
}

func (b *Badge) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	if !b.Shows() {
		return geom.Size{}, geom.Size{}, geom.Size{}
	}
	h := float32(b.ui.Interface.PillHeight())
	w, _ := b.layout().Size()
	size := geom.NewSize(max(h, w+float32(b.ui.Interface.SpaceNear())), h)
	return size, size, size
}

func (b *Badge) draw(gc *unison.Canvas, _ geom.Rect) {
	if !b.Shows() {
		return
	}
	r := b.ContentRect(false)
	painterFor(gc, b.ui).round(r, r.Height/2, b.tone(b.ui.Theme.Tokens()))
	l := b.layout()
	w, h := l.Size()
	l.Draw(gc, r.X+(r.Width-w)/2, r.Y+(r.Height-h)/2)
}

// ProvideAccessibility gives the whole count and its noun, or nothing at all
// when the badge is hidden.
func (b *Badge) ProvideAccessibility(ab *unison.AccessibilityBuilder) {
	n := ab.Node()
	if !b.Shows() {
		n.Ignored = true
		return
	}
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	n.Name = b.Phrase()
}
