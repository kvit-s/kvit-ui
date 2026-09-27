package kvitui

import (
	"math"
	"time"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// ConfirmInPlace is a strip that appears where the reader is, saying what
// just happened and offering to undo it: the undo goes inside the
// confirmation rather than a confirmation coming before the action. When the
// reader does the same thing forty times, a dialog before each costs a press
// every time and stops being read by the third, while a confirmation after,
// with an undo, costs nothing unless something went wrong. That makes it
// right only for an action the application can reverse. It stays until it is
// dismissed or replaced: an undo that expires while the reader is checking
// what happened is an undo they cannot use, and a timed message is a Toast.
type ConfirmInPlace struct {
	unison.Panel
	ui *UI
	// Text says what happened.
	Text string
	// UndoText is the words on the undo; "" draws no undo, for an action that
	// has stopped being reversible.
	UndoText string
	// Shown opens the strip, growing to its height; false closes it.
	Shown bool
	// Affected is how many things the action touched, so the reader can
	// check the count before undoing; below zero says nothing.
	Affected int
	// OnUndo and OnDismiss run when the undo or the close button is pressed.
	OnUndo, OnDismiss func()

	undo   *Button
	close  *IconButton
	open   float32 // how far open the strip is drawn, 0 to 1
	was    bool    // Shown when last laid out
	moving bool
}

// NewConfirmInPlace returns a closed strip saying text, with an undo.
func NewConfirmInPlace(ui *UI, text string) *ConfirmInPlace {
	c := &ConfirmInPlace{ui: ui, Text: text, UndoText: "Undo", Affected: -1}
	c.Self = c
	c.undo = NewButton(ui, "")
	c.undo.Form = ButtonQuiet
	c.undo.OnClick = func() {
		if c.OnUndo != nil {
			c.OnUndo()
		}
	}
	c.close = NewIconButton(ui, "close", "Dismiss")
	c.close.OnClick = func() {
		if c.OnDismiss != nil {
			c.OnDismiss()
		}
	}
	c.AddChild(c.undo)
	c.AddChild(c.close)
	c.SetLayout(syncing{Layout: confirmLayout{c}, sync: c.sync})
	c.DrawCallback = c.draw
	return c
}

func (c *ConfirmInPlace) sync() {
	c.undo.Text = c.UndoText
	c.undo.Hidden = c.UndoText == ""
	if c.Shown != c.was {
		c.was = c.Shown
		if c.Shown {
			unison.AnnounceForAccessibility(c.sentence())
		}
		c.grow()
	}
}

// grow opens or closes the strip over 140 ms, easing out, or at once when
// motion is reduced.
func (c *ConfirmInPlace) grow() {
	target := func() float32 {
		if c.Shown {
			return 1
		}
		return 0
	}
	scale := c.ui.Theme.MotionScale()
	if scale == 0 || c.Window() == nil {
		c.open = target()
		return
	}
	if c.moving {
		return
	}
	c.moving = true
	from, start, to := c.open, time.Now(), target()
	var step func()
	step = func() {
		if t := target(); t != to {
			from, start, to = c.open, time.Now(), t
		}
		s := min(1, float64(time.Since(start))/(140*scale*float64(time.Millisecond)))
		c.open = from + (to-from)*float32(1-math.Pow(1-s, 3))
		if s >= 1 {
			c.open, c.moving = to, false
		}
		c.MarkForLayoutAndRedraw()
		if w := c.Window(); w != nil {
			w.Content().MarkForLayoutRecursively()
		}
		if c.moving {
			unison.InvokeTaskAfter(step, 16*time.Millisecond)
		}
	}
	unison.InvokeTaskAfter(step, 16*time.Millisecond)
}

// phrase is the count of what was touched: "1 item", "250,000 items".
func (c *ConfirmInPlace) phrase() string { return c.ui.CountPhrase(c.Affected, "item", "") }

// sentence is what the strip says, drawn and announced alike, so the two
// cannot disagree about the same number.
func (c *ConfirmInPlace) sentence() string {
	if c.Affected >= 0 {
		return c.Text + " — " + c.phrase()
	}
	return c.Text
}

type confirmLayout struct{ c *ConfirmInPlace }

func (l confirmLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.c.ui.Interface
	h := float32(math.Round(float64(float32(m.RowHeightSlim()+m.Space()) * l.c.open)))
	return geom.NewSize(0, h), geom.NewSize(float32(m.Px(400)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l confirmLayout) PerformLayout(target *unison.Panel) {
	c, m := l.c, l.c.ui.Interface
	r := target.ContentRect(false)
	full := float32(m.RowHeightSlim() + m.Space())
	x := r.Right() - float32(m.SpaceNear())
	for _, p := range []unison.Paneler{c.close, c.undo} {
		if p.AsPanel().Hidden {
			continue
		}
		size := preferred(p)
		x -= size.Width
		p.AsPanel().SetFrameRect(geom.NewRect(x, r.Y+(full-size.Height)/2, size.Width, size.Height))
		x -= float32(m.Space())
	}
}

func (c *ConfirmInPlace) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := c.ui, c.ui.Theme.Tokens(), c.ui.Interface
	p := painterFor(gc, ui)
	r := c.ContentRect(false)
	if r.Height <= 0 {
		return
	}
	radius := float32(m.RadiusControl())
	p.roundTint(r, radius, t.Success, 0.12)
	p.outline(r, radius, float32(m.Hairline()), t.Success)
	full := float32(m.RowHeightSlim() + m.Space())
	x := r.X + float32(m.SpaceLoose())
	s := float32(m.IconSizeSmall())
	if g, ok := icons.Glyph("success"); ok {
		drawGlyph(gc, ui, g, s, t.Success, geom.NewRect(x, r.Y+(full-s)/2, s, s))
	}
	x += s + float32(m.Space())
	right := c.undo.FrameRect().X
	if c.undo.Hidden {
		right = c.close.FrameRect().X
	}
	l := ui.Fonts.Layout([]text.Span{{Text: c.sentence(), Style: ui.Chrome(ui.Size(RoleBody), text.Regular, t.TextPrimary)}},
		text.Options{MaxWidth: max(1, right-float32(m.Space())-x), Elide: true})
	_, h := l.Size()
	l.Draw(gc, x, r.Y+(full-h)/2)
}

// ProvideAccessibility reads the strip as what happened and how much it
// touched, holding its undo and close buttons.
func (c *ConfirmInPlace) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Group
	}
	n.Name = c.Text
	if c.Affected >= 0 {
		n.Name = c.Text + ", " + c.phrase() + " affected"
	}
	if !c.Shown {
		n.Ignored = true
	}
}
