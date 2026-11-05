package kvitui_test

import (
	"image"
	"image/color"
	"runtime"
	"slices"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/palette"
	"github.com/kvit-s/kvit-ui/tokens"
	"github.com/kvit-s/kvit-ui/uitest"
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
		busy.Busy, busy.BusyText = true, "Saving…"
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
	if n := screen.AccessibilityNodeFor(busy); n == nil || n.Name != "Saving…" || n.Disabled {
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

func TestACompactSelectShowsItsValueAndOpensTheSameList(t *testing.T) {
	var s *kvitui.Select
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		s = kvitui.NewSelect(ui, "Period",
			kvitui.Option{Value: "all", Label: "All history"},
			kvitui.Option{Value: "month", Label: "Month"})
		s.Symbol, s.Compact = "calendar", true
		return []unison.Paneler{s}
	})

	// Its natural width holds the symbol and the value, not a field's width.
	screen.Do(func() {
		_, p, _ := s.Sizes(geom.Size{})
		_, field, _ := kvitui.NewSelect(ui, "Period",
			kvitui.Option{Value: "all", Label: "All history"}).Sizes(geom.Size{})
		if p.Width >= field.Width {
			t.Errorf("a compact select is %.1f wide, a field %.1f", p.Width, field.Width)
		}
		if h := s.FrameRect().Height; h < 22 {
			t.Errorf("a compact select is %.1f tall", h)
		}
	})
	screen.Do(func() { s.Focus() })
	screen.KeyPress(unison.KeyDown, 0)
	if s.Current != "month" {
		t.Errorf("Down chose %q", s.Current)
	}
	if n := screen.AccessibilityNodeFor(s); n == nil || n.Role != role.ComboBox || n.Name != "Period" || n.Value != "Month" {
		t.Errorf("the compact select's node: %+v", n)
	}
	screen.KeyPress(unison.KeySpace, 0)
	screen.Sync()
	var w *unison.Window
	screen.Do(func() { w = s.Window() })
	found := false
	for _, n := range screen.AccessibilityTree(w).Nodes {
		if n.Role == role.MenuItem && n.Name == "All history" {
			found = true
		}
	}
	if !found {
		t.Error("Space opened no menu with the options")
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

// plainRows is a list of rows each ending in a plain select, the words of the
// first long enough to need cutting at the cell's width.
func plainRows(ui *kvitui.UI, cell kvitui.Measure, long string) (rows *unison.Panel, selects []*kvitui.Select) {
	options := []kvitui.Option{
		{Value: "project", Label: "Project"},
		{Value: "none", Label: "Not a project"},
		{Value: "long", Label: long},
	}
	var lines []unison.Paneler
	for _, name := range []string{"photo-sync", "recipes", "old-backups"} {
		s := kvitui.NewSelect(ui, name, options...)
		s.Plain = true
		row := kvitui.NewListRow(ui, kvitui.FullWidth(kvitui.NewLabel(ui, name)), kvitui.Width(ui, cell, s))
		row.Label = name
		lines, selects = append(lines, row), append(selects, s)
	}
	return kvitui.Column(ui, kvitui.Px(0), lines...), selects
}

const plainLongWords = "Part of the household accounts archive and everything kept in it"

// inkColumns are the columns of a panel's box, counted from its left edge,
// where some pixel differs from the column just outside its right edge at the
// same height: where it has drawn anything but the ground around it.
func inkColumns(img image.Image, box geom.Rect) []int {
	var out []int
	ground := func(y int) color.Color { return img.At(int(box.Right())+1, y) }
	for x := int(box.X); x < int(box.Right()); x++ {
		for y := int(box.Y); y < int(box.Bottom()); y++ {
			if !samePixel(img.At(x, y), ground(y)) {
				out = append(out, x-int(box.X))
				break
			}
		}
	}
	return out
}

// lastRun is the rightmost unbroken run of columns, as its first and last.
func lastRun(columns []int) (first, last int) {
	if len(columns) == 0 {
		return -1, -1
	}
	first, last = columns[len(columns)-1], columns[len(columns)-1]
	for i := len(columns) - 2; i >= 0 && columns[i] == first-1; i-- {
		first = columns[i]
	}
	return first, last
}

// In all four themes a plain select draws no ground and no border at rest, one
// under the pointer, the selection's tint while pressed and while its list is
// open with the pointer gone from it, and nothing again once the list closes.
func TestAPlainSelectDrawsAGroundOnlyWhenThePointerOrItsListIsOnIt(t *testing.T) {
	for _, theme := range tokens.BuiltInThemes() {
		t.Run(theme, func(t *testing.T) {
			var selects []*kvitui.Select
			s := uitest.Open(t, uitest.Options{Theme: theme}, func(ui *kvitui.UI) unison.Paneler {
				var rows *unison.Panel
				rows, selects = plainRows(ui, kvitui.Px(240), plainLongWords)
				return rows
			})
			box := rootRect(selects[0])
			// The far end of the control, past the words and the arrow, and
			// the row's ground outside it, clear of a focus ring.
			end := image.Pt(int(box.Right())-int(box.Width)/8, int(box.Y+box.Height/2))
			outside := image.Pt(int(box.Right())+2*s.UI.Interface.FocusRingWidth()+1, end.Y)
			colours := s.UI.Theme.Tokens()
			// wants checks the control's far end shows ground, or the row's
			// own ground when it is nil.
			wants := func(state string, ground *palette.Color) {
				t.Helper()
				img := s.Capture()
				got, want := img.At(end.X, end.Y), img.At(outside.X, outside.Y)
				if ground != nil {
					want = kvitui.Color(*ground)
				}
				if !samePixel(got, want) {
					t.Errorf("%s: the control's ground is %v, wanted %v", state, got, want)
				}
			}

			away := geom.NewPoint(box.X, box.Bottom()+float32(100))
			s.Screen.MouseMove(away, 0)
			s.Sync()
			wants("at rest", nil)
			// No border either: the whole edge of the box is the ground around it.
			img := s.Capture()
			for x := int(box.X); x < int(box.Right()); x++ {
				if !samePixel(img.At(x, int(box.Y)), img.At(int(box.Right())+1, int(box.Y))) {
					t.Errorf("at rest: the top edge at %d is drawn", x)
					break
				}
			}

			inside := geom.NewPoint(float32(end.X), float32(end.Y))
			s.Screen.MouseMove(inside, 0)
			s.Sync()
			wants("under the pointer", &colours.HoverTint)

			s.Screen.MouseDown(inside, unison.ButtonLeft, 0)
			s.Sync()
			wants("pressed", &colours.SelectionTint)

			s.Screen.MouseUp(inside, unison.ButtonLeft, 0)
			s.Sync()
			s.Screen.MouseMove(away, 0)
			s.Sync()
			if n := len(s.Window.ShownMenus()); n != 1 {
				t.Fatalf("a press opened %d menus", n)
			}
			wants("list open, pointer gone", &colours.SelectionTint)

			s.Screen.KeyPress(unison.KeyEscape, 0)
			s.Sync()
			wants("list closed", nil)
		})
	}
}

// A press anywhere in the width the layout gave opens the list, a press on the
// select with its list open closes it, and a press on another row's select
// goes straight to that row's list, leaving the first row it was on unmarked.
func TestAPressAnywhereInAPlainSelectOpensItsList(t *testing.T) {
	var selects []*kvitui.Select
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		var rows *unison.Panel
		rows, selects = plainRows(ui, kvitui.Px(240), plainLongWords)
		return rows
	})
	box, second := rootRect(selects[0]), selects[1]
	if want := float32(s.UI.Interface.Px(240)); box.Width != want {
		t.Fatalf("the select is %v wide in a cell of %v", box.Width, want)
	}
	titles := func() []string {
		var out []string
		s.Do(func() {
			for _, m := range s.Window.ShownMenus() {
				out = append(out, m.Title)
			}
		})
		return out
	}
	for name, x := range map[string]float32{"left end": box.X + 1, "words": box.X + box.Width/8, "right end": box.Right() - 2} {
		s.Screen.Click(geom.NewPoint(x, box.Y+box.Height/2))
		s.Sync()
		if got := titles(); len(got) != 1 || got[0] != "photo-sync" {
			t.Errorf("a press at the %s opened %v", name, got)
		}
		s.Screen.KeyPress(unison.KeyEscape, 0)
		s.Sync()
		if got := titles(); len(got) != 0 {
			t.Errorf("Escape left %v open", got)
		}
	}
	// A menu hangs below its row, so the other row pressed is the one above.
	lower := rootRect(second)
	far := geom.NewPoint(lower.Right()-2, lower.Y+lower.Height/2)
	s.Screen.Click(far)
	s.Sync()
	s.Screen.Click(far)
	s.Sync()
	if got := titles(); len(got) != 0 {
		t.Errorf("a press on the select with its list open left %v", got)
	}
	s.Screen.Click(far)
	s.Sync()
	s.Screen.Click(geom.NewPoint(box.Right()-2, box.Y+box.Height/2))
	s.Sync()
	if got := titles(); len(got) != 1 || got[0] != "photo-sync" {
		t.Errorf("a press on the other row's select left %v", got)
	}
	s.Screen.MouseMove(geom.NewPoint(box.X, lower.Bottom()+float32(s.UI.Interface.Px(200))), 0)
	s.Sync()
	img := s.Capture()
	marked := func(r geom.Rect) bool {
		y := int(r.Y) + int(r.Height)/2
		return !samePixel(img.At(int(r.Right())-2, y), img.At(int(r.Right())+2*s.UI.Interface.FocusRingWidth()+1, y))
	}
	if marked(lower) {
		t.Error("the row its list was on stays marked with the other row's list open")
	}
	if !marked(box) {
		t.Error("the row whose list is open is not marked")
	}
}

// Short words have the arrow directly after them with the rest of the width
// empty, the words stand the inset in from the left edge, and long words are
// cut with the same arrow still shown inside the width and the whole words a
// tooltip away. A screen reader still hears the whole words.
func TestAPlainSelectCutsLongWordsAndKeepsItsArrow(t *testing.T) {
	var selects []*kvitui.Select
	cell := kvitui.Px(160)
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		var rows *unison.Panel
		rows, selects = plainRows(ui, cell, plainLongWords)
		return rows
	})
	short, long := selects[0], selects[1]
	s.Do(func() { long.Current = "long"; long.MarkForRedraw() })
	s.Sync()
	shortBox, longBox := rootRect(short), rootRect(long)
	img := s.Capture()
	m := s.UI.Interface
	inset := kvitui.PlainSelectInset.Of(s.UI)

	shortInk, longInk := inkColumns(img, shortBox), inkColumns(img, longBox)
	if len(shortInk) == 0 || len(longInk) == 0 {
		t.Fatal("nothing drawn")
	}
	// Words start the inset in, give or take the first letter's own space.
	if first := shortInk[0]; first < inset || first > inset+m.Px(2) {
		t.Errorf("the words start at %d, the inset is %d", first, inset)
	}

	// The rightmost run of ink is the arrow in both.
	shortFirst, shortLast := lastRun(shortInk)
	longFirst, longLast := lastRun(longInk)
	if width := shortLast - shortFirst + 1; width > m.IconSizeSmall() {
		t.Errorf("the short select's arrow is %d wide, more than an arrow's %d: it has run into the words", width, m.IconSizeSmall())
	}
	if shortLast-shortFirst != longLast-longFirst {
		t.Errorf("the arrow is %d wide after short words and %d after long ones", shortLast-shortFirst+1, longLast-longFirst+1)
	}
	for x := 0; x <= shortLast-shortFirst; x++ {
		for y := 0; y < int(shortBox.Height); y++ {
			a := img.At(int(shortBox.X)+shortFirst+x, int(shortBox.Y)+y)
			b := img.At(int(longBox.X)+longFirst+x, int(longBox.Y)+y)
			if !samePixel(a, b) {
				t.Fatalf("the two arrows differ at %d,%d: %v and %v", x, y, a, b)
			}
		}
	}
	// Beside the short words, with the rest of the cell empty.
	if room := int(shortBox.Width) - shortLast; room < m.Px(40) {
		t.Errorf("the arrow stands %d from the right edge of a wide cell, not beside the words", room)
	}
	// Inside the width, and clear of the words before it, on the long row.
	if longLast >= int(longBox.Width) || int(longBox.Width)-longLast < inset {
		t.Errorf("the arrow ends at %d in a width of %v", longLast, longBox.Width)
	}
	if gap := longInk[slices.Index(longInk, longFirst)-1]; longFirst-gap < m.SpaceNear()/2 {
		t.Errorf("the cut words end at %d, the arrow starts at %d", gap, longFirst)
	}

	if got := long.TooltipWords(); got != plainLongWords {
		t.Errorf("the cut words' tooltip says %q", got)
	}
	if got := short.TooltipWords(); got != "" {
		t.Errorf("words that fit have the tooltip %q", got)
	}
	if n := s.Node(long); n == nil || n.Value != plainLongWords {
		t.Errorf("a screen reader hears %+v", n)
	}
}

// A plain select narrower than the room around its words still shows its arrow
// inside its box, since its words get none of the width rather than all of it.
func TestAPlainSelectNarrowerThanItsRoomStaysInItsBox(t *testing.T) {
	var selects []*kvitui.Select
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		var rows *unison.Panel
		rows, selects = plainRows(ui, kvitui.Px(24), plainLongWords)
		return rows
	})
	box := rootRect(selects[0])
	if ink := inkColumns(s.Capture(), box); len(ink) == 0 {
		t.Errorf("a select %v wide draws nothing", box.Width)
	}
}

// Its preferred width is the widest option's, whichever is chosen, so a row
// does not change width when its answer does; it can be given any width past
// that, and every part of it follows the interface size.
func TestAPlainSelectAsksForItsWidestOptionAndFollowsTheInterfaceSize(t *testing.T) {
	var selects []*kvitui.Select
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		var rows *unison.Panel
		rows, selects = plainRows(ui, kvitui.Px(240), "Part of home-auto")
		return rows
	})
	first := selects[0]
	var chosen, other, largest float32
	s.Do(func() {
		minSize, pref, maxSize := first.Sizes(geom.Size{})
		chosen = pref.Width
		if pref.Height != float32(s.UI.Interface.ControlHeight()) {
			t.Errorf("a plain select is %v tall", pref.Height)
		}
		if maxSize.Width < 1000 {
			t.Errorf("a plain select cannot be wider than %v", maxSize.Width)
		}
		if minSize.Width > pref.Width {
			t.Errorf("its least width %v is past its preferred %v", minSize.Width, pref.Width)
		}
		first.Current = "none"
		_, pref, _ = first.Sizes(geom.Size{})
		other = pref.Width
	})
	if chosen != other {
		t.Errorf("preferred width %v with one option chosen and %v with another", chosen, other)
	}
	s.SetInterfaceSize(tokens.MaxInterfaceSize)
	s.Do(func() {
		_, pref, _ := first.Sizes(geom.Size{})
		largest = pref.Width
		if pref.Height != float32(s.UI.Interface.ControlHeight()) {
			t.Errorf("at the largest interface size a plain select is %v tall", pref.Height)
		}
	})
	if largest <= chosen {
		t.Errorf("preferred width %v at the largest interface size, %v at the default", largest, chosen)
	}
}

// The rest is as any select has it: Up and Down change the choice with no list
// opened, typing chooses, Tab shows the focus ring, a press from a screen
// reader opens the list, Space opens it and choosing from it reports, and the
// words read as a value in a list of rows that has every control named.
func TestAPlainSelectKeepsWhatASelectDoes(t *testing.T) {
	var selects []*kvitui.Select
	var chosen []string
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		var rows *unison.Panel
		rows, selects = plainRows(ui, kvitui.Px(240), plainLongWords)
		selects[0].OnChoose = func(value string) { chosen = append(chosen, value) }
		return rows
	})
	first := selects[0]
	s.CheckNamed()
	s.CheckLabels()

	// Tab moves to the next row's select, which shows the ring, and back.
	s.Screen.KeyPress(unison.KeyTab, 0)
	s.Sync()
	if !selects[1].KeyboardFocus() {
		t.Error("Tab left the next row's select without the ring")
	}
	box := rootRect(selects[1])
	x, y := int(box.X)+int(box.Width)/2, int(box.Y)-1
	if c, ring := s.Capture().At(x, y), kvitui.Color(s.UI.Theme.Tokens().FocusRing); !samePixel(c, ring) {
		t.Errorf("above the focused select the pixel is %v, the ring is %v", c, ring)
	}
	s.Screen.KeyPress(unison.KeyTab, mod.Shift)
	s.Sync()
	if !first.KeyboardFocus() {
		t.Error("Shift+Tab left the first select without the ring")
	}

	s.Screen.KeyPress(unison.KeyDown, 0)
	s.Sync()
	if first.Current != "none" {
		t.Errorf("Down chose %q", first.Current)
	}
	if n := len(s.Window.ShownMenus()); n != 0 {
		t.Errorf("Down opened %d menus", n)
	}
	s.Screen.KeyPress(unison.KeyUp, 0)
	s.Sync()
	if first.Current != "project" {
		t.Errorf("Up chose %q", first.Current)
	}
	s.Screen.Type("n")
	s.Sync()
	if first.Current != "none" {
		t.Errorf("typing n chose %q", first.Current)
	}
	if n := s.Node(first); n == nil || n.Role != role.ComboBox || n.Name != "photo-sync" || n.Value != "Not a project" || !n.Expandable {
		t.Errorf("the select's node: %+v", n)
	}
	if !s.Screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: s.Node(first).ID, Action: accessibility.Press}) {
		t.Error("a screen reader's press did nothing")
	}
	s.Sync()
	if n := len(s.Window.ShownMenus()); n != 1 {
		t.Fatalf("a screen reader's press opened %d menus", n)
	}
	s.Screen.KeyPress(unison.KeyEscape, 0)
	s.Sync()
	s.Screen.KeyPress(unison.KeySpace, 0)
	s.Sync()
	var project *accessibility.Node
	for _, n := range s.Nodes(role.MenuItem) {
		if n.Name == "Project" {
			project = n
		}
	}
	if project == nil {
		t.Fatal("Space opened no menu with the options")
	}
	s.CheckNamed()
	s.Screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: project.ID, Action: accessibility.Press})
	s.Sync()
	if first.Current != "project" || len(s.Window.ShownMenus()) != 0 {
		t.Errorf("choosing Project left %q with %d menus open", first.Current, len(s.Window.ShownMenus()))
	}
	if want := []string{"none", "project", "none", "project"}; !slices.Equal(chosen, want) {
		t.Errorf("reported %v, wanted %v", chosen, want)
	}
}

// A disabled plain select is the same words and arrow, dimmed, and answers
// neither the pointer nor a press.
func TestADisabledPlainSelectStaysQuiet(t *testing.T) {
	var selects []*kvitui.Select
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		var rows *unison.Panel
		rows, selects = plainRows(ui, kvitui.Px(240), plainLongWords)
		selects[0].SetEnabled(false)
		return rows
	})
	box := rootRect(selects[0])
	far := geom.NewPoint(box.Right()-2, box.Y+box.Height/2)
	s.Screen.MouseMove(far, 0)
	s.Sync()
	img := s.Capture()
	if a, b := img.At(int(far.X), int(far.Y)), img.At(int(box.Right())+1, int(far.Y)); !samePixel(a, b) {
		t.Errorf("a disabled select answers the pointer: %v", a)
	}
	s.Screen.Click(far)
	s.Sync()
	if n := len(s.Window.ShownMenus()); n != 0 {
		t.Errorf("a press on a disabled select opened %d menus", n)
	}
}
