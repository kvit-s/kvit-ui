package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
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
