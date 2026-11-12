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
	// OnEscape runs when Escape is pressed while the popover is open, in
	// place of closing it; nil closes it.
	OnEscape func()
	// OnPressOutside runs when the pointer is pressed outside the popover,
	// in place of closing it; nil closes it.
	OnPressOutside func()
	// PressOutsideIgnored makes a press outside the popover do nothing;
	// OnPressOutside is not used.
	PressOutsideIgnored bool
	// ModalPress is what becomes of a press outside a popover opened with
	// OpenModal, once it has closed the popover or run OnPressOutside: it
	// goes no further (the default), or it goes on to what is behind.
	ModalPress ModalPress
	// MayClose, when set, is asked by Close, and by Escape and a press
	// outside where they close the popover: false leaves it open, as while it
	// shows work that cannot be stopped.
	MayClose func() bool
	// FirstFocus answers what takes the keyboard when the popover opens; nil,
	// or answering nil, gives it to the first control inside.
	FirstFocus func() *unison.Panel

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

// drawSurface is DrawSurface with a card's corners.
func drawSurface(gc *unison.Canvas, ui *UI, r geom.Rect) { DrawSurface(gc, ui, r, nil) }

// DrawSurface draws the ground a floating surface stands on: the popup
// colour with a strong hairline edge, its corners rounded by radius. A nil
// radius is a card's corners, which is what a Popover draws.
func DrawSurface(gc *unison.Canvas, ui *UI, r geom.Rect, radius Measure) {
	t, m := ui.Theme.Tokens(), ui.Interface
	p := painterFor(gc, ui)
	round := float32(m.RadiusCard())
	if radius != nil {
		round = float32(radius.Of(ui))
	}
	p.round(r, round, t.PopupBackground)
	p.outline(r, round, float32(m.Hairline()), t.BorderStrong)
}

// Open shows the popover in the anchor's window, where place puts it, and
// gives it the keyboard focus. It does nothing outside a Kvit window.
func (p *Popover) Open(anchor unison.Paneler, place func(bounds geom.Rect, size geom.Size) geom.Rect) {
	p.open(anchor, place, false)
}

// ModalPress is what becomes of a press outside a modal popover.
type ModalPress int

const (
	// ModalPressStops ends the press at the popover: it closes the popover,
	// or runs OnPressOutside, and presses nothing behind.
	ModalPressStops ModalPress = iota
	// ModalPressPasses closes the popover, or runs OnPressOutside, at once,
	// and where that closed the popover the same press goes on to what is
	// behind it. While the popover stays open the scrim takes the press.
	ModalPressPasses
)

// OpenModal shows the popover as Open does, over a scrim that takes every
// press, so nothing behind it can be used until it closes. ModalPress says
// whether the press that closes it from outside reaches what is behind.
func (p *Popover) OpenModal(anchor unison.Paneler, place func(bounds geom.Rect, size geom.Size) geom.Rect) {
	p.open(anchor, place, true)
}

func (p *Popover) open(anchor unison.Paneler, place func(bounds geom.Rect, size geom.Size) geom.Rect, modal bool) {
	if anchor.AsPanel().Window() == nil {
		// Not in a window yet: open once it is.
		unison.InvokeTask(func() {
			if anchor.AsPanel().Window() != nil {
				p.open(anchor, place, modal)
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
	popup := &Popup{Panel: p, Place: sized, Anchor: anchor, HoldsKeyboard: true, Modal: modal}
	popup.OnEscape = func() {
		if p.OnEscape != nil {
			p.OnEscape()
			return
		}
		p.Close()
	}
	// A press outside that is ignored is left to whatever it lands on, where
	// one that closes the popover goes no further.
	if !p.PressOutsideIgnored {
		outside := func() {
			if p.OnPressOutside != nil {
				p.OnPressOutside()
				return
			}
			p.Close()
		}
		popup.OnPressOutside = outside
		if modal && p.ModalPress == ModalPressStops {
			// Once the press is over: closing takes the scrim away, and a
			// press still under way would reach what is behind it.
			popup.OnPressOutside = func() {
				unison.InvokeTask(func() {
					if p.hide != nil {
						outside()
					}
				})
			}
		}
	}
	p.hide = w.Show(popup)
	unison.InvokeTask(func() {
		var f *unison.Panel
		if p.FirstFocus != nil {
			f = p.FirstFocus()
		}
		if f == nil {
			f = p.FirstFocusableChild()
		}
		if f != nil {
			f.RequestFocus()
		}
	})
}

// Opened reports whether the popover is showing.
func (p *Popover) Opened() bool { return p.hide != nil }

// Close takes the popover away and gives the focus back.
func (p *Popover) Close() {
	if p.hide == nil || (p.MayClose != nil && !p.MayClose()) {
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
	OpensPopup(h.trigger)
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
