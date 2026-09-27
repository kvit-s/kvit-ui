package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
	"golang.org/x/text/language"
)

func TestADividerIsOneHairline(t *testing.T) {
	var across, down *kvitui.Divider
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		across = kvitui.NewDivider(ui)
		down = kvitui.NewDivider(ui)
		down.Vertical = true
		return []unison.Paneler{across, down}
	})
	screen.Do(func() {
		h := float32(ui.Interface.Hairline())
		if _, p, _ := across.Sizes(geom.Size{}); p.Height != h {
			t.Errorf("a horizontal rule is %.1f tall, want %.1f", p.Height, h)
		}
		if _, p, _ := down.Sizes(geom.Size{}); p.Width != h {
			t.Errorf("a vertical rule is %.1f wide, want %.1f", p.Width, h)
		}
	})
	if n := screen.AccessibilityNodeFor(across); n == nil || n.Role != role.Separator {
		t.Errorf("a rule's node: %+v", n)
	}
}

// A rule on an edge is drawn inside the panel, in the border colour, and an
// edge without one is the panel's ground.
func TestAPanelDrawsTheRulesItIsAskedFor(t *testing.T) {
	var p *kvitui.Panel
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		p = kvitui.NewPanel(ui)
		p.RuleTop = true
		p.SetLayout(kvitui.AtLeast(ui, kvitui.Px(60), &unison.FlexLayout{Columns: 1}))
		p.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(200, 0)})
		return []unison.Paneler{p}
	})
	var r geom.Rect
	screen.Do(func() { r = p.RectToRoot(p.ContentRect(true)) })
	img := screen.Capture()
	tk := ui.Theme.Tokens()
	is := func(x, y int, want unison.Color) bool {
		pr, pg, pb, _ := img.At(x, y).RGBA()
		return uint32(want.Red()) == pr>>8 && uint32(want.Green()) == pg>>8 && uint32(want.Blue()) == pb>>8
	}
	mid := int(r.X + r.Width/2)
	if !is(mid, int(r.Y), kvitui.Color(tk.Border)) {
		t.Error("the top edge has no rule")
	}
	if !is(mid, int(r.Bottom())-1, kvitui.Color(tk.PanelBackground)) {
		t.Error("the bottom edge is not the panel's ground")
	}
}

func TestAListRowActsOnlyWhenItSaysItDoes(t *testing.T) {
	var opens, layout *kvitui.ListRow
	var inside *kvitui.IconButton
	opened, pressed := 0, 0
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		inside = kvitui.NewIconButton(ui, "pencil", "Edit")
		inside.OnClick = func() { pressed++ }
		opens = kvitui.NewListRow(ui, inside)
		opens.Interactive, opens.Label, opens.OpensLabel = true, "Groceries", "Opens the transaction"
		opens.OnActivate = func() { opened++ }
		layout = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "layout only"))
		layout.Label = "Amount"
		for _, r := range []*kvitui.ListRow{opens, layout} {
			r.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(300, 0)})
		}
		return []unison.Paneler{opens, layout}
	})
	n := screen.AccessibilityNodeFor(opens)
	if n == nil || n.Role != role.ListItem || n.Name != "Groceries" || n.Description != "Opens the transaction" {
		t.Fatalf("a row that opens something: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(layout); n == nil || n.Role != role.Label || n.Name != "Amount" {
		t.Errorf("a row that only lays things out: %+v", n)
	}
	// A click on a row that opens nothing does nothing; on one that opens
	// something it opens it once.
	screen.Click(screen.PanelPoint(layout, geom.NewPoint(5, 5)))
	screen.Click(screen.PanelPoint(opens, geom.NewPoint(250, 10)))
	if opened != 1 {
		t.Errorf("clicks opened the row %d times", opened)
	}
	// Space on the button inside presses the button and leaves the row shut.
	screen.Do(func() { inside.RequestFocus() })
	screen.KeyPress(unison.KeySpace, 0)
	if pressed != 1 || opened != 1 {
		t.Errorf("Space on the button inside: pressed %d, opened %d", pressed, opened)
	}
	screen.Do(func() { opens.RequestFocus() })
	screen.KeyPress(unison.KeyReturn, 0)
	if opened != 2 {
		t.Errorf("Return on the focused row opened it %d times in all", opened)
	}
	var focusable bool
	screen.Do(func() { focusable = layout.Focusable() })
	if focusable {
		t.Error("a row that does nothing takes the focus")
	}
}

func TestASlimRowSaysItsPartsInOrder(t *testing.T) {
	var row *kvitui.SlimRow
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		row = kvitui.NewSlimRow(ui, "kvit-cash")
		row.Kind, row.Phrase, row.Measured = "application", "not started", false
		row.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(500, 0)})
		return []unison.Paneler{row}
	})
	if n := screen.AccessibilityNodeFor(row); n == nil || n.Name != "kvit-cash, application, not started" {
		t.Errorf("the slim row's node: %+v", n)
	}
	// An unmeasured figure is still there to be read, as "not measured".
	var figure *kvitui.Figure
	screen.Do(func() {
		var find func(p *unison.Panel)
		find = func(p *unison.Panel) {
			for _, c := range p.Children() {
				if f, ok := c.Self.(*kvitui.Figure); ok {
					figure = f
				}
				find(c)
			}
		}
		find(row.AsPanel())
	})
	if figure == nil {
		t.Fatal("the slim row has no figure")
	}
	if n := screen.AccessibilityNodeFor(figure); n == nil || n.Name != "not measured" {
		t.Errorf("the unmeasured figure's node: %+v", n)
	}
}

func TestASectionHeadingOpensAndItsActionStaysItsOwn(t *testing.T) {
	var h, fixed *kvitui.SectionHeading
	var toggles []bool
	actions := 0
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.Locale = language.AmericanEnglish
		h = kvitui.NewSectionHeading(ui, "Waiting on me")
		h.Count, h.Counted, h.Action, h.Collapsible = 4, "project", "Hand all to an agent", true
		h.OnToggle = func(open bool) { toggles = append(toggles, open) }
		h.OnAction = func() { actions++ }
		fixed = kvitui.NewSectionHeading(ui, "Waiting on Sam")
		fixed.Count, fixed.Counted = 1200, "account"
		for _, x := range []*kvitui.SectionHeading{h, fixed} {
			x.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(380, 0)})
		}
		return []unison.Paneler{h, fixed}
	})
	if got := fixed.CountPhrase(); got != "1,200 accounts" {
		t.Errorf("the count reads %q", got)
	}
	screen.Do(func() {
		if hgt := h.FrameRect().Height; hgt != float32(ui.Interface.RowHeightCompact()) {
			t.Errorf("a heading is %.1f tall", hgt)
		}
		if fixed.Focusable() || !h.Focusable() {
			t.Error("only the heading that opens should take the focus")
		}
	})
	n := screen.AccessibilityNodeFor(h)
	if n == nil || n.Role != role.Button || !n.Expandable || !n.Expanded || n.Description != "Expanded" {
		t.Fatalf("the heading's node: %+v", n)
	}
	// Return on the heading closes the group; Space on its action runs the
	// action and leaves the group alone.
	screen.Do(func() { h.Focus() })
	screen.KeyPress(unison.KeyReturn, 0)
	if len(toggles) != 1 || toggles[0] || h.Expanded {
		t.Errorf("Return on the heading toggled %v", toggles)
	}
	var link *kvitui.Link
	screen.Do(func() {
		for _, c := range h.Children() {
			if l, ok := c.Self.(*kvitui.Link); ok {
				link = l
			}
		}
	})
	if ln := screen.AccessibilityNodeFor(link); ln == nil || ln.Role != role.Link || ln.Name != "Hand all to an agent" {
		t.Fatalf("the action's node under the heading: %+v", ln)
	}
	screen.Do(func() { link.Focus() })
	screen.KeyPress(unison.KeySpace, 0)
	if actions != 1 || len(toggles) != 1 {
		t.Errorf("Space on the action: %d actions, toggles %v", actions, toggles)
	}
}

func TestACardOpensOnceWhenARowInsideActs(t *testing.T) {
	var card *kvitui.Card
	var row *kvitui.ListRow
	var inert *kvitui.Label
	cards, rows := 0, 0
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		row = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Groceries"))
		row.Interactive = true
		row.OnActivate = func() { rows++ }
		inert = kvitui.NewLabel(ui, "Spending this month")
		card = kvitui.NewCard(ui, inert, row)
		card.Interactive, card.OpensLabel = true, "Opens the month's transactions"
		card.OnActivate = func() { cards++ }
		card.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(300, 0)})
		return []unison.Paneler{card}
	})
	screen.Click(screen.PanelCenter(inert))
	screen.Click(screen.PanelCenter(row))
	if cards != 1 || rows != 1 {
		t.Errorf("a press on the card's text and one on its row opened the card %d times and the row %d", cards, rows)
	}
	if n := screen.AccessibilityNodeFor(card); n == nil || n.Description != "Opens the month's transactions" {
		t.Errorf("the card's node: %+v", n)
	}
}

func TestDisclosuresInAGroupOpenOneAtATime(t *testing.T) {
	var a, b *kvitui.Disclosure
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.Theme.SetReducedMotion(true)
		a = kvitui.NewDisclosure(ui, "First", kvitui.NewLabel(ui, "one"))
		b = kvitui.NewDisclosure(ui, "Second", kvitui.NewLabel(ui, "two"))
		a.Group, b.Group = "faq", "faq"
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(300), kvitui.Column(ui, kvitui.SizeSpace, a, b))}
	})
	var closed, opened float32
	screen.Do(func() { closed = a.FrameRect().Height })
	trigger := func(d *kvitui.Disclosure) *unison.Panel { return d.Children()[0] }
	screen.Click(screen.PanelCenter(trigger(a)))
	screen.Do(func() { opened = a.FrameRect().Height })
	if !a.Expanded() || opened <= closed {
		t.Errorf("clicking the first opened it to %.1f from %.1f", opened, closed)
	}
	screen.Click(screen.PanelCenter(trigger(b)))
	if a.Expanded() || !b.Expanded() {
		t.Errorf("opening the second left the first open (%v) or the second shut (%v)", a.Expanded(), b.Expanded())
	}
	if n := screen.AccessibilityNodeFor(trigger(b)); n == nil || n.Role != role.Button || !n.Expanded || n.Description != "Expanded" {
		t.Errorf("the open trigger's node: %+v", n)
	}
}

func TestAnEmptyStateSaysWhatWouldBeHereAndOffersTheWayIn(t *testing.T) {
	var full, compact *kvitui.EmptyState
	actions := 0
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		full = kvitui.NewEmptyState(ui, "No transactions yet")
		full.Detail, full.Action = "Import a statement.", "Import a statement"
		full.OnAction = func() { actions++ }
		compact = kvitui.NewEmptyState(ui, "No agents")
		compact.Form, compact.Action = kvitui.EmptyCompact, "Start one"
		for _, e := range []*kvitui.EmptyState{full, compact} {
			e.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(360, 0)})
		}
		return []unison.Paneler{full, compact}
	})
	if n := screen.AccessibilityNodeFor(full); n == nil || n.Name != "No transactions yet. Import a statement." {
		t.Errorf("the empty state's node: %+v", n)
	}
	var button *kvitui.Button
	screen.Do(func() {
		if h := compact.FrameRect().Height; h != float32(ui.Interface.RowHeightSlim()) {
			t.Errorf("the compact form is %.1f tall, want a slim row", h)
		}
		for _, c := range full.Children() {
			if b, ok := c.Self.(*kvitui.Button); ok && !b.Hidden {
				button = b
			}
		}
	})
	if button == nil {
		t.Fatal("the full form has no action button")
	}
	// Reachable without a pointer: the only way in for a keyboard reader.
	screen.Do(func() { button.Focus() })
	screen.KeyPress(unison.KeySpace, 0)
	if actions != 1 {
		t.Errorf("Space on the action ran it %d times", actions)
	}
}

func TestAPaneSlidesOverTheRightEdge(t *testing.T) {
	var pane *kvitui.Pane
	var area *unison.Panel
	closes := 0
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.Theme.SetReducedMotion(true)
		pane = kvitui.NewPane(ui, "Transaction", kvitui.NewLabel(ui, "Detail goes here"))
		pane.OnClose = func() { closes++; pane.SetOpen(false) }
		area = kvitui.WithPane(ui, kvitui.NewLabel(ui, "the list"), pane)
		area.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(700, 200)})
		return []unison.Paneler{area}
	})
	screen.Do(func() {
		if got, want := pane.FrameRect().X, area.FrameRect().Width-float32(ui.Interface.PaneWidth()); got != want {
			t.Errorf("the open pane starts at %.1f, want %.1f", got, want)
		}
	})
	close := pane.Children()[0]
	if n := screen.AccessibilityNodeFor(close); n == nil || n.Name != "Close the pane" {
		t.Errorf("the close control's node: %+v", n)
	}
	screen.Click(screen.PanelCenter(close))
	screen.Sync()
	screen.Do(func() {
		if closes != 1 || !pane.Hidden {
			t.Errorf("the close control closed %d times, and the pane is hidden: %v", closes, pane.Hidden)
		}
	})
}
