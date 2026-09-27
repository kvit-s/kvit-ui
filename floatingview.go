package kvitui

import (
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/role"
)

// FloatingView is detail about the one thing the reader selected, held over
// the list they selected it from rather than beside it: a centred card with
// the thing's name and a close control along its top, and the list dimmed
// behind it. A Pane takes its width out of the layout, which on a narrow
// window pushes a table's last columns off the edge (kvit-cash measured a
// ledger left 430 pixels wide at its 880-pixel minimum); a floating view takes
// no width, so the list keeps the whole of it at every window size. It is
// not a Dialog: it is closed rather than answered, floats over one region
// rather than the window, and holds whatever the caller puts in it.
//
// Nothing under it hears a press meant for it, since the dimmed area takes
// every press and the wheel. Its card is the height of the region less the
// view margin, whatever it holds scrolls inside it, and its title stays put
// above what scrolls. Escape closes it from wherever the keyboard is. It is
// kvit-cash's addition to the Qt library.
type FloatingView struct {
	unison.Panel
	ui *UI
	// Title names what is in the view.
	Title string
	// TitleDetail is the whole of what the title names where the title is
	// shortened, shown under the pointer on the title; "" says the title is
	// the whole of it.
	TitleDetail string
	// CloseLabel is what the close control is called on hover and to a
	// screen reader: "Close record" rather than "Close the view", since the
	// reader is closing the record.
	CloseLabel string
	// Closable draws the close control.
	Closable bool
	// ViewWidth is the card's width, at most the region's less the margins;
	// the floating view width unless set.
	ViewWidth Measure
	// CloseOnPressOutside closes the view on a press on the dimmed area. It
	// is off unless set, since a reader half way through typing into the
	// view should not lose it to a press that missed the card.
	CloseOnPressOutside bool
	// OnCloseRequested runs when the reader asks to close the view, by the
	// close control, Escape, or a press outside where that is on.
	OnCloseRequested func()

	body  unison.Paneler
	close *IconButton
	tip   partTip
	hide  func()
}

// NewFloatingView returns a closed view titled title, holding body.
func NewFloatingView(ui *UI, title string, body unison.Paneler) *FloatingView {
	v := &FloatingView{ui: ui, Title: title, CloseLabel: "Close", Closable: true, ViewWidth: SizeFloatingViewWidth, body: body}
	v.Self = v
	v.tip = partTip{ui: ui, owner: v}
	v.close = NewIconButton(ui, "close", "")
	v.close.OnClick = v.request
	v.AddChild(v.close)
	v.AddChild(body)
	v.SetLayout(syncing{Layout: floatingLayout{v}, sync: func() {
		v.close.Label = v.CloseLabel
		v.close.Hidden = !v.Closable
	}})
	v.DrawCallback = v.draw
	v.MouseDownCallback = func(where geom.Point, _, _ int, _ mod.Modifiers) bool {
		if v.CloseOnPressOutside && !where.In(v.card()) {
			v.request()
		}
		return true
	}
	v.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool {
		v.pointAt(where)
		return true
	}
	v.MouseExitCallback = func() bool { v.tip.clear(); return true }
	v.MouseWheelCallback = func(geom.Point, geom.Point, mod.Modifiers) bool { return true }
	return v
}

func (v *FloatingView) request() {
	if v.OnCloseRequested != nil {
		v.OnCloseRequested()
	}
}

// Open floats the view over a region, which must be in a Kvit window.
func (v *FloatingView) Open(over unison.Paneler) {
	w := v.ui.windowOf(over)
	if w == nil || v.hide != nil {
		return
	}
	v.hide = w.Show(&Popup{Panel: v, Place: func(geom.Rect, geom.Size) geom.Rect { return anchorIn(over) },
		OnEscape: v.request, Anchor: over})
}

// Opened reports whether the view is showing.
func (v *FloatingView) Opened() bool { return v.hide != nil }

// Close takes the view away.
func (v *FloatingView) Close() {
	v.tip.clear()
	if v.hide != nil {
		v.hide()
		v.hide = nil
	}
}

// card is the card's box, centred in the region, a view margin in from its
// top and bottom.
func (v *FloatingView) card() geom.Rect {
	m := v.ui.Interface
	r := v.ContentRect(false)
	margin := float32(m.ViewMargin())
	w := min(float32(v.ViewWidth.Of(v.ui)), r.Width-2*margin)
	return geom.NewRect(r.X+(r.Width-w)/2, r.Y+margin, max(0, w), max(0, r.Height-2*margin))
}

// framed is the card less its hairline edge: what the card holds is inside
// the edge rather than drawn over it.
func (v *FloatingView) framed() geom.Rect {
	return v.card().Inset(geom.NewUniformInsets(float32(v.ui.Interface.Hairline())))
}

// head is the row holding the title and the close control.
func (v *FloatingView) head() geom.Rect {
	f := v.framed()
	return geom.NewRect(f.X, f.Y, f.Width, float32(v.ui.Interface.BreadcrumbHeight()))
}

func (v *FloatingView) titleBox() geom.Rect {
	m := v.ui.Interface
	h := v.head()
	right := h.Right() - float32(m.SpaceNear())
	if v.Closable {
		right -= preferred(v.close).Width + float32(m.Space())
	}
	x := h.X + float32(m.ViewMargin())
	return geom.NewRect(x, h.Y, max(0, right-x), h.Height)
}

func (v *FloatingView) titleLayout(width float32) *text.Layout {
	ui := v.ui
	return ui.Fonts.Layout([]text.Span{{Text: v.Title, Style: ui.Chrome(ui.Size(RoleStrong), text.Bold, ui.Theme.Tokens().TextPrimary)}},
		text.Options{MaxWidth: width, Elide: width > 0})
}

type floatingLayout struct{ v *FloatingView }

func (l floatingLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	return geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l floatingLayout) PerformLayout(*unison.Panel) {
	v, m := l.v, l.v.ui.Interface
	h := v.head()
	size := preferred(v.close)
	v.close.SetFrameRect(geom.NewRect(h.Right()-float32(m.SpaceNear())-size.Width, h.Y+(h.Height-size.Height)/2, size.Width, size.Height))
	f := v.framed()
	top := h.Bottom() + float32(m.Hairline())
	v.body.AsPanel().SetFrameRect(geom.NewRect(f.X, top, f.Width, max(0, f.Bottom()-top)))
}

func (v *FloatingView) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := v.ui, v.ui.Theme.Tokens(), v.ui.Interface
	p := painterFor(gc, ui)
	r := v.ContentRect(false)
	// What is behind, dimmed with the wash a spotlight uses, at half its
	// weight: this is a view the reader opened rather than a pointer at
	// something.
	shade := Color(palette.Shade).SetAlphaIntensity(0.28)
	gc.DrawRect(r, shade.Paint(gc, r, paintstyle.Fill))
	card := v.card()
	radius := float32(m.RadiusCard())
	p.round(card, radius, t.PopupBackground)
	p.outline(card, radius, float32(m.Hairline()), t.BorderStrong)
	tb := v.titleBox()
	l := v.titleLayout(tb.Width)
	_, h := l.Size()
	l.Draw(gc, tb.X, tb.Y+(tb.Height-h)/2)
	f := v.framed()
	p.fill(geom.NewRect(f.X, v.head().Bottom(), f.Width, float32(m.Hairline())), t.Border)
}

// pointAt shows the whole title under the pointer, where the title is cut
// short or shortened.
func (v *FloatingView) pointAt(where geom.Point) {
	tb := v.titleBox()
	if !where.In(tb) {
		v.tip.clear()
		return
	}
	whole := ""
	switch {
	case v.TitleDetail != "" && v.TitleDetail != v.Title:
		whole = v.TitleDetail
	case v.titleLayout(tb.Width).Text() != v.Title:
		whole = v.Title
	}
	v.tip.tooltip("title", whole, v.titleBox)
}

// ProvideAccessibility names the view by its title.
func (v *FloatingView) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name = v.Title
	n.Description = v.TitleDetail
}
