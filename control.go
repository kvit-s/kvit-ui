package kvitui

import (
	"slices"
	"time"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/pathop"
	"github.com/richardwilkes/unison/enums/role"
)

// TypeRole names one of the interface's seven type roles: what a run of
// chrome text is, rather than how big it is.
type TypeRole int

// The seven roles, smallest first.
const (
	RoleBody     TypeRole = iota // 12 px at the default size: row text, prose
	RoleCaption                  // 10: kind tags, counts
	RoleSmall                    // 11: chip labels, sub-lines
	RoleStrong                   // 13: a name, an emphasised row
	RoleTitle                    // 15: a section heading
	RoleHeadline                 // 17: a pane title, the wordmark
	RoleDisplay                  // 20: a page title
)

// Size is a role's pixel size at the current interface size.
func (u *UI) Size(r TypeRole) int {
	m := u.Interface
	switch r {
	case RoleCaption:
		return m.Caption()
	case RoleSmall:
		return m.Small()
	case RoleStrong:
		return m.Strong()
	case RoleTitle:
		return m.Title()
	case RoleHeadline:
		return m.Headline()
	case RoleDisplay:
		return m.Display()
	}
	return m.Body()
}

// control is what the interactive components share: hover and press state,
// activation by the pointer and by the keyboard, and a focus ring that shows
// only when the focus arrived from the keyboard. A ring on every click is
// noise; a keyboard user needs to see where they are. The focus a window
// hands its first control when it opens came from neither, and shows no
// ring either, as in .
type control struct {
	unison.Panel
	ui            *UI
	hovered       bool
	pressed       bool
	keyboardFocus bool // the focus came from the keyboard, so the ring shows
	pointerFocus  bool // set while a press is giving the control the focus
	askedFocus    bool // set while Focus is giving the control the focus
	activate      func()
	keys          []unison.KeyCode // the keys that activate the control
	ringRadius    func() float32   // the corner radius of the shape the ring surrounds
	noRing        bool             // shows keyboard focus its own way instead of a ring
	// tip gives the tooltip's words, "" for none: shown beside the control
	// on pointer hover and on keyboard focus.
	tip func() string
	// tipArea, when set, is the part of the control the pointer has to be
	// over for the tooltip, as a row's trailing chevron is.
	tipArea func() geom.Rect
	tipGen  int    // counts tooltip changes, so a late show can tell it is stale
	tipHide func() // takes the tooltip away while one is shown
}

// initControl wires the callbacks. activate runs on a click, on one of keys,
// and on a screen reader's press.
func (c *control) initControl(ui *UI, activate func(), keys ...unison.KeyCode) {
	c.ui = ui
	c.activate = activate
	c.keys = keys
	c.SetFocusable(true)
	// Laid out in a window means shown in it, before anybody can press a key
	// there, which is when the window has to start watching for keys.
	c.FrameChangeCallback = func() {
		if w := c.Window(); w != nil {
			ui.watchWindow(w)
		}
	}
	c.MouseEnterCallback = func(where geom.Point, _ mod.Modifiers) bool {
		c.hovered = true
		c.MarkForRedraw()
		c.pointerAt(where)
		return false
	}
	c.MouseMoveCallback = func(where geom.Point, _ mod.Modifiers) bool {
		c.pointerAt(where)
		return false
	}
	c.MouseExitCallback = func() bool {
		c.hovered = false
		c.MarkForRedraw()
		if !c.KeyboardFocus() {
			c.hideTip()
		}
		return false
	}
	// Outside a Kvit window there is no layer to show a tooltip in, and
	// unison shows its own on pointer hover instead.
	c.UpdateTooltipCallback = func(geom.Point, geom.Rect) geom.Rect {
		c.Tooltip = nil
		if c.tip != nil && ui.windowOf(c) == nil {
			if say := c.tip(); say != "" {
				c.Tooltip = newTooltip(ui, say, "")
			}
		}
		return c.RectToRoot(c.ContentRect(true))
	}
	c.MouseDownCallback = func(_ geom.Point, button, _ int, _ mod.Modifiers) bool {
		if !c.Enabled() || button != unison.ButtonLeft {
			return false
		}
		c.pressed = true
		c.hideTip()
		// The focus arrives after this returns; the marker says where it came
		// from until then. A click also ends any keyboard focus shown.
		c.pointerFocus = true
		c.keyboardFocus = false
		c.RequestFocus()
		c.MarkForRedraw()
		return true
	}
	c.MouseDragCallback = func(where geom.Point, _ int, _ mod.Modifiers) bool {
		if inside := where.In(c.ContentRect(false)); inside != c.pressed {
			c.pressed = inside
			c.MarkForRedraw()
		}
		return true
	}
	c.MouseUpCallback = func(where geom.Point, _ int, _ mod.Modifiers) bool {
		was := c.pressed
		c.pressed = false
		c.pointerFocus = false
		c.MarkForRedraw()
		if was && c.Enabled() && where.In(c.ContentRect(false)) {
			c.fire()
		}
		return true
	}
	c.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if !c.Enabled() || repeat || mods.CommandDown() || mods.OptionDown() || !slices.Contains(c.keys, key) {
			return false
		}
		c.fire()
		return true
	}
	c.GainedFocusCallback = func() {
		c.keyboardFocus = !c.pointerFocus && (c.askedFocus || ui.keyTurn)
		c.pointerFocus, c.askedFocus = false, false
		c.MarkForRedraw()
		if c.keyboardFocus {
			c.showTipSoon()
		}
		if w := c.Window(); w != nil {
			ui.watchWindow(w)
			w.MarkForRedraw()
		}
	}
	c.LostFocusCallback = func() {
		c.keyboardFocus = false
		if !c.hovered {
			c.hideTip()
		}
		c.MarkForRedraw()
		if w := c.Window(); w != nil {
			w.MarkForRedraw()
		}
	}
}

// pointerAt shows the tooltip after the pointer rests, where the pointer is
// over the part that has one.
func (c *control) pointerAt(where geom.Point) {
	if c.tip == nil {
		return
	}
	if c.tipArea != nil && !where.In(c.tipArea()) {
		if !c.KeyboardFocus() {
			c.hideTip()
		}
		return
	}
	if c.tipHide == nil {
		c.showTipSoon()
	}
}

// showTipSoon shows the tooltip after a pause, unless something changes
// before then.
func (c *control) showTipSoon() {
	if c.tip == nil || c.ui.windowOf(c) == nil {
		return
	}
	c.tipGen++
	gen := c.tipGen
	unison.InvokeTaskAfter(func() {
		if gen == c.tipGen && c.tipHide == nil {
			c.showTip()
		}
	}, 500*time.Millisecond)
}

// showTip puts the tooltip beside the control, for long enough to read: at
// least three seconds, and longer for more words. A tooltip that vanishes
// while somebody is reading it is worse than none.
func (c *control) showTip() {
	w := c.ui.windowOf(c)
	say := c.tip()
	if w == nil || say == "" {
		return
	}
	hide := w.Show(&Popup{Panel: newTooltip(c.ui, say, ""), Place: PlaceBeside(c.ui, c), Anchor: c, Passive: true})
	c.tipHide = hide
	gen := c.tipGen
	unison.InvokeTaskAfter(func() {
		if gen == c.tipGen {
			c.hideTip()
		}
	}, time.Duration(max(3000, len([]rune(say))*60))*time.Millisecond)
}

// hideTip takes the tooltip away, and cancels one about to show.
func (c *control) hideTip() {
	c.tipGen++
	if c.tipHide != nil {
		c.tipHide()
		c.tipHide = nil
	}
}

func (c *control) fire() {
	if c.activate != nil {
		c.activate()
	}
}

// Focus gives the control the keyboard focus, with its ring, as a Tab to it
// would. A control not in a window yet takes it once it is.
func (c *control) Focus() {
	if c.Window() != nil {
		c.askedFocus = true
		c.RequestFocus()
		return
	}
	unison.InvokeTask(func() {
		if c.Window() != nil {
			c.askedFocus = true
			c.RequestFocus()
		}
	})
}

// KeyboardFocus reports whether the control holds the focus and it arrived
// from the keyboard, which is when its focus ring shows.
func (c *control) KeyboardFocus() bool { return c.keyboardFocus && c.Focused() }

// Hovered reports whether the pointer is over the control.
func (c *control) Hovered() bool { return c.hovered }

// Pressed reports whether the control is being pressed.
func (c *control) Pressed() bool { return c.pressed }

// ringOwnerKey is the client-data key under which a unison widget, such as
// the field inside a Field, names the Kvit component that draws its ring.
const ringOwnerKey = "kvitui.ringOwner"

// ringed is a component whose keyboard focus the window draws a ring around.
type ringed interface {
	focusRing() (panel *unison.Panel, radius float32, show bool)
}

func (c *control) focusRing() (*unison.Panel, float32, bool) {
	r := float32(c.ui.Interface.RadiusControl())
	if c.ringRadius != nil {
		r = c.ringRadius()
	}
	return c.AsPanel(), r, c.KeyboardFocus() && !c.noRing
}

// windowWatchedKey is the client-data key under which a window records that
// it draws focus rings and watches keys.
const windowWatchedKey = "kvitui.watched"

// watchWindow makes a window draw the focus ring of whichever Kvit control
// holds its keyboard focus, and tell the controls when a key is being
// handled, so a control can tell focus moved by Tab from the focus the window
// hands out when it opens.
//
// The ring sits outside the control's own box, as the library draws it,
// and unison clips each panel's drawing to its bounds, so it is drawn by the
// window over everything, clipped to the area the control's container shows.
//
// The key watch wraps the window's KeyDownCallback. An application that sets
// its own after Kvit controls are shown in the window has to call the one it
// replaces.
func (u *UI) watchWindow(w *unison.Window) {
	// The window is marked in its own client data rather than listed by the
	// UI, which would keep every window that ever closed (menus, popovers
	// and dialogs among them) with all its panels for the life of the
	// program.
	if w.ClientData()[windowWatchedKey] != nil {
		return
	}
	w.ClientData()[windowWatchedKey] = true
	previousKey := w.KeyDownCallback
	w.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		// unison moves the focus for Tab after this returns, within the same
		// key press; the turn ends once the press has been handled.
		u.keyTurn = true
		unison.InvokeTask(func() { u.keyTurn = false })
		if previousKey != nil {
			return previousKey(key, mods, repeat)
		}
		return false
	}
	content := w.Content()
	previous := content.DrawOverCallback
	content.DrawOverCallback = func(gc *unison.Canvas, rect geom.Rect) {
		if previous != nil {
			previous(gc, rect)
		}
		// CurrentFocus rather than Focus, which hands the focus to the first
		// control it finds when nothing holds it.
		focus := w.CurrentFocus()
		if focus == nil {
			return
		}
		r, ok := focus.Self.(ringed)
		if !ok {
			// A unison widget a Kvit component is built around names the
			// component in its client data.
			if r, ok = focus.ClientData()[ringOwnerKey].(ringed); !ok {
				return
			}
		}
		p, radius, show := r.focusRing()
		if !show {
			return
		}
		m := u.Interface
		width := float32(m.FocusRingWidth())
		box := content.RectFromRoot(p.RectToRoot(p.ContentRect(true)))
		ring := box.Inset(geom.NewUniformInsets(-width / 2))
		// Clip to what the control's containers show, grown by the ring.
		clip := box.Inset(geom.NewUniformInsets(-2 * width))
		for a := p.Parent(); a != nil; a = a.Parent() {
			clip = clip.Intersect(content.RectFromRoot(a.RectToRoot(a.ContentRect(false))).Inset(geom.NewUniformInsets(-width)))
		}
		if clip.Empty() {
			return
		}
		gc.Save()
		gc.ClipRect(clip, pathop.Intersect, false)
		paint := Color(u.Theme.Tokens().FocusRing).Paint(gc, ring, paintstyle.Stroke)
		paint.SetStrokeWidth(width)
		gc.DrawRoundedRect(ring, geom.NewSize(radius+width/2, radius+width/2), paint)
		gc.Restore()
	}
}

// joinLines joins the non-empty parts, one to a line.
func joinLines(parts ...string) string {
	out := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		if out != "" {
			out += "\n"
		}
		out += p
	}
	return out
}

// newTooltip is the Kvit tooltip: the label, and the explanation on a line
// of its own after it, in the small role on the popup ground. unison shows a
// panel's tooltip on pointer hover after a delay.
func newTooltip(ui *UI, label, explanation string) *unison.Panel {
	t := ui.Theme.Tokens()
	m := ui.Interface
	body := label
	if explanation != "" {
		if body != "" {
			body += "\n"
		}
		body += explanation
	}
	l := ui.Fonts.Layout([]text.Span{{Text: body, Style: ui.Chrome(ui.Size(RoleSmall), text.Regular, t.TextPrimary)}},
		text.Options{MaxWidth: float32(m.Px(280))})
	pad := float32(m.SpaceNear())
	w, h := l.Size()
	p := unison.NewPanel()
	// Named by its words, as the tooltip it is, for a screen reader that
	// reads what appears.
	p.Accessibility.Role, p.Accessibility.Name = role.Tooltip, body
	size := geom.NewSize(w+2*pad, h+2*pad)
	p.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) { return size, size, size })
	p.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		pt := painterFor(gc, ui)
		b := p.ContentRect(false)
		pt.round(b, float32(m.RadiusControl()), t.PopupBackground)
		pt.outline(b, float32(m.RadiusControl()), float32(m.Hairline()), t.BorderStrong)
		l.Draw(gc, pad, pad)
	}
	return p
}
