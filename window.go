package kvitui

import (
	"math"
	"time"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

// Window is the application shell: a window with a header, an optional
// sidebar, a body and a status bar. Those four are assembled slightly
// differently in each application and the differences are all accidents.
// What it fixes in particular is the reflow: the sidebar collapses to its
// rail below the laptop width, and the window has a floor at the narrowest
// and shortest size the chrome holds, so a view never has to work out for
// itself when it is too narrow. The widths scale with the interface size: at
// twice the size the same screen needs twice the width.
//
// A rail is a strip of symbols with no words, and a reader who meets one
// cannot tell what the marks are. So while the pointer rests on the rail, or
// the keyboard is anywhere inside it, it opens to the full sidebar, over the
// body rather than pushing it aside, so the screen being read does not move
// underneath the reader.
type Window struct {
	*unison.Window
	ui *UI
	// SidebarVisible false hides the sidebar entirely, which is a reader's
	// choice, where a collapsed one is the window being narrow.
	SidebarVisible bool
	// SidebarExpandsOnHover opens the rail while the pointer rests on it or
	// the keyboard is inside it; on unless an application has a reason.
	SidebarExpandsOnHover bool
	// OnKeyDown is offered every key press before the focused control, for
	// the application's own shortcuts; true says it used the key. It is
	// here rather than in the window's KeyDownCallback, which Kvit uses to
	// tell keyboard focus from the rest.
	OnKeyDown func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool

	header, sidebar, body, status unison.Paneler
	panel                         *Panel        // the sidebar's ground, ruled on its right
	stage                         *unison.Panel // holds the body, and draws what spills out of its components
	railHovered, railFocused      bool
	drawnWidth                    float32 // the sidebar's width as drawn, moving towards its target
	easing                        bool
	popups                        []*Popup        // shown above everything, the last on top
	scrim                         *scrim          // behind the topmost modal popup
	spillHosts                    []*unison.Panel // the panels drawing what their components draw past their boxes
}

// NewWindow returns a shell window, the design width wide and 960 design
// pixels tall, with nothing in it yet.
func NewWindow(ui *UI, title string) (*Window, error) {
	uw, err := unison.NewWindow(title)
	if err != nil {
		return nil, err
	}
	w := &Window{Window: uw, ui: ui, SidebarVisible: true, SidebarExpandsOnHover: true}
	w.panel = NewPanel(ui)
	w.panel.RuleRight = true
	w.panel.SetLayout(&unison.FlexLayout{Columns: 1, HAlign: align.Fill, VAlign: align.Fill})
	w.panel.FocusChangeInHierarchyCallback = func(_, to *unison.Panel) {
		inside := false
		for p := to; p != nil; p = p.Parent() {
			if p == w.panel.AsPanel() {
				inside = true
				break
			}
		}
		// The keyboard is in the rail when a key moved the focus there, and
		// stays there while keys move it about inside. The focus unison
		// hands out when the window opens is not the reader's doing.
		w.setRail(w.railHovered, inside && (ui.keyTurn || w.railFocused))
	}
	w.hostSpills(w.panel.AsPanel())
	w.stage = unison.NewPanel()
	w.stage.SetLayout(&unison.FlexLayout{Columns: 1, HAlign: align.Fill, VAlign: align.Fill})
	w.stage.Accessibility.Role = role.None
	w.hostSpills(w.stage)
	content := uw.Content()
	content.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		painterFor(gc, ui).fill(content.ContentRect(true), ui.Theme.Tokens().WindowBackground)
	}
	content.SetLayout(windowLayout{w})
	uw.MinMaxContentSizeCallback = func() (minimum, maximum geom.Size) {
		m := ui.Interface
		return geom.NewSize(float32(m.WidthFloor()), float32(m.HeightFloor())),
			geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
	}
	w.watchPointer()
	uw.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if w.popupKeys(key) {
			return true
		}
		if !repeat && w.contextKeys(key, mods) {
			return true
		}
		return w.OnKeyDown != nil && w.OnKeyDown(key, mods, repeat)
	}
	uw.MouseDownCallback = func(where geom.Point, button, _ int, _ mod.Modifiers) bool {
		at := uw.Content().PointFromRoot(where)
		if w.popupPress(at) {
			return true
		}
		return button == unison.ButtonRight && w.contextPress(at)
	}
	if ui.windows == nil {
		ui.windows = map[*unison.Window]*Window{}
	}
	ui.windows[uw] = w
	uw.WillCloseCallback = func() { delete(ui.windows, uw) }
	ui.watchWindow(uw)
	m := ui.Interface
	uw.SetContentRect(geom.NewRect(0, 0, float32(m.WidthDrawn()), float32(m.Px(960))))
	return w, nil
}

// SetHeader puts the header across the top.
func (w *Window) SetHeader(p unison.Paneler) { w.header = p; w.rebuild() }

// SetSidebar puts the sidebar down the left. A *Sidebar is drawn as its rail
// whenever the window says so; anything else can ask SidebarCollapsed.
func (w *Window) SetSidebar(p unison.Paneler) { w.sidebar = p; w.rebuild() }

// SetBody puts the body in the space the rest leaves.
func (w *Window) SetBody(p unison.Paneler) { w.body = p; w.rebuild() }

// SetStatusBar puts the status bar across the bottom.
func (w *Window) SetStatusBar(p unison.Paneler) { w.status = p; w.rebuild() }

// rebuild puts the parts in the content in the order they are read: the
// header, the sidebar, the body and the status bar. unison draws the first
// child on top, so the sidebar, before the body, covers it when the rail
// opens, and takes a press on the open rail.
func (w *Window) rebuild() {
	content := w.Content()
	content.RemoveAllChildren()
	w.panel.RemoveAllChildren()
	w.stage.RemoveAllChildren()
	if w.header != nil {
		content.AddChild(w.header)
	}
	if w.sidebar != nil {
		w.sidebar.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
		w.panel.AddChild(w.sidebar)
		content.AddChild(w.panel)
	}
	if w.body != nil {
		w.body.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, VAlign: align.Fill, HGrab: true, VGrab: true})
		w.stage.AddChild(w.body)
		content.AddChild(w.stage)
	}
	if w.status != nil {
		content.AddChild(w.status)
	}
	w.addPopups(content)
	content.MarkForLayoutAndRedraw()
}

// SidebarNarrow reports whether the window is narrow enough that the layout
// keeps room for the rail rather than the sidebar. It does not change when
// the rail opens.
func (w *Window) SidebarNarrow() bool {
	return w.Content().FrameRect().Width < float32(w.ui.Interface.WidthLaptop())
}

// SidebarExpanded reports whether the rail is open over the body now.
func (w *Window) SidebarExpanded() bool {
	return w.SidebarNarrow() && w.SidebarVisible && w.SidebarExpandsOnHover && (w.railHovered || w.railFocused)
}

// SidebarCollapsed reports whether the sidebar is drawn as its rail now.
func (w *Window) SidebarCollapsed() bool { return w.SidebarNarrow() && !w.SidebarExpanded() }

// Narrow reports whether the window is below the narrowest width the chrome
// is laid out for.
func (w *Window) Narrow() bool {
	return w.Content().FrameRect().Width < float32(w.ui.Interface.WidthFloor())
}

// setRail records whether the pointer is on the rail and whether the
// keyboard is in it, and redoes the layout when that opens or closes it.
func (w *Window) setRail(hovered, focused bool) {
	before := w.SidebarExpanded()
	w.railHovered, w.railFocused = hovered, focused
	if w.SidebarExpanded() != before {
		w.Content().MarkForLayoutAndRedraw()
	}
}

// watchPointer follows the pointer over the whole window, since unison
// tells only the innermost panel under the pointer that it has arrived, and
// the rail is full of panels.
func (w *Window) watchPointer() {
	uw := w.Window
	over := func(where geom.Point) {
		in := false
		if w.sidebar != nil && w.SidebarVisible {
			in = where.In(w.panel.FrameRect())
		}
		w.setRail(in, w.railFocused)
	}
	enter, move, exit := uw.MouseEnterCallback, uw.MouseMoveCallback, uw.MouseExitCallback
	uw.MouseEnterCallback = func(where geom.Point, mods mod.Modifiers) bool {
		over(where)
		return enter != nil && enter(where, mods)
	}
	uw.MouseMoveCallback = func(where geom.Point, mods mod.Modifiers) bool {
		over(where)
		return move != nil && move(where, mods)
	}
	uw.MouseExitCallback = func() bool {
		w.setRail(false, w.railFocused)
		return exit != nil && exit()
	}
}

// sidebarTarget is the width the sidebar is drawn at when it settles, and
// sidebarRest the width the body is laid out against, which does not change
// when the rail opens.
func (w *Window) sidebarTarget() float32 {
	m := w.ui.Interface
	if w.SidebarCollapsed() {
		return float32(m.RailWidth())
	}
	return float32(m.SidebarWidth())
}

func (w *Window) sidebarRest() float32 {
	m := w.ui.Interface
	if w.sidebar == nil || !w.SidebarVisible {
		return 0
	}
	if w.SidebarNarrow() {
		return float32(m.RailWidth())
	}
	return float32(m.SidebarWidth())
}

// ease moves the drawn width of the sidebar towards its target over 120 ms,
// easing out, because a jump reads as a redraw fault rather than as the
// window adapting. With motion reduced it arrives at once. A target that
// changes partway, as when the pointer leaves an opening rail, is followed
// from wherever the width has got to.
func (w *Window) ease() {
	target := w.sidebarTarget()
	if w.drawnWidth == 0 || w.ui.Theme.MotionScale() == 0 {
		w.drawnWidth = target
		return
	}
	if w.drawnWidth == target || w.easing {
		return
	}
	w.easing = true
	from, start := w.drawnWidth, time.Now()
	var step func()
	step = func() {
		if to := w.sidebarTarget(); to != target {
			target, from, start = to, w.drawnWidth, time.Now()
		}
		duration := float64(120*w.ui.Theme.MotionScale()) * float64(time.Millisecond)
		t := 1.0
		if duration > 0 {
			t = min(1, float64(time.Since(start))/duration)
		}
		w.drawnWidth = from + (target-from)*float32(1-math.Pow(1-t, 3))
		if t >= 1 {
			w.drawnWidth, w.easing = target, false
		}
		w.Content().MarkForLayoutAndRedraw()
		if w.easing {
			unison.InvokeTaskAfter(step, 16*time.Millisecond)
		}
	}
	unison.InvokeTaskAfter(step, 16*time.Millisecond)
}

// windowLayout places the header across the top, the status bar across the
// bottom, the body beside the width the sidebar rests at, and the sidebar
// at the width it is drawn at, which is wider than that while the rail is
// open.
type windowLayout struct{ w *Window }

func (l windowLayout) LayoutSizes(target *unison.Panel, _ geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.w.ui.Interface
	floor := geom.NewSize(float32(m.WidthFloor()), float32(m.HeightFloor()))
	return floor, geom.NewSize(float32(m.WidthDrawn()), float32(m.Px(960))), geom.NewSize(unison.DefaultMaxSize, unison.DefaultMaxSize)
}

func (l windowLayout) PerformLayout(target *unison.Panel) {
	w, m := l.w, l.w.ui.Interface
	r := target.ContentRect(false)
	top, bottom := r.Y, r.Bottom()
	if w.header != nil {
		h := float32(m.HeaderHeight())
		w.header.AsPanel().SetFrameRect(geom.NewRect(r.X, top, r.Width, h))
		top += h
	}
	if w.status != nil {
		_, p, _ := w.status.AsPanel().Sizes(geom.NewSize(r.Width, 0))
		h := max(float32(m.StatusBarHeight()), p.Height)
		bottom -= h
		w.status.AsPanel().SetFrameRect(geom.NewRect(r.X, bottom, r.Width, h))
	}
	if w.sidebar != nil {
		if s, ok := w.sidebar.AsPanel().Self.(*Sidebar); ok {
			s.Collapsed = w.SidebarCollapsed()
		}
		w.ease()
		w.panel.Hidden = !w.SidebarVisible
		w.panel.SetFrameRect(geom.NewRect(r.X, top, w.drawnWidth, bottom-top))
	}
	if w.body != nil {
		left := r.X + w.sidebarRest()
		w.stage.SetFrameRect(geom.NewRect(left, top, r.Right()-left, bottom-top))
	}
	w.layoutPopups(r)
}
