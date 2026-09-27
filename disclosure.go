package kvitui

import (
	"math"
	"time"

	"github.com/kvit-s/kvit-ui/icons"
	"github.com/kvit-s/kvit-ui/text"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
)

// Disclosure is a trigger with a body under it that opens and closes. The
// chevron turns rather than being swapped for another, so the reader sees
// the state change happen, and the body grows open rather than appearing;
// both follow the reduced-motion setting.
//
// Sections that share a Group behave as an accordion: opening one closes the
// others in the same window. The Qt disclosure declares the same property,
// and leaves the closing to its callers.
type Disclosure struct {
	unison.Panel
	ui *UI
	// Title names the section.
	Title string
	// Group makes sections sharing it an accordion; "" stands alone.
	Group string
	// Count is how many things are inside, beside the title; below zero
	// draws none.
	Count int
	// OnToggle runs after the section opened or closed, with the new state.
	OnToggle func(expanded bool)

	expanded bool
	open     float32 // how far open the section is drawn, 0 to 1
	moving   bool
	trigger  *disclosureTrigger
	body     *unison.Panel
}

// disclosureTrigger is the row that opens and closes the section.
type disclosureTrigger struct {
	control
	d *Disclosure
}

// NewDisclosure returns a closed section holding content.
func NewDisclosure(ui *UI, title string, content ...unison.Paneler) *Disclosure {
	d := &Disclosure{ui: ui, Title: title, Count: -1}
	d.Self = d
	t := &disclosureTrigger{d: d}
	t.Self = t
	t.initControl(ui, func() { d.SetExpanded(!d.expanded) }, unison.KeySpace, unison.KeyReturn, unison.KeyNumPadEnter)
	t.ringRadius = func() float32 { return 0 }
	t.DrawCallback = t.draw
	d.trigger = t
	d.body = unison.NewPanel()
	d.body.SetLayout(&unison.FlexLayout{Columns: 1})
	for _, p := range content {
		if p.AsPanel().LayoutData() == nil {
			p.AsPanel().SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		}
		d.body.AddChild(p)
	}
	d.AddChild(t)
	d.AddChild(d.body)
	d.SetLayout(disclosureLayout{d})
	return d
}

// Expanded reports whether the section is open.
func (d *Disclosure) Expanded() bool { return d.expanded }

// SetExpanded opens or closes the section, and closes the others in its
// group when it opens.
func (d *Disclosure) SetExpanded(expanded bool) {
	if d.expanded == expanded {
		return
	}
	d.expanded = expanded
	if expanded && d.Group != "" {
		if w := d.Window(); w != nil {
			closeGroup(w.Content(), d)
		}
	}
	d.animate()
	if d.OnToggle != nil {
		d.OnToggle(expanded)
	}
}

// closeGroup closes every other open section in d's group under p.
func closeGroup(p *unison.Panel, d *Disclosure) {
	for _, c := range p.Children() {
		if o, ok := c.Self.(*Disclosure); ok && o != d && o.Group == d.Group && o.expanded {
			o.SetExpanded(false)
		}
		closeGroup(c, d)
	}
}

// animate moves the drawn openness towards the state over 140 ms, easing
// out; with motion reduced it arrives at once.
func (d *Disclosure) animate() {
	target := float32(0)
	if d.expanded {
		target = 1
	}
	scale := d.ui.Theme.MotionScale()
	if scale == 0 || d.Window() == nil {
		d.open = target
		d.MarkForLayoutAndRedraw()
		d.relayoutWindow()
		return
	}
	if d.moving {
		return
	}
	d.moving = true
	from, start := d.open, time.Now()
	var step func()
	step = func() {
		to := float32(0)
		if d.expanded {
			to = 1
		}
		if to != target {
			target, from, start = to, d.open, time.Now()
		}
		t := min(1, float64(time.Since(start))/(140*scale*float64(time.Millisecond)))
		d.open = from + (target-from)*float32(1-math.Pow(1-t, 3))
		if t >= 1 {
			d.open, d.moving = target, false
		}
		d.relayoutWindow()
		if d.moving {
			unison.InvokeTaskAfter(step, 16*time.Millisecond)
		}
	}
	unison.InvokeTaskAfter(step, 16*time.Millisecond)
}

// relayoutWindow lays the window out again, since the section's height moves
// what is below it.
func (d *Disclosure) relayoutWindow() {
	if w := d.Window(); w != nil {
		w.Content().MarkForLayoutRecursively()
		w.MarkForRedraw()
	}
}

// disclosureLayout is the trigger a compact row tall, and under it, indented
// by the loose space, as much of the body as the section is open.
type disclosureLayout struct{ d *Disclosure }

func (l disclosureLayout) bodyHeight(width float32) float32 {
	_, p, _ := l.d.body.Sizes(geom.NewSize(max(0, width-float32(l.d.ui.Interface.SpaceLoose())), 0))
	return p.Height
}

func (l disclosureLayout) LayoutSizes(_ *unison.Panel, hint geom.Size) (minSize, prefSize, maxSize geom.Size) {
	m := l.d.ui.Interface
	w := float32(m.Px(400))
	if hint.Width > 0 {
		w = hint.Width
	}
	h := float32(m.RowHeightCompact()) + l.d.open*(l.bodyHeight(w)+float32(m.Space()))
	return geom.NewSize(0, h), geom.NewSize(float32(m.Px(400)), h), geom.NewSize(unison.DefaultMaxSize, h)
}

func (l disclosureLayout) PerformLayout(target *unison.Panel) {
	d, m := l.d, l.d.ui.Interface
	r := target.ContentRect(false)
	th := float32(m.RowHeightCompact())
	d.trigger.SetFrameRect(geom.NewRect(r.X, r.Y, r.Width, th))
	indent := float32(m.SpaceLoose())
	bh := l.bodyHeight(r.Width)
	d.body.Hidden = d.open == 0
	d.body.SetFrameRect(geom.NewRect(r.X+indent, r.Y+th+float32(m.SpaceSnug()), max(0, r.Width-indent), bh*d.open))
}

func (t *disclosureTrigger) draw(gc *unison.Canvas, _ geom.Rect) {
	d, ui := t.d, t.ui
	tk, m := ui.Theme.Tokens(), ui.Interface
	b := t.ContentRect(false)
	if t.hovered {
		painterFor(gc, ui).fill(b, tk.HoverTint)
	}
	s := float32(m.IconSizeSmall())
	if g, ok := icons.Glyph("chevron-right"); ok {
		gc.Save()
		gc.Translate(geom.NewPoint(b.X+s/2, b.Y+b.Height/2))
		gc.Rotate(90 * d.open)
		drawGlyph(gc, ui, g, s, tk.TextMuted, geom.NewRect(-s/2, -s/2, s, s))
		gc.Restore()
	}
	x := b.X + s + float32(m.SpaceNear())
	weight := text.Regular
	if d.expanded {
		weight = text.Bold
	}
	title := ui.Fonts.Layout([]text.Span{{Text: d.Title, Style: ui.Chrome(ui.Size(RoleBody), weight, tk.TextPrimary)}}, text.Options{})
	tw, th := title.Size()
	title.Draw(gc, x, b.Y+(b.Height-th)/2)
	if d.Count >= 0 {
		st := ui.Chrome(ui.Size(RoleCaption), text.Regular, tk.TextFaint)
		st.Tabular = true
		count := ui.Fonts.Layout([]text.Span{{Text: ui.Number(d.Count), Style: st}}, text.Options{})
		_, ch := count.Size()
		count.Draw(gc, x+tw+float32(m.SpaceNear()), b.Y+(b.Height-ch)/2)
	}
}

// ProvideAccessibility describes the trigger as a button that says whether
// its section is open, which the chevron says to everybody else.
func (t *disclosureTrigger) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	n := b.Node()
	if n.Role == role.Auto {
		n.Role = role.Button
	}
	n.Name = t.d.Title
	n.Expandable, n.Expanded = true, t.d.expanded
	n.Description = "Collapsed"
	if t.d.expanded {
		n.Description = "Expanded"
	}
	if t.Enabled() {
		n.Actions = n.Actions.With(accessibility.Press)
	}
}

// PerformAccessibilityAction opens or closes the section for a screen reader.
func (t *disclosureTrigger) PerformAccessibilityAction(req accessibility.ActionRequest) bool {
	if req.Action != accessibility.Press || !t.Enabled() {
		return false
	}
	t.fire()
	return true
}
