package kvitui_test

import (
	"strings"
	"testing"

	kvitui "github.com/kvit-s/kvit-ui"
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/accessibility"
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

func TestTheBreadcrumbLinksBackAndCutsTheMiddle(t *testing.T) {
	var short, long *kvitui.Breadcrumb
	var followed []string
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		short = kvitui.NewBreadcrumb(ui, kvitui.Crumb{Label: "Projects", ID: "p"}, kvitui.Crumb{Label: "kvit-ui", ID: "u"})
		short.OnActivate = func(c kvitui.Crumb) { followed = append(followed, c.ID) }
		long = kvitui.NewBreadcrumb(ui,
			kvitui.Crumb{Label: "Projects", ID: "p"}, kvitui.Crumb{Label: "kvit", ID: "k"},
			kvitui.Crumb{Label: "Wave 1", ID: "w"}, kvitui.Crumb{Label: "Tokens", ID: "t"},
			kvitui.Crumb{Label: "Theme", ID: "h"})
		return []unison.Paneler{short, long}
	})
	var labels []string
	for _, c := range long.Shown() {
		labels = append(labels, c.Label)
	}
	if got := strings.Join(labels, " / "); got != "Projects / … / Tokens / Theme" {
		t.Errorf("the long trail shows %q", got)
	}
	screen.Do(func() {
		if h := short.FrameRect().Height; h != float32(ui.Interface.BreadcrumbHeight()) {
			t.Errorf("the breadcrumb is %.1f tall", h)
		}
	})
	// Crumb, chevron, crumb: the first is a link and the last is not.
	kids := short.Children()
	if len(kids) != 3 {
		t.Fatalf("a two-crumb trail has %d children", len(kids))
	}
	first, last := screen.AccessibilityNodeFor(kids[0]), screen.AccessibilityNodeFor(kids[2])
	if first == nil || first.Role != role.Link || last == nil || last.Role != role.Label {
		t.Fatalf("first crumb %+v, last crumb %+v", first, last)
	}
	if !screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: first.ID, Action: accessibility.Press}) {
		t.Fatal("the first crumb did not take a press")
	}
	screen.Sync()
	if len(followed) != 1 || followed[0] != "p" {
		t.Errorf("following the first crumb reported %v", followed)
	}
	// The "…" stands for several places and cannot be followed.
	if n := screen.AccessibilityNodeFor(long.Children()[2]); n == nil || n.Role != role.Label {
		t.Errorf("the elided crumb is %+v", n)
	}
	if n := screen.AccessibilityNodeFor(short); n == nil || n.Role != role.Group || n.Name != "Breadcrumb" {
		t.Errorf("the trail's node: %+v", n)
	}
}

func TestTheHeaderKeepsItsThreePlaces(t *testing.T) {
	var h *kvitui.Header
	var today *kvitui.Tab
	var settings *kvitui.IconButton
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		today = kvitui.NewTab(ui, "Today")
		settings = kvitui.NewIconButton(ui, "settings", "Settings")
		h = kvitui.NewHeader(ui, "kvit", kvitui.Row(ui, kvitui.SizeSpace, today), settings)
		h.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Fill, HGrab: true})
		return []unison.Paneler{h}
	})
	screen.Do(func() {
		m := ui.Interface
		box := h.FrameRect()
		if box.Height != float32(m.HeaderHeight()) {
			t.Errorf("the header is %.1f tall, want %d", box.Height, m.HeaderHeight())
		}
		in := func(p unison.Paneler) geom.Rect {
			return h.RectFromRoot(p.AsPanel().RectToRoot(p.AsPanel().ContentRect(true)))
		}
		margin := float32(m.ViewMargin())
		mark := h.Children()[0].FrameRect()
		if mark.X != margin {
			t.Errorf("the wordmark starts at %.1f, want the view margin %.1f", mark.X, margin)
		}
		if tab := in(today); tab.X < mark.Right()+float32(m.ColumnGap())-0.5 {
			t.Errorf("the navigation starts at %.1f, inside the wordmark's gap", tab.X)
		}
		if right := in(settings).Right(); right < box.Width-margin-0.5 || right > box.Width-margin+0.5 {
			t.Errorf("the actions end at %.1f, want %.1f", right, box.Width-margin)
		}
	})
	if n := screen.AccessibilityNodeFor(h); n == nil || n.Role != role.TabList || n.Name != "Application header" {
		t.Errorf("the header's node: %+v", n)
	}
}

func TestASidebarItemSaysWhereItGoesAndHowMuchIsThere(t *testing.T) {
	var inbox, review *kvitui.SidebarItem
	var rail *kvitui.Sidebar
	pressed := 0
	screen, ui := session(t, func(ui *kvitui.UI) []unison.Paneler {
		ui.Locale = language.AmericanEnglish
		inbox = kvitui.NewSidebarItem(ui, "Inbox", "note")
		inbox.Selected, inbox.Count, inbox.Counted = true, 12, "message"
		inbox.OnClick = func() { pressed++ }
		review = kvitui.NewSidebarItem(ui, "Review", "question")
		review.Count = 214
		rail = kvitui.NewSidebar(ui, review)
		rail.Collapsed = true
		return []unison.Paneler{kvitui.NewSidebar(ui, inbox), rail}
	})
	n := screen.AccessibilityNodeFor(inbox)
	if n == nil || n.Role != role.ListItem || n.Name != "Inbox, 12 messages" || !n.Selected {
		t.Fatalf("the selected item's node: %+v", n)
	}
	screen.PerformAccessibilityAction(accessibility.ActionRequest{Node: n.ID, Action: accessibility.Press})
	screen.Sync()
	if pressed != 1 {
		t.Errorf("a screen reader's press made %d presses", pressed)
	}
	badgeOf := func(it *kvitui.SidebarItem) *kvitui.Badge { return it.Children()[0].Self.(*kvitui.Badge) }
	screen.Do(func() {
		m := ui.Interface
		if !review.Collapsed {
			t.Error("an item in a collapsed sidebar is not collapsed")
		}
		if w := rail.FrameRect().Width; w != float32(m.RailWidth()) {
			t.Errorf("the rail is %.1f wide, want %d", w, m.RailWidth())
		}
		if h := inbox.FrameRect().Height; h != float32(m.RowHeightSlim()) {
			t.Errorf("an item is %.1f tall, want %d", h, m.RowHeightSlim())
		}
		// The rail keeps its count behind unless asked to show it.
		if !badgeOf(inbox).Shows() || badgeOf(review).Shows() {
			t.Error("the badge shows in the rail, or not in the expanded sidebar")
		}
		review.CountInRail = true
		rail.MarkForLayoutRecursively()
		rail.ValidateLayout()
		if !badgeOf(review).Shows() || badgeOf(review).Text() != "99+" {
			t.Errorf("CountInRail shows %q", badgeOf(review).Text())
		}
	})
	if n := screen.AccessibilityNodeFor(rail); n == nil || n.Role != role.List || n.Name != "Navigation" {
		t.Errorf("the sidebar's node: %+v", n)
	}
}
