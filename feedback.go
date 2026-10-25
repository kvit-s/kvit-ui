package kvitui

import (
	"time"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// HoverCard shows more about the thing under the pointer, without a press:
// the full name behind a cut-short one, the exact figure behind a rounded
// one, the date behind "3 days ago". It is read-only and must stay so:
// nothing inside one may be the only route to an action, because a surface
// that appears on hover cannot be reached by keyboard or by touch. A screen
// reader passes over it, since everything in it is a second view of what is
// already on the screen. Show it with Window.Show, or place it yourself.
type HoverCard struct {
	unison.Panel
	ui *UI
}

// NewHoverCard returns a card holding content, stacked top to bottom a tight
// space apart.
func NewHoverCard(ui *UI, content ...unison.Paneler) *HoverCard {
	c := &HoverCard{ui: ui}
	c.Self = c
	for _, p := range content {
		c.AddChild(p)
	}
	c.SetBorder(Padding(ui, SizeSpaceLoose))
	c.SetLayout(&spaced{FlexLayout: unison.FlexLayout{Columns: 1}, ui: ui, gap: SizeSpaceTight})
	c.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) { drawSurface(gc, ui, c.ContentRect(true)) }
	return c
}

// ProvideAccessibility leaves the card and all it holds out of what a screen
// reader is told.
func (c *HoverCard) ProvideAccessibility(b *unison.AccessibilityBuilder) { b.Node().Ignored = true }

// toneSymbol is the symbol a message of a tone carries.
func toneSymbol(t Tone) string {
	switch t {
	case ToneSuccess:
		return "success"
	case ToneWarning:
		return "warning"
	case ToneDanger:
		return "error"
	}
	return "info"
}

// messageTone is a message's tone colour: the accent stands in for info.
func messageTone(t Tone, tk tokens.Tokens) palette.Color {
	if t == ToneNeutral || t == ToneInfo {
		return tk.Accent
	}
	return t.color(tk)
}

// Toast is a message that appears, says one thing and goes away: a
// confirmation of something that already happened, saved, copied, three rows
// archived, and the undo that goes with it. A toast is never the only route
// to a decision, which would disappear after four seconds. One with an action
// stays until it is dismissed: an undo that times out mid-read is one the
// reader cannot use. It is announced to a screen reader when shown.
type Toast struct {
	unison.Panel
	ui *UI
	// Text is the message.
	Text string
	// Tone is info, success, warning or danger.
	Tone Tone
	// Action is the words of the one action; "" for none.
	Action string
	// Timeout is how long a toast without an action stays; 4 s unless set.
	Timeout time.Duration
	// OnAction runs when the action is taken; OnDismiss when the toast goes,
	// by its timeout or its close control.
	OnAction, OnDismiss func()

	button *Button
	close  *IconButton
	gen    int
}

// NewToast returns an info toast saying words.
func NewToast(ui *UI, words string) *Toast {
	t := &Toast{ui: ui, Text: words, Tone: ToneInfo, Timeout: 4 * time.Second}
	t.Self = t
	t.button = NewButton(ui, "")
	t.button.Form = ButtonQuiet
	t.button.OnClick = func() {
		if t.OnAction != nil {
			t.OnAction()
		}
	}
	t.close = NewIconButton(ui, "close", "Dismiss")
	t.close.Size = Px(20)
	t.close.OnClick = t.dismiss
	t.AddChild(t.button)
	t.AddChild(t.close)
	t.SetLayout(syncing{Layout: toastLayout{t}, sync: func() {
		t.button.Text = t.Action
		t.button.Hidden, t.close.Hidden = t.Action == "", t.Action == ""
	}})
	t.DrawCallback = t.draw
	return t
}

// Shown announces the toast to a screen reader and starts its timeout; call
// it when the toast is put on the screen.
func (t *Toast) Shown() {
	unison.AnnounceForAccessibility(t.Text)
	t.gen++
	if t.Action != "" {
		return
	}
	gen := t.gen
	unison.InvokeTaskAfter(func() {
		if gen == t.gen {
			t.dismiss()
		}
	}, t.Timeout)
}

func (t *Toast) dismiss() {
	t.gen++
	if t.OnDismiss != nil {
		t.OnDismiss()
	}
}

func (t *Toast) words() *text.Layout {
	ui := t.ui
	return ui.Fonts.Layout([]text.Span{{Text: t.Text, Style: ui.Chrome(ui.Size(RoleBody), text.Regular, ui.Theme.Tokens().TextPrimary)}}, text.Options{})
}

type toastLayout struct{ t *Toast }

// parts are the widths across the toast: symbol, words, action, close.
func (l toastLayout) parts() (symbol, words, action, close float32) {
	t, m := l.t, l.t.ui.Interface
	symbol = float32(m.IconSizeSmall())
	words, _ = t.words().Size()
	if t.Action != "" {
		_, bp, _ := t.button.Sizes(geom.Size{})
		_, cp, _ := t.close.Sizes(geom.Size{})
		action, close = bp.Width, cp.Width
	}
	return
}

func (l toastLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.t.ui.Interface
	s, w, a, c := l.parts()
	gap := float32(m.Space())
	width := s + gap + w
	if a > 0 {
		width += gap + a + gap + c
	}
	size := geom.NewSize(width+2*float32(m.SpaceLoose()), float32(m.RowHeightSlim()+m.Space()))
	return size, size, size
}

func (l toastLayout) PerformLayout(target *unison.Panel) {
	t, m := l.t, l.t.ui.Interface
	r := target.ContentRect(false)
	s, w, a, c := l.parts()
	gap := float32(m.Space())
	x := r.X + float32(m.SpaceLoose()) + s + gap + w + gap
	mid := func(width float32, p unison.Paneler) {
		_, pref, _ := p.AsPanel().Sizes(geom.Size{})
		p.AsPanel().SetFrameRect(geom.NewRect(x, r.Y+(r.Height-pref.Height)/2, width, pref.Height))
		x += width + gap
	}
	if a > 0 {
		mid(a, t.button)
		mid(c, t.close)
	}
}

func (t *Toast) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := t.ui, t.ui.Theme.Tokens(), t.ui.Interface
	r := t.ContentRect(false)
	drawSurface(gc, ui, r)
	tone := messageTone(t.Tone, tk)
	// The tone stripe down the leading edge, a second channel beside the
	// symbol's colour, stopping short of the rounded corners.
	radius := float32(m.RadiusCard())
	painterFor(gc, ui).round(geom.NewRect(r.X+float32(m.Hairline()), r.Y+radius, float32(m.SpaceTight()), r.Height-2*radius),
		float32(m.RadiusBar()), tone)
	x := r.X + float32(m.SpaceLoose())
	s := float32(m.IconSizeSmall())
	if g, ok := icons.Glyph(toneSymbol(t.Tone)); ok {
		drawGlyph(gc, ui, g, s, tone, geom.NewRect(x, r.Y+(r.Height-s)/2, s, s))
	}
	l := t.words()
	_, h := l.Size()
	l.Draw(gc, x+s+float32(m.Space()), r.Y+(r.Height-h)/2)
}

// ProvideAccessibility reads the toast as its message, with its action.
func (t *Toast) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name, n.Description = t.Text, t.Action
}

// Notice is a message that stays: a banner across the top of a view saying
// something is wrong, out of date or waiting. A toast acknowledges something
// that finished; a notice is a condition still true. It can be dismissed
// only when dismissing it means something: "your licence expires in three
// days" can be, "this file could not be saved" cannot, because the condition
// outlives the reader's having seen it.
type Notice struct {
	unison.Panel
	ui *UI
	// Text is the condition; Detail, optional, says more.
	Text, Detail string
	// Tone is info, success, warning or danger.
	Tone Tone
	// Action is the words of the way out of the condition; "" for none.
	Action string
	// ActionHint says what pressing the action does, where its words cannot
	// (kvit-cash): the action's tooltip.
	ActionHint string
	// Dismissible draws a close control.
	Dismissible bool
	// OnAction and OnDismiss run when the action or the close control is
	// pressed.
	OnAction, OnDismiss func()

	text, detail *Label
	button       *Button
	close        *IconButton
}

// NewNotice returns an info notice.
func NewNotice(ui *UI, condition string) *Notice {
	n := &Notice{ui: ui, Text: condition, Tone: ToneInfo}
	n.Self = n
	n.text = NewLabel(ui, "")
	n.text.Wrap = true
	n.detail = NewLabel(ui, "")
	n.detail.Role, n.detail.Ink, n.detail.Wrap = RoleSmall, InkTextMuted, true
	n.button = NewButton(ui, "")
	n.button.OnClick = func() {
		if n.OnAction != nil {
			n.OnAction()
		}
	}
	n.close = NewIconButton(ui, "close", "Dismiss")
	n.close.OnClick = func() {
		if n.OnDismiss != nil {
			n.OnDismiss()
		}
	}
	for _, p := range []unison.Paneler{n.text, n.detail, n.button, n.close} {
		n.AddChild(p)
	}
	n.SetLayout(syncing{Layout: noticeLayout{n}, sync: func() {
		n.text.Text, n.detail.Text = n.Text, n.Detail
		n.detail.Hidden = n.Detail == ""
		n.button.Text, n.button.Explanation = n.Action, n.ActionHint
		n.button.Hidden = n.Action == ""
		n.close.Hidden = !n.Dismissible
	}})
	n.DrawCallback = n.draw
	return n
}

type noticeLayout struct{ n *Notice }

// arrange places the parts across a width and returns the notice's height.
func (l noticeLayout) arrange(r geom.Rect, set bool) float32 {
	n, m := l.n, l.n.ui.Interface
	gap := float32(m.Space())
	left, right := r.X+float32(m.SpaceLoose()), r.Right()-float32(m.SpaceNear())
	var controls []unison.Paneler
	if !n.button.Hidden {
		controls = append(controls, n.button)
	}
	if !n.close.Hidden {
		controls = append(controls, n.close)
	}
	// From the right: the close control, then the action.
	x := right
	var boxes []geom.Rect
	for i := len(controls) - 1; i >= 0; i-- {
		_, p, _ := controls[i].AsPanel().Sizes(geom.Size{})
		x -= p.Width
		boxes = append([]geom.Rect{geom.NewRect(x, 0, p.Width, p.Height)}, boxes...)
		x -= gap
	}
	textLeft := left + float32(m.IconSizeSmall()) + gap
	width := max(0, x-textLeft)
	_, tp, _ := n.text.Sizes(geom.NewSize(width, 0))
	h := tp.Height
	var dh float32
	if !n.detail.Hidden {
		_, dp, _ := n.detail.Sizes(geom.NewSize(width, 0))
		dh = dp.Height
		h += float32(m.SpaceTight()) + dh
	}
	height := max(float32(m.RowHeightSlim()+m.Space()), h+float32(m.Space()))
	if set {
		top := r.Y + (height-h)/2
		n.text.SetFrameRect(geom.NewRect(textLeft, top, width, tp.Height))
		n.detail.SetFrameRect(geom.NewRect(textLeft, top+tp.Height+float32(m.SpaceTight()), width, dh))
		for i, c := range controls {
			b := boxes[i]
			c.AsPanel().SetFrameRect(geom.NewRect(b.X, r.Y+(height-b.Height)/2, b.Width, b.Height))
		}
	}
	return height
}

func (l noticeLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.n.ui.Interface
	w := float32(m.Px(600))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := l.arrange(geom.NewRect(0, 0, w, 0), false)
	return geom.NewSize(0, h), geom.NewSize(float32(m.Px(600)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l noticeLayout) PerformLayout(target *unison.Panel) { l.arrange(target.ContentRect(false), true) }

func (n *Notice) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := n.ui, n.ui.Theme.Tokens(), n.ui.Interface
	p := painterFor(gc, ui)
	r := n.ContentRect(false)
	radius := float32(m.RadiusControl())
	tone := messageTone(n.Tone, tk)
	p.roundTint(r, radius, tone, 0.12)
	p.outline(r, radius, float32(m.Hairline()), tone)
	s := float32(m.IconSizeSmall())
	if g, ok := icons.Glyph(toneSymbol(n.Tone)); ok {
		// At the top, level with the first line, rather than the middle of a
		// notice several lines tall.
		top := n.text.FrameRect().Y + (n.text.FrameRect().Height-s)/2
		drawGlyph(gc, ui, g, s, tone, geom.NewRect(r.X+float32(m.SpaceLoose()), top, s, s))
	}
}

// ProvideAccessibility reads the notice as its condition and detail.
func (n *Notice) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	node := b.Node()
	if node.Role == role.Auto {
		node.Role = role.Group
	}
	node.Name = n.Text
	if n.Detail != "" {
		node.Name = n.Text + ". " + n.Detail
	}
}
