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
