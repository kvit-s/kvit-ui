package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
	"golang.org/x/text/language"
)

func TestABadgeCapsHidesAtZeroAndSaysItsNoun(t *testing.T) {
	var capped, zero, entries *kvitui.Badge
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.Locale = language.AmericanEnglish
		capped = kvitui.NewBadge(ui, 1204)
		zero = kvitui.NewBadge(ui, 0)
		entries = kvitui.NewBadge(ui, 1204)
		entries.Max, entries.Counted, entries.CountedPlural = 9999, "entry", "entries"
		return []unison.Paneler{capped, zero, entries}
	})
	if got := capped.Text(); got != "99+" {
		t.Errorf("1204 past a cap of 99 reads %q", got)
	}
	if got := entries.Text(); got != "1,204" {
		t.Errorf("1204 under a cap of 9999 reads %q", got)
	}
	// The cap belongs to the pill; a screen reader is told the whole count.
	if n := screen.AccessibilityNodeFor(capped); n == nil || n.Role != role.Label || n.Name != "1,204 items" {
		t.Errorf("the capped badge's node: %+v", n)
	}
	if n := screen.AccessibilityNodeFor(entries); n == nil || n.Name != "1,204 entries" {
		t.Errorf("the counted badge's node: %+v", n)
	}
	entries.Count = 1
	if got := entries.Phrase(); got != "1 entry" {
		t.Errorf("one entry is announced as %q", got)
	}
	screen.Do(func() {
		if _, p, _ := zero.Sizes(geom.Size{}); p.Width != 0 || p.Height != 0 {
			t.Errorf("a badge at zero takes %v", p)
		}
		if _, p, _ := capped.Sizes(geom.Size{}); p.Width < p.Height {
			t.Errorf("a badge is %v, narrower than it is tall", p)
		}
	})
	if n := screen.AccessibilityNodeFor(zero); n != nil && !n.Ignored {
		t.Errorf("a badge at zero is announced: %+v", n)
	}
}
