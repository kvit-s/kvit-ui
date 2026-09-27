package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/role"
)

// Link is an inline destination or action, announced to screen readers as a
// link, never as a button. Hover and keyboard focus both turn it the accent
// colour and underline it; it grows no chevron, ground or border of its own.
// Space, Return and Enter follow it, as in a browser.
type Link struct {
	control
	// Text is what the link says.
	Text string
	// Symbol is an optional meaning name drawn before the text.
	Symbol string
	// Role is the text's type role; RoleBody unless set.
	Role TypeRole
	// Explanation says in a sentence where following the link goes, for a
	// destination its words alone do not describe. It is the tooltip and the
	// accessible description.
	Explanation string
	// OnActivate runs when the link is followed.
	OnActivate func()
}

// NewLink returns a link in the body role.
func NewLink(ui *UI, s string) *Link {
	l := &Link{Text: s, Role: RoleBody}
	l.Self = l
	l.initControl(ui, func() {
		if l.OnActivate != nil {
			l.OnActivate()
		}
	}, unison.KeySpace, unison.KeyReturn, unison.KeyNumPadEnter)
	// Keyboard focus shows as the accent colour and the underline, as hover
	// does; a link draws no ring.
	l.noRing = true
	l.SetSizer(l.sizes)
	l.DrawCallback = l.draw
	l.UpdateCursorCallback = func(geom.Point) *unison.Cursor { return unison.PointingCursor() }
	l.UpdateTooltipCallback = func(geom.Point, geom.Rect) geom.Rect {
		if l.Explanation != "" {
			l.Tooltip = newTooltip(ui, "", l.Explanation)
		} else {
			l.Tooltip = nil
		}
		return l.RectToRoot(l.ContentRect(false))
	}
	return l
}

func (l *Link) active() bool { return l.hovered || l.KeyboardFocus() }

func (l *Link) symbolWidth() float32 {
	if l.Symbol == "" {
		return 0
	}
	m := l.ui.Interface
	return float32(m.IconSizeSmall() + m.SpaceNear())
}

func (l *Link) layout(width float32) *text.Layout {
	t := l.ui.Theme.Tokens()
	ink := t.Link
	if l.active() {
		ink = t.Accent
	}
	st := l.ui.Chrome(l.ui.Size(l.Role), text.Regular, ink)
	st.Underline = l.active()
	return l.ui.Fonts.Layout([]text.Span{{Text: l.Text, Style: st}}, text.Options{MaxWidth: width, Elide: width > 0})
}

func (l *Link) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	w, h := l.layout(0).Size()
	h = max(h, float32(l.ui.Interface.IconSizeSmall()))
	pref := geom.NewSize(l.symbolWidth()+w, h)
	// It can be made narrower, and is then cut short with "…"; given more
	// room than it needs, it keeps to its words on the left.
	return geom.NewSize(l.symbolWidth(), h), pref, geom.NewSize(unison.DefaultMaxSize, h)
}

func (l *Link) draw(gc *unison.Canvas, _ geom.Rect) {
	t := l.ui.Theme.Tokens()
	b := l.ContentRect(false)
	x := float32(0)
	if l.Symbol != "" {
		ink := t.Link
		if l.active() {
			ink = t.Accent
		}
		if g, ok := icons.Glyph(l.Symbol); ok {
			s := float32(l.ui.Interface.IconSizeSmall())
			drawGlyph(gc, l.ui, g, s, ink, geom.NewRect(0, (b.Height-s)/2, s, s))
		}
		x = l.symbolWidth()
	}
	lay := l.layout(b.Width - x)
	_, h := lay.Size()
	lay.Draw(gc, x, (b.Height-h)/2)
}

// ProvideAccessibility describes the link by its words, with the explanation
// as its description.
func (l *Link) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Link
	}
	n.Name = l.Text
	if l.Explanation != "" {
		n.Description = l.Explanation
	}
	if l.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction follows the link for a screen reader.
func (l *Link) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !l.Enabled() {
		return false
	}
	l.fire()
	return true
}
