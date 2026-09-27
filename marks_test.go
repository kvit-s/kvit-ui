package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/kvit-s/kvit-ui/tokens"
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

func TestMarksSayWhatTheyMean(t *testing.T) {
	var chip *kvitui.Chip
	var tag, removable *kvitui.Tag
	var slug *kvitui.Slug
	var silent, dot *kvitui.Dot
	var pip *kvitui.Pip
	removed := 0
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		chip = kvitui.NewChip(ui, "stale")
		chip.Tone, chip.Explanation = kvitui.ToneWarning, "The note changed since this was staged."
		tag = kvitui.NewTag(ui, "groceries")
		tag.SetTint(tokens.ColorPalette()[0])
		removable = kvitui.NewTag(ui, "removable")
		removable.Removable = true
		removable.OnRemove = func() { removed++ }
		slug = kvitui.NewSlug(ui, "a1b2c3d4e5f6a7b8c9d0")
		silent = kvitui.NewDot(ui)
		dot = kvitui.NewDot(ui)
		dot.Shape, dot.Label = kvitui.ShapeDiamond, "stalled"
		pip = kvitui.NewPip(ui, 3, 5)
		return []unison.Paneler{chip, tag, removable, slug, silent, dot, pip}
	})
	for _, c := range []struct {
		p                 unison.Paneler
		role              role.Enum
		name, description string
	}{
		{chip, role.Label, "stale", "The note changed since this was staged."},
		{tag, role.Label, "groceries, Red", ""},
		{slug, role.Label, "a1b2c3d4e5f6a7b8c9d0", ""},
		{dot, role.Image, "stalled", ""},
		{pip, role.ProgressBar, "3 of 5", ""},
	} {
		n := screen.AccessibilityNodeFor(c.p)
		if n == nil || n.Role != c.role || n.Name != c.name || n.Description != c.description {
			t.Errorf("want %v %q (%q), got %+v", c.role, c.name, c.description, n)
		}
	}
	if n := screen.AccessibilityNodeFor(silent); n != nil && !n.Ignored {
		t.Errorf("a dot with no label is announced: %+v", n)
	}
	close := removable.Children()[0]
	if n := screen.AccessibilityNodeFor(close); n == nil || n.Name != "Remove removable" {
		t.Errorf("the close control's node: %+v", n)
	}
	screen.Click(screen.PanelCenter(close))
	if removed != 1 {
		t.Errorf("the close control removed %d times", removed)
	}
}

// A signal draws its number only past one, and grows sideways, not taller.
func TestASignalCountsOnlyPastOne(t *testing.T) {
	var one, many *kvitui.Signal
	screen, _ := session(t, func(ui *kvitui.UI) []unison.Paneler {
		one = kvitui.NewSignal(ui, "1 running")
		many = kvitui.NewSignal(ui, "240 failed")
		many.Count = 240
		return []unison.Paneler{one, many}
	})
	screen.Do(func() {
		_, a, _ := one.Sizes(geom.Size{})
		_, b, _ := many.Sizes(geom.Size{})
		if a.Width != a.Height || b.Height != a.Height || b.Width <= a.Width {
			t.Errorf("one is %v and 240 is %v", a, b)
		}
	})
}
