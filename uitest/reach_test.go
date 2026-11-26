package uitest_test

import (
	"strconv"
	"strings"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/uitest"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// buttons is a column of n buttons, "Item 1" to "Item n".
func buttons(ui *kvitui.UI, n int) *unison.Panel {
	var items []unison.Paneler
	for i := range n {
		items = append(items, kvitui.Left(kvitui.NewButton(ui, "Item "+strconv.Itoa(i+1))))
	}
	return kvitui.Column(ui, kvitui.SizeSpace, items...)
}

// atFloor gives a session's window the smallest size it allows.
func atFloor(s *uitest.Session) {
	s.Do(func() {
		m := s.UI.Interface
		was := s.Window.ContentRect()
		s.Window.SetContentRect(geom.NewRect(was.X, was.Y, float32(m.WidthFloor()), float32(m.HeightFloor())))
	})
	s.Sync()
}

// problemFor is the problem reported for what is named, or "".
func problemFor(problems []uitest.Problem, name string) string {
	for _, p := range problems {
		if strings.Contains(p.Where, quoteName(name)) {
			return p.String()
		}
	}
	return ""
}

func quoteName(s string) string { return `"` + s + `"` }

// A page taller than the window, with nothing to scroll it, leaves its last
// buttons out of reach, and the check names them: the first one cut across,
// and the ones wholly below the window's edge.
func TestTheReachCheckFindsButtonsBelowTheWindowsEdge(t *testing.T) {
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler { return buttons(ui, 40) })
	atFloor(s)
	var problems []uitest.Problem
	s.Do(func() { problems = uitest.Reach(s.Window) })
	if p := problemFor(problems, "Item 1"); p != "" {
		t.Errorf("the first button is in sight, but: %s", p)
	}
	last := problemFor(problems, "Item 40")
	if !strings.HasPrefix(last, "LAY-3") || !strings.Contains(last, "out of sight") {
		t.Errorf("the last button is below the window, and the check says %q", last)
	}
	cut := 0
	for _, p := range problems {
		if strings.Contains(p.What, "its bottom") {
			cut++
		}
	}
	if cut != 1 {
		t.Errorf("%d buttons are reported cut across, want the one at the window's edge: %v", cut, problems)
	}
}

// The same buttons in a region, which scrolls, are all within reach, and so
// are the places of a sidebar that scrolls them under its foot.
func TestTheReachCheckLetsAScrollViewHoldMoreThanItShows(t *testing.T) {
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler { return kvitui.NewLabel(ui, "Body") })
	s.Do(func() {
		var places []unison.Paneler
		for i := range 40 {
			places = append(places, kvitui.NewSidebarItem(s.UI, "Folder "+strconv.Itoa(i+1), "folder"))
		}
		side := kvitui.NewSidebar(s.UI, places...)
		side.AddFooter(kvitui.NewSidebarItem(s.UI, "Trash", "trash"))
		s.Window.SetSidebar(side)
		s.Window.SetBody(kvitui.NewRegion(s.UI, buttons(s.UI, 40)))
	})
	s.Sync()
	s.CheckReach()
}

// CheckReach shrinks the window to its floor and gives it its size back,
// and a known problem keeps a screen that has one passing.
func TestCheckReachRestoresTheWindowAndTakesKnownProblems(t *testing.T) {
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler { return buttons(ui, 40) })
	var before, after geom.Rect
	s.Do(func() { before = s.Window.ContentRect() })
	s.CheckReach("LAY-3")
	s.Do(func() { after = s.Window.ContentRect() })
	if after != before {
		t.Errorf("the window is %v after the check, and was %v", after, before)
	}
}

// An action that changes what is shown passes the response check, and one
// that changes nothing gives no response, even though the press moves the
// pointer onto the button and the keyboard focus to it.
func TestTheResponseCheckTellsAnActionFromNothing(t *testing.T) {
	var copyButton, idle *kvitui.Button
	var state *kvitui.Label
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler {
		state = kvitui.NewLabel(ui, "Not copied yet")
		copyButton = kvitui.NewButton(ui, "Copy")
		copyButton.OnClick = func() {
			state.Text = "Copied"
			state.MarkForLayoutAndRedraw()
		}
		idle = kvitui.NewButton(ui, "Send")
		return kvitui.Column(ui, kvitui.SizeSpace, state, kvitui.Left(copyButton), kvitui.Left(idle))
	})
	s.CheckResponds("pressing Copy", func() { s.Screen.Click(s.Screen.PanelCenter(copyButton)) })
	if got := s.Named("Copied"); got == nil {
		t.Fatal("pressing Copy did not change the label")
	}
	before := s.Tree()
	s.Screen.Click(s.Screen.PanelCenter(idle))
	s.Sync()
	if got := uitest.Responses(before, s.Tree()); len(got) != 0 {
		t.Errorf("a press that does nothing gave %d responses: %+v", len(got), got)
	}
}

// A popup placed past the window's edge is out of reach, while a menu
// hanging from a button that a region has scrolled out of sight is not: it
// comes into view with the button as the region is scrolled to it.
func TestTheReachCheckReadsPopupsByWhereTheyHang(t *testing.T) {
	var last *kvitui.Button
	s := uitest.Open(t, uitest.Options{}, func(ui *kvitui.UI) unison.Paneler { return kvitui.NewLabel(ui, "Body") })
	s.Do(func() {
		column := buttons(s.UI, 40)
		s.Window.SetBody(kvitui.NewRegion(s.UI, column))
		kids := column.Children()
		last = kids[len(kids)-1].Children()[0].Self.(*kvitui.Button)
	})
	atFloor(s)
	s.Do(func() {
		if last.Words() != "Item 40" || last.RectToRoot(last.ContentRect(true)).Y < s.Window.Content().FrameRect().Bottom() {
			t.Errorf("the menu should hang from Item 40, below the window; it hangs from %q at %v", last.Words(), last.RectToRoot(last.ContentRect(true)))
		}
		s.UI.ShowMenuAt(last, geom.Rect{}, "Projects", []kvitui.MenuItem{{Text: "Open"}, {Text: "Archive"}})
	})
	s.Sync()
	var problems []uitest.Problem
	s.Do(func() { problems = uitest.Reach(s.Window) })
	if len(problems) != 0 {
		t.Errorf("a menu beside a button scrolled out of sight was reported: %v", problems)
	}
	s.Do(func() {
		tall := kvitui.NewButton(s.UI, "Tall")
		s.Window.Show(&kvitui.Popup{Panel: tall, Place: func(bounds geom.Rect, size geom.Size) geom.Rect {
			return geom.NewRect(bounds.X, bounds.Bottom()-size.Height/2, size.Width, size.Height)
		}})
	})
	s.Sync()
	s.Do(func() { problems = uitest.Reach(s.Window) })
	if p := problemFor(problems, "Tall"); !strings.Contains(p, "its bottom") || !strings.Contains(p, "the window's edge") {
		t.Errorf("a popup half past the window's bottom edge gave %v", problems)
	}
}
