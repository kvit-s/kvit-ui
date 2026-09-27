package kvitui

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/behavior"
	"github.com/richardwilkes/unison/enums/mod"
)

// Region is a body that takes the height it is given and scrolls what does
// not fit. Its scroll bars sit in strips of their own beside the content
// rather than over it, and the strips are kept whether or not a bar is
// showing: sizing the content to whether the bar shows is a loop for anything
// whose height depends on its width, and content that moves aside the moment
// a list grows past the fold is worse.
//
// It answers the wheel with the distance the wheel was turned, the keys every
// desktop scrolling view answers (Page Up, Page Down, Home, End and the
// arrows) when nothing inside wanted them, and it brings whatever the
// keyboard has just reached into view.
type Region struct {
	unison.Panel
	ui *UI
	// Padding is the space inside the scrolled area, around the content; the
	// view margin unless set, which is what a body directly inside a window
	// wants. A region inside a card that already has padding sets Px(0).
	Padding Measure
	// Horizontal lets the content be wider than the region and scroll
	// sideways too.
	Horizontal bool

	scroll     *unison.ScrollPanel
	holder     *unison.Panel
	vbar, hbar *unison.ScrollBar
	sideways   bool // whether the content was set up to scroll sideways
	ease       wheelScroll
}

// NewRegion returns a region scrolling content, with the view margin around
// it.
func NewRegion(ui *UI, content unison.Paneler) *Region {
	r := &Region{ui: ui, Padding: SizeViewMargin}
	r.Self = r
	r.holder = unison.NewPanel()
	r.holder.SetBorder(regionPadding{r})
	r.holder.SetLayout(&unison.FlexLayout{Columns: 1})
	content.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
	r.holder.AddChild(content)
	r.holder.FocusChangeInHierarchyCallback = func(_, to *unison.Panel) { r.reveal(to) }

	r.scroll = unison.NewScrollPanel()
	r.scroll.SetContent(r.holder, behavior.Follow, behavior.Fill)
	r.scroll.DrawCallback = func(*unison.Canvas, geom.Rect) {}
	for _, own := range []*unison.ScrollBar{r.scroll.Bar(false), r.scroll.Bar(true)} {
		own.Hidden = true
		own.MinimumThickness = 0
	}
	r.ease = wheelScroll{ui: ui, at: r.scroll.Position, moveTo: r.scroll.SetPosition,
		travel: func() (x, y float32) { return r.scroll.Bar(true).MaxValue(), r.scroll.Bar(false).MaxValue() }}
	r.scroll.MouseWheelCallback = r.wheel
	r.scroll.KeyDownCallback = r.keyDown
	r.vbar = NewScrollBar(ui, false)
	r.hbar = NewScrollBar(ui, true)
	r.follow(r.scroll.Bar(false), r.vbar, false)
	r.follow(r.scroll.Bar(true), r.hbar, true)

	r.AddChild(r.scroll)
	r.AddChild(r.vbar)
	r.AddChild(r.hbar)
	r.SetLayout(regionLayout{r})
	return r
}

// follow keeps a visible bar and the scroll panel's hidden one in step, in
// both directions.
func (r *Region) follow(own, shown *unison.ScrollBar, horizontal bool) {
	previous := own.ChangedCallback
	own.ChangedCallback = func() {
		if previous != nil {
			previous()
		}
		shown.SetRange(own.Value(), own.Extent(), own.Max())
	}
	shown.ChangedCallback = func() {
		h, v := r.scroll.Position()
		if horizontal {
			h = shown.Value()
		} else {
			v = shown.Value()
		}
		r.scroll.SetPosition(h, v)
	}
}

// Scrolls reports whether anything in the region is out of sight, which a
// caller may need in order to say so.
func (r *Region) Scrolls() bool {
	b := r.scroll.Bar(false)
	return b.Max() > b.Extent()+1
}

// Position is how far the content is scrolled, across and down.
func (r *Region) Position() (x, y float32) { return r.scroll.Position() }

// ScrollTo scrolls the content so y is at the top of the view, kept within
// the content's ends.
func (r *Region) ScrollTo(y float32) {
	x, _ := r.scroll.Position()
	r.scroll.SetPosition(x, y)
}

// view is the height of the part of the content that shows.
func (r *Region) view() float32 { return r.scroll.Bar(false).Extent() }

// wheel moves the view the distance the wheel was turned; see wheelScroll.
func (r *Region) wheel(_, delta geom.Point, _ mod.Modifiers) bool { return r.ease.wheel(delta) }

// keyDown scrolls for the keys nothing inside the region wanted: a key goes
// to the focused panel first and comes here only if it was not used, so a
// field keeps Home and End for its caret and a list keeps its arrows.
func (r *Region) keyDown(key unison.KeyCode, _ mod.Modifiers, _ bool) bool {
	if !r.Scrolls() {
		return false
	}
	m := r.ui.Interface
	_, y := r.scroll.Position()
	page := max(float32(m.RowHeight()), r.view()-float32(m.RowHeight()))
	switch key {
	case unison.KeyPageDown:
		y += page
	case unison.KeyPageUp:
		y -= page
	case unison.KeyHome:
		y = 0
	case unison.KeyEnd:
		y = r.scroll.Bar(false).MaxValue()
	case unison.KeyDown:
		y += float32(m.RowHeightSlim())
	case unison.KeyUp:
		y -= float32(m.RowHeightSlim())
	default:
		return false
	}
	r.ScrollTo(y)
	return true
}

// reveal brings a panel that has just taken the focus inside the region into
// view, a space clear of the edge: a focus ring below the fold is a focus
// ring nobody can see.
func (r *Region) reveal(p *unison.Panel) {
	if p == nil || !r.Scrolls() {
		return
	}
	inside := false
	for a := p; a != nil; a = a.Parent() {
		if a == r.holder {
			inside = true
			break
		}
	}
	if !inside {
		return
	}
	box := r.holder.RectFromRoot(p.RectToRoot(p.ContentRect(true)))
	_, y := r.scroll.Position()
	space := float32(r.ui.Interface.Space())
	switch {
	case box.Bottom() > y+r.view():
		r.ScrollTo(box.Bottom() - r.view() + space)
	case box.Y < y:
		r.ScrollTo(box.Y - space)
	}
}

// regionPadding is the padding around the content, inside the scrolled area.
type regionPadding struct{ r *Region }

func (p regionPadding) Insets() geom.Insets {
	if p.r.Padding == nil {
		return geom.Insets{}
	}
	return geom.NewUniformInsets(float32(p.r.Padding.Of(p.r.ui)))
}

func (p regionPadding) Draw(*unison.Canvas, geom.Rect) {}

// regionLayout puts the scrolled area beside the strip its vertical bar
// takes, and above the horizontal bar's strip when it scrolls sideways.
type regionLayout struct{ r *Region }

func (l regionLayout) strips() (right, bottom float32) {
	w := float32(l.r.ui.Interface.SpaceWide())
	if l.r.Horizontal {
		return w, w
	}
	return w, 0
}

func (l regionLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	right, bottom := l.strips()
	if l.r.Horizontal != l.r.sideways {
		l.r.sideways = l.r.Horizontal
		across := behavior.Follow
		if l.r.sideways {
			across = behavior.Fill
		}
		l.r.scroll.SetContent(l.r.holder, across, behavior.Fill)
	}
	_, p, _ := l.r.holder.Sizes(geom.NewSize(max(0, hint.Width-right), 0))
	return geom.NewSize(right, bottom), geom.NewSize(p.Width+right, p.Height+bottom),
		geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l regionLayout) PerformLayout(target *unison.Panel) {
	right, bottom := l.strips()
	b := target.ContentRect(false)
	area := geom.NewRect(b.X, b.Y, max(0, b.Width-right), max(0, b.Height-bottom))
	l.r.scroll.SetFrameRect(area)
	l.r.vbar.SetFrameRect(geom.NewRect(area.Right(), b.Y, right, area.Height))
	l.r.hbar.Hidden = !l.r.Horizontal
	l.r.hbar.SetFrameRect(geom.NewRect(b.X, area.Bottom(), area.Width, bottom))
	// The scroll panel sets its bars' ranges as it lays out; the visible
	// bars take theirs from those.
	l.r.scroll.ValidateLayout()
	for _, pair := range [][2]*unison.ScrollBar{{l.r.scroll.Bar(false), l.r.vbar}, {l.r.scroll.Bar(true), l.r.hbar}} {
		pair[1].SetRange(pair[0].Value(), pair[0].Extent(), pair[0].Max())
	}
}
