package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
)

// box is a panel that asks for one size whatever it is offered.
func box(w, h float32) *unison.Panel {
	p := unison.NewPanel()
	p.SetSizer(func(geom.Size) (minSize, prefSize, maxSize geom.Size) {
		return geom.Size{}, geom.NewSize(w, h), geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
	})
	return p
}

// laid gives a panel a frame and lays it out.
func laid(p *unison.Panel, w, h float32) {
	p.SetFrameRect(geom.NewRect(0, 0, w, h))
	p.NeedsLayout = true
	p.ValidateLayout()
}

// layoutSession is a screen and an interface for a test of a layout, which
// needs no window of its own.
func layoutSession(t *testing.T) (*unison.HeadlessScreen, *kvitui.UI) {
	t.Helper()
	return session(t, func(*kvitui.UI) []unison.Paneler { return []unison.Paneler{unison.NewPanel()} })
}

func wantFrame(t *testing.T, what string, p *unison.Panel, x, y, w, h float32) {
	t.Helper()
	if got, want := p.FrameRect(), geom.NewRect(x, y, w, h); got != want {
		t.Errorf("%s is at %v, want %v", what, got, want)
	}
}

func wantSizes(t *testing.T, what string, p *unison.Panel, hint, minWant, prefWant, maxWant geom.Size) {
	t.Helper()
	minSize, prefSize, maxSize := p.Sizes(hint)
	if minSize != minWant || prefSize != prefWant || maxSize != maxWant {
		t.Errorf("%s asks for %v, %v, %v at %v, want %v, %v, %v", what, minSize, prefSize, maxSize, hint, minWant, prefWant, maxWant)
	}
}

// A stack leaves a hidden child out, room and gap, and does not touch its
// frame; the child comes back with its height and a gap when it is shown.
func TestAStackLeavesHiddenChildrenOut(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		gap := float32(ui.Interface.Space())
		a, b, c := box(40, 20), box(60, 30), box(50, 10)
		stack := kvitui.Stack(ui, kvitui.SizeSpace, a, b, c)
		b.Hidden = true
		b.SetFrameRect(geom.NewRect(7, 7, 7, 7))
		wantSizes(t, "a stack with its middle child hidden", stack, geom.NewSize(200, 0),
			geom.NewSize(0, 30+gap), geom.NewSize(50, 30+gap), geom.NewSize(unison.DefaultMaxSize, 30+gap))
		laid(stack, 200, 100)
		wantFrame(t, "the first child", a, 0, 0, 200, 20)
		wantFrame(t, "the hidden child", b, 7, 7, 7, 7)
		wantFrame(t, "the last child", c, 0, 20+gap, 200, 10)
		b.Hidden = false
		wantSizes(t, "the stack with every child shown", stack, geom.NewSize(200, 0),
			geom.NewSize(0, 60+2*gap), geom.NewSize(60, 60+2*gap), geom.NewSize(unison.DefaultMaxSize, 60+2*gap))
		laid(stack, 200, 100)
		wantFrame(t, "the child shown again", b, 0, 20+gap, 200, 30)
		wantFrame(t, "the last child after it", c, 0, 50+2*gap, 200, 10)

		// The gap follows the interface size.
		ui.Interface.SetFontSize(24)
		bigger := float32(ui.Interface.Space())
		if bigger == gap {
			t.Fatal("the space is the same at both interface sizes, so the test shows nothing")
		}
		laid(stack, 200, 200)
		wantFrame(t, "the second child at the larger interface size", b, 0, 20+bigger, 200, 30)
		ui.Interface.SetFontSize(12)
	})
}

func TestAStacksChildrenSayWhereTheyGo(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		start, end := kvitui.OwnWidth(box(40, 20)), kvitui.OwnWidthEnd(box(60, 30))
		// A panel with flex data that would keep it narrow in a column gets
		// the whole width in a stack.
		flex := box(50, 10)
		flex.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Start})
		inset := box(30, 10)
		inset.SetLayoutData(kvitui.StackData{Left: kvitui.Px(20), Right: kvitui.Px(10)})
		stack := kvitui.Stack(ui, nil, start, end, flex, inset)
		stack.SetBorder(kvitui.Insets(ui, kvitui.Px(5), kvitui.Px(10), kvitui.Px(5), kvitui.Px(10)))
		// The insets are counted, and the inset child counts its room.
		wantSizes(t, "a bordered stack", stack, geom.NewSize(220, 0),
			geom.NewSize(20, 80), geom.NewSize(80, 80), geom.NewSize(unison.DefaultMaxSize, 80))
		laid(stack, 220, 80)
		wantFrame(t, "a child at its own width", start.AsPanel(), 10, 5, 40, 20)
		wantFrame(t, "a child at its own width at the end", end.AsPanel(), 150, 25, 60, 30)
		wantFrame(t, "a child with flex layout data", flex, 10, 55, 200, 10)
		wantFrame(t, "a child with room kept beside it", inset, 30, 65, 170, 10)
	})
}

func TestAStackHoldsItsChildrenToAWidthAndSkipsEmptyOnes(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		narrow, wide := box(40, 20), box(60, 30)
		wide.SetLayoutData(kvitui.StackData{Wide: true})
		// Own and End are not consulted once the stack has a width limit.
		own := kvitui.OwnWidthEnd(box(10, 10))
		stack := kvitui.Stack(ui, nil, narrow, wide, own)
		stack.Layout().(*kvitui.StackLayout).MaxWidth = kvitui.Px(120)
		laid(stack, 300, 100)
		wantFrame(t, "a child held to the stack's width limit", narrow, 0, 0, 120, 20)
		wantFrame(t, "a child exempt from the limit", wide, 0, 20, 300, 30)
		wantFrame(t, "a child marked for its own width under a limit", own.AsPanel(), 0, 50, 120, 10)

		first, empty, hidden, last := box(40, 20), box(40, 0), box(40, 20), box(40, 10)
		hidden.Hidden = true
		skipping := kvitui.Stack(ui, kvitui.SizeSpace, first, empty, hidden, last)
		skipping.Layout().(*kvitui.StackLayout).SkipEmpty = true
		gap := float32(ui.Interface.Space())
		wantSizes(t, "a stack that skips its empty child", skipping, geom.NewSize(100, 0),
			geom.NewSize(0, 30+gap), geom.NewSize(40, 30+gap), geom.NewSize(unison.DefaultMaxSize, 30+gap))
		laid(skipping, 100, 100)
		wantFrame(t, "the child after the empty one", last, 0, 20+gap, 100, 10)
		wantFrame(t, "the empty child", empty, 0, 30+gap, 100, 0)
		wantFrame(t, "the hidden child of a stack that skips", hidden, 0, 30+gap, 100, 0)
	})
}

func TestALineSharesTheWidthLeftOver(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		gap := float32(ui.Interface.Space())
		fixed, hidden := kvitui.FixedWidth(box(10, 20), kvitui.Px(50)), box(70, 10)
		one, three := kvitui.Grow(box(30, 10)), kvitui.GrowShare(box(20, 40), 3)
		own := box(40, 30)
		hidden.Hidden = true
		line := kvitui.Line(ui, kvitui.SizeSpace, fixed, one, hidden, three, own)
		spare := 400 - 50 - 40 - 3*gap
		laid(line, 400, 60)
		wantFrame(t, "a child held to a width", fixed.AsPanel(), 0, 20, 50, 20)
		wantFrame(t, "a grower with one share", one.AsPanel(), 50+gap, 25, spare/4, 10)
		wantFrame(t, "a grower with three shares", three.AsPanel(), 50+2*gap+spare/4, 10, spare*3/4, 40)
		wantFrame(t, "a child at its own width", own, 400-40, 15, 40, 30)
		// A child taller than the line is cut to it.
		laid(line, 400, 30)
		wantFrame(t, "a child taller than the line", three.AsPanel(), 50+2*gap+spare/4, 0, spare*3/4, 30)
		// With no width offered a grower reports its own width; the line asks
		// for its children's own widths whatever it is offered.
		want := 50 + 30 + 20 + 40 + 3*gap
		wantSizes(t, "a line offered no width", line, geom.Size{},
			geom.NewSize(0, 40), geom.NewSize(want, 40), geom.NewSize(unison.DefaultMaxSize, 40))
		wantSizes(t, "a line offered a width", line, geom.NewSize(400, 0),
			geom.NewSize(0, 40), geom.NewSize(want, 40), geom.NewSize(unison.DefaultMaxSize, 40))
	})
}

// A wrapping label that grows makes the line as tall as its wrapped text.
func TestALineIsAsTallAsItsWrappedText(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		words := kvitui.NewLabel(ui, "A sentence long enough that it has to be wrapped onto several lines in a narrow line.")
		words.Wrap = true
		line := kvitui.Line(ui, kvitui.SizeSpace, box(40, 10), kvitui.Grow(words))
		_, oneLine, _ := words.Sizes(geom.Size{})
		_, pref, _ := line.Sizes(geom.NewSize(160, 0))
		_, wrapped, _ := words.Sizes(geom.NewSize(160-40-float32(ui.Interface.Space()), 0))
		if pref.Height != wrapped.Height || pref.Height < 2*oneLine.Height {
			t.Errorf("the line is %v tall for text %v tall when wrapped and %v on one line", pref.Height, wrapped.Height, oneLine.Height)
		}
	})
}

func TestALinesOptions(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		gap := float32(ui.Interface.Space())
		options := func(p *unison.Panel) *kvitui.LineLayout { return p.Layout().(*kvitui.LineLayout) }

		// GrowLast: the last shown child grows, the one before a hidden last.
		a, b, c := box(40, 20), box(60, 20), box(30, 20)
		c.Hidden = true
		line := kvitui.Line(ui, kvitui.SizeSpace, a, b, c)
		options(line).GrowLast = true
		laid(line, 300, 20)
		wantFrame(t, "the first child of a line whose last grows", a, 0, 0, 40, 20)
		wantFrame(t, "the last shown child of a line whose last grows", b, 40+gap, 0, 300-40-gap, 20)

		// ClipFixed: a child that does not grow is cut to the room left.
		a, b = box(80, 20), box(80, 20)
		line = kvitui.Line(ui, kvitui.SizeSpace, a, b)
		laid(line, 100, 20)
		wantFrame(t, "a child past the end of a line that does not clip", b, 80+gap, 0, 80, 20)
		options(line).ClipFixed = true
		laid(line, 100, 20)
		wantFrame(t, "a child past the end of a line that clips", b, 80+gap, 0, 100-80-gap, 20)

		// FillsHint: a line with a grower asks for the whole width offered.
		line = kvitui.Line(ui, kvitui.SizeSpace, box(40, 20), kvitui.Grow(box(10, 20)))
		if _, pref, _ := line.Sizes(geom.NewSize(300, 0)); pref.Width != 40+10+gap {
			t.Errorf("a line asks for %v wide, want its children's own %v", pref.Width, 40+10+gap)
		}
		options(line).FillsHint = true
		if _, pref, _ := line.Sizes(geom.NewSize(300, 0)); pref.Width != 300 {
			t.Errorf("a line that fills the width offered asks for %v of 300", pref.Width)
		}

		// InsetsUncounted: the sizes leave the border out; the children are
		// still placed inside it.
		a = box(40, 20)
		line = kvitui.Line(ui, kvitui.SizeSpace, a)
		line.SetBorder(kvitui.Padding(ui, kvitui.Px(10)))
		wantSizes(t, "a bordered line", line, geom.Size{},
			geom.NewSize(20, 40), geom.NewSize(60, 40), geom.NewSize(unison.DefaultMaxSize, 40))
		options(line).InsetsUncounted = true
		wantSizes(t, "a bordered line that does not count its insets", line, geom.Size{},
			geom.NewSize(0, 20), geom.NewSize(40, 20), geom.NewSize(unison.DefaultMaxSize, 20))
		laid(line, 100, 40)
		wantFrame(t, "the child of a bordered line", a, 10, 10, 40, 20)

		// Top: centred on the tallest child from the top, not on the frame.
		a, b = box(40, 20), box(40, 40)
		line = kvitui.Line(ui, kvitui.SizeSpace, a, b)
		laid(line, 200, 100)
		wantFrame(t, "a child centred on the line's height", a, 0, 40, 40, 20)
		options(line).Top = true
		laid(line, 200, 100)
		wantFrame(t, "a short child of a line laid out from the top", a, 0, 10, 40, 20)
		wantFrame(t, "the tallest child of a line laid out from the top", b, 40+gap, 0, 40, 40)

		// NoShrink: the least size is the size asked for, and a child keeps
		// its height in a shorter line.
		a = box(40, 30)
		line = kvitui.Line(ui, kvitui.SizeSpace, a, box(20, 10))
		options(line).NoShrink = true
		wantSizes(t, "a line that does not shrink", line, geom.Size{},
			geom.NewSize(60+gap, 30), geom.NewSize(60+gap, 30), geom.NewSize(unison.DefaultMaxSize, 30))
		laid(line, 200, 20)
		wantFrame(t, "a child of a line that does not shrink, in a shorter line", a, 0, -5, 40, 30)

		// AlignTop: every child at the top of the content at its own height,
		// uncut, over Top; the sizes reported are unchanged.
		a, b = box(40, 40), box(40, 20)
		line = kvitui.Line(ui, kvitui.SizeSpace, a, b)
		line.SetBorder(kvitui.Padding(ui, kvitui.Px(5)))
		_, before, _ := line.Sizes(geom.NewSize(200, 0))
		options(line).AlignTop, options(line).Top = true, true
		if _, after, _ := line.Sizes(geom.NewSize(200, 0)); after != before {
			t.Errorf("a line with its children at the top asks for %v, %v without", after, before)
		}
		laid(line, 200, 100)
		wantFrame(t, "the tall child of a line with its children at the top", a, 5, 5, 40, 40)
		wantFrame(t, "the short child of a line with its children at the top", b, 45+gap, 5, 40, 20)
		laid(line, 200, 30)
		wantFrame(t, "the tall child of a short line with its children at the top", a, 5, 5, 40, 40)
	})
}

func TestTheSmallerLayouts(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		// FillLayout: every child gets the whole content, and nothing is asked.
		a, b := box(40, 20), box(60, 30)
		fill := unison.NewPanel()
		fill.SetLayout(kvitui.FillLayout())
		fill.AddChild(a)
		fill.AddChild(b)
		wantSizes(t, "a filling panel", fill, geom.NewSize(100, 0),
			geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize))
		laid(fill, 120, 80)
		wantFrame(t, "the first child of a filling panel", a, 0, 0, 120, 80)
		wantFrame(t, "the second child of a filling panel", b, 0, 0, 120, 80)

		// OwnSizeLayout: what the first child asks, with each option.
		child := box(40, 20)
		own := unison.NewPanel()
		own.AddChild(child)
		own.SetBorder(kvitui.Padding(ui, kvitui.Px(5)))
		own.SetLayout(kvitui.OwnSizeLayout{})
		wantSizes(t, "a panel at its child's size", own, geom.Size{},
			geom.Size{}, geom.NewSize(40, 20), geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize))
		laid(own, 100, 60)
		wantFrame(t, "the child of a panel at its child's size", child, 5, 5, 90, 50)
		own.SetLayout(kvitui.OwnSizeLayout{CountInsets: true})
		if _, pref, _ := own.Sizes(geom.Size{}); pref != geom.NewSize(50, 30) {
			t.Errorf("a panel counting its insets asks for %v around a child of 40 by 20", pref)
		}
		child.Hidden = true
		own.SetLayout(kvitui.OwnSizeLayout{})
		if _, pref, _ := own.Sizes(geom.Size{}); pref != geom.NewSize(40, 20) {
			t.Errorf("a hidden child keeps its room by default, but the panel asks for %v", pref)
		}
		own.SetLayout(kvitui.OwnSizeLayout{HiddenGone: true})
		wantSizes(t, "a panel whose hidden child takes no room", own, geom.Size{},
			geom.Size{}, geom.Size{}, geom.NewSize(unison.DefaultMaxSize, 0))
		empty := unison.NewPanel()
		empty.SetLayout(kvitui.OwnSizeLayout{})
		if _, pref, _ := empty.Sizes(geom.Size{}); pref != (geom.Size{}) {
			t.Errorf("a panel with no children asks for %v", pref)
		}

		// AtMost: as tall as the content up to the limit, and no taller.
		for _, c := range []struct{ content, want float32 }{{30, 30}, {90, 50}} {
			capped := unison.NewPanel()
			capped.AddChild(box(40, c.content))
			capped.SetLayout(kvitui.AtMost(ui, kvitui.Px(50), kvitui.OwnSizeLayout{}))
			wantSizes(t, "a panel with a height limit", capped, geom.NewSize(100, 0),
				geom.NewSize(0, c.want), geom.NewSize(40, c.want), geom.NewSize(unison.DefaultMaxSize, c.want))
		}

		// NewRowLayout: children added after the layout is made, one column
		// each, centred on the line; a hidden child keeps its room.
		row := unison.NewPanel()
		row.SetLayout(kvitui.NewRowLayout(ui, kvitui.SizeSpace))
		first, middle, last := box(40, 20), box(60, 30), box(50, 10)
		for _, c := range []*unison.Panel{first, middle, last} {
			row.AddChild(c)
		}
		gap := float32(ui.Interface.Space())
		wantSizes(t, "a row laid out by NewRowLayout", row, geom.Size{},
			geom.NewSize(2*gap, 0), geom.NewSize(150+2*gap, 30), geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize))
		laid(row, 300, 30)
		wantFrame(t, "the first child of a row", first, 0, 5, 40, 20)
		wantFrame(t, "the middle child of a row", middle, 40+gap, 0, 60, 30)
		wantFrame(t, "the last child of a row", last, 100+2*gap, 10, 50, 10)
		middle.Hidden = true
		laid(row, 300, 30)
		wantFrame(t, "the last child of a row whose middle child is hidden", last, 100+2*gap, 10, 50, 10)

		// ShownIn and InsetsOf: the children that are not hidden, and the
		// border's insets or none.
		holder := unison.NewPanel()
		shown, hidden := box(1, 1), box(1, 1)
		hidden.Hidden = true
		holder.AddChild(shown)
		holder.AddChild(hidden)
		for _, c := range []struct {
			what   string
			border unison.Border
			want   geom.Insets
		}{
			{"no border", nil, geom.Insets{}},
			{"padding", kvitui.Padding(ui, kvitui.Px(6)), geom.NewUniformInsets(6)},
			{"insets", kvitui.Insets(ui, kvitui.Px(1), kvitui.Px(2), nil, kvitui.Px(4)), geom.Insets{Top: 1, Left: 2, Right: 4}},
		} {
			holder.SetBorder(c.border)
			if got := kvitui.InsetsOf(holder); got != c.want {
				t.Errorf("InsetsOf a panel with %s is %v, want %v", c.what, got, c.want)
			}
			if got := kvitui.ShownIn(holder); len(got) != 1 || got[0] != shown {
				t.Errorf("ShownIn a panel with %s is %v", c.what, got)
			}
		}

		// Relayout marks the panel and what it sits in, and with children
		// what is inside it.
		outer, mid, inner := unison.NewPanel(), unison.NewPanel(), unison.NewPanel()
		outer.AddChild(mid)
		mid.AddChild(inner)
		outer.NeedsLayout, mid.NeedsLayout, inner.NeedsLayout = false, false, false
		kvitui.Relayout(mid, false)
		if !outer.NeedsLayout || !mid.NeedsLayout || inner.NeedsLayout {
			t.Errorf("laying a panel out again marked outer %v, itself %v, inner %v", outer.NeedsLayout, mid.NeedsLayout, inner.NeedsLayout)
		}
		kvitui.Relayout(mid, true)
		if !inner.NeedsLayout {
			t.Error("laying a panel and its children out again did not mark the child")
		}
	})
}

func TestAFlowWrapsItsChildren(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		gap := float32(ui.Interface.Space())
		build := func() (*unison.Panel, *kvitui.FlowLayout, []*unison.Panel) {
			parts := []*unison.Panel{box(60, 20), box(60, 30), box(60, 10), box(60, 10)}
			parts[3].Hidden = true
			p := unison.NewPanel()
			l := kvitui.NewFlowLayout(ui, kvitui.SizeSpace)
			p.SetLayout(l)
			for _, c := range parts {
				p.AddChild(c)
			}
			return p, l, parts
		}
		oneLine := 180 + 2*gap
		flow, opt, parts := build()
		// Offered no width, it asks for one line; offered one, it wraps.
		wantSizes(t, "a flow offered no width", flow, geom.Size{},
			geom.NewSize(0, 30), geom.NewSize(oneLine, 30), geom.NewSize(unison.DefaultMaxSize, 30))
		wantSizes(t, "a flow offered room for two", flow, geom.NewSize(130+gap, 0),
			geom.NewSize(0, 40+gap), geom.NewSize(oneLine, 40+gap), geom.NewSize(unison.DefaultMaxSize, 40+gap))
		laid(flow, 130+gap, 100)
		wantFrame(t, "the first child of a flow", parts[0], 0, 0, 60, 20)
		wantFrame(t, "the second child of a flow", parts[1], 60+gap, 0, 60, 30)
		wantFrame(t, "the child that starts the second line", parts[2], 0, 30+gap, 60, 10)

		// Unasked: the width assumed when none is offered.
		opt.Unasked = kvitui.Px(100)
		if _, pref, _ := flow.Sizes(geom.Size{}); pref != geom.NewSize(oneLine, 60+2*gap) {
			t.Errorf("a flow assuming 100 wide asks for %v, want three lines", pref)
		}
		// FillsHint: the width asked for is the width offered or assumed.
		opt.FillsHint = true
		if _, pref, _ := flow.Sizes(geom.Size{}); pref.Width != 100 {
			t.Errorf("a flow that fills, assuming 100 wide, asks for %v wide", pref.Width)
		}
		if _, pref, _ := flow.Sizes(geom.NewSize(260, 0)); pref.Width != 260 {
			t.Errorf("a flow that fills, offered 260, asks for %v wide", pref.Width)
		}

		// ItemWidth: every child at the width the function answers.
		flow, opt, parts = build()
		opt.ItemWidth = func(lineWidth float32) float32 { return (lineWidth - gap) / 2 }
		laid(flow, 200+gap, 100)
		wantFrame(t, "the first child of a flow of equal boxes", parts[0], 0, 0, 100, 20)
		wantFrame(t, "the second child of a flow of equal boxes", parts[1], 100+gap, 0, 100, 30)
		wantFrame(t, "the third child of a flow of equal boxes", parts[2], 0, 30+gap, 100, 10)
	})
}

// A text area that fits its text asks for the text's size; one that does not
// asks for three rows whatever it holds.
func TestATextAreaFitsItsText(t *testing.T) {
	screen, ui := layoutSession(t)
	screen.Do(func() {
		long := "A paragraph long enough that it wraps onto several lines when the area is narrow, and then some more of it."
		area := kvitui.NewTextArea(ui)
		area.Plain = true
		area.SetText(long + "\n" + long + "\n" + long + "\n" + long)
		_, rows, _ := area.Sizes(geom.NewSize(240, 0))
		if rows.Height != float32(3*ui.Interface.RowHeight()) {
			t.Fatalf("a text area that does not fit asks for %v tall, want three rows", rows.Height)
		}
		edit := area.Edit()

		// Unwrapped: the text's own width and height.
		edit.SetWrap(false)
		_, text, _ := edit.Sizes(geom.Size{})
		area.Fit = kvitui.FitText
		wantSizes(t, "an unwrapped area fitting its text", area.AsPanel(), geom.NewSize(240, 0),
			text, text, geom.NewSize(unison.DefaultMaxSize, text.Height))
		area.Fit = kvitui.FitTextAtLeastOffered
		wantSizes(t, "an unwrapped area fitting its text in a narrower width", area.AsPanel(), geom.NewSize(240, 0),
			text, text, geom.NewSize(unison.DefaultMaxSize, text.Height))
		wide := geom.NewSize(text.Width+500, text.Height)
		wantSizes(t, "an unwrapped area fitting its text in a wider width", area.AsPanel(), geom.NewSize(text.Width+500, 0),
			wide, wide, geom.NewSize(unison.DefaultMaxSize, text.Height))

		// Wrapped: the height of the text at the width offered, and no width
		// of its own.
		edit.SetWrap(true)
		_, wrapped, _ := edit.Sizes(geom.NewSize(240, 0))
		if wrapped.Height <= text.Height {
			t.Fatalf("the text is %v tall wrapped and %v unwrapped, so the test shows nothing", wrapped.Height, text.Height)
		}
		area.Fit = kvitui.FitText
		wantSizes(t, "a wrapped area fitting its text", area.AsPanel(), geom.NewSize(240, 0),
			geom.NewSize(0, wrapped.Height), geom.NewSize(0, wrapped.Height), geom.NewSize(unison.DefaultMaxSize, wrapped.Height))
	})
}
