package kvitui

import (
	"math"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/role"
)

// Tone is what a chip's colour means, rather than the colour itself, which
// is what keeps one call site right in all four themes.
type Tone int

// The tones.
const (
	ToneNeutral Tone = iota
	ToneAccent
	ToneSuccess
	ToneWarning
	ToneDanger
	ToneInfo
)

func (t Tone) color(tk tokens.Tokens) palette.Color {
	switch t {
	case ToneAccent:
		return tk.Accent
	case ToneSuccess:
		return tk.Success
	case ToneWarning:
		return tk.Warning
	case ToneDanger:
		return tk.Danger
	case ToneInfo:
		return tk.Link
	}
	return tk.TextMuted
}

// Chip is a small labelled mark: what kind of thing this is, what state it is
// in, one word about it. Every tone has a second channel besides its hue: the
// outline is drawn in the tone as well as the ground tinted with it, and a
// strong chip is filled rather than tinted, so a reader who cannot tell the
// red tone from the amber one can still tell a warning from a failure.
type Chip struct {
	unison.Panel
	ui *UI
	// Text is the chip's word.
	Text string
	// Tone is what its colour means.
	Tone Tone
	// Strong fills the chip, for the one chip on a row that is its point.
	Strong bool
	// Symbol is an optional meaning name before the word.
	Symbol string
	// Explanation says in a sentence what the word means and what to do
	// about it: the tooltip and the accessible description.
	Explanation string
}

// NewChip returns a neutral chip.
func NewChip(ui *UI, word string) *Chip {
	c := &Chip{ui: ui, Text: word}
	c.Self = c
	c.SetSizer(c.sizes)
	c.DrawCallback = c.draw
	c.UpdateTooltipCallback = func(geom.Point, geom.Rect) geom.Rect {
		c.Tooltip = nil
		if c.Explanation != "" {
			c.Tooltip = newTooltip(ui, "", c.Explanation)
		}
		return c.RectToRoot(c.ContentRect(true))
	}
	return c
}

func (c *Chip) ink(tk tokens.Tokens) palette.Color {
	tone := c.Tone.color(tk)
	switch {
	case c.Strong:
		return tokens.LabelOn(tone)
	case c.Tone == ToneNeutral:
		return tk.TextSecondary
	}
	return tone
}

func (c *Chip) label() *text.Layout {
	ui := c.ui
	return ui.Fonts.Layout([]text.Span{{Text: c.Text, Style: ui.Chrome(ui.Size(RoleCaption), text.Regular, c.ink(ui.Theme.Tokens()))}}, text.Options{})
}

// content is the symbol and the word, a snug space apart.
func (c *Chip) content() float32 {
	m := c.ui.Interface
	w, _ := c.label().Size()
	if c.Symbol != "" {
		w += float32(m.Caption() + m.SpaceSnug())
	}
	return w
}

func (c *Chip) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := c.ui.Interface
	size := geom.NewSize(c.content()+2*float32(m.SpaceNear()), float32(m.ChipHeight()))
	return size, size, size
}

func (c *Chip) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := c.ui, c.ui.Theme.Tokens(), c.ui.Interface
	p := painterFor(gc, ui)
	r := c.ContentRect(false)
	radius := float32(m.RadiusChip())
	tone := c.Tone.color(tk)
	switch {
	case c.Strong:
		p.round(r, radius, tone)
	case c.Tone == ToneNeutral:
		p.round(r, radius, tk.ChipBackground)
	default:
		p.roundTint(r, radius, tone, 0.16)
	}
	if !c.Strong {
		edge := tone
		if c.Tone == ToneNeutral {
			edge = tk.Border
		}
		p.outline(r, radius, float32(m.Hairline()), edge)
	}
	x := r.X + (r.Width-c.content())/2
	if c.Symbol != "" {
		s := float32(m.Caption())
		if g, ok := icons.Glyph(c.Symbol); ok {
			drawGlyph(gc, ui, g, s, c.ink(tk), geom.NewRect(x, r.Y+(r.Height-s)/2, s, s))
		}
		x += s + float32(m.SpaceSnug())
	}
	l := c.label()
	_, h := l.Size()
	l.Draw(gc, x, r.Y+(r.Height-h)/2)
}

// ProvideAccessibility reads the chip as its word, with the explanation.
func (c *Chip) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	n.Name, n.Description = c.Text, c.Explanation
}

// Tag is a label a person put there: a tag, a category, a folder. A chip is
// chrome saying what something is; a tag is content the reader wrote, with a
// colour the reader chose from tokens.ColorPalette, which is the same in every
// theme because the reader picked it. Its words take a colour that contrasts
// with whatever was picked.
type Tag struct {
	unison.Panel
	ui *UI
	// Text is the tag.
	Text string
	// Tint is the colour the reader chose; the zero value is the neutral
	// chip ground.
	Tint palette.Color
	// Tinted says Tint is set.
	Tinted bool
	// Removable draws a close control, which calls OnRemove.
	Removable bool
	// OnRemove runs when the close control is pressed.
	OnRemove func()

	remove *IconButton
}

// NewTag returns an untinted tag.
func NewTag(ui *UI, label string) *Tag {
	t := &Tag{ui: ui, Text: label}
	t.Self = t
	t.remove = NewIconButton(ui, "close", "")
	t.remove.Size = SizeCaption
	t.remove.OnClick = func() {
		if t.OnRemove != nil {
			t.OnRemove()
		}
	}
	t.remove.SetLayoutData(&unison.FlexLayoutData{VAlign: align.Middle})
	t.AddChild(t.remove)
	t.SetLayout(syncing{Layout: tagLayout{t}, sync: func() {
		t.remove.Hidden = !t.Removable
		t.remove.Label = "Remove " + t.Text
	}})
	t.DrawCallback = t.draw
	return t
}

// SetTint gives the tag a colour from tokens.ColorPalette.
func (t *Tag) SetTint(c palette.Color) { t.Tint, t.Tinted = c, true }

func (t *Tag) label() *text.Layout {
	ui, tk := t.ui, t.ui.Theme.Tokens()
	ink := tk.TextSecondary
	if t.Tinted {
		ink = tk.TextPrimary
	}
	return ui.Fonts.Layout([]text.Span{{Text: t.Text, Style: ui.Chrome(ui.Size(RoleCaption), text.Regular, ink)}}, text.Options{})
}

type tagLayout struct{ t *Tag }

func (l tagLayout) content() float32 {
	t, m := l.t, l.t.ui.Interface
	w, _ := t.label().Size()
	if t.Removable {
		w += float32(m.SpaceSnug() + m.Caption())
	}
	return w
}

func (l tagLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.t.ui.Interface
	size := geom.NewSize(l.content()+2*float32(m.SpaceNear()), float32(m.TagHeight()))
	return size, size, size
}

func (l tagLayout) PerformLayout(target *unison.Panel) {
	t, m := l.t, l.t.ui.Interface
	r := target.ContentRect(false)
	s := float32(m.Caption())
	right := r.X + (r.Width+l.content())/2
	t.remove.SetFrameRect(geom.NewRect(right-s, r.Y+(r.Height-s)/2, s, s))
}

func (t *Tag) draw(gc *unison.Canvas, _ geom.Rect) {
	tk, m := t.ui.Theme.Tokens(), t.ui.Interface
	p := painterFor(gc, t.ui)
	r := t.ContentRect(false)
	radius := float32(m.RadiusPill())
	edge := tk.Border
	if t.Tinted {
		p.roundTint(r, radius, t.Tint, 0.22)
		edge = t.Tint
	} else {
		p.round(r, radius, tk.ChipBackground)
	}
	p.outline(r, radius, float32(m.Hairline()), edge)
	l := t.label()
	_, h := l.Size()
	l.Draw(gc, r.X+(r.Width-tagLayout{t}.content())/2, r.Y+(r.Height-h)/2)
}

// ProvideAccessibility reads the tag as its words and its colour's name,
// never the colour's value: "#8a5cc0" tells a reader nothing.
func (t *Tag) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		// A removable tag holds its close control, which a text node would
		// leave out.
		n.Role = role.Label
		if t.Removable {
			n.Role = role.Group
		}
	}
	n.Name = t.Text
	if t.Tinted {
		n.Name = t.Text + ", " + tokens.ColorName(t.Tint.Hex())
	}
}

// Slug is an identifier: a reference, a short hash, a key, a ticket number.
// It is monospace, because what a reader does with an identifier is compare
// it with another, and two proportional strings of the same length are
// different widths. It is cut short in the middle rather than at the end,
// since the end of a hash is what tells it from its neighbours.
type Slug struct {
	unison.Panel
	ui *UI
	// Text is the identifier.
	Text string
	// Ground draws the inline-code ground; true unless turned off, for a slug
	// in a table cell where a ground on every row reads as a column of boxes.
	Ground bool
}

// NewSlug returns a slug on its ground.
func NewSlug(ui *UI, id string) *Slug {
	s := &Slug{ui: ui, Text: id, Ground: true}
	s.Self = s
	s.SetSizer(s.sizes)
	s.DrawCallback = s.draw
	return s
}

func (s *Slug) layout(width float32) *text.Layout {
	ui := s.ui
	return ui.Fonts.Layout([]text.Span{{Text: s.Text, Style: ui.Mono(ui.Size(RoleSmall), ui.Theme.Tokens().TextSecondary)}},
		text.Options{MaxWidth: width, Elide: width > 0, ElideMiddle: true})
}

func (s *Slug) pad() float32 {
	if s.Ground {
		return float32(s.ui.Interface.SpaceNear())
	}
	return 0
}

func (s *Slug) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	w, h := s.layout(0).Size()
	if s.Ground {
		h = float32(s.ui.Interface.ChipHeight())
	}
	pad := 2 * s.pad()
	return geom.NewSize(pad, h), geom.NewSize(w+pad, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (s *Slug) draw(gc *unison.Canvas, _ geom.Rect) {
	r := s.ContentRect(false)
	if s.Ground {
		painterFor(gc, s.ui).round(r, float32(s.ui.Interface.RadiusBar()), s.ui.Theme.Tokens().InlineCodeBackground)
	}
	l := s.layout(max(0, r.Width-2*s.pad()))
	_, h := l.Size()
	l.Draw(gc, r.X+s.pad(), r.Y+(r.Height-h)/2)
}

// ProvideAccessibility reads the slug as the whole identifier, however it is
// cut short on the screen.
func (s *Slug) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	n.Name = s.Text
}

// Shape is a mark's second channel beside its colour: a state told only by
// being red says nothing to a reader who cannot tell red from amber.
type Shape int

// The shapes.
const (
	ShapeCircle Shape = iota
	ShapeSquare
	ShapeDiamond
)

// drawShape draws a filled or hollow shape in a box.
func drawShape(gc *unison.Canvas, ui *UI, r geom.Rect, shape Shape, hollow bool, c palette.Color) {
	m := ui.Interface
	style := paintstyle.Fill
	if hollow {
		style = paintstyle.Stroke
		w := float32(m.Hairline())
		r = r.Inset(geom.NewUniformInsets(w / 2))
	}
	paint := Color(c).Paint(gc, r, style)
	if hollow {
		paint.SetStrokeWidth(float32(m.Hairline()))
	}
	switch shape {
	case ShapeCircle:
		gc.DrawOval(r, paint)
	case ShapeDiamond:
		gc.Save()
		gc.Translate(r.Center())
		gc.Rotate(45)
		// The square turned on its corner. Its box is the turned square's
		// width, so it is drawn at the size of the other shapes' square.
		s := r.Width / float32(math.Sqrt2)
		radius := float32(m.Hairline())
		gc.DrawRoundedRect(geom.NewRect(-s/2, -s/2, s, s), geom.NewSize(radius, radius), paint)
		gc.Restore()
	default:
		radius := float32(m.Hairline())
		gc.DrawRoundedRect(r, geom.NewSize(radius, radius), paint)
	}
}

// Dot is a small shape standing for one thing's state: a health level, a
// severity, whether something is running. It is drawn as geometry rather than
// as a glyph, since a circle six pixels across is exact as geometry and
// approximate as a glyph. Shape is its second channel beside the colour.
type Dot struct {
	unison.Panel
	ui *UI
	// Ink is the colour; the muted text colour unless set.
	Ink Ink
	// Shape is circle, square or diamond.
	Shape Shape
	// Hollow draws the outline only.
	Hollow bool
	// Label says what the dot means, for a screen reader; "" hides it from
	// one.
	Label string
	// Size is the dot's side; a near space unless set.
	Size Measure
}

// NewDot returns a muted circle.
func NewDot(ui *UI) *Dot {
	d := &Dot{ui: ui}
	d.Self = d
	d.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		s := float32(ui.Interface.SpaceNear())
		if d.Size != nil {
			s = float32(d.Size.Of(ui))
		}
		if d.Shape == ShapeDiamond {
			// A square of that side turned on its corner, which Qt draws
			// past the dot's box; unison clips to the box, so the box grows.
			s *= float32(math.Sqrt2)
		}
		size := geom.NewSize(s, s)
		return size, size, size
	})
	d.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		ink := d.Ink
		if ink == nil {
			ink = InkTextMuted
		}
		drawShape(gc, ui, d.ContentRect(false), d.Shape, d.Hollow, ink.Of(ui))
	}
	return d
}

// ProvideAccessibility names the dot by its meaning, or hides it.
func (d *Dot) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if d.Label == "" {
		n.Ignored = true
		return
	}
	if n.Role == role.Auto {
		n.Role = role.Image
	}
	n.Name = d.Label
}

// Signal is a mark saying what state something is in and how many things are
// in it. A dot says which state and carries no number; a badge always draws
// its digits, so a column of rows each with one thing running reads as a
// column of ones. A signal draws the number only past one: at one, the mark
// is the whole statement.
type Signal struct {
	unison.Panel
	ui *UI
	// Count is how many are in the state; zero or less hides the mark.
	Count int
	// Max is the largest count written out; 99 unless set.
	Max int
	// Ink is the state's colour, the application's own idea; the accent
	// unless set.
	Ink Ink
	// Shape is square or circle; a diamond would turn the number with it.
	Shape Shape
	// Hollow draws the outline only, a third distinction without another hue.
	Hollow bool
	// Label says what the mark means, "2 agents running": the colour and the
	// number do not say it.
	Label string
}

// NewSignal returns a square accent signal for one thing.
func NewSignal(ui *UI, label string) *Signal {
	s := &Signal{ui: ui, Count: 1, Max: 99, Shape: ShapeSquare, Label: label}
	s.Self = s
	s.SetSizer(s.sizes)
	s.DrawCallback = s.draw
	return s
}

func (s *Signal) ink() palette.Color {
	if s.Ink == nil {
		return InkAccent.Of(s.ui)
	}
	return s.Ink.Of(s.ui)
}

func (s *Signal) digits() *text.Layout {
	ui := s.ui
	most := s.Max
	if most <= 0 {
		most = 99
	}
	words := ui.Number(s.Count)
	if s.Count > most {
		words = ui.Number(most) + "+"
	}
	c := s.ink()
	if !s.Hollow {
		c = tokens.LabelOn(c)
	}
	st := ui.Chrome(ui.Size(RoleCaption), text.Bold, c)
	st.Tabular = true
	return ui.Fonts.Layout([]text.Span{{Text: words, Style: st}}, text.Options{})
}

func (s *Signal) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	if s.Count <= 0 {
		return geom.Size{}, geom.Size{}, geom.Size{}
	}
	m := s.ui.Interface
	h := float32(m.Caption() + m.SpaceSnug())
	// Measured with the digits whether or not they are drawn, so the mark
	// keeps one width.
	w, _ := s.digits().Size()
	size := geom.NewSize(max(h, w+float32(m.SpaceSnug())), h)
	return size, size, size
}

func (s *Signal) draw(gc *unison.Canvas, _ geom.Rect) {
	if s.Count <= 0 {
		return
	}
	m := s.ui.Interface
	r := s.ContentRect(false)
	style := paintstyle.Fill
	if s.Hollow {
		style = paintstyle.Stroke
	}
	box := r
	if s.Hollow {
		box = r.Inset(geom.NewUniformInsets(float32(m.Hairline()) / 2))
	}
	paint := Color(s.ink()).Paint(gc, box, style)
	paint.SetStrokeWidth(float32(m.Hairline()))
	radius := float32(m.RadiusBar())
	if s.Shape == ShapeCircle {
		radius = r.Height / 2
	}
	gc.DrawRoundedRect(box, geom.NewSize(radius, radius), paint)
	if s.Count > 1 {
		l := s.digits()
		w, h := l.Size()
		l.Draw(gc, r.X+(r.Width-w)/2, r.Y+(r.Height-h)/2)
	}
}

// ProvideAccessibility reads the signal as its meaning, never as the bare
// number inside it.
func (s *Signal) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if s.Count <= 0 {
		n.Ignored = true
		return
	}
	if n.Role == role.Auto {
		n.Role = role.Label
	}
	n.Name = s.Label
}

// Pip is a row of dots standing for a small count: three of five days
// recorded, two of four checks passed. It is for counts a reader takes in
// without counting, up to about seven; past that a Figure is faster to read.
// Filled and hollow dots rather than a bar, because a pip row answers "how
// many of these" and a bar "how far along".
type Pip struct {
	unison.Panel
	ui *UI
	// Filled of Total dots are filled.
	Filled, Total int
	// Ink is the filled dots' colour; the accent unless set.
	Ink Ink
	// Label says what is counted; "3 of 5" unless set.
	Label string
}

// NewPip returns a row of total dots, filled of them filled.
func NewPip(ui *UI, filled, total int) *Pip {
	p := &Pip{ui: ui, Filled: filled, Total: total}
	p.Self = p
	p.SetSizer(func(geom.Size) (geom.Size, geom.Size, geom.Size) {
		m := ui.Interface
		n := max(0, p.Total)
		w := float32(n*m.SpaceSnug() + max(0, n-1)*m.SpaceTight())
		size := geom.NewSize(w, float32(m.SpaceSnug()))
		return size, size, size
	})
	p.DrawCallback = func(gc *unison.Canvas, _ geom.Rect) {
		m, tk := ui.Interface, ui.Theme.Tokens()
		r := p.ContentRect(false)
		s := float32(m.SpaceSnug())
		ink := p.Ink
		if ink == nil {
			ink = InkAccent
		}
		for i := 0; i < p.Total; i++ {
			box := geom.NewRect(r.X+float32(i)*(s+float32(m.SpaceTight())), r.Y+(r.Height-s)/2, s, s)
			if i < p.Filled {
				drawShape(gc, ui, box, ShapeCircle, false, ink.Of(ui))
			} else {
				drawShape(gc, ui, box, ShapeCircle, true, tk.Border)
			}
		}
	}
	return p
}

// ProvideAccessibility reads the row as a progress bar, "3 of 5".
func (p *Pip) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.ProgressBar
	}
	n.Name = p.Label
	if n.Name == "" {
		n.Name = p.ui.Number(p.Filled) + " of " + p.ui.Number(p.Total)
	}
	n.HasNumber = true
	n.Number, n.Min, n.Max = float64(p.Filled), 0, float64(p.Total)
}
