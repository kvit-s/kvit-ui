package kvitui_test

import (
	"testing"
	"time"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// windowSession opens a Kvit window, which is where popups are shown, with
// the panels build returns in a padded column as its body.
func windowSession(t *testing.T, build func(ui *kvitui.UI) []unison.Paneler) (*unison.HeadlessScreen, *kvitui.UI, *kvitui.Window) {
	t.Helper()
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	ui.Theme.SetReducedMotion(true)
	var w *kvitui.Window
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 1500, Height: 1000},
		unison.StartupFinishedCallback(func() {
			if w, err = kvitui.NewWindow(ui, "feedback"); err != nil {
				t.Error(err)
				return
			}
			col := kvitui.Column(ui, kvitui.SizeSpace, build(ui)...)
			col.SetBorder(kvitui.Padding(ui, kvitui.SizeViewMargin))
			w.SetBody(col)
			w.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.EnableAccessibility()
	screen.Sync()
	return screen, ui, w
}

// waitFor runs the event loop until cond holds, or fails after two seconds.
func waitFor(t *testing.T, screen *unison.HeadlessScreen, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		ok := false
		screen.Do(func() { ok = cond() })
		if ok {
			return
		}
		time.Sleep(20 * time.Millisecond)
		screen.Sync()
	}
	t.Fatalf("waited two seconds for %s", what)
}

func TestATooltipShowsForTheKeyboardBesideItsControl(t *testing.T) {
	var edit, remove *kvitui.IconButton
	screen, _, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		edit = kvitui.NewIconButton(ui, "pencil", "Edit")
		remove = kvitui.NewIconButton(ui, "trash", "Delete")
		remove.Explanation = "Moves the note to the bin."
		return []unison.Paneler{kvitui.Row(ui, kvitui.SizeSpace, edit, remove)}
	})
	screen.Do(func() { edit.Focus() })
	screen.KeyPress(unison.KeyTab, 0)
	waitFor(t, screen, "the tooltip", func() bool { return len(w.Popups()) == 1 })
	screen.Do(func() {
		tip := w.Popups()[0].Panel.AsPanel().FrameRect()
		box := w.Content().RectFromRoot(remove.RectToRoot(remove.ContentRect(true)))
		// Beside the control on its trailing side, level with its top.
		if tip.X < box.Right() || tip.Y != box.Y {
			t.Errorf("the tooltip is at %v for a control at %v", tip, box)
		}
	})
	screen.KeyPress(unison.KeyTab, 0)
	waitFor(t, screen, "the tooltip to go", func() bool { return len(w.Popups()) == 0 })
}

func TestAPopoverClosesOnEscapeAndGivesTheFocusBack(t *testing.T) {
	var opener *kvitui.Button
	var pop *kvitui.Popover
	var inside *kvitui.Check
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		opener = kvitui.NewButton(ui, "Filter")
		inside = kvitui.NewCheck(ui, "Settled")
		return []unison.Paneler{kvitui.Left(opener)}
	})
	screen.Do(func() {
		pop = kvitui.NewPopover(ui, "Filter", inside)
		opener.Focus()
		pop.Open(opener, kvitui.PlaceBelow(ui, opener))
	})
	waitFor(t, screen, "the popover to take the focus", func() bool { return inside.Focused() })
	if n := screen.AccessibilityNodeFor(pop); n == nil || n.Role != role.Dialog || n.Name != "Filter" {
		t.Errorf("the popover's node: %+v", n)
	}
	screen.KeyPress(unison.KeyEscape, 0)
	screen.Do(func() {
		if pop.Opened() || len(w.Popups()) != 0 {
			t.Error("Escape left the popover open")
		}
		if !opener.Focused() {
			t.Error("closing did not give the focus back to the button that opened it")
		}
	})
	// A press outside closes it too.
	screen.Do(func() { pop.Open(opener, kvitui.PlaceBelow(ui, opener)) })
	screen.Sync()
	screen.Click(geom.NewPoint(1400, 900))
	screen.Do(func() {
		if pop.Opened() {
			t.Error("a press outside left the popover open")
		}
	})
}

// A tooltip never has the focus, so Escape passes it by: over an open
// popover it closes the popover rather than stopping at the tooltip. A
// popup that is not a tooltip and has no Escape of its own still stops it,
// so the control with the focus gets the key, as the typeahead's field does
// to close its list.
func TestEscapePassesATooltipByToThePopupUnderIt(t *testing.T) {
	var opener *kvitui.Button
	var pop *kvitui.Popover
	var inside *kvitui.Check
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		opener = kvitui.NewButton(ui, "Filter")
		inside = kvitui.NewCheck(ui, "Settled")
		return []unison.Paneler{kvitui.Left(opener)}
	})
	screen.Do(func() {
		pop = kvitui.NewPopover(ui, "Filter", inside)
		opener.Focus()
		pop.Open(opener, kvitui.PlaceBelow(ui, opener))
	})
	waitFor(t, screen, "the popover to take the focus", func() bool { return inside.Focused() })
	screen.Do(func() { ui.ShowTooltip(inside, "Only rows that are settled") })
	waitFor(t, screen, "the tooltip over the popover", func() bool { return len(w.Popups()) == 2 })
	screen.KeyPress(unison.KeyEscape, 0)
	screen.Do(func() {
		if pop.Opened() {
			t.Error("Escape stopped at the tooltip and left the popover open")
		}
	})

	screen.Do(func() { pop.Open(opener, kvitui.PlaceBelow(ui, opener)) })
	waitFor(t, screen, "the popover again", func() bool { return inside.Focused() })
	var hideList func()
	screen.Do(func() {
		list := kvitui.NewLabel(ui, "a list that takes keys")
		hideList = w.Show(&kvitui.Popup{Panel: list, Anchor: inside,
			Place: func(bounds geom.Rect, size geom.Size) geom.Rect { return geom.Rect{Size: size} }})
	})
	screen.KeyPress(unison.KeyEscape, 0)
	screen.Do(func() {
		if !pop.Opened() {
			t.Error("Escape went past a popup that takes keys and closed the popover under it")
		}
		hideList()
		pop.Close()
	})
}

func TestADialogHasToBeAnsweredFirst(t *testing.T) {
	var behind *kvitui.Button
	var dialog *kvitui.Dialog
	pressed, rejected, accepted := 0, 0, 0
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		behind = kvitui.NewButton(ui, "Behind")
		behind.OnClick = func() { pressed++ }
		return []unison.Paneler{kvitui.Left(behind)}
	})
	screen.Do(func() {
		dialog = kvitui.NewDialog(ui, "Delete 40 transactions?")
		dialog.ConfirmText, dialog.Destructive = "Delete them", true
		dialog.OnReject = func() { rejected++ }
		dialog.OnAccept = func() { accepted++ }
		dialog.Open(w)
	})
	screen.Sync()
	screen.Click(screen.PanelCenter(behind))
	if pressed != 0 {
		t.Error("a button behind a modal dialog was pressed")
	}
	// A destructive dialog gives the keyboard no default: Return must not
	// delete anything.
	screen.KeyPress(unison.KeyReturn, 0)
	if accepted != 0 {
		t.Error("Return confirmed a destructive dialog")
	}
	screen.KeyPress(unison.KeyEscape, 0)
	if rejected != 1 || len(w.Popups()) != 0 {
		t.Errorf("Escape rejected %d times and left %d popups", rejected, len(w.Popups()))
	}
	screen.Click(screen.PanelCenter(behind))
	if pressed != 1 {
		t.Error("the button behind could not be pressed once the dialog was answered")
	}
}

func TestAToastWithoutAnActionGoesAndOneWithAnActionStays(t *testing.T) {
	var plain, undo *kvitui.Toast
	gone := map[*kvitui.Toast]bool{}
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		plain = kvitui.NewToast(ui, "Copied")
		undo = kvitui.NewToast(ui, "40 archived")
		undo.Action = "Undo"
		for _, x := range []*kvitui.Toast{plain, undo} {
			x.Timeout = 50 * time.Millisecond
			x.OnDismiss = func() { gone[x] = true }
		}
		return []unison.Paneler{kvitui.Left(plain), kvitui.Left(undo)}
	})
	screen.Do(func() {
		plain.Shown()
		undo.Shown()
	})
	waitFor(t, screen, "the plain toast to go", func() bool { return gone[plain] })
	time.Sleep(100 * time.Millisecond)
	screen.Sync()
	if gone[undo] {
		t.Error("a toast with an undo timed out")
	}
}

func TestANoticeSaysItsConditionAndIsDismissedOnlyWhenItCanBe(t *testing.T) {
	var licence, vault *kvitui.Notice
	screen, _, _ := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		licence = kvitui.NewNotice(ui, "Your licence expires in three days")
		licence.Detail, licence.Dismissible = "Renew before 30 March.", true
		vault = kvitui.NewNotice(ui, "This vault could not be saved")
		return []unison.Paneler{licence, vault}
	})
	if n := screen.AccessibilityNodeFor(licence); n == nil || n.Name != "Your licence expires in three days. Renew before 30 March." {
		t.Errorf("the notice's node: %+v", n)
	}
	count := func(n *kvitui.Notice) (shown int) {
		screen.Do(func() {
			for _, c := range n.Children() {
				if _, ok := c.Self.(*kvitui.IconButton); ok && !c.Hidden {
					shown++
				}
			}
		})
		return shown
	}
	if count(licence) != 1 || count(vault) != 0 {
		t.Errorf("close controls: %d on the dismissible notice, %d on the other", count(licence), count(vault))
	}
}

// A floating view takes every press on what it covers, closes on Escape, and
// says what it holds.
func TestAFloatingViewCoversTheListAndClosesOnEscape(t *testing.T) {
	var view *kvitui.FloatingView
	var list *kvitui.ListRow
	var pressed, asked int
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		list = kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Whole Foods"))
		list.Interactive = true
		list.OnActivate = func() { pressed++ }
		return []unison.Paneler{kvitui.FullWidth(kvitui.Height(ui, kvitui.Px(260), list))}
	})
	screen.Do(func() {
		view = kvitui.NewFloatingView(ui, "Whole Foods · 3 Sep", kvitui.NewLabel(ui, "Payee"))
		view.CloseLabel = "Close record"
		view.OnCloseRequested = func() { asked++; view.Close() }
		view.Open(list.Parent())
	})
	screen.Sync()
	// A press at the list's left end, where the dimmed area is, does not
	// reach the row under it.
	screen.Click(screen.PanelPoint(list, geom.NewPoint(4, 10)))
	if pressed != 0 || asked != 0 || !view.Opened() {
		t.Errorf("a press on the dimmed area pressed the row %d times and asked to close %d times", pressed, asked)
	}
	if n := screen.AccessibilityNodeFor(view); n == nil || n.Name != "Whole Foods · 3 Sep" {
		t.Errorf("the view's node is %+v", n)
	}
	found := false
	for _, n := range screen.AccessibilityTree(w.Window).Nodes {
		found = found || (n.Role == role.Button && n.Name == "Close record")
	}
	if !found {
		t.Error("no close control named Close record")
	}
	screen.KeyPress(unison.KeyEscape, 0)
	if asked != 1 || view.Opened() {
		t.Errorf("Escape asked to close %d times and left the view open: %v", asked, view.Opened())
	}
}
