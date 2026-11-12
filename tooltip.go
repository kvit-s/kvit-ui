package kvitui

import (
	"time"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/mod"
)

// Tooltip is a tooltip on any panel, given with AddTooltip: shown beside the
// panel half a second after the pointer comes to rest on it or the keyboard
// reaches it, placed and drawn as a control's own tooltip is. With every
// field at its zero value it behaves as a control's own does: the keyboard
// focus is watched, and the tooltip goes after long enough to read it and
// when the panel is pressed.
//
// Its words are read when it shows and not again while it is up. A tooltip
// may name something and may not be the only place its meaning lives.
type Tooltip struct {
	// Focus is how the tooltip learns that the keyboard is on the panel.
	Focus TipFocus
	// Keyboard, with TipFocusWatched, answers whether the focus just gained
	// came from the keyboard. Nil uses the rule Kvit's controls use: a key
	// moved the focus or Focus() asked for it, and a press did not.
	Keyboard func() bool
	// FocusAtOnce shows the tooltip for the keyboard without the half-second
	// pause.
	FocusAtOnce bool
	// Stays keeps the tooltip shown for as long as the pointer rests or the
	// keyboard stays. Without it the tooltip also goes after the time it
	// takes to read, at least three seconds, and when a button goes down on
	// the panel.
	Stays bool
	// HoverTold leaves it to the caller to say, with SetHovered, whether the
	// pointer is on the panel, and the panel's own pointer callbacks are not
	// used for it: for a panel that counts the pointer on anything inside it
	// as on itself, as a ListRow with HoverWhole does, from its OnHover.
	HoverTold bool

	ui               *UI
	owner            *unison.Panel
	words            func() string
	hovered, focused bool
	gen              int
	hide             func()
}

// TipFocus is how a Tooltip learns that the keyboard is on its panel.
type TipFocus int

const (
	// TipFocusWatched wraps the panel's focus callbacks.
	TipFocusWatched TipFocus = iota
	// TipFocusTold leaves it to the caller, who says with SetFocused.
	TipFocusTold
	// TipFocusNone never shows the tooltip for the keyboard.
	TipFocusNone
)

// tipSuppressor is a Kvit control, whose own tooltip stays away while one is
// added to it, so there are never two.
type tipSuppressor interface{ suppressOwnTip() }

// AddTooltip gives a panel a tooltip whose words are asked each time it is
// about to show; "" shows none. It wraps the panel's pointer and focus
// callbacks and calls the ones it replaced first, so set the panel's own
// before calling it. On a Kvit control the control's own tooltip no longer
// shows. It shows nothing outside a Kvit window.
func AddTooltip(ui *UI, owner unison.Paneler, words func() string) *Tooltip {
	p := owner.AsPanel()
	t := &Tooltip{ui: ui, owner: p, words: words}
	if s, ok := p.Self.(tipSuppressor); ok {
		s.suppressOwnTip()
	}
	enter, exit, down := p.MouseEnterCallback, p.MouseExitCallback, p.MouseDownCallback
	gained, lost := p.GainedFocusCallback, p.LostFocusCallback
	p.MouseEnterCallback = func(where geom.Point, mods mod.Modifiers) bool {
		if enter != nil {
			enter(where, mods)
		}
		if !t.HoverTold {
			t.hovered = true
			t.update(true)
		}
		return false
	}
	p.MouseExitCallback = func() bool {
		if exit != nil {
			exit()
		}
		if !t.HoverTold {
			t.hovered = false
			t.update(false)
		}
		return false
	}
	p.MouseDownCallback = func(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
		if !t.Stays {
			// Gone until the pointer comes back or the keyboard arrives: the
			// focus the press gives must not bring it straight back.
			t.hovered = false
			t.Hide()
		}
		return down != nil && down(where, button, clicks, mods)
	}
	p.GainedFocusCallback = func() {
		if gained != nil {
			gained()
		}
		if t.Focus != TipFocusWatched {
			return
		}
		t.focused = t.keyboard()
		t.update(t.focused && !t.FocusAtOnce)
	}
	p.LostFocusCallback = func() {
		if lost != nil {
			lost()
		}
		if t.Focus != TipFocusWatched {
			return
		}
		t.focused = false
		t.update(false)
	}
	return t
}

// keyboard reports whether the focus the panel just gained came from the
// keyboard.
func (t *Tooltip) keyboard() bool {
	if t.Keyboard != nil {
		return t.Keyboard()
	}
	if c, ok := t.owner.Self.(interface{ KeyboardFocus() bool }); ok {
		return c.KeyboardFocus()
	}
	return t.ui.keyTurn
}

// update shows or hides the tooltip for the state as it is now, after the
// half-second pause when wait is set. Every change of state cancels a show
// that is waiting.
func (t *Tooltip) update(wait bool) {
	t.gen++
	if !t.hovered && !t.focused {
		if t.hide != nil {
			t.hide()
			t.hide = nil
		}
		return
	}
	if t.hide != nil {
		return
	}
	gen := t.gen
	show := func() {
		if gen != t.gen || t.hide != nil || (!t.hovered && !t.focused) {
			return
		}
		say := t.words()
		if say == "" {
			return
		}
		t.hide = t.ui.ShowTooltip(t.owner, say)
		if t.Stays {
			return
		}
		// For long enough to read: at least three seconds, and longer for
		// more words.
		unison.InvokeTaskAfter(func() {
			if gen == t.gen && !t.Stays {
				t.Hide()
			}
		}, time.Duration(max(3000, len([]rune(say))*60))*time.Millisecond)
	}
	if wait {
		unison.InvokeTaskAfter(show, 500*time.Millisecond)
	} else {
		show()
	}
}

// Shown reports whether the tooltip is up.
func (t *Tooltip) Shown() bool { return t.hide != nil }

// Words are what the tooltip says now, "" for none.
func (t *Tooltip) Words() string { return t.words() }

// SetFocused says whether the keyboard is on the panel, for a tooltip whose
// Focus is TipFocusTold: typically called from the panel's own focus
// callbacks with whether its focus ring shows.
func (t *Tooltip) SetFocused(on bool) {
	t.focused = on
	t.update(on && !t.FocusAtOnce)
}

// SetHovered says whether the pointer is on the panel, for a tooltip with
// HoverTold: it shows half a second after the pointer arrived.
func (t *Tooltip) SetHovered(on bool) {
	t.hovered = on
	t.update(on)
}

// Hide takes the tooltip away, and cancels one about to show. It shows again
// the next time the pointer arrives or the keyboard does.
func (t *Tooltip) Hide() {
	t.gen++
	if t.hide != nil {
		t.hide()
		t.hide = nil
	}
}
