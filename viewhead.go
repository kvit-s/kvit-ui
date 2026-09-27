package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// ViewHead is the strip at the top of a view: what the view is, how much is
// in it, and the controls that act on the whole of it. It belongs to whatever
// the window is showing, where the window's own header belongs to the
// window, and it is where a filter, a sort control and a count go, so a
// reader looking for how to narrow a view down looks in the same place on
// every screen.
type ViewHead struct {
	unison.Panel
	ui *UI
	// Title names the view, and is the heading a screen reader announces.
	Title string
	// Count is how many things the view holds; below zero shows no count.
	Count int
	// Counted is the word for one of the things counted, "item" unless set.
	// A count with no noun beside it does not say what was counted.
	Counted string
	// CountedPlural is the word for several; "" adds an s to Counted.
	CountedPlural string
	// Subtitle is a sentence under the title. It wraps where a row label is
	// cut short, because cutting the end off an explanation loses the half
	// that explains.
	Subtitle string
	// Padding is how far the content sits in from the left and right edges
	// of the region the head belongs to; the view margin unless set.
	// Everything else in a region is inset by that margin, and a head that
	// left it to each caller sat against the edge wherever a caller forgot.
	Padding Measure

	title, count, subtitle *Label
	controls               *unison.Panel
}

// NewViewHead returns a head with a title and, at its right, the controls
// that act on the whole view, a space apart.
func NewViewHead(ui *UI, title string, controls ...unison.Paneler) *ViewHead {
	h := &ViewHead{ui: ui, Title: title, Count: -1, Counted: "item", Padding: SizeViewMargin}
	h.Self = h
	h.title = NewLabel(ui, "")
	h.title.Role, h.title.Weight = RoleTitle, text.Bold
	h.count = NewLabel(ui, "")
	h.count.Role, h.count.Ink, h.count.Tabular = RoleSmall, InkTextFaint, true
	h.subtitle = NewLabel(ui, "")
	h.subtitle.Role, h.subtitle.Ink, h.subtitle.Wrap = RoleSmall, InkTextMuted, true
	h.controls = Row(ui, SizeSpace, controls...)
	h.SetLayout(viewHeadLayout{h})
	h.sync()
	return h
}

// CountPhrase is the count as the head shows it, "1,284 transactions", or ""
// when there is no count.
func (h *ViewHead) CountPhrase() string {
	return h.ui.CountPhrase(h.Count, h.Counted, h.CountedPlural)
}

// sync copies the fields into the labels, and keeps as children only the
// parts that have something to say, so a screen reader is not handed an empty
// count or subtitle.
func (h *ViewHead) sync() {
	h.title.Text = h.Title
	h.count.Text = h.CountPhrase()
	h.subtitle.Text = h.Subtitle
	want := []unison.Paneler{h.title}
	if h.count.Text != "" {
		want = append(want, h.count)
	}
	if h.subtitle.Text != "" {
		want = append(want, h.subtitle)
	}
	if len(h.controls.Children()) > 0 {
		want = append(want, h.controls)
	}
	have := h.Children()
	same := len(have) == len(want)
	for i := 0; same && i < len(want); i++ {
		same = have[i] == want[i].AsPanel()
	}
	if !same {
		h.RemoveAllChildren()
		for _, p := range want {
			h.AddChild(p)
		}
	}
}

// ProvideAccessibility announces the head as the view's heading.
func (h *ViewHead) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Heading
	}
	n.Name = h.Title
}

// viewHeadLayout puts the title, with the count beside it, and the subtitle
// under them in a column on the left, and the controls on the right. Both are
// centred on the head's height, which is at least a row.
type viewHeadLayout struct{ h *ViewHead }

type viewHeadParts struct {
	height                           float32
	title, count, subtitle, controls geom.Rect
}

func (l viewHeadLayout) arrange(width float32) viewHeadParts {
	h, m := l.h, l.h.ui.Interface
	var p viewHeadParts
	left := width
	var cp geom.Size
	if h.controls.Parent() != nil {
		_, cp, _ = h.controls.Sizes(geom.Size{})
		left = max(0, width-cp.Width-float32(m.ColumnGap()))
	}
	_, tp, _ := h.title.Sizes(geom.Size{})
	var np geom.Size
	gap := float32(0)
	if h.count.Parent() != nil {
		_, np, _ = h.count.Sizes(geom.Size{})
		gap = float32(m.Space())
	}
	// The count keeps its width; the title gives way and is cut short.
	tw := min(tp.Width, max(0, left-gap-np.Width))
	line := max(tp.Height, np.Height)
	column := line
	var sh float32
	if h.subtitle.Parent() != nil {
		_, sp, _ := h.subtitle.Sizes(geom.NewSize(left, 0))
		sh = sp.Height
		column += float32(m.SpaceTight()) + sh
	}
	p.height = max(float32(m.RowHeight()), column, cp.Height)
	y := (p.height - column) / 2
	p.title = geom.NewRect(0, y+(line-tp.Height)/2, tw, tp.Height)
	p.count = geom.NewRect(tw+gap, y+(line-np.Height)/2, np.Width, np.Height)
	p.subtitle = geom.NewRect(0, y+line+float32(m.SpaceTight()), left, sh)
	p.controls = geom.NewRect(width-cp.Width, (p.height-cp.Height)/2, cp.Width, cp.Height)
	return p
}

// insets are the border's insets and the padding at each side.
func (l viewHeadLayout) insets(target *unison.Panel) geom.Insets {
	var in geom.Insets
	if b := target.Border(); b != nil {
		in = b.Insets()
	}
	if l.h.Padding != nil {
		pad := float32(l.h.Padding.Of(l.h.ui))
		in.Left += pad
		in.Right += pad
	}
	return in
}

// LayoutSizes implements unison.Layout. The head asks for the width of its
// title, count and controls; the subtitle takes whatever width it is given.
func (l viewHeadLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	h := l.h
	h.sync()
	insets := l.insets(target)
	_, tp, _ := h.title.Sizes(geom.Size{})
	width := tp.Width
	var controls float32
	if h.count.Parent() != nil {
		_, np, _ := h.count.Sizes(geom.Size{})
		width += float32(h.ui.Interface.Space()) + np.Width
	}
	if h.controls.Parent() != nil {
		_, cp, _ := h.controls.Sizes(geom.Size{})
		controls = cp.Width + float32(h.ui.Interface.ColumnGap())
		width += controls
	}
	if hint.Width > 0 {
		width = max(0, hint.Width-insets.Width())
	}
	height := l.arrange(width).height + insets.Height()
	return geom.NewSize(controls+insets.Width(), height),
		geom.NewSize(width+insets.Width(), height),
		geom.NewSize(unison.DefaultMaxSize, height)
}

// PerformLayout implements unison.Layout.
func (l viewHeadLayout) PerformLayout(target *unison.Panel) {
	h := l.h
	h.sync()
	r := target.ContentRect(false)
	pad := float32(0)
	if h.Padding != nil {
		pad = float32(h.Padding.Of(h.ui))
	}
	r.X, r.Width = r.X+pad, max(0, r.Width-2*pad)
	p := l.arrange(r.Width)
	at := func(g geom.Rect) geom.Rect { return geom.NewRect(r.X+g.X, r.Y+g.Y, g.Width, g.Height) }
	h.title.SetFrameRect(at(p.title))
	h.count.SetFrameRect(at(p.count))
	h.subtitle.SetFrameRect(at(p.subtitle))
	h.controls.SetFrameRect(at(p.controls))
}
