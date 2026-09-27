package kvitui

import (
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/role"
)

// Tab is one tab in a row of them: a different view of the same thing. The
// selected tab is marked by an underline as well as by colour and weight, so
// the selection does not rest on colour alone, which is where that rule is
// most often broken: a coloured tab label looks obviously different to
// whoever drew it.
type Tab struct {
	control
	// Text names the view.
	Text string
	// Selected marks the tab whose view is showing.
	Selected bool
	// Count is how many things the view holds, drawn beside the name; below
	// zero draws none.
	Count int
	// Explanation says in one sentence what the view shows, which two or
	// three words on a tab cannot. It is the tooltip and the accessible
	// description, and never a substitute for the name: a tab whose subject
	// is only in its explanation cannot be chosen without hovering it.
	Explanation string
	// OnClick runs when the tab is pressed.
	OnClick func()
}

// NewTab returns an unselected tab with no count.
func NewTab(ui *UI, label string) *Tab {
	t := &Tab{Text: label, Count: -1}
	t.Self = t
	t.initControl(ui, func() {
		if t.OnClick != nil {
			t.OnClick()
		}
	}, unison.KeySpace)
	// The ring surrounds the tab's square ground.
	t.ringRadius = func() float32 { return 0 }
	t.SetSizer(t.sizes)
	t.DrawCallback = t.draw
	t.UpdateTooltipCallback = func(geom.Point, geom.Rect) geom.Rect {
		if t.Explanation != "" {
			t.Tooltip = newTooltip(ui, "", t.Explanation)
		} else {
			t.Tooltip = nil
		}
		return t.RectToRoot(t.ContentRect(false))
	}
	return t
}

// layouts are the name and, when there is one, the count.
func (t *Tab) layouts() (name, count *text.Layout) {
	ui, tk := t.ui, t.ui.Theme.Tokens()
	weight, ink := text.Regular, tk.TextMuted
	if t.Selected {
		weight, ink = text.Bold, tk.TextPrimary
	}
	name = ui.Fonts.Layout([]text.Span{{Text: t.Text, Style: ui.Chrome(ui.Size(RoleBody), weight, ink)}}, text.Options{})
	if t.Count >= 0 {
		st := ui.Chrome(ui.Size(RoleCaption), text.Regular, tk.TextFaint)
		st.Tabular = true
		count = ui.Fonts.Layout([]text.Span{{Text: ui.Number(t.Count), Style: st}}, text.Options{})
	}
	return name, count
}

// content is the name, and the count a near space after it: their width and
// height together, and each laid out.
func (t *Tab) content() (width, height float32, name, count *text.Layout) {
	name, count = t.layouts()
	width, height = name.Size()
	if count != nil {
		cw, ch := count.Size()
		width += float32(t.ui.Interface.SpaceNear()) + cw
		height = max(height, ch)
	}
	return width, height, name, count
}

func (t *Tab) sizes(geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := t.ui.Interface
	w, h, _, _ := t.content()
	size := geom.NewSize(w+2*float32(m.SpaceLoose()), max(float32(m.TabHeight()), h))
	return size, size, size
}

func (t *Tab) draw(gc *unison.Canvas, _ geom.Rect) {
	ui, tk, m := t.ui, t.ui.Theme.Tokens(), t.ui.Interface
	p := painterFor(gc, ui)
	b := t.ContentRect(false)
	if t.hovered && !t.Selected {
		p.fill(b, tk.HoverTint)
	}
	if t.Selected {
		// Two design pixels, so it survives the smallest interface size and
		// does not read as a hairline rule.
		line := float32(m.SpaceTight())
		p.fill(geom.NewRect(b.X, b.Bottom()-line, b.Width, line), tk.Accent)
	}
	w, _, name, count := t.content()
	x := b.X + (b.Width-w)/2
	_, nh := name.Size()
	name.Draw(gc, x, b.Y+(b.Height-nh)/2)
	if count != nil {
		nw, _ := name.Size()
		_, ch := count.Size()
		count.Draw(gc, x+nw+float32(m.SpaceNear()), b.Y+(b.Height-ch)/2)
	}
}

// ProvideAccessibility announces a tab, with its count in its name, "Projects,
// 26 items", grouped as the reader groups digits, and says whether it is the
// selected one.
func (t *Tab) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Tab
	}
	n.Name = t.Text
	if t.Count >= 0 {
		n.Name = t.Text + ", " + t.ui.CountPhrase(t.Count, "item", "")
	}
	n.Description = t.Explanation
	n.Selectable = true
	n.Selected = t.Selected
	if t.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction presses the tab for a screen reader.
func (t *Tab) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !t.Enabled() {
		return false
	}
	t.fire()
	return true
}
