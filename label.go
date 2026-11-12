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
	// Italic sets the text in italics, in the chrome and the monospace
	// family alike.
	Italic bool
	// MaxLines, with Wrap, cuts the text short with "…" after this many
	// lines; 0 is no limit.
	MaxLines int
	// Middle cuts one line short in its middle rather than at its end,
	// keeping the text's start and its end. It applies where the text is
	// not wrapped.
	Middle bool
	// Align places each line across the label's width: at its start (the
	// default), in its middle or at its end.
	Align text.Alignment
	// Name is what a screen reader is told; the text when "".
	Name string
	// Heading describes the label to screen readers as a heading.
	Heading bool
	// Ignored keeps the label from screen readers.
	Ignored bool
	// NoTooltip leaves a label that is cut short without the tooltip of its
	// whole text.
	NoTooltip bool
	// EmptyTakesNoRoom makes a label with no text ask for no size, draw
	// nothing and stay out of a screen reader's way, so a line vanishes when
	// it is emptied. Without it an empty label is one line tall.
	EmptyTakesNoRoom bool
	// Top, with Wrap, draws the text from the top of the label's box, where
	// it is otherwise centred on the box's height.
	Top bool
	// WeightedMono, with Mono, applies Weight, which the monospace family
	// otherwise ignores.
	WeightedMono bool
	// WrapOwnWidth, with Wrap, makes the label ask for its text's width, up
	// to the width offered, where a wrapping label otherwise asks for no
	// width of its own.
	WrapOwnWidth bool
	// Size, when set, is the text's size in place of the Role's. Reach for a
	// Role first: a size of its own is for text that follows something
	// other than the seven roles.
	Size Measure
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
		if !l.NoTooltip && l.Elided() {
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
	size := l.ui.Size(l.Role)
	if l.Size != nil {
		size = l.Size.Of(l.ui)
	}
	var st text.Style
	if l.Mono {
		st = l.ui.Mono(size, ink.Of(l.ui))
		if l.WeightedMono {
			st.Weight = weight
		}
	} else {
		st = l.ui.Chrome(size, weight, ink.Of(l.ui))
	}
	st.Tabular = l.Tabular
	st.Italic = st.Italic || l.Italic
	return st
}

func (l *Label) layout(width float32) *text.Layout {
	opt := text.Options{MaxWidth: width, LineHeight: l.LineHeight, Align: l.Align}
	if !l.Wrap {
		opt.Elide = width > 0
		opt.ElideMiddle = l.Middle
	}
	return l.ui.Fonts.Layout([]text.Span{{Text: l.Text, Style: l.style()}}, opt)
}

// Layout is the label's text as it is laid out in a width, 0 for one
// unbroken line, for a component that draws over part of it, such as an
// underline under a few words. It is the text before MaxLines cuts it short.
func (l *Label) Layout(width float32) *text.Layout { return l.layout(width) }

// limited is the label's text in a width, in two parts where MaxLines cuts it
// short: the lines kept whole, nil when only one line is kept, and under them
// the rest of the text as one line cut short. Otherwise it is the whole text
// and nil.
func (l *Label) limited(width float32) (body, last *text.Layout) {
	body = l.layout(width)
	if !l.Wrap || l.MaxLines <= 0 || width <= 0 || body.LineCount() <= l.MaxLines {
		return body, nil
	}
	start, _, _, _ := body.LineBounds(l.MaxLines - 1)
	runes := []rune(l.Text)
	rest := strings.Join(strings.Fields(string(runes[start:])), " ")
	opt := text.Options{MaxWidth: width, LineHeight: l.LineHeight, Align: l.Align}
	var head *text.Layout
	if l.MaxLines > 1 {
		head = l.ui.Fonts.Layout([]text.Span{{Text: string(runes[:start]), Style: l.style()}}, opt)
	}
	opt.Elide = true
	return head, l.ui.Fonts.Layout([]text.Span{{Text: rest, Style: l.style()}}, opt)
}

// limitedHeight is the height of the text in a width, as limited gives it.
func (l *Label) limitedHeight(width float32) float32 {
	body, last := l.limited(width)
	if last == nil {
		_, h := body.Size()
		return h
	}
	_, h := last.Size()
	if body != nil {
		_, bh := body.Size()
		h += bh
	}
	return h
}

// SetText changes what the label says and, when that changed, has what it
// sits in laid out again and the label drawn.
func (l *Label) SetText(s string) {
	if l.Text == s {
		return
	}
	l.Text = s
	l.MarkForLayoutRecursivelyUpward()
	l.MarkForRedraw()
}

func (l *Label) sizes(hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	var in geom.Insets
	if b := l.Border(); b != nil {
		in = b.Insets()
	}
	if l.EmptyTakesNoRoom && l.Text == "" {
		return geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, 0)
	}
	w, h := l.layout(0).Size()
	if l.Wrap {
		// Wrapping text takes the width it is given and grows downwards; it
		// asks for no width of its own, or one long sentence would widen the
		// whole column it sits in.
		room := hint.Width - in.Width()
		if hint.Width > 0 && room < w {
			h = l.limitedHeight(max(0, room))
		}
		own := float32(0)
		if l.WrapOwnWidth {
			own = w
			if hint.Width > 0 {
				own = min(w, max(0, room))
			}
		}
		return geom.NewSize(in.Width(), h+in.Height()), geom.NewSize(own+in.Width(), h+in.Height()),
			geom.NewSize(unison.DefaultMaxSize, h+in.Height())
	}
	return geom.NewSize(in.Width(), h+in.Height()), geom.NewSize(w+in.Width(), h+in.Height()),
		geom.NewSize(unison.DefaultMaxSize, h+in.Height())
}

func (l *Label) draw(gc *unison.Canvas, _ geom.Rect) {
	if l.EmptyTakesNoRoom && l.Text == "" {
		return
	}
	b := l.ContentRect(false)
	lay, last := l.limited(b.Width)
	var w, h, headHeight float32
	if lay != nil {
		w, h = lay.Size()
		headHeight = h
	}
	if last != nil {
		lw, lh := last.Size()
		w, h = max(w, lw), h+lh
	}
	y := b.Y + (b.Height-h)/2 // vertically centred
	if l.Top && l.Wrap {
		y = b.Y
	}
	if l.Selectable && l.Focused() {
		// Selected whole, in the colour of selected text.
		painterFor(gc, l.ui).fill(geom.NewRect(b.X, y, min(w, b.Width), h), l.ui.Theme.Tokens().SelectionActiveTint)
	}
	if lay != nil {
		lay.Draw(gc, b.X, y)
	}
	if last != nil {
		last.Draw(gc, b.X, y+headHeight)
	}
}

// ProvideAccessibility describes the label to screen readers as static text.
func (l *Label) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Label
		if l.Heading {
			n.Role = role.Heading
		}
	}
	if l.Name != "" {
		n.Name = l.Name
	}
	if n.Name == "" {
		n.Name = l.Text
	}
	if l.Ignored || (l.EmptyTakesNoRoom && n.Name == "") {
		n.Ignored = true
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
