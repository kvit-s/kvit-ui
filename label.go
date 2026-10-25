package kvitui

import (
	"strings"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// Label is a run of chrome text at one of the seven type roles. Every other
// component uses it, which keeps the chrome family, the colour and the
// eliding rule in one place. A label wider than its column is cut short with
// "…", because one that grows pushes whatever is beside it off the screen,
// and then shows its whole text as a tooltip while the pointer rests on it
// (PTR-12).
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
	// Selectable lets content text the reader may want again, a message, an
	// address, an identifier, be selected whole and copied: a press selects
	// it, and Copy, Ctrl+C or the context menu copies it (CTL-18). A
	// control's label is never selectable.
	Selectable bool
}

// NewLabel returns a label in the body role.
func NewLabel(ui *UI, s string) *Label {
	l := &Label{ui: ui, Text: s, Role: RoleBody}
	l.Self = l
	l.SetSizer(l.sizes)
	l.DrawCallback = l.draw
	l.MouseDownCallback = func(geom.Point, int, int, mod.Modifiers) bool {
		if !l.Selectable {
			return false
		}
		l.SetFocusable(true)
		l.RequestFocus()
		return true
	}
	l.GainedFocusCallback = l.MarkForRedraw
	l.LostFocusCallback = func() {
		l.SetFocusable(false)
		l.MarkForRedraw()
	}
	l.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, _ bool) bool {
		if l.Selectable && l.Focused() && Keys(unison.KeyC, ModPrimary).Matches(ui.form, key, mods) {
			unison.ClipboardSetText(l.Text)
			return true
		}
		return false
	}
	// Copy in the Edit menu copies a selected label, as it copies a field.
	l.InstallCmdHandlers(unison.CopyItemID, func(any) bool { return l.Selectable && l.Focused() },
		func(any) { unison.ClipboardSetText(l.Text) })
	ui.SetContextMenu(l, func(geom.Point) (string, []MenuItem) {
		if !l.Selectable {
			return "", nil
		}
		return "", []MenuItem{{Command: &Command{Std: StdCopy, Run: func() { unison.ClipboardSetText(l.Text) }}}}
	})
	l.UpdateTooltipCallback = func(geom.Point, geom.Rect) geom.Rect {
		l.Tooltip = nil
		if l.Elided() {
			l.Tooltip = newTooltip(ui, l.Text, "")
		}
		return l.RectToRoot(l.ContentRect(true))
	}
	return l
}

// Elided reports whether the label is cut short where it is laid out now.
func (l *Label) Elided() bool {
	if l.Wrap || l.Text == "" {
		return false
	}
	w, _ := l.layout(0).Size()
	return w > l.ContentRect(false).Width+1
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
	w, h := lay.Size()
	if l.Selectable && l.Focused() {
		// Selected whole, in the colour of selected text.
		painterFor(gc, l.ui).fill(geom.NewRect(b.X, b.Y+(b.Height-h)/2, min(w, b.Width), h), l.ui.Theme.Tokens().SelectionActiveTint)
	}
	lay.Draw(gc, b.X, b.Y+(b.Height-h)/2) // vertically centred
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
