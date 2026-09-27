package kvitui_test

import (
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/role"
	"golang.org/x/text/language"
)

// The count is grouped by the reader's locale and uses a real plural, as
// test_components.cpp checks for the Qt head.
func TestTheViewHeadCountsInTheReadersLocale(t *testing.T) {
	ui, err := kvitui.New(kvitui.Options{IgnoreDesktop: true})
	if err != nil {
		t.Fatal(err)
	}
	ui.Locale = language.AmericanEnglish
	h := kvitui.NewViewHead(ui, "Transactions")
	h.Counted = "transaction"
	for _, c := range []struct {
		count  int
		plural string
		want   string
	}{
		{1, "", "1 transaction"},
		{250000, "", "250,000 transactions"},
		{250000, "entries", "250,000 entries"},
		{-1, "", ""},
	} {
		h.Count, h.CountedPlural = c.count, c.plural
		if got := h.CountPhrase(); got != c.want {
			t.Errorf("count %d, plural %q: %q, want %q", c.count, c.plural, got, c.want)
		}
	}
	for tag, want := range map[language.Tag]string{language.German: "250.000", language.Und: "250000"} {
		ui.Locale = tag
		if got := ui.Number(250000); got != want {
			t.Errorf("%v writes %q, want %q", tag, got, want)
		}
	}
}

func TestTheViewHeadIsARowTallWithItsControlsAtTheRight(t *testing.T) {
	var plain, full *kvitui.ViewHead
	var export *kvitui.IconButton
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		plain = kvitui.NewViewHead(ui, "Accounts")
		export = kvitui.NewIconButton(ui, "export", "Export")
		full = kvitui.NewViewHead(ui, "Transactions", export)
		full.Count, full.Counted = 1284, "transaction"
		full.Subtitle = "Everything since the account was opened, including the transfers between your own accounts."
		for _, h := range []*kvitui.ViewHead{plain, full} {
			h.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true, SizeHint: geom.NewSize(360, 0)})
		}
		return []unison.Paneler{plain, full}
	})
	screen.Do(func() {
		row := float32(ui.Interface.RowHeight())
		if h := plain.FrameRect().Height; h != row {
			t.Errorf("a head with only a title is %.1f tall, want a row, %.1f", h, row)
		}
		// Only the parts with something to say are there for a screen reader.
		if n := len(plain.Children()); n != 1 {
			t.Errorf("a head with only a title has %d children", n)
		}
		if n := len(full.Children()); n != 4 {
			t.Errorf("a head with every part has %d children", n)
		}
		// Narrow, the subtitle runs onto more lines, and the head grows to hold
		// them rather than letting them spill out of it.
		_, narrow, _ := full.Sizes(geom.NewSize(220, 0))
		if narrow.Height <= row {
			t.Errorf("at 220 wide the head is %.1f tall, no taller than a row", narrow.Height)
		}
		head := full.FrameRect()
		b := full.RectFromRoot(export.RectToRoot(export.ContentRect(true)))
		// The head keeps the view margin inside its own edges, as the rest of
		// a region does (kvit-cash's copy of kvit-ui, d32c373).
		if right, want := b.Right(), head.Width-float32(ui.Interface.ViewMargin()); right < want-0.5 || right > want+0.5 {
			t.Errorf("the control ends at %.1f, want the view margin in from the right edge, %.1f", right, want)
		}
		if mid := b.CenterY(); mid < head.Height/2-1 || mid > head.Height/2+1 {
			t.Errorf("the control is centred at %.1f, not on the head's middle %.1f", mid, head.Height/2)
		}
	})
	n := screen.AccessibilityNodeFor(full)
	if n == nil || n.Role != role.Heading || n.Name != "Transactions" {
		t.Errorf("the head's node: %+v", n)
	}
}
