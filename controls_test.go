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

func TestAFieldTakesTypingAndSaysWhatIsWrong(t *testing.T) {
	var amount *kvitui.Field
	var typed []string
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		amount = kvitui.NewField(ui)
		amount.Label, amount.Placeholder = "Amount", "0.00"
		amount.OnChange = func(s string) { typed = append(typed, s) }
		// Padded, so the ring outside the field is inside the window.
		around := kvitui.Width(ui, kvitui.Px(240), amount)
		around.SetBorder(kvitui.Padding(ui, kvitui.SizeViewMargin))
		return []unison.Paneler{around}
	})
	screen.Click(screen.PanelCenter(amount.Edit()))
	screen.Type("twelve")
	screen.Sync()
	if got := amount.Text(); got != "twelve" {
		t.Fatalf("typing made %q", got)
	}
	if len(typed) == 0 || typed[len(typed)-1] != "twelve" {
		t.Errorf("OnChange saw %v", typed)
	}
	n := screen.AccessibilityNodeFor(amount.Edit())
	if n == nil || n.Role != role.TextField || n.Name != "Amount" || n.Value != "twelve" || n.Placeholder != "0.00" {
		t.Fatalf("the field's node: %+v", n)
	}
	// The ring shows whenever the field holds the focus, a click included.
	var box geom.Rect
	screen.Do(func() { box = amount.Edit().RectToRoot(amount.Edit().ContentRect(true)) })
	img := screen.Capture()
	want := kvitui.Color(ui.Theme.Tokens().FocusRing)
	ringed := false
	for dx := 1; dx <= ui.Interface.FocusRingWidth(); dx++ {
		pr, pg, pb, _ := img.At(int(box.X)-dx, int(box.Y+box.Height/2)).RGBA()
		ringed = ringed || (pr>>8 == uint32(want.Red()) && pg>>8 == uint32(want.Green()) && pb>>8 == uint32(want.Blue()))
	}
	if !ringed {
		t.Error("a focused field has no ring")
	}
	// An error is a message under the field as well as a red outline.
	var before, after float32
	screen.Do(func() {
		before = amount.FrameRect().Height
		amount.Error = "Not a number"
		amount.MarkForLayoutAndRedraw()
		amount.Window().Content().MarkForLayoutRecursively()
		amount.Window().ValidateLayout()
		after = amount.FrameRect().Height
	})
	if before != float32(ui.Interface.ControlHeight()) || after <= before {
		t.Errorf("the field is %.1f tall, and %.1f with an error", before, after)
	}
	if n := screen.AccessibilityNodeFor(amount.Edit()); n == nil || n.Description != "Not a number" || !n.Invalid {
		t.Errorf("the field in error: %+v", n)
	}
}

func TestASearchFieldClearsWithEscapeAndCountsWhatItLeft(t *testing.T) {
	var s *kvitui.SearchField
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.Locale = language.AmericanEnglish
		s = kvitui.NewSearchField(ui)
		s.Matches, s.MatchedNoun, s.MatchedNounPlural = 250000, "entry", "entries"
		return []unison.Paneler{s}
	})
	n := screen.AccessibilityNodeFor(s.Edit())
	if n == nil || n.Name != "Filter" || n.Description != "250,000 entries" {
		t.Fatalf("the search field's node: %+v", n)
	}
	clearShown := func() (shown bool) {
		screen.Do(func() { shown = !s.Children()[0].Hidden })
		return shown
	}
	if clearShown() {
		t.Error("an empty filter shows its clear button")
	}
	screen.Click(screen.PanelCenter(s.Edit()))
	screen.Type("harlow")
	screen.Sync()
	if !clearShown() {
		t.Error("a filter with text has no clear button")
	}
	screen.KeyPress(unison.KeyEscape, 0)
	if got := s.Text(); got != "" {
		t.Errorf("Escape left %q", got)
	}
}
