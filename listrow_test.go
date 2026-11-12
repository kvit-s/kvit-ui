package kvitui_test

import (
	"image"
	"image/color"
	"strconv"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// rowsIn is a session holding list rows 300 wide, one under another.
func rowsIn(t *testing.T, build func(ui *kvitui.UI) []*kvitui.ListRow) (*unison.HeadlessScreen, *kvitui.UI) {
	t.Helper()
	return session(t, func(ui *kvitui.UI) []unison.Paneler {
		rows := build(ui)
		parts := make([]unison.Paneler, len(rows))
		for i, r := range rows {
			parts[i] = r
		}
		col := kvitui.Stack(ui, kvitui.SizeSpace, parts...)
		col.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(300, 0)})
		return []unison.Paneler{col}
	})
}

func TestAListRowsHeightAndPlacement(t *testing.T) {
	var plain, tall, fitted, padded, placed *kvitui.ListRow
	var lines *kvitui.Label
	var placedWith geom.Rect
	var part *unison.Panel
	screen, ui := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		plain = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Plain"))
		tall = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Tall"))
		tall.Height = kvitui.Px(38)
		lines = kvitui.NewLabel(ui, "First line\nSecond line")
		lines.Wrap = true
		fitted = kvitui.NewListRow(ui, kvitui.FullWidth(lines))
		fitted.FitContent = true
		other := kvitui.NewLabel(ui, "First line\nSecond line")
		other.Wrap = true
		padded = kvitui.NewListRow(ui, kvitui.FullWidth(other))
		padded.FitContent, padded.Margins = true, kvitui.Padding(ui, kvitui.Px(4))
		part = unison.NewPanel()
		placed = kvitui.NewListRow(ui, part)
		placed.Height = kvitui.Px(30)
		placed.Place = func(content geom.Rect) {
			placedWith = content
			part.SetFrameRect(geom.NewRect(content.X+12, content.Y+3, 20, 20))
		}
		return []*kvitui.ListRow{plain, tall, fitted, padded, placed}
	})
	screen.Do(func() {
		m := ui.Interface
		if h := plain.FrameRect().Height; h != float32(m.RowHeightSlim()) {
			t.Errorf("a plain row is %v tall, want the slim form's %v", h, m.RowHeightSlim())
		}
		if h := tall.FrameRect().Height; h != 38 {
			t.Errorf("a row with a height of 38 is %v tall", h)
		}
		_, two, _ := lines.Sizes(geom.Size{})
		if h := fitted.FrameRect().Height; h != two.Height || two.Height <= float32(m.RowHeightSlim()) {
			t.Errorf("a row fitting two lines %v tall is %v tall, a slim row %v", two.Height, h, m.RowHeightSlim())
		}
		if h := padded.FrameRect().Height; h != two.Height+8 {
			t.Errorf("a row fitting two lines %v tall inside margins of 4 is %v tall", two.Height, h)
		}
		// Height wins over FitContent.
		fitted.Height = kvitui.Px(50)
		if _, pref, _ := fitted.Sizes(geom.NewSize(300, 0)); pref.Height != 50 {
			t.Errorf("a row with a height and FitContent asks for %v tall", pref.Height)
		}
		// Place is given the box inside the margins, and the row asks for
		// the width it is offered.
		tight := float32(m.SpaceTight())
		if want := geom.NewRect(tight, 0, 300-2*tight, 30); placedWith != want {
			t.Errorf("Place was given %v, want %v", placedWith, want)
		}
		wantFrame(t, "the part a row placed itself", part, tight+12, 3, 20, 20)
		wantSizes(t, "a row that places its own parts", placed.AsPanel(), geom.NewSize(240, 0),
			geom.NewSize(0, 30), geom.NewSize(240, 30), geom.NewSize(unison.DefaultMaxSize, 30))
	})
}

// FitRoom adds room inside the content's box, where a release counts; the
// margins beside it still do not. FitAtLeastForm keeps the form's height
// under a short child. Arrange lays the children out in place of the row's
// line.
func TestAListRowsFitRoomAndArrange(t *testing.T) {
	var roomy, short, high, arranged *kvitui.ListRow
	var words *kvitui.Label
	var grower, fixed *unison.Panel
	opened := 0
	screen, ui := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		words = kvitui.NewLabel(ui, "First line\nSecond line")
		words.Wrap = true
		roomy = kvitui.NewListRow(ui, kvitui.FullWidth(words))
		roomy.Interactive, roomy.NoOpensMark, roomy.FitContent, roomy.FitRoom = true, true, true, kvitui.Px(10)
		roomy.OnActivate = func() { opened++ }
		short = kvitui.NewListRow(ui, box(40, 6))
		short.FitContent, short.FitAtLeastForm = true, true
		high = kvitui.NewListRow(ui, box(40, 60))
		high.FitContent, high.FitAtLeastForm = true, true
		grower, fixed = kvitui.Grow(box(30, 12)), box(50, 16)
		arranged = kvitui.NewListRow(ui, grower, fixed)
		arranged.Arrange = kvitui.NewLineLayout(ui, kvitui.SizeSpace)
		return []*kvitui.ListRow{roomy, short, high, arranged}
	})
	var tight float32
	screen.Do(func() {
		m := ui.Interface
		tight = float32(m.SpaceTight())
		_, two, _ := words.Sizes(geom.NewSize(300-2*tight, 0))
		if h := roomy.FrameRect().Height; h != two.Height+10 {
			t.Errorf("a row fitting two lines %v tall with room of 10 is %v tall", two.Height, h)
		}
		wantFrame(t, "the words of a row with room", words.AsPanel(), tight, 5, 300-2*tight, two.Height)
		slim := float32(m.RowHeightSlim())
		if h := short.FrameRect().Height; h != slim {
			t.Errorf("a row at least its form's height, holding a child 6 tall, is %v tall, want %v", h, slim)
		}
		if h := high.FrameRect().Height; h != 60 {
			t.Errorf("a row at least its form's height, holding a child 60 tall, is %v tall", h)
		}
		// Height still wins over both.
		short.Height = kvitui.Px(12)
		if _, pref, _ := short.Sizes(geom.NewSize(300, 0)); pref.Height != 12 {
			t.Errorf("a row with a height of 12 and FitAtLeastForm asks for %v tall", pref.Height)
		}
		// Arrange: the grower takes the width the other leaves, and the row
		// asks for the line's own width.
		gap := float32(m.Space())
		wantFrame(t, "a child that grows in a row's line", grower.AsPanel(), tight, (slim-12)/2, 300-2*tight-50-gap, 12)
		wantFrame(t, "a child at its own width in a row's line", fixed, 300-tight-50, (slim-16)/2, 50, 16)
		wantSizes(t, "a row arranged by a line", arranged.AsPanel(), geom.NewSize(300, 0),
			geom.NewSize(2*tight, slim), geom.NewSize(30+50+gap+2*tight, slim), geom.NewSize(unison.DefaultMaxSize, slim))
		// With FitContent and FitAtLeastForm it is as tall as the taller of the
		// form and the line, as Kvit Works' rows on a line are.
		arranged.FitContent, arranged.FitAtLeastForm = true, true
		fixed.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
			return geom.Size{}, geom.NewSize(50, 70), geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
		})
		wantSizes(t, "a tall row arranged by a line", arranged.AsPanel(), geom.NewSize(300, 0),
			geom.NewSize(2*tight, 70), geom.NewSize(30+50+gap+2*tight, 70), geom.NewSize(unison.DefaultMaxSize, 70))
		fixed.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
			return geom.Size{}, geom.NewSize(50, 10), geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
		})
		wantSizes(t, "a short row arranged by a line", arranged.AsPanel(), geom.NewSize(300, 0),
			geom.NewSize(2*tight, slim), geom.NewSize(30+50+gap+2*tight, slim), geom.NewSize(unison.DefaultMaxSize, slim))
	})
	// A release in the left margin does not activate the row; one in the room
	// added above the words does. PanelPoint counts from inside the margins.
	screen.Click(screen.PanelPoint(roomy, geom.NewPoint(-tight/2, 2)))
	if opened != 0 {
		t.Errorf("a press in the left margin of a row with room opened it %d times", opened)
	}
	screen.Click(screen.PanelPoint(roomy, geom.NewPoint(150, 2)))
	if opened != 1 {
		t.Errorf("a press in the room above a row's words opened it %d times", opened)
	}
}

func TestAListRowsSlots(t *testing.T) {
	var row *kvitui.ListRow
	lead, title, sub, trail := box(16, 16), box(80, 14), box(120, 12), box(40, 20)
	screen, ui := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		row = kvitui.NewListRow(ui)
		row.Height = kvitui.Px(40)
		row.Leading, row.Title, row.Subtitle, row.Trailing = lead, title, sub, trail
		return []*kvitui.ListRow{row}
	})
	screen.Do(func() {
		m := ui.Interface
		tight, gap := float32(m.SpaceTight()), float32(m.Space())
		if len(row.Children()) != 4 {
			t.Fatalf("a row with four slots holds %d panels", len(row.Children()))
		}
		left, right := tight+16+gap, 300-tight-40-gap
		wantFrame(t, "the leading slot", lead, tight, 12, 16, 16)
		wantFrame(t, "the title", title, left, 7, right-left, 14)
		wantFrame(t, "the subtitle", sub, left, 21, right-left, 12)
		wantFrame(t, "the trailing slot", trail, 300-tight-40, 10, 40, 20)
		// A hidden slot takes no room and no gap.
		lead.Hidden, sub.Hidden = true, true
		row.MarkForLayoutAndRedraw()
	})
	screen.Sync()
	screen.Do(func() {
		m := ui.Interface
		tight, gap := float32(m.SpaceTight()), float32(m.Space())
		wantFrame(t, "the title with no leading slot and no subtitle", title, tight, 13, 300-2*tight-40-gap, 14)
		// A row as tall as its slots.
		row.Height, row.FitContent = nil, true
		lead.Hidden, sub.Hidden = false, false
		if _, pref, _ := row.Sizes(geom.NewSize(300, 0)); pref.Height != 26 || pref.Width != 16+120+40+2*gap+2*tight {
			t.Errorf("a row fitting its slots asks for %v", pref)
		}
	})
}

// near reports whether two colours are within a few levels of one another.
func near(a, b color.Color) bool {
	ar, ag, ab, _ := a.RGBA()
	br, bg, bb, _ := b.RGBA()
	d := func(x, y uint32) uint32 { return max(x, y) - min(x, y) }
	return d(ar, br) < 0x600 && d(ag, bg) < 0x600 && d(ab, bb) < 0x600
}

func TestAListRowsMarks(t *testing.T) {
	t.Setenv(uitest.ScaleVariable, "1")
	var plain, barred, edge, half, anyFocus, round, quiet *kvitui.ListRow
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		row := func(words string) *kvitui.ListRow {
			r := kvitui.NewListRow(ui, kvitui.NewLabel(ui, words))
			r.Interactive, r.NoOpensMark, r.Label, r.Rule = true, true, words, false
			return r
		}
		plain, barred, edge, anyFocus, round, quiet = row("Plain"), row("Bar"), row("Edge"), row("Any focus"), row("Round"), row("Quiet")
		barred.Bar = kvitui.Px(3)
		edge.RingAtEdge = true
		half = row("Half in")
		half.RingHalfIn = true
		anyFocus.RingOnAnyFocus = true
		round.Radius = kvitui.Px(10)
		quiet.NoHoverTint = true
		return kvitui.Width(ui, kvitui.Px(300), kvitui.Stack(ui, kvitui.SizeSpace, plain, barred, edge, half, anyFocus, round, quiet))
	})
	var img image.Image
	shoot := func() {
		s.Do(func() { s.Window.Content().MarkForRedraw() })
		s.Sync()
		img = s.Capture()
	}
	// at is the colour of a pixel of a row, from its top-left corner.
	at := func(r *kvitui.ListRow, x, y float32) color.Color {
		var o geom.Point
		s.Do(func() { o = r.PointToRoot(geom.Point{}) })
		return img.At(int(o.X+x), int(o.Y+y))
	}
	ring := float32(s.UI.Interface.FocusRingWidth())
	mid := float32(s.UI.Interface.RowHeightSlim()) / 2

	// The bar is drawn down the left edge only while the row is selected.
	shoot()
	// Sampled above the words, which start a tight space in.
	const above = 3
	if !near(at(barred, 1, above), at(barred, 200, above)) {
		t.Error("a row that is not selected draws its bar")
	}
	s.Do(func() { plain.Selected, barred.Selected, round.Selected = true, true, true })
	shoot()
	if !near(at(plain, 1, above), at(plain, 200, above)) {
		t.Error("a selected row with no bar is not one colour across")
	}
	if near(at(barred, 1, above), at(barred, 200, above)) || !near(at(barred, 4, above), at(barred, 200, above)) {
		t.Error("a selected row with a bar 3 wide does not draw it at its left edge")
	}
	// A radius rounds the tint: the corner is the ground beside the row.
	if near(at(plain, 0, 0), at(plain, -3, 0)) || !near(at(round, 0, 0), at(round, -3, 0)) {
		t.Error("the selection tint's corner is square with a radius, or round without one")
	}
	s.Do(func() { plain.Selected, barred.Selected, round.Selected = false, false, false })

	// The ring: a ring's width in by default, at the edge when asked.
	// Sampled along the top, clear of the words.
	inner := ring + ring/2
	s.Do(func() { edge.Focus() })
	shoot()
	ringInk := at(edge, 200, ring/2)
	if near(ringInk, at(edge, 200, mid)) || near(at(edge, 200, inner), ringInk) {
		t.Fatal("a row with its ring at the edge does not draw it there")
	}
	s.Do(func() { barred.Focus() })
	shoot()
	if !near(at(barred, 200, inner), ringInk) || near(at(barred, 200, ring/2), ringInk) {
		t.Error("a row focused from the keyboard draws no ring a ring's width inside its edge")
	}
	// Half a ring's width in: the edge's pixel row is clear, the next two are
	// the ring at a ring width of 2, and the row the default ring reaches
	// past them is clear again.
	s.Do(func() { half.Focus() })
	shoot()
	if near(at(half, 200, 0), ringInk) || !near(at(half, 200, ring/2), ringInk) || !near(at(half, 200, ring/2+ring-1), ringInk) ||
		near(at(half, 200, ring/2+ring), ringInk) {
		t.Error("a row with its ring half a width in does not draw it on the second and third pixel rows")
	}
	// A press gives the focus without a ring, unless any focus shows it.
	s.Screen.Click(s.Screen.PanelPoint(plain, geom.NewPoint(200, mid)))
	s.Screen.MouseMove(geom.NewPoint(5, 5), 0)
	shoot()
	if near(at(plain, 200, inner), ringInk) {
		t.Error("a row focused by a press draws its ring")
	}
	s.Screen.Click(s.Screen.PanelPoint(anyFocus, geom.NewPoint(200, mid)))
	s.Screen.MouseMove(geom.NewPoint(5, 5), 0)
	shoot()
	if !near(at(anyFocus, 200, inner), ringInk) {
		t.Error("a row that shows its ring for any focus draws none after a press")
	}

	// The hover tint, and the row that leaves it out.
	s.Do(func() { s.Window.Content().RequestFocus() })
	shoot()
	ground := at(barred, 200, mid)
	s.Screen.MouseMove(s.Screen.PanelPoint(barred, geom.NewPoint(250, mid)), 0)
	shoot()
	if near(at(barred, 200, mid), ground) {
		t.Error("a row under the pointer is not tinted")
	}
	s.Screen.MouseMove(s.Screen.PanelPoint(quiet, geom.NewPoint(250, mid)), 0)
	shoot()
	if !near(at(quiet, 200, mid), ground) {
		t.Error("a row with no hover tint is tinted under the pointer")
	}
	// A disabled row is tinted as an enabled one is, unless told not to be.
	s.Do(func() { barred.SetEnabled(false) })
	s.Screen.MouseMove(s.Screen.PanelPoint(barred, geom.NewPoint(250, mid)), 0)
	shoot()
	if near(at(barred, 200, mid), ground) {
		t.Error("a disabled row is not tinted under the pointer, as it was before the option existed")
	}
	s.Do(func() { barred.NoHoverTintDisabled = true })
	shoot()
	if !near(at(barred, 200, mid), ground) {
		t.Error("a disabled row told to take no tint is tinted under the pointer")
	}
}

// A row told to count the pointer on anything inside it stays hovered over
// its own words, and tells OnHover once for a move from its padding onto
// them.
func TestAListRowHoveredOverItsWords(t *testing.T) {
	var whole, plain *kvitui.ListRow
	var wholeWords, plainWords *kvitui.Label
	var said []bool
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		wholeWords, plainWords = kvitui.NewLabel(ui, "Whole"), kvitui.NewLabel(ui, "Plain")
		whole = kvitui.NewListRow(ui, wholeWords)
		whole.Interactive, whole.HoverWhole = true, true
		whole.OnHover = func(over bool) { said = append(said, over) }
		plain = kvitui.NewListRow(ui, plainWords)
		plain.Interactive = true
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(300), kvitui.Stack(ui, kvitui.SizeSpace, whole, plain))}
	})
	hovered := func(r *kvitui.ListRow) (h bool) {
		screen.Do(func() { h = r.Hovered() })
		return h
	}
	padding := func(r *kvitui.ListRow) geom.Point { return screen.PanelPoint(r, geom.NewPoint(250, 5)) }
	screen.MouseMove(padding(plain), 0)
	if !hovered(plain) {
		t.Fatal("a row is not hovered with the pointer on its padding")
	}
	screen.MouseMove(screen.PanelCenter(plainWords), 0)
	if hovered(plain) {
		t.Error("a plain row stays hovered over its words, which unison does not tell it")
	}
	screen.MouseMove(padding(whole), 0)
	screen.MouseMove(screen.PanelCenter(wholeWords), 0)
	screen.Sync()
	if !hovered(whole) {
		t.Error("a row hovered whole lost the hover over its words")
	}
	if len(said) != 1 || !said[0] {
		t.Errorf("a move from the padding onto the words told OnHover %v", said)
	}
	screen.MouseMove(padding(plain), 0)
	screen.Sync()
	if hovered(whole) || len(said) != 2 || said[1] {
		t.Errorf("after the pointer left, the row is hovered %v and OnHover was told %v", hovered(whole), said)
	}
}

func TestWhenAListRowActs(t *testing.T) {
	cases := []struct {
		name      string
		press     kvitui.RowPress
		anyButton bool
		// How many times the row has acted after: the left button down, its
		// release, a right click, and a left press let go outside the row.
		leftDown, leftUp, right, outside int
	}{
		{"on release, left button", kvitui.PressOnRelease, false, 0, 1, 1, 1},
		{"on release, any button", kvitui.PressOnRelease, true, 0, 1, 2, 2},
		{"on the way down, left button", kvitui.PressOnDown, false, 1, 1, 1, 2},
		{"on the way down, any button", kvitui.PressOnDown, true, 1, 1, 2, 3},
	}
	rows := make([]*kvitui.ListRow, len(cases))
	acted := make([]int, len(cases))
	screen, _ := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		for i, c := range cases {
			rows[i] = kvitui.NewListRow(ui, kvitui.NewLabel(ui, c.name))
			rows[i].Interactive, rows[i].PressOn, rows[i].AnyButton = true, c.press, c.anyButton
			rows[i].OnActivate = func() { acted[i]++ }
		}
		return rows
	})
	for i, c := range cases {
		at := screen.PanelPoint(rows[i], geom.NewPoint(250, 5))
		check := func(what string, want int) {
			t.Helper()
			screen.Sync()
			if acted[i] != want {
				t.Errorf("%s: after %s the row has acted %d times, want %d", c.name, what, acted[i], want)
			}
		}
		screen.MouseMove(at, 0)
		screen.MouseDown(at, unison.ButtonLeft, 0)
		check("the left button going down", c.leftDown)
		screen.MouseUp(at, unison.ButtonLeft, 0)
		check("its release", c.leftUp)
		screen.ClickWith(at, unison.ButtonRight, 0)
		check("a right click", c.right)
		screen.MouseDown(at, unison.ButtonLeft, 0)
		screen.MouseMove(geom.NewPoint(700, 280), 0)
		screen.MouseUp(geom.NewPoint(700, 280), unison.ButtonLeft, 0)
		check("a press let go outside", c.outside)
	}
}

func TestAListRowsPointerHooks(t *testing.T) {
	var row *kvitui.ListRow
	var doubles, opened, dragged int
	var pressedAt geom.Point
	var pressedWith int
	handled := false
	screen, _ := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		row = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Project"))
		row.Interactive = true
		row.OnActivate = func() { opened++ }
		row.OnDoubleClick = func() { doubles++ }
		row.OnPress = func(where geom.Point, button int) { pressedAt, pressedWith = where, button }
		row.OnDrag = func(geom.Point) { dragged++ }
		row.OnRelease = func() bool { return handled }
		return []*kvitui.ListRow{row}
	})
	at := screen.PanelPoint(row, geom.NewPoint(250, 5))
	screen.Click(at)
	if opened != 1 || doubles != 0 || pressedWith != unison.ButtonLeft || abs(pressedAt.Y-5) > 1 {
		t.Errorf("one click opened %d, double-clicked %d, pressed at %v with button %d", opened, doubles, pressedAt, pressedWith)
	}
	// The second press of a double click runs OnDoubleClick, and each
	// release still activates.
	screen.DoubleClick(at)
	if doubles != 1 || opened != 3 {
		t.Errorf("a double click ran OnDoubleClick %d times and opened the row %d times in all", doubles, opened)
	}
	// A drag is reported, and a release that was handled does not activate.
	handled = true
	screen.Drag(at, geom.NewPoint(at.X-40, at.Y), 4)
	if dragged == 0 || opened != 3 {
		t.Errorf("a drag whose release was handled reported %d moves and opened the row %d times in all", dragged, opened)
	}
	handled = false
	screen.Drag(at, geom.NewPoint(at.X-40, at.Y), 4)
	if opened != 4 {
		t.Errorf("a drag let go inside the row, not handled, left it opened %d times in all", opened)
	}
}

func TestAListRowsKeyAndFocusHooks(t *testing.T) {
	var row, static *kvitui.ListRow
	var inside *kvitui.IconButton
	var keys []unison.KeyCode
	var focuses []bool
	take := false
	opened, blurred := 0, 0
	screen, _ := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		before := kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Before"))
		before.Interactive = true
		row = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Track"))
		row.Interactive = true
		row.OnActivate = func() { opened++ }
		row.OnKey = func(key unison.KeyCode, _ mod.Modifiers) bool {
			keys = append(keys, key)
			return take
		}
		row.OnFocus = func(byPointer bool) { focuses = append(focuses, byPointer) }
		row.OnBlur = func() { blurred++ }
		inside = kvitui.NewIconButton(ui, "pencil", "Rename")
		static = kvitui.NewListRow(ui, inside)
		static.OnKey = row.OnKey
		return []*kvitui.ListRow{before, row, static}
	})
	// A press, then Tab from the row before, then Focus().
	screen.Click(screen.PanelPoint(row, geom.NewPoint(250, 5)))
	screen.Do(func() { row.Parent().Children()[0].RequestFocus() })
	screen.KeyPress(unison.KeyTab, 0)
	screen.Do(func() { row.Parent().Children()[0].RequestFocus() })
	screen.Sync()
	screen.Do(func() { row.Focus() })
	screen.Sync()
	if len(focuses) != 3 || !focuses[0] || focuses[1] || focuses[2] {
		t.Errorf("a press, Tab and Focus() told OnFocus byPointer %v, want true, false, false", focuses)
	}
	if blurred != 2 {
		t.Errorf("the row lost the focus twice and OnBlur ran %d times", blurred)
	}
	// OnKey sees a key before the row acts on it, and can take it.
	opened, keys = 0, nil
	screen.KeyPress(unison.KeyReturn, 0)
	if opened != 1 || len(keys) != 1 || keys[0] != unison.KeyReturn {
		t.Errorf("Return, not taken, opened the row %d times after OnKey saw %v", opened, keys)
	}
	take = true
	screen.KeyPress(unison.KeyReturn, 0)
	screen.KeyPress(unison.KeyDown, mod.Control)
	if opened != 1 || len(keys) != 3 {
		t.Errorf("keys OnKey took opened the row %d times in all; it saw %v", opened, keys)
	}
	// A key from a control inside a row nobody presses reaches OnKey too.
	keys, take = nil, false
	screen.Do(func() { inside.Focus() })
	screen.KeyPress(unison.KeyDown, 0)
	if len(keys) != 1 || keys[0] != unison.KeyDown {
		t.Errorf("a key on a control inside the row reached OnKey as %v", keys)
	}

	// OnKeyDown is offered a key first, with whether it repeats; a key it
	// takes reaches neither OnKey nor the row, and one it leaves reaches both.
	var repeats []bool
	takeDown := false
	screen.Do(func() {
		row.OnKeyDown = func(key unison.KeyCode, _ mod.Modifiers, repeat bool) bool {
			if key == unison.KeyReturn {
				repeats = append(repeats, repeat)
			}
			return takeDown
		}
		row.Focus()
	})
	screen.Sync()
	opened, keys, take = 0, nil, false
	screen.KeyDown(unison.KeyReturn, 0)
	screen.KeyDown(unison.KeyReturn, 0)
	screen.KeyUp(unison.KeyReturn, 0)
	if len(repeats) != 2 || repeats[0] || !repeats[1] {
		t.Errorf("Return held down told OnKeyDown repeat %v, want false then true", repeats)
	}
	if opened != 1 || len(keys) != 2 {
		t.Errorf("keys OnKeyDown left reached OnKey %d times and opened the row %d times, want 2 and 1", len(keys), opened)
	}
	opened, keys, takeDown = 0, nil, true
	screen.KeyPress(unison.KeyReturn, 0)
	if opened != 0 || len(keys) != 0 {
		t.Errorf("a key OnKeyDown took reached OnKey %d times and opened the row %d times", len(keys), opened)
	}
}

// A press on a row that shows the keyboard ring ends it, unless the row keeps
// it; a press on a row without the focus gives it the focus without a ring
// either way.
func TestAListRowKeepsItsRingOnAPressWhenTold(t *testing.T) {
	var plain, keeps, other *kvitui.ListRow
	screen, _ := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		row := func(words string) *kvitui.ListRow {
			r := kvitui.NewListRow(ui, kvitui.NewLabel(ui, words))
			r.Interactive, r.NoOpensMark = true, true
			return r
		}
		plain, keeps, other = row("Plain"), row("Keeps"), row("Other")
		keeps.RingKeptOnPress = true
		return []*kvitui.ListRow{plain, keeps, other}
	})
	keyboard := func(r *kvitui.ListRow) (k bool) {
		screen.Do(func() { k = r.KeyboardFocus() })
		return k
	}
	for _, c := range []struct {
		row  *kvitui.ListRow
		keep bool
	}{{plain, false}, {keeps, true}} {
		// From another row, since the first row holds the focus the window
		// handed out, and Focus() on a row that holds it changes nothing.
		screen.Do(func() { other.Focus() })
		screen.Sync()
		screen.Do(func() { c.row.Focus() })
		screen.Sync()
		if !keyboard(c.row) {
			t.Fatalf("a row given the focus with Focus() shows no ring (keeps it on a press: %v)", c.keep)
		}
		screen.Click(screen.PanelPoint(c.row, geom.NewPoint(250, 5)))
		if got := keyboard(c.row); got != c.keep {
			t.Errorf("after a press on a row showing its ring, the ring shows %v; want %v", got, c.keep)
		}
		screen.Do(func() { other.Focus() })
		screen.Sync()
		screen.Click(screen.PanelPoint(c.row, geom.NewPoint(250, 5)))
		var focused bool
		screen.Do(func() { focused = c.row.Focused() })
		if !focused || keyboard(c.row) {
			t.Errorf("a row focused by a press holds the focus %v and shows its ring %v (keeps it on a press: %v)", focused, keyboard(c.row), c.keep)
		}
	}
}

// Focus() on a row that already holds the focus changes nothing then, and
// the row counts its next gain of the focus as one from the keyboard. A row
// told to forgets that request when it loses the focus, so the focus coming
// back by neither a key nor a press shows no ring; a Focus() that moves the
// focus still shows it.
func TestAListRowEndsAFocusRequestOnBlurWhenTold(t *testing.T) {
	var row, other *kvitui.ListRow
	screen, _ := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		newRow := func(words string) *kvitui.ListRow {
			r := kvitui.NewListRow(ui, kvitui.NewLabel(ui, words))
			r.Interactive, r.NoOpensMark, r.RingKeptOnPress = true, true, true
			return r
		}
		row, other = newRow("Row"), newRow("Other")
		return []*kvitui.ListRow{row, other}
	})
	keyboard := func() (k bool) {
		screen.Do(func() { k = row.KeyboardFocus() })
		return k
	}
	press := func() { screen.Click(screen.PanelPoint(row, geom.NewPoint(250, 5))) }
	for _, ends := range []bool{false, true} {
		screen.Do(func() { row.FocusRequestEndsOnBlur = ends })
		for _, pressedAgain := range []bool{true, false} {
			// From the other row, so the press gives the row the focus.
			screen.Do(func() { other.RequestFocus() })
			press()
			var focused bool
			screen.Do(func() { focused = row.Focused() })
			if !focused || keyboard() {
				t.Fatalf("a pressed row holds the focus %v and shows its ring %v", focused, keyboard())
			}
			screen.Do(func() { row.Focus() })
			screen.Sync()
			if keyboard() {
				t.Errorf("Focus() on a row that holds the focus shows its ring at once (ends on blur: %v)", ends)
			}
			if pressedAgain {
				press()
			}
			// The focus leaves and comes back, asked for by code.
			screen.Do(func() { other.RequestFocus() })
			screen.Do(func() { row.RequestFocus() })
			if got := keyboard(); got == ends {
				t.Errorf("with FocusRequestEndsOnBlur %v and pressed again %v, the focus given back by code shows the ring %v",
					ends, pressedAgain, got)
			}
		}
	}
	screen.Do(func() { other.RequestFocus() })
	screen.Do(func() { row.Focus() })
	screen.Sync()
	if !keyboard() {
		t.Error("with FocusRequestEndsOnBlur, Focus() that moves the focus to the row shows no ring")
	}
}

// A row told to leaves Return, Enter and Space pressed with Shift to what it
// sits in: in a region, Shift+Space pages up rather than pressing the row.
func TestAListRowLeavesShiftedKeysWhenTold(t *testing.T) {
	var region *kvitui.Region
	var rows []*kvitui.ListRow
	pressed := 0
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		parts := make([]unison.Paneler, 30)
		for i := range parts {
			r := kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Version "+strconv.Itoa(i+1)))
			r.Interactive, r.NoOpensMark = true, true
			r.OnActivate = func() { pressed++ }
			rows = append(rows, r)
			parts[i] = r
		}
		region = kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.Px(0), parts...))
		region.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(480, 200)})
		return []unison.Paneler{region}
	})
	position := func() (y float32) {
		screen.Do(func() { _, y = region.Position() })
		return y
	}
	// Scrolled to the middle, with a row in view given the focus.
	var row *kvitui.ListRow
	var page float32
	screen.Do(func() {
		if !region.Scrolls() {
			t.Fatal("thirty rows in 200 px do not scroll")
		}
		region.ScrollTo(region.Children()[1].Self.(*unison.ScrollBar).MaxValue() / 2)
		view := region.RectToRoot(region.ContentRect(false))
		for _, r := range rows {
			if box := r.RectToRoot(r.ContentRect(true)); box.Y >= view.Y && box.Bottom() <= view.Bottom() {
				row = r
				break
			}
		}
		page = max(float32(ui.Interface.RowHeight()), region.ViewHeight()-float32(ui.Interface.RowHeight()))
	})
	if row == nil {
		t.Fatal("no row is in view")
	}
	for _, leaves := range []bool{false, true} {
		pressed = 0
		screen.Do(func() {
			row.NoShiftActivate = leaves
			row.Focus()
		})
		screen.Sync()
		y := position()
		if y < page {
			t.Fatalf("the region is scrolled %.1f, less than a page of %.1f", y, page)
		}
		screen.KeyPress(unison.KeySpace, mod.Shift)
		if leaves {
			if got := position(); pressed != 0 || got != y-page {
				t.Errorf("with NoShiftActivate, Shift+Space pressed the row %d times and moved the region from %.1f to %.1f, want to %.1f",
					pressed, y, got, y-page)
			}
		} else if got := position(); pressed != 1 || got != y {
			t.Errorf("Shift+Space pressed the row %d times and moved the region from %.1f to %.1f", pressed, y, got)
		}
	}
	pressed = 0
	screen.KeyPress(unison.KeyReturn, mod.Shift)
	if pressed != 0 {
		t.Errorf("with NoShiftActivate, Shift+Return pressed the row %d times", pressed)
	}
	pressed = 0
	screen.KeyPress(unison.KeySpace, 0)
	if pressed != 1 {
		t.Errorf("with NoShiftActivate, Space without Shift pressed the row %d times, want once", pressed)
	}
}

// A row that takes no focus on a press still acts on it, and the focus stays
// where it was.
func TestAListRowPressedWithoutTakingTheFocus(t *testing.T) {
	var field *kvitui.Field
	var row *kvitui.ListRow
	opened := 0
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		field = kvitui.NewField(ui)
		row = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Version 3"))
		row.Interactive, row.NoOpensMark = true, true
		row.OnActivate = func() { opened++ }
		return []unison.Paneler{kvitui.FullWidth(field), kvitui.Width(ui, kvitui.Px(300), row)}
	})
	fieldFocused := func() (f bool) {
		screen.Do(func() { f = field.Edit().Focused() })
		return f
	}
	for _, noFocus := range []bool{true, false} {
		opened = 0
		screen.Do(func() {
			row.NoFocusOnPress = noFocus
			field.Focus()
		})
		screen.Sync()
		if !fieldFocused() {
			t.Fatal("the field does not hold the focus before the press")
		}
		screen.Click(screen.PanelPoint(row, geom.NewPoint(250, 5)))
		var rowFocused bool
		screen.Do(func() { rowFocused = row.Focused() })
		if opened != 1 || fieldFocused() != noFocus || rowFocused == noFocus {
			t.Errorf("with NoFocusOnPress %v a press opened the row %d times; the field holds the focus %v, the row %v",
				noFocus, opened, fieldFocused(), rowFocused)
		}
	}
	// Tab still reaches it.
	screen.Do(func() {
		row.NoFocusOnPress = true
		field.Focus()
	})
	screen.Sync()
	screen.KeyPress(unison.KeyTab, 0)
	var rowFocused bool
	screen.Do(func() { rowFocused = row.Focused() })
	if !rowFocused {
		t.Error("Tab from the field does not reach a row that takes no focus on a press")
	}
}

// A row that yields its hover is untinted over a part with a hover of its
// own, such as a button, and tinted over its words and its padding.
func TestAListRowYieldsItsHoverToAButton(t *testing.T) {
	var row *kvitui.ListRow
	var words *kvitui.Label
	var button *kvitui.Button
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		words = kvitui.NewLabel(ui, "Words")
		button = kvitui.NewButton(ui, "Stop")
		row = kvitui.NewListRow(ui, words, button)
		row.Interactive, row.NoOpensMark, row.HoverWhole = true, true, true
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(300), kvitui.Stack(ui, nil, row, kvitui.NewLabel(ui, "Below")))}
	})
	hovered := func() (h bool) {
		screen.Do(func() { h = row.Hovered() })
		return h
	}
	for _, yields := range []bool{false, true} {
		screen.Do(func() { row.HoverYields = yields })
		for _, c := range []struct {
			what string
			at   geom.Point
			want bool
		}{
			{"its words", screen.PanelCenter(words), true},
			{"its padding", screen.PanelPoint(row, geom.NewPoint(250, 3)), true},
			{"a button inside it", screen.PanelCenter(button), !yields},
		} {
			screen.MouseMove(geom.NewPoint(5, 400), 0)
			screen.MouseMove(c.at, 0)
			screen.Sync()
			if got := hovered(); got != c.want {
				t.Errorf("with HoverYields %v, the row is hovered %v with the pointer on %s; want %v", yields, got, c.what, c.want)
			}
		}
	}
}

func TestWhatAListRowTellsAScreenReader(t *testing.T) {
	var folder, button *kvitui.ListRow
	screen, _ := rowsIn(t, func(ui *kvitui.UI) []*kvitui.ListRow {
		folder = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "docs"))
		folder.Interactive, folder.NoOpensMark, folder.Label = true, true, "docs"
		folder.Description, folder.Level, folder.Expandable, folder.Expanded = "Folder, 12 files", 2, true, true
		button = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "New track"))
		button.Interactive, button.NoOpensMark, button.Label = true, true, "New track"
		button.Accessibility.Role = role.Button
		return []*kvitui.ListRow{folder, button}
	})
	n := screen.AccessibilityNodeFor(folder)
	if n == nil || n.Role != role.ListItem || n.Name != "docs" || n.Description != "Folder, 12 files" || n.Level != 2 || !n.Expandable || !n.Expanded {
		t.Errorf("a tree's row: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(button); n == nil || n.Role != role.Button || n.Name != "New track" || n.Expandable {
		t.Errorf("a row given the role of a button: %+v", n)
	}
}
