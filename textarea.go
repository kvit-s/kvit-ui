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
