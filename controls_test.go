package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/role"
	"golang.org/x/text/language"
)

func TestATabSaysWhatItHoldsAndWhetherItIsChosen(t *testing.T) {
	var all, plain *kvitui.Tab
	pressed := 0
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.Locale = language.AmericanEnglish
		all = kvitui.NewTab(ui, "All")
		all.Selected, all.Count = true, 1284
		all.Explanation = "Every transaction, whatever its state."
		plain = kvitui.NewTab(ui, "Disputed")
		plain.OnClick = func() { pressed++ }
		return []unison.Paneler{all, plain}
	})
	screen.Do(func() {
		if h := all.FrameRect().Height; h != float32(ui.Interface.TabHeight()) {
			t.Errorf("a tab is %.1f tall, want %d", h, ui.Interface.TabHeight())
		}
		// The count sits a near space after the name, and the tab grows by
		// both.
		_, counted, _ := all.Sizes(geom.Size{})
		all.Count = -1
		_, bare, _ := all.Sizes(geom.Size{})
		all.Count = 1284
		if counted.Width < bare.Width+float32(ui.Interface.SpaceNear()) {
			t.Errorf("the count adds %.1f to the tab's width", counted.Width-bare.Width)
		}
	})
	n := screen.AccessibilityNodeFor(all)
	if n == nil || n.Role != role.Tab || n.Name != "All, 1,284 items" || !n.Selected || n.Description != all.Explanation {
		t.Errorf("the selected tab's node: %+v", n)
	}
	pn := screen.AccessibilityNodeFor(plain)
	if pn == nil || pn.Name != "Disputed" || pn.Selected {
		t.Fatalf("the plain tab's node: %+v", pn)
	}
	screen.Click(screen.PanelCenter(plain))
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: pn.ID, Action: accessibility.Press})
	screen.Sync()
	if pressed != 2 {
		t.Errorf("a click and a screen reader's press made %d presses", pressed)
	}
}
