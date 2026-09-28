package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// ShowTooltip shows a tooltip beside an anchor until the returned function is
// called, placed as a control's own tooltip is. A tooltip may name something
// and may not be the only place its meaning lives: a control explained only
// by its tooltip is unusable on a keyboard, a screen reader or a touch screen.
// It does nothing outside a Kvit window.
func (u *UI) ShowTooltip(anchor unison.Paneler, words string) (hide func()) {
	var hidden bool
	var remove func()
	show := func() {
		if w := u.windowOf(anchor); w != nil && words != "" && !hidden {
			remove = w.Show(&Popup{Panel: newTooltip(u, words, ""), Place: PlaceBeside(u, anchor), Anchor: anchor,
				Passive: true})
		}
	}
	// An anchor not in a window yet shows it once it is.
	if anchor.AsPanel().Window() == nil {
		unison.InvokeTask(show)
	} else {
		show()
	}
	return func() {
		hidden = true
		if remove != nil {
			remove()
		}
	}
}

// Popover is a small surface anchored to a control, holding something the
// reader can act on: a picker, a filter, a short form. It takes the keyboard
// focus and closes on Escape or on a press outside it, which is what sets it
// apart from a HoverCard: a surface a reader can type into has to be
// reachable and dismissible by keyboard, and one that only appears on hover
// can be neither. Closing gives the focus back to where it was before, unless
// the reader has put it somewhere else since.
type Popover struct {
	unison.Panel
	ui *UI
	// Title names the popover for a screen reader.
	Title string
	// Width holds the popover to a width, its height following from it; nil
	// takes the width its content asks for.
	Width Measure
	// OnClose runs after the popover closes.
	OnClose func()

	hide     func()
	previous *unison.Panel // what held the focus when it opened
}

// NewPopover returns a closed popover holding content, stacked top to bottom
// a near space apart.
func NewPopover(ui *UI, title string, content ...unison.Paneler) *Popover {
	p := &Popover{ui: ui, Title: title}
	p.Self = p
	col := &spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: SizeSpaceNear}
	for _, c := range content {
		if c.AsPanel().LayoutData() == nil {
			c.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		}
		p.AddChild(c)
	}
	p.SetBorder(Padding(ui, SizeSpaceLoose))
	p.SetLayout(col)
	p.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) { drawSurface(gc, ui, p.ContentRect(true)) }
	return p
}

// drawSurface is the ground a floating surface stands on: the popup colour, a
// strong hairline edge, a card's corners.
func drawSurface(gc *unison.Canvas, ui *UI, r geom.Rect) {
	t, m := ui.Theme.Tokens(), ui.Interface
	p := painterFor(gc, ui)
	radius := float32(m.RadiusCard())
	p.round(r, radius, t.PopupBackground)
	p.outline(r, radius, float32(m.Hairline()), t.BorderStrong)
}

// Open shows the popover in the anchor's window, where place puts it, and
// gives it the keyboard focus. It does nothing outside a Kvit window.
func (p *Popover) Open(anchor unison.Paneler, place func(bounds geom.Rect, size geom.Size) geom.Rect) {
	if anchor.AsPanel().Window() == nil {
		// Not in a window yet: open once it is.
		unison.InvokeTask(func() {
			if anchor.AsPanel().Window() != nil {
				p.Open(anchor, place)
			}
		})
		return
	}
	w := p.ui.windowOf(anchor)
	if w == nil || p.hide != nil {
		return
	}
	p.previous = w.CurrentFocus()
	sized := func(bounds geom.Rect, size geom.Size) geom.Rect {
		if p.Width != nil {
			width := orZero(p.ui, p.Width)
			_, pref, _ := p.Sizes(geom.NewSize(width, 0))
			size = geom.NewSize(width, pref.Height)
		}
		return place(bounds, size)
	}
	p.hide = w.Show(&Popup{Panel: p, Place: sized, OnEscape: p.Close, OnPressOutside: p.Close, Anchor: anchor})
	unison.InvokeTask(func() {
		if f := p.FirstFocusableChild(); f != nil {
			f.RequestFocus()
		}
	})
}

// Opened reports whether the popover is showing.
func (p *Popover) Opened() bool { return p.hide != nil }

// Close takes the popover away and gives the focus back.
func (p *Popover) Close() {
	if p.hide == nil {
		return
	}
	w := p.Window()
	inside := false
	if w != nil {
		for f := w.CurrentFocus(); f != nil; f = f.Parent() {
			if f == p.AsPanel() {
				inside = true
				break
			}
		}
	}
	p.hide()
	p.hide = nil
	if inside && p.previous != nil && p.previous.Window() != nil {
		p.previous.RequestFocus()
	}
	if p.OnClose != nil {
		p.OnClose()
	}
}

// ProvideAccessibility names the popover as a dialog.
func (p *Popover) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Dialog
	}
	n.Name = p.Title
}

// Hint is an information symbol that opens an explanation too long for a
// tooltip. A press, or Space or Return on it, opens a popover holding the
// label and the explanation, which stays while it is read; Escape or a press
// outside closes it. There is no tooltip first: the label and the explanation
// are one piece of information, reached by one action.
type Hint struct {
	unison.Panel
	ui *UI
	// Label names what is explained; it is the button's name and the
	// popover's title.
	Label string
	// Text is the explanation.
	Text    string
	trigger *IconButton
	pop     *Popover
	title   *Label
	detail  *Label
}

// NewHint returns an information symbol explaining something.
func NewHint(ui *UI, label, explanation string) *Hint {
	h := &Hint{ui: ui, Label: label, Text: explanation}
	h.Self = h
	h.trigger = NewIconButton(ui, "info", label)
	h.trigger.TooltipEnabled = false
	h.trigger.OnClick = func() {
		if h.pop.Opened() {
			h.pop.Close()
		} else {
			h.Open()
		}
	}
	h.title = NewLabel(ui, label)
	h.title.Weight = text.Bold
	h.detail = NewLabel(ui, explanation)
	h.detail.Role, h.detail.Ink, h.detail.Wrap, h.detail.LineHeight = RoleSmall, InkTextSecondary, true, 1.3
	h.pop = NewPopover(ui, label, h.title, h.detail)
	h.pop.Width = Px(280)
	h.AddChild(h.trigger)
	h.SetLayout(&unison.FlexLayout{Columns: 1})
	return h
}

// Open shows the explanation under the symbol, or ending at its right edge
// where it would run off the window.
func (h *Hint) Open() {
	h.title.Text, h.detail.Text, h.pop.Title = h.Label, h.Text, h.Label
	h.pop.Open(h, PlaceBelow(h.ui, h))
}

// Opened reports whether the explanation is showing.
func (h *Hint) Opened() bool { return h.pop.Opened() }

// Close takes the explanation away.
func (h *Hint) Close() { h.pop.Close() }
