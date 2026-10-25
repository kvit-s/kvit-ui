package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// Tip teaches one feature that is not obvious, or is new, beside what it
// explains (HLP-4, decided Q-48): a card with a title, a sentence or two, a
// close button and at most one action. It covers and blocks nothing: the
// window goes on working around it, Escape and Tab pass it by, and a press
// elsewhere leaves it. Inside a scrolling region it closes as soon as what
// it points at moves. Once closed it is not shown again: the settings file
// keeps its ID, and ResetTips brings every tip back. Spotlight is for a
// guided tour.
type Tip struct {
	unison.Panel
	ui *UI
	// ID names the tip in the settings file, which keeps whether it was
	// seen.
	ID string
	// Title says what the feature is; Text says what it does for the reader.
	Title, Text string
	// Action is the words of the tip's one action, "" for none.
	Action string
	// OnAction runs when the action is taken, after the tip closes;
	// OnClose when the tip closes by its close button or its target moving.
	OnAction, OnClose func()

	title, words *Label
	action       *Button
	close        *IconButton
	hide         func()
	anchor       unison.Paneler
	anchoredAt   geom.Rect // where the target was when the tip was shown
}

// tipsSeenKey is the settings key under which the IDs of the tips seen are
// kept.
const tipsSeenKey = "tips.seen"

// NewTip returns a tip, not shown.
func NewTip(ui *UI, id, title, words string) *Tip {
	t := &Tip{ui: ui, ID: id, Title: title, Text: words}
	t.Self = t
	t.title = NewLabel(ui, "")
	t.title.Weight = text.Bold
	t.words = NewLabel(ui, "")
	t.words.Wrap, t.words.Ink = true, InkTextSecondary
	t.action = NewButton(ui, "")
	t.action.OnClick = func() {
		t.dismiss()
		if t.OnAction != nil {
			t.OnAction()
		}
	}
	t.close = NewIconButton(ui, "close", "Close")
	t.close.OnClick = func() {
		t.dismiss()
		if t.OnClose != nil {
			t.OnClose()
		}
	}
	t.SetBorder(Padding(ui, SizeSpaceLoose))
	t.SetLayout(syncing{Layout: tipLayout{t}, sync: func() {
		t.title.Text, t.words.Text, t.action.Text = t.Title, t.Text, t.Action
		parts := []unison.Paneler{t.title, t.close, t.words}
		if t.Action != "" {
			parts = append(parts, t.action)
		}
		showOnly(t.AsPanel(), parts...)
	}})
	t.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		drawSurface(gc, ui, t.ContentRect(true))
		t.watchAnchor()
	}
	return t
}

// Seen reports whether the reader has closed the tip before.
func (t *Tip) Seen() bool { return t.ui.tipSeen(t.ID) }

// ShowBeside shows the tip pointing at what it explains: above it where it
// fits, below it otherwise. A tip already seen is not shown, and false says
// so.
func (t *Tip) ShowBeside(anchor unison.Paneler) bool {
	w := t.ui.windowOf(anchor)
	if w == nil || t.Seen() || t.hide != nil {
		return false
	}
	t.anchor = anchor
	t.anchoredAt = anchorIn(anchor)
	t.hide = w.Show(&Popup{Panel: t, Anchor: anchor, Ambient: true, Place: func(bounds geom.Rect, size geom.Size) geom.Rect {
		a := anchorIn(anchor)
		width := min(float32(t.ui.Interface.Px(300)), bounds.Width)
		_, pref, _ := t.Sizes(geom.NewSize(width, 0))
		gap := float32(t.ui.Interface.Space())
		y := a.Y - gap - pref.Height
		if y < bounds.Y {
			y = a.Bottom() + gap
		}
		x := max(bounds.X, min(a.X, bounds.Right()-width))
		return geom.NewRect(x, y, width, pref.Height)
	}})
	return true
}

// ShowInWindow shows a tip about nothing in particular, centred at the
// bottom of the window. A tip already seen is not shown, and false says so.
func (t *Tip) ShowInWindow(w *Window) bool {
	if w == nil || t.Seen() || t.hide != nil {
		return false
	}
	t.anchor = nil
	t.hide = w.Show(&Popup{Panel: t, Ambient: true, Place: func(bounds geom.Rect, size geom.Size) geom.Rect {
		width := min(float32(t.ui.Interface.Px(360)), bounds.Width)
		_, pref, _ := t.Sizes(geom.NewSize(width, 0))
		margin := float32(t.ui.Interface.ViewMargin())
		return geom.NewRect(bounds.X+(bounds.Width-width)/2, bounds.Bottom()-margin-pref.Height, width, pref.Height)
	}})
	return true
}

// Shown reports whether the tip is on the screen.
func (t *Tip) Shown() bool { return t.hide != nil }

// dismiss takes the tip away and records that it was seen.
func (t *Tip) dismiss() {
	if t.hide != nil {
		t.hide()
		t.hide = nil
	}
	t.ui.markTipSeen(t.ID)
}

// watchAnchor closes the tip once what it points at has moved, as when the
// region holding it scrolls.
func (t *Tip) watchAnchor() {
	if t.anchor == nil || t.hide == nil {
		return
	}
	if a := anchorIn(t.anchor); a.Point != t.anchoredAt.Point {
		unison.InvokeTask(func() {
			if t.hide != nil {
				t.dismiss()
				if t.OnClose != nil {
					t.OnClose()
				}
			}
		})
	}
}

type tipLayout struct{ t *Tip }

// arrange places the title with the close button at its end, the words
// under it, and the action under them, and returns the height.
func (l tipLayout) arrange(r geom.Rect, set bool) float32 {
	t, m := l.t, l.t.ui.Interface
	_, cp, _ := t.close.Sizes(geom.Size{})
	gap := float32(m.SpaceNear())
	_, tp, _ := t.title.Sizes(geom.NewSize(r.Width-cp.Width-gap, 0))
	row := max(tp.Height, cp.Height)
	_, wp, _ := t.words.Sizes(geom.NewSize(r.Width, 0))
	y := r.Y + row + gap
	h := row + gap + wp.Height
	if set {
		t.title.SetFrameRect(geom.NewRect(r.X, r.Y+(row-tp.Height)/2, r.Width-cp.Width-gap, tp.Height))
		t.close.SetFrameRect(geom.NewRect(r.Right()-cp.Width, r.Y+(row-cp.Height)/2, cp.Width, cp.Height))
		t.words.SetFrameRect(geom.NewRect(r.X, y, r.Width, wp.Height))
	}
	if t.Action != "" {
		_, ap, _ := t.action.Sizes(geom.Size{})
		if set {
			t.action.SetFrameRect(geom.NewRect(r.X, y+wp.Height+gap, ap.Width, ap.Height))
		}
		h += gap + ap.Height
	}
	return h
}

func (l tipLayout) LayoutSizes(target *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.t.ui.Interface
	var in geom.Insets
	if b := target.Border(); b != nil {
		in = b.Insets()
	}
	w := float32(m.Px(300))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := l.arrange(geom.NewRect(0, 0, w-in.Width(), 0), false) + in.Height()
	return geom.NewSize(w, h), geom.NewSize(w, h), geom.NewSize(w, h)
}

func (l tipLayout) PerformLayout(target *unison.Panel) { l.arrange(target.ContentRect(false), true) }

// ProvideAccessibility reads the tip as its title and words.
func (t *Tip) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name, n.Description = t.Title, t.Text
}

// tipSeen reports whether a tip was closed before.
func (u *UI) tipSeen(id string) bool {
	if u.Settings != nil {
		if v, ok := u.Settings.Value(tipsSeenKey); ok {
			if ids, ok := v.([]any); ok {
				for _, x := range ids {
					if x == id {
						return true
					}
				}
			}
		}
		return false
	}
	return u.tipsSeen[id]
}

// markTipSeen records that a tip was closed.
func (u *UI) markTipSeen(id string) {
	if u.tipSeen(id) {
		return
	}
	if u.Settings != nil {
		var ids []any
		if v, ok := u.Settings.Value(tipsSeenKey); ok {
			ids, _ = v.([]any)
		}
		u.Settings.SetValue(tipsSeenKey, append(ids, id))
		return
	}
	if u.tipsSeen == nil {
		u.tipsSeen = map[string]bool{}
	}
	u.tipsSeen[id] = true
}

// ResetTips brings every tip back, as a setting a reader can use does.
func (u *UI) ResetTips() {
	u.tipsSeen = nil
	if u.Settings != nil {
		u.Settings.Remove(tipsSeenKey)
	}
}
