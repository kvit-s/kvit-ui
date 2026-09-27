package kvitui

import (
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/role"
)

// TimelineEntry is one thing that happened: when, what, who did it, an
// optional line of detail, and a tone for how it went.
type TimelineEntry struct {
	When, What, Who, Detail string
	Tone                    Tone
}

// Timeline is what happened to something, newest first: a reader opening a
// history is asking what happened recently, and making them scroll to the
// bottom to find out is the commonest thing a history gets wrong. Each entry
// says who did it as well as what and when; where an agent and a person both
// change the same things, who is the column that makes a history worth
// reading. With nothing in it, it says so.
type Timeline struct {
	unison.Panel
	ui *UI
	// Label names the history for a screen reader.
	Label string
	// Entries are the events, newest first.
	Entries []TimelineEntry

	empty *EmptyState
}

// NewTimeline returns a history of entries, newest first.
func NewTimeline(ui *UI, label string, entries ...TimelineEntry) *Timeline {
	t := &Timeline{ui: ui, Label: label, Entries: entries}
	t.Self = t
	t.empty = NewEmptyState(ui, "Nothing has happened yet")
	t.empty.Detail = "Changes will be listed here, newest first, with who made them."
	t.empty.Symbol = "clock"
	t.AddChild(t.empty)
	t.SetLayout(syncing{Layout: timelineLayout{t}, sync: func() { t.empty.Hidden = len(t.Entries) > 0 }})
	t.DrawCallback = t.draw
	return t
}

// entryParts are the laid-out words of one entry.
type entryParts struct {
	what, who, when, detail *text.Layout
	line                    float32 // the headline's height
	height                  float32 // the entry's height
}

func (t *Timeline) parts(e TimelineEntry, width float32) entryParts {
	ui, tk, m := t.ui, t.ui.Theme.Tokens(), t.ui.Interface
	layout := func(s string, r TypeRole, ink palette.Color, tabular bool, wrap float32) *text.Layout {
		st := ui.Chrome(ui.Size(r), text.Regular, ink)
		st.Tabular = tabular
		return ui.Fonts.Layout([]text.Span{{Text: s, Style: st}}, text.Options{MaxWidth: wrap})
	}
	var p entryParts
	p.what = layout(e.What, RoleBody, tk.TextPrimary, false, 0)
	p.who = layout(e.Who, RoleSmall, tk.TextMuted, false, 0)
	p.when = layout(e.When, RoleSmall, tk.TextFaint, true, 0)
	for _, l := range []*text.Layout{p.what, p.who, p.when} {
		_, h := l.Size()
		p.line = max(p.line, h)
	}
	body := p.line
	if e.Detail != "" {
		p.detail = layout(e.Detail, RoleSmall, tk.TextMuted, false, max(1, width-float32(m.Px(28)+m.SpaceNear())))
		_, h := p.detail.Size()
		body += float32(m.SpaceTight()) + h
	}
	p.height = max(float32(m.RowHeightSlim()), body+float32(m.Space()))
	return p
}

type timelineLayout struct{ t *Timeline }

func (l timelineLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	t := l.t
	width := hint.Width
	if width <= 0 {
		width = float32(t.ui.Interface.Px(400))
	}
	if len(t.Entries) == 0 {
		_, p, _ := t.empty.Sizes(geom.NewSize(width, 0))
		return geom.NewSize(0, p.Height), geom.NewSize(p.Width, p.Height), geom.NewSize(unison.DefaultMaxSize, p.Height)
	}
	h := float32(0)
	for _, e := range t.Entries {
		h += t.parts(e, width).height
	}
	return geom.NewSize(0, h), geom.NewSize(width, h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l timelineLayout) PerformLayout(target *unison.Panel) {
	t := l.t
	r := target.ContentRect(false)
	_, p, _ := t.empty.Sizes(geom.NewSize(r.Width, 0))
	t.empty.SetFrameRect(geom.NewRect(r.X, r.Y, r.Width, p.Height))
}

// toneMark is the colour and shape a tone's mark takes: the shape carries
// the tone as well as the colour, so a reader who cannot separate the hues
// still sees that one entry is unlike its neighbours.
func toneMark(ui *UI, tone Tone) (palette.Color, Shape) {
	tk := ui.Theme.Tokens()
	switch tone {
	case ToneSuccess:
		return tk.Success, ShapeCircle
	case ToneWarning:
		return tk.Warning, ShapeSquare
	case ToneDanger:
		return tk.Danger, ShapeDiamond
	}
	return tk.TextMuted, ShapeCircle
}

func (t *Timeline) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := t.ui, t.ui.Theme.Tokens(), t.ui.Interface
	p := painterFor(gc, ui)
	r := t.ContentRect(false)
	y := r.Y
	n := len(t.Entries)
	hair := float32(m.Hairline())
	near, snug := float32(m.SpaceNear()), float32(m.SpaceSnug())
	for i, e := range t.Entries {
		parts := t.parts(e, r.Width)
		// The mark sits on the middle of the entry's first line, and the
		// spine runs from the first mark to the last and no further, since
		// a line past the oldest event says there is more below.
		mark := y + snug + parts.line/2
		top, bottom := y, y+parts.height
		if i == 0 {
			top = mark
		}
		if i == n-1 {
			bottom = mark
		}
		if n > 1 && bottom > top {
			p.fill(geom.NewRect(r.X+near, top, hair, bottom-top), tk.Border)
		}
		ink, shape := toneMark(ui, e.Tone)
		drawShape(gc, ui, geom.NewRect(r.X+snug, mark-near/2, near, near), shape, false, ink)
		x := r.X + float32(m.Px(28))
		for _, l := range []*text.Layout{parts.what, parts.who, parts.when} {
			w, _ := l.Size()
			if w == 0 {
				continue
			}
			l.Draw(gc, x, y+snug)
			x += w + float32(m.Space())
		}
		if parts.detail != nil {
			parts.detail.Draw(gc, r.X+float32(m.Px(28)), y+snug+parts.line+float32(m.SpaceTight()))
		}
		y += parts.height
	}
}

// ProvideAccessibility reads the history as a list whose items say when,
// what and who.
func (t *Timeline) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.List
	}
	n.Name = t.Label
	r := t.ContentRect(false)
	y := r.Y
	for i, e := range t.Entries {
		h := t.parts(e, r.Width).height
		box := geom.NewRect(r.X, y, r.Width, h)
		b.AddVirtualChild(i, func(v *accessibility.Node) {
			v.Role = role.ListItem
			v.Name = e.When + ", " + e.What + ", " + e.Who
			v.Description = e.Detail
			v.RowIndex = i
			v.Bounds = box
		})
		y += h
	}
}
