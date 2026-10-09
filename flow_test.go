package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/mod"
	"github.com/richardwilkes/unison/enums/role"
)

func TestASwitchTurnsOverAtOnceAndSaysSo(t *testing.T) {
	var sw *kvitui.Switch
	var told []bool
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		sw = kvitui.NewSwitch(ui, "Reduce motion")
		sw.OnChange = func(on bool) { told = append(told, on) }
		return []unison.Paneler{kvitui.Left(sw)}
	})
	screen.Click(screen.PanelCenter(sw))
	screen.Do(func() { sw.Focus() })
	screen.KeyPress(unison.KeySpace, 0)
	if len(told) != 2 || !told[0] || told[1] {
		t.Errorf("a press and Space told %v", told)
	}
	screen.Click(screen.PanelCenter(sw))
	if n := screen.AccessibilityNodeFor(sw); n == nil || n.Role != role.CheckBox || n.Name != "Reduce motion" || n.Checked != check.On {
		t.Errorf("the switch's node is %+v", n)
	}
}

func TestASliderMovesByKeyAndByPointer(t *testing.T) {
	var s *kvitui.Slider
	screen, ui, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		s = kvitui.NewSlider(ui, "Opacity", 0.6)
		s.Unit, s.Precision = "%", 1
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(200), s)}
	})
	screen.Do(func() { s.Focus() })
	screen.KeyPress(unison.KeyRight, 0)
	screen.Do(func() {
		if s.Value < 0.69 || s.Value > 0.71 {
			t.Errorf("Right moved the value to %v, want a tenth of the range more", s.Value)
		}
	})
	screen.KeyPress(unison.KeyEnd, 0)
	if n := screen.AccessibilityNodeFor(s); n == nil || n.Role != role.Slider || n.Name != "Opacity" || n.Number != 1 || n.Value != "1.0 %" {
		t.Errorf("the slider's node is %+v", n)
	}
	// A press on the track takes the handle there.
	screen.Click(screen.PanelPoint(s, geom.NewPoint(float32(ui.Interface.Px(7)), float32(ui.Interface.ControlHeight())/2)))
	screen.Do(func() {
		if s.Value != 0 {
			t.Errorf("a press at the start moved the value to %v", s.Value)
		}
	})
}

func TestARadioGroupIsOneStopAndTheArrowsChoose(t *testing.T) {
	var g *kvitui.RadioGroup
	var chosen []string
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		before := kvitui.NewButton(ui, "Before")
		g = kvitui.NewRadioGroup(ui, "Appearance",
			kvitui.RadioOption{Value: "system", Label: "Follow the system", Detail: "Light or dark, whichever the desktop is set to."},
			kvitui.RadioOption{Value: "light", Label: "Always light"},
			kvitui.RadioOption{Value: "dark", Label: "Always dark"})
		g.Current = "system"
		g.OnChoose = func(v string) { chosen = append(chosen, v) }
		after := kvitui.NewButton(ui, "After")
		return []unison.Paneler{kvitui.Left(before), kvitui.FullWidth(g), kvitui.Left(after)}
	})
	for range 4 {
		if n := focusedNode(screen, w); n != nil && n.Role == role.RadioButton {
			break
		}
		screen.KeyPress(unison.KeyTab, 0)
	}
	if n := focusedNode(screen, w); n == nil || n.Name != "Follow the system. Light or dark, whichever the desktop is set to." || n.Checked != check.On {
		t.Fatalf("Tab reached %+v", n)
	}
	screen.KeyPress(unison.KeyDown, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Always light" || n.Checked != check.On || len(chosen) != 1 || chosen[0] != "light" {
		t.Errorf("Down reached %+v and chose %v", n, chosen)
	}
	// Tab leaves the group rather than walking through it.
	screen.KeyPress(unison.KeyTab, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "After" {
		t.Errorf("Tab from the group reached %+v", n)
	}
}

// Each option's circle lies inside the option at every interface size. Where
// the label's line is shorter than the circle, a circle centred on the line
// reached above the option, and the clip to the option's box cut off its top.
func TestARadioGroupsCirclesAreWhole(t *testing.T) {
	var g *kvitui.RadioGroup
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		g = kvitui.NewRadioGroup(ui, "Theme",
			kvitui.RadioOption{Value: "dark", Label: "Dark", Detail: "Dark background"},
			kvitui.RadioOption{Value: "light", Label: "Light"})
		g.Current = "dark"
		return []unison.Paneler{kvitui.FullWidth(g)}
	})
	for size := 10; size <= 24; size++ {
		screen.Do(func() {
			ui.Interface.SetFontSize(size)
			w.Content().MarkForLayoutRecursively()
			w.ValidateLayout()
			for _, opt := range g.Children() {
				box := opt.ContentRect(false)
				circle := opt.Children()[0].FrameRect()
				if circle.Y < box.Y || circle.Bottom() > box.Bottom() {
					t.Errorf("at %d px the circle %v reaches outside its option %v", size, circle, box)
				}
			}
		})
	}
}

func TestAProgressBarSaysHowFarThroughWhat(t *testing.T) {
	var known, unknown *kvitui.Progress
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		known = kvitui.NewProgress(ui, "Importing statement", 0.47)
		unknown = kvitui.NewProgress(ui, "Reconciling", 0)
		unknown.Determinate = false
		return []unison.Paneler{kvitui.FullWidth(known), kvitui.FullWidth(unknown)}
	})
	if n := screen.AccessibilityNodeFor(known); n == nil || n.Role != role.ProgressBar || n.Name != "Importing statement" || n.Description != "47 percent" {
		t.Errorf("the known bar's node is %+v", n)
	}
	if n := screen.AccessibilityNodeFor(unknown); n == nil || n.Description != "in progress" || !n.Busy {
		t.Errorf("the unknown bar's node is %+v", n)
	}
}

func TestASplitViewsHandleMovesWithTheKeys(t *testing.T) {
	var left *kvitui.Panel
	var split *kvitui.SplitView
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		left = kvitui.NewPanel(ui)
		split = kvitui.NewSplitView(ui, left, kvitui.NewPanel(ui))
		split.SetSize(0, kvitui.Px(160))
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(120), split))}
	})
	for range 3 {
		if n := focusedNode(screen, w); n != nil && n.Role == role.Separator {
			break
		}
		screen.KeyPress(unison.KeyTab, 0)
	}
	screen.KeyPress(unison.KeyRight, 0)
	screen.Do(func() {
		if got, want := left.FrameRect().Width, float32(ui.Interface.Px(176)); got != want {
			t.Errorf("Right made the left region %v wide, want %v", got, want)
		}
	})
	// A drag moves it as far as the pointer went.
	var at geom.Point
	screen.Do(func() { at = geom.NewPoint(left.RectToRoot(left.ContentRect(true)).Right()+4, 60) })
	screen.Drag(at, geom.NewPoint(at.X-40, at.Y), 4)
	screen.Do(func() {
		if got, want := left.FrameRect().Width, float32(ui.Interface.Px(176))-40; got < want-1 || got > want+1 {
			t.Errorf("a drag of 40 px made the left region %v wide, want %v", got, want)
		}
	})
}

func TestATreeOpensAndStepsWithTheArrows(t *testing.T) {
	var tree *kvitui.Tree
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		leaf := func(s string) kvitui.TreeNode { return kvitui.TreeNode{Label: s} }
		tree = kvitui.NewTree(ui, "Accounts",
			kvitui.TreeNode{Label: "Everyday", Children: []kvitui.TreeNode{leaf("Checking"), leaf("Cash")}},
			kvitui.TreeNode{Label: "Savings", Children: []kvitui.TreeNode{
				{Label: "Certificates", Children: []kvitui.TreeNode{leaf("18 months")}}}})
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(250), tree))}
	})
	screen.Do(func() { tree.RequestFocus() })
	name := func() string {
		if n := focusedNode(screen, w); n != nil {
			return n.Name
		}
		return ""
	}
	if name() != "Everyday" {
		t.Fatalf("the tree's first node is %q", name())
	}
	// Right opens, Right again steps in, Left steps out, Left again closes.
	screen.KeyPress(unison.KeyRight, 0)
	screen.KeyPress(unison.KeyRight, 0)
	if name() != "Checking" {
		t.Errorf("Right twice reached %q", name())
	}
	screen.KeyPress(unison.KeyLeft, 0)
	if name() != "Everyday" {
		t.Errorf("Left from a child reached %q", name())
	}
	screen.KeyPress(unison.KeyLeft, 0)
	screen.KeyPress(unison.KeyEnd, 0)
	if name() != "Savings" {
		t.Errorf("End with the first node closed reached %q", name())
	}
	if n := focusedNode(screen, w); n == nil || n.Role != role.Row || n.Level != 1 || !n.Expandable || n.Expanded {
		t.Errorf("the closed node's row is %+v", n)
	}
}

func TestATypeAheadOffersMatchesAndTakesOne(t *testing.T) {
	var cat *kvitui.TypeAhead
	var chosen []string
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		cat = kvitui.NewTypeAhead(ui, "Category", false,
			kvitui.Suggestion{Value: "Groceries"}, kvitui.Suggestion{Value: "Transport"}, kvitui.Suggestion{Value: "Travel"})
		cat.OnChoose = func(v string) { chosen = append(chosen, v) }
		return []unison.Paneler{kvitui.Width(ui, kvitui.Px(240), cat)}
	})
	// Filled by the application, the list stays shut.
	screen.Do(func() { cat.SetText("Tr") })
	screen.Sync()
	if n := len(w.Popups()); n != 0 {
		t.Errorf("text set by the application opened %d lists", n)
	}
	screen.Do(func() { cat.SetText("") })
	screen.Click(screen.PanelCenter(cat.Edit()))
	screen.Type("tr")
	if n := len(w.Popups()); n != 1 {
		t.Fatalf("typing opened %d lists", n)
	}
	if n := screen.AccessibilityNodeFor(cat.Edit()); n == nil || n.Description != "2 suggestions" {
		t.Errorf("the field's node is %+v", n)
	}
	screen.KeyPress(unison.KeyDown, 0)
	screen.KeyPress(unison.KeyReturn, 0)
	if len(chosen) != 1 || chosen[0] != "Travel" || cat.Text() != "Travel" || len(w.Popups()) != 0 {
		t.Errorf("Down and Return chose %v, left %q and %d lists", chosen, cat.Text(), len(w.Popups()))
	}
	// Nothing matching, and nothing new allowed: no list, and Return makes nothing.
	screen.Do(func() { cat.SetText("") })
	screen.Type("zz")
	screen.KeyPress(unison.KeyReturn, 0)
	if len(chosen) != 1 || len(w.Popups()) != 0 {
		t.Errorf("an unknown category chose %v with %d lists", chosen, len(w.Popups()))
	}
}

func TestADualListMovesAndOrdersFromTheKeyboard(t *testing.T) {
	var d *kvitui.DualList
	opts := func(names ...string) []kvitui.Option {
		var out []kvitui.Option
		for _, n := range names {
			out = append(out, kvitui.Option{Value: n, Label: n})
		}
		return out
	}
	labels := func(o []kvitui.Option) string {
		s := ""
		for _, x := range o {
			s += x.Label + " "
		}
		return s
	}
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		d = kvitui.NewDualList(ui, opts("Payee", "Account"), opts("Date", "Amount"))
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(220), d))}
	})
	for range 3 {
		if n := focusedNode(screen, w); n != nil && n.Name == "Payee" {
			break
		}
		screen.KeyPress(unison.KeyTab, 0)
	}
	screen.KeyPress(unison.KeyDown, 0)
	screen.KeyPress(unison.KeyReturn, 0)
	if labels(d.Available) != "Payee " || labels(d.Chosen) != "Date Amount Account " {
		t.Errorf("Down and Return left %q and %q", labels(d.Available), labels(d.Chosen))
	}
	// Tab through the buttons to the chosen list, then move its last item earlier.
	for range 6 {
		if n := focusedNode(screen, w); n != nil && n.Name == "Date" {
			break
		}
		screen.KeyPress(unison.KeyTab, 0)
	}
	screen.KeyPress(unison.KeyEnd, 0)
	var earlier *accessibility.Node
	for _, n := range screen.AccessibilityTree(w.Window).Nodes {
		if n.Name == "Move the selected column earlier" {
			earlier = n
		}
	}
	if earlier == nil {
		t.Fatal("no button moves a column earlier")
	}
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: earlier.ID, Action: accessibility.Press})
	screen.Sync()
	if labels(d.Chosen) != "Date Account Amount " {
		t.Errorf("moving the last column earlier left %q", labels(d.Chosen))
	}
}

func TestNumberAndMoneyFieldsSayWhatIsWrong(t *testing.T) {
	var days *kvitui.NumberField
	var fare *kvitui.MoneyField
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		days = kvitui.NewNumberField(ui)
		days.Label, days.Minimum, days.Maximum = "Days", 1, 365
		fare = kvitui.NewMoneyField(ui, "BHD", 3)
		fare.Label = "Amount"
		return []unison.Paneler{kvitui.Left(days), kvitui.Left(fare)}
	})
	screen.Do(func() {
		for text, want := range map[string]string{"14": "", "999": "Must be at most 365", "0": "Must be at least 1", "ten": "Not a number", "": ""} {
			days.SetText(text)
			if days.Error != want {
				t.Errorf("%q says %q, want %q", text, days.Error, want)
			}
		}
		fare.SetText("18.75")
		if units, ok := fare.MinorUnits(); !ok || units != 18750 {
			t.Errorf("18.75 dinar is %d minor units (%v), want 18750", units, ok)
		}
	})
	if n := screen.AccessibilityNodeFor(fare.Edit()); n == nil || n.Name != "Amount in BHD" {
		t.Errorf("the money field's node is %+v", n)
	}
	// Leaving the field writes the value at the currency's places.
	screen.Click(screen.PanelCenter(fare.Edit()))
	screen.KeyPress(unison.KeyTab, 0)
	screen.Do(func() {
		if fare.Text() != "18.750" {
			t.Errorf("leaving the field left %q", fare.Text())
		}
	})
}

func TestAHistoryAndAChangeSayThemselves(t *testing.T) {
	var history *kvitui.Timeline
	var pair, same *kvitui.BeforeAfter
	var done *kvitui.ConfirmInPlace
	var undone int
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		history = kvitui.NewTimeline(ui, "Account history",
			kvitui.TimelineEntry{When: "14:02", What: "Statement imported", Who: "agent", Tone: kvitui.ToneSuccess},
			kvitui.TimelineEntry{When: "Yesterday", What: "Account opened", Who: "you"})
		pair = kvitui.NewBeforeAfter(ui, "Amount", "42.00", "44.50")
		pair.Unit = "GBP"
		same = kvitui.NewBeforeAfter(ui, "Date", "2026-08-14", "2026-08-14")
		done = kvitui.NewConfirmInPlace(ui, "Recategorised as Groceries")
		done.Affected = 40
		done.OnUndo = func() { undone++ }
		return []unison.Paneler{kvitui.FullWidth(history), kvitui.Left(pair), kvitui.Left(same), kvitui.FullWidth(done)}
	})
	screen.Do(func() {
		done.Shown = true
		done.MarkForLayoutAndRedraw()
	})
	screen.Sync()
	n := screen.AccessibilityNodeFor(history)
	if n == nil || n.Role != role.List || len(n.Children) != 2 {
		t.Fatalf("the history's node is %+v", n)
	}
	if first := screen.AccessibilityTree(history.Window()).Nodes[n.Children[0]]; first.Name != "14:02, Statement imported, agent" {
		t.Errorf("the newest entry says %q", first.Name)
	}
	if n := screen.AccessibilityNodeFor(pair); n == nil || n.Name != "Amount, was 42.00 GBP, now 44.50 GBP" {
		t.Errorf("the pair says %+v", n)
	}
	if n := screen.AccessibilityNodeFor(same); n == nil || n.Name != "Date, unchanged, was 2026-08-14" {
		t.Errorf("the unchanged pair says %+v", n)
	}
	if n := screen.AccessibilityNodeFor(done); n == nil || n.Name != "Recategorised as Groceries, 40 items affected" {
		t.Errorf("the strip says %+v", n)
	}
	for _, a := range screen.Announcements() {
		if a == "Recategorised as Groceries — 40 items" {
			undone--
		}
	}
	if undone != -1 {
		t.Errorf("the strip was announced %d times, want once", -undone)
	}
}

func TestASpotlightClosesOnEscape(t *testing.T) {
	var spot *kvitui.Spotlight
	var dismissed int
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		return nil
	})
	var target *kvitui.Button
	screen.Do(func() {
		target = kvitui.NewButton(ui, "Import a statement")
		stage := kvitui.FullWidth(kvitui.At(ui, kvitui.Px(40), kvitui.Px(20), kvitui.Px(160), target))
		w.SetBody(stage)
		spot = kvitui.NewSpotlight(ui, target, "Start here", "Import a statement and the dashboard fills itself in.")
		spot.OnDismiss = func() { dismissed++; spot.Close() }
	})
	screen.Sync()
	screen.Do(func() { spot.Open(target.Parent()) })
	screen.Sync()
	if n := screen.AccessibilityNodeFor(spot); n == nil || n.Role != role.Dialog || n.Name != "Start here" {
		t.Errorf("the spotlight's node is %+v", n)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	if dismissed != 1 || spot.Opened() {
		t.Errorf("Escape dismissed %d times and left it open: %v", dismissed, spot.Opened())
	}
}

// A menu opens with no line lit, the arrow keys move over the lines that can
// be chosen, Return chooses, and the focus goes back where it was.
func TestAMenuIsWorkedFromTheKeyboard(t *testing.T) {
	var button *kvitui.Button
	var chosen []string
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewButton(ui, "Actions")
		return []unison.Paneler{kvitui.Left(button)}
	})
	pick := func(s string) func() { return func() { chosen = append(chosen, s) } }
	screen.Do(func() {
		button.Focus()
		ui.ShowMenu(button, "Actions", []kvitui.MenuItem{
			{Text: "Open", Symbol: "file", OnSelect: pick("open")},
			{Text: "Split", Disabled: true, OnSelect: pick("split")},
			{Separator: true},
			{Text: "Delete", Symbol: "trash", Danger: true, Explanation: "Moves it to the bin.", OnSelect: pick("delete")},
		})
	})
	waitFor(t, screen, "the menu to take the focus", func() bool { return !button.Focused() })
	menu := nodes(screen, w, role.Menu)
	if len(menu) != 1 || menu[0].Name != "Actions" || len(menu[0].Children) != 3 {
		t.Fatalf("the menu's node is %+v", menu)
	}
	// Down lights Open; Down again skips the disabled line and the divider.
	screen.KeyPress(unison.KeyDown, 0)
	screen.KeyPress(unison.KeyDown, 0)
	if n := focusedNode(screen, w); n == nil || n.Role != role.MenuItem || n.Name != "Delete" || n.Description != "Moves it to the bin." {
		t.Errorf("Down twice reached %+v", n)
	}
	screen.KeyPress(unison.KeyReturn, 0)
	screen.Do(func() {
		if len(chosen) != 1 || chosen[0] != "delete" || len(w.Popups()) != 0 || !button.Focused() {
			t.Errorf("Return chose %v, left %d popups, and the button focused: %v", chosen, len(w.Popups()), button.Focused())
		}
	})
}

// A press on a select whose list is open closes the list and leaves it
// closed, and the next press opens it again. A press anywhere else closes
// the list and still reaches what it landed on.
func TestPressingAnOpenSelectClosesItsList(t *testing.T) {
	var s *kvitui.Select
	var other *kvitui.Button
	pressed := 0
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		s = kvitui.NewSelect(ui, "Machine",
			kvitui.Option{Value: "here", Label: "This PC"}, kvitui.Option{Value: "cloud", Label: "Cloud"})
		other = kvitui.NewButton(ui, "Other")
		other.OnClick = func() { pressed++ }
		// Beside the select, where its list does not cover it.
		return []unison.Paneler{kvitui.Left(kvitui.Row(ui, kvitui.SizeSpace, s, other))}
	})
	// The chevron sits at the select's right end, as wide as it is tall.
	var chevron geom.Point
	screen.Do(func() {
		r := s.ContentRect(false)
		chevron = geom.NewPoint(r.Right()-r.Height/2, r.Y+r.Height/2)
	})
	chevron = screen.PanelPoint(s, chevron)
	open := func() bool { return len(nodes(screen, w, role.Menu)) == 1 }
	for i, want := range []bool{true, false, true} {
		screen.Click(chevron)
		screen.Sync()
		if open() != want {
			t.Fatalf("press %d on the chevron left the list open: %v", i+1, open())
		}
	}
	screen.Click(screen.PanelCenter(other))
	screen.Sync()
	if open() || pressed != 1 {
		t.Errorf("a press elsewhere left the list open: %v, and pressed that button %d times", open(), pressed)
	}
}

// The Menu key, Shift+F10 and a right-click open a context menu: a field's
// editing commands, and a table header's column menu.
func TestContextMenusOpenFromTheKeyboardAndThePointer(t *testing.T) {
	var field *kvitui.Field
	var table *kvitui.Table
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		field = kvitui.NewField(ui)
		field.Label = "Payee"
		table = kvitui.NewTable(ui, kvitui.NewBenchmarkTableModel(40))
		return []unison.Paneler{kvitui.Left(field), kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(200), table))}
	})
	items := func() []string {
		var out []string
		for _, n := range nodes(screen, w, role.MenuItem) {
			out = append(out, n.Name)
		}
		return out
	}
	screen.Click(screen.PanelCenter(field.Edit()))
	screen.KeyPress(unison.KeyMenu, 0)
	if got := items(); len(got) != 4 {
		t.Errorf("the Menu key in a field opened %v", got)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	screen.KeyPress(unison.KeyF10, mod.Shift)
	if got := items(); len(got) != 4 {
		t.Errorf("Shift+F10 in a field opened %v", got)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	if len(w.Popups()) != 0 {
		t.Fatal("Escape left the menu open")
	}
	// A right-click on the table's header opens the column's menu.
	screen.ClickWith(screen.PanelPoint(table, geom.NewPoint(20, 10)), unison.ButtonRight, 0)
	found := false
	for _, n := range items() {
		found = found || n == "Hide this column"
	}
	if !found {
		t.Errorf("a right-click on the header opened %v", items())
	}
	screen.KeyPress(unison.KeyEscape, 0)
	// The Menu key on the header opens the menu of the column its cursor is
	// on, under that column.
	for range 4 {
		if n := focusedNode(screen, w); n != nil && n.Role == role.ColumnHeader {
			break
		}
		screen.KeyPress(unison.KeyTab, 0)
	}
	screen.KeyPress(unison.KeyRight, 0)
	screen.KeyPress(unison.KeyMenu, 0)
	menus := nodes(screen, w, role.Menu)
	if len(menus) != 1 || menus[0].Name != "Reference" {
		t.Errorf("the Menu key on the header's second column opened %+v", menus)
	}
}

// A line with items opens a submenu beside it: Right or Return goes in with
// its first line lit, Left comes back out to the line, and choosing inside it
// closes the whole menu and gives the focus back.
func TestASubmenuOpensBesideItsLine(t *testing.T) {
	var button *kvitui.Button
	var chosen []string
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewButton(ui, "Block")
		return []unison.Paneler{kvitui.Left(button)}
	})
	pick := func(s string) func() { return func() { chosen = append(chosen, s) } }
	screen.Do(func() {
		button.Focus()
		ui.ShowMenu(button, "Block", []kvitui.MenuItem{
			{Text: "Duplicate", OnSelect: pick("duplicate")},
			{Text: "Turn into", Items: []kvitui.MenuItem{
				{Text: "Heading", OnSelect: pick("heading")},
				{Text: "Quote", OnSelect: pick("quote")},
			}},
		})
	})
	waitFor(t, screen, "the menu to take the focus", func() bool { return !button.Focused() })
	screen.KeyPress(unison.KeyDown, 0)
	screen.KeyPress(unison.KeyDown, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Turn into" || !n.Expandable || n.Expanded {
		t.Fatalf("the submenu's line is %+v", n)
	}
	screen.KeyPress(unison.KeyRight, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Heading" {
		t.Errorf("Right went in to %+v", n)
	}
	if len(w.Popups()) != 2 {
		t.Errorf("a submenu open shows %d popups, want 2", len(w.Popups()))
	}
	screen.KeyPress(unison.KeyLeft, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Turn into" || len(w.Popups()) != 1 {
		t.Errorf("Left came back out to %+v with %d popups", n, len(w.Popups()))
	}
	screen.KeyPress(unison.KeyReturn, 0)
	screen.KeyPress(unison.KeyDown, 0)
	screen.KeyPress(unison.KeyReturn, 0)
	screen.Do(func() {
		if len(chosen) != 1 || chosen[0] != "quote" || len(w.Popups()) != 0 || !button.Focused() {
			t.Errorf("choosing in the submenu chose %v, left %d popups, and the button focused: %v", chosen, len(w.Popups()), button.Focused())
		}
	})
}

// Resting the pointer on a line with items opens its submenu, and resting on
// another line closes it again.
func TestASubmenuFollowsThePointer(t *testing.T) {
	var button *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewButton(ui, "Block")
		return []unison.Paneler{kvitui.Left(button)}
	})
	screen.Do(func() {
		ui.ShowMenu(button, "Block", []kvitui.MenuItem{
			{Text: "Duplicate"},
			{Text: "Turn into", Items: []kvitui.MenuItem{{Text: "Heading"}}},
		})
	})
	screen.Sync()
	var menu unison.Paneler
	screen.Do(func() { menu = w.Popups()[0].Panel })
	m := ui.Interface
	line := func(i int) geom.Point {
		return screen.PanelPoint(menu, geom.NewPoint(40, float32(m.Hairline()+m.SpaceSnug()+i*m.RowHeightSlim()+m.RowHeightSlim()/2)))
	}
	screen.MouseMove(line(1), 0)
	waitFor(t, screen, "the submenu to open under the pointer", func() bool { return len(w.Popups()) == 2 })
	screen.MouseMove(line(0), 0)
	waitFor(t, screen, "the submenu to close for another line", func() bool { return len(w.Popups()) == 1 })
}
