package kvitui_test

import (
	"strconv"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
)

// A turn begun while the last is still easing in moves at once, in its own
// direction: three fast turns down and one back up moves back up on the
// fourth turn, rather than standing still until the debt down is spun off.
// Without this a fast spin down to the foot of a conversation followed by a
// turn back up goes nowhere until enough turns have been spun off, which is
// the stickiness at the foot of an agent chat.
func TestAWheelTurnWhileEasingAnswersAtOnce(t *testing.T) {
	var region *kvitui.Region
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		// Easing on: the wheel eases in over a few frames rather than
		// arriving at once.
		ui.Theme.SetReducedMotion(false)
		parts := []unison.Paneler{kvitui.NewListRow(ui, kvitui.NewLabel(ui, "First"))}
		for i := range 12 {
			parts = append(parts, kvitui.NewSlimRow(ui, "Row "+strconv.Itoa(i+1)))
		}
		parts = append(parts, kvitui.NewListRow(ui, kvitui.NewLabel(ui, "Last")))
		region = kvitui.NewRegion(ui, kvitui.Column(ui, kvitui.Px(0), parts...))
		region.SetLayoutData(&unison.FlexLayoutData{SizeHint: geom.NewSize(700, 140)})
		return []unison.Paneler{region}
	})
	position := func() (y float32) {
		screen.Do(func() { _, y = region.Position() })
		return y
	}
	// The middle, with room to travel either way.
	screen.Do(func() { region.ScrollTo(10000) })
	screen.Sync()
	var travel float32
	screen.Do(func() { _, travel = region.Position() })
	if travel <= 200 {
		t.Fatalf("fourteen rows in 140 px scroll only %.1f", travel)
	}
	mid := travel / 2
	screen.Do(func() { region.ScrollTo(mid) })
	screen.Sync()
	if y := position(); y != mid {
		t.Fatalf("scrollTo(mid) sits at %.1f, want %.1f", y, mid)
	}
	at := screen.PanelCenter(region.AsPanel())
	// Three fast turns down, with no frame between them for the easing to
	// pay in: only the first moves at once, the rest joins what is owed.
	screen.Wheel(at, geom.NewPoint(0, -1), 0)
	screen.Wheel(at, geom.NewPoint(0, -1), 0)
	screen.Wheel(at, geom.NewPoint(0, -1), 0)
	downs := position()
	if downs <= mid {
		t.Fatalf("three turns down sit at %.1f, never below %.1f", downs, mid)
	}
	// One turn back up, before any frame pays: it moves back up at once.
	screen.Wheel(at, geom.NewPoint(0, 1), 0)
	if up := position(); up >= downs {
		t.Errorf("a turn back up sits at %.1f, never above %.1f", up, downs)
	}
}
