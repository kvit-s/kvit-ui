package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/role"
)

// EmptyForm is how much room an empty state takes.
type EmptyForm int

const (
	// EmptyFull is the block for an empty pane: a symbol, a sentence, a line
	// under it and the way to fill it, centred.
	EmptyFull EmptyForm = iota
	// EmptyCompact is the same words on one line at a slim row's height,
	// starting where rows start, for an empty section in a stack of them.
	EmptyCompact
)

// EmptyState is what a view says when it has nothing to show: what would be
// here and, where there is one, the action that would put something there.
// A blank region says nothing about whether the data is loading, the filter
// is too narrow or nothing has been made yet. It is also the answer for a
// chart with no data, since an axis drawn around zeros says the values are
// zero, which is a different claim from having none.
//
// Dashed makes it the place a file is dropped. The drop is the
// application's; what belongs here is the dashed edge and the action, since
// there is no keyboard gesture for dropping a file and a drop target with no
// other way in is unusable without a mouse.
type EmptyState struct {
	unison.Panel
	ui *UI
	// Form is the full block or the compact line.
	Form EmptyForm
	// Title says what would be here, in one short sentence.
	Title string
	// Detail says why it is not, or what to do about it; optional.
	Detail string
	// Symbol is an optional meaning name.
	Symbol string
	// Action is the words of the action that would fill it; "" for none. It
	// is a button in the full form and a link in the compact one.
	Action string
	// Dashed draws a dashed edge, for a region something can be put into.
	Dashed bool
	// OnAction runs when the action is taken.
	OnAction func()

	button *Button
	link   *Link
}

// NewEmptyState returns a full empty state saying title.
func NewEmptyState(ui *UI, title string) *EmptyState {
	e := &EmptyState{ui: ui, Title: title}
	e.Self = e
	act := func() {
		if e.OnAction != nil {
			e.OnAction()
		}
	}
	e.button = NewButton(ui, "")
	e.button.OnClick = act
	e.link = NewLink(ui, "")
	e.link.Role = RoleSmall
	e.link.OnActivate = act
	e.AddChild(e.button)
	e.AddChild(e.link)
	e.SetLayout(syncing{Layout: emptyLayout{e}, sync: e.sync})
	e.DrawCallback = e.draw
	return e
}

func (e *EmptyState) sync() {
	e.button.Text, e.link.Text = e.Action, e.Action
	e.button.Hidden = e.Action == "" || e.Form != EmptyFull
	e.link.Hidden = e.Action == "" || e.Form != EmptyCompact
}

// fullParts are where the parts of the full block go, and its height.
type fullParts struct {
	icon, title, detail, button geom.Rect
	height                      float32
}

func (e *EmptyState) layouts(detailWidth float32) (title, detail *text.Layout) {
	ui, t := e.ui, e.ui.Theme.Tokens()
	title = ui.Fonts.Layout([]text.Span{{Text: e.Title, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, t.TextSecondary)}}, text.Options{})
	if e.Detail != "" {
		detail = ui.Fonts.Layout([]text.Span{{Text: e.Detail, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, t.TextMuted)}},
			text.Options{MaxWidth: detailWidth, Align: text.AlignMiddle})
	}
	return title, detail
}

// detailWidth is the width the detail wraps at: seven tenths of the width,
// and never more than 320 design pixels.
func (e *EmptyState) detailWidth(width float32) float32 {
	return min(width*0.7, float32(e.ui.Interface.Px(320)))
}

func (e *EmptyState) arrangeFull(width float32) fullParts {
	m := e.ui.Interface
	pad, gap := float32(m.ViewMargin()), float32(m.Space())
	var out fullParts
	y := pad
	centred := func(w, h float32) geom.Rect {
		r := geom.NewRect((width-w)/2, y, w, h)
		y += h + gap
		return r
	}
	if e.Symbol != "" {
		s := float32(m.Px(28))
		out.icon = centred(s, s)
	}
	dw := e.detailWidth(width)
	title, detail := e.layouts(dw)
	out.title = centred(title.Size())
	if detail != nil {
		_, h := detail.Size()
		out.detail = centred(dw, h)
	}
	if e.Action != "" {
		_, p, _ := e.button.Sizes(geom.Size{})
		out.button = centred(p.Width, p.Height)
	}
	out.height = y - gap + pad
	return out
}

// compactParts are where the parts of the one line go.
type compactParts struct {
	icon, title, dot, detail, link geom.Rect
	width                          float32
}

func (e *EmptyState) arrangeCompact(b geom.Rect) compactParts {
	ui, m := e.ui, e.ui.Interface
	gap := float32(m.SpaceNear())
	small := func(s string) float32 {
		w, _ := ui.Fonts.Layout([]text.Span{{Text: s, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, ui.Theme.Tokens().TextPrimary)}}, text.Options{}).Size()
		return w
	}
	var out compactParts
	// The natural widths, then what the title and detail give up to fit:
	// the symbol and the action keep their size, since a symbol at half
	// width is a smudge and a half-drawn action is a control nobody can read.
	var iw, tw, dotW, dw, lw float32
	if e.Symbol != "" {
		iw = float32(m.IconSizeSmall())
	}
	tw = small(e.Title)
	if e.Detail != "" && e.Title != "" {
		dotW = small("·")
	}
	if e.Detail != "" {
		dw = small(e.Detail)
	}
	if e.Action != "" {
		_, p, _ := e.link.Sizes(geom.Size{})
		lw = p.Width
	}
	parts := []float32{iw, tw, dotW, dw}
	n, total := 0, float32(0)
	for _, w := range append(parts, lw) {
		if w > 0 {
			total += w
			n++
		}
	}
	if n > 1 {
		total += gap * float32(n-1)
	}
	if lw > 0 {
		total += float32(m.SpaceSnug())
	}
	out.width = total + 2*float32(m.Space())
	if over := total - (b.Width - 2*float32(m.Space())); over > 0 && tw+dw > 0 {
		cut := min(1, over/(tw+dw))
		tw -= tw * cut
		dw -= dw * cut
	}
	x := b.X + float32(m.Space())
	place := func(w float32, extra float32) geom.Rect {
		if w <= 0 {
			return geom.Rect{}
		}
		x += extra
		r := geom.NewRect(x, b.Y, w, b.Height)
		x += w + gap
		return r
	}
	out.icon = place(iw, 0)
	out.title = place(tw, 0)
	out.dot = place(dotW, 0)
	out.detail = place(dw, 0)
	out.link = place(lw, float32(m.SpaceSnug()))
	return out
}

type emptyLayout struct{ e *EmptyState }

func (l emptyLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	e, m := l.e, l.e.ui.Interface
	if e.Form == EmptyCompact {
		at := e.arrangeCompact(geom.NewRect(0, 0, max(hint.Width, 1e6), 0))
		_, lh := e.ui.Fonts.Layout([]text.Span{{Text: "Hg", Style: e.ui.Chrome(e.ui.Size(RoleSmall), text.Regular, e.ui.Theme.Tokens().TextPrimary)}}, text.Options{}).Size()
		h := max(float32(m.RowHeightSlim()), lh+2*float32(m.SpaceSnug()))
		return geom.NewSize(0, h), geom.NewSize(at.width, h), geom.NewSize(unison.DefaultMaxSize, h)
	}
	w := float32(m.Px(400))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := e.arrangeFull(w).height
	return geom.NewSize(0, h), geom.NewSize(float32(m.Px(400)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l emptyLayout) PerformLayout(target *unison.Panel) {
	e := l.e
	b := target.ContentRect(false)
	if e.Form == EmptyCompact {
		e.link.SetFrameRect(e.centreOn(e.arrangeCompact(b).link, e.link))
		return
	}
	at := e.arrangeFull(b.Width)
	e.button.SetFrameRect(geom.NewRect(b.X+at.button.X, b.Y+at.button.Y, at.button.Width, at.button.Height))
}

// centreOn is a panel's preferred height, centred in a slot of the line.
func (e *EmptyState) centreOn(slot geom.Rect, p unison.Paneler) geom.Rect {
	_, pref, _ := p.AsPanel().Sizes(geom.Size{})
	return geom.NewRect(slot.X, slot.Y+(slot.Height-pref.Height)/2, slot.Width, pref.Height)
}

func (e *EmptyState) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := e.ui, e.ui.Theme.Tokens(), e.ui.Interface
	b := e.ContentRect(false)
	if e.Dashed {
		line := max(1, float32(m.Hairline()))
		r := b.Inset(geom.NewUniformInsets(line / 2))
		paint := Color(t.BorderStrong).Paint(gc, r, paintstyle.Stroke)
		paint.SetStrokeWidth(line)
		paint.SetPathEffect(unison.NewDashPathEffect([]float32{float32(m.Px(5)), float32(m.Px(4))}, 0))
		radius := float32(m.RadiusCard())
		gc.DrawRoundedRect(r, geom.NewSize(radius, radius), paint)
	}
	symbol := func(size float32, box geom.Rect) {
		if g, ok := icons.Glyph(e.Symbol); ok && e.Symbol != "" {
			drawGlyph(gc, ui, g, size, t.TextFaint, box)
		}
	}
	if e.Form == EmptyCompact {
		at := e.arrangeCompact(b)
		symbol(float32(m.IconSizeSmall()), at.icon)
		put := func(s string, ink Ink, box geom.Rect) {
			if s == "" || box.Width <= 0 {
				return
			}
			l := ui.Fonts.Layout([]text.Span{{Text: s, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, ink.Of(ui))}},
				text.Options{MaxWidth: box.Width + 0.5, Elide: true})
			_, h := l.Size()
			l.Draw(gc, box.X, box.Y+(box.Height-h)/2)
		}
		put(e.Title, InkTextSecondary, at.title)
		if at.dot.Width > 0 {
			put("·", InkTextFaint, at.dot)
		}
		put(e.Detail, InkTextMuted, at.detail)
		return
	}
	at := e.arrangeFull(b.Width)
	at.icon.X += b.X
	at.icon.Y += b.Y
	// A box 28 design pixels square, with the glyph no larger than the icon
	// size, as the icon caps it.
	symbol(min(at.icon.Width, float32(m.IconSize())), at.icon)
	title, detail := e.layouts(e.detailWidth(b.Width))
	title.Draw(gc, b.X+at.title.X, b.Y+at.title.Y)
	if detail != nil {
		detail.Draw(gc, b.X+at.detail.X, b.Y+at.detail.Y)
	}
}

// ProvideAccessibility reads the empty state as its sentence, "No
// transactions yet. Import a statement…", holding its action.
func (e *EmptyState) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name = e.Title
	if e.Detail != "" {
		n.Name = e.Title + ". " + e.Detail
	}
}
