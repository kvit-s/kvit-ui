package kvitui

import (
	"strings"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// Label is a run of chrome text at one of the seven type roles. Every other
// component uses it, which keeps the chrome family, the colour and the
// eliding rule in one place. A label wider than its column is cut short with
// "…", because one that grows pushes whatever is beside it off the screen.
type Label struct {
	unison.Panel
	ui *UI
	// Text is what the label says.
	Text string
	// Role is what the text is; its size follows from it.
	Role TypeRole
	// Mono sets the text in the monospace family, for an identifier that has
	// to line up down a column.
	Mono bool
	// Tabular gives every digit the same width, so a column of figures lines
	// up and a changing value does not shift the text beside it.
	Tabular bool
	// Weight is text.Regular unless set.
	Weight int
	// Ink is the colour; InkTextPrimary unless set.
	Ink Ink
	// Wrap lets the text run onto further lines instead of being cut short.
	Wrap bool
	// LineHeight multiplies the lines' height, for a paragraph read in a
	// popover; 0 is 1.
	LineHeight float32
}

// NewLabel returns a label in the body role.
func NewLabel(ui *UI, s string) *Label {
	l := &Label{ui: ui, Text: s, Role: RoleBody}
	l.Self = l
	l.SetSizer(l.sizes)
	l.DrawCallback = l.draw
	return l
}

func (l *Label) style() text.Style {
	ink := l.Ink
	if ink == nil {
		ink = InkTextPrimary
	}
	weight := l.Weight
	if weight == 0 {
		weight = text.Regular
	}
	var st text.Style
	if l.Mono {
		st = l.ui.Mono(l.ui.Size(l.Role), ink.Of(l.ui))
	} else {
		st = l.ui.Chrome(l.ui.Size(l.Role), weight, ink.Of(l.ui))
	}
	st.Tabular = l.Tabular
	return st
}

func (l *Label) layout(width float32) *text.Layout {
	opt := text.Options{MaxWidth: width, LineHeight: l.LineHeight}
	if !l.Wrap {
		opt.Elide = width > 0
	}
	return l.ui.Fonts.Layout([]text.Span{{Text: l.Text, Style: l.style()}}, opt)
}

func (l *Label) sizes(hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	var in geom.Insets
	if b := l.Border(); b != nil {
		in = b.Insets()
	}
	w, h := l.layout(0).Size()
	if l.Wrap {
		// Wrapping text takes the width it is given and grows downwards; it
		// asks for no width of its own, or one long sentence would widen the
		// whole column it sits in.
		if room := hint.Width - in.Width(); hint.Width > 0 && room < w {
			_, h = l.layout(max(0, room)).Size()
		}
		return geom.NewSize(in.Width(), h+in.Height()), geom.NewSize(in.Width(), h+in.Height()),
			geom.NewSize(unison.DefaultMaxSize, h+in.Height())
	}
	return geom.NewSize(in.Width(), h+in.Height()), geom.NewSize(w+in.Width(), h+in.Height()),
		geom.NewSize(unison.DefaultMaxSize, h+in.Height())
}

func (l *Label) draw(gc *unison.Canvas, _ geom.Rect) {
	b := l.ContentRect(false)
	lay := l.layout(b.Width)
	_, h := lay.Size()
	lay.Draw(gc, b.X, b.Y+(b.Height-h)/2) // vertically centred, as the Qt label is
}

// ProvideAccessibility describes the label to screen readers as static text.
func (l *Label) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	if n.Name == "" {
		n.Name = l.Text
	}
}

// wordsIn is the words the labels and figures inside a panel say, in order,
// for a panel that needs a name and has none of its own.
func wordsIn(p *unison.Panel) string {
	var words []string
	var walk func(q *unison.Panel)
	walk = func(q *unison.Panel) {
		for _, c := range q.Children() {
			if c.Hidden {
				continue
			}
			switch v := c.Self.(type) {
			case *Label:
				if v.Text != "" {
					words = append(words, v.Text)
				}
			case *Figure:
				words = append(words, v.Phrase())
			default:
				walk(c)
			}
		}
	}
	walk(p)
	return strings.Join(words, ", ")
}
