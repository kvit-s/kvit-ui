package kvitui

import (
	"slices"

	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/pathop"
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
// noise; a keyboard user needs to see where they are.
type control struct {
	unison.Panel
	ui            *UI
	hovered       bool
	pressed       bool
	keyboardFocus bool // the focus came from the keyboard, so the ring shows
	pointerFocus  bool // set while a press is giving the control the focus
	activate      func()
	keys          []unison.KeyCode // the keys that activate the control
	ringRadius    func() float32   // the corner radius of the shape the ring surrounds
	noRing        bool             // shows keyboard focus its own way instead of a ring
}

// initControl wires the callbacks. activate runs on a click, on one of keys,
// and on a screen reader's press.
func (c *control) initControl(ui *UI, activate func(), keys ...unison.KeyCode) {
	c.ui = ui
	c.activate = activate
	c.keys = keys
	c.SetFocusable(true)
	c.MouseEnterCallback = func(geom.Point, mod.Modifiers) bool {
		c.hovered = true
		c.MarkForRedraw()
		return false
	}
	c.MouseExitCallback = func() bool {
		c.hovered = false
		c.MarkForRedraw()
		return false
	}
	c.MouseDownCallback = func(_ geom.Point, button, _ int, _ mod.Modifiers) bool {
		if !c.Enabled() || button != unison.ButtonLeft {
			return false
		}
		c.pressed = true
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
		c.keyboardFocus = !c.pointerFocus
		c.pointerFocus = false
		c.MarkForRedraw()
		if w := c.Window(); w != nil {
			ui.paintRingsIn(w)
			w.MarkForRedraw()
		}
	}
	c.LostFocusCallback = func() {
		c.keyboardFocus = false
		c.MarkForRedraw()
		if w := c.Window(); w != nil {
			w.MarkForRedraw()
		}
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
		c.RequestFocus()
		return
	}
	unison.InvokeTask(func() {
		if c.Window() != nil {
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

// paintRingsIn makes a window draw the focus ring of whichever Kvit control
// holds its keyboard focus. The ring sits outside the control's own box, as
// the Qt library draws it, and unison clips each panel's drawing to its
// bounds, so it is drawn by the window over everything, clipped to the area
// the control's container shows.
func (u *UI) paintRingsIn(w *unison.Window) {
	if u.ringWindows == nil {
		u.ringWindows = map[*unison.Window]bool{}
	}
	if u.ringWindows[w] {
		return
	}
	u.ringWindows[w] = true
	content := w.Content()
	previous := content.DrawOverCallback
	content.DrawOverCallback = func(gc *unison.Canvas, rect geom.Rect) {
		if previous != nil {
			previous(gc, rect)
		}
		focus := w.Focus()
		if focus == nil {
			return
		}
		r, ok := focus.Self.(ringed)
		if !ok {
			return
		}
		p, radius, show := r.focusRing()
		if !show {
			return
		}
		m := u.Interface
		width := float32(m.FocusRingWidth())
		box := content.RectFromRoot(p.RectToRoot(p.ContentRect(false)))
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
