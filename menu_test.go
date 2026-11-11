package kvitui_test

import (
	"fmt"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"

	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// longMenuItems are lines naming the reader's own things, as a list of
// projects does.
func longMenuItems(n int) []kvitui.MenuItem {
	items := make([]kvitui.MenuItem, n)
	for i := range items {
		items[i] = kvitui.MenuItem{Text: fmt.Sprintf("Project %02d", i+1), Verbatim: true}
	}
	return items
}

// menuPanel is the open Kvit menu's panel.
func menuPanel(w *kvitui.Window) *unison.Panel {
	return w.Popups()[0].Panel.AsPanel()
}

// menuRoot is the menu's bounds in the window.
func menuRoot(t *testing.T, screen *unison.HeadlessScreen, w *kvitui.Window) geom.Rect {
	t.Helper()
	var r geom.Rect
	screen.Do(func() {
		m := menuPanel(w)
		r = m.RectToRoot(m.ContentRect(false))
	})
	return r
}

// contentRoot is the window content's bounds in the window.
func contentRoot(screen *unison.HeadlessScreen, w *kvitui.Window) geom.Rect {
	var r geom.Rect
	screen.Do(func() { r = w.Content().RectToRoot(w.Content().ContentRect(false)) })
	return r
}

func TestALongMenuCapsItsHeightAndScrolls(t *testing.T) {
	var button *kvitui.Button
	var chosen []string
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewButton(ui, "Projects")
		return []unison.Paneler{kvitui.Left(button)}
	})
	// A window 600 design pixels tall, as the request names.
	screen.Do(func() {
		r := w.Window.ContentRect()
		w.Window.SetContentRect(geom.NewRect(r.X, r.Y, r.Width, 600))
	})
	items := longMenuItems(40)
	for i := range items {
		i := i
		items[i].OnSelect = func() { chosen = append(chosen, items[i].Text) }
	}
	screen.Do(func() {
		button.Focus()
		ui.ShowMenu(button, "Projects", items)
	})
	waitFor(t, screen, "the menu to take the focus", func() bool { return !button.Focused() })

	m := ui.Interface
	frame := float32(m.Hairline() + m.SpaceSnug())
	maxInside := float32(15 * m.RowHeightSlim())
	r := menuRoot(t, screen, w)
	if got, want := r.Height, maxInside+2*frame; got != want {
		t.Errorf("a 40-line menu is %.1f tall, want the 15-line maximum %.1f", got, want)
	}
	content := contentRoot(screen, w)
	if r.Y < content.Y-1 || r.Bottom() > content.Bottom()+1 {
		t.Errorf("the menu at %v runs past the window content at %v", r, content)
	}

	// End lights the last line and it is drawn inside the frame.
	screen.KeyPress(unison.KeyEnd, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Project 40" {
		t.Fatalf("End lit %+v, want Project 40", n)
	}
	last := nodes(screen, w, role.MenuItem)[len(nodes(screen, w, role.MenuItem))-1]
	r = menuRoot(t, screen, w)
	if last.Bounds.Y < r.Y-1 || last.Bounds.Bottom() > r.Bottom()+1 {
		t.Errorf("the last line at %v is not drawn inside the menu at %v", last.Bounds, r)
	}
	// The first line has scrolled out of sight.
	first := nodes(screen, w, role.MenuItem)[0]
	if first.Bounds.Y >= r.Y-1 {
		t.Errorf("the first line at %v is still in view with the last line lit", first.Bounds)
	}

	// Home brings the first line back.
	screen.KeyPress(unison.KeyHome, 0)
	if n := focusedNode(screen, w); n == nil || n.Name != "Project 01" {
		t.Fatalf("Home lit %+v, want Project 01", n)
	}

	// A wheel turn moves the view.
	before := nodes(screen, w, role.MenuItem)[0].Bounds.Y
	screen.Wheel(screen.PanelCenter(menuPanel(w)), geom.NewPoint(0, -1), 0)
	screen.Sync()
	after := nodes(screen, w, role.MenuItem)[0].Bounds.Y
	if after >= before {
		t.Errorf("a wheel turn moved the first line from %.1f to %.1f, want it higher", before, after)
	}

	// A press on a line in the scrolled view chooses that line.
	screen.KeyPress(unison.KeyEnd, 0)
	last = nodes(screen, w, role.MenuItem)[len(nodes(screen, w, role.MenuItem))-1]
	screen.Click(last.Bounds.Center())
	screen.Sync()
	screen.Do(func() {
		if len(chosen) != 1 || chosen[0] != "Project 40" {
			t.Errorf("a press on the scrolled view chose %v", chosen)
		}
		if len(w.Popups()) != 0 {
			t.Errorf("the menu stayed open after choosing")
		}
	})
}

func TestAShortMenuOpensWhole(t *testing.T) {
	var button *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewButton(ui, "Actions")
		return []unison.Paneler{kvitui.Left(button)}
	})
	screen.Do(func() {
		button.Focus()
		ui.ShowMenu(button, "Actions", longMenuItems(10))
	})
	waitFor(t, screen, "the menu to take the focus", func() bool { return !button.Focused() })
	m := ui.Interface
	frame := float32(m.Hairline() + m.SpaceSnug())
	want := float32(10*m.RowHeightSlim()) + 2*frame
	if got := menuRoot(t, screen, w).Height; got != want {
		t.Errorf("a 10-line menu is %.1f tall, want its whole height %.1f", got, want)
	}
	// With nothing to scroll the first line stays in view on End.
	screen.KeyPress(unison.KeyEnd, 0)
	r := menuRoot(t, screen, w)
	if first := nodes(screen, w, role.MenuItem)[0]; first.Bounds.Y < r.Y-1 {
		t.Errorf("a 10-line menu scrolled on End: the first line at %v left %v", first.Bounds, r)
	}
	screen.KeyPress(unison.KeyEscape, 0)
}

func TestAMenuWithNoRoomBelowOpensAbove(t *testing.T) {
	var button *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewButton(ui, "Projects")
		// The button at the foot of the window, with no room below it.
		spacer := kvitui.Height(ui, kvitui.Px(700), unison.NewPanel())
		return []unison.Paneler{kvitui.Column(ui, kvitui.SizeSpace, spacer, kvitui.Left(button))}
	})
	screen.Do(func() {
		button.Focus()
		ui.ShowMenu(button, "Projects", longMenuItems(40))
	})
	waitFor(t, screen, "the menu to take the focus", func() bool { return !button.Focused() })
	var anchor, menu geom.Rect
	screen.Do(func() {
		anchor = button.RectToRoot(button.ContentRect(false))
		m := menuPanel(w)
		menu = m.RectToRoot(m.ContentRect(false))
	})
	if menu.Bottom() > anchor.Y+1 {
		t.Errorf("the menu at %v did not open above its anchor at %v", menu, anchor)
	}
	content := contentRoot(screen, w)
	if menu.Y < content.Y-1 || menu.Bottom() > content.Bottom()+1 {
		t.Errorf("the menu at %v runs past the window content at %v", menu, content)
	}
	screen.KeyPress(unison.KeyEscape, 0)
}

func TestAMenuOpensOnItsCheckedLine(t *testing.T) {
	var button *kvitui.Button
	screen, ui, w := windowSession(t, func(ui *kvitui.UI) []unison.Paneler {
		button = kvitui.NewButton(ui, "Projects")
		return []unison.Paneler{kvitui.Left(button)}
	})
	items := longMenuItems(40)
	items[29].Checked = true
	screen.Do(func() {
		button.Focus()
		ui.ShowMenu(button, "Projects", items)
	})
	waitFor(t, screen, "the menu to take the focus", func() bool { return !button.Focused() })
	r := menuRoot(t, screen, w)
	var checked *struct {
		Bounds geom.Rect
	}
	for _, n := range nodes(screen, w, role.MenuItem) {
		if n.Name == "Project 30" {
			b := n.Bounds
			checked = &struct{ Bounds geom.Rect }{Bounds: b}
		}
	}
	if checked == nil {
		t.Fatal("no Project 30 in the menu")
	}
	if checked.Bounds.Y < r.Y-1 || checked.Bounds.Bottom() > r.Bottom()+1 {
		t.Errorf("the checked line at %v is not in view in the menu at %v", checked.Bounds, r)
	}
	screen.KeyPress(unison.KeyEscape, 0)
}
