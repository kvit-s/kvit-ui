package kvitui

import (
	"time"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// partTip is a surface shown for one part of a panel at a time: the tooltip
// of a mark in a row of marks, of a cell of a table cut short, or the hover
// card behind an amount. A control has one tooltip for the whole control; a
// table has one panel for hundreds of cells, so the part the pointer or the
// keyboard is on is named by a key, and the surface changes when the key does.
type partTip struct {
	ui    *UI
	owner unison.Paneler
	key   any
	gen   int    // counts changes, so a late show or hide can tell it is stale
	hide  func() // takes the surface away while one is shown
}

// tooltip shows words beside a part after the pointer rests, as a control's
// tooltip shows, for long enough to read: at least three seconds, and longer
// for more words. box is the part's box in the owner's coordinates, asked on
// every layout so the tooltip follows the part as the owner scrolls. The same
// key again changes nothing; a nil key or no words takes the tooltip away.
func (t *partTip) tooltip(key any, words string, box func() geom.Rect) {
	if key != nil && key == t.key {
		return
	}
	t.clear()
	if key == nil || words == "" {
		return
	}
	t.key = key
	gen := t.gen
	unison.InvokeTaskAfter(func() {
		if gen != t.gen {
			return
		}
		t.show(newTooltip(t.ui, words, ""), placeBeside(t.ui, func() geom.Rect { return partIn(t.owner, box()) }))
		unison.InvokeTaskAfter(func() {
			if gen == t.gen {
				t.clear()
			}
		}, time.Duration(max(3000, len([]rune(words))*60))*time.Millisecond)
	}, 500*time.Millisecond)
}

// card shows a panel under a part at once, its right edge at the part's
// right edge, until the key changes: the second view of a value that a hover
// card is.
func (t *partTip) card(key any, panel func() unison.Paneler, box func() geom.Rect) {
	if key != nil && key == t.key {
		return
	}
	t.clear()
	if key == nil || panel == nil {
		return
	}
	t.key = key
	t.show(panel(), func(_ geom.Rect, size geom.Size) geom.Rect {
		a := partIn(t.owner, box())
		return geom.NewRect(a.Right()-size.Width, a.Bottom(), size.Width, size.Height)
	})
}

func (t *partTip) show(p unison.Paneler, place func(bounds geom.Rect, size geom.Size) geom.Rect) {
	if w := t.ui.windowOf(t.owner); w != nil {
		t.hide = w.Show(&Popup{Panel: p, Place: place, Anchor: t.owner})
	}
}

// clear takes the surface away, and cancels one about to show.
func (t *partTip) clear() {
	t.gen++
	t.key = nil
	if t.hide != nil {
		t.hide()
		t.hide = nil
	}
}
