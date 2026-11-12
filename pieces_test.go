package kvitui_test

import (
	"image"
	"runtime"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

func TestKeyTurnIsTrueWhileAKeyIsHandled(t *testing.T) {
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler { return nil })
	during, calls := false, 0
	screen.Do(func() {
		w.OnKeyDown = func(unison.KeyCode, mod.Modifiers, bool) bool {
			calls++
			during = ui.KeyTurn()
			return false
		}
	})
	if ui.KeyTurn() {
		t.Error("KeyTurn is true with no key pressed")
	}
	screen.KeyPress(unison.KeyA, 0)
	if calls != 1 || !during || ui.KeyTurn() {
		t.Errorf("a key was handled %d times with KeyTurn %v during and %v after", calls, during, ui.KeyTurn())
	}
}

// A window that is not Kvit's is watched once however often it is asked, and
// the key handler it had still runs.
func TestWatchKeysWrapsAWindowOnce(t *testing.T) {
	screen, ui := layoutSession(t)
	var wnd *unison.Window
	earlier, later := 0, 0
	var duringEarlier, duringLater bool
	screen.Do(func() {
		wnd = screen.FocusedWindow()
		wnd.KeyDownCallback = func(unison.KeyCode, mod.Modifiers, bool) bool {
			earlier++
			duringEarlier = ui.KeyTurn()
			return false
		}
		ui.WatchKeys(wnd)
	})
	screen.KeyPress(unison.KeyA, 0)
	if earlier != 1 || !duringEarlier {
		t.Fatalf("the window's earlier key handler ran %d times, with KeyTurn %v", earlier, duringEarlier)
	}
	// A handler set over the watch is not wrapped by a second call, so no
	// turn is reported inside it.
	screen.Do(func() {
		wnd.KeyDownCallback = func(unison.KeyCode, mod.Modifiers, bool) bool {
			later++
			duringLater = ui.KeyTurn()
			return false
		}
		ui.WatchKeys(wnd)
	})
	screen.KeyPress(unison.KeyA, 0)
	if later != 1 || duringLater || earlier != 1 {
		t.Errorf("after a second WatchKeys the later handler ran %d times with KeyTurn %v, the earlier %d times", later, duringLater, earlier)
	}
}

func TestKeyFocusTellsTheKeyboardFromAPress(t *testing.T) {
	var panel, other *unison.Panel
	var focus kvitui.KeyFocus
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		panel, other = box(100, 40), box(100, 40)
		for _, p := range []*unison.Panel{panel, other} {
			p.SetFocusable(true)
		}
		panel.MouseDownCallback = func(geom.Point, int, int, mod.Modifiers) bool {
			focus.Pressed()
			panel.RequestFocus()
			return true
		}
		panel.GainedFocusCallback = func() { focus.Gained(ui) }
		panel.LostFocusCallback = focus.Lost
		return []unison.Paneler{panel, other}
	})
	shows := func() (s bool) {
		screen.Do(func() { s = focus.Shows(panel) })
		return s
	}
	screen.Do(func() { ui.WatchKeys(panel.Window()) })
	var focused bool
	screen.Do(func() { focused = panel.Focused() })
	if !focused || shows() {
		t.Fatalf("as the window opened the panel is focused %v and shows keyboard focus %v", focused, shows())
	}
	screen.Do(func() { other.RequestFocus() })
	screen.KeyPress(unison.KeyTab, 0)
	if !shows() {
		t.Error("focus moved by Tab does not show")
	}
	// A press on the panel while it shows ends it.
	screen.Click(screen.PanelCenter(panel))
	if shows() {
		t.Error("a press on a panel showing keyboard focus leaves it showing")
	}
	screen.Do(func() { other.RequestFocus() })
	screen.Click(screen.PanelCenter(panel))
	if shows() {
		t.Error("focus given by a press shows")
	}
	screen.Do(func() { other.RequestFocus() })
	screen.Sync()
	screen.Do(func() {
		focus.Asked()
		panel.RequestFocus()
	})
	screen.Sync()
	if !shows() {
		t.Error("focus the component asked for does not show")
	}
}

// same reports whether two boxes of a picture, of one size, hold the same
// pixels.
func same(img image.Image, a, b geom.Rect) bool {
	for y := 0; y < int(a.Height); y++ {
		for x := 0; x < int(a.Width); x++ {
			if img.At(int(a.X)+x, int(a.Y)+y) != img.At(int(b.X)+x, int(b.Y)+y) {
				return false
			}
		}
	}
	return true
}

func TestDrawSurfaceDrawsAPopoversGround(t *testing.T) {
	t.Setenv(uitest.ScaleVariable, "1")
	var pop *kvitui.Popover
	var card, round *unison.Panel
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		surface := func(radius kvitui.Measure) *unison.Panel {
			p := unison.NewPanel()
			p.AddChild(box(100, 40))
			p.SetBorder(kvitui.Padding(ui, kvitui.SizeSpaceLoose))
			p.SetLayout(kvitui.OwnSizeLayout{CountInsets: true})
			p.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) { kvitui.DrawSurface(gc, ui, p.ContentRect(true), radius) }
			return p
		}
		pop = kvitui.NewPopover(ui, "Surface", box(100, 40))
		card, round = surface(nil), surface(kvitui.Px(16))
		sized := func(p unison.Paneler) *unison.Panel { return kvitui.Sized(ui, kvitui.Px(160), kvitui.Px(80), p) }
		return kvitui.Stack(ui, kvitui.SizeSpace, sized(pop), sized(card), sized(round))
	})
	img := s.Capture()
	rect := func(p unison.Paneler) (r geom.Rect) {
		s.Do(func() { r = p.AsPanel().RectToRoot(p.AsPanel().ContentRect(true)) })
		return r
	}
	a, b, c := rect(pop), rect(card), rect(round)
	if a.Size != b.Size || a.Size != c.Size || a.Width < 100 {
		t.Fatalf("the surfaces are %v, %v and %v", a.Size, b.Size, c.Size)
	}
	if !same(img, a, b) {
		t.Error("DrawSurface with no radius does not draw what a popover draws")
	}
	// A larger radius changes the corners only.
	corner := func(r geom.Rect) geom.Rect { return geom.NewRect(r.X, r.Y, 16, 16) }
	middle := func(r geom.Rect) geom.Rect { return geom.NewRect(r.X+20, r.Y, r.Width-40, r.Height) }
	if same(img, corner(a), corner(c)) {
		t.Error("a larger radius draws the same corner")
	}
	if !same(img, middle(a), middle(c)) {
		t.Error("a larger radius changes the surface between its corners")
	}
}

func TestRoveStep(t *testing.T) {
	members := []*unison.Panel{box(1, 1), box(1, 1), box(1, 1), box(1, 1)}
	members[2].SetEnabled(false)
	step := func(from int, key unison.KeyCode, vertical, wrap bool) (int, bool) {
		to, moved := kvitui.RoveStep(members, members[from], key, vertical, wrap)
		for i, m := range members {
			if m == to {
				return i, moved
			}
		}
		return -1, moved
	}
	for _, c := range []struct {
		what           string
		from           int
		key            unison.KeyCode
		vertical, wrap bool
		to             int
		moved          bool
	}{
		{"down to the next", 0, unison.KeyDown, true, false, 1, true},
		{"down past a disabled member", 1, unison.KeyDown, true, false, 3, true},
		{"up to the one before", 1, unison.KeyUp, true, false, 0, true},
		{"right in a group that runs across", 0, unison.KeyRight, false, false, 1, true},
		{"an arrow across the group's way", 0, unison.KeyRight, true, false, -1, false},
		{"a key that is no arrow", 0, unison.KeyA, true, false, -1, false},
		{"down at the end of a list that stops", 3, unison.KeyDown, true, false, -1, true},
		{"up at the start of a list that stops", 0, unison.KeyUp, true, false, -1, true},
		{"down at the end of a list that wraps", 3, unison.KeyDown, true, true, 0, true},
		{"up at the start of a list that wraps", 0, unison.KeyUp, true, true, 3, true},
		{"Home", 3, unison.KeyHome, true, false, 0, true},
		{"End", 0, unison.KeyEnd, true, false, 3, true},
	} {
		if to, moved := step(c.from, c.key, c.vertical, c.wrap); to != c.to || moved != c.moved {
			t.Errorf("%s: from %d went to %d, moved %v; want %d, %v", c.what, c.from, to, moved, c.to, c.moved)
		}
	}
	if _, moved := kvitui.RoveStep(members[:1], members[0], unison.KeyDown, true, true); moved {
		t.Error("a member alone moved")
	}
	if _, moved := kvitui.RoveStep(members[:2], members[3], unison.KeyDown, true, true); moved {
		t.Error("a panel that is not a member moved")
	}
}

func TestAPopoverOpenedModal(t *testing.T) {
	var behind, opener *kvitui.Button
	var first, second *kvitui.Check
	var pop *kvitui.Popover
	clicks := 0
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		opener = kvitui.NewButton(ui, "Open")
		behind = kvitui.NewButton(ui, "Behind")
		behind.OnClick = func() { clicks++ }
		return []unison.Paneler{kvitui.Left(opener), kvitui.Left(behind)}
	})
	open := func(set func(p *kvitui.Popover)) {
		screen.Do(func() {
			opener.Focus()
			first, second = kvitui.NewCheck(ui, "First"), kvitui.NewCheck(ui, "Second")
			pop = kvitui.NewPopover(ui, "Options", first, second)
			if set != nil {
				set(pop)
			}
			pop.OpenModal(opener, kvitui.PlaceCentred(nil))
		})
		screen.Sync()
	}
	opened := func() (o bool) {
		screen.Do(func() { o = pop.Opened() })
		return o
	}
	focused := func(p unison.Paneler) (f bool) {
		screen.Do(func() { f = p.AsPanel().Focused() })
		return f
	}
	outside := func() { screen.Click(screen.PanelCenter(behind)) }

	// By default: the first control has the keyboard, Tab stays inside, a
	// press outside closes it and presses nothing, and the focus goes back.
	open(nil)
	if !focused(first) {
		t.Error("the first control of a modal popover does not have the keyboard")
	}
	screen.KeyPress(unison.KeyTab, 0)
	screen.KeyPress(unison.KeyTab, 0)
	if !focused(first) {
		t.Error("Tab twice among two controls did not come back to the first")
	}
	outside()
	if opened() || clicks != 0 || !focused(opener) {
		t.Errorf("a press outside left the popover open %v, pressed the button behind %d times, focus back %v", opened(), clicks, focused(opener))
	}
	open(nil)
	screen.KeyPress(unison.KeyEscape, 0)
	if opened() {
		t.Error("Escape left a modal popover open")
	}

	// The caller's functions run in place of closing.
	escapes, presses := 0, 0
	open(func(p *kvitui.Popover) {
		p.OnEscape = func() { escapes++ }
		p.OnPressOutside = func() { presses++ }
		p.FirstFocus = func() *unison.Panel { return second.AsPanel() }
	})
	if !focused(second) {
		t.Error("the control FirstFocus names does not have the keyboard")
	}
	screen.KeyPress(unison.KeyEscape, 0)
	outside()
	if !opened() || escapes != 1 || presses != 1 || clicks != 0 {
		t.Errorf("with functions for both, the popover is open %v after Escape ran %d and a press outside ran %d, pressing behind %d times", opened(), escapes, presses, clicks)
	}
	screen.Do(func() { pop.Close() })

	// A press outside that is ignored.
	presses = 0
	open(func(p *kvitui.Popover) {
		p.PressOutsideIgnored = true
		p.OnPressOutside = func() { presses++ }
	})
	outside()
	if !opened() || presses != 0 || clicks != 0 {
		t.Errorf("an ignored press outside left the popover open %v, ran OnPressOutside %d times and pressed behind %d times", opened(), presses, clicks)
	}
	screen.Do(func() { pop.Close() })

	// A press outside that passes: it closes the popover and presses what
	// is behind; while a function of the caller's keeps the popover open,
	// the scrim still takes the press.
	open(func(p *kvitui.Popover) { p.ModalPress = kvitui.ModalPressPasses })
	outside()
	if opened() || clicks != 1 {
		t.Errorf("a press outside that passes left the popover open %v and pressed the button behind %d times", opened(), clicks)
	}
	clicks, presses = 0, 0
	open(func(p *kvitui.Popover) {
		p.ModalPress = kvitui.ModalPressPasses
		p.OnPressOutside = func() { presses++ }
	})
	outside()
	if !opened() || presses != 1 || clicks != 0 {
		t.Errorf("a press outside that passes, with a function that keeps the popover open: open %v, function ran %d, pressed behind %d", opened(), presses, clicks)
	}
	screen.Do(func() { pop.Close() })

	// MayClose holds it open against Escape and Close.
	may := false
	open(func(p *kvitui.Popover) { p.MayClose = func() bool { return may } })
	screen.KeyPress(unison.KeyEscape, 0)
	screen.Do(func() { pop.Close() })
	if !opened() {
		t.Error("a popover that may not close closed")
	}
	may = true
	screen.KeyPress(unison.KeyEscape, 0)
	if opened() {
		t.Error("a popover that may close stayed open on Escape")
	}
	if len(w.Popups()) != 0 {
		t.Errorf("%d popups are left", len(w.Popups()))
	}
}

func TestAMoverMovesAPopupByItsBand(t *testing.T) {
	var surface *unison.Panel
	var inBand *kvitui.Button
	var mover *kvitui.Mover
	pressed := 0
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler { return nil })
	screen.Do(func() {
		inBand = kvitui.NewButton(ui, "Close")
		inBand.OnClick = func() { pressed++ }
		surface = unison.NewPanel()
		surface.AddChild(inBand)
		surface.SetLayout(kvitui.FillLayout())
		surface.SetSizer(func(geom.Size) (minSize, prefSize, maxSize geom.Size) {
			return geom.NewSize(300, 200), geom.NewSize(300, 200), geom.NewSize(300, 200)
		})
		surface.SetLayout(nil)
		inBand.SetFrameRect(geom.NewRect(220, 4, 70, 24))
		mover = kvitui.NewMover(ui, surface, func() float32 { return 32 })
		w.Show(&kvitui.Popup{Panel: surface, Modal: true, Place: func(bounds geom.Rect, size geom.Size) geom.Rect {
			return mover.Place(bounds, kvitui.PlaceCentred(nil)(bounds, size))
		}})
	})
	screen.Sync()
	frame := func() (r geom.Rect) {
		screen.Do(func() { r = surface.FrameRect() })
		return r
	}
	var bounds geom.Rect
	screen.Do(func() { bounds = w.Content().ContentRect(false) })
	band := func() geom.Point { return screen.PanelPoint(surface, geom.NewPoint(40, 16)) }
	below := func() geom.Point { return screen.PanelPoint(surface, geom.NewPoint(40, 100)) }
	centred := frame()
	if centred.Size != geom.NewSize(300, 200) || mover.Moved() {
		t.Fatalf("the surface opened at %v, moved %v", centred, mover.Moved())
	}
	// A drag below the band, and a press on the button in it, move nothing.
	at := below()
	screen.Drag(at, geom.NewPoint(at.X-100, at.Y-50), 4)
	screen.Click(screen.PanelCenter(inBand))
	if frame() != centred || pressed != 1 {
		t.Errorf("a drag below the band and a press on the band's button left the surface at %v and pressed the button %d times", frame(), pressed)
	}
	at = band()
	screen.Drag(at, geom.NewPoint(at.X-100, at.Y-50), 4)
	moved := frame()
	if abs(moved.X-(centred.X-100)) > 1 || abs(moved.Y-(centred.Y-50)) > 1 || !mover.Moved() {
		t.Fatalf("dragging the band 100 left and 50 up moved the surface from %v to %v", centred, moved)
	}
	// It stops at the bounds.
	at = band()
	screen.Drag(at, geom.NewPoint(at.X+3000, at.Y+3000), 4)
	if pinned := frame(); abs(pinned.Right()-bounds.Right()) > 1 || abs(pinned.Bottom()-bounds.Bottom()) > 1 {
		t.Errorf("dragged far past the corner, the surface is at %v in %v", pinned, bounds)
	}
	at = band()
	screen.Drag(at, geom.NewPoint(at.X-400, at.Y-300), 4)
	moved = frame()
	// It keeps its place, from the corner of the bounds, when the window is
	// laid out again, and stays inside a smaller one.
	screen.Do(func() { w.Content().MarkForLayoutAndRedraw() })
	screen.Sync()
	if frame() != moved {
		t.Errorf("laid out again, the surface went from %v to %v", moved, frame())
	}
	screen.Do(func() { w.SetContentRect(geom.NewRect(0, 0, 1100, 700)) })
	screen.Sync()
	screen.Do(func() { bounds = w.Content().ContentRect(false) })
	if now := frame(); now.Right() > bounds.Right()+1 || now.Bottom() > bounds.Bottom()+1 || (now.X != moved.X && abs(now.Right()-bounds.Right()) > 1) {
		t.Errorf("in a smaller window the surface went from %v to %v, inside %v", moved, now, bounds)
	}
	// Reset puts it back in the middle.
	screen.Do(func() { mover.Reset() })
	screen.Sync()
	centred = frame()
	if mid := bounds.Center(); abs(centred.Center().X-mid.X) > 1 || abs(centred.Center().Y-mid.Y) > 1 || mover.Moved() {
		t.Errorf("after Reset the surface is at %v in %v, moved %v", centred, bounds, mover.Moved())
	}

	// Escape during a drag does nothing, unless it is to abandon the move.
	for _, abandons := range []bool{false, true} {
		screen.Do(func() {
			mover.EscapeAbandons = abandons
			mover.Reset()
		})
		screen.Sync()
		from := band()
		to := geom.NewPoint(from.X-80, from.Y-40)
		screen.MouseDown(from, unison.ButtonLeft, 0)
		screen.MouseMove(to, 0)
		dragged := frame()
		screen.KeyPress(unison.KeyEscape, 0)
		after := frame()
		screen.MouseUp(to, unison.ButtonLeft, 0)
		if dragged == centred {
			t.Fatal("the drag did not move the surface")
		}
		if !abandons && after != dragged {
			t.Errorf("Escape during a drag moved the surface from %v to %v", dragged, after)
		}
		if abandons && (after != centred || mover.Moved()) {
			t.Errorf("Escape abandoning a drag left the surface at %v, not back at %v", after, centred)
		}
	}
}

func TestDialogButtonOrderFollowsTheForm(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	safe, first, second := box(1, 1), box(1, 1), box(1, 1)
	order := func() []unison.Paneler { return kvitui.DialogButtonOrder(ui, safe, first, second) }
	ui.SetForm(kvitui.FormWindows)
	if got := order(); len(got) != 3 || got[0] != first || got[1] != second || got[2] != safe {
		t.Error("on the Windows form the buttons are not the actions and then the safe button")
	}
	ui.SetForm(kvitui.FormMac)
	if got := order(); len(got) != 3 || got[0] != safe || got[1] != first || got[2] != second {
		t.Error("on the macOS form the buttons are not the safe button and then the actions")
	}
}

func TestAWindowTellsListenersWhereThePointerIs(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	ui.Theme.SetReducedMotion(true)
	var w *kvitui.Window
	var other *unison.Window
	var target *unison.Panel
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 1600, Height: 1000},
		unison.StartupFinishedCallback(func() {
			if w, err = kvitui.NewWindow(ui, "shell"); err != nil {
				t.Error(err)
				return
			}
			target = box(200, 100)
			w.SetBody(kvitui.Stack(ui, nil, kvitui.OwnWidth(target)))
			w.SetContentRect(geom.NewRect(0, 0, float32(ui.Interface.WidthLaptop()-100), 600))
			w.ToFront()
			if other, err = unison.NewWindow("other"); err != nil {
				t.Error(err)
				return
			}
			other.SetFrameRect(geom.NewRect(0, 700, 400, 200))
			other.ToFront()
			w.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.Sync()

	type told struct {
		where  geom.Point
		inside bool
	}
	var first, second, replaced, dropped []told
	owner, otherOwner := unison.NewPanel(), unison.NewPanel()
	screen.Do(func() {
		ui.OnPointer(w.Window, owner, func(where geom.Point, inside bool) { replaced = append(replaced, told{where, inside}) })
		// A second registration for the same owner replaces the first.
		ui.OnPointer(w.Window, owner, func(where geom.Point, inside bool) { first = append(first, told{where, inside}) })
		ui.OnPointer(other, otherOwner, func(where geom.Point, inside bool) { second = append(second, told{where, inside}) })
		// An owner nothing else holds is dropped, and its listener with it.
		ui.OnPointer(w.Window, unison.NewPanel(), func(where geom.Point, inside bool) { dropped = append(dropped, told{where, inside}) })
	})
	at := screen.PanelCenter(target)
	screen.MouseMove(at, 0)
	if len(first) == 0 || !first[len(first)-1].inside || len(second) != 0 || len(replaced) != 0 {
		t.Fatalf("a move in the first window told its listener %v, the replaced one %v, the other window's %v", first, replaced, second)
	}
	var over, on bool
	screen.Do(func() {
		last := first[len(first)-1]
		over, on = kvitui.PointerOver(target, last.where, last.inside), kvitui.PointerOn(target, last.where, last.inside)
	})
	if !over || !on {
		t.Errorf("with the pointer on the panel, PointerOver is %v and PointerOn %v", over, on)
	}

	// Under a popup that covers the panel the pointer is over it, not on it.
	var hide func()
	screen.Do(func() {
		hide = w.Show(&kvitui.Popup{Panel: box(10, 10), Passive: true, Place: func(geom.Rect, geom.Size) geom.Rect {
			return w.Content().RectFromRoot(target.RectToRoot(target.ContentRect(true)))
		}})
	})
	screen.Sync()
	screen.Do(func() {
		last := first[len(first)-1]
		over, on = kvitui.PointerOver(target, last.where, last.inside), kvitui.PointerOn(target, last.where, last.inside)
		hide()
	})
	if !over || on {
		t.Errorf("with the panel under a popup, PointerOver is %v and PointerOn %v", over, on)
	}

	// A dropped owner's listener is removed after a collection.
	before := len(dropped)
	if before == 0 {
		t.Error("the listener of the owner nothing holds was never told")
	}
	runtime.GC()
	runtime.GC()
	screen.MouseMove(geom.NewPoint(at.X+3, at.Y), 0)
	screen.MouseMove(geom.NewPoint(at.X+6, at.Y), 0)
	if len(dropped) != before {
		t.Errorf("the listener of a dropped owner was told %d more times", len(dropped)-before)
	}
	// The other window tells its own listener, and leaving a window is told.
	told1 := len(first)
	screen.MouseMove(geom.NewPoint(200, 800), 0)
	if len(second) == 0 || !second[len(second)-1].inside {
		t.Errorf("a move in the other window told its listener %v", second)
	}
	if len(first) != told1+1 || first[len(first)-1].inside {
		t.Errorf("as the pointer left the first window its listener was told %v", first[told1:])
	}
	runtime.KeepAlive(owner)
	runtime.KeepAlive(otherOwner)
}

// The window's own watch of the pointer, which opens the sidebar's rail,
// goes on working with a listener registered over it.
func TestTheRailStillOpensWithAPointerListener(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	ui.Theme.SetReducedMotion(true)
	var w *kvitui.Window
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 1600, Height: 1000},
		unison.StartupFinishedCallback(func() {
			if w, err = kvitui.NewWindow(ui, "shell"); err != nil {
				t.Error(err)
				return
			}
			w.SetSidebar(kvitui.NewSidebar(ui, kvitui.NewSidebarItem(ui, "Everyday", "wallet"), kvitui.NewSidebarItem(ui, "Savings", "bank")))
			w.SetBody(kvitui.NewRegion(ui, kvitui.NewSlimRow(ui, "Checking")))
			w.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.Sync()
	told := 0
	owner := unison.NewPanel()
	screen.Do(func() {
		w.SetContentRect(geom.NewRect(0, 0, float32(ui.Interface.WidthLaptop()-100), 900))
		w.Content().MarkForLayoutRecursively()
		w.ValidateLayout()
		ui.OnPointer(w.Window, owner, func(geom.Point, bool) { told++ })
	})
	screen.Sync()
	expanded := func() (e bool) {
		screen.Do(func() { e = w.SidebarExpanded() })
		return e
	}
	screen.MouseMove(geom.NewPoint(600, 200), 0)
	screen.Sync()
	told = 0
	if expanded() {
		t.Fatal("the rail is open before the pointer is on it")
	}
	screen.MouseMove(geom.NewPoint(float32(ui.Interface.RailWidth())/2, 200), 0)
	screen.Sync()
	if !expanded() || told == 0 {
		t.Errorf("with the pointer on the rail it is open %v, and the listener was told %d times", expanded(), told)
	}
	screen.MouseMove(geom.NewPoint(600, 200), 0)
	screen.Sync()
	if expanded() {
		t.Error("the rail stayed open after the pointer left it")
	}
	runtime.KeepAlive(owner)
}
