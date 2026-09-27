package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
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

func TestAButtonSaysWhatItDoesInEveryState(t *testing.T) {
	var save, busy, mode, archive *kvitui.Button
	pressed := 0
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		save = kvitui.NewButton(ui, "Save")
		save.Form = kvitui.ButtonPrimary
		save.OnClick = func() { pressed++ }
		busy = kvitui.NewButton(ui, "Save")
		busy.Busy = true
		mode = kvitui.NewButton(ui, "Select region")
		mode.Checkable = true
		archive = kvitui.NewButton(ui, "Archive")
		archive.SetEnabled(false)
		archive.Explanation = "The branch has work that has not been pushed."
		return []unison.Paneler{save, busy, mode, archive}
	})
	screen.Do(func() {
		if h := save.FrameRect().Height; h != float32(ui.Interface.ControlHeight()) {
			t.Errorf("a button is %.1f tall, want %d", h, ui.Interface.ControlHeight())
		}
	})
	screen.Click(screen.PanelCenter(save))
	screen.Do(func() { save.Focus() })
	screen.KeyPress(unison.KeySpace, 0)
	if pressed != 2 {
		t.Errorf("a click and Space made %d presses", pressed)
	}
	// A busy button stays enabled and says what it is doing.
	if n := screen.AccessibilityNodeFor(busy); n == nil || n.Name != "Working…" || n.Disabled {
		t.Errorf("the busy button's node: %+v", n)
	}
	// A disabled button still gives its reason.
	if n := screen.AccessibilityNodeFor(archive); n == nil || !n.Disabled || n.Description != archive.Explanation {
		t.Errorf("the disabled button's node: %+v", n)
	}
	screen.Click(screen.PanelCenter(mode))
	if n := screen.AccessibilityNodeFor(mode); n == nil || !n.HasCheck || n.Checked != check.On {
		t.Errorf("a checkable button pressed once: %+v", n)
	}
}
