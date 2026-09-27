package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
	"github.com/richardwilkes/unison/enums/check"
	"github.com/richardwilkes/unison/enums/role"
)

// session opens a headless window holding the panels build returns, laid
// out in a row, and returns the screen.
func session(t *testing.T, build func(ui *kvitui.UI) []unison.Paneler) (*unison.HeadlessScreen, *kvitui.UI) {
	t.Helper()
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	screen, err := unison.StartHeadless(unison.HeadlessConfig{Width: 800, Height: 300},
		unison.StartupFinishedCallback(func() {
			wnd, err := unison.NewWindow("test")
			if err != nil {
				t.Error(err)
				return
			}
			c := wnd.Content()
			panels := build(ui)
			c.SetLayout(&unison.FlexLayout{Columns: len(panels), HSpacing: float32(ui.Interface.Space())})
			for _, p := range panels {
				c.AddChild(p)
			}
			wnd.SetContentRect(geom.NewRect(0, 0, 800, 300))
			wnd.ToFront()
		}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(screen.Stop)
	screen.EnableAccessibility()
	screen.Sync()
	return screen, ui
}

func TestLabelSizesByRoleAndElides(t *testing.T) {
	var body, display, narrow *kvitui.Label
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		body = kvitui.NewLabel(ui, "Accounts")
		display = kvitui.NewLabel(ui, "Accounts")
		display.Role = kvitui.RoleDisplay
		narrow = kvitui.NewLabel(ui, "A destination whose full name does not fit here")
		narrow.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(120, 0)})
		return []unison.Paneler{body, display, narrow}
	})
	screen.Do(func() {
		_, bp, _ := body.Sizes(geom.Size{})
		_, dp, _ := display.Sizes(geom.Size{})
		if dp.Height <= bp.Height || dp.Width <= bp.Width {
			t.Errorf("display %v is not larger than body %v", dp, bp)
		}
		if w := narrow.FrameRect().Width; w > 120.5 {
			t.Errorf("the narrow label is %.1f wide", w)
		}
	})
	if n := screen.AccessibilityNodeFor(body); n == nil || n.Role != role.Label || n.Name != "Accounts" {
		t.Errorf("a label's node: %+v", n)
	}
}

func TestAnUnknownIconIsMarked(t *testing.T) {
	var known, unknown *kvitui.Icon
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		known = kvitui.NewIcon(ui, "search")
		unknown = kvitui.NewIcon(ui, "not-a-symbol")
		return []unison.Paneler{known, unknown}
	})
	screen.Do(func() {
		if !known.Recognized() || unknown.Recognized() {
			t.Error("recognition is wrong")
		}
	})
	// A decorative icon is hidden from screen readers; a labelled one is an image.
	if n := screen.AccessibilityNodeFor(known); n != nil && !n.Ignored {
		t.Errorf("a decorative icon is visible to screen readers as %v", n.Role)
	}
}

func TestIconButtonActivatesAndDescribesItself(t *testing.T) {
	var edit, pin, disabled *kvitui.IconButton
	clicks := 0
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		edit = kvitui.NewIconButton(ui, "pencil", "Edit")
		edit.Explanation = "Opens the record for editing."
		edit.OnClick = func() { clicks++ }
		pin = kvitui.NewIconButton(ui, "pin", "Pin")
		pin.Checkable = true
		disabled = kvitui.NewIconButton(ui, "copy", "Copy")
		disabled.OnClick = func() { clicks += 100 }
		disabled.SetEnabled(false)
		return []unison.Paneler{edit, pin, disabled}
	})
	screen.Click(screen.PanelCenter(edit))
	screen.Click(screen.PanelCenter(disabled))
	if clicks != 1 {
		t.Errorf("%d clicks counted; a disabled button must ignore its click", clicks)
	}
	// A click gives focus without the ring; Space activates.
	screen.Do(func() {
		if edit.KeyboardFocus() {
			t.Error("a click showed the focus ring")
		}
	})
	screen.KeyPress(unison.KeySpace, 0)
	if clicks != 2 {
		t.Errorf("Space did not press the focused button (%d clicks)", clicks)
	}
	// Tab moves the focus and shows the ring.
	screen.KeyPress(unison.KeyTab, 0)
	screen.Do(func() {
		if !pin.KeyboardFocus() {
			t.Error("Tab did not show the focus ring on the next button")
		}
	})
	screen.KeyPress(unison.KeySpace, 0)
	screen.Do(func() {
		if !pin.Checked {
			t.Error("Space did not check the checkable button")
		}
	})
	n := screen.AccessibilityNodeFor(edit)
	if n == nil || n.Role != role.Button || n.Name != "Edit" || n.Description != "Opens the record for editing." {
		t.Fatalf("the button's node: %+v", n)
	}
	if !screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: n.ID, Action: accessibility.Press}) || clicks != 3 {
		t.Error("a screen reader's press did not press the button")
	}
	if pn := screen.AccessibilityNodeFor(pin); pn == nil || pn.Role != role.ToggleButton || pn.Checked != check.On {
		t.Errorf("the checked button's node: %+v", pn)
	}
}

func TestLinkFollowsAndIsALink(t *testing.T) {
	var link *kvitui.Link
	follows := 0
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		link = kvitui.NewLink(ui, "Privacy policy")
		link.Symbol = "external"
		link.OnActivate = func() { follows++ }
		return []unison.Paneler{link}
	})
	screen.Click(screen.PanelCenter(link))
	screen.KeyPress(unison.KeyReturn, 0)
	screen.KeyPress(unison.KeySpace, 0)
	if follows != 3 {
		t.Errorf("click, Return and Space followed the link %d times", follows)
	}
	if n := screen.AccessibilityNodeFor(link); n == nil || n.Role != role.Link || n.Name != "Privacy policy" {
		t.Errorf("the link's node: %+v", n)
	}
}

// The keyboard focus ring is drawn around a control reached with Tab, outside
// its box, and not around one given the focus by a click.
func TestTheFocusRingShowsForTheKeyboardOnly(t *testing.T) {
	var a, b *kvitui.IconButton
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		a = kvitui.NewIconButton(ui, "pencil", "Edit")
		b = kvitui.NewIconButton(ui, "trash", "Delete")
		row := kvitui.Row(ui, kvitui.SizeSpace, a, b)
		row.SetBorder(kvitui.Padding(ui, kvitui.SizeViewMargin))
		return []unison.Paneler{row}
	})
	ringAround := func(p *kvitui.IconButton) bool {
		var r geom.Rect
		screen.Do(func() { r = p.RectToRoot(p.ContentRect(false)) })
		img := screen.Capture()
		want := kvitui.Color(ui.Theme.Tokens().FocusRing)
		cr, cg, cb := uint32(want.Red()), uint32(want.Green()), uint32(want.Blue())
		// Just left of the box, halfway down, where the ring's stroke is.
		for dx := 1; dx <= ui.Interface.FocusRingWidth(); dx++ {
			pr, pg, pb, _ := img.At(int(r.X)-dx, int(r.Y+r.Height/2)).RGBA()
			if pr>>8 == cr && pg>>8 == cg && pb>>8 == cb {
				return true
			}
		}
		return false
	}
	// A window hands its first control the focus when it becomes active.
	// That came from neither the keyboard nor the pointer, and Qt draws no
	// ring for it.
	var auto bool
	screen.Do(func() { auto = a.Focused() && a.KeyboardFocus() })
	if auto || ringAround(a) {
		t.Error("the focus the window handed out drew the focus ring")
	}
	screen.Click(screen.PanelCenter(a))
	if ringAround(a) {
		t.Error("a click drew the focus ring")
	}
	screen.KeyPress(unison.KeyTab, 0)
	if !ringAround(b) {
		t.Error("Tab did not draw the focus ring")
	}
}
