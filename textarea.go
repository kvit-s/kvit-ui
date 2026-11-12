package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/role"
)

// TextArea is several lines of text the reader types, or reads and selects: a
// message, a note, a commit description, a file's own source. Its outline
// and its error are a Field's, for the same reasons; like a Field, the
// editing is unison's.
type TextArea struct {
	Field
	// Plain draws a document filling a pane: no ground and no outline, since
	// a whole-pane rectangle drawn as a field says there is something beside
	// it. It still draws the focus ring and takes the error message.
	Plain bool
	// Mono draws the text in the monospace family, for source text and
	// anything where a column has to line up.
	Mono bool
	// Underlay draws behind the words, such as a wash behind a marked
	// passage. It is given a function answering where a rune index is: the
	// top-left of the caret before it and the line's height.
	Underlay func(gc *unison.Canvas, at func(index int) geom.Rect)
	// Fit makes the area ask for exactly the size of its text, for read-only
	// text shown whole inside a scrolling region; FitNone asks for three
	// rows whatever it holds.
	Fit TextFit
}

// TextFit is how a text area's size follows its text.
type TextFit int

const (
	// FitNone asks for three rows to start, whatever the area holds.
	FitNone TextFit = iota
	// FitText asks for the text's height at the width offered, and for the
	// text's width when the text does not wrap.
	FitText
	// FitTextAtLeastOffered asks for the size of the text on unwrapped
	// lines, and for the width offered when that is wider.
	FitTextAtLeastOffered
)

// fitSizes are the sizes of an area that fits its text, the field's error
// message under it included.
func (a *TextArea) fitSizes(hint geom.Size, message float32) (minSize, prefSize, maxSize geom.Size) {
	e := a.edit
	var w, h float32
	if a.Fit == FitTextAtLeastOffered {
		_, p, _ := e.Sizes(geom.Size{})
		w, h = max(p.Width, hint.Width), p.Height
	} else {
		width := hint.Width
		if !e.Wrap() {
			width = 0
		}
		_, p, _ := e.Sizes(geom.NewSize(width, 0))
		h = p.Height
		if !e.Wrap() {
			w = p.Width
		}
	}
	h += message
	return geom.NewSize(w, h), geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, h)
}

// NewTextArea returns an empty text area, three rows tall to start.
func NewTextArea(ui *UI) *TextArea {
	a := &TextArea{}
	a.ui = ui
	a.Self = a
	a.area = a
	a.initField(ui, true)
	named := a.edit.Accessibility.Callback
	a.edit.Accessibility.Callback = func(n *accessibility.Node) {
		named(n)
		n.Role = role.TextArea
	}
	return a
}
