package kvitui_test

import (
	"runtime"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/mod"
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

func TestASegmentedControlChoosesOne(t *testing.T) {
	var s *kvitui.Segmented
	var chosen []string
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		s = kvitui.NewSegmented(ui, "Period",
			kvitui.Option{Value: "week", Label: "Week"}, kvitui.Option{Value: "month", Label: "Month"})
		s.OnChoose = func(v string) { chosen = append(chosen, v) }
		return []unison.Paneler{s}
	})
	if n := screen.AccessibilityNodeFor(s); n == nil || n.Role != role.TabList || n.Name != "Period" {
		t.Errorf("the strip's node: %+v", n)
	}
	month := s.Children()[1]
	screen.Click(screen.PanelCenter(month))
	if s.Current != "month" || len(chosen) != 1 || chosen[0] != "month" {
		t.Errorf("clicking Month left %q chosen and reported %v", s.Current, chosen)
	}
	if n := screen.AccessibilityNodeFor(month); n == nil || n.Role != role.Tab || !n.Selected {
		t.Errorf("the chosen segment's node: %+v", n)
	}
}

func TestAStepperStaysInItsRange(t *testing.T) {
	var s *kvitui.Stepper
	var changes []int
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		s = kvitui.NewStepper(ui, "Interface size", 10, 12)
		s.Value, s.Unit = 11, "px"
		s.OnChange = func(v int) { changes = append(changes, v) }
		return []unison.Paneler{s}
	})
	n := screen.AccessibilityNodeFor(s)
	if n == nil || n.Role != role.SpinButton || n.Number != 11 || n.Min != 10 || n.Max != 12 || n.Value != "11 px" {
		t.Fatalf("the stepper's node: %+v", n)
	}
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: n.ID, Action: accessibility.Increment})
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: n.ID, Action: accessibility.Increment})
	screen.Sync()
	if s.Value != 12 || len(changes) != 1 {
		t.Errorf("two increments from 11 to a top of 12 left %d and reported %v", s.Value, changes)
	}
	plus := s.Children()[2]
	screen.Do(func() {
		if plus.Enabled() {
			t.Error("the plus is enabled at the top of the range")
		}
	})
}

func TestACheckHasThreeStates(t *testing.T) {
	var c *kvitui.Check
	var seen []bool
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		c = kvitui.NewCheck(ui, "Some of these")
		c.Partial = true
		c.OnChange = func(on bool) { seen = append(seen, on) }
		return []unison.Paneler{c}
	})
	if n := screen.AccessibilityNodeFor(c); n == nil || n.Role != role.CheckBox || n.Checked != check.Mixed {
		t.Fatalf("a partial check's node: %+v", n)
	}
	// Pressing a partial box checks the whole of it.
	screen.Click(screen.PanelCenter(c))
	screen.Click(screen.PanelCenter(c))
	if len(seen) != 2 || !seen[0] || seen[1] {
		t.Errorf("two presses from partial reported %v", seen)
	}
}

func TestAnUnavailableChipStaysReachableAndDoesNothing(t *testing.T) {
	var behind, filter *kvitui.ChipButton
	pressed := 0
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		behind = kvitui.NewChipButton(ui, "2 behind")
		behind.UnavailableReason = "The remote has not been fetched yet."
		behind.OnActivate = func() { pressed++ }
		filter = kvitui.NewChipButton(ui, "Uncategorised")
		filter.Selectable, filter.Selected = true, true
		return []unison.Paneler{behind, filter}
	})
	screen.Click(screen.PanelCenter(behind))
	var focusable bool
	screen.Do(func() { focusable = behind.Focusable() })
	if pressed != 0 || !focusable {
		t.Errorf("an unavailable chip was pressed %d times, focusable %v", pressed, focusable)
	}
	if n := screen.AccessibilityNodeFor(behind); n == nil || n.Description != behind.UnavailableReason {
		t.Errorf("the unavailable chip's node: %+v", n)
	}
	n := screen.AccessibilityNodeFor(filter)
	if n == nil || n.Role != role.ToggleButton || !n.Pressed || n.Checked != check.On {
		t.Errorf("a quick filter that is on: %+v", n)
	}
	screen.Do(func() {
		_, p, _ := filter.Sizes(geom.Size{})
		if p.Height < 28 {
			t.Errorf("a quick filter is %.1f tall, under a control's height", p.Height)
		}
	})
}

func TestASelectStepsWithTheArrowsAndOpensAMenu(t *testing.T) {
	var s *kvitui.Select
	var chosen []string
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		s = kvitui.NewSelect(ui, "Currency",
			kvitui.Option{Value: "GBP", Label: "GBP"}, kvitui.Option{Value: "EUR", Label: "EUR"}, kvitui.Option{Value: "USD", Label: "USD"})
		s.OnChoose = func(v string) { chosen = append(chosen, v) }
		return []unison.Paneler{s}
	})
	screen.Do(func() { s.Focus() })
	screen.KeyPress(unison.KeyDown, 0)
	if s.Current != "EUR" {
		t.Errorf("Down chose %q", s.Current)
	}
	if n := screen.AccessibilityNodeFor(s); n == nil || n.Role != role.ComboBox || n.Name != "Currency" || n.Value != "EUR" {
		t.Errorf("the select's node: %+v", n)
	}
	screen.KeyPress(unison.KeySpace, 0)
	screen.Sync()
	var w *unison.Window
	screen.Do(func() { w = s.Window() })
	var usd *accessibility.Node
	for _, n := range screen.AccessibilityTree(w).Nodes {
		if n.Role == role.MenuItem && n.Name == "USD" {
			usd = n
		}
	}
	if usd == nil {
		t.Fatal("Space opened no menu with the options")
	}
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: usd.ID, Action: accessibility.Press})
	screen.Sync()
	if s.Current != "USD" || len(chosen) != 2 {
		t.Errorf("choosing USD from the menu left %q, reported %v", s.Current, chosen)
	}
}

func TestAReadOnlyTextAreaTakesNoTyping(t *testing.T) {
	var a *kvitui.TextArea
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		a = kvitui.NewTextArea(ui)
		a.Label, a.ReadOnly = "Source", true
		a.SetText("one\ntwo")
		return []unison.Paneler{a}
	})
	screen.Click(screen.PanelCenter(a.Edit()))
	screen.Type("x")
	screen.KeyPress(unison.KeyBackspace, 0)
	screen.KeyPress(unison.KeyReturn, 0)
	screen.Sync()
	if got := a.Text(); got != "one\ntwo" {
		t.Errorf("a read-only area was changed to %q", got)
	}
	if n := screen.AccessibilityNodeFor(a.Edit()); n == nil || n.Role != role.TextArea || !n.ReadOnly || n.Name != "Source" {
		t.Errorf("the text area's node: %+v", n)
	}
}

// Ctrl+Backspace and Ctrl+Delete take a word in a field and in a text area,
// with OnChange told; plain Backspace takes a character, and a read-only
// field takes nothing.
func TestCtrlBackspaceAndDeleteTakeAWordInAField(t *testing.T) {
	var name, fixed *kvitui.Field
	var area *kvitui.TextArea
	var changes []string
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		name, fixed, area = kvitui.NewField(ui), kvitui.NewField(ui), kvitui.NewTextArea(ui)
		name.SetText("hello big world")
		name.OnChange = func(s string) { changes = append(changes, s) }
		fixed.SetText("stays put")
		fixed.ReadOnly = true
		area.SetText("first line\nsecond")
		return []unison.Paneler{name, fixed, area}
	})
	screen.Click(screen.PanelCenter(name.Edit()))
	screen.Do(func() { name.Edit().SetSelectionToEnd() })
	screen.KeyPress(unison.KeyBackspace, mod.Control)
	screen.Sync()
	if got := name.Text(); got != "hello big " {
		t.Errorf("Ctrl+Backspace left %q", got)
	}
	if len(changes) == 0 || changes[len(changes)-1] != "hello big " {
		t.Errorf("OnChange saw %q", changes)
	}
	screen.KeyPress(unison.KeyBackspace, 0)
	screen.Do(func() { name.Edit().SetSelectionTo(0) })
	screen.KeyPress(unison.KeyDelete, mod.Control)
	screen.Sync()
	if got := name.Text(); got != " big" {
		t.Errorf("Backspace then Ctrl+Delete left %q", got)
	}

	screen.Click(screen.PanelCenter(fixed.Edit()))
	screen.KeyPress(unison.KeyBackspace, mod.Control)
	screen.Sync()
	if got := fixed.Text(); got != "stays put" {
		t.Errorf("a read-only field was changed to %q", got)
	}

	screen.Click(screen.PanelCenter(area.Edit()))
	screen.Do(func() { area.Edit().SetSelectionTo(len("first line\n")) })
	screen.KeyPress(unison.KeyBackspace, mod.Control)
	screen.KeyPress(unison.KeyBackspace, mod.Control)
	screen.Sync()
	if got := area.Text(); got != "first second" {
		t.Errorf("Ctrl+Backspace twice from a line's start left %q", got)
	}
}

// On Windows and Linux, Ctrl+Left and Ctrl+Right move by word in a field,
// with Shift selecting, rather than to the line's start and end.
func TestCtrlArrowsMoveByWordInAField(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("macOS keeps unison's own arrow keys")
	}
	var f *kvitui.Field
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		f = kvitui.NewField(ui)
		f.SetText("one two three")
		return []unison.Paneler{f}
	})
	screen.Click(screen.PanelCenter(f.Edit()))
	selection := func() (start, end int) {
		screen.Do(func() { start, end = f.Edit().Selection() })
		return start, end
	}
	screen.Do(func() { f.Edit().SetSelectionToEnd() })
	screen.KeyPress(unison.KeyLeft, mod.Control)
	if start, end := selection(); start != 8 || end != 8 {
		t.Errorf("Ctrl+Left went to %d-%d, want the start of \"three\"", start, end)
	}
	screen.KeyPress(unison.KeyLeft, mod.Control|mod.Shift)
	if start, end := selection(); start != 4 || end != 8 {
		t.Errorf("Ctrl+Shift+Left selected %d-%d, want \"two \"", start, end)
	}
	screen.Do(func() { f.Edit().SetSelectionTo(0) })
	screen.KeyPress(unison.KeyRight, mod.Control)
	if start, end := selection(); start != 3 || end != 3 {
		t.Errorf("Ctrl+Right went to %d-%d, want the end of \"one\"", start, end)
	}
}
