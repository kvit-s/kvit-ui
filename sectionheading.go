package kvitui

import (
	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/paintstyle"
	"github.com/richardwilkes/unison/enums/role"
)

// SectionHeading is a group heading: a filled bar with a disclosure chevron,
// the group's name, what it holds, the count with the word for what was
// counted, and the one action that applies to every row under it. It is a
// filled bar rather than a label with a rule running off to the right,
// because a rule between the words and the action reads as a divider between
// two unrelated things.
//
// The hoisted action has to be true for every row under the heading; if it is
// not, the group is wrong and wants splitting. It is a control of its own, a
// link showing its words or an icon button showing a symbol, and it takes its
// own keys, so running it never also opens or closes the group. A heading
// that collapses joins the tab order and opens and closes on Return, Enter,
// Space and a screen reader's press, answering a key only while it holds the
// keyboard itself.
type SectionHeading struct {
	control
	// Text names the group.
	Text string
	// Count is how many rows the group holds; below zero draws none. A group
	// of one is still counted, so the count does not come and go as the
	// group changes size.
	Count int
	// Counted is the word for one of them; "" leaves the number bare.
	Counted string
	// CountedPlural is the word for several; "" adds an s to Counted.
	CountedPlural string
	// CountText is a count the caller has already written out, for a count
	// that is not a number, such as "1 · +0 −0". It is drawn as written,
	// beside the name, and replaces Count.
	CountText string
	// Kind says what sort of thing the group holds, beside the name. It is
	// the part that gives way first when the bar is narrow.
	Kind string
	// Action is the hoisted action's words; "" for none.
	Action string
	// ActionSymbol draws the action as a symbol, keeping its words as the
	// button's name and tooltip, for an action a symbol can carry.
	ActionSymbol string
	// ActionExplanation says in a sentence what running the action does.
	ActionExplanation string
	// Strong draws the name bold and the rule above in the strong border
	// colour.
	Strong bool
	// Collapsible lets the heading open and close its group.
	Collapsible bool
	// Expanded is whether the group is open.
	Expanded bool
	// OnToggle runs after the heading opened or closed, with the new state.
	OnToggle func(expanded bool)
	// OnAction runs when the hoisted action is taken.
	OnAction func()

	link   *Link
	button *IconButton
}

// NewSectionHeading returns an open heading that does not collapse.
func NewSectionHeading(ui *UI, label string) *SectionHeading {
	h := &SectionHeading{Text: label, Count: -1, Expanded: true}
	h.Self = h
	h.initControl(ui, h.toggle, unison.KeyReturn, unison.KeyNumPadEnter, unison.KeySpace)
	keyDown := h.KeyDownCallback
	h.KeyDownCallback = func(key unison.KeyCode, mods mod.Modifiers, repeat bool) bool {
		if !h.Collapsible || !h.Focused() {
			return false
		}
		return keyDown(key, mods, repeat)
	}
	down := h.MouseDownCallback
	h.MouseDownCallback = func(where geom.Point, button, clicks int, mods mod.Modifiers) bool {
		return h.Collapsible && down(where, button, clicks, mods)
	}
	// The ring is drawn inside the bar: a heading sits flush against the rows
	// above and below, which would paint over a ring outside it.
	h.noRing = true
	act := func() {
		if h.OnAction != nil {
			h.OnAction()
		}
	}
	h.link = NewLink(ui, "")
	h.link.Role = RoleSmall
	h.link.OnActivate = act
	h.button = NewIconButton(ui, "dot", "")
	h.button.Dense = true
	h.button.OnClick = act
	h.AddChild(h.link)
	h.AddChild(h.button)
	h.SetLayout(syncing{Layout: headingLayout{h}, sync: h.sync})
	h.DrawCallback = h.draw
	h.DrawOverCallback = func(gc *unison.Canvas, _ geom.Rect) {
		if !h.KeyboardFocus() {
			return
		}
		w := float32(ui.Interface.FocusRingWidth())
		ring := h.ContentRect(true).Inset(geom.NewUniformInsets(w + w/2))
		paint := Color(ui.Theme.Tokens().FocusRing).Paint(gc, ring, paintstyle.Stroke)
		paint.SetStrokeWidth(w)
		gc.DrawRect(ring, paint)
	}
	return h
}

func (h *SectionHeading) toggle() {
	if !h.Collapsible {
		return
	}
	h.Expanded = !h.Expanded
	h.MarkForRedraw()
	if h.OnToggle != nil {
		h.OnToggle(h.Expanded)
	}
}

func (h *SectionHeading) sync() {
	h.SetFocusable(h.Collapsible)
	symbolic := h.Action != "" && h.ActionSymbol != ""
	h.link.Hidden = h.Action == "" || symbolic
	h.button.Hidden = !symbolic
	h.link.Text, h.link.Explanation = h.Action, h.ActionExplanation
	h.button.Label, h.button.Explanation = h.Action, h.ActionExplanation
	if symbolic {
		h.button.Symbol = h.ActionSymbol
	}
}

// CountPhrase is the count at the right end: "1 account", "1,200 accounts",
// the grouped number alone when nothing was named, and "" when there is no
// count or the caller wrote one out in CountText.
func (h *SectionHeading) CountPhrase() string {
	switch {
	case h.CountText != "" || h.Count < 0:
		return ""
	case h.Counted == "":
		return h.ui.Number(h.Count)
	}
	return h.ui.CountPhrase(h.Count, h.Counted, h.CountedPlural)
}

// nameStyle is the name's style: the small role, bold when strong.
func (h *SectionHeading) nameStyle() text.Style {
	weight := text.Regular
	if h.Strong {
		weight = text.Bold
	}
	return h.ui.Chrome(h.ui.Size(RoleSmall), weight, h.ui.Theme.Tokens().TextPrimary)
}

// texts are the name, the written count, the kind and the count phrase, at
// their natural widths.
func (h *SectionHeading) texts() (name, countText, kind, count *text.Layout) {
	ui := h.ui
	lay := func(s string, ink Ink, tabular bool) *text.Layout {
		st := ui.Chrome(ui.Size(RoleSmall), text.Regular, ink.Of(ui))
		st.Tabular = tabular
		return ui.Fonts.Layout([]text.Span{{Text: s, Style: st}}, text.Options{})
	}
	name = ui.Fonts.Layout([]text.Span{{Text: h.Text, Style: h.nameStyle()}}, text.Options{})
	return name, lay(h.CountText, InkTextFaint, true), lay(h.Kind, InkTextMuted, false), lay(h.CountPhrase(), InkTextFaint, true)
}

// headingParts is where each part of the bar goes.
type headingParts struct {
	chevron, name, countText, kind, count, action geom.Rect
}

// arrange lays the bar out as Qt's row does: a space in from each end, parts
// a stack gap apart, the empty room (or the kind, when there is one) taking
// what is left, and the action an extra stack gap from what is before it.
func (h *SectionHeading) arrange(b geom.Rect) headingParts {
	m := h.ui.Interface
	gap := float32(m.StackGap())
	name, countText, _, count := h.texts()
	var out headingParts
	x, right := b.X+float32(m.Space()), b.Right()-float32(m.Space())
	mid := func(w, hgt float32) geom.Rect { return geom.NewRect(0, b.Y+(b.Height-hgt)/2, w, hgt) }
	// From the right: the action, then the count.
	if !h.link.Hidden || !h.button.Hidden {
		var size geom.Size
		if !h.button.Hidden {
			s := float32(m.RowHeightCompact() - 2*m.FocusRingWidth())
			size = geom.NewSize(s, s)
		} else {
			_, size, _ = h.link.Sizes(geom.Size{})
		}
		out.action = mid(size.Width, size.Height)
		out.action.X = right - size.Width
		right = out.action.X - 2*gap
	}
	if cw, ch := count.Size(); h.CountPhrase() != "" {
		out.count = mid(cw, ch)
		out.count.X = right - cw
		right = out.count.X - gap
	}
	if h.Collapsible {
		s := float32(m.IconSizeSmall())
		out.chevron = mid(s, s)
		out.chevron.X = x
		x += s + gap
	}
	nw, nh := name.Size()
	var tw, th float32
	if h.CountText != "" {
		tw, th = countText.Size()
	}
	// The room left is shared by the kind and the empty room; the kind takes
	// all of it and is cut short, then the name gives way.
	fixed := nw
	if h.CountText != "" {
		fixed += gap + tw
	}
	room := right - x - fixed - gap
	if room < 0 {
		nw = max(0, nw+room)
		room = 0
	}
	out.name = mid(nw, nh)
	out.name.X = x
	x += nw + gap
	if h.CountText != "" {
		out.countText = mid(tw, th)
		out.countText.X = x
		x += tw + gap
	}
	out.kind = geom.NewRect(x, b.Y, room, b.Height)
	return out
}

type headingLayout struct{ h *SectionHeading }

func (l headingLayout) LayoutSizes(*unison.Panel, geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.h.ui.Interface
	hgt := float32(m.RowHeightCompact())
	return geom.NewSize(0, hgt), geom.NewSize(float32(m.Px(400)), hgt), geom.NewSize(unison.DefaultMaxSize, hgt)
}

func (l headingLayout) PerformLayout(target *unison.Panel) {
	at := l.h.arrange(target.ContentRect(false))
	l.h.link.SetFrameRect(at.action)
	l.h.button.SetFrameRect(at.action)
}

func (h *SectionHeading) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, t, m := h.ui, h.ui.Theme.Tokens(), h.ui.Interface
	p := painterFor(gc, ui)
	b := h.ContentRect(false)
	ground := t.PanelBackground
	if h.Collapsible && h.hovered {
		ground = t.HoverTint
	}
	p.fill(b, ground)
	rule := t.Border
	if h.Strong {
		rule = t.BorderStrong
	}
	p.fill(geom.NewRect(b.X, b.Y, b.Width, float32(m.Hairline())), rule)
	at := h.arrange(b)
	if h.Collapsible {
		symbol := "chevron-right"
		if h.Expanded {
			symbol = "chevron-down"
		}
		if g, ok := icons.Glyph(symbol); ok {
			drawGlyph(gc, ui, g, at.chevron.Width, t.TextFaint, at.chevron)
		}
	}
	_, countText, _, count := h.texts()
	put := func(l *text.Layout, box geom.Rect) {
		_, hgt := l.Size()
		l.Draw(gc, box.X, b.Y+(b.Height-hgt)/2)
	}
	put(ui.Fonts.Layout([]text.Span{{Text: h.Text, Style: h.nameStyle()}}, text.Options{MaxWidth: at.name.Width, Elide: true}), at.name)
	if h.CountText != "" {
		put(countText, at.countText)
	}
	if h.Kind != "" && at.kind.Width > 0 {
		st := ui.Chrome(ui.Size(RoleSmall), text.Regular, t.TextMuted)
		put(ui.Fonts.Layout([]text.Span{{Text: h.Kind, Style: st}}, text.Options{MaxWidth: at.kind.Width, Elide: true}), at.kind)
	}
	if h.CountPhrase() != "" {
		put(count, at.count)
	}
}

// ProvideAccessibility describes the heading. unison treats a heading as a
// line of text and leaves its children out, so a heading that holds its
// action cannot be one: a heading that opens and closes is a button with its
// expanded state, the usual form of a disclosure, and a fixed heading with an
// action is a group named by the heading. Only a fixed heading with no action
// is a heading.
func (h *SectionHeading) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		switch {
		case h.Collapsible:
			n.Role = role.Button
		case h.Action != "":
			n.Role = role.Group
		default:
			n.Role = role.Heading
		}
	}
	n.Name = h.Text
	if h.Collapsible {
		n.Expandable, n.Expanded = true, h.Expanded
		n.Description = "Collapsed"
		if h.Expanded {
			n.Description = "Expanded"
		}
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction opens or closes the group for a screen reader.
func (h *SectionHeading) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !h.Collapsible {
		return false
	}
	h.toggle()
	return true
}
